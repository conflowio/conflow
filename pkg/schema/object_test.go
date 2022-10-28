// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package schema_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"go/format"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/conflowio/conflow/pkg/internal/testhelper"
	"github.com/conflowio/conflow/pkg/schema"
)

var _ schema.Schema = &schema.Object{}

var _ = Describe("Object", func() {
	type TestObject2 struct {
		Qux bool `json:"qux"`
	}

	type TestObject struct {
		Foo         int64       `json:"foo,omitempty"`
		FooPointer  *int64      `json:"foop,omitempty"`
		Bar         interface{} `json:"bar,omitempty"`
		Baz         TestObject2 `json:"baz,omitempty"`
		notExported int8
		Ignored     int8 `json:"-"'`
	}

	type OtherObject struct {
		OtherField string `json:"otherField,omitempty"`
	}

	defaultSchema := func() *schema.Object {
		return &schema.Object{
			Properties: map[string]schema.Schema{
				"foo":  &schema.Integer{},
				"foop": &schema.Integer{Nullable: true},
				"bar":  &schema.String{},
				"baz": &schema.Object{
					Properties: map[string]schema.Schema{
						"qux": &schema.Boolean{},
					},
				},
			},
		}
	}

	DescribeTable("Validate accepts value",
		func(f func(s *schema.Object), value interface{}) {
			s := defaultSchema()
			f(s)
			_, err := s.ValidateValue(value)
			Expect(err).ToNot(HaveOccurred())

			// Let's test the map type
			if _, isMap := value.(map[string]interface{}); !isMap {
				j, err := json.Marshal(value)
				Expect(err).To(Not(HaveOccurred()))
				var m map[string]interface{}
				Expect(json.Unmarshal(j, &m)).ToNot(HaveOccurred())
				_, err = s.ValidateValue(value)
				Expect(err).ToNot(HaveOccurred(), "validating the map type failed")
			}
		},
		Entry(
			"empty object",
			func(*schema.Object) {},
			TestObject{},
		),
		Entry(
			"non-empty object",
			func(*schema.Object) {},
			TestObject{Foo: 1},
		),
		Entry(
			"complex object",
			func(*schema.Object) {},
			TestObject{
				Foo: 1,
				Bar: "value",
				Baz: TestObject2{Qux: true},
			},
		),
		Entry(
			"required value set",
			func(s *schema.Object) {
				s.Required = []string{"foo"}
			},
			TestObject{
				Foo: 1,
			},
		),
		Entry(
			"const value",
			func(s *schema.Object) {
				s.Const = map[string]interface{}{
					"foo": int64(1),
				}
			},
			TestObject{
				Foo: 1,
			},
		),
		Entry(
			"const value - empty object",
			func(s *schema.Object) {
				s.Const = map[string]interface{}{}
			},
			TestObject{},
		),
		Entry(
			"enum value - empty object",
			func(s *schema.Object) {
				s.Enum = []map[string]interface{}{}
			},
			TestObject{},
		),
		Entry(
			"enum value - single",
			func(s *schema.Object) {
				s.Enum = []map[string]interface{}{
					{"foo": int64(1)},
				}
			},
			TestObject{
				Foo: 1,
			},
		),
		Entry(
			"enum value - multiple",
			func(s *schema.Object) {
				s.Enum = []map[string]interface{}{
					{"foo": int64(1)},
					{"foo": int64(2)},
				}
			},
			TestObject{
				Foo: 1,
			},
		),
		Entry(
			"minProperties=1, 1 element",
			func(s *schema.Object) {
				s.MinProperties = 1
			},
			TestObject{
				Foo: 1,
			},
		),
		Entry(
			"minProperties=1, 2 elements",
			func(s *schema.Object) {
				s.MinProperties = 2
			},
			TestObject{
				Foo: 1,
				Bar: "xxx",
			},
		),
		Entry(
			"maxProperties=2, 2 elements",
			func(s *schema.Object) {
				s.MaxProperties = schema.Pointer(int64(2))
			},
			TestObject{
				Foo: 1,
				Bar: "xxx",
			},
		),
		Entry(
			"maxProperties=2, 1 element",
			func(s *schema.Object) {
				s.MaxProperties = schema.Pointer(int64(2))
			},
			TestObject{
				Foo: 1,
			},
		),
		Entry(
			"dependentRequired",
			func(s *schema.Object) {
				s.DependentRequired = map[string][]string{
					"foo": {"bar"},
				}
			},
			TestObject{
				Foo: 1,
				Bar: "xxx",
			},
		),
	)

	DescribeTable("Validate errors",
		func(f func(s *schema.Object), value interface{}, expectedErr error) {
			s := defaultSchema()
			f(s)
			_, err := s.ValidateValue(value)
			Expect(err).To(MatchError(expectedErr))
		},
		Entry(
			"non-object value",
			func(s *schema.Object) {},
			"foo",
			errors.New("must be object"),
		),
		Entry(
			"property does not exist",
			func(s *schema.Object) {},
			OtherObject{OtherField: "foo"},
			schema.NewFieldError("otherField", errors.New("property does not exist")),
		),
		Entry(
			"property does not exist - map input",
			func(s *schema.Object) {},
			map[string]interface{}{"otherField": "foo"},
			schema.NewFieldError("otherField", errors.New("property does not exist")),
		),
		Entry(
			"invalid property value",
			func(s *schema.Object) {},
			TestObject{
				Foo: 1,
				Bar: 123,
			},
			schema.NewFieldError("bar", errors.New("must be string")),
		),
		Entry(
			"invalid property value - map input",
			func(s *schema.Object) {},
			map[string]interface{}{"foo": int64(1), "bar": int64(123)},
			schema.NewFieldError("bar", errors.New("must be string")),
		),
		Entry(
			"required value not set",
			func(s *schema.Object) {
				s.Required = []string{"foo"}
			},
			TestObject{},
			schema.NewFieldError("foo", errors.New("required")),
		),
		Entry(
			"const value",
			func(s *schema.Object) {
				s.Const = map[string]interface{}{
					"foo": int64(1),
				}
			},
			TestObject{Foo: 2},
			errors.New("must be {foo: 1}"),
		),
		Entry(
			"const value - empty object",
			func(s *schema.Object) {
				s.Const = map[string]interface{}{}
			},
			TestObject{Foo: 1},
			errors.New("must be {}"),
		),
		Entry(
			"enum value - empty object",
			func(s *schema.Object) {
				s.Enum = []map[string]interface{}{{}}
			},
			TestObject{Foo: 1},
			errors.New("must be {}"),
		),
		Entry(
			"enum value - single",
			func(s *schema.Object) {
				s.Enum = []map[string]interface{}{
					{"foo": int64(1)},
				}
			},
			TestObject{Foo: 2},
			errors.New("must be {foo: 1}"),
		),
		Entry(
			"enum value - multiple",
			func(s *schema.Object) {
				s.Enum = []map[string]interface{}{
					{"foo": int64(1)},
					{"foo": int64(2)},
				}
			},
			TestObject{Foo: 3},
			errors.New("must be one of {foo: 1}, {foo: 2}"),
		),
		Entry(
			"minProperties: 1, empty",
			func(s *schema.Object) {
				s.MinProperties = 1
			},
			TestObject{},
			errors.New("the object can not be empty"),
		),
		Entry(
			"minProperties: 2, 1 element",
			func(s *schema.Object) {
				s.MinProperties = 2
			},
			TestObject{Foo: 1},
			errors.New("the object must have at least 2 properties defined"),
		),
		Entry(
			"maxProperties: 0, 1 element",
			func(s *schema.Object) {
				s.MaxProperties = schema.Pointer(int64(0))
			},
			TestObject{Foo: 1},
			errors.New("the object must be empty"),
		),
		Entry(
			"maxProperties: 1, 2 elements",
			func(s *schema.Object) {
				s.MaxProperties = schema.Pointer(int64(1))
			},
			TestObject{Foo: 1, Bar: "xxx"},
			errors.New("the object can only have a single property defined"),
		),
		Entry(
			"dependentRequired - one missing",
			func(s *schema.Object) {
				s.DependentRequired = map[string][]string{
					"foo": {"bar"},
				}
			},
			TestObject{Foo: 1},
			schema.NewFieldError("bar", errors.New("required")),
		),
	)

	DescribeTable("GoString prints a valid Go struct",
		func(schema *schema.Object, expected string) {
			str := schema.GoString(map[string]string{})
			Expect(str).To(Equal(expected))
		},
		Entry(
			"empty",
			&schema.Object{},
			`&schema.Object{
}`,
		),
		Entry(
			"const",
			&schema.Object{
				Const: map[string]interface{}{"foo": "bar"},
			},
			`&schema.Object{
	Const: map[string]interface {}{"foo":"bar"},
}`,
		),
		Entry(
			"default",
			&schema.Object{
				Default: map[string]interface{}{"foo": "bar"},
			},
			`&schema.Object{
	Default: map[string]interface {}{"foo":"bar"},
}`,
		),
		Entry(
			"enum",
			&schema.Object{
				Enum: []map[string]interface{}{{"foo": "bar"}},
			},
			`&schema.Object{
	Enum: []map[string]interface {}{map[string]interface {}{"foo":"bar"}},
}`,
		),
		Entry(
			"properties",
			&schema.Object{
				Properties: map[string]schema.Schema{
					"bar": &schema.String{Format: "f1"},
					"foo": &schema.String{Format: "f2"},
				},
			},
			`&schema.Object{
	Properties: map[string]schema.Schema{
		"bar": &schema.String{
			Format: "f1",
		},
		"foo": &schema.String{
			Format: "f2",
		},
	},
}`,
		),
		Entry(
			"required",
			&schema.Object{
				Required: []string{"foo"},
			},
			`&schema.Object{
	Required: []string{"foo"},
}`,
		),
		Entry(
			"fieldNames",
			&schema.Object{
				FieldNames: map[string]string{"myField": "MyField"},
			},
			`&schema.Object{
	FieldNames: map[string]string{"myField":"MyField"},
}`,
		),
		Entry(
			"parameterNames",
			&schema.Object{
				ParameterNames: map[string]string{"myField": "my_field"},
			},
			`&schema.Object{
	ParameterNames: map[string]string{"myField":"my_field"},
}`,
		),
		Entry(
			"minProperties",
			&schema.Object{
				MinProperties: 1,
			},
			`&schema.Object{
	MinProperties: 1,
}`,
		),
		Entry(
			"maxProperties",
			&schema.Object{
				MaxProperties: schema.Pointer(int64(1)),
			},
			`&schema.Object{
	MaxProperties: schema.Pointer(int64(1)),
}`,
		),
		Entry(
			"dependentRequired",
			&schema.Object{
				DependentRequired: map[string][]string{"foo": {"bar"}},
			},
			`&schema.Object{
	DependentRequired: map[string][]string{"foo":[]string{"bar"}},
}`,
		),
	)

	It("should unmarshal/marshal a json", func() {
		testhelper.ExpectConsistentJSONMarshalling(
			`{
				"const": {
					"myField": "val1"
				},
				"default": {
					"myField": "val2"
				},
				"dependentRequired": {
					"foo": ["bar"]
				},
				"enum": [
					{
						"myField": "val3"
					},
					{
						"myField": "val4"
					}
				],
				"type": "object",
				"properties": {
					"myField": {
						"type": "string"
					}
				},
				"x-conflow-fields": {
					"myField": "MyField"
				},
				"minProperties": 1,
				"maxProperties": 2,
				"x-conflow-parameters": {
					"myField": "my_field"
				},
				"required": ["myField"]
			}`,
			&schema.Object{},
		)
	})

	It("should error on an invalid parameter name", func() {
		input := `{
			"type": "object",
			"properties": {
				"myField": {
					"type": "string"
				}
			},
			"x-conflow-parameters": {
				"myField": "MyField"
			}
		}`
		_, err := schema.UnmarshalJSON([]byte(input))
		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fmt.Errorf("invalid parameter name \"MyField\", must match %s", schema.NameRegExp.String())))
	})

	It("should ignore an invalid field name", func() {
		input := `{
			"type": "object",
			"properties": {
				"myField": {
					"type": "string"
				},
				"myField2": {
					"type": "string"
				}
			},
			"x-conflow-fields": {
				"myField": "fielName",
				"myField2": "invalid field name"
			}
		}`
		_, err := schema.UnmarshalJSON([]byte(input))
		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fmt.Errorf("invalid field name \"invalid field name\", must match %s", schema.FieldNameRegexp.String())))
	})

	structWithField := func(def string) string {
		return fmt.Sprintf("struct {\n\t%s\n}", def)
	}

	DescribeTable("GoType prints a valid Go struct",
		func(s *schema.Object, expected string, expectedImports ...map[string]string) {
			imports := map[string]string{
				"github.com/test/test/pkg": "",
			}

			var actual string
			Expect(func() {
				actual = s.GoType(imports)
			}).ToNot(Panic())

			actualFormatted, err := format.Source([]byte("type Test " + actual))
			Expect(err).ToNot(HaveOccurred(), "formatting actual failed, content:\n%s", actual)

			expectedFormatted, err := format.Source([]byte("type Test " + expected))
			Expect(err).ToNot(HaveOccurred(), "formatting expected failed")

			Expect(string(actualFormatted)).To(Equal(string(expectedFormatted)), "Expected\n%s\nActual:\n%s\n", string(expectedFormatted), string(actualFormatted))

			if len(expectedImports) > 0 {
				Expect(imports).To(Equal(expectedImports[0]))
			} else {
				Expect(imports).To(Equal(map[string]string{
					"github.com/test/test/pkg": "",
				}))
			}
		},
		Entry(
			"empty object",
			&schema.Object{},
			`struct {
}`,
		),
		Entry(
			"generates a field correctly, all names are the same",
			&schema.Object{
				Properties: map[string]schema.Schema{
					"testField": &schema.Boolean{},
				},
			},
			structWithField("testField bool"),
		),
		Entry(
			"sets the parameter name",
			&schema.Object{
				Properties: map[string]schema.Schema{
					"testField": &schema.Boolean{},
				},
				ParameterNames: map[string]string{
					"testField": "test_field",
				},
			},
			structWithField("// @name \"test_field\"\n\ttestField bool"),
		),
		Entry(
			"sets the json property name",
			&schema.Object{
				Properties: map[string]schema.Schema{
					"testField": &schema.Boolean{},
				},
				ParameterNames: map[string]string{"testField": "test_field"},
				FieldNames:     map[string]string{"testField": "TestField"},
			},
			structWithField("// @name \"test_field\"\n\tTestField bool `json:\"testField,omitempty\"`"),
		),
		Entry(
			"leaves of omitempty if the field is required",
			&schema.Object{
				Properties: map[string]schema.Schema{
					"testField": &schema.Boolean{},
				},
				ParameterNames: map[string]string{"testField": "test_field"},
				FieldNames:     map[string]string{"testField": "TestField"},
				Required:       []string{"testField"},
			},
			structWithField("// @name \"test_field\"\n\tTestField bool `json:\"testField\"`"),
		),
		Entry(
			"generates a bool field correctly",
			&schema.Object{
				Properties: map[string]schema.Schema{
					"testField": &schema.Boolean{},
				},
			},
			structWithField("testField bool"),
		),
		Entry(
			"generates a bool pointer field correctly",
			&schema.Object{
				Properties: map[string]schema.Schema{
					"testField": &schema.Boolean{Nullable: true},
				},
			},
			structWithField("testField *bool"),
		),
		Entry(
			"generates an integer field correctly",
			&schema.Object{
				Properties: map[string]schema.Schema{
					"testField": &schema.Integer{},
				},
			},
			structWithField("testField int64"),
		),
		Entry(
			"generates an integer pointer field correctly",
			&schema.Object{
				Properties: map[string]schema.Schema{
					"testField": &schema.Integer{Nullable: true},
				},
			},
			structWithField("testField *int64"),
		),
		Entry(
			"generates a number field correctly",
			&schema.Object{
				Properties: map[string]schema.Schema{
					"testField": &schema.Number{},
				},
			},
			structWithField("testField float64"),
		),
		Entry(
			"generates a number pointer field correctly",
			&schema.Object{
				Properties: map[string]schema.Schema{
					"testField": &schema.Number{Nullable: true},
				},
			},
			structWithField("testField *float64"),
		),
		Entry(
			"generates a string field correctly",
			&schema.Object{
				Properties: map[string]schema.Schema{
					"testField": &schema.String{},
				},
			},
			structWithField("testField string"),
		),
		Entry(
			"generates a string pointer field correctly",
			&schema.Object{
				Properties: map[string]schema.Schema{
					"testField": &schema.String{Nullable: true},
				},
			},
			structWithField("testField *string"),
		),
		Entry(
			"generates an array field correctly",
			&schema.Object{
				Properties: map[string]schema.Schema{
					"testField": &schema.Array{
						Items: schema.StringValue(),
					},
				},
			},
			structWithField("testField []string"),
		),
		Entry(
			"generates a map field correctly",
			&schema.Object{
				Properties: map[string]schema.Schema{
					"testField": &schema.Map{
						AdditionalProperties: schema.IntegerValue(),
					},
				},
			},
			structWithField("testField map[string]int64"),
		),
		Entry(
			"generates an object field correctly",
			&schema.Object{
				Properties: map[string]schema.Schema{
					"testField": &schema.Reference{
						Ref: "github.com/test/test/somepackage.Foo",
					},
				},
			},
			`struct {
	testField somepackage.Foo
}`,
			map[string]string{
				"github.com/test/test/pkg":         "",
				"github.com/test/test/somepackage": "somepackage",
			},
		),
		Entry(
			"generates an object pointer field correctly",
			&schema.Object{
				Properties: map[string]schema.Schema{
					"testField": &schema.Reference{
						Ref:      "github.com/test/test/somepackage.Foo",
						Nullable: true,
					},
				},
			},
			`struct {
	testField *somepackage.Foo
}`,
			map[string]string{
				"github.com/test/test/pkg":         "",
				"github.com/test/test/somepackage": "somepackage",
			},
		),
	)

})
