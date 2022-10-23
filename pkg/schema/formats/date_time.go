// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package formats

import (
	"reflect"
	"strings"
	"time"
)

type DateTime struct{}

func (d DateTime) ValidateValue(input string) (interface{}, error) {
	if strings.TrimSpace(input) != input {
		return nil, ErrValueTrimSpace
	}

	return time.Parse(time.RFC3339Nano, input)
}

func (d DateTime) StringValue(input interface{}) (string, bool) {
	switch v := input.(type) {
	case time.Time:
		return v.Format(time.RFC3339Nano), true
	case *time.Time:
		return v.Format(time.RFC3339Nano), true
	default:
		return "", false
	}
}

func (d DateTime) Type() (reflect.Type, bool) {
	return reflect.TypeOf(time.Time{}), true
}
