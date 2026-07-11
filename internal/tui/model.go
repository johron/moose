package tui

import (
	"moose/internal/buffer"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
)

type Keymap = struct {
	left, right, up, down, cursorUp, cursorDown, clearCursors, tab, backtab, newline, delete, quit key.Binding
}

type EditorModel struct {
	Buffer buffer.Buffer
	Viewport viewport.Model
	Keymap Keymap
	Width int
	Height int
}

func NewEditorModel() EditorModel {
	initialBuf := buffer.Buffer{
		Rope: nil,
		CM: buffer.CursorManager{
			Cursors: []buffer.Cursor{{Offset: 0, Goal: 0}},
			PrimaryIdx: 0,
		},
	}

	vp := viewport.New()

	model := EditorModel{
		Buffer:   initialBuf,
		Viewport: vp,
		Keymap: Keymap{
			left: key.NewBinding(
				key.WithKeys("left"),
			),
			right: key.NewBinding(
				key.WithKeys("right"),
			),
			up: key.NewBinding(
				key.WithKeys("up"),
			),
			down: key.NewBinding(
				key.WithKeys("down"),
			),
			cursorUp: key.NewBinding(
				key.WithKeys("shift+up"),
			),
			cursorDown: key.NewBinding(
				key.WithKeys("shift+down"),
			),
			clearCursors: key.NewBinding(
				key.WithKeys("esc"),
			),
			tab: key.NewBinding(
				key.WithKeys("tab"),
			),
			backtab: key.NewBinding(
				key.WithKeys("shift+tab"),
			),
			delete: key.NewBinding(
				key.WithKeys("backspace"),
			),
			newline: key.NewBinding(
				key.WithKeys("enter"),
			),
			quit: key.NewBinding(
				key.WithKeys("ctrl+c"),
			),
		},
	}

	model.Viewport.SetContent(model.renderedContent())
	return model
}

func (m EditorModel) Init() tea.Cmd {
	return nil
}

func (m EditorModel) View() tea.View {
	v := tea.NewView(m.Viewport.View())
	v.AltScreen = true
	return v
}

func (m EditorModel) renderedContent() string {
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