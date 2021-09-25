// Code generated by Basil. DO NOT EDIT.

package directives

import (
	"fmt"
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/schema"
)

// IgnoreInterpreter is the basil interpreter for the Ignore block
type IgnoreInterpreter struct {
	s schema.Schema
}

func (i IgnoreInterpreter) Schema() schema.Schema {
	if i.s == nil {
		i.s = &schema.Object{
			Name: "Ignore",
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

// Create creates a new Ignore block
func (i IgnoreInterpreter) CreateBlock(id basil.ID, blockCtx *basil.BlockContext) basil.Block {
	return &Ignore{
		id: id,
	}
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i IgnoreInterpreter) ValueParamName() basil.ID {
	return ""
}

// ParseContext returns with the parse context for the block
func (i IgnoreInterpreter) ParseContext(ctx *basil.ParseContext) *basil.ParseContext {
	var nilBlock *Ignore
	if b, ok := basil.Block(nilBlock).(basil.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i IgnoreInterpreter) Param(b basil.Block, name basil.ID) interface{} {
	switch name {
	case "id":
		return b.(*Ignore).id
	default:
		panic(fmt.Errorf("unexpected parameter %q in Ignore", name))
	}
}

func (i IgnoreInterpreter) SetParam(block basil.Block, name basil.ID, value interface{}) error {
	return nil
}

func (i IgnoreInterpreter) SetBlock(block basil.Block, name basil.ID, value interface{}) error {
	return nil
}
