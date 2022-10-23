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

type DateTime struct{}

func (d DateTime) ValidateValue(input string) (interface{}, error) {
	if strings.TrimSpace(input) != input {
		return nil, ErrValueTrimSpace
	}

	res, err := time.Parse(time.RFC3339Nano, input)
	if err != nil {
		// Let's try to parse it without a timezone
		res, err = time.Parse("2006-01-02T15:04:05.999999999", input)
	}
	if err != nil {
		// Let's try to parse it as a date
		res, err = time.Parse("2006-01-02", input)
	}

	if err != nil {
		// Errors returned by time.Parse are often meaningless to a user, so we just return a generic message
		return nil, errors.New("must be an RFC 3339 date-time value")
	}

	return res, err
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
