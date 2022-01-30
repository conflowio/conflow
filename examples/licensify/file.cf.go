// Code generated by Conflow. DO NOT EDIT.

package main

import (
	"fmt"
	"github.com/conflowio/conflow/src/conflow"
	"github.com/conflowio/conflow/src/conflow/schema"
)

// FileInterpreter is the conflow interpreter for the File block
type FileInterpreter struct {
	s schema.Schema
}

func (i FileInterpreter) Schema() schema.Schema {
	if i.s == nil {
		i.s = &schema.Object{
			Metadata: schema.Metadata{
				Annotations: map[string]string{"block.conflow.io/type": "configuration"},
			},
			Name: "File",
			Parameters: map[string]schema.Schema{
				"id": &schema.String{
					Metadata: schema.Metadata{
						Annotations: map[string]string{"block.conflow.io/id": "true"},
						ReadOnly:    true,
					},
					Format: "conflow.ID",
				},
				"path": &schema.String{},
			},
		}
	}
	return i.s
}

// Create creates a new File block
func (i FileInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	return &File{
		id: id,
	}
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i FileInterpreter) ValueParamName() conflow.ID {
	return ""
}

// ParseContext returns with the parse context for the block
func (i FileInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *File
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i FileInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "id":
		return b.(*File).id
	case "path":
		return b.(*File).path
	default:
		panic(fmt.Errorf("unexpected parameter %q in File", name))
	}
}

func (i FileInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*File)
	switch name {
	case "path":
		b.path = value.(string)
	}
	return nil
}

func (i FileInterpreter) SetBlock(block conflow.Block, name conflow.ID, value interface{}) error {
	return nil
}
