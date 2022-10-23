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
				annotations.EvalStage: "init",
				annotations.Type:      "directive",
			},
			ID: "github.com/conflowio/conflow/pkg/directives.Run",
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
			"when": &schema.Boolean{
				Metadata: schema.Metadata{
					Annotations: map[string]string{
						annotations.Value: "true",
					},
				},
				Default: schema.Pointer(true),
			},
		},
	})
}

// RunInterpreter is the Conflow interpreter for the Run block
type RunInterpreter struct {
}

func (i RunInterpreter) Schema() schema.Schema {
	s, _ := schema.Get("github.com/conflowio/conflow/pkg/directives.Run")
	return s
}

// Create creates a new Run block
func (i RunInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	b := &Run{}
	b.id = id
	b.when = true
	return b
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i RunInterpreter) ValueParamName() conflow.ID {
	return "when"
}

// ParseContext returns with the parse context for the block
func (i RunInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *Run
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i RunInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "id":
		return b.(*Run).id
	case "when":
		return b.(*Run).when
	default:
		panic(fmt.Errorf("unexpected parameter %q in Run", name))
	}
}

func (i RunInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*Run)
	switch name {
	case "when":
		b.when = value.(bool)
	}
	return nil
}

func (i RunInterpreter) SetBlock(block conflow.Block, name conflow.ID, key string, value interface{}) error {
	return nil
}