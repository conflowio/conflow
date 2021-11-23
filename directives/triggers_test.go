// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives_test

func ExampleTriggers() {
	eval(`
		iterator {
			value = [1, 2, 3]
			a it
		}

		iterator {
			value = [1, 2, 3]
			b it
		}

		delayed_b sleep {
			duration = 10ms
			value := b.value
		}

		@triggers ["a"]
		println {
          value = "a = " + string(a.value) + ", b = " + string(delayed_b.value)
        }
	`)
	// Output:
	// a = 1, b = 1
	// a = 2, b = 1
	// a = 3, b = 1
}
