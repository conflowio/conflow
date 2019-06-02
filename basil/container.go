package basil

import "github.com/opsidian/parsley/parsley"

// Container is a basil object container
type Container interface {
	Identifiable
	Job
	Value() (interface{}, parsley.Error)
}
