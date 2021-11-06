// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package formats

import (
	"fmt"
	"net"
	"strings"

	"github.com/conflowio/conflow/conflow/types"
)

type IPCIDR struct{}

func (i IPCIDR) Parse(input string) (interface{}, error) {
	if strings.TrimSpace(input) != input {
		return nil, ErrValueTrimSpace
	}

	ip, ipNet, err := net.ParseCIDR(input)
	if err != nil {
		return nil, fmt.Errorf("must be a CIDR block: %w", err)
	}

	return &types.CIDR{
		IP:    ip,
		IPNet: ipNet,
	}, nil
}

func (i IPCIDR) Format(input interface{}) string {
	switch v := input.(type) {
	case *types.CIDR:
		return v.String()
	default:
		panic(fmt.Errorf("invalid input type: %T", v))
	}
}
