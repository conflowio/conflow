// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package server

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/conflowio/conflow/pkg/openapi"
	"github.com/conflowio/conflow/pkg/schema"
	"github.com/conflowio/conflow/pkg/util/ptr"
)

type Request interface {
	// Param returns path parameter by name.
	Param(name string) string

	// Request returns `*http.Request`.
	Request() *http.Request
}

func BindParameter[T any](p *openapi.Parameter, r Request, dest interface{}) error {
	var err error
	switch p.In {
	case openapi.ParameterTypePath:
		err = bindPathParameter[T](p, r.Param(p.Name), dest)
	case openapi.ParameterTypeQuery:
		err = bindQueryParameter[T](p, r.Request().URL.Query(), dest)
	case openapi.ParameterTypeCookie:
		err = bindCookieParameter[T](p, r.Request().Cookies(), dest)
	case openapi.ParameterTypeHeader:
		err = bindHeaderParameter[T](p, r.Request().Header, dest)
	default:
		panic(fmt.Errorf("unexpected parameter type: '%s'", p.In))
	}

	if err != nil {
		return NewHTTPError(http.StatusBadRequest, err)
	}

	return nil
}

func BindBody[T any](reqBody *openapi.RequestBody, r Request, dest interface{}) error {
	defer r.Request().Body.Close()

	contentType := r.Request().Header.Get(openapi.HeaderContentType)
	switch {
	case strings.HasPrefix(contentType, openapi.ContentTypeApplicationJSON):
		c, ok := reqBody.Content[openapi.ContentTypeApplicationJSON]
		if !ok {
			return NewHTTPErrorf(http.StatusBadRequest, "unsupported content type: '%s'", contentType)
		}

		bufferedBody := bufio.NewReader(r.Request().Body)
		peek, _ := bufferedBody.Peek(4)

		if len(peek) == 0 || bytes.Equal(peek, []byte("null")) {
			if reqBody.Required {
				return NewHTTPError(http.StatusBadRequest, errors.New("request body can not be empty"))
			}
			return nil
		}

		var res T
		jsonDecoder := json.NewDecoder(bufferedBody)
		jsonDecoder.DisallowUnknownFields()
		if err := jsonDecoder.Decode(&res); err != nil {
			if ute, ok := err.(*json.UnmarshalTypeError); ok {
				return NewHTTPErrorf(http.StatusBadRequest, "failed to decode JSON request body: expected=%v, got=%v, field=%v, offset=%v", ute.Type, ute.Value, ute.Field, ute.Offset)
			} else if se, ok := err.(*json.SyntaxError); ok {
				return NewHTTPErrorf(http.StatusBadRequest, "failed to decode JSON request body: %w (pos %d)", se, se.Offset)
			}
			return NewHTTPErrorf(http.StatusBadRequest, "failed to decode JSON request body: %w", err)
		}

		validatedValue, err := c.Schema.ValidateValue(res)
		if err != nil {
			return NewHTTPError(http.StatusBadRequest, err)
		}

		return setValue[T](validatedValue, dest)
	case strings.HasPrefix(contentType, openapi.ContentTypeTextPlain):
		c, ok := reqBody.Content[openapi.ContentTypeTextPlain]
		if !ok {
			return NewHTTPErrorf(http.StatusBadRequest, "unsupported content type: '%s'", contentType)
		}

		value, err := io.ReadAll(r.Request().Body)
		if err != nil {
			return err
		}

		if len(value) == 0 {
			if reqBody.Required {
				return NewHTTPError(http.StatusBadRequest, errors.New("request body can not be empty"))
			}
			return nil
		}

		res, err := parseLiteralValue(c.Schema, string(value))
		if err != nil {
			return NewHTTPErrorf(http.StatusBadRequest, "failed to parse request body: %w", err)
		}

		validatedValue, err := c.Schema.ValidateValue(res)
		if err != nil {
			return NewHTTPErrorf(http.StatusBadRequest, "failed to parse request body: %w", err)
		}

		return setValue[T](validatedValue, dest)
	default:
		return NewHTTPErrorf(http.StatusBadRequest, "unsupported content type: '%s'", contentType)
	}
}

func bindPathParameter[T any](p *openapi.Parameter, rawValue string, dest interface{}) error {
	if p.Style != "" && p.Style != openapi.ParameterStyleSimple {
		panic("only style=simple is supported on path parameters")
	}
	if ptr.Value(p.Explode) {
		panic("explode=true is not supported on path parameters")
	}

	if rawValue == "" {
		return fmt.Errorf("'%s' path parameter must be set", p.Name)
	}

	value, err := url.PathUnescape(rawValue)
	if err != nil {
		return fmt.Errorf("error unescaping path parameter '%s': %v", p.Name, err)
	}

	return bindValue[T](p, value, dest)
}

func bindQueryParameter[T any](p *openapi.Parameter, values url.Values, dest interface{}) error {
	if p.Style != "" && p.Style != openapi.ParameterStyleForm {
		panic("only style=form is supported on query parameters")
	}
	if p.Explode != nil && !*p.Explode {
		panic("explode=false is not supported on query parameters")
	}

	if len(values[p.Name]) == 0 {
		if ptr.Value(p.Required) {
			return fmt.Errorf("'%s' query parameter must be set", p.Name)
		}

		return nil
	}

	if p.Schema.Type() == schema.TypeArray {
		return bindValues[T](p, values[p.Name], dest)
	} else {
		if len(values[p.Name]) > 1 {
			return fmt.Errorf("multiple values are not allowed for query parameter '%s'", p.Name)
		}

		return bindValue[T](p, values[p.Name][0], dest)
	}
}

func bindCookieParameter[T any](p *openapi.Parameter, cookies []*http.Cookie, dest interface{}) error {
	if p.Style != "" && p.Style != openapi.ParameterStyleForm {
		panic("only style=form is allowed on cookie parameters")
	}
	if ptr.Value(p.Explode) {
		panic("explode=true is not supported on cookie parameters")
	}

	var value *string
	for _, c := range cookies {
		if c.Name == p.Name {
			value = &c.Value
			break
		}
	}

	if value == nil {
		if ptr.Value(p.Required) {
			return fmt.Errorf("'%s' cookie must be set", p.Name)
		}

		return nil
	}

	return bindValue[T](p, *value, dest)
}

func bindHeaderParameter[T any](p *openapi.Parameter, headers http.Header, dest interface{}) error {
	if p.Style != "" && p.Style != openapi.ParameterStyleSimple {
		panic("only style=simple is allowed on header parameters")
	}
	if ptr.Value(p.Explode) {
		panic("explode=true is not supported on header parameters")
	}

	if len(headers.Values(p.Name)) == 0 {
		if ptr.Value(p.Required) {
			return fmt.Errorf("'%s' header must be set", p.Name)
		}

		return nil
	}

	return bindValue[T](p, headers.Get(p.Name), dest)
}

func bindValue[T any](p *openapi.Parameter, stringValue string, dest interface{}) error {
	value, err := parseLiteralValue(p.Schema, stringValue)
	if err != nil {
		return fmt.Errorf("failed to parse '%s' %s parameter: %w", p.Name, p.In, err)
	}

	validatedValue, err := p.Schema.ValidateValue(value)
	if err != nil {
		return fmt.Errorf("failed to parse '%s' %s parameter: %w", p.Name, p.In, err)
	}

	return setValue[T](validatedValue, dest)
}

func bindValues[T any](p *openapi.Parameter, stringValues []string, dest interface{}) error {
	s := p.Schema.(*schema.Array)
	values := make([]interface{}, 0, len(stringValues))
	for _, stringValue := range stringValues {
		value, err := parseLiteralValue(s.Items, stringValue)
		if err != nil {
			return fmt.Errorf("failed to parse '%s' %s parameter: %w", p.Name, p.In, err)
		}
		values = append(values, value.(T))
	}

	validatedValues, err := p.Schema.ValidateValue(values)
	if err != nil {
		return err
	}

	switch s.Items.(type) {
	case *schema.Boolean, *schema.String, *schema.Integer, *schema.Number:
		return setValue[T](typedArray[T](validatedValues.([]interface{})), dest)
	default:
		panic(fmt.Errorf("'%s' array type is not allowed for '%s' %s parameter", s.Items.Type(), p.In, p.Name))
	}
}

func setValue[T any](v interface{}, dest interface{}) error {
	switch d := dest.(type) {
	case *T:
		if vt, ok := v.(T); ok {
			*d = vt
			return nil
		}
	case **T:
		if vt, ok := v.(T); ok {
			*d = &vt
			return nil
		}
	case *[]T:
		if vt, ok := v.([]T); ok {
			*d = vt
			return nil
		}
	}
	return fmt.Errorf("can not use %T value for setting %T variable", v, dest)
}

func parseLiteralValue(s schema.Schema, value string) (interface{}, error) {
	switch s.(type) {
	case *schema.Boolean:
		switch value {
		case "true":
			return true, nil
		case "false":
			return false, nil
		default:
			return nil, fmt.Errorf("must be 'true' or 'false'")
		}
	case *schema.String:
		return value, nil
	case *schema.Integer:
		i, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("must be integer")
		}
		return int64(i), nil
	case *schema.Number:
		i, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, fmt.Errorf("must be number")
		}
		return i, nil
	default:
		panic(fmt.Errorf("'%s' type is not allowed", s.Type()))
	}
}

func typedArray[T any](a []interface{}) []T {
	r := make([]T, 0, len(a))
	for _, v := range a {
		r = append(r, v.(T))
	}
	return r
}
