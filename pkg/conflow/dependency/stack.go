// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package dependency

type stack []*node

func (s *stack) Push(n *node) {
	*s = append(*s, n)
	n.OnStack = true
}

func (s *stack) Pop() *node {
	l := len(*s)
	n := (*s)[l-1]
	n.OnStack = false
	*s = (*s)[:l-1]
	return n
}
