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
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/parsers"
	"github.com/opsidian/basil/test"
)

var _ = Describe("ID", func() {

	p := parsers.ID(basil.IDRegExpPattern)

	DescribeTable("it evaluates the input correctly",
		func(input string, expected interface{}) {
			test.ExpectParserToEvaluate(p)(input, expected)
		},
		test.TableEntry(`a`, basil.ID("a")),
		test.TableEntry(`a_b`, basil.ID("a_b")),
		test.TableEntry(`abcdefghijklmnopqrstuvwxyz_0123456789`, basil.ID("abcdefghijklmnopqrstuvwxyz_0123456789")),
	)

	DescribeTable("it returns a parse error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveParseError(p)(input, expectedErr)
		},
		test.TableEntry(`testkeyword`, errors.New("testkeyword is a reserved keyword at testfile:1:1")),
		test.TableEntry(`a__b`, errors.New("was expecting the end of input at testfile:1:2")),
		test.TableEntry(`_b`, errors.New("was expecting identifier at testfile:1:1")),
		test.TableEntry(`b_`, errors.New("was expecting the end of input at testfile:1:2")),
		test.TableEntry(`0ab`, errors.New("was expecting identifier at testfile:1:1")),
	)

})
