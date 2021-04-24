// Code generated by Basil. DO NOT EDIT.

package blocks

import (
	"fmt"
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/schema"
)

// TickInterpreter is the basil interpreter for the Tick block
type TickInterpreter struct {
	s schema.Schema
}

func (i TickInterpreter) Schema() schema.Schema {
	if i.s == nil {
		i.s = &schema.Object{
			Properties: map[string]schema.Schema{
				"id": &schema.String{
					Metadata: schema.Metadata{
						Annotations: map[string]string{"id": "true"},
						ReadOnly:    true,
					},
					Format: "basil.ID",
				},
				"time": &schema.Time{
					Metadata: schema.Metadata{
						Annotations: map[string]string{"eval_stage": "close"},
						ReadOnly:    true,
					},
				},
			},
		}
	}
	return i.s
}

// Create creates a new Tick block
func (i TickInterpreter) CreateBlock(id basil.ID) basil.Block {
	return &Tick{
		id: id,
	}
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i TickInterpreter) ValueParamName() basil.ID {
	return ""
}

// ParseContext returns with the parse context for the block
func (i TickInterpreter) ParseContext(ctx *basil.ParseContext) *basil.ParseContext {
	var nilBlock *Tick
	if b, ok := basil.Block(nilBlock).(basil.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i TickInterpreter) Param(b basil.Block, name basil.ID) interface{} {
	switch name {
	case "id":
		return b.(*Tick).id
	case "time":
		return b.(*Tick).time
	default:
		panic(fmt.Errorf("unexpected parameter %q in Tick", name))
	}
}

func (i TickInterpreter) SetParam(block basil.Block, name basil.ID, value interface{}) error {
	return nil
}

func (i TickInterpreter) SetBlock(block basil.Block, name basil.ID, value interface{}) error {
	return nil
}
