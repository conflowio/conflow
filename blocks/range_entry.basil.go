// Code generated by Basil. DO NOT EDIT.
package blocks

import (
	"fmt"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/variable"
)

type RangeEntryInterpreter struct{}

// Create creates a new RangeEntry block
func (i RangeEntryInterpreter) CreateBlock(id basil.ID) basil.Block {
	return &RangeEntry{
		id: id,
	}
}

// Params returns with the list of valid parameters
func (i RangeEntryInterpreter) Params() map[basil.ID]basil.ParameterDescriptor {
	return map[basil.ID]basil.ParameterDescriptor{
		"key": {
			Type:       "interface{}",
			EvalStage:  basil.EvalStages["close"],
			IsRequired: false,
			IsOutput:   true,
		},
		"value": {
			Type:       "interface{}",
			EvalStage:  basil.EvalStages["close"],
			IsRequired: false,
			IsOutput:   true,
		},
	}
}

// Blocks returns with the list of valid blocks
func (i RangeEntryInterpreter) Blocks() map[basil.ID]basil.BlockDescriptor {
	return nil
}

// HasForeignID returns true if the block ID is referencing an other block id
func (i RangeEntryInterpreter) HasForeignID() bool {
	return false
}

// HasShortFormat returns true if the block can be defined in the short block format
func (i RangeEntryInterpreter) ValueParamName() basil.ID {
	return ""
}

// ParseContext returns with the parse context for the block
func (i RangeEntryInterpreter) ParseContext(ctx *basil.ParseContext) *basil.ParseContext {
	var nilBlock *RangeEntry
	if b, ok := basil.Block(nilBlock).(basil.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i RangeEntryInterpreter) Param(b basil.Block, name basil.ID) interface{} {
	switch name {
	case "id":
		return b.(*RangeEntry).id
	case "key":
		return b.(*RangeEntry).key
	case "value":
		return b.(*RangeEntry).value
	default:
		panic(fmt.Errorf("unexpected parameter %q in RangeEntry", name))
	}
}

func (i RangeEntryInterpreter) SetParam(block basil.Block, name basil.ID, value interface{}) error {
	var err error
	b := block.(*RangeEntry)
	switch name {
	case "id":
		b.id, err = variable.IdentifierValue(value)
	}
	return err
}

func (i RangeEntryInterpreter) SetBlock(block basil.Block, name basil.ID, value interface{}) error {
	return nil
}

func (i RangeEntryInterpreter) EvalStage() basil.EvalStage {
	var nilBlock *RangeEntry
	if b, ok := basil.Block(nilBlock).(basil.EvalStageAware); ok {
		return b.EvalStage()
	}

	return basil.EvalStageUndefined
}
