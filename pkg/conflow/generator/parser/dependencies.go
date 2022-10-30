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
	DependencyNode           = "node"
	DependencyStdout         = "stdout"
	DependencyUserContext    = "userContext"
)

var validDependencies = []string{
	DependencyBlockPublisher,
	DependencyJobScheduler,
	DependencyLogger,
	DependencyNode,
	DependencyStdout,
	DependencyUserContext,
}

var dependencyTypes = map[string]string{
	DependencyBlockPublisher: "github.com/conflowio/conflow/pkg/conflow.BlockPublisher",
	DependencyJobScheduler:   "github.com/conflowio/conflow/pkg/conflow.JobScheduler",
	DependencyLogger:         "github.com/conflowio/conflow/pkg/conflow.Logger",
	DependencyNode:           "github.com/conflowio/parsley/parsley.Node",
	DependencyStdout:         "io.Writer",
	DependencyUserContext:    "interface{}",
}

type Dependency struct {
	Name      string
	FieldName string
}
