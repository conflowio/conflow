// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parser

const (
	DependencyBlockPublisher = "blockPublisher"
	DependencyJobScheduler   = "jobScheduler"
	DependencyLogger         = "logger"
	DependencyStdout         = "stdout"
	DependencyUserContext    = "userContext"
)

var validDependencies = []string{
	DependencyBlockPublisher,
	DependencyJobScheduler,
	DependencyLogger,
	DependencyStdout,
	DependencyUserContext,
}

var dependencyTypes = map[string]string{
	DependencyBlockPublisher: "github.com/opsidian/basil/basil.BlockPublisher",
	DependencyJobScheduler:   "github.com/opsidian/basil/basil.JobScheduler",
	DependencyLogger:         "github.com/opsidian/basil/basil.Logger",
	DependencyStdout:         "io.Writer",
	DependencyUserContext:    "interface{}",
}

type Dependency struct {
	Name      string
	FieldName string
}
