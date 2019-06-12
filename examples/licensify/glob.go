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
	pattern string   `basil:"required"`
	file    *File    `basil:"generated"`
}

func (g *Glob) ID() basil.ID {
	return g.id
}

func (g *Glob) Main(ctx basil.BlockContext) error {
	regexp, err := regexp.Compile(g.pattern)
	if err != nil {
		return err
	}
	return filepath.Walk(g.path, func(path string, info os.FileInfo, err error) error {
		if !regexp.MatchString(path) {
			return nil
		}
		return ctx.PublishBlock(&File{id: g.file.id, path: path})
	})
}

func (g *Glob) ParseContextOverride() basil.ParseContextOverride {
	return basil.ParseContextOverride{
		BlockTransformerRegistry: block.InterpreterRegistry{
			"file": FileInterpreter{},
		},
	}
}

//go:generate basil generate
type File struct {
	id   basil.ID `basil:"id"`
	path string
}

func (f *File) ID() basil.ID {
	return f.id
}
