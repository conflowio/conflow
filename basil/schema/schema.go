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
	"io"
	"time"

	"github.com/tidwall/gjson"
)

type Schema interface {
	AssignValue(imports map[string]string, valueName, resultName string) string
	CompareValues(a, b interface{}) int
	Copy() Schema
	DefaultValue() interface{}
	GetAnnotation(string) (string, bool)
	GetDeprecated() bool
	GetDescription() string
	GetExamples() []interface{}
	GetPointer() bool
	GetReadOnly() bool
	GetTitle() string
	GetWriteOnly() bool
	GoString() string
	GoType(imports map[string]string) string
	StringValue(interface{}) string
	Type() Type
	TypeString() string
	ValidateValue(interface{}) error
	ValidateSchema(s Schema, compare bool) error
}

type ArrayKind interface {
	GetItems() Schema
}

func IsArray(s Schema) bool {
	_, ok := s.(ArrayKind)
	return ok
}

type ObjectKind interface {
	GetProperties() map[string]Schema
	IsPropertyRequired(name string) bool
	GetRequired() []string
	GetStructProperties() map[string]string
}

func IsObject(s Schema) bool {
	_, ok := s.(ObjectKind)
	return ok
}

type MapKind interface {
	GetAdditionalProperties() Schema
}

func IsMap(s Schema) bool {
	_, ok := s.(MapKind)
	return ok
}

type FunctionKind interface {
	GetAdditionalParameters() *NamedSchema
	GetParameters() Parameters
	GetResult() Schema
	GetResultTypeFrom() string
}

func IsFunction(s Schema) bool {
	_, ok := s.(FunctionKind)
	return ok
}

type Directive interface {
	ApplyToSchema(Schema) error
}

func UnmarshalJSON(b []byte) (Schema, error) {
	switch {
	case bytes.Equal(b, []byte("null")):
		return nil, nil
	case bytes.Equal(b, []byte("false")):
		return False(), nil
	case bytes.Equal(b, []byte("true")), bytes.Equal(b, []byte("{}")):
		return &Untyped{}, nil
	}

	schemaType := gjson.GetBytes(b, "type")
	if !schemaType.Exists() || schemaType.IsArray() {
		s := &Untyped{}
		if err := json.Unmarshal(b, s); err != nil {
			return nil, err
		}
		return s, nil
	}

	var s Schema
	switch Type(gjson.GetBytes(b, "type").String()) {
	case TypeArray:
		s = &Array{}
	case TypeByteStream:
		s = &ByteStream{}
	case TypeBoolean:
		s = &Boolean{}
	case TypeFunction:
		s = &Function{}
	case TypeInteger:
		s = &Integer{}
	case TypeNull:
		s = &Null{}
	case TypeNumber:
		s = &Number{}
	case TypeObject:
		if gjson.GetBytes(b, "properties").Exists() {
			if a := gjson.GetBytes(b, "additionalProperties"); a.Exists() && a.Value() != false {
				return nil, fmt.Errorf("additionalProperties must be false if properties is set")
			}
			s = &Object{}
		} else {
			s = &Map{}
		}
	case TypeString:
		s = &String{}
	case TypeTime:
		s = &Time{}
	case TypeTimeDuration:
		s = &TimeDuration{}
	default:
		return nil, fmt.Errorf("unsupported type %s", gjson.GetBytes(b, "type").String())
	}

	if err := json.Unmarshal(b, s); err != nil {
		return nil, err
	}

	return s, nil
}

type SchemaUnmarshaler struct {
	Schema Schema
}

func (s *SchemaUnmarshaler) UnmarshalJSON(j []byte) error {
	var err error
	s.Schema, err = UnmarshalJSON(j)
	return err
}

type NamedSchema struct {
	Name   string
	Schema Schema
}

func (n NamedSchema) GoString() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("schema.NamedSchema{\n")
	_, _ = fmt.Fprintf(buf, "\tName: %q,\n", n.Name)
	_, _ = fmt.Fprintf(buf, "\tSchema: %s,\n", indent(n.Schema.GoString()))
	buf.WriteRune('}')
	return buf.String()
}

func GetSchemaForValue(value interface{}) (Schema, error) {
	switch v := value.(type) {
	case nil:
		return NullValue(), nil
	case bool:
		return BooleanValue(), nil
	case int64:
		return IntegerValue(), nil
	case float64:
		return NumberValue(), nil
	case string:
		return StringValue(), nil
	case []interface{}:
		items, err := GetSchemaForValues(len(v), func(i int) (Schema, error) {
			return GetSchemaForValue(v[i])
		})
		if err != nil {
			return nil, err
		}

		return &Array{Items: items}, nil
	case map[string]interface{}:
		sortedKeys := getSortedMapKeys(v)
		additionalProperties, err := GetSchemaForValues(len(v), func(i int) (Schema, error) {
			return GetSchemaForValue(v[sortedKeys[i]])
		})
		if err != nil {
			return nil, err
		}

		return &Map{AdditionalProperties: additionalProperties}, nil
	case io.Reader:
		return ByteStreamValue(), nil
	case time.Time:
		return TimeValue(), nil
	case time.Duration:
		return TimeDurationValue(), nil
	default:
		return nil, fmt.Errorf("value type %T is not allowed", v)
	}
}

func GetSchemaForValues(cnt int, s func(i int) (Schema, error)) (Schema, error) {
	if cnt == 0 {
		return NullValue(), nil
	}

	s1, err := s(0)
	if err != nil {
		return nil, err
	}

	for i := 1; i < cnt; i++ {
		s2, err := s(i)
		if err != nil {
			return nil, err
		}
		s1, err = GetCommonSchema(s1, s2)
		if err != nil {
			return nil, err
		}
	}

	return s1, nil
}

func GetCommonSchema(s1, s2 Schema) (Schema, error) {
	switch s1.Type() {
	case TypeArray:
		if s2.Type() == TypeNull {
			return s1, nil
		}

		if s2a, ok := s2.(ArrayKind); ok {
			items, err := GetCommonSchema(s1.(ArrayKind).GetItems(), s2a.GetItems())
			if err != nil {
				return nil, err
			}

			return &Array{
				Items: items,
			}, nil
		}
	case TypeMap:
		if s2.Type() == TypeNull {
			return s1, nil
		}

		if s2m, ok := s2.(MapKind); ok {
			additionalProperties, err := GetCommonSchema(
				s1.(MapKind).GetAdditionalProperties(),
				s2m.GetAdditionalProperties(),
			)
			if err != nil {
				return nil, err
			}

			return &Map{
				AdditionalProperties: additionalProperties,
			}, nil
		}
	case TypeObject:
		panic("GetCommonSchema should not be called for objects")
	case TypeInteger:
		if s2.Type() == TypeInteger {
			return s1, nil
		}

		if s2.Type() == TypeNumber {
			return s2, nil
		}
	case TypeNumber:
		if s2.Type() == TypeNumber || s2.Type() == TypeInteger {
			return s1, nil
		}
	case TypeNull:
		if s2.Type() == TypeNull {
			return s1, nil
		}

		if s2.Type() == TypeArray || s2.Type() == TypeMap {
			return s2, nil
		}
	default:
		if s1.Type() == s2.Type() {
			return s1, nil
		}
	}

	return nil, fmt.Errorf(
		"items must have the same type, but found %s and %s",
		s1.TypeString(),
		s2.TypeString(),
	)
}
