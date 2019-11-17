package blocks

import (
	"io"

	"github.com/opsidian/basil/basil"
)

func NewStream(id basil.ID, reader io.ReadCloser) *Stream {
	return &Stream{
		id:     id,
		Reader: reader,
	}
}

//go:generate basil generate
type Stream struct {
	id     basil.ID      `basil:"id"`
	Reader io.ReadCloser `basil:"name=stream"`
}

func (s Stream) ID() basil.ID {
	return s.id
}
