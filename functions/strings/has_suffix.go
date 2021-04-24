// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package strings

import (
	"strings"
)

// HasSuffix tests whether the string s ends with suffix.
// @function
func HasSuffix(s, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}
