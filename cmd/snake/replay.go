package main

import (
	"fmt"
	"os"

	"snake/internal/engine"
	"snake/internal/renderer"
	"snake/internal/storage"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var (
	replaySpeedMs int
	replayVerify  bool
)

var replayCmd = &cobra.Command{
	Use:   "replay [file]",
	Short: "Playback a recorded game replay file",
	Long:  `Re-run a deterministic recorded game step-by-step or verify determinism. Defaults to latest recorded game if no file is given.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := storage.LatestReplayPath()
		if len(args) > 0 {
			path = args[0]
		}

		rep, err := storage.LoadReplay(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading replay: %v\n", err)
			os.Exit(1)
		}

		cfg, _ := storage.LoadConfig()
		stats, _ := storage.LoadStats()

		if replaySpeedMs > 0 {
			cfg.BaseSpeedMs = replaySpeedMs
		}

		engineCfg := engine.Config{
			Width:           rep.Width,
			Height:          rep.Height,
			Mode:            rep.Mode,
			Seed:            rep.Seed,
			InitialSnakeLen: 3,
			AdaptiveSpeed:   false,
			BaseTickRateMs:  cfg.BaseSpeedMs,
			WithObstacles:   rep.WithObstacles,
		}

		game := engine.New(engineCfg)

		// Map steps by tick for instant lookup
		stepMap := make(map[uint64]engine.Direction, len(rep.Steps))
		for _, s := range rep.Steps {
			stepMap[s.Tick] = s.Direction
		}

		// If verification mode requested, run headless and assert final score
		if replayVerify {
			for !game.GameOver && game.TickCount <= rep.TotalTicks+50 {
				if dir, ok := stepMap[game.TickCount]; ok {
					game.QueueDirection(dir)
				}
				game.Tick()
			}
			if game.Score != rep.FinalScore {
				fmt.Fprintf(os.Stderr, "VERIFICATION FAILED: Expected score %d, got %d\n", rep.FinalScore, game.Score)
				os.Exit(1)
			}
			if len(game.Snake) != rep.FinalLength {
				fmt.Fprintf(os.Stderr, "VERIFICATION FAILED: Expected length %d, got %d\n", rep.FinalLength, len(game.Snake))
				os.Exit(1)
			}
			fmt.Printf("✓ Replay verified successfully! Deterministic match: Score %d, Length %d across %d ticks.\n",
				game.Score, len(game.Snake), game.TickCount)
			return
		}

		// Interactive replay viewer
		model := renderer.NewUIModel(game, cfg, stats, cfg.Theme)
		model.IsBot = true
		model.BotFunc = func(g *engine.Game) engine.Direction {
			if d, ok := stepMap[g.TickCount]; ok {
				return d
			}
			return engine.DirNone
		}

		p := tea.NewProgram(
			model,
			tea.WithAltScreen(),
		)

		if _, err := p.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error in replay viewer: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	replayCmd.Flags().IntVar(&replaySpeedMs, "speed", 0, "Playback tick speed in ms (e.g. 50 for 2x speed)")
	replayCmd.Flags().BoolVar(&replayVerify, "verify", false, "Run headless determinism verification")
	rootCmd.AddCommand(replayCmd)
}
