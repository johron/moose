package editor

import(
	"strings"
)

type ActionManager = struct {
	Common  []Action
	Normal  []Action
	Insert  []Action
	Command []Action
}

type Action = struct {
	Binding  string
	Command  string
	HasArgs  bool
	Callback func(*Model, []string)
}

func DefaultActionManager() ActionManager {
	return ActionManager{
		Common:  []Action{
			Action{
				Binding: "f12",
				Callback: func(m *Model, args []string) {
					m.Mode = ModeNormal
				},
			},
			Action{
				Binding: "left",
				Callback: func(m *Model, args []string) {
					m.BM.Current().MoveHoriz(-1)	
				},
			},
			Action{
				Binding: "right",
				Callback: func(m *Model, args []string) {
					m.BM.Current().MoveHoriz(1)
				},
			},
			Action{
				Binding: "shift+ctrl+Rune[v]",
				Command: "paste",
				HasArgs: false,
				Callback: func(m *Model, args []string) {
					m.Screen.GetClipboard()
				},
			},
			Action{
				Binding: "ctrl+c",
				Command: "q",
				HasArgs: true,
				Callback: func(m *Model, args []string) {
					m.Quit()
				},
			},
		},
		Normal:  []Action{
			Action{
				Binding: "Rune[i]",
				Callback: func(m *Model, args []string) {
					m.Mode = ModeInsert
				},
			},
			Action{
				Binding: "Rune[q]",
				Callback: func(m *Model, args []string) {
					m.Mode = ModeCommand
					m.BM.CommandBuffer.Clear()
					m.BM.CommandBuffer.Insert('/')
				},
			},
		},
		Insert:  []Action{
			Action{
				Binding: "up",
				Callback: func(m *Model, args []string) {
					m.BM.Current().MoveVert(-1)
				},
			},
			Action{
				Binding: "down",
				Callback: func(m *Model, args []string) {
					m.BM.Current().MoveVert(1)
				},
			},
			Action{
				Binding: "shift+up",
				Command: "TODO:",
				HasArgs: false,
				Callback: func(m *Model, args []string) {
					m.BM.Current().AddCursorVert(-1)
				},
			},
			Action{
				Binding: "shift+down",
				Command: "TODO:",
				HasArgs: false,
				Callback: func(m *Model, args []string) {
					m.BM.Current().AddCursorVert(1)
				},
			},
			Action{
				Binding: "esc",
				Command: "TODO:",
				HasArgs: false,
				Callback: func(m *Model, args []string) {
					m.BM.Current().ClearCursors()
				},
			},
			Action{
				Binding: "tab",
				Command: "TODO:",
				HasArgs: false,
				Callback: func(m *Model, args []string) {
					for _ = range 4 {
						m.BM.Current().Insert(' ')
					}
				},
			},
			Action{
				Binding: "shift+tab",
				Command: "TODO:",
				HasArgs: false,
				Callback: func(m *Model, args []string) {
					// TODO: implement m.BM.Current().remove...
					
				},
			},
			Action{
				Binding: "backspace",
				Callback: func(m *Model, args []string) {
					m.BM.Current().Delete()
				},
			},
			Action{
				Binding: "enter",
				Callback: func(m *Model, args []string) {
					m.BM.Current().Insert('\n')
				},
			},
		},
		Command: []Action{
			Action{
				Binding: "backspace",
				Callback: func(m *Model, args []string) {
					m.BM.CommandBuffer.Delete()
				},
			},
			Action{
				Binding: "enter",
				Callback: func(m *Model, _ []string) {
					input := m.BM.CommandBuffer.String()
					args := strings.Split(input, " ")

					if args[0] == "" {
						m.Mode = ModeNormal
						return
					}

					if strings.HasPrefix(args[0], "/") {
						for _, action := range append(m.AM.Common, append(m.AM.Normal, m.AM.Insert...)...) {
							if action.Command != "" && string([]rune(args[0])[1:]) == action.Command {
								m.Mode = ModeNormal
								action.Callback(m, args[1:])
								return
							}
						}
					}

					m.Mode = ModeNormal
					// TODO: Write an error to the commandbuffer 
				},
			},
		},
	}
}