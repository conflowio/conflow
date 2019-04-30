package dependency

type stack []*node

func (s *stack) Push(n *node) {
	*s = append(*s, n)
	n.OnStack = true
}

func (s *stack) Pop() *node {
	l := len(*s)
	n := (*s)[l-1]
	n.OnStack = false
	*s = (*s)[:l-1]
	return n
}
