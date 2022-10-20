// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package util

import (
	"bytes"
	"text/template"

	"github.com/go-openapi/inflect"

	"github.com/conflowio/conflow/src/internal/utils"
)

func GenerateTemplate(tmpl string, params interface{}, funcMap template.FuncMap) ([]byte, error) {
	funcs := template.FuncMap{}
	funcs["asciify"] = inflect.Asciify
	funcs["camelize"] = inflect.Camelize
	funcs["camelizeDownFirst"] = inflect.CamelizeDownFirst
	funcs["capitalize"] = inflect.Capitalize
	funcs["dasherize"] = inflect.Dasherize
	funcs["foreignKey"] = inflect.ForeignKey
	funcs["foreignKeyCondensed"] = inflect.ForeignKeyCondensed
	funcs["humanize"] = inflect.Humanize
	funcs["ordinalize"] = inflect.Ordinalize
	funcs["rarameterize"] = inflect.Parameterize
	funcs["rarameterizeJoin"] = inflect.ParameterizeJoin
	funcs["pluralize"] = inflect.Pluralize
	funcs["singularize"] = inflect.Singularize
	funcs["tableize"] = inflect.Tableize
	funcs["titleize"] = inflect.Titleize
	funcs["typeify"] = inflect.Typeify
	funcs["underscore"] = inflect.Underscore

	funcs["ensureUniqueGoPackageSelector"] = utils.EnsureUniqueGoPackageSelector

	for n, f := range funcMap {
		funcs[n] = f
	}

	t := template.New("main")
	t.Funcs(funcs)
	t.Option("missingkey=error")

	if _, err := t.Parse(tmpl); err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	if err := t.Execute(buf, params); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
