package buffer

type Edit struct {
	Offset   int
	Deleted  []byte
	Inserted []byte
}

func (e Edit) Invert() Edit {
	return Edit{
		Offset:   e.Offset,
		Deleted:  e.Inserted,
		Inserted: e.Deleted,
	}
}

func (buf *Buffer) apply(e Edit) {
	buf.ensureRope()

	if len(e.Deleted) > 0 {
		end := e.Offset + len(e.Deleted)
		buf.Rope.Remove(e.Offset, end)
		buf.LI.RemoveAt(e.Offset, end)
	}
	if len(e.Inserted) > 0 {
		buf.Rope.Insert(e.Offset, e.Inserted)
		buf.LI.InsertAt(e.Offset, e.Inserted)
	}
}

type BufferTransformation struct {
	Edits         []Edit
	CursorsBefore []Cursor
	PrimaryBefore int
	CursorsAfter  []Cursor
	PrimaryAfter  int
}

type UndoStack struct {
	done   []BufferTransformation
	undone []BufferTransformation
}

func (s *UndoStack) push(t BufferTransformation) {
	s.done = append(s.done, t)
	s.undone = s.undone[:0]
}

func (s *UndoStack) CanUndo() bool {
	return len(s.done) > 0
}

func (s *UndoStack) CanRedo() bool {
	return len(s.undone) > 0
}

func cloneCursors(cs []Cursor) []Cursor {
	out := make([]Cursor, len(cs))
	copy(out, cs)
	return out
}

func (buf *Buffer) begin() (cursorsBefore []Cursor, primaryBefore int) {
	return cloneCursors(buf.CM.Cursors), buf.CM.PrimaryIdx
}

func (buf *Buffer) commit(cursorsBefore []Cursor, primaryBefore int, edits []Edit) {
	if len(edits) == 0 {
		return
	}

	buf.History.push(BufferTransformation{
		Edits:         edits,
		CursorsBefore: cursorsBefore,
		PrimaryBefore: primaryBefore,
		CursorsAfter:  cloneCursors(buf.CM.Cursors),
		PrimaryAfter:  buf.CM.PrimaryIdx,
	})
}

func (buf *Buffer) Undo() bool {
	if !buf.History.CanUndo() {
		return false
	}

	t := buf.History.done[len(buf.History.done)-1]
	buf.History.done = buf.History.done[:len(buf.History.done)-1]

	for i := len(t.Edits) - 1; i >= 0; i-- {
		buf.apply(t.Edits[i].Invert())
	}

	buf.CM.Cursors = cloneCursors(t.CursorsBefore)
	buf.CM.PrimaryIdx = t.PrimaryBefore

	buf.History.undone = append(buf.History.undone, t)
	return true
}

func (buf *Buffer) Redo() bool {
	if !buf.History.CanRedo() {
		return false
	}

	t := buf.History.undone[len(buf.History.undone)-1]
	buf.History.undone = buf.History.undone[:len(buf.History.undone)-1]

	for _, e := range t.Edits {
		buf.apply(e)
	}

	buf.CM.Cursors = cloneCursors(t.CursorsAfter)
	buf.CM.PrimaryIdx = t.PrimaryAfter

	buf.History.done = append(buf.History.done, t)
	return true
}