package game

// HeatLevel drives terminal colour escalation
type HeatLevel int

const (
	HeatCold    HeatLevel = 0
	HeatWarm    HeatLevel = 1
	HeatHot     HeatLevel = 2
	HeatCritical HeatLevel = 3
	HeatBurned  HeatLevel = 4
)

type ChoiceType string

const (
	ChoiceSafe    ChoiceType = "safe"
	ChoiceCaution ChoiceType = "caution"
	ChoiceDanger  ChoiceType = "danger"
)

type Choice struct {
	Label      string
	Next       string
	Type       ChoiceType
	HeatCost   int
	TimeCost   int
	Archetype  string
}

type Outcome struct {
	Type string // "success" | "fail" | "gameover"
	Text string
}

type DossierEntry struct {
	Name   string
	Detail string
	Class  string // "hi" | "al" | "danger" | ""
}

type StatusLine struct {
	Text  string
	Class string // "ok" | "warn" | "bad"
}

type Scene struct {
	Prompt     string
	Text       string
	Mystery    string
	Choices    []Choice
	Outcome    *Outcome

	// Dossier updates triggered on entering this scene
	TargetEntries   []DossierEntry
	ContactEntries  []DossierEntry
	SituationLines  []StatusLine
}

type State struct {
	Heat      HeatLevel
	Time      int // 0–100
	Scene     string
	Archetype string
	Log       []LogEntry

	// Dossier accumulates — entries are never removed
	DossierTarget   []DossierEntry
	DossierContacts []DossierEntry
	SituationLines  []StatusLine

	// Track seen scene keys to avoid re-adding dossier entries
	SeenScenes map[string]bool

	// View
	DossierOpen bool
	Quitting    bool
}

type LogEntry struct {
	Text  string
	Level string // "" | "warn" | "bad"
}

func NewState() State {
	return State{
		Heat:       HeatCold,
		Time:       100,
		Scene:      "intro",
		SeenScenes: make(map[string]bool),
	}
}

func (s *State) AddLog(text, level string) {
	s.Log = append(s.Log, LogEntry{text, level})
	if len(s.Log) > 6 {
		s.Log = s.Log[len(s.Log)-6:]
	}
}

func (s *State) ApplyChoice(c Choice) {
	if c.Archetype != "" && s.Archetype == "" {
		s.Archetype = c.Archetype
	}
	if c.HeatCost > 0 {
		newHeat := int(s.Heat) + c.HeatCost
		if newHeat > int(HeatBurned) {
			newHeat = int(HeatBurned)
		}
		s.Heat = HeatLevel(newHeat)
		level := "warn"
		if c.HeatCost >= 2 {
			level = "bad"
		}
		s.AddLog("detection +"+itoa(c.HeatCost), level)
	}
	if c.TimeCost > 0 {
		s.Time -= c.TimeCost
		if s.Time < 0 {
			s.Time = 0
		}
		s.AddLog("-"+itoa(c.TimeCost)+"min", "")
	}
	s.Scene = c.Next
}

func (s *State) EnterScene(scenes map[string]Scene) {
	sc, ok := scenes[s.Scene]
	if !ok || s.SeenScenes[s.Scene] {
		return
	}
	s.SeenScenes[s.Scene] = true

	for _, e := range sc.TargetEntries {
		if !s.hasDossierEntry("target", e.Name) {
			s.DossierTarget = append(s.DossierTarget, e)
		}
	}
	for _, e := range sc.ContactEntries {
		if !s.hasDossierEntry("contact", e.Name) {
			s.DossierContacts = append(s.DossierContacts, e)
		}
	}
	if len(sc.SituationLines) > 0 {
		s.SituationLines = sc.SituationLines
	}
}

func (s *State) hasDossierEntry(section, name string) bool {
	list := s.DossierTarget
	if section == "contact" {
		list = s.DossierContacts
	}
	for _, e := range list {
		if e.Name == name {
			return true
		}
	}
	return false
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	result := ""
	neg := n < 0
	if neg {
		n = -n
	}
	for n > 0 {
		result = string(rune('0'+n%10)) + result
		n /= 10
	}
	if neg {
		result = "-" + result
	}
	return result
}
