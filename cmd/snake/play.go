package main

import (
	"fmt"
	"os"
	"time"

	"snake/internal/engine"
	"snake/internal/renderer"
	"snake/internal/storage"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var (
	playMode  string
	playSeed  int64
	playTheme string
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

		engineCfg := engine.Config{
			Width:           cfg.GridWidth,
			Height:          cfg.GridHeight,
			Mode:            mode,
			Seed:            seed,
			InitialSnakeLen: 3,
			AdaptiveSpeed:   cfg.AdaptiveSpeed,
			BaseTickRateMs:  cfg.BaseSpeedMs,
		}

		game := engine.New(engineCfg)
		model := renderer.NewUIModel(game, cfg, stats, theme)

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
	rootCmd.AddCommand(playCmd)
}
