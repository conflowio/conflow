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
	"github.com/conflowio/parsley/parsley"
	"github.com/conflowio/parsley/text/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"

	"github.com/conflowio/conflow/parsers"
	"github.com/conflowio/conflow/test"
)

var _ = Describe("Or", func() {

	q := combinator.Choice(
		terminal.Bool(schema.BooleanValue(), "true", "false"),
		terminal.Integer(schema.IntegerValue()),
		terminal.Nil(schema.NullValue(), "NULL"),
		test.EvalErrorParser(schema.BooleanValue(), "ERR"),
	).Name("value")

	p := parsers.Or(q)

	DescribeTable("it evaluates the input correctly",
		func(input string, expected interface{}) {
			test.ExpectParserToEvaluate(p)(input, expected)
		},
		test.TableEntry("1", int64(1)),
		test.TableEntry("NULL", nil),
		test.TableEntry("false", false),
		test.TableEntry("true", true),
		test.TableEntry("false || false", false),
		test.TableEntry("false || true", true),
		test.TableEntry("false || false || true", true),
	)

	DescribeTable("it returns a parse error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveParseError(p)(input, expectedErr)
		},
		test.TableEntry("true ||", errors.New("was expecting value at testfile:1:8")),
	)

	DescribeTable("it returns a static check error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveStaticCheckError(p)(input, expectedErr)
		},
		test.TableEntry("1 || 2", errors.New("must be boolean at testfile:1:1")),
		test.TableEntry("NULL || 1", errors.New("must be boolean at testfile:1:1")),
		test.TableEntry("true || 1", errors.New("must be boolean at testfile:1:9")),
	)

	DescribeTable("it returns an eval error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveEvalError(p)(input, expectedErr)
		},
		test.TableEntry("ERR", errors.New("ERR at testfile:1:1")),
		test.TableEntry("true || ERR", errors.New("ERR at testfile:1:9")),
	)

	Context("When there is only one node", func() {
		It("should return the node", func() {
			expectedNode := terminal.NewIntegerNode(schema.IntegerValue(), int64(1), parsley.Pos(1), parsley.Pos(2))
			test.ExpectParserToReturn(p, "1", expectedNode)
		})
	})

})
