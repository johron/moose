package layout

import (
	"github.com/gdamore/tcell/v3"
)

type Rect struct {
	X      int
	Y      int
	Width  int
	Height int
}

func RectDivide(rect Rect, split SplitType, num int) Rect {
	switch split {
	case SplitHorizontal: return Rect{
		X: rect.X,
		Y: rect.Y,
		Width: rect.Width / num,
		Height: rect.Height,
	}
	case SplitVertical: return Rect{
		X: rect.X,
		Y: rect.Y,
		Width: rect.Width,
		Height: rect.Height / num,
	}
	default: panic("[moose-error] impossible split type")
	}
}

func RectDisplace(rect Rect, split SplitType, idx int) Rect {
	switch split {
	case SplitHorizontal: return Rect{
		X: rect.X + (rect.Width * idx),
		Y: rect.Y,
		Width: rect.Width,
		Height: rect.Height,
	}
	case SplitVertical: return Rect{
		X: rect.X,
		Y: rect.Y + (rect.Height * idx),
		Width: rect.Width,
		Height: rect.Height,
	}
	default: panic("[moose-error] impossible split type")
	}
}

func RectFromScren(s tcell.Screen) Rect {
	sWidth, sHeight := s.Size()
	return Rect{
		X: 0,
		Y: 0,
		Width: sWidth,
		Height: sHeight,
	}
}

type LayoutManager struct {
	Workspaces []Workspace
	ActiveIdx  int
}

type Workspace struct {
	RootContainer Container[ContainerBuffers]
}

type SplitType int
const (
	SplitHorizontal SplitType = iota
	SplitVertical
)

type ContainerNode interface {
	ContainerBuffers | Container[ContainerBuffers]
}

type ContainerBuffers struct {
	Buffers   []int
	ActiveIdx int
}

type Container[T ContainerNode] struct {
	Children 	   [2]any
	Split  		   SplitType
	ActiveChildIdx int
}

func NewLayoutManager() LayoutManager {
	return LayoutManager{
		Workspaces: []Workspace{NewWorkspace()},
	}
}

func NewWorkspace() Workspace {
	return Workspace{
		RootContainer: NewContainerEmpty(),
	}
}

func NewContainerEmpty() Container[ContainerBuffers] {
	return Container[ContainerBuffers]{
		Children: 	    [2]any{
			ContainerBuffers{
				Buffers: []int{0},
				ActiveIdx: 0,
			},
			ContainerBuffers{
				Buffers: []int{1},
				ActiveIdx: 0,
			},
		},
		Split:	  	    SplitVertical,
		ActiveChildIdx: 0,
	}
}