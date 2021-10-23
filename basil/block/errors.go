// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package block

import "time"

type retryError struct {
	Duration time.Duration
	Reason   string
}

func (r retryError) Error() string {
	return r.Reason
}
