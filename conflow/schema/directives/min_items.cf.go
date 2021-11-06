// Code generated by Conflow. DO NOT EDIT.

package directives

import (
	"fmt"
	"github.com/conflowio/conflow/conflow"
	"github.com/conflowio/conflow/conflow/schema"
)

// MinItemsInterpreter is the conflow interpreter for the MinItems block
type MinItemsInterpreter struct {
	s schema.Schema
}

func (i MinItemsInterpreter) Schema() schema.Schema {
	if i.s == nil {
		i.s = &schema.Object{
			Name: "MinItems",
			Properties: map[string]schema.Schema{
				"id": &schema.String{
					Metadata: schema.Metadata{
						Annotations: map[string]string{"id": "true"},
						ReadOnly:    true,
					},
					Format: "conflow.ID",
				},
				"value": &schema.Integer{
					Metadata: schema.Metadata{
						Annotations: map[string]string{"value": "true"},
					},
				},
			},
			Required: []string{"value"},
		}
	}
	return i.s
}

// Create creates a new MinItems block
func (i MinItemsInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	return &MinItems{
		id: id,
	}
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i MinItemsInterpreter) ValueParamName() conflow.ID {
	return "value"
}

// ParseContext returns with the parse context for the block
func (i MinItemsInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *MinItems
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i MinItemsInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "id":
		return b.(*MinItems).id
	case "value":
		return b.(*MinItems).value
	default:
		panic(fmt.Errorf("unexpected parameter %q in MinItems", name))
	}
}

func (i MinItemsInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*MinItems)
	switch name {
	case "value":
		b.value = value.(int64)
	}
	return nil
}

func (i MinItemsInterpreter) SetBlock(block conflow.Block, name conflow.ID, value interface{}) error {
	return nil
}