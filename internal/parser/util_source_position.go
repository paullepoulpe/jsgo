// Port of com/google/javascript/jscomp/parsing/parser/util/SourcePosition.java @ closure-compiler v20261021

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

// SourcePosition is a position in a source string: its offset, in UTF-16
// code units. Its line and column come from the file's line table
// (SourceFile.LinePosition).
//
// Go-only: performance — Java's SourcePosition also holds the source file,
// the line and the column. Every token and parse tree has two positions,
// and most are only ever compared or turned into a node's line and column
// once, so they keep the offset alone (4 bytes instead of 12); the code
// that needs lines (comments, errors) uses LinePosition.
type SourcePosition struct {
	Offset int32
}

// GetOffset returns the offset.
//
// Java: SourcePosition.getOffset()
func (s SourcePosition) GetOffset() int {
	return int(s.Offset)
}

// NewSourcePosition returns the position at offset.
//
// Java: SourcePosition(SourceFile, int, int, int)
func NewSourcePosition(offset int) SourcePosition {
	return SourcePosition{Offset: int32(offset)}
}

// String renders the position as "@<offset>".
//
// Java: SourcePosition.toString() — Go-only: language-adapter — the line and
// column need the file (see LinePosition.String).
func (s SourcePosition) String() string {
	return jutil.StringFormat("@%d", int(s.Offset))
}

// LinePosition is a position with its line and column (0-based), as Java's
// SourcePosition has them: what comments and errors keep.
//
// Go-only: see SourcePosition.
type LinePosition struct {
	Offset int32
	Line   int32
	Column int32
}

// Position returns the position without its line and column.
func (p LinePosition) Position() SourcePosition {
	return SourcePosition{Offset: p.Offset}
}

// GetOffset returns the offset.
func (p LinePosition) GetOffset() int {
	return int(p.Offset)
}

// String renders the position as "(<line>, <column>)", 1-based.
//
// Java: SourcePosition.toString() — Go-only: language-adapter — without the
// short source name.
func (p LinePosition) String() string {
	return jutil.StringFormat("(%d, %d)", int(p.Line)+1, int(p.Column)+1)
}
