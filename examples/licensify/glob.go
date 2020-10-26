// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"os"
	"path/filepath"
	"regexp"

	"github.com/opsidian/basil/basil/block"

	"github.com/opsidian/basil/basil"
)

//go:generate basil generate
type Glob struct {
	id      basil.ID `basil:"id"`
	path    string   `basil:"required"`
	include []string
	exclude []string
	file    *File `basil:"generated"`
}

func (g *Glob) ID() basil.ID {
	return g.id
}

func (g *Glob) Main(ctx basil.BlockContext) error {
	includes, err := g.compileRegexps(g.include)
	if err != nil {
		return err
	}

	excludes, err := g.compileRegexps(g.exclude)
	if err != nil {
		return err
	}

	return filepath.Walk(g.path, func(path string, info os.FileInfo, err error) error {
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

		_, perr := ctx.PublishBlock(&File{id: g.file.id, path: path}, nil)
		return perr
	})
}

func (g *Glob) ParseContextOverride() basil.ParseContextOverride {
	return basil.ParseContextOverride{
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

//go:generate basil generate
type File struct {
	id   basil.ID `basil:"id"`
	path string
}

func (f *File) ID() basil.ID {
	return f.id
}
