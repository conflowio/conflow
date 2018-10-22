package parser_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	"github.com/opsidian/ocl/parser"
	"github.com/opsidian/ocl/test"
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/parsley"
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
		test.TableEntry(`[1, foo]`, []interface{}{int64(1), "bar"}),
		test.TableEntry("[]", []interface{}{}),
		test.TableEntry("[nil]", []interface{}{nil}),
		test.TableEntry("foo", "bar"),
		test.TableEntry("map.key1", "value1"),
		test.TableEntry(`map["key1"]`, "value1"),
		test.TableEntry(`map["key2"]["key3"]`, "value3"),
		test.TableEntry(`map.key2["key3"]`, "value3"),
		test.TableEntry(`map["key2"].key3`, "value3"),
		test.TableEntry(`arr[0]`, "value1"),
		test.TableEntry(`arr[intkey]`, []interface{}{"value2"}),
		test.TableEntry(`arr[1][0]`, "value2"),

		// Function
		test.TableEntry(`upper(upper("fOO"))`, "FOO"),

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
		test.TableEntry("a &&", errors.New("was expecting value at testfile:1:5")),
		test.TableEntry("b ||", errors.New("was expecting value at testfile:1:5")),

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

		// And
		test.TableEntry(`true && "a"`, errors.New("unsupported && operation on string at testfile:1:6")),

		// Or
		test.TableEntry(`true || "a"`, errors.New("unsupported || operation on string at testfile:1:6")),

		// Ternary
		test.TableEntry("1 ? 2 : 3", errors.New("expecting bool, got int64 at testfile:1:1")),

		// Variable
		test.TableEntry(`a`, errors.New("variable 'a' does not exist at testfile:1:1")),
		test.TableEntry(`arr[3]`, errors.New("array index out of bounds: 3 (0..2) at testfile:1:5")),
		test.TableEntry(`arr["key"]`, errors.New("invalid non-integer index on array at testfile:1:5")),
		test.TableEntry(`map["nooo"]`, errors.New("variable 'map[nooo]' does not exist at testfile:1:1")),
		test.TableEntry(`map[1]`, errors.New("invalid non-string index on map at testfile:1:5")),

		// Functions
		test.TableEntry(`non_existing()`, errors.New("function does not exist at testfile:1:1")),
	)

	Context("When there is only one node", func() {
		It("should return the node", func() {
			expectedNode := ast.NewTerminalNode("INT", int64(1), parsley.Pos(1), parsley.Pos(2))
			test.ExpectParserToReturn(p, "1", expectedNode)
		})
	})

})
