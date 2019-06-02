package basil

import (
	"sync"
)

// EvalStage means an evaluation stage (default, pre or post)
type EvalStage int8

// Evaluation stages
const (
	EvalStageInit  EvalStage = -1
	EvalStageMain  EvalStage = 0
	EvalStageClose EvalStage = 1
)

type EvalContext interface {
	WithDependencies(dependencies map[ID]BlockContainer) EvalContext
	BlockContext() BlockContext
	SetBlockContext(blockContext BlockContext)
	BlockContainer(id ID) (BlockContainer, bool)
	ScheduleJob(job Job)
	Subscribe(id ID, container *NodeContainer, ready func(c *NodeContainer))
	Publish(c Container)
}

// EvalContext is the evaluation context
type evalContext struct {
	parentCtx    EvalContext
	blockCtx     BlockContext
	scheduler    Scheduler
	dependencies map[ID]BlockContainer
	pubsub       *pubsub
}

// NewEvalContext returns with a new evaluation context
func NewEvalContext(blockCtx BlockContext, scheduler Scheduler) EvalContext {
	return &evalContext{
		blockCtx:  blockCtx,
		scheduler: scheduler,
		pubsub:    newPubSub(),
	}
}

// WithDependencies returns a copy of parent with the given dependencies
func (e *evalContext) WithDependencies(dependencies map[ID]BlockContainer) EvalContext {
	return &evalContext{
		parentCtx:    e,
		blockCtx:     e.blockCtx,
		scheduler:    e.scheduler,
		pubsub:       e.pubsub,
		dependencies: dependencies,
	}
}

// BlockContext returns with the block context
func (e *evalContext) BlockContext() BlockContext {
	return e.blockCtx
}

// SetBlockContext sets the block context
func (e *evalContext) SetBlockContext(blockContext BlockContext) {
	e.blockCtx = blockContext
}

// BlockContainer returns with the given block container instance if it exists
func (e *evalContext) BlockContainer(id ID) (BlockContainer, bool) {
	if container, ok := e.dependencies[id]; ok {
		return container, true
	}

	if e.parentCtx != nil {
		return e.parentCtx.BlockContainer(id)
	}

	return nil, false
}

func (e *evalContext) ScheduleJob(job Job) {
	e.scheduler.Schedule(job)
}

func (e *evalContext) Subscribe(id ID, container *NodeContainer, ready func(c *NodeContainer)) {
	e.pubsub.Subscribe(id, container, ready)
}

func (e *evalContext) Publish(c Container) {
	e.pubsub.Publish(c)
}

type subscription struct {
	container *NodeContainer
	ready     func(c *NodeContainer)
}

type pubsub struct {
	subscriptions     map[ID][]subscription
	subscriptionsLock *sync.RWMutex
}

func newPubSub() *pubsub {
	return &pubsub{
		subscriptions:     make(map[ID][]subscription, 32),
		subscriptionsLock: &sync.RWMutex{},
	}
}

// Subscribe will subscribe the given node container for the given dependency
func (p *pubsub) Subscribe(id ID, container *NodeContainer, ready func(c *NodeContainer)) {
	p.subscriptionsLock.Lock()
	defer p.subscriptionsLock.Unlock()

	p.subscriptions[id] = append(p.subscriptions[id], subscription{
		container: container,
		ready:     ready,
	})
}

// Publish will notify all node containers which are subscribed for the dependency
// The ready function will run on any containers which have all dependencies satisfied
func (p *pubsub) Publish(c Container) {
	p.subscriptionsLock.RLock()
	defer p.subscriptionsLock.RUnlock()

	for _, sub := range p.subscriptions[c.ID()] {
		if sub.container.SetDependency(c) {
			sub.ready(sub.container)
		}
	}
}
