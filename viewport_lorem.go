package main

import (
	"strings"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m model) updateViewportKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.phase = phaseHome
		return m, nil
	case "q", "ctrl+c":
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.loremVP, cmd = m.loremVP.Update(msg)
	return m, cmd
}

func (m *model) initLoremViewport() {
	if m.loremVPReady {
		// Keep scroll position if user re-enters.
		return
	}
	m.loremVP = viewport.New(
		viewport.WithWidth(max(0, m.width-2)),
		viewport.WithHeight(max(0, m.height-4)),
	)
	m.loremVP.Style = lipgloss.NewStyle()
	m.loremVP.SetContent(strings.TrimSpace(loremIpsumText()))
	m.loremVPReady = true
}

func loremIpsumText() string {
	// Intentionally long to demonstrate scrolling.
	return strings.Join([]string{
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed non risus.",
		"Suspendisse lectus tortor, dignissim sit amet, adipiscing nec, ultricies sed, dolor.",
		"Cras elementum ultrices diam. Maecenas ligula massa, varius a, semper congue, euismod non, mi.",
		"Proin porttitor, orci nec nonummy molestie, enim est eleifend mi, non fermentum diam nisl sit amet erat.",
		"Duis semper. Duis arcu massa, scelerisque vitae, consequat in, pretium a, enim.",
		"Pellentesque congue. Ut in risus volutpat libero pharetra tempor. Cras vestibulum bibendum augue.",
		"Praesent egestas leo in pede. Praesent blandit odio eu enim. Pellentesque sed dui ut augue blandit sodales.",
		"Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Aliquam nibh.",
		"Mauris ac mauris sed pede pellentesque fermentum. Maecenas adipiscing ante non diam sodales hendrerit.",
		"",
		"Ut velit mauris, egestas sed, gravida nec, ornare ut, mi. Aenean ut orci vel massa suscipit pulvinar.",
		"Nulla sollicitudin. Fusce varius, ligula non tempus aliquam, nunc turpis ullamcorper nibh, in tempus sapien eros vitae ligula.",
		"Pellentesque rhoncus nunc et augue. Integer id felis. Curabitur aliquet pellentesque diam.",
		"Integer quis metus vitae elit lobortis egestas. Lorem ipsum dolor sit amet, consectetuer adipiscing elit.",
		"Morbi vel erat non mauris convallis vehicula. Nulla et sapien. Integer tortor tellus, aliquam faucibus, convallis id, congue eu, quam.",
		"Mauris ullamcorper felis vitae erat. Proin feugiat, augue non elementum posuere, metus purus iaculis lectus, et tristique ligula justo vitae magna.",
		"",
		"Aliquam convallis sollicitudin purus. Praesent aliquam, enim at fermentum mollis, ligula massa adipiscing nisl, ac euismod nibh nisl eu lectus.",
		"Fusce vulputate sem at sapien. Vivamus leo. Aliquam euismod libero eu enim.",
		"Nullam nec magna. Maecenas odio dolor, vulputate vel, auctor ac, accumsan id, felis.",
		"Donec neque neque, rutrum at, molestie at, tristique et, justo. Praesent mattis, massa quis luctus fermentum, turpis mi volutpat justo, eu volutpat enim diam eget metus.",
		"Maecenas ornare tortor. Donec sed tellus eget sapien fringilla nonummy.",
		"",
		// Repeat a few blocks for comfortable scrolling.
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed non risus.",
		"Suspendisse lectus tortor, dignissim sit amet, adipiscing nec, ultricies sed, dolor.",
		"Cras elementum ultrices diam. Maecenas ligula massa, varius a, semper congue, euismod non, mi.",
		"Proin porttitor, orci nec nonummy molestie, enim est eleifend mi, non fermentum diam nisl sit amet erat.",
		"Duis semper. Duis arcu massa, scelerisque vitae, consequat in, pretium a, enim.",
		"Pellentesque congue. Ut in risus volutpat libero pharetra tempor. Cras vestibulum bibendum augue.",
		"Praesent egestas leo in pede. Praesent blandit odio eu enim. Pellentesque sed dui ut augue blandit sodales.",
		"Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Aliquam nibh.",
		"Mauris ac mauris sed pede pellentesque fermentum. Maecenas adipiscing ante non diam sodales hendrerit.",
		"",
		"Ut velit mauris, egestas sed, gravida nec, ornare ut, mi. Aenean ut orci vel massa suscipit pulvinar.",
		"Nulla sollicitudin. Fusce varius, ligula non tempus aliquam, nunc turpis ullamcorper nibh, in tempus sapien eros vitae ligula.",
		"Pellentesque rhoncus nunc et augue. Integer id felis. Curabitur aliquet pellentesque diam.",
		"Integer quis metus vitae elit lobortis egestas. Lorem ipsum dolor sit amet, consectetuer adipiscing elit.",
		"Morbi vel erat non mauris convallis vehicula. Nulla et sapien. Integer tortor tellus, aliquam faucibus, convallis id, congue eu, quam.",
		"Mauris ullamcorper felis vitae erat. Proin feugiat, augue non elementum posuere, metus purus iaculis lectus, et tristique ligula justo vitae magna.",
	}, "\n")
}

