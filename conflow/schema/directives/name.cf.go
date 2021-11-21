// Code generated by Conflow. DO NOT EDIT.

package directives

import (
	"fmt"
	"github.com/conflowio/conflow/conflow"
	"github.com/conflowio/conflow/conflow/schema"
)

// NameInterpreter is the conflow interpreter for the Name block
type NameInterpreter struct {
	s schema.Schema
}

func (i NameInterpreter) Schema() schema.Schema {
	if i.s == nil {
		i.s = &schema.Object{
			Metadata: schema.Metadata{
				Annotations: map[string]string{"block.conflow.io/type": "directive"},
			},
			JSONPropertyNames: map[string]string{"value": "Value"},
			Name:              "Name",
			Parameters: map[string]schema.Schema{
				"id": &schema.String{
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

// Create creates a new Name block
func (i NameInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	return &Name{
		id: id,
	}
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i NameInterpreter) ValueParamName() conflow.ID {
	return "value"
}

// ParseContext returns with the parse context for the block
func (i NameInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *Name
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i NameInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "id":
		return b.(*Name).id
	case "value":
		return b.(*Name).Value
	default:
		panic(fmt.Errorf("unexpected parameter %q in Name", name))
	}
}

func (i NameInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*Name)
	switch name {
	case "value":
		b.Value = value.(string)
	}
	return nil
}

func (i NameInterpreter) SetBlock(block conflow.Block, name conflow.ID, value interface{}) error {
	return nil
}
