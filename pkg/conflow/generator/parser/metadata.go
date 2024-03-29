// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parser

import (
	"context"
	goast "go/ast"
	"regexp"
	"strings"

	"github.com/conflowio/parsley/ast"
	"github.com/conflowio/parsley/combinator"
	"github.com/conflowio/parsley/parsley"
	"github.com/conflowio/parsley/text"

	"github.com/conflowio/conflow/pkg/conflow"
	"github.com/conflowio/conflow/pkg/conflow/job"
	"github.com/conflowio/conflow/pkg/parsers"
	"github.com/conflowio/conflow/pkg/schema"
	"github.com/conflowio/conflow/pkg/schema/directives"
)

type Metadata struct {
	Description string
	Directives  []schema.Directive
}

var directiveLineRegex = regexp.MustCompile(`^\s*@`)

func ParseMetadataFromComments(comments []*goast.Comment) (*Metadata, error) {
	metadata := &Metadata{}

	var directivesBuilder *strings.Builder
	var descriptionBuilder strings.Builder

	for _, comment := range comments {
		c := strings.TrimPrefix(comment.Text, "//")
		if len(c) > 0 {
			if c[0] != ' ' && c[0] != '\t' {
				continue
			}
			c = c[1:]
		}

		if directivesBuilder == nil {
			if directiveLineRegex.MatchString(c) {
				directivesBuilder = &strings.Builder{}
			} else {
				_, _ = descriptionBuilder.WriteString(c)
				_, _ = descriptionBuilder.WriteRune('\n')
				continue
			}
		}

		_, _ = directivesBuilder.WriteString(c)
		_, _ = directivesBuilder.WriteRune('\n')
	}

	metadata.Description = strings.TrimSpace(descriptionBuilder.String())

	if directivesBuilder != nil {
		idRegistry := conflow.NewIDRegistry(8, 16)
		ctx := conflow.NewParseContext(parsley.NewFileSet(), idRegistry, directives.Registry())
		node, err := conflow.ParseText(ctx, annotationParser, directivesBuilder.String())
		if err != nil {
			return nil, err
		}

		evalCtx := conflow.NewEvalContext(context.Background(), nil, nil, job.SimpleScheduler{}, nil, node)
		val, err := parsley.EvaluateNode(evalCtx, node)
		if err != nil {
			return nil, err
		}

		metadata.Directives = val.([]schema.Directive)
	}

	return metadata, nil
}

var annotationParser = combinator.Sentence(
	text.LeftTrim(
		combinator.Many(
			text.RightTrim(parsers.Directive(parsers.Expression()), text.WsSpacesForceNl),
		).Bind(ast.InterpreterFunc(func(userCtx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
			names := map[conflow.ID]bool{}
			var res []schema.Directive
			for _, n := range node.Children() {
				if names[n.(conflow.BlockNode).BlockType()] {
					return nil, parsley.NewErrorf(
						n.Pos(),
						"%s directive was defined multiple times",
						n.(conflow.BlockNode).BlockType(),
					)
				}
				names[n.(conflow.BlockNode).BlockType()] = true

				evalCtx := conflow.NewEvalContext(context.Background(), nil, nil, job.SimpleScheduler{}, nil, n)
				value, err := parsley.EvaluateNode(evalCtx, n)
				if err != nil {
					return nil, parsley.NewError(n.Pos(), err)
				}
				res = append(res, value.(schema.Directive))
			}
			return res, nil
		})),
		text.WsSpacesNl,
	),
)
