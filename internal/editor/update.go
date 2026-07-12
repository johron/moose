package editor

import (
	//"fmt"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	case tea.MouseWheelMsg:
		return m, nil
	case tea.KeyPressMsg:
		for i := range m.Actions {
			action := m.Actions[i]

			if key.Matches(msg, action.Binding) {
				//fmt.Println(msg, action.Binding)
				cmd := action.Callback(&m, []string{})

				m.Viewport.SetContent(m.renderedContent())
				return m, cmd
			}
		}

		key := msg.Key()
		if key.Text != "" {
			for _, r := range key.Text {
				m.Buffer.Insert(r)
			}
		}

		m.Viewport.SetContent(m.renderedContent())
		return m, nil
	}

	return m, nil
}
