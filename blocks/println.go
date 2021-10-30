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

	"github.com/conflowio/conflow/conflow"
)

// Println will write a string followed by a new line to the standard output
// @block
type Println struct {
	// @id
	id conflow.ID
	// @value
	// @required
	value interface{}
	// @dependency
	stdout io.Writer
}

func (p *Println) ID() conflow.ID {
	return p.id
}

func (p *Println) Run(ctx context.Context) (conflow.Result, error) {
	switch v := p.value.(type) {
	case io.ReadCloser:
		scanner := bufio.NewScanner(v)
		for scanner.Scan() {
			if _, err := fmt.Fprintln(p.stdout, scanner.Text()); err != nil {
				return nil, err
			}
		}
		return nil, nil
	default:
		if _, err := fmt.Fprintln(p.stdout, p.value); err != nil {
			return nil, err
		}
	}
	return nil, nil
}
