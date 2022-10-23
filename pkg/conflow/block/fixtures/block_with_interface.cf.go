// Code generated by Conflow. DO NOT EDIT.

package fixtures

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
				annotations.Type: "configuration",
			},
			ID: "github.com/conflowio/conflow/pkg/conflow/block/fixtures.BlockWithInterface",
		},
		ParameterNames: map[string]string{"Block": "block", "Blocks": "blocks", "IDField": "id_field"},
		Properties: map[string]schema.Schema{
			"Block": &schema.Reference{
				Ref: "github.com/conflowio/conflow/pkg/conflow.Block",
			},
			"Blocks": &schema.Array{
				Items: &schema.Reference{
					Ref: "github.com/conflowio/conflow/pkg/conflow.Block",
				},
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

// BlockWithInterfaceInterpreter is the Conflow interpreter for the BlockWithInterface block
type BlockWithInterfaceInterpreter struct {
}

func (i BlockWithInterfaceInterpreter) Schema() schema.Schema {
	s, _ := schema.Get("github.com/conflowio/conflow/pkg/conflow/block/fixtures.BlockWithInterface")
	return s
}

// Create creates a new BlockWithInterface block
func (i BlockWithInterfaceInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	b := &BlockWithInterface{}
	b.IDField = id
	return b
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i BlockWithInterfaceInterpreter) ValueParamName() conflow.ID {
	return ""
}

// ParseContext returns with the parse context for the block
func (i BlockWithInterfaceInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *BlockWithInterface
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i BlockWithInterfaceInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "id_field":
		return b.(*BlockWithInterface).IDField
	default:
		panic(fmt.Errorf("unexpected parameter %q in BlockWithInterface", name))
	}
}

func (i BlockWithInterfaceInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	return nil
}

func (i BlockWithInterfaceInterpreter) SetBlock(block conflow.Block, name conflow.ID, key string, value interface{}) error {
	b := block.(*BlockWithInterface)
	switch name {
	case "block":
		b.Block = value.(conflow.Block)
	case "blocks":
		b.Blocks = append(b.Blocks, value.(conflow.Block))
	}
	return nil
}