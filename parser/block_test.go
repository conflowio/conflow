package parser_test

import (
	"errors"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/opsidian/basil/basil/block"
	"github.com/opsidian/basil/parser"
	"github.com/opsidian/basil/test"
)

func compareTestBlocks(b1i interface{}, b2i interface{}, input string) {
	b1 := b1i.(*test.TestBlock)
	b2 := b2i.(*test.TestBlock)

	Expect(b1.IDField).To(Equal(b2.IDField), "id does not match, input: %s", input)

	interpreter := test.TestBlockInterpreter{}

	for paramName, _ := range interpreter.Params() {
		v1 := interpreter.Param(b1, paramName)
		v2 := interpreter.Param(b2, paramName)
		if v2 != nil {
			Expect(v1).To(Equal(v2), "%s does not match, input: %s", paramName, input)
		} else {
			Expect(v1).To(BeNil(), "%s does not match, input: %s", paramName, input)
		}
	}

	Expect(len(b1.TestBlock)).To(Equal(len(b2.TestBlock)), "child block count does not match, input: %s", input)

	for _, c1 := range b1.TestBlock {
		found := false
		for _, c2 := range b2.TestBlock {
			if c1.IDField == c2.IDField {
				compareTestBlocks(c1, c2, input)
				found = true
				break
			}
		}
		if !found {
			Fail(fmt.Sprintf("was expecting child block with identifier %q", c1.IDField))
		}
	}
}

var _ = Describe("Block parser", func() {

	p := parser.Block(parser.Expression())

	var registry = block.InterpreterRegistry{
		"testblock": test.TestBlockInterpreter{},
	}

	DescribeTable("it evaluates the block correctly",
		func(input string, expected *test.TestBlock) {
			test.ExpectBlockToEvaluate(p, registry)(input, expected, compareTestBlocks)
		},
		Entry(
			"no id and body",
			`testblock`,
			&test.TestBlock{IDField: "0"},
		),
		Entry(
			"id but no body",
			`testblock foo`,
			&test.TestBlock{IDField: "foo"},
		),
		Entry(
			"no id, empty body",
			`testblock {}`,
			&test.TestBlock{IDField: "0"},
		),
		Entry(
			"id and empty body",
			`testblock foo {}`,
			&test.TestBlock{IDField: "foo"},
		),
		Entry(
			"a parameter defined",
			`testblock {
				field_string = "foo"
			}`,
			&test.TestBlock{IDField: "0", FieldString: "foo"},
		),
		Entry(
			"a custom parameter defined",
			`testblock {
				user_param := "foo"
			}`,
			&test.TestBlock{IDField: "0"},
		),
		Entry(
			"value parameter defined",
			`testblock {
				value = 123
			}`,
			&test.TestBlock{IDField: "0", Value: int64(123)},
		),
		Entry(
			"extra parameter defined",
			`testblock {
				extra_value := 123
			}`,
			&test.TestBlock{IDField: "0"},
		),
		Entry(
			"all parameter types defined",
			`testblock {
				field_string = "a"
				field_int = 1
				field_float = 1.2
				field_bool = true
				field_array = [1.2, "bar"]
				field_map = map{
					"a": 1,
					"b": 1.2,
				}
				field_time_duration = 1h30m
			}`,
			&test.TestBlock{
				IDField:           "0",
				FieldString:       "a",
				FieldInt:          int64(1),
				FieldFloat:        1.2,
				FieldBool:         true,
				FieldArray:        []interface{}{1.2, "bar"},
				FieldMap:          map[string]interface{}{"a": int64(1), "b": 1.2},
				FieldTimeDuration: 1*time.Hour + 30*time.Minute,
			},
		),
		Entry(
			"all parameter types as nil",
			`testblock {
				field_string = nil
				field_int = nil
				field_float = nil
				field_bool = nil
				field_array = nil
				field_map = nil
				field_time_duration = nil
			}`,
			&test.TestBlock{IDField: "0"},
		),
		Entry(
			"parameter with custom name defined",
			`testblock {
				custom_field = "bar"
			}`,
			&test.TestBlock{
				IDField:         "0",
				FieldCustomName: "bar",
			},
		),
		Entry(
			"short format with no id",
			`testblock "value"`,
			&test.TestBlock{IDField: "0", Value: "value"},
		),
		Entry(
			"short format with id",
			`testblock foo "value"`,
			&test.TestBlock{IDField: "foo", Value: "value"},
		),
		Entry(
			"short format with array",
			`testblock foo [1, 2]`,
			&test.TestBlock{IDField: "foo", Value: []interface{}{int64(1), int64(2)}},
		),
		Entry(
			"short format with multiline array",
			`testblock foo [
				1,
				2,
			]`,
			&test.TestBlock{IDField: "foo", Value: []interface{}{int64(1), int64(2)}},
		),
		Entry(
			"short format with map",
			`testblock foo map{
				"a": "b",
			}`,
			&test.TestBlock{IDField: "foo", Value: map[string]interface{}{"a": "b"}},
		),
		Entry(
			"function call in parameter value",
			`testblock {
				field_string = test.func2("foo", "bar")
			}`,
			&test.TestBlock{IDField: "0", FieldString: "foobar"},
		),
	)

	DescribeTable("it evaluates a child block correctly",
		func(input string, expected *test.TestBlock) {
			test.ExpectBlockToEvaluate(p, registry)(input, expected, compareTestBlocks)
		},
		Entry(
			"child block without id and body",
			`testblock {
				testblock
			}`,
			&test.TestBlock{
				IDField: "0",
				TestBlock: []*test.TestBlock{
					{IDField: "1"},
				},
			},
		),
		Entry(
			"child block with id but no body",
			`testblock {
				testblock foo
			}`,
			&test.TestBlock{
				IDField: "0",
				TestBlock: []*test.TestBlock{
					{IDField: "foo"},
				},
			},
		),
		Entry(
			"child block with empty body",
			`testblock {
				testblock {
				}
			}`,
			&test.TestBlock{
				IDField: "0",
				TestBlock: []*test.TestBlock{
					{IDField: "1"},
				},
			},
		),
		Entry(
			"child block with parameter defined",
			`testblock {
				testblock {
					value = 1
				}
			}`,
			&test.TestBlock{
				IDField: "0",
				TestBlock: []*test.TestBlock{
					{IDField: "1", Value: int64(1)},
				},
			},
		),
		Entry(
			"multiple child blocks of the same type",
			`testblock {
				testblock {
					value = 1
				}
				testblock {
					value = 2
				}
			}`,
			&test.TestBlock{
				IDField: "0",
				TestBlock: []*test.TestBlock{
					{IDField: "1", Value: int64(1)},
					{IDField: "2", Value: int64(2)},
				},
			},
		),
		Entry(
			"variable as parameter value",
			`testblock {
				testblock b1 {
					value = 1
				}
				testblock b2 {
					value = b1.value + 1
				}
			}`,
			&test.TestBlock{
				IDField: "0",
				TestBlock: []*test.TestBlock{
					{IDField: "b1", Value: int64(1)},
					{IDField: "b2", Value: int64(2)},
				},
			},
		),
	)

	DescribeTable("it returns a parse error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveParseError(p)(input, expectedErr)
		},
		Entry(
			"no closing }",
			`testblock {
			`,
			errors.New("was expecting \"}\" at testfile:2:4"),
		),
		Entry(
			"no closing } with paramter",
			`testblock {
				a = 1
			`,
			errors.New("was expecting \"}\" at testfile:3:4"),
		),
		Entry(
			"missing , in multiline array",
			`testblock {
				a = [
					1,
					2
				]
			}`,
			errors.New("was expecting \",\" at testfile:4:7"),
		),
		Entry(
			"defining a block body in one line",
			`testblock { param1 = "bar" }`,
			errors.New("was expecting a new line at testfile:1:13"),
		),
		Entry(
			"block body closing } not in new line",
			`testblock {
				param1 = "bar"}`,
			errors.New("was expecting a new line at testfile:2:19"),
		),
	)

	DescribeTable("it returns an parse error",
		func(input string, expectedErr error) {
			test.ExpectBlockToHaveParseError(p, registry)(input, MatchError(expectedErr))
		},
		Entry(
			"unknown block type",
			`unknownblock {}`,
			errors.New("\"unknownblock\" type is invalid or not allowed here at testfile:1:1"),
		),
		Entry(
			"unknown child block type",
			`testblock {
				unknownblock {}
			}`,
			errors.New("\"unknownblock\" type is invalid or not allowed here at testfile:2:5"),
		),
		Entry(
			"same parameter defined multiple times",
			`testblock {
				param1 = "bar"
				param1 = "foo"
			}`,
			errors.New("\"param1\" parameter was defined multiple times at testfile:3:5"),
		),
		Entry(
			"wrong parameter value type",
			`testblock {
				field_string = 1
			}`,
			errors.New("was expecting string at testfile:2:5"),
		),
		Entry(
			"extra parameter using existing parameter name",
			`testblock {
				field_string := "foo"
			}`,
			errors.New("\"field_string\" parameter already exists. Use \"=\" to set the parameter value or use a different name at testfile:2:5"),
		),
		Entry(
			"unknown parameter",
			`testblock {
				unknown = "foo"
			}`,
			errors.New("\"unknown\" parameter does not exist at testfile:2:5"),
		),
		Entry(
			"duplicated identifiers",
			`testblock {
				testblock foo {}
				testblock foo {}
			}`,
			errors.New("\"foo\" is already defined, please use a globally unique identifier at testfile:3:15"),
		),
	)

})
