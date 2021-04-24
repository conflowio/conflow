// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package schema

import (
	"bytes"
	"errors"
	"fmt"
)

var ErrMetadataReadOnly = errors.New("metadata is read-only")

// Metadata contains common metadata for schemas
type Metadata struct {
	Annotations map[string]string `json:"annotations,omitempty"`
	Deprecated  bool              `json:"deprecated,omitempty"`
	Description string            `json:"description,omitempty"`
	Examples    []interface{}     `json:"examples,omitempty"`
	Pointer     bool              `json:"pointer,omitempty"`
	ReadOnly    bool              `json:"readOnly,omitempty"`
	Title       string            `json:"title,omitempty"`
	WriteOnly   bool              `json:"writeOnly,omitempty"`
}

type MetadataAccessor interface {
	Merge(Metadata)
	SetAnnotation(string, string)
	SetDeprecated(bool)
	SetDescription(string)
	SetExamples([]interface{})
	SetPointer(bool)
	SetReadOnly(bool)
	SetTitle(string)
	SetWriteOnly(bool)
}

func (m *Metadata) Merge(m2 Metadata) {
	if len(m2.Annotations) > 0 && m.Annotations == nil {
		m.Annotations = map[string]string{}
	}
	for k, v := range m2.Annotations {
		m.Annotations[k] = v
	}
	if m2.Deprecated {
		m.Deprecated = true
	}
	if m2.Description != "" {
		m.Description = m2.Description
	}
	if m2.Examples != nil {
		m.Examples = m2.Examples
	}
	if m2.Pointer {
		m.Pointer = true
	}
	if m2.ReadOnly {
		m.ReadOnly = true
	}
	if m2.Title != "" {
		m.Title = m2.Title
	}
	if m2.WriteOnly {
		m.WriteOnly = true
	}
}

func (m *Metadata) GetAnnotation(name string) string {
	return m.Annotations[name]
}

func (m *Metadata) SetAnnotation(name, value string) {
	if m.Annotations == nil {
		m.Annotations = map[string]string{}
	}

	if value == "" {
		delete(m.Annotations, name)
		return
	}

	m.Annotations[name] = value
}

func (m *Metadata) GetDeprecated() bool {
	return m.Deprecated
}

func (m *Metadata) SetDeprecated(deprecated bool) {
	m.Deprecated = deprecated
}

func (m *Metadata) GetDescription() string {
	return m.Description
}

func (m *Metadata) SetDescription(description string) {
	m.Description = description
}

func (m Metadata) GetExamples() []interface{} {
	return m.Examples
}

func (m *Metadata) SetExamples(examples []interface{}) {
	m.Examples = examples
}

func (m *Metadata) GetReadOnly() bool {
	return m.ReadOnly
}

func (m *Metadata) SetReadOnly(readOnly bool) {
	m.ReadOnly = readOnly
}

func (m *Metadata) GetPointer() bool {
	return m.Pointer
}

func (m *Metadata) SetPointer(pointer bool) {
	m.Pointer = pointer
}

func (m *Metadata) GetTitle() string {
	return m.Title
}

func (m *Metadata) SetTitle(title string) {
	m.Title = title
}

func (m *Metadata) GetWriteOnly() bool {
	return m.WriteOnly
}

func (m *Metadata) SetWriteOnly(writeOnly bool) {
	m.WriteOnly = writeOnly
}

func (m *Metadata) GoString() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("schema.Metadata{\n")
	if len(m.Annotations) > 0 {
		_, _ = fmt.Fprintf(buf, "\tAnnotations: %#v,\n", m.Annotations)
	}
	if m.Deprecated {
		_, _ = fmt.Fprintf(buf, "\tDeprecated: %#v,\n", m.Deprecated)
	}
	if len(m.Description) > 0 {
		_, _ = fmt.Fprintf(buf, "\tDescription: %#v,\n", m.Description)
	}
	if len(m.Examples) > 0 {
		_, _ = fmt.Fprintf(buf, "\tExamples: %#v,\n", m.Examples)
	}
	if m.Pointer {
		_, _ = fmt.Fprintf(buf, "\tPointer: %#v,\n", m.Pointer)
	}
	if m.ReadOnly {
		buf.WriteString("\tReadOnly: true,\n")
	}
	if len(m.Title) > 0 {
		_, _ = fmt.Fprintf(buf, "\tTitle: %#v,\n", m.Title)
	}
	if m.WriteOnly {
		buf.WriteString("\tWriteOnly: true,\n")
	}

	buf.WriteRune('}')
	return buf.String()
}

type emptyMetadata struct {
}

func (e emptyMetadata) GetAnnotation(string) string {
	return ""
}

func (e emptyMetadata) GetDeprecated() bool {
	return false
}

func (e emptyMetadata) GetDescription() string {
	return ""
}

func (e emptyMetadata) GetExamples() []interface{} {
	return nil
}

func (e emptyMetadata) GetPointer() bool {
	return false
}

func (e emptyMetadata) GetReadOnly() bool {
	return false
}

func (e emptyMetadata) GetTitle() string {
	return ""
}

func (e emptyMetadata) GetWriteOnly() bool {
	return false
}

func HasAnnotationValue(s Schema, name, value string) bool {
	return s.GetAnnotation(name) == value
}

func UpdateMetadata(s Schema, f func(meta MetadataAccessor) error) error {
	if meta, ok := s.(MetadataAccessor); ok {
		return f(meta)
	}

	return ErrMetadataReadOnly
}
