// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package formats

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

type IP struct{}

func (i IP) Parse(input string) (interface{}, error) {
	if strings.TrimSpace(input) != input {
		return nil, ErrValueTrimSpace
	}

	res := net.ParseIP(input)
	if res == nil {
		return nil, errors.New("must be an IP address")
	}
	return &res, nil
}

func (i IP) Format(input interface{}) string {
	switch v := input.(type) {
	case *net.IP:
		return v.String()
	default:
		panic(fmt.Errorf("invalid input type: %T", v))
	}
}
