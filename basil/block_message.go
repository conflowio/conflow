package basil

import "github.com/opsidian/basil/util"

// BlockMessage is a wrapper to send a block in a channel and block until it was processed
type BlockMessage interface {
	Block() Block
	Close(err error)
	Done() <-chan struct{}
	Err() error
}

type blockMessage struct {
	block Block
	done  *util.Done
}

// NewBlockMessage creates a new block message instance
func NewBlockMessage(block Block) BlockMessage {
	return &blockMessage{
		block: block,
		done:  &util.Done{},
	}
}

// Block returns with the block
func (b *blockMessage) Block() Block {
	return b.block
}

// Close will signal the client that the message was processed
func (b *blockMessage) Close(err error) {
	b.done.Close(err)
}

// Done will return with a channel which will return with nil or an error if the message was processed
// If there is no reader on the channel the message won't be sent.
// The channel will be closed after the optional error value.
func (b *blockMessage) Done() <-chan struct{} {
	return b.done.Done()
}

// Error will return the result error if Done() is closed, otherwise returns nil.
func (b *blockMessage) Err() error {
	return b.done.Err()
}
