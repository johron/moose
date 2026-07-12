package editor

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

type Action = struct {
	Binding  key.Binding
	Command  string
	HasArgs  bool
	Callback func(*Model, []string) tea.Cmd
}

//type Actionmap = struct {
//	Left, Right, Up, Down, CursorUp, CursorDown, ClearCursors, Tab, Backtab, Newline, Delete, Quit Action
//}

func DefaultActions() []Action {
	return []Action{
		Action{
			Binding: key.NewBinding(
				key.WithKeys("left"),
			),
			Command: "TODO:",
			HasArgs: false,
			Callback: func(m *Model, args []string) tea.Cmd {
				m.Buffer.MoveHoriz(-1)	
				return nil
			},
		},
		Action{
			Binding: key.NewBinding(
				key.WithKeys("right"),
			),
			Command: "TODO:",
			HasArgs: false,
			Callback: func(m *Model, args []string) tea.Cmd {
				m.Buffer.MoveHoriz(1)
				return nil
			},
		},
		Action{
			Binding: key.NewBinding(
				key.WithKeys("up"),
			),
			Command: "TODO:",
			HasArgs: false,
			Callback: func(m *Model, args []string) tea.Cmd {
				m.Buffer.MoveVert(-1)
				return nil
			},
		},
		Action{
			Binding: key.NewBinding(
				key.WithKeys("down"),
			),
			Command: "TODO:",
			HasArgs: false,
			Callback: func(m *Model, args []string) tea.Cmd {
				m.Buffer.MoveVert(1)
				return nil
			},
		},
		Action{
			Binding: key.NewBinding(
				key.WithKeys("shift+up"),
			),
			Command: "TODO:",
			HasArgs: false,
			Callback: func(m *Model, args []string) tea.Cmd {
				m.Buffer.AddCursorVert(-1)
				return nil
			},
		},
		Action{
			Binding: key.NewBinding(
				key.WithKeys("shift+down"),
			),
			Command: "TODO:",
			HasArgs: false,
			Callback: func(m *Model, args []string) tea.Cmd {
				m.Buffer.AddCursorVert(1)
				return nil
			},
		},
		Action{
			Binding: key.NewBinding(
				key.WithKeys("esc"),
			),
			Command: "TODO:",
			HasArgs: false,
			Callback: func(m *Model, args []string) tea.Cmd {
				m.Buffer.ClearCursors()
				return nil
			},
		},
		Action{
			Binding: key.NewBinding(
				key.WithKeys("tab"),
			),
			Command: "TODO:",
			HasArgs: false,
			Callback: func(m *Model, args []string) tea.Cmd {
				//TODO: implement m.Buffer.
				return nil
			},
		},
		Action{
			Binding: key.NewBinding(
				key.WithKeys("shift+tab"),
			),
			Command: "TODO:",
			HasArgs: false,
			Callback: func(m *Model, args []string) tea.Cmd {
				// TODO: implement m.Buffer.
				return nil
			},
		},
		Action{
			Binding: key.NewBinding(
				key.WithKeys("backspace"),
			),
			Command: "TODO:",
			HasArgs: false,
			Callback: func(m *Model, args []string) tea.Cmd {
				m.Buffer.Delete()
				return nil
			},
		},
		Action{
			Binding: key.NewBinding(
				key.WithKeys("enter"),
			),
			Command: "TODO:",
			HasArgs: false,
			Callback: func(m *Model, args []string) tea.Cmd {
				m.Buffer.Insert('\n')
				return nil
			},
		},
		Action{
			Binding: key.NewBinding(
				key.WithKeys("ctrl+c"),
			),
			Command: "q",
			HasArgs: false,
			Callback: func(m *Model, args []string) tea.Cmd {
				return tea.Quit
			},
		},
	}
}