// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package types

import (
	"encoding/json"
	"errors"
	"net/url"
)

type URL url.URL

func ParseURL(s string) (URL, error) {
	res, err := url.Parse(s)
	if err != nil {
		return URL{}, errors.New("must be a valid URI")
	}

	return (URL)(*res), nil
}

func (u *URL) String() string {
	return (*url.URL)(u).String()
}

func (u URL) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}

func (u *URL) UnmarshalJSON(bytes []byte) error {
	if len(bytes) == 0 {
		return nil
	}

	var s string
	if err := json.Unmarshal(bytes, &s); err != nil {
		return err
	}

	var err error
	*u, err = ParseURL(s)
	return err
}
