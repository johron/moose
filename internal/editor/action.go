package editor

import(
	"strings"
	"strconv"
	"slices"
	"os"
	"github.com/zyedidia/rope"
	"golang.design/x/clipboard"
	"moose/internal/buffer"
	"moose/internal/util"
	"moose/internal/about"
	"moose/internal/layout"
)

type ActionManager = struct {
	Common  []Action
	Normal  []Action
	Insert  []Action
	Palette []Action
	CM		ChordManager
}

type Action = struct {
	Bindings  []Binding
	Commands  []string
	AskArgs   bool
	Callback  func(*Model, []string)
}

type Binding = struct {
	Type   BindingType
	Single string
	Chord  []string
}

type BindingType int
const (
	BindingSingle BindingType = iota
	BindingChord
)

type ChordManager = struct {
	Recording   bool
	Recorded    []string
}

func NewChordManager() ChordManager {
	return ChordManager{
		Recording: false,
		Recorded:  []string{},
	}
}

func DefaultActionManager() ActionManager {
	am := ActionManager{
		CM:	     NewChordManager(),
		Common:  []Action{
			Action{
				Commands: []string{"version"},
				Callback: func(m *Model, args []string) {
					m.BM.PaletteBuffer.Clear()
					m.BM.PaletteBuffer.Insert("moose.info:Moose version '" + about.Version + "'")
				},
			},
			Action{
				Bindings: []Binding{
					Binding{
						Type:   BindingSingle,
						Single: "f12",
					},
				},
				Callback: func(m *Model, args []string) {
					m.Mode = ModeNormal
				},
			},
			Action{
				Bindings: []Binding{
					Binding{
						Type:   BindingSingle,
						Single: "ctrl+q",
					},
				},
				Commands: []string{"q", "quit"},
				AskArgs: true,
				Callback: func(m *Model, args []string) {
					m.Quit()
				},
			},
			Action{
				Bindings: []Binding{
					Binding{
						Type:   BindingSingle,
						Single: "ctrl+s",
					},
				},
				Commands: []string{"s", "save"},
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
				Bindings: []Binding{
					Binding{
						Type:   BindingSingle,
						Single: "ctrl+shift+s",
					},
				},
				Commands: []string{"s", "save"},
				AskArgs: true,
			},
			Action{
				Bindings: []Binding{
					Binding{
						Type:   BindingSingle,
						Single: "ctrl+e",
					},
				},
				Commands: []string{"e", "edit"},
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
					m.BM.Current().History = buffer.UndoStack{}
				},
			},
			Action{
				Bindings: []Binding{
					Binding{
						Type:   BindingSingle,
						Single: "left",
					},
				},
				Callback: func(m *Model, args []string) {
					if m.Mode == ModePalette {
						m.BM.PaletteBuffer.MoveHoriz(-1)
					} else {
						m.BM.Current().MoveHoriz(-1)
					}
				},
			},
			Action{
				Bindings: []Binding{
					Binding{
						Type:   BindingSingle,
						Single: "right",
					},
				},
				Callback: func(m *Model, args []string) {
					if m.Mode == ModePalette {
						m.BM.PaletteBuffer.MoveHoriz(1)					
					} else {
						m.BM.Current().MoveHoriz(1)
					}
				},
			},
			Action{
				Bindings: []Binding{
					Binding{
						Type:   BindingSingle,
						Single: "up",
					},
				},
				Callback: func(m *Model, args []string) {
					if m.Mode != ModePalette {
						m.BM.Current().MoveVert(-1)
					}
				},
			},
			Action{
				Bindings: []Binding{
					Binding{
						Type:   BindingSingle,
						Single: "down",
					},
				},
				Callback: func(m *Model, args []string) {
					if m.Mode != ModePalette {
						m.BM.Current().MoveVert(1)
					}
				},
			},
			Action{
				Bindings: []Binding{
					Binding{
						Type:   BindingSingle,
						Single: "shift+left",
					},
				},
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
				Bindings: []Binding{
					Binding{
						Type:   BindingSingle,
						Single: "shift+right",
					},
				},
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
				Bindings: []Binding{
					Binding{
						Type:   BindingSingle,
						Single: "z",
					},
				},
				Commands: []string{"u", "undo"},
				AskArgs:  false,
				Callback: func(m *Model, args []string) {
					if m.Mode == ModePalette {
						m.BM.PaletteBuffer.Undo()
					} else {
						m.BM.Current().Undo()
					}
				},
			},
			Action{
				Bindings: []Binding{
					Binding{
						Type:   BindingSingle,
						Single: "y",
					},
				},
				Commands: []string{"r", "redo"},
				AskArgs:  false,
				Callback: func(m *Model, args []string) {
					if m.Mode == ModePalette {
						m.BM.PaletteBuffer.Redo()
					} else {
						m.BM.Current().Redo()
					}
				},
			},
			Action{
				Bindings: []Binding{
					Binding{
						Type:  BindingChord,
						Chord: []string{"d", "d"},
					},
				},
				Callback: func(m *Model, args []string) {
					m.BM.Current().DeleteLine()
				},
			},
			Action{
				Bindings: []Binding{
					Binding{
						Type:   BindingSingle,
						Single: "w",
					},
				},
				Callback: func(m *Model, args []string) {
					m.Mode = ModeWrite
				},
			},
			Action{
				Bindings: []Binding{
					Binding{
						Type:   BindingSingle,
						Single: "q",
					},
				},
				Callback: func(m *Model, args []string) {
					m.Mode = ModePalette
					m.BM.PaletteBuffer.Clear()
					m.BM.PaletteBuffer.Insert("/")
				},
			},
			Action{
				Bindings: []Binding{
					Binding{
						Type:   BindingSingle,
						Single: "f",
					},
				},
				AskArgs: true,
				Callback: func(m *Model, args []string) {
					m.Mode = ModePalette
					m.BM.PaletteBuffer.Clear()
					m.BM.PaletteBuffer.Insert("?")
				},
			},
			Action{
				Bindings: []Binding{
					Binding{
						Type:   BindingSingle,
						Single: "v",
					},
				},
				Commands: []string{"p", "paste"},
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
				Bindings: []Binding{
					Binding{
						Type:   BindingSingle,
						Single: "enter",
					},
				},
				Commands: []string{"buf"},
				AskArgs: false,
				Callback: func(m *Model, args []string) {
					m.AddBuffer(false)
				},
			},
			Action{
				Bindings: []Binding{
					Binding{
						Type:  BindingChord,
						Chord: []string{"l", "enter"},
					},
				},
				Commands: []string{"lbuf"},
				AskArgs: false,
				Callback: func(m *Model, args []string) {
					m.AddBuffer(true)
				},
			},
			Action{
				Bindings: []Binding{
					Binding{
						Type:   BindingSingle,
						Single: "ctrl+left",
					},
				},
				Callback: func(m *Model, args []string) {
					m.CycleActiveBuffer(-1.0, 0.0)
				},
			},
			Action{
				Bindings: []Binding{
					Binding{
						Type:   BindingSingle,
						Single: "ctrl+v",
					},
				},
				Callback: func(m *Model, args []string) {
					if m.LM.CurrentSplit == layout.SplitHorizontal {
						m.LM.CurrentSplit = layout.SplitVertical
					} else {
						m.LM.CurrentSplit = layout.SplitHorizontal
					}
				},
			},
			Action{
				Bindings: []Binding{
					Binding{
						Type:   BindingSingle,
						Single: "ctrl+down",
					},
				},
				Callback: func(m *Model, args []string) {
					m.CycleActiveBuffer(0.0, 1.0)
				},
			},
			Action{
				Bindings: []Binding{
					Binding{
						Type:   BindingSingle,
						Single: "ctrl+up",
					},
				},
				Callback: func(m *Model, args []string) {
					m.CycleActiveBuffer(0.0, -1.0)
				},
			},
			Action{
				Bindings: []Binding{
					Binding{
						Type:   BindingSingle,
						Single: "ctrl+right",
					},
				},
				Callback: func(m *Model, args []string) {
					m.CycleActiveBuffer(1.0, 0.0)
				},
			},
			Action{
				Bindings: []Binding{
					Binding{
						Type:   BindingSingle,
						Single: "ctrl+x",
					},
				},
				Callback: func(m *Model, args []string) {
					m.RemoveActiveBuffer()
				},
			},
		},
		Insert:  []Action{
			Action{
				Bindings: []Binding{
					Binding{
						Type:   BindingSingle,
						Single: "shift+up",
					},
				},
				Commands: []string{"curup"},
				Callback: func(m *Model, args []string) {
					m.BM.Current().AddCursorVert(-1)
				},
			},
			Action{
				Bindings: []Binding{
					Binding{
						Type:   BindingSingle,
						Single: "shift+down",
					},
				},
				Commands: []string{"curdown"},
				Callback: func(m *Model, args []string) {
					m.BM.Current().AddCursorVert(1)
				},
			},
			Action{
				Bindings: []Binding{
					Binding{
						Type:   BindingSingle,
						Single: "esc",
					},
				},
				Commands: []string{"clrc", "clearcursors"},
				Callback: func(m *Model, args []string) {
					m.BM.Current().ClearCursors()
				},
			},
			Action{
				Bindings: []Binding{
					Binding{
						Type:   BindingSingle,
						Single: "tab",
					},
				},
				Callback: func(m *Model, args []string) {
					for _ = range 4 {
						m.BM.Current().Insert(" ")
					}
				},
			},
			Action{
				Bindings: []Binding{
					Binding{
						Type:   BindingSingle,
						Single: "shift+tab",
					},
				},
				Callback: func(m *Model, args []string) {
					// TODO: implement m.BM.Current().remove...
					
				},
			},
			Action{
				Bindings: []Binding{
					Binding{
						Type:   BindingSingle,
						Single: "backspace",
					},
				},
				Callback: func(m *Model, args []string) {
					m.BM.Current().Delete()
				},
			},
			Action{
				Bindings: []Binding{
					Binding{
						Type:   BindingSingle,
						Single: "enter",
					},
				},
				Callback: func(m *Model, args []string) {
					m.BM.Current().Insert("\n")
				},
			},
		},
		Palette: []Action{
			Action{
				Bindings: []Binding{
					Binding{
						Type:   BindingSingle,
						Single: "backspace",
					},
				},
				Callback: func(m *Model, args []string) {
					m.BM.PaletteBuffer.Delete()
				},
			},
			Action{
				Bindings: []Binding{
					Binding{
						Type:   BindingSingle,
						Single: "enter",
					},
				},
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

							if len(action.Commands) > 0 &&  slices.Contains(util.StandardizeBindingsArray(action.Commands), util.StandardizeBinding(cmd)) {
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
						m.BM.PaletteBuffer.Insert("moose.error:Unknown palette input \"" + args[0] + "\"")
					}
				},
			},
		},
	}

	for i := range 9 {
		am.Common = append(am.Common, Action{
			Bindings: []Binding{
				Binding{
					Type:   BindingSingle,
					Single: "ctrl+" + strconv.Itoa(i + 1),
				},
			},
			Callback: func(m *Model, args []string) {
				m.Mode = ModeNormal
				m.LM.ActiveIdx = i

				if bufIdx := m.LM.GetActiveBufferIdx(); bufIdx >= 0 && bufIdx < len(m.BM.Buffers) {
					m.BM.CurrentIdx = bufIdx
				} else if len(m.BM.Buffers) > 0 {
					m.BM.CurrentIdx = 0
				}
			},
		})
	}

	return am
}