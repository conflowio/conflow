// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package schema

import (
	"fmt"
	"reflect"
)

type Validatable interface {
	Validate(ctx *Context) error
}

func ValidateAll(ctx *Context, validators ...func(ctx *Context) error) error {
	for _, v := range validators {
		if err := v(ctx); err != nil {
			return err
		}
	}
	return nil
}

func Validate(name string, v Validatable) func(*Context) error {
	return func(ctx *Context) error {
		if v != nil && !reflect.ValueOf(v).IsNil() {
			if err := v.Validate(ctx); err != nil {
				return fmt.Errorf("%s: %w", name, err)
			}
		}
		return nil
	}
}

func ValidateArray[T Validatable](name string, v []T) func(*Context) error {
	return func(ctx *Context) error {
		for i, e := range v {
			if err := e.Validate(ctx); err != nil {
				return fmt.Errorf("%s.%d: %w", name, i, err)
			}
		}
		return nil
	}
}

func ValidateMap[T Validatable](name string, v map[string]T) func(*Context) error {
	return func(ctx *Context) error {
		for k, e := range v {
			if err := e.Validate(ctx); err != nil {
				return fmt.Errorf("%s.%s: %w", name, k, err)
			}
		}
		return nil
	}
}

func validateValue[T any](s Schema, name string, value interface{}) error {
	switch v := value.(type) {
	case *T:
		if v == nil {
			return nil
		}
		value = *v
	default:
		if v == nil || reflect.ValueOf(value).IsNil() {
			return nil
		}
	}

	if _, err := s.ValidateValue(value); err != nil {
		return fmt.Errorf("%s: %w", name, err)
	}
	return nil
}

func validateCommonFields[T any](s Schema, constValue interface{}, defaultValue interface{}, enumValues []T) func(*Context) error {
	return func(ctx *Context) error {
		if err := validateValue[T](s, "const", constValue); err != nil {
			return err
		}

		if err := validateValue[T](s, "default", defaultValue); err != nil {
			return err
		}

		for i, enum := range enumValues {
			if _, err := s.ValidateValue(enum); err != nil {
				return fmt.Errorf("enum.%d is invalid: %w", i, err)
			}
		}

		for i, example := range s.GetExamples() {
			if _, err := s.ValidateValue(example); err != nil {
				return fmt.Errorf("examples.%d is invalid: %w", i, err)
			}
		}

		return nil
	}
}
