// Code generated by Conflow. DO NOT EDIT.

package directives

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
				annotations.Type: "directive",
			},
			ID: "github.com/conflowio/conflow/pkg/schema/directives.ID",
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
		},
	})
}

// NewIDWithDefaults creates a new ID instance with default values
func NewIDWithDefaults() *ID {
	b := &ID{}
	return b
}

// IDInterpreter is the Conflow interpreter for the ID block
type IDInterpreter struct {
}

func (i IDInterpreter) Schema() schema.Schema {
	s, _ := schema.Get("github.com/conflowio/conflow/pkg/schema/directives.ID")
	return s
}

// Create creates a new ID block
func (i IDInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	b := NewIDWithDefaults()
	b.id = id
	return b
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i IDInterpreter) ValueParamName() conflow.ID {
	return ""
}

// ParseContext returns with the parse context for the block
func (i IDInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *ID
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i IDInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "id":
		return b.(*ID).id
	default:
		panic(fmt.Errorf("unexpected parameter %q in ID", name))
	}
}

func (i IDInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	return nil
}

func (i IDInterpreter) SetBlock(block conflow.Block, name conflow.ID, key string, value interface{}) error {
	return nil
}
