// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package utils

import (
	"regexp"
	"strings"
)

var reSnakeCaseFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var reSnakeCaseAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")
var reCamelCaseUnderscore = regexp.MustCompile("(^|_+)[a-zA-Z0-9]")

func ToSnakeCase(name string) string {
	name = reSnakeCaseFirstCap.ReplaceAllString(name, "${1}_${2}")
	name = reSnakeCaseAllCap.ReplaceAllString(name, "${1}_${2}")
	return strings.ToLower(name)
}

func ToCamelCase(name string) string {
	return reCamelCaseUnderscore.ReplaceAllStringFunc(name, func(s string) string {
		return strings.ToUpper(string(s[len(s)-1]))
	})
}
