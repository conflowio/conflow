// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package basil

import (
	"errors"
	"sync"

	"github.com/opsidian/basil/util"
)

// NodeContainer wraps a node and registers the dependencies as they become available
type NodeContainer struct {
	node         Node
	dependencies map[ID]Container
	missingDeps  int
	ready        func(*NodeContainer, []*util.WaitGroup)
	runCount     int
	generated    bool
	waitGroups   []*util.WaitGroup
	mu           *sync.Mutex
}

// NewNodeContainer creates a new node container
func NewNodeContainer(
	node Node,
	dependencies map[ID]Container,
	ready func(*NodeContainer, []*util.WaitGroup),
) *NodeContainer {
	return &NodeContainer{
		node:         node,
		dependencies: dependencies,
		missingDeps:  len(dependencies),
		ready:        ready,
		generated:    node.Generated(),
		mu:           &sync.Mutex{},
	}
}

// ID returns with the node id
func (n *NodeContainer) ID() ID {
	return n.node.ID()
}

// SetDependency stores the given container
func (n *NodeContainer) SetDependency(c Container) {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.dependencies[c.ID()] == nil {
		n.missingDeps--
	}

	n.dependencies[c.ID()] = c

	for _, wg := range c.WaitGroups() {
		wg.Add(1)
		n.waitGroups = append(n.waitGroups, wg)
	}

	if n.missingDeps == 0 {
		n.ready(n, n.waitGroups)
		n.waitGroups = nil
	}
}

// Node returns with the node
func (n *NodeContainer) Node() Node {
	return n.node
}

// Ready returns true if the node doesn't have any unsatisfied dependencies
func (n *NodeContainer) Ready() bool {
	n.mu.Lock()
	ready := n.missingDeps == 0
	n.mu.Unlock()
	return ready
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

func (n *NodeContainer) Close() {
	n.mu.Lock()
	for _, wg := range n.waitGroups {
		wg.Done(errors.New("aborted"))
	}
	n.mu.Unlock()
}
