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

var _ = Describe("Len", func() {

	registry := basilfunction.Registry{
		"test": function.LenInterpreter{},
	}

	DescribeTable("it evaluates the input correctly",
		func(input string, expected interface{}) {
			test.ExpectFunctionToEvaluate(parser.Expression(), registry)(input, expected)
		},
		test.TableEntry(`test("")`, int64(0)),
		test.TableEntry(`test("foo")`, int64(3)),
		test.TableEntry(`test("want some 🍕?")`, int64(12)),
		test.TableEntry(`test([])`, int64(0)),
		test.TableEntry(`test([1, 2])`, int64(2)),
	)

	DescribeTable("it will have a parse error",
		func(input string, expectedErr error) {
			test.ExpectFunctionToHaveParseError(parser.Expression(), registry)(input, expectedErr)
		},
		test.TableEntry(`test()`, errors.New("test expects 1 arguments at testfile:1:1")),
		test.TableEntry(`test(1, 2)`, errors.New("test expects 1 arguments at testfile:1:1")),
		test.TableEntry(`test(1)`, errors.New("was expecting string, array or map at testfile:1:6")),
	)

	It("should return with integer type", func() {
		test.ExpectFunctionNode(parser.Expression(), registry)(
			`test("")`,
			func(userCtx interface{}, node parsley.Node) {
				Expect(node.Type()).To(Equal(variable.TypeInteger))
			},
		)
	})

})
