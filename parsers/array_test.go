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
	"github.com/opsidian/basil/parsers"
	"github.com/opsidian/basil/test"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/text/terminal"
)

var _ = Describe("Array", func() {

	q := combinator.Choice(
		terminal.String(false),
		terminal.Integer(),
		terminal.Nil("nil"),
		test.EvalErrorParser(),
	).Name("value")

	Describe("when new lines are allowed", func() {

		p := parsers.Array(q)

		DescribeTable("it evaluates the input correctly",
			func(input string, expected interface{}) {
				test.ExpectParserToEvaluate(p)(input, expected)
			},
			test.TableEntry("[]", []interface{}{}),
			test.TableEntry("[nil]", []interface{}{nil}),
			test.TableEntry("[1]", []interface{}{int64(1)}),
			test.TableEntry(`[1, "foo"]`, []interface{}{int64(1), "foo"}),
			test.TableEntry("[\n]", []interface{}{}),
			test.TableEntry("[\nnil,\n]", []interface{}{nil}),
			test.TableEntry("[\n1\n]", []interface{}{int64(1)}),
			test.TableEntry("[\n1,\n]", []interface{}{int64(1)}),
			test.TableEntry("[\n1,\n \"foo\",\n]", []interface{}{int64(1), "foo"}),
		)

		DescribeTable("it returns a parse error",
			func(input string, expectedErr error) {
				test.ExpectParserToHaveParseError(p)(input, expectedErr)
			},
			test.TableEntry("[", errors.New("was expecting \"]\" at testfile:1:2")),
			test.TableEntry("[1", errors.New("was expecting \"]\" at testfile:1:3")),
			test.TableEntry("[1,", errors.New("was expecting \"]\" at testfile:1:4")),
		)

		DescribeTable("it returns an eval error",
			func(input string, expectedErr error) {
				test.ExpectParserToHaveEvalError(p)(input, expectedErr)
			},
			test.TableEntry("[ERR, 1]", errors.New("ERR at testfile:1:2")),
			test.TableEntry("[1, ERR]", errors.New("ERR at testfile:1:5")),
		)
	})

})
