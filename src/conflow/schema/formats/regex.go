// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package formats

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type Regex struct{}

func (r Regex) ValidateValue(input string) (interface{}, error) {
	if strings.TrimSpace(input) != input {
		return nil, ErrValueTrimSpace
	}

	if input == "" {
		return nil, errors.New("must be valid regular expression")
	}

	res, err := regexp.Compile(input)
	if err != nil {
		return nil, fmt.Errorf("must be valid regular expression: %w", err)
	}

	return *res, nil
}

func (r Regex) StringValue(input interface{}) (string, bool) {
	switch v := input.(type) {
	case regexp.Regexp:
		return v.String(), true
	case *regexp.Regexp:
		return v.String(), true
	default:
		return "", false
	}
}

func (r Regex) Type() (reflect.Type, bool) {
	return reflect.TypeOf(regexp.Regexp{}), true
}

func (r Regex) PtrFunc() string {
	return "github.com/conflowio/conflow/src/conflow/schema/formats.RegexPtr"
}

func RegexPtr(v regexp.Regexp) *regexp.Regexp {
	return &v
}
