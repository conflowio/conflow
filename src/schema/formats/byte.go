// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package formats

import (
	"encoding/base64"
	"errors"
	"reflect"
	"strings"
)

type Binary struct {
	Default bool
}

func (b Binary) ValidateValue(input string) (interface{}, error) {
	if strings.TrimSpace(input) != input {
		return nil, ErrValueTrimSpace
	}

	res, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return nil, errors.New("must be a base64 encoded string")
	}

	return res, err
}

func (b Binary) StringValue(input interface{}) (string, bool) {
	switch v := input.(type) {
	case []byte:
		return base64.StdEncoding.EncodeToString(v), true
	case *[]byte:
		return base64.StdEncoding.EncodeToString(*v), true
	default:
		return "", false
	}
}

func (b Binary) Type() (reflect.Type, bool) {
	return reflect.TypeOf([]byte(nil)), b.Default
}
