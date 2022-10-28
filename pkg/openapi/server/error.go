// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package server

import (
	"fmt"
	"net/http"
)

type HTTPError struct {
	Code int
	Err  error
}

func NewHTTPError(code int, err error) *HTTPError {
	return &HTTPError{
		Code: code,
		Err:  err,
	}
}

func NewHTTPErrorf(code int, format string, a ...interface{}) *HTTPError {
	return &HTTPError{
		Code: code,
		Err:  fmt.Errorf(format, a...),
	}
}

func (h *HTTPError) Error() string {
	if h.Err == nil {
		return http.StatusText(h.Code)
	}
	return h.Err.Error()
}

func (h *HTTPError) Unwrap() error {
	return h.Err
}
