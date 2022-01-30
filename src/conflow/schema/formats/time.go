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

	"github.com/conflowio/conflow/src/conflow/types"
)

type Time struct{}

func (t Time) ValidateValue(input string) (interface{}, error) {
	if strings.TrimSpace(input) != input {
		return nil, ErrValueTrimSpace
	}

	res, err := time.Parse("15:04:05.999999999Z07:00", input)
	if err != nil {
		// Let's try to parse it without a timezone
		res, err = time.Parse("15:04:05.999999999", input)
	}

	if err != nil {
		// Errors returned by time.Parse are often meaningless to a user, so we just return a generic message
		return nil, errors.New("must be an RFC 3339 time value")
	}

	return types.Time{
		Hour:       res.Hour(),
		Minute:     res.Minute(),
		Second:     res.Second(),
		NanoSecond: res.Nanosecond(),
		Location:   res.Location(),
	}, err
}

func (t Time) StringValue(input interface{}) (string, bool) {
	switch v := input.(type) {
	case types.Time:
		return v.String(), true
	case *types.Time:
		return v.String(), true
	default:
		return "", false
	}
}

func (t Time) Type() (reflect.Type, bool) {
	return reflect.TypeOf(types.Time{}), true
}

func (t Time) PtrFunc() string {
	return "github.com/conflowio/conflow/src/conflow/schema/formats.TimePtr"
}

func TimePtr(v types.Time) *types.Time {
	return &v
}
