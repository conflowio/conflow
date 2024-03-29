// Code generated by Conflow. DO NOT EDIT.

package fixtures

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
			ID: "github.com/conflowio/conflow/pkg/test/fixtures.BlockGenerator",
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
			"items": &schema.Array{
				Items: &schema.Any{},
			},
			"result": &schema.Reference{
				Metadata: schema.Metadata{
					Annotations: map[string]string{
						annotations.EvalStage: "init",
						annotations.Generated: "true",
					},
				},
				Nullable: true,
				Ref:      "github.com/conflowio/conflow/pkg/test/fixtures.BlockGeneratorResult",
			},
		},
		Required: []string{"items", "result"},
	})
}

// NewBlockGeneratorWithDefaults creates a new BlockGenerator instance with default values
func NewBlockGeneratorWithDefaults() *BlockGenerator {
	b := &BlockGenerator{}
	return b
}

// BlockGeneratorInterpreter is the Conflow interpreter for the BlockGenerator block
type BlockGeneratorInterpreter struct {
}

func (i BlockGeneratorInterpreter) Schema() schema.Schema {
	s, _ := schema.Get("github.com/conflowio/conflow/pkg/test/fixtures.BlockGenerator")
	return s
}

// Create creates a new BlockGenerator block
func (i BlockGeneratorInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	b := NewBlockGeneratorWithDefaults()
	b.id = id
	b.blockPublisher = blockCtx.BlockPublisher()
	return b
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i BlockGeneratorInterpreter) ValueParamName() conflow.ID {
	return ""
}

// ParseContext returns with the parse context for the block
func (i BlockGeneratorInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *BlockGenerator
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i BlockGeneratorInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "id":
		return b.(*BlockGenerator).id
	case "items":
		return b.(*BlockGenerator).items
	default:
		panic(fmt.Errorf("unexpected parameter %q in BlockGenerator", name))
	}
}

func (i BlockGeneratorInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*BlockGenerator)
	switch name {
	case "items":
		b.items = value.([]interface{})
	}
	return nil
}

func (i BlockGeneratorInterpreter) SetBlock(block conflow.Block, name conflow.ID, key string, value interface{}) error {
	b := block.(*BlockGenerator)
	switch name {
	case "result":
		b.result = value.(*BlockGeneratorResult)
	}
	return nil
}
