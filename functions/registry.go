// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package functions

import (
	"github.com/conflowio/conflow/conflow/function"
	"github.com/conflowio/conflow/functions/array"
	"github.com/conflowio/conflow/functions/json"
	"github.com/conflowio/conflow/functions/math"
	"github.com/conflowio/conflow/functions/strings"
	"github.com/conflowio/conflow/functions/time"
)

func DefaultRegistry() function.InterpreterRegistry {
	return function.InterpreterRegistry{
		// common
		"len":    LenInterpreter{},
		"string": StringInterpreter{},

		// array
		"arr_contains": array.ContainsInterpreter{},

		// json
		"json_decode": json.DecodeInterpreter{},
		"json_encode": json.EncodeInterpreter{},

		// math
		"abs":   math.AbsInterpreter{},
		"ceil":  math.CeilInterpreter{},
		"floor": math.FloorInterpreter{},
		"max":   math.MaxInterpreter{},
		"min":   math.MinInterpreter{},
		"round": math.RoundInterpreter{},

		// strings
		"str_contains":    strings.ContainsInterpreter{},
		"str_format":      strings.FormatInterpreter{},
		"str_has_prefix":  strings.HasPrefixInterpreter{},
		"str_has_suffix":  strings.HasSuffixInterpreter{},
		"str_join":        strings.JoinInterpreter{},
		"str_lower":       strings.LowerInterpreter{},
		"str_replace":     strings.ReplaceInterpreter{},
		"str_split":       strings.SplitInterpreter{},
		"str_title":       strings.TitleInterpreter{},
		"str_trim_prefix": strings.TrimPrefixInterpreter{},
		"str_trim_space":  strings.TrimSpaceInterpreter{},
		"str_trim_suffix": strings.TrimSuffixInterpreter{},
		"str_upper":       strings.UpperInterpreter{},

		// time
		"time_now": time.NowInterpreter{},
	}
}
