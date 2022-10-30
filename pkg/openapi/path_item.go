// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package openapi

import (
	"context"
	"fmt"

	"github.com/conflowio/conflow/pkg/conflow"
	"github.com/conflowio/conflow/pkg/conflow/block"
	"github.com/conflowio/conflow/pkg/util/validation"
)

// @block "configuration"
type PathItem struct {
	Summary     string     `json:"summary,omitempty"`
	Description string     `json:"description,omitempty"`
	Get         *Operation `json:"get,omitempty"`
	Put         *Operation `json:"put,omitempty"`
	Post        *Operation `json:"post,omitempty"`
	Delete      *Operation `json:"delete,omitempty"`
	Options     *Operation `json:"options,omitempty"`
	Head        *Operation `json:"head,omitempty"`
	Patch       *Operation `json:"patch,omitempty"`
	Trace       *Operation `json:"trace,omitempty"`
	// @name "server"
	Servers []*Server `json:"servers,omitempty"`
	// @name "parameter"
	Parameters []*Parameter `json:"parameters,omitempty"`
}

func (p *PathItem) IterateOperations(f func(method string, op *Operation) error) error {
	ops := map[string]*Operation{
		"GET":     p.Get,
		"PUT":     p.Put,
		"POST":    p.Post,
		"DELETE":  p.Delete,
		"OPTIONS": p.Options,
		"HEAD":    p.Head,
		"PATCH":   p.Patch,
		"TRACE":   p.Trace,
	}
	for _, method := range []string{"GET", "PUT", "POST", "DELETE", "OPTIONS", "HEAD", "PATCH", "TRACE"} {
		if o := ops[method]; o != nil {
			if err := f(method, o); err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *PathItem) ParseContextOverride() conflow.ParseContextOverride {
	return conflow.ParseContextOverride{
		BlockTransformerRegistry: block.InterpreterRegistry{
			"get":       OperationInterpreter{},
			"put":       OperationInterpreter{},
			"post":      OperationInterpreter{},
			"delete":    OperationInterpreter{},
			"options":   OperationInterpreter{},
			"head":      OperationInterpreter{},
			"patch":     OperationInterpreter{},
			"trace":     OperationInterpreter{},
			"server":    ServerInterpreter{},
			"parameter": ParameterInterpreter{},
		},
	}
}

func (p *PathItem) Validate(ctx context.Context) error {
	parameters := map[string]bool{}
	for i, p := range p.Parameters {
		name := fmt.Sprintf("%s,%s", p.In, p.Name)
		if parameters[name] {
			return validation.NewFieldErrorf(fmt.Sprintf("parameters[%d]", i), "%s parameter must be unique", p.In)
		}
		parameters[name] = true
	}
	return nil
}
