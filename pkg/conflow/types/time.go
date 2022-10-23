// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package types

import (
	"encoding/json"
	"errors"
	"time"
)

type Time struct {
	Hour       int
	Minute     int
	Second     int
	NanoSecond int
	Location   *time.Location
}

func NewTime(hour, min, sec, nsec int, loc *time.Location) Time {
	return Time{
		Hour:       hour,
		Minute:     min,
		Second:     sec,
		NanoSecond: nsec,
		Location:   loc,
	}
}

func ParseTime(s string) (Time, error) {
	res, err := time.Parse("15:04:05.999999999Z07:00", s)
	if err != nil {
		// Errors returned by time.Parse are often meaningless to a user, so we just return a generic message
		return Time{}, errors.New("must be an RFC 3339 time value")
	}

	return Time{
		Hour:       res.Hour(),
		Minute:     res.Minute(),
		Second:     res.Second(),
		NanoSecond: res.Nanosecond(),
		Location:   res.Location(),
	}, err
}

func (t *Time) Time(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, t.Hour, t.Minute, t.Second, t.NanoSecond, t.Location)
}

func (t *Time) String() string {
	return time.Date(0, 1, 1, t.Hour, t.Minute, t.Second, t.NanoSecond, t.Location).Format("15:04:05.999999999Z07:00")
}

func (t Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *Time) UnmarshalJSON(bytes []byte) error {
	if len(bytes) == 0 {
		return nil
	}

	var s string
	if err := json.Unmarshal(bytes, &s); err != nil {
		return err
	}

	var err error
	*t, err = ParseTime(s)
	return err
}
