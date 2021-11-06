// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package types

import (
	"bytes"
	"strconv"
)

type RFC3339Duration struct {
	Year   int
	Month  int
	Day    int
	Hour   int
	Minute int
	Second int
	Week   int
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
