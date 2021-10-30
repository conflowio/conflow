// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package test

import (
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/text/terminal"

	"github.com/conflowio/conflow/conflow/schema"
)

// MapParser returns with a parser which will read the "MAP" string but the result will return a sample map
func MapParser(word string, value map[string]interface{}) parser.Func {
	s, err := schema.GetSchemaForValue(value)
	if err != nil {
		panic(err)
	}

	return terminal.Word(s, word, value)
}
