// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package blocks

import (
	"github.com/conflowio/parsley/parsley"

	"github.com/conflowio/conflow/src/conflow"
	"github.com/conflowio/conflow/src/schema"
)

type Root struct {
	params map[conflow.ID]interface{}
}

func (m *Root) ID() conflow.ID {
	return "root"
}

type RootInterpreter struct {
	BlockTransformerRegistry    parsley.NodeTransformerRegistry
	FunctionTransformerRegistry parsley.NodeTransformerRegistry
}

func (r RootInterpreter) Schema() schema.Schema {
	return &schema.Object{
		Metadata: schema.Metadata{
			Annotations: map[string]string{
				conflow.AnnotationType: conflow.BlockTypeRoot,
			},
		},
	}
}

func (r RootInterpreter) CreateBlock(conflow.ID, *conflow.BlockContext) conflow.Block {
	return &Root{
		params: map[conflow.ID]interface{}{},
	}
}

func (r RootInterpreter) SetParam(b conflow.Block, name conflow.ID, value interface{}) error {
	b.(*Root).params[name] = value
	return nil
}

func (r RootInterpreter) SetBlock(b conflow.Block, name conflow.ID, value interface{}) error {
	return nil
}

func (r RootInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	return b.(*Root).params[name]
}

func (r RootInterpreter) ValueParamName() conflow.ID {
	return ""
}

func (r RootInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	return ctx.New(conflow.ParseContextOverride{
		BlockTransformerRegistry:    r.BlockTransformerRegistry,
		FunctionTransformerRegistry: r.FunctionTransformerRegistry,
	})
}
