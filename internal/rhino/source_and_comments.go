// Port of com/google/javascript/rhino/InputId.java @ closure-compiler v20261021
// Port of com/google/javascript/rhino/ErrorReporter.java @ closure-compiler v20261021
// Port of com/google/javascript/rhino/NonJSDocComment.java @ closure-compiler v20261021
// Port of com/google/javascript/rhino/SimpleSourceFile.java @ closure-compiler v20261021

/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * The Original Code is Rhino code, released
 * May 6, 1999.
 *
 * The Initial Developer of the Original Code is
 * Netscape Communications Corporation.
 * Portions created by the Initial Developer are Copyright (C) 1997-1999
 * the Initial Developer. All Rights Reserved.
 *
 * Contributor(s):
 *   John Lenz
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * The Original Code is Rhino code, released
 * May 6, 1999.
 *
 * The Initial Developer of the Original Code is
 * Netscape Communications Corporation.
 * Portions created by the Initial Developer are Copyright (C) 1997-1999
 * the Initial Developer. All Rights Reserved.
 *
 * Contributor(s):
 *   Norris Boyd
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * The Original Code is Rhino code, released
 * May 6, 1999.
 *
 * The Initial Developer of the Original Code is
 * Netscape Communications Corporation.
 * Portions created by the Initial Developer are Copyright (C) 1997-1999
 * the Initial Developer. All Rights Reserved.
 *
 * Contributor(s):
 *   Bob Jervis
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * The Original Code is Rhino code, released
 * May 6, 1999.
 *
 * The Initial Developer of the Original Code is
 * Netscape Communications Corporation.
 * Portions created by the Initial Developer are Copyright (C) 1997-1999
 * the Initial Developer. All Rights Reserved.
 *
 * Contributor(s):
 *   Nick Santos
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

package rhino

import (
	"math"
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
)

// InputId is an id used uniquely identify a CompilerInput.
type InputId struct {
	id   string
	hash int32
}

// Init initializes i in-place with id.
func (i *InputId) Init(id string) {
	i.id = id
	i.hash = jutil.StringHashCode(id)
}

// Java: InputId(String)
func NewInputId(id string) *InputId {
	return &InputId{id: id, hash: jutil.StringHashCode(id)}
}

// Java: InputId.getIdName()
func (i *InputId) GetIdName() string {
	return i.id
}

// Java: InputId.hashCode()
func (i *InputId) HashCode() int32 {
	return i.hash
}

// Java: InputId.equals(Object)
func (i *InputId) Equals(obj any) bool {
	other, ok := obj.(*InputId)
	if !ok || other == nil {
		return false
	}
	return i == other || (i.hash == other.hash && i.id == other.id)
}

// Java: InputId.toString()
func (i *InputId) String() string {
	return "InputId: " + i.GetIdName()
}

// ErrorReporter defines a protocol for the reporting of errors during
// JavaScript translation or execution.
type ErrorReporter interface {
	// Warning reports a warning.
	//
	// The implementing class may choose to ignore the warning if it desires.
	//
	// message: a String describing the warning. sourceName: a String
	// describing the JavaScript source where the warning occurred; typically a
	// filename or URL. line: the line number associated with the warning.
	// lineOffset: the offset into lineSource where problem was detected.
	Warning(message string, sourceName string, line int, lineOffset int)

	// Error reports an error.
	//
	// The implementing class is free to throw an exception if it desires.
	//
	// If execution has not yet begun, the JavaScript engine is free to find
	// additional errors rather than terminating the translation. It will not
	// execute a script that had errors, however.
	//
	// message: a String describing the error. sourceName: a String
	// describing the JavaScript source where the error occurred; typically a
	// filename or URL. line: the line number associated with the error.
	// lineOffset: the offset into lineSource where problem was detected.
	Error(message string, sourceName string, line int, lineOffset int)
}

type errorReporter_nullInstance struct{}

func (errorReporter_nullInstance) Warning(message string, sourceName string, line int, lineOffset int) {
}

func (errorReporter_nullInstance) Error(message string, sourceName string, line int, lineOffset int) {
}

// Java: ErrorReporter.NULL_INSTANCE
var ErrorReporter_NULL_INSTANCE ErrorReporter = errorReporter_nullInstance{}

type errorReporter_alwaysThrowsInstance struct{}

func (errorReporter_alwaysThrowsInstance) Warning(message string, sourceName string, line int, lineOffset int) {
	panic(jutil.NewIllegalStateExceptionMsg(message))
}

func (errorReporter_alwaysThrowsInstance) Error(message string, sourceName string, line int, lineOffset int) {
	panic(jutil.NewIllegalStateExceptionMsg(message))
}

// Java: ErrorReporter.ALWAYS_THROWS_INSTANCE
var ErrorReporter_ALWAYS_THROWS_INSTANCE ErrorReporter = errorReporter_alwaysThrowsInstance{}

// NonJSDocComment is a minimal class holding information about a nonJSDoc
// comment's source location and contents.
type NonJSDocComment struct {
	startPosition     parser.LinePosition
	endPosition       parser.LinePosition
	contents          string
	endsAsLineComment bool
	isInline          bool
}

// Java: NonJSDocComment(SourcePosition, SourcePosition, String)
func NewNonJSDocComment(start parser.LinePosition, end parser.LinePosition, s *string) *NonJSDocComment {
	var contents string
	if s != nil {
		contents = *s
	}
	return &NonJSDocComment{
		startPosition: start,
		endPosition:   end,
		contents:      contents,
	}
}

// Java: NonJSDocComment.getCommentString()
func (c *NonJSDocComment) GetCommentString() string {
	return c.contents
}

// Java: NonJSDocComment.getStartPosition()
func (c *NonJSDocComment) GetStartPosition() parser.LinePosition {
	return c.startPosition
}

// Java: NonJSDocComment.getEndPosition()
func (c *NonJSDocComment) GetEndPosition() parser.LinePosition {
	return c.endPosition
}

// Java: NonJSDocComment.setEndsAsLineComment(boolean)
func (c *NonJSDocComment) SetEndsAsLineComment(endsAsLineComment bool) {
	c.endsAsLineComment = endsAsLineComment
}

// Java: NonJSDocComment.setIsInline(boolean)
func (c *NonJSDocComment) SetIsInline(isInline bool) {
	c.isInline = isInline
}

// Java: NonJSDocComment.isEndingAsLineComment()
func (c *NonJSDocComment) IsEndingAsLineComment() bool {
	return c.endsAsLineComment
}

// Java: NonJSDocComment.isInline()
func (c *NonJSDocComment) IsInline() bool {
	return c.isInline
}

// Java: NonJSDocComment.toString()
func (c *NonJSDocComment) String() string {
	return "NonJSDocComment : " + c.GetCommentString()
}

// SimpleSourceFile is a simple implementation of StaticSourceFile for
// testing.
type SimpleSourceFile struct {
	name string
	kind StaticSourceFile_SourceKind
}

// Java: SimpleSourceFile(String, SourceKind)
func NewSimpleSourceFile(name string, kind StaticSourceFile_SourceKind) *SimpleSourceFile {
	return &SimpleSourceFile{name: name, kind: kind}
}

// Java: SimpleSourceFile.getName()
func (f *SimpleSourceFile) GetName() string {
	return f.name
}

func (f *SimpleSourceFile) GetNamePtr() *string { return &f.name }

// Java: SimpleSourceFile.getKind()
func (f *SimpleSourceFile) GetKind() StaticSourceFile_SourceKind {
	return f.kind
}

// Java: SimpleSourceFile.isClosureUnawareCode()
func (f *SimpleSourceFile) IsClosureUnawareCode() bool {
	return false
}

// Java: SimpleSourceFile.getColumnOfOffset(int)
func (f *SimpleSourceFile) GetColumnOfOffset(offset int) int {
	return 0
}

// Java: SimpleSourceFile.getLineOfOffset(int)
func (f *SimpleSourceFile) GetLineOfOffset(offset int) int {
	return 1
}

// Java: SimpleSourceFile.getLineOffset(int)
func (f *SimpleSourceFile) GetLineOffset(line int) int {
	if line < 1 {
		panic(jutil.NewIllegalStateExceptionMsg(
			"Should not call getLineOffset with line number " + strconv.Itoa(line)))
	}
	return math.MinInt32
}

// Java: SimpleSourceFile.toString()
func (f *SimpleSourceFile) String() string {
	return f.name
}

// Inherited StaticSourceFile default methods.

// Java: StaticSourceFile.isStrong()
func (f *SimpleSourceFile) IsStrong() bool { return StaticSourceFile_IsStrong(f) }

// Java: StaticSourceFile.isWeak()
func (f *SimpleSourceFile) IsWeak() bool { return StaticSourceFile_IsWeak(f) }

// Java: StaticSourceFile.isExtern()
func (f *SimpleSourceFile) IsExtern() bool { return StaticSourceFile_IsExtern(f) }

// Java: StaticSourceFile.isNonCode()
func (f *SimpleSourceFile) IsNonCode() bool { return StaticSourceFile_IsNonCode(f) }

// Java: StaticSourceFile.isTypeScriptSource()
func (f *SimpleSourceFile) IsTypeScriptSource() bool { return StaticSourceFile_IsTypeScriptSource(f) }
