// Code generated by Conflow. DO NOT EDIT.

package fixtures

import (
	"fmt"
	"github.com/conflowio/conflow/src/conflow"
	"github.com/conflowio/conflow/src/schema"
)

// BlockValueRequiredInterpreter is the conflow interpreter for the BlockValueRequired block
type BlockValueRequiredInterpreter struct {
	s schema.Schema
}

func (i BlockValueRequiredInterpreter) Schema() schema.Schema {
	if i.s == nil {
		i.s = &schema.Object{
			Metadata: schema.Metadata{
				Annotations: map[string]string{"block.conflow.io/type": "configuration"},
			},
			JSONPropertyNames: map[string]string{"id_field": "IDField", "value": "Value"},
			Name:              "BlockValueRequired",
			Parameters: map[string]schema.Schema{
				"id_field": &schema.String{
					Metadata: schema.Metadata{
						Annotations: map[string]string{"block.conflow.io/id": "true"},
						ReadOnly:    true,
					},
					Format: "conflow.ID",
				},
				"value": &schema.String{
					Metadata: schema.Metadata{
						Annotations: map[string]string{"block.conflow.io/value": "true"},
					},
				},
			},
			Required: []string{"value"},
		}
	}
	return i.s
}

// Create creates a new BlockValueRequired block
func (i BlockValueRequiredInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	return &BlockValueRequired{
		IDField: id,
	}
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i BlockValueRequiredInterpreter) ValueParamName() conflow.ID {
	return "value"
}

// ParseContext returns with the parse context for the block
func (i BlockValueRequiredInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *BlockValueRequired
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i BlockValueRequiredInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "id_field":
		return b.(*BlockValueRequired).IDField
	case "value":
		return b.(*BlockValueRequired).Value
	default:
		panic(fmt.Errorf("unexpected parameter %q in BlockValueRequired", name))
	}
}

func (i BlockValueRequiredInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*BlockValueRequired)
	switch name {
	case "value":
		b.Value = value.(string)
	}
	return nil
}

func (i BlockValueRequiredInterpreter) SetBlock(block conflow.Block, name conflow.ID, value interface{}) error {
	return nil
}
