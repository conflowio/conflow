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

	"github.com/conflowio/conflow/src/conflow/function"
	"github.com/conflowio/conflow/src/functions/strings"
	"github.com/conflowio/conflow/src/parsers"
	"github.com/conflowio/conflow/src/test"
)

var _ = Describe("Split", func() {

	registry := function.InterpreterRegistry{
		"test": strings.SplitInterpreter{},
	}

	DescribeTable("it evaluates the input correctly",
		func(input string, expected interface{}) {
			test.ExpectFunctionToEvaluate(parsers.Expression(), registry)(input, expected)
		},
		test.TableEntry(`test("", ",")`, []interface{}{""}),
		test.TableEntry(`test("foo", "")`, []interface{}{"f", "o", "o"}),
		test.TableEntry(`test("foo", ",")`, []interface{}{"foo"}),
		test.TableEntry(`test("a,b", ",")`, []interface{}{"a", "b"}),
	)

	DescribeTable("it will have a parse error",
		func(input string, expectedErr error) {
			test.ExpectFunctionToHaveParseError(parsers.Expression(), registry)(input, expectedErr)
		},
		test.TableEntry(`test()`, errors.New("test requires exactly 2 arguments, but got 0 at testfile:1:1")),
		test.TableEntry(`test("foo")`, errors.New("test requires exactly 2 arguments, but got 1 at testfile:1:11")),
		test.TableEntry(`test("foo", "bar", "baz")`, errors.New("test requires exactly 2 arguments, but got 3 at testfile:1:20")),
		test.TableEntry(`test(1, "foo")`, errors.New("must be string at testfile:1:6")),
		test.TableEntry(`test("foo", 1)`, errors.New("must be string at testfile:1:13")),
	)

})
