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
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */
// Port of com/google/javascript/rhino/jstype/SimpleReference.java (generated mechanically by j2go).

package rhino

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java polymorphic class interface: com.google.javascript.rhino.jstype.SimpleReference
type SimpleReference interface {
	Implements_rhino_SimpleReference()
	IsJavaNil_() bool
	GetSimpleReferenceBase() *SimpleReferenceBase
	StaticTypedRef
	String() string
}
func ToIface_SimpleReference[T interface {
	SimpleReference
	comparable
}](v T) SimpleReference {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_SimpleReference[T SimpleReference](v T) SimpleReference {
	var vi SimpleReference = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.rhino.jstype.SimpleReference
type SimpleReferenceBase struct {
	symbol StaticTypedSlot `gson:"symbol"`
	node Node `gson:"node"`
	idHash_ atomic.Int32
}

func (this_ *SimpleReferenceBase) BindSelf_(s any) {
	_ = s
}
func (this_ *SimpleReferenceBase) IsJavaNil_() bool { return this_ == nil }
func (*SimpleReferenceBase) JavaClassInfo_() (string, string) { return "com.google.javascript.rhino.jstype.SimpleReference", "SimpleReference" }
func (this_ *SimpleReferenceBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.SimpleReference", "com.google.javascript.rhino.jstype.SimpleReference", func(v_ any) bool {
		return j2gort.InstanceOf[SimpleReference](v_)
	})
}
func (b *SimpleReferenceBase) GetSimpleReferenceBase() *SimpleReferenceBase { return b }
func (*SimpleReferenceBase) Implements_rhino_SimpleReference() {}
func (*SimpleReferenceBase) Implements_rhino_StaticTypedRef() {}
func (*SimpleReferenceBase) Implements_rhino_StaticRef() {}

func (this_ *SimpleReferenceBase) initFields_SimpleReference() {
	_ = this_
}

// Java constructor: SimpleReference(T symbol, com.google.javascript.rhino.Node node)
func NewSimpleReference(symbol StaticTypedSlot, node Node) SimpleReference {
	this_ := &SimpleReferenceBase{}
	InitSimpleReference(this_, symbol, node)
	return this_
}

func InitSimpleReference(this_ *SimpleReferenceBase, symbol StaticTypedSlot, node Node) {
	_ = symbol
	_ = node
	this_.initFields_SimpleReference()
	this_.symbol = symbol
	this_.node = node
}

// Java method: SimpleReference.getSymbol()
func (this_ *SimpleReferenceBase) GetSymbol() any {
	return this_.symbol
}

// Java method: SimpleReference.getNode()
func (this_ *SimpleReferenceBase) GetNode() Node {
	return this_.node
}

// Java method: SimpleReference.getSourceFile()
func (this_ *SimpleReferenceBase) GetSourceFile() StaticSourceFile {
	return this_.node.GetStaticSourceFile()
}

// Java method: SimpleReference.toString()
func (this_ *SimpleReferenceBase) String() string {
	var sourceName string
	_ = sourceName
	if (this_.node == 0) {
		sourceName = j2gort.NullString
	} else {
		sourceName = this_.node.GetSourceFileName()
	}
	var lineNo int32
	_ = lineNo
	if (this_.node == 0) {
		lineNo = int32(-1)
	} else {
		lineNo = this_.node.GetLineno()
	}
	return (j2gort.StrValString(this_.node.GetQualifiedName()) + "@" + j2gort.StrValString(sourceName) + ":" + j2gort.StrValInt32(lineNo))
}

