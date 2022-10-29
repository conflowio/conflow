// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package schema

import (
	"bytes"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"sync"

	"github.com/conflowio/conflow/pkg/internal/utils"
)

var _ Schema = &Reference{}

//	@block {
//	  type = "configuration"
//	  path = "interpreters"
//	}
type Reference struct {
	Metadata

	Nullable bool `json:"nullable,omitempty"`
	// @value
	// @required
	Ref string `json:"$ref"`

	// @ignore
	schema Schema
	// @ignore
	resolveSchema sync.Once
}

func (r *Reference) AssignValue(imports map[string]string, valueName, resultName string) string {
	return r.mustResolve().AssignValue(imports, valueName, resultName)
}

func (r *Reference) CompareValues(a, b interface{}) int {
	return r.mustResolve().CompareValues(a, b)
}

func (r *Reference) Copy() Schema {
	return &Reference{
		Ref: r.Ref,
	}
}

func (r *Reference) DefaultValue() interface{} {
	return r.mustResolve().DefaultValue()
}

func (r *Reference) GetNullable() bool {
	if n, ok := r.mustResolve().(Nullable); ok && n.GetNullable() {
		return true
	}
	return false
}

func (r *Reference) GoString(imports map[string]string) string {
	buf := bytes.NewBuffer(nil)
	fprintf(buf, "&%sReference{\n", schemaPkg(imports))
	if !reflect.ValueOf(r.Metadata).IsZero() {
		fprintf(buf, "\tMetadata: %s,\n", indent(r.Metadata.GoString(imports)))
	}
	if r.Nullable {
		fprintf(buf, "\tNullable: %#v,\n", r.Nullable)
	}
	fprintf(buf, "\tRef: %q,\n", r.Ref)
	buf.WriteRune('}')
	return buf.String()
}

func (r *Reference) GoType(imports map[string]string) string {
	u, err := url.Parse(r.Ref)
	if err != nil {
		panic(fmt.Errorf("reference %q is invalid: %w", r.Ref, err))
	}
	if u.Scheme != "" || strings.HasPrefix(r.Ref, "#") {
		return r.mustResolve().GoType(imports)
	}

	u.Path = strings.TrimPrefix(u.Path, "/")

	var path string
	typeName := u.Path
	for i := len(u.Path) - 1; i >= 0; i-- {
		if u.Path[i] == '.' {
			path = u.Path[0:i]
			typeName = u.Path[i+1:]
			break
		}
	}

	sel := utils.EnsureUniqueGoPackageSelector(imports, path)

	if r.Nullable {
		return fmt.Sprintf("*%s%s", sel, typeName)
	}

	return fmt.Sprintf("%s%s", sel, typeName)
}

func (r *Reference) SetNullable(nullable bool) {
	r.Nullable = nullable
}

func (r *Reference) StringValue(value interface{}) string {
	return r.mustResolve().StringValue(value)
}

func (r *Reference) Type() Type {
	return r.mustResolve().Type()
}

func (r *Reference) TypeString() string {
	return r.mustResolve().TypeString()
}

func (r *Reference) Validate(ctx *Context) error {
	return r.resolve(ctx)
}

func (r *Reference) mustResolve() Schema {
	if err := r.resolve(nil); err != nil {
		panic(err)
	}
	return r.schema
}

func (r *Reference) resolve(ctx *Context) (resolveErr error) {
	r.resolveSchema.Do(func() {
		var err error
		r.schema, err = ctx.ResolveSchema(r.Ref)
		if err != nil {
			resolveErr = fmt.Errorf("failed to resolve schema %q: %w", r.Ref, err)
			return
		}

		if r.schema == nil {
			resolveErr = fmt.Errorf("schema not found for %q", r.Ref)
		}
	})
	return
}

func (r *Reference) ValidateSchema(s Schema, compare bool) error {
	return r.mustResolve().ValidateSchema(s, compare)
}

func (r *Reference) ValidateValue(value interface{}) (interface{}, error) {
	return r.mustResolve().ValidateValue(value)
}
