// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives

import "github.com/opsidian/basil/basil"

//go:generate basil generate
type Bug struct {
	id          basil.ID `basil:"id"`
	description string   `basil:"value"`
}

func (b Bug) ID() basil.ID {
	return b.id
}

func (b Bug) ApplyDirective(blockCtx basil.BlockContext, container basil.BlockContainer) error {
	return nil
}
