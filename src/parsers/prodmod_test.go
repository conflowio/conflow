// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsers_test

import (
	"errors"

	"github.com/conflowio/parsley/combinator"
	"github.com/conflowio/parsley/parsley"
	"github.com/conflowio/parsley/text/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"

	"github.com/conflowio/conflow/src/parsers"
	"github.com/conflowio/conflow/src/schema"
	"github.com/conflowio/conflow/src/test"
)

var _ = Describe("ProdMod", func() {

	q := combinator.Choice(
		terminal.String(schema.StringValue(), false),
		terminal.Float(schema.NumberValue()),
		terminal.Integer(schema.IntegerValue()),
		terminal.Nil(schema.NullValue(), "NULL"),
		test.EvalErrorParser(schema.IntegerValue(), "ERR"),
	).Name("value")

	p := parsers.ProdMod(q)

	DescribeTable("it evaluates the input correctly",
		func(input string, expected interface{}) {
			test.ExpectParserToEvaluate(p)(input, expected)
		},
		test.TableEntry(`"1"`, "1"),
		test.TableEntry("1", int64(1)),
		test.TableEntry("1 * 2", int64(1*2)),
		test.TableEntry("2 * 1", int64(2*1)),
		test.TableEntry("1 * 2.0", 1*2.0),
		test.TableEntry("1.0 * 2", 1.0*2),
		test.TableEntry("1.0 * 2.0", 1.0*2.0),
		test.TableEntry("2 / 1", int64(2/1)),
		test.TableEntry("1 / 2", 1.0/2.0),
		test.TableEntry("1.0 / 2", 1.0/2),
		test.TableEntry("1 / 2.0", 1/2.0),
		test.TableEntry("1.0 / 2.0", 1.0/2.0),
		test.TableEntry("2 % 5", int64(2%5)),
		test.TableEntry("5 % 2", int64(5%2)),
		test.TableEntry("10 % 5 * 3 % 4", int64(10%5*3%4)),
		test.TableEntry("1 / 0.00000001", 1/0.00000001),
	)

	DescribeTable("it returns a parse error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveParseError(p)(input, expectedErr)
		},
		test.TableEntry("5 *", errors.New("was expecting value at testfile:1:4")),
		test.TableEntry("5 /", errors.New("was expecting value at testfile:1:4")),
		test.TableEntry("5 %", errors.New("was expecting value at testfile:1:4")),
	)

	DescribeTable("it returns a static check error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveStaticCheckError(p)(input, expectedErr)
		},
		test.TableEntry(`"a" * 5`, errors.New("unsupported * operation on string and integer at testfile:1:5")),
		test.TableEntry(`5 * "a"`, errors.New("unsupported * operation on integer and string at testfile:1:3")),
		test.TableEntry(`"a" / 5`, errors.New("unsupported / operation on string and integer at testfile:1:5")),
		test.TableEntry(`5 / "a"`, errors.New("unsupported / operation on integer and string at testfile:1:3")),
		test.TableEntry(`"a" % 5`, errors.New("unsupported % operation on string and integer at testfile:1:5")),
		test.TableEntry(`5 % "a"`, errors.New("unsupported % operation on integer and string at testfile:1:3")),

		test.TableEntry(`5 % 1.0`, errors.New("unsupported % operation on integer and number at testfile:1:3")),
		test.TableEntry(`5.0 % 1`, errors.New("unsupported % operation on number and integer at testfile:1:5")),

		test.TableEntry("NULL * 5", errors.New("unsupported * operation on null and integer at testfile:1:6")),
		test.TableEntry("NULL / 5", errors.New("unsupported / operation on null and integer at testfile:1:6")),
		test.TableEntry("NULL % 5", errors.New("unsupported % operation on null and integer at testfile:1:6")),
		test.TableEntry("NULL * 5.0", errors.New("unsupported * operation on null and number at testfile:1:6")),
		test.TableEntry("NULL / 5.0", errors.New("unsupported / operation on null and number at testfile:1:6")),
		test.TableEntry("NULL % 5.0", errors.New("unsupported % operation on null and number at testfile:1:6")),
		test.TableEntry(`NULL * "foo"`, errors.New("unsupported * operation on null and string at testfile:1:6")),
	)

	DescribeTable("it returns an eval error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveEvalError(p)(input, expectedErr)
		},
		test.TableEntry("1 / 0", errors.New("division by zero at testfile:1:5")),
		test.TableEntry("1 / 0.0", errors.New("division by zero at testfile:1:5")),
		test.TableEntry("1 / 0.0000000001", errors.New("division by zero at testfile:1:5")),
		test.TableEntry("ERR", errors.New("ERR at testfile:1:1")),
		test.TableEntry("ERR * 5", errors.New("ERR at testfile:1:1")),
		test.TableEntry("5 * ERR", errors.New("ERR at testfile:1:5")),
	)

	Context("When there is only one node", func() {
		It("should return the node", func() {
			expectedNode := terminal.NewIntegerNode(schema.IntegerValue(), int64(1), parsley.Pos(1), parsley.Pos(2))
			test.ExpectParserToReturn(p, "1", expectedNode)
		})
	})

})
