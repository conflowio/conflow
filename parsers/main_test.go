// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsers_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"github.com/opsidian/basil/basil/block"
	"github.com/opsidian/basil/parsers"
	"github.com/opsidian/basil/test"
)

var _ = Describe("Main block parser", func() {

	p := parsers.NewMain("main", test.BlockInterpreter{})

	var registry = block.InterpreterRegistry{
		"main":      test.BlockInterpreter{},
		"testblock": test.BlockInterpreter{},
	}

	DescribeTable("it evaluates the main block correctly",
		func(input string, expected *test.Block) {
			test.ExpectBlockToEvaluate(p, registry)(input, expected, func(i interface{}, i2 interface{}, s string) {
				i.(*test.Block).Compare(i2.(*test.Block), input)
			})
		},
		Entry(
			"empty",
			``,
			&test.Block{IDField: "main"},
		),
		Entry(
			"a parameter defined",
			`value = 1`,
			&test.Block{IDField: "main", Value: int64(1)},
		),
		Entry(
			"whitespaces and newlines",
			`
				value = 1

				field_string = "foo"
			`,
			&test.Block{IDField: "main", Value: int64(1), FieldString: "foo"},
		),
		Entry(
			"a block defined",
			`
				c testblock {
					value = 1
				}
			`,
			&test.Block{
				IDField: "main",
				Blocks: []*test.Block{
					{IDField: "c", Value: int64(1)},
				},
			},
		),
		Entry(
			"a block and parameter defined",
			`
				value = 1
				c testblock {
					value = 2
				}
			`,
			&test.Block{
				IDField: "main",
				Value:   int64(1),
				Blocks: []*test.Block{
					{IDField: "c", Value: int64(2)},
				},
			},
		),
	)

	DescribeTable("it evaluates dependencies correctly",
		func(input string, expected *test.Block) {
			test.ExpectBlockToEvaluate(p, registry)(input, expected, func(i interface{}, i2 interface{}, s string) {
				i.(*test.Block).Compare(i2.(*test.Block), input)
			})
		},
		Entry(
			"a parameter depends on an other",
			`
				field_int = main.value + 1
				value = 1
			`,
			&test.Block{
				IDField:  "main",
				Value:    int64(1),
				FieldInt: int64(2),
			},
		),
		Entry(
			"a parameter depends on a custom parameter",
			`
				value = main.user_param + 1
				user_param := 1
			`,
			&test.Block{
				IDField: "main",
				Value:   int64(2),
			},
		),
		Entry(
			"a child block depends on the other",
			`
				c1 testblock {
					value = c2.value + 1
				}
				c2 testblock {
					value = 1
				}
			`,
			&test.Block{
				IDField: "main",
				Blocks: []*test.Block{
					{IDField: "c2", Value: int64(1)},
					{IDField: "c1", Value: int64(2)},
				},
			},
		),
		Entry(
			"a parameter depends on a child block",
			`
				value = c.value + 1
				c testblock {
					value = 1
				}
			`,
			&test.Block{
				IDField: "main",
				Value:   int64(2),
				Blocks: []*test.Block{
					{IDField: "c", Value: int64(1)},
				},
			},
		),
		Entry(
			"a parameter depends on a child block's custom parameter",
			`
				value = c.user_param + 1
				c testblock {
					user_param := 1
				}
			`,
			&test.Block{
				IDField: "main",
				Value:   int64(2),
				Blocks: []*test.Block{
					{IDField: "c"},
				},
			},
		),
		Entry(
			"a child block depends on a parameter",
			`
				c testblock {
					value = main.value + 1
				}
				value = 1
			`,
			&test.Block{
				IDField: "main",
				Value:   int64(1),
				Blocks: []*test.Block{
					{IDField: "c", Value: int64(2)},
				},
			},
		),
		Entry(
			"a child block is referenced on a deeper level",
			`
				c1 testblock {
					c2 testblock {
						value = 1
					}
				}
				value = c2.value + 1
			`,
			&test.Block{
				IDField: "main",
				Value:   int64(2),
				Blocks: []*test.Block{
					{
						IDField: "c1",
						Blocks: []*test.Block{
							{
								IDField: "c2",
								Value:   int64(1),
							},
						},
					},
				},
			},
		),
		Entry(
			"a deeper level references a sibling deeper level",
			`
				c1 testblock {
					c2 testblock {
						value = c4.value + 1
					}
				}
				c3 testblock {
					c4 testblock {
						value = 1
					}
				}
			`,
			&test.Block{
				IDField: "main",
				Blocks: []*test.Block{
					{
						IDField: "c3",
						Blocks: []*test.Block{
							{
								IDField: "c4",
								Value:   int64(1),
							},
						},
					},
					{
						IDField: "c1",
						Blocks: []*test.Block{
							{
								IDField: "c2",
								Value:   int64(2),
							},
						},
					},
				},
			},
		),
		Entry(
			"dependencies are resolved recursively",
			`
				value = c1.value + c2.value
				c1 testblock {
					value = c1.field_int + 1
					field_int = 1
				}
				c2 testblock {
					value = c2.field_int + 2
					field_int = 2
				}
			`,
			&test.Block{
				IDField: "main",
				Value:   int64(6),
				Blocks: []*test.Block{
					{IDField: "c1", Value: int64(2), FieldInt: int64(1)},
					{IDField: "c2", Value: int64(4), FieldInt: int64(2)},
				},
			},
		),
	)

	DescribeTable("it errors on circular dependencies",
		func(input string, errMatcher types.GomegaMatcher) {
			test.ExpectBlockToHaveParseError(p, registry)(input, errMatcher)
		},
		Entry(
			"a parameter depends on itself",
			`
				value = main.value + 1
			`,
			MatchError(errors.New("main.value should not reference itself at testfile:2:13")),
		),
		Entry(
			"parameters depend on each other",
			`
				value = main.field_int + 1
				field_int = main.value + 1
			`,
			Or(
				MatchError(errors.New("circular dependency detected: main.field_int, main.value at testfile:3:5")),
				MatchError(errors.New("circular dependency detected: main.value, main.field_int at testfile:2:5")),
			),
		),
		Entry(
			"a parameter and a block depend on each other",
			`
				value = c.value
				c testblock {
					value = main.value
				}
			`,
			Or(
				MatchError(errors.New("circular dependency detected: c, main.value at testfile:3:5")),
				MatchError(errors.New("circular dependency detected: main.value, c at testfile:2:5")),
			),
		),
		Entry(
			"deeper levels depend on each other",
			`
				c1 testblock {
					c2 testblock {
						value = c4.value
					}
				}
				c3 testblock {
					c4 testblock {
						value = c2.value
					}
				}
			`,
			Or(
				MatchError(errors.New("circular dependency detected: c3, c1 at testfile:7:5")),
				MatchError(errors.New("circular dependency detected: c1, c3 at testfile:2:5")),
			),
		),
	)

	DescribeTable("it returns a parse error",
		func(input string, expectedErr error) {
			test.ExpectBlockToHaveParseError(p, registry)(input, MatchError(expectedErr))
		},
		Entry(
			"unknown block reference",
			`value = nonexisting.value`,
			errors.New("unknown block: \"nonexisting\" at testfile:1:9"),
		),
		Entry(
			"unknown parameter reference",
			`value = main.nonexisting`,
			errors.New("unknown parameter: \"main.nonexisting\" at testfile:1:9"),
		),
	)

})
