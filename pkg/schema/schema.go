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
	"fmt"
	"io"
	"reflect"
	"regexp"
	"time"

	"github.com/tidwall/gjson"
)

// NameRegExpPattern is the regular expression for a valid identifier
const NameRegExpPattern = "[a-z][a-z0-9]*(?:_[a-z0-9]+)*"

// NameRegExp is a compiled regular expression object for a valid identifier
var NameRegExp = regexp.MustCompile("^" + NameRegExpPattern + "$")
var FieldNameRegexp = regexp.MustCompile("^[_a-zA-Z][_a-zA-Z0-9]*$")

type Schema interface {
	AssignValue(imports map[string]string, valueName, resultName string) string
	CompareValues(a, b interface{}) int
	Copy() Schema
	DefaultValue() interface{}
	GetAnnotation(string) string
	GetDeprecated() bool
	GetDescription() string
	GetExamples() []interface{}
	GetID() string
	GetReadOnly() bool
	GetTitle() string
	GetWriteOnly() bool
	GoString(imports map[string]string) string
	GoType(imports map[string]string) string
	StringValue(interface{}) string
	Type() Type
	TypeString() string
	Validate(ctx context.Context) error
	ValidateValue(interface{}) (interface{}, error)
	ValidateSchema(s Schema, compare bool) error
}

type Directive interface {
	ApplyToSchema(Schema) error
}

type SchemaReplacer interface {
	ReplaceSchema(Schema) (Schema, error)
}

type Nullable interface {
	GetNullable() bool
	SetNullable(bool)
}

func UnmarshalJSON(b []byte) (Schema, error) {
	switch {
	case bytes.Equal(b, []byte("null")):
		return nil, nil
	case bytes.Equal(b, []byte("false")):
		return False(), nil
	case bytes.Equal(b, []byte("true")), bytes.Equal(b, []byte("{}")):
		return &Any{}, nil
	}

	schemaType := gjson.GetBytes(b, "type")
	if !schemaType.Exists() {
		switch {
		case gjson.GetBytes(b, "allOf").Exists():
			var s AllOf
			if err := json.Unmarshal(b, &s); err != nil {
				return nil, err
			}
			return &s, nil
		case gjson.GetBytes(b, "oneOf").Exists():
			var s OneOf
			if err := json.Unmarshal(b, &s); err != nil {
				return nil, err
			}
			return &s, nil
		default:
			var s Any
			if err := json.Unmarshal(b, &s); err != nil {
				return nil, err
			}
			return &s, nil
		}
	}

	var s Schema
	switch Type(gjson.GetBytes(b, "type").String()) {
	case TypeArray:
		s = &Array{}
	case TypeBoolean:
		s = &Boolean{}
	case TypeByteStream:
		s = &ByteStream{}
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
	Name string
	Schema
}

func (n NamedSchema) GoString(imports map[string]string) string {
	buf := bytes.NewBuffer(nil)
	fprintf(buf, "%sNamedSchema{\n", schemaPkg(imports))
	fprintf(buf, "\tName: %q,\n", n.Name)
	fprintf(buf, "\tSchema: %s,\n", indent(n.Schema.GoString(imports)))
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
	case time.Duration:
		return &String{Format: FormatDuration}, nil
	default:
		formatName, _, ok := GetFormatForType(getFullyQualifiedTypeName(reflect.TypeOf(value)))
		if ok {
			return &String{
				Format: formatName,
			}, nil
		}

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
	if s1.Type() == TypeAny {
		return s1, nil
	}

	for i := 1; i < cnt; i++ {
		s2, err := s(i)
		if err != nil {
			return nil, err
		}
		s1 = GetCommonSchema(s1, s2)
		if s1.Type() == TypeAny {
			return s1, nil
		}
	}

	return s1, nil
}

func GetCommonSchema(s1, s2 Schema) Schema {
	switch s1.Type() {
	case TypeAny:
		return s1
	case TypeArray:
		if s2.Type() == TypeNull {
			return s1
		}

		if s2a, ok := s2.(*Array); ok {
			return &Array{
				Items: GetCommonSchema(s1.(*Array).Items, s2a.Items),
			}
		}
	case TypeMap:
		if s2.Type() == TypeNull {
			return s1
		}

		if s2m, ok := s2.(*Map); ok {
			return &Map{
				AdditionalProperties: GetCommonSchema(
					s1.(*Map).GetAdditionalProperties(),
					s2m.GetAdditionalProperties(),
				),
			}
		}
	case TypeObject:
		panic("GetCommonSchema should not be called for objects")
	case TypeInteger:
		if s2.Type() == TypeInteger {
			return s1
		}

		if s2.Type() == TypeNumber {
			return s2
		}
	case TypeNumber:
		if s2.Type() == TypeNumber || s2.Type() == TypeInteger {
			return s1
		}
	case TypeNull:
		if s2.Type() == TypeNull {
			return s1
		}

		if s2.Type() == TypeArray || s2.Type() == TypeMap {
			return s2
		}
	default:
		if s1.Type() == s2.Type() {
			return s1
		}
	}

	return AnyValue()
}
