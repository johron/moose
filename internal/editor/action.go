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
				m.Buffer.MoveHoriz(-1)	
			},
		},
		Action{
			Binding: "right",
			Command: "TODO:",
			HasArgs: false,
			Callback: func(m *Model, args []string) {
				m.Buffer.MoveHoriz(1)
			},
		},
		Action{
			Binding: "up",
			Command: "TODO:",
			HasArgs: false,
			Callback: func(m *Model, args []string) {
				m.Buffer.MoveVert(-1)
			},
		},
		Action{
			Binding: "down",
			Command: "TODO:",
			HasArgs: false,
			Callback: func(m *Model, args []string) {
				m.Buffer.MoveVert(1)
			},
		},
		Action{
			Binding: "shift+up",
			Command: "TODO:",
			HasArgs: false,
			Callback: func(m *Model, args []string) {
				m.Buffer.AddCursorVert(-1)
			},
		},
		Action{
			Binding: "shift+down",
			Command: "TODO:",
			HasArgs: false,
			Callback: func(m *Model, args []string) {
				m.Buffer.AddCursorVert(1)
			},
		},
		Action{
			Binding: "esc",
			Command: "TODO:",
			HasArgs: false,
			Callback: func(m *Model, args []string) {
				m.Buffer.ClearCursors()
			},
		},
		Action{
			Binding: "tab",
			Command: "TODO:",
			HasArgs: false,
			Callback: func(m *Model, args []string) {
				//TODO: implement m.Buffer.
			},
		},
		Action{
			Binding: "shift+tab",
			Command: "TODO:",
			HasArgs: false,
			Callback: func(m *Model, args []string) {
				// TODO: implement m.Buffer.
			},
		},
		Action{
			Binding: "backspace",
			Command: "TODO:",
			HasArgs: false,
			Callback: func(m *Model, args []string) {
				m.Buffer.Delete()
			},
		},
		Action{
			Binding: "enter",
			Command: "TODO:",
			HasArgs: false,
			Callback: func(m *Model, args []string) {
				m.Buffer.Insert('\n')
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