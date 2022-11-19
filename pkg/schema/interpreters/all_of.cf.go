// Code generated by Conflow. DO NOT EDIT.

package interpreters

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
				annotations.Type: "configuration",
			},
			ID: "github.com/conflowio/conflow/pkg/schema.AllOf",
		},
		FieldNames:     map[string]string{"allOf": "Schemas"},
		ParameterNames: map[string]string{"allOf": "schema"},
		Properties: map[string]schema.Schema{
			"allOf": &schema.Array{
				Items: &schema.Reference{
					Ref: "github.com/conflowio/conflow/pkg/schema.Schema",
				},
				MinItems: 1,
			},
		},
		Required: []string{"schema"},
	})
}

// NewAllOfWithDefaults creates a new AllOf instance with default values
func NewAllOfWithDefaults() *schema.AllOf {
	b := &schema.AllOf{}
	return b
}

// AllOfInterpreter is the Conflow interpreter for the AllOf block
type AllOfInterpreter struct {
}

func (i AllOfInterpreter) Schema() schema.Schema {
	s, _ := schema.Get("github.com/conflowio/conflow/pkg/schema.AllOf")
	return s
}

// Create creates a new AllOf block
func (i AllOfInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	b := NewAllOfWithDefaults()
	return b
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i AllOfInterpreter) ValueParamName() conflow.ID {
	return ""
}

// ParseContext returns with the parse context for the block
func (i AllOfInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *schema.AllOf
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i AllOfInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	default:
		panic(fmt.Errorf("unexpected parameter %q in AllOf", name))
	}
}

func (i AllOfInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	return nil
}

func (i AllOfInterpreter) SetBlock(block conflow.Block, name conflow.ID, key string, value interface{}) error {
	b := block.(*schema.AllOf)
	switch name {
	case "schema":
		b.Schemas = append(b.Schemas, value.(schema.Schema))
	}
	return nil
}
