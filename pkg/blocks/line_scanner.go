// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package blocks

import (
	"bufio"
	"context"
	"io"
	"strings"

	"github.com/conflowio/conflow/pkg/conflow"
	"github.com/conflowio/conflow/pkg/conflow/block"
)

// @block "generator"
type LineScanner struct {
	// @id
	id conflow.ID
	// @required
	input interface{}
	// @generated
	line *Line
	// @dependency
	blockPublisher conflow.BlockPublisher
}

func (l *LineScanner) ID() conflow.ID {
	return l.id
}

func (l *LineScanner) ParseContextOverride() conflow.ParseContextOverride {
	return conflow.ParseContextOverride{
		BlockTransformerRegistry: block.InterpreterRegistry{
			"line": LineInterpreter{},
		},
	}
}

func (l *LineScanner) Run(ctx context.Context) (conflow.Result, error) {
	var reader io.Reader
	switch v := l.input.(type) {
	case io.ReadCloser:
		reader = v
	case string:
		reader = strings.NewReader(v)
	default:
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		_, err := l.blockPublisher.PublishBlock(&Line{
			id:   l.line.ID(),
			text: scanner.Text(),
		}, nil)
		if err != nil {
			return nil, err
		}

		if ctx.Err() != nil {
			return nil, err
		}
	}
	return nil, nil
}

// @block "configuration"
type Line struct {
	// @id
	id   conflow.ID
	text string
}

func (l *Line) ID() conflow.ID {
	return l.id
}
