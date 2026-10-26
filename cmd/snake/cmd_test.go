package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestCmd_Version(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"--version"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Failed to execute --version: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "snake version") {
		t.Errorf("Expected 'snake version', got: %s", out)
	}
}

func TestCmd_PlayHelpShorthands(t *testing.T) {
	out := playCmd.UsageString()

	// Assert all shorthand flags are present in help
	if !strings.Contains(out, "-m, --mode") {
		t.Errorf("Expected -m, --mode in play help: %s", out)
	}
	if !strings.Contains(out, "-t, --theme") {
		t.Errorf("Expected -t, --theme in play help: %s", out)
	}
	if !strings.Contains(out, "-d, --difficulty") {
		t.Errorf("Expected -d, --difficulty in play help: %s", out)
	}
	if !strings.Contains(out, "-s, --size") {
		t.Errorf("Expected -s, --size in play help: %s", out)
	}
}

func TestCmd_BotHeadlessExecution(t *testing.T) {
	buf := new(bytes.Buffer)
	botCmd.SetOut(buf)
	botCmd.SetArgs([]string{"--headless", "--ticks", "50", "--seed", "12345"})

	if err := botCmd.Execute(); err != nil {
		t.Fatalf("Failed to execute bot --headless: %v", err)
	}
}

func TestCmd_StatsExports(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "snake_cmd_stats_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	// Test stats command execution
	statsCmd.SetArgs([]string{"--json"})
	if err := statsCmd.Execute(); err != nil {
		t.Fatalf("Failed to execute stats --json: %v", err)
	}

	statsCmd.SetArgs([]string{"--csv"})
	if err := statsCmd.Execute(); err != nil {
		t.Fatalf("Failed to execute stats --csv: %v", err)
	}
}

func TestCmd_ConfigView(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "snake_cmd_cfg_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	configCmd.SetArgs([]string{})
	if err := configCmd.Execute(); err != nil {
		t.Fatalf("Failed to execute config: %v", err)
	}
}
