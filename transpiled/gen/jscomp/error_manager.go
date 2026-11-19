/*
 * Copyright 2007 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/ErrorManager.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java interface: com.google.javascript.jscomp.ErrorManager
type ErrorManager interface {
	Implements_jscomp_ErrorManager()
	IsJavaNil_() bool
	ErrorHandler
	GenerateReport()
	GetErrorCount() int32
	GetWarningCount() int32
	GetErrors() j2gort.Ext_com_google_common_collect_ImmutableList
	GetWarnings() j2gort.Ext_com_google_common_collect_ImmutableList
	SetTypedPercent(float64)
	GetTypedPercent() float64
	HasHaltingErrors() bool
	ShouldReportConformanceViolation(*j2gort.Ext_com_google_javascript_jscomp_Requirement, j2gort.Ext_com_google_common_base_Optional, *JSError, *j2gort.Ext_com_google_javascript_jscomp_ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior, bool) bool
}
func ToIface_ErrorManager[T interface {
	ErrorManager
	comparable
}](v T) ErrorManager {
	var z T
	if v == z {
		return nil
	}
	return v
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ErrorManager", "com.google.javascript.jscomp.ErrorManager", func(v_ any) bool {
		return j2gort.InstanceOf[ErrorManager](v_)
	})
}
// Java default interface method: ErrorManager.hasHaltingErrors()
func IfaceDefault_ErrorManager_HasHaltingErrors(this_ ErrorManager) bool {
	_ = this_
	if ((this_.GetErrorCount() == int32(0))) {
		return false
	}
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.GetErrors()))
		for cur_1.HasNext() {
			rawElem_2 := cur_1.Next()
			error_ := j2gort.CastPtr[JSError](rawElem_2)
			_ = error_
			if ((error_.Type().level == StaticGet_CheckLevel_ERROR())) {
				return true
			}
		}
	}
	return false
}

// Java default interface method: ErrorManager.shouldReportConformanceViolation(com.google.javascript.jscomp.Requirement requirement, com.google.common.base.Optional<com.google.javascript.jscomp.RequirementScopeEntry> allowlistEntry, com.google.javascript.jscomp.JSError diagnostic, com.google.javascript.jscomp.ConformanceConfig.LibraryLevelNonAllowlistedConformanceViolationsBehavior behavior, boolean isAllowlisted)
func IfaceDefault_ErrorManager_ShouldReportConformanceViolation(this_ ErrorManager, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement, allowlistEntry j2gort.Ext_com_google_common_base_Optional, diagnostic *JSError, behavior *j2gort.Ext_com_google_javascript_jscomp_ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior, isAllowlisted bool) bool {
	_ = this_
	_ = requirement
	_ = allowlistEntry
	_ = diagnostic
	_ = behavior
	_ = isAllowlisted
	return (behavior != j2gort.Ext_com_google_javascript_jscomp_ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior_F_RECORD_ONLY)
}

