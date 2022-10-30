// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conflow

import (
	"strconv"
	"strings"

	"github.com/conflowio/parsley/parsley"

	"github.com/conflowio/conflow/pkg/schema"
	"github.com/conflowio/conflow/pkg/util/validation"
)

// FindNodeByPath finds the given child node using the path
// A path element can be:
//   - parameter
//   - parameter[1]
//   - parameter["key"]
func FindNodeByPath(node parsley.Node, path []string) (parsley.Node, []string) {
	if len(path) == 0 {
		if p, ok := node.(ParameterNode); ok {
			return p.ValueNode(), nil
		}
		return node, nil
	}

	if key, err := strconv.Atoi(path[0]); err == nil { // numeric index
		if a, ok := node.(*ArrayNode); ok {
			if key < len(a.Children()) {
				return FindNodeByPath(a.Children()[key], path[1:])
			}
		}
	} else if strings.HasPrefix(path[0], `"`) { // string index
		key := strings.Trim(path[0], `"`)
		if m, ok := node.(*MapNode); ok {
			for i, k := range m.Keys() {
				if k == key {
					return FindNodeByPath(m.Children()[i], path[1:])
				}
			}
		}
	} else {
		fieldParts := validation.FieldRegexp.FindStringSubmatch(path[0])
		if len(fieldParts) < 4 {
			return node, path
		}

		field := fieldParts[1]

		numericKey := -1
		if fieldParts[2] != "" {
			if index, err := strconv.Atoi(fieldParts[2]); err == nil {
				numericKey = index
			}
		}
		stringKey := fieldParts[3]

		if b, ok := node.(BlockNode); ok {
			parameterName := b.Schema().(*schema.Object).ParameterName(field)

			var blockIndex int
			for _, c := range b.Children() {
				switch ct := c.(type) {
				case ParameterNode:
					if ct.Name() == ID(parameterName) {
						if numericKey != -1 {
							return FindNodeByPath(ct, append([]string{strconv.Itoa(numericKey)}, path[1:]...))
						} else if stringKey != "" {
							return FindNodeByPath(ct, append([]string{strconv.Quote(stringKey)}, path[1:]...))
						} else {
							return FindNodeByPath(ct, path[1:])
						}
					}
				case BlockNode:
					if ct.ParameterName() == ID(parameterName) {
						if numericKey != -1 {
							if blockIndex == numericKey {
								return FindNodeByPath(ct, path[1:])
							}
						} else if stringKey != "" {
							if ct.Key() != nil && *ct.Key() == stringKey {
								return FindNodeByPath(ct, path[1:])
							}
						} else {
							return FindNodeByPath(ct, path[1:])
						}
						blockIndex++
					}
				}
			}
		}

	}

	return node, path
}
