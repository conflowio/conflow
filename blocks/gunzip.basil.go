// Code generated by Basil. DO NOT EDIT.

package blocks

import (
	"fmt"
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/schema"
	"io"
)

// GunzipInterpreter is the basil interpreter for the Gunzip block
type GunzipInterpreter struct {
	s schema.Schema
}

func (i GunzipInterpreter) Schema() schema.Schema {
	if i.s == nil {
		i.s = &schema.Object{
			Properties: map[string]schema.Schema{
				"id": &schema.String{
					Metadata: schema.Metadata{
						Annotations: map[string]string{"id": "true"},
						ReadOnly:    true,
					},
					Format: "basil.ID",
				},
				"in": &schema.ByteStream{},
				"out": &schema.Reference{
					Metadata: schema.Metadata{
						Annotations: map[string]string{"eval_stage": "init", "generated": "true"},
						Pointer:     true,
					},
					Ref: "http://basil.schema/github.com/opsidian/basil/blocks/Stream",
				},
			},
			Required: []string{"in", "out"},
		}
	}
	return i.s
}

// Create creates a new Gunzip block
func (i GunzipInterpreter) CreateBlock(id basil.ID) basil.Block {
	return &Gunzip{
		id: id,
	}
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i GunzipInterpreter) ValueParamName() basil.ID {
	return ""
}

// ParseContext returns with the parse context for the block
func (i GunzipInterpreter) ParseContext(ctx *basil.ParseContext) *basil.ParseContext {
	var nilBlock *Gunzip
	if b, ok := basil.Block(nilBlock).(basil.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i GunzipInterpreter) Param(b basil.Block, name basil.ID) interface{} {
	switch name {
	case "id":
		return b.(*Gunzip).id
	case "in":
		return b.(*Gunzip).in
	default:
		panic(fmt.Errorf("unexpected parameter %q in Gunzip", name))
	}
}

func (i GunzipInterpreter) SetParam(block basil.Block, name basil.ID, value interface{}) error {
	b := block.(*Gunzip)
	switch name {
	case "in":
		b.in = value.(io.ReadCloser)
	}
	return nil
}

func (i GunzipInterpreter) SetBlock(block basil.Block, name basil.ID, value interface{}) error {
	b := block.(*Gunzip)
	switch name {
	case "out":
		b.out = value.(*Stream)
	}
	return nil
}
