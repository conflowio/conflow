// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package util

import (
	"sync"
	"sync/atomic"
)

// WaitGroup is a slightly different version of sync.WaitGroup with error handling
// and Wait() returns with a channel instead of blocking.
type WaitGroup struct {
	ch  chan struct{}
	mu  sync.Mutex
	err error
	cnt int64
}

func (w *WaitGroup) Add(delta int) {
	atomic.AddInt64(&w.cnt, int64(delta))
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
	if w.err == nil {
		w.err = err
	}
	cnt := atomic.AddInt64(&w.cnt, -1)
	if cnt < 0 {
		panic("negative WaitGroup counter")
	}
	if cnt == 0 {
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
	err := w.err
	w.mu.Unlock()
	return err
}
