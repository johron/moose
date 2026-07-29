package layout

type LayoutManager struct {
	Workspaces   map[int]Workspace
	ActiveIdx    int
	CurrentSplit SplitType
}

type Workspace struct {
	RootContainer Container
}

type SplitType int
const (
	SplitHorizontal SplitType = iota
	SplitVertical
)

func (s *SplitType) SetOpposite()  SplitType{
	if *s == SplitHorizontal {
		*s = SplitVertical
	} else {
		*s = SplitHorizontal
	}
	return *s
}

type ContainerNode interface {
	isContainerNode()
}

type ContainerBuffers struct {
	Buffers   []int
	ActiveIdx int
}

func (ContainerBuffers) isContainerNode() {}

type Container struct {
	Children       [2]ContainerNode
	Split          SplitType
	ActiveChildIdx int
}

func (Container) isContainerNode() {}

func NewLayoutManager() LayoutManager {
	return LayoutManager{
		Workspaces: map[int]Workspace{0: NewWorkspace()},
		CurrentSplit: SplitVertical,
	}
}

func NewWorkspace() Workspace {
	return Workspace{
		RootContainer: NewContainerEmpty(),
	}
}

func NewContainerEmpty() Container {
	return Container{
		Children: 	    [2]ContainerNode{
			ContainerBuffers{
				Buffers: []int{0},
				ActiveIdx: 0,
			},
			ContainerBuffers{
				Buffers: []int{0},
				ActiveIdx: 0,
			},
		},
		Split:	  	    SplitVertical,
		ActiveChildIdx: 0,
	}
}

func (c *Container) WalkAndMutateActive(fn func(cb *ContainerBuffers) ContainerNode) {
	switch child := c.Children[c.ActiveChildIdx].(type) {
	case ContainerBuffers:
		c.Children[c.ActiveChildIdx] = fn(&child)
	case Container:
		child.WalkAndMutateActive(fn)
		c.Children[c.ActiveChildIdx] = child
	default:
		cb := ContainerBuffers{Buffers: []int{}, ActiveIdx: 0}
		c.Children[c.ActiveChildIdx] = fn(&cb)
	}
}

func (lm *LayoutManager) InsertBuffer(bufferIdx int, newNode bool) {
	workspace := lm.Workspaces[lm.ActiveIdx]

	workspace.RootContainer.WalkAndMutateActive(func(cb *ContainerBuffers) ContainerNode {
		if !newNode || len(cb.Buffers) == 0 {
			cb.Buffers = append(cb.Buffers, bufferIdx)
			cb.ActiveIdx = len(cb.Buffers) - 1
			return *cb
		}

		cbNew := ContainerBuffers{
			Buffers:   []int{bufferIdx},
			ActiveIdx: 0,
		}

		return Container{
			Children:       [2]ContainerNode{*cb, cbNew},
			Split:          lm.CurrentSplit.SetOpposite(),
			ActiveChildIdx: 1,
		}
	})

	lm.Workspaces[lm.ActiveIdx] = workspace
}