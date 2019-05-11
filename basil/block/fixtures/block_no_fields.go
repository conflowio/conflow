package fixtures

import "github.com/opsidian/basil/basil"

//go:generate basil generate
type BlockNoFields struct {
	IDField basil.ID `basil:"id"`
}
