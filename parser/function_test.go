package parser_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	"github.com/opsidian/ocl/parser"
	"github.com/opsidian/ocl/test"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/text/terminal"
)

var _ = Describe("Function", func() {

	var p *combinator.Sequence

	q := combinator.Choice(
		terminal.String(false),
		terminal.Nil("nil"),
		parser.Variable(p),
		test.EvalErrorParser(),
	).Name("value")

	p = parser.Function(q)

	DescribeTable("it evaluates the input correctly",
		func(input string, expected interface{}) {
			test.ExpectParserToEvaluate(p)(input, expected)
		},
		test.TableEntry(`rand()`, int64(123)),
		test.TableEntry(`upper("foo")`, "FOO"),
		test.TableEntry(`default(nil, "default")`, "default"),
	)

	DescribeTable("it returns a parse error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveParseError(p)(input, expectedErr)
		},
		test.TableEntry("upper", errors.New("was expecting \"(\" at testfile:1:6")),
		test.TableEntry("upper(", errors.New("was expecting \")\" at testfile:1:7")),
		test.TableEntry(`upper("foo"`, errors.New("was expecting \")\" at testfile:1:12")),
		test.TableEntry(`upper("foo",`, errors.New("was expecting value at testfile:1:13")),
	)

	DescribeTable("it returns an eval error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveEvalError(p)(input, expectedErr)
		},
		test.TableEntry(`nonexisting("foo")`, errors.New("function does not exist at testfile:1:1")),
		test.TableEntry("upper(ERR)", errors.New("ERR at testfile:1:7")),
	)

})
