package basil

import (
	"sync"

	"github.com/opsidian/parsley/parsley"
)

// Container is a basil object container
type Container interface {
	Identifiable
	Job
	Node() Node
	Value() (interface{}, parsley.Error)
	WaitGroups() []*sync.WaitGroup
	Close()
}
