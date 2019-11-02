// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package basil

import (
	"fmt"
	"sync"
)

type subscription struct {
	container *NodeContainer
	next      *subscription
}

type PubSub struct {
	subs map[ID]*subscription
	mu   *sync.RWMutex
}

func NewPubSub() *PubSub {
	return &PubSub{
		subs: make(map[ID]*subscription),
		mu:   &sync.RWMutex{},
	}
}

// Subscribe will subscribe the given node container for the given dependency
func (p *PubSub) Subscribe(c *NodeContainer, id ID) {
	p.mu.Lock()
	p.subs[id] = &subscription{container: c, next: p.subs[id]}
	p.mu.Unlock()
}

// Unsubscribe will unsubscribe the given node container for the given dependency
func (p *PubSub) Unsubscribe(c *NodeContainer, id ID) {
	p.mu.Lock()

	if p.subs[id].container.ID() == c.ID() {
		p.subs[id] = p.subs[id].next
		p.mu.Unlock()
		return
	}

	for sub := p.subs[id]; sub.next != nil; sub = sub.next {
		if sub.next.container.ID() == c.ID() {
			sub.next = sub.next.next
			p.mu.Unlock()
			return
		}
	}

	p.mu.Unlock()
	panic(fmt.Errorf("unsubscribe unsuccessful, %q was never subscribed for %q", c.ID(), id))
}

// Publish will notify all node containers which are subscribed for the dependency
// The ready function will run on any containers which have all dependencies satisfied
func (p *PubSub) Publish(c Container) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	for sub := p.subs[c.Node().ID()]; sub != nil; sub = sub.next {
		sub.container.SetDependency(c)
	}
}
