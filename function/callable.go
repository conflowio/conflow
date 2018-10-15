// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package function

import (
	"github.com/opsidian/parsley/parsley"
)

// Callable is an interface for callable objects (general function interface)
//go:generate counterfeiter . Callable
type Callable interface {
	Call(ctx interface{}, function parsley.Node, params []parsley.Node) (interface{}, parsley.Error)
}

// CallableFunc defines a helper to implement the Callable interface with functions
type CallableFunc func(ctx interface{}, function parsley.Node, nodes []parsley.Node) (interface{}, parsley.Error)

// Call calls the function
func (f CallableFunc) Call(ctx interface{}, function parsley.Node, nodes []parsley.Node) (interface{}, parsley.Error) {
	return f(ctx, function, nodes)
}
