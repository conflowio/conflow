package parser_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	"github.com/opsidian/basil/parser"
	"github.com/opsidian/basil/test"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/text/terminal"
)

var _ = Describe("Function", func() {

	var p *combinator.Sequence

	q := combinator.Choice(
		terminal.String(false),
		terminal.Integer(),
		terminal.Nil("nil"),
		parser.Variable(),
		test.EvalErrorParser(),
	).Name("value")

	p = parser.Function(q)

	DescribeTable("it evaluates the input correctly",
		func(input string, expected interface{}) {
			test.ExpectParserToEvaluate(p)(input, expected)
		},
		test.TableEntry(`test_func0()`, "Hello"),
		test.TableEntry(`test_func1("foo")`, "FOO"),
		test.TableEntry(`test_func2("foo", "bar")`, "foobar"),
	)

	DescribeTable("it returns a parse error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveParseError(p)(input, expectedErr)
		},
		test.TableEntry("test_func1", errors.New("was expecting \"(\" at testfile:1:11")),
		test.TableEntry("test_func1(", errors.New("was expecting \")\" at testfile:1:12")),
		test.TableEntry(`test_func1("foo"`, errors.New("was expecting \")\" at testfile:1:17")),
		test.TableEntry(`test_func1("foo",`, errors.New("was expecting value at testfile:1:18")),
	)

	DescribeTable("it returns a static check error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveStaticCheckError(p)(input, expectedErr)
		},
		test.TableEntry(`test_func0("foo")`, errors.New("test_func0 expects 0 arguments at testfile:1:1")),
		test.TableEntry(`test_func1(5)`, errors.New("was expecting string at testfile:1:12")),
	)

	DescribeTable("it returns an eval error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveEvalError(p)(input, expectedErr)
		},
		test.TableEntry(`nonexisting("foo")`, errors.New("\"nonexisting\" function does not exist at testfile:1:1")),
		test.TableEntry("test_func1(ERR)", errors.New("ERR at testfile:1:12")),
	)

})
