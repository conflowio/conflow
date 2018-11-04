package parser_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	"github.com/opsidian/basil/parser"
	"github.com/opsidian/basil/test"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/text/terminal"
)

var _ = Describe("Variable", func() {

	q := combinator.Choice(
		terminal.String(false),
		terminal.Integer(),
		terminal.Nil("nil"),
		test.EvalErrorParser(),
	).Name("value")

	p := parser.Variable(q)

	DescribeTable("it evaluates the input correctly",
		func(input string, expected interface{}) {
			test.ExpectParserToEvaluate(p)(input, expected)
		},
		test.TableEntry(`foo`, "bar"),
		test.TableEntry(`testmap.key1`, "value1"),
		test.TableEntry(`testmap["key1"]`, "value1"),
		test.TableEntry(`testmap.key2.key3`, "value3"),
		test.TableEntry(`testmap["key2"].key3`, "value3"),
		test.TableEntry(`testmap.key2["key3"]`, "value3"),
		test.TableEntry(`testmap.key4[0]`, "value4"),
		test.TableEntry(`arr[0]`, "value1"),
		test.TableEntry(`arr[1][0]`, "value2"),
		test.TableEntry(`arr[2].key1`, "value3"),
		test.TableEntry(`arr[2]["key1"]`, "value3"),
	)

	DescribeTable("it returns a parse error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveParseError(p)(input, expectedErr)
		},
		test.TableEntry(`testmap.`, errors.New("was expecting identifier at testfile:1:9")),
		test.TableEntry(`testmap.key1.`, errors.New("was expecting identifier at testfile:1:14")),
		test.TableEntry(`testmap[`, errors.New("was expecting value at testfile:1:9")),
		test.TableEntry(`testmap["key1"`, errors.New("was expecting \"]\" at testfile:1:15")),
		test.TableEntry(`testmap[]`, errors.New("was expecting value at testfile:1:9")),
		test.TableEntry(`arr.0`, errors.New("was expecting identifier at testfile:1:5")),
		test.TableEntry(`arr[1].0`, errors.New("was expecting identifier at testfile:1:8")),
		test.TableEntry(`testkeyword`, errors.New("testkeyword is a reserved keyword at testfile:1:1")),
		test.TableEntry(`testkeyword.foo`, errors.New("testkeyword is a reserved keyword at testfile:1:1")),
	)

	DescribeTable("it returns an eval error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveEvalError(p)(input, expectedErr)
		},
		test.TableEntry(`nonexisting`, errors.New("variable 'nonexisting' does not exist at testfile:1:1")),
		test.TableEntry(`nonexisting.key`, errors.New("variable 'nonexisting[key]' does not exist at testfile:1:1")),
		test.TableEntry(`testmap.nonexisting`, errors.New("variable 'testmap[nonexisting]' does not exist at testfile:1:1")),
		test.TableEntry(`arr[3]`, errors.New("array index out of bounds: 3 (0..2) at testfile:1:5")),
		test.TableEntry(`arr[-1]`, errors.New("array index out of bounds: -1 (0..2) at testfile:1:5")),
		test.TableEntry(`testmap.key1.key2`, errors.New("can not get index on string type at testfile:1:14")),
		test.TableEntry(`testmap.key2[0]`, errors.New("invalid non-string index on map at testfile:1:14")),
		test.TableEntry(`arr["string"]`, errors.New("invalid non-integer index on array at testfile:1:5")),
		test.TableEntry(`arr[ERR]`, errors.New("ERR at testfile:1:5")),
		test.TableEntry(`testmap[ERR]`, errors.New("ERR at testfile:1:9")),
	)

})
