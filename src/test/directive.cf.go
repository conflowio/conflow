// Code generated by Conflow. DO NOT EDIT.

package test

import (
	"fmt"
	"github.com/conflowio/conflow/src/conflow"
	"github.com/conflowio/conflow/src/schema"
	"time"
)

// DirectiveInterpreter is the conflow interpreter for the Directive block
type DirectiveInterpreter struct {
	s schema.Schema
}

func (i DirectiveInterpreter) Schema() schema.Schema {
	if i.s == nil {
		i.s = &schema.Object{
			Metadata: schema.Metadata{
				Annotations: map[string]string{"block.conflow.io/type": "directive"},
			},
			JSONPropertyNames: map[string]string{"custom_field": "FieldCustomName", "field_array": "FieldArray", "field_bool": "FieldBool", "field_float": "FieldFloat", "field_int": "FieldInt", "field_map": "FieldMap", "field_string": "FieldString", "field_time_duration": "FieldTimeDuration", "id_field": "IDField", "testblock": "Blocks", "value": "Value"},
			Name:              "Directive",
			Parameters: map[string]schema.Schema{
				"custom_field": &schema.String{},
				"field_array": &schema.Array{
					Items: &schema.Untyped{},
				},
				"field_bool":  &schema.Boolean{},
				"field_float": &schema.Number{},
				"field_int":   &schema.Integer{},
				"field_map": &schema.Map{
					AdditionalProperties: &schema.Untyped{},
				},
				"field_string": &schema.String{},
				"field_time_duration": &schema.String{
					Format: "duration-go",
				},
				"id_field": &schema.String{
					Metadata: schema.Metadata{
						Annotations: map[string]string{"block.conflow.io/id": "true"},
						ReadOnly:    true,
					},
					Format: "conflow.ID",
				},
				"testblock": &schema.Array{
					Items: &schema.Reference{
						Nullable: true,
						Ref:      "http://conflow.schema/github.com/conflowio/conflow/src/test.Block",
					},
				},
				"value": &schema.Untyped{
					Metadata: schema.Metadata{
						Annotations: map[string]string{"block.conflow.io/value": "true"},
					},
				},
			},
		}
	}
	return i.s
}

// Create creates a new Directive block
func (i DirectiveInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	return &Directive{
		IDField: id,
	}
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i DirectiveInterpreter) ValueParamName() conflow.ID {
	return "value"
}

// ParseContext returns with the parse context for the block
func (i DirectiveInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *Directive
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i DirectiveInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "custom_field":
		return b.(*Directive).FieldCustomName
	case "field_array":
		return b.(*Directive).FieldArray
	case "field_bool":
		return b.(*Directive).FieldBool
	case "field_float":
		return b.(*Directive).FieldFloat
	case "field_int":
		return b.(*Directive).FieldInt
	case "field_map":
		return b.(*Directive).FieldMap
	case "field_string":
		return b.(*Directive).FieldString
	case "field_time_duration":
		return b.(*Directive).FieldTimeDuration
	case "id_field":
		return b.(*Directive).IDField
	case "value":
		return b.(*Directive).Value
	default:
		panic(fmt.Errorf("unexpected parameter %q in Directive", name))
	}
}

func (i DirectiveInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*Directive)
	switch name {
	case "custom_field":
		b.FieldCustomName = value.(string)
	case "field_array":
		b.FieldArray = value.([]interface{})
	case "field_bool":
		b.FieldBool = value.(bool)
	case "field_float":
		b.FieldFloat = value.(float64)
	case "field_int":
		b.FieldInt = value.(int64)
	case "field_map":
		b.FieldMap = value.(map[string]interface{})
	case "field_string":
		b.FieldString = value.(string)
	case "field_time_duration":
		b.FieldTimeDuration = value.(time.Duration)
	case "value":
		b.Value = value
	}
	return nil
}

func (i DirectiveInterpreter) SetBlock(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*Directive)
	switch name {
	case "testblock":
		b.Blocks = append(b.Blocks, value.(*Block))
	}
	return nil
}