// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package schema

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/conflowio/conflow/internal/utils"
)

const (
	FormatConflowID = "conflow.ID"
)

type String struct {
	Metadata

	Const     *string  `json:"const,omitempty"`
	Default   *string  `json:"default,omitempty"`
	Enum      []string `json:"enum,omitempty"`
	Format    string   `json:"format,omitempty"`
	MinLength int64    `json:"minLength,omitempty"`
	MaxLength *int64   `json:"maxLength,omitempty"`
}

func (s *String) AssignValue(imports map[string]string, valueName, resultName string) string {
	if s.Pointer {
		schemaPackageName := utils.EnsureUniqueGoPackageName(imports, "github.com/conflowio/conflow/conflow/schema")
		return fmt.Sprintf("%s = %s.StringPtr(%s.(string))", resultName, schemaPackageName, valueName)
	}

	return fmt.Sprintf("%s = %s.(string)", resultName, valueName)
}

func (s *String) CompareValues(v1, v2 interface{}) int {
	s1, ok := v1.(string)
	if !ok {
		return -1
	}

	s2, ok := v2.(string)
	if !ok {
		return 1
	}

	switch {
	case s1 == s2:
		return 0
	case s1 < s2:
		return -1
	default:
		return 1
	}
}

func (s *String) Copy() Schema {
	j, err := json.Marshal(s)
	if err != nil {
		panic(fmt.Errorf("failed to encode schema: %w", err))
	}

	cp := &String{}
	if err := json.Unmarshal(j, cp); err != nil {
		panic(fmt.Errorf("failed to decode schema: %w", err))
	}

	return cp
}

func (s *String) DefaultValue() interface{} {
	if s.Default == nil {
		return nil
	}
	return *s.Default
}

func (s *String) GoString() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("&schema.String{\n")
	if !reflect.ValueOf(s.Metadata).IsZero() {
		_, _ = fmt.Fprintf(buf, "\tMetadata: %s,\n", indent(s.Metadata.GoString()))
	}
	if s.Const != nil {
		_, _ = fmt.Fprintf(buf, "\tConst: schema.StringPtr(%#v),\n", *s.Const)
	}
	if s.Default != nil {
		_, _ = fmt.Fprintf(buf, "\tDefault: schema.StringPtr(%#v),\n", *s.Default)
	}
	if len(s.Enum) > 0 {
		_, _ = fmt.Fprintf(buf, "\tEnum: %#v,\n", s.Enum)
	}
	if len(s.Format) > 0 {
		_, _ = fmt.Fprintf(buf, "\tFormat: %#v,\n", s.Format)
	}
	if s.MinLength > 0 {
		_, _ = fmt.Fprintf(buf, "\tMinLength: %d,\n", s.MinLength)
	}
	if s.MaxLength != nil {
		_, _ = fmt.Fprintf(buf, "\tMaxLength: schema.IntegerPtr(%d),\n", *s.MaxLength)
	}
	buf.WriteRune('}')
	return buf.String()
}

func (s *String) GoType(_ map[string]string) string {
	if s.Pointer {
		return "*string"
	}

	return "string"
}

func (s *String) MarshalJSON() ([]byte, error) {
	type Alias String
	return json.Marshal(struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  string(s.Type()),
		Alias: (*Alias)(s),
	})
}

func (s *String) StringValue(value interface{}) string {
	if v, ok := value.(string); ok {
		return v
	}
	return ""
}

func (s *String) Type() Type {
	return TypeString
}

func (s *String) TypeString() string {
	return string(TypeString)
}

func (s *String) UnmarshalJSON(input []byte) error {
	type Alias String
	return json.Unmarshal(input, &struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  string(s.Type()),
		Alias: (*Alias)(s),
	})
}

func (s *String) ValidateSchema(s2 Schema, _ bool) error {
	if s2.Type() != TypeString {
		return typeError("must be string")
	}

	return nil
}

func (s *String) ValidateValue(value interface{}) error {
	v, ok := value.(string)
	if !ok {
		return errors.New("must be string")
	}

	if s.Const != nil && *s.Const != v {
		return fmt.Errorf("must be %q", s.StringValue(*s.Const))
	}

	if len(s.Enum) == 1 && s.Enum[0] != v {
		return fmt.Errorf("must be %q", s.StringValue(s.Enum[0]))
	}

	if len(s.Enum) > 0 {
		allowed := func() bool {
			for _, e := range s.Enum {
				if e == v {
					return true
				}
			}
			return false
		}
		if !allowed() {
			return fmt.Errorf("must be one of %s", s.join(s.Enum, ", "))
		}
	}

	if s.MaxLength != nil && s.MinLength == *s.MaxLength && len(v) != int(s.MinLength) {
		switch s.MinLength {
		case 0:
			return errors.New("must be empty string")
		case 1:
			return errors.New("must be a single character")
		default:
			return fmt.Errorf("must be exactly %d characters long", s.MinLength)
		}
	}

	if s.MinLength > 0 && int64(utf8.RuneCount([]byte(v))) < s.MinLength {
		switch s.MinLength {
		case 1:
			return errors.New("can not be empty string")
		default:
			return fmt.Errorf("must be at least %d characters long", s.MinLength)
		}
	}

	if s.MaxLength != nil && int64(utf8.RuneCount([]byte(v))) > *s.MaxLength {
		switch *s.MaxLength {
		case 0:
			return errors.New("must be empty string")
		case 1:
			return errors.New("must be empty string or a single character")
		default:
			return fmt.Errorf("must be no more than %d characters long", *s.MaxLength)
		}
	}

	return nil
}

func (s *String) join(elems []string, sep string) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return strconv.Quote(elems[0])
	}

	var b strings.Builder
	b.WriteString(strconv.Quote(elems[0]))
	for _, e := range elems[1:] {
		b.WriteString(sep)
		b.WriteString(strconv.Quote(e))
	}
	return b.String()
}

func StringValue() Schema {
	return stringValueInst
}

var stringValueInst = &stringValue{
	String: &String{},
}

type stringValue struct {
	*String
}

func (s *stringValue) Copy() Schema {
	return stringValueInst
}

func (s *stringValue) GoString() string {
	return "schema.StringValue()"
}

func StringPtr(v string) *string {
	return &v
}
