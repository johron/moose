package editor

import (
	"slices"
	"github.com/gdamore/tcell/v3"
	"moose/internal/util"
)

func (m *Model) HandleKeyInput(ev *tcell.EventKey) {
	for _, action := range append(m.CurrentActionSet(), m.AM.Common...) {
		if len(action.Binding) > 0 && slices.Contains(util.StandardizeBindingsArray(action.Binding), util.StandardizeBinding(ev.Name())) {
			if action.AskArgs == true {
				m.Mode = ModePalette
				m.BM.PaletteBuffer.Clear()
				m.BM.PaletteBuffer.Insert("/" + action.Command[0] + " ")
				return
			}

			action.Callback(m, []string{})
			return
		}
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