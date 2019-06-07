// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parser_test

import (
	"errors"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	"github.com/opsidian/basil/parser"
	"github.com/opsidian/basil/test"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text/terminal"
)

var _ = Describe("Expression", func() {

	p := parser.Expression()

	DescribeTable("it evaluates the input correctly",
		func(input string, expected interface{}) {
			test.ExpectParserToEvaluate(p)(input, expected)
		},
		// Literals
		test.TableEntry("1", int64(1)),
		test.TableEntry("1.23", 1.23),
		test.TableEntry(`"abc"`, "abc"),
		test.TableEntry(`"{{"`, "{{"),
		test.TableEntry(`"}}"`, "}}"),
		test.TableEntry("true", true),
		test.TableEntry("false", false),
		test.TableEntry("nil", nil),
		test.TableEntry("[1, 2]", []interface{}{int64(1), int64(2)}),
		test.TableEntry("[1, 2][1]", int64(2)),
		test.TableEntry(`[1, "foo"]`, []interface{}{int64(1), "foo"}),
		test.TableEntry(`[1, test.field_string]`, []interface{}{int64(1), "bar"}),
		test.TableEntry("[]", []interface{}{}),
		test.TableEntry("[nil]", []interface{}{nil}),
		test.TableEntry("test.field_string", "bar"),
		test.TableEntry(`test.field_map["key1"]`, "value1"),
		test.TableEntry(`test.field_map["key2"]["key3"]`, "value3"),
		test.TableEntry(`test.field_array[0]`, "value1"),
		test.TableEntry(`test.field_array[test.field_int]`, "value2"),
		test.TableEntry(`test.field_array[2][0]`, "value2"),
		test.TableEntry(`1h30m`, time.Hour+30*time.Minute),

		// Function
		test.TableEntry(`test.func1(test.func1("fOO"))`, "FOO"),

		// Boolean not
		test.TableEntry("!true", !true),
		test.TableEntry("! true", !true),
		test.TableEntry("!false", !false),

		// Prod
		test.TableEntry("2 * 3", int64(2*3)),
		test.TableEntry("4 / 3", int64(4/3)),
		test.TableEntry("2 * 3 / 4", int64(2*3/4)),
		test.TableEntry("2 / 3 * 4", int64(2/3*4)),
		test.TableEntry("2 * -3", int64(2*-3)),

		// Modulo
		test.TableEntry("5 % 2", int64(5%2)),
		test.TableEntry("10 % 7 % 2", int64(10%7%2)),

		// Sum
		test.TableEntry("1 + 2", int64(1+2)),
		test.TableEntry("1 - 2", int64(1-2)),
		test.TableEntry("1 + 2 - 3", int64(1+2-3)),
		test.TableEntry("1 - 2 + 3", int64(1-2+3)),
		test.TableEntry("0 - -1", int64(0 - -1)),

		// String concatenation
		test.TableEntry(`"abc" + "def"`, "abcdef"),
		test.TableEntry(`"abc" + "def" + "ghi"`, "abcdefghi"),

		// Compare
		test.TableEntry("1 == 1", 1 == 1),
		test.TableEntry("1 != 1", 1 != 1),
		test.TableEntry("1 > 2", 1 > 2),
		test.TableEntry("2 >= 2", 2 >= 2),
		test.TableEntry("1 < 2", 1 < 2),
		test.TableEntry("2 <= 3", 2 <= 3),
		test.TableEntry("1 == 1.0", true),
		test.TableEntry("1 == 0.9999999", false),
		test.TableEntry("1 == 0.999999999", true),
		test.TableEntry("1 == 1 == true", 1 == 1 == true),
		test.TableEntry("1 == 1 != false", 1 == 1 != false),

		// And
		test.TableEntry("true && true", true),
		test.TableEntry("true && false", true && false),
		test.TableEntry("true && false && true", false),

		// Or
		test.TableEntry("false || true", true),
		test.TableEntry("false || false", false),
		test.TableEntry("false || true || false", true),

		// Ternary
		test.TableEntry("true ? 1 : 2", int64(1)),
		test.TableEntry("false ? 1 : 2", int64(2)),

		// Mixed
		test.TableEntry("1 + 2 * 3 % 4 - 5 / 6", int64(1+2*3%4-5/6)),
		test.TableEntry("1 + 2 * 3 / 4 + 5 - 6 + 7 * 8 + 9 / 10", int64(1+2*3/4+5-6+7*8+9/10)),
		test.TableEntry("true && false || true", true && false || true),
		test.TableEntry("1 == 1 && 3 > 2", 1 == 1 && 3 > 2),
		test.TableEntry("!true && false || false", !true && false || false),
		test.TableEntry("!(true && false) || false", !(true && false) || false),

		// Using parentheses
		test.TableEntry("(1 + 2) * 3", int64((1+2)*3)),
		test.TableEntry("(1 + 2) * 3 / ((4 + 5) -(6 + 7)) * (8 + 9) / 10", int64((1+2)*3/((4+5)-(6+7))*(8+9)/10)),
	)

	DescribeTable("it returns a parse error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveParseError(p)(input, expectedErr)
		},
		// Literals
		test.TableEntry(`"aaa`, errors.New("was expecting '\"' at testfile:1:5")),
		test.TableEntry("`aaa", errors.New("was expecting '`' at testfile:1:5")),

		// Function
		test.TableEntry(`floor(`, errors.New("was expecting \")\" at testfile:1:7")),
		test.TableEntry(`floor(1,`, errors.New("was expecting value at testfile:1:9")),
		test.TableEntry("FLOOR()", errors.New("was expecting value at testfile:1:1")),

		// Variables
		test.TableEntry("VAR", errors.New("was expecting value at testfile:1:1")),

		// Boolean not
		test.TableEntry("!", errors.New("was expecting value at testfile:1:2")),

		// Prod
		test.TableEntry("1 *", errors.New("was expecting value at testfile:1:4")),
		test.TableEntry("1 /", errors.New("was expecting value at testfile:1:4")),

		// Modulo
		test.TableEntry("1 %", errors.New("was expecting value at testfile:1:4")),

		// Sum
		test.TableEntry("1 +", errors.New("was expecting value at testfile:1:4")),
		test.TableEntry("1 -", errors.New("was expecting value at testfile:1:4")),

		// Compare
		test.TableEntry("1 ==", errors.New("was expecting value at testfile:1:5")),
		test.TableEntry("1 !=", errors.New("was expecting value at testfile:1:5")),
		test.TableEntry("1 >", errors.New("was expecting value at testfile:1:4")),
		test.TableEntry("1 >=", errors.New("was expecting value at testfile:1:5")),
		test.TableEntry("1 <", errors.New("was expecting value at testfile:1:4")),
		test.TableEntry("1 <=", errors.New("was expecting value at testfile:1:5")),

		// And/or
		test.TableEntry("true &&", errors.New("was expecting value at testfile:1:8")),
		test.TableEntry("false ||", errors.New("was expecting value at testfile:1:9")),

		// Ternary
		test.TableEntry("true ?", errors.New("was expecting value at testfile:1:7")),
		test.TableEntry("true ? `a`", errors.New("was expecting \":\" at testfile:1:11")),
		test.TableEntry("true ? `a` :", errors.New("was expecting value at testfile:1:13")),

		// Parentheses
		test.TableEntry("(1 + 2", errors.New("was expecting \")\" at testfile:1:7")),

		// Index
		test.TableEntry(`([0, 1])[1]`, errors.New("was expecting the end of input at testfile:1:9")),
		test.TableEntry(`1[1]`, errors.New("was expecting the end of input at testfile:1:2")),
	)

	DescribeTable("it returns an static check error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveStaticCheckError(p)(input, expectedErr)
		},
		// And
		test.TableEntry(`true && "a"`, errors.New("was expecting boolean at testfile:1:9")),

		// Or
		test.TableEntry(`true || "a"`, errors.New("was expecting boolean at testfile:1:9")),

		// Variable
		test.TableEntry(`non.existing`, errors.New("block \"non\" does not exist at testfile:1:1")),
		test.TableEntry(`test.nonexisting`, errors.New("parameter \"nonexisting\" does not exist at testfile:1:6")),

		// Functions
		test.TableEntry(`non_existing()`, errors.New("\"non_existing\" function does not exist at testfile:1:1")),
	)

	DescribeTable("it returns an eval error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveEvalError(p)(input, expectedErr)
		},
		// Not
		test.TableEntry(`!5`, errors.New("unsupported ! operation on int64 at testfile:1:1")),

		// ProdMod
		test.TableEntry(`5 * "a"`, errors.New("unsupported * operation on int64 and string at testfile:1:3")),
		test.TableEntry(`5 / "a"`, errors.New("unsupported / operation on int64 and string at testfile:1:3")),
		test.TableEntry(`5 % "a"`, errors.New("unsupported % operation on int64 and string at testfile:1:3")),

		// Sum
		test.TableEntry(`1 + "a"`, errors.New("unsupported + operation on int64 and string at testfile:1:3")),
		test.TableEntry(`1 - "a"`, errors.New("unsupported - operation on int64 and string at testfile:1:3")),

		// Compare
		test.TableEntry(`1 == "a"`, errors.New("unsupported == operation on int64 and string at testfile:1:3")),
		test.TableEntry(`1 != "a"`, errors.New("unsupported != operation on int64 and string at testfile:1:3")),
		test.TableEntry(`1 > "a"`, errors.New("unsupported > operation on int64 and string at testfile:1:3")),
		test.TableEntry(`1 >= "a"`, errors.New("unsupported >= operation on int64 and string at testfile:1:3")),
		test.TableEntry(`1 < "a"`, errors.New("unsupported < operation on int64 and string at testfile:1:3")),
		test.TableEntry(`1 <= "a"`, errors.New("unsupported <= operation on int64 and string at testfile:1:3")),

		// Ternary
		test.TableEntry("1 ? 2 : 3", errors.New("expecting bool, got int64 at testfile:1:1")),

		// Variable
		test.TableEntry(`test.field_array[3]`, errors.New("array index out of bounds: 3 (0..2) at testfile:1:18")),
		test.TableEntry(`test.field_array["key"]`, errors.New("non-integer index on array at testfile:1:18")),
		test.TableEntry(`test.field_map["nooo"]`, errors.New("key \"nooo\" does not exist on map at testfile:1:16")),
		test.TableEntry(`test.field_map[1]`, errors.New("invalid non-string index on map at testfile:1:16")),
	)

	Context("When there is only one node", func() {
		It("should return the node", func() {
			expectedNode := terminal.NewIntegerNode(int64(1), parsley.Pos(1), parsley.Pos(2))
			test.ExpectParserToReturn(p, "1", expectedNode)
		})
	})

})
