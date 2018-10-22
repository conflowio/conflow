package parser

import (
	"errors"

	"github.com/opsidian/ocl/ocl"
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parser"
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
func Block() parser.Func {
	var block parser.Func
	expr := Expression()

	parameterValue := combinator.Choice(
		expr,
		Array(expr, text.WsSpacesNl),
	)
	parameter := combinator.Seq(
		ID(),
		text.LeftTrim(terminal.Rune('='), text.WsSpaces),
		text.LeftTrim(parameterValue, text.WsSpaces),
	)

	paramOrBlock := combinator.Choice(
		parameter,
		&block,
	).ReturnError("was expecting parameter or block definition")

	block = parser.ReturnError(
		combinator.Seq(
			blockIDs(),
			text.LeftTrim(terminal.Rune('{'), text.WsSpaces),
			combinator.Many(
				text.LeftTrim(paramOrBlock, text.WsSpacesNl),
			),
			text.LeftTrim(terminal.Rune('}'), text.WsSpacesNl),
		).Bind(ast.InterpreterFunc(evalBlock)),
		errors.New("was expecting block definition"),
	)

	return block
}

func blockIDs() *combinator.Recursive {
	lookup := func(i int) parsley.Parser {
		if i == 0 {
			return ID()
		}
		return text.LeftTrim(ID(), text.WsSpaces)
	}
	lenCheck := func(len int) bool {
		return len == 1 || len == 2
	}
	return combinator.NewRecursive("MANY", lookup, lenCheck)
}

func evalBlock(ctx interface{}, nodes []parsley.Node) (interface{}, parsley.Error) {
	blockRegistry := ctx.(ocl.BlockRegistryAware).GetBlockRegistry()

	blockIDs := nodes[0].(*ast.NonTerminalNode)
	typeNode := blockIDs.Children()[0]
	var idNode *ast.TerminalNode
	if len(blockIDs.Children()) == 2 {
		idNode = blockIDs.Children()[1].(*ast.TerminalNode)
	}

	paramCnt := 0
	blockCnt := 0

	paramsBlocks := nodes[2].(*ast.NonTerminalNode).Children()
	for _, paramBlock := range paramsBlocks {
		if len(paramBlock.(*ast.NonTerminalNode).Children()) == 3 {
			paramCnt++
		} else {
			blockCnt++
		}
	}

	paramNodes := make(map[string]parsley.Node, paramCnt)
	blockNodes := make([]parsley.Node, blockCnt)

	var blockIndex int
	for _, paramBlock := range paramsBlocks {
		children := paramBlock.(*ast.NonTerminalNode).Children()
		if len(children) == 3 {
			paramName, _ := children[0].Value(ctx)
			paramNodes[paramName.(string)] = children[2]
		} else {
			blockNodes[blockIndex] = paramBlock
			blockIndex++
		}
	}

	return blockRegistry.CreateBlock(ctx, typeNode, idNode, paramNodes, blockNodes)
}
