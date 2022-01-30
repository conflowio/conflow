// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package formats

import (
	"errors"
	"net"
	"reflect"
	"strings"
)

type IP struct {
	AllowIPv4 bool
	AllowIPv6 bool
	Default   bool
}

func (i IP) ValidateValue(input string) (interface{}, error) {
	if strings.TrimSpace(input) != input {
		return nil, ErrValueTrimSpace
	}

	if !i.AllowIPv4 && strings.Contains(input, ".") {
		return nil, errors.New("must be an IPv6 address")
	}

	if !i.AllowIPv6 && strings.Contains(input, ":") {
		return nil, errors.New("must be an IPv4 address")
	}

	res := net.ParseIP(input)
	if res == nil {
		return nil, errors.New("must be an IP address")
	}
	return res, nil
}

func (i IP) StringValue(input interface{}) (string, bool) {
	switch v := input.(type) {
	case net.IP:
		return v.String(), true
	case *net.IP:
		return v.String(), true
	default:
		return "", false
	}
}

func (i IP) Type() (reflect.Type, bool) {
	return reflect.TypeOf(net.IP{}), i.Default
}

func (i IP) PtrFunc() string {
	return "github.com/conflowio/conflow/src/conflow/schema/formats.IPPtr"
}

func IPPtr(v net.IP) *net.IP {
	return &v
}
