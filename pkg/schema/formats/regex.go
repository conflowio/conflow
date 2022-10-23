// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package formats

import (
	"reflect"
	"strings"

	"github.com/conflowio/conflow/pkg/conflow/types"
)

type Regex struct{}

func (r Regex) ValidateValue(input string) (interface{}, error) {
	if strings.TrimSpace(input) != input {
		return nil, ErrValueTrimSpace
	}

	return types.ParseRegexp(input)
}

func (r Regex) StringValue(input interface{}) (string, bool) {
	switch v := input.(type) {
	case types.Regexp:
		return v.String(), true
	case *types.Regexp:
		return v.String(), true
	default:
		return "", false
	}
}

func (r Regex) Type() (reflect.Type, bool) {
	return reflect.TypeOf(types.Regexp{}), true
}
