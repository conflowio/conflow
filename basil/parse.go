// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package basil

import (
	"fmt"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
)

// ParseText parses the text input with the given parser
func ParseText(ctx *ParseContext, p parsley.Parser, input string) (parsley.Node, error) {
	f := text.NewFile("", []byte(input))
	return parseFile(ctx, p, f)
}

// ParseFile parses a file with the given parser
func ParseFile(ctx *ParseContext, p parsley.Parser, path string) (parsley.Node, error) {
	f, err := text.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s", path)
	}
	return parseFile(ctx, p, f)
}

func parseFile(ctx *ParseContext, p parsley.Parser, f *text.File) (parsley.Node, error) {
	ctx.FileSet().AddFile(f)

	parsleyCtx := parsley.NewContext(ctx.FileSet(), text.NewReader(f))
	parsleyCtx.EnableStaticCheck()
	parsleyCtx.EnableTransformation()
	parsleyCtx.RegisterKeywords(Keywords...)
	parsleyCtx.SetUserContext(ctx)

	return parsley.Parse(parsleyCtx, combinator.Sentence(p))
}

// ParseFiles parses multiple files as one
// The result node will be created using the nodeBuilder function
// The transformation and the static checking will run on the built node
func ParseFiles(
	ctx *ParseContext,
	p parsley.Parser,
	nodeBuilder func(nodes []parsley.Node) parsley.Node,
	paths []string,
) error {
	var children []parsley.Node
	for _, path := range paths {
		f, readErr := text.ReadFile(path)
		if readErr != nil {
			return fmt.Errorf("failed to read %s", path)
		}

		ctx.FileSet().AddFile(f)

		parsleyCtx := parsley.NewContext(ctx.FileSet(), text.NewReader(f))
		parsleyCtx.RegisterKeywords(Keywords...)
		parsleyCtx.SetUserContext(ctx)

		node, parseErr := parsley.Parse(parsleyCtx, combinator.Sentence(p))
		if parseErr != nil {
			return parseErr
		}

		node = node.(*ast.NonTerminalNode).Children()[0]
		children = append(children, node.(*ast.NonTerminalNode).Children()...)
	}

	node := nodeBuilder(children)

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
