package function

import "github.com/opsidian/basil/basil/function"

func Registry() function.InterpreterRegistry {
	return function.InterpreterRegistry{
		"abs":             AbsInterpreter{},
		"array_contains":  ArrayContainsInterpreter{},
		"ceil":            CeilInterpreter{},
		"floor":           FloorInterpreter{},
		"has_prefix":      HasPrefixInterpreter{},
		"has_suffix":      HasSuffixInterpreter{},
		"is_empty":        IsEmptyInterpreter{},
		"join":            JoinInterpreter{},
		"json_decode":     JSONDecodeInterpreter{},
		"json_encode":     JSONEncodeInterpreter{},
		"len":             LenInterpreter{},
		"lower":           LowerInterpreter{},
		"replace":         ReplaceInterpreter{},
		"split":           SplitInterpreter{},
		"string":          StringInterpreter{},
		"string_contains": StringContainsInterpreter{},
		"title":           TitleInterpreter{},
		"trim_prefix":     TrimPrefixInterpreter{},
		"trim_space":      TrimSpaceInterpreter{},
		"trim_suffix":     TrimSuffixInterpreter{},
		"upper":           UpperInterpreter{},
	}
}
