// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package schema

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/conflowio/conflow/pkg/schema/formats"
)

func init() {
	generateRegisteredFormatTypes()
}

//counterfeiter:generate . Format
type Format interface {
	Type() (reflect.Type, bool)
	StringValue(input interface{}) (string, bool)
	ValidateValue(input string) (interface{}, error)
}

const (
	FormatBinary          = "binary"
	FormatByte            = "byte" // Alias to "binary"
	FormatDefault         = ""
	FormatDate            = "date"
	FormatDateTime        = "date-time"
	FormatDurationRFC3339 = "duration"
	FormatDurationGo      = "duration-go"
	FormatEmail           = "email"
	FormatHostname        = "hostname"
	FormatIDNEmail        = "idn-email"
	FormatIDNHostname     = "idn-hostname"
	FormatIP              = "ip"
	FormatIPCIDR          = "ip-cidr"
	FormatIPv4            = "ipv4"
	FormatIPv4CIDR        = "ipv4-cidr"
	FormatIPv6            = "ipv6"
	FormatIPv6CIDR        = "ipv6-cidr"
	FormatIRI             = "iri"
	FormatIRIReference    = "iri-reference"
	FormatRegex           = "regex"
	FormatTime            = "time"
	FormatURI             = "uri"
	FormatURIReference    = "uri-reference"
	FormatUUID            = "uuid"
)

var registeredFormats = map[string]Format{
	FormatBinary:          formats.Binary{Default: true},
	FormatByte:            formats.Binary{},
	FormatDefault:         formats.String{},
	FormatDate:            formats.Date{},
	FormatDateTime:        formats.DateTime{},
	FormatDurationRFC3339: formats.DurationRFC3339{},
	FormatDurationGo:      formats.DurationGo{},
	FormatEmail:           formats.Email{Default: true},
	FormatHostname:        formats.Hostname{},
	FormatIDNEmail:        formats.Email{},
	FormatIDNHostname:     formats.Hostname{},
	FormatIP:              formats.IP{AllowIPv4: true, AllowIPv6: true, Default: true},
	FormatIPCIDR:          formats.IPCIDR{AllowIPv4: true, AllowIPv6: true, Default: true},
	FormatIPv4:            formats.IP{AllowIPv4: true},
	FormatIPv4CIDR:        formats.IPCIDR{AllowIPv4: true},
	FormatIPv6:            formats.IP{AllowIPv6: true},
	FormatIPv6CIDR:        formats.IPCIDR{AllowIPv6: true},
	FormatIRI:             formats.URI{RequireScheme: true},
	FormatIRIReference:    formats.URI{},
	FormatRegex:           formats.Regex{},
	FormatTime:            formats.Time{},
	FormatURI:             formats.URI{Default: true, RequireScheme: true},
	FormatURIReference:    formats.URI{},
	FormatUUID:            formats.UUID{},
}
var registeredFormatTypes map[string][]string
var registeredFormatsLock sync.Mutex

func generateRegisteredFormatTypes() {
	registeredFormatTypes = map[string][]string{}

	for name, f := range registeredFormats {
		t, def := f.Type()
		fqtn := getFullyQualifiedTypeName(t)
		if def {
			if _, exists := registeredFormatTypes[fqtn]; exists {
				panic(fmt.Errorf("multiple default format registered for %q", fqtn))
			}
			registeredFormatTypes[fqtn] = []string{name}
		}
	}
}

func getFullyQualifiedTypeName(t reflect.Type) string {
	if t.PkgPath() == "" {
		if name := t.Name(); name != "" {
			return t.Name()
		}
		return t.String()
	}
	return fmt.Sprintf("%s.%s", t.PkgPath(), t.Name())
}

func RegisterFormat(name string, format Format) {
	if name == FormatDefault {
		panic(errors.New("the default string format can not be overridden (name can not be empty string)"))
	}

	ftype, def := format.Type()
	fqtn := getFullyQualifiedTypeName(ftype)
	if fqtn == "string" && def {
		panic(errors.New(`the default format for strings can not be overridden (Type() should return: "string", false)`))
	}

	registeredFormatsLock.Lock()
	defer registeredFormatsLock.Unlock()

	if f, exists := registeredFormats[name]; exists {
		unregisterFormat(name, f)
	}

	registeredFormats[name] = format
	if def {
		registeredFormatTypes[fqtn] = append(registeredFormatTypes[fqtn], name)
	}
}

func UnregisterFormat(name string) {
	if name == FormatDefault {
		panic(errors.New("the default string format can not be unregistered"))
	}
	registeredFormatsLock.Lock()
	defer registeredFormatsLock.Unlock()

	f, ok := registeredFormats[name]
	if !ok {
		return
	}

	unregisterFormat(name, f)
}

func unregisterFormat(name string, format Format) {
	delete(registeredFormats, name)

	ftype, def := format.Type()
	fqtn := getFullyQualifiedTypeName(ftype)
	if def {
		for i, n := range registeredFormatTypes[fqtn] {
			if n == name {
				registeredFormatTypes[fqtn] = append(registeredFormatTypes[fqtn][0:i], registeredFormatTypes[fqtn][i+1:]...)
				break
			}
		}
	}
}

func GetFormat(name string) (Format, bool) {
	registeredFormatsLock.Lock()
	defer registeredFormatsLock.Unlock()

	f, ok := registeredFormats[name]
	return f, ok
}

func GetFormatForType(path string) (string, Format, bool) {
	registeredFormatsLock.Lock()
	defer registeredFormatsLock.Unlock()

	names, ok := registeredFormatTypes[path]
	if !ok || len(names) == 0 {
		return "", nil, false
	}

	name := names[len(names)-1]
	return name, registeredFormats[name], true
}
