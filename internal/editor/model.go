package editor

import (
	"moose/internal/buffer"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
)

type Model struct {
	Buffer        buffer.Buffer
	CommandBuffer buffer.Buffer
	Viewport      viewport.Model
	Actions       []Action
	Width         int
	Height        int
}

func NewModel() Model {
	initialBuf := buffer.Buffer{
		Rope: nil,
		CM: buffer.CursorManager{
			Cursors:    []buffer.Cursor{{Offset: 0, Goal: 0}},
			PrimaryIdx: 0,
		},
	}

	vp := viewport.New()
	vp.MouseWheelEnabled = false

	model := Model{
		Buffer:   initialBuf,
		Viewport: vp,
		Actions: DefaultActions(),
	}

	model.Viewport.SetContent(model.renderedContent())
	return model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() tea.View {
	v := tea.NewView(m.Viewport.View())
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	return v
}

func (m Model) renderedContent() string {
	content := []rune(m.Buffer.String())

	cursorMap := make(map[int]bool)
	for _, cur := range m.Buffer.CM.Cursors {
		offset := cur.Offset
		if offset < 0 {
			offset = 0
		}
		if offset > len(content) {
			offset = len(content)
		}
		cursorMap[offset] = true
	}

	out := make([]rune, 0, len(content)+len(cursorMap))

	for i, r := range content {
		if cursorMap[i] {
			out = append(out, '█')

			if r == '\n' {
				out = append(out, r)
			}
		} else {
			out = append(out, r)
		}
	}

	if cursorMap[len(content)] {
		out = append(out, '█')
	}

	return string(out)
}
