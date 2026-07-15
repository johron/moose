package editor

import (
	"moose/internal/buffer"
	"github.com/gdamore/tcell/v3"
)

type Model struct {
	Screen	   tcell.Screen
	Style	   tcell.Style
	Mode 	   Mode
	BM		   buffer.BufferManager
	AM         ActionManager
	ShouldQuit bool
}

func NewModel(screen tcell.Screen, style tcell.Style) Model {
	model := Model{
		Screen:	       screen,
		Style:		   style,
		Mode:		   ModeNormal,
		BM:            buffer.BufferManager{
			Buffers:       []buffer.Buffer{buffer.NewBuffer()},
			CurrentIdx:    0,
			PaletteBuffer: buffer.NewBuffer(),
		},
		AM:       DefaultActionManager(),
		ShouldQuit:    false,
	}

	return model
}

func (m *Model) CurrentActionSet() []Action {
	switch m.Mode {
	case ModeNormal:  return m.AM.Normal
	case ModeWrite:  return m.AM.Insert
	case ModePalette: return m.AM.Palette
	}

	return nil
}

func (m *Model) Quit() {
	m.ShouldQuit = true
}