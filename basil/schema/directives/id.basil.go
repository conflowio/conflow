// Code generated by Basil. DO NOT EDIT.

package directives

import (
	"fmt"
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/schema"
)

// IDInterpreter is the basil interpreter for the ID block
type IDInterpreter struct {
	s schema.Schema
}

func (i IDInterpreter) Schema() schema.Schema {
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
			},
		}
	}
	return i.s
}

// Create creates a new ID block
func (i IDInterpreter) CreateBlock(id basil.ID) basil.Block {
	return &ID{
		id: id,
	}
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i IDInterpreter) ValueParamName() basil.ID {
	return ""
}

// ParseContext returns with the parse context for the block
func (i IDInterpreter) ParseContext(ctx *basil.ParseContext) *basil.ParseContext {
	var nilBlock *ID
	if b, ok := basil.Block(nilBlock).(basil.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i IDInterpreter) Param(b basil.Block, name basil.ID) interface{} {
	switch name {
	case "id":
		return b.(*ID).id
	default:
		panic(fmt.Errorf("unexpected parameter %q in ID", name))
	}
}

func (i IDInterpreter) SetParam(block basil.Block, name basil.ID, value interface{}) error {
	return nil
}

func (i IDInterpreter) SetBlock(block basil.Block, name basil.ID, value interface{}) error {
	return nil
}
