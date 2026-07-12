package main

import (
	"fmt"
	"os"

	"moose/internal/editor"

	tea "charm.land/bubbletea/v2"
)

func main() {
	model := editor.NewModel()

	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Printf("moose: error: %v", err)
		os.Exit(1)
	}
}
