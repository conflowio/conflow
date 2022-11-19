// Code generated by Conflow. DO NOT EDIT.

package blocks

import (
	"fmt"

	"github.com/conflowio/conflow/pkg/conflow"
	"github.com/conflowio/conflow/pkg/conflow/annotations"
	"github.com/conflowio/conflow/pkg/schema"
)

func init() {
	schema.Register(&schema.Object{
		Metadata: schema.Metadata{
			Annotations: map[string]string{
				annotations.Type: "configuration",
			},
			ID: "github.com/conflowio/conflow/pkg/blocks.Line",
		},
		Properties: map[string]schema.Schema{
			"id": &schema.String{
				Metadata: schema.Metadata{
					Annotations: map[string]string{
						annotations.ID: "true",
					},
					ReadOnly: true,
				},
				Format: "conflow.ID",
			},
			"text": &schema.String{},
		},
	})
}

// NewLineWithDefaults creates a new Line instance with default values
func NewLineWithDefaults() *Line {
	b := &Line{}
	return b
}

// LineInterpreter is the Conflow interpreter for the Line block
type LineInterpreter struct {
}

func (i LineInterpreter) Schema() schema.Schema {
	s, _ := schema.Get("github.com/conflowio/conflow/pkg/blocks.Line")
	return s
}

// Create creates a new Line block
func (i LineInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	b := NewLineWithDefaults()
	b.id = id
	return b
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
		b.text = schema.Value[string](value)
	}
	return nil
}

func (i LineInterpreter) SetBlock(block conflow.Block, name conflow.ID, key string, value interface{}) error {
	return nil
}
