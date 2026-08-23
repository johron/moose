package editor

import (
	"github.com/gdamore/tcell/v3"
	"moose/internal/util"
	"strconv"
	"strings"
)

func (m *Model) HandleKeyInput(ev *tcell.EventKey) {

	// if m.Mode == ModeNormal && util.StandardizeBinding(ev.Name()) is numeric { // type number in normal mode to repeat a keybind x amount of times
	// m.AM.NumericCM.Recording = true, .... along these lines
	//}

	if m.AM.CM.Recording == false {
		num, err := strconv.Atoi(util.StandardizeBinding(ev.Name()))
		if err != nil {
			goto Out
		}

		if m.Mode == ModeNormal {
			m.AM.RCM.Recording = true
			m.AM.RCM.Recorded = append(m.AM.RCM.Recorded, strconv.Itoa(num))
			return
		}
	}

Out:
	for _, action := range append(m.CurrentActionSet(), m.AM.Common...) {
		for _, binding := range action.Bindings {
			switch binding.Type {
			case BindingSingle:
				if m.AM.CM.Recording {
					continue
				}

				if util.StandardizeBinding(binding.Single) == util.StandardizeBinding(ev.Name()) {
					m.AM.CM = NewChordManager()

					if len(action.Commands) > 0 && action.AskArgs == true {
						m.Mode = ModePalette
						m.BM.PaletteBuffer.Clear()
						m.BM.PaletteBuffer.Insert("/" + action.Commands[0] + " ")
						return
					}

					if m.AM.RCM.Recording == true {
						times, err := strconv.Atoi(strings.Join(m.AM.RCM.Recorded, ""))
						if err != nil {
							action.Callback(m, []string{})
							m.AM.RCM = NewChordManager()
							return
						}

						for range times {
							action.Callback(m, []string{})
						}

						m.AM.RCM = NewChordManager()
						return
					} else {
						action.Callback(m, []string{})
						return
					}
				}
			case BindingChord:
				if m.AM.CM.Recording {
					if len(binding.Chord) > len(m.AM.CM.Recorded) {
						matching := false
						for i := range m.AM.CM.Recorded {
							if m.AM.CM.Recorded[i] == binding.Chord[i] {
								matching = true
							}
						}

						if matching {
							if binding.Chord[len(m.AM.CM.Recorded)] == util.StandardizeBinding(ev.Name()) {
								m.AM.CM.Recorded = append(m.AM.CM.Recorded, util.StandardizeBinding(ev.Name()))

								if len(m.AM.CM.Recorded) == len(binding.Chord) {
									if m.AM.RCM.Recording == true {
										times, err := strconv.Atoi(strings.Join(m.AM.RCM.Recorded, ""))
										if err != nil {
											action.Callback(m, []string{})
											m.AM.RCM = NewChordManager()
											m.AM.CM = NewChordManager()
											return
										}

										for range times {
											action.Callback(m, []string{})
										}
									} else {
										action.Callback(m, []string{})
									}

									m.AM.RCM = NewChordManager()
									m.AM.CM = NewChordManager()
								}

								return
							}
						}
					}
				} else {
					if binding.Chord[0] == util.StandardizeBinding(ev.Name()) {
						m.AM.CM.Recorded = append(m.AM.CM.Recorded, util.StandardizeBinding(ev.Name()))

						if len(m.AM.CM.Recorded) == len(binding.Chord) {
							action.Callback(m, []string{})
							m.AM.RCM = NewChordManager()
							m.AM.CM = NewChordManager()
						} else {
							m.AM.CM.Recording = true
						}

						return
					}
				}
			}
		}
	}

	if m.AM.CM.Recording {
		m.AM.CM = NewChordManager()
	}

	if m.Mode == ModeWrite {
		if ev.Key() == tcell.KeyRune {
			m.BM.Current().Insert(ev.Str())
		}
	}
	if m.Mode == ModePalette {
		if ev.Key() == tcell.KeyRune {
			m.BM.PaletteBuffer.Insert(ev.Str())
		}
	}
}