// Code generated by Basil. DO NOT EDIT.

package blocks

import (
	"fmt"
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/schema"
	"io"
)

// StreamInterpreter is the basil interpreter for the Stream block
type StreamInterpreter struct {
	s schema.Schema
}

func (i StreamInterpreter) Schema() schema.Schema {
	if i.s == nil {
		i.s = &schema.Object{
			Name: "Stream",
			Properties: map[string]schema.Schema{
				"id": &schema.String{
					Metadata: schema.Metadata{
						Annotations: map[string]string{"id": "true"},
						ReadOnly:    true,
					},
					Format: "basil.ID",
				},
				"stream": &schema.ByteStream{},
			},
			PropertyNames: map[string]string{"stream": "Stream"},
		}
	}
	return i.s
}

// Create creates a new Stream block
func (i StreamInterpreter) CreateBlock(id basil.ID, blockCtx *basil.BlockContext) basil.Block {
	return &Stream{
		id: id,
	}
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i StreamInterpreter) ValueParamName() basil.ID {
	return ""
}

// ParseContext returns with the parse context for the block
func (i StreamInterpreter) ParseContext(ctx *basil.ParseContext) *basil.ParseContext {
	var nilBlock *Stream
	if b, ok := basil.Block(nilBlock).(basil.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i StreamInterpreter) Param(b basil.Block, name basil.ID) interface{} {
	switch name {
	case "id":
		return b.(*Stream).id
	case "stream":
		return b.(*Stream).Stream
	default:
		panic(fmt.Errorf("unexpected parameter %q in Stream", name))
	}
}

func (i StreamInterpreter) SetParam(block basil.Block, name basil.ID, value interface{}) error {
	b := block.(*Stream)
	switch name {
	case "stream":
		b.Stream = value.(io.ReadCloser)
	}
	return nil
}

func (i StreamInterpreter) SetBlock(block basil.Block, name basil.ID, value interface{}) error {
	return nil
}
