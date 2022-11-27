// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package interpreters

import "github.com/conflowio/conflow/pkg/conflow/block"

func Registry() block.InterpreterRegistry {
	return block.InterpreterRegistry{
		"all_of":  AllOfInterpreter{},
		"any":     AnyInterpreter{},
		"array":   ArrayInterpreter{},
		"boolean": BooleanInterpreter{},
		"integer": IntegerInterpreter{},
		"map":     MapInterpreter{},
		"number":  NumberInterpreter{},
		"string":  StringInterpreter{},
		"object":  ObjectInterpreter{},
		"one_of":  OneOfInterpreter{},
		"ref":     ReferenceInterpreter{},
	}
}
