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
				annotations.EvalStage: "init",
				annotations.Type:      "directive",
			},
			ID: "github.com/conflowio/conflow/pkg/directives.Retry",
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
			"limit": &schema.Integer{
				Metadata: schema.Metadata{
					Annotations: map[string]string{
						annotations.Value: "true",
					},
				},
				Default: schema.Pointer(int64(-1)),
				Minimum: schema.Pointer(int64(-1)),
				Maximum: schema.Pointer(int64(2147483647)),
			},
		},
		Required: []string{"limit"},
	})
}

// NewRetryWithDefaults creates a new Retry instance with default values
func NewRetryWithDefaults() *Retry {
	b := &Retry{}
	b.limit = -1
	return b
}

// RetryInterpreter is the Conflow interpreter for the Retry block
type RetryInterpreter struct {
}

func (i RetryInterpreter) Schema() schema.Schema {
	s, _ := schema.Get("github.com/conflowio/conflow/pkg/directives.Retry")
	return s
}

// Create creates a new Retry block
func (i RetryInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	b := NewRetryWithDefaults()
	b.id = id
	return b
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i RetryInterpreter) ValueParamName() conflow.ID {
	return "limit"
}

// ParseContext returns with the parse context for the block
func (i RetryInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *Retry
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i RetryInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "id":
		return b.(*Retry).id
	case "limit":
		return b.(*Retry).limit
	default:
		panic(fmt.Errorf("unexpected parameter %q in Retry", name))
	}
}

func (i RetryInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*Retry)
	switch name {
	case "limit":
		b.limit = schema.Value[int64](value)
	}
	return nil
}

func (i RetryInterpreter) SetBlock(block conflow.Block, name conflow.ID, key string, value interface{}) error {
	return nil
}
