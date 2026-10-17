// Port of com/google/javascript/jscomp/parsing/parser/util/SourcePosition.java @ closure-compiler v20260415

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

package parser

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// SourcePosition is a position in a source string - includes offset, line
// and column.
type SourcePosition struct {
	Source *SourceFile
	Offset int
	Line   int
	Column int
}

// GetOffset returns the offset.
//
// Java: SourcePosition.getOffset()
func (s *SourcePosition) GetOffset() int {
	return s.Offset
}

// NewSourcePosition constructs a SourcePosition.
//
// Java: SourcePosition(SourceFile, int, int, int)
func NewSourcePosition(source *SourceFile, offset, line, column int) *SourcePosition {
	s := &SourcePosition{}
	s.Source = source
	s.Offset = offset
	s.Line = line
	s.Column = column
	return s
}

// String renders the position as "<file>(<line>, <column>)", 1-based.
//
// Java: SourcePosition.toString()
func (s *SourcePosition) String() string {
	return jutil.StringFormat("%s(%d, %d)", s.ShortSourceName(), s.Line+1, s.Column+1)
}

// ShortSourceName returns the last '/'-separated segment of the source name.
//
// Java: SourcePosition.shortSourceName()
func (s *SourcePosition) ShortSourceName() string {
	if s.Source == nil {
		return ""
	}

	// Java: Iterables.getLast(Splitter.on('/').split(source.name))
	return s.Source.Name[strings.LastIndexByte(s.Source.Name, '/')+1:]
}
