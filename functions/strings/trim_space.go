// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package strings

import (
	"strings"
)

// TrimSpace returns a slice of the string s, with all leading
// and trailing white space removed, as defined by Unicode.
//go:generate basil generate
func TrimSpace(s string) string {
	return strings.TrimSpace(s)
}
