// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package strings_test

import (
	"errors"

	"github.com/opsidian/conflow/basil/schema"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/opsidian/conflow/basil/function"
	"github.com/opsidian/conflow/functions/strings"
	"github.com/opsidian/conflow/parsers"
	"github.com/opsidian/conflow/test"
	"github.com/opsidian/parsley/parsley"
)

var _ = Describe("Replace", func() {

	registry := function.InterpreterRegistry{
		"test": strings.ReplaceInterpreter{},
	}

	DescribeTable("it evaluates the input correctly",
		func(input string, expected interface{}) {
			test.ExpectFunctionToEvaluate(parsers.Expression(), registry)(input, expected)
		},
		test.TableEntry(`test("", "", "")`, ""),
		test.TableEntry(`test("foo", "", "")`, "foo"),
		test.TableEntry(`test("foo", "", "x")`, "xfxoxox"),
		test.TableEntry(`test("abcd", "abc", "xxx")`, "xxxd"),
		test.TableEntry(`test("abcd", "bcd", "xxx")`, "axxx"),
		test.TableEntry(`test("abcd", "bc", "xx")`, "axxd"),
		test.TableEntry(`test("abcdabcd", "bc", "xx")`, "axxdaxxd"),
		test.TableEntry(`test("abcd", "Bc", "xx")`, "abcd"),
	)

	DescribeTable("it will have a parse error",
		func(input string, expectedErr error) {
			test.ExpectFunctionToHaveParseError(parsers.Expression(), registry)(input, expectedErr)
		},
		test.TableEntry(`test("a", "a")`, errors.New("test requires exactly 3 arguments, but got 2 at testfile:1:14")),
		test.TableEntry(`test("a", "a", "a", "a")`, errors.New("test requires exactly 3 arguments, but got 4 at testfile:1:21")),
		test.TableEntry(`test(1, "a", "a")`, errors.New("must be string at testfile:1:6")),
		test.TableEntry(`test("a", 1, "a")`, errors.New("must be string at testfile:1:11")),
		test.TableEntry(`test("a", "a", 1)`, errors.New("must be string at testfile:1:16")),
	)

	It("should return with string type", func() {
		test.ExpectFunctionNode(parsers.Expression(), registry)(
			`test("", "", "")`,
			func(userCtx interface{}, node parsley.Node) {
				Expect(node.Schema().(schema.Schema).Type()).To(Equal(schema.TypeString))
			},
		)
	})

})
