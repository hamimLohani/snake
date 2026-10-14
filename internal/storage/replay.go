package storage

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"snake/internal/engine"
)

// ReplayStep records an input direction event at a specific engine tick.
type ReplayStep struct {
	Tick      uint64           `json:"tick"`
	Direction engine.Direction `json:"direction"`
}

// ReplayFile encapsulates complete deterministic execution data.
type ReplayFile struct {
	Version       string        `json:"version"`
	Date          time.Time     `json:"date"`
	Mode          engine.Mode   `json:"mode"`
	Seed          int64         `json:"seed"`
	Width         int           `json:"width"`
	Height        int           `json:"height"`
	WithObstacles bool          `json:"with_obstacles"`
	FinalScore    int           `json:"final_score"`
	FinalLength   int           `json:"final_length"`
	TotalTicks    uint64        `json:"total_ticks"`
	Steps         []ReplayStep  `json:"steps"`
}

// ReplayDir returns ~/.config/snake/replays.
func ReplayDir() string {
	return filepath.Join(ConfigDir(), "replays")
}

// LatestReplayPath returns ~/.config/snake/replays/latest.snk.
func LatestReplayPath() string {
	return filepath.Join(ReplayDir(), "latest.snk")
}

// SaveReplay persists a replay file with gzip compression.
func SaveReplay(path string, rep ReplayFile) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.Marshal(rep)
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	gz := gzip.NewWriter(f)
	defer gz.Close()

	_, err = gz.Write(data)
	return err
}

// LoadReplay reads and decompresses a replay file.
func LoadReplay(path string) (ReplayFile, error) {
	var rep ReplayFile
	f, err := os.Open(path)
	if err != nil {
		return rep, fmt.Errorf("failed to open replay file: %w", err)
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		// Try fallback uncompressed JSON
		f2, err2 := os.Open(path)
		if err2 != nil {
			return rep, err
		}
		defer f2.Close()
		if jsonErr := json.NewDecoder(f2).Decode(&rep); jsonErr == nil {
			return rep, nil
		}
		return rep, fmt.Errorf("failed to decompress replay: %w", err)
	}
	defer gz.Close()

	if err := json.NewDecoder(gz).Decode(&rep); err != nil {
		return rep, fmt.Errorf("invalid replay data: %w", err)
	}

	return rep, nil
}
