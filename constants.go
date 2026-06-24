package main

import "time"

const (
	defaultHost = "localhost"
	defaultPort = "23234"

	spinnerInterval   = 80 * time.Millisecond
	loadMinDuration   = 1_800 * time.Millisecond
	shakeFrameCount   = 20
	shakeTickInterval = 45 * time.Millisecond
)

// Horizontal jitter per frame (spaces via lipgloss margin) for a short "screen shake".
var shakeOffsets = []int{0, 2, 1, 3, 0, 2, 4, 1, 0, 3, 2, 5, 1, 0, 2, 3, 1, 4, 0, 2}

// Font selection is ultimately controlled by the user's terminal emulator, but
// some terminals support OSC 50 to request a font change.
var availableFonts = []string{"Fira Code", "Cascadia Mono", "Consolas", "Courier New"}

const loadingBanner = `╺┳╸┏━╸┏━┓┏┳┓╻┏┓╻┏━┓╻     ┏┓ ┏━┓┏━╸┏━┓┏━╸╻ ╻
 ┃ ┣╸ ┣┳┛┃┃┃┃┃┗┫┣━┫┃     ┣┻┓┣┳┛┣╸ ┣━┫┃  ┣━┫
 ╹ ┗━╸╹┗╸╹ ╹╹╹ ╹╹ ╹┗━╸   ┗━┛╹┗╸┗━╸╹ ╹┗━╸╹ ╹`

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

