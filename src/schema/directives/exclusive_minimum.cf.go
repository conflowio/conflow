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
			ID:          "github.com/conflowio/conflow/src/schema/directives.ExclusiveMinimum",
		},
		Name: "ExclusiveMinimum",
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

// ExclusiveMinimumInterpreter is the Conflow interpreter for the ExclusiveMinimum block
type ExclusiveMinimumInterpreter struct {
}

func (i ExclusiveMinimumInterpreter) Schema() schema.Schema {
	s, _ := schema.Get("github.com/conflowio/conflow/src/schema/directives.ExclusiveMinimum")
	return s
}

// Create creates a new ExclusiveMinimum block
func (i ExclusiveMinimumInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	b := &ExclusiveMinimum{}
	b.id = id
	return b
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i ExclusiveMinimumInterpreter) ValueParamName() conflow.ID {
	return "value"
}

// ParseContext returns with the parse context for the block
func (i ExclusiveMinimumInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *ExclusiveMinimum
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i ExclusiveMinimumInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "id":
		return b.(*ExclusiveMinimum).id
	case "value":
		return b.(*ExclusiveMinimum).value
	default:
		panic(fmt.Errorf("unexpected parameter %q in ExclusiveMinimum", name))
	}
}

func (i ExclusiveMinimumInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*ExclusiveMinimum)
	switch name {
	case "value":
		b.value = value
	}
	return nil
}

func (i ExclusiveMinimumInterpreter) SetBlock(block conflow.Block, name conflow.ID, value interface{}) error {
	return nil
}
