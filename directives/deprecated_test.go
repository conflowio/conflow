// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives_test

func ExampleDeprecated() {
	eval(`
		@deprecated "This should not be used"
		println "What, I'm deprecated?!"
	`)
	// Output: What, I'm deprecated?!
}