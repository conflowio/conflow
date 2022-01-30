// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conflow

import "time"

// Retryable is a common interface for a retryable result or error
type Retryable interface {
	Retry() bool
	RetryAfter() time.Duration
}

// Result is an interface for returning the result of a block run
type Result interface {
	Retryable
	RetryReason() string
}

// Retry is a block run result which indicates we need to retry
type Retry string

func (r Retry) Retry() bool {
	return true
}

func (r Retry) RetryAfter() time.Duration {
	return 0
}

func (r Retry) RetryReason() string {
	return string(r)
}

// RetryAfter is a block run result which indicates we need to retry after the given time
func RetryAfter(duration time.Duration, reason string) Result {
	return retryAfter{
		Duration: duration,
		Reason:   reason,
	}
}

type retryAfter struct {
	Duration time.Duration
	Reason   string
}

func (r retryAfter) Retry() bool {
	return true
}

func (r retryAfter) RetryAfter() time.Duration {
	return r.Duration
}

func (r retryAfter) RetryReason() string {
	return r.Reason
}

// RetryableError wraps the given error and makes it retryable
func RetryableError(err error, duration time.Duration) error {
	return retryableError{
		err:      err,
		duration: duration,
	}
}

type retryableError struct {
	err      error
	duration time.Duration
}

func (r retryableError) Retry() bool {
	return true
}

func (r retryableError) RetryAfter() time.Duration {
	return r.duration
}

func (r retryableError) Error() string {
	return r.err.Error()
}

func (r retryableError) UnWrap() error {
	return r.err
}
