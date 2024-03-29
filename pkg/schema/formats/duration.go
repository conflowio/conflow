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

	"github.com/conflowio/conflow/pkg/conflow/types"
)

type Duration struct{}

func (d Duration) ValidateValue(input string) (interface{}, error) {
	if strings.TrimSpace(input) != input {
		return nil, ErrValueTrimSpace
	}

	return types.ParseDuration(input)
}

func (d Duration) StringValue(input interface{}) (string, bool) {
	switch v := input.(type) {
	case time.Duration:
		return v.String(), true
	case *time.Duration:
		return v.String(), true
	case types.Duration:
		return time.Duration(v).String(), true
	case *types.Duration:
		return time.Duration(*v).String(), true
	default:
		return "", false
	}
}

func (d Duration) Type() (reflect.Type, bool) {
	return reflect.TypeOf(types.Duration(0)), true
}
