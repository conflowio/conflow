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

	"github.com/conflowio/conflow/pkg/conflow/function"
	"github.com/conflowio/conflow/pkg/functions/strings"
	"github.com/conflowio/conflow/pkg/parsers"
	"github.com/conflowio/conflow/pkg/schema"
	"github.com/conflowio/conflow/pkg/test"
)

var _ = Describe("TrimSuffix", func() {

	registry := function.InterpreterRegistry{
		"test": strings.TrimSuffixInterpreter{},
	}

	DescribeTable("it evaluates the input correctly",
		func(input string, expected interface{}) {
			test.ExpectFunctionToEvaluate(parsers.Expression(), registry)(input, expected)
		},
		test.TableEntry(`test("", "")`, ""),
		test.TableEntry(`test("foo", "")`, "foo"),
		test.TableEntry(`test("foo", "oo")`, "f"),
		test.TableEntry(`test("foo", "boo")`, "foo"),
		test.TableEntry(`test("foo", "oO")`, "foo"),
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

	It("should return with string type", func() {
		test.ExpectFunctionNode(parsers.Expression(), registry)(
			`test("", "")`,
			func(userCtx interface{}, node parsley.Node) {
				Expect(node.Schema().(schema.Schema).Type()).To(Equal(schema.TypeString))
			},
		)
	})

})
