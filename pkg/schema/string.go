// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package schema

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/conflowio/conflow/pkg/conflow/types"
	"github.com/conflowio/conflow/pkg/internal/utils"
	"github.com/conflowio/conflow/pkg/util/validation"
)

const (
	FormatConflowID = "conflow.ID"
)

//	@block {
//	  type = "configuration"
//	  path = "interpreters"
//	}
type String struct {
	Metadata

	Const   *string  `json:"const,omitempty"`
	Default *string  `json:"default,omitempty"`
	Enum    []string `json:"enum,omitempty"`
	Format  string   `json:"format,omitempty"`
	// @minimum 0
	MinLength int64         `json:"minLength,omitempty"`
	MaxLength *int64        `json:"maxLength,omitempty"`
	Nullable  bool          `json:"nullable,omitempty"`
	Pattern   *types.Regexp `json:"pattern,omitempty"`
}

func (s *String) AssignValue(imports map[string]string, valueName, resultName string) string {
	return fmt.Sprintf("%s = %s(%s)", resultName, assignFuncName(s, imports), valueName)
}

func (s *String) CompareValues(v1, v2 interface{}) int {
	s1, _ := s.valueOf(v1)
	s2, _ := s.valueOf(v2)

	switch {
	case s1 == nil && s2 == nil:
		return 0
	case s1 == nil:
		return -1
	case s2 == nil:
		return 1
	case *s1 == *s2:
		return 0
	case *s1 < *s2:
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
		fprintf(buf, "\tPattern: %sMustCompileRegexp(%q),\n", utils.EnsureUniqueGoPackageSelector(imports, "github.com/conflowio/conflow/pkg/conflow/types"), s.Pattern.String())
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
		s.Pattern = (*types.Regexp)(pattern)
	}

	return nil
}

func (s *String) Validate(ctx context.Context) error {
	return validation.ValidateObject(ctx,
		validation.ValidateField("minLength", validation.ValidatorFunc(func(ctx context.Context) error {
			if s.MinLength < 0 {
				return errors.New("must be greater than or equal to 0")
			}
			return nil
		})),
		validation.ValidatorFunc(func(ctx context.Context) error {
			if s.MaxLength != nil && s.MinLength > *s.MaxLength {
				return errors.New("minLength and maxLength constraints are impossible to fulfil")
			}
			return nil
		}),
		validateCommonFields(s, s.Const, s.Default, s.Enum),
	)
}

func (s *String) ValidateSchema(s2 Schema, _ bool) error {
	if s2.Type() != TypeString {
		return typeError("must be string")
	}

	return nil
}

func (s *String) ValidateValue(value interface{}) (interface{}, error) {
	v, ok := s.valueOf(value)
	if !ok {
		return nil, errors.New("must be string")
	}
	if v == nil {
		return nil, nil
	}

	if s.Const != nil && *s.Const != *v {
		return nil, fmt.Errorf("must be %q", s.StringValue(*s.Const))
	}

	if len(s.Enum) == 1 && s.Enum[0] != *v {
		return nil, fmt.Errorf("must be %q", s.StringValue(s.Enum[0]))
	}

	if len(s.Enum) > 0 {
		allowed := func() bool {
			for _, e := range s.Enum {
				if e == *v {
					return true
				}
			}
			return false
		}
		if !allowed() {
			return nil, fmt.Errorf("must be one of %s", s.join(s.Enum, ", "))
		}
	}

	ve := &validation.Error{}

	if s.MaxLength != nil {
		if s.MinLength == *s.MaxLength && len(*v) != int(s.MinLength) {
			switch s.MinLength {
			case 0:
				ve.AddError(errors.New("must be empty string"))
			case 1:
				ve.AddError(errors.New("must be a single character"))
			default:
				ve.AddErrorf("must be exactly %d characters long", s.MinLength)
			}
		} else {
			if int64(utf8.RuneCount([]byte(*v))) > *s.MaxLength {
				switch *s.MaxLength {
				case 0:
					ve.AddError(errors.New("must be empty string"))
				case 1:
					ve.AddError(errors.New("must be empty string or a single character"))
				default:
					ve.AddErrorf("must be no more than %d characters long", *s.MaxLength)
				}
			}
		}
	}

	if s.MinLength > 0 && int64(utf8.RuneCount([]byte(*v))) < s.MinLength {
		switch s.MinLength {
		case 1:
			ve.AddError(errors.New("can not be empty string"))
		default:
			ve.AddErrorf("must be at least %d characters long", s.MinLength)
		}
	}

	if s.Pattern != nil {
		if !(*regexp.Regexp)(s.Pattern).MatchString(*v) {
			ve.AddErrorf("must match regular expression: %s", s.Pattern.String())
		}
	}

	fv, err := s.format().ValidateValue(*v)
	if err != nil {
		ve.AddError(err)
	}

	if err := ve.ErrOrNil(); err != nil {
		return nil, err
	}

	if s.Nullable {
		return &fv, nil
	}
	return fv, nil
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

func (s *String) valueOf(value interface{}) (*string, bool) {
	switch v := value.(type) {
	case string:
		return &v, true
	case *string:
		return v, true
	default:
		sv, ok := s.format().StringValue(value)
		if !ok {
			return nil, false
		}
		return &sv, true
	}
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
