// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package functions

import (
	"github.com/opsidian/basil/basil/function"
	"github.com/opsidian/basil/functions/json"
	"github.com/opsidian/basil/functions/math"
	"github.com/opsidian/basil/functions/strings"
)

func Registry() function.InterpreterRegistry {
	return function.InterpreterRegistry{
		"array_contains":          ArrayContainsInterpreter{},
		"is_empty":                IsEmptyInterpreter{},
		"len":                     LenInterpreter{},
		"string":                  StringInterpreter{},
		"math.abs":                math.AbsInterpreter{},
		"math.ceil":               math.CeilInterpreter{},
		"math.floor":              math.FloorInterpreter{},
		"strings.has_prefix":      strings.HasPrefixInterpreter{},
		"strings.has_suffix":      strings.HasSuffixInterpreter{},
		"strings.join":            strings.JoinInterpreter{},
		"strings.lower":           strings.LowerInterpreter{},
		"strings.replace":         strings.ReplaceInterpreter{},
		"strings.split":           strings.SplitInterpreter{},
		"strings.string_contains": strings.ContainsInterpreter{},
		"strings.title":           strings.TitleInterpreter{},
		"strings.trim_prefix":     strings.TrimPrefixInterpreter{},
		"strings.trim_space":      strings.TrimSpaceInterpreter{},
		"strings.trim_suffix":     strings.TrimSuffixInterpreter{},
		"strings.upper":           strings.UpperInterpreter{},
		"json_decode":             json.DecodeInterpreter{},
		"json_encode":             json.EncodeInterpreter{},
	}
}
