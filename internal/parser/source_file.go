// Port of com/google/javascript/jscomp/parsing/parser/SourceFile.java @ closure-compiler v20261021

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
	"slices"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// SourceFile is a source file.
//
// Immutable.
type SourceFile struct {
	Name     string
	Contents string

	// Go-only: language-adapter — Java indexes Contents by UTF-16 code unit
	// (char). In the leading pure-ASCII part of Contents a char index is a
	// byte index, so it is read in place; only the rest is converted to
	// UTF-16, once, into utf16Tail (nil for an all-ASCII file).
	asciiPrefix int
	utf16Tail   []uint16
}

// NewSourceFile constructs a SourceFile.
//
// Java: SourceFile(String, String)
func NewSourceFile(name, contents string) *SourceFile {
	f := &SourceFile{Name: name, Contents: contents, asciiPrefix: jutil.ASCIIPrefixLen(contents)}
	if f.asciiPrefix < len(contents) {
		tail := jutil.StringToUTF16(contents[f.asciiPrefix:])
		if cap(tail)-len(tail) > len(tail)/8 {
			tail = slices.Clone(tail) // drop the slack of multi-byte characters
		}
		f.utf16Tail = tail
	}
	return f
}

// Length is contents.length(): the number of UTF-16 code units.
//
// Go-only: language-adapter — Java String API on the contents.
func (f *SourceFile) Length() int {
	return f.asciiPrefix + len(f.utf16Tail)
}

// CharAt is contents.charAt(index), for 0 <= index < Length().
//
// Go-only: language-adapter — Java String API on the contents.
func (f *SourceFile) CharAt(index int) uint16 {
	if index < f.asciiPrefix {
		return uint16(f.Contents[index])
	}
	return f.utf16Tail[index-f.asciiPrefix]
}

// Substring is contents.substring(start, end).
//
// Go-only: language-adapter — Java String API on the contents.
func (f *SourceFile) Substring(start, end int) string {
	if start < 0 || end > f.Length() || start > end {
		panic(jutil.NewStringIndexOutOfBoundsException(start, end))
	}
	p := f.asciiPrefix
	switch {
	case end <= p:
		return f.Contents[start:end]
	case start >= p:
		return jutil.StringFromUTF16(f.utf16Tail[start-p : end-p])
	default:
		return f.Contents[start:p] + jutil.StringFromUTF16(f.utf16Tail[:end-p])
	}
}
