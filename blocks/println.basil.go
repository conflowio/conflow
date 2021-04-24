// Code generated by Basil. DO NOT EDIT.

package blocks

import (
	"fmt"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/schema"
)

// PrintlnInterpreter is the basil interpreter for the Println block
type PrintlnInterpreter struct {
	s schema.Schema
}

func (i PrintlnInterpreter) Schema() schema.Schema {
	if i.s == nil {
		i.s = &schema.Object{
			Metadata: schema.Metadata{
				Description: "It will write a string followed by a new line to the standard output",
			},
			Properties: map[string]schema.Schema{
				"id": &schema.String{
					Metadata: schema.Metadata{
						Annotations: map[string]string{"id": "true"},
						ReadOnly:    true,
					},
					Format: "basil.ID",
				},
				"value": &schema.Untyped{
					Metadata: schema.Metadata{
						Annotations: map[string]string{"value": "true"},
					},
				},
			},
			Required: []string{"value"},
		}
	}
	return i.s
}

// Create creates a new Println block
func (i PrintlnInterpreter) CreateBlock(id basil.ID) basil.Block {
	return &Println{
		id: id,
	}
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i PrintlnInterpreter) ValueParamName() basil.ID {
	return "value"
}

// ParseContext returns with the parse context for the block
func (i PrintlnInterpreter) ParseContext(ctx *basil.ParseContext) *basil.ParseContext {
	var nilBlock *Println
	if b, ok := basil.Block(nilBlock).(basil.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i PrintlnInterpreter) Param(b basil.Block, name basil.ID) interface{} {
	switch name {
	case "id":
		return b.(*Println).id
	case "value":
		return b.(*Println).value
	default:
		panic(fmt.Errorf("unexpected parameter %q in Println", name))
	}
}

func (i PrintlnInterpreter) SetParam(block basil.Block, name basil.ID, value interface{}) error {
	b := block.(*Println)
	switch name {
	case "value":
		b.value = value
	}
	return nil
}

func (i PrintlnInterpreter) SetBlock(block basil.Block, name basil.ID, value interface{}) error {
	return nil
}
