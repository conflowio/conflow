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

	"github.com/conflowio/conflow/internal/utils"
)

type TimeDuration struct {
	Metadata

	Const   *time.Duration  `json:"const,omitempty"`
	Default *time.Duration  `json:"default,omitempty"`
	Enum    []time.Duration `json:"enum,omitempty"`
}

func (s *TimeDuration) AssignValue(imports map[string]string, valueName, resultName string) string {
	timePackageName := utils.EnsureUniqueGoPackageName(imports, "time")
	if s.Pointer {
		schemaPackageName := utils.EnsureUniqueGoPackageName(imports, "github.com/conflowio/conflow/conflow/schema")
		return fmt.Sprintf(
			"%s = %s.TimeDurationPtr(%s.(%s.Duration))",
			resultName,
			schemaPackageName,
			valueName,
			timePackageName,
		)
	}

	return fmt.Sprintf("%s = %s.(%s.Duration)", resultName, valueName, timePackageName)
}

func (s *TimeDuration) CompareValues(v1, v2 interface{}) int {
	t1, ok := v1.(time.Duration)
	if !ok {
		return -1
	}

	t2, ok := v2.(time.Duration)
	if !ok {
		return 1
	}

	switch {
	case t1 == t2:
		return 0
	case t1 < t2:
		return -1
	default:
		return 1
	}
}

func (s *TimeDuration) Copy() Schema {
	j, err := json.Marshal(s)
	if err != nil {
		panic(fmt.Errorf("failed to encode schema: %w", err))
	}

	cp := &TimeDuration{}
	if err := json.Unmarshal(j, cp); err != nil {
		panic(fmt.Errorf("failed to decode schema: %w", err))
	}

	return cp
}

func (s *TimeDuration) DefaultValue() interface{} {
	if s.Default == nil {
		return nil
	}
	return *s.Default
}

func (s *TimeDuration) GoString() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("&schema.TimeDuration{\n")
	if !reflect.ValueOf(s.Metadata).IsZero() {
		_, _ = fmt.Fprintf(buf, "\tMetadata: %s,\n", indent(s.Metadata.GoString()))
	}
	if s.Const != nil {
		_, _ = fmt.Fprintf(buf, "\tConst: schema.TimeDurationPtr(%#v),\n", *s.Const)
	}
	if s.Default != nil {
		_, _ = fmt.Fprintf(buf, "\tDefault: schema.TimeDurationPtr(%#v),\n", *s.Default)
	}
	if len(s.Enum) > 0 {
		_, _ = fmt.Fprintf(buf, "\tEnum: %#v,\n", s.Enum)
	}
	buf.WriteRune('}')
	return buf.String()
}

func (s *TimeDuration) GoType(imports map[string]string) string {
	packageName := utils.EnsureUniqueGoPackageName(imports, "time")

	if s.Pointer {
		return fmt.Sprintf("*%s.Duration", packageName)
	}

	return fmt.Sprintf("%s.Duration", packageName)
}

func (s *TimeDuration) MarshalJSON() ([]byte, error) {
	type Alias TimeDuration
	return json.Marshal(struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  string(s.Type()),
		Alias: (*Alias)(s),
	})
}

func (s *TimeDuration) StringValue(value interface{}) string {
	v, ok := value.(time.Duration)
	if !ok {
		return ""
	}

	return v.String()
}

func (s *TimeDuration) Type() Type {
	return TypeTimeDuration
}

func (s *TimeDuration) TypeString() string {
	return string(TypeTimeDuration)
}

func (s *TimeDuration) ValidateSchema(s2 Schema, _ bool) error {
	if s2.Type() != TypeTimeDuration {
		return typeError("must be time duration")
	}

	return nil
}

func (s *TimeDuration) ValidateValue(value interface{}) error {
	v, ok := value.(time.Duration)
	if !ok {
		return errors.New("must be time duration")
	}

	if s.Const != nil && *s.Const != v {
		return fmt.Errorf("must be %s", s.StringValue(*s.Const))
	}

	if len(s.Enum) == 1 && s.Enum[0] != v {
		return fmt.Errorf("must be %s", s.StringValue(s.Enum[0]))
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

	return nil
}

func (s *TimeDuration) join(elems []time.Duration, sep string) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return s.StringValue(elems[0])
	}

	var b strings.Builder
	b.WriteString(s.StringValue(elems[0]))
	for _, e := range elems[1:] {
		b.WriteString(sep)
		b.WriteString(s.StringValue(e))
	}
	return b.String()
}

func TimeDurationValue() Schema {
	return timeDurationValueInst
}

var timeDurationValueInst = &timeDurationValue{
	TimeDuration: &TimeDuration{},
}

type timeDurationValue struct {
	*TimeDuration
}

func (t *timeDurationValue) Copy() Schema {
	return timeDurationValueInst
}

func (t *timeDurationValue) GoString() string {
	return "schema.TimeDurationValue()"
}

func TimeDurationPtr(v time.Duration) *time.Duration {
	return &v
}
