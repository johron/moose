package buffer

import (
	"cmp"
	"slices"
)

type Cursor struct {
	Offset int
	Goal   int
}

type CursorManager struct {
	Cursors    []Cursor
	PrimaryIdx int
}

func (cm *CursorManager) DeduplicateAndSort() {
	if len(cm.Cursors) <= 1 {
		return
	}

	primaryCursor := cm.Cursors[cm.PrimaryIdx]

	slices.SortFunc(cm.Cursors, func(a, b Cursor) int {
		return cmp.Compare(a.Offset, b.Offset)
	})

	cm.Cursors = slices.CompactFunc(cm.Cursors, func(a, b Cursor) bool {
		return a.Offset == b.Offset
	})

	for i, cur := range cm.Cursors {
		if cur == primaryCursor {
			cm.PrimaryIdx = i
			break
		}
	}
}
