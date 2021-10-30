// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives

import (
	"github.com/opsidian/conflow/conflow/block"
)

func Registry() block.InterpreterRegistry {
	return block.InterpreterRegistry{
		"block":             BlockInterpreter{},
		"const":             ConstInterpreter{},
		"default":           DefaultInterpreter{},
		"dependency":        DependencyInterpreter{},
		"deprecated":        DeprecatedInterpreter{},
		"enum":              EnumInterpreter{},
		"eval_stage":        EvalStageInterpreter{},
		"examples":          ExamplesInterpreter{},
		"exclusive_minimum": ExclusiveMinimumInterpreter{},
		"exclusive_maximum": ExclusiveMaximumInterpreter{},
		"function":          FunctionInterpreter{},
		"generated":         GeneratedInterpreter{},
		"id":                IDInterpreter{},
		"ignore":            IgnoreInterpreter{},
		"max_items":         MaxItemsInterpreter{},
		"maximum":           MaximumInterpreter{},
		"min_items":         MinItemsInterpreter{},
		"minimum":           MinimumInterpreter{},
		"multiple_of":       MultipleOfInterpreter{},
		"name":              NameInterpreter{},
		"read_only":         ReadOnlyInterpreter{},
		"required":          RequiredInterpreter{},
		"result_type":       ResultTypeInterpreter{},
		"title":             TitleInterpreter{},
		"types":             TypesInterpreter{},
		"value":             ValueInterpreter{},
		"write_only":        WriteOnlyInterpreter{},
	}
}
