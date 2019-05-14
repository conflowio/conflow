package parser

import (
	"fmt"

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
			"MAIN_BODY",
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

type Main struct {
	id          basil.ID
	interpreter basil.BlockInterpreter
	p           parsley.Parser
}

func (m *Main) Parse(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
	return m.p.Parse(ctx, leftRecCtx, pos)
}

func (m *Main) ParseFile(ctx *basil.ParseContext, path string) error {
	f, err := text.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read %s", path)
	}

	ctx.FileSet().AddFile(f)

	parsleyCtx := parsley.NewContext(ctx.FileSet(), text.NewReader(f))
	parsleyCtx.EnableStaticCheck()
	parsleyCtx.EnableTransformation()
	parsleyCtx.RegisterKeywords(basil.Keywords...)
	parsleyCtx.SetUserContext(ctx)

	if _, err := parsley.Parse(parsleyCtx, m.p); err != nil {
		return err
	}

	return nil
}

func (m *Main) ParseFiles(ctx *basil.ParseContext, paths ...string) error {
	var children []parsley.Node
	for _, path := range paths {
		f, readErr := text.ReadFile(path)
		if readErr != nil {
			return fmt.Errorf("failed to read %s", path)
		}

		ctx.FileSet().AddFile(f)

		parsleyCtx := parsley.NewContext(ctx.FileSet(), text.NewReader(f))
		parsleyCtx.RegisterKeywords(basil.Keywords...)
		parsleyCtx.SetUserContext(ctx)

		node, parseErr := parsley.Parse(parsleyCtx, m.p)
		if parseErr != nil {
			return parseErr
		}

		children = append(children, node.(*ast.NonTerminalNode).Children()...)
	}

	var node parsley.Node = ast.NewNonTerminalNode("BLOCK_BODY", children, m)

	var transformErr parsley.Error
	node, transformErr = parsley.Transform(ctx, node)
	if transformErr != nil {
		return ctx.FileSet().ErrorWithPosition(transformErr)
	}

	if err := parsley.StaticCheck(ctx, node); err != nil {
		return ctx.FileSet().ErrorWithPosition(err)
	}

	return nil
}

func (m *Main) Eval(userCtx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
	panic("Eval should not be called on a raw block node")
}

func (m *Main) TransformNode(userCtx interface{}, node parsley.Node) (parsley.Node, parsley.Error) {
	return block.TransformMainNode(userCtx, node, m.id, m.interpreter)
}
