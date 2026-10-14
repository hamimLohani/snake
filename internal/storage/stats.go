package storage

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

// ScoreEntry records an individual game result.
type ScoreEntry struct {
	Date        time.Time `json:"date"`
	Mode        string    `json:"mode"`
	Score       int       `json:"score"`
	Length      int       `json:"length"`
	Seed        int64     `json:"seed"`
	DurationSec int       `json:"duration_sec"`
}

// Stats holds cumulative player data, high scores, and leaderboards.
type Stats struct {
	TotalGames       int            `json:"total_games"`
	TotalPlayTimeSec int            `json:"total_play_time_sec"`
	TotalFoodEaten   int            `json:"total_food_eaten"`
	HighScores       map[string]int `json:"high_scores"` // mode -> high score
	Leaderboard      []ScoreEntry   `json:"leaderboard"`
	Achievements     map[string]bool `json:"achievements"`
}

// NewStats initializes an empty stats object.
func NewStats() Stats {
	return Stats{
		HighScores:   make(map[string]int),
		Leaderboard:  make([]ScoreEntry, 0),
		Achievements: make(map[string]bool),
	}
}

// StatsPath returns ~/.config/snake/stats.json.
func StatsPath() string {
	return filepath.Join(ConfigDir(), "stats.json")
}

// LoadStats loads player metrics and scores from disk.
func LoadStats() (Stats, error) {
	st := NewStats()
	path := StatsPath()

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return st, nil
	} else if err != nil {
		return st, err
	}

	if err := json.Unmarshal(data, &st); err != nil {
		return st, err
	}

	if st.HighScores == nil {
		st.HighScores = make(map[string]int)
	}
	if st.Achievements == nil {
		st.Achievements = make(map[string]bool)
	}
	return st, nil
}

// SaveStats writes stats to disk atomically.
func SaveStats(st Stats) error {
	dir := ConfigDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return err
	}

	tmpFile := filepath.Join(dir, fmt.Sprintf("stats.tmp.%d", time.Now().UnixNano()))
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		return err
	}

	return os.Rename(tmpFile, StatsPath())
}

// RecordGame updates stats with a completed game.
func (st *Stats) RecordGame(entry ScoreEntry) (isNewHigh bool, unlockedAchievements []string) {
	st.TotalGames++
	st.TotalPlayTimeSec += entry.DurationSec
	st.TotalFoodEaten += entry.Length - 3 // Approximate food eaten

	currentHigh := st.HighScores[entry.Mode]
	if entry.Score > currentHigh {
		st.HighScores[entry.Mode] = entry.Score
		isNewHigh = true
	}

	st.Leaderboard = append(st.Leaderboard, entry)

	// Keep top 50 scores sorted descending by score
	sort.Slice(st.Leaderboard, func(i, j int) bool {
		return st.Leaderboard[i].Score > st.Leaderboard[j].Score
	})
	if len(st.Leaderboard) > 50 {
		st.Leaderboard = st.Leaderboard[:50]
	}

	// Check achievements
	unlocked := st.checkAchievements(entry)
	return isNewHigh, unlocked
}

func (st *Stats) checkAchievements(entry ScoreEntry) []string {
	var newlyUnlocked []string
	check := func(key, name string, condition bool) {
		if condition && !st.Achievements[key] {
			st.Achievements[key] = true
			newlyUnlocked = append(newlyUnlocked, name)
		}
	}

	check("first_bite", "First Bite: Eat your first food", entry.Score >= 10)
	check("century", "Centurion: Score over 100 points", entry.Score >= 100)
	check("grandmaster", "Grandmaster: Score over 300 points", entry.Score >= 300)
	check("long_boi", "Long Boi: Reach length of 25", entry.Length >= 25)
	check("coffee_break", "Coffee Break: Play for over 5 minutes total", st.TotalPlayTimeSec >= 300)

	return newlyUnlocked
}

// ExportCSV outputs leaderboard as CSV string.
func (st *Stats) ExportCSV() (string, error) {
	var buf os.File
	_ = buf
	out := [][]string{
		{"Rank", "Date", "Mode", "Score", "Length", "DurationSec", "Seed"},
	}

	for i, entry := range st.Leaderboard {
		out = append(out, []string{
			strconv.Itoa(i + 1),
			entry.Date.Format(time.RFC3339),
			entry.Mode,
			strconv.Itoa(entry.Score),
			strconv.Itoa(entry.Length),
			strconv.Itoa(entry.DurationSec),
			strconv.FormatInt(entry.Seed, 10),
		})
	}

	var res string
	for _, row := range out {
		res += fmt.Sprintf("%s\n", row[0]+","+row[1]+","+row[2]+","+row[3]+","+row[4]+","+row[5]+","+row[6])
	}
	_ = csv.NewWriter
	return res, nil
}

// ExportJSON outputs stats as pretty JSON.
func (st *Stats) ExportJSON() (string, error) {
	data, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
