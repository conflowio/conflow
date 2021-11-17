// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsers_test

import (
	"errors"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"

	"github.com/conflowio/conflow/conflow/schema"

	"github.com/conflowio/parsley/combinator"
	pparser "github.com/conflowio/parsley/parser"
	"github.com/conflowio/parsley/parsley"
	"github.com/conflowio/parsley/text/terminal"

	"github.com/conflowio/conflow/parsers"
	"github.com/conflowio/conflow/test"
)

var _ = Describe("Compare", func() {

	time1, _ := time.Parse(time.RFC3339, "2001-01-01T00:00:01Z")
	time2, _ := time.Parse(time.RFC3339, "2001-01-01T00:00:02Z")

	var q pparser.Func
	q = combinator.Choice(
		terminal.String(schema.StringValue(), false),
		terminal.TimeDuration(&schema.String{Format: schema.FormatDurationGo}),
		terminal.Float(schema.NumberValue()),
		terminal.Integer(schema.IntegerValue()),
		terminal.Bool(schema.BooleanValue(), "true", "false"),
		parsers.Array(&q),
		terminal.Nil(schema.NullValue(), "NULL"),
		test.EvalErrorParser(schema.IntegerValue(), "ERR"),
		terminal.Word(&schema.String{Format: schema.FormatDateTime}, "TIME1", time1),
		terminal.Word(&schema.String{Format: schema.FormatDateTime}, "TIME2", time2),
	).Name("value")

	p := parsers.Compare(q)

	DescribeTable("it evaluates the input correctly",
		func(input string, expected interface{}) {
			test.ExpectParserToEvaluate(p)(input, expected)
		},
		test.TableEntry("1", int64(1)),
		test.TableEntry("NULL", nil),
		test.TableEntry("1.0 == 1", true),
		test.TableEntry("0.999999999 == 1", true),
		test.TableEntry("0.999999999 != 1", false),
		test.TableEntry("1 == 0.999999999", true),
		test.TableEntry("1 != 0.999999999", false),
		test.TableEntry("0.99999999 == 1", false),
		test.TableEntry("0.99999999 != 1", true),
		test.TableEntry("[1, 2] == [1, 2]", true),
		test.TableEntry("[1, 2] != [1, 2]", false),
		test.TableEntry("[1, 2] == [2, 1]", false),
		test.TableEntry("[1, 2] != [2, 1]", true),
		test.TableEntry("false == false", true),
		test.TableEntry("false != false", false),
		test.TableEntry("true == true", true),
		test.TableEntry("true != true", false),
		test.TableEntry("false == true", false),
		test.TableEntry("false != true", true),
		test.TableEntry("true == false", false),
		test.TableEntry("true != false", true),
		test.TableEntry("1 == 2 == false", true),

		// Time
		test.TableEntry("TIME1 == TIME1", true),
		test.TableEntry("TIME1 == TIME2", false),
		test.TableEntry("TIME1 != TIME1", false),
		test.TableEntry("TIME1 != TIME2", true),
		test.TableEntry("TIME1 <= TIME1", true),
		test.TableEntry("TIME1 <= TIME2", true),
		test.TableEntry("TIME2 <= TIME1", false),
		test.TableEntry("TIME1 < TIME1", false),
		test.TableEntry("TIME1 < TIME2", true),
		test.TableEntry("TIME2 < TIME1", false),
		test.TableEntry("TIME1 >= TIME1", true),
		test.TableEntry("TIME1 >= TIME2", false),
		test.TableEntry("TIME2 >= TIME1", true),
		test.TableEntry("TIME1 > TIME1", false),
		test.TableEntry("TIME1 > TIME2", false),
		test.TableEntry("TIME2 > TIME1", true),

		// Time duration
		test.TableEntry("1m2s == 1m2s", true),
		test.TableEntry("1m2s == 1m3s", false),
		test.TableEntry("1m2s != 1m2s", false),
		test.TableEntry("1m2s != 1m3s", true),
		test.TableEntry("1m2s <= 1m2s", true),
		test.TableEntry("1m2s <= 1m3s", true),
		test.TableEntry("1m3s <= 1m2s", false),
		test.TableEntry("1m2s < 1m2s", false),
		test.TableEntry("1m2s < 1m3s", true),
		test.TableEntry("1m3s < 1m2s", false),
		test.TableEntry("1m2s >= 1m2s", true),
		test.TableEntry("1m2s >= 1m3s", false),
		test.TableEntry("1m3s >= 1m2s", true),
		test.TableEntry("1m2s > 1m2s", false),
		test.TableEntry("1m2s > 1m3s", false),
		test.TableEntry("1m3s > 1m2s", true),
	)

	DescribeTable("it returns a parse error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveParseError(p)(input, expectedErr)
		},
		test.TableEntry("1 == ", errors.New("was expecting value at testfile:1:6")),
		test.TableEntry("1 != ", errors.New("was expecting value at testfile:1:6")),
		test.TableEntry("1 < ", errors.New("was expecting value at testfile:1:5")),
		test.TableEntry("1 <= ", errors.New("was expecting value at testfile:1:6")),
		test.TableEntry("1 > ", errors.New("was expecting value at testfile:1:5")),
		test.TableEntry("1 >= ", errors.New("was expecting value at testfile:1:6")),
	)

	DescribeTable("it returns a static check error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveStaticCheckError(p)(input, expectedErr)
		},
		test.TableEntry(`NULL == 5`, errors.New("unsupported == operation on null and integer at testfile:1:6")),
		test.TableEntry(`"foo" == 5`, errors.New("unsupported == operation on string and integer at testfile:1:7")),
		test.TableEntry(`"foo" == 5.5`, errors.New("unsupported == operation on string and number at testfile:1:7")),
		test.TableEntry(`5 == "foo"`, errors.New("unsupported == operation on integer and string at testfile:1:3")),
		test.TableEntry("5 == [1,2]", errors.New("unsupported == operation on integer and array(integer) at testfile:1:3")),
		test.TableEntry("1 == true", errors.New("unsupported == operation on integer and boolean at testfile:1:3")),
	)

	DescribeTable("it returns an eval error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveEvalError(p)(input, expectedErr)
		},
		test.TableEntry("ERR", errors.New("ERR at testfile:1:1")),
		test.TableEntry("ERR == 1", errors.New("ERR at testfile:1:1")),
		test.TableEntry("ERR != 1", errors.New("ERR at testfile:1:1")),
		test.TableEntry("ERR < 1", errors.New("ERR at testfile:1:1")),
		test.TableEntry("ERR <= 1", errors.New("ERR at testfile:1:1")),
		test.TableEntry("ERR > 1", errors.New("ERR at testfile:1:1")),
		test.TableEntry("ERR >= 1", errors.New("ERR at testfile:1:1")),
		test.TableEntry("1 == ERR", errors.New("ERR at testfile:1:6")),
		test.TableEntry("1 != ERR", errors.New("ERR at testfile:1:6")),
		test.TableEntry("1 < ERR", errors.New("ERR at testfile:1:5")),
		test.TableEntry("1 <= ERR", errors.New("ERR at testfile:1:6")),
		test.TableEntry("1 > ERR", errors.New("ERR at testfile:1:5")),
		test.TableEntry("1 >= ERR", errors.New("ERR at testfile:1:6")),
	)

	It("should handle all operators correctly with all types", func() {
		type TC struct {
			input    string
			expected interface{}
		}
		values := [][]string{
			{"1", "2"},
			{"1.1", "1.2"},
			{"1.1", "2"},
			{"1", "1.2"},
			{`"ab"`, `"ac"`},
		}
		baseTestCases := []TC{
			{"P1 == P1", true},
			{"P1 == P2", false},
			{"P2 == P1", false},
			{"P1 != P1", false},
			{"P1 != P2", true},
			{"P1 != P2", true},
			{"P1 > P1", false},
			{"P2 > P1", true},
			{"P1 > P2", false},
			{"P1 >= P1", true},
			{"P2 >= P1", true},
			{"P1 >= P2", false},
			{"P1 < P1", false},
			{"P1 < P2", true},
			{"P2 < P1", false},
			{"P1 <= P1", true},
			{"P1 <= P1", true},
			{"P2 <= P1", false},
		}
		for _, valueSet := range values {
			for _, tc := range baseTestCases {
				input := strings.Replace(tc.input, "P1", valueSet[0], -1)
				input = strings.Replace(input, "P2", valueSet[1], -1)
				test.ExpectParserToEvaluate(p)(input, tc.expected)
			}
		}
	})

	Context("When there is only one node", func() {
		It("should return the node", func() {
			expectedNode := terminal.NewIntegerNode(schema.IntegerValue(), int64(1), parsley.Pos(1), parsley.Pos(2))
			test.ExpectParserToReturn(p, "1", expectedNode)
		})
	})

})
