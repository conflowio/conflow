// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package blocks

import (
	"io"

	"github.com/opsidian/basil/basil"
)

func NewStream(id basil.ID, stream io.ReadCloser) *Stream {
	return &Stream{
		id:     id,
		Stream: stream,
	}
}

//go:generate basil generate
type Stream struct {
	id     basil.ID `basil:"id"`
	Stream io.ReadCloser
}

func (s Stream) ID() basil.ID {
	return s.id
}
