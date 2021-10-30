// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package dependency

import "github.com/opsidian/conflow/conflow"

type node struct {
	Node                  conflow.Node
	Index                 int
	LowLink               int
	OnStack               bool
	skip                  bool
	extraDependencies     []conflow.ID
	generatorDependencies []conflow.ID
}
