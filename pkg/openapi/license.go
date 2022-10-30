// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package openapi

import (
	"github.com/conflowio/conflow/pkg/conflow/types"
)

// @block "configuration"
type License struct {
	// @required
	Name       string    `json:"name,omitempty"`
	Identifier string    `json:"identifier,omitempty"`
	URL        types.URL `json:"url"`
}
