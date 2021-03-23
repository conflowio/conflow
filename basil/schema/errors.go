// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package schema

import (
	"fmt"
	"strings"
)

type ValidationError struct {
	Errors []error `json:"errors"`
}

func (v *ValidationError) AddError(field string, err error) {
	v.Errors = append(v.Errors, FieldError{
		Field: field,
		Err:   err,
	})
}

func (v *ValidationError) AddErrorf(field, format string, a ...interface{}) {
	v.Errors = append(v.Errors, FieldError{
		Field: field,
		Err:   fmt.Errorf(format, a...),
	})
}

func (v ValidationError) Error() string {
	var sb strings.Builder
	for _, err := range v.Errors {
		sb.WriteString(" * ")
		sb.WriteString(err.Error())
		sb.WriteRune('\n')
	}
	return sb.String()
}

func (v ValidationError) ErrOrNil() error {
	switch len(v.Errors) {
	case 0:
		return nil
	case 1:
		return v.Errors[0]
	default:
		return v
	}
}

func NewFieldError(field string, err error) FieldError {
	return FieldError{
		Field: field,
		Err:   err,
	}
}

func NewFieldErrorf(field, format string, a ...interface{}) FieldError {
	return FieldError{
		Field: field,
		Err:   fmt.Errorf(format, a...),
	}
}

type FieldError struct {
	Field string `json:"field"`
	Err   error  `json:"error"`
}

func (f FieldError) Error() string {
	return fmt.Sprintf("%s: %s", f.Field, f.Err.Error())
}

type typeError string

func (t typeError) Error() string {
	return string(t)
}

func typeErrorf(format string, a ...interface{}) typeError {
	return typeError(fmt.Sprintf(format, a...))
}
