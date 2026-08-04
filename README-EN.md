# SquadAdmin — Squad Server Administration Panel

Russian version: README-RU.md

A single executable `SquadAdmin.exe`. No installation required, no dependencies.
Run it on the same machine as your Squad server; it opens the panel
in your browser at `http://127.0.0.1:8080`.

---

## Quick Start

1. Copy `SquadAdmin.exe` into a dedicated folder — for example `C:\SquadAdmin\`.
   **Do not place it in your Squad server folder**: the program creates its own files
   next to the exe.
2. Double-click to launch.
3. Windows will show a SmartScreen warning: "Windows protected your PC" —
   this is normal for programs without a paid digital signature.
   Click **"More info" → "Run anyway"**.
4. The browser will open with the panel. If it does not open automatically, navigate
   manually to `http://127.0.0.1:8080`.
5. No login or password is required — the panel opens immediately.
   It is designed for local use on your own machine and listens only on
   `127.0.0.1`. Do not expose its port to the internet (see "Common Tasks").

The console window that appears on launch **must not be closed** — the panel runs as
long as it is open. To stop the panel, close that window or press `Ctrl+C` in it.

---

## Configuration

Open the **Settings** (Настройки) section in the panel. Two things need to be filled in.

### 1. RCON — required for moderation features

Take the values from your server's configuration file
`<server folder>\SquadGame\ServerConfig\Rcon.cfg`:

```
Password=yourpassword   ← RCON password
Port=21114              ← RCON port
```

Fill in the panel fields:

| Field | Value |
|---|---|
| Address | `127.0.0.1` (panel and server on the same machine) |
| Port | from `Rcon.cfg`, usually `21114` |
| Password | from `Rcon.cfg` |

Click **"Test connection"** — your server name should appear.

> If the test fails: make sure the server is running, that a non-empty password is set
> in `Rcon.cfg`, and that the server was started with the `RCONPORT=<port>` argument.
> Squad does not start RCON if the password in `Rcon.cfg` is empty.

### 2. Server folder — for configs, logs, and process management

| Field | What to enter | Example |
|---|---|---|
| Server folder | Root of your Squad server installation | `C:\servers\squad` |
| Configs folder | Leave blank — auto-detected | `…\SquadGame\ServerConfig` |
| Log file | Leave blank — auto-detected | `…\SquadGame\Saved\Logs\SquadGame.log` |
| Server executable | Path to the server exe, if you want to manage startup from the panel | `C:\servers\squad\SquadGameServer.exe` |
| Launch arguments | Same as in your `.bat` file | `Port=7787 QueryPort=27165 RCONPORT=21114` |

Click **"Check folder"** — the panel will report how many configs it found and whether
the log is visible.

Until the folder is set, the **Configs** (Конфиги) and **Logs** (Логи) sections show
a placeholder. Until the executable is set, the **Process** (Процесс) section is also
unavailable. Everything else — players, console, maps, bans, player database,
automation — works over RCON alone.

---

## What the Panel Can Do

| Section | Features |
|---|---|
| **Overview** (Обзор) | Server status, player count, current and next map, live event feed |
| **Players** (Игроки) | Team and squad tree, search, warn, kick, ban, team switch, disband squad, remove commander |
| **Console** (Консоль) | Any RCON commands with autocomplete suggestions and history |
| **Maps** (Карты) | Layer list, change current and next map, end match, restart match |
| **Bans** (Баны) | Local ban database with expiry times, editing, unbanning, and export in `Bans.cfg` format |
| **Player database** (База игроков) | Nickname history, sessions, total playtime, punishments, admin notes |
| **Configs** (Конфиги) | Edit any `.cfg`, structured editor for `Admins.cfg` (groups and permissions via checkboxes), map rotation editor |
| **Logs** (Логи) | Chat, teamkills, joins and leaves, errors; raw log viewer |
| **Mods** (Моды) | Mod list from a user-specified folder, editing mod settings files with backups |
| **Process** (Процесс) | Start, stop, and restart the server, warn players, auto-restart on crash |
| **Automation** (Автоматизация) | Interval and scheduled messages, scheduled commands and restarts, seeding mode |
| **Journal** (Журнал) | Everything the panel has done: commands, bans, config edits, rule triggers |

---

## Files Next to the Exe

```
C:\SquadAdmin\
├── SquadAdmin.exe
├── config.json     — settings (can be edited manually while the panel is closed)
└── data.db         — database: bans, players, notes, rules, journal
```

**Backup** — simply copy `config.json` and `data.db`.
To move to another computer, copy the entire folder.

Before every write to a server config, the panel creates a backup alongside the
original file: `Server.cfg.bak-20260802-194452`. The last 10 backups are kept.

---

## Common Tasks

**Open the panel from another computer on your home network.**
Settings → General → enable "Listen on all interfaces", save, restart
`SquadAdmin.exe`. Access it at `http://<computer-IP>:8080`.
You may need to allow the port through Windows Firewall.
Do not forward this port to the internet: the panel has no password protection,
so anyone who reaches the port gains full access.
Use it only on a trusted local network.

**Change the panel port.**
Settings → General → "Panel port". Takes effect after restarting the program.
If the panel fails to start because the port is already in use, you can edit the port
directly in `config.json` (field `general.panelPort`) with Notepad.

**Launch the panel with Windows.**
Press `Win+R`, type `shell:startup`, and copy a shortcut to `SquadAdmin.exe`
into the folder that opens.

**Auto-restart the server on crash.**
Process section → toggle "Auto-restart on crash". The panel checks the process
every 5 seconds and restarts it. If the server crashes more than 3 times within
10 minutes, auto-restart stops — this indicates a persistent problem that needs
to be investigated in the logs.

**Seeding.**
Automation → Add rule → "Seeding" (Сидинг). Set an activation threshold (e.g.,
fewer than 20 players), a deactivation threshold (e.g., more than 50), a seeding
layer, and a live layer. The panel will switch the map automatically and remind
players about seeding. The thresholds are intentionally different to prevent the
mode from toggling back and forth at the boundary.

---

## Troubleshooting

**The panel does not open in the browser.** Check that the `SquadAdmin.exe` console
window is open and shows no "ERROR" line. Open `http://127.0.0.1:8080` manually.

**Red banner "RCON not connected".** The server is offline, or the port or password
is wrong, or the server was started without `RCONPORT`. Check with the "Test connection"
button.

**The Logs section is empty.** Squad writes the log only while the server is running.
The panel reads the file from the moment it starts and does not parse earlier entries.

**Antivirus deleted the exe.** Go programs without a digital signature can trigger
heuristic detection. Add the folder `C:\SquadAdmin\` to your antivirus exclusions.

**Cannot ban a player.** If the player is offline, the ban is saved in the panel's
database but is not sent to the server. Export the bans using the "Export Bans.cfg"
button in the Bans section and place the file in `SquadGame\ServerConfig\`.

---

## Technical Details

- A single executable with a built-in web interface. No internet connection required:
  no fonts or libraries are loaded from external sources.
- Storage — SQLite in `data.db`.
- No login/password: the panel is designed for local use and listens on `127.0.0.1`.
  Requests are accepted only from the same machine — `Host`, `Origin`, and
  `Sec-Fetch-Site` headers are verified; everything else receives a 403.
  Access control is enforced at the OS and network level.
- RCON polling — once every 20 seconds (configurable); updates are pushed to the
  interface as a stream, no page reload needed.

---

## What's New in 1.0.1

- **Login removed.** The panel opens immediately — it is designed for local use by
  the server host. The login screen, password change, and session storage have been
  removed (along with the question of storing session tokens in the database).
- **Fixed crash on malformed RCON packet.** A packet with an incorrectly declared
  size could previously crash the process; now such a stream simply reconnects.
- **Prevented double server launch.** Manual "Start" and auto-restart on crash can
  no longer start two processes simultaneously.
- **Honest server stop.** The "Stop" button no longer reports success if the process
  is still running.

Version 1.0.1

---

## What's New in 1.0.2

Important findings from a technical audit have been addressed (critical ones were
fixed in 1.0.1).

- **Panel responds only on this computer.** The `Host`, `Origin`, and `Sec-Fetch-Site`
  headers are checked; requests from other machines and third-party sites receive 403.
  This closes DNS-rebinding and CSRF — especially relevant since there is no password
  login. With `bindAll`, the socket still listens on the network but only accepts
  local requests.
- **Configs are written atomically.** Saving `Admins.cfg`, `Bans.cfg`, `Server.cfg`
  goes through a temporary file followed by a rename. An interrupted write (crash,
  antivirus, power loss) will no longer leave a truncated config.
- **Comments in `Admins.cfg` are preserved.** Standalone comment lines between groups
  and admins no longer turn into blank lines when saved through the structured editor.
- **Automation rule cannot fire twice.** True re-entry protection has been added.
  If the rule's state cannot be saved, the scheduled restart is cancelled and an error
  is logged — instead of risking a second restart.
- **Log rotation is detected by file content.** A recreated `SquadGame.log` is
  recognized even if it has grown back to the previous file position between checks;
  events in the first seconds after a restart are no longer lost. The UTF-8 BOM at
  the start of the file is also skipped.
- **Player list tolerates non-standard lines.** Bots/AI and the changed-post-patch
  Squad format no longer fail silently: lines are parsed with a more lenient pattern,
  and unrecognized entries are written to the journal.
- **Database is ready for upgrades.** Schema versioning has been introduced, so future
  structural changes will apply to the existing `data.db` rather than causing a
  "no such column" error.
- **Interface data no longer flickers.** RCON polling no longer runs in parallel with
  itself after a ban or kick.
- **RCON latency is displayed accurately.** Previously the value shown was "poll time / 5";
  now it is the real round-trip time of a single command.

Version 1.0.2

---

## What's New in 1.0.3

Remaining minor technical audit items resolved.

- **Correct `Cache-Control` for static assets.** The header has been changed from
  `no-cache` to `no-store` — responses are no longer stored in any cache.
- **Search works as expected.** The special characters `%` and `_` in search strings
  no longer match arbitrary text: they are escaped before the LIKE query.
- **Seeding rule validation.** When `switchLayer: true`, the `seedLayer` and
  `liveLayer` fields are required — the error is returned at rule creation time,
  not when it fires.
- **Configs are saved with CRLF.** In "raw" `.cfg` editing, line endings are
  normalized to `\r\n` — required by Squad on Windows.
- **Transactions in key player database operations.** `UpsertSeenPlayer`,
  `SessionStart`, `SessionEnd` are now atomic; session and playtime counters are
  updated together with their main records.
- **Dead parameters removed from `doStopWindows`.** `cmd` and `cmdDead` are POSIX
  objects not applicable on Windows; the function signature has been simplified.
- **Broken Squad packet pattern documented.** The comment now contains the exact byte
  signature and a hint on what to check when future game patches arrive.

Version 1.0.3

---

## What's New in 1.0.4

A cosmetic update: the panel has been restyled with a Squad theme. No logic, API,
or behavior changes — a safe update.

- **"Field Command Post" theme.** A daytime military-khaki palette: olive framework
  (header and menu), a light workspace in staff-paper color, a coordinate grid, and
  a muted field-position background.
- **Military typography.** Headings, buttons, and labels use the condensed Oswald
  font; body text uses Inter; console and identifiers use JetBrains Mono. Fonts are
  bundled in the exe including Cyrillic — no internet connection required.
- **Console styled as a field terminal.** Dark background with phosphor-green text,
  commands highlighted in amber.
- **Menu split into groups** — Operations (Оперативно), Moderation (Модерация),
  Server (Сервер), System (Система); symbols replaced with vector icons that do not
  depend on system fonts.
- **Clock in the top bar** and a version label in the menu footer.
- **Minor usability fixes.** Search fields and dropdowns in toolbars no longer
  stretch full-width; checkbox labels are in normal case; link-buttons
  ("Export Bans.cfg") are no longer underlined.

What stayed the same: URLs, settings, `config.json` and `data.db` format, all
sections and key workflows. The update is simply a replacement of `SquadAdmin.exe`.

Version 1.0.4

---

## What's New in 1.0.5

- **Workspace background selection.** Settings → Appearance (Настройки → Оформление)
  now offers ten built-in themed backgrounds: field position (default), command post,
  staff map, comms node, night camp, strongpoint, helicopter pad, column on the march,
  forest line, and ruined city. All images are bundled in the exe — no internet
  connection required.
- **Day / Night mode.** The sun/moon button in the top bar instantly toggles the panel
  between light and dark themes; the same setting is available in Settings →
  Appearance. The dark theme is a "headquarters at night": dark workspace, muted
  lines, readable status indicators.
- **Selection is remembered.** The theme and background are saved in the browser and
  applied before the first render — the panel no longer flashes with the light theme
  on load.
- **File validation before writing.** Everything added or changed through the panel
  is now validated based on file type and extension:
  - only `*.cfg` files can be edited (exe, dll, bat, and others are rejected);
  - binary data, non-UTF-8 content, and files larger than 1 MB are not written;
  - `Admins.cfg` — line-by-line check of `Group=…`/`Admin=…`, SteamID64/EOS
    identifiers, and permission tokens;
  - rotation and exclusion list files — layer name validation;
  - `Bans.cfg` — `ID:expiry` format with optional comment;
  - `Server.cfg` and `Rcon.cfg` — `Parameter=value` format;
  - `RemoteAdminListHosts.cfg` / `RemoteBanListHosts.cfg` — only `http(s)://`
    addresses;
  - free-text files (`MOTD.cfg`, `ServerMessages.cfg`, etc.) — no strict format
    restrictions.
  Errors are returned with the line number and a clear explanation; the file on disk
  is not touched (writes are still atomic with a backup).
- **Path validation in server settings.** "Path to exe" accepts only `.exe` files;
  "Log file" accepts only `.log`/`.txt`; paths must be absolute and must not contain
  `..`.

The `config.json`, `data.db` format, and API did not change. The update is simply
a replacement of `SquadAdmin.exe`.

Version 1.0.5

---

## What's New in 1.0.6

### Interface Language Switching

A language switcher — Russian / English — has been added to the panel.

**How to switch:**
- Quick button **"RU"** / **"EN"** in the top bar next to the theme toggle —
  one click changes the language and reloads the page.
- Or: **Settings → Appearance → Interface language** (Настройки → Оформление →
  Язык интерфейса) — buttons "Русский" / "English".

**What is translated in English mode:**
all interface sections, background preset names, console command suggestions,
time units, error messages, and backend-generated messages — Journal action
entries (including historical ones), automation log details, server
notifications, and game-log events (joins, leaves, teamkills).

**What remains in Russian:**
only free-form user text originally typed in Russian (ban reasons, notes,
broadcast messages) and any strings the panel does not recognize — those are
shown as-is.

**The selection is saved in the browser** (localStorage, key `sa_lang`, default
Russian) — just like the theme and background. The setting is independent on each
browser and computer.

The `config.json`, `data.db` format, and API did not change. The update is simply
a replacement of `SquadAdmin.exe`.

Version 1.0.6

---

## What's New in 1.0.7

### Chat Message Colors by Role

Chat messages can now be highlighted with a color based on the sender's
role — a group from `Admins.cfg` (SuperAdmin, Moderator, etc.).

**How to configure:** **Settings → Chat colors**.
- A master switch "Colorize chat messages by sender role".
- A custom color for every group from `Admins.cfg`: tick the checkbox next to
  a role and pick a color. The group list is read from the file automatically.
- A separate color for the **ChatAdmin** channel (admin chat) — applied to
  messages in that channel when the sender has no role color.

**Where the highlighting is visible:** in the event feed on the dashboard and
in **Logs → Events** (both live messages and history). Hovering over a
highlighted message shows the role name. Messages from players without a role
keep the default color.

The role is resolved from the sender's SteamID/EOS id at the moment of the
message and stored in the journal. Color changes apply immediately in all open
tabs — no page reload needed.

The setting is stored in `config.json` (the `chatColors` section). On first
start a `role` column is added to the event journal in `data.db`
automatically — the update is still just a replacement of `SquadAdmin.exe`.

Version 1.0.7

---

## What's New in 1.0.8

### "Mods" Section

A new **Mods** section for servers running mods from the Steam Workshop and
beyond.

**How to configure:** **Settings → Mods** — set the path to the folder that
contains your mods (e.g. `...\steamapps\workshop\content\393380` or a custom
folder). The **"Check folder"** button immediately shows how many mods were
found. Each subfolder inside is treated as a separate mod.

**What the section shows:**
- A list of all mods: name (read from `.uplugin`/`mod.json` when present),
  folder, version, and settings status.
- A mod with the **"available"** status has settings files — click the mod to
  open the file list, then click any file to edit it.
- The **"will appear after the game starts"** status means the mod has not
  created its settings files yet. Most mods create them on the first server
  launch with the mod — start the server and the panel re-scans the folder
  automatically (the list refreshes on every server process start and stop).

**Editing:** `.cfg`, `.ini`, `.json`, `.txt`, and `.config` files are
supported. Before every write a backup is created next to the file
(`<file>.bak-<time>`, the last 10 are kept), and the write is atomic — same
as for server configs. Escaping the mod folder through the file path is not
possible.

The mods folder path is stored in `config.json` (the `server.modsDir` field).
The update is still just a replacement of `SquadAdmin.exe`.

Version 1.0.8
