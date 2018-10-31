// Code generated by ocl generate. DO NOT EDIT.
package fixtures

import (
	"fmt"

	"github.com/opsidian/ocl/ocl"
	"github.com/opsidian/ocl/util"
	"github.com/opsidian/parsley/parsley"
)

// NewBlockWithBlockFactory creates a new BlockWithBlock block factory
func NewBlockWithBlockFactory(
	typeNode parsley.Node,
	idNode parsley.Node,
	paramNodes map[string]parsley.Node,
	blockNodes []parsley.Node,
) (ocl.BlockFactory, parsley.Error) {
	return &BlockWithBlockFactory{
		typeNode:   typeNode,
		idNode:     idNode,
		paramNodes: paramNodes,
		blockNodes: blockNodes,
	}, nil
}

// BlockWithBlockFactory will create and evaluate a BlockWithBlock block
type BlockWithBlockFactory struct {
	typeNode    parsley.Node
	idNode      parsley.Node
	paramNodes  map[string]parsley.Node
	blockNodes  []parsley.Node
	shortFormat bool
}

// CreateBlock creates a new BlockWithBlock block
func (f *BlockWithBlockFactory) CreateBlock(parentCtx interface{}) (ocl.Block, interface{}, parsley.Error) {
	var err parsley.Error

	block := &BlockWithBlock{}

	if block.IDField, err = util.NodeStringValue(f.idNode, parentCtx); err != nil {
		return nil, nil, err
	}

	ctx := block.Context(parentCtx)

	if len(f.blockNodes) > 0 {
		var childBlockFactory interface{}
		for _, childBlock := range f.blockNodes {
			if childBlockFactory, err = childBlock.Value(ctx); err != nil {
				return nil, nil, err
			}
			switch b := childBlockFactory.(type) {
			case *BlockSimpleFactory:
				block.BlockFactories = append(block.BlockFactories, b)
			default:
				panic(fmt.Sprintf("block type %T is not supported in BlockWithBlock, please open a bug ticket", childBlockFactory))
			}

		}
	}

	return block, ctx, nil
}

// EvalBlock evaluates all fields belonging to the given stage on a BlockWithBlock block
func (f *BlockWithBlockFactory) EvalBlock(ctx interface{}, stage string, res ocl.Block) parsley.Error {
	var err parsley.Error

	if preInterpreter, ok := res.(ocl.BlockPreInterpreter); ok {
		if err = preInterpreter.PreEval(ctx, stage); err != nil {
			return err
		}
	}

	block, ok := res.(*BlockWithBlock)
	if !ok {
		panic("result must be a type of *BlockWithBlock")
	}

	if !f.shortFormat {
		switch stage {
		case "default":
		default:
			panic(fmt.Sprintf("unknown stage: %s", stage))
		}

		switch stage {
		case "default":
			var childBlock ocl.Block
			var childBlockCtx interface{}
			for _, childBlockFactory := range block.BlockFactories {
				if childBlock, childBlockCtx, err = childBlockFactory.CreateBlock(ctx); err != nil {
					return err
				}

				if err = childBlockFactory.EvalBlock(childBlockCtx, stage, childBlock); err != nil {
					return err
				}

				switch b := childBlock.(type) {
				case *BlockSimple:
					block.Blocks = append(block.Blocks, b)
				default:
					panic(fmt.Sprintf("block type %T is not supported in BlockFactories, please open a bug ticket", childBlock))
				}

			}
		default:
			panic(fmt.Sprintf("unknown stage: %s", stage))
		}
	}

	if postInterpreter, ok := res.(ocl.BlockPostInterpreter); ok {
		if err = postInterpreter.PostEval(ctx, stage); err != nil {
			return err
		}
	}

	return nil
}

// HasForeignID returns true if the block ID is referencing an other block id
func (f *BlockWithBlockFactory) HasForeignID() bool {
	return false
}

// HasShortFormat returns true if the block can be defined in the short block format
func (f *BlockWithBlockFactory) HasShortFormat() bool {
	return false
}
