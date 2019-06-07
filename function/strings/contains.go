// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package strings

import "strings"

// Contains reports whether substr is within s.
//go:generate basil generate
func Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
