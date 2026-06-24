package main

import (
	"time"

	tea "charm.land/bubbletea/v2"
)

type tickMsg time.Time

type loadCompleteMsg struct{}

type shakeTickMsg struct{}

type cursorBlinkTickMsg struct{}

func tickSpinner() tea.Cmd {
	return tea.Tick(spinnerInterval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func finishLoading() tea.Cmd {
	return tea.Tick(loadMinDuration, func(time.Time) tea.Msg {
		return loadCompleteMsg{}
	})
}

func shakeTick() tea.Cmd {
	return tea.Tick(shakeTickInterval, func(time.Time) tea.Msg {
		return shakeTickMsg{}
	})
}

const cursorBlinkInterval = 450 * time.Millisecond

func cursorBlinkCmd() tea.Cmd {
	return tea.Tick(cursorBlinkInterval, func(time.Time) tea.Msg {
		return cursorBlinkTickMsg{}
	})
}

