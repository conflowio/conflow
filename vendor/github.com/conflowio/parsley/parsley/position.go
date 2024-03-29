// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsley

import "fmt"

// Position is an interface to translate a file position to a string
//counterfeiter:generate . Position
type Position interface {
	fmt.Stringer
}

type nilPosition int

func (np nilPosition) String() string {
	return "unknown"
}

// NilPosition represents an invalid position
const NilPosition = nilPosition(0)

// Pos is a global offset in a file set which can be translated into a concrete file position
type Pos int

// NilPos represents an invalid position
const NilPos = Pos(0)
