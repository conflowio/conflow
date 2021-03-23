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
	"strings"
	"time"
)

type Time struct {
	Metadata

	Const   *time.Time  `json:"const,omitempty"`
	Default *time.Time  `json:"default,omitempty"`
	Enum    []time.Time `json:"enum,omitempty"`
}

func (s *Time) AssignValue(imports map[string]string, valueName, resultName string) string {
	timePackageName := EnsureUniqueGoPackageName(imports, "time")
	if s.Pointer {
		schemaPackageName := EnsureUniqueGoPackageName(imports, "github.com/opsidian/basil/basil/schema")
		return fmt.Sprintf(
			"%s = %s.TimePtr(%s.(%s.Time))",
			resultName,
			schemaPackageName,
			valueName,
			timePackageName,
		)
	}

	return fmt.Sprintf("%s = %s.(%s.Time)", resultName, valueName, timePackageName)
}

func (s *Time) CompareValues(v1, v2 interface{}) int {
	t1, ok := v1.(time.Time)
	if !ok {
		return -1
	}

	t2, ok := v2.(time.Time)
	if !ok {
		return 1
	}

	switch {
	case t1.Equal(t2):
		return 0
	case t1.Before(t2):
		return -1
	default:
		return 1
	}
}

func (s *Time) Copy() Schema {
	j, err := json.Marshal(s)
	if err != nil {
		panic(fmt.Errorf("failed to encode schema: %w", err))
	}

	cp := &Time{}
	if err := json.Unmarshal(j, cp); err != nil {
		panic(fmt.Errorf("failed to decode schema: %w", err))
	}

	return cp
}

func (s *Time) DefaultValue() interface{} {
	if s.Default == nil {
		return nil
	}
	return *s.Default
}

func (s *Time) GoString() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("&schema.Time{\n")
	if !reflect.ValueOf(s.Metadata).IsZero() {
		_, _ = fmt.Fprintf(buf, "\tMetadata: %s,\n", indent(s.Metadata.GoString()))
	}
	if s.Const != nil {
		_, _ = fmt.Fprintf(buf, "\tConst: schema.TimePtr(%s),\n", s.timeGoString(*s.Const))
	}
	if s.Default != nil {
		_, _ = fmt.Fprintf(buf, "\tDefault: schema.TimePtr(%s),\n", s.timeGoString(*s.Default))
	}
	if len(s.Enum) > 0 {
		_, _ = fmt.Fprintf(buf, "\tEnum: []time.Time{%s},\n", s.join(s.Enum, s.timeGoString, ", "))
	}
	buf.WriteRune('}')
	return buf.String()
}

func (s *Time) GoType(imports map[string]string) string {
	packageName := EnsureUniqueGoPackageName(imports, "time")

	if s.Pointer {
		return fmt.Sprintf("*%s.Time", packageName)
	}

	return fmt.Sprintf("%s.Time", packageName)
}

func (s *Time) MarshalJSON() ([]byte, error) {
	type Alias Time
	return json.Marshal(struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  string(s.Type()),
		Alias: (*Alias)(s),
	})
}

func (s *Time) StringValue(value interface{}) string {
	v, ok := value.(time.Time)
	if !ok {
		return ""
	}

	return v.String()
}

func (s *Time) Type() Type {
	return TypeTime
}

func (s *Time) TypeString() string {
	return string(TypeTime)
}

func (s *Time) ValidateSchema(s2 Schema, _ bool) error {
	if s2.Type() != TypeTime {
		return typeError("must be time")
	}

	return nil
}

func (s *Time) ValidateValue(value interface{}) error {
	v, ok := value.(time.Time)
	if !ok {
		return errors.New("must be date-time")
	}

	if s.Const != nil && !(*s.Const).Equal(v) {
		return fmt.Errorf("must be %s", s.StringValue(*s.Const))
	}

	if len(s.Enum) == 1 && !s.Enum[0].Equal(v) {
		return fmt.Errorf("must be %s", s.StringValue(s.Enum[0]))
	}

	if len(s.Enum) > 0 {
		allowed := func() bool {
			for _, e := range s.Enum {
				if e.Equal(v) {
					return true
				}
			}
			return false
		}
		if !allowed() {
			return fmt.Errorf(
				"must be one of %s",
				s.join(
					s.Enum, func(t time.Time) string {
						return t.String()
					},
					", ",
				),
			)
		}
	}

	return nil
}

func (s *Time) join(elems []time.Time, f func(time.Time) string, sep string) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return f(elems[0])
	}

	var b strings.Builder
	b.WriteString(f(elems[0]))
	for _, e := range elems[1:] {
		b.WriteString(sep)
		b.WriteString(f(e))
	}
	return b.String()
}

func (s *Time) timeGoString(t time.Time) string {
	t = t.UTC()
	return fmt.Sprintf(
		"time.Date(%d, %d, %d, %d, %d, %d, %d, time.UTC)",
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(),
		t.Nanosecond(),
	)
}

func TimeValue() Schema {
	return timeValueInst
}

var timeValueInst = &timeValue{
	Time: &Time{},
}

type timeValue struct {
	*Time
}

func (t *timeValue) Copy() Schema {
	return timeValueInst
}

func (t *timeValue) GoString() string {
	return "schema.TimeValue()"
}

func TimePtr(v time.Time) *time.Time {
	return &v
}
