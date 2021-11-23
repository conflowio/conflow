// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package blocks

import (
	"io"

	"github.com/conflowio/conflow/conflow"
)

func NewStream(id conflow.ID, stream io.ReadCloser) *Stream {
	return &Stream{
		id:     id,
		Stream: stream,
	}
}

// @block "configuration"
type Stream struct {
	// @id
	id     conflow.ID
	Stream io.ReadCloser
}

func (s Stream) ID() conflow.ID {
	return s.id
}
