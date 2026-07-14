package main

import (
	"fmt"
	"os"
	"moose/internal/editor"
	"github.com/gdamore/tcell/v3"
)

func main() {
	s, err := tcell.NewScreen()
	if err != nil {
		fmt.Println("moose: error: %v", err)
		os.Exit(1)
	}
	if err := s.Init(); err != nil {
		fmt.Println("moose: error: %v", err)
		os.Exit(1)
	}

	style := tcell.StyleDefault.Background(tcell.ColorDefault).Foreground(tcell.ColorDefault)
	s.SetStyle(style)

	m := editor.NewModel(s, style)

	s.EnableMouse()
	s.DisablePaste()
	s.Clear()

	quit := func() {
		maybePanic := recover()
		s.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	for {
		s.Clear()

		if m.ShouldQuit {
			return
		}

		ev := <-s.EventQ()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventClipboard:
			m.BM.Current().Paste(string(ev.Data()))
		case *tcell.EventKey:
			m.HandleKeyInput(ev)
		}

		m.Render()

		s.Show()
	}
}
