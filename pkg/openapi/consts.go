// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package openapi

const (
	CharsetUTF8 = "charset=UTF-8"

	ContentTypeApplicationForm            = "application/x-www-form-urlencoded"
	ContentTypeApplicationJSON            = "application/json"
	ContentTypeApplicationJSONCharsetUTF8 = ContentTypeApplicationJSON + "; " + CharsetUTF8
	ContentTypeMultipartForm              = "multipart/form-data"
	ContentTypeOctetStream                = "application/octet-stream"
	ContentTypeTextPlain                  = "text/plain"
	ContentTypeTextPlainCharsetUTF8       = ContentTypeTextPlain + "; " + CharsetUTF8

	HeaderContentType = "Content-Type"

	ParameterStyleDeepObject     = "deepObject"
	ParameterStyleForm           = "form"
	ParameterStyleLabel          = "label"
	ParameterStyleMatrix         = "matrix"
	ParameterStylePipeDelimited  = "pipeDelimited"
	ParameterStyleSimple         = "simple"
	ParameterStyleSpaceDelimited = "spaceDelimited"

	ParameterTypeCookie = "cookie"
	ParameterTypeHeader = "header"
	ParameterTypePath   = "path"
	ParameterTypeQuery  = "query"
)
