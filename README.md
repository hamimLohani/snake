# 🐍 snake — Developer-Grade Terminal Snake CLI

[![Release](https://img.shields.io/github/v/release/hamimLohani/snake?color=brightgreen)](https://github.com/hamimLohani/snake/releases)
[![CI](https://github.com/hamimLohani/snake/actions/workflows/ci.yml/badge.svg)](https://github.com/hamimLohani/snake/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/hamimLohani/snake)](https://goreportcard.com/report/github.com/hamimLohani/snake)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> **A coffee-break game for developers.** Installs in one command, launches in under 5ms, looks polished in any terminal, and engineered with zero dependencies to production-grade Go standards.

```
       ╔══════════════════════════════════════════════════════╗
       ║   [CLASSIC]  SCORE: 0140    HIGH: 0320  LEN: 14      ║
       ║  ┌────────────────────────────────────────────────┐  ║
       ║  │ ·  ·  ·  ·  ·  ·  ·  ·  ·  ·  ·  ·  ·  ·  ·  · │  ║
       ║  │ ·  ·  ·  ·  ·  ★  ·  ·  ·  ·  ·  ·  ·  ·  ·  · │  ║
       ║  │ ·  ·  ·  ·  ·  ·  ·  ·  ·  ●  ●  ●  ▶  ·  ·  · │  ║
       ║  │ ·  ·  ·  ·  ·  ·  ·  ·  ·  ●  ·  ·  ·  ·  ·  · │  ║
       ║  │ ·  ·  ·  ·  ·  ·  ·  ·  ·  ●  ●  ·  ·  ·  ·  · │  ║
       ║  └────────────────────────────────────────────────┘  ║
       ║   [W/A/S/D / Arrows] Move • [P] Pause • [B] Boss Key  ║
       ╚══════════════════════════════════════════════════════╝
```

---

## ✨ Features

- **Blazing Fast**: Sub-5ms startup, zero external runtime dependencies, 0-flicker double-buffered rendering.
- **Terminal Safety**: Full alternate-screen buffer isolation; cleanly cleans up raw mode and restores terminal state on exit, SIGINT, or panic.
- **Strictly Deterministic Engine**: Pure state machine with seeded PRNG; enables 100% reproducible replays, headless simulation, and automated verification.
- **Interactive Main Menu & Live Sample Board**: Sleek two-column launch screen with dynamic live board preview for board sizes (Small, Medium, Classic, Large), visual themes, obstacles, and difficulty.
- **4 Distinct Game Modes**: Classic, Time Attack, Survival (shrinking borders down to 6×6), and Zen (infinite toroidal wrap-around).
- **Dynamic Power-Ups & Avoidable Hazards**: Permanent food with concurrent bonus items: Slow-Mo, Ghost (with 10-second Zen wall-wrap & phasing), 3x Multiplier, Shrink, and avoidable Poison hazard apples.
- **Discreet Boss Key**: Press `b` or `esc` to instantly swap the screen to a convincing fake `htop`, compiler build log, or `git log`.
- **Procedural Obstacles**: Deterministic obstacle mazes generated from any shareable `--seed` with a guaranteed 7×5 snake safety clearance zone.
- **Local Persistence**: Leaderboards, high scores, play time, achievements (`~/.config/snake/stats.json`), and rebindable TOML config (`~/.config/snake/config.toml`).
- **Autonomous AI Bot**: Heuristic BFS autopilot with tail-following survival logic (>7,000 ticks/sec in headless mode).
- **IoT Hardware Mode**: Control the snake over WiFi UDP or USB Serial using an ESP8266 or ESP32 joystick.

---

## 🚀 Installation

### Option 1: Quick Install (Linux & macOS)
Install the latest pre-compiled binary with a single command:
```bash
curl -fsSL https://raw.githubusercontent.com/hamimLohani/snake/main/install.sh | bash
```

### Option 2: Homebrew (macOS & Linux)
```bash
brew tap hamimlohani/tap
brew install snake
```

### Option 3: Direct Binary Downloads
Pre-compiled standalone binaries with zero dependencies:

#### macOS (Apple Silicon & Intel)
```bash
# Apple Silicon (M1/M2/M3/M4):
curl -sL https://github.com/hamimLohani/snake/releases/latest/download/snake-darwin-arm64 -o snake
chmod +x snake && sudo mv snake /usr/local/bin/

# Intel (x86_64):
curl -sL https://github.com/hamimLohani/snake/releases/latest/download/snake-darwin-amd64 -o snake
chmod +x snake && sudo mv snake /usr/local/bin/
```

#### Linux (x86_64 & ARM64)
```bash
# Linux x86_64 (amd64):
curl -sL https://github.com/hamimLohani/snake/releases/latest/download/snake-linux-amd64 -o snake
chmod +x snake && sudo mv snake /usr/local/bin/

# Linux ARM64 (aarch64 / Raspberry Pi / AWS Graviton):
curl -sL https://github.com/hamimLohani/snake/releases/latest/download/snake-linux-arm64 -o snake
chmod +x snake && sudo mv snake /usr/local/bin/
```

#### Windows
Download `snake-windows-amd64.exe` from the [Latest Release](https://github.com/hamimLohani/snake/releases/latest), or download via PowerShell:
```powershell
Invoke-WebRequest -Uri https://github.com/hamimLohani/snake/releases/latest/download/snake-windows-amd64.exe -OutFile snake.exe
```

#### Release Tarball Archive (.tar.gz with Checksums)
```bash
curl -sL https://github.com/hamimLohani/snake/releases/latest/download/snake-v1.0.9-linux-amd64.tar.gz | tar -xz
sudo mv snake-linux-amd64 /usr/local/bin/snake
```

### Option 4: Go Native Package Manager
```bash
go install github.com/hamimLohani/snake/cmd/snake@latest
```

### Option 5: Build From Source
```bash
git clone https://github.com/hamimLohani/snake.git
cd snake
make build
# Binary created at ./bin/snake
sudo make install
```

---

## 🎮 CLI Commands & Options

### Subcommands Overview
| Command | Description |
|---|---|
| `snake` or `snake play` | Launch interactive Main Menu with live sample board preview, or start direct gameplay |
| `snake daily` | Compete on today's unified global seed (identical layout for all players worldwide) |
| `snake stats` | View high scores, personal metrics, achievements, and leaderboard (`--json`, `--csv`) |
| `snake replay` | Replay the latest game or any recorded `.snk` run (`--speed`, `--verify`) |
| `snake bot` | Watch the autonomous BFS AI bot play, or run headless benchmarks (`--headless`) |
| `snake config` | Inspect configuration and keybindings (`--edit` in $EDITOR, `--reset` to defaults) |

### Quick Start Examples
```bash
# Launch interactive Main Menu & Settings (Default)
snake
snake play

# Start playing directly with shorthand flags (bypasses menu)
snake play --direct
snake play -m time-attack -t dark -d hard
snake play -m survival -s small
snake play -m zen -t corporate
snake play -s 30x15 -d insane

# Play with a deterministic shared seed and procedural obstacles
snake play --seed 42109 --obstacles

# Daily challenge (unified global seed for everyone worldwide today)
snake daily
snake daily -t colorblind

# View local leaderboard, player metrics, and unlocked achievements
snake stats
snake stats --json
snake stats --csv > leaderboard.csv

# Replay previous games or verify determinism
snake replay                      # Plays latest game
snake replay my_run.snk --speed 40 # 40ms tick interval
snake replay my_run.snk --verify  # Headless determinism check

# Autonomous AI bot demo & headless stress test
snake bot                         # Watch bot play interactively
snake bot --headless --ticks 1000 # Headless stress test / benchmark (>4,000 ticks/sec)

# Manage configuration
snake config                      # View current config & keybindings
snake config --edit               # Open ~/.config/snake/config.toml in $EDITOR
snake config --reset              # Reset to defaults
```

### CLI Flags Reference
| Flag | Shorthand | Description | Default | Applicable Commands |
|---|---|---|---|---|
| `--mode` | `-m` | Game mode: `classic`, `time-attack`, `survival`, `zen` | `classic` | `play` |
| `--difficulty` | `-d` | Speed preset: `easy` (160ms), `normal` (110ms), `hard` (75ms), `insane` (45ms) | `normal` | `play` |
| `--size` | `-s` | Board size preset (`small`, `medium`, `classic`, `large`) or custom `WxH` (e.g. `30x15`) | `classic` (24×20) | `play` |
| `--theme` | `-t` | Visual theme: `dark`, `light`, `high-contrast`, `colorblind`, `corporate` | `dark` | `play`, `daily`, `bot` |
| `--direct` | | Bypass the interactive start menu and jump straight into gameplay | `false` | `play` |
| `--obstacles` | | Enable deterministic procedural obstacle mazes | `false` | `play` |
| `--seed` | | Seed for pseudo-random number generator (`0` = random) | `0` | `play`, `bot` |
| `--speed` | | Override tick interval directly in milliseconds (e.g. `60`, `150`) | from difficulty | `play`, `replay`, `bot` |
| `--resume` | | Auto-resume previous unfinished game if available | `true` | `play` |
| `--screen-reader` | | Accessible plain-text stream mode | `false` | `play` |
| `--hardware-udp` | | Listen on UDP port (e.g. `:4242`) for ESP8266/ESP32 joystick packets | | `play` |
| `--headless` | | Run AI bot simulation headlessly without terminal UI rendering | `false` | `bot` |
| `--ticks` | | Number of ticks to execute during headless simulation | `500` | `bot` |
| `--verify` | | Verify replay determinism headlessly against recorded state | `false` | `replay` |
| `--json` | | Output stats and leaderboard in JSON format | `false` | `stats` |
| `--csv` | | Export leaderboard in standard CSV format | `false` | `stats` |
| `--edit` | | Open `~/.config/snake/config.toml` in `$EDITOR` | `false` | `config` |
| `--reset` | | Reset configuration to factory defaults | `false` | `config` |

---

## ⌨️ Controls & Keybindings

All keys are fully rebindable in `~/.config/snake/config.toml`:

| Action | Primary Keys | Vim / Alternative | Context |
|---|---|---|---|
| **Turn Up** | `↑` (Arrow Up) | `W`, `K` | Gameplay / Menu Navigation |
| **Turn Down** | `↓` (Arrow Down) | `S`, `J` | Gameplay / Menu Navigation |
| **Turn Left** | `←` (Arrow Left) | `A`, `H` | Gameplay / Adjust Menu Setting |
| **Turn Right** | `→` (Arrow Right) | `D`, `L` | Gameplay / Adjust Menu Setting |
| **Pause / Resume** | `P` | `Space` | While playing |
| **Game Settings Menu** | `M` | `Tab` | Anytime (Active, Paused, Game Over) |
| **Instant Boss Key** | `B` | `Esc` | Instantly swaps screen to fake dashboard |
| **Restart Game** | `R` | — | When Game Over |
| **Start / Confirm** | `Enter` | `Space` | In Main Menu |
| **Quit Game** | `Q` | `Ctrl+C` | Anytime (auto-saves unfinished game) |

---

## 🕹️ Game Modes

> 📖 For detailed mechanics, speed curves, scoring formulas, and collision physics, see the complete [**Modes, Difficulties & Rulebook (`MODES_AND_RULES.md`)**](MODES_AND_RULES.md).

| Mode | Command Flag | Description |
|---|---|---|
| **Classic** | `-m classic` | Standard snake mechanics; eat food, grow, avoid walls, obstacles, and your own tail. |
| **Time Attack** | `-m time-attack` | 60-second starting timer (+3s per food). Base points doubled to 25 pts. Race against the clock! |
| **Survival** | `-m survival` | Arena perimeter contracts inward every 50 ticks down to 6×6! Maneuver in an ever-shrinking arena. |
| **Zen** | `-m zen` | No deaths, no game-overs. Infinite toroidal wall wrap-around and self-overlap phasing. |
| **Daily Challenge** | `snake daily` | Compete on today's worldwide unified seed. Identical obstacle and item layout for all players. |

---

## ⚡ Power-Ups & Special Items
 
One food apple exists at a time and immediately respawns when eaten ("one after another"). Alongside the food, **2–3 featured special items appear simultaneously** — 1 guaranteed poison hazard + 1–2 positive power-ups (5.0 second lifespan each, then rotate to a new location):
 
- ⏱ **Slow-Mo** (Cyan): Halves tick speed for **5.0 seconds**.
- 👻 **Ghost** (Purple): Activates **Zen Mode for 10.0 seconds** — walls do not stop the game (wrap-around), obstacles are phased through, and body overlap is permitted.
- ⚡ **3x Multiplier** (Gold): Triples points earned for **5.0 seconds**.
- ✂ **Shrink** (Orange): Cuts snake body length in half upon eating.
- ☠ **Poison (Hazard - Avoid!)** (Crimson): Deducts -30 points, adds +1 unrewarded length, and triggers an uncontrollable **Toxic Speed Surge** (2x speed) for **5.0 seconds**!
 
For comprehensive details on all items and scoring calculations, read [MODES_AND_RULES.md](MODES_AND_RULES.md).

---

## 👔 The Boss Key

Need to disguise your coffee break? Press `b` or `esc` at any moment:
- Instantly replaces the entire screen with an authentic, animated `htop` dashboard, compiler build stream, or `git log`.
- Automatically pauses the game and preserves exact game state.
- Press `b` again (or any key) to seamlessly resume play.

Configure your preferred screen in `~/.config/snake/config.toml`:
```toml
default_boss_screen = "htop" # "htop", "build", or "git"
```

---

## 🎨 Themes & Accessibility

Switch themes anytime with `-t, --theme <name>` or interactively in the menu:
- `dark`: Premier modern neon palette (Green/Cyan/Purple).
- `light`: Clean, high-legibility theme for bright terminal backgrounds.
- `high-contrast`: Pure black & white accessibility theme.
- `colorblind`: Okabe-Ito verified colorblind-safe palette (Blue/Vermilion/Yellow).
- `corporate`: Muted grayscale aesthetics that blend naturally into enterprise environments.

For screen-reader accessibility, launch with `--screen-reader` for plain-text terminal output.

---

## 🏆 Statistics, Leaderboards & Achievements

All player stats, high scores by mode, and unlocked achievements are locally persisted to `~/.config/snake/stats.json`:

```bash
# View interactive summary
snake stats

# Export leaderboard to CSV
snake stats --csv > leaderboard.csv

# Output full JSON data
snake stats --json
```

### Unlocked Achievements Include:
- **First Bite**: Eat your first food apple.
- **Centurion**: Score over 100 points in a single run.
- **Grandmaster**: Score over 300 points in a single run.
- **Long Boi**: Reach a snake length of 25 segments.
- **Coffee Break**: Accumulate over 5 minutes of total playtime.

---

## 🤖 Autonomous AI Bot & Headless Benchmarking

The built-in AI bot uses an intelligent **Breadth-First Search (BFS)** pathfinder paired with a **Hamiltonian tail-following survival fallback**:
- Prioritizes shortest safe paths toward food apples and positive power-ups.
- Proactively maneuvers around negative poison hazard apples.
- Automatically calculates toroidal wall wrap-around paths when Ghost power-up or Zen mode is active.

```bash
# Watch the bot play live in the terminal
snake bot

# Run a headless benchmark (>4,000 ticks/sec)
snake bot --headless --ticks 1000

# Benchmark a specific seed
snake bot --headless --ticks 5000 --seed 12345
```

---

## 📼 Deterministic Replays & Verification

Every completed game is automatically recorded as a lightweight, compressed `.snk` file stored in `~/.config/snake/replays/`:
- **Deterministic Playback**: Replays record the PRNG seed and input timestamps, reproducing identical games tick-by-tick.
- **Speed Control**: Accelerate or decelerate playback with `--speed <ms>`.
- **Headless Verification**: Ensure game integrity without rendering:
  ```bash
  snake replay my_run.snk --verify
  ```

---

## ⚙️ Configuration File

Settings and keybindings can be customized in `~/.config/snake/config.toml`:

```toml
theme = "dark"
grid_width = 24
grid_height = 20
base_speed_ms = 110
adaptive_speed = true
screen_reader = false
default_boss_screen = "htop"

[keys]
up = ["up", "w", "k"]
down = ["down", "s", "j"]
left = ["left", "a", "h"]
right = ["right", "d", "l"]
pause = ["p", " "]
menu = ["m", "tab"]
boss = ["b", "esc"]
restart = ["r"]
quit = ["q", "ctrl+c"]
```

Manage your configuration easily with:
```bash
snake config         # View current configuration
snake config --edit  # Edit in $EDITOR
snake config --reset # Restore factory defaults
```

---

## 🛠️ IoT Hardware Controller Mode

Play Snake using a physical ESP8266 or ESP32 joystick over WiFi UDP!

```
[ Analog Joystick ] ──> [ ESP8266 NodeMCU ] ──( WiFi / UDP )──> [ snake CLI ]
```

1. **Wiring**: Connect analog joystick VRx to `A0`, button SW to `D5`, VCC to `3.3V`, GND to `GND`.
2. **Firmware**: Flash [`hardware/esp8266_snake.ino`](hardware/esp8266_snake.ino) using Arduino IDE or PlatformIO.
3. **Launch**:
   ```bash
   snake play --hardware-udp :4242
   ```
See [`hardware/PROTOCOL.md`](hardware/PROTOCOL.md) for complete protocol specification and pinout details.

---

## 🏗️ Architecture

The codebase adheres to clean Go architecture with strict separation of concerns:

```
┌────────────────────────────────────────────────────────┐
│                   CLI Wiring (Cobra)                   │
└───────┬─────────────┬─────────────┬────────────┬───────┘
        │             │             │            │
        ▼             ▼             ▼            ▼
┌──────────────┐┌───────────┐┌───────────┐┌──────────────┐
│  Engine      ││ Renderer  ││ Input     ││ Storage      │
│  (Pure logic,││ (BubbleTea││ (Keyboard,││ (TOML config,│
│  no I/O,     ││  Lipgloss,││  AI bot,  ││  stats JSON, │
│  seeded PRNG)││  Themes)  ││  UDP HW)  ││  replays)    │
└──────────────┘└───────────┘└───────────┘└──────────────┘
```

- **Engine** (`internal/engine`): Deterministic state machine. Has zero awareness of terminals, screens, or keypresses.
- **Renderer** (`internal/renderer`): Bubble Tea UI model. Handles terminal layout, double-buffering, and themes.
- **Input** (`internal/input`): Pluggable input interface. Dispatches commands from keyboard, BFS AI bot, or external UDP hardware.
- **Storage** (`internal/storage`): Encapsulates config, save states, stats, and replay codecs.

---

## 🧪 Testing & Verification

The project includes an extreme-level unit, stress, and race-detection test suite:

```bash
# Run all unit and extreme tests with race detector and code coverage
make test

# Run Go static analysis
make lint

# Execute deterministic replay verification
./bin/snake replay --verify

# Run autonomous AI bot in headless benchmark mode (>4,000 ticks/sec)
./bin/snake bot --headless --ticks 1000
```

---

## 👤 Author

**Md. Inzamamul Lohani** ([@hamimLohani](https://github.com/hamimLohani))

---

## 📄 License

MIT License © 2026 Md. Inzamamul Lohani. Crafted with clean engineering standards for developers everywhere.
