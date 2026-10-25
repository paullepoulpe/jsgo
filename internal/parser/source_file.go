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
	"sort"
	"sync"
	"sync/atomic"

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

	// Go-only: performance — the offset at which each line starts, built on
	// first use; positions are offsets (SourcePosition) and their lines and
	// columns come from here.
	linesOnce  sync.Once
	lineStarts []int32
	lastLine   atomic.Int32 // the line of the last lookup: lookups mostly move forward
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

// LinePosition returns pos with its line and column.
//
// Go-only: see SourcePosition.
func (f *SourceFile) LinePosition(pos SourcePosition) LinePosition {
	line, column := f.LineColumn(int(pos.Offset))
	return LinePosition{Offset: pos.Offset, Line: int32(line), Column: int32(column)}
}

// LineRange returns the range from start to end with lines and columns.
func (f *SourceFile) LineRange(start, end int) LineRange {
	return LineRange{
		Start: f.LinePosition(NewSourcePosition(start)),
		End:   f.LinePosition(NewSourcePosition(end)),
	}
}

// Line returns the 0-based line of pos.
func (f *SourceFile) Line(pos SourcePosition) int {
	line, _ := f.LineColumn(int(pos.Offset))
	return line
}

// LineColumn returns the 0-based line and column of offset. A line ends
// after a line terminator: \n, \r, \r\n, U+2028 or U+2029.
func (f *SourceFile) LineColumn(offset int) (line, column int) {
	f.linesOnce.Do(f.buildLineStarts)
	starts := f.lineStarts
	o := int32(offset)
	i := int(f.lastLine.Load())
	if i >= len(starts) || starts[i] > o || (i+1 < len(starts) && starts[i+1] <= o) {
		if i+1 < len(starts) && starts[i+1] <= o && (i+2 >= len(starts) || starts[i+2] > o) {
			i++ // the next line: the common case of a scan going forward
		} else {
			i = sort.Search(len(starts), func(k int) bool { return starts[k] > o }) - 1
		}
		f.lastLine.Store(int32(i))
	}
	return i, offset - int(starts[i])
}

func (f *SourceFile) buildLineStarts() {
	starts := make([]int32, 1, 1+f.Length()/32)
	n := f.Length()
	for i := 0; i < n; i++ {
		switch ch := f.CharAt(i); ch {
		case '\n', 0x2028, 0x2029:
			starts = append(starts, int32(i+1))
		case '\r':
			if i+1 < n && f.CharAt(i+1) == '\n' {
				i++
			}
			starts = append(starts, int32(i+1))
		}
	}
	f.lineStarts = starts
}
