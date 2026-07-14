package buffer

import (
	"slices"
	"unicode/utf8"
	"github.com/zyedidia/rope"
)

type BufferManager struct {
	Buffers       []Buffer
	CurrentIdx    int
	PaletteBuffer Buffer
}

func (bm *BufferManager) Current() *Buffer {
	return &bm.Buffers[bm.CurrentIdx]
}

type Buffer struct {
	Rope *rope.Node
	CM   CursorManager
}

func NewBuffer() Buffer {
	return Buffer{
		Rope: rope.New([]byte{}),
		CM: CursorManager{
			Cursors:    []Cursor{{Offset: 0, Goal: 0}},
			PrimaryIdx: 0,
		},
	}
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
		_, goal := LineCol(buf.Rope, cur.Offset)
		cur.Goal = goal

		delta += shift
	}
}

func (buf *Buffer) Paste(content string) {
	for _, r := range content {
		buf.Insert(r)
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
		_, goal := LineCol(buf.Rope, cur.Offset)
		cur.Goal = goal
	}

	buf.CM.DeduplicateAndSort()
}

func (buf *Buffer) Clear() {
	primaryCursor := &Cursor{
		Offset: 0,
		Goal:   0,
	}

	buf.CM.Cursors = buf.CM.Cursors[:0]
	buf.CM.Cursors = append(buf.CM.Cursors, *primaryCursor)
	buf.CM.PrimaryIdx = 0

	buf.Rope.Remove(0, buf.Rope.Len())
}

func (buf *Buffer) MoveHoriz(dir int) {
	buf.ensureRope()

	for i := range buf.CM.Cursors {
		cur := &buf.CM.Cursors[i]
		switch {
		case dir < 0:
			cur.Offset = prevRuneStart(buf.Rope, cur.Offset)
		case dir > 0:
			cur.Offset = nextRuneEnd(buf.Rope, cur.Offset)
		}
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
		lineEnd := lineContentEnd(buf.Rope, targetLine)
		lineLen := runeCount(buf.Rope, lineStart, lineEnd)

		goal := cur.Goal
		if goal > lineLen {
			goal = lineLen
		}

		cur.Offset = OffsetForLineCol(buf.Rope, targetLine, goal)
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
		lineEnd := lineContentEnd(buf.Rope, targetLine)
		lineLen := runeCount(buf.Rope, lineStart, lineEnd)

		goal := cur.Goal
		if goal > lineLen {
			goal = lineLen
		}

		newCursors = append(newCursors, Cursor{
			Offset: OffsetForLineCol(buf.Rope, targetLine, goal),
			Goal:   goal,
		})
	}

	buf.CM.Cursors = slices.Concat(buf.CM.Cursors, newCursors)

	buf.CM.DeduplicateAndSort()
}

func (buf *Buffer) ClearCursors() {
	primaryCursor := &Cursor{
		Offset: buf.CM.Cursors[buf.CM.PrimaryIdx].Offset,
		Goal:   buf.CM.Cursors[buf.CM.PrimaryIdx].Goal,
	}

	buf.CM.Cursors = buf.CM.Cursors[:0]
	buf.CM.Cursors = append(buf.CM.Cursors, *primaryCursor)
	buf.CM.PrimaryIdx = 0
}

func LineCount(r *rope.Node) int {
	return r.Count(0, r.Len(), []byte{'\n'}) + 1
}

func LineCol(r *rope.Node, offset int) (line, col int) {
	offset = normalizeOffset(r, offset)

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

	col = runeCount(r, lineStart, offset)
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

func OffsetForLineCol(r *rope.Node, line int, col int) int {
	if col <= 0 {
		return OffsetForLine(r, line)
	}

	start := OffsetForLine(r, line)
	end := lineContentEnd(r, line)

	i := start
	for n := 0; i < end && n < col; n++ {
		_, size := utf8.DecodeRune(r.Slice(i, end))
		if size <= 0 {
			size = 1
		}
		i += size
	}

	if i > end {
		return end
	}

	return i
}

func lineContentEnd(r *rope.Node, line int) int {
	nextStart := OffsetForLine(r, line+1)
	if nextStart > 0 && nextStart <= r.Len() && r.At(nextStart-1) == '\n' {
		return nextStart - 1
	}
	return nextStart
}

func runeCount(r *rope.Node, start, end int) int {
	if start < 0 {
		start = 0
	}
	if end < start {
		end = start
	}
	if end > r.Len() {
		end = r.Len()
	}

	return utf8.RuneCount(r.Slice(start, end))
}

func normalizeOffset(r *rope.Node, offset int) int {
	if offset < 0 {
		return 0
	}
	if offset > r.Len() {
		return r.Len()
	}

	for offset > 0 && offset < r.Len() && !utf8.RuneStart(r.At(offset)) {
		offset--
	}

	return offset
}

func prevRuneStart(r *rope.Node, offset int) int {
	offset = normalizeOffset(r, offset)
	if offset <= 0 {
		return 0
	}

	left := r.Slice(0, offset)
	_, size := utf8.DecodeLastRune(left)
	if size <= 0 {
		size = 1
	}

	start := offset - size
	if start < 0 {
		start = 0
	}

	return start
}

func nextRuneEnd(r *rope.Node, offset int) int {
	offset = normalizeOffset(r, offset)
	if offset >= r.Len() {
		return r.Len()
	}

	_, size := utf8.DecodeRune(r.Slice(offset, r.Len()))
	if size <= 0 {
		size = 1
	}

	end := offset + size
	if end > r.Len() {
		end = r.Len()
	}

	return end
}

func (buf Buffer) String() string {
	if buf.Rope == nil {
		return ""
	}

	return string(buf.Rope.Value())
}
