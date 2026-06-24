package ui

import (
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/caleb/terminal-kaleb/internal/game"
)

const (
	mainWidth    = 62
	blinkRate    = 600 * time.Millisecond
)

// ── messages ──────────────────────────────────────────────────────

type tickMsg time.Time
type blinkMsg struct{}

func tickCmd() tea.Cmd {
	return tea.Tick(blinkRate, func(t time.Time) tea.Msg { return blinkMsg{} })
}

// ── Model ──────────────────────────────────────────────────────────

type Model struct {
	state    game.State
	scenes   map[string]game.Scene
	styles   Styles
	blinkOn  bool
	width    int
	height   int

	// Transient system feedback shown after a choice
	feedback     []string
	feedbackType string
	showFeedback bool
}

func NewModel() Model {
	scenes := game.Scenes()
	state := game.NewState()
	m := Model{
		state:  state,
		scenes: scenes,
		styles: NewStyles(game.HeatCold),
		width:  mainWidth,
	}
	m.state.EnterScene(m.scenes)
	return m
}

func (m Model) Init() tea.Cmd {
	return tickCmd()
}

// ── Update ────────────────────────────────────────────────────────

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case blinkMsg:
		m.blinkOn = !m.blinkOn
		return m, tickCmd()

	case tea.KeyMsg:
		// Clear transient feedback on any key
		if m.showFeedback {
			m.showFeedback = false
			return m, nil
		}

		key := msg.String()

		// Dossier toggle
		if key == "d" || key == "D" {
			m.state.DossierOpen = !m.state.DossierOpen
			return m, nil
		}

		// Quit
		if key == "ctrl+c" || key == "q" {
			m.state.Quitting = true
			return m, tea.Quit
		}

		// Restart from outcome screen
		sc := m.currentScene()
		if sc.Outcome != nil || m.state.Heat >= game.HeatBurned || m.state.Time <= 0 {
			if key == "r" || key == "R" {
				return m.restart(), nil
			}
			return m, nil
		}

		// Choice selection: 1/2/3 or a/b/c
		idx := -1
		switch key {
		case "1", "a", "A":
			idx = 0
		case "2", "b", "B":
			idx = 1
		case "3", "c", "C":
			idx = 2
		}

		if idx >= 0 && idx < len(sc.Choices) {
			return m.applyChoice(idx)
		}
	}

	return m, nil
}

func (m Model) applyChoice(idx int) (Model, tea.Cmd) {
	sc := m.currentScene()
	c := sc.Choices[idx]

	// Build system feedback
	fb := m.buildFeedback(c)
	m.state.ApplyChoice(c)
	m.styles = NewStyles(m.state.Heat)
	m.state.EnterScene(m.scenes)

	m.feedback = fb.lines
	m.feedbackType = fb.msgType
	m.showFeedback = len(fb.lines) > 0

	return m, nil
}

type feedbackResult struct {
	lines   []string
	msgType string
}

func (m Model) buildFeedback(c game.Choice) feedbackResult {
	var lines []string
	msgType := "ok"

	// Simulate system response based on choice type and next scene
	switch c.Next {
	case "vpn_exploit":
		lines = []string{
			"> connect vpn-client --exploit CVE-2023-4966",
			"",
			"Injecting payload...",
			".......",
			"ACCESS GRANTED",
			fmt.Sprintf("Session established: %s", "10.42.0.1"),
		}
	case "fast_search":
		lines = []string{
			"> find / -name 'clinical' 2>/dev/null",
			"Scanning shares...",
			"......",
			"SIEM ALERT: anomalous session escalated",
			"IR engineer notified.",
		}
		msgType = "error"
	case "brute_force":
		lines = []string{
			"> brute --wordlist employees.json --target auth.meridian.net",
			"Attempt 1/428... failed",
			"Attempt 2/428... failed",
			"......",
			"Attempt 8/428...",
			"ACCOUNT LOCKED",
			"TRACE INITIATED",
		}
		msgType = "denied"
	case "phish_success":
		lines = []string{
			"> send-phish --target dana.reyes@meridian.net --template it-reset",
			"",
			"Sending...",
			"Delivered.",
			"Awaiting credential capture...",
		}
	case "find_mirror":
		lines = []string{
			"> find /mnt/shares -name '*mirror*' 2>/dev/null",
			"Searching...",
			".........",
			"/archive/legacy_backup/mirror_2019.tar.gz",
			"",
			"> decrypt mirror_2019.tar.gz",
			"Key found in filename.",
			"Decrypting... COMPLETE",
		}
	case "priv_esc":
		lines = []string{
			"> sudo -l && exploit --cve CVE-2023-0386",
			"",
			"Checking sudoers...",
			"Attempting kernel exploit...",
			"ERROR: patch detected (kernel 6.1.72)",
			"Escalation failed — logged to syslog.",
		}
		msgType = "warn"
	}

	// Always append heat cost as trace detection if nonzero
	if c.HeatCost > 0 && len(lines) > 0 {
		pct := int(m.state.Heat)*20 + c.HeatCost*20
		if pct > 100 {
			pct = 100
		}
		lines = append(lines, "")
		lines = append(lines, fmt.Sprintf("TRACE DETECTION: %d%%", pct))
	}

	return feedbackResult{lines, msgType}
}

func (m Model) restart() Model {
	state := game.NewState()
	nm := Model{
		state:  state,
		scenes: m.scenes,
		styles: NewStyles(game.HeatCold),
		width:  m.width,
		height: m.height,
	}
	nm.state.EnterScene(nm.scenes)
	return nm
}

func (m Model) currentScene() game.Scene {
	if sc, ok := m.scenes[m.state.Scene]; ok {
		return sc
	}
	return game.Scene{}
}

// ── View ──────────────────────────────────────────────────────────

func (m Model) View() string {
	if m.state.Quitting {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#3a6a3a")).Render("\n  session terminated.\n")
	}
	if m.state.DossierOpen {
		return m.dossierView()
	}
	return m.mainView()
}

func (m Model) mainView() string {
	w := m.width - 4
	if w < 40 {
		w = 40
	}

	var b strings.Builder

	// HUD
	b.WriteString(m.hudView(w))
	b.WriteString("\n\n")

	// System feedback overlay (transient)
	if m.showFeedback {
		b.WriteString(SystemMessage(m.feedback, m.feedbackType))
		b.WriteString("\n\n")
		b.WriteString(m.styles.HudHint.Render("  press any key to continue"))
		return m.styles.Shell.Width(w).Render(b.String())
	}

	sc := m.currentScene()

	// Terminal caught / time out
	if m.state.Heat >= game.HeatBurned {
		b.WriteString(m.caughtView())
		return m.styles.Shell.Width(w).Render(b.String())
	}
	if m.state.Time <= 0 {
		b.WriteString(m.timeoutView())
		return m.styles.Shell.Width(w).Render(b.String())
	}

	// Outcome screen
	if sc.Outcome != nil {
		b.WriteString(m.outcomeView(sc.Outcome))
		return m.styles.Shell.Width(w).Render(b.String())
	}

	// Normal scene
	b.WriteString(m.sceneView(sc, w))

	return m.styles.Shell.Width(w).Render(b.String())
}

func (m Model) hudView(w int) string {
	// Detection pips
	heat := HeatBar(m.state.Heat)
	heatLabel := m.styles.HudLabel.Render("detection ")

	// Time bar
	timeLabel := m.styles.HudLabel.Render("window ")
	timeBar := TimeBar(m.state.Time, 20)

	// Archetype tag
	arch := m.state.Archetype
	if arch == "" {
		arch = "—"
	}
	atag := m.styles.ATag.Render(" " + arch + " ")

	// Dossier hint
	hint := m.styles.HudHint.Render("[D] dossier  [Q] quit")

	// Trace detection line
	trace := TraceDetection(m.state.Heat)

	hudTop := lipgloss.JoinHorizontal(
		lipgloss.Center,
		heatLabel+heat+"  ",
		timeLabel+timeBar+"  ",
		atag+"  ",
		hint,
	)

	div := lipgloss.NewStyle().Foreground(lipgloss.Color("#1a3a1a")).Render(strings.Repeat("─", w-2))
	return hudTop + "\n" + trace + "\n" + div
}

func (m Model) sceneView(sc game.Scene, w int) string {
	var b strings.Builder

	// Prompt + blinking cursor
	cursor := "▋"
	if !m.blinkOn {
		cursor = " "
	}
	b.WriteString(m.styles.Prompt.Render(sc.Prompt + " " + cursor))
	b.WriteString("\n\n")

	// Scene text
	b.WriteString(m.styles.SceneText.Width(w - 2).Render(sc.Text))
	b.WriteString("\n")

	// Mystery annotation
	if sc.Mystery != "" {
		b.WriteString("\n")
		b.WriteString(m.styles.Mystery.Render(sc.Mystery))
		b.WriteString("\n")
	}

	// Log lines
	if len(m.state.Log) > 0 {
		b.WriteString("\n")
		div := lipgloss.NewStyle().Foreground(lipgloss.Color("#1a3a1a")).Render(strings.Repeat("─", w-2))
		b.WriteString(div + "\n")
		for _, l := range m.state.Log {
			sty := m.styles.LogNormal
			switch l.Level {
			case "warn":
				sty = m.styles.LogWarn
			case "bad":
				sty = m.styles.LogBad
			}
			b.WriteString(sty.Render("> " + l.Text))
			b.WriteString("\n")
		}
	}

	// Choices
	b.WriteString("\n")
	for _, c := range sc.Choices {
		sty := m.styles.ChoiceSafe
		switch c.Type {
		case game.ChoiceCaution:
			sty = m.styles.ChoiceCaut
		case game.ChoiceDanger:
			sty = m.styles.ChoiceDang
		}
		b.WriteString(sty.Width(w - 4).Render(c.Label))
		b.WriteString("\n")
	}

	return b.String()
}

func (m Model) outcomeView(oc *game.Outcome) string {
	sty := m.styles.OutSuccess
	switch oc.Type {
	case "fail":
		sty = m.styles.OutFail
	case "gameover":
		sty = m.styles.OutGameOver
	}
	restart := m.styles.ChoiceSafe.Render("[R] restart   [Q] quit")
	return sty.Render(oc.Text) + "\n\n" + restart
}

func (m Model) caughtView() string {
	lines := []string{
		"DETECTION THRESHOLD EXCEEDED",
		"",
		"Your session has been traced and terminated.",
		"IR team has your exit node IP and session log.",
		"",
		"You have hours before this becomes a legal problem.",
		"",
		"Play cleaner.",
	}
	txt := CorruptText(strings.Join(lines, "\n"))
	restart := m.styles.ChoiceSafe.Render("[R] restart   [Q] quit")
	return txt + "\n\n" + restart
}

func (m Model) timeoutView() string {
	txt := m.styles.OutFail.Render(
		"TIME EXPIRED\n\nIT rotation at 07:00. Routine check flagged the active session.\nYou killed the connection — but not before they got a fingerprint.",
	)
	restart := m.styles.ChoiceSafe.Render("[R] restart   [Q] quit")
	return txt + "\n\n" + restart
}

// ── Dossier view ──────────────────────────────────────────────────

func (m Model) dossierView() string {
	w := m.width - 4
	if w < 40 {
		w = 40
	}
	var b strings.Builder

	title := m.styles.DossierTitle.Render("// kael voss — field dossier")
	closeHint := m.styles.HudHint.Render("[D] close")
	div := lipgloss.NewStyle().Foreground(lipgloss.Color("#1a3a1a")).Render(strings.Repeat("─", w-2))

	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Left, title+"  ", closeHint))
	b.WriteString("\n" + div + "\n\n")

	// Target section
	b.WriteString(m.styles.DossierHead.Render("TARGET") + "\n")
	if len(m.state.DossierTarget) == 0 {
		b.WriteString(m.styles.DossierDetail.Render("  nothing logged yet") + "\n")
	} else {
		for _, e := range m.state.DossierTarget {
			b.WriteString(m.dossierEntry(e))
		}
	}

	b.WriteString("\n" + div + "\n\n")

	// Contacts section
	b.WriteString(m.styles.DossierHead.Render("CONTACTS") + "\n")
	if len(m.state.DossierContacts) == 0 {
		b.WriteString(m.styles.DossierDetail.Render("  no contacts identified") + "\n")
	} else {
		for _, e := range m.state.DossierContacts {
			b.WriteString(m.dossierEntry(e))
		}
	}

	b.WriteString("\n" + div + "\n\n")

	// Situation
	b.WriteString(m.styles.DossierHead.Render("SITUATION") + "\n")
	for _, sl := range m.state.SituationLines {
		sty := m.styles.SitOk
		switch sl.Class {
		case "warn":
			sty = m.styles.SitWarn
		case "bad":
			sty = m.styles.SitBad
		}
		b.WriteString(sty.Width(w-4).Render(sl.Text) + "\n")
	}

	return m.styles.Dossier.Width(w).Render(b.String())
}

func (m Model) dossierEntry(e game.DossierEntry) string {
	nameSty := m.styles.DossierName
	detSty := m.styles.DossierDetail
	switch e.Class {
	case "hi":
		nameSty = m.styles.DossierHiName
		detSty = m.styles.DossierHiDet
	case "al":
		nameSty = m.styles.DossierAlName
		detSty = m.styles.DossierAlDet
	case "danger":
		nameSty = m.styles.DossierDgName
		detSty = m.styles.DossierDgDet
	}
	return "  " + nameSty.Render(e.Name) + "\n" +
		"  " + detSty.Render(e.Detail) + "\n\n"
}
