// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package strings

import "strings"

// Replace returns a copy of the string s with all
// non-overlapping instances of old replaced by new.
//go:generate basil generate
func Replace(s, old, new string) string {
	return strings.Replace(s, old, new, -1)
}
