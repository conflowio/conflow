// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package common

import "github.com/opsidian/basil/basil"

//go:generate basil generate
type Noop struct {
	id basil.ID `basil:"id"`
}
