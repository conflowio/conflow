// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsers_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"

	"github.com/conflowio/conflow/src/conflow"
	"github.com/conflowio/conflow/src/parsers"
	"github.com/conflowio/conflow/src/test"
)

var _ = Describe("Name", func() {

	sep := '.'
	p := parsers.Name(sep)

	DescribeTable("it evaluates the input correctly",
		func(input string, expected interface{}) {
			test.ExpectParserToEvaluate(p)(input, expected)
		},
		test.TableEntry(`a`, conflow.ID("a")),
		test.TableEntry(`a.b`, conflow.ID("a.b")),
	)

	DescribeTable("it returns a parse error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveParseError(p)(input, expectedErr)
		},
		test.TableEntry(`a.`, errors.New("was expecting identifier at testfile:1:3")),
		test.TableEntry(`.b`, errors.New("was expecting identifier at testfile:1:1")),
		test.TableEntry(`testkeyword.a`, errors.New("testkeyword is a reserved keyword at testfile:1:1")),
		test.TableEntry(`a.testkeyword`, errors.New("testkeyword is a reserved keyword at testfile:1:3")),
		test.TableEntry(`a.testkeyword`, errors.New("testkeyword is a reserved keyword at testfile:1:3")),
	)

})
