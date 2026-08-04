// SquadAdmin — панель администрирования Squad-сервера.
// Один исполняемый файл: веб-панель, RCON-модерация, конфиги, логи,
// управление процессом сервера и автоматизация.
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"
	"time"

	"squadadmin/internal/automation"
	"squadadmin/internal/config"
	"squadadmin/internal/logwatch"
	"squadadmin/internal/procman"
	"squadadmin/internal/rcon"
	"squadadmin/internal/store"
	"squadadmin/internal/web"
)

func main() {
	log.SetFlags(log.Ltime)
	log.Printf("SquadAdmin %s — запуск", web.Version)

	base := config.BaseDir()
	log.Printf("рабочая папка: %s", base)

	cfg, err := config.Load()
	if err != nil {
		fatal("не удалось прочитать config.json: %v", err)
	}

	st, err := store.Open(filepath.Join(base, "data.db"))
	if err != nil {
		fatal("не удалось открыть базу data.db: %v", err)
	}
	defer st.Close()

	// ── RCON ──
	rm := rcon.NewManager()
	rm.Configure(cfg.Rcon.Host, cfg.Rcon.Port, cfg.Rcon.Password,
		cfg.Rcon.AutoReconnectSec, cfg.General.PollIntervalSec)

	// ── чтение лога сервера ──
	lw := logwatch.New()
	lw.SetPath(cfg.LogFilePath())
	lw.Teams = teamResolver{rm}

	// ── процесс сервера ──
	pm := procman.New(cfg, rm)

	// ── планировщик автоматизации ──
	// hub появится позже, поэтому уведомления идут через переадресацию
	var noticeSink func(level, text string)
	notice := func(level, text string) {
		if noticeSink != nil {
			noticeSink(level, text)
		}
	}
	sched := automation.New(automation.Deps{Store: st, Rcon: rm, Proc: pm, Notice: notice})

	// ── HTTP-слой ──
	srv := web.New(web.Deps{
		Config: cfg, Store: st, Rcon: rm, Log: lw, Proc: pm, Auto: sched,
	})
	hub := srv.Hub()

	// ── связывание событий ──
	tr := newTracker(st, hub)

	noticeSink = hub.Notice
	pm.SetNotice(hub.Notice)
	pm.SetAudit(st.Audit)

	rm.OnPlayersUpdated = func() {
		tr.sync(rm.Snapshot().Online)
		srv.BroadcastPlayers()
		srv.BroadcastStatus()
	}
	rm.OnStatusChanged = func() { srv.BroadcastStatus() }

	rm.OnChatMessage = func(msg *rcon.ChatMessage) {
		if msg == nil {
			return
		}
		e := &store.LogEvent{
			Type: "chat", Channel: msg.Channel, PlayerName: msg.Name,
			SteamID: msg.SteamID, EosID: msg.EosID, Text: msg.Message,
			Raw: fmt.Sprintf("[%s] %s: %s", msg.Channel, msg.Name, msg.Message),
		}
		st.LogEventAdd(e)
		hub.Broadcast("chat", e)
	}

	lw.OnEvent = func(ev logwatch.Event) {
		// входы/выходы надёжнее считаются по опросу RCON — не дублируем
		if (ev.Type == "join" || ev.Type == "leave") && rm.Connected() {
			return
		}
		e := &store.LogEvent{
			Type: ev.Type, PlayerName: ev.PlayerName, SteamID: ev.SteamID,
			EosID: ev.EosID, Text: ev.Text, Raw: ev.Raw,
		}
		st.LogEventAdd(e)
		hub.Broadcast("log", e)
	}

	// ── запуск фоновых служб ──
	rm.Start()
	lw.Start()
	pm.StartWatcher()
	sched.Start()
	srv.StartBackground()

	// ── HTTP-сервер ──
	host := "127.0.0.1"
	if cfg.General.BindAll {
		host = "0.0.0.0"
	}
	addr := net.JoinHostPort(host, fmt.Sprintf("%d", cfg.General.PanelPort))
	httpSrv := &http.Server{
		Addr:              addr,
		Handler:           srv.Handler(),
		ReadHeaderTimeout: 15 * time.Second,
		WriteTimeout:      0, // 0 — иначе рвётся поток SSE
		IdleTimeout:       120 * time.Second,
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		fatal("не удалось занять адрес %s: %v\n"+
			"Возможно, порт занят другой программой или уже запущена вторая копия SquadAdmin.\n"+
			"Смените порт в файле config.json (поле general.panelPort).", addr, err)
	}

	url := fmt.Sprintf("http://127.0.0.1:%d", cfg.General.PanelPort)
	log.Printf("панель доступна: %s", url)
	if cfg.General.BindAll {
		log.Printf("bindAll включён: сокет слушает 0.0.0.0:%d, но панель принимает запросы "+
			"только с этого компьютера (127.0.0.1 / localhost) — с других машин будет 403",
			cfg.General.PanelPort)
	}
	log.Printf("вход по паролю не требуется — панель работает только на этом ПК")

	go func() {
		if err := httpSrv.Serve(ln); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP-сервер остановлен: %v", err)
		}
	}()

	if cfg.General.OpenBrowser {
		go openBrowser(url)
	}

	// ── ожидание завершения ──
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	log.Printf("завершение работы…")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = httpSrv.Shutdown(ctx)

	srv.StopBackground()
	sched.Stop()
	pm.StopWatcher()
	lw.Stop()
	rm.Stop()
	log.Printf("SquadAdmin остановлен")
}

func fatal(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	log.Printf("ОШИБКА: %s", msg)
	// на Windows окно закрывается мгновенно — дадим прочитать сообщение
	if runtime.GOOS == "windows" {
		fmt.Println("\nНажмите Enter, чтобы закрыть окно…")
		fmt.Scanln()
	}
	os.Exit(1)
}

// openBrowser открывает панель в браузере пользователя.
func openBrowser(url string) {
	time.Sleep(400 * time.Millisecond)
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	if err := cmd.Start(); err != nil {
		log.Printf("не удалось открыть браузер: %v (откройте %s вручную)", err, url)
		return
	}
	go func() { _ = cmd.Wait() }()
}

// ════════════════════════════ склейка базы игроков ════════════════════════════

// tracker сопоставляет онлайн из RCON с локальной базой игроков:
// заводит записи, ведёт сессии и пишет события входа/выхода.
type tracker struct {
	mu       sync.Mutex
	st       *store.Store
	hub      *web.Hub
	sessions map[string]int64 // eosID → id открытой сессии
	names    map[string]string
	primed   bool
}

func newTracker(st *store.Store, hub *web.Hub) *tracker {
	return &tracker{
		st: st, hub: hub,
		sessions: map[string]int64{},
		names:    map[string]string{},
	}
}

// sync вызывается после каждого успешного опроса RCON.
func (t *tracker) sync(online []rcon.OnlinePlayer) {
	t.mu.Lock()
	defer t.mu.Unlock()

	cur := make(map[string]rcon.OnlinePlayer, len(online))
	for _, p := range online {
		if p.EosID == "" {
			continue
		}
		cur[p.EosID] = p
	}

	// вошедшие
	for eos, p := range cur {
		if _, ok := t.sessions[eos]; ok {
			continue
		}
		pid, err := t.st.UpsertSeenPlayer(eos, p.SteamID, p.Name)
		if err != nil || pid == 0 {
			continue
		}
		sid, err := t.st.SessionStart(pid)
		if err == nil {
			t.sessions[eos] = sid
		}
		t.names[eos] = p.Name
		// при первом опросе после старта панели уже находящиеся на сервере
		// игроки не считаются «вошедшими» — не засоряем журнал
		if t.primed {
			t.emit("join", p.Name, p.SteamID, eos, "зашёл на сервер")
		}
	}

	// вышедшие
	for eos, sid := range t.sessions {
		if _, ok := cur[eos]; ok {
			continue
		}
		_ = t.st.SessionEnd(sid)
		delete(t.sessions, eos)
		name := t.names[eos]
		delete(t.names, eos)
		if t.primed {
			t.emit("leave", name, "", eos, "вышел с сервера")
		}
	}

	t.primed = true
}

func (t *tracker) emit(typ, name, steamID, eosID, text string) {
	e := &store.LogEvent{
		Type: typ, PlayerName: name, SteamID: steamID, EosID: eosID,
		Text: text, Raw: fmt.Sprintf("%s: %s", name, text),
	}
	t.st.LogEventAdd(e)
	if t.hub != nil {
		t.hub.Broadcast("log", e)
	}
}

// ════════════════════════════ команды игроков для логов ════════════════════════════

// teamResolver отдаёт logwatch команду игрока из кэша RCON —
// нужно, чтобы отличить тимкилл от обычного убийства.
type teamResolver struct{ m *rcon.Manager }

func (t teamResolver) TeamOfName(name string) (int, bool) {
	for _, p := range t.m.Snapshot().Online {
		if p.Name == name {
			return p.TeamID, true
		}
	}
	return 0, false
}

func (t teamResolver) TeamOfEos(eosID string) (int, bool) {
	if p := t.m.FindOnline(eosID); p != nil {
		return p.TeamID, true
	}
	return 0, false
}

func (t teamResolver) NameOfEos(eosID string) (string, bool) {
	if p := t.m.FindOnline(eosID); p != nil {
		return p.Name, true
	}
	return "", false
}
