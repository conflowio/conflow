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

	"github.com/conflowio/conflow/src/internal/utils"
)

var _ Schema = &Reference{}

type Reference struct {
	Metadata

	Nullable bool `json:"nullable,omitempty"`
	// @value
	// @required
	Ref string `json:"$ref"`

	// @ignore
	schema Schema
}

func (r *Reference) AssignValue(imports map[string]string, valueName, resultName string) string {
	return r.getSchema().AssignValue(imports, valueName, resultName)
}

func (r *Reference) CompareValues(a, b interface{}) int {
	return r.getSchema().CompareValues(a, b)
}

func (r *Reference) Copy() Schema {
	return &Reference{
		Ref: r.Ref,
	}
}

func (r *Reference) DefaultValue() interface{} {
	return r.getSchema().DefaultValue()
}

func (r *Reference) GetNullable() bool {
	return r.Nullable
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

	if path == "" || imports["."] == path {
		if r.Nullable {
			return fmt.Sprintf("*%s", typeName)
		}
		return typeName
	}

	packageName := utils.EnsureUniqueGoPackageName(imports, path)

	if r.Nullable {
		return fmt.Sprintf("*%s.%s", packageName, typeName)
	}

	return fmt.Sprintf("%s.%s", packageName, typeName)
}

func (r *Reference) SetNullable(nullable bool) {
	r.Nullable = nullable
}

func (r *Reference) StringValue(value interface{}) string {
	return r.getSchema().StringValue(value)
}

func (r *Reference) Type() Type {
	return r.getSchema().Type()
}

func (r *Reference) TypeString() string {
	return r.getSchema().TypeString()
}

func (r *Reference) Validate(ctx *Context) error {
	if r.schema != nil {
		return nil
	}

	s, err := ctx.resolver.ResolveSchema(r.Ref)
	if err != nil {
		return fmt.Errorf("failed to resolve schema %q: %w", r.Ref, err)
	}

	if s == nil {
		return fmt.Errorf("schema not found for %q", r.Ref)
	}

	r.schema = s

	return nil
}

func (r *Reference) ValidateSchema(s Schema, compare bool) error {
	return r.getSchema().ValidateSchema(s, compare)
}

func (r *Reference) ValidateValue(value interface{}) (interface{}, error) {
	return r.getSchema().ValidateValue(value)
}

func (r *Reference) getSchema() Schema {
	if err := r.Validate(nil); err != nil {
		panic(err)
	}
	return r.schema
}
