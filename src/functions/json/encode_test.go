// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package json_test

import (
	"errors"

	"github.com/conflowio/parsley/parsley"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/conflowio/conflow/src/conflow/function"
	"github.com/conflowio/conflow/src/functions/json"
	"github.com/conflowio/conflow/src/parsers"
	"github.com/conflowio/conflow/src/schema"
	"github.com/conflowio/conflow/src/test"
)

var _ = Describe("Encode", func() {

	registry := function.InterpreterRegistry{
		"test": json.EncodeInterpreter{},
	}

	DescribeTable("it evaluates the input correctly",
		func(input string, expected interface{}) {
			test.ExpectFunctionToEvaluate(parsers.Expression(), registry)(input, expected)
		},
		test.TableEntry(`test(null)`, "null"),
		test.TableEntry(`test(1)`, "1"),
		test.TableEntry(`test(1.1)`, "1.1"),
		test.TableEntry(`test("foo")`, "\"foo\""),
		test.TableEntry(`test(true)`, "true"),
		test.TableEntry(`test(1m30s)`, "90000000000"),
		test.TableEntry(`test([1, 2])`, "[1,2]"),
	)

	DescribeTable("it will have a parse error",
		func(input string, expectedErr error) {
			test.ExpectFunctionToHaveParseError(parsers.Expression(), registry)(input, expectedErr)
		},
		test.TableEntry(`test()`, errors.New("test requires exactly 1 argument, but got 0 at testfile:1:1")),
		test.TableEntry(`test("a", "a")`, errors.New("test requires exactly 1 argument, but got 2 at testfile:1:11")),
	)

	It("should return with string type", func() {
		test.ExpectFunctionNode(parsers.Expression(), registry)(
			`test("")`,
			func(userCtx interface{}, node parsley.Node) {
				Expect(node.Schema().(schema.Schema).Type()).To(Equal(schema.TypeString))
			},
		)
	})

})
