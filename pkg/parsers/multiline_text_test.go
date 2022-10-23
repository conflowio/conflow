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

	"github.com/conflowio/conflow/pkg/parsers"
	"github.com/conflowio/conflow/pkg/test"
)

var _ = Describe("MultilineText", func() {

	p := parsers.MultilineText()

	DescribeTable("it evaluates the input correctly",
		func(input string, expected interface{}) {
			test.ExpectParserToEvaluate(p)(input, expected)
		},
		test.TableEntry(`'''
'''`, ""),
		test.TableEntry(`'''
a
'''`, "a\n"),
		test.TableEntry(`"""
a
"""`, "a\n"),
		test.TableEntry("```\na\n```", "a\n"),
		test.TableEntry(`'''
abc
def
'''`, "abc\ndef\n"),
		test.TableEntry(`'''
			a
		'''`, "a\n"),
		test.TableEntry(`'''
			abc
			def
		'''`, "abc\ndef\n"),
		test.TableEntry(`'''
			abc

			def
		'''`, "abc\n\ndef\n"),
		test.TableEntry(`'''
			abc
				def
		'''`, "abc\n\tdef\n"),
	)

	DescribeTable("it returns a parse error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveParseError(p)(input, expectedErr)
		},
		test.TableEntry(`''''''`, errors.New("text must start in a new line after ''' at testfile:1:4")),
		test.TableEntry(`'''a''`, errors.New("text must start in a new line after ''' at testfile:1:4")),
		test.TableEntry(`'a'`, errors.New("was expecting multiline text at testfile:1:1")),
		test.TableEntry(`''a''`, errors.New("was expecting multiline text at testfile:1:1")),
		test.TableEntry(`''a'''`, errors.New("was expecting multiline text at testfile:1:1")),
		test.TableEntry(`'''
a
'`, errors.New("was expecting multiline text at testfile:1:1")),
		test.TableEntry(`'''
a
''`, errors.New("was expecting multiline text at testfile:1:1")),
		test.TableEntry(`'''
a'''`, errors.New("the closing ''' must be on a new line at testfile:2:2")),
		test.TableEntry(`'''
	a
b
'''`, errors.New("every line must be empty or start with the same whitespace characters as the first line at testfile:3:1")),
		test.TableEntry(`'''
		a
	b
'''`, errors.New("every line must be empty or start with the same whitespace characters as the first line at testfile:3:2")),
		test.TableEntry(`'''
a
"""`, errors.New("was expecting multiline text at testfile:1:1")),
	)

})
