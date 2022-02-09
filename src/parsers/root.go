// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsers

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/conflowio/parsley/ast"
	"github.com/conflowio/parsley/combinator"
	"github.com/conflowio/parsley/data"
	"github.com/conflowio/parsley/parsley"
	"github.com/conflowio/parsley/text"

	"github.com/conflowio/conflow/src/conflow"
	"github.com/conflowio/conflow/src/conflow/block"
	"github.com/conflowio/conflow/src/schema"
)

// NewRoot returns a parser for parsing the root block (a block body)
//   S     -> (PARAM|BLOCK)*
//   ID    -> /[a-z][a-z0-9]*(?:_[a-z0-9]+)*/
//   PARAM -> ID ("="|":=") P
//   VALUE -> EXPRESSION
//         -> ARRAY
//         -> MAP
func NewRoot(id conflow.ID, interpreter conflow.BlockInterpreter) *Root {
	blockType := interpreter.Schema().GetAnnotation(conflow.AnnotationType)
	if blockType != conflow.BlockTypeRoot && blockType != conflow.BlockTypeConfiguration {
		panic(fmt.Errorf(
			"%T can not be used as a root block, as it is a %s block",
			interpreter.Schema().(schema.ObjectKind).GetName(),
			blockType,
		))
	}

	m := &Root{
		id:          id,
		interpreter: interpreter,
	}

	expr := Expression()

	paramOrBlock := combinator.Choice(
		Parameter(expr, true, true),
		Block(expr),
	).Name("parameter or block definition")

	m.p = text.Trim(
		combinator.Seq(
			block.TokenBlockBody,
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

// Root is the root block parser
// It will parse a block body (list of params and blocks)
// and will return with a block with the given id and the type "root"
type Root struct {
	id          conflow.ID
	interpreter conflow.BlockInterpreter
	p           parsley.Parser
}

// Parse will parse the input into a block
func (r *Root) Parse(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
	return r.p.Parse(ctx, leftRecCtx, pos)
}

// ParseText parses the string input as a root block
func (r *Root) ParseText(ctx *conflow.ParseContext, input string) error {
	_, err := conflow.ParseText(ctx, r.p, input)
	return err
}

// ParseFile parses the given file as a root block
func (r *Root) ParseFile(ctx *conflow.ParseContext, path string) error {
	_, err := conflow.ParseFile(ctx, r.p, path)
	return err
}

func (r *Root) ParseDir(ctx *conflow.ParseContext, dir string) error {
	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("path %q does not exist", dir)
		}
	}
	if !info.IsDir() {
		return fmt.Errorf("path %q is not a directory", dir)
	}

	paths, err := filepath.Glob(path.Join(dir, "*.cf"))
	if err != nil {
		return err
	}

	return r.ParseFiles(ctx, paths...)
}

// ParseFiles parses multiple files as one block
func (r *Root) ParseFiles(ctx *conflow.ParseContext, paths ...string) error {
	nodeBuilder := func(nodes []parsley.Node) parsley.Node {
		return ast.NewNonTerminalNode("BLOCK_BODY", nodes, r)
	}
	return conflow.ParseFiles(ctx, r.p, nodeBuilder, paths)
}

// Eval will panic as it should not be called on a raw block node
func (r *Root) Eval(userCtx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
	panic("Eval should not be called on a raw block node")
}

// TransformNode will transform the parsley node into a conflow block node
func (r *Root) TransformNode(userCtx interface{}, node parsley.Node) (parsley.Node, parsley.Error) {
	return block.TransformRootNode(userCtx, node, r.id, r.interpreter)
}
