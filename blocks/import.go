// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package blocks

import (
	"errors"
	"fmt"

	"github.com/opsidian/conflow/basil"
	"github.com/opsidian/conflow/basil/block"
	"github.com/opsidian/conflow/parsers"
)

// @block {
//   eval_stage = "parse"
// }
type Import struct {
	// @id
	id basil.ID
	// @value
	// @required
	path string
}

func (i *Import) ID() basil.ID {
	return i.id
}

func (i *Import) EvalStage() basil.EvalStage {
	return basil.EvalStageParse
}

func (i *Import) BlockInterpreters(parseCtx *basil.ParseContext) (map[basil.ID]basil.BlockInterpreter, error) {
	parseCtx = parseCtx.NewForModule()
	registry := parseCtx.BlockTransformerRegistry().(block.InterpreterRegistry)

	moduleInterpreter, ok := registry["module"]
	if !ok {
		return nil, errors.New("import is not allowed as the \"module\" block is not registered")
	}

	p := parsers.NewMain("main", moduleInterpreter)

	if err := p.ParseDir(
		parseCtx,
		i.path,
	); err != nil {
		return nil, err
	}

	node, ok := parseCtx.BlockNode("main")
	if !ok {
		return nil, fmt.Errorf("block \"main\" does not exist in module %s", i.path)
	}

	return map[basil.ID]basil.BlockInterpreter{
		i.id: NewModuleInterpreter(node.Interpreter(), node),
	}, nil
}
