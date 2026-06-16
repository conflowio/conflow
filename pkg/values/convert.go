// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package values

import (
	"fmt"
	"reflect"
	"strings"
)

func FromSlice[T any](s []T) *List[T] {
	return NewListFromSlice(s)
}

func FromInterfaceSlice(s []interface{}) (*List[interface{}], error) {
	builder := NewListBuilder[interface{}]()
	for _, v := range s {
		builder.Append(v)
	}
	return builder.Freeze(), nil
}

func FromGoMap[K comparable, V any](m map[K]V) *Map[K, V] {
	return NewMapFromGoMap(m)
}

func FromStringInterfaceMap(m map[string]interface{}) (*Map[string, interface{}], error) {
	return NewMapFromGoMap(m), nil
}

func IsImmutableCollection(value interface{}) bool {
	return IsImmutableList(value) || IsImmutableMap(value)
}

func IsImmutableList(value interface{}) bool {
	if value == nil {
		return false
	}
	t := reflect.TypeOf(value)
	if t.Kind() != reflect.Pointer || t.Elem().Kind() != reflect.Struct {
		return false
	}
	elem := t.Elem()
	return elem.PkgPath() == "github.com/conflowio/conflow/pkg/values" &&
		strings.HasPrefix(elem.Name(), "List[")
}

func IsImmutableMap(value interface{}) bool {
	if value == nil {
		return false
	}
	t := reflect.TypeOf(value)
	if t.Kind() != reflect.Pointer || t.Elem().Kind() != reflect.Struct {
		return false
	}
	elem := t.Elem()
	return elem.PkgPath() == "github.com/conflowio/conflow/pkg/values" &&
		strings.HasPrefix(elem.Name(), "Map[")
}

func AsInterfaceSlice(value interface{}) ([]interface{}, error) {
	switch v := value.(type) {
	case nil:
		return nil, nil
	case []interface{}:
		return v, nil
	case *List[interface{}]:
		return v.Elems(), nil
	default:
		if listElems, ok := listElemsAsInterface(value); ok {
			return listElems, nil
		}
		return nil, fmt.Errorf("values: expected slice or immutable list, got %T", value)
	}
}

func AsStringInterfaceMap(value interface{}) (map[string]interface{}, error) {
	switch v := value.(type) {
	case nil:
		return nil, nil
	case map[string]interface{}:
		return v, nil
	case *Map[string, interface{}]:
		return v.GoMap(), nil
	default:
		if goMap, ok := mapGoMapAsInterface(value); ok {
			return goMap, nil
		}
		return nil, fmt.Errorf("values: expected map or immutable map, got %T", value)
	}
}

func listElemsAsInterface(value interface{}) ([]interface{}, bool) {
	if !IsImmutableList(value) {
		return nil, false
	}
	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Pointer {
		return nil, false
	}
	lenResult := rv.MethodByName("Len").Call(nil)
	if len(lenResult) != 1 {
		return nil, false
	}
	length := int(lenResult[0].Int())
	elems := make([]interface{}, length)
	for i := 0; i < length; i++ {
		atResult := rv.MethodByName("At").Call([]reflect.Value{reflect.ValueOf(i)})
		if len(atResult) != 1 {
			return nil, false
		}
		elems[i] = atResult[0].Interface()
	}
	return elems, true
}

func mapGoMapAsInterface(value interface{}) (map[string]interface{}, bool) {
	if !IsImmutableMap(value) {
		return nil, false
	}
	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Pointer {
		return nil, false
	}
	keysResult := rv.MethodByName("Keys").Call(nil)
	if len(keysResult) != 1 {
		return nil, false
	}
	keysValue := keysResult[0]
	goMap := make(map[string]interface{}, keysValue.Len())
	for i := 0; i < keysValue.Len(); i++ {
		key := keysValue.Index(i).Interface()
		keyStr, ok := key.(string)
		if !ok {
			return nil, false
		}
		getResult := rv.MethodByName("Get").Call([]reflect.Value{reflect.ValueOf(keyStr)})
		if len(getResult) < 2 || !getResult[1].Bool() {
			continue
		}
		goMap[keyStr] = getResult[0].Interface()
	}
	return goMap, true
}
