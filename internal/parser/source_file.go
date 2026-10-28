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
	"strings"
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

// ASCIISlice returns Contents[start:end] and true when [start, end) lies
// entirely within the file's pure-ASCII prefix (where UTF-16 code unit offsets
// equal byte offsets in Contents), avoiding any allocation or UTF-16 decode.
func (f *SourceFile) ASCIISlice(start, end int) (string, bool) {
	if start >= 0 && end <= f.asciiPrefix && start <= end {
		return f.Contents[start:end], true
	}
	return "", false
}

// UTF16Slice returns f's UTF-16 code units for [start, end) without allocating
// when the range lies entirely within f's UTF-16 tail.
func (f *SourceFile) UTF16Slice(start, end int) ([]uint16, bool) {
	p := f.asciiPrefix
	if start >= p && end >= start && end-p <= len(f.utf16Tail) {
		return f.utf16Tail[start-p : end-p], true
	}
	return nil, false
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

// LineStarts builds (if needed) and returns the table of 0-based UTF-16 line
// start offsets for single-threaded callers (such as IRFactory) that maintain
// their own non-atomic cursor.
func (f *SourceFile) LineStarts() []int32 {
	f.linesOnce.Do(f.buildLineStarts)
	return f.lineStarts
}

// FindLine returns the 0-based line index in starts containing offset o
// (the largest k with starts[k] <= o).
func FindLine(starts []int32, o int32) int {
	lo, hi := 0, len(starts)
	for lo < hi {
		mid := int(uint(lo+hi) >> 1)
		if starts[mid] <= o {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	return lo - 1
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
	if i >= len(starts) || o < starts[i] || (i+1 < len(starts) && starts[i+1] <= o) {
		if i+1 < len(starts) && starts[i+1] <= o && (i+2 >= len(starts) || o < starts[i+2]) {
			i++ // the next line: the common case of a scan going forward
		} else if i > 0 && i < len(starts) && o < starts[i] && starts[i-1] <= o {
			i-- // previous line
		} else {
			i = FindLine(starts, o)
		}
		f.lastLine.Store(int32(i))
	}
	return i, offset - int(starts[i])
}

func (f *SourceFile) buildLineStarts() {
	n := f.Length()
	starts := make([]int32, 1, 1+n/32)
	p := f.asciiPrefix
	ascii := f.Contents[:p]
	tail := f.utf16Tail
	tailStart := 0
	if p > 0 {
		if strings.IndexByte(ascii, '\r') < 0 {
			for pos := 0; pos < p; {
				k := strings.IndexByte(ascii[pos:], '\n')
				if k < 0 {
					break
				}
				pos += k + 1
				starts = append(starts, int32(pos))
			}
		} else {
			i := 0
			for ; i < p; i++ {
				switch ascii[i] {
				case '\n':
					starts = append(starts, int32(i+1))
				case '\r':
					if i+1 < n && f.CharAt(i+1) == '\n' {
						i++
					}
					starts = append(starts, int32(i+1))
				}
			}
			tailStart = i - p
		}
	}
	for k := tailStart; k < len(tail); k++ {
		switch tail[k] {
		case '\n', 0x2028, 0x2029:
			starts = append(starts, int32(p+k+1))
		case '\r':
			if k+1 < len(tail) && tail[k+1] == '\n' {
				k++
			}
			starts = append(starts, int32(p+k+1))
		}
	}
	f.lineStarts = starts
}
