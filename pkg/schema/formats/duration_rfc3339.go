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

type DurationRFC3339 struct{}

func (d DurationRFC3339) ValidateValue(input string) (interface{}, error) {
	if strings.TrimSpace(input) != input {
		return nil, ErrValueTrimSpace
	}

	return types.ParseRFC3339Duration(input)
}

func (d DurationRFC3339) StringValue(input interface{}) (string, bool) {
	switch v := input.(type) {
	case types.RFC3339Duration:
		return v.String(), true
	case *types.RFC3339Duration:
		return v.String(), true
	default:
		return "", false
	}
}

func (d DurationRFC3339) Type() (reflect.Type, bool) {
	return reflect.TypeOf(types.RFC3339Duration{}), true
}
