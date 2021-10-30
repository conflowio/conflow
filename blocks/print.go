// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package blocks

import (
	"bufio"
	"context"
	"fmt"
	"io"

	"github.com/opsidian/conflow/conflow"
)

// Print will write a string to the standard output
// @block
type Print struct {
	// @id
	id conflow.ID
	// @value
	// @required
	value interface{}
	// @dependency
	stdout io.Writer
}

func (p *Print) ID() conflow.ID {
	return p.id
}

func (p *Print) Run(ctx context.Context) (conflow.Result, error) {
	switch v := p.value.(type) {
	case io.ReadCloser:
		first := true
		scanner := bufio.NewScanner(v)
		for scanner.Scan() {
			if !first {
				if _, err := fmt.Fprintln(p.stdout); err != nil {
					return nil, err
				}
			}
			if _, err := fmt.Fprint(p.stdout, scanner.Text()); err != nil {
				return nil, err
			}
			first = false
		}
		return nil, nil
	default:
		if _, err := fmt.Fprint(p.stdout, v); err != nil {
			return nil, err
		}
	}
	return nil, nil
}
