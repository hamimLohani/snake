# 🐍 snake — Developer-Grade Terminal Snake CLI

[![Release](https://img.shields.io/github/v/release/hamimLohani/snake?color=brightgreen)](https://github.com/hamimLohani/snake/releases)
[![CI](https://github.com/hamimLohani/snake/actions/workflows/ci.yml/badge.svg)](https://github.com/hamimLohani/snake/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/hamimLohani/snake)](https://goreportcard.com/report/github.com/hamimLohani/snake)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> **A coffee-break game for developers.** Installs in one command, launches instantly, looks polished in any terminal, and built to business-grade engineering standards.

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
- **Terminal Safety**: Full alternate-screen buffer isolation; cleans up raw mode and restores terminal state on exit, SIGINT, or panic.
- **Strictly Deterministic Engine**: Pure state machine with seeded PRNG; enables 100% reproducible replays, headless simulation, and automated verification.
- **Interactive Main Menu & Live Sample Board**: Sleek two-column launch screen with dynamic live board preview for board sizes (Small, Medium, Classic, Large), visual themes, obstacles, and difficulty.
- **4 Distinct Game Modes**: Classic, Time Attack, Survival (shrinking borders), and Zen (infinite wrap-around).
- **Dynamic Power-Ups & Avoidable Hazards**: Permanent food with concurrent bonus items: Slow-Mo, Ghost, 3x Multiplier, Shrink, and avoidable Poison hazard apples.
- **Discreet Boss Key**: Press `b` or `esc` to instantly swap the screen to a convincing fake `htop`, compiler build log, or `git log`.
- **Procedural Obstacles**: Deterministic obstacle mazes generated from any shareable `--seed`.
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

### Option 2: Direct Binary Download (Linux)

#### Linux x86_64 (amd64):
```bash
# Download binary using curl
curl -sL https://github.com/hamimLohani/snake/releases/latest/download/snake-linux-amd64 -o snake
chmod +x snake
sudo mv snake /usr/local/bin/

# Or using wget:
wget -qO snake https://github.com/hamimLohani/snake/releases/latest/download/snake-linux-amd64
chmod +x snake
sudo mv snake /usr/local/bin/
```

#### Linux ARM64 (aarch64 — Raspberry Pi, AWS Graviton, etc.):
```bash
curl -sL https://github.com/hamimLohani/snake/releases/latest/download/snake-linux-arm64 -o snake
chmod +x snake
sudo mv snake /usr/local/bin/
```

#### Linux Tarball Archive (.tar.gz):
```bash
curl -sL https://github.com/hamimLohani/snake/releases/latest/download/snake-v1.0.5-linux-amd64.tar.gz | tar -xz
sudo mv snake-linux-amd64 /usr/local/bin/snake
```

### Option 3: Homebrew (macOS & Linux)
```bash
brew tap hamimlohani/tap
brew install snake
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

## 🎮 CLI Interface

```bash
# Launch interactive Main Menu & Settings
snake
snake play

# Start playing directly with custom board size and difficulty
snake play --direct
snake play --size large --difficulty hard
snake play -s small -d easy
snake play --size 30x15 --mode survival

# Choose a game mode and visual theme
snake play --mode time-attack --theme dark
snake play --mode survival --theme colorblind
snake play --mode zen --theme light

# Play with a deterministic shared seed and obstacles
snake play --seed 42109 --obstacles

# Daily challenge (unified global seed for everyone worldwide today)
snake daily

# View local leaderboard, player metrics, and unlocked achievements
snake stats
snake stats --json
snake stats --csv

# Replay previous games or verify determinism
snake replay                      # Plays latest game
snake replay my_run.snk --speed 40 # 2x playback speed
snake replay my_run.snk --verify  # Headless determinism check

# Autonomous AI bot demo
snake bot                         # Watch bot play interactively
snake bot --headless --ticks 1000 # Headless stress test / benchmark

# Manage configuration
snake config                      # View current config & keybindings
snake config --edit               # Open ~/.config/snake/config.toml in $EDITOR
snake config --reset              # Reset to defaults
```

---

## ⌨️ Controls & Keybindings

All keys are fully rebindable in `~/.config/snake/config.toml`:

| Action | Default Keys |
|--------|--------------|
| **Move Up** | `↑` (Arrow Up), `W`, `K` (Vim) |
| **Move Down** | `↓` (Arrow Down), `S`, `J` (Vim) |
| **Move Left** | `←` (Arrow Left), `A`, `H` (Vim) |
| **Move Right** | `→` (Arrow Right), `D`, `L` (Vim) |
| **Pause / Resume** | `P`, `Space` |
| **Boss Key** | `B`, `Esc` |
| **Restart Game** | `R` (when game is over) |
| **Quit Game** | `Q`, `Ctrl+C` |

---

## 🕹️ Game Modes

> 📖 For detailed mechanics, speed curves, scoring formulas, and collision physics, see the complete [**Modes, Difficulties & Rulebook (`MODES_AND_RULES.md`)**](MODES_AND_RULES.md).

| Mode | Flag | Description |
|------|------|-------------|
| **Classic** | `--mode classic` | Standard snake mechanics; eat food, grow, avoid walls and your own tail. |
| **Time Attack** | `--mode time-attack` | 60-second starting timer. Eating food adds bonus time (+3s). Score as high as possible before the clock expires! |
| **Survival** | `--mode survival` | The walls contract inward every 50 ticks! Maneuver in an ever-shrinking arena. |
| **Zen** | `--mode zen` | No deaths, no game-overs. Snake wraps around edges and passes harmlessly through itself. Perfect for decompression. |

---

## ⚡ Power-Ups & Special Items
 
One food apple exists at a time and immediately respawns when eaten ("one after another"). Alongside the food, **2–3 featured special items appear simultaneously** — 1 guaranteed poison hazard + 1–2 positive power-ups (5.0 second lifespan each, then rotate to a new location):
 
- ⏱ **Slow-Mo** (Cyan): Halves tick speed for **5.0 seconds**.
- 👻 **Ghost** (Purple): Enables passing through your own body segments without crashing for **5.0 seconds**.
- ⚡ **3x Multiplier** (Gold): Triples points earned for **5.0 seconds**.
- ✂ **Shrink** (Orange): Cuts snake body length in half upon eating.
- ☠ **Poison (Hazard - Avoid!)** (Crimson): Deducts -30 points, adds +1 unrewarded length, and triggers an uncontrollable **Toxic Speed Surge** (2x speed) for **5.0 seconds**!
 
For comprehensive details on all items and scoring calculations, read [MODES_AND_RULES.md](file:///Users/Inz_mac/Developer/Shell/Project/snake/MODES_AND_RULES.md).

---

## 👔 The Boss Key

Need to disguise your coffee break? Press `b` or `esc` at any moment:
- Instantly replaces the entire screen with an authentic, animated `htop` dashboard, compiler build stream, or `git log`.
- Automatically pauses the game and preserves exact state.
- Press `b` again (or any key) to seamlessly resume play.

Configure your preferred screen in `~/.config/snake/config.toml`:
```toml
default_boss_screen = "htop" # "htop", "build", or "git"
```

---

## 🎨 Themes

Switch themes with `--theme <name>`:
- `dark`: Premier modern neon palette (Green/Cyan/Purple).
- `light`: Clean, high-legibility theme for bright terminal backgrounds.
- `high-contrast`: Pure black & white accessibility theme.
- `colorblind`: Okabe-Ito verified colorblind-safe palette (Blue/Vermilion/Yellow).
- `corporate`: Muted grayscale aesthetics that blend naturally into enterprise environments.

---

## 🛠️ IoT Hardware Controller Mode

Play Snake using a physical ESP8266 or ESP32 joystick!

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

The codebase adheres to strict separation of concerns:

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

```bash
# Run unit and race-condition test suite
make test

# Run static analysis
make lint

# Execute deterministic replay verification
./bin/snake bot --headless --ticks 500
```

---

## 👤 Author

**Md. Inzamamul Lohani** ([@hamimLohani](https://github.com/hamimLohani))

---

## 📄 License

MIT License © 2026 Md. Inzamamul Lohani. Crafted with clean engineering standards for developers everywhere.
