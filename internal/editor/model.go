package editor

import (
	"strings"
	"moose/internal/buffer"
	"github.com/gdamore/tcell/v3"
)

type Model struct {
	Screen		  tcell.Screen
	Style		  tcell.Style
	Mode		  Mode
	BM			  buffer.BufferManager
	Actions       []Action
	ShouldQuit   bool
}

type Mode int
const (
	ModeNormal Mode = iota
	ModeInsert
	ModeCommand
)

func NewModel(screen tcell.Screen, style tcell.Style) Model {
	model := Model{
		Screen:	       screen,
		Style:		   style,
		Mode:		   ModeNormal,
		BM:            buffer.BufferManager{
			Buffers:       []buffer.Buffer{buffer.NewBuffer()},
			CurrentIdx:     0,
			CommandBuffer: buffer.NewBuffer(),
		},
		Actions:       DefaultActions(),
		ShouldQuit:    false,
	}

	return model
}

func (m *Model) Quit() {
	m.ShouldQuit = true
}

func (m *Model) Current() *buffer.Buffer {
	if m.Mode == ModeCommand {
		return &m.BM.CommandBuffer
	} else {
		return m.BM.Current()
	}
}

func (m Model) Render() {
	len := len(m.BM.Buffers)
	sWidth, _ := m.Screen.Size()
	width := sWidth / len

	for i, buf := range m.BM.Buffers {
		table := strings.SplitAfter(buf.String(), "\n")
		for j, line := range table {
			m.Screen.PutStrStyled(width * i, j, line, m.Style)
		}

		for _, cur := range buf.CM.Cursors {
			line, col := buffer.LineCol(buf.Rope, cur.Offset)
			m.Screen.SetContent(width * i + col, line, ' ', nil, m.Style.Reverse(true))
		}
	}
}

func (m *Model) HandleKeyInput(ev *tcell.EventKey) {
	for i := range m.Actions {
		action := m.Actions[i]

		if strings.ToLower(ev.Name()) == strings.ToLower(action.Binding) {
			action.Callback(m, []string{})
			return
		}
	}

	if ev.Key() == tcell.KeyRune {
		for _, r := range ev.Str() {
			m.Current().Insert(r)
		}
	}
}