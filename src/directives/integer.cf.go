// Code generated by Conflow. DO NOT EDIT.

package directives

import (
	"fmt"
	"github.com/conflowio/conflow/src/conflow"
	"github.com/conflowio/conflow/src/schema"
)

func init() {
	schema.Register(&schema.Object{
		Metadata: schema.Metadata{
			Annotations: map[string]string{"block.conflow.io/eval_stage": "parse", "block.conflow.io/type": "directive"},
			ID:          "github.com/conflowio/conflow/src/directives.Integer",
		},
		FieldNames:        map[string]string{"$id": "ID", "annotations": "Annotations", "const": "Const", "default": "Default", "deprecated": "Deprecated", "description": "Description", "enum": "Enum", "examples": "Examples", "exclusiveMaximum": "ExclusiveMaximum", "exclusiveMinimum": "ExclusiveMinimum", "format": "Format", "maximum": "Maximum", "minimum": "Minimum", "multipleOf": "MultipleOf", "nullable": "Nullable", "readOnly": "ReadOnly", "title": "Title", "writeOnly": "WriteOnly"},
		JSONPropertyNames: map[string]string{"exclusive_maximum": "exclusiveMaximum", "exclusive_minimum": "exclusiveMinimum", "id": "$id", "multiple_of": "multipleOf", "read_only": "readOnly", "write_only": "writeOnly"},
		Name:              "Integer",
		Parameters: map[string]schema.Schema{
			"annotations": &schema.Map{
				AdditionalProperties: &schema.String{},
			},
			"const": &schema.Integer{
				Nullable: true,
			},
			"default": &schema.Integer{
				Nullable: true,
			},
			"deprecated":  &schema.Boolean{},
			"description": &schema.String{},
			"enum": &schema.Array{
				Items: &schema.Integer{},
			},
			"examples": &schema.Array{
				Items: &schema.Untyped{},
			},
			"exclusive_maximum": &schema.Integer{
				Nullable: true,
			},
			"exclusive_minimum": &schema.Integer{
				Nullable: true,
			},
			"format": &schema.String{
				Enum: []string{"int32", "int64"},
			},
			"id": &schema.String{},
			"maximum": &schema.Integer{
				Nullable: true,
			},
			"minimum": &schema.Integer{
				Nullable: true,
			},
			"multiple_of": &schema.Integer{
				Nullable: true,
			},
			"nullable":   &schema.Boolean{},
			"read_only":  &schema.Boolean{},
			"title":      &schema.String{},
			"write_only": &schema.Boolean{},
		},
	})
}

// IntegerInterpreter is the Conflow interpreter for the Integer block
type IntegerInterpreter struct {
}

func (i IntegerInterpreter) Schema() schema.Schema {
	s, _ := schema.Get("github.com/conflowio/conflow/src/directives.Integer")
	return s
}

// Create creates a new Integer block
func (i IntegerInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	b := &Integer{}
	return b
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i IntegerInterpreter) ValueParamName() conflow.ID {
	return ""
}

// ParseContext returns with the parse context for the block
func (i IntegerInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *Integer
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i IntegerInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "annotations":
		return b.(*Integer).Annotations
	case "const":
		return b.(*Integer).Const
	case "default":
		return b.(*Integer).Default
	case "deprecated":
		return b.(*Integer).Deprecated
	case "description":
		return b.(*Integer).Description
	case "enum":
		return b.(*Integer).Enum
	case "examples":
		return b.(*Integer).Examples
	case "exclusive_maximum":
		return b.(*Integer).ExclusiveMaximum
	case "exclusive_minimum":
		return b.(*Integer).ExclusiveMinimum
	case "format":
		return b.(*Integer).Format
	case "id":
		return b.(*Integer).ID
	case "maximum":
		return b.(*Integer).Maximum
	case "minimum":
		return b.(*Integer).Minimum
	case "multiple_of":
		return b.(*Integer).MultipleOf
	case "nullable":
		return b.(*Integer).Nullable
	case "read_only":
		return b.(*Integer).ReadOnly
	case "title":
		return b.(*Integer).Title
	case "write_only":
		return b.(*Integer).WriteOnly
	default:
		panic(fmt.Errorf("unexpected parameter %q in Integer", name))
	}
}

func (i IntegerInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*Integer)
	switch name {
	case "annotations":
		b.Annotations = make(map[string]string, len(value.(map[string]interface{})))
		for valuek, valuev := range value.(map[string]interface{}) {
			b.Annotations[valuek] = valuev.(string)
		}
	case "const":
		b.Const = schema.IntegerPtr(value.(int64))
	case "default":
		b.Default = schema.IntegerPtr(value.(int64))
	case "deprecated":
		b.Deprecated = value.(bool)
	case "description":
		b.Description = value.(string)
	case "enum":
		b.Enum = make([]int64, len(value.([]interface{})))
		for valuek, valuev := range value.([]interface{}) {
			b.Enum[valuek] = valuev.(int64)
		}
	case "examples":
		b.Examples = value.([]interface{})
	case "exclusive_maximum":
		b.ExclusiveMaximum = schema.IntegerPtr(value.(int64))
	case "exclusive_minimum":
		b.ExclusiveMinimum = schema.IntegerPtr(value.(int64))
	case "format":
		b.Format = value.(string)
	case "id":
		b.ID = value.(string)
	case "maximum":
		b.Maximum = schema.IntegerPtr(value.(int64))
	case "minimum":
		b.Minimum = schema.IntegerPtr(value.(int64))
	case "multiple_of":
		b.MultipleOf = schema.IntegerPtr(value.(int64))
	case "nullable":
		b.Nullable = value.(bool)
	case "read_only":
		b.ReadOnly = value.(bool)
	case "title":
		b.Title = value.(string)
	case "write_only":
		b.WriteOnly = value.(bool)
	}
	return nil
}

func (i IntegerInterpreter) SetBlock(block conflow.Block, name conflow.ID, key string, value interface{}) error {
	return nil
}
