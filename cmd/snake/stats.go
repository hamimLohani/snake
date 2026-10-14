package main

import (
	"fmt"
	"os"
	"time"

	"snake/internal/storage"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	statsJSON bool
	statsCSV  bool
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "View your local leaderboard, achievements, and gameplay statistics",
	Long:  `Display overall play metrics, high scores by mode, full leaderboard, and achievements. Supports --json and --csv export.`,
	Run: func(cmd *cobra.Command, args []string) {
		st, err := storage.LoadStats()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading stats: %v\n", err)
			os.Exit(1)
		}

		if statsJSON {
			data, err := st.ExportJSON()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error generating JSON: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(data)
			return
		}

		if statsCSV {
			data, err := st.ExportCSV()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error generating CSV: %v\n", err)
				os.Exit(1)
			}
			fmt.Print(data)
			return
		}

		// Rich Lipgloss view
		titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF87")).Bold(true)
		headerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00E5FF")).Bold(true)
		accentStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD600")).Bold(true)
		dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
		boxStyle := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#7928CA")).Padding(1, 2)

		fmt.Println(titleStyle.Render("=== Snake CLI Player Stats & Leaderboard ==="))
		fmt.Println()

		// Summary metrics
		duration := time.Duration(st.TotalPlayTimeSec) * time.Second
		fmt.Printf("%s %s    %s %s    %s %s\n",
			dimStyle.Render("Total Games:"), accentStyle.Render(fmt.Sprintf("%d", st.TotalGames)),
			dimStyle.Render("Total Time:"), accentStyle.Render(duration.String()),
			dimStyle.Render("Apples Eaten:"), accentStyle.Render(fmt.Sprintf("%d", st.TotalFoodEaten)),
		)
		fmt.Println()

		// High Scores by Mode
		fmt.Println(headerStyle.Render("High Scores by Mode:"))
		modes := []string{"classic", "time-attack", "survival", "zen"}
		for _, m := range modes {
			score := st.HighScores[m]
			fmt.Printf("  • %-14s : %s pts\n", m, accentStyle.Render(fmt.Sprintf("%d", score)))
		}
		fmt.Println()

		// Achievements
		fmt.Println(headerStyle.Render("Unlocked Achievements:"))
		if len(st.Achievements) == 0 {
			fmt.Println(dimStyle.Render("  No achievements unlocked yet. Play more games to unlock!"))
		} else {
			for ach := range st.Achievements {
				fmt.Printf("  🏆 %s\n", lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF87")).Render(ach))
			}
		}
		fmt.Println()

		// Leaderboard Table
		fmt.Println(headerStyle.Render("Top Leaderboard Entries:"))
		if len(st.Leaderboard) == 0 {
			fmt.Println(dimStyle.Render("  No games recorded yet."))
		} else {
			fmt.Printf("%-5s %-12s %-8s %-8s %-8s %s\n", "Rank", "Mode", "Score", "Length", "Time", "Date")
			fmt.Println(dimStyle.Render("------------------------------------------------------------"))
			for i, entry := range st.Leaderboard {
				if i >= 10 {
					break // Top 10 in terminal view
				}
				dateStr := entry.Date.Format("2006-01-02 15:04")
				rankStr := fmt.Sprintf("#%d", i+1)
				if i == 0 {
					rankStr = "🥇 #1"
				} else if i == 1 {
					rankStr = "🥈 #2"
				} else if i == 2 {
					rankStr = "🥉 #3"
				}

				fmt.Printf("%-5s %-12s %-8s %-8d %-8ds %s\n",
					accentStyle.Render(rankStr),
					entry.Mode,
					lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF87")).Bold(true).Render(fmt.Sprintf("%d", entry.Score)),
					entry.Length,
					entry.DurationSec,
					dimStyle.Render(dateStr),
				)
			}
		}

		_ = boxStyle
	},
}

func init() {
	statsCmd.Flags().BoolVar(&statsJSON, "json", false, "Output statistics in JSON format")
	statsCmd.Flags().BoolVar(&statsCSV, "csv", false, "Output leaderboard in CSV format")
	rootCmd.AddCommand(statsCmd)
}
