// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ptr

func To[V any](v V) *V {
	return &v
}

func Value[V any](v *V) (r V) {
	if v != nil {
		r = *v
	}
	return r
}
