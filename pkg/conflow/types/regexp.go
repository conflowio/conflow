// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
)

type Regexp regexp.Regexp

func ParseRegexp(s string) (Regexp, error) {
	if s == "" {
		return Regexp{}, errors.New("must be valid regular expression")
	}

	res, err := regexp.Compile(s)
	if err != nil {
		return Regexp{}, fmt.Errorf("must be valid regular expression: %w", err)
	}

	return (Regexp)(*res), nil
}

func MustCompileRegexp(str string) *Regexp {
	return (*Regexp)(regexp.MustCompile(str))
}

func (r *Regexp) String() string {
	return (*regexp.Regexp)(r).String()
}

func (r Regexp) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

func (r *Regexp) UnmarshalJSON(bytes []byte) error {
	if len(bytes) == 0 {
		return nil
	}

	var s string
	if err := json.Unmarshal(bytes, &s); err != nil {
		return err
	}

	var err error
	*r, err = ParseRegexp(s)
	return err
}
