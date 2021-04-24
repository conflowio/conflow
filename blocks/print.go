// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package blocks

import (
	"bufio"
	"fmt"
	"io"

	"github.com/opsidian/basil/basil"
)

// Print will write a string to the standard output
// @block
type Print struct {
	// @id
	id basil.ID
	// @value
	// @required
	value interface{}
}

func (p *Print) ID() basil.ID {
	return p.id
}

func (p *Print) Main(ctx basil.BlockContext) error {
	switch v := p.value.(type) {
	case io.ReadCloser:
		first := true
		scanner := bufio.NewScanner(v)
		for scanner.Scan() {
			if !first {
				if _, err := fmt.Fprintln(ctx.Stdout()); err != nil {
					return err
				}
			}
			if _, err := fmt.Fprint(ctx.Stdout(), scanner.Text()); err != nil {
				return err
			}
			first = false
		}
		return nil
	default:
		if _, err := fmt.Fprint(ctx.Stdout(), v); err != nil {
			return err
		}
	}
	return nil
}
