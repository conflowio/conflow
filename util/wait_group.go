// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package util

import (
	"sync"

	"github.com/opsidian/basil/util/multierror"
)

// WaitGroup is a slightly different version of sync.WaitGroup with error handling
// and Wait() returns with a channel instead of blocking.
type WaitGroup struct {
	ch  chan struct{}
	mu  sync.Mutex
	err *multierror.Error
	cnt int
}

func (w *WaitGroup) Add(delta int) {
	w.mu.Lock()
	w.cnt = w.cnt + delta
	w.mu.Unlock()
}

func (w *WaitGroup) Wait() <-chan struct{} {
	w.mu.Lock()
	if w.ch == nil {
		w.ch = make(chan struct{})
	}
	ch := w.ch
	w.mu.Unlock()
	return ch
}

func (w *WaitGroup) Done(err error) {
	w.mu.Lock()
	if err != nil {
		w.err.Append(err)
	}
	w.cnt--
	if w.cnt < 0 {
		w.mu.Unlock()
		panic("negative WaitGroup counter")
	}
	if w.cnt == 0 {
		if w.ch == nil {
			w.ch = closedChan
		} else {
			close(w.ch)
		}
	}
	w.mu.Unlock()
}

func (w *WaitGroup) Err() error {
	w.mu.Lock()
	err := w.err.Err()
	w.mu.Unlock()
	return err
}
