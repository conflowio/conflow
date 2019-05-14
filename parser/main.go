package parser

import (
	"github.com/opsidian/parsley/ast"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/block"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
)

// NewMain returns a parser for parsing a main block (a block body)
//   S     -> (PARAM|BLOCK)*
//   ID    -> /[a-z][a-z0-9]*(?:_[a-z0-9]+)*/
//   PARAM -> ID ("="|":=") P
//   VALUE -> EXPRESSION
//         -> ARRAY
//         -> MAP
func NewMain(id basil.ID, interpreter basil.BlockInterpreter) *Main {
	m := &Main{
		id:          id,
		interpreter: interpreter,
	}

	expr := Expression()

	paramOrBlock := combinator.Choice(
		Parameter(expr),
		Block(expr),
	).Name("parameter or block definition")

	m.p = text.Trim(
		combinator.Seq(
			"MAIN",
			func(i int) parsley.Parser {
				if i == 0 {
					return text.LeftTrim(paramOrBlock, text.WsSpacesNl)
				}
				return text.LeftTrim(paramOrBlock, text.WsSpacesForceNl)
			},
			func(int) bool {
				return true
			},
		).Bind(m),
	)

	return m
}

// Main is the main block parser
// It will parse a block body (list of params and blocks)
// and will return with a block with the given id and the type "main"
type Main struct {
	id          basil.ID
	interpreter basil.BlockInterpreter
	p           parsley.Parser
}

// Parse will parse the input into a block
func (m *Main) Parse(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
	return m.p.Parse(ctx, leftRecCtx, pos)
}

// ParseFile parses the given file as a main block
func (m *Main) ParseFile(ctx *basil.ParseContext, path string) error {
	return basil.ParseFile(ctx, m.p, path)
}

// ParseFiles parses multiple files as one block
func (m *Main) ParseFiles(ctx *basil.ParseContext, paths ...string) error {
	nodeBuilder := func(nodes []parsley.Node) parsley.Node {
		return ast.NewNonTerminalNode("BLOCK_BODY", nodes, m)
	}
	return basil.ParseFiles(ctx, m.p, nodeBuilder, paths)
}

// Eval will panic as it should not be called on a raw block node
func (m *Main) Eval(userCtx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
	panic("Eval should not be called on a raw block node")
}

// TransformNode will transform the parsley node into a basil block node
func (m *Main) TransformNode(userCtx interface{}, node parsley.Node) (parsley.Node, parsley.Error) {
	return block.TransformMainNode(userCtx, node, m.id, m.interpreter)
}
