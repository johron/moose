package editor

import(
	"strings"
	"strconv"
	"slices"
	"os"
	"github.com/zyedidia/rope"
	"golang.design/x/clipboard"
	"moose/internal/util"
)

type ActionManager = struct {
	Common  []Action
	Normal  []Action
	Insert  []Action
	Palette []Action
}

type Action = struct {
	Binding  []string
	Command   []string
	AskArgs   bool
	Callback  func(*Model, []string)
}

func DefaultActionManager() ActionManager {
	return ActionManager{
		Common:  []Action{
			Action{
				Binding: []string{"f12"},
				Callback: func(m *Model, args []string) {
					m.Mode = ModeNormal
				},
			},
			Action{
				Binding: []string{"ctrl+q"},
				Command: []string{"q", "quit"},
				AskArgs: true,
				Callback: func(m *Model, args []string) {
					m.Quit()
				},
			},
			Action{
				Binding: []string{"ctrl+s"},
				Command: []string{"s", "save"},
				AskArgs: false,
				Callback: func(m *Model, args []string) {
					path := ""
					if len(args) < 1 || args[0] == "" {
						if m.BM.Current().Path == "" {
							m.Mode = ModeNormal
							m.BM.PaletteBuffer.Clear()
							m.BM.PaletteBuffer.Insert("moose.error:Missing filename argument for save")
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
						m.BM.PaletteBuffer.Insert("moose.error:Could not save to file \"" + path + "\": " + err.Error())
						return
					}

					m.BM.PaletteBuffer.Clear()
					m.BM.PaletteBuffer.Insert("moose.info:Wrote " + strconv.Itoa(len(data)) + " bytes to \"" + path + "\"")
					return
				},
			},
			Action{
				Binding: []string{"shift+ctrl+Rune[s]"},
				Command: []string{"s", "save"},
				AskArgs: true,
			},
			Action{
				Binding: []string{"ctrl+e"},
				Command: []string{"e", "edit"},
				AskArgs: true,
				Callback: func(m *Model, args []string) {
					if len(args) < 1 {
						m.Mode = ModeNormal
						m.BM.PaletteBuffer.Clear()
						m.BM.PaletteBuffer.Insert("moose.error:Did not specify path for edit")
						return
					}

					content, err := os.ReadFile(args[0])
					if err != nil {
						m.Mode = ModeNormal
						m.BM.PaletteBuffer.Clear()
						m.BM.PaletteBuffer.Insert("moose.error:Could not edit file \"" + args[0] + "\": " + err.Error())
						return
					}

					m.BM.Current().Clear()
					m.BM.Current().Rope = rope.New(content)
					m.BM.Current().LI.Rebuild(m.BM.Current().Rope)
					m.BM.Current().Path = args[0]
				},
			},
			Action{
				Binding: []string{"left"},
				Callback: func(m *Model, args []string) {
					if m.Mode == ModePalette {
						m.BM.PaletteBuffer.MoveHoriz(-1)
					} else {
						m.BM.Current().MoveHoriz(-1)
					}
				},
			},
			Action{
				Binding: []string{"right"},
				Callback: func(m *Model, args []string) {
					if m.Mode == ModePalette {
						m.BM.PaletteBuffer.MoveHoriz(1)					
					} else {
						m.BM.Current().MoveHoriz(1)
					}
				},
			},
			Action{
				Binding: []string{"up"},
				Callback: func(m *Model, args []string) {
					if m.Mode != ModePalette {
						m.BM.Current().MoveVert(-1)
					}
				},
			},
			Action{
				Binding: []string{"down"},
				Callback: func(m *Model, args []string) {
					if m.Mode != ModePalette {
						m.BM.Current().MoveVert(1)
					}
				},
			},
			Action{
				Binding: []string{"shift+left"},
				AskArgs: false,
				Callback: func(m *Model, args []string) {
					if m.Mode == ModePalette {
						m.BM.PaletteBuffer.MoveWordHoriz(-1)
					} else {
						m.BM.Current().MoveWordHoriz(-1)
					}
				},
			},
			Action{
				Binding: []string{"shift+right"},
				AskArgs: false,
				Callback: func(m *Model, args []string) {
					if m.Mode == ModePalette {
						m.BM.PaletteBuffer.MoveWordHoriz(1)
					} else {
						m.BM.Current().MoveWordHoriz(1)
					}
				},
			},
		},
		Normal:  []Action{
			Action{
				Binding: []string{"Rune[w]"},
				Callback: func(m *Model, args []string) {
					m.Mode = ModeWrite
				},
			},
			Action{
				Binding: []string{"Rune[q]"},
				Callback: func(m *Model, args []string) {
					m.Mode = ModePalette
					m.BM.PaletteBuffer.Clear()
					m.BM.PaletteBuffer.Insert("/")
				},
			},
			Action{
				Binding: []string{"Rune[f]"},
				AskArgs: true,
				Callback: func(m *Model, args []string) {
					m.Mode = ModePalette
					m.BM.PaletteBuffer.Clear()
					m.BM.PaletteBuffer.Insert("?")
				},
			},
			Action{
				Binding: []string{"Rune[v]"},
				Command: []string{"p", "paste"},
				AskArgs: false,
				Callback: func(m *Model, args []string) {
					text := clipboard.Read(clipboard.FmtText)

					if string(text) != "" {
						if m.Mode == ModePalette {
							m.BM.PaletteBuffer.Insert(string(text))
						} else {
							m.BM.Current().Insert(string(text))
						}
					}
				},
			},
			Action{
				Binding: []string{"Rune[j]"},
				Callback: func(m *Model, args []string) {
					if m.Mode == ModePalette {
						m.BM.PaletteBuffer.MoveHoriz(-1)
					} else {
						m.BM.Current().MoveHoriz(-1)
					}
				},
			},
		},
		Insert:  []Action{
			Action{
				Binding: []string{"shift+up"},
				Command: []string{"TODO:"},
				AskArgs: false,
				Callback: func(m *Model, args []string) {
					m.BM.Current().AddCursorVert(-1)
				},
			},
			Action{
				Binding: []string{"shift+down"},
				Command: []string{"TODO:"},
				AskArgs: false,
				Callback: func(m *Model, args []string) {
					m.BM.Current().AddCursorVert(1)
				},
			},
			Action{
				Binding: []string{"esc"},
				Command: []string{"clrc", "clearcursors"},
				AskArgs: false,
				Callback: func(m *Model, args []string) {
					m.BM.Current().ClearCursors()
				},
			},
			Action{
				Binding: []string{"tab"},
				Command: []string{"TODO:"},
				AskArgs: false,
				Callback: func(m *Model, args []string) {
					for _ = range 4 {
						m.BM.Current().Insert(" ")
					}
				},
			},
			Action{
				Binding: []string{"shift+tab"},
				Command: []string{"TODO:"},
				AskArgs: false,
				Callback: func(m *Model, args []string) {
					// TODO: implement m.BM.Current().remove...
					
				},
			},
			Action{
				Binding: []string{"backspace"},
				Callback: func(m *Model, args []string) {
					m.BM.Current().Delete()
				},
			},
			Action{
				Binding: []string{"enter"},
				Callback: func(m *Model, args []string) {
					m.BM.Current().Insert("\n")
				},
			},
		},
		Palette: []Action{
			Action{
				Binding: []string{"backspace"},
				Callback: func(m *Model, args []string) {
					m.BM.PaletteBuffer.Delete()
				},
			},
			Action{
				Binding: []string{"enter"},
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

							if len(action.Command) > 0 &&  slices.Contains(util.StandardizeBindingsArray(action.Command), util.StandardizeBinding(cmd)) {
								m.Mode = ModeNormal
								action.Callback(m, args[1:])
								return
							}
						}

						m.Mode = ModeNormal
						m.BM.PaletteBuffer.Clear()
						m.BM.PaletteBuffer.Insert("moose.error:Unknown command \"" + cmd + "\"")
					} else if strings.HasPrefix(args[0], "?") {
						m.Mode = ModeFind
						m.BM.PaletteBuffer.Clear()
						m.BM.PaletteBuffer.Insert("moose.error:Find mode unimplemented \"" + cmd + "\"")
					} else {
						m.Mode = ModeNormal
						m.BM.PaletteBuffer.Clear()
						m.BM.PaletteBuffer.Insert("moose.error:Unknown palette input \"" + cmd + "\"")
					}
				},
			},
		},
	}
}