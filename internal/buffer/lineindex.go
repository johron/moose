package buffer

import (
	"sort"

	"github.com/zyedidia/rope"
)

type LineIndex struct {
	starts []int
}

func NewLineIndex() *LineIndex {
	return &LineIndex{starts: []int{0}}
}

func NewLineIndexFromRope(r *rope.Node) *LineIndex {
	li := NewLineIndex()
	li.Rebuild(r)
	return li
}

func (li *LineIndex) Rebuild(r *rope.Node) {
	starts := []int{0}
	if r != nil {
		data := r.Value()
		for i, b := range data {
			if b == '\n' {
				starts = append(starts, i+1)
			}
		}
	}
	li.starts = starts
}

func (li *LineIndex) Count() int {
	return len(li.starts)
}

func (li *LineIndex) LineForOffset(offset int) int {
	// first index whose start is > offset; the line we want is one before that.
	i := sort.Search(len(li.starts), func(i int) bool {
		return li.starts[i] > offset
	})
	if i == 0 {
		return 0
	}
	return i - 1
}

func (li *LineIndex) OffsetForLine(line int) int {
	if line < 0 {
		return li.starts[0]
	}
	if line >= len(li.starts) {
		return li.starts[len(li.starts)-1]
	}
	return li.starts[line]
}

func (li *LineIndex) InsertAt(pos int, data []byte) {
	shift := len(data)
	if shift == 0 {
		return
	}

	lowIdx := sort.Search(len(li.starts), func(i int) bool {
		return li.starts[i] > pos
	})

	var newStarts []int
	for j := range data {
		if data[j] == '\n' {
			newStarts = append(newStarts, pos+j+1)
		}
	}

	result := make([]int, 0, len(li.starts)+len(newStarts))
	result = append(result, li.starts[:lowIdx]...)
	result = append(result, newStarts...)
	for _, o := range li.starts[lowIdx:] {
		result = append(result, o+shift)
	}

	li.starts = result
}

func (li *LineIndex) RemoveAt(start, end int) {
	if end <= start {
		return
	}
	removed := end - start

	lowIdx := sort.Search(len(li.starts), func(i int) bool {
		return li.starts[i] > start
	})
	highIdx := sort.Search(len(li.starts), func(i int) bool {
		return li.starts[i] > end
	})

	result := make([]int, 0, lowIdx+(len(li.starts)-highIdx))
	result = append(result, li.starts[:lowIdx]...)
	for _, o := range li.starts[highIdx:] {
		result = append(result, o-removed)
	}

	li.starts = result
}
