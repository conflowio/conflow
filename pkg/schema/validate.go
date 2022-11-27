// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package schema

import (
	"context"
	"fmt"
	"reflect"

	"github.com/conflowio/conflow/pkg/util/validation"
)

func validateCommonFields[T any](s Schema, constValue interface{}, defaultValue interface{}, enumValues []T) validation.ValidatorFunc {
	return func(ctx context.Context) error {
		ve := &validation.Error{}
		if err := validateValue[T](s, constValue); err != nil {
			ve.AddFieldError("const", err)
		}

		if err := validateValue[T](s, defaultValue); err != nil {
			ve.AddFieldError("default", err)
		}

		for i, enum := range enumValues {
			if _, err := s.ValidateValue(enum); err != nil {
				ve.AddFieldError(fmt.Sprintf("enum[%d]", i), err)
			}
		}

		for i, example := range s.GetExamples() {
			if _, err := s.ValidateValue(example); err != nil {
				ve.AddFieldError(fmt.Sprintf("examples[%d]", i), err)
			}
		}

		return ve.ErrOrNil()
	}
}

func validateValue[T any](s Schema, value interface{}) error {
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
		return err
	}
	return nil
}
