package renderer

import (
	"fmt"
	"strings"
	"time"

	"snake/internal/engine"
	"snake/internal/storage"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tickMsg time.Time

type hardwareDirMsg engine.Direction

// UIModel is the Bubble Tea model controlling the interactive game.
type UIModel struct {
	Game               *engine.Game
	Theme              Theme
	Config             storage.Config
	Stats              storage.Stats
	TermWidth          int
	TermHeight         int
	Paused             bool
	BossScreen         bool
	BossType           string
	StartTime          time.Time
	FinalRecorded      bool
	IsBot              bool
	BotFunc            func(g *engine.Game) engine.Direction
	RecordedMoves      []storage.ReplayStep
	HardwareChan       chan engine.Direction
	InMenu             bool
	MenuIndex          int
	SelectedMode       int
	SelectedDifficulty int
	SelectedTheme      int
	SelectedObstacles  bool
	IsInitialStart     bool
	HasSaveGame        bool
}

// NewUIModel constructs an initialized Bubble Tea model.
func NewUIModel(game *engine.Game, cfg storage.Config, stats storage.Stats, themeName string, startInMenu bool) UIModel {
	th := GetTheme(themeName)
	if themeName == "" {
		th = GetTheme(cfg.Theme)
	}

	bossType := cfg.DefaultBossKey
	if bossType == "" {
		bossType = "htop"
	}

	high := stats.HighScores[string(game.Config.Mode)]
	game.Highscore = high

	// Initial mode index
	modes := []engine.Mode{engine.ModeClassic, engine.ModeTimeAttack, engine.ModeSurvival, engine.ModeZen}
	modeIdx := 0
	for i, m := range modes {
		if m == game.Config.Mode {
			modeIdx = i
			break
		}
	}

	// Initial difficulty index
	diffIdx := 1 // normal
	switch {
	case game.Config.BaseTickRateMs <= 45:
		diffIdx = 3 // insane
	case game.Config.BaseTickRateMs <= 75:
		diffIdx = 2 // hard
	case game.Config.BaseTickRateMs <= 110:
		diffIdx = 1 // normal
	default:
		diffIdx = 0 // easy
	}

	// Initial theme index
	themes := []string{"dark", "light", "high-contrast", "colorblind", "corporate"}
	themeIdx := 0
	for i, t := range themes {
		if t == th.Name {
			themeIdx = i
			break
		}
	}

	return UIModel{
		Game:               game,
		Theme:              th,
		Config:             cfg,
		Stats:              stats,
		BossType:           bossType,
		StartTime:          time.Now(),
		RecordedMoves:      make([]storage.ReplayStep, 0),
		SelectedMode:       modeIdx,
		SelectedDifficulty: diffIdx,
		SelectedTheme:      themeIdx,
		SelectedObstacles:  game.Config.WithObstacles,
		InMenu:             startInMenu,
		Paused:             startInMenu,
		IsInitialStart:     startInMenu,
		HasSaveGame:        storage.HasSaveGame(),
		MenuIndex:          0,
	}
}

// Init triggers the first tick schedule and tea subscriptions.
func (m UIModel) Init() tea.Cmd {
	cmds := []tea.Cmd{
		scheduleTick(m.Game.TickDuration()),
	}
	if m.HardwareChan != nil {
		cmds = append(cmds, listenHardware(m.HardwareChan))
	}
	return tea.Batch(cmds...)
}

func listenHardware(ch chan engine.Direction) tea.Cmd {
	return func() tea.Msg {
		if ch == nil {
			return nil
		}
		dir, ok := <-ch
		if !ok {
			return nil
		}
		return hardwareDirMsg(dir)
	}
}

func scheduleTick(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Update processes Bubble Tea events (ticks, keypresses, resize, blur).
func (m UIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case hardwareDirMsg:
		dir := engine.Direction(msg)
		if !m.Paused && !m.Game.GameOver {
			if m.Game.QueueDirection(dir) {
				m.RecordedMoves = append(m.RecordedMoves, storage.ReplayStep{
					Tick:      m.Game.TickCount,
					Direction: dir,
				})
			}
		}
		return m, listenHardware(m.HardwareChan)
	case tea.WindowSizeMsg:
		m.TermWidth = msg.Width
		m.TermHeight = msg.Height
		return m, nil

	case tea.FocusMsg:
		// Focus restored
		return m, nil

	case tea.BlurMsg:
		// Focus lost: auto-pause to prevent accidental death
		if !m.Game.GameOver && !m.BossScreen {
			m.Paused = true
		}
		return m, nil

	case tickMsg:
		if m.Game.GameOver {
			if !m.FinalRecorded {
				m.FinalRecorded = true
				duration := int(time.Since(m.StartTime).Seconds())
				m.Stats.RecordGame(storage.ScoreEntry{
					Date:        time.Now(),
					Mode:        string(m.Game.Config.Mode),
					Score:       m.Game.Score,
					Length:      len(m.Game.Snake),
					Seed:        m.Game.Config.Seed,
					DurationSec: duration,
				})
				_ = storage.SaveStats(m.Stats)
				_ = storage.ClearSaveGame()

				// Auto-record replay
				rep := storage.ReplayFile{
					Version:       "1.0.0",
					Date:          time.Now(),
					Mode:          m.Game.Config.Mode,
					Seed:          m.Game.Config.Seed,
					Width:         m.Game.Config.Width,
					Height:        m.Game.Config.Height,
					WithObstacles: m.Game.Config.WithObstacles,
					FinalScore:    m.Game.Score,
					FinalLength:   len(m.Game.Snake),
					TotalTicks:    m.Game.TickCount,
					Steps:         m.RecordedMoves,
				}
				_ = storage.SaveReplay(storage.LatestReplayPath(), rep)
			}
			return m, nil
		}

		if m.Paused || m.BossScreen {
			// Continue ticker so game resumes immediately upon unpause
			return m, scheduleTick(100 * time.Millisecond)
		}

		// If AI bot is active, compute move
		if m.IsBot && m.BotFunc != nil {
			dir := m.BotFunc(m.Game)
			if dir != engine.DirNone {
				m.Game.QueueDirection(dir)
			}
		}

		m.Game.Tick()

		return m, scheduleTick(m.Game.TickDuration())

	case tea.KeyMsg:
		keyStr := msg.String()

		// Boss key toggle
		if containsKey(m.Config.Keys.Boss, keyStr) {
			m.BossScreen = !m.BossScreen
			if m.BossScreen {
				m.Paused = true
				_ = storage.SaveActiveGame(m.Game)
			}
			return m, nil
		}

		// If on boss screen, any key returns to game
		if m.BossScreen {
			m.BossScreen = false
			return m, nil
		}

		// Quit
		if containsKey(m.Config.Keys.Quit, keyStr) || keyStr == "ctrl+c" {
			if !m.Game.GameOver {
				_ = storage.SaveActiveGame(m.Game)
			}
			return m, tea.Quit
		}

		// In-Game Menu toggle (M or Tab)
		if containsKey(m.Config.Keys.Menu, keyStr) || keyStr == "m" || keyStr == "M" || keyStr == "tab" {
			if m.IsInitialStart {
				// If at initial start, pressing M or Tab starts the game
				return m, m.applyAndRestart()
			}
			m.InMenu = !m.InMenu
			if m.InMenu {
				m.Paused = true
				m.MenuIndex = 0
			} else if !m.Game.GameOver {
				m.Paused = false
			}
			return m, nil
		}

		// Handle Menu Navigation if InMenu is open
		if m.InMenu {
			hasSave := m.IsInitialStart && m.HasSaveGame
			totalMenuRows := 6
			if hasSave || !m.IsInitialStart {
				totalMenuRows = 7
			}

			switch {
			case containsKey(m.Config.Keys.Up, keyStr) || keyStr == "up" || keyStr == "w" || keyStr == "k":
				m.MenuIndex = (m.MenuIndex - 1 + totalMenuRows) % totalMenuRows
				return m, nil
			case containsKey(m.Config.Keys.Down, keyStr) || keyStr == "down" || keyStr == "s" || keyStr == "j":
				m.MenuIndex = (m.MenuIndex + 1) % totalMenuRows
				return m, nil
			case containsKey(m.Config.Keys.Left, keyStr) || keyStr == "left" || keyStr == "a" || keyStr == "h":
				if sIdx := m.currentSettingIndex(); sIdx >= 0 {
					m.adjustMenuSetting(-1, sIdx)
				}
				return m, nil
			case containsKey(m.Config.Keys.Right, keyStr) || keyStr == "right" || keyStr == "d" || keyStr == "l":
				if sIdx := m.currentSettingIndex(); sIdx >= 0 {
					m.adjustMenuSetting(1, sIdx)
				}
				return m, nil
			case keyStr == "enter" || keyStr == " ":
				if sIdx := m.currentSettingIndex(); sIdx >= 0 {
					if keyStr == " " {
						if m.IsInitialStart {
							return m, m.applyAndRestart()
						}
						m.InMenu = false
						if !m.Game.GameOver {
							m.Paused = false
						}
						return m, nil
					}
					m.adjustMenuSetting(1, sIdx)
					return m, nil
				}

				if m.IsInitialStart {
					if hasSave {
						switch m.MenuIndex {
						case 0: // Resume Saved Game
							m.InMenu = false
							m.Paused = false
							m.IsInitialStart = false
							return m, scheduleTick(m.Game.TickDuration())
						case 1: // Start New Game
							return m, m.applyAndRestart()
						case 6: // Quit Game
							return m, tea.Quit
						}
					} else {
						switch m.MenuIndex {
						case 0: // Start Game
							return m, m.applyAndRestart()
						case 5: // Quit Game
							return m, tea.Quit
						}
					}
				} else {
					switch m.MenuIndex {
					case 0: // Resume Current Game
						m.InMenu = false
						if !m.Game.GameOver {
							m.Paused = false
						}
						return m, nil
					case 5: // Apply & Restart New Game
						return m, m.applyAndRestart()
					case 6: // Quit Game
						return m, tea.Quit
					}
				}
				return m, nil

			case keyStr == "esc" || keyStr == "p":
				if m.IsInitialStart {
					return m, m.applyAndRestart()
				}
				m.InMenu = false
				if !m.Game.GameOver {
					m.Paused = false
				}
				return m, nil

			case keyStr == "q" || keyStr == "Q":
				return m, tea.Quit
			}
			return m, nil
		}

		// Pause
		if containsKey(m.Config.Keys.Pause, keyStr) {
			if !m.Game.GameOver {
				m.Paused = !m.Paused
			}
			return m, nil
		}

		// Restart
		if containsKey(m.Config.Keys.Restart, keyStr) {
			if m.Game.GameOver {
				seed := time.Now().UnixNano()
				cfg := m.Game.Config
				cfg.Seed = seed
				m.Game = engine.New(cfg)
				m.Game.Highscore = m.Stats.HighScores[string(cfg.Mode)]
				m.StartTime = time.Now()
				m.FinalRecorded = false
				m.Paused = false
				m.RecordedMoves = make([]storage.ReplayStep, 0)
				return m, scheduleTick(m.Game.TickDuration())
			}
		}

		// Movement inputs
		if !m.Paused && !m.Game.GameOver && !m.IsBot {
			var targetDir engine.Direction
			switch {
			case containsKey(m.Config.Keys.Up, keyStr):
				targetDir = engine.DirUp
			case containsKey(m.Config.Keys.Down, keyStr):
				targetDir = engine.DirDown
			case containsKey(m.Config.Keys.Left, keyStr):
				targetDir = engine.DirLeft
			case containsKey(m.Config.Keys.Right, keyStr):
				targetDir = engine.DirRight
			}

			if targetDir != engine.DirNone {
				queued := m.Game.QueueDirection(targetDir)
				if queued {
					m.RecordedMoves = append(m.RecordedMoves, storage.ReplayStep{
						Tick:      m.Game.TickCount,
						Direction: targetDir,
					})
				}
			}
		}
	}

	return m, nil
}

func containsKey(bindings []string, key string) bool {
	for _, b := range bindings {
		if strings.EqualFold(b, key) {
			return true
		}
	}
	return false
}

// View constructs the full visual terminal interface.
func (m UIModel) View() string {
	if m.BossScreen {
		return RenderBossScreen(m.BossType, m.TermWidth, m.TermHeight)
	}

	if m.Config.ScreenReader {
		return m.renderScreenReader()
	}

	return m.renderRichGame()
}

func (m UIModel) renderScreenReader() string {
	head := m.Game.Snake[0]
	return fmt.Sprintf("Snake CLI. Mode: %s. Score: %d. Length: %d. Head at (%d, %d). Food at (%d, %d). Direction: %s. Game Over: %t.\n",
		m.Game.Config.Mode, m.Game.Score, len(m.Game.Snake), head.X, head.Y, m.Game.Food.X, m.Game.Food.Y, m.Game.Direction.String(), m.Game.GameOver)
}

func (m UIModel) currentSettingIndex() int {
	if m.IsInitialStart && m.HasSaveGame {
		if m.MenuIndex >= 2 && m.MenuIndex <= 5 {
			return m.MenuIndex - 2
		}
		return -1
	}
	if m.MenuIndex >= 1 && m.MenuIndex <= 4 {
		return m.MenuIndex - 1
	}
	return -1
}

func (m *UIModel) adjustMenuSetting(delta int, settingIdx int) {
	switch settingIdx {
	case 0: // Mode
		totalModes := 4
		m.SelectedMode = (m.SelectedMode + delta + totalModes) % totalModes
	case 1: // Difficulty
		totalDiffs := 4
		m.SelectedDifficulty = (m.SelectedDifficulty + delta + totalDiffs) % totalDiffs
		if m.SelectedDifficulty == 3 {
			m.SelectedObstacles = true
		}
	case 2: // Obstacles
		m.SelectedObstacles = !m.SelectedObstacles
	case 3: // Theme
		totalThemes := 5
		m.SelectedTheme = (m.SelectedTheme + delta + totalThemes) % totalThemes
		themes := []string{"dark", "light", "high-contrast", "colorblind", "corporate"}
		m.Theme = GetTheme(themes[m.SelectedTheme])
		m.Config.Theme = themes[m.SelectedTheme]
	}
}

func (m *UIModel) applyAndRestart() tea.Cmd {
	modes := []engine.Mode{engine.ModeClassic, engine.ModeTimeAttack, engine.ModeSurvival, engine.ModeZen}
	chosenMode := modes[m.SelectedMode]

	diffs := []string{"easy", "normal", "hard", "insane"}
	chosenDiff := diffs[m.SelectedDifficulty]

	var speed int
	obstacles := m.SelectedObstacles
	switch chosenDiff {
	case "easy":
		speed = 160
	case "normal":
		speed = 110
	case "hard":
		speed = 75
	case "insane":
		speed = 45
		obstacles = true
	}

	themes := []string{"dark", "light", "high-contrast", "colorblind", "corporate"}
	chosenTheme := themes[m.SelectedTheme]
	m.Theme = GetTheme(chosenTheme)
	m.Config.Theme = chosenTheme

	seed := time.Now().UnixNano()
	cfg := engine.Config{
		Width:           m.Config.GridWidth,
		Height:          m.Config.GridHeight,
		Mode:            chosenMode,
		Seed:            seed,
		InitialSnakeLen: 3,
		AdaptiveSpeed:   m.Config.AdaptiveSpeed,
		BaseTickRateMs:  speed,
		WithObstacles:   obstacles,
	}

	m.Game = engine.New(cfg)
	m.Game.Highscore = m.Stats.HighScores[string(chosenMode)]
	m.StartTime = time.Now()
	m.FinalRecorded = false
	m.Paused = false
	m.InMenu = false
	m.IsInitialStart = false
	m.HasSaveGame = false
	m.RecordedMoves = make([]storage.ReplayStep, 0)
	_ = storage.ClearSaveGame()

	return scheduleTick(m.Game.TickDuration())
}

func (m UIModel) renderMenuScreen() string {
	modes := []string{"CLASSIC", "TIME ATTACK", "SURVIVAL", "ZEN"}
	diffs := []string{"EASY (160ms)", "NORMAL (110ms)", "HARD (75ms)", "INSANE (45ms)"}
	themes := []string{"DARK", "LIGHT", "HIGH-CONTRAST", "COLORBLIND", "CORPORATE"}

	obsStr := "OFF"
	if m.SelectedObstacles {
		obsStr = "ON"
	}

	selectedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00E5FF")).Bold(true)
	cursorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF87")).Bold(true)
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true)
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#7f8c8d"))
	valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD600")).Bold(true)
	titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF87")).Bold(true)
	headerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#7928CA")).Bold(true)

	formatRow := func(idx int, label, val string) string {
		cursor := "  "
		displayLabel := labelStyle.Render(label)
		if m.MenuIndex == idx {
			cursor = cursorStyle.Render("▶ ")
			displayLabel = selectedStyle.Render(label)
		}
		arrowL := dimStyle.Render("◀ ")
		arrowR := dimStyle.Render(" ▶")
		if m.MenuIndex == idx {
			arrowL = selectedStyle.Render("◀ ")
			arrowR = selectedStyle.Render(" ▶")
		}
		return fmt.Sprintf("%s%-14s %s%-16s%s", cursor, displayLabel, arrowL, valueStyle.Render(val), arrowR)
	}

	formatButton := func(idx int, text string, primary bool) string {
		cursor := "  "
		var btnText string
		if m.MenuIndex == idx {
			cursor = cursorStyle.Render("▶ ")
			if primary {
				btnText = lipgloss.NewStyle().
					Background(lipgloss.Color("#00FF87")).
					Foreground(lipgloss.Color("#000000")).
					Bold(true).
					Padding(0, 2).
					Render(text)
			} else {
				btnText = selectedStyle.Render("[ " + text + " ]")
			}
		} else {
			if primary {
				btnText = lipgloss.NewStyle().
					Foreground(lipgloss.Color("#00FF87")).
					Bold(true).
					Render("[ " + text + " ]")
			} else {
				btnText = dimStyle.Render("[ " + text + " ]")
			}
		}
		return fmt.Sprintf("%s%s", cursor, btnText)
	}

	var titleBanner string
	var lines []string

	hasSave := m.IsInitialStart && m.HasSaveGame

	if m.IsInitialStart {
		titleBanner = titleStyle.Render("🐍  S N A K E   C L I  🐍")
		subBanner := headerStyle.Render("─── MAIN MENU & SETTINGS ───")
		if hasSave {
			lines = []string{
				titleBanner,
				subBanner,
				"",
				formatButton(0, fmt.Sprintf("▶ RESUME SAVED GAME (Score: %d)", m.Game.Score), true),
				formatButton(1, "Start New Game", false),
				"",
				formatRow(2, "Mode:", modes[m.SelectedMode]),
				formatRow(3, "Difficulty:", diffs[m.SelectedDifficulty]),
				formatRow(4, "Obstacles:", obsStr),
				formatRow(5, "Theme:", themes[m.SelectedTheme]),
				"",
				formatButton(6, "Quit Game", false),
				"",
				dimStyle.Render(" [↑/↓] Navigate  •  [←/→] Change Setting"),
				dimStyle.Render(" [Enter] Select  •  [Space] Play Now  •  [Q] Quit"),
			}
		} else {
			lines = []string{
				titleBanner,
				subBanner,
				"",
				formatButton(0, "▶ START GAME", true),
				"",
				formatRow(1, "Mode:", modes[m.SelectedMode]),
				formatRow(2, "Difficulty:", diffs[m.SelectedDifficulty]),
				formatRow(3, "Obstacles:", obsStr),
				formatRow(4, "Theme:", themes[m.SelectedTheme]),
				"",
				formatButton(5, "Quit Game", false),
				"",
				dimStyle.Render(" [↑/↓] Navigate  •  [←/→] Change Setting"),
				dimStyle.Render(" [Enter] Select  •  [Space] Play Now  •  [Q] Quit"),
			}
		}
	} else {
		titleBanner = titleStyle.Render("🐍  G A M E   P A U S E D  🐍")
		subBanner := headerStyle.Render("─── SETTINGS & OPTIONS ───")
		lines = []string{
			titleBanner,
			subBanner,
			"",
			formatButton(0, "▶ RESUME GAME", true),
			"",
			formatRow(1, "Mode:", modes[m.SelectedMode]),
			formatRow(2, "Difficulty:", diffs[m.SelectedDifficulty]),
			formatRow(3, "Obstacles:", obsStr),
			formatRow(4, "Theme:", themes[m.SelectedTheme]),
			"",
			formatButton(5, "Apply & Restart New Game", false),
			formatButton(6, "Quit Game", false),
			"",
			dimStyle.Render(" [↑/↓] Navigate  •  [←/→] Change Setting"),
			dimStyle.Render(" [Enter] Select  •  [Esc / P] Resume  •  [Q] Quit"),
		}
	}

	innerContent := strings.Join(lines, "\n")
	w := (m.Game.Bounds.MaxX - m.Game.Bounds.MinX) * 2
	h := m.Game.Bounds.MaxY - m.Game.Bounds.MinY

	placed := lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center, innerContent)
	return m.Theme.BorderStyle.Render(placed)
}

func (m UIModel) renderRichGame() string {
	var sb strings.Builder

	// 1. HUD
	modeTag := strings.ToUpper(string(m.Game.Config.Mode))
	modeBadge := lipgloss.NewStyle().
		Background(lipgloss.Color("#7928CA")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true).
		Padding(0, 1).
		Render(modeTag)

	hudLeft := fmt.Sprintf("%s  %s %s",
		modeBadge,
		m.Theme.HUDLabelStyle.Render("SCORE:"),
		m.Theme.ScoreStyle.Render(fmt.Sprintf("%04d", m.Game.Score)),
	)

	highVal := m.Game.Highscore
	if m.Game.Score > highVal {
		highVal = m.Game.Score
	}
	hudRight := fmt.Sprintf("%s %s  %s %s",
		m.Theme.HUDLabelStyle.Render("HIGH:"),
		m.Theme.HighScoreStyle.Render(fmt.Sprintf("%04d", highVal)),
		m.Theme.HUDLabelStyle.Render("LEN:"),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF87")).Bold(true).Render(fmt.Sprintf("%d", len(m.Game.Snake))),
	)

	if m.Game.Config.Mode == engine.ModeTimeAttack {
		timerBadge := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0055")).Bold(true).Render(fmt.Sprintf(" ⏱ %ds", m.Game.TimeRemainingSec))
		hudLeft += timerBadge
	}

	if m.Game.ActivePowerUp != nil {
		pBadge := lipgloss.NewStyle().
			Background(lipgloss.Color("#FF0055")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true).
			Padding(0, 1).
			Render(fmt.Sprintf("%s %dt", m.Game.ActivePowerUp.Type.String(), m.Game.ActivePowerUp.ExpiresAtTick-m.Game.TickCount))
		hudLeft += "  " + pBadge
	}

	hudRow := lipgloss.JoinHorizontal(lipgloss.Center, hudLeft, "    ", hudRight)

	// 2. Grid Canvas
	var board strings.Builder
	for y := m.Game.Bounds.MinY; y < m.Game.Bounds.MaxY; y++ {
		for x := m.Game.Bounds.MinX; x < m.Game.Bounds.MaxX; x++ {
			pt := engine.Point{X: x, Y: y}
			switch {
			case m.Game.IsSnakeHead(pt):
				board.WriteString(m.Theme.HeadStyle.Render(m.Theme.HeadGlyph(m.Game.Direction.String())))
			case m.Game.IsSnakeBody(pt):
				board.WriteString(m.Theme.BodyStyle.Render(m.Theme.BodyGlyph))
			case pt == m.Game.Food:
				switch m.Game.FoodType {
				case engine.ItemSlowMo:
					board.WriteString(m.Theme.SlowMoStyle.Render(m.Theme.SlowMoGlyph))
				case engine.ItemGhost:
					board.WriteString(m.Theme.GhostStyle.Render(m.Theme.GhostGlyph))
				case engine.ItemMultiplier:
					board.WriteString(m.Theme.MultiplierStyle.Render(m.Theme.MultiplierGlyph))
				case engine.ItemShrink:
					board.WriteString(m.Theme.ShrinkStyle.Render(m.Theme.ShrinkGlyph))
				default:
					board.WriteString(m.Theme.FoodStyle.Render(m.Theme.FoodGlyph))
				}
			case m.Game.Obstacles[pt]:
				board.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#555555")).Render("██"))
			default:
				// Subtle grid dots
				board.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#222222")).Render("· "))
			}
		}
		if y < m.Game.Bounds.MaxY-1 {
			board.WriteString("\n")
		}
	}

	borderedBoard := m.Theme.BorderStyle.Render(board.String())

	// 3. Status Overlays
	var statusOverlay string
	if m.Game.GameOver {
		gameOverBox := m.Theme.ModalStyle.Render(fmt.Sprintf(
			"╔═══════════════════════════════════╗\n"+
				"║             GAME OVER             ║\n"+
				"║  Reason: %-24s ║\n"+
				"║  Final Score: %-19d ║\n"+
				"║  Length:      %-19d ║\n"+
				"║                                   ║\n"+
				"║  [R] Restart                      ║\n"+
				"║  [M] Change Mode / Difficulty     ║\n"+
				"║  [Q] Quit                         ║\n"+
				"╚═══════════════════════════════════╝",
			m.Game.GameOverReason, m.Game.Score, len(m.Game.Snake),
		))
		statusOverlay = gameOverBox
	} else if m.Paused {
		pauseBox := m.Theme.ModalStyle.Render(
			"╔═══════════════════════════════════╗\n" +
				"║              PAUSED               ║\n" +
				"║                                   ║\n" +
				"║  Press [P] or Space to resume     ║\n" +
				"║  Press [M] for Mode & Difficulty  ║\n" +
				"║  Press [B] for Boss Screen        ║\n" +
				"║  Press [Q] to quit                ║\n" +
				"╚═══════════════════════════════════╝",
		)
		statusOverlay = pauseBox
	}

	// 4. Instructions Footer
	footer := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666")).Render(
		"[W/A/S/D / Arrows] Move  •  [P] Pause  •  [M] Menu  •  [B] Boss Key  •  [R] Restart  •  [Q] Quit",
	)
	if m.InMenu {
		footer = lipgloss.NewStyle().Foreground(lipgloss.Color("#666666")).Render(
			"[↑/↓ / W/S] Navigate  •  [←/→ / A/D] Change Setting  •  [Enter / Space] Play  •  [Q] Quit",
		)
	}

	// Combine components
	var content string
	if m.InMenu {
		content = lipgloss.JoinVertical(lipgloss.Center,
			hudRow,
			m.renderMenuScreen(),
			footer,
		)
	} else if statusOverlay != "" {
		content = lipgloss.JoinVertical(lipgloss.Center,
			hudRow,
			borderedBoard,
			statusOverlay,
			footer,
		)
	} else {
		content = lipgloss.JoinVertical(lipgloss.Center,
			hudRow,
			borderedBoard,
			footer,
		)
	}

	sb.WriteString(content)

	// Center view in terminal if dimensions are known
	if m.TermWidth > 0 && m.TermHeight > 0 {
		return lipgloss.Place(m.TermWidth, m.TermHeight, lipgloss.Center, lipgloss.Center, sb.String())
	}

	return sb.String()
}
