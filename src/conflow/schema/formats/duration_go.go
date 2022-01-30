// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package formats

import (
	"reflect"
	"time"
)

type DurationGo struct{}

func (d DurationGo) ValidateValue(input string) (interface{}, error) {
	return time.ParseDuration(input)
}

func (d DurationGo) StringValue(input interface{}) (string, bool) {
	switch v := input.(type) {
	case time.Duration:
		return v.String(), true
	case *time.Duration:
		return v.String(), true
	default:
		return "", false
	}
}

func (d DurationGo) Type() (reflect.Type, bool) {
	return reflect.TypeOf(time.Duration(0)), true
}

func (d DurationGo) PtrFunc() string {
	return "github.com/conflowio/conflow/src/conflow/schema/formats.DurationGoPtr"
}

func DurationGoPtr(v time.Duration) *time.Duration {
	return &v
}
