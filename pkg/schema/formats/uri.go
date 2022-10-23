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

	"github.com/conflowio/conflow/pkg/conflow/types"
)

type URI struct {
	Default       bool
	RequireScheme bool
}

func (u URI) ValidateValue(input string) (interface{}, error) {
	if strings.TrimSpace(input) != input {
		return nil, ErrValueTrimSpace
	}

	res, err := types.ParseURL(input)
	if err != nil || (u.RequireScheme && res.Scheme == "") {
		return nil, errors.New("must be a valid URI")
	}

	return res, nil
}

func (u URI) StringValue(input interface{}) (string, bool) {
	switch v := input.(type) {
	case types.URL:
		return v.String(), true
	case *types.URL:
		return v.String(), true
	default:
		return "", false
	}
}

func (u URI) Type() (reflect.Type, bool) {
	return reflect.TypeOf(types.URL{}), u.Default
}
