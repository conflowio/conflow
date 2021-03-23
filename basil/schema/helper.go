// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package schema

import (
	"fmt"
	"sort"
	"strings"
)

func getSortedMapKeys(v map[string]interface{}) []string {
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func indent(s string) string {
	return strings.ReplaceAll(s, "\n", "\n\t")
}

func EnsureUniqueGoPackageName(imports map[string]string, path string) string {
	packageParts := strings.Split(path, "/")
	packageName := packageParts[len(packageParts)-1]

	n := 1
	for {
		if p, ok := imports[packageName]; !ok || p == path {
			break
		}

		n++
		packageName = fmt.Sprintf("%s%d", packageParts[len(packageParts)-1], n)
	}

	imports[packageName] = path

	return packageName
}
