package directives

import (
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/schema"
)

//go:generate basil generate
type Description struct {
	// @id
	id basil.ID
	// @value
	// @required
	value string
}

func (d *Description) ID() basil.ID {
	return d.id
}

func (d *Description) ApplyToSchema(s schema.Schema) error {
	return schema.UpdateMetadata(s, func(meta schema.MetadataAccessor) error {
		meta.SetDescription(d.value)
		return nil
	})
}
