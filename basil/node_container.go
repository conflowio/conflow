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

	"github.com/opsidian/basil/util"
)

// NodeContainer wraps a node and registers the dependencies as they become available
type NodeContainer struct {
	ctx           *EvalContext
	parent        BlockContainer
	node          Node
	scheduler     JobScheduler
	runtimeConfig RuntimeConfig
	dependencies  map[ID]Container
	missingDeps   int
	nilDeps       int
	pending       bool
	waitGroups    []WaitGroup
	mu            *sync.Mutex
}

// NewNodeContainer creates a new node container
func NewNodeContainer(
	ctx *EvalContext,
	parent BlockContainer,
	node Node,
	scheduler JobScheduler,
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
		scheduler:    scheduler,
		dependencies: dependencies,
		missingDeps:  len(dependencies),
		mu:           &sync.Mutex{},
	}

	var err parsley.Error

	if err = n.evaluateDirectives(EvalStageResolve); err != nil {
		return nil, err
	}

	for id := range n.dependencies {
		ctx.Subscribe(n, id)
	}

	return n, nil
}

// Node returns with the node
func (n *NodeContainer) Node() Node {
	return n.node
}

// SetDependency stores the given container
// If all dependencies are set on the node then it will schedule the node for running.
func (n *NodeContainer) SetDependency(dep Container) {
	n.mu.Lock()
	defer n.mu.Unlock()

	prevDep, ok := n.dependencies[dep.Node().ID()]
	if !ok {
		panic(fmt.Errorf("unknown dependency: %s", dep.Node().ID()))
	}
	n.dependencies[dep.Node().ID()] = dep

	isTrigger := n.runtimeConfig.IsTrigger(dep.Node().ID())
	run := isTrigger && n.missingDeps == 0

	if prevDep == nil {
		n.missingDeps--
		if n.missingDeps == 0 {
			run = true
		}
	}

	isSkipped := n.nilDeps > 0
	n.calculateNilDeps(prevDep, dep)
	if n.nilDeps > 0 {
		if isSkipped {
			n.setNilChild()
		}
		return
	}

	if isTrigger {
		for _, wg := range dep.WaitGroups() {
			wg.Add(1)
			n.waitGroups = append(n.waitGroups, wg)
		}
	}

	if run && n.parent.EvalStage() == n.node.EvalStage() {
		if err := n.run(); err != nil {
			n.parent.SetError(err)
		}
	}
}

func (n *NodeContainer) calculateNilDeps(prevDep, newDep Container) {
	var wasNil bool
	if prevDep != nil {
		prevVal, _ := prevDep.Value()
		wasNil = prevVal == nil
	}
	newVal, _ := newDep.Value()
	isNil := newVal == nil
	if !wasNil && isNil {
		n.nilDeps++
	} else if wasNil && !isNil {
		n.nilDeps--
	}
}

func (n *NodeContainer) setNilChild() {
	nilContainer := NewNilContainer(n.node, n.waitGroups, n.pending)
	n.parent.SetChild(nilContainer)
	n.waitGroups = nil
	n.pending = false
}

// Run will schedule the node if it's ready.
// If the node is not ready, then it will return with pending true
func (n *NodeContainer) Run() (pending bool, err parsley.Error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.missingDeps > 0 {
		n.pending = true
		return true, nil
	}

	if err := n.run(); err != nil {
		return false, err
	}

	return false, nil
}

func (n *NodeContainer) run() parsley.Error {
	var value interface{}
	if paramNode, ok := n.node.(ParameterNode); ok {
		if override, ok := n.ctx.InputParams[paramNode.Name()]; ok {
			value = override
		}
	}

	container, err := n.CreateContainer(value, n.waitGroups)
	if err != nil {
		return err
	}

	if container == nil {
		n.setNilChild()
		return nil
	}

	if err := n.scheduler.ScheduleJob(container); err != nil {
		return parsley.NewError(0, err)
	}

	n.waitGroups = nil
	n.pending = false

	return nil
}

func (n *NodeContainer) CreateContainer(value interface{}, wgs []WaitGroup) (JobContainer, parsley.Error) {
	if err := n.evaluateDirectives(EvalStageInit); err != nil {
		return nil, err
	}

	if util.BoolValue(n.runtimeConfig.Skip) {
		return nil, nil
	}

	ctx := n.createEvalContext(util.TimeDurationValue(n.runtimeConfig.Timeout))
	return n.node.CreateContainer(ctx, n.runtimeConfig, n.parent, value, wgs, n.pending), nil
}

// CreateEvalContext returns with a new evaluation context
func (n *NodeContainer) createEvalContext(timeout time.Duration) *EvalContext {
	dependencies := make(map[ID]BlockContainer, len(n.dependencies))
	for id, cont := range n.dependencies {
		if cont == nil {
			continue
		}
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

func (n *NodeContainer) evaluateDirectives(evalStage EvalStage) parsley.Error {
	for _, d := range n.node.Directives() {
		if d.EvalStage() != evalStage {
			continue
		}

		directive, err := parsley.EvaluateNode(n.createEvalContext(0), d)
		if err != nil {
			return err
		}

		opt, ok := directive.(RuntimeConfigOption)
		if ok {
			opt.ApplyToRuntimeConfig(&n.runtimeConfig)
		}
	}
	return nil
}
