// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package strings

import "strings"

// Join concatenates the elements of an array to create a single string. The separator string
// sep is placed between elements in the resulting string.
// @function
func Join(a []string, sep string) string {
	return strings.Join(a, sep)
}
