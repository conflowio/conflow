// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package blocks

import (
	"context"

	"github.com/opsidian/basil/basil"

	"github.com/opsidian/parsley/parsley"
)

func NewModule(id basil.ID, interpreter basil.BlockInterpreter, node parsley.Node) basil.Block {
	return &module{
		id:          id,
		interpreter: interpreter,
		node:        node,
		params:      map[basil.ID]interface{}{},
	}
}

type module struct {
	id          basil.ID `basil:"id"`
	interpreter basil.BlockInterpreter
	node        parsley.Node
	params      map[basil.ID]interface{}
}

func (m *module) ID() basil.ID {
	return m.id
}

func (m *module) Main(blockCtx basil.BlockContext) error {
	ctx, cancel := context.WithCancel(blockCtx)
	defer cancel()

	evalContext := basil.NewEvalContext(
		ctx, blockCtx.UserContext(), blockCtx.Logger(), blockCtx.JobScheduler(), nil,
	)
	evalContext.InputParams = m.params

	value, err := m.node.Value(evalContext)
	if err != nil {
		return err
	}

	for paramName, param := range m.interpreter.Params() {
		if param.IsOutput {
			m.params[paramName] = m.interpreter.Param(value.(basil.Block), paramName)
		}
	}

	return nil
}
