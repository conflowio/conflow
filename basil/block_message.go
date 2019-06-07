// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package basil

import "github.com/opsidian/basil/util"

// BlockMessage is a wrapper to send a block in a channel and block until it was processed
type BlockMessage interface {
	Block() Block
	WaitGroup() *util.WaitGroup
}

type blockMessage struct {
	block Block
	wg    *util.WaitGroup
}

// NewBlockMessage creates a new block message instance
func NewBlockMessage(block Block) BlockMessage {
	return &blockMessage{
		block: block,
		wg:    &util.WaitGroup{},
	}
}

// Block returns with the block
func (b *blockMessage) Block() Block {
	return b.block
}

// WaitGroup returns with the wait group
func (b *blockMessage) WaitGroup() *util.WaitGroup {
	return b.wg
}
