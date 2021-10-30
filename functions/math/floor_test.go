// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package math_test

import (
	"errors"

	"github.com/opsidian/conflow/conflow/schema"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/opsidian/conflow/conflow/function"
	"github.com/opsidian/conflow/functions/math"
	"github.com/opsidian/conflow/parsers"
	"github.com/opsidian/conflow/test"
	"github.com/opsidian/parsley/parsley"
)

var _ = Describe("Floor", func() {

	registry := function.InterpreterRegistry{
		"test": math.FloorInterpreter{},
	}

	DescribeTable("it evaluates the input correctly",
		func(input string, expected interface{}) {
			test.ExpectFunctionToEvaluate(parsers.Expression(), registry)(input, expected)
		},
		test.TableEntry("test(1)", int64(1)),
		test.TableEntry("test(1.0)", int64(1)),
		test.TableEntry("test(1.1)", int64(1)),
		test.TableEntry("test(1.9)", int64(1)),
		test.TableEntry("test(-1)", int64(-1)),
		test.TableEntry("test(-1.0)", int64(-1)),
		test.TableEntry("test(-1.1)", int64(-2)),
		test.TableEntry("test(-1.9)", int64(-2)),
	)

	DescribeTable("it will have a parse error",
		func(input string, expectedErr error) {
			test.ExpectFunctionToHaveParseError(parsers.Expression(), registry)(input, expectedErr)
		},
		test.TableEntry(`test()`, errors.New("test requires exactly 1 argument, but got 0 at testfile:1:1")),
		test.TableEntry(`test(1, 2)`, errors.New("test requires exactly 1 argument, but got 2 at testfile:1:9")),
		test.TableEntry(`test("nan")`, errors.New("was expecting integer or number at testfile:1:6")),
	)

	It("should return with integer type", func() {
		test.ExpectFunctionNode(parsers.Expression(), registry)(
			"test(1)",
			func(userCtx interface{}, node parsley.Node) {
				Expect(node.Schema().(schema.Schema).Type()).To(Equal(schema.TypeInteger))
			},
		)
	})

})
