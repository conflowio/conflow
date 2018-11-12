package parser_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	"github.com/opsidian/basil/parser"
	"github.com/opsidian/basil/test"
	"github.com/opsidian/parsley/combinator"
	pparser "github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/terminal"
)

var _ = Describe("Element", func() {

	var q pparser.Func
	q = combinator.Choice(
		terminal.String(false),
		terminal.Integer(),
		terminal.Nil("nil"),
		parser.Array(&q, text.WsSpaces),
		test.EvalErrorParser(),
		test.MapParser(),
	).Name("value")

	index := combinator.Choice(
		terminal.String(false),
		terminal.Integer(),
		test.EvalErrorParser(),
	).Name("value")

	p := parser.Element(q, index)

	DescribeTable("it evaluates the input correctly",
		func(input string, expected interface{}) {
			test.ExpectParserToEvaluate(p)(input, expected)
		},
		test.TableEntry(`nil`, nil),
		test.TableEntry(`1`, int64(1)),
		test.TableEntry(`[1, 2, 3][0]`, int64(1)),
		test.TableEntry(`[1, 2, 3][1]`, int64(2)),
		test.TableEntry(`[1, [2, 3, 4], 5][1][1]`, int64(3)),
		test.TableEntry(`MAP["a"]`, int64(1)),
		test.TableEntry(`MAP["c"]["d"]`, int64(2)),
		test.TableEntry(`MAP["d"][1]`, "bar"),
	)

	DescribeTable("it returns a parse error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveParseError(p)(input, expectedErr)
		},
		test.TableEntry(`MAP[`, errors.New("was expecting value at testfile:1:5")),
		test.TableEntry(`MAP["key1"`, errors.New("was expecting \"]\" at testfile:1:11")),
		test.TableEntry(`MAP[]`, errors.New("was expecting value at testfile:1:5")),
	)

	DescribeTable("it returns an eval error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveEvalError(p)(input, expectedErr)
		},
		test.TableEntry(`MAP[1]`, errors.New("invalid non-string index on map at testfile:1:5")),
		test.TableEntry(`MAP["non-existing"]`, errors.New("index 'non-existing' does not exist at testfile:1:5")),
		test.TableEntry(`[1, 2, 3]["string"]`, errors.New("invalid non-integer index on array at testfile:1:11")),
		test.TableEntry(`"string"[0]`, errors.New("can not get index on string type at testfile:1:10")),
		test.TableEntry(`nil[0]`, errors.New("can not get index on <nil> type at testfile:1:5")),
		test.TableEntry(`[1, 2, 3][ERR]`, errors.New("ERR at testfile:1:11")),
		test.TableEntry(`ERR[1]`, errors.New("ERR at testfile:1:1")),
		test.TableEntry(`[1, 2, 3][3]`, errors.New("array index out of bounds: 3 (0..2) at testfile:1:11")),
	)

	Context("When there is only one node", func() {
		It("should return the node", func() {
			expectedNode := terminal.NewIntegerNode(int64(1), parsley.Pos(1), parsley.Pos(2))
			test.ExpectParserToReturn(p, "1", expectedNode)
		})
	})

})
