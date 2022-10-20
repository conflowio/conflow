// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package openapi

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/conflowio/conflow/src/conflow"
	"github.com/conflowio/conflow/src/schema"
	schemainterpreters "github.com/conflowio/conflow/src/schema/interpreters"
)

// @block "main"
type OpenAPI struct {
	// @required
	// @name "openapi"
	OpenAPI string `json:"openapi"`
	// @required
	Info *Info `json:"info"`
	// @name "server"
	Servers []*Server `json:"servers,omitempty"`
	// @name "path"
	Paths map[string]*PathItem `json:"paths"`
	Tags  []string             `json:"tags,omitempty"`

	// @name "schema"
	Schemas map[string]schema.Schema `json:"-"`
	// @name "response"
	Responses map[string]*Response `json:"-"`
	// @name "parameter"
	Parameters map[string]*Parameter `json:"-"`
	// @name "request_body"
	RequestBodies map[string]*RequestBody `json:"-"`
}

func (o *OpenAPI) ParseContextOverride() conflow.ParseContextOverride {
	registry := schemainterpreters.Registry()
	registry["info"] = InfoInterpreter{}
	registry["server"] = ServerInterpreter{}
	registry["response"] = ResponseInterpreter{}
	registry["parameter"] = ParameterInterpreter{}
	registry["request_body"] = RequestBodyInterpreter{}
	registry["path"] = PathItemInterpreter{}

	return conflow.ParseContextOverride{
		BlockTransformerRegistry: registry,
	}
}

func (o *OpenAPI) MarshalJSON() ([]byte, error) {
	type Components struct {
		Schemas       map[string]schema.Schema `json:"schemas,omitempty"`
		Responses     map[string]*Response     `json:"response,omitempty"`
		Parameters    map[string]*Parameter    `json:"parameters,omitempty"`
		RequestBodies map[string]*RequestBody  `json:"requestBodies,omitempty"`
	}
	type Alias OpenAPI

	return json.Marshal(struct {
		*Alias
		Components Components `json:"components,omitempty"`
	}{
		Alias: (*Alias)(o),
		Components: Components{
			Schemas:       o.Schemas,
			Responses:     o.Responses,
			Parameters:    o.Parameters,
			RequestBodies: o.RequestBodies,
		},
	})
}

func (o *OpenAPI) Validate(ctx *schema.Context) error {
	return schema.ValidateAll(
		ctx,
		schema.Validate("info", o.Info),
		schema.ValidateArray("servers", o.Servers),
		schema.ValidateMap("paths", o.Paths),
		schema.ValidateMap("schemas", o.Schemas),
		schema.ValidateMap("responses", o.Responses),
		schema.ValidateMap("parameters", o.Parameters),
		schema.ValidateMap("requestBodies", o.RequestBodies),
	)
}

func (o *OpenAPI) Run(context.Context) (conflow.Result, error) {
	return nil, o.Validate(schema.NewContext().WithResolver(o))
}

func (o *OpenAPI) ResolveSchema(uri string) (schema.Schema, error) {
	if strings.HasPrefix(uri, "#/components/schemas/") {
		name := strings.TrimPrefix(uri, "#/components/schemas/")
		return o.Schemas[name], nil
	}
	return nil, nil
}
