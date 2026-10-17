// Port of com/google/debugging/sourcemap/SourceMapSection.java @ closure-compiler v20261015

/*
 * Copyright 2011 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package sourcemap

import "github.com/paullepoulpe/jsgo/internal/jutil"

// SourceMapSection is a class representing a partial source map (a Java
// record).
type SourceMapSection struct {
	type_  SourceMapSection_SectionType
	value  *string
	line   int
	column int
}

// SourceMapSection_SectionType is SourceMapSection.SectionType.
type SourceMapSection_SectionType int32

const (
	SourceMapSection_SectionType_URL SourceMapSection_SectionType = iota
	SourceMapSection_SectionType_MAP
)

var sourceMapSection_SectionType_names = [...]string{"URL", "MAP"}

// SourceMapSection_SectionType_values is SectionType.values().
func SourceMapSection_SectionType_values() []SourceMapSection_SectionType {
	return []SourceMapSection_SectionType{SourceMapSection_SectionType_URL, SourceMapSection_SectionType_MAP}
}

// String is Enum.name().
func (t SourceMapSection_SectionType) String() string {
	return sourceMapSection_SectionType_names[t]
}

// NewSourceMapSection is the record's canonical constructor.
//
// Java: SourceMapSection(SectionType, String, int, int)
func NewSourceMapSection(type_ SourceMapSection_SectionType, value *string, line, column int) *SourceMapSection {
	return &SourceMapSection{type_: type_, value: value, line: line, column: column}
}

// SourceMapSection_ForMap creates a section holding a source map.
//
// Java: SourceMapSection.forMap(String, int, int)
func SourceMapSection_ForMap(value string, line, column int) *SourceMapSection {
	return NewSourceMapSection(SourceMapSection_SectionType_MAP, &value, line, column)
}

// SourceMapSection_ForURL creates a section referring to a source map by
// URL.
//
// Java: SourceMapSection.forURL(String, int, int)
func SourceMapSection_ForURL(value string, line, column int) *SourceMapSection {
	return NewSourceMapSection(SourceMapSection_SectionType_URL, &value, line, column)
}

// Type is the record accessor.
//
// Java: SourceMapSection.type()
func (s *SourceMapSection) Type() SourceMapSection_SectionType {
	return s.type_
}

// Value is the record accessor.
//
// Java: SourceMapSection.value()
func (s *SourceMapSection) Value() *string {
	return s.value
}

// Line is the record accessor.
//
// Java: SourceMapSection.line()
func (s *SourceMapSection) Line() int {
	return s.line
}

// Column is the record accessor.
//
// Java: SourceMapSection.column()
func (s *SourceMapSection) Column() int {
	return s.column
}

// GetSectionType returns the section type.
//
// Java: SourceMapSection.getSectionType()
func (s *SourceMapSection) GetSectionType() SourceMapSection_SectionType {
	return s.type_
}

// GetSectionValue returns the value that represents the map for this
// section. This is either a URL for a valid source map file that represents
// a section of a generate source file such as when multiple files are
// concatenated together, or a string representing a JSON object that is a
// source map.
//
// Java: SourceMapSection.getSectionValue()
func (s *SourceMapSection) GetSectionValue() *string {
	return s.value
}

// GetLine returns the starting line for this section.
//
// Java: SourceMapSection.getLine()
func (s *SourceMapSection) GetLine() int {
	return s.line
}

// GetColumn returns the starting column for this section.
//
// Java: SourceMapSection.getColumn()
func (s *SourceMapSection) GetColumn() int {
	return s.column
}

// Equals is the record's component-wise equals. (The record's hashCode is
// not ported: it hashes the enum constant by identity, which is not
// reproducible.)
//
// Java: SourceMapSection.equals(Object)
func (s *SourceMapSection) Equals(o any) bool {
	other, ok := o.(*SourceMapSection)
	if !ok || other == nil {
		return false
	}
	return s.type_ == other.type_ && stringsEqual(s.value, other.value) &&
		s.line == other.line && s.column == other.column
}

// Go-only: language-adapter — Objects.equals on nullable Strings.
func stringsEqual(a, b *string) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

// String is the record's toString.
//
// Java: SourceMapSection.toString()
func (s *SourceMapSection) String() string {
	return "SourceMapSection[type=" + s.type_.String() + ", value=" + jutil.ToJavaString(s.value) +
		", line=" + jutil.ToJavaString(s.line) + ", column=" + jutil.ToJavaString(s.column) + "]"
}
