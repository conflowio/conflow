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

	"github.com/opsidian/parsley/parsley"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/job"
	"github.com/opsidian/basil/basil/parameter"
	"github.com/opsidian/basil/basil/schema"
	"github.com/opsidian/basil/util"
)

const (
	_ int64 = iota
	containerStatePending
	containerStateStart
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
	evalStage   basil.EvalStage
	state       int64
	stateChan   chan int64
	resultChan  chan basil.Container
	errChan     chan parsley.Error
	jobTracker  *job.Tracker
	jobID       int
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
	pending bool,
) *Container {
	var block basil.Block
	if value != nil {
		var ok bool
		if block, ok = value.(basil.Block); !ok {
			panic(fmt.Errorf("was expecting block, got %T", value))
		}
	}
	var state int64
	if pending {
		state = containerStatePending
	}

	return &Container{
		evalCtx: evalCtx,
		node:    node,
		parent:  parent,
		block:   block,
		wgs:     wgs,
		state:   state,
	}
}

// Node returns with the block node
func (c *Container) Node() basil.Node {
	return c.node
}

// SetJobID sets the job id
func (c *Container) JobName() basil.ID {
	return c.node.ID()
}

// JobID returns with the job id
func (c *Container) JobID() int {
	return c.jobID
}

// SetJobID sets the job id
func (c *Container) SetJobID(id int) {
	c.jobID = id
}

func (c *Container) Cancel() bool {
	return c.evalCtx.Cancel()
}

func (c *Container) Lightweight() bool {
	return true
}

func (c *Container) EvalStage() basil.EvalStage {
	return c.node.EvalStage()
}

func (c *Container) Pending() bool {
	return c.state == containerStatePending
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
	s := c.node.Interpreter().Schema()
	if p, ok := s.(schema.ObjectKind).GetProperties()[string(name)]; ok {
		if !IsBlockSchema(p) {
			return c.node.Interpreter().Param(c.block, name)
		}
	}

	return c.extraParams[name]
}

func (c *Container) Run() {
	defer func() {
		if c.parent != nil {
			c.parent.SetChild(c)
		}

		c.evalCtx.Cancel()
	}()

	if !c.evalCtx.Run() {
		return
	}

	if c.block == nil {
		c.block = c.node.Interpreter().CreateBlock(c.node.ID())
	}
	c.blockCtx = basil.NewBlockContext(c.evalCtx, c)

	c.jobTracker = job.NewTracker(c.node.ID(), c.evalCtx.JobScheduler(), c.evalCtx.Logger)
	c.stateChan = make(chan int64, 1)
	c.resultChan = make(chan basil.Container, 8)
	c.errChan = make(chan parsley.Error, 1)

	c.stateChan <- containerStateStart
	c.mainLoop()

	for _, container := range c.children {
		container.Close()
	}

	if c.jobTracker.Stop() > 0 {
		c.shutdownLoop()
	}
}

func (c *Container) mainLoop() {
	for {
		select {
		case state := <-c.stateChan:
			c.setState(state)

		case child := <-c.resultChan:
			if j, ok := child.(basil.Job); ok {
				c.jobTracker.Finished(j.JobID())
			}
			if nc, ok := child.(basil.NilContainer); ok && nc.Pending() {
				c.jobTracker.RemovePending(1)
			}

			err := c.setChild(child)
			child.Close()

			if err != nil {
				c.err = parsley.NewError(c.node.Pos(), err)
				c.logError(err).ID("childID", child.Node().ID()).Msg("evaluating child failed")
				c.setState(containerStateErrored)
				return
			}

			if c.jobTracker.ActiveJobCount() == 0 && c.state < containerStateFinished {
				c.stateChan <- containerStateNext
			}

		case err := <-c.errChan:
			c.setError(err)
			return

		case <-c.evalCtx.Done():
			switch c.evalCtx.Err() {
			case context.DeadlineExceeded:
				c.err = parsley.NewError(c.node.Pos(), fmt.Errorf("timeout reached in %s.%s", c.node.BlockType(), c.node.ID()))
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
			if j, ok := child.(basil.Job); ok {
				c.jobTracker.Finished(j.JobID())
			}
			if nc, ok := child.(basil.NilContainer); ok && nc.Pending() {
				c.jobTracker.RemovePending(1)
			}

			child.Close()

			if _, err := child.Value(); err != nil {
				c.logError(err).ID("childID", child.Node().ID()).Msg("evaluating child failed")
			}
		case err := <-c.errChan:
			c.logError(err).Err(err)
		case <-shutdownTimer.C:
			return
		}

		if c.jobTracker.ActiveJobCount() == 0 {
			return
		}
	}
}

func (c *Container) SetError(err parsley.Error) {
	go func() {
		c.errChan <- err
	}()
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
	case containerStateStart:
		c.children = make(map[basil.ID]*basil.NodeContainer, len(c.node.Children()))
		for _, node := range c.node.Children() {
			container, err := basil.NewNodeContainer(c.evalCtx, c, node, c.jobTracker)
			if err != nil {
				c.setError(err)
				return
			}

			c.children[node.ID()] = container
		}
		c.evalStage = basil.EvalStageInit
		if err := c.evaluateChildren(); err != nil {
			c.setError(err)
			return
		}
	case containerStateInit:
		c.updateEvalContext()

		if _, ok := c.block.(basil.BlockInitialiser); ok {
			err := c.jobTracker.ScheduleJob(newContainerStage(
				c,
				"init",
				basil.EvalStageInit,
				false,
				c.runInitStage,
			))
			if err != nil {
				c.setError(parsley.NewError(c.node.Pos(), fmt.Errorf("failed to run the init stage: %s", err)))
			}
		} else {
			c.stateChan <- containerStateNext
		}
	case containerStatePreMain:
		c.evalStage = basil.EvalStageMain
		if err := c.evaluateChildren(); err != nil {
			c.setError(err)
			return
		}
	case containerStateMain:
		if _, ok := c.block.(basil.BlockRunner); ok {
			err := c.jobTracker.ScheduleJob(newContainerStage(
				c,
				"main",
				basil.EvalStageMain,
				len(c.node.Generates()) > 0,
				c.runMainStage,
			))
			if err != nil {
				c.setError(parsley.NewError(c.node.Pos(), fmt.Errorf("failed to run the main stage: %s", err)))
			}
		} else {
			c.stateChan <- containerStateNext
		}
	case containerStatePreClose:
		c.evalStage = basil.EvalStageClose
		if err := c.evaluateChildren(); err != nil {
			c.setError(err)
			return
		}
	case containerStateClose:
		if _, ok := c.block.(basil.BlockCloser); ok {
			err := c.jobTracker.ScheduleJob(newContainerStage(
				c,
				"close",
				basil.EvalStageClose,
				false,
				c.runCloseStage,
			))
			if err != nil {
				c.setError(parsley.NewError(c.node.Pos(), fmt.Errorf("failed to run the close stage: %s", err)))
			}
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
	var pending, total int
	var isPending bool
	var err parsley.Error
	for _, child := range c.children {
		if child.Node().EvalStage() == c.evalStage {
			total++
			isPending, err = child.Run()
			if err != nil {
				break
			}
			if isPending {
				pending++
			}
		}
	}

	c.jobTracker.AddPending(pending)

	if err != nil {
		return err
	}

	if total == 0 {
		c.stateChan <- containerStateNext
		return nil
	}

	if pending == total {
		return parsley.NewErrorf(c.node.Pos(), "%q is deadlocked as no children could be evaluated, pending: %d", c.node.ID(), pending)
	}

	return nil
}

// SetChild sets a child container in a non-blocking way
func (c *Container) SetChild(container basil.Container) {
	go func() {
		c.resultChan <- container
	}()
}

func (c *Container) setChild(result basil.Container) parsley.Error {
	value, err := result.Value()
	if err != nil {
		return err
	}

	if value == nil {
		c.evalCtx.Publish(result)
		return nil
	}

	switch r := result.(type) {
	case *parameter.Container:
		node := r.Node().(basil.ParameterNode)
		name := node.Name()

		s := c.node.Interpreter().Schema()
		if _, ok := s.(schema.ObjectKind).GetProperties()[string(name)]; ok {
			if err := c.node.Interpreter().SetParam(c.block, node.Name(), value); err != nil {
				return parsley.NewError(r.Node().Pos(), err)
			}
		} else {
			if c.extraParams == nil {
				c.extraParams = make(map[basil.ID]interface{})
			}

			if _, exists := c.extraParams[name]; exists {
				panic(fmt.Errorf("%q parameter was already set on %q block", name, c.node.ID()))
			}

			c.extraParams[name] = value
		}

	case *Container:
		node := r.Node().(basil.BlockNode)
		name, p := getNameSchemaForChildBlock(c.Node().Schema().(*schema.Object), node)

		if err := c.node.Interpreter().SetBlock(c.block, name, value); err != nil {
			return parsley.NewError(r.Node().Pos(), err)
		}

		// Do not publish the empty generated node
		if p != nil && p.GetAnnotation("generated") == "true" && c.evalStage == basil.EvalStageInit {
			return nil
		}

	default:
		panic(fmt.Errorf("unknown container type: %T", result))
	}

	c.evalCtx.Publish(result)

	return nil
}

func (c *Container) PublishBlock(block basil.Block, f func() error) (bool, error) {
	b, ok := block.(basil.Identifiable)
	if !ok {
		return false, fmt.Errorf("%T block must implement the basil.Identifiable interface", block)
	}
	blockID := b.ID()

	nodeContainer, ok := c.children[blockID]
	propertyName := string(nodeContainer.Node().(basil.BlockNode).BlockType())
	if !ok || c.Node().Schema().(*schema.Object).Properties[propertyName].GetAnnotation("generated") != "true" {
		return false, fmt.Errorf("%q block does not exist or is not marked as generated", blockID)
	}

	if !c.evalCtx.HasSubscribers(blockID) {
		return false, nil
	}

	wg := &util.WaitGroup{}
	wg.Add(1)
	container, err := nodeContainer.CreateContainer(block, []basil.WaitGroup{wg})
	if err != nil {
		return false, err
	}
	if container == nil {
		return false, nil
	}

	if err := c.jobTracker.ScheduleJob(container); err != nil {
		return false, fmt.Errorf("failed to schedule %q: %s", blockID, err)
	}

	if f != nil {
		if err := f(); err != nil {
			return true, err
		}
	}

	select {
	case <-wg.Wait():
		return true, wg.Err()
	case <-c.evalCtx.Done():
		return true, fmt.Errorf("%q was aborted", c.node.ID())
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
		e = e.ID("id", c.node.ID()).
			ID("type", c.node.BlockType()).
			Uint8("state", uint8(c.state)).
			Int("active_jobs", c.jobTracker.ActiveJobCount())
	}
	return e
}

func (c *Container) logError(err error) basil.LogEvent {
	e := c.evalCtx.Logger.Error()
	if !e.Enabled() {
		return e
	}
	e.ID("id", c.node.ID()).
		ID("type", c.node.BlockType()).
		Uint8("state", uint8(c.state))

	if c.jobTracker != nil {
		e.Int("active_jobs", c.jobTracker.ActiveJobCount())
	}
	return e.Err(err)
}
