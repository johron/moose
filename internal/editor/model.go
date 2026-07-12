package editor

import (
	"moose/internal/buffer"
	"strings"
	"unicode/utf8"

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
	content := []byte(m.Buffer.String())

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

	var out strings.Builder
	out.Grow(len(content) + len(cursorMap))

	for i := 0; i < len(content); {
		r, size := utf8.DecodeRune(content[i:])
		if r == utf8.RuneError && size == 1 {
			if cursorMap[i] {
				out.WriteRune('█')
				i++
				continue
			}

			out.WriteByte(content[i])
			i++
			continue
		}

		if cursorMap[i] {
			out.WriteRune('█')
			if r == '\n' {
				out.WriteRune('\n')
			}
			i += size
			continue
		}

		out.WriteRune(r)
		i += size
	}

	if cursorMap[len(content)] {
		out.WriteRune('█')
	}

	return out.String()
}
