package parser_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	"github.com/opsidian/ocl/parser"
	"github.com/opsidian/ocl/test"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/text/terminal"
)

var _ = Describe("Map", func() {

	q := combinator.Choice(
		terminal.String(false),
		terminal.Integer(),
		terminal.Word("nil", nil),
		test.EvalErrorParser(),
	).Name("value")

	p := parser.Map(q)

	DescribeTable("it evaluates the input correctly",
		func(input string, expected interface{}) {
			test.ExpectParserToEvaluate(p)(input, expected)
		},
		test.TableEntry(
			`map{}`,
			map[string]interface{}{},
		),
		test.TableEntry(
			`map {}`,
			map[string]interface{}{},
		),
		test.TableEntry(
			`map{
			}`,
			map[string]interface{}{},
		),
		test.TableEntry(
			`map{
				"foo": "bar",
			}`,
			map[string]interface{}{
				"foo": "bar",
			},
		),
		test.TableEntry(
			`map {
				"foo": "bar",
			}`,
			map[string]interface{}{
				"foo": "bar",
			},
		),
		test.TableEntry(
			`map{
				"a": "b",
				"c": 4,
			}`,
			map[string]interface{}{
				"a": "b",
				"c": int64(4),
			},
		),
		test.TableEntry(
			`map{
				"nil": nil,
			}`,
			map[string]interface{}{
				"nil": nil,
			},
		),
	)

	DescribeTable("it returns a parse error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveParseError(p)(input, expectedErr)
		},
		test.TableEntry(
			`map{
			`,
			errors.New("was expecting \"}\" at testfile:2:4"),
		),
		test.TableEntry(
			`map{
				"a": 1,
			`,
			errors.New("was expecting \"}\" at testfile:3:4"),
		),
		test.TableEntry(
			`map{ 
				"a": 1
			}`,
			errors.New("was expecting \",\" at testfile:2:11"),
		),
		test.TableEntry(
			`map{ "a": "b" }`,
			errors.New("was expecting a new line at testfile:1:6"),
		),
		test.TableEntry(
			`map{ 
				"a": "b",}
			`,
			errors.New("was expecting a new line at testfile:2:14"),
		),

		test.TableEntry(
			`map
			{ 
				"a": "b",
			}`,
			errors.New("was expecting \"{\" at testfile:1:4"),
		),

		test.TableEntry(
			`{ 
				"a": "b",
			}`,
			errors.New("was expecting map at testfile:1:1"),
		),
	)

	DescribeTable("it returns an eval error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveEvalError(p)(input, expectedErr)
		},
		test.TableEntry(
			`map{
				"a": ERR,
			}`,
			errors.New("ERR at testfile:2:10"),
		),
	)

})
