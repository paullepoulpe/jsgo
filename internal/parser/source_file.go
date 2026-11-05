// Port of com/google/javascript/jscomp/parsing/parser/SourceFile.java @ closure-compiler v20261028

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
	"sync"
	"sync/atomic"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

const (
	utf16CheckpointShift  = 5
	utf16CheckpointStride = 1 << utf16CheckpointShift
)

// SourceFile is a source file.
//
// Immutable.
type SourceFile struct {
	Name     string
	Contents string

	// Go-only: language-adapter — Java indexes Contents by UTF-16 code unit
	// (char). Both the leading pure-ASCII prefix [0, asciiPrefix) and trailing
	// pure-ASCII suffix [asciiSuffixUTF16Start, utf16Len) (starting at byte
	// asciiSuffixByteStart in Contents) map in O(1); any non-ASCII middle span
	// uses an O(1)-indexed table of byte offsets every 32 UTF-16 code units
	// (checkpoints[k] = byteOff, with bit 31 set if asciiPrefix + 32*k lands on
	// the low surrogate of the 4-byte astral UTF-8 rune starting at byteOff).
	asciiPrefix           int
	asciiSuffixUTF16Start int
	asciiSuffixByteStart  int
	utf16Len              int
	checkpoints           []uint32

	// Go-only: performance — the offset at which each line starts, built on
	// first use; positions are offsets (SourcePosition) and their lines and
	// columns come from here.
	linesOnce  sync.Once
	lineStarts []int32
	lastLine   atomic.Int32 // the line of the last lookup: lookups mostly move forward
}

func decodeWTF8At(s string, i int) (r rune, size int) {
	c := s[i]
	if c < 0x80 {
		return rune(c), 1
	}
	// Generalized UTF-8 surrogate: ED A0..BF xx
	if c == 0xED && i+3 <= len(s) && s[i+1] >= 0xA0 && s[i+1] <= 0xBF && s[i+2]&0xC0 == 0x80 {
		return rune(c&0x0F)<<12 | rune(s[i+1]&0x3F)<<6 | rune(s[i+2]&0x3F), 3
	}
	return utf8.DecodeRuneInString(s[i:])
}

func utf16LenWTF8(s string) int {
	p := jutil.ASCIIPrefixLen(s)
	if p == len(s) {
		return p
	}
	n := p
	for i := p; i < len(s); {
		for i+8 <= len(s) {
			if (s[i] | s[i+1] | s[i+2] | s[i+3] | s[i+4] | s[i+5] | s[i+6] | s[i+7]) >= 0x80 {
				break
			}
			i += 8
			n += 8
		}
		if i < len(s) && s[i] < 0x80 {
			i++
			n++
			continue
		}
		if i >= len(s) {
			break
		}
		r, size := decodeWTF8At(s, i)
		i += size
		if r >= 0x10000 {
			n += 2
		} else {
			n++
		}
	}
	return n
}

func buildUTF16Checkpoints(contents string, startByte, endByte int) (midUTF16Len int, checkpoints []uint32) {
	spanLen := endByte - startByte
	if spanLen <= utf16CheckpointStride {
		return utf16LenWTF8(contents[startByte:endByte]), nil
	}
	checkpoints = make([]uint32, 1, 1+spanLen/utf16CheckpointStride)
	checkpoints[0] = uint32(startByte)
	uRel := 0
	nextU := utf16CheckpointStride
	for i := startByte; i < endByte; {
		for i+8 <= endByte && uRel+8 <= nextU {
			if (contents[i] | contents[i+1] | contents[i+2] | contents[i+3] |
				contents[i+4] | contents[i+5] | contents[i+6] | contents[i+7]) >= 0x80 {
				break
			}
			i += 8
			uRel += 8
			if uRel == nextU {
				checkpoints = append(checkpoints, uint32(i))
				nextU += utf16CheckpointStride
			}
		}
		for i < endByte && uRel < nextU && contents[i] < 0x80 {
			i++
			uRel++
			if uRel == nextU {
				checkpoints = append(checkpoints, uint32(i))
				nextU += utf16CheckpointStride
			}
		}
		if i >= endByte {
			break
		}
		if contents[i] < 0x80 {
			continue
		}
		r, size := decodeWTF8At(contents, i)
		if r >= 0x10000 {
			if uRel+1 == nextU {
				checkpoints = append(checkpoints, uint32(i)|0x80000000)
				nextU += utf16CheckpointStride
			}
			i += size
			uRel += 2
			if uRel == nextU {
				checkpoints = append(checkpoints, uint32(i))
				nextU += utf16CheckpointStride
			}
		} else {
			i += size
			uRel++
			if uRel == nextU {
				checkpoints = append(checkpoints, uint32(i))
				nextU += utf16CheckpointStride
			}
		}
	}
	return uRel, checkpoints
}

func findASCIISuffixByteStart(s string, minIdx int) int {
	i := len(s)
	for i-minIdx >= 8 {
		b := s[i-8 : i]
		if (b[0] | b[1] | b[2] | b[3] | b[4] | b[5] | b[6] | b[7]) >= 0x80 {
			break
		}
		i -= 8
	}
	for i > minIdx && s[i-1] < 0x80 {
		i--
	}
	return i
}

// NewSourceFile constructs a SourceFile.
//
// Java: SourceFile(String, String)
func NewSourceFile(name, contents string) *SourceFile {
	p := jutil.ASCIIPrefixLen(contents)
	if p == len(contents) {
		return &SourceFile{
			Name:                  name,
			Contents:              contents,
			asciiPrefix:           p,
			asciiSuffixUTF16Start: p,
			asciiSuffixByteStart:  p,
			utf16Len:              p,
		}
	}
	suffixByteStart := findASCIISuffixByteStart(contents, p+1)
	midUTF16Len, checkpoints := buildUTF16Checkpoints(contents, p, suffixByteStart)
	suffixUTF16Start := p + midUTF16Len
	return &SourceFile{
		Name:                  name,
		Contents:              contents,
		asciiPrefix:           p,
		asciiSuffixUTF16Start: suffixUTF16Start,
		asciiSuffixByteStart:  suffixByteStart,
		utf16Len:              suffixUTF16Start + (len(contents) - suffixByteStart),
		checkpoints:           checkpoints,
	}
}

// byteOffsetForUTF16 maps a UTF-16 offset uOff in [0, f.utf16Len] to
// (byteOff, lowSurrogate): if uOff lands on the trailing low surrogate of a
// 4-byte supplementary rune, lowSurrogate != 0 and byteOff points at the lead
// 0xF0..0xF4 byte of that 4-byte sequence; otherwise lowSurrogate == 0 and
// byteOff is the exact byte boundary in f.Contents.
func (f *SourceFile) byteOffsetForUTF16(uOff int) (byteOff int, lowSurrogate uint16) {
	if uOff <= f.asciiPrefix {
		return uOff, 0
	}
	if uOff >= f.asciiSuffixUTF16Start {
		return f.asciiSuffixByteStart + (uOff - f.asciiSuffixUTF16Start), 0
	}
	uPos := f.asciiPrefix
	bPos := f.asciiPrefix
	s := f.Contents
	if cps := f.checkpoints; len(cps) > 0 {
		k := (uOff - f.asciiPrefix) >> utf16CheckpointShift
		entry := cps[k]
		bPos = int(entry & 0x7fffffff)
		uPos = f.asciiPrefix + (k << utf16CheckpointShift)
		if entry&0x80000000 != 0 {
			if uPos == uOff {
				r, _ := decodeWTF8At(s, bPos)
				_, loSurr := utf16.EncodeRune(r)
				return bPos, uint16(loSurr)
			}
			bPos += 4
			uPos++
		}
	}
	rem := uOff - uPos
	if rem == 0 {
		return bPos, 0
	}
	if bPos+rem <= len(s) && jutil.IsASCII(s[bPos:bPos+rem]) {
		return bPos + rem, 0
	}
	for uPos < uOff {
		if s[bPos] < 0x80 {
			uPos++
			bPos++
			continue
		}
		r, size := decodeWTF8At(s, bPos)
		if r >= 0x10000 {
			if uPos+1 == uOff {
				_, loSurr := utf16.EncodeRune(r)
				return bPos, uint16(loSurr)
			}
			bPos += size
			uPos += 2
		} else {
			bPos += size
			uPos++
		}
	}
	return bPos, 0
}

// Length is contents.length(): the number of UTF-16 code units.
//
// Go-only: language-adapter — Java String API on the contents.
func (f *SourceFile) Length() int {
	return f.utf16Len
}

// CharAt is contents.charAt(index), for 0 <= index < Length().
//
// Go-only: language-adapter — Java String API on the contents.
func (f *SourceFile) CharAt(index int) uint16 {
	if index < f.asciiPrefix {
		return uint16(f.Contents[index])
	}
	if index >= f.asciiSuffixUTF16Start {
		return uint16(f.Contents[f.asciiSuffixByteStart+(index-f.asciiSuffixUTF16Start)])
	}
	bOff, loSurr := f.byteOffsetForUTF16(index)
	if loSurr != 0 {
		return loSurr
	}
	if f.Contents[bOff] < 0x80 {
		return uint16(f.Contents[bOff])
	}
	r, _ := decodeWTF8At(f.Contents, bOff)
	if r >= 0x10000 {
		hi, _ := utf16.EncodeRune(r)
		return uint16(hi)
	}
	return uint16(r)
}

// ASCIISlice returns the zero-copy substring of Contents for [start, end) and
// true when [start, end) consists solely of ASCII bytes (< 0x80).
func (f *SourceFile) ASCIISlice(start, end int) (string, bool) {
	if start < 0 || start > end || end > f.utf16Len {
		return "", false
	}
	if start == end {
		return "", true
	}
	if end <= f.asciiPrefix {
		return f.Contents[start:end], true
	}
	if start >= f.asciiSuffixUTF16Start {
		base := f.asciiSuffixByteStart + (start - f.asciiSuffixUTF16Start)
		return f.Contents[base : base+(end-start)], true
	}
	bStart, loStart := f.byteOffsetForUTF16(start)
	if loStart != 0 {
		return "", false
	}
	n := end - start
	if bStart+n > len(f.Contents) {
		return "", false
	}
	sub := f.Contents[bStart : bStart+n]
	if !jutil.IsASCII(sub) {
		return "", false
	}
	return sub, true
}

// Substring is contents.substring(start, end).
//
// Go-only: language-adapter — Java String API on the contents.
func (f *SourceFile) Substring(start, end int) string {
	if start < 0 || end > f.utf16Len || start > end {
		panic(jutil.NewStringIndexOutOfBoundsException(start, end))
	}
	if start == end {
		return ""
	}
	if end <= f.asciiPrefix {
		return f.Contents[start:end]
	}
	if start >= f.asciiSuffixUTF16Start {
		base := f.asciiSuffixByteStart + (start - f.asciiSuffixUTF16Start)
		return f.Contents[base : base+(end-start)]
	}
	bStart, loStart := f.byteOffsetForUTF16(start)
	if loStart == 0 {
		n := end - start
		if bStart+n <= len(f.Contents) && jutil.IsASCII(f.Contents[bStart:bStart+n]) {
			return f.Contents[bStart : bStart+n]
		}
	}
	bEnd, loEnd := f.byteOffsetForUTF16(end)
	if loStart == 0 && loEnd == 0 {
		sub := f.Contents[bStart:bEnd]
		if utf8.ValidString(sub) {
			return sub
		}
	}
	units := make([]uint16, 0, end-start)
	bPos := bStart
	uRem := end - start
	if loStart != 0 {
		units = append(units, loStart)
		bPos += 4
		uRem--
	}
	for uRem > 0 {
		if f.Contents[bPos] < 0x80 {
			units = append(units, uint16(f.Contents[bPos]))
			bPos++
			uRem--
			continue
		}
		r, size := decodeWTF8At(f.Contents, bPos)
		bPos += size
		if r >= 0x10000 {
			hi, lo := utf16.EncodeRune(r)
			units = append(units, uint16(hi))
			uRem--
			if uRem > 0 {
				units = append(units, uint16(lo))
				uRem--
			}
		} else {
			units = append(units, uint16(r))
			uRem--
		}
	}
	return jutil.StringFromUTF16(units)
}

// LineStarts builds (if needed) and returns the table of 0-based UTF-16 line
// start offsets for single-threaded callers (such as IRFactory) that maintain
// their own non-atomic cursor.
func (f *SourceFile) LineStarts() []int32 {
	f.linesOnce.Do(f.buildLineStarts)
	return f.lineStarts
}

// FindLine returns the 0-based line index in starts containing offset o
// (the largest k with starts[k] <= o, or 0 when o <= starts[0]).
func FindLine(starts []int32, o int32) int {
	if len(starts) == 0 || o <= starts[0] {
		return 0
	}
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
	if offset < 0 {
		return 0, offset
	}
	f.linesOnce.Do(f.buildLineStarts)
	starts := f.lineStarts
	o := int32(offset)
	i := int(f.lastLine.Load())
	if i < 0 || i >= len(starts) || o < starts[i] || (i+1 < len(starts) && starts[i+1] <= o) {
		if i >= 0 && i+1 < len(starts) && starts[i+1] <= o && (i+2 >= len(starts) || o < starts[i+2]) {
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
	n := f.utf16Len
	starts := make([]int32, 1, 1+n/32)
	s := f.Contents
	p := f.asciiPrefix
	sb := f.asciiSuffixByteStart
	sufU := f.asciiSuffixUTF16Start
	bPos := 0
	uPos := 0
	if p > 0 {
		ascii := s[:p]
		if strings.IndexByte(ascii, '\r') < 0 {
			for bPos < p {
				k := strings.IndexByte(ascii[bPos:], '\n')
				if k < 0 {
					break
				}
				bPos += k + 1
				starts = append(starts, int32(bPos))
			}
			bPos = p
			uPos = p
		} else {
			for bPos < p {
				switch s[bPos] {
				case '\n':
					bPos++
					uPos++
					starts = append(starts, int32(uPos))
				case '\r':
					bPos++
					uPos++
					if bPos < len(s) && s[bPos] == '\n' {
						bPos++
						uPos++
					}
					starts = append(starts, int32(uPos))
				default:
					bPos++
					uPos++
				}
			}
		}
	}
	for bPos < sb {
		for bPos+8 <= sb {
			if uint8(s[bPos]-0x20) >= 0x60 || uint8(s[bPos+1]-0x20) >= 0x60 ||
				uint8(s[bPos+2]-0x20) >= 0x60 || uint8(s[bPos+3]-0x20) >= 0x60 ||
				uint8(s[bPos+4]-0x20) >= 0x60 || uint8(s[bPos+5]-0x20) >= 0x60 ||
				uint8(s[bPos+6]-0x20) >= 0x60 || uint8(s[bPos+7]-0x20) >= 0x60 {
				break
			}
			bPos += 8
			uPos += 8
		}
		if bPos >= sb {
			break
		}
		c := s[bPos]
		if c < 0x80 {
			bPos++
			uPos++
			switch c {
			case '\n':
				starts = append(starts, int32(uPos))
			case '\r':
				if bPos < len(s) && s[bPos] == '\n' {
					bPos++
					uPos++
				}
				starts = append(starts, int32(uPos))
			}
			continue
		}
		r, size := decodeWTF8At(s, bPos)
		bPos += size
		if r >= 0x10000 {
			uPos += 2
		} else {
			uPos++
			if r == 0x2028 || r == 0x2029 {
				starts = append(starts, int32(uPos))
			}
		}
	}
	if bPos < len(s) {
		suffix := s[sb:]
		rel := bPos - sb
		if strings.IndexByte(suffix[rel:], '\r') < 0 {
			for rel < len(suffix) {
				k := strings.IndexByte(suffix[rel:], '\n')
				if k < 0 {
					break
				}
				rel += k + 1
				starts = append(starts, int32(sufU+rel))
			}
		} else {
			for rel < len(suffix) {
				switch suffix[rel] {
				case '\n':
					rel++
					starts = append(starts, int32(sufU+rel))
				case '\r':
					rel++
					if rel < len(suffix) && suffix[rel] == '\n' {
						rel++
					}
					starts = append(starts, int32(sufU+rel))
				default:
					rel++
				}
			}
		}
	}
	f.lineStarts = starts
}
