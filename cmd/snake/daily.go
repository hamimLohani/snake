package main

import (
	"fmt"
	"hash/fnv"
	"os"
	"time"

	"snake/internal/engine"
	"snake/internal/renderer"
	"snake/internal/storage"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var dailyTheme string

func getDailySeed(t time.Time) int64 {
	dateStr := t.UTC().Format("2006-01-02")
	h := fnv.New64a()
	h.Write([]byte(dateStr))
	return int64(h.Sum64())
}

var dailyCmd = &cobra.Command{
	Use:   "daily",
	Short: "Play the daily challenge with today's global seed",
	Long:  `Play a competitive daily challenge. Everyone worldwide receives the exact same seed, food spawns, and obstacles for the day.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := storage.LoadConfig()
		stats, _ := storage.LoadStats()

		today := time.Now().UTC()
		seed := getDailySeed(today)
		dateStr := today.Format("2006-01-02")

		theme := dailyTheme
		if theme == "" {
			theme = cfg.Theme
		}

		engineCfg := engine.Config{
			Width:           cfg.GridWidth,
			Height:          cfg.GridHeight,
			Mode:            engine.ModeClassic,
			Seed:            seed,
			InitialSnakeLen: 3,
			AdaptiveSpeed:   cfg.AdaptiveSpeed,
			BaseTickRateMs:  cfg.BaseSpeedMs,
			WithObstacles:   true, // Daily challenge features procedural obstacles!
		}

		game := engine.New(engineCfg)
		model := renderer.NewUIModel(game, cfg, stats, theme, false)

		p := tea.NewProgram(
			model,
			tea.WithAltScreen(),
			tea.WithReportFocus(),
		)

		fmt.Printf("Starting Daily Challenge for UTC %s (Seed: %d)...\n", dateStr, seed)

		if _, err := p.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error in daily challenge: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	dailyCmd.Flags().StringVar(&dailyTheme, "theme", "", "Visual theme override")
	rootCmd.AddCommand(dailyCmd)
}
