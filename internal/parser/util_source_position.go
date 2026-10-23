// Port of com/google/javascript/jscomp/parsing/parser/util/SourcePosition.java @ closure-compiler v20261015

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

import "github.com/paullepoulpe/jsgo/internal/jutil"

// SourcePosition is a position in a source string - includes offset, line
// and column.
//
// Go-only: performance — a small pointer-free value (int32 fields, no
// SourceFile): a parse creates one or two per token, and the file is always
// known from context (the scanner, the parser, the error reporter).
type SourcePosition struct {
	Offset int32
	Line   int32
	Column int32
}

// GetOffset returns the offset.
//
// Java: SourcePosition.getOffset()
func (s SourcePosition) GetOffset() int {
	return int(s.Offset)
}

// NewSourcePosition constructs a SourcePosition.
//
// Java: SourcePosition(SourceFile, int, int, int)
func NewSourcePosition(offset, line, column int) SourcePosition {
	return SourcePosition{Offset: int32(offset), Line: int32(line), Column: int32(column)}
}

// String renders the position as "(<line>, <column>)", 1-based.
//
// Java: SourcePosition.toString() — Go-only: language-adapter — without the
// short source name (see SourcePosition).
func (s SourcePosition) String() string {
	return jutil.StringFormat("(%d, %d)", int(s.Line)+1, int(s.Column)+1)
}
