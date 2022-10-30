// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package openapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/conflowio/conflow/pkg/conflow"
	"github.com/conflowio/conflow/pkg/schema"
	schemainterpreters "github.com/conflowio/conflow/pkg/schema/interpreters"
	"github.com/conflowio/conflow/pkg/util/validation"
)

// @block "main"
type OpenAPI struct {
	// @name "openapi"
	// @default "3.1.0"
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

	// @dependency
	userContext interface{} `json:"-"`
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

func (o *OpenAPI) Validate(ctx context.Context) error {
	operationIDs := map[string]bool{}
	for pathName, p := range o.Paths {
		if err := p.IterateOperations(func(method string, op *Operation) error {
			if operationIDs[op.OperationID] {
				return validation.NewFieldError(fmt.Sprintf("paths[%q].%s.operationId", pathName, strings.ToLower(method)), errors.New("operation id must be unique"))
			}
			operationIDs[op.OperationID] = true
			return nil
		}); err != nil {
			return err
		}
	}
	return nil
}

func (o *OpenAPI) Run(ctx context.Context) (conflow.Result, error) {
	var schemaCtx *schema.Context
	if v, ok := ctx.Value(schema.GoContextKey).(*schema.Context); ok {
		schemaCtx = v
	} else if v, ok := o.userContext.(schema.ContextAware); ok {
		schemaCtx = v.SchemaContext()
	}
	schemaCtx.SetResolver(o)

	return nil, validation.Validate(ctx, o)
}

func (o *OpenAPI) ResolveSchema(ctx context.Context, uri string) (schema.Schema, error) {
	if strings.HasPrefix(uri, "#/components/schemas/") {
		name := strings.TrimPrefix(uri, "#/components/schemas/")
		return o.Schemas[name], nil
	}
	return nil, nil
}
