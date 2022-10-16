// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package openapi

import "github.com/conflowio/conflow/src/schema"

// @block "configuration"
type ServerVariable struct {
	// @min_items 1
	Enum []string `json:"enum,omitempty"`
	// @required
	Default     string `json:"default"`
	Description string `json:"description,omitempty"`
}

func (s *ServerVariable) Validate(*schema.Context) error {
	return nil
}
