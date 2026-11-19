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
// Port of com/google/javascript/jscomp/parsing/parser/LineNumberScanner.java (generated mechanically by j2go).

package parser

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.LineNumberScanner
type LineNumberScanner struct {
	sourceFile     *SourceFile    `gson:"sourceFile"`
	contents       string         `gson:"contents"`
	cachedPosition SourcePosition `gson:"cachedPosition"`
	sourceLength   int32          `gson:"sourceLength"`
	lastLine       int32          `gson:"lastLine"`
	lastLineStart  int32          `gson:"lastLineStart"`
	nextLineStart  int32          `gson:"nextLineStart"`
	cachedOffset   int32          `gson:"cachedOffset"`
	idHash_        atomic.Int32
}

func (this_ *LineNumberScanner) BindSelf_(s any) {
	_ = s
}
func (this_ *LineNumberScanner) IsJavaNil_() bool { return this_ == nil }
func (*LineNumberScanner) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.LineNumberScanner", "LineNumberScanner"
}
func (this_ *LineNumberScanner) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*LineNumberScanner) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.LineNumberScanner", "com.google.javascript.jscomp.parsing.parser.LineNumberScanner", func(v_ any) bool {
		return j2gort.InstanceOf[*LineNumberScanner](v_)
	})
}

func (this_ *LineNumberScanner) initFields_LineNumberScanner() {
	_ = this_
	this_.contents = j2gort.NullString
	this_.lastLine = int32(-1)
	this_.lastLineStart = int32(-1)

	this_.cachedOffset = int32(-1)
}

// Java constructor: LineNumberScanner(com.google.javascript.jscomp.parsing.parser.SourceFile sourceFile)
func NewLineNumberScanner(sourceFile *SourceFile) *LineNumberScanner {
	this_ := &LineNumberScanner{}
	InitLineNumberScanner(this_, sourceFile)
	return this_
}

func InitLineNumberScanner(this_ *LineNumberScanner, sourceFile *SourceFile) {
	_ = sourceFile
	this_.initFields_LineNumberScanner()
	this_.sourceFile = sourceFile
	this_.contents = sourceFile.contents
	this_.sourceLength = j2gort.JStrLen(this_.contents)
}

// Java method: LineNumberScanner.getSourcePosition(int offset)
func (this_ *LineNumberScanner) GetSourcePosition(offset int32) SourcePosition {
	_ = offset
	if (offset == this_.cachedOffset) && (this_.cachedPosition != (SourcePosition{})) {
		return this_.cachedPosition
	}
	if !(offset >= this_.lastLineStart) {
		j2gort.CheckArgumentFmt(false, "Must call rewindTo before calling getSourcePosition for an earlier line (%s < %s)", j2gort.ValueOfInt(offset), j2gort.ValueOfInt(this_.lastLineStart))
	}
	for offset >= this_.nextLineStart {
		this_.priv_LineNumberScanner_advanceLine_7181()
	}
	var pos SourcePosition = NewSourcePosition(this_.sourceFile, offset, this_.lastLine, (offset - this_.lastLineStart))
	_ = pos
	this_.cachedOffset = offset
	this_.cachedPosition = pos
	return pos
}

// Java method: LineNumberScanner.getSourceRange(int startOffset, int endOffset)
func (this_ *LineNumberScanner) GetSourceRange(startOffset int32, endOffset int32) SourceRange {
	_ = startOffset
	_ = endOffset
	return NewSourceRange(this_.GetSourcePosition(startOffset), this_.GetSourcePosition(endOffset))
}

// Java method: LineNumberScanner.rewindTo(com.google.javascript.jscomp.parsing.parser.util.SourcePosition position)
func (this_ *LineNumberScanner) RewindTo(position SourcePosition) {
	_ = position
	if !(position.source == this_.sourceFile) {
		j2gort.CheckArgument(false)
	}
	if position.offset < this_.lastLineStart {
		this_.lastLine = (position.line - int32(1))
		this_.nextLineStart = (position.offset - position.column)
		this_.cachedOffset = int32(-1)
		this_.cachedPosition = (SourcePosition{})
		this_.priv_LineNumberScanner_advanceLine_7181()
	}
}

// Java method: LineNumberScanner.advanceLine()
func (this_ *LineNumberScanner) priv_LineNumberScanner_advanceLine_7181() {
	this_.lastLine = this_.lastLine + 1
	this_.lastLineStart = this_.nextLineStart
	{
		var index int32 = this_.lastLineStart
		_ = index
		for ; index < this_.sourceLength; index = index + 1 {
			var ch uint16 = j2gort.JStrCharAt(this_.contents, index)
			_ = ch
			if LineNumberScanner_IsLineTerminator(ch) {
				if ((int32(ch) == int32(uint16(13))) && ((index + int32(1)) < this_.sourceLength)) && (int32(j2gort.JStrCharAt(this_.contents, (index+int32(1)))) == int32(uint16(10))) {
					index = index + 1
				}
				this_.nextLineStart = (index + int32(1))
				return
			}
		}
	}
	this_.nextLineStart = int32(2147483647)
}

// Java static method: LineNumberScanner.isLineTerminator(char ch)
func LineNumberScanner_IsLineTerminator(ch uint16) bool {
	_ = ch
	switch ch {
	case uint16(10), uint16(13), uint16(8232), uint16(8233):
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}
