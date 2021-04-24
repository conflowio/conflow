// Code generated by Basil. DO NOT EDIT.

package directives

import (
	"fmt"
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/schema"
)

// SkipInterpreter is the basil interpreter for the Skip block
type SkipInterpreter struct {
	s schema.Schema
}

func (i SkipInterpreter) Schema() schema.Schema {
	if i.s == nil {
		i.s = &schema.Object{
			Metadata: schema.Metadata{
				Annotations: map[string]string{"eval_stage": "init"},
			},
			Properties: map[string]schema.Schema{
				"id": &schema.String{
					Metadata: schema.Metadata{
						Annotations: map[string]string{"id": "true"},
						ReadOnly:    true,
					},
					Format: "basil.ID",
				},
				"when": &schema.Boolean{
					Metadata: schema.Metadata{
						Annotations: map[string]string{"value": "true"},
					},
					Default: schema.BooleanPtr(true),
				},
			},
		}
	}
	return i.s
}

// Create creates a new Skip block
func (i SkipInterpreter) CreateBlock(id basil.ID) basil.Block {
	return &Skip{
		id:   id,
		when: true,
	}
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i SkipInterpreter) ValueParamName() basil.ID {
	return "when"
}

// ParseContext returns with the parse context for the block
func (i SkipInterpreter) ParseContext(ctx *basil.ParseContext) *basil.ParseContext {
	var nilBlock *Skip
	if b, ok := basil.Block(nilBlock).(basil.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i SkipInterpreter) Param(b basil.Block, name basil.ID) interface{} {
	switch name {
	case "id":
		return b.(*Skip).id
	case "when":
		return b.(*Skip).when
	default:
		panic(fmt.Errorf("unexpected parameter %q in Skip", name))
	}
}

func (i SkipInterpreter) SetParam(block basil.Block, name basil.ID, value interface{}) error {
	b := block.(*Skip)
	switch name {
	case "when":
		b.when = value.(bool)
	}
	return nil
}

func (i SkipInterpreter) SetBlock(block basil.Block, name basil.ID, value interface{}) error {
	return nil
}
