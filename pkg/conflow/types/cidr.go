// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package types

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
)

type CIDR struct {
	IP    net.IP
	IPNet *net.IPNet
}

func NewCIDR(ip net.IP, ipNet *net.IPNet) CIDR {
	return CIDR{
		IP:    ip,
		IPNet: ipNet,
	}
}

func ParseCIDR(input string, allowIPv4, allowIPv6 bool) (CIDR, error) {
	if !allowIPv4 && strings.Contains(input, ".") {
		return CIDR{}, errors.New("must be an IPv6 CIDR block")
	}

	if !allowIPv6 && strings.Contains(input, ":") {
		return CIDR{}, errors.New("must be an IPv4 CIDR block")
	}

	ip, ipNet, err := net.ParseCIDR(input)
	if err != nil {
		return CIDR{}, fmt.Errorf("must be a CIDR block: %w", err)
	}

	return CIDR{
		IP:    ip,
		IPNet: ipNet,
	}, nil
}

func (c *CIDR) Equals(c2 CIDR) bool {
	if !c.IP.Equal(c2.IP) {
		return false
	}

	if !c.IPNet.IP.Equal(c2.IPNet.IP) {
		return false
	}

	if !bytes.Equal(c.IPNet.Mask, c.IPNet.Mask) {
		return false
	}

	return true
}

func (c *CIDR) String() string {
	ones, _ := c.IPNet.Mask.Size()
	return fmt.Sprintf("%s/%d", c.IP.String(), ones)
}

func (c CIDR) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

func (c *CIDR) UnmarshalJSON(bytes []byte) error {
	if len(bytes) == 0 {
		return nil
	}

	var s string
	if err := json.Unmarshal(bytes, &s); err != nil {
		return err
	}

	var err error
	*c, err = ParseCIDR(s, true, true)
	return err
}
