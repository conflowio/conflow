// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives_test

func ExampleTimeout() {
	eval(`
		@timeout 10ms
		sleep s {
			duration = 20ms
		}
	`)
	// Output:
	// Example errored: timeout reached in sleep.s at 3:3
}
