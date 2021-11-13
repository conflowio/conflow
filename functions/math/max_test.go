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

var _ = Describe("Min", func() {

	registry := function.InterpreterRegistry{
		"test": math.MaxInterpreter{},
	}

	DescribeTable("it evaluates the input correctly",
		func(input string, expected interface{}) {
			test.ExpectFunctionToEvaluate(parsers.Expression(), registry)(input, expected)
		},
		test.TableEntry("test(0)", 0.0),
		test.TableEntry("test(1)", 1.0),
		test.TableEntry("test(1.0)", 1.0),
		test.TableEntry("test(2, 3, 1)", 3.0),
		test.TableEntry("test(2.0, 3.0, 1.0)", 3.0),
		test.TableEntry("test(3.0, 2.999999999)", 3.0),
	)

	DescribeTable("it will have a parse error",
		func(input string, expectedErr error) {
			test.ExpectFunctionToHaveParseError(parsers.Expression(), registry)(input, expectedErr)
		},
		test.TableEntry(`test()`, errors.New("test requires at least 1 argument, but got 0 at testfile:1:1")),
		test.TableEntry(`test("foo")`, errors.New("was expecting integer or number at testfile:1:6")),
	)

})
