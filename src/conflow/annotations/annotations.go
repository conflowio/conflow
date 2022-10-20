// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package annotations

import (
	"fmt"
	"strconv"

	"github.com/conflowio/conflow/src/internal/utils"
)

const (
	EvalStage   = "block.conflow.io/eval_stage"
	Generated   = "block.conflow.io/generated"
	ID          = "block.conflow.io/id"
	Key         = "block.conflow.io/key"
	Type        = "block.conflow.io/type"
	UserDefined = "block.conflow.io/user_defined"
	Value       = "block.conflow.io/value"
)

var annotations = map[string]string{
	EvalStage:   "EvalStage",
	Generated:   "Generated",
	ID:          "ID",
	Key:         "Key",
	Type:        "Type",
	UserDefined: "UserDefined",
	Value:       "Value",
}

func GoString(value string, imports map[string]string) string {
	if k, ok := annotations[value]; ok {
		sel := utils.EnsureUniqueGoPackageSelector(imports, "github.com/conflowio/conflow/src/conflow/annotations")
		return fmt.Sprintf("%s%s", sel, k)
	}
	return strconv.Quote(value)
}
