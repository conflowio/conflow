package fixtures

import (
	"context"
	"time"

	"github.com/opsidian/basil/basil"
)

//go:generate basil generate
type BlockWithContext struct {
	IDField basil.ID `basil:"id"`
}

func (b *BlockWithContext) Context(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, 1*time.Second)
}
