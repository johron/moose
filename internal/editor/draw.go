package editor

import (
	"moose/internal/layout"
	"moose/internal/buffer"
	"fmt"
	"strings"
	"github.com/gdamore/tcell/v3/color"
)

func (m *Model) Draw() {
	m.DrawWorkspace(&m.LM.Workspaces[m.LM.ActiveIdx], layout.RectFromScren(m.Screen))
}

func (m *Model) DrawWorkspace(w *layout.Workspace, rect layout.Rect) {
	m.DrawContainer(&w.RootContainer, rect)
}

func (m *Model) DrawContainer(c *layout.Container[layout.ContainerBuffers], rect layout.Rect) {
	length := len(c.Children)
	rect = layout.RectDivide(rect, c.Split, length)

	for i, child := range c.Children {
		switch child.(type) {
		case layout.ContainerBuffers: {
			cb := child.(layout.ContainerBuffers)
			m.DrawContainerBuffers(&cb, layout.RectDisplace(rect, c.Split, i))
		}
		case layout.Container[layout.ContainerBuffers]: {
			c := child.(layout.Container[layout.ContainerBuffers])
			m.DrawContainer(&c, layout.RectDisplace(rect, c.Split, i))
		}
		}
	}
}

func (m *Model) DrawContainerBuffers(c *layout.ContainerBuffers, rect layout.Rect) {
	m.DrawBuffer(&m.BM.Buffers[c.Buffers[c.ActiveIdx]], c.Buffers[c.ActiveIdx] == m.BM.CurrentIdx, rect)
}

func (m *Model) DrawBuffer(buf *buffer.Buffer, isActive bool, rect layout.Rect) {
	primary := buf.CM.Cursors[buf.CM.PrimaryIdx]
	curLine, _ := buffer.LineCol(buf, primary.Offset)
	buf.ScrollToShow(curLine, rect.Height)

	startOffset := buffer.OffsetForLine(buf, buf.TopLine)
	endOffset := buf.Rope.Len()
	if endLine := buf.TopLine + rect.Height; endLine < buffer.LineCount(buf) {
		endOffset = buffer.OffsetForLine(buf, endLine)
	}

	visible := string(buf.Rope.Slice(startOffset, endOffset))
	table := strings.SplitAfter(visible, "\n")
	for j, line := range table {
		if j + 1 > rect.Height { break }

		nums := fmt.Sprintf("%4d ", j + buf.TopLine)
		r := []rune(nums + expandTabs(line))
		if len(r) + 1 > rect.Width {
			r = r[:rect.Width]
		}

		m.Screen.PutStrStyled(rect.X, rect.Y + j, string(r), m.Style)
	}

	if isActive {
		for _, cur := range buf.CM.Cursors {
			line, col := buffer.LineCol(buf, cur.Offset)
			screenLine := line - buf.TopLine
			if screenLine < 0 || screenLine + 1 > rect.Height { continue }

			visCol := visualCol(buffer.LineText(buf, line), col)
			if visCol + 1 > rect.Width { continue }

			ch := buffer.RuneAt(buf, cur.Offset)

			if m.Mode == ModeWrite {
				m.Screen.SetContent(rect.X + visCol + 5, rect.Y + screenLine, ch, nil, m.Style.Reverse(true))
			} else {
				m.Screen.SetContent(rect.X + visCol + 5, rect.Y + screenLine, ch, nil, m.Style.Reverse(true).Foreground(color.Gray))	
			}
		}
	}
}