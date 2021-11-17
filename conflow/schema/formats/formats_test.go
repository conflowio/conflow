// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package formats_test

import (
	. "github.com/onsi/gomega"

	"github.com/conflowio/conflow/conflow/schema"
)

func expectFormatToParse(format schema.Format) func(string, interface{}, string) {
	return func(input string, output interface{}, formattedExpected string) {
		res, err := format.ValidateValue(input)
		Expect(err).ToNot(HaveOccurred())
		Expect(res).To(Equal(output), "output mismatch")

		formatted, _ := format.StringValue(output)
		Expect(formatted).To(Equal(formattedExpected), "format mismatch")
	}
}
