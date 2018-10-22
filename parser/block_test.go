package parser_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	"github.com/opsidian/ocl/parser"
	"github.com/opsidian/ocl/test"
)

var _ = FDescribe("Block", func() {

	p := parser.Block()

	DescribeTable("it evaluates the input correctly",
		func(input string, expected interface{}) {
			test.ExpectParserToEvaluate(p)(input, expected)
		},
		test.TableEntry(`testblock foo { 
			param1 = "bar"
		}`, test.NewTestBlock("foo", "bar", nil)),
		test.TableEntry(`testblock { 
			param1 = "bar"
		}`, test.NewTestBlock("", "bar", nil)),
	)

	DescribeTable("it returns a parse error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveParseError(p)(input, expectedErr)
		},
	)

	DescribeTable("it returns an eval error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveEvalError(p)(input, expectedErr)
		},
	)

})
