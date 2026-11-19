/*
 * Copyright 2014 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/CheckConformance.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_check_conformance_1 j2gort.Ext_java_util_function_Function
var lambdaSingleton_check_conformance_2 CheckConformance_Precondition

func init() {
	lambdaSingleton_check_conformance_1 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_5 any) any {
		e := j2gort.AsIface[j2gort.Ext_java_util_Map_Entry](lp_5)
		_ = e
		return NewCheckConformance_Category(j2gort.AsIface[CheckConformance_Precondition](e.GetKey()), j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](e.GetValue())))
	})
	lambdaSingleton_check_conformance_2 = NewLambda_CheckConformance_Precondition(func(n pkg_rhino.Node) bool {
		_ = n
		return (n.IsClass() || NodeUtil_IsConstructor(n))
	})
}

// Java class: com.google.javascript.jscomp.CheckConformance
type CheckConformance struct {
	compiler        AbstractCompiler                                                                                                                                             `gson:"compiler"`
	categories      []*CheckConformance_Category                                                                                                                                 `gson:"categories"`
	mergedBehaviors j2gort.Ext_java_util_Map                                                                                                                                     `gson:"mergedBehaviors"`
	ruleToBehavior  j2gort.OrderedMap[CheckConformance_Rule, *j2gort.Ext_com_google_javascript_jscomp_ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior] `gson:"ruleToBehavior"`
	idHash_         atomic.Int32
}

func (this_ *CheckConformance) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckConformance) IsJavaNil_() bool { return this_ == nil }
func (*CheckConformance) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CheckConformance", "CheckConformance"
}
func (this_ *CheckConformance) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CheckConformance) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CheckConformance", "com.google.javascript.jscomp.CheckConformance", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckConformance](v_)
	})
}
func (*CheckConformance) Implements_jscomp_NodeTraversal_Callback() {}
func (*CheckConformance) Implements_jscomp_CompilerPass()           {}

var static_CheckConformance_CONFORMANCE_ERROR *DiagnosticType

func StaticGet_CheckConformance_CONFORMANCE_ERROR() *DiagnosticType {
	clinit_CheckConformance()
	return static_CheckConformance_CONFORMANCE_ERROR
}
func StaticSet_CheckConformance_CONFORMANCE_ERROR(v *DiagnosticType) {
	clinit_CheckConformance()
	static_CheckConformance_CONFORMANCE_ERROR = v
}

var static_CheckConformance_CONFORMANCE_VIOLATION *DiagnosticType

func StaticGet_CheckConformance_CONFORMANCE_VIOLATION() *DiagnosticType {
	clinit_CheckConformance()
	return static_CheckConformance_CONFORMANCE_VIOLATION
}
func StaticSet_CheckConformance_CONFORMANCE_VIOLATION(v *DiagnosticType) {
	clinit_CheckConformance()
	static_CheckConformance_CONFORMANCE_VIOLATION = v
}

var static_CheckConformance_CONFORMANCE_POSSIBLE_VIOLATION *DiagnosticType

func StaticGet_CheckConformance_CONFORMANCE_POSSIBLE_VIOLATION() *DiagnosticType {
	clinit_CheckConformance()
	return static_CheckConformance_CONFORMANCE_POSSIBLE_VIOLATION
}
func StaticSet_CheckConformance_CONFORMANCE_POSSIBLE_VIOLATION(v *DiagnosticType) {
	clinit_CheckConformance()
	static_CheckConformance_CONFORMANCE_POSSIBLE_VIOLATION = v
}

var static_CheckConformance_INVALID_REQUIREMENT_SPEC *DiagnosticType

func StaticGet_CheckConformance_INVALID_REQUIREMENT_SPEC() *DiagnosticType {
	clinit_CheckConformance()
	return static_CheckConformance_INVALID_REQUIREMENT_SPEC
}
func StaticSet_CheckConformance_INVALID_REQUIREMENT_SPEC(v *DiagnosticType) {
	clinit_CheckConformance()
	static_CheckConformance_INVALID_REQUIREMENT_SPEC = v
}

var static_CheckConformance_EXTENDABLE_FIELDS j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_CheckConformance_EXTENDABLE_FIELDS() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_CheckConformance()
	return static_CheckConformance_EXTENDABLE_FIELDS
}
func StaticSet_CheckConformance_EXTENDABLE_FIELDS(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_CheckConformance()
	static_CheckConformance_EXTENDABLE_FIELDS = v
}

var clinitDone_CheckConformance uint32
var clinitState_CheckConformance j2gort.ClinitState

func clinit_CheckConformance() {
	if atomic.LoadUint32(&clinitDone_CheckConformance) != 0 {
		return
	}
	clinitSlow_CheckConformance()
}

func clinitSlow_CheckConformance() {
	if !j2gort.EnterClinit(&clinitState_CheckConformance, &clinitDone_CheckConformance) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CheckConformance, &clinitDone_CheckConformance)
	static_CheckConformance_CONFORMANCE_ERROR = DiagnosticType_Error("JSC_CONFORMANCE_ERROR", "Violation: {0}{1}{2}")
	static_CheckConformance_CONFORMANCE_VIOLATION = DiagnosticType_Warning("JSC_CONFORMANCE_VIOLATION", "Violation: {0}{1}{2}")
	static_CheckConformance_CONFORMANCE_POSSIBLE_VIOLATION = DiagnosticType_Warning("JSC_CONFORMANCE_POSSIBLE_VIOLATION", "Possible violation: {0}{1}{2}")
	static_CheckConformance_INVALID_REQUIREMENT_SPEC = DiagnosticType_Error("JSC_INVALID_REQUIREMENT_SPEC", "Invalid requirement. Reason: {0}\nRequirement spec:\n{1}")
	static_CheckConformance_EXTENDABLE_FIELDS = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_7("config_file", "error_message", "extends", "only_apply_to", "only_apply_to_regexp", "whitelist", []any{"whitelist_regexp", "allowlist", "allowlist_regexp", "value"})
}

func (this_ *CheckConformance) initFields_CheckConformance() {
	_ = this_
	this_.mergedBehaviors = j2gort.NewExt_java_util_LinkedHashMap()
	this_.ruleToBehavior = j2gort.NewOrderedMap[CheckConformance_Rule, *j2gort.

		// Java constructor: CheckConformance(com.google.javascript.jscomp.AbstractCompiler compiler)
		Ext_com_google_javascript_jscomp_ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior]()
}

func NewCheckConformance_1(compiler AbstractCompiler) *CheckConformance {
	clinit_CheckConformance()
	this_ := &CheckConformance{}
	InitCheckConformance_1(this_, compiler)
	return this_
}

func InitCheckConformance_1(this_ *CheckConformance, compiler AbstractCompiler) {
	_ = compiler
	InitCheckConformance_3(this_, compiler, j2gort.Ext_com_google_common_collect_ImmutableList_Of_0(), CompilerOptions_ConformanceReportingMode(0))
}

// Java constructor: CheckConformance(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.common.collect.ImmutableList<com.google.javascript.jscomp.ConformanceConfig> configs, com.google.javascript.jscomp.CompilerOptions.ConformanceReportingMode reportingMode)
func NewCheckConformance_3(compiler AbstractCompiler, configs j2gort.Ext_com_google_common_collect_ImmutableList, reportingMode CompilerOptions_ConformanceReportingMode) *CheckConformance {
	clinit_CheckConformance()
	this_ := &CheckConformance{}
	InitCheckConformance_3(this_, compiler, configs, reportingMode)
	return this_
}

func InitCheckConformance_3(this_ *CheckConformance, compiler AbstractCompiler, configs j2gort.Ext_com_google_common_collect_ImmutableList, reportingMode CompilerOptions_ConformanceReportingMode) {
	_ = compiler
	_ = configs
	_ = reportingMode
	this_.initFields_CheckConformance()
	this_.compiler = compiler
	this_.categories = j2gort.ImmListToSliceNullable[*CheckConformance_Category](this_.priv_CheckConformance_initRules_1048(compiler, configs, reportingMode))
}

// Java method: CheckConformance.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *CheckConformance) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	if len(this_.categories) != 0 {
		NodeTraversal_TraverseRoots(this_.compiler, this_, externs, root)
	}
}

// Java method: CheckConformance.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckConformance) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	return (!(n.IsScript()) || (this_.priv_CheckConformance_isScriptOfInterest_1047(t.GetInput().GetSourceFile()) && !(javalang.StringEndsWith(j2gort.CheckNotNullString(t.GetSourceName()), "tsmes.closure.js"))))
}

// Java method: CheckConformance.isScriptOfInterest(com.google.javascript.jscomp.SourceFile sf)
func (this_ *CheckConformance) priv_CheckConformance_isScriptOfInterest_1047(sf *SourceFile) bool {
	_ = sf
	return (!(sf.IsWeak()) && !(sf.IsExtern()))
}

// Java method: CheckConformance.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckConformance) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	{
		var c int32 = (int32(len(this_.categories)) - int32(1))
		_ = c
		for ; c >= int32(0); c = c - 1 {
			var category *CheckConformance_Category = this_.categories[c]
			_ = category
			if category.precondition.ShouldCheck(n) {
				{
					var r int32 = (int32(len(category.rules)) - int32(1))
					_ = r
					for ; r >= int32(0); r = r - 1 {
						var rule CheckConformance_Rule = category.rules[r]
						_ = rule
						var behavior *j2gort.Ext_com_google_javascript_jscomp_ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior = j2gort.CheckNotNullPtr(this_.ruleToBehavior.GetOrDefault(rule, j2gort.Ext_com_google_javascript_jscomp_ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior_F_UNSPECIFIED))
						_ = behavior
						rule.Check(t, n, behavior)
					}
				}
			}
		}
	}
}

// Java method: CheckConformance.initRules(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.common.collect.ImmutableList<com.google.javascript.jscomp.ConformanceConfig> configs, com.google.javascript.jscomp.CompilerOptions.ConformanceReportingMode reportingMode)
func (this_ *CheckConformance) priv_CheckConformance_initRules_1048(compiler AbstractCompiler, configs j2gort.Ext_com_google_common_collect_ImmutableList, reportingMode CompilerOptions_ConformanceReportingMode) j2gort.Ext_com_google_common_collect_ImmutableList {
	_ = compiler
	_ = configs
	_ = reportingMode
	var builder *j2gort.Ext_com_google_common_collect_HashMultimap = j2gort.Ext_com_google_common_collect_HashMultimap_Create_0()
	_ = builder
	var isLibraryLevelReportingMode bool = ((reportingMode != 0) && func() bool {
		re_1 := reportingMode
		ae_2 := StaticGet_CompilerOptions_ConformanceReportingMode_RESPECT_LIBRARY_LEVEL_BEHAVIOR_SPECIFIED_IN_CONFIG()
		if re_1 == 0 {
			j2gort.ThrowNPE()
		}
		return re_1 == ae_2
	}())
	_ = isLibraryLevelReportingMode
	if isLibraryLevelReportingMode {
		if !(CheckConformance_ValidateBehaviorSettingOfConfigs(compiler, configs, reportingMode)) {
			return j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
		}
	}
	var mergedRequirements j2gort.Ext_java_util_List = this_.MergeRequirements(compiler, j2gort.AsIface[j2gort.Ext_java_util_List](configs), isLibraryLevelReportingMode)
	_ = mergedRequirements
	{
		cur_3 := j2gort.NewIterCursor(j2gort.CheckNotNull(mergedRequirements))
		for cur_3.HasNext() {
			rawElem_4 := cur_3.Next()
			requirement := j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_Requirement](rawElem_4)
			_ = requirement
			var behavior *j2gort.Ext_com_google_javascript_jscomp_ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior = j2gort.Ext_com_google_javascript_jscomp_ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior_F_UNSPECIFIED
			_ = behavior
			var rule CheckConformance_Rule = CheckConformance_InitRule(compiler, requirement)
			_ = rule
			if rule != nil && !rule.IsJavaNil_() {
				builder.Put(rule.GetPrecondition(), rule)
				if isLibraryLevelReportingMode && this_.mergedBehaviors.ContainsKey(j2gort.PtrToAny(requirement)) {
					behavior = j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior](this_.mergedBehaviors.Get(j2gort.PtrToAny(requirement)))
					j2gort.CheckNotNullPtr(behavior)
				}
				this_.ruleToBehavior.Put(rule, behavior)
			}
		}
	}
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableList](j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.AsIface[j2gort.Ext_java_util_Map](builder.AsMap()).EntrySet()).Stream().Map(lambdaSingleton_check_conformance_1).Collect_1(j2gort.Ext_com_google_common_collect_ImmutableList_ToImmutableList()))
}

// Java static method: CheckConformance.validateBehaviorSettingOfConfigs(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.common.collect.ImmutableList<com.google.javascript.jscomp.ConformanceConfig> configs, com.google.javascript.jscomp.CompilerOptions.ConformanceReportingMode reportingMode)
func CheckConformance_ValidateBehaviorSettingOfConfigs(compiler AbstractCompiler, configs j2gort.Ext_com_google_common_collect_ImmutableList, reportingMode CompilerOptions_ConformanceReportingMode) bool {
	clinit_CheckConformance()
	_ = compiler
	_ = configs
	_ = reportingMode
	{
		var iifeCond_1 bool
		{

			re_6 := reportingMode
			ae_7 := StaticGet_CompilerOptions_ConformanceReportingMode_RESPECT_LIBRARY_LEVEL_BEHAVIOR_SPECIFIED_IN_CONFIG()
			if re_6 == 0 {
				j2gort.ThrowNPE()
			}
			iifeCond_1 = re_6 == ae_7
		}
		if !(iifeCond_1) {
			j2gort.CheckState(false)
		}
	}
	var baseRequirementBehaviors j2gort.OrderedMap[string, *j2gort.Ext_com_google_javascript_jscomp_ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior] = j2gort.NewOrderedMap[string, *j2gort.Ext_com_google_javascript_jscomp_ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior]()
	_ = baseRequirementBehaviors
	{
		cur_8 := j2gort.NewIterCursor(j2gort.CheckNotNull(configs))
		for cur_8.HasNext() {
			rawElem_9 := cur_8.Next()
			config := j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_ConformanceConfig](rawElem_9)
			_ = config
			if !(config.HasLibraryLevelNonAllowlistedConformanceViolationsBehavior()) {
				continue
			}
			var behavior *j2gort.Ext_com_google_javascript_jscomp_ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior = config.GetLibraryLevelNonAllowlistedConformanceViolationsBehavior()
			_ = behavior
			{
				cur_10 := j2gort.NewIterCursor(j2gort.CheckNotNull(config.GetRequirementList()))
				for cur_10.HasNext() {
					rawElem_11 := cur_10.Next()
					requirement := j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_Requirement](rawElem_11)
					_ = requirement
					if requirement.HasRuleId() {
						baseRequirementBehaviors.Put(requirement.GetRuleId(), behavior)
					}
				}
			}
		}
	}
	{
		cur_12 := j2gort.NewIterCursor(j2gort.CheckNotNull(configs))
		for cur_12.HasNext() {
			rawElem_13 := cur_12.Next()
			config_2 := j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_ConformanceConfig](rawElem_13)
			_ = config_2
			if !(config_2.HasLibraryLevelNonAllowlistedConformanceViolationsBehavior()) {
				continue
			}
			{
				cur_14 := j2gort.NewIterCursor(j2gort.CheckNotNull(config_2.GetRequirementList()))
				for cur_14.HasNext() {
					rawElem_15 := cur_14.Next()
					requirement_2 := j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_Requirement](rawElem_15)
					_ = requirement_2
					if requirement_2.HasExtends() {
						var extendingBehavior *j2gort.Ext_com_google_javascript_jscomp_ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior = config_2.GetLibraryLevelNonAllowlistedConformanceViolationsBehavior()
						_ = extendingBehavior
						var baseRuleId string = requirement_2.GetExtends()
						_ = baseRuleId
						if baseRequirementBehaviors.ContainsKey(baseRuleId) {
							var baseBehavior *j2gort.Ext_com_google_javascript_jscomp_ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior = baseRequirementBehaviors.Get(baseRuleId)
							_ = baseBehavior
							if !(j2gort.RefEquals(baseBehavior, extendingBehavior)) {
								CheckConformance_ReportInvalidRequirement(compiler, requirement_2, "extending rule's config may not specify a different value of 'library_level_non_allowlisted_conformance_violations_behavior' than the base rule's config. Skipping all conformance checks.")
								return false
							}
						}
					}
				}
			}
		}
	}
	return true
}

// Java method: CheckConformance.mergeRequirements(com.google.javascript.jscomp.AbstractCompiler compiler, java.util.List<com.google.javascript.jscomp.ConformanceConfig> configs, boolean isLibraryLevelConformanceReportingMode)
func (this_ *CheckConformance) MergeRequirements(compiler AbstractCompiler, configs j2gort.Ext_java_util_List, isLibraryLevelConformanceReportingMode bool) j2gort.Ext_java_util_List {
	_ = compiler
	_ = configs
	_ = isLibraryLevelConformanceReportingMode
	var rootRequirements j2gort.OrderedMap[*j2gort.Ext_com_google_javascript_jscomp_Requirement_Builder, *j2gort.Ext_com_google_javascript_jscomp_ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior] = j2gort.NewOrderedMap[*j2gort.Ext_com_google_javascript_jscomp_Requirement_Builder, *j2gort.Ext_com_google_javascript_jscomp_ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior]()
	_ = rootRequirements
	var extendable j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_LinkedHashMap()
	_ = extendable
	{
		cur_16 := j2gort.NewIterCursor(j2gort.CheckNotNull(configs))
		for cur_16.HasNext() {
			rawElem_17 := cur_16.Next()
			config := j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_ConformanceConfig](rawElem_17)
			_ = config
			{
				cur_18 := j2gort.NewIterCursor(j2gort.CheckNotNull(config.GetRequirementList()))
				for cur_18.HasNext() {
					rawElem_19 := cur_18.Next()
					requirement := j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_Requirement](rawElem_19)
					_ = requirement
					var builder *j2gort.Ext_com_google_javascript_jscomp_Requirement_Builder = j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_Requirement_Builder](requirement.ToBuilder())
					_ = builder
					if requirement.HasRuleId() {
						if !(CheckConformance_ValidateExtendableRequirement(compiler, requirement, extendable)) {
							continue
						}
						j2gort.MapPutStr(extendable, requirement.GetRuleId(), j2gort.PtrToAny(builder))
					}
					var hasBehavior bool = config.HasLibraryLevelNonAllowlistedConformanceViolationsBehavior()
					_ = hasBehavior
					var behavior *j2gort.Ext_com_google_javascript_jscomp_ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior = config.GetLibraryLevelNonAllowlistedConformanceViolationsBehavior()
					_ = behavior
					if !(requirement.HasExtends()) {
						if !(isLibraryLevelConformanceReportingMode) {
							rootRequirements.Put(builder, j2gort.Ext_com_google_javascript_jscomp_ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior_F_UNSPECIFIED)
						} else {
							rootRequirements.Put(builder, func() *j2gort.Ext_com_google_javascript_jscomp_ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior {
								if hasBehavior {
									return behavior
								}
								return j2gort.Ext_com_google_javascript_jscomp_ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior_F_UNSPECIFIED
							}())
						}
					}
					var libraryLevelConformanceEnabled bool = ((isLibraryLevelConformanceReportingMode && hasBehavior) && (behavior == j2gort.Ext_com_google_javascript_jscomp_ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior_F_RECORD_ONLY))
					_ = libraryLevelConformanceEnabled
					if (!(requirement.HasRuleId()) && !(requirement.HasExtends())) && libraryLevelConformanceEnabled {
						CheckConformance_ReportInvalidRequirement(compiler, requirement, "Library-level conformance requirements require their rule_id or extends field set.")
					}
				}
			}
		}
	}
	{
		cur_20 := j2gort.NewIterCursor(j2gort.CheckNotNull(configs))
		for cur_20.HasNext() {
			rawElem_21 := cur_20.Next()
			config_2 := j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_ConformanceConfig](rawElem_21)
			_ = config_2
			{
				cur_22 := j2gort.NewIterCursor(j2gort.CheckNotNull(config_2.GetRequirementList()))
				for cur_22.HasNext() {
					rawElem_23 := cur_22.Next()
					requirement_2 := j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_Requirement](rawElem_23)
					_ = requirement_2
					if requirement_2.HasExtends() {
						var extendableRequirement *j2gort.Ext_com_google_javascript_jscomp_Requirement_Builder = j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_Requirement_Builder](j2gort.MapGetStr(extendable, requirement_2.GetExtends()))
						_ = extendableRequirement
						if !(CheckConformance_ValidateExtendingRequirement(compiler, requirement_2, extendableRequirement)) {
							continue
						}
						if config_2.HasLibraryLevelNonAllowlistedConformanceViolationsBehavior() {
							rootRequirements.Put(extendableRequirement, config_2.GetLibraryLevelNonAllowlistedConformanceViolationsBehavior())
						}
						CheckConformance_MergeExtendingRequirementIntoExtended(extendableRequirement, requirement_2)
					}
				}
			}
		}
	}
	var cleanedUpRootRequirements []*j2gort.Ext_com_google_javascript_jscomp_Requirement = j2gort.MakeNonNilSlice[*j2gort.Ext_com_google_javascript_jscomp_Requirement](rootRequirements.Size())
	_ = cleanedUpRootRequirements
	{
		cur_24 := rootRequirements.IterEntries()
		for cur_24.NextEntry() {
			var builder_2 *j2gort.Ext_com_google_javascript_jscomp_Requirement_Builder = cur_24.Key()
			_ = builder_2
			var behavior_2 *j2gort.Ext_com_google_javascript_jscomp_ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior = cur_24.Value()
			_ = behavior_2
			CheckConformance_RemoveDuplicates(builder_2)
			var requirement_3 *j2gort.Ext_com_google_javascript_jscomp_Requirement = j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_Requirement](builder_2.Build())
			_ = requirement_3
			if (isLibraryLevelConformanceReportingMode && (behavior_2 != nil)) && (behavior_2 != j2gort.Ext_com_google_javascript_jscomp_ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior_F_UNSPECIFIED) {
				this_.mergedBehaviors.Put_2(j2gort.PtrToAny(requirement_3), j2gort.PtrToAny(behavior_2))
			}
			cleanedUpRootRequirements = append(cleanedUpRootRequirements, requirement_3)
		}
	}
	return j2gort.WrapTypedArrayList(cleanedUpRootRequirements)
}

// Java static method: CheckConformance.mergeExtendingRequirementIntoExtended(com.google.javascript.jscomp.Requirement.Builder extendedRequirement, com.google.javascript.jscomp.Requirement extendingRequirement)
func CheckConformance_MergeExtendingRequirementIntoExtended(extendedRequirement *j2gort.Ext_com_google_javascript_jscomp_Requirement_Builder, extendingRequirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) {
	clinit_CheckConformance()
	_ = extendedRequirement
	_ = extendingRequirement
	{
		var iifeCond_1 bool
		{

			sr_25 := extendingRequirement.GetExtends()
			se_26 := extendedRequirement.GetRuleId()
			iifeCond_1 = j2gort.CheckNotNullString(sr_25) == se_26
		}
		if !(iifeCond_1) {
			j2gort.CheckStateStr(false, "Attempting to merge an requirement with the wrong extended requirement")
		}
	}
	extendedRequirement.AddAllWhitelist(j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](extendingRequirement.GetWhitelistList())).AddAllWhitelistRegexp(j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](extendingRequirement.GetWhitelistRegexpList())).AddAllAllowlist(j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](extendingRequirement.GetAllowlistList())).AddAllAllowlistRegexp(j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](extendingRequirement.GetAllowlistRegexpList())).AddAllOnlyApplyTo(j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](extendingRequirement.GetOnlyApplyToList())).AddAllOnlyApplyToRegexp(j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](extendingRequirement.GetOnlyApplyToRegexpList())).AddAllWhitelistEntry(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](extendingRequirement.GetWhitelistEntryList())).AddAllAllowlistEntry(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](extendingRequirement.GetAllowlistEntryList())).AddAllValue(j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](extendingRequirement.GetValueList())).AddAllConfigFile(j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](extendingRequirement.GetConfigFileList()))
	if extendingRequirement.HasErrorMessage() && extendedRequirement.HasErrorMessage() {
		extendedRequirement.SetErrorMessage((j2gort.StrValString(extendedRequirement.GetErrorMessage()) + " " + j2gort.StrValString(extendingRequirement.GetErrorMessage())))
	}
}

// Java static method: CheckConformance.validateExtendableRequirement(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.Requirement requirement, java.util.Map<java.lang.String,com.google.javascript.jscomp.Requirement.Builder> extendable)
func CheckConformance_ValidateExtendableRequirement(compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement, extendable j2gort.Ext_java_util_Map) bool {
	clinit_CheckConformance()
	_ = compiler
	_ = requirement
	_ = extendable
	if !(requirement.HasRuleId()) {
		j2gort.CheckStateStr(false, "Extendable requirement must have a rule_id")
	}
	if j2gort.CheckNotNullString(requirement.GetRuleId()) == "" {
		CheckConformance_ReportInvalidRequirement(compiler, requirement, "empty rule_id")
		return false
	}
	if j2gort.MapContainsKeyStr(extendable, requirement.GetRuleId()) {
		CheckConformance_ReportInvalidRequirement(compiler, requirement, ("two requirements with the same rule_id: " + j2gort.StrValString(requirement.GetRuleId())))
		return false
	}
	return true
}

// Java static method: CheckConformance.validateExtendingRequirement(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.Requirement extendingRequirement, com.google.javascript.jscomp.Requirement.Builder extendableRequirement)
func CheckConformance_ValidateExtendingRequirement(compiler AbstractCompiler, extendingRequirement *j2gort.Ext_com_google_javascript_jscomp_Requirement, extendableRequirement *j2gort.Ext_com_google_javascript_jscomp_Requirement_Builder) bool {
	clinit_CheckConformance()
	_ = compiler
	_ = extendingRequirement
	_ = extendableRequirement
	if !(extendingRequirement.HasExtends()) {
		j2gort.CheckStateStr(false, "Extending requirement must have an extends")
	}
	if extendableRequirement == nil {
		CheckConformance_ReportInvalidRequirement(compiler, extendingRequirement, ("no extendable requirement with rule_id: " + j2gort.StrValString(extendingRequirement.GetExtends())))
		return false
	}
	{
		cur_27 := j2gort.NewMapKeyCursor(extendingRequirement.GetAllFields())
		for cur_27.HasNext() {
			rawElem_28 := cur_27.Next()
			field := j2gort.CastPtr[j2gort.Ext_com_google_protobuf_Descriptors_FieldDescriptor](rawElem_28)
			_ = field
			if !(j2gort.SetContainsStr(StaticGet_CheckConformance_EXTENDABLE_FIELDS(), field.GetName())) {
				CheckConformance_ReportInvalidRequirement(compiler, extendingRequirement, ("extending rules allow only " + j2gort.StrVal(StaticGet_CheckConformance_EXTENDABLE_FIELDS())))
			}
		}
	}
	if (extendingRequirement.GetValueCount() > int32(0)) && !(extendableRequirement.GetAllowExtendingValue()) {
		CheckConformance_ReportInvalidRequirement(compiler, extendingRequirement, "extending rule may not specify 'value' if extendable rule does not allow it")
	}
	return true
}

// Java static method: CheckConformance.removeDuplicates(com.google.javascript.jscomp.Requirement.Builder requirement)
func CheckConformance_RemoveDuplicates(requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement_Builder) {
	clinit_CheckConformance()
	_ = requirement
	var list1 j2gort.Ext_com_google_common_collect_ImmutableSet = j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](requirement.GetWhitelistList()))
	_ = list1
	requirement.ClearWhitelist().AddAllWhitelist(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](list1))
	var allowlist j2gort.Ext_com_google_common_collect_ImmutableSet = j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](requirement.GetAllowlistList()))
	_ = allowlist
	requirement.ClearAllowlist().AddAllAllowlist(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](allowlist))
	var list2 j2gort.Ext_com_google_common_collect_ImmutableSet = j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](requirement.GetWhitelistRegexpList()))
	_ = list2
	requirement.ClearWhitelistRegexp().AddAllWhitelistRegexp(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](list2))
	var allowlistRegexp j2gort.Ext_com_google_common_collect_ImmutableSet = j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](requirement.GetAllowlistRegexpList()))
	_ = allowlistRegexp
	requirement.ClearAllowlistRegexp().AddAllAllowlistRegexp(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](allowlistRegexp))
	var list3 j2gort.Ext_com_google_common_collect_ImmutableSet = j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](requirement.GetOnlyApplyToList()))
	_ = list3
	requirement.ClearOnlyApplyTo().AddAllOnlyApplyTo(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](list3))
	var list4 j2gort.Ext_com_google_common_collect_ImmutableSet = j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](requirement.GetOnlyApplyToRegexpList()))
	_ = list4
	requirement.ClearOnlyApplyToRegexp().AddAllOnlyApplyToRegexp(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](list4))
}

// Java static method: CheckConformance.initRule(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.Requirement requirement)
func CheckConformance_InitRule(compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) CheckConformance_Rule {
	clinit_CheckConformance()
	_ = compiler
	_ = requirement
	var didRet_29 bool
	var retVal_30 CheckConformance_Rule
	func() {
		defer func() {
			if rec_31 := recover(); rec_31 != nil {
				jex_32 := javalang.RecoverException(rec_31)
				_ = jex_32
				if j2gort.InstanceOf[*CheckConformance_InvalidRequirementSpec](jex_32) {
					e := j2gort.AsIface[*CheckConformance_InvalidRequirementSpec](jex_32)
					_ = e
					CheckConformance_ReportInvalidRequirement(compiler, requirement, j2gort.ThrowableGetMessage(e))
					retVal_30 = nil
					didRet_29 = true
					return
				} else {
					panic(rec_31)
				}
			}
		}()
		switch j2gort.CheckNotNullPtr(requirement.GetType()).Ordinal() {
		case int32(0) /* Type.CUSTOM */ :
			retVal_30 = NewConformanceRules_CustomRuleProxy(compiler, requirement)
			didRet_29 = true
			return
		case int32(1) /* Type.NO_OP */ :
			retVal_30 = NewConformanceRules_NoOp(compiler, requirement)
			didRet_29 = true
			return
		case int32(11) /* Type.BANNED_CODE_PATTERN */ :
			retVal_30 = NewConformanceRules_BannedCodePattern(compiler, requirement)
			didRet_29 = true
			return
		case int32(2) /* Type.BANNED_DEPENDENCY */ :
			retVal_30 = NewConformanceRules_BannedDependency(compiler, requirement)
			didRet_29 = true
			return
		case int32(3) /* Type.BANNED_DEPENDENCY_REGEX */ :
			retVal_30 = NewConformanceRules_BannedDependencyRegex(compiler, requirement)
			didRet_29 = true
			return
		case int32(4) /* Type.BANNED_ENHANCE */ :
			retVal_30 = NewConformanceRules_BannedEnhance(compiler, requirement)
			didRet_29 = true
			return
		case int32(17) /* Type.BANNED_MODS_REGEX */ :
			retVal_30 = NewConformanceRules_BannedModsRegex(compiler, requirement)
			didRet_29 = true
			return
		case int32(5) /* Type.BANNED_NAME */, int32(14) /* Type.BANNED_NAME_CALL */ :
			retVal_30 = NewConformanceRules_BannedName(compiler, requirement)
			didRet_29 = true
			return
		case int32(6) /* Type.BANNED_PROPERTY */, int32(7) /* Type.BANNED_PROPERTY_READ */, int32(8) /* Type.BANNED_PROPERTY_WRITE */, int32(13) /* Type.BANNED_PROPERTY_NON_CONSTANT_WRITE */, int32(12) /* Type.BANNED_PROPERTY_CALL */ :
			retVal_30 = NewConformanceRules_BannedProperty(compiler, requirement)
			didRet_29 = true
			return
		case int32(9) /* Type.RESTRICTED_NAME_CALL */ :
			retVal_30 = NewConformanceRules_RestrictedNameCall(compiler, requirement)
			didRet_29 = true
			return
		case int32(10) /* Type.RESTRICTED_METHOD_CALL */ :
			retVal_30 = NewConformanceRules_RestrictedMethodCall(compiler, requirement)
			didRet_29 = true
			return
		case int32(15) /* Type.RESTRICTED_PROPERTY_WRITE */ :
			retVal_30 = NewConformanceRules_RestrictedPropertyWrite(compiler, requirement)
			didRet_29 = true
			return
		case int32(16) /* Type.BANNED_STRING_REGEX */ :
			retVal_30 = NewConformanceRules_BannedStringRegex(compiler, requirement)
			didRet_29 = true
			return
		default:
			panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError", "Unexpected requirement type: "+j2gort.StrVal(requirement.GetType()))))
		}
	}()
	if didRet_29 {
		return retVal_30
	}
	panic("unreachable")
}

// Java static method: CheckConformance.reportInvalidRequirement(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.Requirement requirement, java.lang.String reason)
func CheckConformance_ReportInvalidRequirement(compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement, reason string) {
	clinit_CheckConformance()
	_ = compiler
	_ = requirement
	_ = reason
	compiler.Report_1(JSError_Make_2(StaticGet_CheckConformance_INVALID_REQUIREMENT_SPEC(), []string{reason, j2gort.Ext_com_google_protobuf_TextFormat_S_Printer().PrintToString_1_Ext_com_google_protobuf_MessageOrBuilder(j2gort.AsIface[j2gort.Ext_com_google_protobuf_MessageOrBuilder](requirement))}))
}

// Java interface: com.google.javascript.jscomp.CheckConformance.Rule
type CheckConformance_Rule interface {
	Implements_jscomp_CheckConformance_Rule()
	IsJavaNil_() bool
	GetPrecondition() CheckConformance_Precondition
	Check(*NodeTraversal, pkg_rhino.Node, *j2gort.Ext_com_google_javascript_jscomp_ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior)
}

func ToIface_CheckConformance_Rule[T interface {
	CheckConformance_Rule
	comparable
}](v T) CheckConformance_Rule {
	var z T
	if v == z {
		return nil
	}
	return v
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CheckConformance.Rule", "com.google.javascript.jscomp.CheckConformance$Rule", func(v_ any) bool {
		return j2gort.InstanceOf[CheckConformance_Rule](v_)
	})
}

// Java default interface method: Rule.getPrecondition()
func IfaceDefault_CheckConformance_Rule_GetPrecondition(this_ CheckConformance_Rule) CheckConformance_Precondition {
	_ = this_
	return StaticGet_CheckConformance_Precondition_CHECK_ALL()
}

// Java interface: com.google.javascript.jscomp.CheckConformance.Precondition
type CheckConformance_Precondition interface {
	Implements_jscomp_CheckConformance_Precondition()
	IsJavaNil_() bool
	ShouldCheck(pkg_rhino.Node) bool
}

func ToIface_CheckConformance_Precondition[T interface {
	CheckConformance_Precondition
	comparable
}](v T) CheckConformance_Precondition {
	var z T
	if v == z {
		return nil
	}
	return v
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CheckConformance.Precondition", "com.google.javascript.jscomp.CheckConformance$Precondition", func(v_ any) bool {
		return j2gort.InstanceOf[CheckConformance_Precondition](v_)
	})
}

var static_CheckConformance_Precondition_CHECK_ALL CheckConformance_Precondition

func StaticGet_CheckConformance_Precondition_CHECK_ALL() CheckConformance_Precondition {
	clinit_CheckConformance_Precondition()
	return static_CheckConformance_Precondition_CHECK_ALL
}
func StaticSet_CheckConformance_Precondition_CHECK_ALL(v CheckConformance_Precondition) {
	clinit_CheckConformance_Precondition()
	static_CheckConformance_Precondition_CHECK_ALL = v
}

var static_CheckConformance_Precondition_IS_CONSTRUCTOR_OR_CLASS CheckConformance_Precondition

func StaticGet_CheckConformance_Precondition_IS_CONSTRUCTOR_OR_CLASS() CheckConformance_Precondition {
	clinit_CheckConformance_Precondition()
	return static_CheckConformance_Precondition_IS_CONSTRUCTOR_OR_CLASS
}
func StaticSet_CheckConformance_Precondition_IS_CONSTRUCTOR_OR_CLASS(v CheckConformance_Precondition) {
	clinit_CheckConformance_Precondition()
	static_CheckConformance_Precondition_IS_CONSTRUCTOR_OR_CLASS = v
}

var clinitDone_CheckConformance_Precondition uint32
var clinitState_CheckConformance_Precondition j2gort.ClinitState

func clinit_CheckConformance_Precondition() {
	if atomic.LoadUint32(&clinitDone_CheckConformance_Precondition) != 0 {
		return
	}
	clinitSlow_CheckConformance_Precondition()
}

func clinitSlow_CheckConformance_Precondition() {
	if !j2gort.EnterClinit(&clinitState_CheckConformance_Precondition, &clinitDone_CheckConformance_Precondition) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CheckConformance_Precondition, &clinitDone_CheckConformance_Precondition)
	static_CheckConformance_Precondition_CHECK_ALL = NewCheckConformance__Anon_9()
	static_CheckConformance_Precondition_IS_CONSTRUCTOR_OR_CLASS = lambdaSingleton_check_conformance_2
}

type lambda_CheckConformance_Precondition struct {
	fn func(pkg_rhino.Node) bool
}

func (l *lambda_CheckConformance_Precondition) IsJavaNil_() bool { return l == nil }
func NewLambda_CheckConformance_Precondition(fn func(pkg_rhino.Node) bool) CheckConformance_Precondition {
	if fn == nil {
		return nil
	}
	return &lambda_CheckConformance_Precondition{fn: fn}
}

func (*lambda_CheckConformance_Precondition) Implements_jscomp_CheckConformance_Precondition() {}
func (l *lambda_CheckConformance_Precondition) ShouldCheck(p0_ pkg_rhino.Node) bool {
	return l.fn(p0_)
}

// Java class: com.google.javascript.jscomp.CheckConformance.Category
type CheckConformance_Category struct {
	precondition CheckConformance_Precondition `gson:"precondition"`
	rules        []CheckConformance_Rule       `gson:"rules"`
	idHash_      atomic.Int32
}

func (this_ *CheckConformance_Category) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckConformance_Category) IsJavaNil_() bool { return this_ == nil }
func (*CheckConformance_Category) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CheckConformance$Category", "Category"
}
func (this_ *CheckConformance_Category) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CheckConformance_Category) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CheckConformance.Category", "com.google.javascript.jscomp.CheckConformance$Category", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckConformance_Category](v_)
	})
}

func (this_ *CheckConformance_Category) initFields_CheckConformance_Category() {
	_ = this_
}

// Java constructor: Category(com.google.javascript.jscomp.CheckConformance.Precondition precondition, com.google.common.collect.ImmutableList<com.google.javascript.jscomp.CheckConformance.Rule> rules)
func NewCheckConformance_Category(precondition CheckConformance_Precondition, rules j2gort.Ext_com_google_common_collect_ImmutableList) *CheckConformance_Category {
	this_ := &CheckConformance_Category{}
	InitCheckConformance_Category(this_, precondition, rules)
	return this_
}

func InitCheckConformance_Category(this_ *CheckConformance_Category, precondition CheckConformance_Precondition, rules j2gort.Ext_com_google_common_collect_ImmutableList) {
	_ = precondition
	_ = rules
	this_.initFields_CheckConformance_Category()
	this_.precondition = precondition
	this_.rules = j2gort.ImmListToSliceNullable[

	// Java class: com.google.javascript.jscomp.CheckConformance.InvalidRequirementSpec
	CheckConformance_Rule](rules)
}

type CheckConformance_InvalidRequirementSpec struct {
	javalang.JavaException
	idHash_ atomic.Int32
}

func (this_ *CheckConformance_InvalidRequirementSpec) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckConformance_InvalidRequirementSpec) IsJavaNil_() bool { return this_ == nil }
func (*CheckConformance_InvalidRequirementSpec) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CheckConformance$InvalidRequirementSpec", "InvalidRequirementSpec"
}
func (this_ *CheckConformance_InvalidRequirementSpec) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CheckConformance.InvalidRequirementSpec", "com.google.javascript.jscomp.CheckConformance$InvalidRequirementSpec", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckConformance_InvalidRequirementSpec](v_)
	})
}

func (this_ *CheckConformance_InvalidRequirementSpec) initFields_CheckConformance_InvalidRequirementSpec() {
	_ = this_
}

// Java constructor: InvalidRequirementSpec(java.lang.String message)
func NewCheckConformance_InvalidRequirementSpec_1(message string) *CheckConformance_InvalidRequirementSpec {
	this_ := &CheckConformance_InvalidRequirementSpec{}
	InitCheckConformance_InvalidRequirementSpec_1(this_, message)
	return this_
}

func InitCheckConformance_InvalidRequirementSpec_1(this_ *CheckConformance_InvalidRequirementSpec, message string) {
	_ = message
	j2gort.InitUserThrowable(this_, "com.google.javascript.jscomp.CheckConformance.InvalidRequirementSpec", j2gort.BoxString(message))
	this_.initFields_CheckConformance_InvalidRequirementSpec()
}

// Java constructor: InvalidRequirementSpec(java.lang.String message, java.lang.Throwable cause)
func NewCheckConformance_InvalidRequirementSpec_2(message string, cause javalang.Throwable) *CheckConformance_InvalidRequirementSpec {
	this_ := &CheckConformance_InvalidRequirementSpec{}
	InitCheckConformance_InvalidRequirementSpec_2(this_, message, cause)
	return this_
}

func InitCheckConformance_InvalidRequirementSpec_2(this_ *CheckConformance_InvalidRequirementSpec, message string, cause javalang.Throwable) {
	_ = message
	_ = cause
	j2gort.InitUserThrowable(this_, "com.google.javascript.jscomp.CheckConformance.InvalidRequirementSpec", j2gort.BoxString(message), cause)
	this_.initFields_CheckConformance_InvalidRequirementSpec()
}

// Java class:
type CheckConformance__Anon_9 struct {
	idHash_ atomic.Int32
}

func (this_ *CheckConformance__Anon_9) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckConformance__Anon_9) IsJavaNil_() bool { return this_ == nil }
func (*CheckConformance__Anon_9) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CheckConformance$Precondition$1", ""
}
func (this_ *CheckConformance__Anon_9) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CheckConformance__Anon_9) HasIdentityEquals_()                              {}
func (*CheckConformance__Anon_9) Implements_jscomp_CheckConformance_Precondition() {}

func (this_ *CheckConformance__Anon_9) initFields_CheckConformance__Anon_9() {
	_ = this_
}

func NewCheckConformance__Anon_9() *CheckConformance__Anon_9 {
	this_ := &CheckConformance__Anon_9{}
	this_.initFields_CheckConformance__Anon_9()
	return this_
}

// Java method: .shouldCheck(com.google.javascript.rhino.Node n)
func (this_ *CheckConformance__Anon_9) ShouldCheck(n pkg_rhino.Node) bool {
	_ = n
	return true
}
