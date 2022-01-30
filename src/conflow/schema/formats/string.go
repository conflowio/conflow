// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package formats

import "reflect"

type String struct{}

func (s String) ValidateValue(input string) (interface{}, error) {
	return input, nil
}

func (s String) StringValue(input interface{}) (string, bool) {
	switch v := input.(type) {
	case string:
		return v, true
	case *string:
		return *v, true
	default:
		return "", false
	}
}

func (s String) Type() (reflect.Type, bool) {
	return reflect.TypeOf(""), true
}

func (s String) PtrFunc() string {
	return "github.com/conflowio/conflow/src/conflow/schema.StringPtr"
}
