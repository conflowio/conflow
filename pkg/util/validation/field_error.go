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

func NewFieldError(field string, err error) error {
	switch t := err.(type) {
	case *Error:
		ve := &Error{}
		for _, e := range t.errors {
			ve.errors = append(ve.errors, NewFieldError(field, e))
		}
		return ve
	case FieldError:
		if strings.HasPrefix(t.field, "[") {
			field = fmt.Sprintf("%s%s", field, t.field)
		} else {
			field = fmt.Sprintf("%s.%s", field, t.field)
		}

		return FieldError{
			field: field,
			err:   t.err,
		}
	default:
		return FieldError{
			field: field,
			err:   err,
		}
	}

}

func NewFieldErrorf(field, format string, a ...interface{}) FieldError {
	return FieldError{
		field: field,
		err:   fmt.Errorf(format, a...),
	}
}

type FieldError struct {
	field string
	err   error
}

func (f FieldError) Field() string {
	return f.field
}

func (f FieldError) Err() error {
	return f.err
}

func (f FieldError) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Field string `json:"field"`
		Err   error  `json:"error"`
	}{
		Field: f.field,
		Err:   f.err,
	})
}

func (f FieldError) Error() string {
	return fmt.Sprintf("%s: %s", f.field, f.err.Error())
}

func (f FieldError) TransformError(tr func(path string, err error) error) error {
	return tr(f.field, f.err)
}
