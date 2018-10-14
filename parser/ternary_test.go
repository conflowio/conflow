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

var _ = Describe("TernaryIf", func() {

	q := combinator.Choice(
		terminal.Bool("true", "false"),
		terminal.Integer(),
		terminal.Word("nil", nil),
		test.EvalErrorParser(),
	).ReturnError("was expecting value")

	p := parser.TernaryIf(q)

	DescribeTable("it evaluates the input correctly",
		func(input string, expected interface{}) {
			test.ExpectParserToEvaluate(p)(input, expected)
		},
		test.TableEntry("1", int64(1)),
		test.TableEntry("nil", nil),
		test.TableEntry("true ? 1 : 2", int64(1)),
		test.TableEntry("false ? 1 : 2", int64(2)),
		test.TableEntry("true ? nil : 2", nil),
	)

	DescribeTable("it returns a parse error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveParseError(p)(input, expectedErr)
		},
		test.TableEntry("true ?", errors.New("was expecting value at testfile:1:7")),
		test.TableEntry("true ? 1", errors.New("was expecting \":\" at testfile:1:9")),
		test.TableEntry("true ? 1 :", errors.New("was expecting value at testfile:1:11")),
	)

	DescribeTable("it returns an eval error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveEvalError(p)(input, expectedErr)
		},
		test.TableEntry("1 ? 1 : 2", errors.New("expecting bool, got int64 at testfile:1:1")),
		test.TableEntry("ERR", errors.New("ERR at testfile:1:1")),
		test.TableEntry("ERR ? 1 : 2", errors.New("ERR at testfile:1:1")),
		test.TableEntry("true ? ERR : 2", errors.New("ERR at testfile:1:8")),
		test.TableEntry("false ? 1 : ERR", errors.New("ERR at testfile:1:13")),
	)

	Context("When there is only one node", func() {
		It("should return the node", func() {
			expectedNode := ast.NewTerminalNode("INT", int64(1), parsley.Pos(1), parsley.Pos(2))
			test.ExpectParserToReturn(p, "1", expectedNode)
		})
	})

})
