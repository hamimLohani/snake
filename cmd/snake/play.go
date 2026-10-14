package main

import (
	"fmt"
	"os"
	"time"

	"snake/internal/engine"
	"snake/internal/input"
	"snake/internal/renderer"
	"snake/internal/storage"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var (
	playMode         string
	playSeed         int64
	playTheme        string
	playObstacles    bool
	playResume       bool
	playScreenReader bool
	playSpeed        int
	playDifficulty   string
	playHardwareUDP  string
)

var playCmd = &cobra.Command{
	Use:   "play",
	Short: "Start playing a game of Snake",
	Long:  `Launch an interactive game of Snake with customizable mode, seed, and visual theme.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := storage.LoadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not load config: %v\n", err)
		}

		stats, err := storage.LoadStats()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not load stats: %v\n", err)
		}

		// Mode resolution
		mode := engine.ModeClassic
		switch playMode {
		case "time-attack":
			mode = engine.ModeTimeAttack
		case "survival":
			mode = engine.ModeSurvival
		case "zen":
			mode = engine.ModeZen
		case "classic", "":
			mode = engine.ModeClassic
		default:
			fmt.Fprintf(os.Stderr, "Unknown mode '%s'. Valid modes: classic, time-attack, survival, zen\n", playMode)
			os.Exit(1)
		}

		seed := playSeed
		if seed == 0 {
			seed = time.Now().UnixNano()
		}

		theme := playTheme
		if theme == "" {
			theme = cfg.Theme
		}

		switch playDifficulty {
		case "easy":
			cfg.BaseSpeedMs = 160
		case "normal", "medium":
			cfg.BaseSpeedMs = 110
		case "hard":
			cfg.BaseSpeedMs = 75
		case "insane":
			cfg.BaseSpeedMs = 45
			playObstacles = true
		case "":
			// Keep existing or config speed
		default:
			fmt.Fprintf(os.Stderr, "Unknown difficulty '%s'. Valid options: easy, normal, hard, insane\n", playDifficulty)
			os.Exit(1)
		}

		if playSpeed > 0 {
			cfg.BaseSpeedMs = playSpeed
		}
		if playScreenReader {
			cfg.ScreenReader = true
		}

		var game *engine.Game
		if playResume && storage.HasSaveGame() {
			savedGame, err := storage.LoadSavedGame()
			if err == nil {
				game = savedGame
			}
		}

		if game == nil {
			engineCfg := engine.Config{
				Width:           cfg.GridWidth,
				Height:          cfg.GridHeight,
				Mode:            mode,
				Seed:            seed,
				InitialSnakeLen: 3,
				AdaptiveSpeed:   cfg.AdaptiveSpeed,
				BaseTickRateMs:  cfg.BaseSpeedMs,
				WithObstacles:   playObstacles,
			}
			game = engine.New(engineCfg)
		}
		model := renderer.NewUIModel(game, cfg, stats, theme)

		if playHardwareUDP != "" {
			hwChan := make(chan engine.Direction, 16)
			stopHw := make(chan struct{})
			defer close(stopHw)

			listener := input.NewHardwareInputListener(playHardwareUDP)
			if err := listener.StartUDP(hwChan, stopHw); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to start hardware listener: %v\n", err)
			} else {
				model.HardwareChan = hwChan
			}
		}

		p := tea.NewProgram(
			model,
			tea.WithAltScreen(),
			tea.WithReportFocus(),
		)

		if _, err := p.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error running game: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	playCmd.Flags().StringVar(&playMode, "mode", "classic", "Game mode: classic, time-attack, survival, zen")
	playCmd.Flags().Int64Var(&playSeed, "seed", 0, "Deterministic RNG seed (0 = random)")
	playCmd.Flags().StringVar(&playTheme, "theme", "", "Visual theme: dark, light, high-contrast, colorblind, corporate")
	playCmd.Flags().BoolVar(&playObstacles, "obstacles", false, "Generate procedural obstacles from seed")
	playCmd.Flags().BoolVar(&playResume, "resume", true, "Auto-resume previous unfinished game if available")
	playCmd.Flags().BoolVar(&playScreenReader, "screen-reader", false, "Accessible screen-reader output mode")
	playCmd.Flags().IntVar(&playSpeed, "speed", 0, "Override tick speed in milliseconds (e.g. 150 for slower, 70 for faster)")
	playCmd.Flags().StringVarP(&playDifficulty, "difficulty", "d", "", "Difficulty preset: easy, normal, hard, insane")
	playCmd.Flags().StringVar(&playHardwareUDP, "hardware-udp", "", "Listen for ESP8266/ESP32 inputs on UDP addr (e.g. ':4242')")
	rootCmd.AddCommand(playCmd)
}
