// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsers_test

import (
	"errors"
	"time"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/parsers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	"github.com/opsidian/basil/test"
)

var _ = Describe("KeyValuePairs", func() {

	var p = parsers.KeyValuePairs()

	DescribeTable("it evaluates the input correctly",
		func(input string, expected interface{}) {
			test.ExpectParserToEvaluate(p)(input, expected)
		},
		test.TableEntry("", map[basil.ID]interface{}{}),
		test.TableEntry(`key_1="foo bar"`, map[basil.ID]interface{}{"key_1": "foo bar"}),
		test.TableEntry("key_1=2", map[basil.ID]interface{}{"key_1": int64(2)}),
		test.TableEntry("key_1=1.2", map[basil.ID]interface{}{"key_1": 1.2}),
		test.TableEntry("key_1=1h30m", map[basil.ID]interface{}{"key_1": 90 * time.Minute}),
		test.TableEntry("key_1=true", map[basil.ID]interface{}{"key_1": true}),
		test.TableEntry("key_1=false", map[basil.ID]interface{}{"key_1": false}),
		test.TableEntry("key_1=[]", map[basil.ID]interface{}{"key_1": []interface{}{}}),
		test.TableEntry("key_1=[1]", map[basil.ID]interface{}{"key_1": []interface{}{int64(1)}}),
		test.TableEntry("key_1=[1,2]", map[basil.ID]interface{}{"key_1": []interface{}{int64(1), int64(2)}}),
		test.TableEntry("key_1=2,key_2=4", map[basil.ID]interface{}{"key_1": int64(2), "key_2": int64(4)}),
		test.TableEntry("key_1 = 1 , key_2 = [1 , 2]", map[basil.ID]interface{}{"key_1": int64(1), "key_2": []interface{}{int64(1), int64(2)}}),
	)

	DescribeTable("it returns a parse error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveParseError(p)(input, expectedErr)
		},
		test.TableEntry(`key_1`, errors.New("was expecting \"=\" at testfile:1:6")),
		test.TableEntry(`key_1=`, errors.New("was expecting value at testfile:1:7")),
		test.TableEntry(`key_1=1,,key_2=2`, errors.New("was expecting parameter name and value pair at testfile:1:9")),
	)

	DescribeTable("it returns an eval error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveEvalError(p)(input, expectedErr)
		},
		test.TableEntry("key_1=1,key_1=2", errors.New("parameter \"key_1\" was already defined at testfile:1:9")),
	)
})
