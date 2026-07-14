package editor

import (
	"strings"
	"moose/internal/buffer"
	"github.com/gdamore/tcell/v3"
)

func (m Model) Render() {
	sWidth, sHeight := m.Screen.Size()

	maxHeight := sHeight - 2 // TODO: saturating sub?

	bLen := len(m.BM.Buffers)
	bufWidth := sWidth / bLen

	for i, buf := range m.BM.Buffers {
		table := strings.SplitAfter(buf.String(), "\n")
		for j, line := range table {
			if j + 1 > maxHeight { break }

			r := []rune(line)
			if len(r) + 1 > bufWidth {
				r = r[:bufWidth]
			}

			m.Screen.PutStrStyled(bufWidth*i, j, string(r), m.Style)
		}

		for _, cur := range buf.CM.Cursors {
			line, col := buffer.LineCol(buf.Rope, cur.Offset)
			if line + 1 > maxHeight { continue }
			if col + 1 > bufWidth { continue }

			if i == m.BM.CurrentIdx && m.Mode != ModePalette {
				m.Screen.SetContent(bufWidth * i + col, line, ' ', nil, m.Style.Reverse(true))
			} else {
				m.Screen.SetContent(bufWidth * i + col, line, ' ', nil, m.Style.Reverse(true).Foreground(tcell.ColorGray))				
			}
		}
	}

	for col := 0; col < sWidth; col++ {
		m.Screen.SetContent(col, maxHeight, ' ', nil, m.Style.Reverse(true))
	}

	modeStr := m.Mode.String()
	m.Screen.PutStrStyled(sWidth - (len(modeStr) + 1), maxHeight, strings.ToUpper(modeStr), m.Style.Reverse(true))

	if m.Mode == ModePalette {
		m.Screen.PutStrStyled(0, maxHeight + 1, m.BM.PaletteBuffer.String(), m.Style)

		for _, cur := range m.BM.PaletteBuffer.CM.Cursors {
			line, col := buffer.LineCol(m.BM.PaletteBuffer.Rope, cur.Offset)
			if line + maxHeight != maxHeight { continue }
			if col + 1 > sWidth { continue }

			m.Screen.SetContent(col, maxHeight + 1, ' ', nil, m.Style.Reverse(true))
		}
	} else if strings.HasPrefix(m.BM.PaletteBuffer.String(), "moose.error:") {
		m.Screen.PutStrStyled(0, maxHeight + 1, string([]rune(m.BM.PaletteBuffer.String())[12:]), m.Style.Foreground(tcell.ColorRed))
	}
}