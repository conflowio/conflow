// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives

import (
	"github.com/opsidian/basil/basil"
)

//go:generate basil generate
type Run struct {
	id       basil.ID `basil:"id"`
	when     bool     `basil:"value"`
	triggers []interface{}
}

func (r Run) ID() basil.ID {
	return r.id
}

func (r Run) ApplyDirective(blockCtx basil.BlockContext, container basil.BlockContainer) error {
	if !r.when {
		container.Skip()
	}

	trigger := container.Trigger()
	if trigger == "" {
		return nil
	}
	for _, id := range r.triggers {
		if trigger == id {
			return nil
		}
	}

	container.Skip()

	return nil
}
