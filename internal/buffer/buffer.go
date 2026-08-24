package buffer

import (
	"fmt"
	"os"
	"slices"
	"unicode"
	"unicode/utf8"

	"github.com/zyedidia/rope"
)

type BufferManager struct {
	Buffers       []Buffer
	CurrentIdx    int
	PaletteBuffer Buffer
}

func (bm *BufferManager) Current() *Buffer {
	if len(bm.Buffers) == 0 {
		return nil
	}

	if bm.CurrentIdx < 0 || bm.CurrentIdx >= len(bm.Buffers) {
		bm.CurrentIdx = len(bm.Buffers) - 1
	}

	return &bm.Buffers[bm.CurrentIdx]
}

// TODO: implement BufferType (BufferNormal, BufferVisual, BufferInteractive).
// BufferVisual (completely readonly) and BufferInteractive (readonly for user) will only be navigable the Visual mode.
// The visual mode will have abilities to modify the buffer, through managed functions which eg. a extension has defined.
// For example pressing space to check or uncheck a radio button.
// BufferVisual may be redundant, just use BufferInteractive without any defined functionality? Rename BufferInteractive to BufferVisual?
// Used for interactive things: file explorer (like emacs dired), ...

type Buffer struct {
	Rope    *rope.Node
	LI      *LineIndex
	CM      CursorManager
	Path    string
	TopLine int
	History UndoStack
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

func NewBufferFromPath(path string) Buffer {
	b := NewBuffer()
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("[moose-error] %v", err)
		os.Exit(1)
	}

	b.Clear()
	b.Rope = rope.New(content)
	b.LI.Rebuild(b.Rope)
	b.Path = path
	b.History = UndoStack{}

	return b
}

func (buf *Buffer) ensureRope() {
	if buf == nil {
		return
	}

	if buf.Rope == nil {
		buf.Rope = rope.New([]byte{})
	}
	if buf.LI == nil {
		buf.LI = NewLineIndexFromRope(buf.Rope)
	}
}

func (buf *Buffer) Insert(content string) {
	if buf == nil {
		return
	}

	buf.ensureRope()
	buf.CM.DeduplicateAndSort()

	cursorsBefore, primaryBefore := buf.begin()

	delta := 0
	data := []byte(content)
	shift := len(data)

	edits := make([]Edit, 0, len(buf.CM.Cursors))

	for i := range buf.CM.Cursors {
		cur := &buf.CM.Cursors[i]

		pos := max(cur.Offset+delta, 0)
		if pos > buf.Rope.Len() {
			pos = buf.Rope.Len()
		}

		edit := Edit{Offset: pos, Inserted: data}
		buf.apply(edit)
		edits = append(edits, edit)

		cur.Offset = pos + shift
		_, goal := LineCol(buf, cur.Offset)
		cur.Goal = goal

		delta += shift
	}

	buf.commit(cursorsBefore, primaryBefore, edits)
}

func (buf *Buffer) Delete() {
	buf.ensureRope()

	buf.CM.DeduplicateAndSort()

	cursorsBefore, primaryBefore := buf.begin()

	delta := 0
	var edits []Edit

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

		start := max(pos-size, 0)

		deleted := append([]byte{}, buf.Rope.Slice(start, pos)...)
		edit := Edit{Offset: start, Deleted: deleted}
		buf.apply(edit)
		edits = append(edits, edit)

		delta -= len(deleted)

		cur.Offset = start
		_, goal := LineCol(buf, cur.Offset)
		cur.Goal = goal
	}

	buf.CM.DeduplicateAndSort()
	buf.commit(cursorsBefore, primaryBefore, edits)
}

func (buf *Buffer) DeleteLine() {
	buf.ensureRope()
	buf.CM.DeduplicateAndSort()

	if buf.Rope.Len() == 0 || len(buf.CM.Cursors) == 0 {
		return
	}

	cursorsBefore, primaryBefore := buf.begin()

	primary := buf.CM.Cursors[buf.CM.PrimaryIdx]
	line, _ := LineCol(buf, primary.Offset)

	start := OffsetForLine(buf, line)
	end := buf.Rope.Len()

	if line+1 < LineCount(buf) {
		end = OffsetForLine(buf, line+1)
	} else if line > 0 && start > 0 && buf.Rope.At(start-1) == '\n' {
		start--
	}

	if end <= start {
		return
	}

	deleted := append([]byte{}, buf.Rope.Slice(start, end)...)
	edit := Edit{Offset: start, Deleted: deleted}
	buf.apply(edit)

	newLine := line
	if newLine >= LineCount(buf) {
		newLine = LineCount(buf) - 1
	}
	if newLine < 0 {
		newLine = 0
	}
	newOffset := OffsetForLine(buf, newLine)

	for i := range buf.CM.Cursors {
		buf.CM.Cursors[i].Offset = newOffset
		buf.CM.Cursors[i].Goal = 0
	}

	buf.CM.DeduplicateAndSort()
	buf.commit(cursorsBefore, primaryBefore, []Edit{edit})
}

func (buf *Buffer) Clear() {
	buf.ensureRope()

	cursorsBefore, primaryBefore := buf.begin()

	primaryCursor := &Cursor{
		Offset: 0,
		Goal:   0,
	}

	buf.CM.Cursors = buf.CM.Cursors[:0]
	buf.CM.Cursors = append(buf.CM.Cursors, *primaryCursor)
	buf.CM.PrimaryIdx = 0

	var edits []Edit
	if buf.Rope.Len() > 0 {
		deleted := append([]byte{}, buf.Rope.Value()...)
		edit := Edit{Offset: 0, Deleted: deleted}
		buf.apply(edit)
		edits = append(edits, edit)
	}

	buf.commit(cursorsBefore, primaryBefore, edits)
}

func isWordRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}

func (buf *Buffer) MoveWordHoriz(dir int) {
	buf.ensureRope()

	for i := range buf.CM.Cursors {
		cur := &buf.CM.Cursors[i]
		pos := normalizeOffset(buf.Rope, cur.Offset)

		switch {
		case dir > 0:
			for pos < buf.Rope.Len() {
				r, size := utf8.DecodeRune(buf.Rope.Slice(pos, buf.Rope.Len()))
				if size <= 0 {
					size = 1
				}
				if !isWordRune(r) {
					break
				}
				pos += size
			}

			for pos < buf.Rope.Len() {
				r, size := utf8.DecodeRune(buf.Rope.Slice(pos, buf.Rope.Len()))
				if size <= 0 {
					size = 1
				}
				if isWordRune(r) {
					break
				}
				pos += size
			}

		case dir < 0:
			for pos > 0 {
				start := prevRuneStart(buf.Rope, pos)
				r, _ := utf8.DecodeRune(buf.Rope.Slice(start, pos))
				if isWordRune(r) {
					break
				}
				pos = start
			}

			for pos > 0 {
				start := prevRuneStart(buf.Rope, pos)
				r, _ := utf8.DecodeRune(buf.Rope.Slice(start, pos))
				if !isWordRune(r) {
					break
				}
				pos = start
			}
		}

		cur.Offset = pos
		_, goal := LineCol(buf, cur.Offset)
		cur.Goal = goal
	}

	buf.CM.DeduplicateAndSort()
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

		goal := min(cur.Goal, lineLen)

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

		goal := buf.CM.Cursors[buf.CM.PrimaryIdx].Goal

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
	offset = max(normalizeOffset(buf.Rope, offset), 0)
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

func LineText(buf *Buffer, line int) string {
	start := OffsetForLine(buf, line)
	end := max(lineContentEnd(buf, line), start)

	return string(buf.Rope.Slice(start, end))
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

	start := max(offset-size, 0)

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

func RuneAt(buf *Buffer, offset int) rune {
	if buf.Rope == nil || offset < 0 || offset >= buf.Rope.Len() {
		return ' '
	}

	r, size := utf8.DecodeRune(buf.Rope.Slice(offset, buf.Rope.Len()))
	if size <= 0 || r == '\n' || r == '\t' {
		return ' '
	}

	return r
}

func (buf Buffer) String() string {
	if buf.Rope == nil {
		return ""
	}

	return string(buf.Rope.Value())
}
