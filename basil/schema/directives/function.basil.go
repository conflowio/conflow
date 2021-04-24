// Code generated by Basil. DO NOT EDIT.

package directives

import (
	"fmt"
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/schema"
)

// FunctionInterpreter is the basil interpreter for the Function block
type FunctionInterpreter struct {
	s schema.Schema
}

func (i FunctionInterpreter) Schema() schema.Schema {
	if i.s == nil {
		i.s = &schema.Object{
			Metadata: schema.Metadata{
				Description: "It is the directive for marking functions as basil functions",
			},
			Properties: map[string]schema.Schema{
				"id": &schema.String{
					Metadata: schema.Metadata{
						Annotations: map[string]string{"id": "true"},
						ReadOnly:    true,
					},
					Format: "basil.ID",
				},
				"path": &schema.String{},
			},
			StructProperties: map[string]string{"path": "Path"},
		}
	}
	return i.s
}

// Create creates a new Function block
func (i FunctionInterpreter) CreateBlock(id basil.ID) basil.Block {
	return &Function{
		id: id,
	}
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i FunctionInterpreter) ValueParamName() basil.ID {
	return ""
}

// ParseContext returns with the parse context for the block
func (i FunctionInterpreter) ParseContext(ctx *basil.ParseContext) *basil.ParseContext {
	var nilBlock *Function
	if b, ok := basil.Block(nilBlock).(basil.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i FunctionInterpreter) Param(b basil.Block, name basil.ID) interface{} {
	switch name {
	case "id":
		return b.(*Function).id
	case "path":
		return b.(*Function).Path
	default:
		panic(fmt.Errorf("unexpected parameter %q in Function", name))
	}
}

func (i FunctionInterpreter) SetParam(block basil.Block, name basil.ID, value interface{}) error {
	b := block.(*Function)
	switch name {
	case "path":
		b.Path = value.(string)
	}
	return nil
}

func (i FunctionInterpreter) SetBlock(block basil.Block, name basil.ID, value interface{}) error {
	return nil
}
