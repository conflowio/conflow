// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives_test

func ExampleSkip_ShortFormat() {
	eval(`
		println "This will run"

		@skip
		println "This won't run"
	`)
	// Output: This will run
}

// FIXME: is deadlocked as no children could be evaluated at
//func ExampleSkip_UsingAVariable() {
//	eval(`
//		skip_second := true
//
//		println "This will run"
//
//		@skip main.skip_second
//		println "This won't run"
//	`)
//	// Output: This will run
//}
