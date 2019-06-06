package block

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/opsidian/basil/util"

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
	ctx               basil.EvalContext
	node              basil.BlockNode
	block             basil.Block
	err               parsley.Error
	extraParams       map[basil.ID]interface{}
	state             containerState
	evalStage         basil.EvalStage
	stateChan         chan containerState
	childrenChan      chan basil.Container
	errChan           chan parsley.Error
	remainingJobs     int64
	children          map[basil.ID]*basil.NodeContainer
	generatedChildren map[basil.ID]*basil.NodeContainer
	cancel            context.CancelFunc
	wgs               []*util.WaitGroup
}

// NewContainer creates a new block container instance
func NewContainer(
	ctx basil.EvalContext,
	node basil.BlockNode,
	block basil.Block,
	wgs []*util.WaitGroup,
) *Container {
	if block == nil {
		block = node.Interpreter().CreateBlock(node.ID())
	}

	return &Container{
		ctx:          ctx,
		node:         node,
		block:        block,
		wgs:          wgs,
		evalStage:    basil.EvalStageInit,
		stateChan:    make(chan containerState, 1),
		childrenChan: make(chan basil.Container, 0), // We need tight synchronisation here
		errChan:      make(chan parsley.Error, 1),
	}
}

// ID returns with the block id
func (c *Container) ID() basil.ID {
	return c.node.ID()
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
	c.ctx.Logger().Debug().Fields(map[string]interface{}{"blockID": c.ID()}).Msg("started")

	defer func() {
		if c.cancel != nil {
			c.cancel()
		}
	}()

	c.stateChan <- containerStateNext
	c.mainLoop()

	// TODO: call context cancel here, but we need to initialise our own context+cancel

	c.node.Interpreter().CloseChannels(c)

	// Check the jobs that we never started
	if atomic.LoadInt64(&c.remainingJobs) > 0 {
		for _, container := range c.children {
			if container.Node().EvalStage() == c.evalStage && container.RunCount() == 0 {
				atomic.AddInt64(&c.remainingJobs, -1)
				container.Close()
			}
		}
	}

	if atomic.LoadInt64(&c.remainingJobs) > 0 {
		c.shutdownLoop()
	}
}

func (c *Container) mainLoop() {
	for {
		select {
		case state := <-c.stateChan:
			c.setState(state)

		case child := <-c.childrenChan:
			atomic.AddInt64(&c.remainingJobs, -1)

			err := c.setChild(child)
			child.Close()

			if err != nil {
				c.setState(containerStateErrored)
				c.err = parsley.NewError(c.node.Pos(), err)
				return
			}

			if atomic.LoadInt64(&c.remainingJobs) == 0 && c.state < containerStateFinished {
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
			atomic.AddInt64(&c.remainingJobs, -1)
			child.Close()

			if _, err := child.Value(); err != nil {
				// TODO: what to do with errors while shutting down?
			}
		case <-c.errChan:
			// TODO: what to do with additional errors?
		case <-shutdownTimer.C:
			return
		}

		if atomic.LoadInt64(&c.remainingJobs) == 0 {
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
			if node.Generated() {
				if c.generatedChildren == nil {
					c.generatedChildren = make(map[basil.ID]*basil.NodeContainer, 1)
				}
				c.generatedChildren[node.(basil.BlockNode).BlockType()] = c.createNodeContainer(node)
			} else {
				c.children[node.ID()] = c.createNodeContainer(node)
			}
		}
		c.node.Interpreter().ProcessChannels(c)

		c.evalStage = basil.EvalStageInit
		c.evaluateChildren()
	case containerStateInit:
		c.setBlockContext()

		if _, ok := c.block.(basil.BlockInitialiser); ok {
			atomic.AddInt64(&c.remainingJobs, 1)
			c.ctx.ScheduleJob(basil.NewJob(c.ctx, c.ID(), basil.JobFunc(c.init)))
		} else {
			c.stateChan <- containerStateNext
		}
	case containerStatePreMain:
		c.evalStage = basil.EvalStageMain
		c.evaluateChildren()
	case containerStateMain:
		// TODO: if the block is a generator, then run the main job in a goroutine instead on the main job queue
		if _, ok := c.block.(basil.BlockRunner); ok {
			atomic.AddInt64(&c.remainingJobs, 1)
			c.ctx.ScheduleJob(basil.NewJob(c.ctx, c.ID(), basil.JobFunc(c.main)))
		} else {
			c.stateChan <- containerStateNext
		}
	case containerStatePreClose:
		c.evalStage = basil.EvalStageClose
		c.evaluateChildren()
	case containerStateClose:
		if _, ok := c.block.(basil.BlockCloser); ok {
			atomic.AddInt64(&c.remainingJobs, 1)
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

	container := basil.NewNodeContainer(node, dependencies, c.scheduleChildJob)

	for id := range dependencies {
		c.ctx.Subscribe(container, id)
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

	logger := c.ctx.BlockContext().Logger()
	if b, ok := c.block.(basil.Loggerer); ok {
		logger = b.Logger(logger)
	}

	c.ctx.SetBlockContext(basil.NewBlockContext(ctx, userCtx, logger))
}

func (c *Container) init() {
	defer func() {
		if r := recover(); r != nil {
			c.errChan <- parsley.NewErrorf(c.node.Pos(), "init stage panicked in %q: %s", c.ID(), r)
		}
	}()

	c.ctx.Logger().Debug().Fields(map[string]interface{}{"blockID": c.ID()}).Msg("init stage started")

	start, err := c.block.(basil.BlockInitialiser).Init(c.ctx.BlockContext())
	atomic.AddInt64(&c.remainingJobs, -1)

	c.ctx.Logger().Debug().Fields(map[string]interface{}{"blockID": c.ID()}).Msg("init stage finished")

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
	defer func() {
		if r := recover(); r != nil {
			c.errChan <- parsley.NewErrorf(c.node.Pos(), "main stage panicked in %q: %s", c.ID(), r)
		}
	}()

	c.ctx.Logger().Debug().Fields(map[string]interface{}{"blockID": c.ID()}).Msg("main stage started")

	err := c.block.(basil.BlockRunner).Main(c.ctx.BlockContext())
	atomic.AddInt64(&c.remainingJobs, -1)

	c.ctx.Logger().Debug().Fields(map[string]interface{}{"blockID": c.ID()}).Msg("main stage finished")

	if err != nil {
		c.errChan <- parsley.NewError(c.node.Pos(), err)
		return
	}

	c.stateChan <- containerStateNext
}

func (c *Container) close() {
	defer func() {
		if r := recover(); r != nil {
			c.errChan <- parsley.NewErrorf(c.node.Pos(), "close stage panicked in %q: %s", c.ID(), r)
		}
	}()

	c.ctx.Logger().Debug().Fields(map[string]interface{}{"blockID": c.ID()}).Msg("close stage started")

	err := c.block.(basil.BlockCloser).Close(c.ctx.BlockContext())
	atomic.AddInt64(&c.remainingJobs, -1)

	c.ctx.Logger().Debug().Fields(map[string]interface{}{"blockID": c.ID()}).Msg("close stage finished")

	if err != nil {
		c.errChan <- parsley.NewError(c.node.Pos(), err)
		return
	}

	c.stateChan <- containerStateNext
}

func (c *Container) evaluateChildren() {
	for _, container := range c.children {
		if container.Node().EvalStage() == c.evalStage && container.Ready() {
			// TODO: check if we started no jobs
			c.scheduleChildJob(container, nil)
		}
	}

	if atomic.LoadInt64(&c.remainingJobs) == 0 {
		c.stateChan <- containerStateNext
	}
}

func (c *Container) scheduleChildJob(nodeContainer *basil.NodeContainer, wgs []*util.WaitGroup) {
	if nodeContainer.Node().EvalStage() != c.evalStage {
		return
	}

	ctx := nodeContainer.EvalContext(c.ctx)
	var container basil.Container

	switch n := nodeContainer.Node().(type) {
	case basil.BlockNode:
		container = NewContainer(ctx, n, nil, wgs)
	case basil.ParameterNode:
		container = parameter.NewContainer(ctx, n, c)
	default:
		panic(fmt.Errorf("invalid node type: %T", n))
	}

	atomic.AddInt64(&c.remainingJobs, 1)
	nodeContainer.IncRunCount()

	c.ctx.Logger().Debug().Fields(map[string]interface{}{
		"blockID":       c.ID(),
		"id":            nodeContainer.ID(),
		"remainingJobs": atomic.LoadInt64(&c.remainingJobs),
	}).Msg("scheduled")

	// A block's main loop is lightweight, and evaluating parameters is cheap, so we start a goroutine only
	// Also we don't want to block the main loop or the main job queue
	go func() {
		container.Run()
		c.childrenChan <- container
	}()
}

func (c *Container) setChild(result basil.Container) parsley.Error {
	c.ctx.Logger().Debug().Fields(map[string]interface{}{
		"blockID":       c.ID(),
		"id":            result.ID(),
		"remainingJobs": atomic.LoadInt64(&c.remainingJobs),
	}).Msg("set child")

	value, err := result.Value()
	if err != nil {
		return err
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
		if r.state == containerStateSkipped {
			return nil
		}
		name := node.BlockType()

		if err := c.node.Interpreter().SetBlock(c.block, name, value); err != nil {
			return parsley.NewError(r.Node().Pos(), err)
		}
	default:
		panic(fmt.Errorf("unknown container type: %T", result))
	}

	c.ctx.Publish(result)

	return nil
}

func (c *Container) PublishBlock(blockType basil.ID, blockMessage basil.BlockMessage) {
	atomic.AddInt64(&c.remainingJobs, 1)

	nodeContainer := c.generatedChildren[blockType]
	nodeContainer.IncRunCount()

	ctx := nodeContainer.EvalContext(c.ctx)

	wgs := []*util.WaitGroup{blockMessage.WaitGroup()}
	blockMessage.WaitGroup().Add(1)
	container := NewContainer(ctx, nodeContainer.Node().(basil.BlockNode), blockMessage.Block(), wgs)
	container.Run()

	c.childrenChan <- container

	blockMessage.WaitGroup().Wait()
}

// Close does nothing
func (c *Container) Close() {
	for _, wg := range c.wgs {
		wg.Done(c.err)
	}
	c.ctx.Logger().Debug().Fields(map[string]interface{}{"blockID": c.ID()}).Msg("finished")
}

// WaitGroups returns nil
func (c *Container) WaitGroups() []*util.WaitGroup {
	return c.wgs
}
