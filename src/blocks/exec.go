// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package blocks

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"syscall"

	"github.com/conflowio/conflow/src/conflow"
	"github.com/conflowio/conflow/src/conflow/block"
	"github.com/conflowio/conflow/src/util"
	"github.com/conflowio/conflow/src/util/multierror"
)

// @block "task"
type Exec struct {
	// @id
	id conflow.ID
	// @required
	cmd    string
	params []string
	dir    string
	env    []string
	// @read_only
	exitCode int64
	// @generated
	stdout *Stream
	// @generated
	stderr *Stream
	// @dependency
	blockPublisher conflow.BlockPublisher
}

func (e *Exec) ID() conflow.ID {
	return e.id
}

func (e *Exec) Run(ctx context.Context) (conflow.Result, error) {
	cmd := exec.CommandContext(ctx, e.cmd, e.params...)

	if e.dir != "" {
		cmd.Dir = e.dir
	}
	if len(e.env) > 0 {
		cmd.Env = e.env
	}

	var err error
	e.stdout.Stream, err = cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout reader: %w", err)
	}
	e.stderr.Stream, err = cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr reader: %w", err)
	}

	if startErr := cmd.Start(); startErr != nil {
		return nil, fmt.Errorf("Failed to start command: %w", startErr)
	}

	wg := &util.WaitGroup{}
	wg.Run(func() error {
		_, err := e.blockPublisher.PublishBlock(e.stdout, nil)
		return err
	})
	wg.Run(func() error {
		_, err := e.blockPublisher.PublishBlock(e.stderr, nil)
		return err
	})

	var retErr error
	select {
	case <-wg.Wait():
		retErr = wg.Err()
	case <-ctx.Done():
		return nil, errors.New("aborted")
	}

	// Wait shouldn't be called until both stdin/stderr was read
	if err := cmd.Wait(); err != nil {
		e.exitCode = 256 // If we can't get the exit code at least return with a custom one
		if exitErr, ok := err.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				e.exitCode = int64(status.ExitStatus())
			}
		}
		retErr = multierror.Append(retErr, err)
	}

	return nil, retErr
}

func (e *Exec) ParseContextOverride() conflow.ParseContextOverride {
	return conflow.ParseContextOverride{
		BlockTransformerRegistry: block.InterpreterRegistry{
			"stdout": StreamInterpreter{},
			"stderr": StreamInterpreter{},
		},
	}
}
