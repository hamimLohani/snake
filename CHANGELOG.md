# Changelog

All notable changes to `snake` will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-10-14

### Added
- **Core Engine**: Pure, 100% deterministic Snake game logic with tick-based state transitions and seeded PRNG.
- **Game Modes**:
  - `classic`: Standard authentic Snake gameplay with wall and self-collisions.
  - `time-attack`: 60-second countdown with bonus time awarded for apples consumed.
  - `survival`: Dynamic boundary shrinking every 50 ticks to test agility.
  - `zen`: Boundary wrapping and pass-through ghost mechanics for zero stress.
- **Power-Ups**:
  - `Slow-Mo`: Halves tick speed for 60 ticks.
  - `Ghost`: Enables phasing through body segments for 50 ticks.
  - `Score Multiplier`: 3x points bonus for 70 ticks.
  - `Shrink`: Halves snake length on consumption.
- **Boss Key Concealment**: One keypress (`b` or `esc`) switches to authentic fake screens:
  - `htop`: Multi-core CPU bars, memory usage, load averages, and realistic process table.
  - `build`: Compiler test execution and build log stream.
  - `git`: Formatted commit history with graph decoration.
- **Themes**:
  - `dark`: High-vibrancy modern cyber aesthetic.
  - `light`: Clean, crisp contrast for light terminals.
  - `high-contrast`: Pure monochrome max-contrast rendering.
  - `colorblind`: Okabe-Ito accessible palette.
  - `corporate`: Low-saturation subdued styling.
- **Persistent Storage**:
  - Config: `~/.config/snake/config.toml` with rebindable keybindings.
  - Stats & Leaderboard: `~/.config/snake/stats.json` with CSV/JSON exports.
  - Auto-Save & Resume: Automatically persists in-progress game on exit/blur and resumes seamlessly.
- **Daily Challenge**:
  - `snake daily`: Deterministic global seed derived from UTC date (`YYYY-MM-DD`).
- **Autonomous AI Bot**:
  - `snake bot`: Heuristic BFS pathfinding with tail-escape survival logic.
  - Headless benchmark mode supporting >7,000 ticks/second.
- **Replay System**:
  - Gzip-compressed `.snk` recording and verification.
  - Headless determinism verification (`snake replay --verify`).
- **IoT Hardware Controller Mode**:
  - UDP and Serial listener for external hardware inputs.
  - Arduino sketch for ESP8266/ESP32 analog joystick and button.
  - Binary and ASCII protocol specification (`hardware/PROTOCOL.md`).
- **Distribution & Quality**:
  - Homebrew formula (`Formula/snake.rb`).
  - GitHub Actions CI matrix builds for macOS (ARM64/AMD64) and Linux (ARM64/AMD64).
  - Makefile with targets for build, test, lint, install, and clean.
