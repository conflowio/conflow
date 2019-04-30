// Code generated by Basil. DO NOT EDIT.
package fixtures

import (
	"fmt"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/variable"
	"github.com/opsidian/parsley/parsley"
)

type BlockSimpleInterpreter struct{}

// Create creates a new BlockSimple block
func (i BlockSimpleInterpreter) Create(ctx *basil.EvalContext, node basil.BlockNode) basil.Block {
	return &BlockSimple{
		IDField: node.ID(),
	}
}

// Params returns with the list of valid parameters
func (i BlockSimpleInterpreter) Params() map[basil.ID]string {
	return map[basil.ID]string{
		"value": "interface{}",
	}
}

// RequiredParams returns with the list of required parameters
func (i BlockSimpleInterpreter) RequiredParams() map[basil.ID]bool {
	return nil
}

// HasForeignID returns true if the block ID is referencing an other block id
func (i BlockSimpleInterpreter) HasForeignID() bool {
	return false
}

// HasShortFormat returns true if the block can be defined in the short block format
func (i BlockSimpleInterpreter) ValueParamName() basil.ID {
	return "value"
}

// ParseContext returns with the parse context for the block
func (i BlockSimpleInterpreter) ParseContext(parentCtx *basil.ParseContext) *basil.ParseContext {
	var nilBlock *BlockSimple
	if b, ok := basil.Block(nilBlock).(basil.ParseContextAware); ok {
		return b.ParseContext(parentCtx)
	}

	return parentCtx
}

func (i BlockSimpleInterpreter) Param(block basil.Block, name basil.ID) interface{} {
	b := block.(*BlockSimple)

	switch name {
	case "id":
		return b.IDField
	case "value":
		return b.Value
	default:
		panic(fmt.Errorf("unexpected parameter %q in BlockSimple", name))
	}
}

func (i BlockSimpleInterpreter) SetParam(ctx *basil.EvalContext, block basil.Block, name basil.ID, node parsley.Node) parsley.Error {
	b := block.(*BlockSimple)

	switch name {
	case "id":
		var err parsley.Error
		b.IDField, err = variable.NodeIdentifierValue(node, ctx)
		return err
	case "value":
		var err parsley.Error
		b.Value, err = variable.NodeAnyValue(node, ctx)
		return err
	default:
		panic(fmt.Errorf("unexpected parameter or block %q in BlockSimple", name))
	}
}
