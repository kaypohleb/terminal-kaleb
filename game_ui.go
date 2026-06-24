package main

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

// Agents.md §2: main viewport layout — header, narrative log, choices, CLI, dossier.

const (
	colorOpsGreen  = "#33ff33"
	colorCaution   = "#ffb300"
	colorCritical  = "#ff3333"
	colorObsCyan   = "#00ffff"
	colorObsAccent = "#6272ff"
)

type gameChoice struct {
	Key  string
	Text string
	Hint string
}

type gameScene struct {
	Narrative   string
	Observation string
	Choices     []gameChoice
	Intel       []string
	Contacts    []string
	Situation   []string
	CLILog      string
}

func demoGameScene() gameScene {
	return gameScene{
		Narrative: strings.TrimSpace(`
Meridian Dynamics edge node is live. Bounty wire shows $35,000 — deadline 07:00.
You are one hop from the ARGO trial partition. Logs show biometric re-auth every ninety seconds.`),
		Observation: "They already know someone is probing; the question is whether they can prove it's you.",
		Choices: []gameChoice{
			{Key: "A", Text: "phish an employee", Hint: "spoofed IT reset. slow, clean. [detection +2]"},
			{Key: "B", Text: "lateral move via stale VPN cert", Hint: "reuse contractor session. [time −15m]"},
			{Key: "C", Text: "hold and watch heartbeat timing", Hint: "map the guard cycle. [detection +0]"},
		},
		Intel: []string{
			"Meridian Dynamics — clinical trials, ARGO compound.",
			"Leak ref: MER-INT-09 (partial): trial anomalies flagged, not disclosed.",
		},
		Contacts: []string{
			"Marcus Webb — procurement — trust: low",
			"Sandra Okoye — Chief Clinical Strategy — trust: hostile",
		},
		Situation: []string{
			"Internal: Don't torch the cover on a cheap file.",
			"Heat: elevated — assume mailbox monitoring.",
		},
		CLILog: "> nmap -sV edge.meridian.internal\n  open 443/tcp https\n  [detection +1]",
	}
}

func (m model) renderGameUI() string {
	if m.width <= 0 {
		return ""
	}

	dossierW := min(32, max(24, m.width/3))
	mainW := max(20, m.width-dossierW-1)

	header := m.renderGameHeader(mainW + dossierW + 1)
	mainCol := m.renderGameMainColumn(mainW)
	dossierCol := m.renderGameDossier(dossierW)

	body := lipgloss.JoinHorizontal(lipgloss.Top, mainCol, lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("│"), dossierCol)

	return lipgloss.JoinVertical(lipgloss.Left, header, body)
}

func (m model) renderGameHeader(totalW int) string {
	// Detection: 5 segments (higher fill = more exposed / less stealth).
	segs := m.gameDetectionSegs
	if segs < 0 {
		segs = 0
	}
	if segs > 5 {
		segs = 5
	}
	filled := strings.Repeat("█", segs)
	empty := strings.Repeat("░", 5-segs)
	detStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(colorOpsGreen))
	if segs >= 3 {
		detStyle = detStyle.Foreground(lipgloss.Color(colorCaution))
	}
	if segs >= 4 {
		detStyle = detStyle.Foreground(lipgloss.Color(colorCritical))
	}
	det := detStyle.Render(filled + empty)

	tmax := max(1, m.gameTimeMaxMin)
	tcur := clamp(m.gameTimeCurMin, 0, tmax)
	pct := float64(tcur) / float64(tmax)
	barW := max(8, min(28, totalW/3))
	filledBar := int(float64(barW) * pct)
	if filledBar > barW {
		filledBar = barW
	}
	bar := "[" + strings.Repeat("─", filledBar) + strings.Repeat("·", barW-filledBar) + "]"
	timeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(colorOpsGreen))
	if pct < 0.35 {
		timeStyle = timeStyle.Foreground(lipgloss.Color(colorCaution))
	}
	if pct < 0.15 {
		timeStyle = timeStyle.Foreground(lipgloss.Color(colorCritical))
	}

	left := lipgloss.NewStyle().Foreground(lipgloss.Color(colorOpsGreen)).
		Render("detection ") + det
	mid := lipgloss.NewStyle().Foreground(lipgloss.Color("7")).
		Render("  time remaining ") + timeStyle.Render(fmt.Sprintf("%s %dm/%dm", bar, tcur, tmax))

	pill := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#000000")).
		Background(lipgloss.Color("#33ff33")).
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#33ff33")).
		Render(m.gameAlias)

	line := left + mid
	// Push pill to the right within totalW (approximate; lipgloss width aware).
	pad := max(0, totalW-lipgloss.Width(line)-lipgloss.Width(pill))
	line = line + strings.Repeat(" ", pad) + pill

	return lipgloss.NewStyle().
		Width(totalW).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 0, 1, 0).
		Render(line)
}

func (m model) renderGameMainColumn(w int) string {
	body := lipgloss.NewStyle().
		Width(w).
		Foreground(lipgloss.Color(colorOpsGreen)).
		Render(wrapHard(m.gameScene.Narrative, w))

	obs := m.renderSystemObservation(w)

	choices := m.renderChoiceMenu(w)

	cursor := "_"
	if !m.cursorBlink {
		cursor = " "
	}
	prompt := lipgloss.NewStyle().Foreground(lipgloss.Color(colorOpsGreen)).
		Render("root@ghost:~$ ") +
		m.itemStyle.Render(m.gamePrompt) +
		lipgloss.NewStyle().Foreground(lipgloss.Color(colorOpsGreen)).Bold(true).Render(cursor)

	cliLog := ""
	if strings.TrimSpace(m.gameScene.CLILog) != "" {
		cliLog = "\n" + lipgloss.NewStyle().Width(w).Foreground(lipgloss.Color("8")).Render(m.gameScene.CLILog)
	}

	block := lipgloss.JoinVertical(lipgloss.Left, body, obs, choices, "\n"+prompt+cliLog)
	return lipgloss.NewStyle().Width(w).Render(block)
}

func (m model) renderSystemObservation(w int) string {
	if strings.TrimSpace(m.gameScene.Observation) == "" {
		return ""
	}
	line := lipgloss.NewStyle().Foreground(lipgloss.Color(colorObsCyan)).Render("// " + m.gameScene.Observation)
	accent := lipgloss.NewStyle().
		Background(lipgloss.Color(colorObsAccent)).
		Width(1).
		Render(" ")
	return "\n" + lipgloss.JoinHorizontal(lipgloss.Top, accent, " ", lipgloss.NewStyle().Width(max(1, w-3)).Render(line))
}

func (m model) renderChoiceMenu(w int) string {
	var lines []string
	for i, c := range m.gameScene.Choices {
		box := lipgloss.NewStyle().
			Width(w).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1)
		if i == m.gameChoiceIdx {
			box = box.BorderForeground(lipgloss.Color(colorOpsGreen)).
				Background(lipgloss.Color("#0d1f0d"))
		}
		lines = append(lines, box.Render(
			lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(colorOpsGreen)).Render("["+c.Key+"] ") +
				lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render(c.Text+" — ") +
				lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Render(c.Hint),
		))
	}
	return "\n" + strings.Join(lines, "\n")
}

func (m model) renderGameDossier(w int) string {
	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(colorCaution)).Render("DOSSIER")
	sec := func(name string, lines []string) string {
		h := lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Bold(true).Render(name)
		body := lipgloss.NewStyle().Width(w - 2).Foreground(lipgloss.Color("7")).Render(strings.Join(lines, "\n"))
		return lipgloss.JoinVertical(lipgloss.Left, h, body)
	}
	block := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		sec("Intel", m.gameScene.Intel),
		"",
		sec("Contacts", m.gameScene.Contacts),
		"",
		sec("Situation", m.gameScene.Situation),
	)
	return lipgloss.NewStyle().
		Width(w).
		Padding(0, 0, 0, 1).
		BorderStyle(lipgloss.NormalBorder()).
		BorderLeft(true).
		BorderForeground(lipgloss.Color("240")).
		Render(block)
}

func wrapHard(s string, w int) string {
	if w <= 0 {
		return s
	}
	var b strings.Builder
	for _, para := range strings.Split(s, "\n") {
		para = strings.TrimSpace(para)
		if para == "" {
			b.WriteByte('\n')
			continue
		}
		for len(para) > w {
			b.WriteString(para[:w])
			b.WriteByte('\n')
			para = strings.TrimSpace(para[w:])
		}
		b.WriteString(para)
		b.WriteByte('\n')
	}
	return strings.TrimRight(b.String(), "\n")
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
