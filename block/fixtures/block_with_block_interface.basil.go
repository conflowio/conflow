// Code generated by Basil. DO NOT EDIT.
package fixtures

import (
	"fmt"

	"github.com/opsidian/basil/basil"
	basilblock "github.com/opsidian/basil/block"
	"github.com/opsidian/parsley/parsley"
)

type BlockWithBlockInterfaceInterpreter struct{}

func (i BlockWithBlockInterfaceInterpreter) StaticCheck(ctx interface{}, node basil.BlockNode) (string, parsley.Error) {
	validParamNames := map[basil.ID]struct{}{
		"block_nodes": struct{}{},
	}

	for paramName, paramNode := range node.ParamNodes() {
		if _, valid := validParamNames[paramName]; !valid {
			return "", parsley.NewError(paramNode.Pos(), fmt.Errorf("%q parameter does not exist", paramName))
		}
	}

	requiredParamNames := []basil.ID{}

	for _, paramName := range requiredParamNames {
		if _, set := node.ParamNodes()[paramName]; !set {
			return "", parsley.NewError(node.Pos(), fmt.Errorf("%s parameter is required", paramName))
		}
	}

	return "*BlockWithBlockInterface", nil
}

// CreateBlock creates a new BlockWithBlockInterface block
func (i BlockWithBlockInterfaceInterpreter) Eval(parentCtx interface{}, node basil.BlockNode) (basil.Block, parsley.Error) {
	block := &BlockWithBlockInterface{
		IDField: node.ID(),
	}

	ctx := block.Context(parentCtx)

	for _, blockNode := range node.BlockNodes() {
		switch b := blockNode.(type) {
		case basil.BlockNode:
			block.BlockNodes = append(block.BlockNodes, b)
		}
	}

	if err := i.EvalBlock(ctx, node, "default", block); err != nil {
		return nil, err
	}

	return block, nil
}

// EvalBlock evaluates all fields belonging to the given stage on a BlockWithBlockInterface block
func (i BlockWithBlockInterfaceInterpreter) EvalBlock(ctx interface{}, node basil.BlockNode, stage string, res basil.Block) parsley.Error {
	var err parsley.Error

	if preInterpreter, ok := res.(basil.BlockPreInterpreter); ok {
		if err = preInterpreter.PreEval(ctx, stage); err != nil {
			return err
		}
	}

	block, ok := res.(*BlockWithBlockInterface)
	if !ok {
		panic("result must be a type of *BlockWithBlockInterface")
	}

	switch stage {
	case "default":
	default:
		panic(fmt.Sprintf("unknown stage: %s", stage))
	}

	switch stage {
	case "default":
		var childBlock interface{}
		for _, childBlockNode := range node.BlockNodes() {
			if childBlock, err = childBlockNode.Value(ctx); err != nil {
				return err
			}

			switch b := childBlock.(type) {
			case BlockInterface:
				block.Blocks = append(block.Blocks, b)
			default:
				panic(fmt.Sprintf("block type %T is not supported in BlockWithBlockInterface, please open a bug ticket", childBlock))
			}

		}
	default:
		panic(fmt.Sprintf("unknown stage: %s", stage))
	}

	if postInterpreter, ok := res.(basil.BlockPostInterpreter); ok {
		if err = postInterpreter.PostEval(ctx, stage); err != nil {
			return err
		}
	}

	return nil
}

// HasForeignID returns true if the block ID is referencing an other block id
func (i BlockWithBlockInterfaceInterpreter) HasForeignID() bool {
	return false
}

// HasShortFormat returns true if the block can be defined in the short block format
func (i BlockWithBlockInterfaceInterpreter) ValueParamName() basil.ID {
	return ""
}

func (i BlockWithBlockInterfaceInterpreter) NodeTransformer(name string) (parsley.NodeTransformer, bool) {
	var block basil.Block = &BlockWithBlockInterface{}
	if b, ok := block.(basilblock.RegistryAware); ok {
		return b.Registry().NodeTransformer(name)
	}

	return nil, false
}
