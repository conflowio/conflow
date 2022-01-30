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
	"reflect"
	"strings"

	"github.com/conflowio/conflow/src/conflow/types"
)

type IPCIDR struct {
	AllowIPv4 bool
	AllowIPv6 bool
	Default   bool
}

func (i IPCIDR) ValidateValue(input string) (interface{}, error) {
	if strings.TrimSpace(input) != input {
		return nil, ErrValueTrimSpace
	}

	if !i.AllowIPv4 && strings.Contains(input, ".") {
		return nil, errors.New("must be an IPv6 CIDR block")
	}

	if !i.AllowIPv6 && strings.Contains(input, ":") {
		return nil, errors.New("must be an IPv4 CIDR block")
	}

	ip, ipNet, err := net.ParseCIDR(input)
	if err != nil {
		return nil, fmt.Errorf("must be a CIDR block: %w", err)
	}

	return types.CIDR{
		IP:    ip,
		IPNet: ipNet,
	}, nil
}

func (i IPCIDR) StringValue(input interface{}) (string, bool) {
	switch v := input.(type) {
	case types.CIDR:
		return v.String(), true
	case *types.CIDR:
		return v.String(), true
	default:
		return "", false
	}
}

func (i IPCIDR) Type() (reflect.Type, bool) {
	return reflect.TypeOf(types.CIDR{}), i.Default
}

func (i IPCIDR) PtrFunc() string {
	return "github.com/conflowio/conflow/src/conflow/schema/formats.IPCIDRPtr"
}

func IPCIDRPtr(v types.CIDR) *types.CIDR {
	return &v
}
