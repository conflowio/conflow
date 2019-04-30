package dependency

import (
	"github.com/opsidian/basil/basil"
)

type node struct {
	Node    basil.Node
	Index   int
	LowLink int
	OnStack bool
}
