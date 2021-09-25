// Code generated by Basil. DO NOT EDIT.

package directives

import (
	"fmt"
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/schema"
)

// DeprecatedInterpreter is the basil interpreter for the Deprecated block
type DeprecatedInterpreter struct {
	s schema.Schema
}

func (i DeprecatedInterpreter) Schema() schema.Schema {
	if i.s == nil {
		i.s = &schema.Object{
			Metadata: schema.Metadata{
				Annotations: map[string]string{"eval_stage": "ignore"},
			},
			Name: "Deprecated",
			Properties: map[string]schema.Schema{
				"description": &schema.String{
					Metadata: schema.Metadata{
						Annotations: map[string]string{"value": "true"},
					},
				},
				"id": &schema.String{
					Metadata: schema.Metadata{
						Annotations: map[string]string{"id": "true"},
						ReadOnly:    true,
					},
					Format: "basil.ID",
				},
			},
			Required: []string{"description"},
		}
	}
	return i.s
}

// Create creates a new Deprecated block
func (i DeprecatedInterpreter) CreateBlock(id basil.ID, blockCtx *basil.BlockContext) basil.Block {
	return &Deprecated{
		id: id,
	}
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i DeprecatedInterpreter) ValueParamName() basil.ID {
	return "description"
}

// ParseContext returns with the parse context for the block
func (i DeprecatedInterpreter) ParseContext(ctx *basil.ParseContext) *basil.ParseContext {
	var nilBlock *Deprecated
	if b, ok := basil.Block(nilBlock).(basil.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i DeprecatedInterpreter) Param(b basil.Block, name basil.ID) interface{} {
	switch name {
	case "description":
		return b.(*Deprecated).description
	case "id":
		return b.(*Deprecated).id
	default:
		panic(fmt.Errorf("unexpected parameter %q in Deprecated", name))
	}
}

func (i DeprecatedInterpreter) SetParam(block basil.Block, name basil.ID, value interface{}) error {
	b := block.(*Deprecated)
	switch name {
	case "description":
		b.description = value.(string)
	}
	return nil
}

func (i DeprecatedInterpreter) SetBlock(block basil.Block, name basil.ID, value interface{}) error {
	return nil
}
