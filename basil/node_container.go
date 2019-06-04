package basil

// NodeContainer wraps a node and registers the dependencies as they become available
type NodeContainer struct {
	node         Node
	dependencies map[ID]Container
	missingDeps  int
	ready        func(c *NodeContainer)
	runCount     int
}

// NewNodeContainer creates a new node container
func NewNodeContainer(node Node, dependencies map[ID]Container, ready func(c *NodeContainer)) *NodeContainer {
	return &NodeContainer{
		node:         node,
		dependencies: dependencies,
		missingDeps:  len(dependencies),
		ready:        ready,
	}
}

// ID returns with the node id
func (n *NodeContainer) ID() ID {
	return n.node.ID()
}

// SetDependency stores the given container
// It returns true if the node has all dependencies satisfied
func (n *NodeContainer) SetDependency(c Container) {
	if n.dependencies[c.ID()] == nil {
		n.missingDeps--
	}

	n.dependencies[c.ID()] = c

	if n.missingDeps == 0 {
		n.ready(n)
	}
}

// Node returns with the node
func (n *NodeContainer) Node() Node {
	return n.node
}

// Ready returns true if the node doesn't have any unsatisfied dependencies
func (n *NodeContainer) Ready() bool {
	return n.missingDeps == 0
}

// EvalContext returns with a new evaluation context
func (n *NodeContainer) EvalContext(ctx EvalContext) EvalContext {
	dependencies := make(map[ID]BlockContainer, len(n.dependencies))
	for id, cont := range n.dependencies {
		switch c := cont.(type) {
		case BlockContainer:
			dependencies[id] = c
		case ParameterContainer:
			dependencies[c.Parent().ID()] = c.Parent()
		}
	}

	return ctx.WithDependencies(dependencies)
}

// RunCount will return with the run count
func (n *NodeContainer) RunCount() int {
	return n.runCount
}

// IncRunCount will increase the run count by one
func (n *NodeContainer) IncRunCount() {
	n.runCount++
}
