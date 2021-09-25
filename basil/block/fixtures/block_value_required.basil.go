// Code generated by Basil. DO NOT EDIT.

package fixtures

import (
	"fmt"
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/schema"
)

// BlockValueRequiredInterpreter is the basil interpreter for the BlockValueRequired block
type BlockValueRequiredInterpreter struct {
	s schema.Schema
}

func (i BlockValueRequiredInterpreter) Schema() schema.Schema {
	if i.s == nil {
		i.s = &schema.Object{
			Name: "BlockValueRequired",
			Properties: map[string]schema.Schema{
				"id_field": &schema.String{
					Metadata: schema.Metadata{
						Annotations: map[string]string{"id": "true"},
						ReadOnly:    true,
					},
					Format: "basil.ID",
				},
				"value": &schema.String{
					Metadata: schema.Metadata{
						Annotations: map[string]string{"value": "true"},
					},
				},
			},
			PropertyNames: map[string]string{"id_field": "IDField", "value": "Value"},
			Required:      []string{"value"},
		}
	}
	return i.s
}

// Create creates a new BlockValueRequired block
func (i BlockValueRequiredInterpreter) CreateBlock(id basil.ID, blockCtx *basil.BlockContext) basil.Block {
	return &BlockValueRequired{
		IDField: id,
	}
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i BlockValueRequiredInterpreter) ValueParamName() basil.ID {
	return "value"
}

// ParseContext returns with the parse context for the block
func (i BlockValueRequiredInterpreter) ParseContext(ctx *basil.ParseContext) *basil.ParseContext {
	var nilBlock *BlockValueRequired
	if b, ok := basil.Block(nilBlock).(basil.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i BlockValueRequiredInterpreter) Param(b basil.Block, name basil.ID) interface{} {
	switch name {
	case "id_field":
		return b.(*BlockValueRequired).IDField
	case "value":
		return b.(*BlockValueRequired).Value
	default:
		panic(fmt.Errorf("unexpected parameter %q in BlockValueRequired", name))
	}
}

func (i BlockValueRequiredInterpreter) SetParam(block basil.Block, name basil.ID, value interface{}) error {
	b := block.(*BlockValueRequired)
	switch name {
	case "value":
		b.Value = value.(string)
	}
	return nil
}

func (i BlockValueRequiredInterpreter) SetBlock(block basil.Block, name basil.ID, value interface{}) error {
	return nil
}
