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

var _ = Describe("Function", func() {

	var p *combinator.Sequence

	q := combinator.Choice(
		terminal.String(false),
		terminal.Integer(),
		terminal.Nil("nil", variable.TypeNil),
		parsers.Variable(),
		test.EvalErrorParser("ERR", variable.TypeUnknown),
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
		test.TableEntry(`test.func0("foo")`, errors.New("test.func0 expects 0 arguments at testfile:1:1")),
		test.TableEntry(`test.func1(5)`, errors.New("was expecting string at testfile:1:12")),
		test.TableEntry(`nonexisting("foo")`, errors.New("\"nonexisting\" function does not exist at testfile:1:1")),
	)

	DescribeTable("it returns an eval error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveEvalError(p)(input, expectedErr)
		},
		test.TableEntry("test.func1(ERR)", errors.New("ERR at testfile:1:12")),
	)

})
