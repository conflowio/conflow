// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package server

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/conflowio/conflow/src/openapi"
	"github.com/conflowio/conflow/src/schema"
	"github.com/conflowio/conflow/src/util/ptr"
)

const reservedChars = ":/?#[]@!$&'()*+,;="

type RequestContext interface {
	// Param returns path parameter by name.
	Param(name string) string

	// Request returns `*http.Request`.
	Request() *http.Request
}

func BindParameter[T any](p *openapi.Parameter, ctx RequestContext) (T, error) {
	switch p.In {
	case "path":
		return bindPathParameter[T](p, ctx.Param(p.Name))
	case "query":
		return bindQueryParameter[T](p, ctx.Request().URL.Query())
	case "cookie":
		return bindCookieParameter[T](p, ctx.Request().Cookies())
	case "header":
		return bindHeaderParameter[T](p, ctx.Request().Header)
	default:
		panic(fmt.Errorf("unexpected parameter type: %s", p.In))
	}
}

func BindParameterPtr[T any](p *openapi.Parameter, ctx RequestContext) (*T, error) {
	res, err := BindParameter[T](p, ctx)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func BindParameterArray[T any](p *openapi.Parameter, ctx RequestContext) ([]T, error) {
	switch p.In {
	case "path":
		panic(fmt.Errorf("array type is not supported for %s %s parameter", p.Name, p.In))
	case "query":
		return bindQueryParameterArray[T](p, ctx.Request().URL.Query())
	case "cookie":
		panic(fmt.Errorf("array type is not supported for %s %s parameter", p.Name, p.In))
	case "header":
		panic(fmt.Errorf("array type is not supported for %s %s parameter", p.Name, p.In))
	default:
		panic(fmt.Errorf("unexpected parameter type: %s", p.In))
	}
}

func bindPathParameter[T any](p *openapi.Parameter, rawValue string) (T, error) {
	var emptyVal T

	if p.Style != "" && p.Style != "simple" {
		panic("only style=simple is supported on path parameters")
	}
	if ptr.Value(p.Explode) {
		panic("explode=true is not supported on path parameters")
	}

	if rawValue == "" {
		return emptyVal, fmt.Errorf("%s path parameter must be set", p.Name)
	}

	value, err := url.PathUnescape(rawValue)
	if err != nil {
		return emptyVal, fmt.Errorf("error unescaping path parameter %s: %v", p.Name, err)
	}

	return bindValue[T](p, value)
}

func bindQueryParameter[T any](p *openapi.Parameter, values url.Values) (T, error) {
	var emptyVal T

	unescapedValues, err := getQueryValues(p, values)
	if err != nil {
		return emptyVal, err
	}

	if len(unescapedValues) > 1 {
		return emptyVal, fmt.Errorf("multiple values are not allowed for query parameter %s", p.Name)
	}

	return bindValue[T](p, unescapedValues[0])
}

func bindQueryParameterArray[T any](p *openapi.Parameter, values url.Values) ([]T, error) {
	unescapedValues, err := getQueryValues(p, values)
	if err != nil {
		return nil, err
	}

	return bindValues[T](p, unescapedValues)
}

func getQueryValues(p *openapi.Parameter, values url.Values) ([]string, error) {
	if p.Style != "" && p.Style != "form" {
		panic("only style=form is supported on query parameters")
	}
	if p.Explode != nil && !*p.Explode {
		panic("explode=false is not supported on query parameters")
	}

	if ptr.Value(p.Required) && !values.Has(p.Name) {
		return nil, fmt.Errorf("%s query parameter must be set", p.Name)
	}

	if len(values[p.Name]) == 0 {
		return nil, nil
	}

	unescapedValues := make([]string, 0, len(values[p.Name]))
	for _, v := range values[p.Name] {
		if !p.AllowReserved && strings.ContainsAny(v, reservedChars) {
			return nil, fmt.Errorf("%s query parameter can not contain unescaped reserved characters ('%q')", p.Name, reservedChars)
		}

		uv, err := url.QueryUnescape(v)
		if err != nil {
			return nil, fmt.Errorf("error unescaping query parameter %s: %v", p.Name, err)
		}
		unescapedValues = append(unescapedValues, uv)
	}

	return unescapedValues, nil
}

func bindCookieParameter[T any](p *openapi.Parameter, cookies []*http.Cookie) (T, error) {
	var emptyVal T

	if p.Style != "" && p.Style != "form" {
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

	if ptr.Value(p.Required) && value == nil {
		return emptyVal, fmt.Errorf("%s cookie must be set", p.Name)
	}

	return bindValue[T](p, *value)
}

func bindHeaderParameter[T any](p *openapi.Parameter, headers http.Header) (T, error) {
	var emptyVal T

	if p.Style != "" && p.Style != "simple" {
		panic("only style=simple is allowed on header parameters")
	}
	if ptr.Value(p.Explode) {
		panic("explode=true is not supported on header parameters")
	}

	if ptr.Value(p.Required) && len(headers.Values(p.Name)) == 0 {
		return emptyVal, fmt.Errorf("%s header must be set", p.Name)
	}

	return bindValue[T](p, headers.Get(p.Name))
}

func bindValue[T any](p *openapi.Parameter, stringValue string) (T, error) {
	var emptyVal T

	value, err := parseLiteralValue(p, stringValue)
	if err != nil {
		return emptyVal, err
	}

	validatedValue, err := p.Schema.ValidateValue(value)
	if err != nil {
		return emptyVal, err
	}

	return validatedValue.(T), nil
}

func bindValues[T any](p *openapi.Parameter, stringValues []string) ([]T, error) {
	s := p.Schema.(*schema.Array)
	values := make([]interface{}, 0, len(stringValues))
	for _, stringValue := range stringValues {
		value, err := parseLiteralValue(p, stringValue)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}

	validatedValues, err := p.Schema.ValidateValue(values)
	if err != nil {
		return nil, err
	}

	switch s.Items.(type) {
	case *schema.Boolean, *schema.String, *schema.Integer, *schema.Number:
		return typedArray[T](validatedValues.([]interface{})), nil
	default:
		panic(fmt.Errorf("%s array type is not allowed for %s %s parameter", s.Type(), p.In, p.Name))
	}
}

func parseLiteralValue(p *openapi.Parameter, value string) (interface{}, error) {
	switch s := p.Schema.(type) {
	case *schema.Boolean:
		switch value {
		case "true":
			return true, nil
		case "false":
			return false, nil
		default:
			return nil, fmt.Errorf("invalid value, was expecting 'true' or 'false'")
		}
	case *schema.String:
		return value, nil
	case *schema.Integer:
		i, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("invalid value, was expecting integer")
		}
		return i, nil
	case *schema.Number:
		i, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value, was expecting number")
		}
		return i, nil
	default:
		panic(fmt.Errorf("%s type is not allowed for %s %s parameter", s.Type(), p.In, p.Name))
	}
}

func typedArray[T any](a []interface{}) []T {
	r := make([]T, 0, len(a))
	for _, v := range a {
		r = append(r, v.(T))
	}
	return r
}
