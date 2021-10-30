// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package functions_test

import (
	"errors"

	"github.com/conflowio/conflow/conflow/schema"

	conflowfunction "github.com/conflowio/conflow/conflow/function"
	"github.com/conflowio/conflow/functions"
	"github.com/conflowio/conflow/parsers"
	"github.com/conflowio/conflow/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/opsidian/parsley/parsley"
)

var _ = Describe("ArrayContains", func() {

	registry := conflowfunction.InterpreterRegistry{
		"test": functions.ArrayContainsInterpreter{},
	}

	DescribeTable("it evaluates the input correctly",
		func(input string, expected interface{}) {
			test.ExpectFunctionToEvaluate(parsers.Expression(), registry)(input, expected)
		},
		test.TableEntry("test([], 1)", false),
		test.TableEntry("test([1, 2], 1)", true),
		test.TableEntry("test([1, 2], 1.0)", true),
		test.TableEntry("test([1.0, 2.0], 1)", true),
		test.TableEntry("test([1, 2], 3)", false),
		test.TableEntry("test([[1, 2]], [1, 2])", true),
		test.TableEntry("test([[1, 2]], [1])", false),
		test.TableEntry("test([[1, 3]], [1])", false),
	)

	DescribeTable("it will have a parse error",
		func(input string, expectedErr error) {
			test.ExpectFunctionToHaveParseError(parsers.Expression(), registry)(input, expectedErr)
		},
		test.TableEntry(`test()`, errors.New("test requires exactly 2 arguments, but got 0 at testfile:1:1")),
		test.TableEntry(`test([], 1, 2)`, errors.New("test requires exactly 2 arguments, but got 3 at testfile:1:13")),
		test.TableEntry(`test("foo", 1)`, errors.New("must be array at testfile:1:6")),
	)

	It("should return with boolean type", func() {
		test.ExpectFunctionNode(parsers.Expression(), registry)(
			"test([], 1)",
			func(userCtx interface{}, node parsley.Node) {
				Expect(node.Schema().(schema.Schema).Type()).To(Equal(schema.TypeBoolean))
			},
		)
	})

})
