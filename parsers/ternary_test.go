// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsers_test

import (
	"errors"

	"github.com/opsidian/basil/basil/schema"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	"github.com/opsidian/basil/parsers"
	"github.com/opsidian/basil/test"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text/terminal"
)

var _ = Describe("TernaryIf", func() {

	q := combinator.Choice(
		terminal.Bool(schema.BooleanValue(), "true", "false"),
		terminal.Integer(schema.IntegerValue()),
		terminal.Float(schema.NumberValue()),
		terminal.String(schema.StringValue(), false),
		terminal.Nil(schema.NullValue(), "NULL"),
		parsers.Array(terminal.String(schema.StringValue(), false)),
		parsers.Map(terminal.String(schema.StringValue(), false)),
		test.EvalErrorParser(schema.BooleanValue(), "ERR_BOOL"),
		test.EvalErrorParser(schema.IntegerValue(), "ERR_INT"),
	).Name("value")

	p := parsers.TernaryIf(q)

	DescribeTable("it evaluates the input correctly",
		func(input string, expected interface{}) {
			test.ExpectParserToEvaluate(p)(input, expected)
		},
		test.TableEntry("1", int64(1)),
		test.TableEntry("NULL", nil),
		test.TableEntry("true ? 1 : 2", int64(1)),
		test.TableEntry("false ? 1 : 2", int64(2)),
		test.TableEntry("true ? 1 : 2.0", int64(1)),
		test.TableEntry("false ? 1 : 2.0", float64(2)),
		test.TableEntry("false ? NULL : NULL", nil),
		test.TableEntry(`true ? ["foo"] : NULL`, []interface{}{"foo"}),
		test.TableEntry(`false ? NULL : ["foo"]`, []interface{}{"foo"}),
		test.TableEntry(`true ? map{"foo":"bar"} : NULL`, map[string]interface{}{"foo": "bar"}),
		test.TableEntry(`false ? NULL : map{"foo":"bar"}`, map[string]interface{}{"foo": "bar"}),
	)

	DescribeTable("it returns a parse error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveParseError(p)(input, expectedErr)
		},
		test.TableEntry("true ?", errors.New("was expecting value at testfile:1:7")),
		test.TableEntry("true ? 1", errors.New("was expecting \":\" at testfile:1:9")),
		test.TableEntry("true ? 1 :", errors.New("was expecting value at testfile:1:11")),
	)

	DescribeTable("it returns a static check error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveStaticCheckError(p)(input, expectedErr)
		},
		test.TableEntry("1 ? 1 : 2", errors.New("must be boolean at testfile:1:1")),
		test.TableEntry(`true ? "foo" : 1`, errors.New("both expressions must have the same type, but got string and integer at testfile:1:8")),
		test.TableEntry("true ? NULL : 2", errors.New("both expressions must have the same type, but got null and integer at testfile:1:8")),
	)

	DescribeTable("it returns an eval error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveEvalError(p)(input, expectedErr)
		},
		test.TableEntry("ERR_BOOL", errors.New("ERR at testfile:1:1")),
		test.TableEntry("ERR_BOOL ? 1 : 2", errors.New("ERR at testfile:1:1")),
		test.TableEntry("true ? ERR_INT : 2", errors.New("ERR at testfile:1:8")),
		test.TableEntry("false ? 1 : ERR_INT", errors.New("ERR at testfile:1:13")),
	)

	Context("When there is only one node", func() {
		It("should return the node", func() {
			expectedNode := terminal.NewIntegerNode(schema.IntegerValue(), int64(1), parsley.Pos(1), parsley.Pos(2))
			test.ExpectParserToReturn(p, "1", expectedNode)
		})
	})

})
