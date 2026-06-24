package main

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m model) homeContent() string {
	var b strings.Builder
	b.WriteString(m.probeLabel.Render("Color probe (raw ANSI)") + "\n")
	b.WriteString(ansiColorProbe() + "\n")
	b.WriteString(
		m.swatchRed.Render("lipgloss") + " " +
			m.swatchGreen.Render("RGB") + " " +
			m.swatchBlue.Render("swatches") + "\n",
	)
	b.WriteString(
		"selection sample: " + m.selStyle.Render(" white bg / black fg ") + "\n\n",
	)
	b.WriteString(m.txtStyle.Render("Choose an option (↑ / ↓, enter to open)") + "\n\n")
	for i, label := range m.choices {
		if i == m.cursor {
			b.WriteString(m.selStyle.Render("› "+label) + "\n")
		} else {
			b.WriteString(m.itemStyle.Render("  "+label) + "\n")
		}
	}
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf(
		"Term %s · %dx%d · %s · %s\n\n",
		m.term, m.width, m.height, m.bg, m.profile,
	))
	b.WriteString(m.quitStyle.Render("Press 'q' to quit\n"))
	return b.String()
}

func (m model) messageModalView() string {
	msg := strings.Join([]string{
		"msg_001  03:12",
		"",
		"Package confirmed. Window open.",
		"",
		"Don't leave traces.",
	}, "\n")
	inner := m.itemStyle.Render(msg) + "\n\n" + m.quitStyle.Render("esc / enter — close")
	box := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#ffffff")).
		Padding(1, 2).
		Render(inner)
	if m.width <= 0 || m.height <= 0 {
		return box
	}
	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		box,
		lipgloss.WithWhitespaceStyle(lipgloss.NewStyle()),
	)
}

func ansiColorProbe() string {
	basic := "\x1b[31m16-color\x1b[0m"
	xterm256 := "\x1b[38;5;214m256-color\x1b[0m"
	direct := "\x1b[38;2;255;121;198mtruecolor\x1b[0m"
	reverse := "\x1b[7mreverse\x1b[0m"
	return basic + "  " + xterm256 + "  " + direct + "  " + reverse
}

func (m model) View() tea.View {
	var b strings.Builder

	switch m.phase {
	case phaseLoading:
		var loading strings.Builder
		frame := spinnerFrames[m.spinnerFrame%len(spinnerFrames)]
		loading.WriteString(m.itemStyle.Render(loadingBanner) + "\n\n")
		loading.WriteString(m.spinnerStyle.Render(frame) + " ")
		loading.WriteString(m.txtStyle.Render("Loading") + "\n\n")
		loading.WriteString(m.pendingStyle.Render("Preparing your session…"))
		if m.width > 0 && m.height > 0 {
			b.WriteString(lipgloss.Place(
				m.width, m.height,
				lipgloss.Center, lipgloss.Center,
				loading.String(),
				lipgloss.WithWhitespaceStyle(lipgloss.NewStyle()),
			))
		} else {
			b.WriteString(loading.String())
		}

	case phaseHome:
		home := m.homeContent()
		off := 0
		if m.shakeActive && m.shakeFrame < len(shakeOffsets) {
			off = shakeOffsets[m.shakeFrame]
		}
		if off > 0 {
			home = lipgloss.NewStyle().MarginLeft(off).Render(home)
		}
		b.WriteString(home)

	case phaseDashboard:
		b.WriteString(m.txtStyle.Render("Dashboard") + "\n\n")
		b.WriteString(m.itemStyle.Render("Overview and status go here.\n\n"))
		b.WriteString(m.quitStyle.Render("esc — menu · q — quit\n"))

	case phaseSettings:
		b.WriteString(m.txtStyle.Render("Settings") + "\n\n")
		current := "n/a"
		if len(availableFonts) > 0 {
			current = availableFonts[m.fontIdx%len(availableFonts)]
		}
		b.WriteString(m.itemStyle.Render("Font: "+current) + "\n")
		b.WriteString(m.pendingStyle.Render("←/→ (or h/l) to change · enter cycles") + "\n\n")
		b.WriteString(m.pendingStyle.Render("Note: font is controlled by your terminal. This tries OSC 50; unsupported terminals will ignore it.") + "\n\n")
		b.WriteString(m.quitStyle.Render("esc — menu · q — quit\n"))

	case phaseLogs:
		b.WriteString(m.txtStyle.Render("Logs") + "\n\n")
		b.WriteString(m.itemStyle.Render("Log output goes here.\n\n"))
		b.WriteString(m.quitStyle.Render("esc — menu · q — quit\n"))

	case phaseAbout:
		b.WriteString(m.txtStyle.Render("About") + "\n\n")
		b.WriteString(m.itemStyle.Render("SSH + Bubble Tea + Wish.\n\n"))
		b.WriteString(m.quitStyle.Render("esc — menu · q — quit\n"))

	case phaseMessageModal:
		b.WriteString(m.messageModalView())

	case phaseLoremViewport:
		title := m.txtStyle.Render("Lorem Ipsum")
		help := m.quitStyle.Render("↑/↓ scroll · pgup/pgdn · home/end · esc — menu · q — quit")
		body := m.loremVP.View()
		b.WriteString(title + "\n")
		b.WriteString(body + "\n")
		b.WriteString(help + "\n")

	case phaseGame:
		b.WriteString(m.renderGameUI())
		b.WriteString("\n")
		b.WriteString(m.quitStyle.Render("↑/↓ · enter to confirm · a–d jump · type at prompt · esc — menu"))

	case phaseSplitView:
		b.WriteString(m.renderSplitView())
	}

	v := tea.NewView(b.String())
	v.AltScreen = true
	return v
}

