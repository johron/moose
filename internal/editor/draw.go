package editor

import (
	"fmt"
	"math"
	"moose/internal/buffer"
	"moose/internal/layout"
	"moose/internal/util"
	"slices"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v3"
)

func (m *Model) Draw() {
	screen := layout.RectFromScren(m.Screen)
	main := layout.Rect{
		X:      screen.X,
		Y:      screen.Y,
		Width:  screen.Width,
		Height: screen.Height - 2,
	}
	palette := layout.Rect{
		X:      screen.X,
		Y:      main.Height,
		Width:  screen.Width,
		Height: 2,
	}

	copy := m.LM.Workspaces[m.LM.ActiveIdx]
	m.DrawWorkspace(&copy, main)
	m.DrawPalette(palette)
}

func (m *Model) generateChordStr() string {
	chordStr := ""

	if len(m.AM.RCM.Recorded) > 0 {
		chordStr += strings.Join(m.AM.RCM.Recorded, "") + " + "
	}

	if len(m.AM.CM.Recorded) > 0 {
		chordStr += strings.Join(m.AM.CM.Recorded, "+") + " +"
	}

	return chordStr
}

func (m *Model) DrawPalette(rect layout.Rect) {
	for col := 0; col < rect.Width; col++ {
		m.Screen.SetContent(rect.X+col, rect.Y, ' ', nil, m.Config.StyleDefault.Background(tcell.GetColor(m.Config.Colors.PaletteBarBackground)).Foreground(tcell.GetColor(m.Config.Colors.PaletteBarForeground)))
	}

	modeStr := m.Mode.String()
	if m.Mode == ModePalette {
		if strings.HasPrefix(m.BM.PaletteBuffer.String(), "/") {
			modeStr += " (command)"
		} else if strings.HasPrefix(m.BM.PaletteBuffer.String(), "?") {
			modeStr += " (find)"
		} else if strings.HasPrefix(m.BM.PaletteBuffer.String(), "=") {
			modeStr += " (replace)"
		}
	}
	m.Screen.PutStrStyled(rect.Width-(len(modeStr)+1), rect.Y, strings.ToUpper(modeStr), m.Config.StyleDefault.Foreground(tcell.GetColor(m.Config.Colors.PaletteBarForeground)).Background(tcell.GetColor(m.Config.Colors.PaletteBarBackground)))

	splitStr := ""
	if m.LM.CurrentSplit == layout.SplitHorizontal {
		splitStr = "H"
	} else {
		splitStr = "V"
	}
	m.Screen.PutStrStyled(rect.Width-(len(modeStr)+1)-(m.Config.Properties.GutterWidth+1), rect.Y, splitStr, m.Config.StyleDefault.Foreground(tcell.GetColor(m.Config.Colors.PaletteBarForeground)).Background(tcell.GetColor(m.Config.Colors.PaletteBarBackground)))

	populatedWorkspaces := []int{m.LM.ActiveIdx + 1}
	for workspaceIdx, workspace := range m.LM.Workspaces {
		if workspaceIdx == m.LM.ActiveIdx {
			continue
		}

		if !workspace.IsEmpty() {
			populatedWorkspaces = append(populatedWorkspaces, workspaceIdx+1)
		}
	}

	slices.Sort(populatedWorkspaces)
	strWorkspaces := make([]string, len(populatedWorkspaces))

	for i, v := range populatedWorkspaces {
	    strWorkspaces[i] = strconv.Itoa(v)
	}

	workspacesStr := strings.Join(strWorkspaces, " ")

	m.Screen.PutStrStyled(1, rect.Y, workspacesStr, m.Config.StyleDefault.Foreground(tcell.GetColor(m.Config.Colors.WorkspaceForeground)).Background(tcell.GetColor(m.Config.Colors.WorkspaceBackground)))

	for idx, workspace := range populatedWorkspaces {
		if workspace == m.LM.ActiveIdx + 1 {
			m.Screen.SetContent(idx * 2 + 1, rect.Y, '0' + rune(m.LM.ActiveIdx + 1), nil, m.Config.StyleDefault.Background(tcell.GetColor(m.Config.Colors.WorkspaceBackgroundActive)).Foreground(tcell.GetColor(m.Config.Colors.WorkspaceForegroundActive)))
		}
	}

	chordStr := m.generateChordStr()
	m.Screen.PutStrStyled(len(workspacesStr)+2, rect.Y, chordStr, m.Config.StyleDefault.Foreground(tcell.GetColor(m.Config.Colors.PaletteBarForeground)).Background(tcell.GetColor(m.Config.Colors.PaletteBarBackground)))
	m.Screen.PutStrStyled(len(workspacesStr)+2+len(chordStr)+1, rect.Y, m.DebugLog, m.Config.StyleDefault.Foreground(tcell.GetColor(m.Config.Colors.PaletteBarForeground)).Background(tcell.GetColor(m.Config.Colors.PaletteBarBackground)))

	if m.Mode == ModePalette {
		m.Screen.PutStrStyled(0, rect.Y+1, m.BM.PaletteBuffer.String(), m.Config.StyleDefault.Foreground(tcell.GetColor(m.Config.Colors.PaletteInputForeground)).Background(tcell.GetColor(m.Config.Colors.PaletteInputBackground)))

		for _, cur := range m.BM.PaletteBuffer.CM.Cursors {
			line, col := buffer.LineCol(&m.BM.PaletteBuffer, cur.Offset)
			if line+rect.Y != rect.Y {
				continue
			}
			if col+1 > rect.Width {
				continue
			}

			ch := buffer.RuneAt(&m.BM.PaletteBuffer, cur.Offset)

			m.Screen.SetContent(col, rect.Y+1, ch, nil, m.Config.StyleDefault.Background(tcell.GetColor(m.Config.Colors.CursorColorWrite)))
		}
	} else if strings.HasPrefix(m.BM.PaletteBuffer.String(), "moose.info:") {
		m.Screen.PutStrStyled(0, rect.Y+1, string([]rune(m.BM.PaletteBuffer.String())[11:]), m.Config.StyleDefault.Foreground(tcell.GetColor(m.Config.Colors.InfoMsgForeground)).Background(tcell.GetColor(m.Config.Colors.InfoMsgBackground)))
	} else if strings.HasPrefix(m.BM.PaletteBuffer.String(), "moose.warn:") {
		m.Screen.PutStrStyled(0, rect.Y+1, string([]rune(m.BM.PaletteBuffer.String())[11:]), m.Config.StyleDefault.Foreground(tcell.GetColor(m.Config.Colors.WarnMsgForeground)).Background(tcell.GetColor(m.Config.Colors.WarnMsgBackground)))
	} else if strings.HasPrefix(m.BM.PaletteBuffer.String(), "moose.error:") {
		m.Screen.PutStrStyled(0, rect.Y+1, string([]rune(m.BM.PaletteBuffer.String())[12:]), m.Config.StyleDefault.Foreground(tcell.GetColor(m.Config.Colors.ErrorMsgForeground)).Background(tcell.GetColor(m.Config.Colors.ErrorMsgBackground)))
	}
}

func (m *Model) DrawWorkspace(w *layout.Workspace, rect layout.Rect) {
	m.DrawContainer(&w.RootContainer, rect)
}

func (m *Model) DrawContainer(c *layout.Container, rect layout.Rect) {
	length := util.NonNilLen(c.Children[:])
	if length == 0 {
		return
	}

	rect = layout.Rect{
		X:      rect.X,
		Y:      rect.Y,
		Width:  rect.Width,
		Height: rect.Height,
	}
	rect = layout.RectDivide(rect, c.Split, length)

	for i, child := range c.Children {
		main := layout.RectDisplace(rect, c.Split, i)

		switch child.(type) {
		case layout.ContainerBuffers:
			{
				cb := child.(layout.ContainerBuffers)
				tabs := layout.Rect{
					X:      main.X,
					Y:      main.Y,
					Width:  main.Width,
					Height: 1,
				}

				main.Y += 1
				main.Height -= 1

				m.DrawContainerTabs(&cb, tabs)
				m.DrawContainerBuffers(&cb, main)
			}
		case layout.Container:
			{
				c := child.(layout.Container)
				m.DrawContainer(&c, main)
			}
		}
	}
}

func (m *Model) DrawContainerTabs(c *layout.ContainerBuffers, rect layout.Rect) {
	tabs := []string{}
	activeTabIdx := -1
	for _, bufIdx := range c.Buffers {
		if bufIdx >= len(m.BM.Buffers) {
			continue
		}

		if c.ActiveIdx < len(c.Buffers) && c.Buffers[c.ActiveIdx] == bufIdx {
			activeTabIdx = len(tabs)
		}

		buf := m.BM.Buffers[bufIdx]
		if buf.Path == "" {
			tabs = append(tabs, "Buffer "+strconv.Itoa(bufIdx))
		} else {
			tabs = append(tabs, buf.Path)
		}
	}

	length := 0
	for i, tab := range tabs {
		if length+len(tab)+2 > rect.Width {
			m.Screen.PutStrStyled(rect.X+length, rect.Y, ">", m.Config.StyleDefault.Foreground(tcell.GetColor(m.Config.Colors.TabForeground)).Background(tcell.GetColor(m.Config.Colors.TabBackground)))
			return
		}

		style := m.Config.StyleDefault.Background(tcell.GetColor(m.Config.Colors.TabBackground)).Foreground(tcell.GetColor(m.Config.Colors.TabForeground))
		if i == activeTabIdx {
			style = style.Foreground(tcell.GetColor(m.Config.Colors.TabForegroundActive)).Background(tcell.GetColor(m.Config.Colors.TabBackgroundActive))
		}

		tabStr := strconv.Itoa(int(math.Abs(float64((activeTabIdx - i))))) + ": " + tab

		m.Screen.PutStrStyled(rect.X+length, rect.Y, tabStr, style)
		length += len(tabStr) + 1
	}
}

func (m *Model) DrawContainerBuffers(c *layout.ContainerBuffers, rect layout.Rect) {
	if c.ActiveIdx >= len(c.Buffers) {
		return
	}
	if c.Buffers[c.ActiveIdx] >= len(m.BM.Buffers) {
		return
	}

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
		if j+1 > rect.Height {
			break
		}

		lineNum := j + buf.TopLine
		relLine := lineNum - curLine
		nums := fmt.Sprintf("%*d ", m.Config.Properties.GutterWidth, int(math.Abs(float64(relLine))))
		r := []rune(nums + expandTabs(line))
		if len(r)+1 > rect.Width {
			r = r[:rect.Width]
		}

		m.Screen.PutStrStyled(rect.X, rect.Y+j, string(r), m.Config.StyleDefault.Foreground(tcell.GetColor(m.Config.Colors.MainForeground)).Background(tcell.GetColor(m.Config.Colors.MainBackground)))
	}

	if isActive {
		for _, cur := range buf.CM.Cursors {
			line, col := buffer.LineCol(buf, cur.Offset)
			screenLine := line - buf.TopLine
			if screenLine < 0 || screenLine+1 > rect.Height {
				continue
			}

			visCol := visualCol(buffer.LineText(buf, line), col)
			if visCol+1 > rect.Width {
				continue
			}

			ch := buffer.RuneAt(buf, cur.Offset)

			if m.Mode == ModeWrite {
				m.Screen.SetContent(rect.X+visCol+m.Config.Properties.GutterWidth+1, rect.Y+screenLine, ch, nil, m.Config.StyleDefault.Background(tcell.GetColor(m.Config.Colors.CursorColorWrite)))
			} else {
				m.Screen.SetContent(rect.X+visCol+m.Config.Properties.GutterWidth+1, rect.Y+screenLine, ch, nil, m.Config.StyleDefault.Background(tcell.GetColor(m.Config.Colors.CursorColor)))
			}
		}
	}
}

const tabWidth = 4

func expandTabs(s string) string {
	var b strings.Builder
	col := 0
	for _, r := range s {
		if r == '\t' {
			spaces := tabWidth - (col % tabWidth)
			b.WriteString(strings.Repeat(" ", spaces))
			col += spaces
		} else {
			b.WriteRune(r)
			col++
		}
	}
	return b.String()
}

func visualCol(lineText string, runeCol int) int {
	col := 0
	n := 0
	for _, r := range lineText {
		if n >= runeCol {
			break
		}
		if r == '\t' {
			col += tabWidth - (col % tabWidth)
		} else {
			col++
		}
		n++
	}
	return col
}
