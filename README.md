# 🐍 snake — Developer-Grade Terminal Snake CLI

[![CI](https://github.com/developer/snake/actions/workflows/ci.yml/badge.svg)](https://github.com/developer/snake/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/developer/snake)](https://goreportcard.com/report/github.com/developer/snake)
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
- **4 Distinct Game Modes**: Classic, Time Attack, Survival (shrinking borders), and Zen (infinite wrap-around).
- **Dynamic Power-Ups**: Slow-Mo, Ghost (pass through own body), 3x Score Multiplier, and Shrink.
- **Discreet Boss Key**: Press `b` or `esc` to instantly swap the screen to a convincing fake `htop`, compiler build log, or `git log`.
- **Procedural Obstacles**: Deterministic obstacle mazes generated from any shareable `--seed`.
- **Local Persistence**: Leaderboards, high scores, play time, achievements (`~/.config/snake/stats.json`), and rebindable TOML config (`~/.config/snake/config.toml`).
- **Autonomous AI Bot**: Heuristic BFS autopilot with tail-following survival logic (>7,000 ticks/sec in headless mode).
- **IoT Hardware Mode**: Control the snake over WiFi UDP or USB Serial using an ESP8266 or ESP32 joystick.

---

## 🚀 Installation

### Option 1: Homebrew (macOS & Linux)
```bash
brew tap developer/snake
brew install snake
```

### Option 2: Go Native Package Manager
```bash
go install snake/cmd/snake@latest
```

### Option 3: From Source
```bash
git clone https://github.com/developer/snake.git
cd snake
make build
# Binary created at ./bin/snake
```

---

## 🎮 CLI Interface

```bash
# Play standard classic mode
snake play

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

| Mode | Flag | Description |
|------|------|-------------|
| **Classic** | `--mode classic` | Standard snake mechanics; eat food, grow, avoid walls and your own tail. |
| **Time Attack** | `--mode time-attack` | 60-second starting timer. Eating food adds bonus time (+3s). Score as high as possible before the clock expires! |
| **Survival** | `--mode survival` | The walls contract inward every 50 ticks! Maneuver in an ever-shrinking arena. |
| **Zen** | `--mode zen` | No deaths, no game-overs. Snake wraps around edges and passes harmlessly through itself. Perfect for decompression. |

---

## ⚡ Power-Ups

During gameplay, special items appear on the grid:

- ⏱ **Slow-Mo** (Cyan): Halves the snake's tick speed for 60 ticks.
- 👻 **Ghost** (Purple): Enables passing through your own body segments without crashing for 50 ticks.
- ⚡ **3x Multiplier** (Gold): Triples points earned for 70 ticks.
- ✂ **Shrink** (Orange): Cuts snake body length in half upon eating.

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

## 📄 License

MIT License. Crafted with clean engineering standards for developers everywhere.
