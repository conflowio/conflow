// Code generated by Conflow. DO NOT EDIT.

package blocks

import (
	"fmt"
	"github.com/conflowio/conflow/conflow"
	"github.com/conflowio/conflow/conflow/schema"
)

// LineInterpreter is the conflow interpreter for the Line block
type LineInterpreter struct {
	s schema.Schema
}

func (i LineInterpreter) Schema() schema.Schema {
	if i.s == nil {
		i.s = &schema.Object{
			Metadata: schema.Metadata{
				Annotations: map[string]string{"block.conflow.io/type": "configuration"},
			},
			Name: "Line",
			Parameters: map[string]schema.Schema{
				"id": &schema.String{
					Metadata: schema.Metadata{
						Annotations: map[string]string{"block.conflow.io/id": "true"},
						ReadOnly:    true,
					},
					Format: "conflow.ID",
				},
				"text": &schema.String{},
			},
		}
	}
	return i.s
}

// Create creates a new Line block
func (i LineInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	return &Line{
		id: id,
	}
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i LineInterpreter) ValueParamName() conflow.ID {
	return ""
}

// ParseContext returns with the parse context for the block
func (i LineInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *Line
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i LineInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "id":
		return b.(*Line).id
	case "text":
		return b.(*Line).text
	default:
		panic(fmt.Errorf("unexpected parameter %q in Line", name))
	}
}

func (i LineInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*Line)
	switch name {
	case "text":
		b.text = value.(string)
	}
	return nil
}

func (i LineInterpreter) SetBlock(block conflow.Block, name conflow.ID, value interface{}) error {
	return nil
}
