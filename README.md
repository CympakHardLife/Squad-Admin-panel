<div align="center">

# 🎖️ SquadAdmin

### Squad Server Administration Panel — one `.exe`, zero dependencies

**English** · [Русский](README.ru.md)

[![CI](https://github.com/CympakHardLife/Squad-Admin-panel/actions/workflows/ci.yml/badge.svg)](https://github.com/CympakHardLife/Squad-Admin-panel/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-22c55e.svg?style=for-the-badge)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.25-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev)
[![Platform](https://img.shields.io/badge/Platform-Windows-0078D6?style=for-the-badge&logo=windows&logoColor=white)]()
[![Release](https://img.shields.io/badge/Release-v1.0.8-8b5cf6?style=for-the-badge)](../../releases)
[![Language](https://img.shields.io/badge/UI-RU%20%2F%20EN-f59e0b?style=for-the-badge)]()

A complete web-based control panel for your **Squad** dedicated server.
Players, RCON console, maps, bans, configs, logs, mods, process control and automation —
all in a single portable executable that runs on the same machine as your server.

</div>

---

## ✨ Highlights

<table>
<tr>
<td width="33%" valign="top">

### 📦 One file
`SquadAdmin.exe` — no installer, no runtime, no Node, no Python.
Drop it in a folder and double-click.

</td>
<td width="33%" valign="top">

### 🔌 Fully offline
Fonts, styles and scripts are embedded in the binary.
Nothing is loaded from the internet.

</td>
<td width="33%" valign="top">

### 🌍 RU / EN interface
Switch languages with one click in the top bar.
Journal, logs and notifications are translated too.

</td>
</tr>
</table>

---

## 📸 Screenshots

<div align="center">
<img src="docs/screenshots/1.png" width="900" alt="SquadAdmin — overview">
<br><br>
<img src="docs/screenshots/2.png" width="900" alt="SquadAdmin — players">
<br><br>
<img src="docs/screenshots/3.png" width="900" alt="SquadAdmin — configs">
</div>

---

## 🚀 Quick Start

```text
1. Copy SquadAdmin.exe into its own folder, e.g. C:\SquadAdmin\
   WARNING: do NOT put it inside your Squad server folder — it creates its own files.
2. Double-click to run.
3. SmartScreen warning → "More info" → "Run anyway"
   (normal for apps without a paid code-signing certificate).
4. Your browser opens at  http://127.0.0.1:8080
5. No login required — the panel listens on localhost only.
```

> ⚠️ **Keep the console window open** — it is the running panel.
> Close it or press `Ctrl+C` to stop.

**Updating from an earlier version:** stop the panel, replace `SquadAdmin.exe`, start it again.
`config.json` and `data.db` are migrated automatically — your bans, notes and settings stay in place.
Copy `data.db` somewhere safe first if you like to keep a way back (see the note in 1.0.7 below).

---

## ⚙️ Configuration

Everything is set up inside **Settings** in the panel. Only two things are required.

<details>
<summary><b>1 · RCON — required for all moderation features</b></summary>
<br>

Take the values from `<server folder>\SquadGame\ServerConfig\Rcon.cfg`:

```ini
Password=yourpassword   ← RCON password
Port=21114              ← RCON port
```

| Field | Value |
|---|---|
| Address | `127.0.0.1` (panel and server on the same machine) |
| Port | from `Rcon.cfg`, usually `21114` |
| Password | from `Rcon.cfg` |

Click **Test connection** — your server name should appear.

> Squad does **not** start RCON if the password in `Rcon.cfg` is empty, and the server
> must be launched with the `RCONPORT=<port>` argument.

</details>

<details>
<summary><b>2 · Server folder — for configs, logs and process control</b></summary>
<br>

| Field | What to enter | Example |
|---|---|---|
| Server folder | Root of the Squad server install | `C:\servers\squad` |
| Configs folder | Leave blank — auto-detected | `…\SquadGame\ServerConfig` |
| Log file | Leave blank — auto-detected | `…\SquadGame\Saved\Logs\SquadGame.log` |
| Server executable | Needed for start/stop from the panel | `C:\servers\squad\SquadGameServer.exe` |
| Launch arguments | Same as in your `.bat` | `Port=7787 QueryPort=27165 RCONPORT=21114` |

Click **Check folder** — the panel reports how many configs it found and whether the log is visible.

Everything else — players, console, maps, bans, player database, automation —
works over RCON alone.

</details>

<details>
<summary><b>3 · Mods folder — optional, for the Mods section</b></summary>
<br>

Settings → **Mods** → **Mods folder**. This path is **not** auto-detected: everyone keeps
their mods somewhere else, so you point the panel at the folder yourself.

| Setup | Typical path |
|---|---|
| Steam Workshop | `…\steamapps\workshop\content\393380` |
| Manual / custom folder | whatever folder you download mods into |

Every subfolder inside is treated as a separate mod. Click **Check folder** — the panel
immediately reports how many mods it found.

Leave this empty if your server runs vanilla; the **Mods** section will simply invite
you to set the path.

</details>

<details>
<summary><b>4 · Chat colors — optional, cosmetic</b></summary>
<br>

Settings → **Chat colors**. Highlights chat messages by the sender's role, so an admin
line no longer looks like every other line in the feed.

| Setting | What it does |
|---|---|
| Colorize chat messages by sender role | Master switch for the whole feature |
| Roles from `Admins.cfg` | Tick a group and pick its color — the group list is read from the file |
| ChatAdmin channel (admin chat) | Fallback color for the admin channel when the sender has no role color |

Requires the configs folder to be set (that is where `Admins.cfg` lives). If the list of
groups is empty, the panel says so: *"No groups found: set the server folder and fill in
`Admins.cfg`"*.

</details>

---

## 🧩 Features

| Section | What it does |
|---|---|
| 📊 **Overview** | Server status, player count, current and next map, live event feed |
| 👥 **Players** | Team and squad tree, search, warn, kick, ban, team switch, disband squad, remove commander |
| 💻 **Console** | Any RCON command with autocomplete suggestions and history |
| 🗺️ **Maps** | Layer list, change current and next layer, end match, restart match |
| 🚫 **Bans** | Local ban database with expiry, editing, unban and export in `Bans.cfg` format |
| 🗂️ **Player database** | Nickname history, sessions, total playtime, punishments, admin notes |
| 📝 **Configs** | Edit any `.cfg`, structured `Admins.cfg` editor (groups & permissions via checkboxes), rotation editor |
| 📜 **Logs** | Chat, teamkills, joins and leaves, errors, plus a raw log viewer — chat lines colored by admin role |
| 🧱 **Mods** | List of mods from your mods folder with names and versions, plus an editor for their settings files with backups |
| ⚡ **Process** | Start / stop / restart the server, warn players, auto-restart on crash |
| 🤖 **Automation** | Interval and scheduled messages, scheduled commands and restarts, seeding mode |
| 🧾 **Journal** | Every action the panel took: commands, bans, config edits, rule triggers |

---

## 📁 Files Created Next to the Exe

```text
C:\SquadAdmin\
├── SquadAdmin.exe
├── config.json     — settings (editable by hand while the panel is closed)
└── data.db         — SQLite: bans, players, notes, rules, journal
```

**Backup** = copy `config.json` + `data.db`.

Before every write to a server config the panel makes a backup next to the original
(`Server.cfg.bak-20260802-194452`); the last 10 are kept. Mod settings files work exactly
the same way, except the backups appear next to the mod's own file inside the mods folder.

---

## 🛠️ Common Tasks

<details>
<summary><b>Open the panel from another PC on the local network</b></summary>
<br>

Settings → General → enable **Listen on all interfaces** → save → restart `SquadAdmin.exe`.
Open `http://<computer-IP>:8080`. You may need to allow the port in Windows Firewall.

> 🔒 **Never forward this port to the internet.** The panel has no password —
> anyone who reaches the port gets full control. Trusted LAN only.

</details>

<details>
<summary><b>Change the panel port</b></summary>
<br>

Settings → General → **Panel port**, then restart the program.
If the panel won't start because the port is busy, edit `general.panelPort` in
`config.json` with Notepad.

</details>

<details>
<summary><b>Launch the panel with Windows</b></summary>
<br>

`Win+R` → `shell:startup` → put a shortcut to `SquadAdmin.exe` in the folder that opens.

</details>

<details>
<summary><b>Auto-restart the server on crash</b></summary>
<br>

Process → toggle **Auto-restart on crash**. The panel checks the process every 5 seconds.
If the server crashes more than 3 times in 10 minutes, auto-restart stops — that means a
persistent problem worth investigating in the logs.

</details>

<details>
<summary><b>Seeding mode</b></summary>
<br>

Automation → Add rule → **Seeding**. Set an activation threshold (e.g. under 20 players),
a deactivation threshold (e.g. over 50), a seeding layer and a live layer. The panel
switches the map automatically and reminds players. The two thresholds differ on purpose
so the mode does not flip back and forth at the boundary.

</details>

<details>
<summary><b>Make admin messages stand out in the chat feed</b></summary>
<br>

1. Make sure the configs folder is set and `Admins.cfg` contains your groups.
2. Settings → **Chat colors** → enable **Colorize chat messages by sender role**.
3. Tick the groups you care about (e.g. `SuperAdmin`, `Moderator`) and pick a color for each.
4. Optionally set a color for the **ChatAdmin** channel — it applies to admin-chat messages
   from senders without a role color.

Colors apply instantly in every open tab. Hover a highlighted message to see the role name.

</details>

<details>
<summary><b>Edit a mod's settings file</b></summary>
<br>

Settings → Mods → set the **Mods folder** → open the **Mods** section.

- A mod marked **available** already has settings files — click the mod, then the file.
- A mod marked **will appear after the game starts** has not generated its files yet.
  Most mods create them on the first server launch with the mod loaded. Start the server,
  then come back — the list refreshes on every server start and stop, or press **Refresh**.

`.cfg`, `.ini`, `.json`, `.txt` and `.config` files can be edited. The panel never creates
new files: the mod itself must generate them.

</details>

---

## 🩺 Troubleshooting

| Problem | Solution |
|---|---|
| Panel doesn't open in the browser | Check the console window is open with no `ERROR` line; open `http://127.0.0.1:8080` manually |
| Red banner "RCON not connected" | Server offline, wrong port/password, or started without `RCONPORT`. Use **Test connection** |
| Logs section is empty | Squad writes the log only while running; the panel reads from its own start, not earlier entries |
| Antivirus deleted the exe | Unsigned Go binaries trigger heuristics — add `C:\SquadAdmin\` to exclusions |
| Can't ban an offline player | The ban is stored locally but not sent to the server — use **Export Bans.cfg** and place it in `SquadGame\ServerConfig\` |
| Mods section says the path is not set | Settings → Mods → **Mods folder**, then **Check folder** — it must report the number of mods found |
| A mod shows "will appear after the game starts" | Normal. The mod has not created its settings files yet — start the server with the mod loaded, then press **Refresh** |
| "No mods found in folder" | You pointed at the mod itself instead of the folder that *contains* mods. Each subfolder = one mod |
| Chat colors don't apply | The configs folder is not set, `Admins.cfg` has no groups, or the group is unticked in Settings → Chat colors |
| Panel won't start: "data.db was created by a newer version" | You rolled back to an older build after running 1.0.7+. Restore your `data.db` backup or go back to the newer exe |

---

## 🔬 Technical Details

- **Go 1.25**, single static binary with an embedded web interface.
- **Storage** — SQLite (`data.db`) via the pure-Go driver `modernc.org/sqlite` — no CGO, no DLLs.
  The schema is versioned (`PRAGMA user_version`) and migrated forward on start.
- **No login/password by design** — bound to `127.0.0.1`; `Host`, `Origin` and `Sec-Fetch-Site`
  headers are validated, everything else gets `403`. Access control belongs to the OS/network layer.
- **RCON polling** every 20 seconds (configurable); updates are streamed to the UI — no page reloads.
- **Safe config writes** — validation with line numbers, atomic writes, automatic backups,
  strict path validation (`.exe` / `.log` only, absolute paths, no `..`).
- **Role resolution** — `Admins.cfg` is parsed lazily and re-read only when its size or
  modification time changes; sender SteamID64 / EOS id is matched against `Admin=<id>:<group>`.
- **Mod scanning is bounded** — depth 6, up to 64 settings files per mod, up to 500 mods, 2 MB per
  editable file; `Content`, `Binaries`, `Intermediate` and `DerivedDataCache` are skipped so a
  multi-gigabyte Workshop folder never stalls the panel.

---

## 🆕 What's New in 1.0.8

### 🧱 Mods section

A new **Mods** section for servers running mods from the Steam Workshop and beyond.

**How to configure:** **Settings → Mods** — set the path to the folder that contains your mods
(e.g. `…\steamapps\workshop\content\393380` or a custom folder). The **Check folder** button
immediately shows how many mods were found. Each subfolder inside is treated as a separate mod.

**What the section shows:**

- A list of all mods: name (read from `.uplugin` / `mod.json` when present), folder, version,
  and settings status.
- A mod with the **available** status has settings files — click the mod to open the file list,
  then click any file to edit it.
- The **will appear after the game starts** status means the mod has not created its settings
  files yet. Most mods create them on the first server launch with the mod — start the server and
  the panel re-scans the folder automatically (the list refreshes on every server process start
  and stop).

**Editing:** `.cfg`, `.ini`, `.json`, `.txt` and `.config` files are supported. Before every write
a backup is created next to the file (`<file>.bak-<time>`, the last 10 are kept), and the write is
atomic — same as for server configs. Escaping the mod folder through the file path is not possible,
and the panel never creates new files: the mod itself must generate them.

The mods folder path is stored in `config.json` (the `server.modsDir` field).
The update is still just a replacement of `SquadAdmin.exe`.

> 📌 **Coming from 1.0.6?** Version 1.0.7 was folded into this release — it added
> chat message colors by admin role. Expand the section below to see what you also get.

<details>
<summary><h3>🎨 What's New in 1.0.7 — Chat message colors by role</h3></summary>
<br>

Chat messages can now be highlighted with a color based on the sender's role — a group from
`Admins.cfg` (`SuperAdmin`, `Moderator`, etc.).

**How to configure:** **Settings → Chat colors**.

- A master switch *"Colorize chat messages by sender role"*.
- A custom color for every group from `Admins.cfg`: tick the checkbox next to a role and pick a
  color. The group list is read from the file automatically.
- A separate color for the **ChatAdmin** channel (admin chat) — applied to messages in that
  channel when the sender has no role color.

**Where the highlighting is visible:** in the event feed on the dashboard and in **Logs → Events**
(both live messages and history). Hovering over a highlighted message shows the role name.
Messages from players without a role keep the default color.

The role is resolved from the sender's SteamID / EOS id at the moment of the message and stored in
the journal. Color changes apply immediately in all open tabs — no page reload needed.

The setting is stored in `config.json` (the `chatColors` section). On first start a `role` column
is added to the event journal in `data.db` automatically — the update is still just a replacement
of `SquadAdmin.exe`.

> ⚠️ **One-way database upgrade.** Once 1.0.7 or newer has opened your `data.db`, builds up to
> 1.0.6 will refuse to start on it ("created by a newer version of the panel"). Copy `data.db`
> before updating if you want a guaranteed way back.

</details>

<details>
<summary><h3>🌍 What's New in 1.0.6 — Interface language switching</h3></summary>
<br>

A **RU / EN** switch was added next to the theme toggle in the top bar (also in
Settings → Appearance → Interface language).

- **Translated:** all UI sections, background preset names, console command suggestions, time
  units, error messages, and backend-generated text — journal entries (including historical ones),
  automation log details, server notifications and game-log events (joins, leaves, teamkills).
- **Stays in Russian:** free-form text you typed yourself (ban reasons, notes, broadcasts).
- **Saved in the browser** (`localStorage`, key `sa_lang`, default Russian), per browser and per
  machine.

`config.json`, the `data.db` format and the API are unchanged.

</details>

---

## 👤 Author

<div align="center">

### **Cympak** — [@CympakHardLife](https://github.com/CympakHardLife)

*"I'm not a master at anything, I'm a simple person without education, I do what I like."*

**This program was written by me with the help of [Fable 5](https://www.anthropic.com).**

</div>

---

## 📄 License

Released under the **MIT License** — see [LICENSE](LICENSE).

You may freely use, copy, modify, merge, publish, distribute and even sell this software,
including in commercial projects. The only requirement is that the copyright notice and this
license text remain in copies. The software is provided **"as is"**, without warranty of any kind.

```text
Copyright (c) 2026 Cympak (CympakHardLife)
```

---

<div align="center">

⭐ **If SquadAdmin helps you run your server, consider starring the repository.**

</div>
