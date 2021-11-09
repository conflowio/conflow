// Code generated by Conflow. DO NOT EDIT.

package blocks

import (
	"fmt"
	"github.com/conflowio/conflow/conflow"
	"github.com/conflowio/conflow/conflow/schema"
	"io"
)

// GunzipInterpreter is the conflow interpreter for the Gunzip block
type GunzipInterpreter struct {
	s schema.Schema
}

func (i GunzipInterpreter) Schema() schema.Schema {
	if i.s == nil {
		i.s = &schema.Object{
			Name: "Gunzip",
			Properties: map[string]schema.Schema{
				"id": &schema.String{
					Metadata: schema.Metadata{
						Annotations: map[string]string{"block.conflow.io/id": "true"},
						ReadOnly:    true,
					},
					Format: "conflow.ID",
				},
				"in": &schema.ByteStream{},
				"out": &schema.Reference{
					Metadata: schema.Metadata{
						Annotations: map[string]string{"block.conflow.io/eval_stage": "init", "block.conflow.io/generated": "true"},
						Pointer:     true,
					},
					Ref: "http://conflow.schema/github.com/conflowio/conflow/blocks.Stream",
				},
			},
			Required: []string{"in", "out"},
		}
	}
	return i.s
}

// Create creates a new Gunzip block
func (i GunzipInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	return &Gunzip{
		id:             id,
		blockPublisher: blockCtx.BlockPublisher(),
	}
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i GunzipInterpreter) ValueParamName() conflow.ID {
	return ""
}

// ParseContext returns with the parse context for the block
func (i GunzipInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *Gunzip
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i GunzipInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "id":
		return b.(*Gunzip).id
	case "in":
		return b.(*Gunzip).in
	default:
		panic(fmt.Errorf("unexpected parameter %q in Gunzip", name))
	}
}

func (i GunzipInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*Gunzip)
	switch name {
	case "in":
		b.in = value.(io.ReadCloser)
	}
	return nil
}

func (i GunzipInterpreter) SetBlock(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*Gunzip)
	switch name {
	case "out":
		b.out = value.(*Stream)
	}
	return nil
}
