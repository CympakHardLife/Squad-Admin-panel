# Squad Admin Panel — Internal API Contract

Russian version: API.md

Application: a single Go binary (`SquadAdmin.exe` for Windows + a binary for Linux VPS).
Web panel at `http://127.0.0.1:8080` (port is configurable).
Interface language: **Russian** (starting from version 1.0.6 an RU/EN language toggle is available). No authentication: the panel runs locally by the hosting owner and listens on `127.0.0.1`. Requests are accepted only from the same machine — see "Access". Additional access restriction is handled at the OS/network level (do not expose the port externally).

All responses are JSON. Error: HTTP 4xx/5xx + `{"error": "текст на русском"}`.
The `error` field text is **always transmitted in Russian** — regardless of the selected UI language. In EN mode the frontend translates the error text on the client side (exact-match dictionary + known prefixes); the wire protocol does not change. API clients should not rely on the specific error text for programmatic handling — the `code` field (see examples below) is provided for that purpose.
Success with no data: `{"ok": true}`.

## Degradation Modes

Scenario: the panel runs **on the same machine** as the Squad server (home hosting, Windows).

| Module | Requires |
|---|---|
| Players, console, maps, bans, player database, automation | RCON only |
| Server configs, logs | path to the Squad server folder on the same disk |
| Start, stop, restart, auto-restart on crash | path to server exe (local process) |

If the server path is not configured — `/api/configs/*` and `/api/logs/*` return 409
`{"error":"Путь к папке Squad-сервера не настроен","code":"server_dir_not_configured"}` <!-- "Squad server folder path is not configured" -->
`/api/process/*` — 409 `{"code":"process_not_configured"}`.
The UI hides the corresponding sections (see `GET /api/status` → `remote.configured`, `process.supported`).

---

## Access

Login/password authentication has been removed (local launch). There are no dedicated authentication endpoints; all routes are available immediately. The author label `admin` is applied to actions in the audit log, bans, and punishments.

Since there is no password, the only barrier is request origin. All routes (including static files and SSE) pass a check before routing:

| Condition | Response |
|---|---|
| `Host` is not `127.0.0.1` / `localhost` / `::1` | 403 `{"code":"host_not_allowed"}` |
| `Origin` is set and is not local (or is `null`) | 403 `{"code":"origin_not_allowed"}` |
| `Sec-Fetch-Site: cross-site` | 403 `{"code":"cross_site_denied"}` |

This mitigates DNS-rebinding (a foreign domain resolving to `127.0.0.1`) and CSRF from third-party sites. Practical implication: with `bindAll: true` the socket listens on `0.0.0.0`, but requests from other machines by IP receive 403 — the panel works only on the PC where it is running.

---

## Status and Events

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
`server` = `null` if RCON is not connected.

### GET /api/events  (SSE, text/event-stream)
Events (`event:` + JSON in `data:`):

| event | data |
|---|---|
| `status` | body of `GET /api/status` |
| `players` | body of `GET /api/players` |
| `console` | `{"source":"rcon","command":"...","response":"...","at":"..."}` |
| `log` | `LogEvent` (see below) |
| `notice` | `{"level":"info\|warn\|error","text":"..."}` |
| `chat` | `LogEvent` of type `chat` |

---

## Players (live, RCON)

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
                             "players": [ /* same structure as online */ ] } ] } ]
}
```
`knownId` — record ID in the player database (for linking), may be `0`.

### POST /api/players/action
```json
{ "action": "warn|kick|ban|forceTeamChange|removeFromSquad|disbandSquad|demoteCommander",
  "target": "eosID|steamID", "playerName": "for history", "reason": "text",
  "duration": "1d" }
```
`duration` for `ban`: `0` (permanent) / `1h` / `1d` / `3d` / `7d` / `30d` / `1M` / `перм` <!-- "перм" = permanent, Russian alias accepted by server -->.
→ `{"ok":true, "response":"ответ RCON"}` <!-- "ответ RCON" = RCON response -->. Writes to audit log and punishment history.

### POST /api/players/broadcast  `{"message":"текст"}` <!-- "текст" = message text -->

---

## Console

### POST /api/rcon/execute  `{"command":"ListPlayers"}` → `{"response":"..."}`
Requires `console.execute` permission. Written to audit log.

### GET /api/rcon/commands → command reference for autocomplete
`[{"cmd":"AdminBroadcast","args":"<сообщение>","desc":"Сообщение всем"}]`
<!-- "сообщение" = message, "Сообщение всем" = "Broadcast to all" — these are server-side values -->

---

## Maps

### GET /api/layers → `{"layers":["Gorodok_AAS_v1", "..."], "cachedAt":"..."}`
### GET /api/layers/current → `{"current":{"level":"...","layer":"..."},"next":{...}}`
### POST /api/layers/next     `{"layer":"..."}`
### POST /api/layers/change   `{"layer":"..."}`
### POST /api/match/end  |  POST /api/match/restart

---

## Bans (local database + push to server)

### GET /api/bans?active=true&query=&limit=50&offset=0
```json
{ "total": 120, "items": [ { "id": 7, "playerName": "...", "steamID": "...", "eosID": "...",
    "reason": "...", "durationLabel": "7 дней", "expiresAt": "...|null",
    "active": true, "createdBy": "admin", "createdAt": "...", "note": "" } ] }
```
<!-- "7 дней" = "7 days" — server-side Russian value for durationLabel -->
### POST /api/bans   `{ "steamID|eosID", "playerName", "reason", "duration", "pushToServer": true }`
### PUT /api/bans/{id}     — edit reason/duration
### DELETE /api/bans/{id}  — unban (including `AdminRemoveBanById` if a server-side ban ID exists)
### GET /api/bans/export.txt — `Bans.cfg` format: `76561198000000000:0 //Ник — причина`
<!-- "Ник — причина" = "Nickname — reason" — this is the actual file format -->

---

## Player Database

### GET /api/playerdb?query=&sort=lastSeen&limit=50&offset=0
```json
{ "total": 900, "items": [ { "id": 45, "name": "Ник", "steamID": "...", "eosID": "...",
  "firstSeen": "...", "lastSeen": "...", "playtimeSec": 360000, "sessions": 42,
  "punishments": 3, "banned": false, "noteCount": 1 } ] }
```
<!-- "Ник" = player nickname — actual value from server -->
### GET /api/playerdb/{id}
`{ "player": {...}, "names": [{"name":"...","seenAt":"..."}], "sessions":[{"startedAt","endedAt","durationSec"}], "punishments":[{"type":"ban|kick|warn","reason","by","at","durationLabel"}], "notes":[{"id","text","by","at"}] }`
### POST /api/playerdb/{id}/notes  `{"text":"..."}`   ### DELETE /api/playerdb/notes/{id}

---

## Server Configs (local files)

### GET /api/configs → `{"files":[{"name":"Admins.cfg","path":"...","size":1234,"modifiedAt":"...","exists":true}]}`
### GET /api/configs/file?name=Server.cfg → `{"name":"...","content":"...","modifiedAt":"..."}`
### PUT /api/configs/file  `{"name":"Server.cfg","content":"..."}` → creates backup `.bak-<ts>` and writes
### GET /api/configs/admins → structured
```json
{ "groups": [ { "name": "Admin", "permissions": ["kick","ban","chat"] } ],
  "admins": [ { "id": "76561198000000000", "group": "Admin", "comment": "Ник" } ],
  "raw": "исходный текст" }
```
<!-- "Ник" = nickname; "исходный текст" = raw source text -->
### PUT /api/configs/admins — same body (without `raw`), rebuilds the file preserving the header comments and individual comment lines between groups/admins
### GET /api/configs/rotation → `{"layers":["..."]}`  ### PUT /api/configs/rotation `{"layers":[...]}`

Known Squad group permissions: `changemap, canseeadminchat, balance, chat, kick, ban, config,
cameraman, immune, manageserver, teamchange, forceteamchange, reserve, demos, debug, startvote,
private, clientdemos`.

---

## Logs (local file tail)

### GET /api/logs/events?type=all|chat|teamkill|join|leave|admin|error&query=&limit=200
→ `{"items":[LogEvent]}`
```json
LogEvent = { "id": 1, "at": "2026-08-02T15:00:00Z", "type": "chat",
  "channel": "ChatAll|ChatTeam|ChatSquad|ChatAdmin", "playerName": "...",
  "steamID": "...", "eosID": "...", "text": "...", "raw": "исходная строка" }
```
<!-- "исходная строка" = raw log line -->
### GET /api/logs/raw?lines=300 → `{"lines":["..."]}`
### GET /api/logs/state → `{"tailing":true,"offset":123456,"file":"SquadGame.log","lastReadAt":"..."}`

---

## Server Process (local)

### GET /api/process/status
→ `{"supported":true,"running":true,"pid":1234,"uptimeSec":3600,"autoRestart":true,"detail":""}`
Definition of "running": a process started by the panel is alive; or a process is found by exe name.
### POST /api/process/start | /stop | /restart
`{"warnPlayers":true,"warnSeconds":60,"message":"Рестарт через 60 секунд"}`
<!-- "Рестарт через 60 секунд" = "Restart in 60 seconds" — example message sent via AdminBroadcast -->
→ `{"ok":true,"output":"..."}`
`stop`/`restart` with `warnPlayers` send an AdminBroadcast and wait `warnSeconds`, then act.
### PUT /api/process/autorestart `{"enabled":true}` — auto-restart on crash (monitors, restarts, writes a notice)

---

## Automation

### GET /api/automation/rules
```json
{ "rules": [ { "id": 1, "name": "Правила сервера", "type": "broadcast_interval",
  "enabled": true, "config": { "message": "...", "intervalMin": 15 },
  "lastRunAt": "...", "nextRunAt": "...", "lastResult": "ok" } ] }
```
<!-- "Правила сервера" = "Server rules" — example rule name, user-defined value -->
Types and `config`:
- `broadcast_interval` — `{ "message": "...", "intervalMin": 15 }`
- `broadcast_schedule` — `{ "message": "...", "times": ["18:00","21:00"] }`
- `command_schedule` — `{ "command": "AdminRestartMatch", "times": ["06:00"] }`
- `restart_schedule` — `{ "times": ["06:00"], "warnMinutes": [10,5,1], "message": "Рестарт через {n} мин" }` <!-- "Рестарт через {n} мин" = "Restart in {n} min" — template, {n} is substituted --> (requires a configured process)
- `seeding` — `{ "startBelow": 20, "stopAbove": 50, "seedLayer": "Sumari_Seed_v1", "liveLayer": "", "message": "Идёт сидинг...", "intervalMin": 5, "switchLayer": true }` <!-- "Идёт сидинг..." = "Seeding in progress..." -->

### POST /api/automation/rules  ### PUT /api/automation/rules/{id}  ### DELETE /api/automation/rules/{id}
### POST /api/automation/rules/{id}/run — run immediately
### GET /api/automation/log?limit=100 → `{"items":[{"at","ruleName","result","detail"}]}`

---

## Audit Log

### GET /api/audit?limit=100&offset=0 → `{"total":N,"items":[{"at","action","details"}]}`
(log of panel actions: commands, bans, config edits, automation triggers)

---

## Settings

### GET /api/settings → (passwords are masked as `"********"`)
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
`configDir`/`logFile` empty = derived from `dir` (`<dir>/SquadGame/ServerConfig`,
`<dir>/SquadGame/Saved/Logs/SquadGame.log`).
### PUT /api/settings/general | /api/settings/rcon | /api/settings/server
Empty password or `********` = keep the stored value unchanged.
### POST /api/settings/rcon/test → `{"ok":true,"detail":"Подключено, сервер: ..."}` <!-- "Подключено, сервер: ..." = "Connected, server: ..." -->
### POST /api/settings/server/test → `{"ok":true,"detail":"Папка найдена, конфиги: 5 файлов, лог найден"}` <!-- "Папка найдена, конфиги: 5 файлов, лог найден" = "Folder found, configs: 5 files, log found" -->

---

## Utility
### GET /healthz → `ok`
### Static files: `/` — SPA; all unknown paths (except `/api/*`) serve `index.html`.
