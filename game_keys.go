package main

import (
	"fmt"
	"strings"
	"unicode/utf8"

	tea "charm.land/bubbletea/v2"
)

func (m model) updateGameKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.phase = phaseHome
		return m, nil
	case "q", "ctrl+c":
		return m, tea.Quit
	case "up":
		if m.gameChoiceIdx > 0 {
			m.gameChoiceIdx--
		}
	case "down":
		if m.gameChoiceIdx < len(m.gameScene.Choices)-1 {
			m.gameChoiceIdx++
		}
	case "enter":
		if len(m.gameScene.Choices) == 0 {
			return m, nil
		}
		i := clamp(m.gameChoiceIdx, 0, len(m.gameScene.Choices)-1)
		c := m.gameScene.Choices[i]
		cmd := fmt.Sprintf("> choose %s — %s", c.Key, c.Text)
		m.gameScene.CLILog = strings.TrimSpace(m.gameScene.CLILog + "\n" + cmd + "\n  [state] branch pending (demo)")
		m.gamePrompt = ""
	case "backspace":
		if len(m.gamePrompt) > 0 {
			_, sz := utf8.DecodeLastRuneInString(m.gamePrompt)
			if sz > 0 {
				m.gamePrompt = m.gamePrompt[:len(m.gamePrompt)-sz]
			}
		}
	default:
		if kp, ok := msg.(tea.KeyPressMsg); ok {
			t := kp.Key().Text
			if len(t) == 1 {
				switch strings.ToLower(t) {
				case "a", "b", "c", "d":
					for i, ch := range m.gameScene.Choices {
						if strings.EqualFold(ch.Key, t) {
							m.gameChoiceIdx = i
							break
						}
					}
					return m, nil
				}
			}
			if t != "" && len(m.gamePrompt) < 120 {
				m.gamePrompt += t
			}
		}
	}
	return m, cursorBlinkCmd()
}
