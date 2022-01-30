// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package block

import (
	"github.com/conflowio/parsley/parsley"

	"github.com/conflowio/conflow/src/conflow"
)

// InterpreterRegistry contains a list of block interpreters and behaves as a node transformer registry
type InterpreterRegistry map[string]conflow.BlockInterpreter

// NodeTransformer returns with the named node transformer
func (i InterpreterRegistry) NodeTransformer(name string) (parsley.NodeTransformer, bool) {
	interpreter, exists := i[name]
	if !exists {
		return nil, false
	}

	return parsley.NodeTransformFunc(func(userCtx interface{}, node parsley.Node) (parsley.Node, parsley.Error) {
		return TransformNode(userCtx, node, interpreter)
	}), true
}

func (i InterpreterRegistry) Register(name string, interpreter conflow.BlockInterpreter) InterpreterRegistry {
	i[name] = interpreter
	return i
}
