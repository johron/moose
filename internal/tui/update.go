package tui

import (
    tea "charm.land/bubbletea/v2"
)

func (m EditorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.Width = msg.Width
        m.Height = msg.Height
        m.Viewport.Width = msg.Width
        m.Viewport.Height = msg.Height - 2 // Space for footer status bar
        return m, nil

    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c":
            return m, tea.Quit
        default:
            m.Document.Table.Insert(msg.String(), m.Document.Table.Size())
            m.Viewport.SetContent(m.RenderBufferToString())
        }
    }
    return m, nil
}
