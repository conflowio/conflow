// Copyright (c) 2018 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package block

import "github.com/opsidian/parsley/parsley"

// RegistryAware defines an interface to get a block registry
type RegistryAware interface {
	Registry() Registry
}

// Registry contains a list of block interpreters and behaves as a node transformer registry
type Registry map[string]Interpreter

// NodeTransformer returns with the named node transformer
func (r Registry) NodeTransformer(name string) (parsley.NodeTransformer, bool) {
	interpreter, exists := r[name]
	if !exists {
		return nil, false
	}

	return parsley.NodeTransformFunc(func(node parsley.Node) (parsley.Node, parsley.Error) {
		return transformNode(node, interpreter)
	}), true
}
