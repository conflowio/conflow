package block

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/opsidian/parsley/parsley"

	"github.com/opsidian/basil/basil/parameter"

	"github.com/opsidian/basil/basil"
)

type containerState uint8

const (
	containerStateNext containerState = iota
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
)

// ContainerGracefulTimeoutSec is the graceful timeout in seconds
var ContainerGracefulTimeoutSec = 10

var _ basil.BlockContainer = &Container{}

// Container is a block container
type Container struct {
	ctx           basil.EvalContext
	node          basil.BlockNode
	block         basil.Block
	err           parsley.Error
	extraParams   map[basil.ID]interface{}
	state         containerState
	evalStage     basil.EvalStage
	stateChan     chan containerState
	childrenChan  chan basil.Container
	errChan       chan parsley.Error
	remainingJobs int
	children      map[basil.ID]*basil.NodeContainer
	cancel        context.CancelFunc
}

// NewContainer creates a new block container instance
func NewContainer(
	ctx basil.EvalContext,
	node basil.BlockNode,
	block basil.Block,
) *Container {
	if block == nil {
		block = node.Interpreter().CreateBlock(node.ID())
	}

	return &Container{
		ctx:       ctx,
		node:      node,
		block:     block,
		evalStage: basil.EvalStageInit,
	}
}

// ID returns with the block id
func (c *Container) ID() basil.ID {
	return c.node.ID()
}

// Node returns with the block node
func (c *Container) Node() basil.BlockNode {
	return c.node
}

// Value returns with the block or the error if any occurred
func (c *Container) Value() (interface{}, parsley.Error) {
	if c.err != nil {
		return nil, c.err
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
		if c.cancel != nil {
			c.cancel()
		}
	}()

	c.stateChan = make(chan containerState, 1)
	c.childrenChan = make(chan basil.Container, 1)
	c.errChan = make(chan parsley.Error, 1)

	c.stateChan <- containerStateNext
	c.mainLoop()

	// Check the jobs that we never started
	if c.remainingJobs > 0 {
		for _, container := range c.children {
			if container.Node().EvalStage() == c.evalStage && !container.Ready() {
				c.remainingJobs--
			}
		}
	}

	if c.remainingJobs > 0 {
		c.shutdownLoop()
	}
}

func (c *Container) mainLoop() {
	for {
		select {
		case state := <-c.stateChan:
			c.setState(state)

		case child := <-c.childrenChan:
			if err := c.setChild(child); err != nil {
				c.setState(containerStateErrored)
				c.err = parsley.NewError(c.node.Pos(), err)
				return
			}

			if c.remainingJobs == 0 && c.state < containerStateFinished {
				c.stateChan <- containerStateNext
			}

		case err := <-c.errChan:
			c.setState(containerStateErrored)
			c.err = err
			return

		case <-c.ctx.BlockContext().Context().Done():
			c.setState(containerStateAborted)
			c.err = parsley.NewError(c.node.Pos(), errors.New("aborted"))
			return
		}

		if c.state >= containerStateFinished {
			return
		}
	}
}

func (c *Container) shutdownLoop() {
	shutdownTimer := time.NewTimer(time.Duration(ContainerGracefulTimeoutSec) * time.Second)

	for {
		select {
		case child := <-c.childrenChan:
			if _, err := child.Value(); err != nil {
				// TODO: what to do with errors while shutting down?
			}
			c.remainingJobs--
		case <-c.errChan:
			// TODO: what to do with additional errors?
		case <-shutdownTimer.C:
			return
		}

		if c.remainingJobs == 0 {
			return
		}
	}
}

func (c *Container) setState(state containerState) {
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
			c.children[node.ID()] = c.createNodeContainer(node)
		}

		c.evalStage = basil.EvalStageInit
		c.evaluateChildren()
	case containerStateInit:
		c.setBlockContext()

		if _, ok := c.block.(basil.BlockInitialiser); ok {
			c.ctx.ScheduleJob(basil.NewJob(c.ctx, c.ID(), basil.JobFunc(c.init)))
		} else {
			c.stateChan <- containerStateNext
		}
	case containerStatePreMain:
		c.evalStage = basil.EvalStageMain
		c.evaluateChildren()
	case containerStateMain:
		if _, ok := c.block.(basil.BlockRunner); ok {
			c.ctx.ScheduleJob(basil.NewJob(c.ctx, c.ID(), basil.JobFunc(c.main)))
		} else {
			c.stateChan <- containerStateNext
		}
	case containerStatePreClose:
		c.evalStage = basil.EvalStageClose
		c.evaluateChildren()
	case containerStateClose:
		if _, ok := c.block.(basil.BlockCloser); ok {
			c.ctx.ScheduleJob(basil.NewJob(c.ctx, c.ID(), basil.JobFunc(c.close)))
		} else {
			c.stateChan <- containerStateNext
		}
	case containerStateFinished:
	case containerStateSkipped:
		// TODO: notify block about skipped task
	case containerStateErrored:
		// TODO: notify block about error
	case containerStateAborted:
		// TODO: notify block about abort
	default:
		panic("invalid container state")
	}
}

func (c *Container) createNodeContainer(node basil.Node) *basil.NodeContainer {
	dependencies := make(map[basil.ID]basil.Container, len(node.Dependencies()))
	for _, v := range node.Dependencies() {
		if _, ok := c.Node().Dependencies()[v.ID()]; ok {
			continue
		}
		if c.ID() == v.ParentID() {
			dependencies[v.ID()] = nil
		} else {
			dependencies[v.ParentID()] = nil
		}
	}

	container := basil.NewNodeContainer(node, dependencies)

	for id := range dependencies {
		c.ctx.Subscribe(id, container, c.scheduleChildJob)
	}

	return container
}

func (c *Container) setBlockContext() {
	var ctx context.Context
	if b, ok := c.block.(basil.Contexter); ok {
		ctx, c.cancel = b.Context(c.ctx.BlockContext().Context())
	} else {
		ctx, c.cancel = context.WithCancel(c.ctx.BlockContext().Context())
	}

	userCtx := c.ctx.BlockContext().UserContext()
	if b, ok := c.block.(basil.UserContexter); ok {
		userCtx = b.UserContext(userCtx)
	}

	c.ctx.SetBlockContext(basil.NewBlockContext(ctx, userCtx))
}

func (c *Container) init() {
	start, err := c.block.(basil.BlockInitialiser).Init(c.ctx.BlockContext())
	if err != nil {
		c.errChan <- parsley.NewError(c.node.Pos(), err)
		return
	}
	if !start {
		c.stateChan <- containerStateSkipped
		return
	}

	c.stateChan <- containerStateNext
}

func (c *Container) main() {
	if err := c.block.(basil.BlockRunner).Main(c.ctx.BlockContext()); err != nil {
		c.errChan <- parsley.NewError(c.node.Pos(), err)
		return
	}

	c.stateChan <- containerStateNext
}

func (c *Container) close() {
	if err := c.block.(basil.BlockCloser).Close(c.ctx.BlockContext()); err != nil {
		c.errChan <- parsley.NewError(c.node.Pos(), err)
		return
	}

	c.stateChan <- containerStateNext
}

func (c *Container) evaluateChildren() {
	for _, container := range c.children {
		if container.Node().EvalStage() == c.evalStage {
			c.remainingJobs++
			c.scheduleChildJob(container)
		}
	}

	if c.remainingJobs == 0 {
		c.stateChan <- containerStateNext
	}
}

func (c *Container) scheduleChildJob(nodeContainer *basil.NodeContainer) {
	// We have to make sure nodes are only evaluated in the correct evaluation stage
	if nodeContainer.Node().EvalStage() != c.evalStage || !nodeContainer.Ready() {
		return
	}

	ctx := nodeContainer.EvalContext(c.ctx)
	var container basil.Container

	switch n := nodeContainer.Node().(type) {
	case basil.BlockNode:
		container = NewContainer(ctx, n, nil)
	case basil.ParameterNode:
		container = parameter.NewContainer(ctx, n, c)
	default:
		panic(fmt.Errorf("invalid node type: %T", n))
	}

	// A block's main loop is lightweight, and evaluating parameters is cheap, so we start a goroutine only
	// Also we don't want to block the main loop or the main job queue
	go func() {
		container.Run()
		c.childrenChan <- container
	}()
}

func (c *Container) setChild(result basil.Container) parsley.Error {
	c.remainingJobs--

	value, err := result.Value()
	if err != nil {
		return err
	}

	switch r := result.(type) {
	case *parameter.Container:
		name := r.Node().Name()

		if r.Node().IsDeclaration() {
			if c.extraParams == nil {
				c.extraParams = make(map[basil.ID]interface{})
			}

			if _, exists := c.extraParams[name]; exists {
				panic(fmt.Errorf("%q parameter was already set on %q block", name, c.ID()))
			}

			c.extraParams[name] = value
		} else {
			if err := c.node.Interpreter().SetParam(c.block, r.Node().Name(), value); err != nil {
				return parsley.NewError(r.Node().Pos(), err)
			}
		}
	case *Container:
		if r.state == containerStateSkipped {
			return nil
		}
		name := r.Node().BlockType()

		if err := c.node.Interpreter().SetBlock(c.block, name, value); err != nil {
			return parsley.NewError(r.Node().Pos(), err)
		}
	default:
		panic(fmt.Errorf("unknown container type: %T", result))
	}

	c.ctx.Publish(result)

	return nil
}
