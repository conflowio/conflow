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
			ID:          "github.com/conflowio/conflow/src/schema/directives.Minimum",
		},
		Name: "Minimum",
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

// MinimumInterpreter is the Conflow interpreter for the Minimum block
type MinimumInterpreter struct {
}

func (i MinimumInterpreter) Schema() schema.Schema {
	s, _ := schema.Get("github.com/conflowio/conflow/src/schema/directives.Minimum")
	return s
}

// Create creates a new Minimum block
func (i MinimumInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	b := &Minimum{}
	b.id = id
	return b
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i MinimumInterpreter) ValueParamName() conflow.ID {
	return "value"
}

// ParseContext returns with the parse context for the block
func (i MinimumInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *Minimum
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i MinimumInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "id":
		return b.(*Minimum).id
	case "value":
		return b.(*Minimum).value
	default:
		panic(fmt.Errorf("unexpected parameter %q in Minimum", name))
	}
}

func (i MinimumInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*Minimum)
	switch name {
	case "value":
		b.value = value
	}
	return nil
}

func (i MinimumInterpreter) SetBlock(block conflow.Block, name conflow.ID, key string, value interface{}) error {
	return nil
}
