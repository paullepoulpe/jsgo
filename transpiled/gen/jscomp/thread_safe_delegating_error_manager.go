/*
 * Copyright 2017 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/ThreadSafeDelegatingErrorManager.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java polymorphic class interface: com.google.javascript.jscomp.ThreadSafeDelegatingErrorManager
type ThreadSafeDelegatingErrorManager interface {
	Implements_jscomp_ThreadSafeDelegatingErrorManager()
	IsJavaNil_() bool
	GetThreadSafeDelegatingErrorManagerBase() *ThreadSafeDelegatingErrorManagerBase
	ErrorManager
}
func ToIface_ThreadSafeDelegatingErrorManager[T interface {
	ThreadSafeDelegatingErrorManager
	comparable
}](v T) ThreadSafeDelegatingErrorManager {
	var z T
	if v == z {
		return nil
	}
	return v
}

// Java class: com.google.javascript.jscomp.ThreadSafeDelegatingErrorManager
type ThreadSafeDelegatingErrorManagerBase struct {
	delegated ErrorManager `gson:"delegated"`
	idHash_ atomic.Int32
}

func (this_ *ThreadSafeDelegatingErrorManagerBase) BindSelf_(s any) {
	_ = s
}
func (this_ *ThreadSafeDelegatingErrorManagerBase) IsJavaNil_() bool { return this_ == nil }
func (*ThreadSafeDelegatingErrorManagerBase) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.ThreadSafeDelegatingErrorManager", "ThreadSafeDelegatingErrorManager" }
func (this_ *ThreadSafeDelegatingErrorManagerBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ThreadSafeDelegatingErrorManager", "com.google.javascript.jscomp.ThreadSafeDelegatingErrorManager", func(v_ any) bool {
		return j2gort.InstanceOf[ThreadSafeDelegatingErrorManager](v_)
	})
}
func (b *ThreadSafeDelegatingErrorManagerBase) GetThreadSafeDelegatingErrorManagerBase() *ThreadSafeDelegatingErrorManagerBase { return b }
func (*ThreadSafeDelegatingErrorManagerBase) Implements_jscomp_ThreadSafeDelegatingErrorManager() {}
func (*ThreadSafeDelegatingErrorManagerBase) Implements_jscomp_ErrorManager() {}
func (*ThreadSafeDelegatingErrorManagerBase) Implements_jscomp_ErrorHandler() {}

func (this_ *ThreadSafeDelegatingErrorManagerBase) initFields_ThreadSafeDelegatingErrorManager() {
	_ = this_
}

// Java constructor: ThreadSafeDelegatingErrorManager(com.google.javascript.jscomp.ErrorManager delegated)
func NewThreadSafeDelegatingErrorManager(delegated ErrorManager) ThreadSafeDelegatingErrorManager {
	this_ := &ThreadSafeDelegatingErrorManagerBase{}
	InitThreadSafeDelegatingErrorManager(this_, delegated)
	return this_
}

func InitThreadSafeDelegatingErrorManager(this_ *ThreadSafeDelegatingErrorManagerBase, delegated ErrorManager) {
	_ = delegated
	this_.initFields_ThreadSafeDelegatingErrorManager()
	this_.delegated = delegated
}

// Java method: ThreadSafeDelegatingErrorManager.report(com.google.javascript.jscomp.CheckLevel level, com.google.javascript.jscomp.JSError error)
func (this_ *ThreadSafeDelegatingErrorManagerBase) Report_2(level *CheckLevel, error_ *JSError) {
	unlock_ := j2gort.LockMonitor(this_)
	defer unlock_()
	_ = level
	_ = error_
	this_.delegated.Report_2(level, error_)
}

// Java method: ThreadSafeDelegatingErrorManager.generateReport()
func (this_ *ThreadSafeDelegatingErrorManagerBase) GenerateReport() {
	unlock_ := j2gort.LockMonitor(this_)
	defer unlock_()
	this_.delegated.GenerateReport()
}

// Java method: ThreadSafeDelegatingErrorManager.hasHaltingErrors()
func (this_ *ThreadSafeDelegatingErrorManagerBase) HasHaltingErrors() bool {
	return this_.delegated.HasHaltingErrors()
}

// Java method: ThreadSafeDelegatingErrorManager.getErrorCount()
func (this_ *ThreadSafeDelegatingErrorManagerBase) GetErrorCount() int32 {
	unlock_ := j2gort.LockMonitor(this_)
	defer unlock_()
	return this_.delegated.GetErrorCount()
}

// Java method: ThreadSafeDelegatingErrorManager.getWarningCount()
func (this_ *ThreadSafeDelegatingErrorManagerBase) GetWarningCount() int32 {
	unlock_ := j2gort.LockMonitor(this_)
	defer unlock_()
	return this_.delegated.GetWarningCount()
}

// Java method: ThreadSafeDelegatingErrorManager.getErrors()
func (this_ *ThreadSafeDelegatingErrorManagerBase) GetErrors() j2gort.Ext_com_google_common_collect_ImmutableList {
	unlock_ := j2gort.LockMonitor(this_)
	defer unlock_()
	return this_.delegated.GetErrors()
}

// Java method: ThreadSafeDelegatingErrorManager.getWarnings()
func (this_ *ThreadSafeDelegatingErrorManagerBase) GetWarnings() j2gort.Ext_com_google_common_collect_ImmutableList {
	unlock_ := j2gort.LockMonitor(this_)
	defer unlock_()
	return this_.delegated.GetWarnings()
}

// Java method: ThreadSafeDelegatingErrorManager.setTypedPercent(double typedPercent)
func (this_ *ThreadSafeDelegatingErrorManagerBase) SetTypedPercent(typedPercent float64) {
	unlock_ := j2gort.LockMonitor(this_)
	defer unlock_()
	_ = typedPercent
	this_.delegated.SetTypedPercent(typedPercent)
}

// Java method: ThreadSafeDelegatingErrorManager.getTypedPercent()
func (this_ *ThreadSafeDelegatingErrorManagerBase) GetTypedPercent() float64 {
	unlock_ := j2gort.LockMonitor(this_)
	defer unlock_()
	return this_.delegated.GetTypedPercent()
}

// Java method: ThreadSafeDelegatingErrorManager.shouldReportConformanceViolation(com.google.javascript.jscomp.Requirement requirement, com.google.common.base.Optional<com.google.javascript.jscomp.RequirementScopeEntry> allowlistEntry, com.google.javascript.jscomp.JSError diagnostic, com.google.javascript.jscomp.ConformanceConfig.LibraryLevelNonAllowlistedConformanceViolationsBehavior behavior, boolean isAllowlisted)
func (this_ *ThreadSafeDelegatingErrorManagerBase) ShouldReportConformanceViolation(requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement, allowlistEntry j2gort.Ext_com_google_common_base_Optional, diagnostic *JSError, behavior *j2gort.Ext_com_google_javascript_jscomp_ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior, isAllowlisted bool) bool {
	unlock_ := j2gort.LockMonitor(this_)
	defer unlock_()
	_ = requirement
	_ = allowlistEntry
	_ = diagnostic
	_ = behavior
	_ = isAllowlisted
	return this_.delegated.ShouldReportConformanceViolation(requirement, allowlistEntry, diagnostic, behavior, isAllowlisted)
}

