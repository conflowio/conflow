// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives

import (
	"fmt"

	"github.com/conflowio/conflow/pkg/conflow"
	"github.com/conflowio/conflow/pkg/schema"
)

// @block "directive"
type Format struct {
	// @id
	id conflow.ID
	// @value
	// @required
	// @enum [
	//   "date",
	//   "date-time",
	//   "duration",
	//   "email",
	//   "hostname",
	//   "idn-email",
	//   "idn-hostname",
	//   "ip",
	//   "ip-cidr",
	//   "ipv4",
	//   "ipv4-cidr",
	//   "ipv6",
	//   "ipv6-cidr",
	//   "iri",
	//   "iri-reference",
	//   "regex",
	//   "time",
	//   "uri",
	//   "uri-reference",
	//   "uri-template",
	//   "uuid",
	// ]
	value string
}

func (f *Format) ID() conflow.ID {
	return f.id
}

func (f *Format) ApplyToSchema(s schema.Schema) error {
	switch st := s.(type) {
	case *schema.String:
		st.Format = f.value
		return nil
	default:
		return fmt.Errorf("format directive can not be applied to %T", s)
	}
}
