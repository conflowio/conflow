// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives_test

func ExampleSkip_ShortFormat() {
	eval(`
		println "First will run"

		@skip
		println "Second won't run"
	`)
	// Output: First will run
}

func ExampleSkip_LongFormat() {
	eval(`
		println "First will run"

		@skip {
          when = true
        }
		println "Second won't run"
	`)
	// Output: First will run
}

func ExampleSkip_UsingAVariable() {
	eval(`
		skip_second := true

		println "First will run"

		@skip main.skip_second
		println "Second won't run"
	`)
	// Output: First will run
}
