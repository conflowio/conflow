// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package fixtures

import (
	"time"

	"github.com/conflowio/conflow/conflow"
)

// @block "configuration"
type Block struct {
	// @id
	IDField         conflow.ID
	FieldInterface  interface{}
	FieldArray      []interface{}
	FieldBool       bool
	FieldFloat      float64
	FieldIdentifier conflow.ID
	FieldInteger    int64
	FieldMap        map[string]interface{}
	// @types ["integer", "number"]
	FieldNumber       interface{}
	FieldString       string
	FieldStringArray  []string
	FieldTime         time.Time
	FieldTimeDuration time.Duration
}

func (t *Block) ID() conflow.ID {
	return t.IDField
}
