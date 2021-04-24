// Code generated by Basil. DO NOT EDIT.

package common

import (
	"fmt"
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/schema"
)

// IteratorInterpreter is the basil interpreter for the Iterator block
type IteratorInterpreter struct {
	s schema.Schema
}

func (i IteratorInterpreter) Schema() schema.Schema {
	if i.s == nil {
		i.s = &schema.Object{
			Properties: map[string]schema.Schema{
				"count": &schema.Integer{},
				"id": &schema.String{
					Metadata: schema.Metadata{
						Annotations: map[string]string{"id": "true"},
						ReadOnly:    true,
					},
					Format: "basil.ID",
				},
				"it": &schema.Reference{
					Metadata: schema.Metadata{
						Annotations: map[string]string{"eval_stage": "init", "generated": "true"},
						Pointer:     true,
					},
					Ref: "http://basil.schema/github.com/opsidian/basil/examples/common/It",
				},
			},
			Required: []string{"count", "it"},
		}
	}
	return i.s
}

// Create creates a new Iterator block
func (i IteratorInterpreter) CreateBlock(id basil.ID) basil.Block {
	return &Iterator{
		id: id,
	}
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i IteratorInterpreter) ValueParamName() basil.ID {
	return ""
}

// ParseContext returns with the parse context for the block
func (i IteratorInterpreter) ParseContext(ctx *basil.ParseContext) *basil.ParseContext {
	var nilBlock *Iterator
	if b, ok := basil.Block(nilBlock).(basil.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i IteratorInterpreter) Param(b basil.Block, name basil.ID) interface{} {
	switch name {
	case "count":
		return b.(*Iterator).count
	case "id":
		return b.(*Iterator).id
	default:
		panic(fmt.Errorf("unexpected parameter %q in Iterator", name))
	}
}

func (i IteratorInterpreter) SetParam(block basil.Block, name basil.ID, value interface{}) error {
	b := block.(*Iterator)
	switch name {
	case "count":
		b.count = value.(int64)
	}
	return nil
}

func (i IteratorInterpreter) SetBlock(block basil.Block, name basil.ID, value interface{}) error {
	b := block.(*Iterator)
	switch name {
	case "it":
		b.it = value.(*It)
	}
	return nil
}
