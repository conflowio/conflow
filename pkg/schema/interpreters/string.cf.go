// Code generated by Conflow. DO NOT EDIT.

package interpreters

import (
	"fmt"
	"regexp"

	"github.com/conflowio/conflow/pkg/conflow"
	"github.com/conflowio/conflow/pkg/conflow/annotations"
	"github.com/conflowio/conflow/pkg/schema"
)

func init() {
	schema.Register(&schema.Object{
		Metadata: schema.Metadata{
			Annotations: map[string]string{
				annotations.Type: "configuration",
			},
			ID: "github.com/conflowio/conflow/pkg/schema.String",
		},
		FieldNames:     map[string]string{"$id": "ID", "const": "Const", "default": "Default", "deprecated": "Deprecated", "description": "Description", "enum": "Enum", "examples": "Examples", "format": "Format", "maxLength": "MaxLength", "minLength": "MinLength", "nullable": "Nullable", "pattern": "Pattern", "readOnly": "ReadOnly", "title": "Title", "writeOnly": "WriteOnly", "x-annotations": "Annotations"},
		ParameterNames: map[string]string{"$id": "id", "maxLength": "max_length", "minLength": "min_length", "readOnly": "read_only", "writeOnly": "write_only", "x-annotations": "annotations"},
		Properties: map[string]schema.Schema{
			"$id": &schema.String{},
			"const": &schema.String{
				Nullable: true,
			},
			"default": &schema.String{
				Nullable: true,
			},
			"deprecated":  &schema.Boolean{},
			"description": &schema.String{},
			"enum": &schema.Array{
				Items: &schema.String{},
			},
			"examples": &schema.Array{
				Items: &schema.Any{},
			},
			"format": &schema.String{},
			"maxLength": &schema.Integer{
				Nullable: true,
			},
			"minLength": &schema.Integer{},
			"nullable":  &schema.Boolean{},
			"pattern": &schema.String{
				Format:   "regex",
				Nullable: true,
			},
			"readOnly":  &schema.Boolean{},
			"title":     &schema.String{},
			"writeOnly": &schema.Boolean{},
			"x-annotations": &schema.Map{
				AdditionalProperties: &schema.String{},
			},
		},
	})
}

// StringInterpreter is the Conflow interpreter for the String block
type StringInterpreter struct {
}

func (i StringInterpreter) Schema() schema.Schema {
	s, _ := schema.Get("github.com/conflowio/conflow/pkg/schema.String")
	return s
}

// Create creates a new String block
func (i StringInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	b := &schema.String{}
	return b
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i StringInterpreter) ValueParamName() conflow.ID {
	return ""
}

// ParseContext returns with the parse context for the block
func (i StringInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *schema.String
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i StringInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "id":
		return b.(*schema.String).ID
	case "const":
		return b.(*schema.String).Const
	case "default":
		return b.(*schema.String).Default
	case "deprecated":
		return b.(*schema.String).Deprecated
	case "description":
		return b.(*schema.String).Description
	case "enum":
		return b.(*schema.String).Enum
	case "examples":
		return b.(*schema.String).Examples
	case "format":
		return b.(*schema.String).Format
	case "max_length":
		return b.(*schema.String).MaxLength
	case "min_length":
		return b.(*schema.String).MinLength
	case "nullable":
		return b.(*schema.String).Nullable
	case "pattern":
		return b.(*schema.String).Pattern
	case "read_only":
		return b.(*schema.String).ReadOnly
	case "title":
		return b.(*schema.String).Title
	case "write_only":
		return b.(*schema.String).WriteOnly
	case "annotations":
		return b.(*schema.String).Annotations
	default:
		panic(fmt.Errorf("unexpected parameter %q in String", name))
	}
}

func (i StringInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*schema.String)
	switch name {
	case "id":
		b.ID = value.(string)
	case "const":
		b.Const = schema.Pointer(value.(string))
	case "default":
		b.Default = schema.Pointer(value.(string))
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
	case "max_length":
		b.MaxLength = schema.Pointer(value.(int64))
	case "min_length":
		b.MinLength = value.(int64)
	case "nullable":
		b.Nullable = value.(bool)
	case "pattern":
		b.Pattern = schema.Pointer(value.(regexp.Regexp))
	case "read_only":
		b.ReadOnly = value.(bool)
	case "title":
		b.Title = value.(string)
	case "write_only":
		b.WriteOnly = value.(bool)
	case "annotations":
		b.Annotations = make(map[string]string, len(value.(map[string]interface{})))
		for valuek, valuev := range value.(map[string]interface{}) {
			b.Annotations[valuek] = valuev.(string)
		}
	}
	return nil
}

func (i StringInterpreter) SetBlock(block conflow.Block, name conflow.ID, key string, value interface{}) error {
	return nil
}