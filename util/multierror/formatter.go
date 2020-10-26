// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package multierror

import (
	"fmt"
	"strings"
)

type Formatter interface {
	FormatMultiError([]error) string
}

type FormatterFunc func([]error) string

func (f FormatterFunc) FormatMultiError(errors []error) string {
	return f(errors)
}

func DefaultFormatter(errors []error) string {
	if len(errors) == 1 {
		return errors[0].Error()
	}

	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%d errors occurred:", len(errors)))
	for _, err := range errors {
		sb.WriteString("\n\t* ")
		sb.WriteString(err.Error())
	}
	sb.WriteString("\n")
	return sb.String()
}
