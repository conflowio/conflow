// Code generated by Conflow. DO NOT EDIT.

package blocks

import (
	"fmt"
	"github.com/conflowio/conflow/conflow"
	"github.com/conflowio/conflow/conflow/schema"
)

// RangeEntryInterpreter is the conflow interpreter for the RangeEntry block
type RangeEntryInterpreter struct {
	s schema.Schema
}

func (i RangeEntryInterpreter) Schema() schema.Schema {
	if i.s == nil {
		i.s = &schema.Object{
			Name: "RangeEntry",
			Parameters: map[string]schema.Schema{
				"id": &schema.String{
					Metadata: schema.Metadata{
						Annotations: map[string]string{"block.conflow.io/id": "true"},
						ReadOnly:    true,
					},
					Format: "conflow.ID",
				},
				"key": &schema.Untyped{
					Metadata: schema.Metadata{
						Annotations: map[string]string{"block.conflow.io/eval_stage": "close"},
						ReadOnly:    true,
					},
				},
				"value": &schema.Untyped{
					Metadata: schema.Metadata{
						Annotations: map[string]string{"block.conflow.io/eval_stage": "close"},
						ReadOnly:    true,
					},
				},
			},
		}
	}
	return i.s
}

// Create creates a new RangeEntry block
func (i RangeEntryInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	return &RangeEntry{
		id: id,
	}
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i RangeEntryInterpreter) ValueParamName() conflow.ID {
	return ""
}

// ParseContext returns with the parse context for the block
func (i RangeEntryInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *RangeEntry
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i RangeEntryInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "id":
		return b.(*RangeEntry).id
	case "key":
		return b.(*RangeEntry).key
	case "value":
		return b.(*RangeEntry).value
	default:
		panic(fmt.Errorf("unexpected parameter %q in RangeEntry", name))
	}
}

func (i RangeEntryInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	return nil
}

func (i RangeEntryInterpreter) SetBlock(block conflow.Block, name conflow.ID, value interface{}) error {
	return nil
}
