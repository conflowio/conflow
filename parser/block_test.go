package parser_test

import (
	"errors"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	"github.com/opsidian/ocl/parser"
	"github.com/opsidian/ocl/test"
)

var _ = Describe("Block", func() {

	p := parser.Block()

	DescribeTable("it evaluates the input correctly",
		func(input string, expected *test.TestBlock) {
			test.ExpectBlockToEvaluate(p)(input, expected)
		},
		test.TableEntry(
			`testblock`,
			&test.TestBlock{IDField: "0"},
		),
		test.TableEntry(
			`testblock foo`,
			&test.TestBlock{IDField: "foo"},
		),
		test.TableEntry(
			`testblock {}`,
			&test.TestBlock{IDField: "0"},
		),
		test.TableEntry(
			`testblock foo {}`,
			&test.TestBlock{IDField: "foo"},
		),
		test.TableEntry(
			`testblock {
				value = 123
			}`,
			&test.TestBlock{IDField: "0", Value: int64(123)},
		),
		test.TableEntry(
			`testblock foo {
				value = 123
			}`,
			&test.TestBlock{IDField: "foo", Value: int64(123)},
		),
		test.TableEntry(
			`testblock foo {
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
				IDField:           "foo",
				FieldString:       "a",
				FieldInt:          int64(1),
				FieldFloat:        1.2,
				FieldBool:         true,
				FieldArray:        []interface{}{1.2, "bar"},
				FieldMap:          map[string]interface{}{"a": int64(1), "b": 1.2},
				FieldTimeDuration: 1*time.Hour + 30*time.Minute,
			},
		),
		test.TableEntry(
			`testblock foo {
				custom_field = "bar"
			}`,
			&test.TestBlock{
				IDField:         "foo",
				FieldCustomName: "bar",
			},
		),
		test.TableEntry(
			`testblock a {
				value = 123
				testblock b {
					value = 234
				}
				testblock {
					value = 345
				}
			}`,
			&test.TestBlock{
				IDField: "a",
				Value:   int64(123),
				Blocks: []*test.TestBlock{
					&test.TestBlock{IDField: "b", Value: int64(234)},
					&test.TestBlock{IDField: "0", Value: int64(345)},
				},
			},
		),
		test.TableEntry(
			`testblock 123`,
			&test.TestBlock{IDField: "0", Value: int64(123)},
		),
		test.TableEntry(
			`testblock foo 123`,
			&test.TestBlock{IDField: "foo", Value: int64(123)},
		),
		// test.TableEntry(
		// 	`testblock foo 5`,
		// 	test.NewTestBlock("foo", int64(5), map[string]*test.TestBlock{}),
		// ),
		// test.TableEntry(
		// 	`testblock foo 5.6`,
		// 	test.NewTestBlock("foo", 5.6, map[string]*test.TestBlock{}),
		// ),
		// test.TableEntry(
		// 	`testblock foo true`,
		// 	test.NewTestBlock("foo", true, map[string]*test.TestBlock{}),
		// ),
		// test.TableEntry(
		// 	`testblock true`,
		// 	test.NewTestBlock("0", true, map[string]*test.TestBlock{}),
		// ),
		// test.TableEntry(
		// 	`testblock foo [1, 2]`,
		// 	test.NewTestBlock("foo", []interface{}{int64(1), int64(2)}, map[string]*test.TestBlock{}),
		// ),
		// test.TableEntry(
		// 	`testblock foo [
		// 		1,
		// 		2,
		// 	]`,
		// 	test.NewTestBlock("foo", []interface{}{int64(1), int64(2)}, map[string]*test.TestBlock{}),
		// ),
		// test.TableEntry(
		// 	`testblock foo map{
		// 		"a": "b",
		// 	}`,
		// 	test.NewTestBlock("foo", map[string]interface{}{"a": "b"}, map[string]*test.TestBlock{}),
		// ),
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
				param1 = "bar"}`,
			errors.New("was expecting a new line at testfile:2:19"),
		),
	)

	DescribeTable("it returns an eval error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveEvalError(p)(input, expectedErr)
		},
		test.TableEntry(
			`unknownblock {}`,
			errors.New("\"unknownblock\" type is invalid or not allowed here at testfile:1:1"),
		),
		// test.TableEntry(
		// 	`testblock {
		// 		testblock foo {}
		// 		testblock foo {}
		// 	}`,
		// 	errors.New("duplicated id at testfile:3:15"),
		// ),
	)

})
