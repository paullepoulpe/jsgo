// Port of com/google/javascript/jscomp/DiagnosticType.java @ closure-compiler v20261021

/*
 * Copyright 2008 The Closure Compiler Authors.
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

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// DiagnosticType is the type of a compile or analysis error.
type DiagnosticType struct {
	// The error type. Used as the BugPattern and BugInstance types by
	// BugBot's XML
	key string

	// The default way to format errors. The style of format is
	// java.text.MessageFormat.
	format string

	// The default reporting level for this diagnostic
	level CheckLevel

	msgFormat atomic.Pointer[jutil.MessageFormat]
}

// Java: DiagnosticType.serialVersionUID
const DiagnosticType_serialVersionUID int64 = 1

// DiagnosticType_Error creates a DiagnosticType at level CheckLevel.ERROR.
//
// Java: DiagnosticType.error(String, String)
func DiagnosticType_Error(name string, descriptionFormat string) *DiagnosticType {
	return DiagnosticType_Make(name, CheckLevel_ERROR, descriptionFormat)
}

// DiagnosticType_Warning creates a DiagnosticType at level
// CheckLevel.WARNING.
//
// Java: DiagnosticType.warning(String, String)
func DiagnosticType_Warning(name string, descriptionFormat string) *DiagnosticType {
	return DiagnosticType_Make(name, CheckLevel_WARNING, descriptionFormat)
}

// DiagnosticType_Disabled creates a DiagnosticType at level CheckLevel.OFF.
//
// Java: DiagnosticType.disabled(String, String)
func DiagnosticType_Disabled(name string, descriptionFormat string) *DiagnosticType {
	return DiagnosticType_Make(name, CheckLevel_OFF, descriptionFormat)
}

// DiagnosticType_Make creates a DiagnosticType at a given CheckLevel.
//
// Java: DiagnosticType.make(String, CheckLevel, String)
func DiagnosticType_Make(name string, level CheckLevel, descriptionFormat string) *DiagnosticType {
	return newDiagnosticType(name, level, descriptionFormat)
}

// Create a DiagnosticType. Private to force use of static factory methods.
func newDiagnosticType(key string, level CheckLevel, format string) *DiagnosticType {
	return &DiagnosticType{key: key, level: level, format: format}
}

// Java: DiagnosticType.format(String...)
func (t *DiagnosticType) Format(arguments ...string) string {
	mf := t.msgFormat.Load()
	if mf == nil {
		mf = jutil.NewMessageFormat(t.format)
		t.msgFormat.Store(mf)
	}
	return mf.FormatStrings(arguments...)
}

// Java: DiagnosticType.equals(Object)
func (t *DiagnosticType) Equals(typ any) bool {
	diagnosticType, ok := typ.(*DiagnosticType)
	return ok && diagnosticType != nil && diagnosticType.key == t.key
}

// Java: DiagnosticType.hashCode()
func (t *DiagnosticType) HashCode() int32 {
	return jutil.StringHashCode(t.key)
}

// Java: DiagnosticType.compareTo(DiagnosticType)
func (t *DiagnosticType) CompareTo(diagnosticType *DiagnosticType) int {
	return jutil.StringCompareTo(t.key, diagnosticType.key)
}

// Java: DiagnosticType.toString()
func (t *DiagnosticType) String() string {
	return t.key + ": " + t.format
}
