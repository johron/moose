package editor

type Mode int
const (
	ModeNormal Mode = iota
	ModeInsert
	ModePalette
)

func (mode Mode) String() string {
	switch mode {
	case ModeNormal:  return "normal"
	case ModeInsert:  return "insert"
	case ModePalette: return "palette"
	}

	return "unknown"
}