// Port of com/google/javascript/rhino/StaticSourceFile.java @ closure-compiler v20261021
// Port of com/google/javascript/rhino/SourcePosition.java @ closure-compiler v20261021

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

package rhino

import (
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// StaticSourceFile contains information about a compiler input.
//
// Java default methods are declared in the interface (implementations may
// override them, as in Java) and their default bodies are the
// StaticSourceFile_X functions, which implementations call to inherit them.
type StaticSourceFile interface {
	// GetName is the name of the file. Must be unique across all files in the
	// compilation.
	GetName() string

	// GetKind is the source kind.
	GetKind() StaticSourceFile_SourceKind

	// IsStrong reports whether the source kind is STRONG.
	IsStrong() bool

	// IsWeak reports whether the source kind is WEAK.
	IsWeak() bool

	// IsExtern reports whether the source kind is EXTERN.
	IsExtern() bool

	// IsNonCode reports whether the source kind is NON_CODE.
	IsNonCode() bool

	IsTypeScriptSource() bool

	// IsClosureUnawareCode reports whether this source file contains code
	// that should be optimized using a minimal set of known-safe
	// optimizations. In particular, code in this StaticSourceFile is unlikely
	// to be well-typed (or have types annotated at all), and is likely written
	// by someone unaware of the limitations of this compiler.
	IsClosureUnawareCode() bool

	// GetLineOffset returns the offset of the given line number relative to
	// the file start. Line number should be 1-based.
	//
	// If the source file doesn't have line information, it should return
	// Integer.MIN_VALUE. The negative offsets will make it more obvious what
	// happened.
	//
	// Throws IllegalArgumentException if lineno is less than 1 or greater
	// than the number of lines in the source.
	GetLineOffset(lineNumber int) int

	// GetLineOfOffset gets the 1-based line number of the given source offset.
	// The behavior is undefined if this offset does not exist in the source
	// file.
	GetLineOfOffset(offset int) int

	// GetColumnOfOffset gets the 0-based column number of the given source
	// offset. The behavior is undefined if this offset does not exist in the
	// source file.
	GetColumnOfOffset(offset int) int
}

// StaticSourceFileFast is an optional Go-only fast-path interface implemented
// by SourceFile and SimpleSourceFile so Node.GetSourceFileName avoids *string
// heap boxing.
type StaticSourceFileFast interface {
	GetNamePtr() *string
}

// StaticSourceFile_SourceKind enumerates source kinds.
//
// Java: StaticSourceFile.SourceKind
type StaticSourceFile_SourceKind int32

const (
	// A file whose contents are necessary both for type checking and emitting
	// code.
	StaticSourceFile_SourceKind_STRONG StaticSourceFile_SourceKind = iota
	// A file whose contents are necessary for type checking only.
	StaticSourceFile_SourceKind_WEAK
	// A file whose contents are extern declarations.
	StaticSourceFile_SourceKind_EXTERN
	// A file that is not part of the compilation unit and is not strong, weak,
	// or an extern
	StaticSourceFile_SourceKind_NON_CODE
)

var staticSourceFile_SourceKind_names = [...]string{"STRONG", "WEAK", "EXTERN", "NON_CODE"}

// String is Enum.name().
func (k StaticSourceFile_SourceKind) String() string { return staticSourceFile_SourceKind_names[k] }

// Java: StaticSourceFile.isStrong() (default method)
func StaticSourceFile_IsStrong(f StaticSourceFile) bool {
	return f.GetKind() == StaticSourceFile_SourceKind_STRONG
}

// Java: StaticSourceFile.isWeak() (default method)
func StaticSourceFile_IsWeak(f StaticSourceFile) bool {
	return f.GetKind() == StaticSourceFile_SourceKind_WEAK
}

// Java: StaticSourceFile.isExtern() (default method)
func StaticSourceFile_IsExtern(f StaticSourceFile) bool {
	return f.GetKind() == StaticSourceFile_SourceKind_EXTERN
}

// Java: StaticSourceFile.isNonCode() (default method)
func StaticSourceFile_IsNonCode(f StaticSourceFile) bool {
	return f.GetKind() == StaticSourceFile_SourceKind_NON_CODE
}

// Java: StaticSourceFile.isTypeScriptSource() (default method)
func StaticSourceFile_IsTypeScriptSource(f StaticSourceFile) bool {
	sourceName := f.GetName()
	return strings.HasSuffix(sourceName, ".closure.js") || strings.HasSuffix(sourceName, ".tsx.cl.js")
}

// SourcePosition represents a position in some piece of source code, with an
// associated item of type T found at that position.
//
// Go-only: performance — collapsed from an 8-type interface/subclass hierarchy
// into a single concrete generic struct so JSDocInfo_Marker embeds positions
// inline by value with zero per-position heap allocations.
type SourcePosition[T comparable] struct {
	item        T
	startLineno int32
	startCharno int32
	endLineno   int32
	endCharno   int32
	brackets    bool
	valid       bool
}

// IsValid reports whether this SourcePosition has been populated.
func (p *SourcePosition[T]) IsValid() bool {
	return p != nil && p.valid
}

// SetItem sets the item that this source position references.
//
// Java: SourcePosition.setItem(T)
func (p *SourcePosition[T]) SetItem(item T) {
	p.item = item
	p.valid = true
}

func checkTrimmedStringItem(item string) {
	if len(item) == 0 {
		jutil.StringCharAt(item, 0)
	}
	jutil.CheckArgumentMsg(
		item[0] != ' ' && item[len(item)-1] != ' ',
		"String has leading or trailing whitespace")
}

// SetTrimmedItem sets a string item after validating that it has no leading or
// trailing space.
//
// Java: JSDocInfo.TrimmedStringPosition.setItem(String)
func SetTrimmedItem(p *SourcePosition[string], item string) {
	checkTrimmedStringItem(item)
	p.item = item
	p.valid = true
}

// HasBrackets returns whether a type position has curly braces around it.
//
// Java: JSDocInfo.TypePosition.hasBrackets()
func (p *SourcePosition[T]) HasBrackets() bool {
	return p.brackets
}

// SetHasBrackets records whether a type position has curly braces around it.
//
// Java: JSDocInfo.TypePosition.setHasBrackets(boolean)
func (p *SourcePosition[T]) SetHasBrackets(newVal bool) {
	p.brackets = newVal
}

// SetPositionInformation sets the position information contained in this
// source position.
//
// Java: SourcePosition.setPositionInformation(int, int, int, int)
func (p *SourcePosition[T]) SetPositionInformation(startLineno int, startCharno int, endLineno int, endCharno int) {
	if startLineno > endLineno {
		panic(jutil.NewIllegalStateExceptionMsg(
			"Recorded bad position information\n" +
				"start-line: " + strconv.Itoa(startLineno) + "\n" +
				"end-line: " + strconv.Itoa(endLineno)))
	} else if startLineno == endLineno && startCharno >= endCharno {
		panic(jutil.NewIllegalStateExceptionMsg(
			"Recorded bad position information\n" +
				"line: " + strconv.Itoa(startLineno) + "\n" +
				"start-char: " + strconv.Itoa(startCharno) + "\n" +
				"end-char: " + strconv.Itoa(endCharno)))
	}

	p.startLineno = int32(startLineno)
	p.startCharno = int32(startCharno)
	p.endLineno = int32(endLineno)
	p.endCharno = int32(endCharno)
	p.valid = true
}

// GetItem returns the item found at this source position.
//
// Java: SourcePosition.getItem()
func (p *SourcePosition[T]) GetItem() T {
	return p.item
}

// GetStartLine returns the starting line number of this position.
//
// Java: SourcePosition.getStartLine()
func (p *SourcePosition[T]) GetStartLine() int {
	return int(p.startLineno)
}

// GetPositionOnStartLine returns the character position on the starting line.
//
// Java: SourcePosition.getPositionOnStartLine()
func (p *SourcePosition[T]) GetPositionOnStartLine() int {
	return int(p.startCharno)
}

// GetEndLine returns the ending line number of this position.
//
// Java: SourcePosition.getEndLine()
func (p *SourcePosition[T]) GetEndLine() int {
	return int(p.endLineno)
}

// GetPositionOnEndLine returns the character position on the ending line.
//
// Java: SourcePosition.getPositionOnEndLine()
func (p *SourcePosition[T]) GetPositionOnEndLine() int {
	return int(p.endCharno)
}

// IsSamePositionAs compares the four position fields.
//
// Java: SourcePosition.isSamePositionAs(SourcePosition<?>)
func (p *SourcePosition[T]) IsSamePositionAs(that *SourcePosition[T]) bool {
	return p.startLineno == that.startLineno &&
		p.startCharno == that.startCharno &&
		p.endLineno == that.endLineno &&
		p.endCharno == that.endCharno
}
