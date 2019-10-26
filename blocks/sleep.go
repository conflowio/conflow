// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package blocks

import (
	"errors"
	"time"

	"github.com/opsidian/basil/basil"
)

//go:generate basil generate
type Sleep struct {
	id       basil.ID      `basil:"id"`
	duration time.Duration `basil:"value,required"`
}

func (s *Sleep) ID() basil.ID {
	return s.id
}

func (s *Sleep) Main(blockCtx basil.BlockContext) error {
	select {
	case <-time.After(s.duration):
		return nil
	case <-blockCtx.Done():
		return errors.New("aborted")
	}
}
