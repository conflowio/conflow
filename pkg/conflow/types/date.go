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

type Date struct {
	Year  int
	Month time.Month
	Day   int
}

func NewDate(year int, month time.Month, day int) Date {
	return Date{
		Year:  year,
		Month: month,
		Day:   day,
	}
}

func ParseDate(s string) (Date, error) {
	res, err := time.Parse("2006-01-02", s)
	if err != nil {
		// Errors returned by time.Parse are often meaningless to a user, so we just return a generic message
		return Date{}, errors.New("must be an RFC 3339 date value")
	}

	return Date{
		Year:  res.Year(),
		Month: res.Month(),
		Day:   res.Day(),
	}, nil
}

func (d *Date) Time(hour, min, sec, nsec int, loc *time.Location) time.Time {
	return time.Date(d.Year, d.Month, d.Day, hour, min, sec, nsec, loc)
}

func (d *Date) String() string {
	return time.Date(d.Year, d.Month, d.Day, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
}

func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Date) UnmarshalJSON(bytes []byte) error {
	if len(bytes) == 0 {
		return nil
	}

	var s string
	if err := json.Unmarshal(bytes, &s); err != nil {
		return err
	}

	var err error
	*d, err = ParseDate(s)
	return err
}
