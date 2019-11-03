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
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text/terminal"
)

var _ = Describe("Not", func() {

	q := combinator.Choice(
		terminal.Bool("true", "false"),
		terminal.Integer(),
		terminal.Nil("nil"),
		test.EvalErrorParser(),
	).Name("value")

	p := parsers.Not(q)

	DescribeTable("it evaluates the input correctly",
		func(input string, expected interface{}) {
			test.ExpectParserToEvaluate(p)(input, expected)
		},
		test.TableEntry("1", int64(1)),
		test.TableEntry("nil", nil),
		test.TableEntry("! false", true),
		test.TableEntry("! true", false),
	)

	DescribeTable("it returns an eval error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveEvalError(p)(input, expectedErr)
		},
		test.TableEntry("! 1", errors.New("unsupported ! operation on int64 at testfile:1:1")),
		test.TableEntry("! nil", errors.New("unsupported ! operation on <nil> at testfile:1:1")),
		test.TableEntry("ERR", errors.New("ERR at testfile:1:1")),
		test.TableEntry("! ERR", errors.New("ERR at testfile:1:3")),
	)

	Context("When there is only one node", func() {
		It("should return the node", func() {
			expectedNode := terminal.NewIntegerNode(int64(1), parsley.Pos(1), parsley.Pos(2))
			test.ExpectParserToReturn(p, "1", expectedNode)
		})
	})

})
