package layout

type Rect struct {
	x      int
	y      int
	width  int
	height int
}

type LayoutManager struct {
	Workspaces []Workspace
	ActiveIdx  int
}

type Workspace struct {
	RootContainer Container[[]int]
}

type SplitType int
const (
	SplitHorizontal SplitType = iota
	SplitVertical
)

type ContainerNode interface {
	[]int | Container[[]int]
}

type Container[T ContainerNode] struct {
	Children 	 [2]any
	Split  		 SplitType
	ActiveBuffer *int
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

func NewContainerEmpty() Container[[]int] {
	return Container[[]int]{
		Children: 	  [2]any{},
		Split:	  	  SplitHorizontal,
		ActiveBuffer: nil,
	}
}