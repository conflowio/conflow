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
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/conflowio/conflow/src/internal/utils"
)

const (
	FormatConflowID = "conflow.ID"
)

type String struct {
	Metadata

	Const     *string        `json:"const,omitempty"`
	Default   *string        `json:"default,omitempty"`
	Enum      []string       `json:"enum,omitempty"`
	Format    string         `json:"format,omitempty"`
	MinLength int64          `json:"minLength,omitempty"`
	MaxLength *int64         `json:"maxLength,omitempty"`
	Nullable  bool           `json:"nullable,omitempty"`
	Pattern   *regexp.Regexp `json:"pattern,omitempty"`
}

func (s *String) AssignValue(imports map[string]string, valueName, resultName string) string {
	formatType, _ := s.format().Type()
	goType := utils.GoType(imports, formatType, false)

	if s.Nullable {
		return fmt.Sprintf(
			"%s = %sPointer(%s.(%s))",
			resultName,
			schemaPkg(imports),
			valueName,
			goType,
		)
	}

	return fmt.Sprintf(
		"%s = %s.(%s)",
		resultName,
		valueName,
		goType,
	)
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

func (s *String) GetNullable() bool {
	return s.Nullable
}

func (s *String) GoString(imports map[string]string) string {
	pkg := schemaPkg(imports)
	buf := bytes.NewBuffer(nil)
	fprintf(buf, "&%sString{\n", pkg)
	if !reflect.ValueOf(s.Metadata).IsZero() {
		fprintf(buf, "\tMetadata: %s,\n", indent(s.Metadata.GoString(imports)))
	}
	if s.Const != nil {
		fprintf(buf, "\tConst: %sPointer(%#v),\n", pkg, *s.Const)
	}
	if s.Default != nil {
		fprintf(buf, "\tDefault: %sPointer(%#v),\n", pkg, *s.Default)
	}
	if len(s.Enum) > 0 {
		fprintf(buf, "\tEnum: %#v,\n", s.Enum)
	}
	if len(s.Format) > 0 {
		fprintf(buf, "\tFormat: %#v,\n", s.Format)
	}
	if s.MinLength > 0 {
		fprintf(buf, "\tMinLength: %d,\n", s.MinLength)
	}
	if s.MaxLength != nil {
		fprintf(buf, "\tMaxLength: %sPointer(int64(%d)),\n", pkg, *s.MaxLength)
	}
	if s.Nullable {
		fprintf(buf, "\tNullable: %#v,\n", s.Nullable)
	}
	if s.Pattern != nil {
		pkgName := utils.EnsureUniqueGoPackageName(imports, "regexp")
		fprintf(buf, "\tPattern: %s.MustCompile(%q),\n", pkgName, s.Pattern.String())
	}
	buf.WriteRune('}')
	return buf.String()
}

func (s *String) GoType(imports map[string]string) string {
	formatType, _ := s.format().Type()
	return utils.GoType(imports, formatType, s.Nullable)
}

func (s *String) MarshalJSON() ([]byte, error) {
	type Alias String
	var pattern string
	if s.Pattern != nil {
		pattern = s.Pattern.String()
	}
	return json.Marshal(struct {
		Type    string `json:"type"`
		Pattern string `json:"pattern,omitempty"`
		*Alias
	}{
		Type:    string(s.Type()),
		Pattern: pattern,
		Alias:   (*Alias)(s),
	})
}

func (s *String) SetNullable(nullable bool) {
	s.Nullable = nullable
}

func (s *String) StringValue(value interface{}) string {
	if s, ok := value.(string); ok {
		return s
	}

	res, ok := s.format().StringValue(value)
	if !ok {
		panic(fmt.Errorf("invalid value %T in StringValue", value))
	}
	return res
}

func (s *String) Type() Type {
	return TypeString
}

func (s *String) TypeString() string {
	return string(TypeString)
}

func (s *String) UnmarshalJSON(input []byte) error {
	type Alias String
	v := &struct {
		Type    string  `json:"type"`
		Pattern *string `json:"pattern,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(input, v); err != nil {
		return err
	}

	if v.Pattern != nil {
		pattern, err := regexp.Compile(*v.Pattern)
		if err != nil {
			return fmt.Errorf("pattern is not a valid regular expression: %w", err)
		}
		s.Pattern = pattern
	}

	return nil
}

func (s *String) Validate(ctx *Context) error {
	return validateCommonFields(s, s.Const, s.Default, s.Enum)(ctx)
}

func (s *String) ValidateSchema(s2 Schema, _ bool) error {
	if s2.Type() != TypeString {
		return typeError("must be string")
	}

	return nil
}

func (s *String) ValidateValue(value interface{}) (interface{}, error) {
	v, ok := value.(string)
	if !ok {
		v, ok = s.format().StringValue(value)
		if !ok {
			return nil, errors.New("must be string")
		}
	}

	if s.Const != nil && *s.Const != v {
		return nil, fmt.Errorf("must be %q", s.StringValue(*s.Const))
	}

	if len(s.Enum) == 1 && s.Enum[0] != v {
		return nil, fmt.Errorf("must be %q", s.StringValue(s.Enum[0]))
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
			return nil, fmt.Errorf("must be one of %s", s.join(s.Enum, ", "))
		}
	}

	if s.MaxLength != nil && s.MinLength == *s.MaxLength && len(v) != int(s.MinLength) {
		switch s.MinLength {
		case 0:
			return nil, errors.New("must be empty string")
		case 1:
			return nil, errors.New("must be a single character")
		default:
			return nil, fmt.Errorf("must be exactly %d characters long", s.MinLength)
		}
	}

	if s.MinLength > 0 && int64(utf8.RuneCount([]byte(v))) < s.MinLength {
		switch s.MinLength {
		case 1:
			return nil, errors.New("can not be empty string")
		default:
			return nil, fmt.Errorf("must be at least %d characters long", s.MinLength)
		}
	}

	if s.MaxLength != nil && int64(utf8.RuneCount([]byte(v))) > *s.MaxLength {
		switch *s.MaxLength {
		case 0:
			return nil, errors.New("must be empty string")
		case 1:
			return nil, errors.New("must be empty string or a single character")
		default:
			return nil, fmt.Errorf("must be no more than %d characters long", *s.MaxLength)
		}
	}

	if s.Pattern != nil {
		if !s.Pattern.MatchString(v) {
			return nil, fmt.Errorf("must match regular expression: %s", s.Pattern.String())
		}
	}

	return s.format().ValidateValue(v)
}

func (s *String) format() Format {
	if f, ok := registeredFormats[s.Format]; ok {
		return f
	}
	return registeredFormats[FormatDefault]
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

func (s *stringValue) GoString(imports map[string]string) string {
	return fmt.Sprintf("%sStringValue()", schemaPkg(imports))
}
