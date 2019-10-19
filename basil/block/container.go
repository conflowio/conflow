// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package block

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/job"
	"github.com/opsidian/basil/basil/parameter"
	"github.com/opsidian/basil/util"
	"github.com/opsidian/parsley/parsley"
)

const (
	containerStateWaiting int64 = iota
	containerStateRunning
	containerStatePreInit
	containerStateInit
	containerStatePreMain
	containerStateMain
	containerStatePreClose
	containerStateClose
	containerStateFinished
	containerStateSkipped
	containerStateErrored
	containerStateAborted
	containerStateCancelled
	containerStateNext
)

const (
	initJobID  = "@init"
	mainJobID  = "@main"
	closeJobID = "@close"
)

// ContainerGracefulTimeoutSec is the graceful timeout in seconds
var ContainerGracefulTimeoutSec = 10

var _ basil.BlockContainer = &Container{}

// Container is a block container
type Container struct {
	evalCtx     *basil.EvalContext
	blockCtx    basil.BlockContext
	node        basil.BlockNode
	parent      basil.BlockContainer
	block       basil.Block
	err         parsley.Error
	extraParams map[basil.ID]interface{}
	state       int64
	evalStage   basil.EvalStage
	stateChan   chan int64
	resultChan  chan basil.Container
	errChan     chan parsley.Error
	jobManager  *job.Manager
	jobID       basil.ID
	children    map[basil.ID]*basil.NodeContainer
	wgs         []*util.WaitGroup
}

// NewContainer creates a new block container instance
func NewContainer(
	evalCtx *basil.EvalContext,
	jobID basil.ID,
	node basil.BlockNode,
	parent basil.BlockContainer,
	block basil.Block,
	wgs []*util.WaitGroup,
) *Container {
	if block == nil {
		block = node.Interpreter().CreateBlock(node.ID())
	}

	return &Container{
		evalCtx:    evalCtx,
		node:       node,
		parent:     parent,
		block:      block,
		wgs:        wgs,
		evalStage:  basil.EvalStageInit,
		stateChan:  make(chan int64, 1),
		resultChan: make(chan basil.Container, 8),
		errChan:    make(chan parsley.Error, 1),
		jobManager: job.NewManager(node.ID(), evalCtx.Scheduler(), evalCtx.Logger),
		jobID:      jobID,
	}
}

// ID returns with the block id
func (c *Container) ID() basil.ID {
	return c.node.ID()
}

// JobID returns with the job id
func (c *Container) JobID() basil.ID {
	return c.jobID
}

// Node returns with the block node
func (c *Container) Node() basil.Node {
	return c.node
}

// Block returns with the block instance
func (c *Container) Block() basil.Block {
	return c.block
}

// Value returns with the block or the error if any occurred
// If the block was skipped then a nil value is returned
func (c *Container) Value() (interface{}, parsley.Error) {
	if c.err != nil {
		return nil, c.err
	}

	if c.state == containerStateSkipped {
		return nil, nil
	}

	return c.block, nil
}

// Param returns with the parameter value
func (c *Container) Param(name basil.ID) interface{} {
	if val, ok := c.extraParams[name]; ok {
		return val
	}

	return c.node.Interpreter().Param(c.block, name)
}

func (c *Container) Run() {
	if !atomic.CompareAndSwapInt64(&c.state, containerStateWaiting, containerStateRunning) {
		return
	}

	c.blockCtx = basil.NewBlockContext(c.evalCtx, c)

	c.stateChan <- containerStateNext
	c.mainLoop()

	c.evalCtx.Cancel()

	for _, container := range c.children {
		container.Close(c.evalCtx)
	}

	if c.jobManager.Stop() > 0 {
		c.shutdownLoop()
	}

	if c.parent != nil {
		c.parent.SetChild(c)
	}
}

func (c *Container) Cancel() bool {
	c.evalCtx.Cancel()
	return atomic.CompareAndSwapInt64(&c.state, containerStateWaiting, containerStateCancelled)
}

func (c *Container) Lightweight() bool {
	return true
}

func (c *Container) mainLoop() {
	for {
		select {
		case state := <-c.stateChan:
			c.setState(state)

		case child := <-c.resultChan:
			c.jobManager.Finished(child.JobID())

			err := c.setChild(child)
			child.Close()

			if err != nil {
				c.err = parsley.NewError(c.node.Pos(), err)
				c.logError(err).ID("childID", child.ID()).Msg("evaluating child failed")
				c.setState(containerStateErrored)
				return
			}

			if c.jobManager.ActiveJobCount() == 0 && c.state < containerStateFinished {
				c.stateChan <- containerStateNext
			}

		case err := <-c.errChan:
			c.err = err
			c.setState(containerStateErrored)
			return

		case <-c.evalCtx.Context.Done():
			switch c.evalCtx.Context.Err() {
			case context.DeadlineExceeded:
				c.err = parsley.NewError(c.node.Pos(), fmt.Errorf("timeout reached in %s.%s", c.Node().Type(), c.ID()))
			default:
				c.err = parsley.NewError(c.node.Pos(), errors.New("aborted"))
			}
			c.setState(containerStateAborted)
			return
		}

		if c.state >= containerStateFinished {
			return
		}
	}
}

func (c *Container) shutdownLoop() {
	c.debug().Msg("graceful shutdown")

	shutdownTimer := time.NewTimer(time.Duration(ContainerGracefulTimeoutSec) * time.Second)

	for {
		select {
		case <-c.stateChan:
		case child := <-c.resultChan:
			c.jobManager.Finished(child.JobID())

			child.Close()

			if _, err := child.Value(); err != nil {
				c.logError(err).ID("childID", child.ID()).Msg("evaluating child failed")
			}
		case err := <-c.errChan:
			c.logError(err).Err(err)
		case <-shutdownTimer.C:
			return
		}

		if c.jobManager.ActiveJobCount() == 0 {
			return
		}
	}
}

func (c *Container) setState(state int64) {
	if c.state < containerStateFinished {
		if state == containerStateNext {
			c.state++
		} else {
			c.state = state
		}
	}

	switch c.state {
	case containerStatePreInit:
		c.children = make(map[basil.ID]*basil.NodeContainer, len(c.node.Children()))
		for _, node := range c.node.Children() {
			c.children[node.ID()] = basil.NewNodeContainer(c.evalCtx, c, node)
		}
		c.evalStage = basil.EvalStageInit
		c.evaluateChildren()
	case containerStateInit:
		c.updateEvalContext()

		if _, ok := c.block.(basil.BlockInitialiser); ok {
			c.jobManager.Schedule(basil.NewJob(
				c.evalCtx,
				initJobID,
				false,
				c.runInitStage,
			), false)
		} else {
			c.stateChan <- containerStateNext
		}
	case containerStatePreMain:
		c.evalStage = basil.EvalStageMain
		c.evaluateChildren()
	case containerStateMain:
		if _, ok := c.block.(basil.BlockRunner); ok {
			c.jobManager.Schedule(basil.NewJob(
				c.evalCtx,
				mainJobID,
				len(c.node.Generates()) > 0,
				c.runMainStage,
			), false)
		} else {
			c.stateChan <- containerStateNext
		}
	case containerStatePreClose:
		c.evalStage = basil.EvalStageClose
		c.evaluateChildren()
	case containerStateClose:
		if _, ok := c.block.(basil.BlockCloser); ok {
			c.jobManager.Schedule(basil.NewJob(
				c.evalCtx,
				closeJobID,
				false,
				c.runCloseStage,
			), false)
		} else {
			c.stateChan <- containerStateNext
		}
	case containerStateFinished:
	case containerStateSkipped:
		c.debug().Msg("skipped")
		// TODO: notify block about skipped task
	case containerStateErrored:
		c.logError(c.err).Msg("failed")
		// TODO: notify block about error
	case containerStateAborted:
		c.logError(c.err).Msg("aborted")
		// TODO: notify block about abort
	default:
		panic("invalid container state")
	}
}

func (c *Container) updateEvalContext() {
	if b, ok := c.block.(basil.Contexter); ok {
		c.evalCtx.Context, c.evalCtx.Cancel = b.Context(c.evalCtx.Context)
	}

	if b, ok := c.block.(basil.UserContexter); ok {
		c.evalCtx.UserContext = b.UserContext(c.evalCtx.UserContext)
	}

	if b, ok := c.block.(basil.Loggerer); ok {
		c.evalCtx.Logger = b.Logger(c.evalCtx.Logger)
	}
}

func (c *Container) runInitStage() {
	defer func() {
		if r := recover(); r != nil {
			c.errChan <- parsley.NewErrorf(c.node.Pos(), "init stage panicked in %q: %s", c.ID(), r)
		}
	}()

	skipped, err := c.block.(basil.BlockInitialiser).Init(c.blockCtx)
	c.jobManager.Finished(initJobID)
	if err != nil {
		c.errChan <- parsley.NewError(c.node.Pos(), err)
		return
	}

	if !skipped {
		c.stateChan <- containerStateNext
	} else {
		c.stateChan <- containerStateSkipped
	}
}

func (c *Container) runMainStage() {
	defer func() {
		if r := recover(); r != nil {
			c.errChan <- parsley.NewErrorf(c.node.Pos(), "main stage panicked in %q: %s", c.ID(), r)
		}
	}()

	err := c.block.(basil.BlockRunner).Main(c.blockCtx)
	c.jobManager.Finished(mainJobID)
	if err != nil {
		c.errChan <- parsley.NewError(c.node.Pos(), err)
		return
	}

	c.stateChan <- containerStateNext
}

func (c *Container) runCloseStage() {
	defer func() {
		if r := recover(); r != nil {
			c.errChan <- parsley.NewErrorf(c.node.Pos(), "close stage panicked in %q: %s", c.ID(), r)
		}
	}()

	err := c.block.(basil.BlockCloser).Close(c.blockCtx)
	c.jobManager.Finished(closeJobID)
	if err != nil {
		c.errChan <- parsley.NewError(c.node.Pos(), err)
		return
	}

	c.stateChan <- containerStateNext
}

func (c *Container) evaluateChildren() {
	var pending, running int
	for _, container := range c.children {
		if container.Node().EvalStage() == c.evalStage {
			if container.Run() {
				running++
			} else {
				container.SetPending()
				pending++
			}
		}
	}

	if pending > 0 {
		c.jobManager.AddPending(pending)
		if running == 0 {
			c.errChan <- parsley.NewErrorf(c.node.Pos(), "%q is deadlocked as no children could be evaluated", c.ID())
			return
		}
	} else if running == 0 {
		c.stateChan <- containerStateNext
	}
}

func (c *Container) EvaluateChild(nodeContainer *basil.NodeContainer) bool {
	if nodeContainer.Node().EvalStage() != c.evalStage {
		return false
	}

	c.evaluateChild(nodeContainer, nil, nodeContainer.WaitGroups())

	return true
}

func (c *Container) evaluateChild(nodeContainer *basil.NodeContainer, value interface{}, wgs []*util.WaitGroup) {
	ctx := nodeContainer.CreateEvalContext(c.evalCtx)
	var container basil.Container
	jobID := c.jobManager.GenerateJobID(nodeContainer.ID())

	switch n := nodeContainer.Node().(type) {
	case basil.BlockNode:
		var block basil.Block
		if value != nil {
			var ok bool
			if block, ok = value.(basil.Block); !ok {
				panic(fmt.Errorf("was expecting block, got %T", value))
			}
		}
		container = NewContainer(ctx, jobID, n, c, block, wgs)
	case basil.ParameterNode:
		container = parameter.NewContainer(ctx, jobID, n, c)
	default:
		panic(fmt.Errorf("invalid node type: %T", n))
	}

	c.jobManager.Schedule(container, nodeContainer.RemovePending())
}

func (c *Container) SetChild(container basil.Container) {
	c.resultChan <- container
}

func (c *Container) setChild(result basil.Container) parsley.Error {
	value, err := result.Value()
	if err != nil {
		return err
	}

	if value == nil {
		return nil
	}

	switch r := result.(type) {
	case *parameter.Container:
		node := r.Node().(basil.ParameterNode)
		name := node.Name()

		if node.IsDeclaration() {
			if c.extraParams == nil {
				c.extraParams = make(map[basil.ID]interface{})
			}

			if _, exists := c.extraParams[name]; exists {
				panic(fmt.Errorf("%q parameter was already set on %q block", name, c.ID()))
			}

			c.extraParams[name] = value
		} else {
			if err := c.node.Interpreter().SetParam(c.block, node.Name(), value); err != nil {
				return parsley.NewError(r.Node().Pos(), err)
			}
		}
	case *Container:
		node := r.Node().(basil.BlockNode)
		name := node.BlockType()

		if err := c.node.Interpreter().SetBlock(c.block, name, value); err != nil {
			return parsley.NewError(r.Node().Pos(), err)
		}
	default:
		panic(fmt.Errorf("unknown container type: %T", result))
	}

	// Do not publish the empty generated node
	if result.Node().Generated() && c.evalStage == basil.EvalStageInit {
		return nil
	}

	c.evalCtx.Publish(result)

	return nil
}

func (c *Container) PublishBlock(block basil.Block) error {
	nodeContainer, ok := c.children[block.ID()]
	if !ok || !nodeContainer.Node().Generated() {
		return fmt.Errorf("%q block does not exist or is not marked as generated", block.ID())
	}

	wg := &util.WaitGroup{}
	wg.Add(1)
	c.evaluateChild(nodeContainer, block, []*util.WaitGroup{wg})

	select {
	case <-wg.Wait():
		return wg.Err()
	case <-c.evalCtx.Context.Done():
		return fmt.Errorf("%q was aborted", c.ID())
	}
}

// Close notifies all wait groups
func (c *Container) Close() {
	for _, wg := range c.wgs {
		wg.Done(c.err)
	}
}

// WaitGroups returns nil
func (c *Container) WaitGroups() []*util.WaitGroup {
	return c.wgs
}

func (c *Container) debug() basil.LogEvent {
	e := c.evalCtx.Logger.Debug()
	if e.Enabled() {
		e = e.ID("id", c.ID()).
			ID("type", c.node.BlockType()).
			Uint8("state", uint8(c.state)).
			Int("active_jobs", c.jobManager.ActiveJobCount())
	}
	return e
}

func (c *Container) logError(err error) basil.LogEvent {
	e := c.evalCtx.Logger.Error()
	if e.Enabled() {
		e = e.ID("id", c.ID()).
			ID("type", c.node.BlockType()).
			Uint8("state", uint8(c.state)).
			Int("active_jobs", c.jobManager.ActiveJobCount()).
			Err(err)
	}
	return e
}
