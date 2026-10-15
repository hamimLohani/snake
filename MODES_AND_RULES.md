# 🐍 Snake CLI — Game Modes, Difficulties & Rulebook

This document provides a comprehensive guide to all game modes, difficulty curves, power-ups, scoring formulas, and collision mechanics in **Snake CLI**.

---

## 📑 Table of Contents
1. [Game Modes](#-game-modes)
   - [Classic Mode](#1-classic-mode)
   - [Time Attack Mode](#2-time-attack-mode)
   - [Survival Mode](#3-survival-mode)
   - [Zen Mode](#4-zen-mode)
   - [Daily Challenge](#5-daily-challenge)
2. [Difficulty Presets & Speed Curves](#-difficulty-presets--speed-curves)
   - [Difficulty Presets](#difficulty-presets)
   - [Adaptive Speed Scaling](#adaptive-speed-scaling)
3. [Power-Ups & Special Items](#-power-ups--special-items)
4. [Obstacles & Procedural Mazes](#-obstacles--procedural-mazes)
5. [Physics & Collision Mechanics](#-physics--collision-mechanics)
6. [Scoring & Multiplier Calculations](#-scoring--multiplier-calculations)
7. [In-Game Controls & Shortcuts](#-in-game-controls--shortcuts)

---

## 🕹️ Game Modes

| Mode | Command | Objective | Death Conditions | Best For |
|------|---------|-----------|------------------|----------|
| **Classic** | `snake play -m classic` | Grow as long as possible | Wall hit, Self collision, Obstacle hit | Traditional Snake enthusiasts |
| **Time Attack** | `snake play -m time-attack` | Maximize score within 60s (+3s per food) | Timer reaches 0, Wall/Tail collision | High-intensity sprint play |
| **Survival** | `snake play -m survival` | Outlast contracting perimeter walls | Wall hit, Tail collision, Obstacle hit | Reflex & claustrophobic agility |
| **Zen** | `snake play -m zen` | Relaxing, infinite gameplay | **None** (cannot die) | Casual coffee breaks & unwinding |
| **Daily** | `snake daily` | Compete on worldwide unified seed | Same as Classic | Global leaderboards & fair competition |

---

### 1. Classic Mode
> **Flag:** `--mode classic` (Default)

The quintessential arcade Snake experience built to competitive deterministic standards.

#### Rules & Mechanics
- **Starting Grid:** 24 × 20 cells (configurable in `config.toml`).
- **Initial Length:** 3 body segments moving Right at the center of the board.
- **Growth:** Consuming food adds +1 body segment. Growth is delayed by one tick, meaning the tail does not pop off on the move when food is consumed.
- **Game Over Conditions:**
  - Head moves outside the grid boundaries (`X < 0 || X >= Width || Y < 0 || Y >= Height`).
  - Head collides with any segment of the snake’s own body (excluding the vacating tail tip unless growth is pending).
  - Head collides with a procedural obstacle block (`██`).

---

### 2. Time Attack Mode
> **Flag:** `--mode time-attack`

A high-pressure race against a relentless countdown clock.

#### Rules & Mechanics
- **Starting Timer:** 60 seconds.
- **Time Replenishment:** Every apple/food item consumed awards **+3 seconds** to the clock.
- **Base Points:** Base point value is doubled to **25 points per food** (vs. 10 in Classic).
- **Game Over Conditions:**
  - Clock reaches `00s` ("Time is up!").
  - Crashing into walls, tail, or obstacles.
- **Strategy:** Aggressive routing is rewarded. Sacrificing neat tail packing for rapid food intercept is necessary to keep the clock from expiring.

---

### 3. Survival Mode
> **Flag:** `--mode survival`

The arena itself is your enemy. The walls actively close in on you over time.

#### Rules & Mechanics
- **Wall Contraction:** Every **50 ticks** (~5.5 seconds at Normal speed), the arena perimeter shrinks by **1 unit on all four sides**:
  - `MinX += 1`, `MaxX -= 1`
  - `MinY += 1`, `MaxY -= 1`
- **Floor Limit:** Contraction halts once the arena reaches a tight **6 × 6** box.
- **Auditory/Visual Warning:** An in-game alert banner flashes: *"Warning! The walls are closing in!"*
- **Game Over Conditions:**
  - Hitting any of the contracted perimeter borders.
  - Biting own tail in the increasingly constrained space.
- **Strategy:** Avoid looping near the exterior edges as 50-tick milestones approach. Keep the snake centralized in the inner 10×10 safe radius.

---

### 4. Zen Mode
> **Flag:** `--mode zen`

A meditative, stress-free mode designed for taking a mental break without frustration.

#### Rules & Mechanics
- **Infinite Boundary Wrapping:** Moving past the right border warps the snake instantly to the left border; moving past the top warps to the bottom.
- **Ghost Pass-Through:** The snake passes harmlessly through its own body segments without dying.
- **No Game Over:** It is mathematically impossible to trigger a Game Over screen in Zen mode.
- **Exit Strategy:** Quit or pause whenever you are ready using `Q` or `P`.

---

### 5. Daily Challenge
> **Flag:** `snake daily`

A unified global competition updated every 24 hours.

#### Rules & Mechanics
- **Deterministic Global Seed:** Calculated using an FNV-64a hash of the current UTC date string (`UTC YYYY-MM-DD`):
  $$\text{Seed} = \text{Hash}(\text{"YYYY-MM-DD"})$$
- **Fair Play:** Every player running `snake daily` anywhere in the world receives:
  - The exact same starting snake spawn.
  - The exact same procedural obstacle maze layout.
  - The exact same coordinates and sequence of food and power-up spawns.
- **Leaderboard Integration:** Daily high scores are cataloged in `~/.config/snake/stats.json`.

---

## ⚡ Difficulty Presets & Speed Curves

Difficulty in Snake is determined by the **tick duration** (the time in milliseconds between each discrete step of the game engine).

### Difficulty Presets

You can select a difficulty via the in-game menu (`M`), via the CLI flag (`-d, --difficulty`), or via custom millisecond values (`--speed`):

| Preset | Flag | Tick Interval | Speed (Steps/Sec) | Characteristics |
|--------|------|---------------|-------------------|-----------------|
| **Easy** | `-d easy` | **160 ms** | 6.2 steps/sec | Forgiving pace; ideal for beginners or relaxed sessions. |
| **Normal** | `-d normal` | **110 ms** | 9.1 steps/sec | Default balanced speed; standard tournament pace. |
| **Hard** | `-d hard` | **75 ms** | 13.3 steps/sec | Fast-paced; demands tight reflexes and early turn queuing. |
| **Insane** | `-d insane` | **45 ms** | 22.2 steps/sec | Hyper-speed with procedural obstacles automatically enabled. |

### Adaptive Speed Scaling

When `adaptive_speed = true` (default in `~/.config/snake/config.toml`), the game dynamically accelerates as your score increases:

$$\text{Effective Speed} = \max\left(40\text{ ms},\, \text{BaseSpeed} - \left\lfloor\frac{\text{Score}}{50}\right\rfloor \times 4\text{ ms}\right)$$

- For every **50 points** scored, tick duration decreases by **4 milliseconds**.
- Scaling caps at a **40 ms floor** (25 steps per second) to preserve human playability.
- Consuming a **Slow-Mo power-up** overrides this calculation by doubling tick interval (+100% time per step).

---

## 🎁 Power-Ups & Special Items

Snake CLI features a **1 food + multiple concurrent special items** model:
1. **Single Food (One After Another):** One food apple exists at a time. The moment it is eaten, another immediately respawns elsewhere — always present, always rotating.
2. **2–3 Featured Special Items Simultaneously:** Alongside the food, **2–3 special items coexist on the board at the same time**: 1 guaranteed poison hazard + 1–2 positive power-ups (small boards get 2, medium/large boards get 3). Each special item lasts 5.0 seconds on the board before despawning if uneaten — then a fresh one appears in a new location.
3. **Calibrated 5-Second Durations:** All active power-up effects, item despawn lifespans, and hazard surges last **5.0 seconds**, with an exact live second countdown in the HUD (`Slow-Mo 5.0s`, `⚠️ POISON 5.0s`).
4. **Negative Impact Hazard Food (Avoid!):** A poison hazard apple (`☠`) is **always present concurrently** with the power-up(s), requiring players to actively steer around it at all times.

| Item | Glyph | Palette Color | Duration | Gameplay Effect |
|------|-------|---------------|----------|-----------------|
| **Food (Single)** | `★` / `🍎` | Red / Neon Green | Permanent | Awards standard points (+10 pts, or +25 pts in Time Attack) and grows snake length by 1. One apple at a time — immediately respawns upon being eaten ("one after another"). |
| **Slow-Mo** | `⏱` | Cyan (`#00E5FF`) | **5.0 seconds** | **Halves engine speed** (doubles tick duration). Awards +20 bonus pts. |
| **Ghost Mode** | `👻` | Purple (`#D500F9`) | **10.0 seconds** | **Temporary Zen Mode**: For the next 10 seconds, walls do not stop the game (toroidal wrap-around), obstacles are safely phased through, and snake head passes harmlessly through its own body (full overlap allowed). Awards +20 bonus pts. |
| **3× Multiplier** | `⚡` | Gold (`#FFD600`) | **5.0 seconds** | **3× Score Multiplier**: Triples score for all food consumed while active. Awards +30 bonus pts. |
| **Shrink** | `✂` | Orange (`#FF6D00`) | Instant | **Halves snake body length** (minimum 3 segments). Awards +15 bonus pts. |
| **Poison (Hazard)** | `☠` / `X` / `!` | Crimson (`#FF0055`) | **5.0 seconds** surge | **HAZARD - MUST AVOID!** Incur **-30 pts penalty** (min 0), **+1 unrewarded length**, and triggers a **Toxic Speed Surge** (doubles speed: `TickRate / 2`). |

> [!WARNING]
> **Avoid the Poison (`☠`)!** Eating a poison apple inflicts an immediate -30 point penalty, grows your body without score reward, and causes an uncontrollable 5.0-second speed surge that makes tight cornering extremely dangerous.

> [!NOTE]
> Power-up countdown is measured in real gameplay seconds calibrated to engine tick rates (`TicksForDuration(5 * time.Second)`). If the game is paused or you enter the Boss screen, remaining power-up time is preserved.

---

## 🧱 Obstacles & Procedural Mazes

When enabled via the in-game menu (`M`), `--obstacles`, or in `insane` difficulty, procedural obstacles are placed on the board:

- **Safety Zone:** The engine enforces a 7×5 cell clearance box around the snake's initial starting coordinate (`centerX`, `centerY`). You will never spawn facing an immediate wall.
- **Symmetry Patterns:**
  - **Pillars:** Symmetrical 2×2 blocks situated near the 4 corners.
  - **Barriers:** Horizontal split barriers with central passage slots.
  - **Cross Corners:** L-shaped defensive barriers along quadrant perimeters.
- **RNG Determinism:** Obstacle placement is 100% deterministic based on the provided `--seed`.

---

## 📐 Physics & Collision Mechanics

### Coordinate Grid
- The board uses a 2D Cartesian grid with `(0, 0)` at the **top-left corner**.
- `X` increases towards the Right (`+1`).
- `Y` increases downwards (`+1`).

### Direction Buffer (Anti-180° Suicide Protection)
In arcade Snake, tapping two arrow keys in rapid succession (e.g., `Right` then `Down` while moving `Left`) can cause the snake to turn 180° into its own neck within a single frame.

**The Solution in Snake CLI:**
- The engine uses a 2-event FIFO **Direction Queue** (`DirQueue []Direction`).
- Rapid key inputs are buffered and validated against the *planned* next direction, not just the currently rendered position.
- Instant 180° opposite reversals (`DirRight` -> `DirLeft`) are rejected immediately, protecting your run.

### Tail Vacating Rule
When the snake moves:
1. Head advances to `NewHead = Head + Delta`.
2. `NewHead` is checked against all body segments *excluding the tail tip*.
3. Because the tail vacates the cell at the exact moment the head moves, **moving into your current tail's position is completely safe** unless you are currently growing from an apple eaten in the previous frame.

---

## 🧮 Scoring & Multiplier Calculations

$$\text{Points Earned} = \text{Base Points} \times \text{Active Multiplier}$$

| Action / Item | Base Score | With 3× Multiplier Active |
|---------------|------------|----------------------------|
| Standard Food (Classic / Survival) | +10 pts | +30 pts |
| Standard Food (Time Attack) | +25 pts | +75 pts |
| Slow-Mo Pickup | +20 pts | +60 pts |
| Ghost Pickup | +20 pts | +60 pts |
| Multiplier Pickup | +30 pts | +90 pts |
| Shrink Pickup | +15 pts | +45 pts |
| Poison Hazard (Avoid!) | -30 pts | -30 pts (floored at 0) |

---

## ⌨️ In-Game Controls & Shortcuts

| Action | Primary Keys | Vim / Alternative | Context |
|--------|--------------|-------------------|---------|
| **Start / Play Game** | `Enter` | `Space` | In Main Menu |
| **Turn Up** | `↑` (Up Arrow) | `W` / `K` | Playing / In Menu |
| **Turn Down** | `↓` (Down Arrow) | `S` / `J` | Playing / In Menu |
| **Turn Left** | `←` (Left Arrow) | `A` / `H` | Playing / Adjust Setting |
| **Turn Right** | `→` (Right Arrow) | `D` / `L` | Playing / Adjust Setting |
| **Pause / Resume** | `P` | `Space` | While playing |
| **Game Settings Menu** | `M` | `Tab` | Anytime (Active, Paused, Game Over) |
| **Instant Boss Key** | `B` | `Esc` | Instantly swaps screen to fake `htop` |
| **Restart Game** | `R` | — | When Game Over |
| **Quit** | `Q` | `Ctrl+C` | Anytime (auto-saves progress) |

---

## 📋 Main Menu & Settings Screen

When running `snake` or `snake play` without flags, the game opens directly into the interactive **Main Menu & Settings** screen before gameplay starts:

- **Start Game:** Press `Enter` on `[ START GAME ]` or press `Space` anywhere in the menu to start playing immediately.
- **Resume Saved Game:** If an unfinished session was saved (via Boss Key `B` or Quit `Q`), the menu offers `[ RESUME SAVED GAME (Score: X) ]`.
- **Adjust Settings:** Navigate with `↑`/`↓` and cycle `Mode`, `Difficulty`, `Board Size`, `Obstacles`, and `Theme` with `←`/`→`.
- **Live Sample Board Preview:** As you adjust `Board Size`, `Theme`, or `Obstacles`, the right side of the menu live-renders an interactive sample arena:
  - **Board Size & Proportions:** Dynamically scales the preview box between Small (20×12), Medium (28×16), Classic (24×20), and Large (36×20) with cell counts and aspect ratios.
  - **Visual Theme:** Live-transforms arena borders (rounded, normal, thick), snake head/body colors and glyphs, food apple, and bonus items across Dark, Light, High-Contrast, Colorblind-Safe, and Corporate themes.
  - **Procedural Obstacles:** Shows obstacle blocks (`██`) when Obstacles are toggled `ON` and hides them when `OFF`.
  - **Glyph Legend:** Displays the exact in-game symbols for Head, Body, Food, and Obstacles.
- **Direct Play (`--direct`):** Bypass the menu completely by passing `--direct` or any specific flags (e.g., `snake play --mode survival`, `snake -d hard`, or `snake --size large`).

---

## 📐 Board Size Presets

| Size Preset | Grid Dimensions (Cells) | Terminal Footprint (Columns × Rows) | Best For |
|-------------|-------------------------|-------------------------------------|----------|
| **Small**   | `20 × 12` | 40 cols × 12 rows | Compact split panes & small laptop windows |
| **Medium**  | `28 × 16` | 56 cols × 16 rows | Balanced arena with quick reflex dynamics |
| **Classic** | `24 × 20` | 48 cols × 20 rows | Standard retro terminal aspect ratio (Default) |
| **Large**   | `36 × 20` | 72 cols × 20 rows | Spacious board for mega lengths and high scores |

> **CLI Flag:** `-s, --size <preset|WxH>`
>
> Examples:
> - `snake play --size small`
> - `snake play -s large`
> - `snake play --size 30x15` (custom Width × Height)

---

## 👤 Author

**Md. Inzamamul Lohani** ([@hamimLohani](https://github.com/hamimLohani))

