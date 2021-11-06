// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives

import (
	"github.com/conflowio/conflow/conflow/block"
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
		"exclusive_maximum": ExclusiveMaximumInterpreter{},
		"exclusive_minimum": ExclusiveMinimumInterpreter{},
		"format":            FormatInterpreter{},
		"function":          FunctionInterpreter{},
		"generated":         GeneratedInterpreter{},
		"id":                IDInterpreter{},
		"ignore":            IgnoreInterpreter{},
		"max_items":         MaxItemsInterpreter{},
		"max_length":        MaxLengthInterpreter{},
		"maximum":           MaximumInterpreter{},
		"min_items":         MinItemsInterpreter{},
		"min_length":        MinLengthInterpreter{},
		"minimum":           MinimumInterpreter{},
		"multiple_of":       MultipleOfInterpreter{},
		"name":              NameInterpreter{},
		"pattern":           PatternInterpreter{},
		"read_only":         ReadOnlyInterpreter{},
		"required":          RequiredInterpreter{},
		"result_type":       ResultTypeInterpreter{},
		"title":             TitleInterpreter{},
		"types":             TypesInterpreter{},
		"unique_items":      UniqueItemsInterpreter{},
		"value":             ValueInterpreter{},
		"write_only":        WriteOnlyInterpreter{},
	}
}
