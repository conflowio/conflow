// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsers_test

import (
	"errors"

	"github.com/conflowio/parsley/combinator"
	"github.com/conflowio/parsley/text/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"

	"github.com/conflowio/conflow/src/parsers"
	"github.com/conflowio/conflow/src/schema"
	"github.com/conflowio/conflow/src/test"
)

var _ = Describe("Map", func() {

	q := combinator.Choice(
		terminal.String(schema.StringValue(), false),
		terminal.Integer(schema.IntegerValue()),
		terminal.Nil(schema.NullValue(), "NULL"),
		test.EvalErrorParser(schema.UntypedValue(), "ERR"),
	).Name("value")

	p := parsers.Map(q)

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
			`map { "a": "b" }`,
			map[string]interface{}{"a": "b"},
		),
		test.TableEntry(
			`map { "a": "b", "c": "d" }`,
			map[string]interface{}{"a": "b", "c": "d"},
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
				"foo": "bar"
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
				"c": "d",
			}`,
			map[string]interface{}{
				"a": "b",
				"c": "d",
			},
		),
		test.TableEntry(
			`map{
				"a": NULL,
			}`,
			map[string]interface{}{
				"a": nil,
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
			`map
			{ 
				"a": "b",
			}`,
			errors.New("new line is not allowed at testfile:1:4"),
		),

		test.TableEntry(
			`{ 
				"a": "b",
			}`,
			errors.New("was expecting map at testfile:1:1"),
		),
	)

	DescribeTable("it returns a static check error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveStaticCheckError(p)(input, expectedErr)
		},
		test.TableEntry(
			`map{
				"a": "b",
				"c": 1,
			}`,
			errors.New("items must have the same type, but found string and integer at testfile:1:1"),
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
