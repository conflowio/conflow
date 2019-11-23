// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"errors"
	"os/exec"
	"syscall"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/block"
	"github.com/opsidian/basil/blocks"
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
	exitCode int64          `basil:"output"`
	stdout   *blocks.Stream `basil:"generated"`
	stderr   *blocks.Stream `basil:"generated"`
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
	e.stdout.Reader, err = cmd.StdoutPipe()
	if err != nil {
		return xerrors.Errorf("failed to create stdout reader: %v", err)
	}
	e.stderr.Reader, err = cmd.StderrPipe()
	if err != nil {
		return xerrors.Errorf("failed to create stderr reader: %v", err)
	}

	if startErr := cmd.Start(); startErr != nil {
		return xerrors.Errorf("Failed to start command: %v", startErr)
	}

	wg := &util.WaitGroup{}

	wg.Run(func() error {
		return ctx.PublishBlock(e.stdout, nil)
	})

	wg.Run(func() error {
		return ctx.PublishBlock(e.stderr, nil)
	})

	wg.Run(func() error {
		err := cmd.Wait()
		if err != nil {
			e.exitCode = 256 // If we can't get the exit code at least return with a custom one
			if exitErr, ok := err.(*exec.ExitError); ok {
				if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
					e.exitCode = int64(status.ExitStatus())
				}
			}
		}
		return err
	})

	select {
	case <-wg.Wait():
		if err := wg.Err(); err != nil {
			return err
		}
	case <-ctx.Done():
		return errors.New("aborted")
	}

	return nil
}

func (e *Exec) ParseContextOverride() basil.ParseContextOverride {
	return basil.ParseContextOverride{
		BlockTransformerRegistry: block.InterpreterRegistry{
			"stdout": blocks.StreamInterpreter{},
			"stderr": blocks.StreamInterpreter{},
		},
	}
}
