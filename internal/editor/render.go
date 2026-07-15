package editor

import (
	"strings"
	"moose/internal/buffer"
	"github.com/gdamore/tcell/v3/color"
)

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

			r := []rune(line)
			if len(r) + 1 > bufWidth {
				r = r[:bufWidth]
			}

			m.Screen.PutStrStyled(bufWidth*i, j, string(r), m.Style)
		}

		if i == m.BM.CurrentIdx || m.Mode == ModePalette {
			for _, cur := range buf.CM.Cursors {
				line, col := buffer.LineCol(buf, cur.Offset)
				screenLine := line - buf.TopLine
				if screenLine < 0 || screenLine + 1 > maxHeight { continue }
				if col + 1 > bufWidth { continue }

				if m.Mode == ModeWrite {
					m.Screen.SetContent(bufWidth * i + col, screenLine, ' ', nil, m.Style.Reverse(true))				
				} else {
					m.Screen.SetContent(bufWidth * i + col, screenLine, ' ', nil, m.Style.Reverse(true).Foreground(color.Gray))	
				}
			}
		}
	}

	for col := 0; col < sWidth; col++ {
		m.Screen.SetContent(col, maxHeight, ' ', nil, m.Style.Reverse(true))
	}

	modeStr := m.Mode.String()
	if m.Mode == ModePalette {
		if strings.HasPrefix(m.BM.PaletteBuffer.String(), "/") {
			modeStr = "command"
		} else if strings.HasPrefix(m.BM.PaletteBuffer.String(), "?") {
			modeStr = "find"
		} else if strings.HasPrefix(m.BM.PaletteBuffer.String(), "=") {
			modeStr = "replace"
		}
	}

	m.Screen.PutStrStyled(sWidth - (len(modeStr) + 1), maxHeight, strings.ToUpper(modeStr), m.Style.Reverse(true))

	if m.Mode == ModePalette {
		m.Screen.PutStrStyled(0, maxHeight + 1, m.BM.PaletteBuffer.String(), m.Style)

		for _, cur := range m.BM.PaletteBuffer.CM.Cursors {
			line, col := buffer.LineCol(&m.BM.PaletteBuffer, cur.Offset)
			if line + maxHeight != maxHeight { continue }
			if col + 1 > sWidth { continue }

			m.Screen.SetContent(col, maxHeight + 1, ' ', nil, m.Style.Reverse(true))
		}
	} else if strings.HasPrefix(m.BM.PaletteBuffer.String(), "moose.info:") {
		m.Screen.PutStrStyled(0, maxHeight + 1, string([]rune(m.BM.PaletteBuffer.String())[11:]), m.Style.Foreground(color.LightGray))
	} else if strings.HasPrefix(m.BM.PaletteBuffer.String(), "moose.warn:") {
		m.Screen.PutStrStyled(0, maxHeight + 1, string([]rune(m.BM.PaletteBuffer.String())[11:]), m.Style.Foreground(color.Orange))
	} else if strings.HasPrefix(m.BM.PaletteBuffer.String(), "moose.error:") {
		m.Screen.PutStrStyled(0, maxHeight + 1, string([]rune(m.BM.PaletteBuffer.String())[12:]), m.Style.Foreground(color.Red))
	}
}