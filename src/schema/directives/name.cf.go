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
			ID: "github.com/conflowio/conflow/src/schema/directives.Name",
		},
		JSONPropertyNames: map[string]string{"value": "Value"},
		ParameterNames:    map[string]string{"Value": "value"},
		Properties: map[string]schema.Schema{
			"Value": &schema.String{
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
		Required: []string{"value"},
	})
}

// NameInterpreter is the Conflow interpreter for the Name block
type NameInterpreter struct {
}

func (i NameInterpreter) Schema() schema.Schema {
	s, _ := schema.Get("github.com/conflowio/conflow/src/schema/directives.Name")
	return s
}

// Create creates a new Name block
func (i NameInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	b := &Name{}
	b.id = id
	return b
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i NameInterpreter) ValueParamName() conflow.ID {
	return "value"
}

// ParseContext returns with the parse context for the block
func (i NameInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *Name
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i NameInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "value":
		return b.(*Name).Value
	case "id":
		return b.(*Name).id
	default:
		panic(fmt.Errorf("unexpected parameter %q in Name", name))
	}
}

func (i NameInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*Name)
	switch name {
	case "value":
		b.Value = value.(string)
	}
	return nil
}

func (i NameInterpreter) SetBlock(block conflow.Block, name conflow.ID, key string, value interface{}) error {
	return nil
}
