// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package testhelper

import (
	"bytes"
	"encoding/json"

	. "github.com/onsi/gomega"
)

// ExpectConsistentJSONMarshalling unmarshalls the given input JSON into target, encodes target back into JSON,
// then checks whether the input and output jsons are the same
func ExpectConsistentJSONMarshalling(input string, target interface{}) {
	unmarshaler := json.NewDecoder(bytes.NewReader([]byte(input)))
	unmarshaler.DisallowUnknownFields()
	err := unmarshaler.Decode(&target)
	Expect(err).ToNot(HaveOccurred())

	output, err := json.Marshal(target)
	Expect(err).ToNot(HaveOccurred())

	m1 := map[string]interface{}{}
	m2 := map[string]interface{}{}
	Expect(json.Unmarshal([]byte(input), &m1)).ToNot(HaveOccurred())
	Expect(json.Unmarshal(output, &m2)).ToNot(HaveOccurred())
	Expect(m1).To(Equal(m2))
}
