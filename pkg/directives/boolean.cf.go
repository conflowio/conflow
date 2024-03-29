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
			ID: "github.com/conflowio/conflow/pkg/directives.Boolean",
		},
		FieldNames:     map[string]string{"$id": "ID", "const": "Const", "default": "Default", "deprecated": "Deprecated", "description": "Description", "enum": "Enum", "examples": "Examples", "nullable": "Nullable", "readOnly": "ReadOnly", "title": "Title", "writeOnly": "WriteOnly", "x-annotations": "Annotations"},
		ParameterNames: map[string]string{"$id": "id", "readOnly": "read_only", "writeOnly": "write_only", "x-annotations": "annotations"},
		Properties: map[string]schema.Schema{
			"$id": &schema.String{},
			"const": &schema.Boolean{
				Nullable: true,
			},
			"default": &schema.Boolean{
				Nullable: true,
			},
			"deprecated":  &schema.Boolean{},
			"description": &schema.String{},
			"enum": &schema.Array{
				Items:    &schema.Boolean{},
				MaxItems: schema.Pointer(int64(1)),
			},
			"examples": &schema.Array{
				Items: &schema.Any{},
			},
			"nullable":  &schema.Boolean{},
			"readOnly":  &schema.Boolean{},
			"title":     &schema.String{},
			"writeOnly": &schema.Boolean{},
			"x-annotations": &schema.Map{
				AdditionalProperties: &schema.String{},
			},
		},
	})
}

// NewBooleanWithDefaults creates a new Boolean instance with default values
func NewBooleanWithDefaults() *Boolean {
	b := &Boolean{}
	return b
}

// BooleanInterpreter is the Conflow interpreter for the Boolean block
type BooleanInterpreter struct {
}

func (i BooleanInterpreter) Schema() schema.Schema {
	s, _ := schema.Get("github.com/conflowio/conflow/pkg/directives.Boolean")
	return s
}

// Create creates a new Boolean block
func (i BooleanInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	b := NewBooleanWithDefaults()
	return b
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
	case "id":
		return b.(*Boolean).ID
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
	case "nullable":
		return b.(*Boolean).Nullable
	case "read_only":
		return b.(*Boolean).ReadOnly
	case "title":
		return b.(*Boolean).Title
	case "write_only":
		return b.(*Boolean).WriteOnly
	case "annotations":
		return b.(*Boolean).Annotations
	default:
		panic(fmt.Errorf("unexpected parameter %q in Boolean", name))
	}
}

func (i BooleanInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*Boolean)
	switch name {
	case "id":
		b.ID = schema.Value[string](value)
	case "const":
		b.Const = schema.PointerValue[bool](value)
	case "default":
		b.Default = schema.PointerValue[bool](value)
	case "deprecated":
		b.Deprecated = schema.Value[bool](value)
	case "description":
		b.Description = schema.Value[string](value)
	case "enum":
		b.Enum = make([]bool, len(value.([]interface{})))
		for valuek, valuev := range value.([]interface{}) {
			b.Enum[valuek] = schema.Value[bool](valuev)
		}
	case "examples":
		b.Examples = value.([]interface{})
	case "nullable":
		b.Nullable = schema.Value[bool](value)
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

func (i BooleanInterpreter) SetBlock(block conflow.Block, name conflow.ID, key string, value interface{}) error {
	return nil
}
