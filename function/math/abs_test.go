package math_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/opsidian/basil/basil/function"
	"github.com/opsidian/basil/basil/variable"
	"github.com/opsidian/basil/function/math"
	"github.com/opsidian/basil/parser"
	"github.com/opsidian/basil/test"
	"github.com/opsidian/parsley/parsley"
)

var _ = Describe("Abs", func() {

	registry := function.InterpreterRegistry{
		"test": math.AbsInterpreter{},
	}

	DescribeTable("it evaluates the input correctly",
		func(input string, expected interface{}) {
			test.ExpectFunctionToEvaluate(parser.Expression(), registry)(input, expected)
		},
		test.TableEntry("test(0)", int64(0)),
		test.TableEntry("test(1)", int64(1)),
		test.TableEntry("test(-1)", int64(1)),
		test.TableEntry("test(0.0)", 0.0),
		test.TableEntry("test(1.0)", 1.0),
		test.TableEntry("test(-1.0)", 1.0),
		test.TableEntry("test(1) + 2", int64(3)),
		test.TableEntry("test(1.1) + 2.2", 3.3),
	)

	DescribeTable("it will have a parse error",
		func(input string, expectedErr error) {
			test.ExpectFunctionToHaveParseError(parser.Expression(), registry)(input, expectedErr)
		},
		test.TableEntry(`test()`, errors.New("test expects 1 arguments at testfile:1:1")),
		test.TableEntry(`test(1, 2)`, errors.New("test expects 1 arguments at testfile:1:1")),
		test.TableEntry(`test("foo")`, errors.New("was expecting number at testfile:1:6")),
	)

	DescribeTable("it should keep the type of the first argument",
		func(input string, expectedType string) {
			test.ExpectFunctionNode(parser.Expression(), registry)(
				input,
				func(userCtx interface{}, node parsley.Node) {
					Expect(node.Type()).To(Equal(expectedType))
				},
			)
		},
		test.TableEntry(`test(1)`, variable.TypeInteger),
		test.TableEntry(`test(1.1)`, variable.TypeFloat),
	)

})
