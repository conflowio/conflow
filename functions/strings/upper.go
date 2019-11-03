// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package strings

import (
	"strings"
)

// Upper returns a copy of the string s with all Unicode letters mapped to their upper case.
//go:generate basil generate
func Upper(s string) string {
	return strings.ToUpper(s)
}
