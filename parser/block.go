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
	var p parser.Func
	expr := Expression()

	parameterValue := combinator.Choice(
		Array(expr, text.WsSpacesNl),
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

	emptyBlock := combinator.SeqOf(
		blockIDs(),
		text.LeftTrim(terminal.Rune('{'), text.WsSpaces),
		text.LeftTrim(terminal.Rune('}'), text.WsSpacesNl),
	).Token("BLOCK").Bind(ast.InterpreterFunc(evalBlock))

	parsers := []parsley.Parser{
		blockIDs(),
		text.LeftTrim(terminal.Rune('{'), text.WsSpaces),
		combinator.Many(
			text.LeftTrim(paramOrBlock, text.WsSpacesForceNl),
		),
		text.LeftTrim(terminal.Rune('}'), text.WsSpacesForceNl),
	}

	lookup := func(i int) parsley.Parser {
		if i < len(parsers) {
			return parsers[i]
		}
		return nil
	}
	l := len(parsers)
	lenCheck := func(len int) bool {
		return len == 1 || len == l
	}

	p = combinator.Choice(
		emptyBlock,
		combinator.Seq("BLOCK", lookup, lenCheck).Bind(ast.InterpreterFunc(evalBlock)),
	).Name("block definition")

	return p
}

func blockIDs() *combinator.Sequence {
	id := ID()
	idtrim := text.LeftTrim(id, text.WsSpaces)

	lookup := func(i int) parsley.Parser {
		if i == 0 {
			return id
		}
		return idtrim
	}
	lenCheck := func(len int) bool {
		return len == 1 || len == 2
	}
	return combinator.Seq("MANY", lookup, lenCheck)
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

	if len(nodes) < 4 {
		return blockRegistry.CreateBlock(ctx, typeNode, idNode, nil, nil)
	}

	paramCnt := 0
	blockCnt := 0

	blockChildren := nodes[2].(*ast.NonTerminalNode).Children()
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
