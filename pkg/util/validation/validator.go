// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package validation

import (
	"context"
	"fmt"
	"reflect"
	"strings"
)

type Validator interface {
	Validate(ctx context.Context) error
}

type ValidatorFunc func(ctx context.Context) error

func (f ValidatorFunc) Validate(ctx context.Context) error {
	return f(ctx)
}

func Validate(ctx context.Context, target interface{}) error {
	ve := &Error{}
	v := reflect.ValueOf(target)

	switch v.Kind() {
	case reflect.Pointer, reflect.Interface:
		if v.IsNil() {
			return nil
		}
		if err := Validate(ctx, v.Elem().Interface()); err != nil {
			ve.AddError(err)
		}
	case reflect.Map:
		if v.IsNil() {
			return nil
		}

		iter := v.MapRange()
		for iter.Next() {
			k := iter.Key()
			v := iter.Value()
			if err := Validate(ctx, v.Interface()); err != nil {
				ve.AddFieldError(fmt.Sprintf("[%q]", k.String()), err)
			}
		}
	case reflect.Slice:
		if v.IsNil() {
			return nil
		}
		fallthrough
	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			if err := Validate(ctx, v.Index(i).Interface()); err != nil {
				ve.AddFieldError(fmt.Sprintf("[%d]", i), err)
			}
		}
	case reflect.Struct:
		t := v.Type()
		for i := 0; i < t.NumField(); i++ {
			ft := t.Field(i)
			if !ft.IsExported() {
				continue
			}

			name := ft.Name
			if jsonTag, ok := ft.Tag.Lookup("json"); ok {
				parts := strings.Split(jsonTag, ",")
				// TODO: ignore "-"
				if parts[0] != "" && parts[0] != "-" {
					name = parts[0]
				}
			}

			if err := Validate(ctx, v.Field(i).Interface()); err != nil {
				ve.AddFieldError(name, err)
			}
		}
	}

	if v, ok := target.(Validator); ok {
		if err := v.Validate(ctx); err != nil {
			ve.AddError(err)
		}
	}

	return ve.ErrOrNil()
}

func ValidateObject(ctx context.Context, validators ...Validator) error {
	ve := &Error{}
	for _, v := range validators {
		if err := v.Validate(ctx); err != nil {
			ve.AddError(err)
		}
	}

	return ve.ErrOrNil()
}

func ValidateField(name string, v Validator) Validator {
	return ValidatorFunc(func(ctx context.Context) error {
		if v != nil && !reflect.ValueOf(v).IsNil() {
			if err := v.Validate(ctx); err != nil {
				return NewFieldError(name, err)
			}
		}
		return nil
	})
}

func ValidateArray[T Validator](name string, v []T) Validator {
	return ValidatorFunc(func(ctx context.Context) error {
		ve := &Error{}
		for i, e := range v {
			if err := e.Validate(ctx); err != nil {
				ve.AddFieldError(fmt.Sprintf("%s[%d]", name, i), err)
			}
		}
		return ve.ErrOrNil()
	})
}

func ValidateMap[T Validator](name string, v map[string]T) Validator {
	return ValidatorFunc(func(ctx context.Context) error {
		ve := &Error{}
		for k, e := range v {
			if err := e.Validate(ctx); err != nil {
				ve.AddFieldError(fmt.Sprintf("%s[%q]", name, k), err)
			}
		}
		return ve.ErrOrNil()
	})
}
