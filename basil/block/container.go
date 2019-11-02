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
	"time"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/job"
	"github.com/opsidian/basil/basil/parameter"
	"github.com/opsidian/basil/util"
	"github.com/opsidian/parsley/parsley"
)

const (
	_ int64 = iota
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
	containerStateNext
)

// ContainerGracefulTimeoutSec is the graceful timeout in seconds
var ContainerGracefulTimeoutSec = 10

var _ basil.BlockContainer = &Container{}

// Container is a block container
type Container struct {
	evalCtx     *basil.EvalContext
	node        basil.BlockNode
	blockCtx    basil.BlockContext
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
	wgs         []basil.WaitGroup
	retry       basil.Retryable
}

// NewContainer creates a new block container instance
func NewContainer(
	evalCtx *basil.EvalContext,
	node basil.BlockNode,
	parent basil.BlockContainer,
	value interface{},
	wgs []basil.WaitGroup,
) *Container {
	var block basil.Block
	if value == nil {
		block = node.Interpreter().CreateBlock(node.ID())
	} else {
		var ok bool
		if block, ok = value.(basil.Block); !ok {
			panic(fmt.Errorf("was expecting block, got %T", value))
		}
	}

	jobManager := job.NewManager(node.ID(), evalCtx.Scheduler, evalCtx.Logger)

	return &Container{
		evalCtx:    evalCtx,
		node:       node,
		parent:     parent,
		block:      block,
		wgs:        wgs,
		evalStage:  basil.EvalStageInit,
		jobManager: jobManager,
	}
}

// ID returns with the block id
func (c *Container) ID() basil.ID {
	return c.node.ID()
}

// SetJobID sets the job id
func (c *Container) JobName() basil.ID {
	return c.node.ID()
}

// JobID returns with the job id
func (c *Container) JobID() basil.ID {
	return c.jobID
}

// SetJobID sets the job id
func (c *Container) SetJobID(id basil.ID) {
	c.jobID = id
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
	defer func() {
		c.evalCtx.Cancel()

		if c.parent != nil {
			c.parent.SetChild(c)
		}
	}()

	if !c.evalCtx.Run() {
		return
	}

	c.blockCtx = basil.NewBlockContext(c.evalCtx, c)
	c.stateChan = make(chan int64, 1)
	c.resultChan = make(chan basil.Container, 8)
	c.errChan = make(chan parsley.Error, 1)

	c.stateChan <- containerStateNext
	c.mainLoop()

	for _, container := range c.children {
		container.Close()
	}

	if c.jobManager.Stop() > 0 {
		c.shutdownLoop()
	}
}

func (c *Container) Cancel() bool {
	return c.evalCtx.Cancel()
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
				c.logError(err).ID("childID", child.Node().ID()).Msg("evaluating child failed")
				c.setState(containerStateErrored)
				return
			}

			if c.jobManager.ActiveJobCount() == 0 && c.state < containerStateFinished {
				c.stateChan <- containerStateNext
			}

		case err := <-c.errChan:
			c.setError(err)
			return

		case <-c.evalCtx.Done():
			switch c.evalCtx.Err() {
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
				c.logError(err).ID("childID", child.Node().ID()).Msg("evaluating child failed")
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

func (c *Container) SetError(err parsley.Error) {
	c.errChan <- err
}

func (c *Container) setError(err parsley.Error) {
	c.err = err
	c.setState(containerStateErrored)
}

func (c *Container) setState(state int64) {
	if c.state < containerStateFinished {
		if state == containerStateNext {
			c.state = c.state + 1
		} else {
			c.state = state
		}
	}

	switch c.state {
	case containerStatePreInit:
		c.children = make(map[basil.ID]*basil.NodeContainer, len(c.node.Children()))
		var err parsley.Error
		for _, node := range c.node.Children() {
			c.children[node.ID()], err = basil.NewNodeContainer(c.evalCtx, c, node)
			if err != nil {
				c.setError(err)
				return
			}
		}
		c.evalStage = basil.EvalStageInit
		if err := c.evaluateChildren(); err != nil {
			c.setError(err)
		}
	case containerStateInit:
		c.updateEvalContext()

		if _, ok := c.block.(basil.BlockInitialiser); ok {
			c.jobManager.ScheduleJob(newContainerStage(
				c,
				"init",
				false,
				c.runInitStage,
			), false)
		} else {
			c.stateChan <- containerStateNext
		}
	case containerStatePreMain:
		c.evalStage = basil.EvalStageMain
		if err := c.evaluateChildren(); err != nil {
			c.setError(err)
		}
	case containerStateMain:
		if _, ok := c.block.(basil.BlockRunner); ok {
			c.jobManager.ScheduleJob(newContainerStage(
				c,
				"main",
				len(c.node.Generates()) > 0,
				c.runMainStage,
			), false)
		} else {
			c.stateChan <- containerStateNext
		}
	case containerStatePreClose:
		c.evalStage = basil.EvalStageClose
		if err := c.evaluateChildren(); err != nil {
			c.setError(err)
		}
	case containerStateClose:
		if _, ok := c.block.(basil.BlockCloser); ok {
			c.jobManager.ScheduleJob(newContainerStage(
				c,
				"close",
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
	if b, ok := c.block.(basil.UserContexter); ok {
		c.evalCtx.UserContext = b.UserContext(c.evalCtx.UserContext)
	}

	if b, ok := c.block.(basil.Loggerer); ok {
		c.evalCtx.Logger = b.Logger(c.evalCtx.Logger)
	}
}

func (c *Container) runInitStage() (int64, error) {
	skipped, err := c.block.(basil.BlockInitialiser).Init(c.blockCtx)
	if err != nil {
		return 0, err
	}
	if skipped {
		return containerStateSkipped, nil
	}
	return containerStateNext, nil
}

func (c *Container) runMainStage() (int64, error) {
	return containerStateNext, c.block.(basil.BlockRunner).Main(c.blockCtx)
}

func (c *Container) runCloseStage() (int64, error) {
	return containerStateNext, c.block.(basil.BlockCloser).Close(c.blockCtx)
}

func (c *Container) evaluateChildren() parsley.Error {
	var pending, running int
	var err parsley.Error
	var ran bool
	for _, child := range c.children {
		if child.Node().EvalStage() == c.evalStage {
			ran, err = child.Run()
			if err != nil {
				break
			}
			if ran {
				running++
			} else {
				pending++
			}
		}
	}

	c.jobManager.AddPending(pending)

	if err != nil {
		return err
	}

	if pending > 0 && running == 0 {
		return parsley.NewErrorf(c.node.Pos(), "%q is deadlocked as no children could be evaluated", c.ID())
	}

	if pending+running == 0 {
		c.stateChan <- containerStateNext
	}

	return nil
}

func (c *Container) ScheduleChild(child basil.Container, pending bool) bool {
	if child.Node().EvalStage() != c.evalStage {
		return false
	}
	c.jobManager.ScheduleJob(child, pending)
	return true
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
	container, err := nodeContainer.CreateContainer(block, []basil.WaitGroup{wg})
	if err != nil {
		return err
	}
	if container == nil {
		return nil
	}

	c.jobManager.ScheduleJob(container, false)

	select {
	case <-wg.Wait():
		return wg.Err()
	case <-c.evalCtx.Done():
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
func (c *Container) WaitGroups() []basil.WaitGroup {
	return c.wgs
}

func (c *Container) RetryCount() int {
	if c.retry != nil {
		return c.retry.RetryCount()
	}
	return 0
}

func (c *Container) RetryDelay(retryIndex int) time.Duration {
	if c.retry != nil {
		return c.retry.RetryDelay(retryIndex)
	}
	return 0
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
	if !e.Enabled() {
		return e
	}
	e.ID("id", c.ID()).
		ID("type", c.node.BlockType()).
		Uint8("state", uint8(c.state))

	if c.jobManager != nil {
		e.Int("active_jobs", c.jobManager.ActiveJobCount())
	}
	return e.Err(err)
}
