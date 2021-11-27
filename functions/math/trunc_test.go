// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package math_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"

	"github.com/conflowio/conflow/conflow/function"
	"github.com/conflowio/conflow/functions/math"
	"github.com/conflowio/conflow/parsers"
	"github.com/conflowio/conflow/test"
)

var _ = Describe("Trunc", func() {

	registry := function.InterpreterRegistry{
		"test": math.TruncInterpreter{},
	}

	DescribeTable("it evaluates the input correctly",
		func(input string, expected interface{}) {
			test.ExpectFunctionToEvaluate(parsers.Expression(), registry)(input, expected)
		},
		test.TableEntry("test(0)", 0.0),
		test.TableEntry("test(1)", 1.0),
		test.TableEntry("test(0.5)", 0.0),
		test.TableEntry("test(0.49)", 0.0),
		test.TableEntry("test(-0.5)", 0.0),
	)

	DescribeTable("it will have a parse error",
		func(input string, expectedErr error) {
			test.ExpectFunctionToHaveParseError(parsers.Expression(), registry)(input, expectedErr)
		},
		test.TableEntry(`test()`, errors.New("test requires exactly 1 argument, but got 0 at testfile:1:1")),
		test.TableEntry(`test(1, 2)`, errors.New("test requires exactly 1 argument, but got 2 at testfile:1:9")),
		test.TableEntry(`test("foo")`, errors.New("must be number at testfile:1:6")),
	)

})
