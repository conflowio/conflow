// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package openapi

import (
	"encoding/json"
	"net/mail"
	"net/url"
	"strings"

	"github.com/conflowio/conflow/pkg/schema"
)

// @block "configuration"
type Contact struct {
	Name  string       `json:"name,omitempty"`
	URL   url.URL      `json:"url,omitempty"`
	Email mail.Address `json:"email,omitempty"`
}

func (c Contact) MarshalJSON() ([]byte, error) {
	email := c.Email.String()
	if c.Email.Name == "" {
		email = strings.Trim(email, "<>")
	}

	return json.Marshal(struct {
		Name  string `json:"name,omitempty"`
		URL   string `json:"url,omitempty"`
		Email string `json:"email,omitempty"`
	}{
		Name:  c.Name,
		URL:   c.URL.String(),
		Email: email,
	})
}

func (c *Contact) Validate(*schema.Context) error {
	return nil
}
