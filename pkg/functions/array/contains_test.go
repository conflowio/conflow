// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package array_test

import (
	"errors"

	"github.com/conflowio/parsley/parsley"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	conflowfunction "github.com/conflowio/conflow/pkg/conflow/function"
	"github.com/conflowio/conflow/pkg/functions/array"
	"github.com/conflowio/conflow/pkg/parsers"
	"github.com/conflowio/conflow/pkg/schema"
	"github.com/conflowio/conflow/pkg/test"
)

var _ = Describe("Contains", func() {

	registry := conflowfunction.InterpreterRegistry{
		"test": array.ContainsInterpreter{},
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
