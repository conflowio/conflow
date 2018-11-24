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

var _ = Describe("Replace", func() {

	registry := basilfunction.InterpreterRegistry{
		"test": function.ReplaceInterpreter{},
	}

	DescribeTable("it evaluates the input correctly",
		func(input string, expected interface{}) {
			test.ExpectFunctionToEvaluate(parser.Expression(), registry)(input, expected)
		},
		test.TableEntry(`test("", "", "")`, ""),
		test.TableEntry(`test("foo", "", "")`, "foo"),
		test.TableEntry(`test("foo", "", "x")`, "xfxoxox"),
		test.TableEntry(`test("abcd", "abc", "xxx")`, "xxxd"),
		test.TableEntry(`test("abcd", "bcd", "xxx")`, "axxx"),
		test.TableEntry(`test("abcd", "bc", "xx")`, "axxd"),
		test.TableEntry(`test("abcdabcd", "bc", "xx")`, "axxdaxxd"),
		test.TableEntry(`test("abcd", "Bc", "xx")`, "abcd"),
	)

	DescribeTable("it will have a parse error",
		func(input string, expectedErr error) {
			test.ExpectFunctionToHaveParseError(parser.Expression(), registry)(input, expectedErr)
		},
		test.TableEntry(`test("a", "a")`, errors.New("test expects 3 arguments at testfile:1:1")),
		test.TableEntry(`test("a", "a", "a", "a")`, errors.New("test expects 3 arguments at testfile:1:1")),
		test.TableEntry(`test(1, "a", "a")`, errors.New("was expecting string at testfile:1:6")),
		test.TableEntry(`test("a", 1, "a")`, errors.New("was expecting string at testfile:1:11")),
		test.TableEntry(`test("a", "a", 1)`, errors.New("was expecting string at testfile:1:16")),
	)

	It("should return with string type", func() {
		test.ExpectFunctionNode(parser.Expression(), registry)(
			`test("", "", "")`,
			func(userCtx interface{}, node parsley.Node) {
				Expect(node.Type()).To(Equal(variable.TypeString))
			},
		)
	})

})
