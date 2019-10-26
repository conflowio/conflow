// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives_test

func ExampleRun_ShortFormat() {
	eval(`
		println "This will run"

		@run false
		println "This won't run"
	`)
	// Output: This will run
}

func ExampleRun_LongFormat() {
	eval(`
		println "This will run"

		@run {
          when = false
        }
		println "This won't run"
	`)
	// Output: This will run
}
