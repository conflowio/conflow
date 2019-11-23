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
