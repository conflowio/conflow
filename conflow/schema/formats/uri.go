// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package formats

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

type URI struct{}

func (u URI) Parse(input string) (interface{}, error) {
	if strings.TrimSpace(input) != input {
		return nil, ErrValueTrimSpace
	}

	res, err := url.Parse(input)
	if err != nil || res.Scheme == "" {
		return nil, errors.New("must be a valid URI")
	}

	return res, nil
}

func (u URI) Format(input interface{}) string {
	switch v := input.(type) {
	case *url.URL:
		return v.String()
	default:
		panic(fmt.Errorf("invalid input type: %T", v))
	}
}