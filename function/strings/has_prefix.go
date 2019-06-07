// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package strings

import (
	"strings"
)

// HasPrefix tests whether the string s begins with prefix.
//go:generate basil generate
func HasPrefix(s string, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}
