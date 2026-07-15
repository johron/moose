package editor

import (
	"strings"
	"github.com/gdamore/tcell/v3"
)

func (m *Model) HandleKeyInput(ev *tcell.EventKey) {
	for _, action := range append(m.CurrentActionSet(), m.AM.Common...) {
		if action.Binding != "" && strings.ToLower(ev.Name()) == strings.ToLower(action.Binding) {
			if action.Command != "" && action.AskArgs == true {
				m.Mode = ModePalette
				m.BM.PaletteBuffer.Clear()
				m.BM.PaletteBuffer.Paste("/" + action.Command)
				return
			} else if action.Command == "" && action.AskArgs == true {
				m.Mode = ModePalette
			}

			action.Callback(m, []string{})
			return
		}
	}

	if m.Mode == ModeWrite {
		if ev.Key() == tcell.KeyRune {
			for _, r := range ev.Str() {
				m.BM.Current().Insert(r)
			}
		}
	}
	if m.Mode == ModePalette {
		if ev.Key() == tcell.KeyRune {
			for _, r := range ev.Str() {
				m.BM.PaletteBuffer.Insert(r)
			}
		}
	}
}