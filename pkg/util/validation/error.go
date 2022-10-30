// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package validation

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Error struct {
	errors []error
}

func NewError(errs ...error) *Error {
	e := &Error{}
	for _, err := range errs {
		e.AddError(err)
	}
	return e
}

func (e *Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Errors []error `json:"errors"`
	}{
		Errors: e.errors,
	})
}

func (e *Error) AddError(err error) {
	e.addError(err)
}

func (e *Error) AddErrorf(format string, a ...interface{}) {
	e.errors = append(e.errors, fmt.Errorf(format, a...))
}

func (e *Error) AddFieldError(field string, err error) {
	e.addError(NewFieldError(field, err))
}

func (e *Error) AddFieldErrorf(field, format string, a ...interface{}) {
	e.errors = append(e.errors, NewFieldErrorf(field, format, a...))
}

func (e *Error) addError(err error) {
	switch t := err.(type) {
	case *Error:
		for _, ee := range t.errors {
			e.AddError(ee)
		}
	default:
		e.errors = append(e.errors, t)
	}
}

func (e *Error) Error() string {
	if len(e.errors) == 1 {
		return e.errors[0].Error()
	}

	var sb strings.Builder
	for i, err := range e.errors {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(err.Error())
	}
	return sb.String()
}

func (e *Error) ErrOrNil() error {
	switch len(e.errors) {
	case 0:
		return nil
	case 1:
		return e.errors[0]
	default:
		return e
	}
}
