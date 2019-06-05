package common

import "github.com/opsidian/basil/basil"

//go:generate basil generate
type Noop struct {
	id basil.ID `basil:"id"`
}
