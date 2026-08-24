package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gdamore/tcell/v3"
	"golang.design/x/clipboard"
	"moose/internal/editor"
	"moose/internal/extension"
	"os"
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

	m := editor.NewModel(s)
	m.AddBuffer(false)

	init_extensions(&m)

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

func init_extensions(m *editor.Model) {
	em := extension.NewExtensionManager(m)
	if err := em.LoadString("test", "print(ms)"); err != nil {
		m.Mode = editor.ModeNormal
		m.BM.PaletteBuffer.Clear()
		m.BM.PaletteBuffer.Insert("moose.error:Lua error " + err.Error())
	}
}
