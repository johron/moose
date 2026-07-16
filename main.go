package main

import (
	"fmt"
	"os"
	"moose/internal/editor"
	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

func main() {
	s, err := tcell.NewScreen()
	if err != nil {
		fmt.Printf("moose: error: %v", err)
		os.Exit(1)
	}
	if err := s.Init(); err != nil {
		fmt.Printf("moose: error: %v", err)
		os.Exit(1)
	}


	style := tcell.StyleDefault.Background(color.Reset).Foreground(color.Reset)
	s.SetStyle(style)

	m := editor.NewModel(s, style)

	s.EnableMouse()
	s.EnablePaste()
	s.Clear()

	quit := func() {
		maybePanic := recover()
		s.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	isPasting := false
	pasteBuffer := ""

	s.Clear()
	m.Render()
	s.Show()

	for {
		if m.ShouldQuit {
			return
		}

		ev := <-s.EventQ()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventPaste:
			if ev.Start() {
				isPasting = true
				pasteBuffer = "" 
			} else if ev.End() {
				isPasting = false
				m.BM.Current().Insert(pasteBuffer)
				
				s.Clear()
				m.Render()
				s.Show()
			}

		case *tcell.EventKey:
			if isPasting {
				if ev.Key() == tcell.KeyEnter || ev.Key() == tcell.KeyCtrlJ {
					pasteBuffer += "\n"
				} else if ev.Key() == tcell.KeyTab {
					pasteBuffer += "\t"
				} else if ev.Key() == tcell.KeyRune {
					pasteBuffer += ev.Str()
				}
			} else {
				m.HandleKeyInput(ev)
				
				s.Clear()
				m.Render()
				s.Show()
			}
		}
	}

}