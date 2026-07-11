package buffer

import (
	"unicode/utf8"

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
	data := []byte(string(content))
	shift := len(data)

	for i := range buf.CM.Cursors {
		cur := &buf.CM.Cursors[i]
		if cur.Offset < 0 {
			cur.Offset = 0
		}
		if cur.Offset > buf.Rope.Len() {
			cur.Offset = buf.Rope.Len()
		}

		buf.Rope.Insert(cur.Offset, data)
		cur.Offset += shift
	}

	buf.CM.DeduplicateAndSort()
}

func (buf *Buffer) Delete() {
	buf.ensureRope()

	for i := range buf.CM.Cursors {
		cur := &buf.CM.Cursors[i]
		if cur.Offset <= 0 {
			continue
		}

		if cur.Offset > buf.Rope.Len() {
			cur.Offset = buf.Rope.Len()
		}

		left := buf.Rope.Slice(0, cur.Offset)
		_, size := utf8.DecodeLastRune(left)
		if size <= 0 {
			size = 1
		}

		start := cur.Offset - size
		if start < 0 {
			start = 0
		}

		buf.Rope.Remove(start, cur.Offset)
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
}

func (buf *Buffer) MoveVert(dir int) {
	buf.ensureRope()

	for i := range buf.CM.Cursors {
		cur := &buf.CM.Cursors[i]
		line, _ := LineCol(buf.Rope, cur.Offset)

		if line + dir < 0 || line + dir > LineCount(buf.Rope) {
			continue
		}

		line += dir
		cur.Offset = OffsetForLine(buf.Rope, line) + cur.Goal
	}
}

func LineCount(r *rope.Node) int {
    if r.Len() == 0 {
        return 0
    }

    lines := r.Count(0, r.Len(), []byte{'\n'})

    if r.At(r.Len()-1) != '\n' {
        lines++
    }

    return lines
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
