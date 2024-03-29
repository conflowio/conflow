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

	"github.com/conflowio/parsley/parsley"

	"github.com/conflowio/conflow/pkg/conflow"
	"github.com/conflowio/conflow/pkg/conflow/annotations"
	"github.com/conflowio/conflow/pkg/conflow/job"
	"github.com/conflowio/conflow/pkg/conflow/parameter"
	"github.com/conflowio/conflow/pkg/schema"
	"github.com/conflowio/conflow/pkg/util"
	"github.com/conflowio/conflow/pkg/util/ptr"
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

var _ conflow.BlockContainer = &Container{}

// This will result in approximately 10 tries in the first minute (assuming no meaningful delay in the job queue of course)
var defaultBackoff = job.ExponentialRetryBackoff(1.57, 1*time.Second, 15*time.Minute)

// Container is a block container
type Container struct {
	evalCtx       *conflow.EvalContext
	node          conflow.BlockNode
	parent        conflow.BlockContainer
	block         conflow.Block
	err           parsley.Error
	extraParams   map[conflow.ID]interface{}
	evalStage     conflow.EvalStage
	state         int64
	stateChan     chan int64
	resultChan    chan conflow.Container
	errChan       chan parsley.Error
	jobTracker    *job.Tracker
	jobID         int
	children      map[conflow.ID]*conflow.NodeContainer
	wgs           []conflow.WaitGroup
	runtimeConfig conflow.RuntimeConfig
}

// NewContainer creates a new block container instance
func NewContainer(
	evalCtx *conflow.EvalContext,
	runtimeConfig conflow.RuntimeConfig,
	node conflow.BlockNode,
	parent conflow.BlockContainer,
	value interface{},
	wgs []conflow.WaitGroup,
	pending bool,
) *Container {
	var block conflow.Block
	if value != nil {
		var ok bool
		if block, ok = value.(conflow.Block); !ok {
			panic(fmt.Errorf("was expecting block, got %T", value))
		}
	}
	var state int64
	if pending {
		state = containerStatePending
	}

	return &Container{
		evalCtx:       evalCtx,
		runtimeConfig: runtimeConfig,
		node:          node,
		parent:        parent,
		block:         block,
		wgs:           wgs,
		state:         state,
	}
}

// Node returns with the block node
func (c *Container) Node() conflow.Node {
	return c.node
}

// SetJobID sets the job id
func (c *Container) JobName() conflow.ID {
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

func (c *Container) EvalStage() conflow.EvalStage {
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
func (c *Container) Param(name conflow.ID) interface{} {
	s := c.node.Interpreter().Schema()
	if p, ok := s.(*schema.Object).PropertyByParameterName(string(name)); ok {
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
		c.block = c.node.Interpreter().CreateBlock(c.node.ID(), conflow.NewBlockContext(c.evalCtx, c))
	}

	c.jobTracker = job.NewTracker(c.node.ID(), c.evalCtx.JobScheduler(), c.evalCtx.Logger, defaultBackoff)
	c.stateChan = make(chan int64, 1)
	c.resultChan = make(chan conflow.Container, 8)
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
			if j, ok := child.(conflow.Job); ok {
				c.jobTracker.Succeeded(j.JobID())
			}
			if nc, ok := child.(conflow.NilContainer); ok && nc.Pending() {
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
			if j, ok := child.(conflow.Job); ok {
				c.jobTracker.Succeeded(j.JobID())
			}
			if nc, ok := child.(conflow.NilContainer); ok && nc.Pending() {
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
		c.children = make(map[conflow.ID]*conflow.NodeContainer, len(c.node.Children()))
		for _, node := range c.node.Children() {
			container, err := conflow.NewNodeContainer(c.evalCtx, c, node, c.jobTracker)
			if err != nil {
				c.setError(err)
				return
			}

			c.children[node.ID()] = container
		}
		c.evalStage = conflow.EvalStageInit
		if err := c.evaluateChildren(); err != nil {
			c.setError(err)
			return
		}
	case containerStateInit:
		c.updateEvalContext()

		if _, ok := c.block.(conflow.BlockInitialiser); ok {
			err := c.jobTracker.ScheduleJob(newContainerStage(
				c,
				"init",
				conflow.EvalStageInit,
				false,
				conflow.RetryConfig{},
				c.runInitStage,
			))
			if err != nil {
				c.setError(parsley.NewError(c.node.Pos(), fmt.Errorf("failed to run the init stage: %s", err)))
			}
		} else {
			c.stateChan <- containerStateNext
		}
	case containerStatePreMain:
		c.evalStage = conflow.EvalStageMain
		if err := c.evaluateChildren(); err != nil {
			c.setError(err)
			return
		}
	case containerStateMain:
		if _, ok := c.block.(conflow.BlockRunner); ok {
			retryConfig := conflow.RetryConfig{Limit: -1}
			if c.runtimeConfig.RetryConfig != nil {
				retryConfig = *c.runtimeConfig.RetryConfig
			}

			err := c.jobTracker.ScheduleJob(newContainerStage(
				c,
				"main",
				conflow.EvalStageMain,
				len(c.node.Generates()) > 0,
				retryConfig,
				c.runMainStage,
			))
			if err != nil {
				c.setError(parsley.NewError(c.node.Pos(), fmt.Errorf("failed to run the main stage: %s", err)))
			}
		} else {
			c.stateChan <- containerStateNext
		}
	case containerStatePreClose:
		c.evalStage = conflow.EvalStageClose
		if err := c.evaluateChildren(); err != nil {
			c.setError(err)
			return
		}
	case containerStateClose:
		if _, ok := c.block.(conflow.BlockCloser); ok {
			err := c.jobTracker.ScheduleJob(newContainerStage(
				c,
				"close",
				conflow.EvalStageClose,
				false,
				conflow.RetryConfig{},
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
	if b, ok := c.block.(conflow.UserContexter); ok {
		c.evalCtx.UserContext = b.UserContext(c.evalCtx.UserContext)
	}

	if b, ok := c.block.(conflow.Loggerer); ok {
		c.evalCtx.Logger = b.Logger(c.evalCtx.Logger)
	}
}

func (c *Container) runInitStage() (int64, error) {
	skipped, err := c.block.(conflow.BlockInitialiser).Init(c.evalCtx)
	if err != nil {
		return 0, err
	}
	if skipped {
		return containerStateSkipped, nil
	}
	return containerStateNext, nil
}

func (c *Container) runMainStage() (int64, error) {
	result, err := c.block.(conflow.BlockRunner).Run(c.evalCtx)
	if err != nil {
		if re, ok := err.(conflow.Retryable); ok {
			if re.Retry() || re.RetryAfter() > 0 {
				return containerStateNext, retryError{
					Duration: re.RetryAfter(),
					Reason:   err.Error(),
				}
			}
		}
		return 0, err
	}

	if result != nil && (result.Retry() || result.RetryAfter() > 0) {
		return 0, retryError{
			Duration: result.RetryAfter(),
			Reason:   result.RetryReason(),
		}
	}

	return containerStateNext, nil
}

func (c *Container) runCloseStage() (int64, error) {
	return containerStateNext, c.block.(conflow.BlockCloser).Close(c.evalCtx)
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
func (c *Container) SetChild(container conflow.Container) {
	go func() {
		c.resultChan <- container
	}()
}

func (c *Container) setChild(result conflow.Container) parsley.Error {
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
		node := r.Node().(conflow.ParameterNode)
		name := node.Name()

		s := c.node.Interpreter().Schema()
		if _, ok := s.(*schema.Object).PropertyByParameterName(string(name)); ok {
			if err := c.node.Interpreter().SetParam(c.block, node.Name(), value); err != nil {
				return parsley.NewError(r.Node().Pos(), err)
			}
		} else {
			if c.extraParams == nil {
				c.extraParams = make(map[conflow.ID]interface{})
			}

			if _, exists := c.extraParams[name]; exists {
				panic(fmt.Errorf("%q parameter was already set on %q block", name, c.node.ID()))
			}

			c.extraParams[name] = value
		}

	case *Container:
		node := r.Node().(conflow.BlockNode)
		name, p := getNameSchemaForChildBlock(c.Node().Schema().(*schema.Object), node)

		if err := c.node.Interpreter().SetBlock(c.block, name, ptr.Value(node.Key()), value); err != nil {
			return parsley.NewError(r.Node().Pos(), err)
		}

		// Do not publish the empty generated node
		if p != nil && p.GetAnnotation(annotations.Generated) == "true" && c.evalStage == conflow.EvalStageInit {
			return nil
		}

	default:
		panic(fmt.Errorf("unknown container type: %T", result))
	}

	c.evalCtx.Publish(result)

	return nil
}

func (c *Container) PublishBlock(block conflow.Block, f func() error) (bool, error) {
	b, ok := block.(conflow.Identifiable)
	if !ok {
		return false, fmt.Errorf("%T block must implement the conflow.IDentifiable interface", block)
	}
	blockID := b.ID()

	nodeContainer, ok := c.children[blockID]
	parameterName := string(nodeContainer.Node().(conflow.BlockNode).ParameterName())
	s := c.Node().Schema().(*schema.Object)
	if !ok || s.Properties[s.JSONPropertyName(parameterName)].GetAnnotation(annotations.Generated) != "true" {
		return false, fmt.Errorf("%q block does not exist or is not marked as generated", blockID)
	}

	if !c.evalCtx.HasSubscribers(blockID) {
		return false, nil
	}

	wg := &util.WaitGroup{}
	wg.Add(1)
	container, err := nodeContainer.CreateContainer(block, []conflow.WaitGroup{wg})
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
func (c *Container) WaitGroups() []conflow.WaitGroup {
	return c.wgs
}

func (c *Container) debug() conflow.LogEvent {
	e := c.evalCtx.Logger.Debug()
	if e.Enabled() {
		e = e.ID("id", c.node.ID()).
			ID("type", c.node.BlockType()).
			Uint8("state", uint8(c.state)).
			Int("active_jobs", c.jobTracker.ActiveJobCount())
	}
	return e
}

func (c *Container) logError(err error) conflow.LogEvent {
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
