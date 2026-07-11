package buffer

import (
    "github.com/gboncoffee/gopiecetable"
)

type Document struct {
	Table *gopiecetable.PieceTable[rune]
	CM CursorManager
}