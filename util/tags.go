// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package util

import (
	"fmt"
	"reflect"
	"strings"
)

// FieldTags contains the meta tags for a field
type FieldTags map[string]string

// StructTags contains the meta tags for a struct
type StructTags map[string]FieldTags

func (f FieldTags) Get(name string) (string, bool) {
	val, ok := f[strings.ToLower(name)]
	return val, ok
}

func (f FieldTags) GetBool(name string) bool {
	val, _ := f[strings.ToLower(name)]
	return val == "true"
}

func (f FieldTags) GetWithDefault(name string, def string) string {
	if val, ok := f[strings.ToLower(name)]; ok {
		return val
	}

	return def
}

func (f FieldTags) Keys() []string {
	keys := make([]string, 0, len(f))
	for k := range f {
		keys = append(keys, k)
	}

	return keys
}

// GetTags scans the given struct and returns with the associated values for all fields for the given meta tag
// The tag can contain "key:value" pairs, or simple values separated by a comma.
// Simple values will be returned as "value":"true" entry in the result.
// Whitespaces are ignored.
//
// Example:
//  field string `flow:"field1:value1,field2:value2,values3"`
func GetTags(obj interface{}, tagName string) StructTags {
	objType := reflect.TypeOf(obj)
	for objType.Kind() == reflect.Ptr {
		objType = objType.Elem()
	}
	if objType.Kind() != reflect.Struct {
		panic(fmt.Sprintf("GetTags was called with %s type instead of a struct", objType.Kind()))
	}
	if objType.NumField() == 0 {
		return nil
	}
	tags := make(StructTags)
	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		tags[field.Name] = ParseFieldTag(field.Tag, tagName, field.Name)
	}

	return tags
}

// ParseFieldTag parses a field's tag value, looks for the given tag name and returns with key-value pairs
func ParseFieldTag(tag reflect.StructTag, tagName string, fieldName string) FieldTags {
	tagVal := strings.TrimSpace(tag.Get(tagName))
	if tagVal == "" {
		return nil
	}
	fields := make(FieldTags)
	for _, tagField := range strings.Split(tagVal, ",") {
		keyValue := strings.SplitN(tagField, "=", 2)
		key := strings.ToLower(strings.TrimSpace(keyValue[0]))
		if len(keyValue) == 2 {
			if key == "" {
				panic(fmt.Sprintf("Struct tag key can not be empty for field %s", fieldName))
			}
			fields[key] = strings.TrimSpace(keyValue[1])
		} else {
			if key != "" {
				fields[key] = "true"
			}
		}
	}
	return fields
}

// FilterFieldsByTags returns with the list of fields which has the given tag with the given values
func FilterFieldsByTags(tags StructTags, filterTag string, filterValues ...string) []string {
	if len(filterValues) == 0 {
		panic("at least one tag value should be passed")
	}

	var res []string
	for field, fieldTags := range tags {
		for tag, value := range fieldTags {
			if tag == filterTag {
				for _, filterValue := range filterValues {
					if filterValue == value {
						res = append(res, field)
					}
				}
			}
		}
	}
	return res
}
