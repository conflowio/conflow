// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/opsidian/basil/basil"
)

//go:generate basil generate
type Licensify struct {
	id      basil.ID `basil:"id"`
	path    string   `basil:"required"`
	license string   `basil:"required"`
}

func (l *Licensify) ID() basil.ID {
	return l.id
}

func (l *Licensify) Main(ctx basil.BlockContext) error {
	content, err := ioutil.ReadFile(l.path)
	if err != nil {
		return err
	}

	if bytes.HasPrefix(content, []byte("// Code generated")) {
		return nil
	}

	if bytes.Compare(content[0:len(l.license)], []byte(l.license)) != 0 {
		buf := bytes.NewBufferString(l.license)
		buf.Write(content)
		if err := ioutil.WriteFile(l.path, buf.Bytes(), 0644); err != nil {
			return err
		}
		fmt.Printf("%s was updated\n", l.path)
	}

	return nil
}
