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

// @block "generator"
type FileWalker struct {
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

func (f *FileWalker) ID() conflow.ID {
	return f.id
}

func (f *FileWalker) Run(ctx context.Context) (conflow.Result, error) {
	includes, err := f.compileRegexps(f.include)
	if err != nil {
		return nil, err
	}

	excludes, err := f.compileRegexps(f.exclude)
	if err != nil {
		return nil, err
	}

	return nil, filepath.Walk(f.path, func(path string, info os.FileInfo, err error) error {
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

		_, perr := f.blockPublisher.PublishBlock(&File{id: f.file.id, path: path}, nil)
		return perr
	})
}

func (f *FileWalker) ParseContextOverride() conflow.ParseContextOverride {
	return conflow.ParseContextOverride{
		BlockTransformerRegistry: block.InterpreterRegistry{
			"file": FileInterpreter{},
		},
	}
}

func (f *FileWalker) compileRegexps(exprs []string) ([]*regexp.Regexp, error) {
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

// @block "configuration"
type File struct {
	// @id
	id   conflow.ID
	path string
}

func (f *File) ID() conflow.ID {
	return f.id
}
