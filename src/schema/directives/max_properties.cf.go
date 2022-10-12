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
			ID:          "github.com/conflowio/conflow/src/schema/directives.MaxProperties",
		},
		Name: "MaxProperties",
		Parameters: map[string]schema.Schema{
			"id": &schema.String{
				Metadata: schema.Metadata{
					Annotations: map[string]string{"block.conflow.io/id": "true"},
					ReadOnly:    true,
				},
				Format: "conflow.ID",
			},
			"value": &schema.Integer{
				Metadata: schema.Metadata{
					Annotations: map[string]string{"block.conflow.io/value": "true"},
				},
			},
		},
		Required: []string{"value"},
	})
}

// MaxPropertiesInterpreter is the Conflow interpreter for the MaxProperties block
type MaxPropertiesInterpreter struct {
}

func (i MaxPropertiesInterpreter) Schema() schema.Schema {
	s, _ := schema.Get("github.com/conflowio/conflow/src/schema/directives.MaxProperties")
	return s
}

// Create creates a new MaxProperties block
func (i MaxPropertiesInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	b := &MaxProperties{}
	b.id = id
	return b
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i MaxPropertiesInterpreter) ValueParamName() conflow.ID {
	return "value"
}

// ParseContext returns with the parse context for the block
func (i MaxPropertiesInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *MaxProperties
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i MaxPropertiesInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "id":
		return b.(*MaxProperties).id
	case "value":
		return b.(*MaxProperties).value
	default:
		panic(fmt.Errorf("unexpected parameter %q in MaxProperties", name))
	}
}

func (i MaxPropertiesInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*MaxProperties)
	switch name {
	case "value":
		b.value = value.(int64)
	}
	return nil
}

func (i MaxPropertiesInterpreter) SetBlock(block conflow.Block, name conflow.ID, value interface{}) error {
	return nil
}
