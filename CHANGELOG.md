# Changelog

All notable changes to `snake` will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.4] - 2026-10-15

### Added
- **Live Sample Board Preview in Main Menu**: The interactive menu now features a dynamic side-by-side **Sample Board Preview** that live-updates as settings are cycled:
  - **Dynamic Sizing & Proportions**: The sample arena visually resizes and updates its dimensions badge across `SMALL (20×12)`, `MEDIUM (28×16)`, `CLASSIC (24×20)`, and `LARGE (36×20)` presets.
  - **Theme Visualization**: Real-time rendering of the active theme's borders (rounded, normal, thick), color palettes, snake head/body glyphs, food, and bonus items (Dark, Light, High-Contrast, Colorblind-Safe, Corporate).
  - **Obstacle Toggle Feedback**: Procedural obstacle blocks (`██`) appear in the sample arena when Obstacles are toggled `ON` and disappear when `OFF`.
  - **Glyph Legend**: Instant reference showing exact Head, Body, Food, and Obstacle symbols.
  - **Responsive Geometry**: Two-column layout fits within standard 80×24 terminal boundaries without scrolling, with automatic graceful fallback for narrow terminals.

## [1.0.3] - 2026-10-15

### Added
- **Board Size Selection**: Choose from 4 curated board size presets: `SMALL (20×12)` for compact split screens, `MEDIUM (28×16)` for balanced play, `CLASSIC (24×20)` retro standard, and `LARGE (36×20)` spacious arena. Cycle interactively in the main/pause menu or pass CLI flags via `--size` / `-s` (supporting named presets or custom `WxH` like `30x15`).
- **Interactive Main Menu on Launch**: Running `snake` or `snake play` now opens directly into a sleek, full-screen **Main Menu & Settings** screen before gameplay starts.
- **Instant Play / Skip Menu (`--direct`)**: Added `--direct` flag to bypass the main menu and begin playing immediately. Passing explicit configuration flags (e.g. `--mode`, `-d`, `-s`, `--speed`, `--obstacles`, `--seed`) also bypasses the menu for immediate gameplay.
- **Smart Saved-Game Resume**: When an unfinished game exists from a previous session or Boss Key exit, the menu prominently offers `[ ▶ RESUME SAVED GAME ]` with score preview, as well as `[ Start New Game ]`.
- **Keyboard-First Menu UX**: Press `Enter` or `Space` anywhere in the start menu to start playing immediately. Use `Up`/`Down` (`W`/`S`) to navigate, `Left`/`Right` (`A`/`D`) to cycle options, and `Q` to quit.

## [1.0.2] - 2026-10-15

### Added
- **In-Game Settings & Options Menu**: Press `M` or `Tab` at any time (or from Pause / Game Over screens) to interactively change **Game Mode**, **Difficulty**, **Procedural Obstacles**, and **Visual Theme** with live preview and instant restart.

## [1.0.1] - 2026-10-14

### Added
- **Difficulty Presets**: Added `--difficulty` (`-d`) flag supporting `easy` (160ms/tick), `normal` (110ms/tick), `hard` (75ms/tick), and `insane` (45ms/tick with procedural obstacles).
- **Manual Pages**: Added troff/groff manual pages in Section 6 (`snake.6`) and Section 1 (`snake.1`), along with built-in `snake man` CLI pager.

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
