// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives

import "github.com/conflowio/conflow/pkg/conflow/block"

func DefaultRegistry() block.InterpreterRegistry {
	return block.InterpreterRegistry{
		"array":      ArrayInterpreter{},
		"boolean":    BooleanInterpreter{},
		"bug":        BugInterpreter{},
		"deprecated": DeprecatedInterpreter{},
		"doc":        DocInterpreter{},
		"input":      InputInterpreter{},
		"integer":    IntegerInterpreter{},
		"map":        MapInterpreter{},
		"number":     NumberInterpreter{},
		"output":     OutputInterpreter{},
		"retry":      RetryInterpreter{},
		"run":        RunInterpreter{},
		"skip":       SkipInterpreter{},
		"string":     StringInterpreter{},
		"timeout":    TimeoutInterpreter{},
		"todo":       TodoInterpreter{},
		"triggers":   TriggersInterpreter{},
	}
}

var schemaRegistry = block.InterpreterRegistry{
	"array":   ArrayInterpreter{},
	"boolean": BooleanInterpreter{},
	"integer": IntegerInterpreter{},
	"map":     MapInterpreter{},
	"number":  NumberInterpreter{},
	"string":  StringInterpreter{},
}
