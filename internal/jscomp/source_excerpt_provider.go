// Port of com/google/javascript/jscomp/SourceExcerptProvider.java @ closure-compiler v20261021

/*
 * Copyright 2007 The Closure Compiler Authors.
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

package jscomp

import "github.com/paullepoulpe/jsgo/internal/sourcemap"

// SourceExcerptProvider is responsible for building source code excerpt of
// specific locations, such as a specific line or a region around a given
// line number.
//
// Go-only: language-adapter — Java Strings that may be null (source names,
// lines, formatted excerpts) are *string; nil is null.
type SourceExcerptProvider interface {
	// GetSourceLine gets the line indicated by the line number (1 being the
	// first line of the file). This call will return only the specific line,
	// or nil if it does not exist.
	GetSourceLine(sourceName *string, lineNumber int) *string

	// GetSourceLines gets the specific lines returned by the beginning and
	// excerpt length. Must not start or end with a carriage return.
	// Implementations may decide to truncate lines but will always include
	// the first line and never be longer than the specified length.
	//
	// length is the desired length of the excerpt (in chars). If -1, returns
	// just a single line. Otherwise, returns as many lines as needed
	// (including the remainder of the last line). The result is nil if it
	// does not exist.
	GetSourceLines(sourceName *string, lineNumber int, length int) Region

	// GetSourceRegion gets a region around the indicated line number. The
	// exact definition of a region is implementation specific, but it must
	// contain the line indicated by the line number. A region must not start
	// or end by a carriage return. The result is nil if it does not exist.
	GetSourceRegion(sourceName *string, lineNumber int) Region

	// GetSourceMapping does a lookup for the original file position based on
	// this source map, if the source file is associated with an
	// inputSourceMap.
	GetSourceMapping(sourceName *string, lineNumber int, columnNumber int) *sourcemap.Mapping_OriginalMapping
}

// SourceExcerptProvider_SourceExcerpt is the source excerpt variety.
type SourceExcerptProvider_SourceExcerpt int32

const (
	// Line excerpt.
	SourceExcerptProvider_SourceExcerpt_LINE SourceExcerptProvider_SourceExcerpt = iota
	// Multiple lines excerpt.
	SourceExcerptProvider_SourceExcerpt_FULL
	// Region excerpt.
	SourceExcerptProvider_SourceExcerpt_REGION
)

var sourceExcerptProvider_SourceExcerpt_names = [...]string{"LINE", "FULL", "REGION"}

// SourceExcerptProvider_SourceExcerpt_values is SourceExcerpt.values().
func SourceExcerptProvider_SourceExcerpt_values() []SourceExcerptProvider_SourceExcerpt {
	return []SourceExcerptProvider_SourceExcerpt{
		SourceExcerptProvider_SourceExcerpt_LINE, SourceExcerptProvider_SourceExcerpt_FULL,
		SourceExcerptProvider_SourceExcerpt_REGION,
	}
}

// String is Enum.name().
func (e SourceExcerptProvider_SourceExcerpt) String() string {
	return sourceExcerptProvider_SourceExcerpt_names[e]
}

// Get gets a source excerpt string based on the type of the source excerpt.
//
// Java: SourceExcerptProvider.SourceExcerpt.get(SourceExcerptProvider, String, int, ExcerptFormatter)
func (e SourceExcerptProvider_SourceExcerpt) Get(
	source SourceExcerptProvider, sourceName *string, lineNumber int, formatter SourceExcerptProvider_ExcerptFormatter) *string {
	return e.GetWithLength(source, sourceName, lineNumber, -1, formatter)
}

// GetWithLength gets a source excerpt string based on the type of the
// source excerpt.
//
// Java: SourceExcerptProvider.SourceExcerpt.get(SourceExcerptProvider, String, int, int, ExcerptFormatter)
func (e SourceExcerptProvider_SourceExcerpt) GetWithLength(
	source SourceExcerptProvider, sourceName *string, lineNumber int, length int,
	formatter SourceExcerptProvider_ExcerptFormatter) *string {
	switch e {
	case SourceExcerptProvider_SourceExcerpt_LINE:
		return formatter.FormatLine(source.GetSourceLine(sourceName, lineNumber), lineNumber)
	case SourceExcerptProvider_SourceExcerpt_FULL:
		return formatter.FormatRegion(source.GetSourceLines(sourceName, lineNumber, length))
	default: // REGION
		return formatter.FormatRegion(source.GetSourceRegion(sourceName, lineNumber))
	}
}

// SourceExcerptProvider_ExcerptFormatter is responsible of formatting
// source excerpts.
type SourceExcerptProvider_ExcerptFormatter interface {
	// FormatLine formats a line excerpt.
	FormatLine(line *string, lineNumber int) *string

	// FormatRegion formats a region excerpt.
	FormatRegion(region Region) *string
}
