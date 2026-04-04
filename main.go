package main

// SSH server that runs a Bubble Tea TUI per session (Wish bubbletea middleware).

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"charm.land/log/v2"
	"charm.land/wish/v2"
	"charm.land/wish/v2/activeterm"
	"charm.land/wish/v2/bubbletea"
	"charm.land/wish/v2/logging"
	"github.com/charmbracelet/colorprofile"
	"github.com/charmbracelet/ssh"
)

const (
	defaultHost = "localhost"
	defaultPort = "23234"

	spinnerInterval = 80 * time.Millisecond
	loadMinDuration = 1_800 * time.Millisecond
)

type phase int

const (
	phaseLoading phase = iota
	phaseHome
	phaseDashboard
	phaseSettings
	phaseLogs
	phaseAbout
)

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

type tickMsg time.Time

type loadCompleteMsg struct{}

func main() {
	host := getenv("SSH_HOST", defaultHost)
	port := getenv("SSH_PORT", defaultPort)

	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("could not create server", "error", err)
		os.Exit(1)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("starting SSH server", "host", host, "port", port)
	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("server error", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("shutdown error", "error", err)
		os.Exit(1)
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, _ := s.Pty()
	m := model{
		term:      pty.Term,
		width:     pty.Window.Width,
		height:    pty.Window.Height,
		txtStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("10")),
		quitStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		itemStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("7")),
		selStyle: lipgloss.NewStyle().
			Background(lipgloss.Color("#ffffff")).
			Foreground(lipgloss.Color("#000000")).
			Bold(true),
		probeLabel:  lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Bold(true),
		swatchRed:   lipgloss.NewStyle().Foreground(lipgloss.Color("#ff5555")),
		swatchGreen: lipgloss.NewStyle().Foreground(lipgloss.Color("#50fa7b")),
		swatchBlue:  lipgloss.NewStyle().Foreground(lipgloss.Color("#6272ff")),
		spinnerStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("213")).Bold(true),
		pendingStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		bg:           "light",
		phase:        phaseLoading,
		choices: []string{
			"Dashboard",
			"Settings",
			"Logs",
			"About",
		},
	}
	// SSH sessions often inherit the server env; Bubble Tea may otherwise detect
	// a no-color profile and strip ANSI. True color matches modern clients.
	opts := []tea.ProgramOption{
		tea.WithColorProfile(colorprofile.TrueColor),
	}
	return m, opts
}

type model struct {
	term         string
	profile      string
	width        int
	height       int
	bg           string
	choices      []string
	cursor       int
	phase        phase
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
}

func tickSpinner() tea.Cmd {
	return tea.Tick(spinnerInterval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func finishLoading() tea.Cmd {
	return tea.Tick(loadMinDuration, func(t time.Time) tea.Msg {
		return loadCompleteMsg{}
	})
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
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.phase != phaseHome {
				break
			}
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.phase != phaseHome {
				break
			}
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			if m.phase != phaseHome {
				break
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
			}
		case "esc":
			if m.phase == phaseDashboard || m.phase == phaseSettings || m.phase == phaseLogs || m.phase == phaseAbout {
				m.phase = phaseHome
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func ansiColorProbe() string {
	const reset = "\x1b[0m"
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
		frame := spinnerFrames[m.spinnerFrame%len(spinnerFrames)]
		b.WriteString(m.spinnerStyle.Render(frame) + " ")
		b.WriteString(m.txtStyle.Render("Loading") + "\n\n")
		b.WriteString(m.pendingStyle.Render("Preparing your session…"))
	case phaseHome:
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
	case phaseDashboard:
		b.WriteString(m.txtStyle.Render("Dashboard") + "\n\n")
		b.WriteString(m.itemStyle.Render("Overview and status go here.\n\n"))
		b.WriteString(m.quitStyle.Render("esc — menu · q — quit\n"))
	case phaseSettings:
		b.WriteString(m.txtStyle.Render("Settings") + "\n\n")
		b.WriteString(m.itemStyle.Render("Configuration goes here.\n\n"))
		b.WriteString(m.quitStyle.Render("esc — menu · q — quit\n"))
	case phaseLogs:
		b.WriteString(m.txtStyle.Render("Logs") + "\n\n")
		b.WriteString(m.itemStyle.Render("Log output goes here.\n\n"))
		b.WriteString(m.quitStyle.Render("esc — menu · q — quit\n"))
	case phaseAbout:
		b.WriteString(m.txtStyle.Render("About") + "\n\n")
		b.WriteString(m.itemStyle.Render("SSH + Bubble Tea + Wish.\n\n"))
		b.WriteString(m.quitStyle.Render("esc — menu · q — quit\n"))
	}

	v := tea.NewView(b.String())
	v.AltScreen = true
	return v
}
