package main

import (
	"fmt"
	"os"
	"strings"
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
	playDirect       bool
	playSize         string
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

		// Detect if explicit options/flags were provided on CLI
		hasExplicitFlags := cmd.Flags().Changed("mode") ||
			cmd.Flags().Changed("difficulty") ||
			cmd.Flags().Changed("speed") ||
			cmd.Flags().Changed("obstacles") ||
			cmd.Flags().Changed("seed") ||
			cmd.Flags().Changed("theme") ||
			cmd.Flags().Changed("screen-reader") ||
			cmd.Flags().Changed("hardware-udp") ||
			cmd.Flags().Changed("direct") ||
			cmd.Flags().Changed("size")

		startInMenu := !playDirect && !hasExplicitFlags

		if playSize != "" {
			cleanSize := strings.ToLower(playSize)
			switch cleanSize {
			case "small":
				cfg.GridWidth = 20
				cfg.GridHeight = 12
			case "medium":
				cfg.GridWidth = 28
				cfg.GridHeight = 16
			case "classic":
				cfg.GridWidth = 24
				cfg.GridHeight = 20
			case "large":
				cfg.GridWidth = 36
				cfg.GridHeight = 20
			default:
				var w, h int
				if n, _ := fmt.Sscanf(cleanSize, "%dx%d", &w, &h); n == 2 && w >= 10 && h >= 8 {
					cfg.GridWidth = w
					cfg.GridHeight = h
				} else {
					fmt.Fprintf(os.Stderr, "Unknown size '%s'. Valid options: small, medium, classic, large, or WxH (e.g. 30x15)\n", playSize)
					os.Exit(1)
				}
			}
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
		resumeAllowed := playResume && (!hasExplicitFlags || cmd.Flags().Changed("resume"))
		if resumeAllowed && storage.HasSaveGame() {
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
		model := renderer.NewUIModel(game, cfg, stats, theme, startInMenu)

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
	playCmd.Flags().StringVarP(&playMode, "mode", "m", "classic", "Game mode: classic, time-attack, survival, zen")
	playCmd.Flags().Int64Var(&playSeed, "seed", 0, "Deterministic RNG seed (0 = random)")
	playCmd.Flags().StringVarP(&playTheme, "theme", "t", "", "Visual theme: dark, light, high-contrast, colorblind, corporate")
	playCmd.Flags().BoolVar(&playObstacles, "obstacles", false, "Generate procedural obstacles from seed")
	playCmd.Flags().BoolVar(&playResume, "resume", true, "Auto-resume previous unfinished game if available")
	playCmd.Flags().BoolVar(&playScreenReader, "screen-reader", false, "Accessible screen-reader output mode")
	playCmd.Flags().IntVar(&playSpeed, "speed", 0, "Override tick speed in milliseconds (e.g. 150 for slower, 70 for faster)")
	playCmd.Flags().StringVarP(&playDifficulty, "difficulty", "d", "", "Difficulty preset: easy, normal, hard, insane")
	playCmd.Flags().StringVarP(&playSize, "size", "s", "", "Board size preset (small, medium, classic, large) or WxH (e.g. 30x15)")
	playCmd.Flags().StringVar(&playHardwareUDP, "hardware-udp", "", "Listen for ESP8266/ESP32 inputs on UDP addr (e.g. ':4242')")
	playCmd.Flags().BoolVar(&playDirect, "direct", false, "Start game directly without showing the main menu")
	rootCmd.AddCommand(playCmd)
	rootCmd.Flags().AddFlagSet(playCmd.Flags())
}
