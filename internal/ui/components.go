package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/caleb/terminal-kaleb/internal/game"
	"github.com/charmbracelet/lipgloss"
)

// HeatBar renders detection pips
func HeatBar(heat game.HeatLevel) string {
	pips := make([]string, 5)
	for i := 0; i < 5; i++ {
		if i < int(heat) {
			switch {
			case heat <= game.HeatWarm:
				pips[i] = lipgloss.NewStyle().Foreground(lipgloss.Color("#ef9f27")).Render("█")
			case heat == game.HeatHot:
				pips[i] = lipgloss.NewStyle().Foreground(lipgloss.Color("#e24b4b")).Render("█")
			default:
				pips[i] = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff2020")).Render("█")
			}
		} else {
			pips[i] = lipgloss.NewStyle().Foreground(lipgloss.Color("#1a2a1a")).Render("░")
		}
	}
	return strings.Join(pips, " ")
}

// TimeBar renders a progress bar for the time window
func TimeBar(pct int, width int) string {
	if width < 4 {
		width = 20
	}
	filled := (pct * width) / 100
	if filled < 0 {
		filled = 0
	}
	if filled > width {
		filled = width
	}

	barCol := lipgloss.Color("#1d9e75")
	if pct <= 50 {
		barCol = lipgloss.Color("#ef9f27")
	}
	if pct <= 25 {
		barCol = lipgloss.Color("#e24b4a")
	}

	bar := lipgloss.NewStyle().Foreground(barCol).Render(strings.Repeat("━", filled))
	empty := lipgloss.NewStyle().Foreground(lipgloss.Color("#1a2a1a")).Render(strings.Repeat("─", width-filled))
	return "[" + bar + empty + "]"
}

// TraceDetection renders a percentage line — heat-reactive
func TraceDetection(heat game.HeatLevel) string {
	pct := int(heat) * 20
	col := lipgloss.Color("#3a9a3a")
	switch {
	case pct >= 80:
		col = lipgloss.Color("#ff2020")
	case pct >= 60:
		col = lipgloss.Color("#e24b4b")
	case pct >= 40:
		col = lipgloss.Color("#ef9f27")
	case pct >= 20:
		col = lipgloss.Color("#c8c860")
	}
	label := lipgloss.NewStyle().Foreground(lipgloss.Color("#3a6a3a")).Render("TRACE DETECTION: ")
	val := lipgloss.NewStyle().Foreground(col).Bold(true).Render(fmt.Sprintf("%d%%", pct))
	return label + val
}

// ProgressBar renders a fake loading bar with a label
type ProgressState struct {
	Label    string
	Current  int
	Total    int
	Done     bool
	Failed   bool
	Width    int
}

func (p ProgressState) View() string {
	width := p.Width
	if width == 0 {
		width = 30
	}
	pct := 0
	if p.Total > 0 {
		pct = (p.Current * width) / p.Total
	}
	if p.Done {
		pct = width
	}

	col := lipgloss.Color("#1d9e75")
	char := "█"
	if p.Failed {
		col = lipgloss.Color("#cc2a0a")
		char = "▒"
	}

	bar := lipgloss.NewStyle().Foreground(col).Render(strings.Repeat(char, pct))
	empty := lipgloss.NewStyle().Foreground(lipgloss.Color("#1a2a1a")).Render(strings.Repeat("░", width-pct))
	suffix := ""
	if p.Done {
		suffix = lipgloss.NewStyle().Foreground(lipgloss.Color("#3a9a3a")).Render(" COMPLETE")
	} else if p.Failed {
		suffix = lipgloss.NewStyle().Foreground(lipgloss.Color("#cc2a0a")).Render(" FAILED")
	} else {
		suffix = fmt.Sprintf(" %d%%", (p.Current*100)/max(p.Total, 1))
	}

	label := lipgloss.NewStyle().Foreground(lipgloss.Color("#3a7a3a")).Render(p.Label)
	return label + "\n[" + bar + empty + "]" + suffix
}

// Countdown renders a lockdown timer
func Countdown(d time.Duration) string {
	if d <= 0 {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#ff2020")).Bold(true).Render("LOCKDOWN IN: 00:00")
	}
	mins := int(d.Minutes())
	secs := int(d.Seconds()) % 60
	col := lipgloss.Color("#ef9f27")
	if d < 30*time.Second {
		col = lipgloss.Color("#ff2020")
	}
	label := lipgloss.NewStyle().Foreground(lipgloss.Color("#7a3a1a")).Render("LOCKDOWN IN: ")
	val := lipgloss.NewStyle().Foreground(col).Bold(true).Render(fmt.Sprintf("%02d:%02d", mins, secs))
	return label + val
}

// SystemMessage renders a terminal system response block
func SystemMessage(lines []string, msgType string) string {
	col := lipgloss.Color("#3a9a3a")
	border := lipgloss.Color("#1a4a1a")
	switch msgType {
	case "warn":
		col = lipgloss.Color("#ef9f27")
		border = lipgloss.Color("#4a3a0a")
	case "error":
		col = lipgloss.Color("#cc2a0a")
		border = lipgloss.Color("#4a1a1a")
	case "denied":
		col = lipgloss.Color("#ff2020")
		border = lipgloss.Color("#6a0000")
	}
	styled := lipgloss.NewStyle().
		Foreground(col).
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(border).
		PaddingLeft(1)
	return styled.Render(strings.Join(lines, "\n"))
}

// CorruptText renders glitched text for crisis moments
func CorruptText(text string) string {
	// Intersperse zero-width or combining characters for visual corruption
	corrupt := lipgloss.Color("#cc2a0a")
	return lipgloss.NewStyle().Foreground(corrupt).Bold(true).Render(text)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
