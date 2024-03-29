// Code generated by Conflow. DO NOT EDIT.

package directives

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
				annotations.Type: "directive",
			},
			ID: "github.com/conflowio/conflow/pkg/schema/directives.EvalStage",
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
			"value": &schema.String{
				Metadata: schema.Metadata{
					Annotations: map[string]string{
						annotations.Value: "true",
					},
				},
				Enum: []string{"init", "main", "close"},
			},
		},
		Required: []string{"value"},
	})
}

// NewEvalStageWithDefaults creates a new EvalStage instance with default values
func NewEvalStageWithDefaults() *EvalStage {
	b := &EvalStage{}
	return b
}

// EvalStageInterpreter is the Conflow interpreter for the EvalStage block
type EvalStageInterpreter struct {
}

func (i EvalStageInterpreter) Schema() schema.Schema {
	s, _ := schema.Get("github.com/conflowio/conflow/pkg/schema/directives.EvalStage")
	return s
}

// Create creates a new EvalStage block
func (i EvalStageInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	b := NewEvalStageWithDefaults()
	b.id = id
	return b
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i EvalStageInterpreter) ValueParamName() conflow.ID {
	return "value"
}

// ParseContext returns with the parse context for the block
func (i EvalStageInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *EvalStage
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i EvalStageInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "id":
		return b.(*EvalStage).id
	case "value":
		return b.(*EvalStage).value
	default:
		panic(fmt.Errorf("unexpected parameter %q in EvalStage", name))
	}
}

func (i EvalStageInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*EvalStage)
	switch name {
	case "value":
		b.value = schema.Value[string](value)
	}
	return nil
}

func (i EvalStageInterpreter) SetBlock(block conflow.Block, name conflow.ID, key string, value interface{}) error {
	return nil
}
