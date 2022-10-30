// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package schema

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Resolver interface {
	ResolveSchema(ctx context.Context, uri string) (Schema, error)
}

type ResolverFunc func(ctx context.Context, uri string) (Schema, error)

func (f ResolverFunc) ResolveSchema(ctx context.Context, uri string) (Schema, error) {
	return f(ctx, uri)
}

func ResolveFromFile(ctx context.Context, uri string) (Schema, error) {
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

func ResolveFromURL(ctx context.Context, uri string) (Schema, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema: %w", err)
	}

	response, err := http.DefaultClient.Do(req)
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
