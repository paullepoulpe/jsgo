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
// Port of com/google/javascript/rhino/SimpleSourceFile.java (generated mechanically by j2go).

package rhino

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.rhino.SimpleSourceFile
type SimpleSourceFile struct {
	name      string `gson:"name"`
	_sfHandle atomic.Uint32
	idHash_   atomic.Int32
	kind      StaticSourceFile_SourceKind `gson:"kind"`
}

func (this_ *SimpleSourceFile) BindSelf_(s any) {
	_ = s
}
func (this_ *SimpleSourceFile) IsJavaNil_() bool { return this_ == nil }
func (*SimpleSourceFile) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.SimpleSourceFile", "SimpleSourceFile"
}
func (this_ *SimpleSourceFile) SourceFileHandlePtr_() *atomic.Uint32 { return &this_._sfHandle }
func (this_ *SimpleSourceFile) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*SimpleSourceFile) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.SimpleSourceFile", "com.google.javascript.rhino.SimpleSourceFile", func(v_ any) bool {
		return j2gort.InstanceOf[*SimpleSourceFile](v_)
	})
}
func (*SimpleSourceFile) Implements_rhino_StaticSourceFile() {}

func (this_ *SimpleSourceFile) initFields_SimpleSourceFile() {
	_ = this_
	this_.name = j2gort.NullString
}

// Java constructor: SimpleSourceFile(java.lang.String name, com.google.javascript.rhino.StaticSourceFile.SourceKind kind)
func NewSimpleSourceFile(name string, kind StaticSourceFile_SourceKind) *SimpleSourceFile {
	this_ := &SimpleSourceFile{}
	InitSimpleSourceFile(this_, name, kind)
	return this_
}

func InitSimpleSourceFile(this_ *SimpleSourceFile, name string, kind StaticSourceFile_SourceKind) {
	_ = name
	_ = kind
	this_.initFields_SimpleSourceFile()
	this_.name = name
	this_.kind = kind
}

// Java method: SimpleSourceFile.getName()
func (this_ *SimpleSourceFile) GetName() string {
	return this_.name
}

// Java method: SimpleSourceFile.getKind()
func (this_ *SimpleSourceFile) GetKind() StaticSourceFile_SourceKind {
	return this_.kind
}

// Java method: SimpleSourceFile.isClosureUnawareCode()
func (this_ *SimpleSourceFile) IsClosureUnawareCode() bool {
	return false
}

// Java method: SimpleSourceFile.getColumnOfOffset(int offset)
func (this_ *SimpleSourceFile) GetColumnOfOffset(offset int32) int32 {
	_ = offset
	return int32(0)
}

// Java method: SimpleSourceFile.getLineOfOffset(int offset)
func (this_ *SimpleSourceFile) GetLineOfOffset(offset int32) int32 {
	_ = offset
	return int32(1)
}

// Java method: SimpleSourceFile.getLineOffset(int line)
func (this_ *SimpleSourceFile) GetLineOffset(line int32) int32 {
	_ = line
	if line < int32(1) {
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Should not call getLineOffset with line number "+j2gort.StrValInt32(line))))
	}
	return int32(-2147483648)
}

// Java method: SimpleSourceFile.toString()
func (this_ *SimpleSourceFile) String() string {
	return this_.name
}

func (this_ *SimpleSourceFile) IsStrong() bool {
	return IfaceDefault_StaticSourceFile_IsStrong(this_)
}

func (this_ *SimpleSourceFile) IsWeak() bool {
	return IfaceDefault_StaticSourceFile_IsWeak(this_)
}

func (this_ *SimpleSourceFile) IsExtern() bool {
	return IfaceDefault_StaticSourceFile_IsExtern(this_)
}

func (this_ *SimpleSourceFile) IsNonCode() bool {
	return IfaceDefault_StaticSourceFile_IsNonCode(this_)
}

func (this_ *SimpleSourceFile) IsTypeScriptSource() bool {
	return IfaceDefault_StaticSourceFile_IsTypeScriptSource(this_)
}
