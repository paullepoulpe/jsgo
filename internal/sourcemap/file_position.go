// Port of com/google/debugging/sourcemap/FilePosition.java @ closure-compiler v20261015

/*
 * Copyright 2009 The Closure Compiler Authors.
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

// FilePosition represents a position in a source file.
type FilePosition struct {
	line   int
	column int
}

// NewFilePosition constructs a FilePosition.
//
// Java: FilePosition(int, int)
func NewFilePosition(line, column int) *FilePosition {
	p := &FilePosition{}
	p.line = line
	p.column = column
	return p
}

// GetLine returns the line number of this position. Note: The v1 and v2
// source maps use a line number with the first line being 1, whereas the v3
// source map corrects this and uses a first line number of 0 to be
// consistent with the column representation.
//
// Java: FilePosition.getLine()
func (p *FilePosition) GetLine() int {
	return p.line
}

// GetColumn returns the character index on the line of this position, with
// the first column being 0.
//
// Java: FilePosition.getColumn()
func (p *FilePosition) GetColumn() int {
	return p.column
}

// String is for debugging only.
//
// Java: FilePosition.toString()
func (p *FilePosition) String() string {
	// TODO(nickreid): Make this class an @AutoValue and delete this. Too many
	// places were using the public constructor to do so when this method was
	// added.
	return jutil.MoreObjects_ToStringHelper("FilePosition").Add("line", p.line).Add("column", p.column).String()
}
