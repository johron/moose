package editor

type Action = struct {
	Binding  string
	Command  string
	HasArgs  bool
	Callback func(*Model, []string)
}

func DefaultActions() []Action {
	return []Action{
		Action{
			Binding: "left",
			Command: "TODO:",
			HasArgs: false,
			Callback: func(m *Model, args []string) {
				m.Current().MoveHoriz(-1)	
			},
		},
		Action{
			Binding: "right",
			Command: "TODO:",
			HasArgs: false,
			Callback: func(m *Model, args []string) {
				m.Current().MoveHoriz(1)
			},
		},
		Action{
			Binding: "up",
			Command: "TODO:",
			HasArgs: false,
			Callback: func(m *Model, args []string) {
				m.Current().MoveVert(-1)
			},
		},
		Action{
			Binding: "down",
			Command: "TODO:",
			HasArgs: false,
			Callback: func(m *Model, args []string) {
				m.Current().MoveVert(1)
			},
		},
		Action{
			Binding: "shift+up",
			Command: "TODO:",
			HasArgs: false,
			Callback: func(m *Model, args []string) {
				m.Current().AddCursorVert(-1)
			},
		},
		Action{
			Binding: "shift+down",
			Command: "TODO:",
			HasArgs: false,
			Callback: func(m *Model, args []string) {
				m.Current().AddCursorVert(1)
			},
		},
		Action{
			Binding: "esc",
			Command: "TODO:",
			HasArgs: false,
			Callback: func(m *Model, args []string) {
				m.Current().ClearCursors()
			},
		},
		Action{
			Binding: "tab",
			Command: "TODO:",
			HasArgs: false,
			Callback: func(m *Model, args []string) {
				//TODO: implement m.Current().
				for _ = range 4 {
					m.Current().Insert(' ')
				}
			},
		},
		Action{
			Binding: "shift+tab",
			Command: "TODO:",
			HasArgs: false,
			Callback: func(m *Model, args []string) {
				// TODO: implement m.Current().
				
			},
		},
		Action{
			Binding: "backspace",
			Command: "TODO:",
			HasArgs: false,
			Callback: func(m *Model, args []string) {
				m.Current().Delete()
			},
		},
		Action{
			Binding: "enter",
			Command: "TODO:",
			HasArgs: false,
			Callback: func(m *Model, args []string) {
				m.Current().Insert('\n')
			},
		},
		Action{
			Binding: "shift+ctrl+Rune[v]",
			Command: "paste",
			HasArgs: false,
			Callback: func(m *Model, args []string) {
				m.Screen.GetClipboard()
				
				//for _, r := range m.Screen.GetClipboard() {
				//	m.Current().Insert(r)
				//}
			},
		},
		Action{
			Binding: "ctrl+c",
			Command: "q",
			HasArgs: false,
			Callback: func(m *Model, args []string) {
				m.Quit()
			},
		},
	}
}