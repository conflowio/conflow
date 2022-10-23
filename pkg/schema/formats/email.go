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

type Email struct {
	Default bool
}

func (e Email) ValidateValue(input string) (interface{}, error) {
	if strings.TrimSpace(input) != input {
		return nil, ErrValueTrimSpace
	}

	return types.ParseEmail(input)
}

func (e Email) StringValue(input interface{}) (string, bool) {
	switch v := input.(type) {
	case types.Email:
		return v.String(), true
	case *types.Email:
		return v.String(), true
	default:
		return "", false
	}
}

func (e Email) Type() (reflect.Type, bool) {
	return reflect.TypeOf(types.Email{}), e.Default
}
