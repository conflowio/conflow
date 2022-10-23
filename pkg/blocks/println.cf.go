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
			Description: "It will write a string followed by a new line to the standard output",
			ID:          "github.com/conflowio/conflow/pkg/blocks.Println",
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

// PrintlnInterpreter is the Conflow interpreter for the Println block
type PrintlnInterpreter struct {
}

func (i PrintlnInterpreter) Schema() schema.Schema {
	s, _ := schema.Get("github.com/conflowio/conflow/pkg/blocks.Println")
	return s
}

// Create creates a new Println block
func (i PrintlnInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	b := &Println{}
	b.id = id
	b.stdout = blockCtx.Stdout()
	return b
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i PrintlnInterpreter) ValueParamName() conflow.ID {
	return "value"
}

// ParseContext returns with the parse context for the block
func (i PrintlnInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *Println
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i PrintlnInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "id":
		return b.(*Println).id
	case "value":
		return b.(*Println).value
	default:
		panic(fmt.Errorf("unexpected parameter %q in Println", name))
	}
}

func (i PrintlnInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*Println)
	switch name {
	case "value":
		b.value = value
	}
	return nil
}

func (i PrintlnInterpreter) SetBlock(block conflow.Block, name conflow.ID, key string, value interface{}) error {
	return nil
}