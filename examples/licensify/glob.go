// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"context"
	"os"
	"path/filepath"
	"regexp"

	"github.com/conflowio/conflow/conflow"
	"github.com/conflowio/conflow/conflow/block"
)

// @block
type Glob struct {
	// @id
	id conflow.ID
	// @required
	path    string
	include []string
	exclude []string
	// @generated
	file *File
	// @dependency
	blockPublisher conflow.BlockPublisher
}

func (g *Glob) ID() conflow.ID {
	return g.id
}

func (g *Glob) Run(ctx context.Context) (conflow.Result, error) {
	includes, err := g.compileRegexps(g.include)
	if err != nil {
		return nil, err
	}

	excludes, err := g.compileRegexps(g.exclude)
	if err != nil {
		return nil, err
	}

	return nil, filepath.Walk(g.path, func(path string, info os.FileInfo, err error) error {
		match := len(includes) == 0

		for _, re := range includes {
			if re.MatchString(path) {
				match = true
				break
			}
		}

		if !match {
			return nil
		}

		for _, re := range excludes {
			if re.MatchString(path) {
				return nil
			}
		}

		_, perr := g.blockPublisher.PublishBlock(&File{id: g.file.id, path: path}, nil)
		return perr
	})
}

func (g *Glob) ParseContextOverride() conflow.ParseContextOverride {
	return conflow.ParseContextOverride{
		BlockTransformerRegistry: block.InterpreterRegistry{
			"file": FileInterpreter{},
		},
	}
}

func (g *Glob) compileRegexps(exprs []string) ([]*regexp.Regexp, error) {
	var res []*regexp.Regexp
	for _, expr := range exprs {
		r, err := regexp.Compile(expr)
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	return res, nil
}

// @block
type File struct {
	// @id
	id   conflow.ID
	path string
}

func (f *File) ID() conflow.ID {
	return f.id
}
