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
	botHeadless bool
	botTicks    int
	botSeed     int64
	botSpeedMs  int
	botTheme    string
)

var botCmd = &cobra.Command{
	Use:   "bot",
	Short: "Watch the AI bot play itself (interactive or headless)",
	Long:  `Run the autonomous heuristic BFS solver. Can run interactively with full TUI or headless for automated testing and benchmarks.`,
	Run: func(cmd *cobra.Command, args []string) {
		seed := botSeed
		if seed == 0 {
			seed = time.Now().UnixNano()
		}

		cfg, _ := storage.LoadConfig()
		stats, _ := storage.LoadStats()

		engineCfg := engine.Config{
			Width:           cfg.GridWidth,
			Height:          cfg.GridHeight,
			Mode:            engine.ModeClassic,
			Seed:            seed,
			InitialSnakeLen: 3,
			AdaptiveSpeed:   false,
			BaseTickRateMs:  50,
			WithObstacles:   false,
		}

		game := engine.New(engineCfg)
		solver := input.NewBotSolver()

		// Headless mode for CI, testing and benchmarks
		if botHeadless {
			maxTicks := botTicks
			if maxTicks <= 0 {
				maxTicks = 1500
			}

			start := time.Now()
			for !game.GameOver && game.TickCount < uint64(maxTicks) {
				nextDir := solver.NextMove(game)
				if nextDir != engine.DirNone {
					game.QueueDirection(nextDir)
				}
				game.Tick()
			}
			elapsed := time.Since(start)

			fmt.Printf("=== Headless AI Bot Simulation Complete ===\n")
			fmt.Printf("Seed:         %d\n", seed)
			fmt.Printf("Final Score:  %d\n", game.Score)
			fmt.Printf("Final Length: %d\n", len(game.Snake))
			fmt.Printf("Total Ticks:  %d\n", game.TickCount)
			fmt.Printf("Sim Time:     %s (%.0f ticks/sec)\n", elapsed, float64(game.TickCount)/elapsed.Seconds())
			fmt.Printf("Game Over:    %t (%s)\n", game.GameOver, game.GameOverReason)
			return
		}

		// Interactive visual mode
		if botSpeedMs > 0 {
			game.Config.BaseTickRateMs = botSpeedMs
		} else {
			game.Config.BaseTickRateMs = 65 // Fast, smooth demo speed
		}

		theme := botTheme
		if theme == "" {
			theme = cfg.Theme
		}

		model := renderer.NewUIModel(game, cfg, stats, theme, false)
		model.IsBot = true
		model.BotFunc = solver.NextMove

		p := tea.NewProgram(
			model,
			tea.WithAltScreen(),
		)

		if _, err := p.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error in AI bot demo: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	botCmd.Flags().BoolVar(&botHeadless, "headless", false, "Run in headless mode without TUI at full CPU speed")
	botCmd.Flags().IntVar(&botTicks, "ticks", 1500, "Maximum ticks to simulate in headless mode")
	botCmd.Flags().Int64Var(&botSeed, "seed", 0, "Deterministic RNG seed")
	botCmd.Flags().IntVar(&botSpeedMs, "speed", 0, "Override visual tick speed in ms (default 65ms)")
	botCmd.Flags().StringVarP(&botTheme, "theme", "t", "", "Visual theme override")
	rootCmd.AddCommand(botCmd)
}
