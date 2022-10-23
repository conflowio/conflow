// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package types

import (
	"encoding/json"
	"errors"
	"net/mail"
)

type Email mail.Address

func ParseEmail(s string) (Email, error) {
	res, err := mail.ParseAddress(s)
	if err != nil {
		return Email{}, errors.New("must be a valid email address")
	}

	return Email(*res), nil
}

func (e *Email) String() string {
	if e.Name == "" {
		return e.Address
	}

	return (*mail.Address)(e).String()
}

func (e Email) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

func (e *Email) UnmarshalJSON(bytes []byte) error {
	if len(bytes) == 0 {
		return nil
	}

	var s string
	if err := json.Unmarshal(bytes, &s); err != nil {
		return err
	}

	var err error
	*e, err = ParseEmail(s)
	return err
}
