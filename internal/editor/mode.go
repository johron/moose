package editor

type Mode int
const (
	ModeNormal Mode = iota
	ModeInsert
	ModeCommand
)

func (mode Mode) String() string {
	switch mode {
	case ModeNormal:  return "normal"
	case ModeInsert:  return "insert"
	case ModeCommand: return "command"
	}

	return "unknown"
}