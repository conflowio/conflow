// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package openapi

import (
	"encoding/json"
	"net/url"

	"github.com/conflowio/conflow/src/schema"
)

// @block "configuration"
type License struct {
	// @required
	Name       string  `json:"name,omitempty"`
	Identifier string  `json:"identifier,omitempty"`
	URL        url.URL `json:"url"`
}

func (l License) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name       string `json:"name,omitempty"`
		Identifier string `json:"identifier,omitempty"`
		URL        string `json:"url"`
	}{
		Name:       l.Name,
		Identifier: l.Identifier,
		URL:        l.URL.String(),
	})
}

func (l *License) Validate(*schema.Context) error {
	return nil
}
