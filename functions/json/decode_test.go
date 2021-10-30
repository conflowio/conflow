// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package json_test

import (
	"errors"

	"github.com/conflowio/conflow/conflow/schema"

	"github.com/conflowio/parsley/parsley"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/conflowio/conflow/conflow/function"
	"github.com/conflowio/conflow/functions/json"
	"github.com/conflowio/conflow/parsers"
	"github.com/conflowio/conflow/test"
)

var _ = Describe("Decode", func() {

	registry := function.InterpreterRegistry{
		"test": json.DecodeInterpreter{},
	}

	DescribeTable("it evaluates the input correctly",
		func(input string, expected interface{}) {
			test.ExpectFunctionToEvaluate(parsers.Expression(), registry)(input, expected)
		},
		test.TableEntry(`test("null")`, nil),
		test.TableEntry(`test("1")`, int64(1)),
		test.TableEntry(`test("1.1")`, 1.1),
		test.TableEntry(`test("\"foo\"")`, "foo"),
		test.TableEntry(`test("true")`, true),
		test.TableEntry(`test("[1, \"foo\"]")`, []interface{}{int64(1), "foo"}),
		test.TableEntry(`test("{\"a\": 1, \"b\": [1, \"foo\"]}")`, map[string]interface{}{
			"a": int64(1),
			"b": []interface{}{int64(1), "foo"},
		}),
	)

	DescribeTable("it will have a parse error",
		func(input string, expectedErr error) {
			test.ExpectFunctionToHaveParseError(parsers.Expression(), registry)(input, expectedErr)
		},
		test.TableEntry(`test()`, errors.New("test requires exactly 1 argument, but got 0 at testfile:1:1")),
		test.TableEntry(`test("a", "a")`, errors.New("test requires exactly 1 argument, but got 2 at testfile:1:11")),
		test.TableEntry(`test(1)`, errors.New("must be string at testfile:1:6")),
	)

	DescribeTable("it will have an eval error",
		func(input string, expectedErr error) {
			test.ExpectFunctionToHaveEvalError(parsers.Expression(), registry)(input, expectedErr)
		},
		test.TableEntry(`test("\"a")`, errors.New("decoding JSON failed: unexpected EOF at testfile:1:6")),
	)

	It("should return with interface{} type", func() {
		test.ExpectFunctionNode(parsers.Expression(), registry)(
			`test("")`,
			func(userCtx interface{}, node parsley.Node) {
				Expect(node.Schema().(schema.Schema).Type()).To(Equal(schema.TypeUntyped))
			},
		)
	})

})
