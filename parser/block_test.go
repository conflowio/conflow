package parser_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	"github.com/opsidian/ocl/parser"
	"github.com/opsidian/ocl/test"
)

var _ = Describe("Block", func() {

	p := parser.Block()

	DescribeTable("it evaluates the input correctly",
		func(input string, expected interface{}) {
			test.ExpectParserToEvaluate(p)(input, expected)
		},
		test.TableEntry(
			`testblock`,
			test.NewTestBlock("0", "", map[string]*test.TestBlock{}),
		),
		test.TableEntry(
			`testblock foo`,
			test.NewTestBlock("foo", "", map[string]*test.TestBlock{}),
		),
		test.TableEntry(
			`testblock {}`,
			test.NewTestBlock("0", "", map[string]*test.TestBlock{}),
		),
		test.TableEntry(
			`testblock foo {}`,
			test.NewTestBlock("foo", "", map[string]*test.TestBlock{}),
		),
		test.TableEntry(
			`testblock foo { 
				param1 = "bar"
			}`,
			test.NewTestBlock("foo", "bar", map[string]*test.TestBlock{}),
		),
		test.TableEntry(
			`testblock { 
				param1 = "bar"
			}`,
			test.NewTestBlock("0", "bar", map[string]*test.TestBlock{}),
		),
		test.TableEntry(
			`testblock a { 
				param1 = "bar"
				testblock b {
					param1 = "foo"
				}
				testblock {
					param1 = "baz"
				}
			}`,
			test.NewTestBlock("a", "bar", map[string]*test.TestBlock{
				"b": test.NewTestBlock("b", "foo", map[string]*test.TestBlock{}),
				"0": test.NewTestBlock("0", "baz", map[string]*test.TestBlock{}),
			}),
		),
		test.TableEntry(
			`testblock { 
				a = 1
				b = 1.2
				c = "string"
				d = true
				e = nil
				f = [1, 2, "three"]
				g = [
					1,
					2,
					"three",
				]
			}`,
			test.NewTestBlock("0", "", map[string]*test.TestBlock{}),
		),
	)

	DescribeTable("it returns a parse error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveParseError(p)(input, expectedErr)
		},
		test.TableEntry(
			`testblock {
			`,
			errors.New("was expecting \"}\" at testfile:2:4"),
		),
		test.TableEntry(
			`testblock {
				a = 1
			`,
			errors.New("was expecting \"}\" at testfile:3:4"),
		),
		test.TableEntry(
			`testblock { 
				a = [
					1,
					2
				]
			}`,
			errors.New("was expecting \",\" at testfile:4:7"),
		),
		test.TableEntry(
			`testblock { param1 = "bar" }`,
			errors.New("was expecting a new line at testfile:1:13"),
		),
		test.TableEntry(
			`testblock { 
				param1 = "bar"}
			`,
			errors.New("was expecting a new line at testfile:2:19"),
		),
	)

	DescribeTable("it returns an eval error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveEvalError(p)(input, expectedErr)
		},
		test.TableEntry(
			`unknownblock {}`,
			errors.New("unknown block type at testfile:1:1"),
		),
		test.TableEntry(
			`testblock {
				testblock foo {}
				testblock foo {}
			}`,
			errors.New("duplicated id at testfile:3:15"),
		),
	)

})
