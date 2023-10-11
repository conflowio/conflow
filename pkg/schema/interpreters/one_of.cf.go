// Code generated by Conflow. DO NOT EDIT.

package interpreters

import (
	"fmt"

	"github.com/conflowio/conflow/pkg/conflow"
	"github.com/conflowio/conflow/pkg/schema"
)

// OneOfInterpreter is the Conflow interpreter for the OneOf block
type OneOfInterpreter struct {
}

func (i OneOfInterpreter) Schema() schema.Schema {
	s, _ := schema.Get("github.com/conflowio/conflow/pkg/schema.OneOf")
	return s
}

// Create creates a new OneOf block
func (i OneOfInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	b := schema.NewOneOfWithDefaults()
	return b
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i OneOfInterpreter) ValueParamName() conflow.ID {
	return ""
}

// ParseContext returns with the parse context for the block
func (i OneOfInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *schema.OneOf
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i OneOfInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "id":
		return b.(*schema.OneOf).ID
	case "const":
		return b.(*schema.OneOf).Const
	case "default":
		return b.(*schema.OneOf).Default
	case "deprecated":
		return b.(*schema.OneOf).Deprecated
	case "description":
		return b.(*schema.OneOf).Description
	case "enum":
		return b.(*schema.OneOf).Enum
	case "examples":
		return b.(*schema.OneOf).Examples
	case "nullable":
		return b.(*schema.OneOf).Nullable
	case "read_only":
		return b.(*schema.OneOf).ReadOnly
	case "title":
		return b.(*schema.OneOf).Title
	case "write_only":
		return b.(*schema.OneOf).WriteOnly
	case "annotations":
		return b.(*schema.OneOf).Annotations
	default:
		panic(fmt.Errorf("unexpected parameter %q in OneOf", name))
	}
}

func (i OneOfInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*schema.OneOf)
	switch name {
	case "id":
		b.ID = schema.Value[string](value)
	case "const":
		b.Const = value
	case "default":
		b.Default = value
	case "deprecated":
		b.Deprecated = schema.Value[bool](value)
	case "description":
		b.Description = schema.Value[string](value)
	case "enum":
		b.Enum = value.([]interface{})
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

func (i OneOfInterpreter) SetBlock(block conflow.Block, name conflow.ID, key string, value interface{}) error {
	b := block.(*schema.OneOf)
	switch name {
	case "schema":
		b.Schemas = append(b.Schemas, value.(schema.Schema))
	}
	return nil
}