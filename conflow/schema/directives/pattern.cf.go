// Code generated by Conflow. DO NOT EDIT.

package directives

import (
	"fmt"
	"github.com/conflowio/conflow/conflow"
	"github.com/conflowio/conflow/conflow/schema"
)

// PatternInterpreter is the conflow interpreter for the Pattern block
type PatternInterpreter struct {
	s schema.Schema
}

func (i PatternInterpreter) Schema() schema.Schema {
	if i.s == nil {
		i.s = &schema.Object{
			Name: "Pattern",
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
					Format: "regex",
				},
			},
			Required: []string{"value"},
		}
	}
	return i.s
}

// Create creates a new Pattern block
func (i PatternInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	return &Pattern{
		id: id,
	}
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i PatternInterpreter) ValueParamName() conflow.ID {
	return "value"
}

// ParseContext returns with the parse context for the block
func (i PatternInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *Pattern
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i PatternInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "id":
		return b.(*Pattern).id
	case "value":
		return b.(*Pattern).value
	default:
		panic(fmt.Errorf("unexpected parameter %q in Pattern", name))
	}
}

func (i PatternInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*Pattern)
	switch name {
	case "value":
		b.value = value.(string)
	}
	return nil
}

func (i PatternInterpreter) SetBlock(block conflow.Block, name conflow.ID, value interface{}) error {
	return nil
}
