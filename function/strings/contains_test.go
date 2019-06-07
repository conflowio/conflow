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
	. "github.com/onsi/gomega"
	"github.com/opsidian/basil/basil/function"
	"github.com/opsidian/basil/basil/variable"
	"github.com/opsidian/basil/function/strings"
	"github.com/opsidian/basil/parser"
	"github.com/opsidian/basil/test"
	"github.com/opsidian/parsley/parsley"
)

var _ = Describe("Contains", func() {

	registry := function.InterpreterRegistry{
		"test": strings.ContainsInterpreter{},
	}

	DescribeTable("it evaluates the input correctly",
		func(input string, expected interface{}) {
			test.ExpectFunctionToEvaluate(parser.Expression(), registry)(input, expected)
		},
		test.TableEntry(`test("", "")`, true),
		test.TableEntry(`test("foo", "")`, true),
		test.TableEntry(`test("abcd", "abc")`, true),
		test.TableEntry(`test("abcd", "bcd")`, true),
		test.TableEntry(`test("abcd", "bc")`, true),
		test.TableEntry(`test("abcd", "foo")`, false),
		test.TableEntry(`test("abcd", "Bc")`, false),
	)

	DescribeTable("it will have a parse error",
		func(input string, expectedErr error) {
			test.ExpectFunctionToHaveParseError(parser.Expression(), registry)(input, expectedErr)
		},
		test.TableEntry(`test("foo")`, errors.New("test expects 2 arguments at testfile:1:1")),
		test.TableEntry(`test("a", "a", "a")`, errors.New("test expects 2 arguments at testfile:1:1")),
		test.TableEntry(`test(1, "a")`, errors.New("was expecting string at testfile:1:6")),
		test.TableEntry(`test("a", 1)`, errors.New("was expecting string at testfile:1:11")),
	)

	It("should return with boolean type", func() {
		test.ExpectFunctionNode(parser.Expression(), registry)(
			`test("", "")`,
			func(userCtx interface{}, node parsley.Node) {
				Expect(node.Type()).To(Equal(variable.TypeBool))
			},
		)
	})

})
