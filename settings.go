package main

import (
	tea "charm.land/bubbletea/v2"
)

func (m model) updateSettingsKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "left", "h":
		if len(availableFonts) > 0 {
			m.fontIdx = (m.fontIdx - 1 + len(availableFonts)) % len(availableFonts)
			return m, setFontCmd(availableFonts[m.fontIdx])
		}
	case "right", "l", "enter":
		if len(availableFonts) > 0 {
			m.fontIdx = (m.fontIdx + 1) % len(availableFonts)
			return m, setFontCmd(availableFonts[m.fontIdx])
		}
	case "esc":
		m.phase = phaseHome
		return m, nil
	case "q", "ctrl+c":
		return m, tea.Quit
	}
	return m, nil
}

func setFontCmd(font string) tea.Cmd {
	// OSC 50: Set font. Supported varies by terminal; safe to ignore.
	// Terminated with BEL.
	return tea.Printf("\x1b]50;%s\x07", font)
}

