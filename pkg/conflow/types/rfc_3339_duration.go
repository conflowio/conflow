// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var rfc3339DurationRegex = regexp.MustCompile(`^P((\d+Y)?(\d+M)?(\d+D)?(T(\d+H)?(\d+M)?(\d+S)?)?|\d+W)$`)

type RFC3339Duration struct {
	Year   int
	Month  int
	Day    int
	Hour   int
	Minute int
	Second int
	Week   int
}

func NewRFC3339Duration(year, month, day, hour, minute, second, week int) RFC3339Duration {
	return RFC3339Duration{
		Year:   year,
		Month:  month,
		Day:    day,
		Hour:   hour,
		Minute: minute,
		Second: second,
		Week:   week,
	}
}

func ParseRFC3339Duration(s string) (RFC3339Duration, error) {
	if s == "P" || s == "PT" || !rfc3339DurationRegex.MatchString(s) {
		return RFC3339Duration{}, fmt.Errorf("must be an RFC 3339 time duration")
	}

	var res RFC3339Duration

	isDate := true
	var amount string
	for _, ch := range s {
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
			return RFC3339Duration{}, fmt.Errorf("must be an RFC 3339 time duration")
		}
	}

	return res, nil
}

func (r RFC3339Duration) String() string {
	buf := bytes.NewBuffer([]byte{})
	buf.WriteRune('P')
	if r.Week > 0 {
		buf.WriteString(strconv.Itoa(r.Week))
		buf.WriteRune('W')
	}
	if r.Year > 0 {
		buf.WriteString(strconv.Itoa(r.Year))
		buf.WriteRune('Y')
	}
	if r.Month > 0 {
		buf.WriteString(strconv.Itoa(r.Month))
		buf.WriteRune('M')
	}
	if r.Day > 0 {
		buf.WriteString(strconv.Itoa(r.Day))
		buf.WriteRune('D')
	}
	if r.Hour > 0 || r.Minute > 0 || r.Second > 0 {
		buf.WriteRune('T')
	}
	if r.Hour > 0 {
		buf.WriteString(strconv.Itoa(r.Hour))
		buf.WriteRune('H')
	}
	if r.Minute > 0 {
		buf.WriteString(strconv.Itoa(r.Minute))
		buf.WriteRune('M')
	}
	if r.Second > 0 {
		buf.WriteString(strconv.Itoa(r.Second))
		buf.WriteRune('S')
	}
	return buf.String()
}

func (r RFC3339Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

func (r *RFC3339Duration) UnmarshalJSON(bytes []byte) error {
	if len(bytes) == 0 {
		return nil
	}

	var s string
	if err := json.Unmarshal(bytes, &s); err != nil {
		return err
	}

	var err error
	*r, err = ParseRFC3339Duration(s)
	return err
}
