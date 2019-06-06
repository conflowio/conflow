package basil

import (
	"github.com/opsidian/basil/util"
	"github.com/opsidian/parsley/parsley"
)

// Container is a basil object container
type Container interface {
	Identifiable
	Job
	Node() Node
	Value() (interface{}, parsley.Error)
	WaitGroups() []*util.WaitGroup
	Close()
}
