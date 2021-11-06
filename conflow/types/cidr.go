// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package types

import (
	"fmt"
	"net"
)

type CIDR struct {
	IP    net.IP
	IPNet *net.IPNet
}

func (c CIDR) String() string {
	ones, _ := c.IPNet.Mask.Size()
	return fmt.Sprintf("%s/%d", c.IP.String(), ones)
}
