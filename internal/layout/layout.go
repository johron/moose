package layout

type LayoutManager struct {
	Workspaces []Workspace
	ActiveIdx  int
}

type Workspace struct {
	RootContainer Container
}

type SplitType int
const (
	SplitHorizontal SplitType = iota
	SplitVertical
)

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
		Workspaces: []Workspace{NewWorkspace()},
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
		Split:	  	    SplitHorizontal,
		ActiveChildIdx: 0,
	}
}

// NB: does not set the buffer manager's CurrentIdx, only the layouts active indecies!
func (c *Container) InsertBufferInContainer(bufferIdx int) {
	for i := range c.Children {
		switch child := c.Children[i].(type) {
		case ContainerBuffers:
			child.Buffers = append(child.Buffers, bufferIdx)
			child.ActiveIdx = len(child.Buffers) - 1

			c.Children[i] = child
			c.ActiveChildIdx = i
			return

		case Container:
			child.InsertBufferInContainer(bufferIdx)
			c.Children[i] = child
			c.ActiveChildIdx = i
			return
		}
	}

	c.Children[0] = ContainerBuffers{
		Buffers:   []int{bufferIdx},
		ActiveIdx: 0,
	}
	c.ActiveChildIdx = 0
}

// ny funksjon som setter in i aktiv
func (lm *LayoutManager) InsertBuffer(bufferIdx int, newNode bool) {
	// follow active child trail to find the active buffer container, if not newNode insert in there using InsertBufferInContainer,
	// else: make new containerbuffers in that current active container, if it's full then split the containerbuffers into a new containernode
	// and then a new containerbuffers is made and the bufferidx is put there
}