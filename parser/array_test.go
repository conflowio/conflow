package parser_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	"github.com/opsidian/ocl/parser"
	"github.com/opsidian/ocl/test"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/terminal"
)

var _ = Describe("Array", func() {

	q := combinator.Choice(
		terminal.String(false),
		terminal.Integer(),
		terminal.Word("nil", nil),
		test.EvalErrorParser(),
	).ReturnError("was expecting value")

	Describe("when new lines are not allowed", func() {

		p := parser.Array(q, text.WsSpaces)

		DescribeTable("it evaluates the input correctly",
			func(input string, expected interface{}) {
				test.ExpectParserToEvaluate(p)(input, expected)
			},
			test.TableEntry("[]", []interface{}{}),
			test.TableEntry("[nil]", []interface{}{nil}),
			test.TableEntry("[1]", []interface{}{int64(1)}),
			test.TableEntry(`[1, "foo"]`, []interface{}{int64(1), "foo"}),
		)

		DescribeTable("it returns a parse error",
			func(input string, expectedErr error) {
				test.ExpectParserToHaveParseError(p)(input, expectedErr)
			},
			test.TableEntry("[", errors.New("was expecting \"]\" at testfile:1:2")),
			test.TableEntry("[1", errors.New("was expecting \"]\" at testfile:1:3")),
			test.TableEntry("[1,", errors.New("was expecting value at testfile:1:4")),
		)

		DescribeTable("it returns an eval error",
			func(input string, expectedErr error) {
				test.ExpectParserToHaveEvalError(p)(input, expectedErr)
			},
			test.TableEntry("[ERR, 1]", errors.New("ERR at testfile:1:2")),
			test.TableEntry("[1, ERR]", errors.New("ERR at testfile:1:5")),
		)
	})

})
