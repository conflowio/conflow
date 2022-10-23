// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives_test

func ExampleRun_ShortFormat() {
	eval(`
		println "First will run"

		@run false
		println "Second won't run"
	`)
	// Output: First will run
}

func ExampleRun_LongFormat() {
	eval(`
		println "First will run"

		@run {
          when = false
        }
		println "Second won't run"
	`)
	// Output: First will run
}

func ExampleRun_UsingAVariable() {
	eval(`
		run_second := false

		println "First will run"

		@run main.run_second
		println "Second won't run"
	`)
	// Output: First will run
}
