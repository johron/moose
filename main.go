package main

import (
	"fmt"
	"os"
	"moose/internal/buffer"
	"moose/internal/editor"
	"strings"
	"github.com/gdamore/tcell/v3"
)

func main() {
	m := editor.NewModel()

	s, err := tcell.NewScreen()
	if err != nil {
		fmt.Println("moose: error: %v", err)
		os.Exit(1)
	}
	if err := s.Init(); err != nil {
		fmt.Println("moose: error: %v", err)
		os.Exit(1)
	}

	defStyle := tcell.StyleDefault.Background(tcell.ColorDefault).Foreground(tcell.ColorDefault)
	s.SetStyle(defStyle)

	s.EnableMouse()
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

		outer:
		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			for i := range m.Actions {
				action := m.Actions[i]

				if strings.ToLower(ev.Name()) == strings.ToLower(action.Binding) {
					action.Callback(&m, []string{})
					break outer
				}
			}

			if ev.Key() == tcell.KeyRune {
				for _, r := range ev.Str() {
					m.Buffer.Insert(r)
				}
			}
		}

		table := strings.SplitAfter(m.Buffer.String(), "\n")
		for i, line := range table {
			s.PutStrStyled(0, i, line, defStyle)
		}

		for _, cur := range m.Buffer.CM.Cursors {
			line, col := buffer.LineCol(m.Buffer.Rope, cur.Offset)
			s.SetContent(col, line, ' ', nil, tcell.StyleDefault.Reverse(true))
		}

		s.Show()
	}
}
