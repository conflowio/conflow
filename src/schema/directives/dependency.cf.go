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
			ID:          "github.com/conflowio/conflow/src/schema/directives.Dependency",
		},
		JSONPropertyNames: map[string]string{"name": "Name"},
		Name:              "Dependency",
		Parameters: map[string]schema.Schema{
			"id": &schema.String{
				Metadata: schema.Metadata{
					Annotations: map[string]string{"block.conflow.io/id": "true"},
					ReadOnly:    true,
				},
				Format: "conflow.ID",
			},
			"name": &schema.String{
				Metadata: schema.Metadata{
					Annotations: map[string]string{"block.conflow.io/value": "true"},
				},
			},
		},
	})
}

// DependencyInterpreter is the Conflow interpreter for the Dependency block
type DependencyInterpreter struct {
}

func (i DependencyInterpreter) Schema() schema.Schema {
	s, _ := schema.Get("github.com/conflowio/conflow/src/schema/directives.Dependency")
	return s
}

// Create creates a new Dependency block
func (i DependencyInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	b := &Dependency{}
	b.id = id
	return b
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i DependencyInterpreter) ValueParamName() conflow.ID {
	return "name"
}

// ParseContext returns with the parse context for the block
func (i DependencyInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *Dependency
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i DependencyInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "id":
		return b.(*Dependency).id
	case "name":
		return b.(*Dependency).Name
	default:
		panic(fmt.Errorf("unexpected parameter %q in Dependency", name))
	}
}

func (i DependencyInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*Dependency)
	switch name {
	case "name":
		b.Name = value.(string)
	}
	return nil
}

func (i DependencyInterpreter) SetBlock(block conflow.Block, name conflow.ID, value interface{}) error {
	return nil
}
