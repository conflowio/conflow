// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package schema

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Resolver interface {
	ResolveSchema(uri string) (Schema, error)
}

type ResolverFunc func(uri string) (Schema, error)

func (f ResolverFunc) ResolveSchema(uri string) (Schema, error) {
	return f(uri)
}

func ResolveFromFile(uri string) (Schema, error) {
	path := strings.TrimPrefix(uri, "file://")
	json, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema: %w", err)
	}

	s, err := UnmarshalJSON(json)
	if err != nil {
		return nil, fmt.Errorf("failed to parse schema: %w", err)
	}

	return s, nil
}

func ResolveFromURL(uri string) (Schema, error) {
	response, err := http.Get(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema: %w", err)
	}

	defer response.Body.Close()
	content, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema: %w", err)
	}

	s, err := UnmarshalJSON(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse schema: %w", err)
	}

	return s, nil
}
