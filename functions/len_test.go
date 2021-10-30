// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package functions_test

import (
	"errors"

	"github.com/opsidian/conflow/basil/schema"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	basilfunction "github.com/opsidian/conflow/basil/function"
	"github.com/opsidian/conflow/functions"
	"github.com/opsidian/conflow/parsers"
	"github.com/opsidian/conflow/test"
	"github.com/opsidian/parsley/parsley"
)

var _ = Describe("Len", func() {

	registry := basilfunction.InterpreterRegistry{
		"test": functions.LenInterpreter{},
	}

	DescribeTable("it evaluates the input correctly",
		func(input string, expected interface{}) {
			test.ExpectFunctionToEvaluate(parsers.Expression(), registry)(input, expected)
		},
		test.TableEntry(`test("")`, int64(0)),
		test.TableEntry(`test("foo")`, int64(3)),
		test.TableEntry(`test("want some 🍕?")`, int64(12)),
		test.TableEntry(`test([])`, int64(0)),
		test.TableEntry(`test([1, 2])`, int64(2)),
	)

	DescribeTable("it will have a parse error",
		func(input string, expectedErr error) {
			test.ExpectFunctionToHaveParseError(parsers.Expression(), registry)(input, expectedErr)
		},
		test.TableEntry(`test()`, errors.New("test requires exactly 1 argument, but got 0 at testfile:1:1")),
		test.TableEntry(`test(1, 2)`, errors.New("test requires exactly 1 argument, but got 2 at testfile:1:9")),
		test.TableEntry(`test(1)`, errors.New("was expecting string, array or map at testfile:1:6")),
	)

	It("should return with integer type", func() {
		test.ExpectFunctionNode(parsers.Expression(), registry)(
			`test("")`,
			func(userCtx interface{}, node parsley.Node) {
				Expect(node.Schema().(schema.Schema).Type()).To(Equal(schema.TypeInteger))
			},
		)
	})

})
