// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package utils

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

var versionRegex = regexp.MustCompile(`^v\d+$`)

func EnsureUniqueGoPackageSelector(imports map[string]string, path string) string {
	selector, ok := imports[path]

	if !ok {
		packageParts := strings.Split(path, "/")
		selectorBase := packageParts[len(packageParts)-1]
		if versionRegex.MatchString(selectorBase) && len(packageParts) > 1 {
			selectorBase = packageParts[len(packageParts)-2]
		}
		selector = selectorBase

		n := 1
		for {
			var found bool
			for _, s := range imports {
				if s == selector {
					found = true
					break
				}
			}
			if !found {
				break
			}

			n++
			selector = fmt.Sprintf("%s%d", selectorBase, n)
		}

		imports[path] = selector
	}

	if selector == "" || selector == "." {
		return ""
	}

	return fmt.Sprintf("%s.", selector)
}

func GoType(imports map[string]string, t reflect.Type, nullable bool) string {
	pkgPath := t.PkgPath()
	if pkgPath == "" {
		if nullable {
			return fmt.Sprintf("*%s", t.Name())
		}
		return t.Name()
	}

	sel := EnsureUniqueGoPackageSelector(imports, pkgPath)
	if nullable {
		return fmt.Sprintf("*%s%s", sel, t.Name())
	}

	return fmt.Sprintf("%s%s", sel, t.Name())
}
