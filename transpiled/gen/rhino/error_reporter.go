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

// API class
// Port of com/google/javascript/rhino/ErrorReporter.java (generated mechanically by j2go).

package rhino

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java interface: com.google.javascript.rhino.ErrorReporter
type ErrorReporter interface {
	Implements_rhino_ErrorReporter()
	IsJavaNil_() bool
	Warning(string, string, int32, int32)
	M_Error(string, string, int32, int32)
}

func ToIface_ErrorReporter[T interface {
	ErrorReporter
	comparable
}](v T) ErrorReporter {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_ErrorReporter[T ErrorReporter](v T) ErrorReporter {
	var vi ErrorReporter = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.ErrorReporter", "com.google.javascript.rhino.ErrorReporter", func(v_ any) bool {
		return j2gort.InstanceOf[ErrorReporter](v_)
	})
}

var static_ErrorReporter_NULL_INSTANCE ErrorReporter

func StaticGet_ErrorReporter_NULL_INSTANCE() ErrorReporter {
	clinit_ErrorReporter()
	return static_ErrorReporter_NULL_INSTANCE
}
func StaticSet_ErrorReporter_NULL_INSTANCE(v ErrorReporter) {
	clinit_ErrorReporter()
	static_ErrorReporter_NULL_INSTANCE = v
}

var static_ErrorReporter_ALWAYS_THROWS_INSTANCE ErrorReporter

func StaticGet_ErrorReporter_ALWAYS_THROWS_INSTANCE() ErrorReporter {
	clinit_ErrorReporter()
	return static_ErrorReporter_ALWAYS_THROWS_INSTANCE
}
func StaticSet_ErrorReporter_ALWAYS_THROWS_INSTANCE(v ErrorReporter) {
	clinit_ErrorReporter()
	static_ErrorReporter_ALWAYS_THROWS_INSTANCE = v
}

var clinitDone_ErrorReporter uint32
var clinitState_ErrorReporter j2gort.ClinitState

func clinit_ErrorReporter() {
	if atomic.LoadUint32(&clinitDone_ErrorReporter) != 0 {
		return
	}
	clinitSlow_ErrorReporter()
}

func clinitSlow_ErrorReporter() {
	if !j2gort.EnterClinit(&clinitState_ErrorReporter, &clinitDone_ErrorReporter) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ErrorReporter, &clinitDone_ErrorReporter)
	static_ErrorReporter_NULL_INSTANCE = NewErrorReporter__Anon_153()
	static_ErrorReporter_ALWAYS_THROWS_INSTANCE = NewErrorReporter__Anon_154()
}

// Java class:
type ErrorReporter__Anon_153 struct {
	idHash_ atomic.Int32
}

func (this_ *ErrorReporter__Anon_153) BindSelf_(s any) {
	_ = s
}
func (this_ *ErrorReporter__Anon_153) IsJavaNil_() bool { return this_ == nil }
func (*ErrorReporter__Anon_153) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.ErrorReporter$1", ""
}
func (this_ *ErrorReporter__Anon_153) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ErrorReporter__Anon_153) HasIdentityEquals_()             {}
func (*ErrorReporter__Anon_153) Implements_rhino_ErrorReporter() {}

func (this_ *ErrorReporter__Anon_153) initFields_ErrorReporter__Anon_153() {
	_ = this_
}

func NewErrorReporter__Anon_153() *ErrorReporter__Anon_153 {
	this_ := &ErrorReporter__Anon_153{}
	this_.initFields_ErrorReporter__Anon_153()
	return this_
}

// Java method: .warning(java.lang.String message, java.lang.String sourceName, int line, int lineOffset)
func (this_ *ErrorReporter__Anon_153) Warning(message string, sourceName string, line int32, lineOffset int32) {
	_ = message
	_ = sourceName
	_ = line
	_ = lineOffset
}

// Java method: .error(java.lang.String message, java.lang.String sourceName, int line, int lineOffset)
func (this_ *ErrorReporter__Anon_153) M_Error(message string, sourceName string, line int32, lineOffset int32) {
	_ = message
	_ = sourceName
	_ = line
	_ = lineOffset
}

// Java class:
type ErrorReporter__Anon_154 struct {
	idHash_ atomic.Int32
}

func (this_ *ErrorReporter__Anon_154) BindSelf_(s any) {
	_ = s
}
func (this_ *ErrorReporter__Anon_154) IsJavaNil_() bool { return this_ == nil }
func (*ErrorReporter__Anon_154) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.ErrorReporter$2", ""
}
func (this_ *ErrorReporter__Anon_154) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ErrorReporter__Anon_154) HasIdentityEquals_()             {}
func (*ErrorReporter__Anon_154) Implements_rhino_ErrorReporter() {}

func (this_ *ErrorReporter__Anon_154) initFields_ErrorReporter__Anon_154() {
	_ = this_
}

func NewErrorReporter__Anon_154() *ErrorReporter__Anon_154 {
	this_ := &ErrorReporter__Anon_154{}
	this_.initFields_ErrorReporter__Anon_154()
	return this_
}

// Java method: .warning(java.lang.String message, java.lang.String sourceName, int line, int lineOffset)
func (this_ *ErrorReporter__Anon_154) Warning(message string, sourceName string, line int32, lineOffset int32) {
	_ = message
	_ = sourceName
	_ = line
	_ = lineOffset
	panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", j2gort.BoxString(message))))
}

// Java method: .error(java.lang.String message, java.lang.String sourceName, int line, int lineOffset)
func (this_ *ErrorReporter__Anon_154) M_Error(message string, sourceName string, line int32, lineOffset int32) {
	_ = message
	_ = sourceName
	_ = line
	_ = lineOffset
	panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", j2gort.BoxString(message))))
}
