// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives_test

func ExampleDoc_SingleLine() {
	eval(`
		@doc "This will print Hello World! to stdout"
		println "Hello World!"
	`)
	// Output: Hello World!
}

func ExampleDoc_Multiline() {
	eval(`
		@doc '''
			This will print Hello World! to stdout
			What else should I say?
		'''
		println "Hello World!"
	`)
	// Output: Hello World!
}
