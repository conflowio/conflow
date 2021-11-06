// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package formats

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"
)

type Email struct{}

func (e Email) Parse(input string) (interface{}, error) {
	if strings.TrimSpace(input) != input {
		return nil, ErrValueTrimSpace
	}

	res, err := mail.ParseAddress(input)
	if err != nil {
		return nil, errors.New("must be a valid email address")
	}

	return res, nil
}

func (e Email) Format(input interface{}) string {
	switch v := input.(type) {
	case *mail.Address:
		if v.Name == "" {
			return v.Address
		}

		return v.String()
	default:
		panic(fmt.Errorf("invalid input type: %T", v))
	}
}
