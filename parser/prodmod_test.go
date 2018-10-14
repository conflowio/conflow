package parser_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	"github.com/opsidian/ocl/parser"
	"github.com/opsidian/ocl/test"
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text/terminal"
)

var _ = Describe("ProdMod", func() {

	q := combinator.Choice(
		terminal.String(false),
		terminal.Float(),
		terminal.Integer(),
		terminal.Word("nil", nil),
		test.EvalErrorParser(),
	).ReturnError("was expecting value")

	p := parser.ProdMod(q)

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
		test.TableEntry("1 / 2", int64(1/2)),
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

	DescribeTable("it returns an eval error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveEvalError(p)(input, expectedErr)
		},
		test.TableEntry("1 / 0", errors.New("divison by zero at testfile:1:5")),
		test.TableEntry("1 / 0.0", errors.New("divison by zero at testfile:1:5")),
		test.TableEntry("1 / 0.000000001", errors.New("divison by zero at testfile:1:5")),
		test.TableEntry("nil * 5", errors.New("unsupported * operation on <nil> and int64 at testfile:1:5")),
		test.TableEntry("nil / 5", errors.New("unsupported / operation on <nil> and int64 at testfile:1:5")),
		test.TableEntry("nil % 5", errors.New("unsupported % operation on <nil> and int64 at testfile:1:5")),
		test.TableEntry("nil * 5.0", errors.New("unsupported * operation on <nil> and float64 at testfile:1:5")),
		test.TableEntry("nil / 5.0", errors.New("unsupported / operation on <nil> and float64 at testfile:1:5")),
		test.TableEntry("nil % 5.0", errors.New("unsupported % operation on <nil> and float64 at testfile:1:5")),
		test.TableEntry(`nil * "foo"`, errors.New("unsupported * operation on <nil> and string at testfile:1:5")),
		test.TableEntry("ERR", errors.New("ERR at testfile:1:1")),
		test.TableEntry("ERR * 5", errors.New("ERR at testfile:1:1")),
		test.TableEntry("5 * ERR", errors.New("ERR at testfile:1:5")),
	)

	Context("When there is only one node", func() {
		It("should return the node", func() {
			expectedNode := ast.NewTerminalNode("INT", int64(1), parsley.Pos(1), parsley.Pos(2))
			test.ExpectParserToReturn(p, "1", expectedNode)
		})
	})

})
