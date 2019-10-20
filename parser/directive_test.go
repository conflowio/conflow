// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parser_test

import (
	"errors"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/opsidian/basil/parser"
	"github.com/opsidian/basil/test"
)

var _ = Describe("Directive parser", func() {

	p := parser.Directive(parser.Expression())

	DescribeTable("it evaluates the block correctly",
		func(input string, expected *test.Directive) {
			test.ExpectBlockToEvaluate(p, nil)(input, expected, func(i interface{}, i2 interface{}, s string) {
				i.(*test.Directive).Compare(i2.(*test.Directive), input)
			})
		},
		Entry(
			"no body",
			`@testdirective`,
			&test.Directive{IDField: "0"},
		),
		Entry(
			"empty body",
			`@testdirective {}`,
			&test.Directive{IDField: "0"},
		),
		Entry(
			"a parameter defined",
			`@testdirective {
				field_string = "foo"
			}`,
			&test.Directive{IDField: "0", FieldString: "foo"},
		),
		Entry(
			"value parameter defined",
			`@testdirective {
				value = 123
			}`,
			&test.Directive{IDField: "0", Value: int64(123)},
		),
		Entry(
			"all parameter types defined",
			`@testdirective {
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
			&test.Directive{
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
			`@testdirective {
				field_string = nil
				field_int = nil
				field_float = nil
				field_bool = nil
				field_array = nil
				field_map = nil
				field_time_duration = nil
			}`,
			&test.Directive{IDField: "0"},
		),
		Entry(
			"short format",
			`@testdirective "value"`,
			&test.Directive{IDField: "0", Value: "value"},
		),
		Entry(
			"short format with array",
			`@testdirective [1, 2]`,
			&test.Directive{IDField: "0", Value: []interface{}{int64(1), int64(2)}},
		),
		Entry(
			"short format with multiline array",
			`@testdirective [
				1,
				2,
			]`,
			&test.Directive{IDField: "0", Value: []interface{}{int64(1), int64(2)}},
		),
		Entry(
			"short format with map",
			`@testdirective map{
				"a": "b",
			}`,
			&test.Directive{IDField: "0", Value: map[string]interface{}{"a": "b"}},
		),
		Entry(
			"function call in parameter value",
			`@testdirective {
				field_string = test.func2("foo", "bar")
			}`,
			&test.Directive{IDField: "0", FieldString: "foobar"},
		),
	)

	DescribeTable("it evaluates a child block correctly",
		func(input string, expected *test.Directive) {
			test.ExpectBlockToEvaluate(p, nil)(input, expected, func(i interface{}, i2 interface{}, s string) {
				i.(*test.Directive).Compare(i2.(*test.Directive), input)
			})
		},
		Entry(
			"child block without body",
			`@testdirective {
				testblock
			}`,
			&test.Directive{
				IDField: "0",
				Blocks: []*test.Block{
					{IDField: "1"},
				},
			},
		),
		Entry(
			"child block with empty body",
			`@testdirective {
				testblock {
				}
			}`,
			&test.Directive{
				IDField: "0",
				Blocks: []*test.Block{
					{IDField: "1"},
				},
			},
		),
		Entry(
			"child block with parameter defined",
			`@testdirective {
				testblock {
					value = 1
				}
			}`,
			&test.Directive{
				IDField: "0",
				Blocks: []*test.Block{
					{IDField: "1", Value: int64(1)},
				},
			},
		),
		Entry(
			"multiple child blocks of the same type",
			`@testdirective {
				testblock {
					value = 1
				}
				testblock {
					value = 2
				}
			}`,
			&test.Directive{
				IDField: "0",
				Blocks: []*test.Block{
					{IDField: "1", Value: int64(1)},
					{IDField: "2", Value: int64(2)},
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
			`@testdirective {
			`,
			errors.New("was expecting \"}\" at testfile:2:4"),
		),
		Entry(
			"no closing } with paramter",
			`@testdirective {
				a = 1
			`,
			errors.New("was expecting \"}\" at testfile:3:4"),
		),
		Entry(
			"missing , in multiline array",
			`@testdirective {
				a = [
					1,
					2
				]
			}`,
			errors.New("was expecting \",\" at testfile:4:7"),
		),
		Entry(
			"defining a block body in one line",
			`@testdirective { param1 = "bar" }`,
			errors.New("was expecting a new line at testfile:1:18"),
		),
		Entry(
			"block body closing } not in new line",
			`@testdirective {
				param1 = "bar"}`,
			errors.New("was expecting a new line at testfile:2:19"),
		),
		Entry(
			"defining an id",
			`@testdirective foo {}`,
			errors.New("was expecting the end of input at testfile:1:16"),
		),
		Entry(
			"defining a custom parameter",
			`@testdirective {
				extra_param := "foo"
			}`,
			errors.New("was expecting a new line at testfile:2:17"),
		),
		Entry(
			"defining a directive as a child block",
			`@testdirective {
				@testdirective
			}`,
			errors.New("was expecting \"}\" at testfile:2:5"),
		),
	)

	DescribeTable("it returns an static check error",
		func(input string, expectedErr error) {
			test.ExpectBlockToHaveParseError(p, nil)(input, MatchError(expectedErr))
		},
		Entry(
			"unknown block type",
			`@unknowndirective {}`,
			errors.New("\"@unknowndirective\" directive is unknown or not allowed at testfile:1:1"),
		),
		Entry(
			"unknown child block type",
			`@testdirective {
				unknownblock {}
			}`,
			errors.New("\"unknownblock\" block is unknown or not allowed at testfile:2:5"),
		),
		Entry(
			"same parameter defined multiple times",
			`@testdirective {
				param1 = "bar"
				param1 = "foo"
			}`,
			errors.New("\"param1\" parameter was defined multiple times at testfile:3:5"),
		),
		Entry(
			"wrong parameter value type",
			`@testdirective {
				field_string = 1
			}`,
			errors.New("was expecting string at testfile:2:5"),
		),
		Entry(
			"unknown parameter",
			`@testdirective {
				unknown = "foo"
			}`,
			errors.New("\"unknown\" parameter does not exist at testfile:2:5"),
		),
	)

})
