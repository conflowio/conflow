// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package block

import (
	"fmt"

	"github.com/conflowio/conflow/src/conflow"
	"github.com/conflowio/conflow/src/schema"
)

type block struct {
	id     conflow.ID
	params map[conflow.ID]interface{}
	blocks map[conflow.ID]interface{}
}

func (b *block) ID() conflow.ID {
	return b.id
}

func NewInterpreter(s schema.ObjectKind) conflow.BlockInterpreter {
	interpreterRegistry := InterpreterRegistry{}

	var valueParamName string
	for name, p := range s.GetParameters() {
		if p.GetAnnotation(conflow.AnnotationValue) == "true" {
			valueParamName = name
		}

		switch pt := p.(type) {
		case schema.ObjectKind:
			p.(schema.MetadataAccessor).SetID(fmt.Sprintf("%s/%s", s.GetID(), name))
			interpreterRegistry[name] = NewInterpreter(pt)
			s.GetParameters()[name] = &schema.Reference{Ref: p.GetID()}
		case schema.ArrayKind:
			if o, ok := pt.GetItems().(schema.ObjectKind); ok {
				pt.GetItems().(schema.MetadataAccessor).SetID(fmt.Sprintf("%s/%s", s.GetID(), name))
				interpreterRegistry[name] = NewInterpreter(o)
				s.GetParameters()[name].(*schema.Array).Items = &schema.Reference{Ref: p.GetID()}
			}
		}
	}

	return interpreter{
		schema:              s,
		valueParamName:      conflow.ID(valueParamName),
		interpreterRegistry: interpreterRegistry,
	}
}

type interpreter struct {
	schema              schema.ObjectKind
	valueParamName      conflow.ID
	interpreterRegistry InterpreterRegistry
}

func (g interpreter) Schema() schema.Schema {
	return g.schema
}

func (g interpreter) CreateBlock(id conflow.ID, _ *conflow.BlockContext) conflow.Block {
	return &block{
		id:     id,
		params: map[conflow.ID]interface{}{},
		blocks: map[conflow.ID]interface{}{},
	}
}

func (g interpreter) SetParam(b conflow.Block, name conflow.ID, value interface{}) error {
	b.(*block).params[name] = value

	return nil
}

func (g interpreter) SetBlock(b conflow.Block, name conflow.ID, value interface{}) error {
	s := g.schema.GetParameters()[string(name)]
	if s.Type() == schema.TypeArray {
		if b.(*block).blocks[name] == nil {
			b.(*block).blocks[name] = []interface{}{}
		}
		b.(*block).blocks[name] = append(b.(*block).blocks[name].([]interface{}), value)
	} else {
		b.(*block).blocks[name] = value
	}

	return nil
}

func (g interpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	return b.(*block).params[name]
}

func (g interpreter) ValueParamName() conflow.ID {
	return g.valueParamName
}

func (g interpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	return ctx.New(conflow.ParseContextOverride{
		BlockTransformerRegistry: g.interpreterRegistry,
	})
}
