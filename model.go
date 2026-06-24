package main

import (
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type model struct {
	term    string
	profile string
	width   int
	height  int
	bg      string

	choices []string
	cursor  int
	phase   phase

	spinnerFrame int

	txtStyle     lipgloss.Style
	quitStyle    lipgloss.Style
	itemStyle    lipgloss.Style
	selStyle     lipgloss.Style
	probeLabel   lipgloss.Style
	swatchRed    lipgloss.Style
	swatchGreen  lipgloss.Style
	swatchBlue   lipgloss.Style
	spinnerStyle lipgloss.Style
	pendingStyle lipgloss.Style

	shakeActive bool
	shakeFrame  int

	loremVP      viewport.Model
	loremVPReady bool

	fontIdx int

	// Agents.md game viewport (demo scene).
	gameScene          gameScene
	gameDetectionSegs  int
	gameTimeCurMin     int
	gameTimeMaxMin     int
	gameAlias          string
	gamePrompt         string
	gameChoiceIdx      int
	cursorBlink        bool

	splitFocus       paneFocus
	splitMainVP      viewport.Model
	splitSidebarVP   viewport.Model
	splitVPReady     bool
}

func newModel(term string, width, height int) model {
	return model{
		term:   term,
		width:  width,
		height: height,

		txtStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("10")),
		quitStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		itemStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("7")),
		selStyle: lipgloss.NewStyle().
			Background(lipgloss.Color("#ffffff")).
			Foreground(lipgloss.Color("#000000")).
			Bold(true),
		probeLabel:   lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Bold(true),
		swatchRed:    lipgloss.NewStyle().Foreground(lipgloss.Color("#ff5555")),
		swatchGreen:  lipgloss.NewStyle().Foreground(lipgloss.Color("#50fa7b")),
		swatchBlue:   lipgloss.NewStyle().Foreground(lipgloss.Color("#6272ff")),
		spinnerStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("213")).Bold(true),
		pendingStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("8")),

		bg:    "light",
		phase: phaseLoading,
		choices: []string{
			"Dashboard",
			"Settings",
			"Logs",
			"About",
			"Shake screen",
			"Secure message",
			"Lorem viewport",
			"Mission",
			"Split view",
		},
		fontIdx: 0,

		gameScene:         demoGameScene(),
		gameDetectionSegs: 2,
		gameTimeCurMin:    180,
		gameTimeMaxMin:    240,
		gameAlias:         "ghost",
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		tea.RequestBackgroundColor,
		tea.RequestCapability("RGB"),
		tea.RequestCapability("Tc"),
		tickSpinner(),
		finishLoading(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case loadCompleteMsg:
		if m.phase == phaseLoading {
			m.phase = phaseHome
			return m, nil
		}

	case cursorBlinkTickMsg:
		if m.phase == phaseGame {
			m.cursorBlink = !m.cursorBlink
			return m, cursorBlinkCmd()
		}

	case shakeTickMsg:
		if m.shakeActive {
			m.shakeFrame++
			if m.shakeFrame >= shakeFrameCount {
				m.shakeActive = false
				m.shakeFrame = 0
				return m, nil
			}
			return m, shakeTick()
		}

	case tickMsg:
		if m.phase == phaseLoading {
			m.spinnerFrame++
			if m.spinnerFrame >= len(spinnerFrames) {
				m.spinnerFrame = 0
			}
			return m, tickSpinner()
		}

	case tea.ColorProfileMsg:
		m.profile = msg.String()

	case tea.BackgroundColorMsg:
		if msg.IsDark() {
			m.bg = "dark"
		}

	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		if m.loremVPReady {
			m.loremVP.SetWidth(max(0, m.width-2))
			m.loremVP.SetHeight(max(0, m.height-4))
		}
		if m.splitVPReady {
			m.layoutSplitViewports()
		}

	case tea.KeyMsg:
		if m.phase == phaseGame {
			return m.updateGameKeys(msg)
		}

		if m.phase == phaseSplitView {
			return m.updateSplitViewKeys(msg)
		}

		// Settings mode consumes keys for font cycling.
		if m.phase == phaseSettings {
			return m.updateSettingsKeys(msg)
		}

		// Viewport mode consumes scrolling keys.
		if m.phase == phaseLoremViewport {
			return m.updateViewportKeys(msg)
		}

		// Default navigation.
		switch msg.String() {
		case "up":
			if m.phase == phaseHome && m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.phase == phaseHome && m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			return m.handleEnter()
		case "esc":
			return m.handleEsc()
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) handleEnter() (tea.Model, tea.Cmd) {
	if m.phase == phaseMessageModal {
		m.phase = phaseHome
		return m, nil
	}

	if m.phase != phaseHome {
		return m, nil
	}

	switch m.cursor {
	case 0:
		m.phase = phaseDashboard
	case 1:
		m.phase = phaseSettings
	case 2:
		m.phase = phaseLogs
	case 3:
		m.phase = phaseAbout
	case 4:
		m.shakeActive = true
		m.shakeFrame = 0
		return m, shakeTick()
	case 5:
		m.phase = phaseMessageModal
	case 6:
		m.phase = phaseLoremViewport
		m.initLoremViewport()
	case 7:
		m.phase = phaseGame
		m.cursorBlink = true
		return m, cursorBlinkCmd()
	case 8:
		m.phase = phaseSplitView
		m.splitFocus = focusMain
		m.initSplitViewports()
	}
	return m, nil
}

func (m model) handleEsc() (tea.Model, tea.Cmd) {
	switch m.phase {
	case phaseMessageModal:
		m.phase = phaseHome
	case phaseDashboard, phaseSettings, phaseLogs, phaseAbout, phaseGame, phaseSplitView:
		m.phase = phaseHome
	}
	return m, nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

