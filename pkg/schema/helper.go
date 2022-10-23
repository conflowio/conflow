// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package schema

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/conflowio/conflow/pkg/internal/utils"
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

func fprintf(w io.Writer, format string, a ...any) {
	_, err := fmt.Fprintf(w, format, a...)
	if err != nil {
		panic(err)
	}
}

func schemaPkg(imports map[string]string) string {
	return utils.EnsureUniqueGoPackageSelector(imports, "github.com/conflowio/conflow/pkg/schema")
}
