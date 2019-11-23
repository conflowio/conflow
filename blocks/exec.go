// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package blocks

import (
	"errors"
	"os/exec"
	"syscall"

	"github.com/opsidian/basil/util/multierror"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/block"
	"github.com/opsidian/basil/util"
	"golang.org/x/xerrors"
)

//go:generate basil generate
type Exec struct {
	id       basil.ID `basil:"id"`
	cmd      string   `basil:"required"`
	params   []string
	dir      string
	env      []string
	exitCode int64   `basil:"output"`
	stdout   *Stream `basil:"generated"`
	stderr   *Stream `basil:"generated"`
}

func (e *Exec) ID() basil.ID {
	return e.id
}

func (e *Exec) Main(ctx basil.BlockContext) error {
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
		return xerrors.Errorf("failed to create stdout reader: %v", err)
	}
	e.stderr.Stream, err = cmd.StderrPipe()
	if err != nil {
		return xerrors.Errorf("failed to create stderr reader: %v", err)
	}

	if startErr := cmd.Start(); startErr != nil {
		return xerrors.Errorf("Failed to start command: %v", startErr)
	}

	wg := &util.WaitGroup{}
	wg.Run(func() error {
		_, err := ctx.PublishBlock(e.stdout, nil)
		return err
	})
	wg.Run(func() error {
		_, err := ctx.PublishBlock(e.stderr, nil)
		return err
	})

	var retErr error
	select {
	case <-wg.Wait():
		retErr = wg.Err()
	case <-ctx.Done():
		return errors.New("aborted")
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

	return retErr
}

func (e *Exec) ParseContextOverride() basil.ParseContextOverride {
	return basil.ParseContextOverride{
		BlockTransformerRegistry: block.InterpreterRegistry{
			"stdout": StreamInterpreter{},
			"stderr": StreamInterpreter{},
		},
	}
}
