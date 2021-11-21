// Code generated by Conflow. DO NOT EDIT.

package main

import (
	"fmt"
	"github.com/conflowio/conflow/conflow"
	"github.com/conflowio/conflow/conflow/schema"
)

// FileWalkerInterpreter is the conflow interpreter for the FileWalker block
type FileWalkerInterpreter struct {
	s schema.Schema
}

func (i FileWalkerInterpreter) Schema() schema.Schema {
	if i.s == nil {
		i.s = &schema.Object{
			Metadata: schema.Metadata{
				Annotations: map[string]string{"block.conflow.io/type": "generator"},
			},
			Name: "FileWalker",
			Parameters: map[string]schema.Schema{
				"exclude": &schema.Array{
					Items: &schema.String{},
				},
				"file": &schema.Reference{
					Metadata: schema.Metadata{
						Annotations: map[string]string{"block.conflow.io/eval_stage": "init", "block.conflow.io/generated": "true"},
					},
					Nullable: true,
					Ref:      "http://conflow.schema/github.com/conflowio/conflow/examples/licensify.File",
				},
				"id": &schema.String{
					Metadata: schema.Metadata{
						Annotations: map[string]string{"block.conflow.io/id": "true"},
						ReadOnly:    true,
					},
					Format: "conflow.ID",
				},
				"include": &schema.Array{
					Items: &schema.String{},
				},
				"path": &schema.String{},
			},
			Required: []string{"path", "file"},
		}
	}
	return i.s
}

// Create creates a new FileWalker block
func (i FileWalkerInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	return &FileWalker{
		id:             id,
		blockPublisher: blockCtx.BlockPublisher(),
	}
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i FileWalkerInterpreter) ValueParamName() conflow.ID {
	return ""
}

// ParseContext returns with the parse context for the block
func (i FileWalkerInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *FileWalker
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i FileWalkerInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "exclude":
		return b.(*FileWalker).exclude
	case "id":
		return b.(*FileWalker).id
	case "include":
		return b.(*FileWalker).include
	case "path":
		return b.(*FileWalker).path
	default:
		panic(fmt.Errorf("unexpected parameter %q in FileWalker", name))
	}
}

func (i FileWalkerInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*FileWalker)
	switch name {
	case "exclude":
		b.exclude = make([]string, len(value.([]interface{})))
		for valuek, valuev := range value.([]interface{}) {
			b.exclude[valuek] = valuev.(string)
		}
	case "include":
		b.include = make([]string, len(value.([]interface{})))
		for valuek, valuev := range value.([]interface{}) {
			b.include[valuek] = valuev.(string)
		}
	case "path":
		b.path = value.(string)
	}
	return nil
}

func (i FileWalkerInterpreter) SetBlock(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*FileWalker)
	switch name {
	case "file":
		b.file = value.(*File)
	}
	return nil
}
