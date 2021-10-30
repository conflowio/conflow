// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package schema

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

type Function struct {
	Metadata

	AdditionalParameters *NamedSchema `json:"additionalParameters,omitempty"`
	Parameters           Parameters   `json:"parameters,omitempty"`
	Result               Schema       `json:"result,omitempty"`
	ResultTypeFrom       string       `json:"result_type_from,omitempty"`
}

func (f *Function) AssignValue(_ map[string]string, _, _ string) string {
	panic("AssignValue should not be called on a function schema")
}

func (f *Function) CompareValues(_, _ interface{}) int {
	return -1
}

func (f *Function) Copy() Schema {
	j, err := json.Marshal(f)
	if err != nil {
		panic(fmt.Errorf("failed to encode schema: %w", err))
	}

	cp := &Function{}
	if err := json.Unmarshal(j, cp); err != nil {
		panic(fmt.Errorf("failed to decode schema: %w", err))
	}

	return cp
}

func (f *Function) DefaultValue() interface{} {
	return nil
}

func (f *Function) GetAdditionalParameters() *NamedSchema {
	return f.AdditionalParameters
}

func (f *Function) GetParameters() Parameters {
	return f.Parameters
}

func (f *Function) GetResult() Schema {
	return f.Result
}

func (f *Function) GetResultTypeFrom() string {
	return f.ResultTypeFrom
}

func (f *Function) GoType(imports map[string]string) string {
	sb := &strings.Builder{}

	sb.WriteString("func(")

	for i, p := range f.Parameters {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(p.Name)
		sb.WriteRune(' ')
		sb.WriteString(p.Schema.GoType(imports))
	}

	if f.AdditionalParameters != nil {
		sb.WriteString(", ")
		sb.WriteString(f.AdditionalParameters.Name)
		sb.WriteString("...")
		sb.WriteString(f.AdditionalParameters.GoString())
	}

	sb.WriteRune(')')

	if f.Result != nil {
		sb.WriteRune(' ')
		sb.WriteString(f.Result.GoType(imports))
	}
	return sb.String()
}

func (f *Function) MarshalJSON() ([]byte, error) {
	type Alias Function
	return json.Marshal(struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  string(f.Type()),
		Alias: (*Alias)(f),
	})
}

func (f *Function) GoString() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("&schema.Function{\n")
	if !reflect.ValueOf(f.Metadata).IsZero() {
		_, _ = fmt.Fprintf(buf, "\tMetadata: %s,\n", indent(f.Metadata.GoString()))
	}
	if f.AdditionalParameters != nil {
		_, _ = fmt.Fprintf(buf, "\tAdditionalParameters: &%s,\n", indent(f.AdditionalParameters.GoString()))
	}
	if len(f.Parameters) > 0 {
		_, _ = fmt.Fprintf(buf, "\tParameters: %s,\n", indent(f.Parameters.GoString()))
	}
	if f.Result != nil {
		_, _ = fmt.Fprintf(buf, "\tResult: %s,\n", indent(f.Result.GoString()))
	}
	if f.ResultTypeFrom != "" {
		_, _ = fmt.Fprintf(buf, "\tResultTypeFrom: %q,\n", f.ResultTypeFrom)
	}
	buf.WriteRune('}')
	return buf.String()
}

func (f *Function) StringValue(interface{}) string {
	return "<function>"
}

func (f *Function) Type() Type {
	return TypeFunction
}

func (f *Function) TypeString() string {
	sb := &strings.Builder{}
	sb.WriteString("function(")
	for i, p := range f.Parameters {
		if i > 0 {
			sb.WriteString(", ")
		}
		_, _ = fmt.Fprintf(sb, "%s %s", p.Name, p.Schema.TypeString())
	}
	if f.AdditionalParameters != nil {
		_, _ = fmt.Fprintf(sb, "%s ...%s", f.AdditionalParameters.Name, f.AdditionalParameters.Schema.TypeString())
	}
	sb.WriteRune(')')

	if f.Result != nil {
		sb.WriteRune(' ')
		sb.WriteString(f.Result.TypeString())
	}

	return sb.String()
}

func (f *Function) UnmarshalJSON(j []byte) error {
	type Alias Function

	type namedSchema struct {
		Name   string             `json:"name"`
		Schema *SchemaUnmarshaler `json:"schema"`
	}

	v := struct {
		*Alias
		AdditionalParameters *namedSchema       `json:"additionalParameters,omitempty"`
		Result               *SchemaUnmarshaler `json:"result,omitempty"`
	}{
		Alias: (*Alias)(f),
	}
	if err := json.Unmarshal(j, &v); err != nil {
		return err
	}

	if v.AdditionalParameters != nil {
		f.AdditionalParameters = &NamedSchema{
			Name:   v.AdditionalParameters.Name,
			Schema: v.AdditionalParameters.Schema.Schema,
		}
	}

	if v.Result != nil {
		f.Result = v.Result.Schema
	}

	return nil
}

func (f *Function) ValidateSchema(s Schema, compare bool) error {
	panic("ValidateSchema on functions should not be called")
}

func (f *Function) ValidateValue(value interface{}) error {
	panic("ValidateValue on functions should not be called")
}

type Parameters []NamedSchema

func (p Parameters) MarshalJSON() ([]byte, error) {
	if len(p) == 0 {
		return []byte("{}"), nil
	}

	sb := bytes.NewBuffer([]byte{})

	sb.WriteRune('{')
	for i, param := range p {
		if i > 0 {
			sb.WriteRune(',')
		}
		sb.WriteString(strconv.Quote(param.Name))
		sb.WriteRune(':')
		if err := json.NewEncoder(sb).Encode(param.Schema); err != nil {
			return nil, err
		}
	}
	sb.WriteRune('}')

	return sb.Bytes(), nil
}

func (p Parameters) GoString() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("schema.Parameters{\n")
	for _, param := range p {
		_, _ = fmt.Fprintf(buf, "\t%s,\n", indent(param.GoString()))
	}
	buf.WriteRune('}')
	return buf.String()
}

func (p *Parameters) UnmarshalJSON(j []byte) error {
	if bytes.Equal(j, []byte("null")) || bytes.Equal(j, []byte("{}")) {
		return nil
	}

	o := gjson.GetBytes(j, "@this")
	if !o.IsObject() {
		return fmt.Errorf("was expecting object")
	}
	var err error
	o.ForEach(func(key, value gjson.Result) bool {
		k := key.String()

		exists := func() bool {
			for _, param := range *p {
				if param.Name == k {
					return true
				}
			}
			return false
		}()

		if exists {
			err = fmt.Errorf("duplicate parameter: %s", k)
			return false
		}

		s, serr := UnmarshalJSON([]byte(value.Raw))
		if serr != nil {
			err = serr
			return false
		}

		*p = append(*p, NamedSchema{
			Name:   k,
			Schema: s,
		})

		return true
	})

	return err
}
