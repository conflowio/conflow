// Code generated by Conflow. DO NOT EDIT.

package interpreters

import (
	"fmt"

	"github.com/conflowio/conflow/pkg/conflow"
	"github.com/conflowio/conflow/pkg/conflow/types"
	"github.com/conflowio/conflow/pkg/schema"
)

// StringInterpreter is the Conflow interpreter for the String block
type StringInterpreter struct {
}

func (i StringInterpreter) Schema() schema.Schema {
	s, _ := schema.Get("github.com/conflowio/conflow/pkg/schema.String")
	return s
}

// Create creates a new String block
func (i StringInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	b := schema.NewStringWithDefaults()
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
		b.ID = schema.Value[string](value)
	case "const":
		b.Const = schema.PointerValue[string](value)
	case "default":
		b.Default = schema.PointerValue[string](value)
	case "deprecated":
		b.Deprecated = schema.Value[bool](value)
	case "description":
		b.Description = schema.Value[string](value)
	case "enum":
		b.Enum = make([]string, len(value.([]interface{})))
		for valuek, valuev := range value.([]interface{}) {
			b.Enum[valuek] = schema.Value[string](valuev)
		}
	case "examples":
		b.Examples = value.([]interface{})
	case "format":
		b.Format = schema.Value[string](value)
	case "max_length":
		b.MaxLength = schema.PointerValue[int64](value)
	case "min_length":
		b.MinLength = schema.Value[int64](value)
	case "nullable":
		b.Nullable = schema.Value[bool](value)
	case "pattern":
		b.Pattern = schema.PointerValue[types.Regexp](value)
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

func (i StringInterpreter) SetBlock(block conflow.Block, name conflow.ID, key string, value interface{}) error {
	return nil
}
