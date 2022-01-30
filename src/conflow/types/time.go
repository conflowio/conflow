// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package types

import "time"

type Time struct {
	Hour       int
	Minute     int
	Second     int
	NanoSecond int
	Location   *time.Location
}

func (t Time) Time(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, t.Hour, t.Minute, t.Second, t.NanoSecond, t.Location)
}

func (t Time) String() string {
	return time.Date(0, 1, 1, t.Hour, t.Minute, t.Second, t.NanoSecond, t.Location).Format("15:04:05.999999999Z07:00")
}
