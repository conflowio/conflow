package fixtures

import "github.com/opsidian/basil/basil"

//go:generate basil generate
type BlockValueRequired struct {
	IDField basil.ID    `basil:"id"`
	Value   interface{} `basil:"value,required"`
}
