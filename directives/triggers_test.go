// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives_test

func ExampleTriggers() {
	eval(`
		range {
			value = [1, 2, 3]
			entry a
		}

		range {
			value = [1, 2, 3]
			entry b
		}

		sleep delayed_b {
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
