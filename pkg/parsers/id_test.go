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

	"github.com/conflowio/conflow/pkg/conflow"
	"github.com/conflowio/conflow/pkg/parsers"
	"github.com/conflowio/conflow/pkg/test"
)

var _ = Describe("ID", func() {

	p := parsers.ID()

	DescribeTable("it evaluates the input correctly",
		func(input string, expected interface{}) {
			test.ExpectParserToEvaluate(p)(input, expected)
		},
		test.TableEntry(`a`, conflow.ID("a")),
		test.TableEntry(`a_b`, conflow.ID("a_b")),
		test.TableEntry(`abcdefghijklmnopqrstuvwxyz_0123456789`, conflow.ID("abcdefghijklmnopqrstuvwxyz_0123456789")),
	)

	DescribeTable("it returns a parse error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveParseError(p)(input, expectedErr)
		},
		test.TableEntry(`testkeyword`, errors.New("testkeyword is a reserved keyword at testfile:1:1")),
		test.TableEntry(`0ab`, errors.New("was expecting identifier at testfile:1:1")),
	)

	DescribeTable("it returns a static check error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveStaticCheckError(p)(input, expectedErr)
		},
		test.TableEntry(`a__b`, errors.New("invalid identifier \"a__b\" (did you mean \"a_b\"?) at testfile:1:1")),
		test.TableEntry(`_b`, errors.New("invalid identifier \"_b\" (did you mean \"b\"?) at testfile:1:1")),
		test.TableEntry(`b_`, errors.New("invalid identifier \"b_\" (did you mean \"b\"?) at testfile:1:1")),
		test.TableEntry(`UPPER`, errors.New("invalid identifier \"UPPER\" (did you mean \"upper\"?) at testfile:1:1")),
	)

})
