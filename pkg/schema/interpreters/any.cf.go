// Code generated by Conflow. DO NOT EDIT.

package interpreters

import (
	"fmt"

	"github.com/conflowio/conflow/pkg/conflow"
	"github.com/conflowio/conflow/pkg/schema"
)

// AnyInterpreter is the Conflow interpreter for the Any block
type AnyInterpreter struct {
}

func (i AnyInterpreter) Schema() schema.Schema {
	s, _ := schema.Get("github.com/conflowio/conflow/pkg/schema.Any")
	return s
}

// Create creates a new Any block
func (i AnyInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	b := schema.NewAnyWithDefaults()
	return b
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i AnyInterpreter) ValueParamName() conflow.ID {
	return ""
}

// ParseContext returns with the parse context for the block
func (i AnyInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *schema.Any
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i AnyInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "id":
		return b.(*schema.Any).ID
	case "const":
		return b.(*schema.Any).Const
	case "default":
		return b.(*schema.Any).Default
	case "deprecated":
		return b.(*schema.Any).Deprecated
	case "description":
		return b.(*schema.Any).Description
	case "enum":
		return b.(*schema.Any).Enum
	case "examples":
		return b.(*schema.Any).Examples
	case "nullable":
		return b.(*schema.Any).Nullable
	case "read_only":
		return b.(*schema.Any).ReadOnly
	case "title":
		return b.(*schema.Any).Title
	case "write_only":
		return b.(*schema.Any).WriteOnly
	case "annotations":
		return b.(*schema.Any).Annotations
	default:
		panic(fmt.Errorf("unexpected parameter %q in Any", name))
	}
}

func (i AnyInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*schema.Any)
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

func (i AnyInterpreter) SetBlock(block conflow.Block, name conflow.ID, key string, value interface{}) error {
	return nil
}
