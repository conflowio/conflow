// Code generated by Basil. DO NOT EDIT.

package directives

import (
	"fmt"
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/schema"
)

// BlockInterpreter is the basil interpreter for the Block block
type BlockInterpreter struct {
	s schema.Schema
}

func (i BlockInterpreter) Schema() schema.Schema {
	if i.s == nil {
		i.s = &schema.Object{
			Metadata: schema.Metadata{
				Description: "It is the directive for marking structs as basil blocks",
			},
			Name: "Block",
			Properties: map[string]schema.Schema{
				"eval_stage": &schema.String{
					Enum: []string{"ignore", "init", "parse", "resolve"},
				},
				"id": &schema.String{
					Metadata: schema.Metadata{
						Annotations: map[string]string{"id": "true"},
						ReadOnly:    true,
					},
					Format: "basil.ID",
				},
				"path": &schema.String{},
			},
			PropertyNames: map[string]string{"eval_stage": "EvalStage", "path": "Path"},
		}
	}
	return i.s
}

// Create creates a new Block block
func (i BlockInterpreter) CreateBlock(id basil.ID, blockCtx *basil.BlockContext) basil.Block {
	return &Block{
		id: id,
	}
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i BlockInterpreter) ValueParamName() basil.ID {
	return ""
}

// ParseContext returns with the parse context for the block
func (i BlockInterpreter) ParseContext(ctx *basil.ParseContext) *basil.ParseContext {
	var nilBlock *Block
	if b, ok := basil.Block(nilBlock).(basil.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i BlockInterpreter) Param(b basil.Block, name basil.ID) interface{} {
	switch name {
	case "eval_stage":
		return b.(*Block).EvalStage
	case "id":
		return b.(*Block).id
	case "path":
		return b.(*Block).Path
	default:
		panic(fmt.Errorf("unexpected parameter %q in Block", name))
	}
}

func (i BlockInterpreter) SetParam(block basil.Block, name basil.ID, value interface{}) error {
	b := block.(*Block)
	switch name {
	case "eval_stage":
		b.EvalStage = value.(string)
	case "path":
		b.Path = value.(string)
	}
	return nil
}

func (i BlockInterpreter) SetBlock(block basil.Block, name basil.ID, value interface{}) error {
	return nil
}
