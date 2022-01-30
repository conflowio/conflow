// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package util

// StringSliceContains tells whether a contains x.
func StringSliceContains(slice []string, search string) bool {
	for _, str := range slice {
		if str == search {
			return true
		}
	}
	return false
}
