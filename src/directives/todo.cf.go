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
				annotations.EvalStage: "ignore",
				annotations.Type:      "directive",
			},
			ID: "github.com/conflowio/conflow/src/directives.Todo",
		},
		Properties: map[string]schema.Schema{
			"description": &schema.String{
				Metadata: schema.Metadata{
					Annotations: map[string]string{
						annotations.Value: "true",
					},
				},
			},
			"id": &schema.String{
				Metadata: schema.Metadata{
					Annotations: map[string]string{
						annotations.ID: "true",
					},
					ReadOnly: true,
				},
				Format: "conflow.ID",
			},
		},
		Required: []string{"description"},
	})
}

// TodoInterpreter is the Conflow interpreter for the Todo block
type TodoInterpreter struct {
}

func (i TodoInterpreter) Schema() schema.Schema {
	s, _ := schema.Get("github.com/conflowio/conflow/src/directives.Todo")
	return s
}

// Create creates a new Todo block
func (i TodoInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	b := &Todo{}
	b.id = id
	return b
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i TodoInterpreter) ValueParamName() conflow.ID {
	return "description"
}

// ParseContext returns with the parse context for the block
func (i TodoInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *Todo
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i TodoInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "description":
		return b.(*Todo).description
	case "id":
		return b.(*Todo).id
	default:
		panic(fmt.Errorf("unexpected parameter %q in Todo", name))
	}
}

func (i TodoInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*Todo)
	switch name {
	case "description":
		b.description = value.(string)
	}
	return nil
}

func (i TodoInterpreter) SetBlock(block conflow.Block, name conflow.ID, key string, value interface{}) error {
	return nil
}
