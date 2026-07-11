package tui

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
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
        case key.Matches(msg, m.Keymap.left):
            m.Buffer.MoveHoriz(-1)
        case key.Matches(msg, m.Keymap.right):
            m.Buffer.MoveHoriz(1)
        case key.Matches(msg, m.Keymap.down):
            m.Buffer.MoveVert(1)
        case key.Matches(msg, m.Keymap.up):
            m.Buffer.MoveVert(-1)
		case key.Matches(msg, m.Keymap.newline):
			m.Buffer.Insert('\n')
		case key.Matches(msg, m.Keymap.delete):
			m.Buffer.Delete()
		default:
			key := msg.Key()
			if key.Text != "" {
				for _, r := range key.Text {
					m.Buffer.Insert(r)
				}
			}
		}

		m.Viewport.SetContent(m.Buffer.String())
        
	}

	return m, nil
}
