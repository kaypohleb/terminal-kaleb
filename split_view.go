package main

import (
	"strings"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type paneFocus int

const (
	focusMain paneFocus = iota
	focusSidebar
)

const (
	splitSidebarWide   = 32
	splitSidebarNarrow = 12
	splitHelpLines     = 2 // blank + help
)

func (m model) splitSidebarWidth() int {
	if m.splitFocus == focusSidebar {
		return min(splitSidebarWide, max(24, m.width/3))
	}
	return splitSidebarNarrow
}

func (m model) splitPaneHeight() int {
	return max(1, m.height-splitHelpLines)
}

func (m model) splitMainPaneStyle(w int) lipgloss.Style {
	borderColor := lipgloss.Color("240")
	if m.splitFocus == focusMain {
		borderColor = lipgloss.Color("#33ff33")
	}
	return lipgloss.NewStyle().
		Width(w).
		Height(m.splitPaneHeight()).
		Padding(0, 1).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(borderColor)
}

func (m model) splitSidebarPaneStyle(w int) lipgloss.Style {
	borderColor := lipgloss.Color("240")
	if m.splitFocus == focusSidebar {
		borderColor = lipgloss.Color("#ffb300")
	}
	return lipgloss.NewStyle().
		Width(w).
		Height(m.splitPaneHeight()).
		Padding(0, 0, 0, 1).
		BorderStyle(lipgloss.NormalBorder()).
		BorderLeft(true).
		BorderForeground(borderColor)
}

func (m *model) initSplitViewports() {
	lorem := strings.TrimSpace(loremIpsumText())

	if !m.splitVPReady {
		m.splitMainVP = viewport.New()
		m.splitMainVP.SetContent("MAIN\n\n" + lorem)
		m.splitMainVP.SoftWrap = true
		m.splitMainVP.FillHeight = true

		m.splitSidebarVP = viewport.New()
		m.splitSidebarVP.SetContent("SIDEBAR\n\n" + lorem)
		m.splitSidebarVP.SoftWrap = true
		m.splitSidebarVP.FillHeight = true

		m.splitVPReady = true
	}

	m.layoutSplitViewports()
}

func (m *model) layoutSplitViewports() {
	if !m.splitVPReady {
		return
	}

	sidebarW := m.splitSidebarWidth()
	mainW := max(20, m.width-sidebarW-1)
	paneH := m.splitPaneHeight()

	m.splitMainVP.SetWidth(mainW)
	m.splitMainVP.SetHeight(paneH)
	m.splitMainVP.Style = m.splitMainPaneStyle(mainW)

	if sidebarW > splitSidebarNarrow {
		m.splitSidebarVP.SetWidth(sidebarW)
		m.splitSidebarVP.SetHeight(paneH)
		m.splitSidebarVP.Style = m.splitSidebarPaneStyle(sidebarW)
	}
}

func (m model) renderSplitView() string {
	if m.width <= 0 {
		return ""
	}

	sidebarW := m.splitSidebarWidth()
	mainW := max(20, m.width-sidebarW-1)

	main := m.renderSplitMain(mainW)
	sidebar := m.renderSplitSidebar(sidebarW)

	sep := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("│")
	body := lipgloss.JoinHorizontal(lipgloss.Top, main, sep, sidebar)

	help := m.quitStyle.Render("↑/↓ scroll focused pane · ←/→ switch focus · esc — menu · q — quit")
	return lipgloss.JoinVertical(lipgloss.Left, body, "", help)
}

func (m model) renderSplitMain(w int) string {
	if !m.splitVPReady {
		return lipgloss.NewStyle().Width(w).Height(m.splitPaneHeight()).Render(m.itemStyle.Render("loading…"))
	}
	return m.splitMainVP.View()
}

func (m model) renderSplitSidebar(w int) string {
	if w <= splitSidebarNarrow+2 {
		borderColor := lipgloss.Color("240")
		if m.splitFocus == focusSidebar {
			borderColor = lipgloss.Color("#ffb300")
		}
		return lipgloss.NewStyle().
			Width(w).
			Height(m.splitPaneHeight()).
			BorderStyle(lipgloss.NormalBorder()).
			BorderLeft(true).
			BorderForeground(borderColor).
			Render(lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render("›"))
	}
	if !m.splitVPReady {
		return lipgloss.NewStyle().Width(w).Height(m.splitPaneHeight()).Render(m.itemStyle.Render("loading…"))
	}
	return m.splitSidebarVP.View()
}

func (m model) updateSplitViewKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.phase = phaseHome
		return m, nil
	case "q", "ctrl+c":
		return m, tea.Quit
	case "left", "h":
		if m.splitFocus != focusMain {
			m.splitFocus = focusMain
			m.layoutSplitViewports()
		}
		return m, nil
	case "right", "l":
		if m.splitFocus != focusSidebar {
			m.splitFocus = focusSidebar
			m.layoutSplitViewports()
		}
		return m, nil
	}

	var cmd tea.Cmd
	switch m.splitFocus {
	case focusMain:
		m.splitMainVP, cmd = m.splitMainVP.Update(msg)
	case focusSidebar:
		if m.splitSidebarWidth() > splitSidebarNarrow {
			m.splitSidebarVP, cmd = m.splitSidebarVP.Update(msg)
		}
	}
	return m, cmd
}
