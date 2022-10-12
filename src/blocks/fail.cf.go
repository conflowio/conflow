// Code generated by Conflow. DO NOT EDIT.

package blocks

import (
	"fmt"
	"github.com/conflowio/conflow/src/conflow"
	"github.com/conflowio/conflow/src/schema"
)

func init() {
	schema.Register(&schema.Object{
		Metadata: schema.Metadata{
			Annotations: map[string]string{"block.conflow.io/type": "task"},
			ID:          "github.com/conflowio/conflow/src/blocks.Fail",
		},
		Name: "Fail",
		Parameters: map[string]schema.Schema{
			"id": &schema.String{
				Metadata: schema.Metadata{
					Annotations: map[string]string{"block.conflow.io/id": "true"},
					ReadOnly:    true,
				},
				Format: "conflow.ID",
			},
			"msg": &schema.String{
				Metadata: schema.Metadata{
					Annotations: map[string]string{"block.conflow.io/value": "true"},
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
	s, _ := schema.Get("github.com/conflowio/conflow/src/blocks.Fail")
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
		b.msg = value.(string)
	}
	return nil
}

func (i FailInterpreter) SetBlock(block conflow.Block, name conflow.ID, value interface{}) error {
	return nil
}
