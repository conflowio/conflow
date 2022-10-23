// Code generated by Conflow. DO NOT EDIT.

package test

import (
	"fmt"
	"time"

	"github.com/conflowio/conflow/pkg/conflow"
	"github.com/conflowio/conflow/pkg/conflow/annotations"
	"github.com/conflowio/conflow/pkg/schema"
)

func init() {
	schema.Register(&schema.Object{
		Metadata: schema.Metadata{
			Annotations: map[string]string{
				annotations.Type: "directive",
			},
			ID: "github.com/conflowio/conflow/pkg/test.Directive",
		},
		ParameterNames: map[string]string{"Blocks": "testblock", "FieldArray": "field_array", "FieldBool": "field_bool", "FieldCustomName": "custom_field", "FieldFloat": "field_float", "FieldInt": "field_int", "FieldMap": "field_map", "FieldString": "field_string", "FieldTimeDuration": "field_time_duration", "IDField": "id_field", "Value": "value"},
		Properties: map[string]schema.Schema{
			"Blocks": &schema.Array{
				Items: &schema.Reference{
					Nullable: true,
					Ref:      "github.com/conflowio/conflow/pkg/test.Block",
				},
			},
			"FieldArray": &schema.Array{
				Items: &schema.Any{},
			},
			"FieldBool":       &schema.Boolean{},
			"FieldCustomName": &schema.String{},
			"FieldFloat":      &schema.Number{},
			"FieldInt":        &schema.Integer{},
			"FieldMap": &schema.Map{
				AdditionalProperties: &schema.Any{},
			},
			"FieldString": &schema.String{},
			"FieldTimeDuration": &schema.String{
				Format: "duration-go",
			},
			"IDField": &schema.String{
				Metadata: schema.Metadata{
					Annotations: map[string]string{
						annotations.ID: "true",
					},
					ReadOnly: true,
				},
				Format: "conflow.ID",
			},
			"Value": &schema.Any{
				Metadata: schema.Metadata{
					Annotations: map[string]string{
						annotations.Value: "true",
					},
				},
			},
		},
	})
}

// DirectiveInterpreter is the Conflow interpreter for the Directive block
type DirectiveInterpreter struct {
}

func (i DirectiveInterpreter) Schema() schema.Schema {
	s, _ := schema.Get("github.com/conflowio/conflow/pkg/test.Directive")
	return s
}

// Create creates a new Directive block
func (i DirectiveInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	b := &Directive{}
	b.IDField = id
	return b
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
	case "field_array":
		return b.(*Directive).FieldArray
	case "field_bool":
		return b.(*Directive).FieldBool
	case "custom_field":
		return b.(*Directive).FieldCustomName
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
	case "field_array":
		b.FieldArray = value.([]interface{})
	case "field_bool":
		b.FieldBool = value.(bool)
	case "custom_field":
		b.FieldCustomName = value.(string)
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

func (i DirectiveInterpreter) SetBlock(block conflow.Block, name conflow.ID, key string, value interface{}) error {
	b := block.(*Directive)
	switch name {
	case "testblock":
		b.Blocks = append(b.Blocks, value.(*Block))
	}
	return nil
}