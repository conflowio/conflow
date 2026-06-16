// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package test

import ginkgo "github.com/onsi/ginkgo/v2"

// TableEntry creates an custom entry for table driven tests where the input is the description
func TableEntry(input string, parameters ...interface{}) ginkgo.TableEntry {
	return ginkgo.Entry(input, append([]interface{}{input}, parameters...)...)
}

// FTableEntry creates an custom focused entry for table driven tests where the input is the description
func FTableEntry(input string, parameters ...interface{}) ginkgo.TableEntry {
	return ginkgo.FEntry(input, append([]interface{}{input}, parameters...)...)
}
