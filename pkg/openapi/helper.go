// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package openapi

import (
	"fmt"
	"io"
	"strings"

	"github.com/conflowio/conflow/pkg/internal/utils"
)

func indent(s string) string {
	return strings.ReplaceAll(s, "\n", "\n\t")
}

func fprintf(w io.Writer, format string, a ...any) {
	_, err := fmt.Fprintf(w, format, a...)
	if err != nil {
		panic(err)
	}
}

func openapiPkg(imports map[string]string) string {
	return utils.EnsureUniqueGoPackageSelector(imports, "github.com/conflowio/conflow/pkg/openapi")
}

func schemaPkg(imports map[string]string) string {
	return utils.EnsureUniqueGoPackageSelector(imports, "github.com/conflowio/conflow/pkg/schema")
}
