// Port of com/google/javascript/jscomp/parsing/parser/LineNumberScanner.java @ closure-compiler v20260415

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
func (l *LineNumberScanner) GetSourcePosition(offset int) *SourcePosition {
	jutil.CheckArgumentMsg(
		offset >= l.lastLineStart,
		"Must call rewindTo before calling getSourcePosition for an earlier line (%s < %s)",
		offset,
		l.lastLineStart)
	for offset >= l.nextLineStart {
		l.AdvanceLine()
	}
	return NewSourcePosition(l.sourceFile, offset, l.lastLine, offset-l.lastLineStart)
}

// GetSourceRange returns the range between two offsets.
//
// Java: LineNumberScanner.getSourceRange(int, int)
func (l *LineNumberScanner) GetSourceRange(startOffset, endOffset int) *SourceRange {
	return NewSourceRange(l.GetSourcePosition(startOffset), l.GetSourcePosition(endOffset))
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
