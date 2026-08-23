package layout

type LayoutManager struct {
	Workspaces   map[int]Workspace
	ActiveIdx    int
	CurrentSplit SplitType
}

type Workspace struct { // TODO: this is not necessary..
	RootContainer Container
}

func (w Workspace) IsEmpty() bool {
    return !isPopulated(w.RootContainer)
}

type SplitType int
const (
	SplitHorizontal SplitType = iota
	SplitVertical
)

func (s *SplitType) SetOpposite()  SplitType{
	if *s == SplitHorizontal {
		*s = SplitVertical
        return SplitHorizontal
	} else {
		*s = SplitHorizontal
        return SplitVertical
	}
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
		CurrentSplit: SplitHorizontal,
	}
}

func NewWorkspace() Workspace {
	return Workspace{
		RootContainer: NewContainerEmpty(),
	}
}

func NewContainerEmpty() Container {
    return Container{
        Children: [2]ContainerNode{
            nil,
            nil,
        },
        Split:          SplitHorizontal,
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

func (lm *LayoutManager) CycleActiveBuffer(horisontal float32, vertical float32) {
    if horisontal == 0 && vertical == 0 {
        return
    }

    ws := lm.Workspaces[lm.ActiveIdx]

    if newRoot, moved := ws.RootContainer.NavigateContainer(horisontal, vertical); moved {
        ws.RootContainer = newRoot
        lm.Workspaces[lm.ActiveIdx] = ws
    }
}

func (c Container) NavigateContainer(horisontal float32, vertical float32) (Container, bool) {
    activeChild := c.Children[c.ActiveChildIdx]

    if childContainer, ok := activeChild.(Container); ok {
        updatedChild, moved := childContainer.NavigateContainer(horisontal, vertical)
        if moved {
            c.Children[c.ActiveChildIdx] = updatedChild
            return c, true
        }
    }

    targetIdx := -1

    if horisontal != 0 && c.Split == SplitHorizontal {
        if horisontal < 0 && c.ActiveChildIdx == 1 {
            targetIdx = 0 // Try moving Left
        } else if horisontal > 0 && c.ActiveChildIdx == 0 {
            targetIdx = 1 // Try moving Right
        }
    } else if vertical != 0 && c.Split == SplitVertical {
        if vertical < 0 && c.ActiveChildIdx == 1 {
            targetIdx = 0 // Try moving Up
        } else if vertical > 0 && c.ActiveChildIdx == 0 {
            targetIdx = 1 // Try moving Down
        }
    }

    if targetIdx != -1 && isPopulated(c.Children[targetIdx]) {
        c.ActiveChildIdx = targetIdx
        return c, true
    }

    return c, false
}

func isPopulated(node ContainerNode) bool {
    if node == nil {
        return false
    }

    switch n := node.(type) {
    case ContainerBuffers:
        return len(n.Buffers) > 0
    case Container:
        return isPopulated(n.Children[0]) || isPopulated(n.Children[1])
    default:
        return false
    }
}

func (c Container) GetActiveBufferIdx() int {
    switch child := c.Children[c.ActiveChildIdx].(type) {
    case ContainerBuffers:
        if len(child.Buffers) > 0 && child.ActiveIdx < len(child.Buffers) {
            return child.Buffers[child.ActiveIdx]
        }
    case Container:
        return child.GetActiveBufferIdx()
    }
    return -1
}

func (lm LayoutManager) GetActiveBufferIdx() int {
    ws, ok := lm.Workspaces[lm.ActiveIdx]
    if !ok {
        return -1
    }
    return ws.RootContainer.GetActiveBufferIdx()
}

func (c Container) RemoveActive() (ContainerNode, bool) {
    activeChild := c.Children[c.ActiveChildIdx]

    switch child := activeChild.(type) {
    case ContainerBuffers:
        if len(child.Buffers) == 0 {
            return nil, true
        }

        child.Buffers = append(child.Buffers[:child.ActiveIdx], child.Buffers[child.ActiveIdx+1:]...)

        if child.ActiveIdx >= len(child.Buffers) && len(child.Buffers) > 0 {
            child.ActiveIdx = len(child.Buffers) - 1
        }

        if len(child.Buffers) == 0 {
            siblingIdx := 1 - c.ActiveChildIdx
            if siblingIdx < 0 || siblingIdx >= len(c.Children) || c.Children[siblingIdx] == nil {
                return nil, true
            }
            return c.Children[siblingIdx], false
        }

        c.Children[c.ActiveChildIdx] = child
        return c, false

    case Container:
        updatedChild, childIsEmpty := child.RemoveActive()
        if childIsEmpty {
            siblingIdx := 1 - c.ActiveChildIdx
            if siblingIdx < 0 || siblingIdx >= len(c.Children) || c.Children[siblingIdx] == nil {
                return nil, true
            }
            return c.Children[siblingIdx], false
        }

        c.Children[c.ActiveChildIdx] = updatedChild
        return c, false

    default:
        return nil, true
    }
}

func (c Container) ContainsBufferIdx(target int) bool {
    switch child := c.Children[c.ActiveChildIdx].(type) {
    case ContainerBuffers:
        for _, idx := range child.Buffers {
            if idx == target {
                return true
            }
        }
    case Container:
        if child.ContainsBufferIdx(target) {
            return true
        }
    }

    for _, node := range c.Children {
        switch child := node.(type) {
        case ContainerBuffers:
            for _, idx := range child.Buffers {
                if idx == target {
                    return true
                }
            }
        case Container:
            if child.ContainsBufferIdx(target) {
                return true
            }
        }
    }

    return false
}

func (c *Container) ActivateBufferIdx(target int) bool {
    for i, child := range c.Children {
        switch n := child.(type) {
        case ContainerBuffers:
            for j, idx := range n.Buffers {
                if idx == target {
                    c.ActiveChildIdx = i
                    n.ActiveIdx = j
                    c.Children[i] = n
                    return true
                }
            }
        case Container:
            if n.ActivateBufferIdx(target) {
                c.Children[i] = n
                return true
            }
        }
    }
    return false
}

func (lm *LayoutManager) RemoveActiveBufferAndReindex(removedIdx int) int {
    ws := lm.Workspaces[lm.ActiveIdx]

    if !ws.RootContainer.ContainsBufferIdx(removedIdx) {
        return ws.RootContainer.GetActiveBufferIdx()
    }

    if !ws.RootContainer.ActivateBufferIdx(removedIdx) {
        return ws.RootContainer.GetActiveBufferIdx()
    }

    newRoot, isEmpty := ws.RootContainer.RemoveActive()
    if newRoot == nil || isEmpty {
        ws.RootContainer = NewContainerEmpty()
        lm.Workspaces[lm.ActiveIdx] = ws
        return -1
    }

    switch root := newRoot.(type) {
    case Container:
        ws.RootContainer = root
    case ContainerBuffers:
        ws.RootContainer = Container{
            Children:       [2]ContainerNode{root, nil},
            Split:          ws.RootContainer.Split,
            ActiveChildIdx: 0,
        }
    }

    ws.RootContainer.DecrementIndicesAbove(removedIdx)
    lm.Workspaces[lm.ActiveIdx] = ws

    return ws.RootContainer.GetActiveBufferIdx()
}

func (c Container) DecrementIndicesAbove(threshold int) {
    for i := range c.Children {
        if c.Children[i] == nil {
            continue
        }
        switch child := c.Children[i].(type) {
        case ContainerBuffers:
            for bIdx := range child.Buffers {
                if child.Buffers[bIdx] > threshold {
                    child.Buffers[bIdx]--
                }
            }
            c.Children[i] = child
        case Container:
            child.DecrementIndicesAbove(threshold)
            c.Children[i] = child
        }
    }
}