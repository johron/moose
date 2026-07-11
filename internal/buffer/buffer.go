package buffer

import (
	"fmt"
	"strconv"
    "github.com/gboncoffee/gopiecetable"
)

type Buffer struct {
	Table *gopiecetable.PieceTable[rune]
	CM CursorManager
}

func (buf Buffer) Insert(content rune) {
	for i := range buf.CM.Cursors {
		cur := &buf.CM.Cursors[i]
		err := buf.Table.Insert(cur.Offset, content)
		if err != nil {
			fmt.Println("moose: error" + err.Error() + strconv.Itoa(cur.Offset))
		}

		cur.Offset++
		buf.CM.DeduplicateAndSort()
	}
}

func (buf *Buffer) Delete() {
	for i := range buf.CM.Cursors {
		cur := &buf.CM.Cursors[i]
		
		if cur.Offset > 0 {
			cur.Offset--
		}

		err := buf.Table.Delete(cur.Offset)
		if err != nil {
		}

		buf.CM.DeduplicateAndSort()
	}
}

func (buf *Buffer) Undo() {
    num, _ := buf.Table.Undo()
    for i := range buf.CM.Cursors {
        buf.CM.Cursors[i].Offset = num
    }

    buf.CM.DeduplicateAndSort()
}
