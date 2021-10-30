// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package schema

import (
	"errors"
)

func False() Schema {
	return falseSchema
}

var falseSchema = &falseImpl{}

type falseImpl struct {
	emptyMetadata
}

func (f falseImpl) AssignValue(_ map[string]string, _, _ string) string {
	panic("AssignValue should not be called on a false schema")
}

func (f falseImpl) CompareValues(_, _ interface{}) int {
	return 0
}

func (f falseImpl) Copy() Schema {
	return f
}

func (f falseImpl) DefaultValue() interface{} {
	return nil
}

func (f falseImpl) MarshalJSON() ([]byte, error) {
	return []byte("false"), nil
}

func (f falseImpl) GoString() string {
	return "schema.False()"
}

func (f falseImpl) GoType(_ map[string]string) string {
	panic("GoType should not be called on a false schema")
}

func (f falseImpl) StringValue(interface{}) string {
	return ""
}

func (f falseImpl) Type() Type {
	return TypeFalse
}

func (f falseImpl) TypeString() string {
	return ""
}

func (f falseImpl) ValidateSchema(Schema, bool) error {
	return errors.New("no value is allowed")
}

func (f falseImpl) ValidateValue(interface{}) error {
	return errors.New("no value is allowed")
}
