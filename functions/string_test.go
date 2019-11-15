// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package functions_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	basilfunction "github.com/opsidian/basil/basil/function"
	"github.com/opsidian/basil/basil/variable"
	"github.com/opsidian/basil/functions"
	"github.com/opsidian/basil/parsers"
	"github.com/opsidian/basil/test"
	"github.com/opsidian/parsley/parsley"
)

var _ = Describe("String", func() {

	registry := basilfunction.InterpreterRegistry{
		"test": functions.StringInterpreter{},
	}

	DescribeTable("it evaluates the input correctly",
		func(input string, expected interface{}) {
			test.ExpectFunctionToEvaluate(parsers.Expression(), registry)(input, expected)
		},
		test.TableEntry(`test(1)`, "1"),
		test.TableEntry(`test(1.1)`, "1.1"),
		test.TableEntry(`test(false)`, "false"),
		test.TableEntry(`test(true)`, "true"),
		test.TableEntry(`test("foo")`, "foo"),
		test.TableEntry(`test(1m30s)`, "1m30s"),
	)

	DescribeTable("it will have a parse error",
		func(input string, expectedErr error) {
			test.ExpectFunctionToHaveParseError(parsers.Expression(), registry)(input, expectedErr)
		},
		test.TableEntry(`test()`, errors.New("test expects 1 arguments at testfile:1:1")),
		test.TableEntry(`test(1, 2)`, errors.New("test expects 1 arguments at testfile:1:1")),
		test.TableEntry(`test([])`, errors.New("was expecting any basic type at testfile:1:6")),
	)

	It("should return with string type", func() {
		test.ExpectFunctionNode(parsers.Expression(), registry)(
			"test(1)",
			func(userCtx interface{}, node parsley.Node) {
				Expect(node.Type()).To(Equal(variable.TypeString))
			},
		)
	})

})