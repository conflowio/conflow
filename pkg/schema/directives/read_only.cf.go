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
			ID: "github.com/conflowio/conflow/pkg/schema/directives.ReadOnly",
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

// ReadOnlyInterpreter is the Conflow interpreter for the ReadOnly block
type ReadOnlyInterpreter struct {
}

func (i ReadOnlyInterpreter) Schema() schema.Schema {
	s, _ := schema.Get("github.com/conflowio/conflow/pkg/schema/directives.ReadOnly")
	return s
}

// Create creates a new ReadOnly block
func (i ReadOnlyInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	b := &ReadOnly{}
	b.id = id
	return b
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i ReadOnlyInterpreter) ValueParamName() conflow.ID {
	return ""
}

// ParseContext returns with the parse context for the block
func (i ReadOnlyInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *ReadOnly
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i ReadOnlyInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "id":
		return b.(*ReadOnly).id
	default:
		panic(fmt.Errorf("unexpected parameter %q in ReadOnly", name))
	}
}

func (i ReadOnlyInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	return nil
}

func (i ReadOnlyInterpreter) SetBlock(block conflow.Block, name conflow.ID, key string, value interface{}) error {
	return nil
}