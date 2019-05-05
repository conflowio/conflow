// Copyright (c) 2018 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package block

import (
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/parsley/parsley"
)

// Interpreter defines an interpreter for blocks
//go:generate counterfeiter . Interpreter
type Interpreter interface {
	Create(ctx *basil.EvalContext, node basil.BlockNode) basil.Block
	SetParam(ctx *basil.EvalContext, b basil.Block, name basil.ID, node parsley.Node) parsley.Error
	Param(block basil.Block, name basil.ID) interface{}
	Params() map[basil.ID]string
	RequiredParams() map[basil.ID]bool
	ValueParamName() basil.ID
	HasForeignID() bool
	basil.ParseContextAware
}

// InterpreterRegistry contains a list of block interpreters and behaves as a node transformer registry
type InterpreterRegistry map[string]Interpreter

// NodeTransformer returns with the named node transformer
func (i InterpreterRegistry) NodeTransformer(name string) (parsley.NodeTransformer, bool) {
	interpreter, exists := i[name]
	if !exists {
		return nil, false
	}

	if name == basil.MainID {
		return parsley.NodeTransformFunc(func(userCtx interface{}, node parsley.Node) (parsley.Node, parsley.Error) {
			return transformMainNode(userCtx, node, interpreter)
		}), true
	}

	return parsley.NodeTransformFunc(func(userCtx interface{}, node parsley.Node) (parsley.Node, parsley.Error) {
		return transformNode(userCtx, node, interpreter)
	}), true
}
