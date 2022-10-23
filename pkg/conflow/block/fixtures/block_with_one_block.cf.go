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
			ID: "github.com/conflowio/conflow/pkg/conflow/block/fixtures.BlockWithOneBlock",
		},
		ParameterNames: map[string]string{"BlockSimple": "block_simple", "IDField": "id_field"},
		Properties: map[string]schema.Schema{
			"BlockSimple": &schema.Reference{
				Nullable: true,
				Ref:      "github.com/conflowio/conflow/pkg/conflow/block/fixtures.BlockSimple",
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

// BlockWithOneBlockInterpreter is the Conflow interpreter for the BlockWithOneBlock block
type BlockWithOneBlockInterpreter struct {
}

func (i BlockWithOneBlockInterpreter) Schema() schema.Schema {
	s, _ := schema.Get("github.com/conflowio/conflow/pkg/conflow/block/fixtures.BlockWithOneBlock")
	return s
}

// Create creates a new BlockWithOneBlock block
func (i BlockWithOneBlockInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	b := &BlockWithOneBlock{}
	b.IDField = id
	return b
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i BlockWithOneBlockInterpreter) ValueParamName() conflow.ID {
	return ""
}

// ParseContext returns with the parse context for the block
func (i BlockWithOneBlockInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *BlockWithOneBlock
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i BlockWithOneBlockInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "id_field":
		return b.(*BlockWithOneBlock).IDField
	default:
		panic(fmt.Errorf("unexpected parameter %q in BlockWithOneBlock", name))
	}
}

func (i BlockWithOneBlockInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	return nil
}

func (i BlockWithOneBlockInterpreter) SetBlock(block conflow.Block, name conflow.ID, key string, value interface{}) error {
	b := block.(*BlockWithOneBlock)
	switch name {
	case "block_simple":
		b.BlockSimple = value.(*BlockSimple)
	}
	return nil
}