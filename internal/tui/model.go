package tui

import (
    "moose/internal/buffer"

    "github.com/gboncoffee/gopiecetable"
    tea "charm.land/bubbletea/v2"
    "charm.land/bubbles/v2/key"
    "charm.land/bubbles/v2/viewport"
)

type Keymap = struct {
	left, right, up, down, tab, backtab, newline, delete, undo, redo, quit key.Binding
}

func (m *EditorModel) updateKeybindings() {
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
        Table: gopiecetable.FromString(""),
        CM: buffer.CursorManager{
            Cursors: []buffer.Cursor{{Offset: 0}},
            PrimaryIdx: 0,
        },
    }

    vp := viewport.New()
    vp.SetContent(gopiecetable.String(initialBuf.Table))

    return EditorModel{
        Buffer: initialBuf,
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
            undo: key.NewBinding(
                key.WithKeys("ctrl+z"),
            ),
            redo: key.NewBinding(
                key.WithKeys("ctrl+y"),
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