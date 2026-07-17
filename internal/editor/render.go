package editor

import (
	"strings"
	"strconv"
	"fmt"
	"moose/internal/buffer"
	"github.com/gdamore/tcell/v3/color"
)

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

func (m Model) Render() {
	sWidth, sHeight := m.Screen.Size()

	maxHeight := sHeight - 2 // TODO: saturating sub?

	bLen := len(m.BM.Buffers)
	bufWidth := sWidth / bLen

	for i := range m.BM.Buffers {
		buf := &m.BM.Buffers[i]

		primary := buf.CM.Cursors[buf.CM.PrimaryIdx]
		curLine, _ := buffer.LineCol(buf, primary.Offset)
		buf.ScrollToShow(curLine, maxHeight)

		startOffset := buffer.OffsetForLine(buf, buf.TopLine)
		endOffset := buf.Rope.Len()
		if endLine := buf.TopLine + maxHeight; endLine < buffer.LineCount(buf) {
			endOffset = buffer.OffsetForLine(buf, endLine)
		}

		visible := string(buf.Rope.Slice(startOffset, endOffset))
		table := strings.SplitAfter(visible, "\n")
		for j, line := range table {
			if j + 1 > maxHeight { break }

			nums := fmt.Sprintf("%4d ", j + buf.TopLine)
			r := []rune(nums + expandTabs(line))
			if len(r) + 1 > bufWidth {
				r = r[:bufWidth]
			}

			m.Screen.PutStrStyled(bufWidth*i, j, string(r), m.Style)
		}

		if i == m.BM.CurrentIdx {
			for _, cur := range buf.CM.Cursors {
				line, col := buffer.LineCol(buf, cur.Offset)
				screenLine := line - buf.TopLine
				if screenLine < 0 || screenLine + 1 > maxHeight { continue }

				visCol := visualCol(buffer.LineText(buf, line), col)
				if visCol + 1 > bufWidth { continue }

				ch := buffer.RuneAt(buf, cur.Offset)

				if m.Mode == ModeWrite {
					m.Screen.SetContent(bufWidth * i + visCol + 5, screenLine, ch, nil, m.Style.Reverse(true))
				} else {
					m.Screen.SetContent(bufWidth * i + visCol + 5, screenLine, ch, nil, m.Style.Reverse(true).Foreground(color.Gray))	
				}
			}
		}
	}

	for col := 0; col < sWidth; col++ {
		m.Screen.SetContent(col, maxHeight, ' ', nil, m.Style.Background(color.Black))
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

	m.Screen.PutStrStyled(sWidth - (len(modeStr) + 1), maxHeight, strings.ToUpper(modeStr), m.Style.Background(color.Black))
	
	logStr := m.DebugLog + ", " + strings.Join(m.AM.CM.Recorded, "+") + ", " + strconv.FormatBool(m.AM.CM.Recording)
	m.Screen.PutStrStyled(1, maxHeight, logStr, m.Style.Background(color.Black))

	if m.Mode == ModePalette {
		m.Screen.PutStrStyled(0, maxHeight + 1, m.BM.PaletteBuffer.String(), m.Style)

		for _, cur := range m.BM.PaletteBuffer.CM.Cursors {
			line, col := buffer.LineCol(&m.BM.PaletteBuffer, cur.Offset)
			if line + maxHeight != maxHeight { continue }
			if col + 1 > sWidth { continue }

			ch := buffer.RuneAt(&m.BM.PaletteBuffer, cur.Offset)

			m.Screen.SetContent(col, maxHeight + 1, ch, nil, m.Style.Reverse(true))
		}
	} else if strings.HasPrefix(m.BM.PaletteBuffer.String(), "moose.info:") {
		m.Screen.PutStrStyled(0, maxHeight + 1, string([]rune(m.BM.PaletteBuffer.String())[11:]), m.Style.Foreground(color.LightGray))
	} else if strings.HasPrefix(m.BM.PaletteBuffer.String(), "moose.warn:") {
		m.Screen.PutStrStyled(0, maxHeight + 1, string([]rune(m.BM.PaletteBuffer.String())[11:]), m.Style.Foreground(color.Orange))
	} else if strings.HasPrefix(m.BM.PaletteBuffer.String(), "moose.error:") {
		m.Screen.PutStrStyled(0, maxHeight + 1, string([]rune(m.BM.PaletteBuffer.String())[12:]), m.Style.Foreground(color.Red))
	}
}