// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conflow

import (
	"github.com/conflowio/parsley/parsley"
)

// Container is a conflow object container
//counterfeiter:generate . Container
type Container interface {
	Node() Node
	Value() (interface{}, parsley.Error)
	WaitGroups() []WaitGroup
	Close()
}

// NilContainer is a container which evaluates to nil
type NilContainer struct {
	node    Node
	wgs     []WaitGroup
	pending bool
}

// NewNilContainer creates a new nil container
func NewNilContainer(node Node, wgs []WaitGroup, pending bool) Container {
	return NilContainer{
		node:    node,
		wgs:     wgs,
		pending: pending,
	}
}

func (n NilContainer) Node() Node {
	return n.node
}

func (n NilContainer) Value() (interface{}, parsley.Error) {
	return nil, nil
}

func (n NilContainer) WaitGroups() []WaitGroup {
	return n.wgs
}

func (n NilContainer) Close() {
	for _, wg := range n.wgs {
		wg.Done(nil)
	}
}

func (n NilContainer) Pending() bool {
	return n.pending
}
