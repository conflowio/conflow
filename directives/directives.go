// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives

import "github.com/opsidian/basil/basil/block"

func DefaultRegistry() block.InterpreterRegistry {
	return block.InterpreterRegistry{
		"bug":        BugInterpreter{},
		"deprecated": DeprecatedInterpreter{},
		"doc":        DocInterpreter{},
		"input":      InputInterpreter{},
		"output":     OutputInterpreter{},
		"retry":      RetryInterpreter{},
		"run":        RunInterpreter{},
		"skip":       SkipInterpreter{},
		"timeout":    TimeoutInterpreter{},
		"todo":       TodoInterpreter{},
		"triggers":   TriggersInterpreter{},
	}
}
