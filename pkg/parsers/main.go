// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsers

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"

	"github.com/conflowio/parsley/ast"
	"github.com/conflowio/parsley/combinator"
	"github.com/conflowio/parsley/data"
	"github.com/conflowio/parsley/parsley"
	"github.com/conflowio/parsley/text"

	"github.com/conflowio/conflow/pkg/conflow"
	"github.com/conflowio/conflow/pkg/conflow/annotations"
	"github.com/conflowio/conflow/pkg/conflow/block"
	"github.com/conflowio/conflow/pkg/functions/strings"
)

// NewMain returns a parser for parsing a main block (a block body)
//
//	S     -> (PARAM|BLOCK)*
//	ID    -> /[a-z][a-z0-9]*(?:_[a-z0-9]+)*/
//	PARAM -> ID ("="|":=") P
//	VALUE -> EXPRESSION
//	      -> ARRAY
//	      -> MAP
func NewMain(id conflow.ID, interpreter conflow.BlockInterpreter) *Main {
	blockType := interpreter.Schema().GetAnnotation(annotations.Type)
	if blockType != conflow.BlockTypeMain && blockType != conflow.BlockTypeConfiguration {
		panic(fmt.Errorf(
			"%T can not be used as a main block, as it is a %s block",
			interpreter.Schema().GetID(),
			blockType,
		))
	}

	m := &Main{
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

// Main is the main block parser
// It will parse a block body (list of params and blocks)
// and will return with a block with the given id and the type "main"
type Main struct {
	id          conflow.ID
	interpreter conflow.BlockInterpreter
	p           parsley.Parser
}

// Parse will parse the input into a block
func (m *Main) Parse(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
	return m.p.Parse(ctx, leftRecCtx, pos)
}

// ParseText parses the string input as a main block
func (m *Main) ParseText(ctx *conflow.ParseContext, input string) error {
	_, err := conflow.ParseText(ctx, m.p, input)
	return err
}

// ParseFile parses the given file as a main block
func (m *Main) ParseFile(ctx *conflow.ParseContext, path string) error {
	_, err := conflow.ParseFile(ctx, m.p, path)
	return err
}

func (m *Main) ParseDir(ctx *conflow.ParseContext, dir string, recursive bool) error {
	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("path %q does not exist", dir)
		}
	}
	if !info.IsDir() {
		return fmt.Errorf("path %q is not a directory", dir)
	}

	var paths []string
	if !recursive {
		var err error
		if paths, err = filepath.Glob(path.Join(dir, "*.cf")); err != nil {
			return err
		}
	} else {
		walker := func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if strings.HasSuffix(path, ".cf") {
				paths = append(paths, path)
			}

			return nil
		}

		if err := filepath.WalkDir(dir, walker); err != nil {
			return err
		}
	}

	if len(paths) == 0 {
		return fmt.Errorf("there are no .cf files in %s", dir)
	}

	return m.ParseFiles(ctx, paths...)
}

// ParseFiles parses multiple files as one block
func (m *Main) ParseFiles(ctx *conflow.ParseContext, paths ...string) error {
	nodeBuilder := func(nodes []parsley.Node) parsley.Node {
		return ast.NewNonTerminalNode("BLOCK_BODY", nodes, m)
	}
	return conflow.ParseFiles(ctx, m.p, nodeBuilder, paths)
}

// Eval will panic as it should not be called on a raw block node
func (m *Main) Eval(userCtx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
	panic("Eval should not be called on a raw block node")
}

// TransformNode will transform the parsley node into a conflow block node
func (m *Main) TransformNode(userCtx interface{}, node parsley.Node) (parsley.Node, parsley.Error) {
	return block.TransformMainNode(userCtx, node, m.id, m.interpreter)
}
