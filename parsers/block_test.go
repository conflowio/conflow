// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsers_test

import (
	"errors"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/opsidian/basil/basil/block"
	"github.com/opsidian/basil/parsers"
	"github.com/opsidian/basil/test"
)

var _ = Describe("Block parser", func() {

	p := parsers.Block(parsers.Expression())

	var registry = block.InterpreterRegistry{
		"testblock": test.BlockInterpreter{},
	}

	DescribeTable("it evaluates the block correctly",
		func(input string, expected *test.Block) {
			test.ExpectBlockToEvaluate(p, registry)(input, expected, func(i interface{}, i2 interface{}, s string) {
				i.(*test.Block).Compare(i2.(*test.Block), input)
			})
		},
		Entry(
			"no id and body",
			`testblock`,
			&test.Block{IDField: "0"},
		),
		Entry(
			"id but no body",
			`testblock foo`,
			&test.Block{IDField: "foo"},
		),
		Entry(
			"no id, empty body",
			`testblock {}`,
			&test.Block{IDField: "0"},
		),
		Entry(
			"id and empty body",
			`testblock foo {}`,
			&test.Block{IDField: "foo"},
		),
		Entry(
			"a parameter defined",
			`testblock {
				field_string = "foo"
			}`,
			&test.Block{IDField: "0", FieldString: "foo"},
		),
		Entry(
			"a custom parameter defined",
			`testblock {
				user_param := "foo"
			}`,
			&test.Block{IDField: "0"},
		),
		Entry(
			"value parameter defined",
			`testblock {
				value = 123
			}`,
			&test.Block{IDField: "0", Value: int64(123)},
		),
		Entry(
			"extra parameter defined",
			`testblock {
				extra_value := 123
			}`,
			&test.Block{IDField: "0"},
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
			&test.Block{
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
			&test.Block{IDField: "0"},
		),
		Entry(
			"parameter with custom name defined",
			`testblock {
				custom_field = "bar"
			}`,
			&test.Block{
				IDField:         "0",
				FieldCustomName: "bar",
			},
		),
		Entry(
			"short format with no id",
			`testblock "value"`,
			&test.Block{IDField: "0", Value: "value"},
		),
		Entry(
			"short format with id",
			`testblock foo "value"`,
			&test.Block{IDField: "foo", Value: "value"},
		),
		Entry(
			"short format with array",
			`testblock foo [1, 2]`,
			&test.Block{IDField: "foo", Value: []interface{}{int64(1), int64(2)}},
		),
		Entry(
			"short format with multiline array",
			`testblock foo [
				1,
				2,
			]`,
			&test.Block{IDField: "foo", Value: []interface{}{int64(1), int64(2)}},
		),
		Entry(
			"short format with map",
			`testblock foo map{
				"a": "b",
			}`,
			&test.Block{IDField: "foo", Value: map[string]interface{}{"a": "b"}},
		),
		Entry(
			"function call in parameter value",
			`testblock {
				field_string = test.func2("foo", "bar")
			}`,
			&test.Block{IDField: "0", FieldString: "foobar"},
		),
	)

	DescribeTable("it evaluates a child block correctly",
		func(input string, expected *test.Block) {
			test.ExpectBlockToEvaluate(p, registry)(input, expected, func(i interface{}, i2 interface{}, input string) {
				i.(*test.Block).Compare(i2.(*test.Block), input)
			})
		},
		Entry(
			"child block without id and body",
			`testblock {
				testblock
			}`,
			&test.Block{
				IDField: "0",
				Blocks: []*test.Block{
					{IDField: "1"},
				},
			},
		),
		Entry(
			"child block with id but no body",
			`testblock {
				testblock foo
			}`,
			&test.Block{
				IDField: "0",
				Blocks: []*test.Block{
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
			&test.Block{
				IDField: "0",
				Blocks: []*test.Block{
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
			&test.Block{
				IDField: "0",
				Blocks: []*test.Block{
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
			&test.Block{
				IDField: "0",
				Blocks: []*test.Block{
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
			&test.Block{
				IDField: "0",
				Blocks: []*test.Block{
					{IDField: "b1", Value: int64(1)},
					{IDField: "b2", Value: int64(2)},
				},
			},
		),
	)

	DescribeTable("it evaluates directives correctly",
		func(input string, expected *test.Block) {
			test.ExpectBlockToEvaluate(p, registry)(input, expected, func(i interface{}, i2 interface{}, input string) {
				i.(*test.Block).Compare(i2.(*test.Block), input)
			})
		},
		Entry(
			"a directive with no body",
			`@testdirective
				testblock {}`,
			&test.Block{IDField: "1"},
		),
		Entry(
			"a directive with a value",
			`@testdirective "foo"
				testblock {}`,
			&test.Block{IDField: "1"},
		),
		Entry(
			"a directive with an empty body",
			`@testdirective {}
				testblock {}`,
			&test.Block{IDField: "1"},
		),
		Entry(
			"a directive with a body",
			`@testdirective {
				value = "foo"
			}
			testblock {}`,
			&test.Block{IDField: "1"},
		),
		Entry(
			"multiple directives",
			`@testdirective "foo"
			@testdirective2 "bar"
			testblock {}`,
			&test.Block{IDField: "2"},
		),
		Entry(
			"a parameter directive with no body",
			`testblock {
				@testdirective
				value = "bar"
			}`,
			&test.Block{IDField: "0", Value: "bar"},
		),
		Entry(
			"a parameter directive with an value",
			`testblock {
				@testdirective "foo"
				value = "bar"
			}`,
			&test.Block{IDField: "0", Value: "bar"},
		),
		Entry(
			"a parameter directive with an empty body",
			`testblock {
				@testdirective {}
				value = "bar"
			}`,
			&test.Block{IDField: "0", Value: "bar"},
		),
		Entry(
			"a parameter directive with a body",
			`testblock {
				@testdirective {
				  value = "foo"
				}
				value = "bar"
			}`,
			&test.Block{IDField: "0", Value: "bar"},
		),
		Entry(
			"a parameter with multiple directives",
			`testblock {
				@testdirective "foo"
				@testdirective2 "foo2"
				value = "bar"
			}`,
			&test.Block{IDField: "0", Value: "bar"},
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
			errors.New("\"unknownblock\" block is unknown or not allowed at testfile:1:1"),
		),
		Entry(
			"unknown child block type",
			`testblock {
				unknownblock {}
			}`,
			errors.New("\"unknownblock\" block is unknown or not allowed at testfile:2:5"),
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
			errors.New("was expecting string at testfile:2:20"),
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
