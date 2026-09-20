package editor

import (
	"os"

	"moose/internal/buffer"
	"moose/internal/layout"
	"moose/internal/highlight"

	"github.com/gdamore/tcell/v3"
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
	Highlighter *highlight.Highlighter
}

func NewModel(screen tcell.Screen) Model {
	// blank := buffer.NewBuffer()

	model := Model{
		Screen: screen,
		Config: DefaultConfig(),
		Mode:   ModeNormal,
		BM: buffer.BufferManager{
			Buffers:       []buffer.Buffer{},
			CurrentIdx:    0,
			PaletteBuffer: buffer.NewBuffer(),
		},
		AM:         DefaultActionManager(),
		LM:         layout.NewLayoutManager(),
		ShouldQuit: false,
	}

	highlighter, err := highlight.NewHighlighter(DefaultTheme)
	if err != nil {
		os.Exit(1)
	}
	model.Highlighter = highlighter

	model.ReloadConfig()

	return model
}

func (m *Model) ReloadConfig() {
	m.Config.StyleDefault = m.Config.StyleDefault.Background(tcell.GetColor(m.Config.Colors.MainBackground)).Foreground(tcell.GetColor(m.Config.Colors.MainForeground))
	m.Screen.SetStyle(m.Config.StyleDefault)
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
