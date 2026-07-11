package tui

import (
    "moose/internal/buffer"

	"github.com/gboncoffee/gopiecetable"
	tea "charm.land/bubbletea/v2"
)

type EditorModel struct {
    Document buffer.Document
    Viewport viewport.Model
    Width int
    Height int
}

func NewEditorModel() EditorModel {
	initialDoc := buffer.Document{
		Table: gopiecetable.FromString("Welcome to your modular multi-cursor editor.\nPress Ctrl+A to spawn cursors.\nType freely to modify text blocks."),
		CM: buffer.CursorManager{
			Cursors: []buffer.Cursor{{Offset: 0}},
			PrimaryIdx: 0,
		},
	}

	vp := viewport.New(0, 0)

	return EditorModel{
		Document: initialDoc,
		Viewport: vp,
	}
}

func (m EditorModel) Init() tea.Cmd {
    return nil
}

func (m EditorModel) View() tea.View {
    return tea.NewView(gopiecetable.String(m.Document.Table))
}