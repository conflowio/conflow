// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package basil

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/opsidian/parsley/parsley"
)

// NodeContainer wraps a node and registers the dependencies as they become available
type NodeContainer struct {
	ctx           *EvalContext
	parent        BlockContainer
	node          Node
	runtimeConfig RuntimeConfig
	dependencies  map[ID]Container
	missingDeps   int
	pending       bool
	waitGroups    []WaitGroup
	mu            *sync.Mutex
}

// NewNodeContainer creates a new node container
func NewNodeContainer(
	ctx *EvalContext,
	parent BlockContainer,
	node Node,
) (*NodeContainer, parsley.Error) {
	dependencies := make(map[ID]Container, len(node.Dependencies()))
	parentDependencies := parent.Node().Dependencies()
	for _, v := range node.Dependencies() {
		if _, ok := parentDependencies[v.ID()]; ok {
			continue
		}
		if parent.Node().ID() == v.ParentID() {
			dependencies[v.ID()] = nil
		} else {
			dependencies[v.ParentID()] = nil
		}
	}

	n := &NodeContainer{
		ctx:          ctx,
		parent:       parent,
		node:         node,
		dependencies: dependencies,
		missingDeps:  len(dependencies),
		mu:           &sync.Mutex{},
	}

	var err parsley.Error
	if n.runtimeConfig, err = n.evaluateDirectives(EvalStageResolve); err != nil {
		return nil, err
	}

	for id := range n.dependencies {
		ctx.Subscribe(n, id)
	}

	return n, nil
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

	isTrigger := n.runtimeConfig.IsTrigger(c.Node().ID())
	run := isTrigger

	val, ok := n.dependencies[c.Node().ID()]
	if !ok {
		panic(fmt.Errorf("unknown dependency: %s", c.Node().ID()))
	}

	if val == nil {
		n.missingDeps--
		if n.missingDeps == 0 {
			run = true
		}
	}
	n.dependencies[c.Node().ID()] = c

	if isTrigger {
		for _, wg := range c.WaitGroups() {
			wg.Add(1)
			n.waitGroups = append(n.waitGroups, wg)
		}
	}

	if run {
		_, err := n.run()
		if err != nil {
			n.parent.SetError(err)
		}
	}
}

// Run will schedule the node if it's ready. If it was then it returns true.
// If the node is not ready, then it will set to a pending status
func (n *NodeContainer) Run() (bool, parsley.Error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	ran, err := n.run()
	if err != nil {
		return false, err
	}
	if !ran {
		n.pending = true
	}
	return ran, nil
}

func (n *NodeContainer) run() (bool, parsley.Error) {
	if n.missingDeps == 0 {
		container, err := n.CreateContainer(nil, n.waitGroups)
		if err != nil {
			return false, err
		}
		if container == nil || n.parent.ScheduleChild(container, n.pending) {
			n.waitGroups = nil
			n.pending = false
			return true, nil
		}
	}
	return false, nil
}

func (n *NodeContainer) CreateContainer(value interface{}, wgs []WaitGroup) (Container, parsley.Error) {
	runtimeConfig, err := n.evaluateDirectives(EvalStageInit)
	if err != nil {
		return nil, err
	}

	if runtimeConfig.Skip {
		return nil, nil
	}

	ctx := n.createEvalContext(runtimeConfig.Timeout)
	return n.node.CreateContainer(ctx, n.parent, value, wgs), nil
}

// CreateEvalContext returns with a new evaluation context
func (n *NodeContainer) createEvalContext(timeout time.Duration) *EvalContext {
	dependencies := make(map[ID]BlockContainer, len(n.dependencies))
	for id, cont := range n.dependencies {
		switch c := cont.(type) {
		case BlockContainer:
			dependencies[id] = c
		case ParameterContainer:
			dependencies[c.BlockContainer().Node().ID()] = c.BlockContainer()
		default:
			panic(fmt.Errorf("Unexpected dependency type: %T", cont))
		}
	}
	var ctx context.Context
	var cancel context.CancelFunc
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}

	return n.ctx.New(ctx, cancel, dependencies)
}

func (n *NodeContainer) Close() {
	for id := range n.dependencies {
		n.ctx.Unsubscribe(n, id)
	}

	for _, wg := range n.waitGroups {
		wg.Done(errors.New("aborted"))
	}
}

func (n *NodeContainer) evaluateDirectives(evalStage EvalStage) (RuntimeConfig, parsley.Error) {
	r := n.runtimeConfig
	for _, d := range n.node.Directives() {
		if d.EvalStage() != evalStage {
			continue
		}

		ctx, cancel := context.WithCancel(n.ctx)
		evalCtx := n.ctx.New(ctx, cancel, nil)
		directive, err := d.Value(evalCtx)
		if err != nil {
			return RuntimeConfig{}, err
		}

		r = r.Merge(directive.(Directive).RuntimeConfig())
	}
	return r, nil
}
