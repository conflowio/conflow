package basil

// NodeContainer wraps a node and registers the dependencies as they become available
type NodeContainer struct {
	node         Node
	dependencies map[ID]Container
	missingDeps  int
}

// NewNodeContainer creates a new node container
func NewNodeContainer(node Node, dependencies map[ID]Container) *NodeContainer {
	return &NodeContainer{
		node:         node,
		dependencies: dependencies,
		missingDeps:  len(dependencies),
	}
}

// SetDependency stores the given container
// It returns true if the node has all dependencies satisfied
func (n *NodeContainer) SetDependency(c Container) bool {
	if n.dependencies[c.ID()] == nil {
		n.missingDeps--
	}

	n.dependencies[c.ID()] = c

	return n.missingDeps == 0
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
