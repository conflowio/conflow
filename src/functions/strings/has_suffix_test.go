// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package strings_test

import (
	"errors"

	"github.com/conflowio/parsley/parsley"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/conflowio/conflow/src/conflow/function"
	"github.com/conflowio/conflow/src/functions/strings"
	"github.com/conflowio/conflow/src/parsers"
	"github.com/conflowio/conflow/src/schema"
	"github.com/conflowio/conflow/src/test"
)

var _ = Describe("HasSuffix", func() {

	registry := function.InterpreterRegistry{
		"test": strings.HasSuffixInterpreter{},
	}

	DescribeTable("it evaluates the input correctly",
		func(input string, expected interface{}) {
			test.ExpectFunctionToEvaluate(parsers.Expression(), registry)(input, expected)
		},
		test.TableEntry(`test("", "")`, true),
		test.TableEntry(`test("foo", "")`, true),
		test.TableEntry(`test("foo", "oo")`, true),
		test.TableEntry(`test("foo", "boo")`, false),
		test.TableEntry(`test("foo", "oO")`, false),
	)

	DescribeTable("it will have a parse error",
		func(input string, expectedErr error) {
			test.ExpectFunctionToHaveParseError(parsers.Expression(), registry)(input, expectedErr)
		},
		test.TableEntry(`test("foo")`, errors.New("test requires exactly 2 arguments, but got 1 at testfile:1:11")),
		test.TableEntry(`test("a", "a", "a")`, errors.New("test requires exactly 2 arguments, but got 3 at testfile:1:16")),
		test.TableEntry(`test(1, "a")`, errors.New("must be string at testfile:1:6")),
		test.TableEntry(`test("a", 1)`, errors.New("must be string at testfile:1:11")),
	)

	It("should return with boolean type", func() {
		test.ExpectFunctionNode(parsers.Expression(), registry)(
			`test("", "")`,
			func(userCtx interface{}, node parsley.Node) {
				Expect(node.Schema().(schema.Schema).Type()).To(Equal(schema.TypeBoolean))
			},
		)
	})

})
