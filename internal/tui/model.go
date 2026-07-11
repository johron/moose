package tui

import (
	"moose/internal/buffer"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
)

type Keymap = struct {
	left, right, up, down, tab, backtab, newline, delete, quit key.Binding
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
			Cursors: []buffer.Cursor{{Offset: 0}},
			PrimaryIdx: 0,
		},
	}

	vp := viewport.New()
	vp.SetContent(initialBuf.String())

	return EditorModel{
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
}

func (m EditorModel) Init() tea.Cmd {
	return nil
}

func (m EditorModel) View() tea.View {
	v := tea.NewView(m.Viewport.View())
	v.AltScreen = true
	return v
}
