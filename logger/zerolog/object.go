// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

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
