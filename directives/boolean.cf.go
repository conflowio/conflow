// Code generated by Conflow. DO NOT EDIT.

package directives

import (
	"fmt"
	"github.com/conflowio/conflow/conflow"
	"github.com/conflowio/conflow/conflow/schema"
)

// BooleanInterpreter is the conflow interpreter for the Boolean block
type BooleanInterpreter struct {
	s schema.Schema
}

func (i BooleanInterpreter) Schema() schema.Schema {
	if i.s == nil {
		i.s = &schema.Object{
			Metadata: schema.Metadata{
				Annotations: map[string]string{"eval_stage": "parse"},
			},
			Name: "Boolean",
			Properties: map[string]schema.Schema{
				"annotations": &schema.Map{
					AdditionalProperties: &schema.String{},
				},
				"const": &schema.Boolean{
					Metadata: schema.Metadata{
						Pointer: true,
					},
				},
				"default": &schema.Boolean{
					Metadata: schema.Metadata{
						Pointer: true,
					},
				},
				"deprecated":  &schema.Boolean{},
				"description": &schema.String{},
				"enum": &schema.Array{
					Items: &schema.Boolean{},
				},
				"examples": &schema.Array{
					Items: &schema.Untyped{},
				},
				"pointer":    &schema.Boolean{},
				"read_only":  &schema.Boolean{},
				"title":      &schema.String{},
				"write_only": &schema.Boolean{},
			},
			PropertyNames: map[string]string{"annotations": "Annotations", "const": "Const", "default": "Default", "deprecated": "Deprecated", "description": "Description", "enum": "Enum", "examples": "Examples", "pointer": "Pointer", "read_only": "ReadOnly", "title": "Title", "write_only": "WriteOnly"},
		}
	}
	return i.s
}

// Create creates a new Boolean block
func (i BooleanInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	return &Boolean{}
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i BooleanInterpreter) ValueParamName() conflow.ID {
	return ""
}

// ParseContext returns with the parse context for the block
func (i BooleanInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *Boolean
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i BooleanInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "annotations":
		return b.(*Boolean).Annotations
	case "const":
		return b.(*Boolean).Const
	case "default":
		return b.(*Boolean).Default
	case "deprecated":
		return b.(*Boolean).Deprecated
	case "description":
		return b.(*Boolean).Description
	case "enum":
		return b.(*Boolean).Enum
	case "examples":
		return b.(*Boolean).Examples
	case "pointer":
		return b.(*Boolean).Pointer
	case "read_only":
		return b.(*Boolean).ReadOnly
	case "title":
		return b.(*Boolean).Title
	case "write_only":
		return b.(*Boolean).WriteOnly
	default:
		panic(fmt.Errorf("unexpected parameter %q in Boolean", name))
	}
}

func (i BooleanInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*Boolean)
	switch name {
	case "annotations":
		b.Annotations = make(map[string]string, len(value.(map[string]interface{})))
		for valuek, valuev := range value.(map[string]interface{}) {
			b.Annotations[valuek] = valuev.(string)
		}
	case "const":
		b.Const = schema.BooleanPtr(value.(bool))
	case "default":
		b.Default = schema.BooleanPtr(value.(bool))
	case "deprecated":
		b.Deprecated = value.(bool)
	case "description":
		b.Description = value.(string)
	case "enum":
		b.Enum = make([]bool, len(value.([]interface{})))
		for valuek, valuev := range value.([]interface{}) {
			b.Enum[valuek] = valuev.(bool)
		}
	case "examples":
		b.Examples = value.([]interface{})
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

func (i BooleanInterpreter) SetBlock(block conflow.Block, name conflow.ID, value interface{}) error {
	return nil
}