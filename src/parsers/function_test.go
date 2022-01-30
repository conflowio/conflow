// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsers_test

import (
	"errors"

	"github.com/conflowio/conflow/src/conflow/schema"

	"github.com/conflowio/parsley/combinator"
	"github.com/conflowio/parsley/text/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"

	"github.com/conflowio/conflow/src/parsers"
	"github.com/conflowio/conflow/src/test"
)

var _ = Describe("Function", func() {

	var p *combinator.Sequence

	q := combinator.Choice(
		terminal.String(schema.StringValue(), false),
		terminal.Integer(schema.IntegerValue()),
		terminal.Nil(schema.NullValue(), "NULL"),
		parsers.Variable(),
		test.EvalErrorParser(schema.StringValue(), "ERR"),
	).Name("value")

	p = parsers.Function(q)

	DescribeTable("it evaluates the input correctly",
		func(input string, expected interface{}) {
			test.ExpectParserToEvaluate(p)(input, expected)
		},
		test.TableEntry(`test.func0()`, "Hello"),
		test.TableEntry(`test.func1("foo")`, "FOO"),
		test.TableEntry(`test.func2("foo", "bar")`, "foobar"),
	)

	DescribeTable("it returns a parse error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveParseError(p)(input, expectedErr)
		},
		test.TableEntry("test.func1", errors.New("was expecting \"(\" at testfile:1:11")),
		test.TableEntry("test.func1(", errors.New("was expecting \")\" at testfile:1:12")),
		test.TableEntry(`test.func1("foo"`, errors.New("was expecting \")\" at testfile:1:17")),
		test.TableEntry(`test.func1("foo",`, errors.New("was expecting \")\" at testfile:1:18")),
	)

	DescribeTable("it returns a static check error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveStaticCheckError(p)(input, expectedErr)
		},
		test.TableEntry(`test.func0("foo")`, errors.New("test.func0 requires exactly 0 argument, but got 1 at testfile:1:12")),
		test.TableEntry(`test.func1(5)`, errors.New("must be string at testfile:1:12")),
		test.TableEntry(`nonexisting("foo")`, errors.New("\"nonexisting\" function does not exist at testfile:1:1")),
	)

	DescribeTable("it returns an eval error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveEvalError(p)(input, expectedErr)
		},
		test.TableEntry("test.func1(ERR)", errors.New("ERR at testfile:1:12")),
	)

})
