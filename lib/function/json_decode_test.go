package function_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	basilfunction "github.com/opsidian/basil/basil/function"
	"github.com/opsidian/basil/basil/variable"
	"github.com/opsidian/basil/lib/function"
	"github.com/opsidian/basil/parser"
	"github.com/opsidian/basil/test"
	"github.com/opsidian/parsley/parsley"
)

var _ = Describe("JSONDecode", func() {

	registry := basilfunction.InterpreterRegistry{
		"test": function.JSONDecodeInterpreter{},
	}

	DescribeTable("it evaluates the input correctly",
		func(input string, expected interface{}) {
			test.ExpectFunctionToEvaluate(parser.Expression(), registry)(input, expected)
		},
		test.TableEntry(`test("null")`, nil),
		test.TableEntry(`test("1")`, int64(1)),
		test.TableEntry(`test("1.1")`, 1.1),
		test.TableEntry(`test("\"foo\"")`, "foo"),
		test.TableEntry(`test("true")`, true),
		test.TableEntry(`test("[1, \"foo\"]")`, []interface{}{int64(1), "foo"}),
		test.TableEntry(`test("{\"a\": 1, \"b\": [1, \"foo\"]}")`, map[string]interface{}{
			"a": int64(1),
			"b": []interface{}{int64(1), "foo"},
		}),
	)

	DescribeTable("it will have a parse error",
		func(input string, expectedErr error) {
			test.ExpectFunctionToHaveParseError(parser.Expression(), registry)(input, expectedErr)
		},
		test.TableEntry(`test()`, errors.New("test expects 1 arguments at testfile:1:1")),
		test.TableEntry(`test("a", "a")`, errors.New("test expects 1 arguments at testfile:1:1")),
		test.TableEntry(`test(1)`, errors.New("was expecting string at testfile:1:6")),
	)

	DescribeTable("it will have an eval error",
		func(input string, expectedErr error) {
			test.ExpectFunctionToHaveEvalError(parser.Expression(), registry)(input, expectedErr)
		},
		test.TableEntry(`test("\"a")`, errors.New("decoding JSON failed: unexpected EOF at testfile:1:6")),
	)

	It("should return with interface{} type", func() {
		test.ExpectFunctionNode(parser.Expression(), registry)(
			`test("")`,
			func(userCtx interface{}, node parsley.Node) {
				Expect(node.Type()).To(Equal(variable.TypeAny))
			},
		)
	})

})
