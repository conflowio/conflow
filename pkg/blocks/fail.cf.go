// Code generated by Conflow. DO NOT EDIT.

package blocks

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
			ID: "github.com/conflowio/conflow/pkg/blocks.Fail",
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
			"msg": &schema.String{
				Metadata: schema.Metadata{
					Annotations: map[string]string{
						annotations.Value: "true",
					},
				},
			},
		},
		Required: []string{"msg"},
	})
}

// FailInterpreter is the Conflow interpreter for the Fail block
type FailInterpreter struct {
}

func (i FailInterpreter) Schema() schema.Schema {
	s, _ := schema.Get("github.com/conflowio/conflow/pkg/blocks.Fail")
	return s
}

// Create creates a new Fail block
func (i FailInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	b := &Fail{}
	b.id = id
	return b
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i FailInterpreter) ValueParamName() conflow.ID {
	return "msg"
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
	case "msg":
		return b.(*Fail).msg
	default:
		panic(fmt.Errorf("unexpected parameter %q in Fail", name))
	}
}

func (i FailInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*Fail)
	switch name {
	case "msg":
		b.msg = schema.Value[string](value)
	}
	return nil
}

func (i FailInterpreter) SetBlock(block conflow.Block, name conflow.ID, key string, value interface{}) error {
	return nil
}
