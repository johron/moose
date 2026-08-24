package editor

import (
	"github.com/gdamore/tcell/v3"
	"moose/internal/buffer"
	"moose/internal/layout"
)

type Model struct {
	Screen     tcell.Screen
	Config     Config
	Mode       Mode
	BM         buffer.BufferManager
	AM         ActionManager
	LM         layout.LayoutManager
	ShouldQuit bool
	DebugLog   string
}

func NewModel(screen tcell.Screen) Model {
	blank := buffer.NewBuffer()

	model := Model{
		Screen: screen,
		Config: DefaultConfig(),
		Mode:   ModeNormal,
		BM: buffer.BufferManager{
			Buffers:       []buffer.Buffer{blank},
			CurrentIdx:    0,
			PaletteBuffer: buffer.NewBuffer(),
		},
		AM:         DefaultActionManager(),
		LM:         layout.NewLayoutManager(),
		ShouldQuit: false,
	}

	return model
}
func (m *Model) CurrentActionSet() []Action {
	switch m.Mode {
	case ModeNormal:
		return m.AM.Normal
	case ModeWrite:
		return m.AM.Insert
	case ModePalette:
		return m.AM.Palette
	}

	return nil
}

func (m *Model) Quit() {
	m.ShouldQuit = true
}
