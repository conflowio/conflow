// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsers_test

import (
	"errors"

	"github.com/opsidian/basil/basil/variable"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	"github.com/opsidian/basil/parsers"
	"github.com/opsidian/basil/test"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/text/terminal"
)

var _ = Describe("Map", func() {

	q := combinator.Choice(
		terminal.String(false),
		terminal.Integer(),
		terminal.Nil("nil", variable.TypeNil),
		test.EvalErrorParser("ERR", variable.TypeUnknown),
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
