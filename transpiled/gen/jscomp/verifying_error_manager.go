/*
 * Copyright 2026 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/VerifyingErrorManager.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.VerifyingErrorManager
type VerifyingErrorManager struct {
	ThreadSafeDelegatingErrorManagerBase
	unmatchedExpectations j2gort.Ext_java_util_List    `gson:"unmatchedExpectations"`
	formatter             *LightweightMessageFormatter `gson:"formatter"`
}

func (this_ *VerifyingErrorManager) BindSelf_(s any) {
	_ = s
}
func (this_ *VerifyingErrorManager) IsJavaNil_() bool { return this_ == nil }
func (*VerifyingErrorManager) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.VerifyingErrorManager", "VerifyingErrorManager"
}
func (this_ *VerifyingErrorManager) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ThreadSafeDelegatingErrorManagerBase.IdentityHash_()
}
func (*VerifyingErrorManager) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.VerifyingErrorManager", "com.google.javascript.jscomp.VerifyingErrorManager", func(v_ any) bool {
		return j2gort.InstanceOf[*VerifyingErrorManager](v_)
	})
}
func (*VerifyingErrorManager) Implements_jscomp_ThreadSafeDelegatingErrorManager() {}
func (*VerifyingErrorManager) Implements_jscomp_ErrorManager()                     {}
func (*VerifyingErrorManager) Implements_jscomp_ErrorHandler()                     {}

func (this_ *VerifyingErrorManager) initFields_VerifyingErrorManager() {
	_ = this_
}

// Java constructor: VerifyingErrorManager(com.google.javascript.jscomp.ErrorManager delegate, java.util.List<java.lang.String> expectedDiagnostics)
func NewVerifyingErrorManager(delegate ErrorManager, expectedDiagnostics j2gort.Ext_java_util_List) *VerifyingErrorManager {
	this_ := &VerifyingErrorManager{}
	InitVerifyingErrorManager(this_, delegate, expectedDiagnostics)
	return this_
}

func InitVerifyingErrorManager(this_ *VerifyingErrorManager, delegate ErrorManager, expectedDiagnostics j2gort.Ext_java_util_List) {
	_ = delegate
	_ = expectedDiagnostics
	InitThreadSafeDelegatingErrorManager(&this_.ThreadSafeDelegatingErrorManagerBase, delegate)
	this_.initFields_VerifyingErrorManager()
	this_.formatter = LightweightMessageFormatter_WithoutSource()
	this_.unmatchedExpectations = j2gort.NewExt_java_util_ArrayList()
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(expectedDiagnostics))
		for cur_1.HasNext() {
			exp := cur_1.NextStr()
			_ = exp
			this_.unmatchedExpectations.Add_1(j2gort.PtrToAny(j2gort.Ext_java_util_regex_Pattern_Compile_1(exp)))
		}
	}
}

// Java method: VerifyingErrorManager.report(com.google.javascript.jscomp.CheckLevel level, com.google.javascript.jscomp.JSError error)
func (this_ *VerifyingErrorManager) Report_2(level *CheckLevel, error_ *JSError) {
	unlock_ := j2gort.LockMonitor(this_)
	defer unlock_()
	_ = level
	_ = error_
	var formattedError string = this_.priv_VerifyingErrorManager_formatErrorForMatching_5893(level, error_)
	_ = formattedError
	var matchedPattern *j2gort.Ext_java_util_regex_Pattern = nil
	_ = matchedPattern
	var iterator j2gort.Ext_java_util_Iterator = j2gort.AsIface[j2gort.Ext_java_util_Iterator](this_.unmatchedExpectations.Iterator())
	_ = iterator
	for iterator.HasNext() {
		var p *j2gort.Ext_java_util_regex_Pattern = j2gort.CastPtr[j2gort.Ext_java_util_regex_Pattern](iterator.Next())
		_ = p
		if p.Matcher(j2gort.BoxString(formattedError)).Find_0() {
			if matchedPattern == nil {
				matchedPattern = p
				iterator.Remove()
				continue
			}
			{
				var iifeCond_1 bool
				{

					sr_2 := matchedPattern.String()
					se_3 := p.String()
					iifeCond_1 = j2gort.CheckNotNullString(sr_2) == se_3
				}
				if iifeCond_1 {
					continue
				}
			}
			this_.ThreadSafeDelegatingErrorManagerBase.Report_2(StaticGet_CheckLevel_ERROR(), JSError_Make_2(StaticGet_AbstractCommandLineRunner_AMBIGUOUS_EXPECTATION(), []string{formattedError, matchedPattern.String(), p.String()}))
			iterator.Remove()
		}
	}
	if matchedPattern == nil {
		this_.ThreadSafeDelegatingErrorManagerBase.Report_2(level, error_)
	}
}

// Java method: VerifyingErrorManager.generateReport()
func (this_ *VerifyingErrorManager) GenerateReport() {
	unlock_ := j2gort.LockMonitor(this_)
	defer unlock_()
	this_.priv_VerifyingErrorManager_reportMissingExpectations_5892()
	this_.ThreadSafeDelegatingErrorManagerBase.GenerateReport()
}

// Java method: VerifyingErrorManager.reportMissingExpectations()
func (this_ *VerifyingErrorManager) priv_VerifyingErrorManager_reportMissingExpectations_5892() {
	unlock_ := j2gort.LockMonitor(this_)
	defer unlock_()
	{
		cur_4 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.unmatchedExpectations))
		for cur_4.HasNext() {
			rawElem_5 := cur_4.Next()
			p := j2gort.CastPtr[j2gort.Ext_java_util_regex_Pattern](rawElem_5)
			_ = p
			this_.ThreadSafeDelegatingErrorManagerBase.Report_2(StaticGet_CheckLevel_ERROR(), JSError_Make_2(StaticGet_AbstractCommandLineRunner_EXPECTED_DIAGNOSTIC_NOT_FOUND(), []string{p.String()}))
		}
	}
	this_.unmatchedExpectations.Clear()
}

// Java method: VerifyingErrorManager.formatErrorForMatching(com.google.javascript.jscomp.CheckLevel level, com.google.javascript.jscomp.JSError error)
func (this_ *VerifyingErrorManager) priv_VerifyingErrorManager_formatErrorForMatching_5893(level *CheckLevel, error_ *JSError) string {
	_ = level
	_ = error_
	if level == StaticGet_CheckLevel_ERROR() {
		return this_.formatter.FormatError(error_)
	} else {
		return this_.formatter.FormatWarning(error_)
	}
}
