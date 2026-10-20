// Port of com/google/javascript/jscomp/parsing/parser/SourceFile.java @ closure-compiler v20261015

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

// SourceFile is a source file.
//
// Immutable.
type SourceFile struct {
	Name     string
	Contents string

	// Go-only: language-adapter — Contents as UTF-16 code units (Java chars),
	// converted once, so that the scanner and every offset work in Java char
	// indices. Read-only.
	ContentsUTF16 []uint16
	// Go-only: performance — length of the leading pure-ASCII prefix of
	// Contents where UTF-16 offsets and Go byte offsets are identical.
	asciiPrefix int
}

// NewSourceFile constructs a SourceFile.
//
// Java: SourceFile(String, String)
func NewSourceFile(name, contents string) *SourceFile {
	utf16, asciiPrefix := jutil.StringToUTF16AndASCIIPrefix(contents)
	return &SourceFile{
		Name:          name,
		Contents:      contents,
		ContentsUTF16: utf16,
		asciiPrefix:   asciiPrefix,
	}
}

// NewSourceFileInArena constructs a SourceFile using arena's reusable UTF-16
// buffer.
func NewSourceFileInArena(name, contents string, arena *ParseArena) *SourceFile {
	if arena == nil {
		return NewSourceFile(name, contents)
	}
	utf16, asciiPrefix := arena.decodeUTF16(contents)
	return &SourceFile{
		Name:          name,
		Contents:      contents,
		ContentsUTF16: utf16,
		asciiPrefix:   asciiPrefix,
	}
}
