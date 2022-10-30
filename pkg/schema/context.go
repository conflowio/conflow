// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package schema

import (
	"context"

	"github.com/conflowio/conflow/pkg/util/validation"
)

type GoContextKeyType string

const GoContextKey GoContextKeyType = "github.com/conflowio/conflow/pkg/schema.Context"

type Context struct {
	resolver                 Resolver
	delayedValidators        []validation.Validator
	runningDelatedValidators bool
}

func NewContext() *Context {
	return &Context{}
}

func (c *Context) SchemaContext() *Context {
	return c
}

func (c *Context) SetResolver(resolver Resolver) *Context {
	c.resolver = resolver
	return c
}

func (c *Context) ResolveSchema(ctx context.Context, uri string) (Schema, error) {
	if c != nil && c.resolver != nil {
		s, err := c.resolver.ResolveSchema(ctx, uri)
		if s != nil || err != nil {
			return s, err
		}
	}

	return Get(uri)
}

func (c *Context) AddDelayedValidator(v validation.Validator) {
	c.delayedValidators = append(c.delayedValidators, v)
}

func (c *Context) IsRunningDelayedValidators() bool {
	if c == nil {
		return true
	}

	return c.runningDelatedValidators
}

func (c *Context) RunDelayedValidators(ctx context.Context) error {
	ve := &validation.Error{}
	for _, v := range c.delayedValidators {
		if err := v.Validate(ctx); err != nil {
			ve.AddError(err)
		}
	}
	return ve.ErrOrNil()
}

type ContextAware interface {
	SchemaContext() *Context
}
