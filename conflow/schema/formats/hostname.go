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
	"regexp"
	"strings"
)

var hostnameRegex = regexp.MustCompile(`^([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])(\.([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]{0,61}[a-zA-Z0-9]))*$`)

type Hostname struct{}

func (h Hostname) Parse(input string) (interface{}, error) {
	if strings.TrimSpace(input) != input {
		return nil, ErrValueTrimSpace
	}

	res := input
	if hostnameRegex.MatchString(input) {
		return &res, nil
	}

	if ip := net.ParseIP(input); ip != nil {
		return &res, nil
	}
	return nil, errors.New("must be a valid hostname")
}

func (h Hostname) Format(input interface{}) string {
	switch v := input.(type) {
	case *string:
		return *v
	default:
		panic(fmt.Errorf("invalid input type: %T", v))
	}
}
