// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsers_test

import (
	"errors"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"

	"github.com/conflowio/conflow/pkg/conflow"
	"github.com/conflowio/conflow/pkg/parsers"
	"github.com/conflowio/conflow/pkg/test"
)

var _ = Describe("KeyValuePairs", func() {

	var p = parsers.KeyValuePairs()

	DescribeTable("it evaluates the input correctly",
		func(input string, expected interface{}) {
			test.ExpectParserToEvaluate(p)(input, expected)
		},
		test.TableEntry("", map[conflow.ID]interface{}{}),
		test.TableEntry(`key_1="foo bar"`, map[conflow.ID]interface{}{"key_1": "foo bar"}),
		test.TableEntry("key_1=2", map[conflow.ID]interface{}{"key_1": int64(2)}),
		test.TableEntry("key_1=1.2", map[conflow.ID]interface{}{"key_1": 1.2}),
		test.TableEntry("key_1=1h30m", map[conflow.ID]interface{}{"key_1": 90 * time.Minute}),
		test.TableEntry("key_1=true", map[conflow.ID]interface{}{"key_1": true}),
		test.TableEntry("key_1=false", map[conflow.ID]interface{}{"key_1": false}),
		test.TableEntry("key_1=[]", map[conflow.ID]interface{}{"key_1": []interface{}{}}),
		test.TableEntry("key_1=[1]", map[conflow.ID]interface{}{"key_1": []interface{}{int64(1)}}),
		test.TableEntry("key_1=[1,2]", map[conflow.ID]interface{}{"key_1": []interface{}{int64(1), int64(2)}}),
		test.TableEntry("key_1=2,key_2=4", map[conflow.ID]interface{}{"key_1": int64(2), "key_2": int64(4)}),
		test.TableEntry("key_1 = 1 , key_2 = [1 , 2]", map[conflow.ID]interface{}{"key_1": int64(1), "key_2": []interface{}{int64(1), int64(2)}}),
	)

	DescribeTable("it returns a parse error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveParseError(p)(input, expectedErr)
		},
		test.TableEntry(`key_1`, errors.New("was expecting \"=\" at testfile:1:6")),
		test.TableEntry(`key_1=`, errors.New("was expecting value at testfile:1:7")),
		test.TableEntry(`key_1=1,,key_2=2`, errors.New("was expecting the end of input at testfile:1:9")),
	)

	DescribeTable("it returns an eval error",
		func(input string, expectedErr error) {
			test.ExpectParserToHaveEvalError(p)(input, expectedErr)
		},
		test.TableEntry("key_1=1,key_1=2", errors.New("parameter \"key_1\" was already defined at testfile:1:9")),
	)
})
