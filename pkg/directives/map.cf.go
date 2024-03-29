// Code generated by Conflow. DO NOT EDIT.

package directives

import (
	"fmt"

	"github.com/conflowio/conflow/pkg/conflow"
	"github.com/conflowio/conflow/pkg/conflow/annotations"
	"github.com/conflowio/conflow/pkg/schema"
)

func init() {
	schema.Register(&schema.Object{
		Metadata: schema.Metadata{
			Annotations: map[string]string{
				annotations.EvalStage: "parse",
				annotations.Type:      "directive",
			},
			ID: "github.com/conflowio/conflow/pkg/directives.Map",
		},
		FieldNames:     map[string]string{"$id": "ID", "additionalProperties": "AdditionalProperties", "const": "Const", "default": "Default", "deprecated": "Deprecated", "description": "Description", "enum": "Enum", "examples": "Examples", "maxProperties": "MaxProperties", "minProperties": "MinProperties", "readOnly": "ReadOnly", "title": "Title", "writeOnly": "WriteOnly", "x-annotations": "Annotations"},
		ParameterNames: map[string]string{"$id": "id", "additionalProperties": "additional_properties", "maxProperties": "max_properties", "minProperties": "min_properties", "readOnly": "read_only", "writeOnly": "write_only", "x-annotations": "annotations"},
		Properties: map[string]schema.Schema{
			"$id": &schema.String{},
			"additionalProperties": &schema.Reference{
				Ref: "github.com/conflowio/conflow/pkg/schema.Schema",
			},
			"const": &schema.Map{
				AdditionalProperties: &schema.Any{},
			},
			"default": &schema.Map{
				AdditionalProperties: &schema.Any{},
			},
			"deprecated":  &schema.Boolean{},
			"description": &schema.String{},
			"enum": &schema.Array{
				Items: &schema.Map{
					AdditionalProperties: &schema.Any{},
				},
			},
			"examples": &schema.Array{
				Items: &schema.Any{},
			},
			"maxProperties": &schema.Integer{
				Nullable: true,
			},
			"minProperties": &schema.Integer{},
			"readOnly":      &schema.Boolean{},
			"title":         &schema.String{},
			"writeOnly":     &schema.Boolean{},
			"x-annotations": &schema.Map{
				AdditionalProperties: &schema.String{},
			},
		},
	})
}

// NewMapWithDefaults creates a new Map instance with default values
func NewMapWithDefaults() *Map {
	b := &Map{}
	return b
}

// MapInterpreter is the Conflow interpreter for the Map block
type MapInterpreter struct {
}

func (i MapInterpreter) Schema() schema.Schema {
	s, _ := schema.Get("github.com/conflowio/conflow/pkg/directives.Map")
	return s
}

// Create creates a new Map block
func (i MapInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	b := NewMapWithDefaults()
	return b
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i MapInterpreter) ValueParamName() conflow.ID {
	return ""
}

// ParseContext returns with the parse context for the block
func (i MapInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *Map
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i MapInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "id":
		return b.(*Map).ID
	case "const":
		return b.(*Map).Const
	case "default":
		return b.(*Map).Default
	case "deprecated":
		return b.(*Map).Deprecated
	case "description":
		return b.(*Map).Description
	case "enum":
		return b.(*Map).Enum
	case "examples":
		return b.(*Map).Examples
	case "max_properties":
		return b.(*Map).MaxProperties
	case "min_properties":
		return b.(*Map).MinProperties
	case "read_only":
		return b.(*Map).ReadOnly
	case "title":
		return b.(*Map).Title
	case "write_only":
		return b.(*Map).WriteOnly
	case "annotations":
		return b.(*Map).Annotations
	default:
		panic(fmt.Errorf("unexpected parameter %q in Map", name))
	}
}

func (i MapInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*Map)
	switch name {
	case "id":
		b.ID = schema.Value[string](value)
	case "const":
		b.Const = value.(map[string]interface{})
	case "default":
		b.Default = value.(map[string]interface{})
	case "deprecated":
		b.Deprecated = schema.Value[bool](value)
	case "description":
		b.Description = schema.Value[string](value)
	case "enum":
		b.Enum = make([]map[string]interface{}, len(value.([]interface{})))
		for valuek, valuev := range value.([]interface{}) {
			b.Enum[valuek] = valuev.(map[string]interface{})
		}
	case "examples":
		b.Examples = value.([]interface{})
	case "max_properties":
		b.MaxProperties = schema.PointerValue[int64](value)
	case "min_properties":
		b.MinProperties = schema.Value[int64](value)
	case "read_only":
		b.ReadOnly = schema.Value[bool](value)
	case "title":
		b.Title = schema.Value[string](value)
	case "write_only":
		b.WriteOnly = schema.Value[bool](value)
	case "annotations":
		b.Annotations = make(map[string]string, len(value.(map[string]interface{})))
		for valuek, valuev := range value.(map[string]interface{}) {
			b.Annotations[valuek] = schema.Value[string](valuev)
		}
	}
	return nil
}

func (i MapInterpreter) SetBlock(block conflow.Block, name conflow.ID, key string, value interface{}) error {
	b := block.(*Map)
	switch name {
	case "additional_properties":
		b.AdditionalProperties = value.(schema.Schema)
	}
	return nil
}
