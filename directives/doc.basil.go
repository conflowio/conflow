// Code generated by Basil. DO NOT EDIT.
package directives

import (
	"fmt"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/variable"
)

type DocInterpreter struct{}

// Create creates a new Doc block
func (i DocInterpreter) CreateBlock(id basil.ID) basil.Block {
	return &Doc{
		id: id,
	}
}

// Params returns with the list of valid parameters
func (i DocInterpreter) Params() map[basil.ID]basil.ParameterDescriptor {
	return map[basil.ID]basil.ParameterDescriptor{
		"description": {
			Type:       "string",
			EvalStage:  basil.EvalStages["main"],
			IsRequired: false,
			IsOutput:   false,
		},
	}
}

// Blocks returns with the list of valid blocks
func (i DocInterpreter) Blocks() map[basil.ID]basil.BlockDescriptor {
	return nil
}

// HasForeignID returns true if the block ID is referencing an other block id
func (i DocInterpreter) HasForeignID() bool {
	return false
}

// HasShortFormat returns true if the block can be defined in the short block format
func (i DocInterpreter) ValueParamName() basil.ID {
	return "description"
}

// ParseContext returns with the parse context for the block
func (i DocInterpreter) ParseContext(ctx *basil.ParseContext) *basil.ParseContext {
	var nilBlock *Doc
	if b, ok := basil.Block(nilBlock).(basil.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i DocInterpreter) Param(b basil.Block, name basil.ID) interface{} {
	switch name {
	case "id":
		return b.(*Doc).id
	case "description":
		return b.(*Doc).description
	default:
		panic(fmt.Errorf("unexpected parameter %q in Doc", name))
	}
}

func (i DocInterpreter) SetParam(block basil.Block, name basil.ID, value interface{}) error {
	var err error
	b := block.(*Doc)
	switch name {
	case "id":
		b.id, err = variable.IdentifierValue(value)
	case "description":
		b.description, err = variable.StringValue(value)
	}
	return err
}

func (i DocInterpreter) SetBlock(block basil.Block, name basil.ID, value interface{}) error {
	return nil
}

func (i DocInterpreter) EvalStage() basil.EvalStage {
	var nilBlock *Doc
	if b, ok := basil.Block(nilBlock).(basil.EvalStageAware); ok {
		return b.EvalStage()
	}

	return basil.EvalStageUndefined
}