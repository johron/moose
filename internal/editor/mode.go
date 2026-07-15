package editor

type Mode int
const (
	ModeNormal Mode = iota
	ModeWrite
	ModePalette
	ModeFind
)

func (mode Mode) String() string {
	switch mode {
	case ModeNormal:  return "normal"
	case ModeWrite:   return "write"
	case ModePalette: return "palette"
	case ModeFind:    return "find"
	}

	return "unknown"
}