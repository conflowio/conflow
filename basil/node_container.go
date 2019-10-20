// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package basil

import (
	"errors"
	"sync"
	"sync/atomic"

	"github.com/opsidian/basil/util"
)

// NodeContainer wraps a node and registers the dependencies as they become available
type NodeContainer struct {
	parent       BlockContainer
	node         Node
	dependencies map[ID]Container
	missingDeps  int
	pending      uint64
	waitGroups   []*util.WaitGroup
	mu           *sync.Mutex
}

// NewNodeContainer creates a new node container
func NewNodeContainer(
	ctx *EvalContext,
	parent BlockContainer,
	node Node,
) *NodeContainer {
	dependencies := make(map[ID]Container, len(node.Dependencies()))
	for _, v := range node.Dependencies() {
		if _, ok := parent.Node().Dependencies()[v.ID()]; ok {
			continue
		}
		if parent.ID() == v.ParentID() {
			dependencies[v.ID()] = nil
		} else {
			dependencies[v.ParentID()] = nil
		}
	}

	n := &NodeContainer{
		parent:       parent,
		node:         node,
		dependencies: dependencies,
		missingDeps:  len(dependencies),
		mu:           &sync.Mutex{},
	}

	for id := range n.dependencies {
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

func (n *NodeContainer) WaitGroups() []*util.WaitGroup {
	return n.waitGroups
}

// SetDependency stores the given container
// If all dependencies are set on the node then it will schedule the node for running.
func (n *NodeContainer) SetDependency(c Container) {
	n.mu.Lock()
	defer n.mu.Unlock()

	trigger := c.ID()

	if n.dependencies[c.ID()] == nil {
		n.missingDeps--
		if n.missingDeps == 0 {
			trigger = ""
		}
	}
	n.dependencies[c.ID()] = c

	for _, wg := range c.WaitGroups() {
		wg.Add(1)
		n.waitGroups = append(n.waitGroups, wg)
	}

	n.run(trigger)
}

// Run will schedule the node for running if it's ready. If it is then it returns true.
func (n *NodeContainer) Run() bool {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.run("")
}

func (n *NodeContainer) run(trigger ID) bool {
	if n.missingDeps == 0 {
		if n.parent.EvaluateChild(n, trigger) {
			n.waitGroups = nil
			return true
		}
	}
	return false
}

// CreateEvalContext returns with a new evaluation context
func (n *NodeContainer) CreateEvalContext(ctx *EvalContext) *EvalContext {
	dependencies := make(map[ID]BlockContainer, len(n.dependencies))
	for id, cont := range n.dependencies {
		switch c := cont.(type) {
		case BlockContainer:
			dependencies[id] = c
		case ParameterContainer:
			dependencies[c.BlockContainer().ID()] = c.BlockContainer()
		}
	}

	return ctx.New(dependencies)
}

func (n *NodeContainer) Close(ctx *EvalContext) {
	for id := range n.dependencies {
		ctx.Unsubscribe(n, id)
	}

	for _, wg := range n.waitGroups {
		wg.Done(errors.New("aborted"))
	}
}

func (n *NodeContainer) SetPending() {
	atomic.StoreUint64(&n.pending, 1)
}

func (n *NodeContainer) RemovePending() bool {
	return atomic.CompareAndSwapUint64(&n.pending, 1, 0)
}
