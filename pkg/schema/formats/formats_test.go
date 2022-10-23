// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package formats_test

import (
	"encoding/json"

	. "github.com/onsi/gomega"

	"github.com/conflowio/conflow/pkg/conflow"
	"github.com/conflowio/conflow/pkg/conflow/annotations"
	"github.com/conflowio/conflow/pkg/internal/testhelper"
	"github.com/conflowio/conflow/pkg/schema"
)

func expectFormatToParse[T any](format schema.Format, equals ...func(T, T) bool) func(string, interface{}, string, ...bool) {
	return func(input string, output interface{}, formattedExpected string, skips ...bool) {
		res, err := format.ValidateValue(input)
		Expect(err).ToNot(HaveOccurred())
		Expect(res).To(Equal(output), "output mismatch")

		formatted, _ := format.StringValue(output)
		Expect(formatted).To(Equal(formattedExpected), "format mismatch")

		s, err := json.Marshal(input)
		Expect(err).ToNot(HaveOccurred())

		if len(skips) == 0 || skips[0] == false {
			expectConsistentJSONMarshalling[T](s)
		}
	}
}

func expectGoStructToHaveStringSchema(source string, format string, nullable bool) {
	testhelper.ExpectGoStructToHaveSchema(source, &schema.Object{
		Metadata: schema.Metadata{
			ID: "test.Foo",
			Annotations: map[string]string{
				annotations.Type: conflow.BlockTypeConfiguration,
			},
		},
		Properties: map[string]schema.Schema{
			"v": &schema.String{
				Format:   format,
				Nullable: nullable,
			},
		},
	})
}

func expectConsistentJSONMarshalling[T any](s []byte, equals ...func(T, T) bool) {
	var v T
	err := json.Unmarshal(s, &v)
	Expect(err).ToNot(HaveOccurred(), "input: %s", s)

	s2, err := json.Marshal(v)
	Expect(err).ToNot(HaveOccurred(), "input: %s", s)

	Expect(string(s2)).To(Equal(string(s)))

	var v2 T
	err = json.Unmarshal(s2, &v2)
	Expect(err).ToNot(HaveOccurred(), "input: %v", v)

	if len(equals) == 0 {
		Expect(v).To(Equal(v2))
	} else {
		for _, eq := range equals {
			Expect(eq(v, v2)).To(BeTrue(), "mismatch: %v vs %s", v, v2)
		}
	}
}
