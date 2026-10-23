// Port of com/google/javascript/rhino/SimpleSourceFile.java @ closure-compiler v20261015

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
)

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
