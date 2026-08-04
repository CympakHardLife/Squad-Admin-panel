<div align="center">

# 🎖️ SquadAdmin

### Squad Server Administration Panel — one `.exe`, zero dependencies

**English** · [Русский](README.ru.md)

[![License: MIT](https://img.shields.io/badge/License-MIT-22c55e.svg?style=for-the-badge)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.25-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev)
[![Platform](https://img.shields.io/badge/Platform-Windows-0078D6?style=for-the-badge&logo=windows&logoColor=white)]()
[![Release](https://img.shields.io/badge/Release-v1.0.6-8b5cf6?style=for-the-badge)](../../releases)
[![Language](https://img.shields.io/badge/UI-RU%20%2F%20EN-f59e0b?style=for-the-badge)]()

A complete web-based control panel for your **Squad** dedicated server.
Players, RCON console, maps, bans, configs, logs, process control and automation —
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

## 🚀 Quick Start

```text
1. Copy SquadAdmin.exe into its own folder, e.g. C:\SquadAdmin\
   ⚠️  Do NOT put it inside your Squad server folder — it creates its own files.
2. Double-click to run.
3. SmartScreen warning → "More info" → "Run anyway"
   (normal for apps without a paid code-signing certificate).
4. Your browser opens at  http://127.0.0.1:8080
5. No login required — the panel listens on localhost only.
```

> ⚠️ **Keep the console window open** — it is the running panel.
> Close it or press `Ctrl+C` to stop.

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
| 📜 **Logs** | Chat, teamkills, joins and leaves, errors, plus a raw log viewer |
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
(`Server.cfg.bak-20260802-194452`); the last 10 are kept.

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

---

## 🩺 Troubleshooting

| Problem | Solution |
|---|---|
| Panel doesn't open in the browser | Check the console window is open with no `ERROR` line; open `http://127.0.0.1:8080` manually |
| Red banner "RCON not connected" | Server offline, wrong port/password, or started without `RCONPORT`. Use **Test connection** |
| Logs section is empty | Squad writes the log only while running; the panel reads from its own start, not earlier entries |
| Antivirus deleted the exe | Unsigned Go binaries trigger heuristics — add `C:\SquadAdmin\` to exclusions |
| Can't ban an offline player | The ban is stored locally but not sent to the server — use **Export Bans.cfg** and place it in `SquadGame\ServerConfig\` |

---

## 🔬 Technical Details

- **Go 1.25**, single static binary with an embedded web interface.
- **Storage** — SQLite (`data.db`) via the pure-Go driver `modernc.org/sqlite` — no CGO, no DLLs.
- **No login/password by design** — bound to `127.0.0.1`; `Host`, `Origin` and `Sec-Fetch-Site`
  headers are validated, everything else gets `403`. Access control belongs to the OS/network layer.
- **RCON polling** every 20 seconds (configurable); updates are streamed to the UI — no page reloads.
- **Safe config writes** — validation with line numbers, atomic writes, automatic backups,
  strict path validation (`.exe` / `.log` only, absolute paths, no `..`).

---

## 🆕 What's New in 1.0.6

**Interface language switching.** A **RU / EN** switch was added next to the theme toggle in the
top bar (also in Settings → Appearance → Interface language).

- **Translated:** all UI sections, background preset names, console command suggestions, time units,
  error messages, and backend-generated text — journal entries (including historical ones),
  automation log details, server notifications and game-log events (joins, leaves, teamkills).
- **Stays in Russian:** free-form text you typed yourself (ban reasons, notes, broadcasts).
- **Saved in the browser** (`localStorage`, key `sa_lang`, default Russian), per browser and per machine.

`config.json`, the `data.db` format and the API are unchanged — updating is simply replacing `SquadAdmin.exe`.

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
