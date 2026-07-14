package editor

import (
	"strings"
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
			Buffers:       []buffer.Buffer{buffer.NewBuffer(), buffer.NewBuffer()},
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
	case ModeInsert:  return m.AM.Insert
	case ModePalette: return m.AM.Palette
	}

	return nil
}

func (m *Model) Quit() {
	m.ShouldQuit = true
}

func (m Model) Render() {
	sWidth, sHeight := m.Screen.Size()

	maxHeight := sHeight - 2 // TODO: saturating sub?

	bLen := len(m.BM.Buffers)
	bufWidth := sWidth / bLen

	for i, buf := range m.BM.Buffers {
		table := strings.SplitAfter(buf.String(), "\n")
		for j, line := range table {
			if j + 1 > maxHeight { break }
			//if len(line) + 1 > bufWidth { continue }

			r := []rune(line)
			if len(r) + 1 > bufWidth {
				r = r[:bufWidth]
			}

			m.Screen.PutStrStyled(bufWidth*i, j, string(r), m.Style)
		}

		for _, cur := range buf.CM.Cursors {
			line, col := buffer.LineCol(buf.Rope, cur.Offset)
			if line + 1 > maxHeight { continue }
			if col + 1 > bufWidth { continue }

			m.Screen.SetContent(bufWidth * i + col, line, ' ', nil, m.Style.Reverse(true))
		}
	}

	for col := 0; col < sWidth; col++ {
		m.Screen.SetContent(col, maxHeight, ' ', nil, m.Style.Reverse(true))
	}

	modeStr := m.Mode.String()
	m.Screen.PutStrStyled(sWidth - (len(modeStr) + 1), maxHeight, strings.ToUpper(modeStr), m.Style.Reverse(true))

	if m.Mode == ModePalette {
		m.Screen.PutStrStyled(0, maxHeight + 1, m.BM.PaletteBuffer.String(), m.Style)
	} else if strings.HasPrefix(m.BM.PaletteBuffer.String(), "moose.error:") {
		m.Screen.PutStrStyled(0, maxHeight + 1, string([]rune(m.BM.PaletteBuffer.String())[12:]), m.Style.Foreground(tcell.ColorRed))
	}
}

func (m *Model) HandleKeyInput(ev *tcell.EventKey) {
	for _, action := range append(m.AM.Common, m.CurrentActionSet()...) {
		if action.Binding != "" && strings.ToLower(ev.Name()) == strings.ToLower(action.Binding) {
			if action.Command != "" && action.HasArgs == true {
				m.Mode = ModePalette
				m.BM.PaletteBuffer.Paste(action.Command)
				return
			}

			action.Callback(m, []string{})
			return
		}
	}

	if m.Mode == ModeInsert {
		if ev.Key() == tcell.KeyRune {
			for _, r := range ev.Str() {
				m.BM.Current().Insert(r)
			}
		}
	}
	if m.Mode == ModePalette {
		if ev.Key() == tcell.KeyRune {
			for _, r := range ev.Str() {
				m.BM.PaletteBuffer.Insert(r)
			}
		}
	}
}