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
		"split": function.SplitInterpreter{},
	}

	DescribeTable("it evaluates the input correctly",
		func(input string, expected interface{}) {
			test.ExpectFunctionToEvaluate(parser.Expression(), registry)(input, expected)
		},
		test.TableEntry(`split("", ",")`, []interface{}{""}),
		test.TableEntry(`split("foo", "")`, []interface{}{"f", "o", "o"}),
		test.TableEntry(`split("foo", ",")`, []interface{}{"foo"}),
		test.TableEntry(`split("a,b", ",")`, []interface{}{"a", "b"}),
	)

	DescribeTable("it will have a parse error",
		func(input string, expectedErr error) {
			test.ExpectFunctionToHaveParseError(parser.Expression(), registry)(input, expectedErr)
		},
		test.TableEntry(`split()`, errors.New("split expects 2 arguments at testfile:1:1")),
		test.TableEntry(`split("foo")`, errors.New("split expects 2 arguments at testfile:1:1")),
		test.TableEntry(`split("foo", "bar", "baz")`, errors.New("split expects 2 arguments at testfile:1:1")),
		test.TableEntry(`split(1, "foo")`, errors.New("was expecting string at testfile:1:7")),
		test.TableEntry(`split("foo", 1)`, errors.New("was expecting string at testfile:1:14")),
	)

})
