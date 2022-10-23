// Code generated by Conflow. DO NOT EDIT.

package blocks

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
				annotations.Type: "task",
			},
			Description: "It will write a string to the standard output",
			ID:          "github.com/conflowio/conflow/pkg/blocks.Print",
		},
		Properties: map[string]schema.Schema{
			"id": &schema.String{
				Metadata: schema.Metadata{
					Annotations: map[string]string{
						annotations.ID: "true",
					},
					ReadOnly: true,
				},
				Format: "conflow.ID",
			},
			"value": &schema.Any{
				Metadata: schema.Metadata{
					Annotations: map[string]string{
						annotations.Value: "true",
					},
				},
			},
		},
		Required: []string{"value"},
	})
}

// PrintInterpreter is the Conflow interpreter for the Print block
type PrintInterpreter struct {
}

func (i PrintInterpreter) Schema() schema.Schema {
	s, _ := schema.Get("github.com/conflowio/conflow/pkg/blocks.Print")
	return s
}

// Create creates a new Print block
func (i PrintInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	b := &Print{}
	b.id = id
	b.stdout = blockCtx.Stdout()
	return b
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i PrintInterpreter) ValueParamName() conflow.ID {
	return "value"
}

// ParseContext returns with the parse context for the block
func (i PrintInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *Print
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i PrintInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "id":
		return b.(*Print).id
	case "value":
		return b.(*Print).value
	default:
		panic(fmt.Errorf("unexpected parameter %q in Print", name))
	}
}

func (i PrintInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*Print)
	switch name {
	case "value":
		b.value = value
	}
	return nil
}

func (i PrintInterpreter) SetBlock(block conflow.Block, name conflow.ID, key string, value interface{}) error {
	return nil
}