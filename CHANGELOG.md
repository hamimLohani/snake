# Changelog

All notable changes to `snake` will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.8] - 2026-10-15

### Added
- **10s Zen Mode for Ghost Power-Up (`👻`)**: Eating a Ghost power-up now activates full **Zen Mode for 10.0 seconds**:
  - **Wall Wrapping**: Running into walls no longer stops the game; the snake smoothly wraps around the arena boundaries.
  - **Body Overlap Phasing**: The snake passes harmlessly through its own body segments without self-collision.
  - **Obstacle Phasing**: Safe passage directly through procedural obstacles.
  - **Live HUD Badge**: Purple `👻 ZEN MODE 10.0s` live countdown displayed in the HUD.
  - **AI Autopilot Awareness**: Autonomous BFS bot recognizes walls and body overlaps as safe during the 10-second Zen window.

## [1.0.7] - 2026-10-15

### Added
- **Single Food (One After Another)**: Exactly one permanent food apple appears on the board at a time and respawns immediately upon consumption.
- **Concurrent Featured Specials (2–3 Items Simultaneously)**:
  - 1 Poison hazard (`☠`) is continuously present alongside the power-up(s) to actively dodge.
  - 1–2 positive power-ups (`⏱` Slow-Mo, `👻` Ghost, `⚡` Multiplier, `✂` Shrink) appear simultaneously with the poison and the food (2 specials on small boards, 3 specials on medium/large boards).
  - Special featured items despawn and rotate every 5.0 seconds if uneaten.
- **Updated Sample Board Preview**: Menu preview updated to reflect the 1 food + 2-3 concurrent specials model with poison hazard indicator.

## [1.0.6] - 2026-10-15

### Added
- **5-Second Calibrated Durations**: All power-up active durations (`Slow-Mo`, `Ghost Mode`, `3× Multiplier`), board despawn lifespans, and toxic speed surges are calibrated to exactly **5.0 seconds** via `TicksForDuration(5 * time.Second)`.
- **Multiple Simultaneous Foods**: The arena now continuously maintains multiple permanent food apples (default 3) on the board simultaneously. Consuming any apple immediately replenishes another so food is always abundant.
- **Concurrent Extra Items & Guaranteed Poison Hazard**: A positive power-up and a negative poison hazard apple (`☠`) coexist concurrently alongside multiple permanent food apples at all times (5 total items active).
- **HUD Live Countdown**: Replaced raw tick count with a real-time countdown in seconds (`Slow-Mo 5.0s`, `⚠️ POISON 5.0s`, `Ghost 4.8s`).
- **AI Autopilot Multiple-Food Pathfinding**: BFS bot solver dynamically scans all active food apples and power-ups to pursue the shortest safe path while actively steering around all negative poison hazard cells.
- **Save State Multi-Item Serialization**: `savegame.go` now serializes and restores the full slice of multiple concurrent food and hazard items.

## [1.0.5] - 2026-10-15

### Added
- **Permanent Food Spawning (One After Another)**: Standard food apples are permanent on the board until eaten and immediately respawn upon consumption.
- **Concurrent Extra Power Items**: Power-ups and hazard items spawn concurrently alongside permanent food rather than replacing it, remaining on the board for up to 80 ticks before despawning if uneaten.
- **One-Third Reduced Power-Up Durations**: Rebalanced positive power-up active durations by one-third for tighter, faster-paced gameplay:
  - `Slow-Mo`: Reduced from 60 to **40 ticks**.
  - `Ghost Mode`: Reduced from 50 to **34 ticks**.
  - `3× Multiplier`: Reduced from 70 to **46 ticks**.
- **Negative Impact Hazard Food (`ItemPoison` / `☠`)**: Introduced avoidable negative food that players must steer around:
  - **Score Penalty**: Deducts **-30 points** upon ingestion (floored at 0).
  - **Toxic Speed Surge**: Triggers an active negative status effect accelerating tick rate by **35%** for 25 ticks, drastically increasing crash hazards.
  - **Unrewarded Body Growth**: Increases snake length by +1 segment without awarding points.
  - **Theme Glyphs**: Rendered as `☠ ` (Dark), `X ` (Light), `! ` (High-Contrast), `✖ ` (Colorblind-Safe), and `! ` (Corporate).
- **AI Bot Hazard Avoidance**: Updated autonomous BFS autopilot solver to identify poison hazard food as an unsafe cell and navigate around it.
- **UI & HUD Enhancements**: Added crimson `⚠️ POISON <ticks>` badge to active HUD status display and included the Hazard glyph (`Hazard:☠`) in the live Sample Board preview legend.
- **State Persistence**: Savegame snapshot now serializes and restores concurrent bonus point, type, active status, and remaining board lifespan.

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
