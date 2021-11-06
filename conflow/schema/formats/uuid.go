// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package formats

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type UUID struct{}

func (u UUID) Parse(input string) (interface{}, error) {
	if strings.TrimSpace(input) != input {
		return nil, ErrValueTrimSpace
	}

	res, err := uuid.Parse(input)
	if err != nil {
		return nil, fmt.Errorf("must be a valid UUID: %w", err)
	}
	return &res, nil
}

func (u UUID) Format(input interface{}) string {
	switch v := input.(type) {
	case *uuid.UUID:
		return v.String()
	default:
		panic(fmt.Errorf("invalid input type: %T", v))
	}
}
