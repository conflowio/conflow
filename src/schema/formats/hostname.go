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
	"regexp"
	"strings"
)

var hostnameRegex = regexp.MustCompile(`^([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])(\.([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]{0,61}[a-zA-Z0-9]))*$`)

type Hostname struct{}

func (h Hostname) ValidateValue(input string) (interface{}, error) {
	if strings.TrimSpace(input) != input {
		return nil, ErrValueTrimSpace
	}

	if hostnameRegex.MatchString(input) {
		return input, nil
	}

	if ip := net.ParseIP(input); ip != nil {
		return input, nil
	}
	return nil, errors.New("must be a valid hostname")
}

func (h Hostname) StringValue(input interface{}) (string, bool) {
	switch v := input.(type) {
	case string:
		return v, true
	case *string:
		return *v, true
	default:
		return "", false
	}
}

func (h Hostname) Type() (reflect.Type, bool) {
	return reflect.TypeOf(""), false
}

func (h Hostname) PtrFunc() string {
	return "github.com/conflowio/conflow/src/schema/StringPtr"
}
