// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package basil

import (
	"github.com/opsidian/basil/util"
	"github.com/opsidian/parsley/parsley"
)

// Container is a basil object container
type Container interface {
	Job
	Node() Node
	Value() (interface{}, parsley.Error)
	WaitGroups() []*util.WaitGroup
	Close()
}
