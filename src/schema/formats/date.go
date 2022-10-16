// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package formats

import (
	"errors"
	"reflect"
	"strings"
	"time"
)

type Date struct{}

func (d Date) ValidateValue(input string) (interface{}, error) {
	if strings.TrimSpace(input) != input {
		return nil, ErrValueTrimSpace
	}

	res, err := time.Parse("2006-01-02", input)
	if err != nil {
		// Errors returned by time.Parse are often meaningless to a user, so we just return a generic message
		return nil, errors.New("must be an RFC 3339 date value")
	}

	return res, err
}

func (d Date) StringValue(input interface{}) (string, bool) {
	switch v := input.(type) {
	case time.Time:
		return v.Format("2006-01-02"), true
	case *time.Time:
		return v.Format("2006-01-02"), true
	default:
		return "", false
	}
}

func (d Date) Type() (reflect.Type, bool) {
	return reflect.TypeOf(time.Time{}), false
}
