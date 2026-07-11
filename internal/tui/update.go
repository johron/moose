package tui

import (
    "github.com/gboncoffee/gopiecetable"
    tea "charm.land/bubbletea/v2"
    "charm.land/bubbles/v2/key"
)

func (m EditorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.Width = msg.Width
        m.Height = msg.Height
        m.Viewport.SetWidth(msg.Width)

        h := msg.Height - 2
        if h < 0 {
            h = 0
        }
        m.Viewport.SetHeight(h)
        return m, nil

    case tea.KeyPressMsg:
        switch {
        case key.Matches(msg, m.Keymap.quit):
            return m, tea.Quit
        case key.Matches(msg, m.Keymap.newline):
            m.Buffer.Insert('\n')
        case key.Matches(msg, m.Keymap.delete):
            m.Buffer.Delete()
        case key.Matches(msg, m.Keymap.undo):
            m.Buffer.Undo()
        case key.Matches(msg, m.Keymap.redo):
            m.Buffer.Table.Redo()
        default:
            key := msg.Key()
            if key.Text != "" {
                for _, r := range key.Text {
                    //_ = m.Buffer.Table.Insert(m.Buffer.Table.Size(), r)
                    m.Buffer.Insert(r)
                }
            }
        }

        m.Viewport.SetContent(gopiecetable.String(m.Buffer.Table))
    }

    return m, nil
}