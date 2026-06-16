// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package bind

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/conflowio/conflow/pkg/schema"
	"github.com/conflowio/conflow/pkg/util"
	"github.com/conflowio/conflow/pkg/values"
)

func BindValue(s schema.Schema, value interface{}) (interface{}, error) {
	if value == nil {
		return nil, nil
	}

	debugBind(s.Type(), value)

	switch s.Type() {
	case schema.TypeBoolean, schema.TypeInteger, schema.TypeNumber, schema.TypeString:
		return value, nil
	case schema.TypeAny:
		return bindAnyValue(value)
	case schema.TypeArray:
		return bindArray(s.(*schema.Array), value)
	case schema.TypeMap:
		return bindMap(s.(*schema.Map), value)
	case schema.TypeObject:
		return bindObject(s.(*schema.Object), value)
	default:
		return nil, fmt.Errorf("bind: unsupported schema type %s", s.Type())
	}
}

func bindAnyValue(value interface{}) (interface{}, error) {
	switch v := value.(type) {
	case bool, int64, float64, string:
		return v, nil
	case []interface{}:
		return bindArray(&schema.Array{Items: schema.AnyValue()}, v)
	case map[string]interface{}:
		return bindMap(&schema.Map{}, v)
	default:
		if isValuesList(value) {
			return bindImmutableList(value, schema.AnyValue())
		}
		if isValuesMap(value) {
			return bindImmutableMap(value, schema.AnyValue())
		}

		rv := reflect.ValueOf(value)
		if rv.Kind() == reflect.Slice {
			return bindSliceValue(rv, schema.AnyValue())
		}
		if rv.Kind() == reflect.Map {
			return bindMapValue(rv, schema.AnyValue())
		}

		return value, nil
	}
}

func bindArray(a *schema.Array, value interface{}) (interface{}, error) {
	if frozen, ok := freezeListBuilder(value); ok {
		return bindArray(a, frozen)
	}

	if isValuesList(value) {
		return bindImmutableList(value, a.Items)
	}

	if items, ok := value.([]interface{}); ok {
		builder := values.NewListBuilder[interface{}]()
		for _, item := range items {
			bound, err := BindValue(a.Items, item)
			if err != nil {
				return nil, err
			}
			builder.Append(bound)
		}
		return builder.Freeze(), nil
	}

	return bindSliceValue(reflect.ValueOf(value), a.Items)
}

func bindSliceValue(rv reflect.Value, itemSchema schema.Schema) (interface{}, error) {
	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return nil, nil
		}
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Slice {
		return nil, fmt.Errorf("bind: expected slice, got %T", rv.Interface())
	}

	cloner, err := sliceElementCloner(itemSchema)
	if err != nil {
		return nil, err
	}

	result := reflect.MakeSlice(rv.Type(), rv.Len(), rv.Len())
	for i := 0; i < rv.Len(); i++ {
		cloned, err := cloner(rv.Index(i).Interface())
		if err != nil {
			return nil, err
		}
		result.Index(i).Set(reflect.ValueOf(cloned))
	}

	return result.Interface(), nil
}

func bindMap(m *schema.Map, value interface{}) (interface{}, error) {
	if frozen, ok := freezeMapBuilder(value); ok {
		return bindMap(m, frozen)
	}

	if isValuesMap(value) {
		return bindImmutableMap(value, m.GetAdditionalProperties())
	}

	if goMap, ok := value.(map[string]interface{}); ok {
		builder := values.NewMapBuilder[string, interface{}]()
		for k, v := range goMap {
			bound, err := BindValue(m.GetAdditionalProperties(), v)
			if err != nil {
				return nil, err
			}
			builder.Set(k, bound)
		}
		return builder.Freeze(), nil
	}

	return bindMapValue(reflect.ValueOf(value), m.GetAdditionalProperties())
}

func bindMapValue(rv reflect.Value, valueSchema schema.Schema) (interface{}, error) {
	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return nil, nil
		}
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Map {
		return nil, fmt.Errorf("bind: expected map, got %T", rv.Interface())
	}

	cloner, err := mapValueCloner(valueSchema)
	if err != nil {
		return nil, err
	}

	result := reflect.MakeMap(rv.Type())
	iter := rv.MapRange()
	for iter.Next() {
		cloned, err := cloner(iter.Value().Interface())
		if err != nil {
			return nil, err
		}
		result.SetMapIndex(iter.Key(), reflect.ValueOf(cloned))
	}

	return result.Interface(), nil
}

func bindObject(o *schema.Object, value interface{}) (interface{}, error) {
	goMap, ok := value.(map[string]interface{})
	if !ok {
		rv := reflect.ValueOf(value)
		if rv.Kind() == reflect.Pointer && !rv.IsNil() {
			rv = rv.Elem()
		}
		if rv.Kind() == reflect.Struct {
			return bindStructObject(o, rv)
		}
		return nil, fmt.Errorf("bind: expected object map, got %T", value)
	}

	result := make(map[string]interface{}, len(goMap))
	for k, v := range goMap {
		propSchema, ok := o.Properties[k]
		if !ok {
			propSchema = schema.AnyValue()
		}
		bound, err := BindValue(propSchema, v)
		if err != nil {
			return nil, err
		}
		result[k] = bound
	}

	return result, nil
}

func bindStructObject(o *schema.Object, rv reflect.Value) (interface{}, error) {
	result := make(map[string]interface{}, len(o.Properties))
	for jsonName, propSchema := range o.Properties {
		fieldName := o.FieldName(jsonName)
		field := rv.FieldByName(fieldName)
		if !field.IsValid() {
			continue
		}
		if field.Kind() == reflect.Pointer && field.IsNil() {
			result[jsonName] = nil
			continue
		}
		fieldValue := field.Interface()
		if field.Kind() == reflect.Pointer {
			fieldValue = field.Elem().Interface()
		}
		bound, err := BindValue(propSchema, fieldValue)
		if err != nil {
			return nil, err
		}
		result[jsonName] = bound
	}
	return result, nil
}

type valueCloner func(interface{}) (interface{}, error)

func sliceElementCloner(itemSchema schema.Schema) (valueCloner, error) {
	switch itemSchema.Type() {
	case schema.TypeBoolean:
		return cloneScalar[bool](itemSchema)
	case schema.TypeInteger:
		return cloneScalar[int64](itemSchema)
	case schema.TypeNumber:
		return cloneScalar[float64](itemSchema)
	case schema.TypeString:
		return cloneScalar[string](itemSchema)
	default:
		return func(v interface{}) (interface{}, error) {
			return BindValue(itemSchema, v)
		}, nil
	}
}

func cloneScalar[T ~bool | ~int64 | ~float64 | ~string](itemSchema schema.Schema) (valueCloner, error) {
	return func(v interface{}) (interface{}, error) {
		if typed, ok := v.(T); ok {
			return util.CloneValue(typed), nil
		}
		return BindValue(itemSchema, v)
	}, nil
}

func mapValueCloner(valueSchema schema.Schema) (valueCloner, error) {
	return sliceElementCloner(valueSchema)
}

func bindImmutableList(value interface{}, itemSchema schema.Schema) (interface{}, error) {
	rv := reflect.ValueOf(value)
	changed := false
	boundElems := make([]interface{}, rv.MethodByName("Len").Call(nil)[0].Int())

	for i := range boundElems {
		elem := rv.MethodByName("At").Call([]reflect.Value{reflect.ValueOf(i)})[0].Interface()
		bound, err := BindValue(itemSchema, elem)
		if err != nil {
			return nil, err
		}
		if !reflect.DeepEqual(bound, elem) {
			changed = true
		}
		boundElems[i] = bound
	}

	if !changed {
		return value, nil
	}

	list, err := values.FromInterfaceSlice(boundElems)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func bindImmutableMap(value interface{}, valueSchema schema.Schema) (interface{}, error) {
	rv := reflect.ValueOf(value)
	keysValue := rv.MethodByName("Keys").Call(nil)[0]

	changed := false
	builder := values.NewMapBuilder[string, interface{}]()
	for i := 0; i < keysValue.Len(); i++ {
		key := keysValue.Index(i).Interface()
		keyStr, ok := key.(string)
		if !ok {
			return nil, fmt.Errorf("bind: immutable map key type %T not supported", key)
		}

		getResult := rv.MethodByName("Get").Call([]reflect.Value{reflect.ValueOf(keyStr)})
		if len(getResult) < 2 || !getResult[1].Bool() {
			continue
		}
		elem := getResult[0].Interface()
		bound, err := BindValue(valueSchema, elem)
		if err != nil {
			return nil, err
		}
		if !reflect.DeepEqual(bound, elem) {
			changed = true
		}
		builder.Set(keyStr, bound)
	}

	if !changed {
		return value, nil
	}

	return builder.Freeze(), nil
}

func isValuesList(value interface{}) bool {
	return isValuesType(value, "List[")
}

func isValuesMap(value interface{}) bool {
	return isValuesType(value, "Map[")
}

func isValuesType(value interface{}, prefix string) bool {
	if value == nil {
		return false
	}
	t := reflect.TypeOf(value)
	if t.Kind() != reflect.Pointer || t.Elem().Kind() != reflect.Struct {
		return false
	}
	elem := t.Elem()
	return elem.PkgPath() == "github.com/conflowio/conflow/pkg/values" && strings.HasPrefix(elem.Name(), prefix)
}

func freezeListBuilder(value interface{}) (interface{}, bool) {
	return freezeBuilder(value, "ListBuilder[")
}

func freezeMapBuilder(value interface{}) (interface{}, bool) {
	return freezeBuilder(value, "MapBuilder[")
}

func freezeBuilder(value interface{}, prefix string) (interface{}, bool) {
	if value == nil {
		return nil, false
	}
	t := reflect.TypeOf(value)
	if t.Kind() != reflect.Pointer || t.Elem().Kind() != reflect.Struct {
		return nil, false
	}
	elem := t.Elem()
	if elem.PkgPath() != "github.com/conflowio/conflow/pkg/values" || !strings.HasPrefix(elem.Name(), prefix) {
		return nil, false
	}
	frozen := reflect.ValueOf(value).MethodByName("Freeze").Call(nil)
	if len(frozen) == 0 {
		return nil, false
	}
	return frozen[0].Interface(), true
}
