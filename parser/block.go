package parser

import (
	"errors"

	"github.com/opsidian/ocl/ocl"
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/terminal"
)

// Block return a parser for parsing blocks
//   S -> ID ID? {
//          (ATTR|S)*
//        }
//   ID -> /[a-z][a-z0-9]*(?:_[a-z0-9]+)*/
//   ATTR -> ID "=" P
func Block() *combinator.Sequence {
	var p combinator.Sequence
	expr := Expression()

	parameterValue := combinator.Choice(
		Array(expr, text.WsSpacesNl),
		Map(expr),
		expr,
	)
	parameter := combinator.SeqOf(
		ID(),
		text.LeftTrim(terminal.Rune('='), text.WsSpaces),
		text.LeftTrim(parameterValue, text.WsSpaces),
	).Token("PARAMETER")

	paramOrBlock := combinator.Choice(
		parameter,
		&p,
	).Name("parameter or block definition")

	emptyBlockValue := combinator.SeqOf(
		terminal.Rune('{'),
		text.LeftTrim(terminal.Rune('}'), text.WsSpacesNl),
	).Token("BLOCK_BODY")

	nonEmptyBlockValue := combinator.SeqOf(
		terminal.Rune('{'),
		combinator.Many(
			text.LeftTrim(paramOrBlock, text.WsSpacesForceNl),
		),
		text.LeftTrim(terminal.Rune('}'), text.WsSpacesForceNl),
	).Token("BLOCK_BODY")

	blockValue := combinator.Choice(
		emptyBlockValue,
		nonEmptyBlockValue,
		terminal.TimeDuration(),
		terminal.Float(),
		terminal.Integer(),
		terminal.String(true),
		terminal.Bool("true", "false"),
		Array(expr, text.WsSpaces),
		Array(expr, text.WsSpacesNl),
		Map(expr),
	)

	p = *combinator.SeqTry(
		combinator.SeqTry(ID(), text.LeftTrim(ID(), text.WsSpaces)),
		text.LeftTrim(blockValue, text.WsSpaces),
	).Name("block definition").Token("BLOCK").Bind(ast.InterpreterFunc(evalBlock))

	return &p
}

func evalBlock(ctx interface{}, nodes []parsley.Node) (interface{}, parsley.Error) {
	blockRegistry := ctx.(ocl.BlockRegistryAware).GetBlockRegistry()

	blockIDNodes := nodes[0].(*ast.NonTerminalNode).Children()
	typeNode := blockIDNodes[0]
	var idNode *IDNode
	if len(blockIDNodes) == 2 {
		idNode = blockIDNodes[1].(*IDNode)
	} else {
		idNode = NewIDNode("", typeNode.ReaderPos(), typeNode.ReaderPos())
	}

	if len(nodes) == 1 {
		return blockRegistry.CreateBlock(ctx, typeNode, idNode, nil, nil)
	}

	blockValueNode := nodes[1]

	// We have an expression as the value of the block
	if blockValueNode.Token() != "BLOCK_BODY" {
		paramNodes := map[string]parsley.Node{
			"value": blockValueNode,
		}
		return blockRegistry.CreateBlock(ctx, typeNode, idNode, paramNodes, nil)
	}

	blockValueChildren := blockValueNode.(*ast.NonTerminalNode).Children()

	// We have an empty block
	if len(blockValueChildren) == 2 {
		return blockRegistry.CreateBlock(ctx, typeNode, idNode, nil, nil)
	}

	blockChildren := blockValueChildren[1].(*ast.NonTerminalNode).Children()

	paramCnt := 0
	blockCnt := 0
	for _, blockChild := range blockChildren {
		if blockChild.Token() == "BLOCK" {
			blockCnt++
		} else {
			paramCnt++
		}
	}

	paramNodes := make(map[string]parsley.Node, paramCnt)
	blockNodes := make([]parsley.Node, blockCnt)

	var blockIndex int
	for _, blockChild := range blockChildren {
		children := blockChild.(*ast.NonTerminalNode).Children()
		if blockChild.Token() == "BLOCK" {
			blockNodes[blockIndex] = blockChild
			blockIndex++
		} else {
			paramName, _ := children[0].Value(ctx)
			if _, alreadyExists := paramNodes[paramName.(string)]; alreadyExists {
				return nil, parsley.NewError(children[0].Pos(), errors.New("parameter was already set"))
			}
			paramNodes[paramName.(string)] = children[2]
		}
	}

	return blockRegistry.CreateBlock(ctx, typeNode, idNode, paramNodes, blockNodes)
}
