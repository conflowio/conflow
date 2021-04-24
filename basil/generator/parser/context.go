// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parser

import (
	"go/ast"
	"go/token"
)

type Context struct {
	WorkDir string
	FileSet *token.FileSet
	Parent  ast.Node
	File    *ast.File
}

func (c *Context) WithFile(f *ast.File) *Context {
	return &Context{
		WorkDir: c.WorkDir,
		FileSet: c.FileSet,
		Parent:  c.Parent,
		File:    f,
	}
}

func (c *Context) WithParent(parent ast.Node) *Context {
	return &Context{
		WorkDir: c.WorkDir,
		FileSet: c.FileSet,
		Parent:  parent,
		File:    c.File,
	}
}

func (c *Context) WithWorkdir(workDir string) *Context {
	return &Context{
		WorkDir: workDir,
		FileSet: c.FileSet,
		Parent:  c.Parent,
		File:    c.File,
	}
}
