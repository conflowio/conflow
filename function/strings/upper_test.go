package strings_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/opsidian/basil/basil/function"
	"github.com/opsidian/basil/basil/variable"
	"github.com/opsidian/basil/function/strings"
	"github.com/opsidian/basil/parser"
	"github.com/opsidian/basil/test"
	"github.com/opsidian/parsley/parsley"
)

var _ = Describe("Upper", func() {

	registry := function.InterpreterRegistry{
		"test": strings.UpperInterpreter{},
	}

	DescribeTable("it evaluates the input correctly",
		func(input string, expected interface{}) {
			test.ExpectFunctionToEvaluate(parser.Expression(), registry)(input, expected)
		},
		test.TableEntry(`test("")`, ""),
		test.TableEntry(`test("FOO")`, "FOO"),
		test.TableEntry(`test("Foo")`, "FOO"),
		test.TableEntry(`test("Ármányos Ödön")`, "ÁRMÁNYOS ÖDÖN"),
	)

	DescribeTable("it will have a parse error",
		func(input string, expectedErr error) {
			test.ExpectFunctionToHaveParseError(parser.Expression(), registry)(input, expectedErr)
		},
		test.TableEntry(`test()`, errors.New("test expects 1 arguments at testfile:1:1")),
		test.TableEntry(`test("a", "a")`, errors.New("test expects 1 arguments at testfile:1:1")),
		test.TableEntry(`test(1)`, errors.New("was expecting string at testfile:1:6")),
	)

	It("should return with string type", func() {
		test.ExpectFunctionNode(parser.Expression(), registry)(
			`test("")`,
			func(userCtx interface{}, node parsley.Node) {
				Expect(node.Type()).To(Equal(variable.TypeString))
			},
		)
	})

})
