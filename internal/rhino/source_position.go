// Port of com/google/javascript/rhino/SourcePosition.java @ closure-compiler v20261015

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

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// SourcePosition represents a position in some piece of source code, with an
// associated item of type T found at that position.
//
// Go-only: language-adapter — an abstract class with subclasses (all in
// JSDocInfo): the interface below plus SourcePositionBase, which the
// subclasses embed.
type SourcePosition[T any] interface {
	SetItem(item T)
	SetPositionInformation(startLineno int, startCharno int, endLineno int, endCharno int)
	GetItem() T
	GetStartLine() int
	GetPositionOnStartLine() int
	GetEndLine() int
	GetPositionOnEndLine() int
	IsSamePositionAs(that SourcePosition[T]) bool
	GetSourcePositionBase() *SourcePositionBase[T]
}

// SourcePositionBase holds SourcePosition's fields and methods.
type SourcePositionBase[T any] struct {
	// The (well typed) item found at the source position.
	item T

	// The starting line number.
	startLineno int

	// The character position on the starting line.
	startCharno int

	// The ending line number.
	endLineno int

	// The character position on the ending line.
	endCharno int
}

// GetSourcePositionBase gives access to the base fields through the
// interface.
func (p *SourcePositionBase[T]) GetSourcePositionBase() *SourcePositionBase[T] { return p }

// SetItem sets the item that this source position references.
//
// Java: SourcePosition.setItem(T)
func (p *SourcePositionBase[T]) SetItem(item T) {
	p.item = item
}

// SetPositionInformation sets the position information contained in this
// source position.
//
// Java: SourcePosition.setPositionInformation(int, int, int, int)
func (p *SourcePositionBase[T]) SetPositionInformation(startLineno int, startCharno int, endLineno int, endCharno int) {
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

	p.startLineno = startLineno
	p.startCharno = startCharno
	p.endLineno = endLineno
	p.endCharno = endCharno
}

// GetItem returns the item found at this source position.
//
// Java: SourcePosition.getItem()
func (p *SourcePositionBase[T]) GetItem() T {
	return p.item
}

// GetStartLine returns the starting line number of this position.
//
// Java: SourcePosition.getStartLine()
func (p *SourcePositionBase[T]) GetStartLine() int {
	return p.startLineno
}

// GetPositionOnStartLine returns the character position on the starting
// line.
//
// Java: SourcePosition.getPositionOnStartLine()
func (p *SourcePositionBase[T]) GetPositionOnStartLine() int {
	return p.startCharno
}

// GetEndLine returns the ending line number of this position.
//
// Java: SourcePosition.getEndLine()
func (p *SourcePositionBase[T]) GetEndLine() int {
	return p.endLineno
}

// GetPositionOnEndLine returns the character position on the ending line.
//
// Java: SourcePosition.getPositionOnEndLine()
func (p *SourcePositionBase[T]) GetPositionOnEndLine() int {
	return p.endCharno
}

// IsSamePositionAs compares the four position fields.
//
// Go-only: language-adapter — Java takes a SourcePosition<?>; Go generics
// cannot express the wildcard, and every caller compares positions of the
// same item type.
//
// Java: SourcePosition.isSamePositionAs(SourcePosition<?>)
func (p *SourcePositionBase[T]) IsSamePositionAs(that SourcePosition[T]) bool {
	b := that.GetSourcePositionBase()
	return p.startLineno == b.startLineno &&
		p.startCharno == b.startCharno &&
		p.endLineno == b.endLineno &&
		p.endCharno == b.endCharno
}
