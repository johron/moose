package main

import (
	"fmt"
	"os"
	"moose/internal/editor"
	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
	"golang.design/x/clipboard"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	s, err := tcell.NewScreen()
	if err != nil {
		fmt.Printf("[moose-error] %v", err)
		os.Exit(1)
	}
	if err := s.Init(); err != nil {
		fmt.Printf("[moose-error] %v", err)
		os.Exit(1)
	}

	err = clipboard.Init()
	if err != nil {
		fmt.Printf("[moose-error] %v", err)
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

		spew.Dump(m.LM.Workspaces)
		spew.Dump(m.BM)
		fmt.Println(m.DebugLog)
	}
	defer quit()

	isPasting := false

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
			} else if ev.End() {
				isPasting = false
			}
		case *tcell.EventKey:
			if !isPasting {
				m.HandleKeyInput(ev)
			}
		}

		s.Clear()
		//m.Render()
		m.Draw()
		s.Show()
	}
}