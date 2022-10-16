// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package openapi

import (
	"encoding/json"
	"net/url"

	"github.com/conflowio/conflow/src/conflow"
	"github.com/conflowio/conflow/src/conflow/block"
)

// @block "configuration"
type Server struct {
	URL         url.URL `json:"url"`
	Description string  `json:"description,omitempty"`
	// @name "variable"
	Variables map[string]*ServerVariable `json:"variables,omitempty"`
}

func (s *Server) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		URL         string                     `json:"url"`
		Description string                     `json:"description,omitempty"`
		Variables   map[string]*ServerVariable `json:"variables,omitempty"`
	}{
		URL:         s.URL.String(),
		Description: s.Description,
		Variables:   s.Variables,
	})
}

func (s *Server) ParseContextOverride() conflow.ParseContextOverride {
	return conflow.ParseContextOverride{
		BlockTransformerRegistry: block.InterpreterRegistry{
			"variable": ServerVariableInterpreter{},
		},
	}
}
