// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package formats

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/conflowio/conflow/pkg/conflow/types"
)

var durationRegex = regexp.MustCompile(`^P((\d+Y)?(\d+M)?(\d+D)?(T(\d+H)?(\d+M)?(\d+S)?)?|\d+W)$`)

type DurationRFC3339 struct{}

func (d DurationRFC3339) ValidateValue(input string) (interface{}, error) {
	if strings.TrimSpace(input) != input {
		return nil, ErrValueTrimSpace
	}

	if input == "P" || input == "PT" || !durationRegex.MatchString(input) {
		return nil, fmt.Errorf("must be an RFC 3339 time duration")
	}

	var res types.RFC3339Duration

	isDate := true
	var amount string
	for _, ch := range input {
		var num int
		if strings.ContainsRune("YMDHSW", ch) {
			num, _ = strconv.Atoi(amount)
			amount = ""
		}
		switch ch {
		case 'P':
		case 'T':
			isDate = false
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			amount = amount + string(ch)
		case 'Y':
			res.Year = num
		case 'M':
			if isDate {
				res.Month = num
			} else {
				res.Minute = num
			}
		case 'D':
			res.Day = num
		case 'H':
			res.Hour = num
		case 'S':
			res.Second = num
		case 'W':
			res.Week = num
		default:
			return nil, fmt.Errorf("must be an RFC 3339 time duration")
		}
	}

	return res, nil
}

func (d DurationRFC3339) StringValue(input interface{}) (string, bool) {
	switch v := input.(type) {
	case types.RFC3339Duration:
		return v.String(), true
	case *types.RFC3339Duration:
		return v.String(), true
	default:
		return "", false
	}
}

func (d DurationRFC3339) Type() (reflect.Type, bool) {
	return reflect.TypeOf(types.RFC3339Duration{}), true
}
