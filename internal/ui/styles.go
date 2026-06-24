package ui

import (
	"github.com/caleb/terminal-kaleb/internal/game"
	"github.com/charmbracelet/lipgloss"
)

// Palette — terminal greens shifting to red-amber under heat
var (
	// Base greens
	colGreenDim    = lipgloss.Color("#2a5a2a")
	colGreenMid    = lipgloss.Color("#3a9a3a")
	colGreenBright = lipgloss.Color("#7afa7a")
	colTextMain    = lipgloss.Color("#b8e8b8")
	colBorder      = lipgloss.Color("#1a3a1a")
	colBg          = lipgloss.Color("#0d0d0d")
	colBgDark      = lipgloss.Color("#080808")
	colBlue        = lipgloss.Color("#7aaafa")
	colBlueDim     = lipgloss.Color("#1a3a7a")

	// Heat ramp
	colWarm1  = lipgloss.Color("#c8c860")
	colWarm2  = lipgloss.Color("#d4a060")
	colHot1   = lipgloss.Color("#e07050")
	colHot2   = lipgloss.Color("#cc2a0a")
	colRed    = lipgloss.Color("#fa7a7a")
	colRedBrt = lipgloss.Color("#ff2020")

	// Semantic
	colAmber    = lipgloss.Color("#facc7a")
	colAmberDim = lipgloss.Color("#ba7517")
	colTeal     = lipgloss.Color("#5DCAA5")
	colFail     = lipgloss.Color("#f09595")

	// Dossier
	colDossierName    = lipgloss.Color("#7aba7a")
	colDossierDetail  = lipgloss.Color("#4a7a4a")
	colDossierHiName  = lipgloss.Color("#7aaafa")
	colDossierHiDet   = lipgloss.Color("#3a5a8a")
	colDossierAlName  = lipgloss.Color("#facc7a")
	colDossierAlDet   = lipgloss.Color("#7a6030")
	colDossierDgName  = lipgloss.Color("#fa7a7a")
	colDossierDgDet   = lipgloss.Color("#7a3030")
	colSitOk          = lipgloss.Color("#3a9a5a")
	colSitWarn        = lipgloss.Color("#9a7a20")
	colSitBad         = lipgloss.Color("#9a3a3a")
)

// Styles holds all computed lipgloss styles
type Styles struct {
	Heat game.HeatLevel

	Shell       lipgloss.Style
	Prompt      lipgloss.Style
	SceneText   lipgloss.Style
	Mystery     lipgloss.Style
	ChoiceSafe  lipgloss.Style
	ChoiceCaut  lipgloss.Style
	ChoiceDang  lipgloss.Style
	HudLabel    lipgloss.Style
	HudHint     lipgloss.Style
	ATag        lipgloss.Style
	LogNormal   lipgloss.Style
	LogWarn     lipgloss.Style
	LogBad      lipgloss.Style
	OutSuccess  lipgloss.Style
	OutFail     lipgloss.Style
	OutGameOver lipgloss.Style
	Blink       lipgloss.Style

	// Dossier
	Dossier       lipgloss.Style
	DossierHead   lipgloss.Style
	DossierTitle  lipgloss.Style
	DossierName   lipgloss.Style
	DossierDetail lipgloss.Style
	DossierHiName lipgloss.Style
	DossierHiDet  lipgloss.Style
	DossierAlName lipgloss.Style
	DossierAlDet  lipgloss.Style
	DossierDgName lipgloss.Style
	DossierDgDet  lipgloss.Style
	SitOk         lipgloss.Style
	SitWarn       lipgloss.Style
	SitBad        lipgloss.Style
	SitBorder     lipgloss.Style
}

func NewStyles(heat game.HeatLevel) Styles {
	s := Styles{Heat: heat}

	// Shell border reacts to heat
	borderCol := colBorder
	switch heat {
	case game.HeatWarm:
		borderCol = lipgloss.Color("#3a3a1a")
	case game.HeatHot:
		borderCol = lipgloss.Color("#6a3a0a")
	case game.HeatCritical, game.HeatBurned:
		borderCol = lipgloss.Color("#8a2a0a")
	}

	s.Shell = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderCol).
		Padding(0)

	// Prompt colour escalates
	promptCol := colGreenMid
	switch heat {
	case game.HeatWarm:
		promptCol = lipgloss.Color("#7a7a1a")
	case game.HeatHot:
		promptCol = lipgloss.Color("#9a4a1a")
	case game.HeatCritical, game.HeatBurned:
		promptCol = colHot2
	}
	s.Prompt = lipgloss.NewStyle().Foreground(promptCol)

	// Scene text colour shifts warm under heat
	textCol := colTextMain
	switch heat {
	case game.HeatWarm:
		textCol = colWarm1
	case game.HeatHot:
		textCol = colWarm2
	case game.HeatCritical, game.HeatBurned:
		textCol = colHot1
	}
	s.SceneText = lipgloss.NewStyle().Foreground(textCol)

	// Mystery annotation
	mystCol := colBlue
	mystBorder := colBlueDim
	if heat >= game.HeatHot {
		mystCol = colAmber
		mystBorder = lipgloss.Color("#5a3a0a")
	}
	if heat >= game.HeatCritical {
		mystCol = lipgloss.Color("#fa6a2a")
		mystBorder = lipgloss.Color("#6a2a00")
	}
	s.Mystery = lipgloss.NewStyle().
		Foreground(mystCol).
		BorderLeft(true).
		BorderStyle(lipgloss.Border{Left: "│"}).
		BorderForeground(mystBorder).
		PaddingLeft(1).
		Italic(true)

	// Choices
	choiceBase := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1)

	safeCol := colGreenBright
	safeBord := colGreenDim
	if heat >= game.HeatHot {
		safeCol = colAmber
		safeBord = lipgloss.Color("#3a2a0a")
	}
	s.ChoiceSafe = choiceBase.Copy().
		Foreground(safeCol).
		BorderForeground(safeBord)

	s.ChoiceCaut = choiceBase.Copy().
		Foreground(colAmber).
		BorderForeground(lipgloss.Color("#5a4a1a"))

	dangCol := colRed
	if heat >= game.HeatCritical {
		dangCol = colRedBrt
	}
	s.ChoiceDang = choiceBase.Copy().
		Foreground(dangCol).
		BorderForeground(lipgloss.Color("#5a2a2a"))

	// HUD
	hudCol := colGreenMid
	if heat >= game.HeatHot {
		hudCol = lipgloss.Color("#7a3a1a")
	}
	s.HudLabel = lipgloss.NewStyle().Foreground(hudCol)

	hintCol := colGreenDim
	if heat >= game.HeatHot {
		hintCol = lipgloss.Color("#5a2a1a")
	}
	if heat >= game.HeatCritical {
		hintCol = lipgloss.Color("#7a2a1a")
	}
	s.HudHint = lipgloss.NewStyle().Foreground(hintCol)

	atagBg := lipgloss.Color("#0a1a0a")
	atagFg := lipgloss.Color("#3a9a3a")
	atagBord := lipgloss.Color("#2a4a2a")
	if heat >= game.HeatHot {
		atagBg = lipgloss.Color("#1a0a0a")
		atagFg = lipgloss.Color("#9a3a3a")
		atagBord = lipgloss.Color("#4a1a1a")
	}
	s.ATag = lipgloss.NewStyle().
		Background(atagBg).
		Foreground(atagFg).
		Border(lipgloss.NormalBorder()).
		BorderForeground(atagBord).
		Padding(0, 1)

	// Log lines
	s.LogNormal = lipgloss.NewStyle().Foreground(lipgloss.Color("#3a6a3a"))
	s.LogWarn = lipgloss.NewStyle().Foreground(colAmberDim)
	s.LogBad = lipgloss.NewStyle().Foreground(lipgloss.Color("#cc2a0a"))

	// Outcomes
	s.OutSuccess = lipgloss.NewStyle().
		Foreground(colTeal).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#1d9e75")).
		Background(lipgloss.Color("#04140e")).
		Padding(0, 1)

	s.OutFail = lipgloss.NewStyle().
		Foreground(colFail).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#a32d2d")).
		Background(lipgloss.Color("#14040a")).
		Padding(0, 1)

	s.OutGameOver = lipgloss.NewStyle().
		Foreground(colRedBrt).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#cc0000")).
		Background(lipgloss.Color("#1a0000")).
		Padding(0, 1)

	s.Blink = lipgloss.NewStyle().Foreground(promptCol)

	// ── Dossier ──────────────────────────────────────────────────
	s.Dossier = lipgloss.NewStyle().
		Background(colBgDark).
		Padding(1)

	s.DossierHead = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#2a6a2a")).
		Bold(true)

	s.DossierTitle = lipgloss.NewStyle().
		Foreground(colGreenMid)

	s.DossierName = lipgloss.NewStyle().Foreground(colDossierName).Bold(true)
	s.DossierDetail = lipgloss.NewStyle().Foreground(colDossierDetail)
	s.DossierHiName = lipgloss.NewStyle().Foreground(colDossierHiName).Bold(true)
	s.DossierHiDet = lipgloss.NewStyle().Foreground(colDossierHiDet)
	s.DossierAlName = lipgloss.NewStyle().Foreground(colDossierAlName).Bold(true)
	s.DossierAlDet = lipgloss.NewStyle().Foreground(colDossierAlDet)
	s.DossierDgName = lipgloss.NewStyle().Foreground(colDossierDgName).Bold(true)
	s.DossierDgDet = lipgloss.NewStyle().Foreground(colDossierDgDet)

	s.SitOk = lipgloss.NewStyle().
		Foreground(colSitOk).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#1a4a2a")).
		Padding(0, 1)

	s.SitWarn = lipgloss.NewStyle().
		Foreground(colSitWarn).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#4a3a0a")).
		Padding(0, 1)

	s.SitBad = lipgloss.NewStyle().
		Foreground(colSitBad).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#4a1a1a")).
		Padding(0, 1)

	return s
}
