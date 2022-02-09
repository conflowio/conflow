// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package blocks

import (
	"errors"
	"fmt"

	"github.com/conflowio/conflow/src/conflow"
	"github.com/conflowio/conflow/src/conflow/block"
	"github.com/conflowio/conflow/src/parsers"
)

// @block {
//   type = "task"
//   eval_stage = "parse"
// }
type Import struct {
	// @id
	id conflow.ID
	// @value
	// @required
	path string
}

func (i *Import) ID() conflow.ID {
	return i.id
}

func (i *Import) EvalStage() conflow.EvalStage {
	return conflow.EvalStageParse
}

func (i *Import) BlockInterpreters(parseCtx *conflow.ParseContext) (map[conflow.ID]conflow.BlockInterpreter, error) {
	parseCtx = parseCtx.NewForModule()
	registry := parseCtx.BlockTransformerRegistry().(block.InterpreterRegistry)

	moduleInterpreter, ok := registry["module"]
	if !ok {
		return nil, errors.New("import is not allowed as the \"module\" block is not registered")
	}

	p := parsers.NewRoot("root", moduleInterpreter)

	if err := p.ParseDir(
		parseCtx,
		i.path,
	); err != nil {
		return nil, err
	}

	node, ok := parseCtx.BlockNode("root")
	if !ok {
		return nil, fmt.Errorf("block \"root\" does not exist in module %s", i.path)
	}

	return map[conflow.ID]conflow.BlockInterpreter{
		i.id: NewModuleInterpreter(node.Interpreter(), node),
	}, nil
}
