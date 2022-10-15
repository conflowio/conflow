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
				annotations.Type: "directive",
			},
			ID: "github.com/conflowio/conflow/src/schema/directives.Types",
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
			"value": &schema.Array{
				Metadata: schema.Metadata{
					Annotations: map[string]string{
						annotations.Value: "true",
					},
				},
				Items: &schema.String{},
			},
		},
		Required: []string{"value"},
	})
}

// TypesInterpreter is the Conflow interpreter for the Types block
type TypesInterpreter struct {
}

func (i TypesInterpreter) Schema() schema.Schema {
	s, _ := schema.Get("github.com/conflowio/conflow/src/schema/directives.Types")
	return s
}

// Create creates a new Types block
func (i TypesInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	b := &Types{}
	b.id = id
	return b
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i TypesInterpreter) ValueParamName() conflow.ID {
	return "value"
}

// ParseContext returns with the parse context for the block
func (i TypesInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *Types
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i TypesInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "id":
		return b.(*Types).id
	case "value":
		return b.(*Types).value
	default:
		panic(fmt.Errorf("unexpected parameter %q in Types", name))
	}
}

func (i TypesInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*Types)
	switch name {
	case "value":
		b.value = make([]string, len(value.([]interface{})))
		for valuek, valuev := range value.([]interface{}) {
			b.value[valuek] = valuev.(string)
		}
	}
	return nil
}

func (i TypesInterpreter) SetBlock(block conflow.Block, name conflow.ID, key string, value interface{}) error {
	return nil
}
