// Port of com/google/javascript/jscomp/parsing/parser/LineNumberScanner.java @ closure-compiler v20261015

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
	"math"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// LineNumberScanner is a utility for finding line and column offsets
// within a source file.
type LineNumberScanner struct {
	sourceFile    *SourceFile
	contents      []uint16
	sourceLength  int
	lastLine      int
	lastLineStart int
	nextLineStart int
}

// NewLineNumberScanner constructs a LineNumberScanner.
//
// Java: LineNumberScanner(SourceFile)
func NewLineNumberScanner(sourceFile *SourceFile) *LineNumberScanner {
	l := &LineNumberScanner{lastLine: -1, lastLineStart: -1, nextLineStart: 0}
	l.sourceFile = sourceFile
	l.contents = sourceFile.ContentsUTF16
	l.sourceLength = len(l.contents)
	return l
}

// GetSourcePosition returns the source position of character offset
// offset. This class expects this method to be called with increasing
// values of offset, more or less. RewindTo must be called before backing up
// to a previous line.
//
// Java: LineNumberScanner.getSourcePosition(int)
func (l *LineNumberScanner) GetSourcePosition(offset int) SourcePosition {
	if offset < l.lastLineStart {
		jutil.CheckArgumentMsg2(
			false,
			"Must call rewindTo before calling getSourcePosition for an earlier line (%s < %s)",
			offset,
			l.lastLineStart)
	}
	for offset >= l.nextLineStart {
		l.AdvanceLine()
	}
	return NewSourcePosition(offset, l.lastLine, offset-l.lastLineStart)
}

// GetSourceRange returns the range between two offsets.
//
// Java: LineNumberScanner.getSourceRange(int, int)
func (l *LineNumberScanner) GetSourceRange(startOffset, endOffset int) SourceRange {
	return NewSourceRange(l.GetSourcePosition(startOffset), l.GetSourcePosition(endOffset))
}

// RewindTo rewinds the scanner to an earlier position in the source file.
// This is necessary if backing up to a previous line.
//
// Java: LineNumberScanner.rewindTo(SourcePosition)
func (l *LineNumberScanner) RewindTo(position SourcePosition) {
	// Java: checkArgument(position.source == sourceFile) — Go-only:
	// language-adapter — positions carry no source (see SourcePosition).
	if offset := int(position.Offset); offset < l.lastLineStart {
		l.lastLine = int(position.Line) - 1
		l.nextLineStart = offset - int(position.Column)
		l.AdvanceLine()
	}
}

// AdvanceLine moves to the next line.
//
// Java: LineNumberScanner.advanceLine()
func (l *LineNumberScanner) AdvanceLine() {
	l.lastLine++

	l.lastLineStart = l.nextLineStart
	contents := l.contents
	sourceLength := l.sourceLength
	for index := l.lastLineStart; index < sourceLength; index++ {
		ch := contents[index]
		if ch == '\n' {
			l.nextLineStart = index + 1
			return
		}
		if ch == '\r' || ch == 0x2028 || ch == 0x2029 {
			if ch == '\r' && index+1 < sourceLength && contents[index+1] == '\n' {
				index++
			}
			l.nextLineStart = index + 1
			return
		}
	}
	l.nextLineStart = math.MaxInt32
}

// LineNumberScanner_IsLineTerminator reports whether ch terminates a line.
//
// Java: LineNumberScanner.isLineTerminator(char)
func LineNumberScanner_IsLineTerminator(ch uint16) bool {
	switch ch {
	case '\n', // Line Feed
		'\r',   // Carriage Return
		0x2028, // Line Separator
		0x2029: // Paragraph Separator
		return true
	default:
		return false
	}
}
