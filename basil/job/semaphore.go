// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package job

import "sync/atomic"

type Semaphore int64

func (s *Semaphore) Run() bool {
	return atomic.CompareAndSwapInt64((*int64)(s), 0, 1)
}

func (s *Semaphore) Cancel() bool {
	return atomic.CompareAndSwapInt64((*int64)(s), 0, 2)
}
