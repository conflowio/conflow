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

var _ = Describe("Join", func() {

	registry := basilfunction.Registry{
		"join": function.JoinInterpreter{},
	}

	DescribeTable("it evaluates the input correctly",
		func(input string, expected interface{}) {
			test.ExpectFunctionToEvaluate(parser.Expression(), registry)(input, expected)
		},
		test.TableEntry(`join(["foo", "bar"], ",")`, string("foo,bar")),
	)

	DescribeTable("it will have a parse error",
		func(input string, expectedErr error) {
			test.ExpectFunctionToHaveParseError(parser.Expression(), registry)(input, expectedErr)
		},
		test.TableEntry(`join()`, errors.New("join expects 2 arguments at testfile:1:1")),
		test.TableEntry(`join([])`, errors.New("join expects 2 arguments at testfile:1:1")),
		test.TableEntry(`join([], "foo", "bar")`, errors.New("join expects 2 arguments at testfile:1:1")),
		test.TableEntry(`join("not array", "foo")`, errors.New("was expecting string array at testfile:1:6")),
		test.TableEntry(`join(["a", 2], "foo")`, errors.New("was expecting string at testfile:1:12")),
		test.TableEntry(`join([], 5)`, errors.New("was expecting string at testfile:1:10")),
	)

})
