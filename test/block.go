// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package test

import "github.com/opsidian/basil/basil"

//go:generate counterfeiter . BlockWithInit
type BlockWithInit interface {
	basil.Block
	basil.BlockInitialiser
}

//go:generate counterfeiter . BlockWithMain
type BlockWithMain interface {
	basil.Block
	basil.BlockRunner
}

//go:generate counterfeiter . BlockWithClose
type BlockWithClose interface {
	basil.Block
	basil.BlockCloser
}
