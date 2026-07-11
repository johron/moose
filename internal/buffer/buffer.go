package buffer

import (
	"unicode/utf8"
	"slices"

	"github.com/zyedidia/rope"
)

type Buffer struct {
	Rope *rope.Node
	CM CursorManager
}

func (buf *Buffer) ensureRope() {
	if buf.Rope == nil {
		buf.Rope = rope.New([]byte{})
	}
}

func (buf *Buffer) Insert(content rune) {
	buf.ensureRope()
	buf.CM.DeduplicateAndSort()

	delta := 0
	data := []byte(string(content))
	shift := len(data)

	for i := range buf.CM.Cursors {
		cur := &buf.CM.Cursors[i]

		pos := cur.Offset + delta

		if pos < 0 {
			pos = 0
		}
		if pos > buf.Rope.Len() {
			pos = buf.Rope.Len()
		}

		buf.Rope.Insert(pos, data)
		cur.Offset = pos + shift

		delta += shift
	}
}

func (buf *Buffer) Delete() {
	buf.ensureRope()

	buf.CM.DeduplicateAndSort()

	delta := 0

	for i := range buf.CM.Cursors {
		cur := &buf.CM.Cursors[i]

		pos := cur.Offset + delta

		if pos <= 0 {
			cur.Offset = 0
			continue
		}

		if pos > buf.Rope.Len() {
			pos = buf.Rope.Len()
		}

		left := buf.Rope.Slice(0, pos)
		_, size := utf8.DecodeLastRune(left)
		if size <= 0 {
			size = 1
		}

		start := pos - size
		if start < 0 {
			start = 0
		}

		buf.Rope.Remove(start, pos)

		deleted := pos - start
		delta -= deleted

		cur.Offset = start
	}

	buf.CM.DeduplicateAndSort()
}

func (buf *Buffer) MoveHoriz(dir int) {
	buf.ensureRope()

	for i := range buf.CM.Cursors {
		cur := &buf.CM.Cursors[i]
		if cur.Offset + dir < 0 || cur.Offset + dir > buf.Rope.Len() {
			continue
		}

		cur.Offset += dir
		_, goal := LineCol(buf.Rope, cur.Offset)
		cur.Goal = goal
	}

	buf.CM.DeduplicateAndSort()
}

func (buf *Buffer) MoveVert(dir int) {
	buf.ensureRope()

	for i := range buf.CM.Cursors {
		cur := &buf.CM.Cursors[i]
		line, _ := LineCol(buf.Rope, cur.Offset)

		targetLine := line + dir
		if targetLine < 0 || targetLine >= LineCount(buf.Rope) {
			continue
		}

		lineStart := OffsetForLine(buf.Rope, targetLine)
		lineEnd := OffsetForLine(buf.Rope, targetLine+1)
		lineLen := lineEnd - lineStart
		if lineLen < 0 {
			lineLen = 0
		}

		goal := cur.Goal
		if goal > lineLen {
			goal = lineLen
		}

		cur.Offset = lineStart + goal
	}

	buf.CM.DeduplicateAndSort()
}

func (buf *Buffer) AddCursorVert(dir int) {
	buf.ensureRope()

	var newCursors []Cursor 

	for i := range buf.CM.Cursors {
		cur := &buf.CM.Cursors[i]

		line, _ := LineCol(buf.Rope, cur.Offset)

		targetLine := line + dir
		if targetLine < 0 || targetLine >= LineCount(buf.Rope) {
			continue
		}

		lineStart := OffsetForLine(buf.Rope, targetLine)
		lineEnd := OffsetForLine(buf.Rope, targetLine+1)
		lineLen := lineEnd - lineStart
		if lineLen < 0 {
			lineLen = 0
		}

		goal := cur.Goal
		if goal > lineLen {
			goal = lineLen
		}

		newCursors = append(newCursors, Cursor{
			Offset: lineStart + goal,
			Goal: goal,
		})
	}

	buf.CM.Cursors = slices.Concat(buf.CM.Cursors, newCursors)

	buf.CM.DeduplicateAndSort()
}

func (buf *Buffer) ClearCursors() {
	primaryCursor := &Cursor{
		Offset: buf.CM.Cursors[buf.CM.PrimaryIdx].Offset,
		Goal: buf.CM.Cursors[buf.CM.PrimaryIdx].Goal,
	}

	buf.CM.Cursors = buf.CM.Cursors[:0]
	buf.CM.Cursors = append(buf.CM.Cursors, *primaryCursor)
	buf.CM.PrimaryIdx = 0
}

func LineCount(r *rope.Node) int {
	return r.Count(0, r.Len(), []byte{'\n'}) + 1
}

func LineCol(r *rope.Node, offset int) (line, col int) {
    if offset < 0 {
        offset = 0
    }
    if offset > r.Len() {
        offset = r.Len()
    }

    line = r.Count(0, offset, []byte{'\n'})

    lineStart := 0
    for i := offset - 1; i >= 0; i-- {
        if r.At(i) == '\n' {
            lineStart = i + 1
            break
        }
    }

    col = offset - lineStart
    return
}

func OffsetForLine(r *rope.Node, targetLine int) int {
    if targetLine <= 0 {
        return 0
    }

    line := 0
    for i := 0; i < r.Len(); i++ {
        if r.At(i) == '\n' {
            line++
            if line == targetLine {
                return i + 1
            }
        }
    }

    return r.Len()
}

func (buf Buffer) String() string {
	if buf.Rope == nil {
		return ""
	}

	return string(buf.Rope.Value())
}
