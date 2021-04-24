// Code generated by Basil. DO NOT EDIT.

package directives

import (
	"fmt"
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/schema"
)

// StringInterpreter is the basil interpreter for the String block
type StringInterpreter struct {
	s schema.Schema
}

func (i StringInterpreter) Schema() schema.Schema {
	if i.s == nil {
		i.s = &schema.Object{
			Metadata: schema.Metadata{
				Annotations: map[string]string{"eval_stage": "parse"},
			},
			Properties: map[string]schema.Schema{
				"annotations": &schema.Map{
					AdditionalProperties: &schema.String{},
				},
				"const": &schema.String{
					Metadata: schema.Metadata{
						Pointer: true,
					},
				},
				"default": &schema.String{
					Metadata: schema.Metadata{
						Pointer: true,
					},
				},
				"deprecated":  &schema.Boolean{},
				"description": &schema.String{},
				"enum": &schema.Array{
					Items: &schema.String{},
				},
				"examples": &schema.Array{
					Items: &schema.Untyped{},
				},
				"format":     &schema.String{},
				"pointer":    &schema.Boolean{},
				"read_only":  &schema.Boolean{},
				"title":      &schema.String{},
				"write_only": &schema.Boolean{},
			},
			StructProperties: map[string]string{"annotations": "Annotations", "const": "Const", "default": "Default", "deprecated": "Deprecated", "description": "Description", "enum": "Enum", "examples": "Examples", "format": "Format", "pointer": "Pointer", "read_only": "ReadOnly", "title": "Title", "write_only": "WriteOnly"},
		}
	}
	return i.s
}

// Create creates a new String block
func (i StringInterpreter) CreateBlock(id basil.ID) basil.Block {
	return &String{}
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i StringInterpreter) ValueParamName() basil.ID {
	return ""
}

// ParseContext returns with the parse context for the block
func (i StringInterpreter) ParseContext(ctx *basil.ParseContext) *basil.ParseContext {
	var nilBlock *String
	if b, ok := basil.Block(nilBlock).(basil.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i StringInterpreter) Param(b basil.Block, name basil.ID) interface{} {
	switch name {
	case "annotations":
		return b.(*String).Annotations
	case "const":
		return b.(*String).Const
	case "default":
		return b.(*String).Default
	case "deprecated":
		return b.(*String).Deprecated
	case "description":
		return b.(*String).Description
	case "enum":
		return b.(*String).Enum
	case "examples":
		return b.(*String).Examples
	case "format":
		return b.(*String).Format
	case "pointer":
		return b.(*String).Pointer
	case "read_only":
		return b.(*String).ReadOnly
	case "title":
		return b.(*String).Title
	case "write_only":
		return b.(*String).WriteOnly
	default:
		panic(fmt.Errorf("unexpected parameter %q in String", name))
	}
}

func (i StringInterpreter) SetParam(block basil.Block, name basil.ID, value interface{}) error {
	b := block.(*String)
	switch name {
	case "annotations":
		b.Annotations = make(map[string]string, len(value.(map[string]interface{})))
		for valuek, valuev := range value.(map[string]interface{}) {
			b.Annotations[valuek] = valuev.(string)
		}
	case "const":
		b.Const = schema.StringPtr(value.(string))
	case "default":
		b.Default = schema.StringPtr(value.(string))
	case "deprecated":
		b.Deprecated = value.(bool)
	case "description":
		b.Description = value.(string)
	case "enum":
		b.Enum = make([]string, len(value.([]interface{})))
		for valuek, valuev := range value.([]interface{}) {
			b.Enum[valuek] = valuev.(string)
		}
	case "examples":
		b.Examples = value.([]interface{})
	case "format":
		b.Format = value.(string)
	case "pointer":
		b.Pointer = value.(bool)
	case "read_only":
		b.ReadOnly = value.(bool)
	case "title":
		b.Title = value.(string)
	case "write_only":
		b.WriteOnly = value.(bool)
	}
	return nil
}

func (i StringInterpreter) SetBlock(block basil.Block, name basil.ID, value interface{}) error {
	return nil
}
