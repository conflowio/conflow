package zerolog

import (
	"github.com/opsidian/basil/basil"
	"github.com/rs/zerolog"
)

type ObjectMarshalerWrapper struct {
	obj basil.LogObjectMarshaler
}

func (o *ObjectMarshalerWrapper) MarshalZerologObject(e *zerolog.Event) {
	o.obj.MarshalLogObject(&Event{e: e})
}
