package parser_test

import (
	"errors"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/opsidian/basil/block"
	"github.com/opsidian/basil/parser"
	"github.com/opsidian/basil/test"
)

func compareTestBlocks(b1i interface{}, b2i interface{}, input string) {
	b1 := b1i.(*test.TestBlock)
	b2 := b2i.(*test.TestBlock)
	Expect(b1.IDField).To(Equal(b2.IDField), "IDField does not match, input: %s", input)
	if b2.Value != nil {
		Expect(b1.Value).To(Equal(b2.Value), "Value does not match, input: %s", input)
	} else {
		Expect(b1.Value).To(BeNil(), "Value does not match, input: %s", input)
	}
	Expect(b1.FieldString).To(Equal(b2.FieldString), "FieldString does not match, input: %s", input)
	Expect(b1.FieldInt).To(Equal(b2.FieldInt), "FieldInt does not match, input: %s", input)
	Expect(b1.FieldFloat).To(Equal(b2.FieldFloat), "FieldFloat does not match, input: %s", input)
	Expect(b1.FieldBool).To(Equal(b2.FieldBool), "FieldBool does not match, input: %s", input)
	Expect(b1.FieldArray).To(Equal(b2.FieldArray), "FieldArray does not match, input: %s", input)
	Expect(b1.FieldMap).To(Equal(b2.FieldMap), "FieldMap does not match, input: %s", input)
	Expect(b1.FieldTimeDuration).To(Equal(b2.FieldTimeDuration), "FieldTimeDuration does not match, input: %s", input)
	Expect(b1.FieldCustomName).To(Equal(b2.FieldCustomName), "FieldCustomName does not match, input: %s", input)

	Expect(len(b1.Blocks)).To(Equal(len(b2.Blocks)), "child block count does not match, input: %s", input)

	for i, block := range b1.Blocks {
		compareTestBlocks(block, b2.Blocks[i], input)
	}
}

var _ = Describe("Block parser", func() {

	p := parser.Block()

	var registry = block.Registry{
		"testblock": test.TestBlockInterpreter{},
	}

	DescribeTable("it evaluates the input correctly",
		func(input string, expected *test.TestBlock) {
			test.ExpectBlockToEvaluate(p, registry)(input, expected, compareTestBlocks)
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
				field_string = nil
			}`,
			&test.TestBlock{IDField: "foo"},
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
				testblock {
					value = 456
				}
			}`,
			&test.TestBlock{
				IDField: "a",
				Value:   int64(123),
				Blocks: []*test.TestBlock{
					&test.TestBlock{IDField: "b", Value: int64(234)},
					&test.TestBlock{IDField: "0", Value: int64(345)},
					&test.TestBlock{IDField: "1", Value: int64(456)},
				},
			},
		),
		test.TableEntry(
			`testblock 123`,
			&test.TestBlock{IDField: "0", Value: int64(123)},
		),
		test.TableEntry(
			`testblock foo "bar"`,
			&test.TestBlock{IDField: "foo", Value: "bar"},
		),
		test.TableEntry(
			`testblock foo 5`,
			&test.TestBlock{IDField: "foo", Value: int64(5)},
		),
		test.TableEntry(
			`testblock foo 5.6`,
			&test.TestBlock{IDField: "foo", Value: 5.6},
		),
		test.TableEntry(
			`testblock foo true`,
			&test.TestBlock{IDField: "foo", Value: true},
		),
		test.TableEntry(
			`testblock true`,
			&test.TestBlock{IDField: "0", Value: true},
		),
		test.TableEntry(
			`testblock foo [1, 2]`,
			&test.TestBlock{IDField: "foo", Value: []interface{}{int64(1), int64(2)}},
		),
		test.TableEntry(
			`testblock foo [
				1,
				2,
			]`,
			&test.TestBlock{IDField: "foo", Value: []interface{}{int64(1), int64(2)}},
		),
		test.TableEntry(
			`testblock foo map{
				"a": "b",
			}`,
			&test.TestBlock{IDField: "foo", Value: map[string]interface{}{"a": "b"}},
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
				param1 = "bar"}`,
			errors.New("was expecting a new line at testfile:2:19"),
		),
	)

	DescribeTable("it returns an parse error",
		func(input string, expectedErr error) {
			test.ExpectBlockToHaveParseError(p, registry)(input, expectedErr)
		},
		test.TableEntry(
			`unknownblock {}`,
			errors.New("\"unknownblock\" type is invalid or not allowed here at testfile:1:1"),
		),
		test.TableEntry(
			`testblock {
				unknownblock {}
			}`,
			errors.New("\"unknownblock\" type is invalid or not allowed here at testfile:2:5"),
		),
		test.TableEntry(
			`testblock {
				param1 = "bar"
				param1 = "foo"
			}`,
			errors.New("\"param1\" parameter was defined multiple times at testfile:3:5"),
		),
	)

	DescribeTable("block returns an eval error",
		func(input string, expectedErr error) {
			test.ExpectBlockToHaveCheckError(p, registry)(input, expectedErr)
		},
		test.TableEntry(
			`testblock {
				unknown = "foo"
			}`,
			errors.New("\"unknown\" parameter does not exist at testfile:2:5"),
		),
		test.TableEntry(
			`testblock {
				testblock foo {}
				testblock foo {}
			}`,
			errors.New("\"testblock.foo\" was defined multiple times at testfile:3:5"),
		),
	)

})
