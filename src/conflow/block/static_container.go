// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package block

import (
	"context"
	"fmt"

	"github.com/conflowio/parsley/parsley"

	"github.com/conflowio/conflow/src/conflow"
	"github.com/conflowio/conflow/src/conflow/parameter"
	"github.com/conflowio/conflow/src/schema"
)

// StaticContainer is a container for blocks where there is no dynamic child evaluation required
type StaticContainer struct {
	evalCtx *conflow.EvalContext
	node    conflow.BlockNode
	block   conflow.Block
	err     parsley.Error
	jobID   int
}

// NewStaticContainer creates a new static block container instance
func NewStaticContainer(
	evalCtx *conflow.EvalContext,
	node conflow.BlockNode,
) *StaticContainer {
	return &StaticContainer{
		evalCtx: evalCtx,
		node:    node,
	}
}

// Node returns with the block node
func (s *StaticContainer) Node() conflow.Node {
	return s.node
}

func (s *StaticContainer) JobName() conflow.ID {
	return s.node.ID()
}

func (s *StaticContainer) JobID() int {
	return s.jobID
}

func (s *StaticContainer) SetJobID(jobID int) {
	s.jobID = jobID
}

func (s *StaticContainer) Cancel() bool {
	return s.evalCtx.Cancel()
}

func (s *StaticContainer) EvalStage() conflow.EvalStage {
	return s.node.EvalStage()
}

func (s *StaticContainer) Run() {
	defer s.evalCtx.Cancel()

	if !s.evalCtx.Run() {
		return
	}

	s.block = s.node.Interpreter().CreateBlock(s.node.ID(), conflow.NewBlockContext(s.evalCtx, nil))
	for _, child := range s.node.Children() {
		if err := s.evaluateChild(child); err != nil {
			s.err = err
			return
		}
	}
}

func (s *StaticContainer) Lightweight() bool {
	return true
}

func (s *StaticContainer) Pending() bool {
	return false
}

// Value returns with the block or the error if any occurred
// If the block was skipped then a nil value is returned
func (s *StaticContainer) Value() (interface{}, parsley.Error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.block, nil
}

func (s *StaticContainer) createContainer(node conflow.Node) conflow.JobContainer {
	ctx, cancel := context.WithCancel(context.Background())
	childCtx := s.evalCtx.New(ctx, cancel, nil)
	switch n := node.(type) {
	case conflow.BlockNode:
		return NewStaticContainer(childCtx, n)
	case conflow.ParameterNode:
		return parameter.NewContainer(childCtx, n, nil, nil, nil, false)
	default:
		panic(fmt.Errorf("unknown node type: %T", node))
	}
}

func (s *StaticContainer) evaluateChild(node conflow.Node) parsley.Error {
	container := s.createContainer(node)
	container.Run()
	value, err := container.Value()
	if err != nil {
		return err
	}

	if value == nil {
		return nil
	}

	switch r := container.(type) {
	case conflow.ParameterContainer:
		node := r.Node().(conflow.ParameterNode)
		if err := s.node.Interpreter().SetParam(s.block, node.Name(), value); err != nil {
			return parsley.NewError(r.Node().Pos(), err)
		}
	case *StaticContainer:
		node := r.Node().(conflow.BlockNode)
		name, _ := getNameSchemaForChildBlock(s.Node().Schema().(*schema.Object), node)

		if err := s.node.Interpreter().SetBlock(s.block, name, value); err != nil {
			return parsley.NewError(r.Node().Pos(), err)
		}
	default:
		panic(fmt.Errorf("unknown container type: %T", container))
	}

	return nil
}

func (s *StaticContainer) WaitGroups() []conflow.WaitGroup {
	return nil
}

func (s *StaticContainer) Close() {
}
