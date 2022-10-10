// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package openapi

import (
	"github.com/conflowio/conflow/src/conflow"
	"github.com/conflowio/conflow/src/conflow/block"
)

// @block "configuration"
type Operation struct {
	Tags        []string `json:"tags,omitempty"`
	Summary     string   `json:"summary,omitempty"`
	Description string   `json:"description,omitempty"`
	// @required
	OperationID string `json:"operationId"`
	// @name "parameter"
	Parameters  []*Parameter `json:"parameters,omitempty"`
	RequestBody *RequestBody `json:"requestBody,omitempty"`
	// @name "response"
	Responses  map[string]*Response `json:"responses,omitempty"`
	Deprecated bool                 `json:"deprecated,omitempty"`
	// @name "server"
	Servers []*Server `json:"servers,omitempty"`
}

func (o *Operation) ParseContextOverride() conflow.ParseContextOverride {
	return conflow.ParseContextOverride{
		BlockTransformerRegistry: block.InterpreterRegistry{
			"parameter":    ParameterInterpreter{},
			"request_body": RequestBodyInterpreter{},
			"response":     ResponseInterpreter{},
			"server":       ServerInterpreter{},
		},
	}
}
