package editor

import (
	"moose/internal/buffer"
)

func (m *Model) AddBuffer(newNode bool) {
	b := buffer.NewBuffer()
	idx := len(m.BM.Buffers)
	m.BM.Buffers = append(m.BM.Buffers, b)
	m.BM.CurrentIdx = idx

	//m.LM.Workspaces[m.LM.ActiveIdx].RootContainer.InsertBufferInContainer(idx)
	m.LM.InsertBuffer(idx, newNode)
}

func (m *Model) AddBufferFromPath(path string, newNode bool) {
	b := buffer.NewBufferFromPath(path)
	idx := len(m.BM.Buffers)
	m.BM.Buffers = append(m.BM.Buffers, b)
	m.BM.CurrentIdx = idx

	m.LM.InsertBuffer(idx, newNode)
}