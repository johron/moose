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
	Rope    *rope.Node
	LI      *LineIndex
	CM      CursorManager
	Path    string
	TopLine int
}

func NewBuffer() Buffer {
	return Buffer{
		Rope: rope.New([]byte{}),
		LI:   NewLineIndex(),
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
	if buf.LI == nil {
		buf.LI = NewLineIndexFromRope(buf.Rope)
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
		buf.LI.InsertAt(pos, data)

		cur.Offset = pos + shift
		_, goal := LineCol(buf, cur.Offset)
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
		buf.LI.RemoveAt(start, pos)

		deleted := pos - start
		delta -= deleted

		cur.Offset = start
		_, goal := LineCol(buf, cur.Offset)
		cur.Goal = goal
	}

	buf.CM.DeduplicateAndSort()
}

func (buf *Buffer) Clear() {
	buf.ensureRope()

	primaryCursor := &Cursor{
		Offset: 0,
		Goal:   0,
	}

	buf.CM.Cursors = buf.CM.Cursors[:0]
	buf.CM.Cursors = append(buf.CM.Cursors, *primaryCursor)
	buf.CM.PrimaryIdx = 0

	buf.Rope.Remove(0, buf.Rope.Len())
	buf.LI = NewLineIndex()
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
		_, goal := LineCol(buf, cur.Offset)
		cur.Goal = goal
	}

	buf.CM.DeduplicateAndSort()
}

func (buf *Buffer) MoveVert(dir int) {
	buf.ensureRope()

	for i := range buf.CM.Cursors {
		cur := &buf.CM.Cursors[i]
		line, _ := LineCol(buf, cur.Offset)

		targetLine := line + dir
		if targetLine < 0 || targetLine >= LineCount(buf) {
			continue
		}

		lineStart := OffsetForLine(buf, targetLine)
		lineEnd := lineContentEnd(buf, targetLine)
		lineLen := runeCount(buf.Rope, lineStart, lineEnd)

		goal := cur.Goal
		if goal > lineLen {
			goal = lineLen
		}

		cur.Offset = OffsetForLineCol(buf, targetLine, goal)
	}

	buf.CM.DeduplicateAndSort()
}

func (buf *Buffer) AddCursorVert(dir int) {
	buf.ensureRope()

	var newCursors []Cursor

	for i := range buf.CM.Cursors {
		cur := &buf.CM.Cursors[i]

		line, _ := LineCol(buf, cur.Offset)

		targetLine := line + dir
		if targetLine < 0 || targetLine >= LineCount(buf) {
			continue
		}

		lineStart := OffsetForLine(buf, targetLine)
		lineEnd := lineContentEnd(buf, targetLine)
		lineLen := runeCount(buf.Rope, lineStart, lineEnd)

		goal := cur.Goal
		if goal > lineLen {
			goal = lineLen
		}

		newCursors = append(newCursors, Cursor{
			Offset: OffsetForLineCol(buf, targetLine, goal),
			Goal:   goal,
		})
	}

	buf.CM.Cursors = slices.Concat(buf.CM.Cursors, newCursors)

	buf.CM.DeduplicateAndSort()
}

func (buf *Buffer) ScrollToShow(line, maxHeight int) {
	if maxHeight <= 0 {
		return
	}
	if line < buf.TopLine {
		buf.TopLine = line
	} else if line >= buf.TopLine+maxHeight {
		buf.TopLine = line - maxHeight + 1
	}
	if buf.TopLine < 0 {
		buf.TopLine = 0
	}
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

func LineCount(buf *Buffer) int {
	return buf.LI.Count()
}

func LineCol(buf *Buffer, offset int) (line, col int) {
	offset = normalizeOffset(buf.Rope, offset)

	if offset < 0 {
		offset = 0
	}
	if offset > buf.Rope.Len() {
		offset = buf.Rope.Len()
	}

	line = buf.LI.LineForOffset(offset)
	lineStart := buf.LI.OffsetForLine(line)

	col = runeCount(buf.Rope, lineStart, offset)
	return
}

func OffsetForLine(buf *Buffer, targetLine int) int {
	if targetLine <= 0 {
		return 0
	}
	if targetLine >= buf.LI.Count() {
		return buf.Rope.Len()
	}

	return buf.LI.OffsetForLine(targetLine)
}

func OffsetForLineCol(buf *Buffer, line int, col int) int {
	if col <= 0 {
		return OffsetForLine(buf, line)
	}

	start := OffsetForLine(buf, line)
	end := lineContentEnd(buf, line)

	i := start
	for n := 0; i < end && n < col; n++ {
		_, size := utf8.DecodeRune(buf.Rope.Slice(i, end))
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

func lineContentEnd(buf *Buffer, line int) int {
	nextStart := OffsetForLine(buf, line+1)
	if nextStart > 0 && nextStart <= buf.Rope.Len() && buf.Rope.At(nextStart-1) == '\n' {
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