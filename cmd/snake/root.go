package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "1.0.1"
	commit  = "none"
	date    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "snake",
	Short: "snake - A modern, developer-grade terminal Snake game",
	Long: `snake is a professional, production-quality terminal Snake game built for
developers taking a coffee break. Features classic and modern game modes,
power-ups, boss key, replays, seed sharing, and customizable themes.`,
	Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date),
	Run: func(cmd *cobra.Command, args []string) {
		// Default to running 'snake play' if no subcommand is given
		playCmd.Run(playCmd, args)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetVersionTemplate("snake version {{.Version}}\n")
}
