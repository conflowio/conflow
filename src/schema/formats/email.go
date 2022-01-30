// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package formats

import (
	"errors"
	"net/mail"
	"reflect"
	"strings"
)

type Email struct {
	Default bool
}

func (e Email) ValidateValue(input string) (interface{}, error) {
	if strings.TrimSpace(input) != input {
		return nil, ErrValueTrimSpace
	}

	res, err := mail.ParseAddress(input)
	if err != nil {
		return nil, errors.New("must be a valid email address")
	}

	return *res, nil
}

func (e Email) StringValue(input interface{}) (string, bool) {
	switch v := input.(type) {
	case mail.Address:
		if v.Name == "" {
			return v.Address, true
		}

		return v.String(), true
	case *mail.Address:
		if v.Name == "" {
			return v.Address, true
		}

		return v.String(), true
	default:
		return "", false
	}
}

func (e Email) Type() (reflect.Type, bool) {
	return reflect.TypeOf(mail.Address{}), e.Default
}

func (e Email) PtrFunc() string {
	return "github.com/conflowio/conflow/src/schema/formats.EmailPtr"
}

func EmailPtr(v mail.Address) *mail.Address {
	return &v
}
