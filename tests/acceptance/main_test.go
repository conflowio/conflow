// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package acceptance_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	uzerolog "github.com/rs/zerolog"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/job"
	"github.com/opsidian/basil/examples/common"
	"github.com/opsidian/basil/loggers/zerolog"
	"github.com/opsidian/basil/parsers"
	"github.com/opsidian/basil/tests/acceptance"
	"github.com/opsidian/basil/util"
)

type testCase struct {
	fixture []byte
	out     []byte
}

type testCases map[string]*testCase

func (t testCases) get(name string) *testCase {
	tc, ok := t[name]
	if ok {
		return tc
	}

	tc = &testCase{}
	t[name] = tc
	return tc
}

func readTestCases(dir string, name string, testCases testCases) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		if f.IsDir() {
			if err := readTestCases(path.Join(dir, f.Name()), path.Join(name, f.Name()), testCases); err != nil {
				return err
			}
		}
		if !f.IsDir() {
			switch {
			case strings.HasSuffix(f.Name(), ".basil"):
				tcName := path.Join(name, strings.TrimSuffix(f.Name(), ".basil"))
				tc := testCases.get(tcName)
				tc.fixture, err = ioutil.ReadFile(path.Join(dir, f.Name()))
				if err != nil {
					return err
				}
			case strings.HasSuffix(f.Name(), ".out"):
				tcName := path.Join(name, strings.TrimSuffix(f.Name(), ".out"))
				tc := testCases.get(tcName)
				tc.out, err = ioutil.ReadFile(path.Join(dir, f.Name()))
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

var _ = Describe("Acceptance tests", func() {

	var logger basil.Logger
	var scheduler *job.Scheduler

	runTest := func(input string) (string, error) {
		ctx, cancel := util.CreateDefaultContext()
		defer cancel()

		parseCtx := common.NewParseContext()

		p := parsers.NewMain("main", acceptance.MainInterpreter{})

		if err := p.ParseText(parseCtx, input); err != nil {
			return "", err
		}

		node, ok := parseCtx.BlockNode("main")
		if !ok {
			return "", fmt.Errorf("block main does not exist")
		}

		stdout := bytes.NewBuffer(make([]byte, 0, 256))

		evalContext := basil.NewEvalContext(ctx, nil, logger, scheduler, nil)
		evalContext.SetStdout(stdout)

		_, err := node.Value(evalContext)
		if err != nil {
			return "", parseCtx.FileSet().ErrorWithPosition(err)
		}

		return stdout.String(), nil
	}

	BeforeSuite(func() {
		level := uzerolog.InfoLevel
		if envLevel := os.Getenv("BASIL_LOG"); envLevel != "" {
			var err error
			level, err = uzerolog.ParseLevel(envLevel)
			if err != nil {
				panic(fmt.Errorf("invalid log level %q", envLevel))
			}
		}

		logger = zerolog.NewConsoleLogger(level)
		scheduler = job.NewScheduler(logger, runtime.NumCPU()*2, 100)
		scheduler.Start()
	})

	AfterSuite(func() {
		scheduler.Stop()
	})

	testCases := map[string]*testCase{}
	_, filename, _, _ := runtime.Caller(0)
	err := readTestCases(path.Join(path.Dir(filename), "fixtures"), "", testCases)
	if err != nil {
		panic(err)
	}

	for name, tc := range testCases {
		if tc.out != nil {
			It(name, func() {
				out, err := runTest(string(tc.fixture))
				if err != nil {
					Fail(err.Error())
				}
				Expect(out).To(Equal(string(tc.out)))
			})
		}
	}
})
