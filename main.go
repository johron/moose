package main

import (
    "os"
    "fmt"

    "moose/internal/tui"
    tea "charm.land/bubbletea/v2"
)

func main() {
    model := tui.NewEditorModel()

    p := tea.NewProgram(model, tea.WithAltScreen())
    if _, err := p.Run(); err != nil {
        fmt.Printf("moose: error: %v", err)
        os.Exit(1)
    }
}
