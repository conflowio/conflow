// Code generated by Conflow. DO NOT EDIT.

package acceptance

import (
	"fmt"
	"github.com/conflowio/conflow/conflow"
	"github.com/conflowio/conflow/conflow/schema"
)

// MainInterpreter is the conflow interpreter for the Main block
type MainInterpreter struct {
	s schema.Schema
}

func (i MainInterpreter) Schema() schema.Schema {
	if i.s == nil {
		i.s = &schema.Object{
			Name: "Main",
			Parameters: map[string]schema.Schema{
				"id": &schema.String{
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

// Create creates a new Main block
func (i MainInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	return &Main{
		id: id,
	}
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i MainInterpreter) ValueParamName() conflow.ID {
	return ""
}

// ParseContext returns with the parse context for the block
func (i MainInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *Main
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i MainInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "id":
		return b.(*Main).id
	default:
		panic(fmt.Errorf("unexpected parameter %q in Main", name))
	}
}

func (i MainInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	return nil
}

func (i MainInterpreter) SetBlock(block conflow.Block, name conflow.ID, value interface{}) error {
	return nil
}
