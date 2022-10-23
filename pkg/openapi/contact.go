// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package openapi

import (
	"github.com/conflowio/conflow/pkg/conflow/types"
	"github.com/conflowio/conflow/pkg/schema"
)

// @block "configuration"
type Contact struct {
	Name  string      `json:"name,omitempty"`
	URL   types.URL   `json:"url,omitempty"`
	Email types.Email `json:"email,omitempty"`
}

func (c *Contact) Validate(*schema.Context) error {
	return nil
}
