// Code generated by Basil. DO NOT EDIT.
package fixtures

import (
	"fmt"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/variable"
)

type BlockRequiredFieldInterpreter struct{}

// Create creates a new BlockRequiredField block
func (i BlockRequiredFieldInterpreter) CreateBlock(id basil.ID) basil.Block {
	return &BlockRequiredField{
		IDField: id,
	}
}

// Params returns with the list of valid parameters
func (i BlockRequiredFieldInterpreter) Params() map[basil.ID]basil.ParameterDescriptor {
	return map[basil.ID]basil.ParameterDescriptor{
		"required": {
			Type:       "interface{}",
			EvalStage:  basil.EvalStages["main"],
			IsRequired: true,
			IsOutput:   false,
		},
	}
}

// Blocks returns with the list of valid blocks
func (i BlockRequiredFieldInterpreter) Blocks() map[basil.ID]basil.BlockDescriptor {
	return nil
}

// HasForeignID returns true if the block ID is referencing an other block id
func (i BlockRequiredFieldInterpreter) HasForeignID() bool {
	return false
}

// HasShortFormat returns true if the block can be defined in the short block format
func (i BlockRequiredFieldInterpreter) ValueParamName() basil.ID {
	return ""
}

// ParseContext returns with the parse context for the block
func (i BlockRequiredFieldInterpreter) ParseContext(ctx *basil.ParseContext) *basil.ParseContext {
	var nilBlock *BlockRequiredField
	if b, ok := basil.Block(nilBlock).(basil.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i BlockRequiredFieldInterpreter) Param(b basil.Block, name basil.ID) interface{} {
	switch name {
	case "id":
		return b.(*BlockRequiredField).IDField
	case "required":
		return b.(*BlockRequiredField).required
	default:
		panic(fmt.Errorf("unexpected parameter %q in BlockRequiredField", name))
	}
}

func (i BlockRequiredFieldInterpreter) SetParam(block basil.Block, name basil.ID, value interface{}) error {
	var err error
	b := block.(*BlockRequiredField)
	switch name {
	case "id":
		b.IDField, err = variable.IdentifierValue(value)
	case "required":
		b.required, err = variable.AnyValue(value)
	}
	return err
}

func (i BlockRequiredFieldInterpreter) SetBlock(block basil.Block, name basil.ID, value interface{}) error {
	return nil
}

func (i BlockRequiredFieldInterpreter) EvalStage() basil.EvalStage {
	var nilBlock *BlockRequiredField
	if b, ok := basil.Block(nilBlock).(basil.EvalStageAware); ok {
		return b.EvalStage()
	}

	return basil.EvalStageUndefined
}