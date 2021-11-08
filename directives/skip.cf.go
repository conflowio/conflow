// Code generated by Conflow. DO NOT EDIT.

package directives

import (
	"fmt"
	"github.com/conflowio/conflow/conflow"
	"github.com/conflowio/conflow/conflow/schema"
)

// SkipInterpreter is the conflow interpreter for the Skip block
type SkipInterpreter struct {
	s schema.Schema
}

func (i SkipInterpreter) Schema() schema.Schema {
	if i.s == nil {
		i.s = &schema.Object{
			Metadata: schema.Metadata{
				Annotations: map[string]string{"block.conflow.io/eval_stage": "init"},
			},
			Name: "Skip",
			Properties: map[string]schema.Schema{
				"id": &schema.String{
					Metadata: schema.Metadata{
						Annotations: map[string]string{"block.conflow.io/id": "true"},
						ReadOnly:    true,
					},
					Format: "conflow.ID",
				},
				"when": &schema.Boolean{
					Metadata: schema.Metadata{
						Annotations: map[string]string{"block.conflow.io/value": "true"},
					},
					Default: schema.BooleanPtr(true),
				},
			},
		}
	}
	return i.s
}

// Create creates a new Skip block
func (i SkipInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	return &Skip{
		id:   id,
		when: true,
	}
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i SkipInterpreter) ValueParamName() conflow.ID {
	return "when"
}

// ParseContext returns with the parse context for the block
func (i SkipInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *Skip
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i SkipInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "id":
		return b.(*Skip).id
	case "when":
		return b.(*Skip).when
	default:
		panic(fmt.Errorf("unexpected parameter %q in Skip", name))
	}
}

func (i SkipInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*Skip)
	switch name {
	case "when":
		b.when = value.(bool)
	}
	return nil
}

func (i SkipInterpreter) SetBlock(block conflow.Block, name conflow.ID, value interface{}) error {
	return nil
}
