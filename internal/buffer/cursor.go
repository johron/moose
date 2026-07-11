package buffer

import (
	"os"
	"fmt"
)

type Cursor struct {
    Offset int
}

type CursorManager struct {
	Cursors []Cursor
	PrimaryIdx int
}

func (cm *CursorManager) AddCursor(cur Cursor) {
    cm.Cursors = append(cm.Cursors, cur)
    cm.DeduplicateAndSort()
}

func (cm *CursorManager) DeduplicateAndSort() {
    fmt.Printf("moose: todo: implement DeduplicateAndSort()")
	os.Exit(1)
}
