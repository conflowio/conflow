// Code generated by Conflow. DO NOT EDIT.

package directives

import (
	"fmt"

	"github.com/conflowio/conflow/pkg/conflow"
	"github.com/conflowio/conflow/pkg/conflow/annotations"
	"github.com/conflowio/conflow/pkg/conflow/types"
	"github.com/conflowio/conflow/pkg/schema"
)

func init() {
	schema.Register(&schema.Object{
		Metadata: schema.Metadata{
			Annotations: map[string]string{
				annotations.EvalStage: "init",
				annotations.Type:      "directive",
			},
			ID: "github.com/conflowio/conflow/pkg/directives.Timeout",
		},
		Properties: map[string]schema.Schema{
			"duration": &schema.String{
				Metadata: schema.Metadata{
					Annotations: map[string]string{
						annotations.Value: "true",
					},
				},
				Format: "duration-go",
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
		Required: []string{"duration"},
	})
}

// NewTimeoutWithDefaults creates a new Timeout instance with default values
func NewTimeoutWithDefaults() *Timeout {
	b := &Timeout{}
	return b
}

// TimeoutInterpreter is the Conflow interpreter for the Timeout block
type TimeoutInterpreter struct {
}

func (i TimeoutInterpreter) Schema() schema.Schema {
	s, _ := schema.Get("github.com/conflowio/conflow/pkg/directives.Timeout")
	return s
}

// Create creates a new Timeout block
func (i TimeoutInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	b := NewTimeoutWithDefaults()
	b.id = id
	return b
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i TimeoutInterpreter) ValueParamName() conflow.ID {
	return "duration"
}

// ParseContext returns with the parse context for the block
func (i TimeoutInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *Timeout
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i TimeoutInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "duration":
		return b.(*Timeout).duration
	case "id":
		return b.(*Timeout).id
	default:
		panic(fmt.Errorf("unexpected parameter %q in Timeout", name))
	}
}

func (i TimeoutInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*Timeout)
	switch name {
	case "duration":
		b.duration = schema.Value[types.Duration](value)
	}
	return nil
}

func (i TimeoutInterpreter) SetBlock(block conflow.Block, name conflow.ID, key string, value interface{}) error {
	return nil
}
