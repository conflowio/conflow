package util

import "sync"

// closedChan is a reusable closed channel.
var closedChan = make(chan struct{})

func init() {
	close(closedChan)
}

// Done is a lazily created channel which can be used to signal clients when an action is done
// It's similar to context.Context's Done function
type Done struct {
	ch     chan struct{}
	mu     sync.Mutex
	err    error
	closed bool
}

func (d *Done) Done() <-chan struct{} {
	d.mu.Lock()
	if d.ch == nil {
		d.ch = make(chan struct{})
	}
	ch := d.ch
	d.mu.Unlock()
	return ch
}

func (d *Done) Close(err error) {
	d.mu.Lock()
	if d.err == nil {
		d.err = err
	}
	if !d.closed {
		if d.ch == nil {
			d.ch = closedChan
		} else {
			close(d.ch)
		}
	}
	d.closed = true
	d.mu.Unlock()
}

func (d *Done) Err() error {
	d.mu.Lock()
	err := d.err
	d.mu.Unlock()
	return err
}
