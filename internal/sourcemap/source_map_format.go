// Port of com/google/debugging/sourcemap/SourceMapFormat.java @ closure-compiler v20261021

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

// SourceMapFormat is a list of currently support SourceMap format
// revisions.
type SourceMapFormat int32

const (
	// The latest "stable" format
	SourceMapFormat_DEFAULT SourceMapFormat = iota

	// V3: A nice compact format
	SourceMapFormat_V3
)

var sourceMapFormat_names = [...]string{"DEFAULT", "V3"}

// SourceMapFormat_values is SourceMapFormat.values().
func SourceMapFormat_values() []SourceMapFormat {
	return []SourceMapFormat{SourceMapFormat_DEFAULT, SourceMapFormat_V3}
}

// String is Enum.name().
func (f SourceMapFormat) String() string { return sourceMapFormat_names[f] }
