// Code generated by Conflow. DO NOT EDIT.

package blocks

import (
	"fmt"
	"github.com/conflowio/conflow/src/conflow"
	"github.com/conflowio/conflow/src/schema"
)

func init() {
	schema.Register(&schema.Object{
		Metadata: schema.Metadata{
			Annotations: map[string]string{"block.conflow.io/type": "configuration"},
			ID:          "github.com/conflowio/conflow/src/schema/blocks.Map",
		},
		FieldNames:        map[string]string{"$id": "ID", "additionalProperties": "AdditionalProperties", "annotations": "Annotations", "const": "Const", "default": "Default", "deprecated": "Deprecated", "description": "Description", "enum": "Enum", "examples": "Examples", "maxProperties": "MaxProperties", "minProperties": "MinProperties", "readOnly": "ReadOnly", "title": "Title", "writeOnly": "WriteOnly"},
		JSONPropertyNames: map[string]string{"additional_properties": "additionalProperties", "id": "$id", "max_properties": "maxProperties", "min_properties": "minProperties", "read_only": "readOnly", "write_only": "writeOnly"},
		Name:              "Map",
		Parameters: map[string]schema.Schema{
			"additional_properties": &schema.Reference{
				Ref: "github.com/conflowio/conflow/src/schema.Schema",
			},
			"annotations": &schema.Map{
				AdditionalProperties: &schema.String{},
			},
			"const": &schema.Map{
				AdditionalProperties: &schema.Untyped{},
			},
			"default": &schema.Map{
				AdditionalProperties: &schema.Untyped{},
			},
			"deprecated":  &schema.Boolean{},
			"description": &schema.String{},
			"enum": &schema.Array{
				Items: &schema.Map{
					AdditionalProperties: &schema.Untyped{},
				},
			},
			"examples": &schema.Array{
				Items: &schema.Untyped{},
			},
			"id": &schema.String{},
			"max_properties": &schema.Integer{
				Nullable: true,
			},
			"min_properties": &schema.Integer{},
			"read_only":      &schema.Boolean{},
			"title":          &schema.String{},
			"write_only":     &schema.Boolean{},
		},
		Required: []string{"additional_properties"},
	})
}

// MapInterpreter is the Conflow interpreter for the Map block
type MapInterpreter struct {
}

func (i MapInterpreter) Schema() schema.Schema {
	s, _ := schema.Get("github.com/conflowio/conflow/src/schema/blocks.Map")
	return s
}

// Create creates a new Map block
func (i MapInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	b := &Map{}
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
	case "annotations":
		return b.(*Map).Annotations
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
	case "id":
		return b.(*Map).ID
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
	default:
		panic(fmt.Errorf("unexpected parameter %q in Map", name))
	}
}

func (i MapInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*Map)
	switch name {
	case "annotations":
		b.Annotations = make(map[string]string, len(value.(map[string]interface{})))
		for valuek, valuev := range value.(map[string]interface{}) {
			b.Annotations[valuek] = valuev.(string)
		}
	case "const":
		b.Const = value.(map[string]interface{})
	case "default":
		b.Default = value.(map[string]interface{})
	case "deprecated":
		b.Deprecated = value.(bool)
	case "description":
		b.Description = value.(string)
	case "enum":
		b.Enum = make([]map[string]interface{}, len(value.([]interface{})))
		for valuek, valuev := range value.([]interface{}) {
			b.Enum[valuek] = valuev.(map[string]interface{})
		}
	case "examples":
		b.Examples = value.([]interface{})
	case "id":
		b.ID = value.(string)
	case "max_properties":
		b.MaxProperties = schema.IntegerPtr(value.(int64))
	case "min_properties":
		b.MinProperties = value.(int64)
	case "read_only":
		b.ReadOnly = value.(bool)
	case "title":
		b.Title = value.(string)
	case "write_only":
		b.WriteOnly = value.(bool)
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