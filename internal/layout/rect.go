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