// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package function_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	basilfunction "github.com/opsidian/basil/basil/function"
	"github.com/opsidian/basil/basil/variable"
	"github.com/opsidian/basil/function"
	"github.com/opsidian/basil/parser"
	"github.com/opsidian/basil/test"
	"github.com/opsidian/parsley/parsley"
)

var _ = Describe("IsEmpty", func() {

	registry := basilfunction.InterpreterRegistry{
		"test": function.IsEmptyInterpreter{},
	}

	DescribeTable("it evaluates the input correctly",
		func(input string, expected interface{}) {
			test.ExpectFunctionToEvaluate(parser.Expression(), registry)(input, expected)
		},
		test.TableEntry(`test(nil)`, true),
		test.TableEntry(`test("")`, true),
		test.TableEntry(`test("a")`, false),
		test.TableEntry(`test(0)`, true),
		test.TableEntry(`test(1)`, false),
		test.TableEntry(`test(0.0)`, true),
		test.TableEntry(`test(0.1)`, false),
		test.TableEntry(`test([])`, true),
		test.TableEntry(`test([1])`, false),
		test.TableEntry(`test(0s)`, true),
		test.TableEntry(`test(1s)`, false),
	)

	DescribeTable("it will have a parse error",
		func(input string, expectedErr error) {
			test.ExpectFunctionToHaveParseError(parser.Expression(), registry)(input, expectedErr)
		},
		test.TableEntry(`test()`, errors.New("test expects 1 arguments at testfile:1:1")),
		test.TableEntry(`test(1, 2)`, errors.New("test expects 1 arguments at testfile:1:1")),
	)

	It("should return with bool type", func() {
		test.ExpectFunctionNode(parser.Expression(), registry)(
			`test("")`,
			func(userCtx interface{}, node parsley.Node) {
				Expect(node.Type()).To(Equal(variable.TypeBool))
			},
		)
	})

})
