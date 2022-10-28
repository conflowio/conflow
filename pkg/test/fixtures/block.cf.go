// Code generated by Conflow. DO NOT EDIT.

package fixtures

import (
	"fmt"
	"time"

	"github.com/conflowio/conflow/pkg/conflow"
	"github.com/conflowio/conflow/pkg/conflow/annotations"
	"github.com/conflowio/conflow/pkg/conflow/types"
	"github.com/conflowio/conflow/pkg/schema"
)

func init() {
	schema.Register(&schema.Object{
		Metadata: schema.Metadata{
			Annotations: map[string]string{
				annotations.Type: "configuration",
			},
			ID: "github.com/conflowio/conflow/pkg/test/fixtures.Block",
		},
		ParameterNames: map[string]string{"FieldArray": "field_array", "FieldBool": "field_bool", "FieldFloat": "field_float", "FieldIdentifier": "field_identifier", "FieldInteger": "field_integer", "FieldInterface": "field_interface", "FieldMap": "field_map", "FieldNumber": "field_number", "FieldString": "field_string", "FieldStringArray": "field_string_array", "FieldTime": "field_time", "FieldTimeDuration": "field_time_duration", "IDField": "id_field"},
		Properties: map[string]schema.Schema{
			"FieldArray": &schema.Array{
				Items: &schema.Any{},
			},
			"FieldBool":  &schema.Boolean{},
			"FieldFloat": &schema.Number{},
			"FieldIdentifier": &schema.String{
				Metadata: schema.Metadata{
					Annotations: map[string]string{
						annotations.EvalStage: "close",
					},
					ReadOnly: true,
				},
				Format: "conflow.ID",
			},
			"FieldInteger":   &schema.Integer{},
			"FieldInterface": &schema.Any{},
			"FieldMap": &schema.Map{
				AdditionalProperties: &schema.Any{},
			},
			"FieldNumber": &schema.Any{
				Types: []string{"integer", "number"},
			},
			"FieldString": &schema.String{},
			"FieldStringArray": &schema.Array{
				Items: &schema.String{},
			},
			"FieldTime": &schema.String{
				Format: "date-time",
			},
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
		},
	})
}

// BlockInterpreter is the Conflow interpreter for the Block block
type BlockInterpreter struct {
}

func (i BlockInterpreter) Schema() schema.Schema {
	s, _ := schema.Get("github.com/conflowio/conflow/pkg/test/fixtures.Block")
	return s
}

// Create creates a new Block block
func (i BlockInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	b := &Block{}
	b.IDField = id
	return b
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i BlockInterpreter) ValueParamName() conflow.ID {
	return ""
}

// ParseContext returns with the parse context for the block
func (i BlockInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *Block
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i BlockInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "field_array":
		return b.(*Block).FieldArray
	case "field_bool":
		return b.(*Block).FieldBool
	case "field_float":
		return b.(*Block).FieldFloat
	case "field_identifier":
		return b.(*Block).FieldIdentifier
	case "field_integer":
		return b.(*Block).FieldInteger
	case "field_interface":
		return b.(*Block).FieldInterface
	case "field_map":
		return b.(*Block).FieldMap
	case "field_number":
		return b.(*Block).FieldNumber
	case "field_string":
		return b.(*Block).FieldString
	case "field_string_array":
		return b.(*Block).FieldStringArray
	case "field_time":
		return b.(*Block).FieldTime
	case "field_time_duration":
		return b.(*Block).FieldTimeDuration
	case "id_field":
		return b.(*Block).IDField
	default:
		panic(fmt.Errorf("unexpected parameter %q in Block", name))
	}
}

func (i BlockInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*Block)
	switch name {
	case "field_array":
		b.FieldArray = value.([]interface{})
	case "field_bool":
		b.FieldBool = schema.Value[bool](value)
	case "field_float":
		b.FieldFloat = schema.Value[float64](value)
	case "field_integer":
		b.FieldInteger = schema.Value[int64](value)
	case "field_interface":
		b.FieldInterface = value
	case "field_map":
		b.FieldMap = value.(map[string]interface{})
	case "field_number":
		b.FieldNumber = value
	case "field_string":
		b.FieldString = schema.Value[string](value)
	case "field_string_array":
		b.FieldStringArray = make([]string, len(value.([]interface{})))
		for valuek, valuev := range value.([]interface{}) {
			b.FieldStringArray[valuek] = schema.Value[string](valuev)
		}
	case "field_time":
		b.FieldTime = schema.Value[time.Time](value)
	case "field_time_duration":
		b.FieldTimeDuration = schema.Value[types.Duration](value)
	}
	return nil
}

func (i BlockInterpreter) SetBlock(block conflow.Block, name conflow.ID, key string, value interface{}) error {
	return nil
}
