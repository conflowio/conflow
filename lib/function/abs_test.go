package function_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	basilfunction "github.com/opsidian/basil/function"
	"github.com/opsidian/basil/lib/function"
	"github.com/opsidian/basil/parser"
	"github.com/opsidian/basil/test"
	"github.com/opsidian/basil/variable"
	"github.com/opsidian/parsley/parsley"
)

var _ = Describe("Abs", func() {

	registry := basilfunction.Registry{
		"abs": function.AbsInterpreter{},
	}

	DescribeTable("it evaluates the input correctly",
		func(input string, expected interface{}) {
			res := variable.NewNumber(expected)
			test.ExpectFunctionToEvaluate(parser.Expression(), registry)(input, res)
		},
		test.TableEntry("abs(0)", int64(0)),
		test.TableEntry("abs(1)", int64(1)),
		test.TableEntry("abs(-1)", int64(1)),
		test.TableEntry("abs(0.0)", 0.0),
		test.TableEntry("abs(1.0)", 1.0),
		test.TableEntry("abs(-1.0)", 1.0),
	)

	DescribeTable("it will have a parse error",
		func(input string, expectedErr error) {
			test.ExpectFunctionToHaveParseError(parser.Expression(), registry)(input, expectedErr)
		},
		test.TableEntry(`abs()`, errors.New("abs expects 1 arguments at testfile:1:1")),
		test.TableEntry(`abs(1, 2)`, errors.New("abs expects 1 arguments at testfile:1:1")),
		test.TableEntry(`abs("foo")`, errors.New("was expecting number at testfile:1:5")),
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
		test.TableEntry(`abs(1)`, variable.TypeInteger),
		test.TableEntry(`abs(1.1)`, variable.TypeFloat),
	)

})
