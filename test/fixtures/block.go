// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package fixtures

import (
	"time"

	"github.com/opsidian/basil/basil/variable"

	"github.com/opsidian/basil/basil"
)

//go:generate basil generate
type Block struct {
	IDField           basil.ID `basil:"id"`
	FieldInterface    interface{}
	FieldArray        []interface{}
	FieldBasic        *variable.Basic
	FieldBool         bool
	FieldFloat        float64
	FieldIdentifier   basil.ID
	FieldInteger      int64
	FieldMap          map[string]interface{}
	FieldNumber       *variable.Number
	FieldString       string
	FieldStringArray  []string
	FieldTime         time.Time
	FieldTimeDuration time.Duration
}

func (t *Block) ID() basil.ID {
	return t.IDField
}
