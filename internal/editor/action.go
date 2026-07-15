package editor

import(
	"strings"
	"strconv"
	"os"
)

type ActionManager = struct {
	Common  []Action
	Normal  []Action
	Insert  []Action
	Palette []Action
}

type Action = struct {
	Binding   string
	Command   string
	AskArgs bool
	Callback  func(*Model, []string)
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
				AskArgs: false,
				Callback: func(m *Model, args []string) {
					m.Screen.GetClipboard()
				},
			},
			Action{
				Binding: "ctrl+q",
				Command: "q",
				AskArgs: true,
				Callback: func(m *Model, args []string) {
					m.Quit()
				},
			},
			Action{
				Binding: "ctrl+s",
				Command: "s",
				AskArgs: false,
				Callback: func(m *Model, args []string) {
					path := ""
					if len(args) < 1 || args[0] == "" {
						if m.BM.Current().Path == "" {
							m.Mode = ModeNormal
							m.BM.PaletteBuffer.Clear()
							m.BM.PaletteBuffer.Paste("moose.error:Missing filename argument for save")
							return
						}
					}

					if len(args) > 0 && args[0] != "" {
						m.BM.Current().Path = args[0]
					}

					path = m.BM.Current().Path

					data := []byte(m.BM.Current().String())
					err := os.WriteFile(path, data, 0644)
					if err != nil {
						m.Mode = ModeNormal
						m.BM.PaletteBuffer.Clear()
						m.BM.PaletteBuffer.Paste("moose.error:Could not save to file \"" + path + "\": " + err.Error())
						return
					}

					m.BM.PaletteBuffer.Clear()
					m.BM.PaletteBuffer.Paste("moose.info:Wrote " + strconv.Itoa(len(data)) + " bytes to \"" + path + "\"")
					return
				},
			},
			Action{
				Binding: "shift+ctrl+Rune[s]",
				Command: "s",
				AskArgs: true,
			},
		},
		Normal:  []Action{
			Action{
				Binding: "Rune[w]",
				Callback: func(m *Model, args []string) {
					m.Mode = ModeWrite
				},
			},
			Action{
				Binding: "Rune[q]",
				Callback: func(m *Model, args []string) {
					m.Mode = ModePalette
					m.BM.PaletteBuffer.Clear()
					m.BM.PaletteBuffer.Insert('/')
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
				AskArgs: false,
				Callback: func(m *Model, args []string) {
					m.BM.Current().AddCursorVert(-1)
				},
			},
			Action{
				Binding: "shift+down",
				Command: "TODO:",
				AskArgs: false,
				Callback: func(m *Model, args []string) {
					m.BM.Current().AddCursorVert(1)
				},
			},
			Action{
				Binding: "esc",
				Command: "TODO:",
				AskArgs: false,
				Callback: func(m *Model, args []string) {
					m.BM.Current().ClearCursors()
				},
			},
			Action{
				Binding: "tab",
				Command: "TODO:",
				AskArgs: false,
				Callback: func(m *Model, args []string) {
					for _ = range 4 {
						m.BM.Current().Insert(' ')
					}
				},
			},
			Action{
				Binding: "shift+tab",
				Command: "TODO:",
				AskArgs: false,
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
		Palette: []Action{
			Action{
				Binding: "backspace",
				Callback: func(m *Model, args []string) {
					m.BM.PaletteBuffer.Delete()
				},
			},
			Action{
				Binding: "left",
				Callback: func(m *Model, args []string) {
					m.BM.PaletteBuffer.MoveHoriz(-1)
				},
			},
			Action{
				Binding: "right",
				Callback: func(m *Model, args []string) {
					m.BM.PaletteBuffer.MoveHoriz(1)
				},
			},
			Action{
				Binding: "enter",
				Callback: func(m *Model, _ []string) {
					input := m.BM.PaletteBuffer.String()
					args := strings.Split(input, " ")

					if args[0] == "" {
						m.Mode = ModeNormal
						return
					}

					cmd := string([]rune(args[0])[1:])
					if strings.HasPrefix(args[0], "/") {
						for _, action := range append(m.AM.Common, append(m.AM.Normal, m.AM.Insert...)...) {
							if action.Callback == nil { continue }

							if action.Command != "" && cmd == action.Command {
								m.Mode = ModeNormal
								action.Callback(m, args[1:])
								return
							}
						}

						m.Mode = ModeNormal
						m.BM.PaletteBuffer.Clear()
						m.BM.PaletteBuffer.Paste("moose.error:Unknown command \"" + cmd + "\"")
					} else {
						m.Mode = ModeNormal
						m.BM.PaletteBuffer.Clear()
						m.BM.PaletteBuffer.Paste("moose.error:Unknown palette input \"" + cmd + "\"")
					}
				},
			},
		},
	}
}