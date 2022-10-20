// Code generated by Conflow. DO NOT EDIT.

package directives

import (
	"fmt"
	"github.com/conflowio/conflow/src/conflow"
	"github.com/conflowio/conflow/src/conflow/annotations"
	"github.com/conflowio/conflow/src/schema"
)

func init() {
	schema.Register(&schema.Object{
		Metadata: schema.Metadata{
			Annotations: map[string]string{
				annotations.EvalStage: "parse",
				annotations.Type:      "directive",
			},
			ID: "github.com/conflowio/conflow/src/directives.Output",
		},
		ParameterNames: map[string]string{"schema": "type"},
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
			"schema": &schema.Reference{
				Ref: "github.com/conflowio/conflow/src/schema.Schema",
			},
		},
		Required: []string{"type"},
	})
}

// OutputInterpreter is the Conflow interpreter for the Output block
type OutputInterpreter struct {
}

func (i OutputInterpreter) Schema() schema.Schema {
	s, _ := schema.Get("github.com/conflowio/conflow/src/directives.Output")
	return s
}

// Create creates a new Output block
func (i OutputInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	b := &Output{}
	b.id = id
	return b
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i OutputInterpreter) ValueParamName() conflow.ID {
	return ""
}

// ParseContext returns with the parse context for the block
func (i OutputInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *Output
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i OutputInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "id":
		return b.(*Output).id
	default:
		panic(fmt.Errorf("unexpected parameter %q in Output", name))
	}
}

func (i OutputInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	return nil
}

func (i OutputInterpreter) SetBlock(block conflow.Block, name conflow.ID, key string, value interface{}) error {
	b := block.(*Output)
	switch name {
	case "type":
		b.schema = value.(schema.Schema)
	}
	return nil
}
