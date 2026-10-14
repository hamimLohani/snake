package renderer

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// RenderBossScreen returns a realistic fake terminal screen (htop, build log, git log).
func RenderBossScreen(screenType string, width, height int) string {
	switch screenType {
	case "build":
		return renderBuildScreen(width, height)
	case "git":
		return renderGitScreen(width, height)
	case "htop":
		fallthrough
	default:
		return renderHtopScreen(width, height)
	}
}

func renderHtopScreen(width, height int) string {
	var b strings.Builder

	green := lipgloss.NewStyle().Foreground(lipgloss.Color("#2ecc71"))
	blue := lipgloss.NewStyle().Foreground(lipgloss.Color("#3498db"))
	cyan := lipgloss.NewStyle().Foreground(lipgloss.Color("#1abc9c"))
	red := lipgloss.NewStyle().Foreground(lipgloss.Color("#e74c3c"))
	white := lipgloss.NewStyle().Foreground(lipgloss.Color("#ecf0f1"))
	dim := lipgloss.NewStyle().Foreground(lipgloss.Color("#7f8c8d"))
	barBg := lipgloss.NewStyle().Foreground(lipgloss.Color("#2c3e50"))

	makeBar := func(cpuNum int, pct float64, barWidth int) string {
		fill := int(float64(barWidth) * pct / 100.0)
		if fill > barWidth {
			fill = barWidth
		}
		empty := barWidth - fill
		barStr := green.Render(strings.Repeat("|", fill)) + barBg.Render(strings.Repeat(" ", empty))
		return fmt.Sprintf("%-2d[%s%5.1f%%] ", cpuNum, barStr, pct)
	}

	bw := 20
	if width < 80 {
		bw = 10
	}

	b.WriteString(makeBar(1, 23.4, bw) + makeBar(5, 38.1, bw) + "\n")
	b.WriteString(makeBar(2, 45.2, bw) + makeBar(6, 12.8, bw) + "\n")
	b.WriteString(makeBar(3, 19.8, bw) + makeBar(7, 62.4, bw) + "\n")
	b.WriteString(makeBar(4, 31.0, bw) + makeBar(8, 18.5, bw) + "\n")

	b.WriteString(cyan.Render("Mem ") + "[||||||||||||||||||             9.42G/32.0G]  ")
	b.WriteString(white.Render("Tasks: ") + green.Render("124") + ", " + white.Render("481 thr; 1 running\n"))

	b.WriteString(cyan.Render("Swp ") + "[                               0K/8.00G]     ")
	b.WriteString(white.Render("Load average: ") + green.Render("1.84 2.12 1.98\n"))

	nowStr := time.Now().Format("15:04:05")
	b.WriteString(white.Render("Uptime: ") + green.Render("6 days, 04:22:18  ") + dim.Render(nowStr) + "\n\n")

	// Table header
	headerStyle := lipgloss.NewStyle().Background(lipgloss.Color("#1abc9c")).Foreground(lipgloss.Color("#000000")).Bold(true)
	header := fmt.Sprintf("%-7s %-8s %-4s %-4s %-7s %-7s %-6s %-2s %-6s %-6s %s",
		"PID", "USER", "PRI", "NI", "VIRT", "RES", "SHR", "S", "CPU%", "MEM%", "Command")
	if len(header) > width && width > 0 {
		header = header[:width]
	}
	b.WriteString(headerStyle.Render(header) + "\n")

	processes := []struct {
		pid, user, pri, ni, virt, res, shr, s, cpu, mem, cmd string
	}{
		{"10482", "staff", "20", "0", "14.2G", "1.8G", "112M", "S", "14.2", "5.6", "node /usr/local/bin/tsc --watch"},
		{"8912", "staff", "20", "0", "4280M", "680M", "92M", "S", "9.8", "2.1", "go build ./cmd/..."},
		{"3419", "staff", "20", "0", "2140M", "340M", "64M", "S", "4.5", "1.1", "/usr/bin/rust-analyzer"},
		{"512", "root", "20", "0", "840M", "120M", "32M", "S", "2.1", "0.4", "dockerd --config-file /etc/docker/daemon.json"},
		{"1420", "staff", "20", "0", "1820M", "240M", "48M", "S", "1.8", "0.8", "nvim internal/engine/engine.go"},
		{"412", "staff", "20", "0", "410M", "45M", "18M", "S", "0.9", "0.2", "tmux new-session -s prod"},
		{"9981", "staff", "20", "0", "312M", "38M", "12M", "R", "0.7", "0.1", "htop"},
		{"1", "root", "20", "0", "180M", "12M", "8M", "S", "0.0", "0.0", "/sbin/launchd"},
	}

	for _, p := range processes {
		line := fmt.Sprintf("%-7s %-8s %-4s %-4s %-7s %-7s %-6s %-2s %-6s %-6s %s",
			p.pid, p.user, p.pri, p.ni, p.virt, p.res, p.shr, p.s,
			red.Render(p.cpu), blue.Render(p.mem), white.Render(p.cmd))
		b.WriteString(line + "\n")
	}

	b.WriteString("\n" + dim.Render("F1:Help  F2:Setup  F3:Search  F4:Filter  F5:Tree  F6:SortBy  F9:Kill  F10:Quit (Press Boss Key to resume)"))
	return b.String()
}

func renderBuildScreen(width, height int) string {
	var b strings.Builder
	green := lipgloss.NewStyle().Foreground(lipgloss.Color("#2ecc71")).Bold(true)
	cyan := lipgloss.NewStyle().Foreground(lipgloss.Color("#00E5FF"))
	dim := lipgloss.NewStyle().Foreground(lipgloss.Color("#7f8c8d"))
	yellow := lipgloss.NewStyle().Foreground(lipgloss.Color("#f1c40f"))

	b.WriteString(cyan.Render("➜  snake-core (main) ") + "go test -v ./... -race -cover\n")
	b.WriteString(dim.Render("=== RUN   TestEngineMovement\n"))
	b.WriteString(green.Render("--- PASS: TestEngineMovement (0.01s)\n"))
	b.WriteString(dim.Render("=== RUN   TestDeterminismReplay\n"))
	b.WriteString(green.Render("--- PASS: TestDeterminismReplay (0.03s)\n"))
	b.WriteString(dim.Render("=== RUN   TestPowerUpExpiration\n"))
	b.WriteString(green.Render("--- PASS: TestPowerUpExpiration (0.00s)\n"))
	b.WriteString(dim.Render("=== RUN   TestSurvivalShrinkBounds\n"))
	b.WriteString(green.Render("--- PASS: TestSurvivalShrinkBounds (0.02s)\n"))
	b.WriteString(yellow.Render("coverage: 92.4% of statements\n"))
	b.WriteString(green.Render("ok  \tgithub.com/engine/v2\t0.312s\n\n"))
	b.WriteString(cyan.Render("[100%] Built target libcore_sim.a\n"))
	b.WriteString(dim.Render("Waiting for file changes... (Press Boss Key to return)"))
	return b.String()
}

func renderGitScreen(width, height int) string {
	var b strings.Builder
	yellow := lipgloss.NewStyle().Foreground(lipgloss.Color("#f1c40f"))
	cyan := lipgloss.NewStyle().Foreground(lipgloss.Color("#00E5FF"))
	dim := lipgloss.NewStyle().Foreground(lipgloss.Color("#7f8c8d"))

	b.WriteString(yellow.Render("commit 8f7e6d1c9204bf3dbea2f6d9fbf69c20") + cyan.Render(" (HEAD -> main, origin/main)\n"))
	b.WriteString("Author: Senior Engineer <eng@corp.internal>\n")
	b.WriteString("Date:   " + time.Now().Format(time.RFC1123) + "\n\n")
	b.WriteString("    feat(engine): optimize tick deterministic state scheduler\n\n")
	b.WriteString("    - Add zero-alloc coordinate delta arithmetic\n")
	b.WriteString("    - Validate boundary clamping across virtual viewports\n")
	b.WriteString("    - Wire continuous CI regression matrix\n\n")
	b.WriteString(dim.Render(": (Press Boss Key to return)"))
	return b.String()
}
