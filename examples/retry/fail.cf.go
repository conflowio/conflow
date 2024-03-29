// Code generated by Conflow. DO NOT EDIT.

package main

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
				annotations.Type: "task",
			},
			Description: "It will error for the given tries",
			ID:          "github.com/conflowio/conflow/examples/retry.Fail",
		},
		ParameterNames: map[string]string{"triesRequired": "tries_required"},
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
			"tries": &schema.Integer{
				Metadata: schema.Metadata{
					Annotations: map[string]string{
						annotations.EvalStage: "close",
					},
					ReadOnly: true,
				},
			},
			"triesRequired": &schema.Integer{},
		},
		Required: []string{"tries_required"},
	})
}

// NewFailWithDefaults creates a new Fail instance with default values
func NewFailWithDefaults() *Fail {
	b := &Fail{}
	return b
}

// FailInterpreter is the Conflow interpreter for the Fail block
type FailInterpreter struct {
}

func (i FailInterpreter) Schema() schema.Schema {
	s, _ := schema.Get("github.com/conflowio/conflow/examples/retry.Fail")
	return s
}

// Create creates a new Fail block
func (i FailInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	b := NewFailWithDefaults()
	b.id = id
	return b
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i FailInterpreter) ValueParamName() conflow.ID {
	return ""
}

// ParseContext returns with the parse context for the block
func (i FailInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *Fail
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i FailInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "id":
		return b.(*Fail).id
	case "tries":
		return b.(*Fail).tries
	case "tries_required":
		return b.(*Fail).triesRequired
	default:
		panic(fmt.Errorf("unexpected parameter %q in Fail", name))
	}
}

func (i FailInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*Fail)
	switch name {
	case "tries_required":
		b.triesRequired = schema.Value[int64](value)
	}
	return nil
}

func (i FailInterpreter) SetBlock(block conflow.Block, name conflow.ID, key string, value interface{}) error {
	return nil
}
