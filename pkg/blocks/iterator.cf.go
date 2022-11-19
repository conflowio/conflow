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
				annotations.Type: "generator",
			},
			ID: "github.com/conflowio/conflow/pkg/blocks.Iterator",
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
			"it": &schema.Reference{
				Metadata: schema.Metadata{
					Annotations: map[string]string{
						annotations.EvalStage: "init",
						annotations.Generated: "true",
					},
				},
				Nullable: true,
				Ref:      "github.com/conflowio/conflow/pkg/blocks.It",
			},
			"value": &schema.Any{},
		},
		Required: []string{"it"},
	})
}

// NewIteratorWithDefaults creates a new Iterator instance with default values
func NewIteratorWithDefaults() *Iterator {
	b := &Iterator{}
	return b
}

// IteratorInterpreter is the Conflow interpreter for the Iterator block
type IteratorInterpreter struct {
}

func (i IteratorInterpreter) Schema() schema.Schema {
	s, _ := schema.Get("github.com/conflowio/conflow/pkg/blocks.Iterator")
	return s
}

// Create creates a new Iterator block
func (i IteratorInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	b := NewIteratorWithDefaults()
	b.id = id
	b.blockPublisher = blockCtx.BlockPublisher()
	return b
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i IteratorInterpreter) ValueParamName() conflow.ID {
	return ""
}

// ParseContext returns with the parse context for the block
func (i IteratorInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *Iterator
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i IteratorInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "id":
		return b.(*Iterator).id
	case "value":
		return b.(*Iterator).value
	default:
		panic(fmt.Errorf("unexpected parameter %q in Iterator", name))
	}
}

func (i IteratorInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*Iterator)
	switch name {
	case "value":
		b.value = value
	}
	return nil
}

func (i IteratorInterpreter) SetBlock(block conflow.Block, name conflow.ID, key string, value interface{}) error {
	b := block.(*Iterator)
	switch name {
	case "it":
		b.it = value.(*It)
	}
	return nil
}
