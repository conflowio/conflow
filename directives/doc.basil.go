// Code generated by Basil. DO NOT EDIT.

package directives

import (
	"fmt"
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/schema"
)

// DocInterpreter is the basil interpreter for the Doc block
type DocInterpreter struct {
	s schema.Schema
}

func (i DocInterpreter) Schema() schema.Schema {
	if i.s == nil {
		i.s = &schema.Object{
			Metadata: schema.Metadata{
				Annotations: map[string]string{"eval_stage": "ignore"},
			},
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

// Create creates a new Doc block
func (i DocInterpreter) CreateBlock(id basil.ID) basil.Block {
	return &Doc{
		id: id,
	}
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i DocInterpreter) ValueParamName() basil.ID {
	return "description"
}

// ParseContext returns with the parse context for the block
func (i DocInterpreter) ParseContext(ctx *basil.ParseContext) *basil.ParseContext {
	var nilBlock *Doc
	if b, ok := basil.Block(nilBlock).(basil.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i DocInterpreter) Param(b basil.Block, name basil.ID) interface{} {
	switch name {
	case "description":
		return b.(*Doc).description
	case "id":
		return b.(*Doc).id
	default:
		panic(fmt.Errorf("unexpected parameter %q in Doc", name))
	}
}

func (i DocInterpreter) SetParam(block basil.Block, name basil.ID, value interface{}) error {
	b := block.(*Doc)
	switch name {
	case "description":
		b.description = value.(string)
	}
	return nil
}

func (i DocInterpreter) SetBlock(block basil.Block, name basil.ID, value interface{}) error {
	return nil
}
