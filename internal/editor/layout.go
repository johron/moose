package editor

import (
	"moose/internal/buffer"
	"slices"
)

func (m *Model) AddBuffer(newNode bool) {
	b := buffer.NewBuffer()
	idx := len(m.BM.Buffers)
	m.BM.Buffers = append(m.BM.Buffers, b)
	m.BM.CurrentIdx = idx

	m.LM.InsertBuffer(idx, newNode)
}

func (m *Model) AddBufferFromPath(path string, newNode bool) {
	b := buffer.NewBufferFromPath(path)
	idx := len(m.BM.Buffers)
	m.BM.Buffers = append(m.BM.Buffers, b)
	m.BM.CurrentIdx = idx

	m.LM.InsertBuffer(idx, newNode)
}

func (m *Model) CycleActiveBuffer(horisontal float32, vertical float32) {
	m.LM.CycleActiveBuffer(horisontal, vertical)
	bufIdx := m.LM.GetActiveBufferIdx()
	m.BM.CurrentIdx = bufIdx
}

func (m *Model) RemoveActiveBuffer() {
    if len(m.BM.Buffers) == 0 {
        m.BM.Buffers = []buffer.Buffer{buffer.NewBuffer()}
        m.BM.CurrentIdx = 0
        return
    }

    removedIdx := m.BM.CurrentIdx
    if removedIdx < 0 || removedIdx >= len(m.BM.Buffers) {
        return
    }

    ws, ok := m.LM.Workspaces[m.LM.ActiveIdx]
    if !ok || !ws.RootContainer.ContainsBufferIdx(removedIdx) {
        return
    }

    newActiveIdx := m.LM.RemoveActiveBufferAndReindex(removedIdx)
    m.BM.Buffers = slices.Delete(m.BM.Buffers, removedIdx, removedIdx+1)

    if len(m.BM.Buffers) == 0 {
        m.BM.Buffers = []buffer.Buffer{buffer.NewBuffer()}
        m.BM.CurrentIdx = 0
        return
    }

    if newActiveIdx < 0 || newActiveIdx >= len(m.BM.Buffers) {
        newActiveIdx = len(m.BM.Buffers) - 1
    }

    m.BM.CurrentIdx = newActiveIdx
}