// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package server_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
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
	"github.com/conflowio/conflow/pkg/util/validation"
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

func ExpectHttp400Error(err error, wrappedErr interface{}) {
	Expect(err).To(BeAssignableToTypeOf(&server.HTTPError{}))
	Expect(err.(*server.HTTPError).Code).To(Equal(http.StatusBadRequest))
	if wrappedErr != nil {
		Expect(err).To(MatchError(wrappedErr))
	}
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

var _ = Describe("BindBody", func() {
	type Obj2 struct {
		Field2 string `json:"field2,omitempty"`
	}

	type Obj struct {
		Field    string  `json:"field,omitempty"`
		FieldPtr *string `json:"fieldPtr,omitempty"`
		Obj2     Obj2    `json:"obj2,omitempty"`
	}

	var requestBody *openapi.RequestBody
	var req *testRequest
	objSchema := func() *schema.Object {
		return &schema.Object{
			Properties: map[string]schema.Schema{
				"field":    &schema.String{},
				"fieldPtr": &schema.String{Nullable: true},
				"obj2": &schema.Object{
					Properties: map[string]schema.Schema{
						"field2": &schema.String{},
					},
				},
			},
		}
	}

	BeforeEach(func() {
		req = &testRequest{
			Req: &http.Request{
				Header: map[string][]string{},
			},
		}
		requestBody = &openapi.RequestBody{
			Content: map[string]*openapi.MediaType{
				openapi.ContentTypeApplicationJSON: {Schema: objSchema()},
			},
		}
	})

	Context("json", func() {
		BeforeEach(func() {
			req.Req.Header.Set(openapi.HeaderContentType, openapi.ContentTypeApplicationJSON)
		})

		It("binds valid input", func() {
			req.Req.Body = io.NopCloser(bytes.NewBuffer([]byte(`{"field": "foo", "fieldPtr": "bar", "obj2": {"field2": "baz"}}`)))
			var t Obj
			Expect(server.BindBody[Obj](requestBody, req, &t)).ToNot(HaveOccurred())
			Expect(t).To(Equal(Obj{
				Field:    "foo",
				FieldPtr: ptr.To("bar"),
				Obj2: Obj2{
					Field2: "baz",
				},
			}))
		})

		It("binds valid input to pointer target", func() {
			req.Req.Body = io.NopCloser(bytes.NewBuffer([]byte(`{"field": "foo", "fieldPtr": "bar", "obj2": {"field2": "baz"}}`)))
			var t *Obj
			Expect(server.BindBody[Obj](requestBody, req, &t)).ToNot(HaveOccurred())
			Expect(t).To(HaveValue(Equal(Obj{
				Field:    "foo",
				FieldPtr: ptr.To("bar"),
				Obj2: Obj2{
					Field2: "baz",
				},
			})))
		})

		It("handles empty input", func() {
			req.Req.Body = io.NopCloser(bytes.NewBuffer(nil))
			var t Obj
			Expect(server.BindBody[Obj](requestBody, req, &t)).ToNot(HaveOccurred())
			Expect(t).To(Equal(Obj{}))
		})

		It("handles empty input when using a pointer target", func() {
			req.Req.Body = io.NopCloser(bytes.NewBuffer(nil))
			var t *Obj
			Expect(server.BindBody[Obj](requestBody, req, &t)).ToNot(HaveOccurred())
			Expect(t).To(BeNil())
		})

		It("handles null", func() {
			req.Req.Body = io.NopCloser(bytes.NewBuffer([]byte(`null`)))
			var t Obj
			Expect(server.BindBody[Obj](requestBody, req, &t)).ToNot(HaveOccurred())
			Expect(t).To(Equal(Obj{}))
		})

		It("handles null when using a pointer target", func() {
			req.Req.Body = io.NopCloser(bytes.NewBuffer([]byte(`null`)))
			var t *Obj
			Expect(server.BindBody[Obj](requestBody, req, &t)).ToNot(HaveOccurred())
			Expect(t).To(BeNil())
		})

		It("handles empty object", func() {
			req.Req.Body = io.NopCloser(bytes.NewBuffer([]byte(`{}`)))
			var t Obj
			Expect(server.BindBody[Obj](requestBody, req, &t)).ToNot(HaveOccurred())
			Expect(t).To(Equal(Obj{}))
		})

		It("handles empty object when using a pointer target", func() {
			req.Req.Body = io.NopCloser(bytes.NewBuffer([]byte(`{}`)))
			var t *Obj
			Expect(server.BindBody[Obj](requestBody, req, &t)).ToNot(HaveOccurred())
			Expect(t).To(HaveValue(Equal(Obj{})))
		})

		It("handles JSON syntax error", func() {
			req.Req.Body = io.NopCloser(bytes.NewBuffer([]byte(`x`)))
			var t Obj
			err := server.BindBody[Obj](requestBody, req, &t)
			ExpectHttp400Error(err, "failed to decode JSON request body: invalid character 'x' looking for beginning of value (pos 1)")
		})

		It("handles generic JSON error", func() {
			req.Req.Body = io.NopCloser(bytes.NewBuffer([]byte(`{"xxx": `)))
			var t Obj
			err := server.BindBody[Obj](requestBody, req, &t)
			ExpectHttp400Error(err, "failed to decode JSON request body: unexpected EOF")
		})

		It("handles JSON object unmarshalling error", func() {
			req.Req.Body = io.NopCloser(bytes.NewBuffer([]byte(`{"foo": `)))
			var t Obj
			err := server.BindBody[Obj](requestBody, req, &t)
			ExpectHttp400Error(err, "failed to decode JSON request body: unexpected EOF")
		})

		It("handles JSON object unknown fields", func() {
			req.Req.Body = io.NopCloser(bytes.NewBuffer([]byte(`{"foo": "bar"}`)))
			var t Obj
			err := server.BindBody[Obj](requestBody, req, &t)
			ExpectHttp400Error(err, "failed to decode JSON request body: json: unknown field \"foo\"")
		})

		It("errors on empty input if body is required", func() {
			req.Req.Body = io.NopCloser(bytes.NewBuffer(nil))
			requestBody.Required = true
			var t Obj
			err := server.BindBody[Obj](requestBody, req, &t)
			ExpectHttp400Error(err, "request body can not be empty")
		})

		It("errors on null input if body is required", func() {
			req.Req.Body = io.NopCloser(bytes.NewBuffer([]byte("null")))
			requestBody.Required = true
			var t Obj
			err := server.BindBody[Obj](requestBody, req, &t)
			ExpectHttp400Error(err, "request body can not be empty")
		})

		It("validates JSON input using the schema", func() {
			req.Req.Body = io.NopCloser(bytes.NewBuffer([]byte(`{}`)))
			s := objSchema()
			s.Required = []string{"field"}
			requestBody.Content = map[string]*openapi.MediaType{
				openapi.ContentTypeApplicationJSON: {Schema: s},
			}
			var t Obj
			err := server.BindBody[Obj](requestBody, req, &t)
			ExpectHttp400Error(err, validation.NewFieldError("field", errors.New("required")).Error())
		})
	})

	Context("plain text", func() {
		BeforeEach(func() {
			req.Req.Header.Set(openapi.HeaderContentType, openapi.ContentTypeTextPlain)
			requestBody.Content = map[string]*openapi.MediaType{
				openapi.ContentTypeTextPlain: {Schema: &schema.String{}},
			}
		})

		It("binds valid input", func() {
			req.Req.Body = io.NopCloser(bytes.NewBuffer([]byte(`foo`)))
			var t string
			Expect(server.BindBody[string](requestBody, req, &t)).ToNot(HaveOccurred())
			Expect(t).To(Equal("foo"))
		})

		It("binds valid input to pointer target", func() {
			req.Req.Body = io.NopCloser(bytes.NewBuffer([]byte(`foo`)))
			var t *string
			Expect(server.BindBody[string](requestBody, req, &t)).ToNot(HaveOccurred())
			Expect(t).To(HaveValue(Equal("foo")))
		})

		It("handles empty input", func() {
			req.Req.Body = io.NopCloser(bytes.NewBuffer(nil))
			var t string
			Expect(server.BindBody[string](requestBody, req, &t)).ToNot(HaveOccurred())
			Expect(t).To(Equal(""))
		})

		It("handles empty input when using a pointer target", func() {
			req.Req.Body = io.NopCloser(bytes.NewBuffer(nil))
			var t *string
			Expect(server.BindBody[string](requestBody, req, &t)).ToNot(HaveOccurred())
			Expect(t).To(BeNil())
		})

		It("handles and validates a string format input", func() {
			req.Req.Body = io.NopCloser(bytes.NewBuffer([]byte("1m2s")))
			requestBody.Content = map[string]*openapi.MediaType{
				openapi.ContentTypeTextPlain: {Schema: &schema.String{Format: schema.FormatDuration}},
			}
			var t types.Duration
			Expect(server.BindBody[types.Duration](requestBody, req, &t)).ToNot(HaveOccurred())
			Expect(t).To(Equal(types.Duration(1*time.Minute + 2*time.Second)))
		})

		It("validates the input", func() {
			req.Req.Body = io.NopCloser(bytes.NewBuffer([]byte("foo")))
			requestBody.Content = map[string]*openapi.MediaType{
				openapi.ContentTypeTextPlain: {Schema: &schema.String{MinLength: 4}},
			}
			var t string
			err := server.BindBody[string](requestBody, req, &t)
			ExpectHttp400Error(err, "failed to parse request body: must be at least 4 characters long")
		})

		It("errors on empty input if body is required", func() {
			req.Req.Body = io.NopCloser(bytes.NewBuffer(nil))
			requestBody.Required = true
			var t string
			err := server.BindBody[string](requestBody, req, &t)
			ExpectHttp400Error(err, "request body can not be empty")
		})
	})

})
