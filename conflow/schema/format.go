// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package schema

import (
	"github.com/conflowio/conflow/conflow/schema/formats"
)

type Format interface {
	Parse(input string) (interface{}, error)
	Format(input interface{}) string
}

var formatCheckers = map[string]Format{
	"date":          formats.Date{},
	"date-time":     formats.DateTime{},
	"duration":      formats.Duration{},
	"email":         formats.Email{},
	"hostname":      formats.Hostname{},
	"idn-email":     formats.Email{},
	"idn-hostname":  formats.Hostname{},
	"ip":            formats.IP{},
	"ip-cidr":       formats.IPCIDR{},
	"ipv4":          formats.IPv4{},
	"ipv4-cidr":     formats.IPv4CIDR{},
	"ipv6":          formats.IPv6{},
	"ipv6-cidr":     formats.IPv6CIDR{},
	"iri":           formats.URI{},
	"iri-reference": formats.URIReference{},
	"regex":         formats.Regex{},
	"time":          formats.Time{},
	"uri":           formats.URI{},
	"uri-reference": formats.URIReference{},
	"uri-template":  formats.URIReference{},
	"uuid":          formats.UUID{},
}
