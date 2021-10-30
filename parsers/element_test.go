// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsers_test

import (
	"errors"

	"github.com/conflowio/conflow/conflow/schema"

	"github.com/conflowio/parsley/combinator"
	pparser "github.com/conflowio/parsley/parser"
	"github.com/conflowio/parsley/parsley"
	"github.com/conflowio/parsley/text/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"

	"github.com/conflowio/conflow/parsers"
	"github.com/conflowio/conflow/test"
)

var _ = Describe("Element", func() {

	var q pparser.Func
	q = combinator.Choice(
		terminal.String(schema.StringValue(), false),
		terminal.Integer(schema.IntegerValue()),
		terminal.Nil(schema.NullValue(), "NULL"),
		parsers.Array(&q),
		test.EvalErrorParser(&schema.Array{Items: schema.StringValue()}, "ERR_ARRAY"),
		test.EvalErrorParser(&schema.Map{AdditionalProperties: schema.StringValue()}, "ERR_MAP"),
		test.MapParser("MAP", map[string]interface{}{
			"a": []interface{}{int64(1), int64(2)},
			"b": []interface{}{int64(3), int64(4)},
		}),
	).Name("value")

	index := combinator.Choice(
		terminal.String(schema.StringValue(), false),
		terminal.Integer(schema.IntegerValue()),
		test.EvalErrorParser(schema.IntegerValue(), "ERR_INT"),
		test.EvalErrorParser(schema.StringValue(), "ERR_STRING"),
	).Name("value")

	p := parsers.Element(q, index)

	DescribeTable("it evaluates the input correctly",
		func(input string, expected interface{}) {
			test.ExpectParserToEvaluate(p)(input, expected)
		},
		test.TableEntry(`NULL`, nil),
		test.TableEntry(`1`, int64(1)),
		test.TableEntry(`[1, 2, 3][0]`, int64(1)),
		test.TableEntry(`[1, 2, 3][1]`, int64(2)),
		test.TableEntry(`[[1], [2, 3, 4], [5]][1][1]`, int64(3)),
		test.TableEntry(`MAP["a"]`, []interface{}{int64(1), int64(2)}),
		test.TableEntry(`MAP["a"][0]`, int64(1)),
	)

	DescribeTable("it returns a parse error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveParseError(p)(input, expectedErr)
		},
		test.TableEntry(`MAP[`, errors.New("was expecting value at testfile:1:5")),
		test.TableEntry(`MAP["key1"`, errors.New("was expecting \"]\" at testfile:1:11")),
		test.TableEntry(`MAP[]`, errors.New("was expecting value at testfile:1:5")),
	)

	DescribeTable("it returns a static check error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveStaticCheckError(p)(input, expectedErr)
		},
		test.TableEntry(`MAP[1]`, errors.New("must be string at testfile:1:5")),
		test.TableEntry(`[1, 2, 3]["string"]`, errors.New("must be integer at testfile:1:11")),
		test.TableEntry(`"string"[0]`, errors.New("can not get index on string type at testfile:1:10")),
		test.TableEntry(`NULL[0]`, errors.New("can not get index on null type at testfile:1:6")),
	)

	DescribeTable("it returns an eval error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveEvalError(p)(input, expectedErr)
		},

		test.TableEntry(`MAP["non-existing"]`, errors.New("key \"non-existing\" does not exist on map at testfile:1:5")),
		test.TableEntry(`[1, 2, 3][ERR_INT]`, errors.New("ERR at testfile:1:11")),
		test.TableEntry(`MAP[ERR_STRING]`, errors.New("ERR at testfile:1:5")),
		test.TableEntry(`ERR_ARRAY[1]`, errors.New("ERR at testfile:1:1")),
		test.TableEntry(`ERR_MAP["foo"]`, errors.New("ERR at testfile:1:1")),
		test.TableEntry(`[1, 2, 3][3]`, errors.New("array index out of bounds: 3 (0..2) at testfile:1:11")),
	)

	Context("When there is only one node", func() {
		It("should return the node", func() {
			expectedNode := terminal.NewIntegerNode(schema.IntegerValue(), int64(1), parsley.Pos(1), parsley.Pos(2))
			test.ExpectParserToReturn(p, "1", expectedNode)
		})
	})

})
