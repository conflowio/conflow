package block

import (
	"fmt"

	"github.com/opsidian/basil/basil"
)

var _ basil.BlockContainerRegistry = ContainerRegistry{}

// ContainerRegistry stores block container instances
type ContainerRegistry map[basil.ID]basil.BlockContainer

// NewContainerRegistry creates a new block container registry
func NewContainerRegistry() ContainerRegistry {
	return ContainerRegistry(make(map[basil.ID]basil.BlockContainer, 0))
}

// BlockContainer returns with the given block container instance if it exists
func (c ContainerRegistry) BlockContainer(id basil.ID) (basil.BlockContainer, bool) {
	node, ok := c[id]
	return node, ok
}

// AddBlockContainer adds a new block container instance
// It returns with an error if a block with the same id was already registered
func (c ContainerRegistry) AddBlockContainer(b basil.BlockContainer) error {
	id := b.Block().ID()
	if _, exists := c[id]; exists {
		return fmt.Errorf("duplicated identifier: %q", id)
	}

	c[id] = b

	return nil
}
