// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package schema_test

import (
	"encoding/json"

	"github.com/conflowio/conflow/conflow/schema"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ schema.Schema = &schema.Function{}
var _ schema.FunctionKind = &schema.Function{}

var _ = Describe("Function", func() {

	DescribeTable("GoString prints a valid Go struct",
		func(schema *schema.Function, expected string) {
			str := schema.GoString(map[string]string{})
			Expect(str).To(Equal(expected))
		},
		Entry(
			"empty",
			&schema.Function{},
			`&schema.Function{
}`,
		),
		Entry(
			"additionalParameters",
			&schema.Function{
				AdditionalParameters: &schema.NamedSchema{
					Name: "foo",
					Schema: &schema.String{
						Format: "foo",
					},
				},
			},
			`&schema.Function{
	AdditionalParameters: &schema.NamedSchema{
		Name: "foo",
		Schema: &schema.String{
			Format: "foo",
		},
	},
}`,
		),
		Entry(
			"parameters",
			&schema.Function{
				Parameters: schema.Parameters{
					{
						Name:   "bar",
						Schema: &schema.String{Format: "f1"},
					},
					{
						Name:   "foo",
						Schema: &schema.String{Format: "f2"},
					},
				},
			},
			`&schema.Function{
	Parameters: schema.Parameters{
		schema.NamedSchema{
			Name: "bar",
			Schema: &schema.String{
				Format: "f1",
			},
		},
		schema.NamedSchema{
			Name: "foo",
			Schema: &schema.String{
				Format: "f2",
			},
		},
	},
}`,
		),
		Entry(
			"result",
			&schema.Function{Result: &schema.String{Format: "foo"}},
			`&schema.Function{
	Result: &schema.String{
		Format: "foo",
	},
}`,
		),
	)

	It("should marshal/unmarshal", func() {
		s := &schema.Function{
			AdditionalParameters: &schema.NamedSchema{
				Name: "foo",
				Schema: &schema.String{
					Metadata: schema.Metadata{
						Description: "foodesc",
					},
				},
			},
			Parameters: schema.Parameters{
				{
					Name:   "p1",
					Schema: &schema.String{Format: "f1"},
				},
				{
					Name:   "p2",
					Schema: &schema.String{Format: "f2"},
				},
			},
			Result: &schema.String{
				Metadata: schema.Metadata{
					Description: "resdesc",
				},
				Format: "res",
			},
		}
		j, err := json.Marshal(s)
		Expect(err).ToNot(HaveOccurred())

		s2 := &schema.Function{}
		err = json.Unmarshal(j, &s2)
		Expect(err).ToNot(HaveOccurred())
		Expect(s2).To(Equal(s))
	})

})
