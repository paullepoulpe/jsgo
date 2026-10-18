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

	// Go-only: performance — the positions and ranges of the scanner's
	// tokens (see getTokenSourcePosition).
	// Last token position only: token/tree locations are immutable and nearby
	// queries often repeat the same offset. Comments keep independent positions.
	lastTokenPosition *SourcePosition
	tokenPositions    jutil.Slab[SourcePosition]
	tokenRanges       jutil.Slab[SourceRange]
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
func (l *LineNumberScanner) GetSourcePosition(offset int) *SourcePosition {
	return l.getSourcePosition(offset, false)
}

// getSourcePosition is getSourcePosition(int); the position comes from the
// tokenPositions slab if forToken is set.
//
// Go-only: performance — see getTokenSourcePosition.
func (l *LineNumberScanner) getSourcePosition(offset int, forToken bool) *SourcePosition {
	jutil.CheckArgumentMsg2(
		offset >= l.lastLineStart,
		"Must call rewindTo before calling getSourcePosition for an earlier line (%s < %s)",
		offset,
		l.lastLineStart)
	for offset >= l.nextLineStart {
		l.AdvanceLine()
	}
	if forToken {
		if p := l.lastTokenPosition; p != nil && p.Offset == offset {
			return p
		}
		p := l.tokenPositions.New()
		p.Source, p.Offset, p.Line, p.Column = l.sourceFile, offset, l.lastLine, offset-l.lastLineStart
		l.lastTokenPosition = p
		return p
	}
	return NewSourcePosition(l.sourceFile, offset, l.lastLine, offset-l.lastLineStart)
}

// GetSourceRange returns the range between two offsets.
//
// Java: LineNumberScanner.getSourceRange(int, int)
func (l *LineNumberScanner) GetSourceRange(startOffset, endOffset int) *SourceRange {
	return NewSourceRange(l.GetSourcePosition(startOffset), l.GetSourcePosition(endOffset))
}

// getTokenSourcePosition is GetSourcePosition for the Scanner's tokens and
// for the positions it hands the Parser.
//
// Go-only: performance — a parse creates two positions and a range per
// token, which were ~half of the front end's allocations. These come from
// slabs instead (jutil.Slab): they become garbage with the tokens and parse
// trees, after IRFactory has read them. Comments keep using
// GetSourcePosition and GetSourceRange: IDE mode and NonJSDocComment keep
// their positions for the whole compilation, which would pin whole chunks.
func (l *LineNumberScanner) getTokenSourcePosition(offset int) *SourcePosition {
	return l.getSourcePosition(offset, true)
}

// getTokenSourceRange is GetSourceRange for the Scanner's tokens.
//
// Go-only: performance — see getTokenSourcePosition.
func (l *LineNumberScanner) getTokenSourceRange(startOffset, endOffset int) *SourceRange {
	start := l.getTokenSourcePosition(startOffset)
	end := l.getTokenSourcePosition(endOffset)
	r := l.tokenRanges.New()
	r.Start, r.End = start, end
	return r
}

// RewindTo rewinds the scanner to an earlier position in the source file.
// This is necessary if backing up to a previous line.
//
// Java: LineNumberScanner.rewindTo(SourcePosition)
func (l *LineNumberScanner) RewindTo(position *SourcePosition) {
	jutil.CheckArgument(position.Source == l.sourceFile)
	if position.Offset < l.lastLineStart {
		l.lastLine = position.Line - 1
		l.nextLineStart = position.Offset - position.Column
		l.AdvanceLine()
	}
}

// AdvanceLine moves to the next line.
//
// Java: LineNumberScanner.advanceLine()
func (l *LineNumberScanner) AdvanceLine() {
	l.lastLine++

	l.lastLineStart = l.nextLineStart
	for index := l.lastLineStart; index < l.sourceLength; index++ {
		ch := l.contents[index]
		if LineNumberScanner_IsLineTerminator(ch) {
			if ch == '\r' && index+1 < l.sourceLength && l.contents[index+1] == '\n' {
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
