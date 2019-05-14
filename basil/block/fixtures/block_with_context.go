package fixtures

import (
	"context"
	"time"

	"github.com/opsidian/basil/basil"
)

//go:generate basil generate
type BlockWithContext struct {
	IDField basil.ID           `basil:"id"`
	cancel  context.CancelFunc `basil:"ignore"`
}

func (b *BlockWithContext) Close(blockCtx basil.BlockContext) error {
	b.cancel()
	return nil
}

func (b *BlockWithContext) Context(blockCtx basil.BlockContext) basil.BlockContextOverride {
	ctx, _ := context.WithTimeout(blockCtx.Context(), 1*time.Second)
	return basil.BlockContextOverride{
		Context: ctx,
	}
}
