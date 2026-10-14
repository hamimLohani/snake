package main

import (
	"fmt"
	"os"
	"os/exec"

	"snake/internal/storage"

	"github.com/charmbracelet/lipgloss"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

var (
	configEdit  bool
	configReset bool
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "View or edit snake configuration",
	Long:  `View current configuration, key bindings, or open ~/.config/snake/config.toml in your editor.`,
	Run: func(cmd *cobra.Command, args []string) {
		path := storage.ConfigPath()

		if configReset {
			cfg := storage.DefaultConfig()
			if err := storage.SaveConfig(cfg); err != nil {
				fmt.Fprintf(os.Stderr, "Error resetting config: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Successfully reset config to defaults at %s\n", path)
			return
		}

		if configEdit {
			editor := os.Getenv("EDITOR")
			if editor == "" {
				editor = "nano"
			}
			c := exec.Command(editor, path)
			c.Stdin = os.Stdin
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			if err := c.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Error running editor %s: %v\n", editor, err)
				os.Exit(1)
			}
			return
		}

		// View current config
		cfg, err := storage.LoadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}

		titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF87")).Bold(true)
		dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#7f8c8d"))
		keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00E5FF")).Bold(true)

		fmt.Println(titleStyle.Render("=== Snake CLI Configuration ==="))
		fmt.Printf("%s: %s\n\n", dimStyle.Render("Config File"), path)

		data, _ := toml.Marshal(cfg)
		fmt.Println(string(data))

		fmt.Println(dimStyle.Render("To edit, run: ") + keyStyle.Render("snake config --edit"))
	},
}

func init() {
	configCmd.Flags().BoolVarP(&configEdit, "edit", "e", false, "Open configuration file in default editor")
	configCmd.Flags().BoolVar(&configReset, "reset", false, "Reset configuration to defaults")
	rootCmd.AddCommand(configCmd)
}
