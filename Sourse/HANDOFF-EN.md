# SquadAdmin — Technical Project Documentation

Russian version: HANDOFF.md

This document is for whoever will continue working on the project. The user manual is in `README-RU.md`;
the API contract is in `API.md` (the single source of truth for all endpoints).

## What It Is

An administration panel for the Squad game server. A single `SquadAdmin.exe` file: no installation,
no external dependencies, no internet required. It runs on the same Windows machine as the Squad
server, exposes a web panel at `http://127.0.0.1:8080`, and stores its data alongside itself
(`config.json` + `data.db`).

The project started as a request to "build an exe from SQCP source code" (Vue+Node+MySQL), but per
the client's clarifications it became a standalone Go application. The SQCP and SquadJS source code
was used only as a reference for the RCON protocol — the parser regular expressions were ported
from there verbatim.

## Client Requirements (locked)

- Run on the same machine as the Squad server (home hosting, Windows).
- No login/password: the panel is local, listens on `127.0.0.1`, access is restricted at the
  OS/network level (authentication was removed in 1.0.1 — see changelog).
- RCON moderation, config editor, server process management, log parsing, local player database,
  automation.
- Russian is the native UI language; starting from 1.0.6, an RU/EN language switcher is available.
- **Target platform is Windows only.** The Linux build was used during development for functional
  testing and is not included in the release package.

## Architecture

- Go 1.25, stdlib only + `modernc.org/sqlite` (pure Go, NO cgo — essential for cross-compiling the
  exe). The `golang.org/x/crypto/bcrypt` dependency was removed along with authentication in 1.0.1.
- Web server on `net/http`, routing via `http.ServeMux` with `METHOD /path/{id}` patterns (Go 1.22+).
- Frontend embedded via `embed` from `internal/web/static/`: plain HTML/CSS/JS with no frameworks
  or CDN, SPA with hash-based routing.
- Live updates via SSE at `/api/events`.

## File Overview

| File | Lines | Purpose |
|---|---|---|
| `main.go` | ~300 | Module wiring, merging player database with RCON, graceful shutdown |
| `internal/config/config.go` | 150 | `config.json` alongside the exe, defaults, deriving paths from the server folder |
| `internal/store/store.go` | ~760 | SQLite: bans, player database, punishments, notes, rules, audit, log events |
| `internal/rcon/rcon.go` | ~375 | Source RCON client: packets, multipacket, Squad broken-packet workaround, corrupt-size guard |
| `internal/rcon/squad.go` | 281 | Parsers for ListPlayers, ListSquads, maps, layers, chat, ShowServerInfo |
| `internal/rcon/manager.go` | 341 | Background loop: auto-connect, reconnect, polling, cache snapshot |
| `internal/confedit/confedit.go` | 284 | List and edit `.cfg` files with backups, structured Admins.cfg, map rotation |
| `internal/logwatch/logwatch.go` | 290 | Tail `SquadGame.log`, join/leave/error/teamkill events |
| `internal/procman/procman.go` | ~800 | Start/stop/restart server (with opMu against double start), process detection and clean shutdown, auto-restart |
| `internal/automation/automation.go` | 1088 | Scheduler, 5 rule types from API.md |
| `internal/web/server.go` | ~960 | HTTP skeleton: routes, SSE hub, status, players, console, maps, static (auth removed) |
| `internal/web/handlers.go` | ~995 | Bans, player database, configs, logs, process, automation, journal, settings |
| `internal/web/static/app.js` | ~4000 | SPA: i18n layer (first ~600 lines) + core (router, SSE, modals, toasts) + 12 sections |
| `internal/web/static/index.html` | ~185 | Page skeleton, SVG icon sprite, clock |
| `internal/web/static/style.css` | ~800 | "Field Command Post" theme (military khaki) |
| `internal/web/static/fonts/` | 16 files | Oswald / Inter / JetBrains Mono, woff2, latin + cyrillic |
| `internal/web/static/img/` | 3 files | Field position background, map texture, favicon |

## Key Technical Decisions (do not change without good reason)

- **SQLite strictly via `modernc.org/sqlite`**, `db.SetMaxOpenConns(1)`. Switching to
  `mattn/go-sqlite3` will break the cross-compiled exe: that library requires cgo.
- **RCON**: END packet used as end-of-response marker; the Squad broken-empty-packet workaround
  (`size=10` → skip 21 bytes) is mandatory — without it the client will desynchronize.
  The declared size of incoming packets is range-checked (`14…maximumPacketSize`): a corrupt or
  negative size tears down the connection with a reconnect instead of panicking; `readLoop` is
  additionally wrapped in `recover`.
- **Player actions** prefer `...ById` commands using the in-game ID from the online cache, falling
  back to name/SteamID variants only when the player is not currently online.
- **Teamkills** are best-effort: the log provides the victim's name and the killer's EOS ID;
  commands are taken from the RCON cache via the `logwatch.TeamResolver` interface (implemented
  in `main.go` on top of the manager snapshot).
- **Joins and leaves** are computed as a diff of the RCON online list (`tracker` in `main.go`),
  not from the log — this way they work even without access to the log file. join/leave events
  from `logwatch` are suppressed while RCON is connected to avoid duplicates.
- **Passwords in settings** are masked as `********`; an empty string or the mask value on save
  means "do not change".
- **Config backups** are taken before every write, up to 10 per file.
- **Cross-compilation**: no build tags, no `_windows.go` files. OS branching is done exclusively
  via `runtime.GOOS` at runtime. `syscall.Kill` and `syscall.SysProcAttr{Setpgid}` must not be
  used — they do not exist on Windows.
- **`WriteTimeout: 0`** on the HTTP server — otherwise the SSE stream is broken.
- **`SQUADADMIN_HOME`** — environment variable to relocate the working directory, needed for tests.

## Release Build

```
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o SquadAdmin.exe .
```

Output is approximately 11 MB, x86-64, no external dependencies.

## Testing

- `go test ./...` — unit tests for RCON parsers, mock RCON server, `confedit`, `store`, and `procman`.
- `internal/rcon/mockserver_test.go` starts a fake Squad server on a `net.Listener` and verifies,
  among other things, multipacket handling and the broken-packet workaround.
- `internal/rcon/framing_test.go` — regression for a critical bug: a negative or oversized packet
  is rejected without panic and the connection is torn down.
- `internal/procman/procman_test.go` — `windowsTaskkillArgs` (never `/PID 0`) and no double start
  under concurrent `Start()` calls.
- Functional verification during development was done with a Linux binary using `SQUADADMIN_HOME`:
  curl-based API checks (bans, 409 degradation, automation, settings, SSE) plus a run against a
  Python mock RCON. No login or password-change endpoints exist anymore.
- The frontend was tested under a custom DOM stub in Node: rendering of all 12 sections,
  degradation, SSE event handling, and HTML-escaping of nicknames and chat messages.

## What Was Not Covered

- **No live Squad server was available during development.** The parsers were ported from
  SquadJS/SQCP and verified against real-format fixtures, but production verification can only
  be done on a real server. The most likely source of discrepancies is the JSON keys in
  `ShowServerInfo`: they change between game versions; the parser is written to be tolerant
  (missing keys produce zero values).
- **Windows was not tested** — the development container had neither Windows nor wine. Only a
  successful cross-compilation was confirmed. Windows-specific areas to check first:
  `tasklist`/`taskkill` in `procman`, opening the browser via `cmd /c start`, backslash paths
  in `confedit`.
- Lifting a ban on the actual server (`AdminRemoveBanById`) is not implemented: the panel does
  not know the server-side ban identifier. The ban is lifted in the panel's database; a
  `Bans.cfg` export is provided for the server.

## Internationalisation (i18n)

Starting with version 1.0.6, the frontend supports two interface languages: Russian (the original)
and English. The backend (Go) was not changed — all text on the wire remains Russian; translation
is done entirely on the client.

### Architecture

At the top of `app.js` sits an i18n layer implemented as an IIFE. It runs before any application
code and registers the function `T` and several of its methods in the global scope.

**Dictionary.** The object `EN` contains approximately 500 entries in the form
`'Russian string': 'English text'`. The keys are the verbatim Russian strings used in the code;
Russian remains the source language and requires no separate dictionary of its own.

**`T(key)`** is the primary function for translating UI strings. In `ru` mode it returns the key
as-is (no dictionary lookup). In `en` mode it looks for an exact match in `EN`; if no entry is
found it returns the key unchanged. This provides safe degradation: a string that was omitted from
the dictionary simply stays in Russian.

**`T.err(msg)`** translates error messages coming from the backend. It first looks for an exact
match in `EN`. Then it tries the pattern `строка N: text` and translates it recursively
(`line N: <translated text>`). If neither matches, it iterates over known prefixes of the form
`"Не удалось …: "` (keys in `EN` that end with a space), translates the prefix, and recursively
translates the remainder. Prefixes are sorted by descending length so that longer matches win.

**`T.msg(msg)`** translates backend-generated messages (journal, automation log, SSE notices,
process operation results). It extends `T.err`: after an exact match it applies template rules
from the `RULES` array (regular expressions for parameterized messages), then prefixes, then
piecewise translation across the “, ” separator. See the “Backend messages” subsection below.

**`T.dom(root)`** translates the static markup in `index.html` (topbar, sidebar navigation). It
walks the DOM tree of `root` using a `TreeWalker`, translates text nodes by exact match
(preserving surrounding whitespace), and also translates the `title`, `placeholder`, and
`aria-label` attributes of all elements within `root`.

**`T.lang`** holds the current language (`'ru'` or `'en'`), read from `localStorage` (key
`sa_lang`, default `'ru'`). **`T.set(l)`** saves the choice to `localStorage`; the new language
takes effect after `location.reload()`.

**Integration with the rest of the code.** All Russian string literals in the application portion
of `app.js` are wrapped in `T('…')` by an automated transformer — the key strings remain Russian.
In the HTTP client `api()`, the `data.error` field from the server response is passed through
`T.err()`. In the Console section, the args/desc fields of commands returned by the backend are
translated via `T()` at render time.

**Initialisation.** `app.js` is linked at the end of `<body>`, so `T.dom(document.body)` is
called as soon as the script runs. At the same time, `document.documentElement.lang` is set to
`'ru'` or `'en'`.

**Language toggle entry points.** The `#lang-toggle` button in the topbar (`index.html`); its
style was added to `style.css` alongside `#theme-toggle`. The second path is the "Interface
language" segment in Settings → Appearance (`renderSettings` in `app.js`).

**Backend messages.** Action journal entries, automation log entries, server notifications
(SSE type `notice`), process operation results, settings check summaries and game-log events
are produced by the backend in Russian and stored in the database as-is. On the client they
are translated by `T.msg(msg)`: exact dictionary match → template rules (the `RULES` array in
the i18n layer: regular expressions for parameterized messages — PIDs, rule names in «quotes»,
player counters, times) → prefixes (“Не удалось …: ”) → piecewise translation across the “, ”
separator (for composite folder-check summaries). An unrecognized message is returned
unchanged — safe degradation. Because translation happens at render time, historical journal
entries are shown in English too, while the database stays canonically Russian. `T.msg` call
sites in `app.js`: the journal (details), the automation log (detail, lastResult), the SSE
notice toast, process status details, start/stop/restart results, RCON/folder check toasts,
and event texts in Logs → Events and the live feed on the Overview page.

### Rule for New Strings

Whenever a future change introduces a new UI string:

1. Write the string in Russian and wrap it in `T('…')`.
2. Add a matching entry to the `EN` dictionary at the top of `app.js`:
   `'Russian string': 'English text'`.

If step 2 is skipped, the string will simply remain Russian in EN mode — that is safe degradation,
not a bug. The dictionary is intentionally kept in a single location (the top of `app.js`) so it
can be extended without searching through the file.

## Changes in 1.0.1

Bug fixes following the technical audit of 1.0.0.

- **Authentication removed (per client requirement — local deployment).** The login screen,
  password change, sessions, and the `sa_session` cookie were removed. Removed from the backend:
  `/api/auth/*` and `/api/account/password` endpoints, `auth` and `ensureAdmin` middleware,
  KV keys `admin_hash`/`must_change_password`/`last_login_at`, and the `sessions` table.
  `guard` was kept as a pass-through wrapper to avoid rewriting the route table. The
  `x/crypto/bcrypt` dependency was removed from `go.mod`. This also resolves the audit finding
  about storing the session token in plaintext — there are no tokens anymore.
- **Critical: panic on abnormal RCON packet** (`rcon.go`). A range check on the packet size was
  added to `nextPacket`; `recover` was added to `readLoop`. A corrupt stream reconnects instead
  of crashing the process.
- **Critical: double server start** (`procman.go`). An `opMu` mutex was introduced; the "already
  running" check and `cmd.Start()` are now atomic under it — concurrent `Start()` calls launch
  exactly one process.
- **Critical: false "server stopped" + `taskkill /PID 0`** (`procman.go`). When the PID is
  unknown, the stop proceeds by image name (`/IM`); success is verified by checking for the
  process, and an error is returned on failure. The word "SIGKILL" was removed from Windows
  messages; the BOM was stripped from `tasklist` output.
- Other important and minor audit findings were deliberately excluded from this release and
  scheduled separately.

## Changes in 1.0.2

Closed all "IMPORTANT" findings from the 1.0.0 audit. Two of the ten were already resolved along
with the removal of authentication in 1.0.1 (login throttling and session invalidation on password
change) — their subject matter no longer exists. The remaining eight are addressed below.

- **Request origin check** (`web/server.go`). `Handler()` now returns a `mux` wrapped in
  `localOnly`: `Host` must be a loopback address, `Origin` (if present) must be as well,
  and `Sec-Fetch-Site: cross-site` is rejected. Error codes: `host_not_allowed`,
  `origin_not_allowed`, `cross_site_denied`. The check runs before routing, so it covers both
  static assets and SSE. The audit finding about CSRF/DNS-rebinding after password removal
  became the only remaining barrier — hence the strictness: `bindAll` no longer grants access
  from other machines (the socket still listens, but requests are rejected).
- **Atomic config writes** (`confedit/confedit.go`). Added `writeFileAtomic`:
  `os.CreateTemp` in the same directory → `Write` → `Sync` → `Chmod` → `Rename` with five
  retries to handle antivirus blocking, plus a `defer` to remove the temporary file. `WriteFile`
  uses this instead of calling `os.WriteFile` directly. The `.bak-<ts>` backup is still taken
  beforehand.
- **Preserving standalone comments in `Admins.cfg`** (`confedit/confedit.go`). Added
  `adminsLayout` and `parseAdminsLayout`: comments are attached to the entry that follows them
  (`byGroup` keyed by lowercase group name, `byAdmin` by ID) and are restored to their positions
  on serialisation. Comments belonging to entries deleted via the panel are not lost — they move
  to the end of the file.
- **Automation rule re-entrancy guard** (`automation/automation.go`). Added
  `running map[int64]bool` with `tryAcquire`/`release`; the flag is set in `tick()` **before**
  the goroutine is launched. `Store.RuleSetState` and `RuleMarkRun` now return `error`; calls are
  wrapped in `saveState`, which reports failures via `Notice`. For `restart_schedule`, a failed
  state save **cancels** the restart — otherwise the next tick 30 seconds later would restart the
  server again.
- **Log rotation by content** (`logwatch/logwatch.go`). The first 512 bytes are stored as a
  fingerprint (`headProbe`). Because the log is append-only, the old fingerprint must remain a
  prefix of the new content; if it does not, the file was replaced and is read from the
  beginning. A fingerprint was chosen over `os.SameFile` because on Windows the file identifier
  is resolved lazily by path, and the result is unpredictable for an already-replaced file.
  A UTF-8 BOM is also skipped.
- **Lenient `ListPlayers` parser** (`rcon/squad.go`). `EOS: (\w{32})` →
  `([^\s|]+)`, `steam:` is now optional, `Team ID` accepts `N/A`, minutes/seconds in
  "Since Disconnect" accept any length. Unrecognised lines that resemble player entries go to
  `UnparsedLogger` (defaults to `log`) with rate limiting: one log entry per line per 10 minutes,
  map cleared after 200 entries.
- **Database schema versioning** (`store/store.go`). `PRAGMA user_version` + a `migrations`
  list, each step in a transaction. The base schema was extracted into `schemaV1`. A database
  created by a newer version of the panel is rejected with a clear error message. An idempotent
  `ensureColumn` helper is available for future steps (SQLite does not support
  `ADD COLUMN IF NOT EXISTS`).
- **One RCON poll at a time** (`rcon/manager.go`). A `polling` flag and `tryBeginPoll`/`endPoll`
  are placed at the start of `pollOnce`, so a `ForcePoll()` triggered after a ban or kick does
  not interleave with the background polling cycle. Additionally, `LatencyMs` is now measured as
  the round-trip time for a single `ListPlayers` command, not "total poll time / 5".

Not included in 1.0.2 (remaining from the audit "minor items"): hard-coding of the Squad broken
packet workaround to a byte pattern, `Cache-Control: no-cache` for static assets, escaping `%`/`_`
in `SanitizeQuery`, validating `SeedLayer`/`LiveLayer` in `ValidateRule`, potential CRLF→LF
conversion when raw-editing other `.cfg` files, lack of transactions in multi-step `store`
operations. Tests for `automation`, `procman`, `logwatch`, `web`, and `config` are still absent.

## Changes in 1.0.3

Closed the remaining minor items from the 1.0.0 audit.

- **`Cache-Control: no-store`** (`web/server.go`). Static assets are now served with `no-store`
  instead of `no-cache` — the response is not stored in any cache.
- **LIKE escaping** (`store/store.go`). `SanitizeQuery` now replaces `\`, `%`, and `_` with
  escape sequences; all three groups of LIKE queries received `ESCAPE '\'`. Searching for `%`
  or `_` no longer matches everything.
- **ValidateRule for seeding** (`automation/automation.go`). When `switchLayer: true`, the
  presence of `seedLayer` and `liveLayer` is validated — the error is returned at rule creation
  time rather than silently failing at trigger time.
- **CRLF normalisation in configs** (`confedit/confedit.go`). `WriteFile` now calls
  `normalizeCRLF`: bare `\n` → `\r\n`. Squad on Windows requires CRLF; browsers send LF when
  performing raw edits. The `confedit_test.go` test was updated: the expected round-trip value
  is now `\r\n`.
- **Transactions in multi-step operations** (`store/store.go`). `UpsertSeenPlayer`,
  `SessionStart`, and `SessionEnd` are wrapped in `sql.Tx`. With `SetMaxOpenConns(1)` there is
  no real concurrency, but an explicit transaction guarantees atomicity if the connection limit
  is ever raised and also speeds up bulk snapshot inserts.
- **`doStopWindows` signature** (`procman/procman.go`). `cmd *exec.Cmd` and `cmdDead bool` were
  removed — POSIX constructs that do not apply on Windows (formerly `_ = cmd; _ = cmdDead`).
  The call site was updated: the Windows branch now passes only `pid` and `graceSec`.
- **Documentation of the Squad broken-packet pattern** (`rcon/rcon.go`). The comment on
  `nextPacket` now contains the exact byte signature (with a reference to SquadJS) and a note:
  if the format changes after a Squad patch, this is the block to check.

## Changes in 1.0.4 (visual only)

Go code was barely touched: only `Version` and the MIME/cache branch in `hStatic` were changed.
Everything else is static assets.

- **New theme** (`static/style.css`, rewritten in full). CSS tokens in `:root`: olive scale
  `--olive-*` for the frame, sand scale `--sand-*` for the workspace, `--ink*` for text,
  status colours `--green/--amber/--red/--blue`, terminal palette `--term-*`. Corner radii of
  2 px (angular, not rounded). The alias `--text3` is kept — it is used by several inline styles
  in `app.js`.
- **Local fonts** (`static/fonts/`). Oswald 500/600 (headings, buttons, labels), Inter 400–700
  (body text), JetBrains Mono 400/600 (console, IDs). Loaded via the generated `fonts/fonts.css`
  with `unicode-range` (latin + cyrillic). No CDN — the panel works without internet.
- **Background assets** (`static/img/`). `bg-field.jpg` — page background (opacity .18,
  grayscale + blur, with a veil and coordinate-grid overlay via gradients); `bg-map.jpg` —
  header texture (7 %, `mix-blend-mode: luminosity`); `favicon.svg`.
- **Page skeleton** (`static/index.html`). Added an SVG sprite with `<symbol id="i-*">`
  (12 sections + broadcast/lock/ban), the navigation is split into `.nav-section` groups, a
  sidebar footer `#sidebar-foot` and a clock `#clock` were added (an isolated script at the
  end of the file, unrelated to `app.js`).
  The DOM contract is preserved: `#app`, `#topbar`, `#rcon-status/.dot/.label/#rcon-host/
  #rcon-latency`, `#sse-indicator/#sse-label`, `#burger-btn`, `#sidebar`,
  `#sidebar-nav .nav-link[data-section]`, `#main`, `#view`, `#toast-container`,
  `#modal-container`.
- **`app.js`** — three targeted fixes: the emoji `📢`, `🔒`, `⛔` were replaced with
  `<svg class="ico"><use href="#i-*"/></svg>` (emoji fonts are not universally available and
  clashed with the visual style). Logic was not changed.
- **`web/server.go`** — `hStatic` now sets Content-Type explicitly for `.woff2`, `.svg`,
  `.jpg`, `.png` (on Windows `mime.TypeByExtension` reads the registry and may return garbage)
  and serves `fonts/` and `img/` with `Cache-Control: public, max-age=604800`; everything else
  keeps the previous `no-store`.
- **Exe size** grew by approximately 0.8 MB due to fonts and images.

Verification: `gofmt -l`, `go vet ./...`, `go test ./...` — clean; all 12 sections, the modal,
the context menu, and toasts were screenshotted in headless Chrome against a demo fixture API;
a real binary was also checked (RCON-not-configured state and static-serving headers).

---

## Changes in 1.0.5 (visual + validation)

- `internal/web/static/style.css`: colour palette extracted into CSS variables; dark theme is
  the `html[data-theme="night"]` block. Workspace background controlled by `--bg-image`/`--bg-pos`,
  set from JS.
- `internal/web/static/app.js`: `SA.appearance` module (theme, background, localStorage keys
  `sa_theme`/`sa_bg`/`sa_bg_url`/`sa_bg_pos`), a `#theme-toggle` button in the topbar, and an
  "Appearance" card in the settings. Anti-flash: an inline script in `index.html` applies the
  theme before `app.js` loads.
- `internal/web/static/img/backgrounds/` — 9 new 1920 px backgrounds + thumbnails (`thumbs/`),
  all bundled into the exe via go:embed.
- `internal/confedit/validate.go` — validation of files before writing: `CheckExtension`
  (.cfg files only), `ValidateContent` (general + per-filename checks), `ValidateServerPaths`
  (exe/log extensions, absolute paths, no `..`). Errors are `*ValidationError`; the HTTP layer
  returns them as 422 (`writeConfeditErr` in handlers.go). The injection point is
  `confedit.WriteFile`, which is called by both raw edits and structural editors.
- Server settings (`PUT /api/settings/server`) validate paths before saving. No client-side API
  changes are required: the frontend displays the error text from the response as before.

## Changes in 1.0.6 (internationalisation)

The backend was not changed. All changes are frontend and static assets.

- **i18n layer in `app.js`** (first ~600 lines). An IIFE adds the global function `T` with the
  `EN` dictionary (~570 entries), `T.err()` for API errors, `T.msg()` for backend-generated
  messages (journal, automation log, SSE notices — translated at render time, including
  historical entries), and `T.dom()` for static markup.
  See the "Internationalisation (i18n)" section above for details.
- **All strings in `app.js` wrapped in `T('…')`** by an automated transformer; the key strings
  remain Russian.
- **Language switcher** — the `#lang-toggle` button in the topbar (`index.html`) and the
  "Interface language" segment in Settings → Appearance. The choice is stored in `localStorage`
  (`sa_lang`) and applied via `location.reload()`.
- **`style.css`** — the `#lang-toggle` button style was added alongside `#theme-toggle`.
- **`const Version`** in `internal/web/server.go` = `"1.0.6"`.
- `index.html` footer updated to version 1.0.6.
