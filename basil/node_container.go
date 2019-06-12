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
	run          func(*NodeContainer, []*util.WaitGroup) bool
	runCount     int
	waitGroups   []*util.WaitGroup
	mu           *sync.Mutex
}

// NewNodeContainer creates a new node container
func NewNodeContainer(
	ctx *EvalContext,
	node Node,
	dependencies map[ID]Container,
	run func(*NodeContainer, []*util.WaitGroup) bool,
) *NodeContainer {
	n := &NodeContainer{
		node:         node,
		dependencies: dependencies,
		missingDeps:  len(dependencies),
		run:          run,
		mu:           &sync.Mutex{},
	}

	for id := range dependencies {
		ctx.Subscribe(n, id)
	}

	return n
}

// ID returns with the node id
func (n *NodeContainer) ID() ID {
	return n.node.ID()
}

// Node returns with the node
func (n *NodeContainer) Node() Node {
	return n.node
}

// SetDependency stores the given container
// If all dependencies are set on the node then it will schedule the node for running.
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
		if n.run(n, n.waitGroups) {
			n.waitGroups = nil
		}
	}
}

// Run will schedule the node for running if it's ready. If it is then it returns true.
func (n *NodeContainer) Run() bool {
	n.mu.Lock()

	var run bool
	if n.missingDeps == 0 {
		if n.run(n, n.waitGroups) {
			run = true
			n.waitGroups = nil
		}
	}

	n.mu.Unlock()
	return run
}

// EvalContext returns with a new evaluation context
func (n *NodeContainer) EvalContext(ctx *EvalContext) *EvalContext {
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

func (n *NodeContainer) Close(ctx *EvalContext) {
	for id := range n.dependencies {
		ctx.Unsubscribe(n, id)
	}

	for _, wg := range n.waitGroups {
		wg.Done(errors.New("aborted"))
	}
}
