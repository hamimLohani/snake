package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

//go:embed man_embed.txt
var manPageText string

var manCmd = &cobra.Command{
	Use:   "man",
	Short: "View the manual page for snake",
	Long:  `Display the formatted UNIX manual page for the snake CLI.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Try invoking 'man snake' if installed
		if c := exec.Command("man", "snake"); c.Run() == nil {
			return
		}

		// 2. Try looking in local relative directory
		execPath, err := os.Executable()
		var candidatePaths []string
		if err == nil {
			base := filepath.Dir(execPath)
			candidatePaths = append(candidatePaths,
				filepath.Join(base, "..", "man", "snake.6"),
				filepath.Join(base, "man", "snake.6"),
				"man/snake.6",
			)
		} else {
			candidatePaths = append(candidatePaths, "man/snake.6")
		}

		for _, p := range candidatePaths {
			if _, statErr := os.Stat(p); statErr == nil {
				c := exec.Command("man", p)
				c.Stdin = os.Stdin
				c.Stdout = os.Stdout
				c.Stderr = os.Stderr
				if c.Run() == nil {
					return
				}
			}
		}

		// 3. Fallback: pipe embedded man text to pager or stdout
		pager := os.Getenv("PAGER")
		if pager == "" {
			pager = "less"
		}

		pagerCmd := exec.Command(pager)
		pagerCmd.Stdin = os.Stdin
		pagerCmd.Stdout = os.Stdout
		pagerCmd.Stderr = os.Stderr

		stdin, pipeErr := pagerCmd.StdinPipe()
		if pipeErr == nil && pagerCmd.Start() == nil {
			_, _ = stdin.Write([]byte(manPageText))
			_ = stdin.Close()
			_ = pagerCmd.Wait()
			return
		}

		fmt.Print(manPageText)
	},
}

func init() {
	rootCmd.AddCommand(manCmd)
}
