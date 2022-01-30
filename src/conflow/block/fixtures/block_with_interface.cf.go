// Code generated by Conflow. DO NOT EDIT.

package fixtures

import (
	"fmt"
	"github.com/conflowio/conflow/src/conflow"
	"github.com/conflowio/conflow/src/schema"
)

// BlockWithInterfaceInterpreter is the conflow interpreter for the BlockWithInterface block
type BlockWithInterfaceInterpreter struct {
	s schema.Schema
}

func (i BlockWithInterfaceInterpreter) Schema() schema.Schema {
	if i.s == nil {
		i.s = &schema.Object{
			Metadata: schema.Metadata{
				Annotations: map[string]string{"block.conflow.io/type": "configuration"},
			},
			JSONPropertyNames: map[string]string{"block": "Block", "blocks": "Blocks", "id_field": "IDField"},
			Name:              "BlockWithInterface",
			Parameters: map[string]schema.Schema{
				"block": &schema.Reference{
					Ref: "http://conflow.schema/github.com/conflowio/conflow/src/conflow.Block",
				},
				"blocks": &schema.Array{
					Items: &schema.Reference{
						Ref: "http://conflow.schema/github.com/conflowio/conflow/src/conflow.Block",
					},
				},
				"id_field": &schema.String{
					Metadata: schema.Metadata{
						Annotations: map[string]string{"block.conflow.io/id": "true"},
						ReadOnly:    true,
					},
					Format: "conflow.ID",
				},
			},
		}
	}
	return i.s
}

// Create creates a new BlockWithInterface block
func (i BlockWithInterfaceInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	return &BlockWithInterface{
		IDField: id,
	}
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i BlockWithInterfaceInterpreter) ValueParamName() conflow.ID {
	return ""
}

// ParseContext returns with the parse context for the block
func (i BlockWithInterfaceInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *BlockWithInterface
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i BlockWithInterfaceInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "id_field":
		return b.(*BlockWithInterface).IDField
	default:
		panic(fmt.Errorf("unexpected parameter %q in BlockWithInterface", name))
	}
}

func (i BlockWithInterfaceInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	return nil
}

func (i BlockWithInterfaceInterpreter) SetBlock(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*BlockWithInterface)
	switch name {
	case "block":
		b.Block = value.(conflow.Block)
	case "blocks":
		b.Blocks = append(b.Blocks, value.(conflow.Block))
	}
	return nil
}
