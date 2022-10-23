// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package server_test

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/conflowio/conflow/pkg/conflow/types"
	"github.com/conflowio/conflow/pkg/openapi"
	"github.com/conflowio/conflow/pkg/openapi/server"
	"github.com/conflowio/conflow/pkg/schema"
	"github.com/conflowio/conflow/pkg/util/ptr"
)

type testRequest struct {
	Params map[string]string
	Req    *http.Request
}

func (t *testRequest) Param(name string) string {
	return t.Params[name]
}

// Request returns `*http.Request`.
func (t *testRequest) Request() *http.Request {
	return t.Req
}

var _ = Describe("BindParameter", func() {
	type target struct {
		FieldBool    bool
		FieldBoolPtr *bool
		FieldBoolArr []bool

		FieldInt    int64
		FieldIntPtr *int64
		FieldIntArr []int64

		FieldNumber    float64
		FieldNumberPtr *float64
		FieldNumberArr []float64

		FieldString    string
		FieldStringPtr *string
		FieldStringArr []string

		FieldDuration    types.Duration
		FieldDurationPtr *types.Duration
	}

	var t target
	var p *openapi.Parameter

	var pathParams map[string]string
	var headers http.Header
	var values url.Values
	var req *testRequest

	BeforeEach(func() {
		t = target{}
		p = &openapi.Parameter{}
		pathParams = map[string]string{
			"booltrue":  "true",
			"boolfalse": "false",
			"int":       "1",
			"float":     "1.2",
			"string":    "foo",
			"duration":  "10s",
		}
		values = url.Values{
			"booltrue":    {"true"},
			"boolfalse":   {"false"},
			"int":         {"1"},
			"float":       {"1.2"},
			"string":      {"foo"},
			"duration":    {"10s"},
			"boolarr":     {"false", "true"},
			"intarr":      {"1", "2"},
			"floatarr":    {"1", "1.1"},
			"stringarr":   {"foo", "bar"},
			"durationarr": {"1s", "2s"},
		}
		headers = http.Header{
			"Booltrue":  {"true"},
			"Boolfalse": {"false"},
			"Int":       {"1"},
			"Float":     {"1.2"},
			"String":    {"foo"},
			"Duration":  {"10s"},
			"Cookie": {
				"booltrue=true",
				"boolfalse=false",
				"int=1",
				"float=1.2",
				"string=foo",
				"duration=10s",
			},
		}
	})

	JustBeforeEach(func() {
		url, err := url.Parse("?" + values.Encode())
		Expect(err).ToNot(HaveOccurred())

		req = &testRequest{
			Params: pathParams,
			Req: &http.Request{
				Header: headers,
				URL:    url,
			},
		}
	})

	for _, ptype := range []string{
		openapi.ParameterTypePath,
		openapi.ParameterTypeQuery,
		openapi.ParameterTypeHeader,
		openapi.ParameterTypeCookie,
	} {
		ptype := ptype
		Context(fmt.Sprintf("%s parameter", ptype), func() {
			BeforeEach(func() {
				p.In = ptype
			})

			It("parses a bool true parameter", func() {
				p.Name = "booltrue"
				p.Schema = &schema.Boolean{}
				Expect(server.BindParameter[bool](p, req, &t.FieldBool)).ToNot(HaveOccurred())
				Expect(t.FieldBool).To(Equal(true))
			})

			It("binds a bool false parameter", func() {
				p.Name = "boolfalse"
				p.Schema = &schema.Boolean{}
				Expect(server.BindParameter[bool](p, req, &t.FieldBool)).ToNot(HaveOccurred())
				Expect(t.FieldBool).To(Equal(false))
			})

			It("binds into a bool pointer field", func() {
				p.Name = "booltrue"
				p.Schema = &schema.Boolean{}
				Expect(server.BindParameter[bool](p, req, &t.FieldBoolPtr)).ToNot(HaveOccurred())
				Expect(t.FieldBoolPtr).To(HaveValue(Equal(true)))
			})

			It("errors if not bool value", func() {
				p.Name = "string"
				p.Schema = &schema.Boolean{}
				err := server.BindParameter[bool](p, req, nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(MatchRegexp("failed to parse 'string' .*: must be 'true' or 'false'"))
			})

			It("parses an int parameter", func() {
				p.Name = "int"
				p.Schema = &schema.Integer{}
				Expect(server.BindParameter[int64](p, req, &t.FieldInt)).ToNot(HaveOccurred())
				Expect(t.FieldInt).To(Equal(int64(1)))
			})

			It("binds into an int pointer field", func() {
				p.Name = "int"
				p.Schema = &schema.Integer{}
				Expect(server.BindParameter[int64](p, req, &t.FieldIntPtr)).ToNot(HaveOccurred())
				Expect(t.FieldIntPtr).To(HaveValue(Equal(int64(1))))
			})

			It("errors if not int value", func() {
				p.Name = "string"
				p.Schema = &schema.Integer{}
				err := server.BindParameter[int](p, req, nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(MatchRegexp("failed to parse 'string' .*: must be integer"))
			})

			It("parses a number parameter", func() {
				p.Name = "float"
				p.Schema = &schema.Number{}
				Expect(server.BindParameter[float64](p, req, &t.FieldNumber)).ToNot(HaveOccurred())
				Expect(t.FieldNumber).To(Equal(1.2))
			})

			It("binds into a number pointer field", func() {
				p.Name = "float"
				p.Schema = &schema.Number{}
				Expect(server.BindParameter[float64](p, req, &t.FieldNumberPtr)).ToNot(HaveOccurred())
				Expect(t.FieldNumberPtr).To(HaveValue(Equal(1.2)))
			})

			It("errors if not number value", func() {
				p.Name = "string"
				p.Schema = &schema.Number{}
				err := server.BindParameter[float64](p, req, nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(MatchRegexp("failed to parse 'string' .*: must be number"))
			})

			It("parses a string parameter", func() {
				p.Name = "string"
				p.Schema = &schema.String{}
				Expect(server.BindParameter[string](p, req, &t.FieldString)).ToNot(HaveOccurred())
				Expect(t.FieldString).To(Equal("foo"))
			})

			It("binds into a string pointer field", func() {
				p.Name = "string"
				p.Schema = &schema.String{}
				Expect(server.BindParameter[string](p, req, &t.FieldStringPtr)).ToNot(HaveOccurred())
				Expect(t.FieldStringPtr).To(HaveValue(Equal("foo")))
			})

			It("parses a string format parameter", func() {
				p.Name = "duration"
				p.Schema = &schema.String{Format: schema.FormatDuration}
				Expect(server.BindParameter[types.Duration](p, req, &t.FieldDuration)).ToNot(HaveOccurred())
				Expect(t.FieldDuration).To(Equal(types.Duration(10 * time.Second)))
			})

			It("binds into a string format pointer field", func() {
				p.Name = "duration"
				p.Schema = &schema.String{Format: schema.FormatDuration}
				Expect(server.BindParameter[types.Duration](p, req, &t.FieldDurationPtr)).ToNot(HaveOccurred())
				Expect(t.FieldDurationPtr).To(HaveValue(Equal(types.Duration(10 * time.Second))))
			})

			It("errors if not string format value", func() {
				p.Name = "string"
				p.Schema = &schema.String{Format: schema.FormatDuration}
				err := server.BindParameter[types.Duration](p, req, nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(MatchRegexp("failed to parse 'string' .*: time: invalid duration \"foo\""))
			})

			It("returns an error if parameter is empty", func() {
				p.Name = "missing"
				p.Schema = &schema.String{}
				p.Required = ptr.To(true)
				err := server.BindParameter[string](p, req, &t.FieldString)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(MatchRegexp(fmt.Sprintf("'missing' %s.*must be set", ptype)))
			})
		})
	}

	Context("query parameter as array", func() {
		BeforeEach(func() {
			p.In = openapi.ParameterTypeQuery
		})

		It("parses a bool array parameter", func() {
			p.Name = "boolarr"
			p.Schema = &schema.Array{Items: &schema.Boolean{}}
			Expect(server.BindParameter[bool](p, req, &t.FieldBoolArr)).ToNot(HaveOccurred())
			Expect(t.FieldBoolArr).To(Equal([]bool{false, true}))
		})

		It("parses an int array parameter", func() {
			p.Name = "intarr"
			p.Schema = &schema.Array{Items: &schema.Integer{}}
			Expect(server.BindParameter[int64](p, req, &t.FieldIntArr)).ToNot(HaveOccurred())
			Expect(t.FieldIntArr).To(Equal([]int64{1, 2}))
		})

		It("parses a number array parameter", func() {
			p.Name = "floatarr"
			p.Schema = &schema.Array{Items: &schema.Number{}}
			Expect(server.BindParameter[float64](p, req, &t.FieldNumberArr)).ToNot(HaveOccurred())
			Expect(t.FieldNumberArr).To(Equal([]float64{1, 1.1}))
		})

		It("parses a string array parameter", func() {
			p.Name = "stringarr"
			p.Schema = &schema.Array{Items: &schema.String{}}
			Expect(server.BindParameter[string](p, req, &t.FieldStringArr)).ToNot(HaveOccurred())
			Expect(t.FieldStringArr).To(Equal([]string{"foo", "bar"}))
		})
	})

})
