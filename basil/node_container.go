package basil

import "sync"

// NodeContainer wraps a node and registers the dependencies as they become available
type NodeContainer struct {
	node         Node
	dependencies map[ID]Container
	missingDeps  int
	ready        func(*NodeContainer)
	runCount     int
	generated    bool
	waitGroups   []*sync.WaitGroup
}

// NewNodeContainer creates a new node container
func NewNodeContainer(
	node Node,
	dependencies map[ID]Container,
	ready func(*NodeContainer),
) *NodeContainer {
	return &NodeContainer{
		node:         node,
		dependencies: dependencies,
		missingDeps:  len(dependencies),
		ready:        ready,
		generated:    node.Generated(),
	}
}

// ID returns with the node id
func (n *NodeContainer) ID() ID {
	return n.node.ID()
}

// SetDependency stores the given container
func (n *NodeContainer) SetDependency(c Container) {
	if n.dependencies[c.ID()] == nil {
		n.missingDeps--
	}

	n.dependencies[c.ID()] = c

	for _, wg := range c.WaitGroups() {
		n.waitGroups = append(n.waitGroups, wg)
	}

	if n.missingDeps == 0 {
		n.ready(n)
		n.waitGroups = nil
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

// Generated returns true if the node is generated (either directly or indirectly)
func (n *NodeContainer) Generated() bool {
	return n.generated
}

// EvalContext returns with a new evaluation context
func (n *NodeContainer) EvalContext(ctx EvalContext) EvalContext {
	dependencies := make(map[ID]BlockContainer, len(n.dependencies))
	for id, cont := range n.dependencies {
		switch c := cont.(type) {
		case BlockContainer:
			dependencies[id] = c
		case ParameterContainer:
			dependencies[c.BlockContainer().ID()] = c.BlockContainer()
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

// WaitGroups returns with the registered wait groups
func (n *NodeContainer) WaitGroups() []*sync.WaitGroup {
	return n.waitGroups
}
