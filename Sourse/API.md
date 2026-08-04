# Squad Admin Panel — внутренний контракт

English version: API-EN.md

Приложение: один Go-бинарник (`SquadAdmin.exe` для Windows + бинарник для Linux VPS).
Веб-панель на `http://127.0.0.1:8080` (порт настраивается).
Язык интерфейса: **русский** (с версии 1.0.6 поддерживается переключатель RU/EN). Аутентификация отсутствует: панель запускается локально автором хостинга и слушает `127.0.0.1`. Обращения принимаются только с этого же компьютера — см. «Доступ». Дополнительное ограничение доступа — на стороне ОС/сети (не пробрасывать порт наружу).

Все ответы — JSON. Ошибка: HTTP 4xx/5xx + `{"error": "текст на русском"}`.
Текст поля `error` **всегда передаётся на русском языке** — независимо от выбранного языка интерфейса. В режиме EN фронтенд переводит текст ошибки на стороне клиента (словарь точных совпадений + известные префиксы); «на проводе» язык не меняется. Клиентам API не следует полагаться на конкретный текст ошибки для программной обработки — для этого предназначено поле `code` (см. примеры ниже).
Успех без данных: `{"ok": true}`.

## Режимы деградации

Сценарий: панель работает **на той же машине**, что и Squad-сервер (домашний хостинг, Windows).

| Модуль | Требует |
|---|---|
| Игроки, консоль, карты, баны, база игроков, автоматизация | только RCON |
| Конфиги сервера, логи | путь к папке Squad-сервера на этом же диске |
| Запуск, остановка, рестарт, авторестарт при падении | путь к exe сервера (локальный процесс) |

Если путь к серверу не указан — `/api/configs/*` и `/api/logs/*` возвращают 409
`{"error":"Путь к папке Squad-сервера не настроен","code":"server_dir_not_configured"}`,
`/api/process/*` — 409 `{"code":"process_not_configured"}`.
UI прячет соответствующие разделы (см. `GET /api/status` → `remote.configured`, `process.supported`).

---

## Доступ

Вход по логину/паролю удалён (локальный запуск). Отдельных эндпоинтов
аутентификации нет; все маршруты доступны сразу. Автором действий в журнале,
банах и наказаниях проставляется метка `admin`.

Поскольку пароля нет, единственный барьер — происхождение запроса. Все маршруты
(включая статику и SSE) проходят проверку до маршрутизации:

| Условие | Ответ |
|---|---|
| `Host` не `127.0.0.1` / `localhost` / `::1` | 403 `{"code":"host_not_allowed"}` |
| `Origin` задан и не локальный (или `null`) | 403 `{"code":"origin_not_allowed"}` |
| `Sec-Fetch-Site: cross-site` | 403 `{"code":"cross_site_denied"}` |

Это закрывает DNS-rebinding (чужой домен, резолвящийся в `127.0.0.1`) и CSRF со
стороннего сайта. Практическое следствие: при `bindAll: true` сокет слушает
`0.0.0.0`, но обращения с других машин по IP получают 403 — панель работает
только на том ПК, где запущена.

---

## Статус и события

### GET /api/status
```json
{
  "rcon": { "configured": true, "connected": true, "host": "1.2.3.4", "port": 21114,
            "lastError": "", "latencyMs": 42, "lastUpdate": "2026-08-02T15:00:00Z" },
  "remote": { "configured": false, "lastError": "" },
  "server": {
    "name": "MyServer #1", "playerCount": 78, "maxPlayers": 100,
    "publicQueue": 3, "reserveQueue": 0,
    "currentLayer": "Gorodok_AAS_v1", "nextLayer": "Yehorivka_RAAS_v1",
    "gameMode": "AAS", "teamOne": "RUS", "teamTwo": "USA",
    "matchTimeoutSec": 3600, "playtimeSec": 1200
  },
  "process": { "supported": false, "running": false },
  "version": "1.0.2", "uptimeSec": 3600
}
```
`server` = `null`, если RCON не подключён.

### GET /api/events  (SSE, text/event-stream)
События (`event:` + JSON в `data:`):

| event | data |
|---|---|
| `status` | тело `GET /api/status` |
| `players` | тело `GET /api/players` |
| `console` | `{"source":"rcon","command":"...","response":"...","at":"..."}` |
| `log` | `LogEvent` (см. ниже) |
| `notice` | `{"level":"info\|warn\|error","text":"..."}` |
| `chat` | `LogEvent` типа `chat` |

---

## Игроки (live, RCON)

### GET /api/players
```json
{
  "updatedAt": "2026-08-02T15:00:00Z",
  "online": [ { "playerID": 12, "eosID": "0002abc...", "steamID": "76561198000000000",
                "name": "Nickname", "teamID": 1, "squadID": 3, "isLeader": true,
                "role": "RUS_Rifleman_01", "knownId": 45 } ],
  "disconnected": [ { "playerID": 3, "eosID": "...", "steamID": "...", "name": "...",
                      "sinceDisconnect": "02m 13s" } ],
  "teams": [ { "teamID": 1, "name": "Russian Ground Forces",
               "squads": [ { "squadID": 1, "name": "SQUAD 1", "size": 9, "locked": false,
                             "players": [ /* как в online */ ] } ] } ]
}
```
`knownId` — id записи в базе игроков (для ссылки), может быть `0`.

### POST /api/players/action
```json
{ "action": "warn|kick|ban|forceTeamChange|removeFromSquad|disbandSquad|demoteCommander",
  "target": "eosID|steamID", "playerName": "для истории", "reason": "текст",
  "duration": "1d" }
```
`duration` для `ban`: `0`(навсегда)/`1h`/`1d`/`3d`/`7d`/`30d`/`1M`/`перм`.
→ `{"ok":true, "response":"ответ RCON"}`. Пишет в аудит и историю наказаний.

### POST /api/players/broadcast  `{"message":"текст"}`

---

## Консоль

### POST /api/rcon/execute  `{"command":"ListPlayers"}` → `{"response":"..."}`
Требует `console.execute`. Пишется в аудит.

### GET /api/rcon/commands → справочник команд для автодополнения
`[{"cmd":"AdminBroadcast","args":"<сообщение>","desc":"Сообщение всем"}]`

---

## Карты

### GET /api/layers → `{"layers":["Gorodok_AAS_v1", "..."], "cachedAt":"..."}`
### GET /api/layers/current → `{"current":{"level":"...","layer":"..."},"next":{...}}`
### POST /api/layers/next     `{"layer":"..."}`
### POST /api/layers/change   `{"layer":"..."}`
### POST /api/match/end  |  POST /api/match/restart

---

## Баны (локальная база + отправка на сервер)

### GET /api/bans?active=true&query=&limit=50&offset=0
```json
{ "total": 120, "items": [ { "id": 7, "playerName": "...", "steamID": "...", "eosID": "...",
    "reason": "...", "durationLabel": "7 дней", "expiresAt": "...|null",
    "active": true, "createdBy": "admin", "createdAt": "...", "note": "" } ] }
```
### POST /api/bans   `{ "steamID|eosID", "playerName", "reason", "duration", "pushToServer": true }`
### PUT /api/bans/{id}     — правка причины/срока
### DELETE /api/bans/{id}  — снять бан (в т.ч. `AdminRemoveBanById`, если есть id сервера)
### GET /api/bans/export.txt — формат `Bans.cfg`: `76561198000000000:0 //Ник — причина`

---

## База игроков

### GET /api/playerdb?query=&sort=lastSeen&limit=50&offset=0
```json
{ "total": 900, "items": [ { "id": 45, "name": "Ник", "steamID": "...", "eosID": "...",
  "firstSeen": "...", "lastSeen": "...", "playtimeSec": 360000, "sessions": 42,
  "punishments": 3, "banned": false, "noteCount": 1 } ] }
```
### GET /api/playerdb/{id}
`{ "player": {...}, "names": [{"name":"...","seenAt":"..."}], "sessions":[{"startedAt","endedAt","durationSec"}], "punishments":[{"type":"ban|kick|warn","reason","by","at","durationLabel"}], "notes":[{"id","text","by","at"}] }`
### POST /api/playerdb/{id}/notes  `{"text":"..."}`   ### DELETE /api/playerdb/notes/{id}

---

## Конфиги сервера (локальные файлы)

### GET /api/configs → `{"files":[{"name":"Admins.cfg","path":"...","size":1234,"modifiedAt":"...","exists":true}]}`
### GET /api/configs/file?name=Server.cfg → `{"name":"...","content":"...","modifiedAt":"..."}`
### PUT /api/configs/file  `{"name":"Server.cfg","content":"..."}` → бэкап `.bak-<ts>` и запись
### GET /api/configs/admins → структурировано
```json
{ "groups": [ { "name": "Admin", "permissions": ["kick","ban","chat"] } ],
  "admins": [ { "id": "76561198000000000", "group": "Admin", "comment": "Ник" } ],
  "raw": "исходный текст" }
```
### PUT /api/configs/admins — тем же телом (без `raw`), пересобирает файл, сохраняя комментарии-шапку и отдельные строки-комментарии между группами/админами
### GET /api/configs/rotation → `{"layers":["..."]}`  ### PUT /api/configs/rotation `{"layers":[...]}`

Известные права групп Squad: `changemap, canseeadminchat, balance, chat, kick, ban, config,
cameraman, immune, manageserver, teamchange, forceteamchange, reserve, demos, debug, startvote,
private, clientdemos`.

---

## Логи (локальный tail файла)

### GET /api/logs/events?type=all|chat|teamkill|join|leave|admin|error&query=&limit=200
→ `{"items":[LogEvent]}`
```json
LogEvent = { "id": 1, "at": "2026-08-02T15:00:00Z", "type": "chat",
  "channel": "ChatAll|ChatTeam|ChatSquad|ChatAdmin", "playerName": "...",
  "steamID": "...", "eosID": "...", "text": "...", "raw": "исходная строка" }
```
### GET /api/logs/raw?lines=300 → `{"lines":["..."]}`
### GET /api/logs/state → `{"tailing":true,"offset":123456,"file":"SquadGame.log","lastReadAt":"..."}`

---

## Процесс сервера (локальный)

### GET /api/process/status
→ `{"supported":true,"running":true,"pid":1234,"uptimeSec":3600,"autoRestart":true,"detail":""}`
Определение "запущен": процесс, стартованный панелью, жив; либо найден процесс по имени exe.
### POST /api/process/start | /stop | /restart
`{"warnPlayers":true,"warnSeconds":60,"message":"Рестарт через 60 секунд"}`
→ `{"ok":true,"output":"..."}`
`stop`/`restart` при `warnPlayers` шлют AdminBroadcast и ждут `warnSeconds`, затем действуют.
### PUT /api/process/autorestart `{"enabled":true}` — авторестарт при падении (следит, перезапускает, пишет notice)

---

## Автоматизация

### GET /api/automation/rules
```json
{ "rules": [ { "id": 1, "name": "Правила сервера", "type": "broadcast_interval",
  "enabled": true, "config": { "message": "...", "intervalMin": 15 },
  "lastRunAt": "...", "nextRunAt": "...", "lastResult": "ok" } ] }
```
Типы и `config`:
- `broadcast_interval` — `{ "message": "...", "intervalMin": 15 }`
- `broadcast_schedule` — `{ "message": "...", "times": ["18:00","21:00"] }`
- `command_schedule` — `{ "command": "AdminRestartMatch", "times": ["06:00"] }`
- `restart_schedule` — `{ "times": ["06:00"], "warnMinutes": [10,5,1], "message": "Рестарт через {n} мин" }` (нужен настроенный процесс)
- `seeding` — `{ "startBelow": 20, "stopAbove": 50, "seedLayer": "Sumari_Seed_v1", "liveLayer": "", "message": "Идёт сидинг...", "intervalMin": 5, "switchLayer": true }`

### POST /api/automation/rules  ### PUT /api/automation/rules/{id}  ### DELETE /api/automation/rules/{id}
### POST /api/automation/rules/{id}/run — выполнить немедленно
### GET /api/automation/log?limit=100 → `{"items":[{"at","ruleName","result","detail"}]}`

---

## Журнал действий

### GET /api/audit?limit=100&offset=0 → `{"total":N,"items":[{"at","action","details"}]}`
(журнал того, что делала панель: команды, баны, правки конфигов, срабатывания автоматизации)

---

## Настройки

### GET /api/settings → (пароли маскируются как `"********"`)
```json
{ "general": { "panelPort": 8080, "bindAll": false, "pollIntervalSec": 20,
               "openBrowser": true },
  "rcon": { "host": "127.0.0.1", "port": 21114, "password": "********",
            "autoReconnectSec": 10, "keepAliveSec": 30 },
  "server": { "dir": "C:\\servers\\SquadServer",
              "configDir": "", "logFile": "",
              "exePath": "C:\\servers\\SquadServer\\SquadGameServer.exe",
              "exeArgs": "Port=7787 QueryPort=27165 RCONPORT=21114",
              "stopGraceSec": 10, "autoRestart": false } }
```
`configDir`/`logFile` пустые = вывести из `dir` (`<dir>/SquadGame/ServerConfig`,
`<dir>/SquadGame/Saved/Logs/SquadGame.log`).
### PUT /api/settings/general | /api/settings/rcon | /api/settings/server
Пустой пароль или `********` = не менять сохранённый.
### POST /api/settings/rcon/test → `{"ok":true,"detail":"Подключено, сервер: ..."}`
### POST /api/settings/server/test → `{"ok":true,"detail":"Папка найдена, конфиги: 5 файлов, лог найден"}`

---

## Служебное
### GET /healthz → `ok`
### Статика: `/` — SPA, всё неизвестное (кроме `/api/*`) отдаёт `index.html`.
