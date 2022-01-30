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

type ReferenceResolver interface {
	ResolveSchemaReference(string) Schema
}

type Reference struct {
	Metadata

	Nullable bool              `json:"nullable,omitempty"`
	Ref      string            `json:"ref,omitempty"`
	Resolver ReferenceResolver `json:"-"`
	schema   Schema
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

func (r *Reference) GoString(map[string]string) string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("&schema.Reference{\n")
	if !reflect.ValueOf(r.Metadata).IsZero() {
		_, _ = fmt.Fprintf(buf, "\tMetadata: %s,\n", indent(r.Metadata.GoString()))
	}
	if r.Nullable {
		_, _ = fmt.Fprintf(buf, "\tNullable: %#v,\n", r.Nullable)
	}
	_, _ = fmt.Fprintf(buf, "\tRef: %q,\n", r.Ref)
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
	if r.schema == nil && r.Resolver == nil {
		return TypeReference
	}

	return r.getSchema().Type()
}

func (r *Reference) TypeString() string {
	return r.getSchema().TypeString()
}

func (r *Reference) ValidateSchema(s Schema, compare bool) error {
	return r.getSchema().ValidateSchema(s, compare)
}

func (r *Reference) ValidateValue(value interface{}) (interface{}, error) {
	return r.getSchema().ValidateValue(value)
}

func (r *Reference) getSchema() Schema {
	if r.schema == nil {
		r.schema = r.Resolver.ResolveSchemaReference(r.Ref).Copy()
		r.schema.(MetadataAccessor).Merge(r.Metadata)
	}

	return r.schema
}
