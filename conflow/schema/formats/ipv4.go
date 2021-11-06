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

type IPv4 struct{}

func (i IPv4) Parse(input string) (interface{}, error) {
	if strings.TrimSpace(input) != input {
		return nil, ErrValueTrimSpace
	}

	if !strings.Contains(input, ".") {
		return nil, errors.New("must be an IPv4 address")
	}

	res := net.ParseIP(input)
	if res == nil {
		return nil, errors.New("must be an IPv4 address")
	}

	return &res, nil
}

func (i IPv4) Format(input interface{}) string {
	switch v := input.(type) {
	case *net.IP:
		return v.String()
	default:
		panic(fmt.Errorf("invalid input type: %T", v))
	}
}
