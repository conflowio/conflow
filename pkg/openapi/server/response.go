// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package server

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"reflect"
	"strings"

	"github.com/conflowio/conflow/pkg/openapi"
)

type ErrUnexpectedResponseFormat struct {
	ExpectedContentType string
	ActualType          reflect.Type
}

func (e ErrUnexpectedResponseFormat) Error() string {
	return fmt.Sprintf("server returned an unexpected response format %s for %s", e.ActualType.String(), e.ExpectedContentType)
}

func WriteResponse(w http.ResponseWriter, contentType string, response interface{}) error {
	switch contentType {
	case openapi.ContentTypeApplicationJSON:
		return json.NewEncoder(w).Encode(response)
	case openapi.ContentTypeOctetStream:
		r, ok := response.([]byte)
		if !ok {
			return ErrUnexpectedResponseFormat{ExpectedContentType: contentType, ActualType: reflect.TypeOf(response)}
		}
		_, err := w.Write(r)
		return err
	}

	if strings.HasPrefix(contentType, "text/") {
		r, ok := response.(string)
		if !ok {
			return ErrUnexpectedResponseFormat{ExpectedContentType: contentType, ActualType: reflect.TypeOf(response)}
		}
		_, err := w.Write([]byte(r))
		return err
	}

	switch r := response.(type) {
	case []byte:
		_, err := w.Write(r)
		return err
	case fs.File:
		defer r.Close()
		_, err := io.Copy(w, r)
		return err
	default:
		return ErrUnexpectedResponseFormat{ExpectedContentType: contentType, ActualType: reflect.TypeOf(response)}
	}
}
