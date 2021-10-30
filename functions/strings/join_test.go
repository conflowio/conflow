// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package strings_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	"github.com/opsidian/conflow/conflow/function"
	"github.com/opsidian/conflow/functions/strings"
	"github.com/opsidian/conflow/parsers"
	"github.com/opsidian/conflow/test"
)

var _ = Describe("Join", func() {

	registry := function.InterpreterRegistry{
		"test": strings.JoinInterpreter{},
	}

	DescribeTable("it evaluates the input correctly",
		func(input string, expected interface{}) {
			test.ExpectFunctionToEvaluate(parsers.Expression(), registry)(input, expected)
		},
		test.TableEntry(`test(["foo", "bar"], ",")`, string("foo,bar")),
	)

	DescribeTable("it will have a parse error",
		func(input string, expectedErr error) {
			test.ExpectFunctionToHaveParseError(parsers.Expression(), registry)(input, expectedErr)
		},
		test.TableEntry(`test()`, errors.New("test requires exactly 2 arguments, but got 0 at testfile:1:1")),
		test.TableEntry(`test([])`, errors.New("test requires exactly 2 arguments, but got 1 at testfile:1:8")),
		test.TableEntry(`test([], "foo", "bar")`, errors.New("test requires exactly 2 arguments, but got 3 at testfile:1:17")),
		test.TableEntry(`test("not array", "foo")`, errors.New("must be array(string) at testfile:1:6")),
		test.TableEntry(`test(["a", 2], "foo")`, errors.New("items must have the same type, but found string and integer at testfile:1:6")),
		test.TableEntry(`test([], 5)`, errors.New("must be string at testfile:1:10")),
	)

})
