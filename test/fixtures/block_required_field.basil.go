// Code generated by Basil. DO NOT EDIT.

package fixtures

import (
	"fmt"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/schema"
)

// BlockRequiredFieldInterpreter is the basil interpreter for the BlockRequiredField block
type BlockRequiredFieldInterpreter struct {
	s schema.Schema
}

func (i BlockRequiredFieldInterpreter) Schema() schema.Schema {
	if i.s == nil {
		i.s = &schema.Object{
			Properties: map[string]schema.Schema{
				"id_field": &schema.String{
					Metadata: schema.Metadata{
						Annotations: map[string]string{"id": "true"},
						ReadOnly:    true,
					},
					Format: "basil.ID",
				},
				"required": &schema.Untyped{},
			},
			Required:         []string{"required"},
			StructProperties: map[string]string{"id_field": "IDField"},
		}
	}
	return i.s
}

// Create creates a new BlockRequiredField block
func (i BlockRequiredFieldInterpreter) CreateBlock(id basil.ID) basil.Block {
	return &BlockRequiredField{
		IDField: id,
	}
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i BlockRequiredFieldInterpreter) ValueParamName() basil.ID {
	return ""
}

// ParseContext returns with the parse context for the block
func (i BlockRequiredFieldInterpreter) ParseContext(ctx *basil.ParseContext) *basil.ParseContext {
	var nilBlock *BlockRequiredField
	if b, ok := basil.Block(nilBlock).(basil.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i BlockRequiredFieldInterpreter) Param(b basil.Block, name basil.ID) interface{} {
	switch name {
	case "id_field":
		return b.(*BlockRequiredField).IDField
	case "required":
		return b.(*BlockRequiredField).required
	default:
		panic(fmt.Errorf("unexpected parameter %q in BlockRequiredField", name))
	}
}

func (i BlockRequiredFieldInterpreter) SetParam(block basil.Block, name basil.ID, value interface{}) error {
	b := block.(*BlockRequiredField)
	switch name {
	case "required":
		b.required = value
	}
	return nil
}

func (i BlockRequiredFieldInterpreter) SetBlock(block basil.Block, name basil.ID, value interface{}) error {
	return nil
}
