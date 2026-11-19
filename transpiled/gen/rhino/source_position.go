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
// Port of com/google/javascript/rhino/SourcePosition.java (generated mechanically by j2go).

package rhino

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java polymorphic class interface: com.google.javascript.rhino.SourcePosition
type SourcePosition interface {
	Implements_rhino_SourcePosition()
	IsJavaNil_() bool
	GetSourcePositionBase() *SourcePositionBase
	SetItem(any)
	SetPositionInformation(int32, int32, int32, int32)
	GetItem() any
	GetStartLine() int32
	GetPositionOnStartLine() int32
	GetEndLine() int32
	GetPositionOnEndLine() int32
	IsSamePositionAs(SourcePosition) bool
}
func ToIface_SourcePosition[T interface {
	SourcePosition
	comparable
}](v T) SourcePosition {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_SourcePosition[T SourcePosition](v T) SourcePosition {
	var vi SourcePosition = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.rhino.SourcePosition
type SourcePositionBase struct {
	item any `gson:"item"`
	startLineno int32 `gson:"startLineno"`
	startCharno int32 `gson:"startCharno"`
	endLineno int32 `gson:"endLineno"`
	endCharno int32 `gson:"endCharno"`
	idHash_ atomic.Int32
}

func (this_ *SourcePositionBase) BindSelf_(s any) {
	_ = s
}
func (this_ *SourcePositionBase) IsJavaNil_() bool { return this_ == nil }
func (*SourcePositionBase) JavaClassInfo_() (string, string) { return "com.google.javascript.rhino.SourcePosition", "SourcePosition" }
func (this_ *SourcePositionBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.SourcePosition", "com.google.javascript.rhino.SourcePosition", func(v_ any) bool {
		return j2gort.InstanceOf[SourcePosition](v_)
	})
}
func (b *SourcePositionBase) GetSourcePositionBase() *SourcePositionBase { return b }
func (*SourcePositionBase) Implements_rhino_SourcePosition() {}

func (this_ *SourcePositionBase) initFields_SourcePosition() {
	_ = this_
	this_.item = nil
	this_.startLineno = int32(0)
	this_.startCharno = int32(0)
	this_.endLineno = int32(0)
	this_.endCharno = int32(0)
}

func InitSourcePosition(this_ *SourcePositionBase) {
	this_.initFields_SourcePosition()
}

// Java method: SourcePosition.setItem(T item)
func (this_ *SourcePositionBase) SetItem(arg_1 any) {
	item := arg_1
	_ = item
	this_.item = item
}

// Java method: SourcePosition.setPositionInformation(int startLineno, int startCharno, int endLineno, int endCharno)
func (this_ *SourcePositionBase) SetPositionInformation(startLineno int32, startCharno int32, endLineno int32, endCharno int32) {
	_ = startLineno
	_ = startCharno
	_ = endLineno
	_ = endCharno
	if ((startLineno > endLineno)) {
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", j2gort.BoxString(("Recorded bad position information\nstart-line: " + j2gort.StrValInt32(startLineno) + "\nend-line: " + j2gort.StrValInt32(endLineno))))))
	} else if (((startLineno == endLineno) && (startCharno >= endCharno))) {
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", j2gort.BoxString(("Recorded bad position information\nline: " + j2gort.StrValInt32(startLineno) + "\nstart-char: " + j2gort.StrValInt32(startCharno) + "\nend-char: " + j2gort.StrValInt32(endCharno))))))
	}
	this_.startLineno = startLineno
	this_.startCharno = startCharno
	this_.endLineno = endLineno
	this_.endCharno = endCharno
}

// Java method: SourcePosition.getItem()
func (this_ *SourcePositionBase) GetItem() any {
	return this_.item
}

// Java method: SourcePosition.getStartLine()
func (this_ *SourcePositionBase) GetStartLine() int32 {
	return this_.startLineno
}

// Java method: SourcePosition.getPositionOnStartLine()
func (this_ *SourcePositionBase) GetPositionOnStartLine() int32 {
	return this_.startCharno
}

// Java method: SourcePosition.getEndLine()
func (this_ *SourcePositionBase) GetEndLine() int32 {
	return this_.endLineno
}

// Java method: SourcePosition.getPositionOnEndLine()
func (this_ *SourcePositionBase) GetPositionOnEndLine() int32 {
	return this_.endCharno
}

// Java method: SourcePosition.isSamePositionAs(com.google.javascript.rhino.SourcePosition<?> that)
func (this_ *SourcePositionBase) IsSamePositionAs(that SourcePosition) bool {
	_ = that
	return ((((this_.startLineno == that.GetSourcePositionBase().startLineno) && (this_.startCharno == that.GetSourcePositionBase().startCharno)) && (this_.endLineno == that.GetSourcePositionBase().endLineno)) && (this_.endCharno == that.GetSourcePositionBase().endCharno))
}

