// Code generated by Conflow. DO NOT EDIT.

package directives

import (
	"fmt"
	"github.com/conflowio/conflow/src/conflow"
	"github.com/conflowio/conflow/src/schema"
)

func init() {
	schema.Register(&schema.Object{
		Metadata: schema.Metadata{
			Annotations: map[string]string{"block.conflow.io/type": "directive"},
			ID:          "github.com/conflowio/conflow/src/schema/directives.ExclusiveMaximum",
		},
		Name: "ExclusiveMaximum",
		Parameters: map[string]schema.Schema{
			"id": &schema.String{
				Metadata: schema.Metadata{
					Annotations: map[string]string{"block.conflow.io/id": "true"},
					ReadOnly:    true,
				},
				Format: "conflow.ID",
			},
			"value": &schema.Untyped{
				Metadata: schema.Metadata{
					Annotations: map[string]string{"block.conflow.io/value": "true"},
				},
			},
		},
		Required: []string{"value"},
	})
}

// ExclusiveMaximumInterpreter is the Conflow interpreter for the ExclusiveMaximum block
type ExclusiveMaximumInterpreter struct {
}

func (i ExclusiveMaximumInterpreter) Schema() schema.Schema {
	s, _ := schema.Get("github.com/conflowio/conflow/src/schema/directives.ExclusiveMaximum")
	return s
}

// Create creates a new ExclusiveMaximum block
func (i ExclusiveMaximumInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	b := &ExclusiveMaximum{}
	b.id = id
	return b
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i ExclusiveMaximumInterpreter) ValueParamName() conflow.ID {
	return "value"
}

// ParseContext returns with the parse context for the block
func (i ExclusiveMaximumInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *ExclusiveMaximum
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i ExclusiveMaximumInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "id":
		return b.(*ExclusiveMaximum).id
	case "value":
		return b.(*ExclusiveMaximum).value
	default:
		panic(fmt.Errorf("unexpected parameter %q in ExclusiveMaximum", name))
	}
}

func (i ExclusiveMaximumInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*ExclusiveMaximum)
	switch name {
	case "value":
		b.value = value
	}
	return nil
}

func (i ExclusiveMaximumInterpreter) SetBlock(block conflow.Block, name conflow.ID, key string, value interface{}) error {
	return nil
}
