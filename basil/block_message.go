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
