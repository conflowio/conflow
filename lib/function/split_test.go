package function_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	basilfunction "github.com/opsidian/basil/function"
	"github.com/opsidian/basil/lib/function"
	"github.com/opsidian/basil/parser"
	"github.com/opsidian/basil/test"
)

var _ = Describe("Split", func() {

	registry := basilfunction.Registry{
		"test": function.SplitInterpreter{},
	}

	DescribeTable("it evaluates the input correctly",
		func(input string, expected interface{}) {
			test.ExpectFunctionToEvaluate(parser.Expression(), registry)(input, expected)
		},
		test.TableEntry(`test("", ",")`, []interface{}{""}),
		test.TableEntry(`test("foo", "")`, []interface{}{"f", "o", "o"}),
		test.TableEntry(`test("foo", ",")`, []interface{}{"foo"}),
		test.TableEntry(`test("a,b", ",")`, []interface{}{"a", "b"}),
	)

	DescribeTable("it will have a parse error",
		func(input string, expectedErr error) {
			test.ExpectFunctionToHaveParseError(parser.Expression(), registry)(input, expectedErr)
		},
		test.TableEntry(`test()`, errors.New("test expects 2 arguments at testfile:1:1")),
		test.TableEntry(`test("foo")`, errors.New("test expects 2 arguments at testfile:1:1")),
		test.TableEntry(`test("foo", "bar", "baz")`, errors.New("test expects 2 arguments at testfile:1:1")),
		test.TableEntry(`test(1, "foo")`, errors.New("was expecting string at testfile:1:6")),
		test.TableEntry(`test("foo", 1)`, errors.New("was expecting string at testfile:1:13")),
	)

})
