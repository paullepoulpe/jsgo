/*
 * Copyright 2018 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/ConformanceAllowlister.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

var lambdaSingleton_conformance_allowlister_1 j2gort.Ext_java_util_function_Function
var lambdaSingleton_conformance_allowlister_2 j2gort.Ext_java_util_function_Function

func init() {
	lambdaSingleton_conformance_allowlister_1 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any { return j2gort.BoxString(j2gort.CastPtr[JSError](mrp0_).SourceName()) })
	lambdaSingleton_conformance_allowlister_2 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any { return j2gort.BoxNode(j2gort.CastPtr[JSError](mrp0_).Node()) })
}

// Java class: com.google.javascript.jscomp.ConformanceAllowlister
type ConformanceAllowlister struct {
	idHash_ atomic.Int32
}

func (this_ *ConformanceAllowlister) BindSelf_(s any) {
	_ = s
}
func (this_ *ConformanceAllowlister) IsJavaNil_() bool { return this_ == nil }
func (*ConformanceAllowlister) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConformanceAllowlister", "ConformanceAllowlister"
}
func (this_ *ConformanceAllowlister) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ConformanceAllowlister) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConformanceAllowlister", "com.google.javascript.jscomp.ConformanceAllowlister", func(v_ any) bool {
		return j2gort.InstanceOf[*ConformanceAllowlister](v_)
	})
}

func (this_ *ConformanceAllowlister) initFields_ConformanceAllowlister() {
	_ = this_
}

// Java constructor: ConformanceAllowlister()
func NewConformanceAllowlister() *ConformanceAllowlister {
	this_ := &ConformanceAllowlister{}
	InitConformanceAllowlister(this_)
	return this_
}

func InitConformanceAllowlister(this_ *ConformanceAllowlister) {
	this_.initFields_ConformanceAllowlister()
}

// Java static method: ConformanceAllowlister.getViolatingPaths(com.google.javascript.jscomp.Compiler compiler, com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node ast, com.google.javascript.jscomp.Requirement requirement)
func ConformanceAllowlister_GetViolatingPaths(compiler *Compiler, externs pkg_rhino.Node, ast pkg_rhino.Node, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) j2gort.Ext_com_google_common_collect_ImmutableSet {
	_ = compiler
	_ = externs
	_ = ast
	_ = requirement
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](ConformanceAllowlister_GetConformanceErrors(compiler, externs, ast, requirement).Stream().Map(lambdaSingleton_conformance_allowlister_1).Collect_1(j2gort.Ext_com_google_common_collect_ImmutableSet_ToImmutableSet()))
}

// Java static method: ConformanceAllowlister.getViolatingNodes(com.google.javascript.jscomp.Compiler compiler, com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node ast, com.google.javascript.jscomp.Requirement requirement)
func ConformanceAllowlister_GetViolatingNodes(compiler *Compiler, externs pkg_rhino.Node, ast pkg_rhino.Node, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) j2gort.Ext_com_google_common_collect_ImmutableSet {
	_ = compiler
	_ = externs
	_ = ast
	_ = requirement
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](ConformanceAllowlister_GetConformanceErrors(compiler, externs, ast, requirement).Stream().Map(lambdaSingleton_conformance_allowlister_2).Collect_1(j2gort.Ext_com_google_common_collect_ImmutableSet_ToImmutableSet()))
}

// Java static method: ConformanceAllowlister.getConformanceErrors(com.google.javascript.jscomp.Compiler compiler, com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node ast, com.google.javascript.jscomp.Requirement requirement)
func ConformanceAllowlister_GetConformanceErrors(compiler *Compiler, externs pkg_rhino.Node, ast pkg_rhino.Node, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) j2gort.Ext_com_google_common_collect_ImmutableList {
	_ = compiler
	_ = externs
	_ = ast
	_ = requirement
	var cleanedRequirement *j2gort.Ext_com_google_javascript_jscomp_Requirement = j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_Requirement](j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_Requirement_Builder](requirement.ToBuilder()).ClearWhitelist().ClearAllowlist().ClearWhitelistRegexp().ClearAllowlistRegexp().ClearWhitelistEntry().ClearAllowlistEntry().SetSeverity(j2gort.Ext_com_google_javascript_jscomp_Requirement_Severity_F_ERROR).Build())
	_ = cleanedRequirement
	var cleanedConfig *j2gort.Ext_com_google_javascript_jscomp_ConformanceConfig = j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_ConformanceConfig](j2gort.Ext_com_google_javascript_jscomp_ConformanceConfig_NewBuilder_0().AddRequirement_1_Ext_com_google_javascript_jscomp_Requirement(cleanedRequirement).Build())
	_ = cleanedConfig
	var oldErrorManager ErrorManager = compiler.GetErrorManager()
	_ = oldErrorManager
	var errors *j2gort.Ext_com_google_common_collect_ImmutableList_Builder = (&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{})
	_ = errors
	func() {
		defer func() {
			compiler.SetErrorManager(oldErrorManager)
		}()
		compiler.SetErrorManager(NewConformanceAllowlister__Anon_15(errors, oldErrorManager))
		clinit_CheckConformance()
		var val_1 CheckConformance
		InitCheckConformance_3(&val_1, ToIface_AbstractCompiler(compiler), j2gort.ImmListOfPtr(cleanedConfig), CompilerOptions_ConformanceReportingMode(0))
		var check *CheckConformance = &val_1
		_ = check
		check.Process_2(externs, ast)
	}()
	return errors.BuildImmList_()
}

// Java class:
type ConformanceAllowlister__Anon_15 struct {
	ThreadSafeDelegatingErrorManagerBase
	cap_0_errors *j2gort.Ext_com_google_common_collect_ImmutableList_Builder
}

func (this_ *ConformanceAllowlister__Anon_15) BindSelf_(s any) {
	_ = s
}
func (this_ *ConformanceAllowlister__Anon_15) IsJavaNil_() bool { return this_ == nil }
func (*ConformanceAllowlister__Anon_15) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConformanceAllowlister$1", ""
}
func (this_ *ConformanceAllowlister__Anon_15) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ThreadSafeDelegatingErrorManagerBase.IdentityHash_()
}
func (*ConformanceAllowlister__Anon_15) HasIdentityEquals_()                                 {}
func (*ConformanceAllowlister__Anon_15) Implements_jscomp_ThreadSafeDelegatingErrorManager() {}
func (*ConformanceAllowlister__Anon_15) Implements_jscomp_ErrorManager()                     {}
func (*ConformanceAllowlister__Anon_15) Implements_jscomp_ErrorHandler()                     {}

func (this_ *ConformanceAllowlister__Anon_15) initFields_ConformanceAllowlister__Anon_15() {
	_ = this_
}

func NewConformanceAllowlister__Anon_15(cap_p_errors *j2gort.Ext_com_google_common_collect_ImmutableList_Builder, sp0_ ErrorManager) *ConformanceAllowlister__Anon_15 {
	this_ := &ConformanceAllowlister__Anon_15{}
	this_.cap_0_errors = cap_p_errors
	InitThreadSafeDelegatingErrorManager(&this_.ThreadSafeDelegatingErrorManagerBase, sp0_)
	this_.initFields_ConformanceAllowlister__Anon_15()
	return this_
}

// Java method: .shouldReportConformanceViolation(com.google.javascript.jscomp.Requirement requirement, com.google.common.base.Optional<com.google.javascript.jscomp.RequirementScopeEntry> whitelistEntry, com.google.javascript.jscomp.JSError diagnostic, com.google.javascript.jscomp.ConformanceConfig.LibraryLevelNonAllowlistedConformanceViolationsBehavior behavior, boolean isAllowlisted)
func (this_ *ConformanceAllowlister__Anon_15) ShouldReportConformanceViolation(requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement, whitelistEntry j2gort.Ext_com_google_common_base_Optional, diagnostic *JSError, behavior *j2gort.Ext_com_google_javascript_jscomp_ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior, isAllowlisted bool) bool {
	unlock_ := j2gort.LockMonitor(this_)
	defer unlock_()
	_ = requirement
	_ = whitelistEntry
	_ = diagnostic
	_ = behavior
	_ = isAllowlisted
	this_.cap_0_errors.AddVoid_(j2gort.PtrToAny(diagnostic))
	return false
}
