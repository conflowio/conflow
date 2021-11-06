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

	"github.com/conflowio/conflow/conflow/types"
)

type IPv6CIDR struct{}

func (i IPv6CIDR) Parse(input string) (interface{}, error) {
	if strings.TrimSpace(input) != input {
		return nil, ErrValueTrimSpace
	}

	if !strings.Contains(input, ":") {
		return nil, errors.New("must be an IPv6 CIDR block")
	}

	ip, ipNet, err := net.ParseCIDR(input)
	if err != nil {
		return nil, fmt.Errorf("must be an IPv6 CIDR block: %w", err)
	}

	return &types.CIDR{
		IP:    ip,
		IPNet: ipNet,
	}, nil
}

func (i IPv6CIDR) Format(input interface{}) string {
	switch v := input.(type) {
	case *types.CIDR:
		ones, _ := v.IPNet.Mask.Size()
		return fmt.Sprintf("%s/%d", v.IP.String(), ones)
	default:
		panic(fmt.Errorf("invalid input type: %T", v))
	}
}
