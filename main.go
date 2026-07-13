package main

import (
	"fmt"
	"os"

	"moose/internal/editor"

	//"unicode/utf8"
	"strings"
	"github.com/gdamore/tcell/v3"

	//tea "charm.land/bubbletea/v2"
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
			s.PutStrStyled(0, 2, "tried to quit", defStyle)
			return
		}

		ev := <-s.EventQ()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			for i := range m.Actions {
				action := m.Actions[i]

				if strings.ToLower(ev.Name()) == strings.ToLower(action.Binding) {
					s.PutStrStyled(0, 1, "found match", defStyle)
					action.Callback(&m, []string{})
					break
				}
			}

			//if utf8.RuneCountInString(ev.Name()) == 1 {
			if ev.Name() != "" {
				for _, r := range ev.Name() {
					m.Buffer.Insert(r)
				}
			}

			s.PutStrStyled(0, 0, strings.ToLower(ev.Name()), defStyle)
		}
		s.Show()
	}

	//model := editor.NewModel()
//
	//p := tea.NewProgram(model)
	//if _, err := p.Run(); err != nil {
	//	fmt.Printf("moose: error: %v", err)
	//	os.Exit(1)
	//}
}
