// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package blocks

import (
	"errors"
	"fmt"

	"github.com/opsidian/basil/basil/block"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/parsers"
)

//go:generate basil generate
type Import struct {
	id   basil.ID `basil:"id"`
	path string   `basil:"value,required"`
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

	p := parsers.NewMain("main", moduleInterpreter.(basil.BlockInterpreter))

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
