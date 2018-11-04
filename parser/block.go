package parser

import (
	"errors"

	"github.com/opsidian/basil/block"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/identifier"
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
//     -> ID ID? VALUE
//   ID -> /[a-z][a-z0-9]*(?:_[a-z0-9]+)*/
//   ATTR -> ID "=" P
//   VALUE -> STRING
//         -> INT
//         -> FLOAT
//         -> BOOL
//         -> TIME_DURATION
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
	blockRegistry := ctx.(basil.BlockRegistryAware).GetBlockRegistry()
	idRegistry := ctx.(basil.IDRegistryAware).GetIDRegistry()

	blockIDNodes := nodes[0].(parsley.NonTerminalNode).Children()
	typeNode := blockIDNodes[0]
	var idNode *identifier.Node
	if len(blockIDNodes) == 2 {
		idNode = blockIDNodes[1].(*identifier.Node)
	} else {
		id := idRegistry.GenerateID()
		idNode = identifier.NewNode(id, true, typeNode.ReaderPos(), typeNode.ReaderPos())
	}

	var paramNodes map[string]parsley.Node
	var blockNodes []parsley.Node
	var isShortFormat bool

	if len(nodes) > 1 {
		blockValueNode := nodes[1]

		if blockValueNode.Token() == "BLOCK_BODY" {
			blockValueChildren := blockValueNode.(parsley.NonTerminalNode).Children()

			if len(blockValueChildren) > 2 {
				blockChildren := blockValueChildren[1].(parsley.NonTerminalNode).Children()

				paramCnt := 0
				blockCnt := 0
				for _, blockChild := range blockChildren {
					if blockChild.Token() == "BLOCK" {
						blockCnt++
					} else {
						paramCnt++
					}
				}

				if paramCnt > 0 {
					paramNodes = make(map[string]parsley.Node, paramCnt)
				}
				if blockCnt > 0 {
					blockNodes = make([]parsley.Node, 0, blockCnt)
				}

				for _, blockChild := range blockChildren {
					if blockChild.Token() == "BLOCK" {
						blockNodes = append(blockNodes, blockChild)
					} else {
						children := blockChild.(parsley.NonTerminalNode).Children()
						paramName, _ := children[0].Value(ctx)
						if _, alreadyExists := paramNodes[paramName.(string)]; alreadyExists {
							return nil, parsley.NewError(children[0].Pos(), errors.New("parameter was already defined"))
						}
						paramNodes[paramName.(string)] = block.NewParamNode(children[0], children[2])
					}
				}
			}
		} else { // We have an expression as the value of the block
			isShortFormat = true
			paramNodes = map[string]parsley.Node{
				"_value": blockValueNode,
			}
		}
	}

	factory, err := blockRegistry.CreateBlockFactory(ctx, typeNode, idNode, paramNodes, blockNodes)
	if err != nil {
		return nil, err
	}

	if isShortFormat && !factory.HasShortFormat() {
		parsley.NewError(typeNode.Pos(), errors.New("this block does not support short format"))
	}

	if !factory.HasForeignID() {
		if !idNode.IsGenerated() {
			id, err := idNode.Value(ctx)
			if err != nil {
				return nil, err
			}
			if err := idRegistry.RegisterID(id.(string)); err != nil {
				return nil, parsley.NewError(idNode.Pos(), err)
			}
		}
	} else {
		if idNode.IsGenerated() {
			return nil, parsley.NewError(idNode.Pos(), errors.New("identifier must be set"))
		}
	}

	return factory, nil
}
