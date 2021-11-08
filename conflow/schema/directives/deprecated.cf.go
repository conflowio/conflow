// Code generated by Conflow. DO NOT EDIT.

package directives

import (
	"fmt"
	"github.com/conflowio/conflow/conflow"
	"github.com/conflowio/conflow/conflow/schema"
)

// DeprecatedInterpreter is the conflow interpreter for the Deprecated block
type DeprecatedInterpreter struct {
	s schema.Schema
}

func (i DeprecatedInterpreter) Schema() schema.Schema {
	if i.s == nil {
		i.s = &schema.Object{
			Name: "Deprecated",
			Properties: map[string]schema.Schema{
				"id": &schema.String{
					Metadata: schema.Metadata{
						Annotations: map[string]string{"block.conflow.io/id": "true"},
						ReadOnly:    true,
					},
					Format: "conflow.ID",
				},
			},
		}
	}
	return i.s
}

// Create creates a new Deprecated block
func (i DeprecatedInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	return &Deprecated{
		id: id,
	}
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i DeprecatedInterpreter) ValueParamName() conflow.ID {
	return ""
}

// ParseContext returns with the parse context for the block
func (i DeprecatedInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *Deprecated
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i DeprecatedInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "id":
		return b.(*Deprecated).id
	default:
		panic(fmt.Errorf("unexpected parameter %q in Deprecated", name))
	}
}

func (i DeprecatedInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	return nil
}

func (i DeprecatedInterpreter) SetBlock(block conflow.Block, name conflow.ID, value interface{}) error {
	return nil
}
