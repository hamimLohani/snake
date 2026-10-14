package storage

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

// KeyConfig defines rebindable keyboard mappings.
type KeyConfig struct {
	Up      []string `toml:"up"`
	Down    []string `toml:"down"`
	Left    []string `toml:"left"`
	Right   []string `toml:"right"`
	Pause   []string `toml:"pause"`
	Boss    []string `toml:"boss"`
	Restart []string `toml:"restart"`
	Quit    []string `toml:"quit"`
}

// Config represents user-facing settings.
type Config struct {
	Theme          string    `toml:"theme"`
	GridWidth      int       `toml:"grid_width"`
	GridHeight     int       `toml:"grid_height"`
	BaseSpeedMs    int       `toml:"base_speed_ms"`
	AdaptiveSpeed  bool      `toml:"adaptive_speed"`
	ScreenReader   bool      `toml:"screen_reader"`
	DefaultBossKey string    `toml:"default_boss_screen"` // "htop", "git", "build"
	Keys           KeyConfig `toml:"keys"`
}

// DefaultConfig returns clean defaults.
func DefaultConfig() Config {
	return Config{
		Theme:          "dark",
		GridWidth:      24,
		GridHeight:     20,
		BaseSpeedMs:    110,
		AdaptiveSpeed:  true,
		ScreenReader:   false,
		DefaultBossKey: "htop",
		Keys: KeyConfig{
			Up:      []string{"up", "w", "k"},
			Down:    []string{"down", "s", "j"},
			Left:    []string{"left", "a", "h"},
			Right:   []string{"right", "d", "l"},
			Pause:   []string{"p", " "},
			Boss:    []string{"b", "esc"},
			Restart: []string{"r"},
			Quit:    []string{"q", "ctrl+c"},
		},
	}
}

// ConfigDir returns ~/.config/snake.
func ConfigDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, ".config", "snake")
}

// ConfigPath returns the path to config.toml.
func ConfigPath() string {
	return filepath.Join(ConfigDir(), "config.toml")
}

// LoadConfig loads the config from ~/.config/snake/config.toml, creating it if absent.
func LoadConfig() (Config, error) {
	cfg := DefaultConfig()
	dir := ConfigDir()
	path := ConfigPath()

	if err := os.MkdirAll(dir, 0755); err != nil {
		return cfg, err
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		// Save default configuration with comments
		if saveErr := SaveConfig(cfg); saveErr != nil {
			return cfg, saveErr
		}
		return cfg, nil
	} else if err != nil {
		return cfg, err
	}

	if err := toml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

// SaveConfig persists the config to ~/.config/snake/config.toml.
func SaveConfig(cfg Config) error {
	dir := ConfigDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	var buf bytes.Buffer
	enc := toml.NewEncoder(&buf)
	if err := enc.Encode(cfg); err != nil {
		return err
	}

	return os.WriteFile(ConfigPath(), buf.Bytes(), 0644)
}
