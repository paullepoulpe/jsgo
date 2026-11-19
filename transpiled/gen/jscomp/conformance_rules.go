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
// Port of com/google/javascript/jscomp/ConformanceRules.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parsing "github.com/paullepoulpe/jsgo/transpiled/gen/parsing"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_conformance_rules_1 CheckConformance_Precondition
var lambdaSingleton_conformance_rules_2 ConformanceRules_CustomRuleProxy_CustomRuleFactory
var lambdaSingleton_conformance_rules_3 ConformanceRules_CustomRuleProxy_CustomRuleFactory
var lambdaSingleton_conformance_rules_4 ConformanceRules_CustomRuleProxy_CustomRuleFactory
var lambdaSingleton_conformance_rules_5 ConformanceRules_CustomRuleProxy_CustomRuleFactory
var lambdaSingleton_conformance_rules_6 ConformanceRules_CustomRuleProxy_CustomRuleFactory
var lambdaSingleton_conformance_rules_7 ConformanceRules_CustomRuleProxy_CustomRuleFactory
var lambdaSingleton_conformance_rules_8 ConformanceRules_CustomRuleProxy_CustomRuleFactory
var lambdaSingleton_conformance_rules_9 ConformanceRules_CustomRuleProxy_CustomRuleFactory
var lambdaSingleton_conformance_rules_10 ConformanceRules_CustomRuleProxy_CustomRuleFactory
var lambdaSingleton_conformance_rules_11 ConformanceRules_CustomRuleProxy_CustomRuleFactory
var lambdaSingleton_conformance_rules_12 ConformanceRules_CustomRuleProxy_CustomRuleFactory
var lambdaSingleton_conformance_rules_13 ConformanceRules_CustomRuleProxy_CustomRuleFactory
var lambdaSingleton_conformance_rules_14 ConformanceRules_CustomRuleProxy_CustomRuleFactory
var lambdaSingleton_conformance_rules_15 ConformanceRules_CustomRuleProxy_CustomRuleFactory
var lambdaSingleton_conformance_rules_16 ConformanceRules_CustomRuleProxy_CustomRuleFactory
var lambdaSingleton_conformance_rules_17 ConformanceRules_CustomRuleProxy_CustomRuleFactory
var lambdaSingleton_conformance_rules_18 ConformanceRules_CustomRuleProxy_CustomRuleFactory
var lambdaSingleton_conformance_rules_19 ConformanceRules_CustomRuleProxy_CustomRuleFactory
var lambdaSingleton_conformance_rules_20 ConformanceRules_CustomRuleProxy_CustomRuleFactory
var lambdaSingleton_conformance_rules_21 ConformanceRules_CustomRuleProxy_CustomRuleFactory
var lambdaSingleton_conformance_rules_22 ConformanceRules_CustomRuleProxy_CustomRuleFactory
var lambdaSingleton_conformance_rules_23 ConformanceRules_CustomRuleProxy_CustomRuleFactory
var lambdaSingleton_conformance_rules_24 ConformanceRules_CustomRuleProxy_CustomRuleFactory
var lambdaSingleton_conformance_rules_25 ConformanceRules_CustomRuleProxy_CustomRuleFactory
var lambdaSingleton_conformance_rules_26 ConformanceRules_CustomRuleProxy_CustomRuleFactory
var lambdaSingleton_conformance_rules_27 ConformanceRules_CustomRuleProxy_CustomRuleFactory
var lambdaSingleton_conformance_rules_28 ConformanceRules_CustomRuleProxy_CustomRuleFactory
var lambdaSingleton_conformance_rules_29 ConformanceRules_CustomRuleProxy_CustomRuleFactory
var lambdaSingleton_conformance_rules_30 ConformanceRules_CustomRuleProxy_CustomRuleFactory
var lambdaSingleton_conformance_rules_31 CheckConformance_Precondition
var lambdaSingleton_conformance_rules_32 CheckConformance_Precondition
var lambdaSingleton_conformance_rules_33 j2gort.Ext_java_util_function_Supplier
var lambdaSingleton_conformance_rules_34 j2gort.Ext_java_util_function_Supplier
var lambdaSingleton_conformance_rules_35 j2gort.Ext_java_util_function_Function
var lambdaSingleton_conformance_rules_36 j2gort.Ext_java_util_function_Function

func init() {
	lambdaSingleton_conformance_rules_1 = NewLambda_CheckConformance_Precondition(func(n pkg_rhino.Node) bool {
		_ = n
		return func() bool {
			switch n.TokenSwitchOrd_() {
			case int32(26) /* Token.GETPROP */ :
				return n.GetFirstChild().IsQualifiedName()
			case int32(32) /* Token.NAME */ :
				return !(j2gort.CheckNotNullString(n.GetString()) == "")
			default:
				return false
			}
			panic("unreachable switch expression")
		}()
	})
	lambdaSingleton_conformance_rules_2 = NewLambda_ConformanceRules_CustomRuleProxy_CustomRuleFactory(func(mrp0_ AbstractCompiler, mrp1_ *j2gort.Ext_com_google_javascript_jscomp_Requirement) CheckConformance_Rule {
		return NewConformanceRules_BanStaticThis(mrp0_, mrp1_)
	})
	lambdaSingleton_conformance_rules_3 = NewLambda_ConformanceRules_CustomRuleProxy_CustomRuleFactory(func(mrp0_ AbstractCompiler, mrp1_ *j2gort.Ext_com_google_javascript_jscomp_Requirement) CheckConformance_Rule {
		return NewConformanceRules_BanExecCommand(mrp0_, mrp1_)
	})
	lambdaSingleton_conformance_rules_4 = NewLambda_ConformanceRules_CustomRuleProxy_CustomRuleFactory(func(mrp0_ AbstractCompiler, mrp1_ *j2gort.Ext_com_google_javascript_jscomp_Requirement) CheckConformance_Rule {
		return NewConformanceRules_BanSettingAttributes(mrp0_, mrp1_)
	})
	lambdaSingleton_conformance_rules_5 = NewLambda_ConformanceRules_CustomRuleProxy_CustomRuleFactory(func(mrp0_ AbstractCompiler, mrp1_ *j2gort.Ext_com_google_javascript_jscomp_Requirement) CheckConformance_Rule {
		return NewConformanceRules_BanElementSetAttribute(mrp0_, mrp1_)
	})
	lambdaSingleton_conformance_rules_6 = NewLambda_ConformanceRules_CustomRuleProxy_CustomRuleFactory(func(mrp0_ AbstractCompiler, mrp1_ *j2gort.Ext_com_google_javascript_jscomp_Requirement) CheckConformance_Rule {
		return NewConformanceRules_BanCreateDom(mrp0_, mrp1_)
	})
	lambdaSingleton_conformance_rules_7 = NewLambda_ConformanceRules_CustomRuleProxy_CustomRuleFactory(func(mrp0_ AbstractCompiler, mrp1_ *j2gort.Ext_com_google_javascript_jscomp_Requirement) CheckConformance_Rule {
		return NewConformanceRules_BanCreateElement(mrp0_, mrp1_)
	})
	lambdaSingleton_conformance_rules_8 = NewLambda_ConformanceRules_CustomRuleProxy_CustomRuleFactory(func(mrp0_ AbstractCompiler, mrp1_ *j2gort.Ext_com_google_javascript_jscomp_Requirement) CheckConformance_Rule {
		return NewConformanceRules_BannedModsRegex(mrp0_, mrp1_)
	})
	lambdaSingleton_conformance_rules_9 = NewLambda_ConformanceRules_CustomRuleProxy_CustomRuleFactory(func(mrp0_ AbstractCompiler, mrp1_ *j2gort.Ext_com_google_javascript_jscomp_Requirement) CheckConformance_Rule {
		return NewConformanceRules_BannedEnhance(mrp0_, mrp1_)
	})
	lambdaSingleton_conformance_rules_10 = NewLambda_ConformanceRules_CustomRuleProxy_CustomRuleFactory(func(mrp0_ AbstractCompiler, mrp1_ *j2gort.Ext_com_google_javascript_jscomp_Requirement) CheckConformance_Rule {
		return NewConformanceRules_BanGlobalVars(mrp0_, mrp1_)
	})
	lambdaSingleton_conformance_rules_11 = NewLambda_ConformanceRules_CustomRuleProxy_CustomRuleFactory(func(mrp0_ AbstractCompiler, mrp1_ *j2gort.Ext_com_google_javascript_jscomp_Requirement) CheckConformance_Rule {
		return NewConformanceRules_StrictBanUnresolvedType(mrp0_, mrp1_)
	})
	lambdaSingleton_conformance_rules_12 = NewLambda_ConformanceRules_CustomRuleProxy_CustomRuleFactory(func(mrp0_ AbstractCompiler, mrp1_ *j2gort.Ext_com_google_javascript_jscomp_Requirement) CheckConformance_Rule {
		return NewConformanceRules_BanUnresolvedType(mrp0_, mrp1_)
	})
	lambdaSingleton_conformance_rules_13 = NewLambda_ConformanceRules_CustomRuleProxy_CustomRuleFactory(func(mrp0_ AbstractCompiler, mrp1_ *j2gort.Ext_com_google_javascript_jscomp_Requirement) CheckConformance_Rule {
		return NewConformanceRules_BanUnknownTypedClassPropsReferences(mrp0_, mrp1_)
	})
	lambdaSingleton_conformance_rules_14 = NewLambda_ConformanceRules_CustomRuleProxy_CustomRuleFactory(func(mrp0_ AbstractCompiler, mrp1_ *j2gort.Ext_com_google_javascript_jscomp_Requirement) CheckConformance_Rule {
		return NewConformanceRules_BanUnknownDirectThisPropsReferences(mrp0_, mrp1_)
	})
	lambdaSingleton_conformance_rules_15 = NewLambda_ConformanceRules_CustomRuleProxy_CustomRuleFactory(func(mrp0_ AbstractCompiler, mrp1_ *j2gort.Ext_com_google_javascript_jscomp_Requirement) CheckConformance_Rule {
		return NewConformanceRules_BanUnknownThis(mrp0_, mrp1_)
	})
	lambdaSingleton_conformance_rules_16 = NewLambda_ConformanceRules_CustomRuleProxy_CustomRuleFactory(func(mrp0_ AbstractCompiler, mrp1_ *j2gort.Ext_com_google_javascript_jscomp_Requirement) CheckConformance_Rule {
		return NewConformanceRules_BanNullDeref(mrp0_, mrp1_)
	})
	lambdaSingleton_conformance_rules_17 = NewLambda_ConformanceRules_CustomRuleProxy_CustomRuleFactory(func(mrp0_ AbstractCompiler, mrp1_ *j2gort.Ext_com_google_javascript_jscomp_Requirement) CheckConformance_Rule {
		return NewConformanceRules_BanThrowOfNonErrorTypes(mrp0_, mrp1_)
	})
	lambdaSingleton_conformance_rules_18 = NewLambda_ConformanceRules_CustomRuleProxy_CustomRuleFactory(func(mrp0_ AbstractCompiler, mrp1_ *j2gort.Ext_com_google_javascript_jscomp_Requirement) CheckConformance_Rule {
		return NewConformanceRules_RequireUseStrict(mrp0_, mrp1_)
	})
	lambdaSingleton_conformance_rules_19 = NewLambda_ConformanceRules_CustomRuleProxy_CustomRuleFactory(func(mrp0_ AbstractCompiler, mrp1_ *j2gort.Ext_com_google_javascript_jscomp_Requirement) CheckConformance_Rule {
		return NewConformanceRules_BanForOf(mrp0_, mrp1_)
	})
	lambdaSingleton_conformance_rules_20 = NewLambda_ConformanceRules_CustomRuleProxy_CustomRuleFactory(func(mrp0_ AbstractCompiler, mrp1_ *j2gort.Ext_com_google_javascript_jscomp_Requirement) CheckConformance_Rule {
		return NewConformanceRules_BannedCodePattern(mrp0_, mrp1_)
	})
	lambdaSingleton_conformance_rules_21 = NewLambda_ConformanceRules_CustomRuleProxy_CustomRuleFactory(func(mrp0_ AbstractCompiler, mrp1_ *j2gort.Ext_com_google_javascript_jscomp_Requirement) CheckConformance_Rule {
		return NewConformanceRules_RestrictedPropertyWrite(mrp0_, mrp1_)
	})
	lambdaSingleton_conformance_rules_22 = NewLambda_ConformanceRules_CustomRuleProxy_CustomRuleFactory(func(mrp0_ AbstractCompiler, mrp1_ *j2gort.Ext_com_google_javascript_jscomp_Requirement) CheckConformance_Rule {
		return NewConformanceRules_RestrictedMethodCall(mrp0_, mrp1_)
	})
	lambdaSingleton_conformance_rules_23 = NewLambda_ConformanceRules_CustomRuleProxy_CustomRuleFactory(func(mrp0_ AbstractCompiler, mrp1_ *j2gort.Ext_com_google_javascript_jscomp_Requirement) CheckConformance_Rule {
		return NewConformanceRules_RestrictedNameCall(mrp0_, mrp1_)
	})
	lambdaSingleton_conformance_rules_24 = NewLambda_ConformanceRules_CustomRuleProxy_CustomRuleFactory(func(mrp0_ AbstractCompiler, mrp1_ *j2gort.Ext_com_google_javascript_jscomp_Requirement) CheckConformance_Rule {
		return NewConformanceRules_BannedStringRegex(mrp0_, mrp1_)
	})
	lambdaSingleton_conformance_rules_25 = NewLambda_ConformanceRules_CustomRuleProxy_CustomRuleFactory(func(mrp0_ AbstractCompiler, mrp1_ *j2gort.Ext_com_google_javascript_jscomp_Requirement) CheckConformance_Rule {
		return NewConformanceRules_BannedProperty(mrp0_, mrp1_)
	})
	lambdaSingleton_conformance_rules_26 = NewLambda_ConformanceRules_CustomRuleProxy_CustomRuleFactory(func(mrp0_ AbstractCompiler, mrp1_ *j2gort.Ext_com_google_javascript_jscomp_Requirement) CheckConformance_Rule {
		return NewConformanceRules_BannedName(mrp0_, mrp1_)
	})
	lambdaSingleton_conformance_rules_27 = NewLambda_ConformanceRules_CustomRuleProxy_CustomRuleFactory(func(mrp0_ AbstractCompiler, mrp1_ *j2gort.Ext_com_google_javascript_jscomp_Requirement) CheckConformance_Rule {
		return NewConformanceRules_BannedDependencyRegex(mrp0_, mrp1_)
	})
	lambdaSingleton_conformance_rules_28 = NewLambda_ConformanceRules_CustomRuleProxy_CustomRuleFactory(func(mrp0_ AbstractCompiler, mrp1_ *j2gort.Ext_com_google_javascript_jscomp_Requirement) CheckConformance_Rule {
		return NewConformanceRules_BannedDependency(mrp0_, mrp1_)
	})
	lambdaSingleton_conformance_rules_29 = NewLambda_ConformanceRules_CustomRuleProxy_CustomRuleFactory(func(mrp0_ AbstractCompiler, mrp1_ *j2gort.Ext_com_google_javascript_jscomp_Requirement) CheckConformance_Rule {
		return NewConformanceRules_InferredConstCheck(mrp0_, mrp1_)
	})
	lambdaSingleton_conformance_rules_30 = NewLambda_ConformanceRules_CustomRuleProxy_CustomRuleFactory(func(mrp0_ AbstractCompiler, mrp1_ *j2gort.Ext_com_google_javascript_jscomp_Requirement) CheckConformance_Rule {
		return NewConformanceRules_NoOp(mrp0_, mrp1_)
	})
	lambdaSingleton_conformance_rules_31 = NewLambda_CheckConformance_Precondition(func(mrp0_ pkg_rhino.Node) bool { return mrp0_.IsScript() })
	lambdaSingleton_conformance_rules_32 = NewLambda_CheckConformance_Precondition(func(mrp0_ pkg_rhino.Node) bool { return mrp0_.IsScript() })
	lambdaSingleton_conformance_rules_33 = j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
		return nil
	})
	lambdaSingleton_conformance_rules_34 = j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
		return nil
	})
	lambdaSingleton_conformance_rules_35 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_49 any) any {
		s := j2gort.CastString(lp_49)
		_ = s
		return j2gort.BoxString(javalang.StringToLowerCaseASCII(s))
	})
	lambdaSingleton_conformance_rules_36 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_61 any) any {
		v := j2gort.CastString(lp_61)
		_ = v
		return j2gort.BoxString(javalang.StringToLowerCaseASCII(v))
	})
}

// Java class: com.google.javascript.jscomp.ConformanceRules
type ConformanceRules struct {
	idHash_ atomic.Int32
}

func (this_ *ConformanceRules) BindSelf_(s any) {
	_ = s
}
func (this_ *ConformanceRules) IsJavaNil_() bool { return this_ == nil }
func (*ConformanceRules) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConformanceRules", "ConformanceRules"
}
func (this_ *ConformanceRules) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ConformanceRules) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConformanceRules", "com.google.javascript.jscomp.ConformanceRules", func(v_ any) bool {
		return j2gort.InstanceOf[*ConformanceRules](v_)
	})
}

var static_ConformanceRules_ALL_TS_ALLOWLIST *ConformanceRules_AllowList

func StaticGet_ConformanceRules_ALL_TS_ALLOWLIST() *ConformanceRules_AllowList {
	clinit_ConformanceRules()
	return static_ConformanceRules_ALL_TS_ALLOWLIST
}
func StaticSet_ConformanceRules_ALL_TS_ALLOWLIST(v *ConformanceRules_AllowList) {
	clinit_ConformanceRules()
	static_ConformanceRules_ALL_TS_ALLOWLIST = v
}

var static_ConformanceRules_ON_PROTOTYPE *j2gort.Ext_com_google_common_base_Splitter

func StaticGet_ConformanceRules_ON_PROTOTYPE() *j2gort.Ext_com_google_common_base_Splitter {
	clinit_ConformanceRules()
	return static_ConformanceRules_ON_PROTOTYPE
}
func StaticSet_ConformanceRules_ON_PROTOTYPE(v *j2gort.Ext_com_google_common_base_Splitter) {
	clinit_ConformanceRules()
	static_ConformanceRules_ON_PROTOTYPE = v
}

var static_ConformanceRules_ON_DOT *j2gort.Ext_com_google_common_base_Splitter

func StaticGet_ConformanceRules_ON_DOT() *j2gort.Ext_com_google_common_base_Splitter {
	clinit_ConformanceRules()
	return static_ConformanceRules_ON_DOT
}
func StaticSet_ConformanceRules_ON_DOT(v *j2gort.Ext_com_google_common_base_Splitter) {
	clinit_ConformanceRules()
	static_ConformanceRules_ON_DOT = v
}

var static_ConformanceRules_GOOG_DOM pkg_rhino.QualifiedName

func StaticGet_ConformanceRules_GOOG_DOM() pkg_rhino.QualifiedName {
	clinit_ConformanceRules()
	return static_ConformanceRules_GOOG_DOM
}
func StaticSet_ConformanceRules_GOOG_DOM(v pkg_rhino.QualifiedName) {
	clinit_ConformanceRules()
	static_ConformanceRules_GOOG_DOM = v
}

var static_ConformanceRules_GOOG_DOM_TAGNAME pkg_rhino.QualifiedName

func StaticGet_ConformanceRules_GOOG_DOM_TAGNAME() pkg_rhino.QualifiedName {
	clinit_ConformanceRules()
	return static_ConformanceRules_GOOG_DOM_TAGNAME
}
func StaticSet_ConformanceRules_GOOG_DOM_TAGNAME(v pkg_rhino.QualifiedName) {
	clinit_ConformanceRules()
	static_ConformanceRules_GOOG_DOM_TAGNAME = v
}

var clinitDone_ConformanceRules uint32
var clinitState_ConformanceRules j2gort.ClinitState

func clinit_ConformanceRules() {
	if atomic.LoadUint32(&clinitDone_ConformanceRules) != 0 {
		return
	}
	clinitSlow_ConformanceRules()
}

func clinitSlow_ConformanceRules() {
	if !j2gort.EnterClinit(&clinitState_ConformanceRules, &clinitDone_ConformanceRules) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ConformanceRules, &clinitDone_ConformanceRules)
	static_ConformanceRules_ALL_TS_ALLOWLIST = ConformanceRules_CreateTsAllowlist()
	static_ConformanceRules_ON_PROTOTYPE = j2gort.Ext_com_google_common_base_Splitter_On_1_string(".prototype.")
	static_ConformanceRules_ON_DOT = j2gort.Ext_com_google_common_base_Splitter_On_1_string(".")
	static_ConformanceRules_GOOG_DOM = pkg_rhino.QualifiedName_Of("goog.dom")
	static_ConformanceRules_GOOG_DOM_TAGNAME = pkg_rhino.QualifiedName_Of("goog.dom.TagName")
}

func (this_ *ConformanceRules) initFields_ConformanceRules() {
	_ = this_
}

// Java constructor: ConformanceRules()
func NewConformanceRules() *ConformanceRules {
	clinit_ConformanceRules()
	this_ := &ConformanceRules{}
	InitConformanceRules(this_)
	return this_
}

func InitConformanceRules(this_ *ConformanceRules) {
	this_.initFields_ConformanceRules()
}

// Java static method: ConformanceRules.createTsAllowlist()
func ConformanceRules_CreateTsAllowlist() *ConformanceRules_AllowList {
	clinit_ConformanceRules()
	var didRet_1 bool
	var retVal_2 *ConformanceRules_AllowList
	func() {
		defer func() {
			if rec_3 := recover(); rec_3 != nil {
				jex_4 := javalang.RecoverException(rec_3)
				_ = jex_4
				if j2gort.InstanceOf[javalang.Throwable](jex_4) {
					t := j2gort.AsIface[javalang.Throwable](jex_4)
					_ = t
					panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError", t)))
				} else {
					panic(rec_3)
				}
			}
		}()
		retVal_2 = NewConformanceRules_AllowList_2(j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()), j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.ImmListOfStr2(".*\\.closure\\.js", ".*\\.tsx\\.cl\\.js")))
		didRet_1 = true
		return
	}()
	if didRet_1 {
		return retVal_2
	}
	panic("unreachable")
}

// Java static method: ConformanceRules.buildPattern(java.util.List<java.lang.String> reqPatterns)
func ConformanceRules_BuildPattern(reqPatterns j2gort.Ext_java_util_List) *j2gort.Ext_java_util_regex_Pattern {
	clinit_ConformanceRules()
	_ = reqPatterns
	if (reqPatterns == nil) || reqPatterns.IsEmpty() {
		return nil
	}
	{
		cur_5 := j2gort.NewIterCursor(j2gort.CheckNotNull(reqPatterns))
		for cur_5.HasNext() {
			reqPattern := cur_5.NextStr()
			_ = reqPattern
			func() {
				defer func() {
					if rec_6 := recover(); rec_6 != nil {
						jex_7 := javalang.RecoverException(rec_6)
						_ = jex_7
						if j2gort.InstanceOf[j2gort.Ext_java_util_regex_PatternSyntaxException](jex_7) {
							e := j2gort.AsIface[j2gort.Ext_java_util_regex_PatternSyntaxException](jex_7)
							_ = e
							panic(NewCheckConformance_InvalidRequirementSpec_2("invalid regex pattern", j2gort.AsIface[javalang.Throwable](e)))
						} else {
							panic(rec_6)
						}
					}
				}()
				j2gort.Ext_java_util_regex_Pattern_Compile_1(reqPattern)
			}()
		}
	}
	var pattern *j2gort.Ext_java_util_regex_Pattern = nil
	_ = pattern
	func() {
		defer func() {
			if rec_8 := recover(); rec_8 != nil {
				jex_9 := javalang.RecoverException(rec_8)
				_ = jex_9
				if j2gort.InstanceOf[j2gort.Ext_java_util_regex_PatternSyntaxException](jex_9) {
					e_2 := j2gort.AsIface[j2gort.Ext_java_util_regex_PatternSyntaxException](jex_9)
					_ = e_2
					panic(j2gort.AsIface[javalang.IRuntimeException](j2gort.ConstructShimThrowable("java.lang.RuntimeException", "bad joined regexp", e_2)))
				} else {
					panic(rec_8)
				}
			}
		}()
		var jointRegExp string = ("(" + j2gort.StrValString(j2gort.Ext_com_google_common_base_Joiner_On_1_string("|").Join_1_Ext_java_lang_Iterable(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](reqPatterns))) + ")")
		_ = jointRegExp
		pattern = j2gort.Ext_java_util_regex_Pattern_Compile_1(jointRegExp)
	}()
	return pattern
}

// Java static method: ConformanceRules.isTaggedTemplateLitString(com.google.javascript.rhino.Node node)
func ConformanceRules_IsTaggedTemplateLitString(node pkg_rhino.Node) bool {
	clinit_ConformanceRules()
	_ = node
	if !(node.IsTemplateLitString()) {
		j2gort.CheckState(false)
	}
	return ((node.GetGrandparent() != 0) && node.GetGrandparent().IsTaggedTemplateLit())
}

// Java static method: ConformanceRules.isLooseObject(com.google.javascript.rhino.jstype.JSType type, com.google.javascript.rhino.jstype.JSTypeRegistry registry)
func ConformanceRules_IsLooseObject(type_ pkg_rhino.JSType, registry *pkg_rhino.JSTypeRegistry) bool {
	clinit_ConformanceRules()
	_ = type_
	_ = registry
	return type_.Equals(registry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_OBJECT_TYPE()))
}

// Java class: com.google.javascript.jscomp.ConformanceRules.ConformanceResult
type ConformanceRules_ConformanceResult struct {
	note    string `gson:"note"`
	idHash_ atomic.Int32
	level   ConformanceRules_ConformanceLevel `gson:"level"`
}

func (this_ *ConformanceRules_ConformanceResult) BindSelf_(s any) {
	_ = s
}
func (this_ *ConformanceRules_ConformanceResult) IsJavaNil_() bool { return this_ == nil }
func (*ConformanceRules_ConformanceResult) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConformanceRules$ConformanceResult", "ConformanceResult"
}
func (this_ *ConformanceRules_ConformanceResult) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ConformanceRules_ConformanceResult) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConformanceRules.ConformanceResult", "com.google.javascript.jscomp.ConformanceRules$ConformanceResult", func(v_ any) bool {
		return j2gort.InstanceOf[*ConformanceRules_ConformanceResult](v_)
	})
}

var static_ConformanceRules_ConformanceResult_CONFORMANCE *ConformanceRules_ConformanceResult

func StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE() *ConformanceRules_ConformanceResult {
	clinit_ConformanceRules_ConformanceResult()
	return static_ConformanceRules_ConformanceResult_CONFORMANCE
}
func StaticSet_ConformanceRules_ConformanceResult_CONFORMANCE(v *ConformanceRules_ConformanceResult) {
	clinit_ConformanceRules_ConformanceResult()
	static_ConformanceRules_ConformanceResult_CONFORMANCE = v
}

var static_ConformanceRules_ConformanceResult_POSSIBLE_VIOLATION *ConformanceRules_ConformanceResult

func StaticGet_ConformanceRules_ConformanceResult_POSSIBLE_VIOLATION() *ConformanceRules_ConformanceResult {
	clinit_ConformanceRules_ConformanceResult()
	return static_ConformanceRules_ConformanceResult_POSSIBLE_VIOLATION
}
func StaticSet_ConformanceRules_ConformanceResult_POSSIBLE_VIOLATION(v *ConformanceRules_ConformanceResult) {
	clinit_ConformanceRules_ConformanceResult()
	static_ConformanceRules_ConformanceResult_POSSIBLE_VIOLATION = v
}

var static_ConformanceRules_ConformanceResult_POSSIBLE_VIOLATION_DUE_TO_LOOSE_TYPES *ConformanceRules_ConformanceResult

func StaticGet_ConformanceRules_ConformanceResult_POSSIBLE_VIOLATION_DUE_TO_LOOSE_TYPES() *ConformanceRules_ConformanceResult {
	clinit_ConformanceRules_ConformanceResult()
	return static_ConformanceRules_ConformanceResult_POSSIBLE_VIOLATION_DUE_TO_LOOSE_TYPES
}
func StaticSet_ConformanceRules_ConformanceResult_POSSIBLE_VIOLATION_DUE_TO_LOOSE_TYPES(v *ConformanceRules_ConformanceResult) {
	clinit_ConformanceRules_ConformanceResult()
	static_ConformanceRules_ConformanceResult_POSSIBLE_VIOLATION_DUE_TO_LOOSE_TYPES = v
}

var static_ConformanceRules_ConformanceResult_VIOLATION *ConformanceRules_ConformanceResult

func StaticGet_ConformanceRules_ConformanceResult_VIOLATION() *ConformanceRules_ConformanceResult {
	clinit_ConformanceRules_ConformanceResult()
	return static_ConformanceRules_ConformanceResult_VIOLATION
}
func StaticSet_ConformanceRules_ConformanceResult_VIOLATION(v *ConformanceRules_ConformanceResult) {
	clinit_ConformanceRules_ConformanceResult()
	static_ConformanceRules_ConformanceResult_VIOLATION = v
}

var clinitDone_ConformanceRules_ConformanceResult uint32
var clinitState_ConformanceRules_ConformanceResult j2gort.ClinitState

func clinit_ConformanceRules_ConformanceResult() {
	if atomic.LoadUint32(&clinitDone_ConformanceRules_ConformanceResult) != 0 {
		return
	}
	clinitSlow_ConformanceRules_ConformanceResult()
}

func clinitSlow_ConformanceRules_ConformanceResult() {
	if !j2gort.EnterClinit(&clinitState_ConformanceRules_ConformanceResult, &clinitDone_ConformanceRules_ConformanceResult) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ConformanceRules_ConformanceResult, &clinitDone_ConformanceRules_ConformanceResult)
	static_ConformanceRules_ConformanceResult_CONFORMANCE = NewConformanceRules_ConformanceResult_1(StaticGet_ConformanceRules_ConformanceLevel_CONFORMANCE())
	static_ConformanceRules_ConformanceResult_POSSIBLE_VIOLATION = NewConformanceRules_ConformanceResult_1(StaticGet_ConformanceRules_ConformanceLevel_POSSIBLE_VIOLATION())
	static_ConformanceRules_ConformanceResult_POSSIBLE_VIOLATION_DUE_TO_LOOSE_TYPES = NewConformanceRules_ConformanceResult_2(StaticGet_ConformanceRules_ConformanceLevel_POSSIBLE_VIOLATION(), "The type information available for this expression is too loose to ensure conformance.")
	static_ConformanceRules_ConformanceResult_VIOLATION = NewConformanceRules_ConformanceResult_1(StaticGet_ConformanceRules_ConformanceLevel_VIOLATION())
}

func (this_ *ConformanceRules_ConformanceResult) initFields_ConformanceRules_ConformanceResult() {
	_ = this_
	this_.note = j2gort.NullString
}

// Java constructor: ConformanceResult(com.google.javascript.jscomp.ConformanceRules.ConformanceLevel level)
func NewConformanceRules_ConformanceResult_1(level ConformanceRules_ConformanceLevel) *ConformanceRules_ConformanceResult {
	clinit_ConformanceRules_ConformanceResult()
	this_ := &ConformanceRules_ConformanceResult{}
	InitConformanceRules_ConformanceResult_1(this_, level)
	return this_
}

func InitConformanceRules_ConformanceResult_1(this_ *ConformanceRules_ConformanceResult, level ConformanceRules_ConformanceLevel) {
	_ = level
	InitConformanceRules_ConformanceResult_2(this_, level, "")
}

// Java constructor: ConformanceResult(com.google.javascript.jscomp.ConformanceRules.ConformanceLevel level, java.lang.String note)
func NewConformanceRules_ConformanceResult_2(level ConformanceRules_ConformanceLevel, note string) *ConformanceRules_ConformanceResult {
	clinit_ConformanceRules_ConformanceResult()
	this_ := &ConformanceRules_ConformanceResult{}
	InitConformanceRules_ConformanceResult_2(this_, level, note)
	return this_
}

func InitConformanceRules_ConformanceResult_2(this_ *ConformanceRules_ConformanceResult, level ConformanceRules_ConformanceLevel, note string) {
	_ = level
	_ = note
	this_.initFields_ConformanceRules_ConformanceResult()
	this_.level = level
	this_.note = note
}

// Java enum: com.google.javascript.jscomp.ConformanceRules.ConformanceLevel
type ConformanceRules_ConformanceLevel uint8

const (
	enum_ConformanceRules_ConformanceLevel_CONFORMANCE        ConformanceRules_ConformanceLevel = 1
	enum_ConformanceRules_ConformanceLevel_POSSIBLE_VIOLATION ConformanceRules_ConformanceLevel = 2
	enum_ConformanceRules_ConformanceLevel_VIOLATION          ConformanceRules_ConformanceLevel = 3
)

var enumBaseTable_ConformanceRules_ConformanceLevel = [...]j2gort.EnumBase{
	{Name_: "CONFORMANCE", Ordinal_: 0},
	{Name_: "POSSIBLE_VIOLATION", Ordinal_: 1},
	{Name_: "VIOLATION", Ordinal_: 2},
}
var enumValues_ConformanceRules_ConformanceLevel = []ConformanceRules_ConformanceLevel{1, 2, 3}

func StaticGet_ConformanceRules_ConformanceLevel_CONFORMANCE() ConformanceRules_ConformanceLevel {
	return enum_ConformanceRules_ConformanceLevel_CONFORMANCE
}
func StaticGet_ConformanceRules_ConformanceLevel_POSSIBLE_VIOLATION() ConformanceRules_ConformanceLevel {
	return enum_ConformanceRules_ConformanceLevel_POSSIBLE_VIOLATION
}
func StaticGet_ConformanceRules_ConformanceLevel_VIOLATION() ConformanceRules_ConformanceLevel {
	return enum_ConformanceRules_ConformanceLevel_VIOLATION
}

func ConformanceRules_ConformanceLevel_Values() []ConformanceRules_ConformanceLevel {
	return j2gort.CloneSlice(enumValues_ConformanceRules_ConformanceLevel)
}

func ConformanceRules_ConformanceLevel_ValueOf(name string) ConformanceRules_ConformanceLevel {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_ConformanceRules_ConformanceLevel {
		if enumBaseTable_ConformanceRules_ConformanceLevel[i_].Name_ == name {
			return ConformanceRules_ConformanceLevel(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.ConformanceRules.ConformanceLevel." + name))
}

func (this_ ConformanceRules_ConformanceLevel) IsJavaNil_() bool                              { return this_ == 0 }
func (ConformanceRules_ConformanceLevel) HasIdentityEquals_()                                 {}
func (ConformanceRules_ConformanceLevel) Implements_j2gort_Ext_java_lang_Enum()               {}
func (ConformanceRules_ConformanceLevel) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ ConformanceRules_ConformanceLevel) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (ConformanceRules_ConformanceLevel) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConformanceRules$ConformanceLevel", "ConformanceLevel"
}
func (this_ ConformanceRules_ConformanceLevel) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_ConformanceRules_ConformanceLevel[this_-1]
}
func (this_ ConformanceRules_ConformanceLevel) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ ConformanceRules_ConformanceLevel) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ ConformanceRules_ConformanceLevel) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ConformanceRules_ConformanceLevel[this_-1].Name_
}
func (this_ ConformanceRules_ConformanceLevel) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_ConformanceRules_ConformanceLevel[this_-1].Name_
}
func (this_ ConformanceRules_ConformanceLevel) Equals(other_ any) bool {
	o_, ok_ := other_.(ConformanceRules_ConformanceLevel)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ ConformanceRules_ConformanceLevel) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_ConformanceRules_ConformanceLevel[this_-1].HashCode()
}
func (this_ ConformanceRules_ConformanceLevel) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ConformanceRules_ConformanceLevel[this_-1].CompareTo(other_)
}
func (this_ ConformanceRules_ConformanceLevel) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_ConformanceRules_ConformanceLevel[this_-1].Name_)
}
func init() {
	vals_ := []any{ConformanceRules_ConformanceLevel(1), ConformanceRules_ConformanceLevel(2), ConformanceRules_ConformanceLevel(3)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.ConformanceRules.ConformanceLevel", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.ConformanceRules$ConformanceLevel", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConformanceRules.ConformanceLevel", "com.google.javascript.jscomp.ConformanceRules$ConformanceLevel", func(v_ any) bool {
		return j2gort.InstanceOf[ConformanceRules_ConformanceLevel](v_)
	})
}

// Java class: com.google.javascript.jscomp.ConformanceRules.AllowList
type ConformanceRules_AllowList struct {
	prefixes       []string                                                       `gson:"prefixes"`
	regexp         *j2gort.Ext_java_util_regex_Pattern                            `gson:"regexp"`
	allowlistEntry *j2gort.Ext_com_google_javascript_jscomp_RequirementScopeEntry `gson:"allowlistEntry"`
	idHash_        atomic.Int32
}

func (this_ *ConformanceRules_AllowList) BindSelf_(s any) {
	_ = s
}
func (this_ *ConformanceRules_AllowList) IsJavaNil_() bool { return this_ == nil }
func (*ConformanceRules_AllowList) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConformanceRules$AllowList", "AllowList"
}
func (this_ *ConformanceRules_AllowList) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ConformanceRules_AllowList) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConformanceRules.AllowList", "com.google.javascript.jscomp.ConformanceRules$AllowList", func(v_ any) bool {
		return j2gort.InstanceOf[*ConformanceRules_AllowList](v_)
	})
}

func (this_ *ConformanceRules_AllowList) initFields_ConformanceRules_AllowList() {
	_ = this_
}

// Java constructor: AllowList(java.util.List<java.lang.String> prefixes, java.util.List<java.lang.String> regexps)
func NewConformanceRules_AllowList_2(prefixes j2gort.Ext_java_util_List, regexps j2gort.Ext_java_util_List) *ConformanceRules_AllowList {
	this_ := &ConformanceRules_AllowList{}
	InitConformanceRules_AllowList_2(this_, prefixes, regexps)
	return this_
}

func InitConformanceRules_AllowList_2(this_ *ConformanceRules_AllowList, prefixes j2gort.Ext_java_util_List, regexps j2gort.Ext_java_util_List) {
	_ = prefixes
	_ = regexps
	this_.initFields_ConformanceRules_AllowList()
	this_.prefixes = j2gort.ImmListCopyToSlice[string](j2gort.AsIface[j2gort.Ext_java_util_Collection](prefixes))
	this_.regexp = ConformanceRules_BuildPattern(regexps)
	this_.allowlistEntry = nil
}

// Java constructor: AllowList(com.google.javascript.jscomp.RequirementScopeEntry allowlistEntry)
func NewConformanceRules_AllowList_1(allowlistEntry *j2gort.Ext_com_google_javascript_jscomp_RequirementScopeEntry) *ConformanceRules_AllowList {
	this_ := &ConformanceRules_AllowList{}
	InitConformanceRules_AllowList_1(this_, allowlistEntry)
	return this_
}

func InitConformanceRules_AllowList_1(this_ *ConformanceRules_AllowList, allowlistEntry *j2gort.Ext_com_google_javascript_jscomp_RequirementScopeEntry) {
	_ = allowlistEntry
	this_.initFields_ConformanceRules_AllowList()
	this_.prefixes = j2gort.ImmListCopyToSlice[string](j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](allowlistEntry.GetPrefixList()))
	this_.regexp = ConformanceRules_BuildPattern(j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](allowlistEntry.GetRegexpList()))
	this_.allowlistEntry = allowlistEntry
}

// Java method: AllowList.matches(java.lang.String path)
func (this_ *ConformanceRules_AllowList) Matches(path string) bool {
	_ = path
	var tsPath string
	_ = tsPath
	if javalang.StringEndsWith(j2gort.CheckNotNullString(path), ".closure.js") {
		tsPath = (j2gort.StrValString(j2gort.JStrSub(path, int32(0), (j2gort.JStrLen(path)-j2gort.JStrLen(".closure.js")))) + ".ts")
	} else if javalang.StringEndsWith(j2gort.CheckNotNullString(path), ".tsx.cl.js") {
		tsPath = j2gort.JStrSub(path, int32(0), (j2gort.JStrLen(path) - j2gort.JStrLen(".cl.js")))
	} else {
		tsPath = j2gort.NullString
	}
	if this_.prefixes != nil {
		{

			for _, prefix := range this_.prefixes {

				_ = prefix
				if !(j2gort.CheckNotNullString(path) == "") && (javalang.StringStartsWith(j2gort.CheckNotNullString(path), j2gort.CheckNotNullString(prefix)) || ((tsPath != j2gort.NullString) && javalang.StringStartsWith(j2gort.CheckNotNullString(tsPath), j2gort.CheckNotNullString(prefix)))) {
					return true
				}
			}
		}
	}
	return ((this_.regexp != nil) && (this_.regexp.Matcher(j2gort.BoxString(path)).Find_0() || ((tsPath != j2gort.NullString) && this_.regexp.Matcher(j2gort.BoxString(tsPath)).Find_0())))
}

// Java polymorphic class interface: com.google.javascript.jscomp.ConformanceRules.AbstractRule
type ConformanceRules_AbstractRule interface {
	Implements_jscomp_ConformanceRules_AbstractRule()
	IsJavaNil_() bool
	GetConformanceRules_AbstractRuleBase() *ConformanceRules_AbstractRuleBase
	CheckConformance_Rule
	TsIsAllowlisted() bool
	CheckConformance_2(*NodeTraversal, pkg_rhino.Node) *ConformanceRules_ConformanceResult
	Report_3(pkg_rhino.Node, *ConformanceRules_ConformanceResult, *j2gort.Ext_com_google_javascript_jscomp_ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior)
}

func ToIface_ConformanceRules_AbstractRule[T interface {
	ConformanceRules_AbstractRule
	comparable
}](v T) ConformanceRules_AbstractRule {
	var z T
	if v == z {
		return nil
	}
	return v
}

// Java class: com.google.javascript.jscomp.ConformanceRules.AbstractRule
type ConformanceRules_AbstractRuleBase struct {
	self_                     ConformanceRules_AbstractRule
	compiler                  AbstractCompiler                                              `gson:"compiler"`
	message                   string                                                        `gson:"message"`
	severity                  *j2gort.Ext_com_google_javascript_jscomp_Requirement_Severity `gson:"severity"`
	allowlists                []*ConformanceRules_AllowList                                 `gson:"allowlists"`
	onlyApplyTo               *ConformanceRules_AllowList                                   `gson:"onlyApplyTo"`
	requirement               *j2gort.Ext_com_google_javascript_jscomp_Requirement          `gson:"requirement"`
	idHash_                   atomic.Int32
	reportLooseTypeViolations bool                 `gson:"reportLooseTypeViolations"`
	typeMatchingStrategy      TypeMatchingStrategy `gson:"typeMatchingStrategy"`
}

func (this_ *ConformanceRules_AbstractRuleBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(ConformanceRules_AbstractRule)
}
func (this_ *ConformanceRules_AbstractRuleBase) IsJavaNil_() bool { return this_ == nil }
func (*ConformanceRules_AbstractRuleBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConformanceRules$AbstractRule", "AbstractRule"
}
func (this_ *ConformanceRules_AbstractRuleBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConformanceRules.AbstractRule", "com.google.javascript.jscomp.ConformanceRules$AbstractRule", func(v_ any) bool {
		return j2gort.InstanceOf[ConformanceRules_AbstractRule](v_)
	})
}
func (b *ConformanceRules_AbstractRuleBase) GetConformanceRules_AbstractRuleBase() *ConformanceRules_AbstractRuleBase {
	return b
}
func (*ConformanceRules_AbstractRuleBase) Implements_jscomp_ConformanceRules_AbstractRule() {}
func (*ConformanceRules_AbstractRuleBase) Implements_jscomp_CheckConformance_Rule()         {}

func (this_ *ConformanceRules_AbstractRuleBase) initFields_ConformanceRules_AbstractRule() {
	_ = this_
	this_.message = j2gort.NullString
}

func InitConformanceRules_AbstractRule(this_ *ConformanceRules_AbstractRuleBase, compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) {
	_ = compiler
	_ = requirement
	this_.initFields_ConformanceRules_AbstractRule()
	if !(requirement.HasErrorMessage()) {
		panic(NewCheckConformance_InvalidRequirementSpec_1("missing message"))
	}
	this_.compiler = compiler
	var message string = requirement.GetErrorMessage()
	_ = message
	if requirement.GetConfigFileCount() > int32(0) {
		message = (j2gort.StrValString(message) + "\n  defined in " + j2gort.StrValString(j2gort.Ext_com_google_common_base_Joiner_On_1_string("\n  extended by ").Join_1_Ext_java_lang_Iterable(j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](requirement.GetConfigFileList()))))
	}
	this_.message = message
	if requirement.GetSeverity() == j2gort.Ext_com_google_javascript_jscomp_Requirement_Severity_F_UNSPECIFIED {
		this_.severity = j2gort.Ext_com_google_javascript_jscomp_Requirement_Severity_F_WARNING
	} else {
		this_.severity = requirement.GetSeverity()
	}
	var allowlistsBuilder []*ConformanceRules_AllowList = []*ConformanceRules_AllowList{}
	_ = allowlistsBuilder
	{
		cur_11 := j2gort.NewIterCursor(j2gort.CheckNotNull(requirement.GetWhitelistEntryList()))
		for cur_11.HasNext() {
			rawElem_12 := cur_11.Next()
			entry := j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_RequirementScopeEntry](rawElem_12)
			_ = entry
			allowlistsBuilder = append(allowlistsBuilder, NewConformanceRules_AllowList_1(entry))
		}
	}
	{
		cur_13 := j2gort.NewIterCursor(j2gort.CheckNotNull(requirement.GetAllowlistEntryList()))
		for cur_13.HasNext() {
			rawElem_14 := cur_13.Next()
			entry_2 := j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_RequirementScopeEntry](rawElem_14)
			_ = entry_2
			allowlistsBuilder = append(allowlistsBuilder, NewConformanceRules_AllowList_1(entry_2))
		}
	}
	if this_.self_.TsIsAllowlisted() {
		allowlistsBuilder = append(allowlistsBuilder, j2gort.CheckNotNullElem(StaticGet_ConformanceRules_ALL_TS_ALLOWLIST()))
	}
	if (requirement.GetWhitelistCount() > int32(0)) || (requirement.GetWhitelistRegexpCount() > int32(0)) {
		var val_15 ConformanceRules_AllowList
		InitConformanceRules_AllowList_2(&val_15, j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](requirement.GetWhitelistList()), j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](requirement.GetWhitelistRegexpList()))
		var allowlist *ConformanceRules_AllowList = &val_15
		_ = allowlist
		allowlistsBuilder = append(allowlistsBuilder, j2gort.CheckNotNullElem(allowlist))
	}
	if (requirement.GetAllowlistCount() > int32(0)) || (requirement.GetAllowlistRegexpCount() > int32(0)) {
		var val_16 ConformanceRules_AllowList
		InitConformanceRules_AllowList_2(&val_16, j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](requirement.GetAllowlistList()), j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](requirement.GetAllowlistRegexpList()))
		var allowlist_2 *ConformanceRules_AllowList = &val_16
		_ = allowlist_2
		allowlistsBuilder = append(allowlistsBuilder, j2gort.CheckNotNullElem(allowlist_2))
	}
	this_.allowlists = allowlistsBuilder
	if (requirement.GetOnlyApplyToCount() > int32(0)) || (requirement.GetOnlyApplyToRegexpCount() > int32(0)) {
		this_.onlyApplyTo = NewConformanceRules_AllowList_2(j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](requirement.GetOnlyApplyToList()), j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](requirement.GetOnlyApplyToRegexpList()))
	} else {
		this_.onlyApplyTo = nil
	}
	this_.reportLooseTypeViolations = requirement.GetReportLooseTypeViolations()
	this_.typeMatchingStrategy = ConformanceRules_AbstractRule_GetTypeMatchingStrategy(requirement)
	this_.requirement = requirement
}

// Java method: AbstractRule.tsIsAllowlisted()
func (this_ *ConformanceRules_AbstractRuleBase) TsIsAllowlisted() bool {
	return false
}

// Java static method: AbstractRule.getTypeMatchingStrategy(com.google.javascript.jscomp.Requirement requirement)
func ConformanceRules_AbstractRule_GetTypeMatchingStrategy(requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) TypeMatchingStrategy {
	_ = requirement
	switch j2gort.CheckNotNullPtr(requirement.GetTypeMatchingStrategy()).Ordinal() {
	case int32(1) /* TypeMatchingStrategy.LOOSE */ :
		return StaticGet_TypeMatchingStrategy_LOOSE()
	case int32(2) /* TypeMatchingStrategy.STRICT_NULLABILITY */ :
		return StaticGet_TypeMatchingStrategy_STRICT_NULLABILITY()
	case int32(3) /* TypeMatchingStrategy.SUBTYPES */ :
		return StaticGet_TypeMatchingStrategy_SUBTYPES()
	case int32(4) /* TypeMatchingStrategy.EXACT */ :
		return StaticGet_TypeMatchingStrategy_EXACT()
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unknown TypeMatchingStrategy")))
	}
	panic("unreachable switch expression")
}

// Java method: AbstractRule.checkConformance(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *ConformanceRules_AbstractRuleBase) CheckConformance_2(t *NodeTraversal, n pkg_rhino.Node) *ConformanceRules_ConformanceResult {
	_ = t
	_ = n
	panic("abstract or native method invoked")
}

// Java method: AbstractRule.findAllowListForPath(java.lang.String path)
func (this_ *ConformanceRules_AbstractRuleBase) priv_AbstractRule_findAllowListForPath_2140(path string) *ConformanceRules_AllowList {
	_ = path
	var pathRegex j2gort.Ext_com_google_common_base_Optional = this_.compiler.GetOptions().GetConformanceRemoveRegexFromPath()
	_ = pathRegex
	if pathRegex.IsPresent() {
		path = j2gort.CastPtr[j2gort.Ext_java_util_regex_Pattern](pathRegex.Get()).Matcher(j2gort.BoxString(path)).ReplaceFirst_1_string("")
	}
	{

		for _, allowlist := range this_.allowlists {

			_ = allowlist
			if allowlist.Matches(path) {
				return allowlist
			}
		}
	}
	return nil
}

// Java method: AbstractRule.check(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.jscomp.ConformanceConfig.LibraryLevelNonAllowlistedConformanceViolationsBehavior behavior)
func (this_ *ConformanceRules_AbstractRuleBase) Check(t *NodeTraversal, n pkg_rhino.Node, behavior *j2gort.Ext_com_google_javascript_jscomp_ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior) {
	_ = t
	_ = n
	_ = behavior
	var result *ConformanceRules_ConformanceResult = this_.self_.CheckConformance_2(t, n)
	_ = result
	if result.level != StaticGet_ConformanceRules_ConformanceLevel_CONFORMANCE() {
		this_.Report_3(n, result, behavior)
	}
}

// Java method: AbstractRule.report(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.ConformanceRules.ConformanceResult result, com.google.javascript.jscomp.ConformanceConfig.LibraryLevelNonAllowlistedConformanceViolationsBehavior behavior)
func (this_ *ConformanceRules_AbstractRuleBase) Report_3(n pkg_rhino.Node, result *ConformanceRules_ConformanceResult, behavior *j2gort.Ext_com_google_javascript_jscomp_ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior) {
	_ = n
	_ = result
	_ = behavior
	var msg *DiagnosticType = nil
	_ = msg
	if this_.severity == j2gort.Ext_com_google_javascript_jscomp_Requirement_Severity_F_ERROR {
		msg = StaticGet_CheckConformance_CONFORMANCE_ERROR()
	} else {
		if result.level == StaticGet_ConformanceRules_ConformanceLevel_VIOLATION() {
			msg = StaticGet_CheckConformance_CONFORMANCE_VIOLATION()
		} else {
			msg = StaticGet_CheckConformance_CONFORMANCE_POSSIBLE_VIOLATION()
		}
	}
	var separator string
	_ = separator
	if j2gort.CheckNotNullString(result.note) == "" {
		separator = ""
	} else {
		separator = "\n"
	}
	var err *JSError = JSError_Make_4_Ext_com_google_javascript_jscomp_Requirement_Node_DiagnosticType_slice_string(this_.requirement, n, msg, []string{this_.message, separator, result.note})
	_ = err
	var path string = NodeUtil_GetSourceName(n)
	_ = path
	var allowlist *ConformanceRules_AllowList
	_ = allowlist
	if path != j2gort.NullString {
		allowlist = this_.priv_AbstractRule_findAllowListForPath_2140(path)
	} else {
		allowlist = nil
	}
	var isAllowlisted bool = !((allowlist == nil) && ((this_.onlyApplyTo == nil) || this_.onlyApplyTo.Matches(path)))
	_ = isAllowlisted
	var shouldReport bool = this_.compiler.GetErrorManager().ShouldReportConformanceViolation(this_.requirement, func() j2gort.Ext_com_google_common_base_Optional {
		if allowlist != nil {
			return j2gort.Ext_com_google_common_base_Optional_FromNullable(j2gort.PtrToAny(allowlist.allowlistEntry))
		}
		return j2gort.Ext_com_google_common_base_Optional_Absent()
	}(), err, behavior, isAllowlisted)
	_ = shouldReport
	if shouldReport && !(isAllowlisted) {
		this_.compiler.Report_1(err)
	}
}

func (this_ *ConformanceRules_AbstractRuleBase) GetPrecondition() CheckConformance_Precondition {
	return IfaceDefault_CheckConformance_Rule_GetPrecondition(this_.self_)
}

// Java class: com.google.javascript.jscomp.ConformanceRules.NoOp
type ConformanceRules_NoOp struct {
	ConformanceRules_AbstractRuleBase
}

func (this_ *ConformanceRules_NoOp) BindSelf_(s any) {
	_ = s
	this_.ConformanceRules_AbstractRuleBase.BindSelf_(s)
}
func (this_ *ConformanceRules_NoOp) IsJavaNil_() bool { return this_ == nil }
func (*ConformanceRules_NoOp) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConformanceRules$NoOp", "NoOp"
}
func (this_ *ConformanceRules_NoOp) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ConformanceRules_AbstractRuleBase.IdentityHash_()
}
func (*ConformanceRules_NoOp) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConformanceRules.NoOp", "com.google.javascript.jscomp.ConformanceRules$NoOp", func(v_ any) bool {
		return j2gort.InstanceOf[*ConformanceRules_NoOp](v_)
	})
}
func (*ConformanceRules_NoOp) Implements_jscomp_ConformanceRules_AbstractRule() {}
func (*ConformanceRules_NoOp) Implements_jscomp_CheckConformance_Rule()         {}

func (this_ *ConformanceRules_NoOp) initFields_ConformanceRules_NoOp() {
	_ = this_
}

// Java constructor: NoOp(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.Requirement requirement)
func NewConformanceRules_NoOp(compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) *ConformanceRules_NoOp {
	this_ := &ConformanceRules_NoOp{}
	this_.BindSelf_(this_)
	InitConformanceRules_NoOp(this_, compiler, requirement)
	return this_
}

func InitConformanceRules_NoOp(this_ *ConformanceRules_NoOp, compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) {
	_ = compiler
	_ = requirement
	InitConformanceRules_AbstractRule(&this_.ConformanceRules_AbstractRuleBase, compiler, requirement)
	this_.initFields_ConformanceRules_NoOp()
}

// Java method: NoOp.checkConformance(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *ConformanceRules_NoOp) CheckConformance_2(t *NodeTraversal, n pkg_rhino.Node) *ConformanceRules_ConformanceResult {
	_ = t
	_ = n
	return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
}

// Java polymorphic class interface: com.google.javascript.jscomp.ConformanceRules.AbstractTypeRestrictionRule
type ConformanceRules_AbstractTypeRestrictionRule interface {
	Implements_jscomp_ConformanceRules_AbstractTypeRestrictionRule()
	IsJavaNil_() bool
	GetConformanceRules_AbstractTypeRestrictionRuleBase() *ConformanceRules_AbstractTypeRestrictionRuleBase
	ConformanceRules_AbstractRule
	IsAllowlistedType(pkg_rhino.Node) bool
	IsNativeObjectType(pkg_rhino.Node) bool
	Union(j2gort.Ext_java_util_List) pkg_rhino.JSType
	IsAssertionCall(pkg_rhino.Node) bool
	IsTypeImmediatelyTightened(pkg_rhino.Node) bool
	IsUsed(pkg_rhino.Node) bool
}

func ToIface_ConformanceRules_AbstractTypeRestrictionRule[T interface {
	ConformanceRules_AbstractTypeRestrictionRule
	comparable
}](v T) ConformanceRules_AbstractTypeRestrictionRule {
	var z T
	if v == z {
		return nil
	}
	return v
}

// Java class: com.google.javascript.jscomp.ConformanceRules.AbstractTypeRestrictionRule
type ConformanceRules_AbstractTypeRestrictionRuleBase struct {
	ConformanceRules_AbstractRuleBase
	nativeObjectType   pkg_rhino.JSType                          `gson:"nativeObjectType"`
	allowlistedTypes   pkg_rhino.JSType                          `gson:"allowlistedTypes"`
	assertionFunctions *CodingConvention_AssertionFunctionLookup `gson:"assertionFunctions"`
}

func (this_ *ConformanceRules_AbstractTypeRestrictionRuleBase) BindSelf_(s any) {
	_ = s
	this_.ConformanceRules_AbstractRuleBase.BindSelf_(s)
}
func (this_ *ConformanceRules_AbstractTypeRestrictionRuleBase) IsJavaNil_() bool { return this_ == nil }
func (*ConformanceRules_AbstractTypeRestrictionRuleBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConformanceRules$AbstractTypeRestrictionRule", "AbstractTypeRestrictionRule"
}
func (this_ *ConformanceRules_AbstractTypeRestrictionRuleBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ConformanceRules_AbstractRuleBase.IdentityHash_()
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConformanceRules.AbstractTypeRestrictionRule", "com.google.javascript.jscomp.ConformanceRules$AbstractTypeRestrictionRule", func(v_ any) bool {
		return j2gort.InstanceOf[ConformanceRules_AbstractTypeRestrictionRule](v_)
	})
}
func (b *ConformanceRules_AbstractTypeRestrictionRuleBase) GetConformanceRules_AbstractTypeRestrictionRuleBase() *ConformanceRules_AbstractTypeRestrictionRuleBase {
	return b
}
func (*ConformanceRules_AbstractTypeRestrictionRuleBase) Implements_jscomp_ConformanceRules_AbstractTypeRestrictionRule() {
}
func (*ConformanceRules_AbstractTypeRestrictionRuleBase) Implements_jscomp_ConformanceRules_AbstractRule() {
}
func (*ConformanceRules_AbstractTypeRestrictionRuleBase) Implements_jscomp_CheckConformance_Rule() {}

func (this_ *ConformanceRules_AbstractTypeRestrictionRuleBase) initFields_ConformanceRules_AbstractTypeRestrictionRule() {
	_ = this_
}

func InitConformanceRules_AbstractTypeRestrictionRule(this_ *ConformanceRules_AbstractTypeRestrictionRuleBase, compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) {
	_ = compiler
	_ = requirement
	InitConformanceRules_AbstractRule(&this_.ConformanceRules_AbstractRuleBase, compiler, requirement)
	this_.initFields_ConformanceRules_AbstractTypeRestrictionRule()
	this_.nativeObjectType = compiler.GetTypeRegistry().GetNativeType(pkg_rhino.StaticGet_JSTypeNative_OBJECT_TYPE())
	var allowlistedTypeNames j2gort.Ext_java_util_List = j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](requirement.GetValueList())
	_ = allowlistedTypeNames
	this_.allowlistedTypes = this_.Union(allowlistedTypeNames)
	this_.assertionFunctions = CodingConvention_AssertionFunctionLookup_Of(j2gort.AsIface[j2gort.Ext_java_util_Collection](compiler.GetCodingConvention().GetAssertionFunctions()))
}

// Java method: AbstractTypeRestrictionRule.isAllowlistedType(com.google.javascript.rhino.Node n)
func (this_ *ConformanceRules_AbstractTypeRestrictionRuleBase) IsAllowlistedType(n pkg_rhino.Node) bool {
	_ = n
	if (this_.allowlistedTypes != nil && !this_.allowlistedTypes.IsJavaNil_()) && !(j2gort.IfaceEquals[pkg_rhino.JSType](n.GetJSType(), nil)) {
		var targetType pkg_rhino.JSType = n.GetJSType().RestrictByNotNullOrUndefined()
		_ = targetType
		if targetType.IsSubtypeOf_1(this_.allowlistedTypes) {
			return true
		}
	}
	return false
}

// Java static method: AbstractTypeRestrictionRule.isKnown(com.google.javascript.rhino.Node n)
func ConformanceRules_AbstractTypeRestrictionRule_IsKnown(n pkg_rhino.Node) bool {
	_ = n
	return ((!(ConformanceRules_AbstractTypeRestrictionRule_IsUnknown(n)) && !(ConformanceRules_AbstractTypeRestrictionRule_IsBottom(n))) && !(ConformanceRules_AbstractTypeRestrictionRule_IsTemplateType(n)))
}

// Java method: AbstractTypeRestrictionRule.isNativeObjectType(com.google.javascript.rhino.Node n)
func (this_ *ConformanceRules_AbstractTypeRestrictionRuleBase) IsNativeObjectType(n pkg_rhino.Node) bool {
	_ = n
	var type_ pkg_rhino.JSType = n.GetJSType().RestrictByNotNullOrUndefined()
	_ = type_
	return type_.Equals(this_.nativeObjectType)
}

// Java static method: AbstractTypeRestrictionRule.isTop(com.google.javascript.rhino.Node n)
func ConformanceRules_AbstractTypeRestrictionRule_IsTop(n pkg_rhino.Node) bool {
	_ = n
	var type_ pkg_rhino.JSType = n.GetJSType()
	_ = type_
	return ((type_ != nil && !type_.IsJavaNil_()) && type_.IsAllType())
}

// Java static method: AbstractTypeRestrictionRule.isUnknown(com.google.javascript.rhino.Node n)
func ConformanceRules_AbstractTypeRestrictionRule_IsUnknown(n pkg_rhino.Node) bool {
	_ = n
	var type_ pkg_rhino.JSType = n.GetJSType()
	_ = type_
	return ((type_ == nil || type_.IsJavaNil_()) || type_.IsUnknownType())
}

// Java static method: AbstractTypeRestrictionRule.isTemplateType(com.google.javascript.rhino.Node n)
func ConformanceRules_AbstractTypeRestrictionRule_IsTemplateType(n pkg_rhino.Node) bool {
	_ = n
	var type_ pkg_rhino.JSType = n.GetJSType()
	_ = type_
	if type_.IsUnionType() {
		{
			cur_19 := j2gort.NewIterCursor(j2gort.CheckNotNull(type_.GetUnionMembers()))
			for cur_19.HasNext() {
				rawElem_20 := cur_19.Next()
				member := j2gort.AsIface[pkg_rhino.JSType](rawElem_20)
				_ = member
				if member.IsTemplateType() {
					return true
				}
			}
		}
	}
	return type_.IsTemplateType()
}

// Java static method: AbstractTypeRestrictionRule.isBottom(com.google.javascript.rhino.Node n)
func ConformanceRules_AbstractTypeRestrictionRule_IsBottom(n pkg_rhino.Node) bool {
	_ = n
	var type_ pkg_rhino.JSType = n.GetJSType().RestrictByNotNullOrUndefined()
	_ = type_
	return type_.IsEmptyType()
}

// Java method: AbstractTypeRestrictionRule.union(java.util.List<java.lang.String> typeNames)
func (this_ *ConformanceRules_AbstractTypeRestrictionRuleBase) Union(typeNames j2gort.Ext_java_util_List) pkg_rhino.JSType {
	_ = typeNames
	var registry *pkg_rhino.JSTypeRegistry = this_.GetConformanceRules_AbstractRuleBase().compiler.GetTypeRegistry()
	_ = registry
	var types j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList()
	_ = types
	{
		cur_21 := j2gort.NewIterCursor(j2gort.CheckNotNull(typeNames))
		for cur_21.HasNext() {
			typeName := cur_21.NextStr()
			_ = typeName
			var type_ pkg_rhino.JSType = registry.GetGlobalType(typeName)
			_ = type_
			if type_ == nil || type_.IsJavaNil_() {
				type_ = registry.ResolveViaClosureNamespace(typeName)
			}
			if type_ != nil && !type_.IsJavaNil_() {
				types.Add_1(type_)
			}
		}
	}
	if types.IsEmpty() {
		return nil
	} else {
		return registry.CreateUnionType_1_Ext_java_util_List(types)
	}
}

// Java method: AbstractTypeRestrictionRule.isAssertionCall(com.google.javascript.rhino.Node n)
func (this_ *ConformanceRules_AbstractTypeRestrictionRuleBase) IsAssertionCall(n pkg_rhino.Node) bool {
	_ = n
	if n.IsCall() && n.GetFirstChild().IsQualifiedName() {
		var target pkg_rhino.Node = n.GetFirstChild()
		_ = target
		return !(j2gort.IfaceEquals[CodingConvention_AssertionFunctionSpec](this_.assertionFunctions.LookupByCallee(target), nil))
	}
	return false
}

// Java method: AbstractTypeRestrictionRule.isTypeImmediatelyTightened(com.google.javascript.rhino.Node n)
func (this_ *ConformanceRules_AbstractTypeRestrictionRuleBase) IsTypeImmediatelyTightened(n pkg_rhino.Node) bool {
	_ = n
	return ((this_.IsAssertionCall(n.GetParent()) || n.GetParent().IsTypeOf()) || !(j2gort.IfaceEquals[pkg_rhino.JSType](n.GetJSTypeBeforeCast(), nil)))
}

// Java method: AbstractTypeRestrictionRule.isUsed(com.google.javascript.rhino.Node n)
func (this_ *ConformanceRules_AbstractTypeRestrictionRuleBase) IsUsed(n pkg_rhino.Node) bool {
	_ = n
	if n.GetParent().IsName() || NodeUtil_IsLhsByDestructuring(n) {
		return false
	}
	if NodeUtil_IsAssignmentOp(n.GetParent()) {
		return NodeUtil_IsExpressionResultUsed(n.GetParent())
	}
	return NodeUtil_IsExpressionResultUsed(n)
}

// Java class: com.google.javascript.jscomp.ConformanceRules.InferredConstCheck
type ConformanceRules_InferredConstCheck struct {
	ConformanceRules_AbstractRuleBase
}

func (this_ *ConformanceRules_InferredConstCheck) BindSelf_(s any) {
	_ = s
	this_.ConformanceRules_AbstractRuleBase.BindSelf_(s)
}
func (this_ *ConformanceRules_InferredConstCheck) IsJavaNil_() bool { return this_ == nil }
func (*ConformanceRules_InferredConstCheck) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConformanceRules$InferredConstCheck", "InferredConstCheck"
}
func (this_ *ConformanceRules_InferredConstCheck) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ConformanceRules_AbstractRuleBase.IdentityHash_()
}
func (*ConformanceRules_InferredConstCheck) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConformanceRules.InferredConstCheck", "com.google.javascript.jscomp.ConformanceRules$InferredConstCheck", func(v_ any) bool {
		return j2gort.InstanceOf[*ConformanceRules_InferredConstCheck](v_)
	})
}
func (*ConformanceRules_InferredConstCheck) Implements_jscomp_ConformanceRules_AbstractRule() {}
func (*ConformanceRules_InferredConstCheck) Implements_jscomp_CheckConformance_Rule()         {}

func (this_ *ConformanceRules_InferredConstCheck) initFields_ConformanceRules_InferredConstCheck() {
	_ = this_
}

// Java constructor: InferredConstCheck(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.Requirement requirement)
func NewConformanceRules_InferredConstCheck(compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) *ConformanceRules_InferredConstCheck {
	this_ := &ConformanceRules_InferredConstCheck{}
	this_.BindSelf_(this_)
	InitConformanceRules_InferredConstCheck(this_, compiler, requirement)
	return this_
}

func InitConformanceRules_InferredConstCheck(this_ *ConformanceRules_InferredConstCheck, compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) {
	_ = compiler
	_ = requirement
	InitConformanceRules_AbstractRule(&this_.ConformanceRules_AbstractRuleBase, compiler, requirement)
	this_.initFields_ConformanceRules_InferredConstCheck()
}

// Java method: InferredConstCheck.checkConformance(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *ConformanceRules_InferredConstCheck) CheckConformance_2(t *NodeTraversal, n pkg_rhino.Node) *ConformanceRules_ConformanceResult {
	_ = t
	_ = n
	var jsDoc *pkg_rhino.JSDocInfo = n.GetJSDocInfo()
	_ = jsDoc
	if ((jsDoc != nil) && jsDoc.HasConstAnnotation()) && (jsDoc.GetType() == (pkg_rhino.JSTypeExpression{})) {
		if n.IsAssign() {
			n = n.GetFirstChild()
		}
		var type_ pkg_rhino.JSType = n.GetJSType()
		_ = type_
		if ((type_ != nil && !type_.IsJavaNil_()) && type_.IsUnknownType()) && !(NodeUtil_IsNamespaceDecl(n)) {
			return StaticGet_ConformanceRules_ConformanceResult_VIOLATION()
		}
	}
	return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
}

// Java class: com.google.javascript.jscomp.ConformanceRules.BannedDependency
type ConformanceRules_BannedDependency struct {
	ConformanceRules_AbstractRuleBase
	paths j2gort.Ext_java_util_List `gson:"paths"`
}

func (this_ *ConformanceRules_BannedDependency) BindSelf_(s any) {
	_ = s
	this_.ConformanceRules_AbstractRuleBase.BindSelf_(s)
}
func (this_ *ConformanceRules_BannedDependency) IsJavaNil_() bool { return this_ == nil }
func (*ConformanceRules_BannedDependency) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConformanceRules$BannedDependency", "BannedDependency"
}
func (this_ *ConformanceRules_BannedDependency) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ConformanceRules_AbstractRuleBase.IdentityHash_()
}
func (*ConformanceRules_BannedDependency) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConformanceRules.BannedDependency", "com.google.javascript.jscomp.ConformanceRules$BannedDependency", func(v_ any) bool {
		return j2gort.InstanceOf[*ConformanceRules_BannedDependency](v_)
	})
}
func (*ConformanceRules_BannedDependency) Implements_jscomp_ConformanceRules_AbstractRule() {}
func (*ConformanceRules_BannedDependency) Implements_jscomp_CheckConformance_Rule()         {}

func (this_ *ConformanceRules_BannedDependency) initFields_ConformanceRules_BannedDependency() {
	_ = this_
}

// Java constructor: BannedDependency(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.Requirement requirement)
func NewConformanceRules_BannedDependency(compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) *ConformanceRules_BannedDependency {
	this_ := &ConformanceRules_BannedDependency{}
	this_.BindSelf_(this_)
	InitConformanceRules_BannedDependency(this_, compiler, requirement)
	return this_
}

func InitConformanceRules_BannedDependency(this_ *ConformanceRules_BannedDependency, compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) {
	_ = compiler
	_ = requirement
	InitConformanceRules_AbstractRule(&this_.ConformanceRules_AbstractRuleBase, compiler, requirement)
	this_.initFields_ConformanceRules_BannedDependency()
	this_.paths = j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](requirement.GetValueList())
	if this_.paths.IsEmpty() {
		panic(NewCheckConformance_InvalidRequirementSpec_1("missing value"))
	}
}

// Java method: BannedDependency.checkConformance(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *ConformanceRules_BannedDependency) CheckConformance_2(t *NodeTraversal, n pkg_rhino.Node) *ConformanceRules_ConformanceResult {
	_ = t
	_ = n
	if n.IsScript() {
		var srcFile string = n.GetSourceFileName()
		_ = srcFile
		{
			var i int32 = int32(0)
			_ = i
			for ; i < this_.paths.Size(); i = i + 1 {
				var path string = j2gort.ListGetStr(this_.paths, i)
				_ = path
				if javalang.StringStartsWith(j2gort.CheckNotNullString(srcFile), j2gort.CheckNotNullString(path)) {
					return StaticGet_ConformanceRules_ConformanceResult_VIOLATION()
				}
			}
		}
	}
	return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
}

// Java class: com.google.javascript.jscomp.ConformanceRules.BannedDependencyRegex
type ConformanceRules_BannedDependencyRegex struct {
	ConformanceRules_AbstractRuleBase
	pathRegexp *j2gort.Ext_java_util_regex_Pattern `gson:"pathRegexp"`
}

func (this_ *ConformanceRules_BannedDependencyRegex) BindSelf_(s any) {
	_ = s
	this_.ConformanceRules_AbstractRuleBase.BindSelf_(s)
}
func (this_ *ConformanceRules_BannedDependencyRegex) IsJavaNil_() bool { return this_ == nil }
func (*ConformanceRules_BannedDependencyRegex) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConformanceRules$BannedDependencyRegex", "BannedDependencyRegex"
}
func (this_ *ConformanceRules_BannedDependencyRegex) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ConformanceRules_AbstractRuleBase.IdentityHash_()
}
func (*ConformanceRules_BannedDependencyRegex) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConformanceRules.BannedDependencyRegex", "com.google.javascript.jscomp.ConformanceRules$BannedDependencyRegex", func(v_ any) bool {
		return j2gort.InstanceOf[*ConformanceRules_BannedDependencyRegex](v_)
	})
}
func (*ConformanceRules_BannedDependencyRegex) Implements_jscomp_ConformanceRules_AbstractRule() {}
func (*ConformanceRules_BannedDependencyRegex) Implements_jscomp_CheckConformance_Rule()         {}

func (this_ *ConformanceRules_BannedDependencyRegex) initFields_ConformanceRules_BannedDependencyRegex() {
	_ = this_
}

// Java constructor: BannedDependencyRegex(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.Requirement requirement)
func NewConformanceRules_BannedDependencyRegex(compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) *ConformanceRules_BannedDependencyRegex {
	this_ := &ConformanceRules_BannedDependencyRegex{}
	this_.BindSelf_(this_)
	InitConformanceRules_BannedDependencyRegex(this_, compiler, requirement)
	return this_
}

func InitConformanceRules_BannedDependencyRegex(this_ *ConformanceRules_BannedDependencyRegex, compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) {
	_ = compiler
	_ = requirement
	InitConformanceRules_AbstractRule(&this_.ConformanceRules_AbstractRuleBase, compiler, requirement)
	this_.initFields_ConformanceRules_BannedDependencyRegex()
	var pathRegexpList j2gort.Ext_java_util_List = j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](requirement.GetValueList())
	_ = pathRegexpList
	if pathRegexpList.IsEmpty() {
		panic(NewCheckConformance_InvalidRequirementSpec_1("missing value (no banned dependency regexps)"))
	}
	this_.pathRegexp = ConformanceRules_BuildPattern(pathRegexpList)
}

// Java method: BannedDependencyRegex.checkConformance(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *ConformanceRules_BannedDependencyRegex) CheckConformance_2(t *NodeTraversal, n pkg_rhino.Node) *ConformanceRules_ConformanceResult {
	_ = t
	_ = n
	if n.IsScript() {
		var srcFile string = n.GetSourceFileName()
		_ = srcFile
		if (this_.pathRegexp != nil) && this_.pathRegexp.Matcher(j2gort.BoxString(srcFile)).Find_0() {
			return StaticGet_ConformanceRules_ConformanceResult_VIOLATION()
		}
	}
	return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
}

// Java class: com.google.javascript.jscomp.ConformanceRules.BannedName
type ConformanceRules_BannedName struct {
	ConformanceRules_AbstractRuleBase
	requirementType *j2gort.Ext_com_google_javascript_jscomp_Requirement_Type `gson:"requirementType"`
	qualifiedNames  []pkg_rhino.Node                                          `gson:"qualifiedNames"`
	shortNames      j2gort.Ext_com_google_common_collect_ImmutableSet         `gson:"shortNames"`
}

func (this_ *ConformanceRules_BannedName) BindSelf_(s any) {
	_ = s
	this_.ConformanceRules_AbstractRuleBase.BindSelf_(s)
}
func (this_ *ConformanceRules_BannedName) IsJavaNil_() bool { return this_ == nil }
func (*ConformanceRules_BannedName) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConformanceRules$BannedName", "BannedName"
}
func (this_ *ConformanceRules_BannedName) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ConformanceRules_AbstractRuleBase.IdentityHash_()
}
func (*ConformanceRules_BannedName) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConformanceRules.BannedName", "com.google.javascript.jscomp.ConformanceRules$BannedName", func(v_ any) bool {
		return j2gort.InstanceOf[*ConformanceRules_BannedName](v_)
	})
}
func (*ConformanceRules_BannedName) Implements_jscomp_ConformanceRules_AbstractRule() {}
func (*ConformanceRules_BannedName) Implements_jscomp_CheckConformance_Rule()         {}

var static_ConformanceRules_BannedName_IS_CANDIDATE_NODE CheckConformance_Precondition

func StaticGet_ConformanceRules_BannedName_IS_CANDIDATE_NODE() CheckConformance_Precondition {
	clinit_ConformanceRules_BannedName()
	return static_ConformanceRules_BannedName_IS_CANDIDATE_NODE
}
func StaticSet_ConformanceRules_BannedName_IS_CANDIDATE_NODE(v CheckConformance_Precondition) {
	clinit_ConformanceRules_BannedName()
	static_ConformanceRules_BannedName_IS_CANDIDATE_NODE = v
}

var clinitDone_ConformanceRules_BannedName uint32
var clinitState_ConformanceRules_BannedName j2gort.ClinitState

func clinit_ConformanceRules_BannedName() {
	if atomic.LoadUint32(&clinitDone_ConformanceRules_BannedName) != 0 {
		return
	}
	clinitSlow_ConformanceRules_BannedName()
}

func clinitSlow_ConformanceRules_BannedName() {
	if !j2gort.EnterClinit(&clinitState_ConformanceRules_BannedName, &clinitDone_ConformanceRules_BannedName) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ConformanceRules_BannedName, &clinitDone_ConformanceRules_BannedName)
	static_ConformanceRules_BannedName_IS_CANDIDATE_NODE = lambdaSingleton_conformance_rules_1
}

func (this_ *ConformanceRules_BannedName) initFields_ConformanceRules_BannedName() {
	_ = this_
}

// Java constructor: BannedName(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.Requirement requirement)
func NewConformanceRules_BannedName(compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) *ConformanceRules_BannedName {
	this_ := &ConformanceRules_BannedName{}
	this_.BindSelf_(this_)
	InitConformanceRules_BannedName(this_, compiler, requirement)
	return this_
}

func InitConformanceRules_BannedName(this_ *ConformanceRules_BannedName, compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) {
	_ = compiler
	_ = requirement
	InitConformanceRules_AbstractRule(&this_.ConformanceRules_AbstractRuleBase, compiler, requirement)
	this_.initFields_ConformanceRules_BannedName()
	if requirement.GetValueCount() == int32(0) {
		panic(NewCheckConformance_InvalidRequirementSpec_1("missing value"))
	}
	this_.requirementType = requirement.GetType()
	var qualifiedBuilder []pkg_rhino.Node = []pkg_rhino.Node{}
	_ = qualifiedBuilder
	var shortBuilder j2gort.Ext_com_google_common_collect_ImmutableSet_Builder = j2gort.Ext_com_google_common_collect_ImmutableSet_S_Builder()
	_ = shortBuilder
	{
		cur_22 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](requirement.GetValueList())))
		for cur_22.HasNext() {
			name := cur_22.NextStr()
			_ = name
			var qualifiedName pkg_rhino.Node = NodeUtil_NewQName_2_AbstractCompiler_string(compiler, name)
			_ = qualifiedName
			qualifiedBuilder = append(qualifiedBuilder, j2gort.CheckNotNullElem(qualifiedName))
			j2gort.ImmSetBuilderAddStrVoid(shortBuilder, qualifiedName.GetString())
		}
	}
	this_.qualifiedNames = qualifiedBuilder
	this_.shortNames = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](shortBuilder.Build())
}

// Java method: BannedName.getPrecondition()
func (this_ *ConformanceRules_BannedName) GetPrecondition() CheckConformance_Precondition {
	return StaticGet_ConformanceRules_BannedName_IS_CANDIDATE_NODE()
}

// Java method: BannedName.checkConformance(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *ConformanceRules_BannedName) CheckConformance_2(t *NodeTraversal, n pkg_rhino.Node) *ConformanceRules_ConformanceResult {
	_ = t
	_ = n
	if (this_.requirementType == j2gort.Ext_com_google_javascript_jscomp_Requirement_Type_F_BANNED_NAME_CALL) && !(ConformanceRules_ConformanceUtil_IsCallTarget(n)) {
		return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
	}
	if !(j2gort.SetContainsStr(this_.shortNames, n.GetString())) {
		return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
	}
	if NodeUtil_IsInSyntheticScript(n) || !(ConformanceRules_BannedName_IsRootOfQualifiedNameGlobal(t, n)) {
		return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
	}
	{
		var i int32 = (int32(len(this_.qualifiedNames)) - int32(1))
		_ = i
		for ; i >= int32(0); i = i - 1 {
			if n.MatchesQualifiedName_1_Node(this_.qualifiedNames[i]) {
				return StaticGet_ConformanceRules_ConformanceResult_VIOLATION()
			}
		}
	}
	return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
}

// Java static method: BannedName.isRootOfQualifiedNameGlobal(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func ConformanceRules_BannedName_IsRootOfQualifiedNameGlobal(t *NodeTraversal, n pkg_rhino.Node) bool {
	_ = t
	_ = n
	var rootName string = NodeUtil_GetRootOfQualifiedName_1_Node(n).GetQualifiedName()
	_ = rootName
	var v *Var = j2gort.CastPtr[Var](t.GetScope().GetVar(rootName))
	_ = v
	return ((v != nil) && v.IsGlobal())
}

// Java class: com.google.javascript.jscomp.ConformanceRules.BannedProperty
type ConformanceRules_BannedProperty struct {
	ConformanceRules_AbstractRuleBase
	registry                *pkg_rhino.JSTypeRegistry                                  `gson:"registry"`
	props                   *j2gort.Ext_com_google_common_collect_ImmutableSetMultimap `gson:"props"`
	requirementPrecondition *ConformanceRules_BannedProperty_RequirementPrecondition   `gson:"requirementPrecondition"`
}

func (this_ *ConformanceRules_BannedProperty) BindSelf_(s any) {
	_ = s
	this_.ConformanceRules_AbstractRuleBase.BindSelf_(s)
}
func (this_ *ConformanceRules_BannedProperty) IsJavaNil_() bool { return this_ == nil }
func (*ConformanceRules_BannedProperty) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConformanceRules$BannedProperty", "BannedProperty"
}
func (this_ *ConformanceRules_BannedProperty) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ConformanceRules_AbstractRuleBase.IdentityHash_()
}
func (*ConformanceRules_BannedProperty) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConformanceRules.BannedProperty", "com.google.javascript.jscomp.ConformanceRules$BannedProperty", func(v_ any) bool {
		return j2gort.InstanceOf[*ConformanceRules_BannedProperty](v_)
	})
}
func (*ConformanceRules_BannedProperty) Implements_jscomp_ConformanceRules_AbstractRule() {}
func (*ConformanceRules_BannedProperty) Implements_jscomp_CheckConformance_Rule()         {}

func (this_ *ConformanceRules_BannedProperty) initFields_ConformanceRules_BannedProperty() {
	_ = this_
}

// Java constructor: BannedProperty(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.Requirement requirement)
func NewConformanceRules_BannedProperty(compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) *ConformanceRules_BannedProperty {
	this_ := &ConformanceRules_BannedProperty{}
	this_.BindSelf_(this_)
	InitConformanceRules_BannedProperty(this_, compiler, requirement)
	return this_
}

func InitConformanceRules_BannedProperty(this_ *ConformanceRules_BannedProperty, compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) {
	_ = compiler
	_ = requirement
	InitConformanceRules_AbstractRule(&this_.ConformanceRules_AbstractRuleBase, compiler, requirement)
	this_.initFields_ConformanceRules_BannedProperty()
	if requirement.GetValueCount() == int32(0) {
		panic(NewCheckConformance_InvalidRequirementSpec_1("missing value"))
	}
	this_.requirementPrecondition = func() *ConformanceRules_BannedProperty_RequirementPrecondition {
		switch j2gort.CheckNotNullPtr(requirement.GetType()).Ordinal() {
		case int32(6) /* Type.BANNED_PROPERTY */ :
			return StaticGet_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY()
		case int32(7) /* Type.BANNED_PROPERTY_READ */ :
			return StaticGet_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_READ()
		case int32(8) /* Type.BANNED_PROPERTY_WRITE */ :
			return StaticGet_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_WRITE()
		case int32(13) /* Type.BANNED_PROPERTY_NON_CONSTANT_WRITE */ :
			return StaticGet_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_NON_CONSTANT_WRITE()
		case int32(12) /* Type.BANNED_PROPERTY_CALL */ :
			return StaticGet_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_CALL()
		default:
			panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError", j2gort.PtrToAny(requirement.GetType()))))
		}
		panic("unreachable switch expression")
	}()
	this_.registry = compiler.GetTypeRegistry()
	var builder *j2gort.Ext_com_google_common_collect_ImmutableSetMultimap_Builder = j2gort.Ext_com_google_common_collect_ImmutableSetMultimap_S_Builder()
	_ = builder
	{
		cur_23 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](requirement.GetValueList())))
		for cur_23.HasNext() {
			value := cur_23.NextStr()
			_ = value
			var typename string = ConformanceRules_ConformanceUtil_GetClassFromDeclarationName(value)
			_ = typename
			var property string = ConformanceRules_ConformanceUtil_GetPropertyFromDeclarationName(value)
			_ = property
			if (typename == j2gort.NullString) || (property == j2gort.NullString) {
				panic(NewCheckConformance_InvalidRequirementSpec_1("bad prop value"))
			}
			var type_ pkg_rhino.JSType = this_.registry.GetGlobalType(typename)
			_ = type_
			if type_ == nil || type_.IsJavaNil_() {
				type_ = this_.registry.ResolveViaClosureNamespace(typename)
			}
			if ((type_ == nil || type_.IsJavaNil_()) || type_.IsUnknownType()) || type_.IsEmptyType() {
				continue
			}
			j2gort.MultimapBuilderPutStr(builder, property, type_)
		}
	}
	this_.props = j2gort.CastPtr[j2gort.Ext_com_google_common_collect_ImmutableSetMultimap](builder.Build())
}

// Java method: BannedProperty.getPrecondition()
func (this_ *ConformanceRules_BannedProperty) GetPrecondition() CheckConformance_Precondition {
	return ToIface_CheckConformance_Precondition(this_.requirementPrecondition)
}

// Java method: BannedProperty.checkConformance(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *ConformanceRules_BannedProperty) CheckConformance_2(t *NodeTraversal, n pkg_rhino.Node) *ConformanceRules_ConformanceResult {
	_ = t
	_ = n
	var checkTypes j2gort.Ext_com_google_common_collect_ImmutableSet = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](j2gort.MultimapGetStr(this_.props, this_.priv_BannedProperty_extractName_2151(n)))
	_ = checkTypes
	if checkTypes.IsEmpty() {
		return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
	}
	var foundType pkg_rhino.JSType = this_.priv_BannedProperty_extractType_2150(n)
	_ = foundType
	if foundType == nil || foundType.IsJavaNil_() {
		return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
	}
	foundType = foundType.RestrictByNotNullOrUndefined()
	{
		cur_24 := j2gort.NewIterCursor(j2gort.CheckNotNull(checkTypes))
		for cur_24.HasNext() {
			rawElem_25 := cur_24.Next()
			checkType := j2gort.AsIface[pkg_rhino.JSType](rawElem_25)
			_ = checkType
			var result *ConformanceRules_ConformanceResult = this_.priv_BannedProperty_matchTypes_2148(foundType, checkType)
			_ = result
			if result.level != StaticGet_ConformanceRules_ConformanceLevel_CONFORMANCE() {
				return result
			}
		}
	}
	return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
}

// Java method: BannedProperty.matchTypes(com.google.javascript.rhino.jstype.JSType foundType, com.google.javascript.rhino.jstype.JSType checkType)
func (this_ *ConformanceRules_BannedProperty) priv_BannedProperty_matchTypes_2148(foundType pkg_rhino.JSType, checkType pkg_rhino.JSType) *ConformanceRules_ConformanceResult {
	_ = foundType
	_ = checkType
	var foundObj pkg_rhino.ObjectType = foundType.ToMaybeObjectType()
	_ = foundObj
	if foundObj != nil && !foundObj.IsJavaNil_() {
		if foundObj.IsFunctionPrototypeType() {
			var ownerFun *pkg_rhino.FunctionType = foundObj.GetOwnerFunction()
			_ = ownerFun
			if ownerFun.IsConstructor() {
				foundType = j2gort.AsIface[pkg_rhino.JSType](ownerFun.GetInstanceType())
			}
		} else if foundObj.IsTemplatizedType() {
			foundType = j2gort.AsIface[pkg_rhino.JSType](foundObj.GetRawType())
		}
	}
	if (((foundType.IsUnknownType() || foundType.IsTemplateType()) || foundType.IsEmptyType()) || foundType.IsAllType()) || ConformanceRules_IsLooseObject(foundType, this_.registry) {
		if this_.GetConformanceRules_AbstractRuleBase().reportLooseTypeViolations {
			return StaticGet_ConformanceRules_ConformanceResult_POSSIBLE_VIOLATION_DUE_TO_LOOSE_TYPES()
		}
	} else if foundType.IsSubtypeOf_1(checkType) {
		return StaticGet_ConformanceRules_ConformanceResult_VIOLATION()
	} else if checkType.IsSubtypeWithoutStructuralTyping(foundType) {
		if this_.priv_BannedProperty_matchesPrototype_2149(checkType, foundType) {
			return StaticGet_ConformanceRules_ConformanceResult_VIOLATION()
		} else if this_.GetConformanceRules_AbstractRuleBase().reportLooseTypeViolations {
			return StaticGet_ConformanceRules_ConformanceResult_POSSIBLE_VIOLATION_DUE_TO_LOOSE_TYPES()
		}
	}
	return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
}

// Java method: BannedProperty.matchesPrototype(com.google.javascript.rhino.jstype.JSType type, com.google.javascript.rhino.jstype.JSType maybePrototype)
func (this_ *ConformanceRules_BannedProperty) priv_BannedProperty_matchesPrototype_2149(type_ pkg_rhino.JSType, maybePrototype pkg_rhino.JSType) bool {
	_ = type_
	_ = maybePrototype
	var methodClassObjectType pkg_rhino.ObjectType = type_.ToMaybeObjectType()
	_ = methodClassObjectType
	if methodClassObjectType != nil && !methodClassObjectType.IsJavaNil_() {
		if methodClassObjectType.GetImplicitPrototype().Equals(maybePrototype) {
			return true
		}
	}
	return false
}

// Java method: BannedProperty.extractType(com.google.javascript.rhino.Node n)
func (this_ *ConformanceRules_BannedProperty) priv_BannedProperty_extractType_2150(n pkg_rhino.Node) pkg_rhino.JSType {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(27) /* Token.GETELEM */, int32(26) /* Token.GETPROP */ :
		return n.GetFirstChild().GetJSType()
	case int32(100) /* Token.STRING_KEY */, int32(129) /* Token.COMPUTED_PROP */ :
		var parent pkg_rhino.Node = n.GetParent()
		_ = parent
		return func() pkg_rhino.JSType {
			switch parent.TokenSwitchOrd_() {
			case int32(103) /* Token.OBJECT_PATTERN */, int32(47) /* Token.OBJECTLIT */ :
				return parent.GetJSType()
			case int32(106) /* Token.CLASS_MEMBERS */ :
				return nil
			default:
				panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError")))
			}
			panic("unreachable switch expression")
		}()
	default:
		return nil
	}
	panic("unreachable switch expression")
}

// Java method: BannedProperty.extractName(com.google.javascript.rhino.Node n)
func (this_ *ConformanceRules_BannedProperty) priv_BannedProperty_extractName_2151(n pkg_rhino.Node) string {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(26) /* Token.GETPROP */, int32(100) /* Token.STRING_KEY */ :
		return n.GetString()
	case int32(27) /* Token.GETELEM */ :
		var string_ pkg_rhino.Node = n.GetSecondChild()
		_ = string_
		return func() string {
			if string_.IsStringLit() {
				return string_.GetString()
			}
			return j2gort.NullString
		}()
	case int32(129) /* Token.COMPUTED_PROP */ :
		var string__2 pkg_rhino.Node = n.GetFirstChild()
		_ = string__2
		return func() string {
			if string__2.IsStringLit() {
				return string__2.GetString()
			}
			return j2gort.NullString
		}()
	default:
		return j2gort.NullString
	}
	panic("unreachable switch expression")
}

// Java enum: com.google.javascript.jscomp.ConformanceRules.BannedProperty.RequirementPrecondition
type ConformanceRules_BannedProperty_RequirementPrecondition struct {
	j2gort.EnumBase
	fn_ShouldCheck func(*ConformanceRules_BannedProperty_RequirementPrecondition, pkg_rhino.Node) bool
}

func (this_ *ConformanceRules_BannedProperty_RequirementPrecondition) GetEnumBase() *j2gort.EnumBase {
	if this_ == nil {
		return nil
	}
	return &this_.EnumBase
}
func (this_ *ConformanceRules_BannedProperty_RequirementPrecondition) EnumOrdinal_() int32 {
	if this_ == nil {
		return -1
	}
	return this_.EnumBase.Ordinal_
}
func (*ConformanceRules_BannedProperty_RequirementPrecondition) Implements_j2gort_Ext_java_lang_Enum() {
}
func (*ConformanceRules_BannedProperty_RequirementPrecondition) Implements_jscomp_CheckConformance_Precondition() {
}
func (*ConformanceRules_BannedProperty_RequirementPrecondition) Implements_j2gort_Ext_java_lang_constant_Constable() {
}
func (this_ *ConformanceRules_BannedProperty_RequirementPrecondition) initFields_ConformanceRules_BannedProperty_RequirementPrecondition() {
	_ = this_
}

func initEnumConformanceRules_BannedProperty_RequirementPrecondition(this_ *ConformanceRules_BannedProperty_RequirementPrecondition, enumName_ string, enumOrd_ int32) *ConformanceRules_BannedProperty_RequirementPrecondition {
	this_.EnumBase.Name_ = enumName_
	this_.EnumBase.Ordinal_ = enumOrd_
	this_.initFields_ConformanceRules_BannedProperty_RequirementPrecondition()
	return this_
}

var enum_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY *ConformanceRules_BannedProperty_RequirementPrecondition

func StaticGet_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY() *ConformanceRules_BannedProperty_RequirementPrecondition {
	clinit_ConformanceRules_BannedProperty_RequirementPrecondition()
	return enum_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY
}

var enum_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_WRITE *ConformanceRules_BannedProperty_RequirementPrecondition

func StaticGet_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_WRITE() *ConformanceRules_BannedProperty_RequirementPrecondition {
	clinit_ConformanceRules_BannedProperty_RequirementPrecondition()
	return enum_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_WRITE
}

var enum_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_NON_CONSTANT_WRITE *ConformanceRules_BannedProperty_RequirementPrecondition

func StaticGet_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_NON_CONSTANT_WRITE() *ConformanceRules_BannedProperty_RequirementPrecondition {
	clinit_ConformanceRules_BannedProperty_RequirementPrecondition()
	return enum_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_NON_CONSTANT_WRITE
}

var enum_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_READ *ConformanceRules_BannedProperty_RequirementPrecondition

func StaticGet_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_READ() *ConformanceRules_BannedProperty_RequirementPrecondition {
	clinit_ConformanceRules_BannedProperty_RequirementPrecondition()
	return enum_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_READ
}

var enum_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_CALL *ConformanceRules_BannedProperty_RequirementPrecondition

func StaticGet_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_CALL() *ConformanceRules_BannedProperty_RequirementPrecondition {
	clinit_ConformanceRules_BannedProperty_RequirementPrecondition()
	return enum_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_CALL
}

var clinitDone_ConformanceRules_BannedProperty_RequirementPrecondition uint32
var clinitState_ConformanceRules_BannedProperty_RequirementPrecondition j2gort.ClinitState

func clinit_ConformanceRules_BannedProperty_RequirementPrecondition() {
	if atomic.LoadUint32(&clinitDone_ConformanceRules_BannedProperty_RequirementPrecondition) != 0 {
		return
	}
	clinitSlow_ConformanceRules_BannedProperty_RequirementPrecondition()
}

func clinitSlow_ConformanceRules_BannedProperty_RequirementPrecondition() {
	if !j2gort.EnterClinit(&clinitState_ConformanceRules_BannedProperty_RequirementPrecondition, &clinitDone_ConformanceRules_BannedProperty_RequirementPrecondition) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ConformanceRules_BannedProperty_RequirementPrecondition, &clinitDone_ConformanceRules_BannedProperty_RequirementPrecondition)
	enum_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY = &ConformanceRules_BannedProperty_RequirementPrecondition{}
	enum_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY.fn_ShouldCheck = func(this_ *ConformanceRules_BannedProperty_RequirementPrecondition, n pkg_rhino.Node) bool {
		_ = n
		switch n.TokenSwitchOrd_() {
		case int32(100) /* Token.STRING_KEY */, int32(26) /* Token.GETPROP */, int32(27) /* Token.GETELEM */, int32(129) /* Token.COMPUTED_PROP */ :
			return true
		default:
			return false
		}
		panic("unreachable switch expression")
	}
	initEnumConformanceRules_BannedProperty_RequirementPrecondition(enum_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY, "BANNED_PROPERTY", int32(0))
	enum_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_WRITE = &ConformanceRules_BannedProperty_RequirementPrecondition{}
	enum_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_WRITE.fn_ShouldCheck = func(this_ *ConformanceRules_BannedProperty_RequirementPrecondition, n pkg_rhino.Node) bool {
		_ = n
		return NodeUtil_IsLValue(n)
	}
	initEnumConformanceRules_BannedProperty_RequirementPrecondition(enum_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_WRITE, "BANNED_PROPERTY_WRITE", int32(1))
	enum_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_NON_CONSTANT_WRITE = &ConformanceRules_BannedProperty_RequirementPrecondition{}
	enum_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_NON_CONSTANT_WRITE.fn_ShouldCheck = func(this_ *ConformanceRules_BannedProperty_RequirementPrecondition, n pkg_rhino.Node) bool {
		_ = n
		if !(NodeUtil_IsLValue(n)) {
			return false
		}
		if NodeUtil_IsLhsOfAssign(n) && (NodeUtil_IsLiteralValue(n.GetNext(), false) || NodeUtil_IsSomeCompileTimeConstStringValue(n.GetNext())) {
			return false
		}
		return true
	}
	initEnumConformanceRules_BannedProperty_RequirementPrecondition(enum_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_NON_CONSTANT_WRITE, "BANNED_PROPERTY_NON_CONSTANT_WRITE", int32(2))
	enum_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_READ = &ConformanceRules_BannedProperty_RequirementPrecondition{}
	enum_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_READ.fn_ShouldCheck = func(this_ *ConformanceRules_BannedProperty_RequirementPrecondition, n pkg_rhino.Node) bool {
		_ = n
		return (!(NodeUtil_IsLValue(n)) && NodeUtil_IsExpressionResultUsed(n))
	}
	initEnumConformanceRules_BannedProperty_RequirementPrecondition(enum_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_READ, "BANNED_PROPERTY_READ", int32(3))
	enum_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_CALL = &ConformanceRules_BannedProperty_RequirementPrecondition{}
	enum_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_CALL.fn_ShouldCheck = func(this_ *ConformanceRules_BannedProperty_RequirementPrecondition, n pkg_rhino.Node) bool {
		_ = n
		return ConformanceRules_ConformanceUtil_IsCallTarget(n)
	}
	initEnumConformanceRules_BannedProperty_RequirementPrecondition(enum_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_CALL, "BANNED_PROPERTY_CALL", int32(4))
	enumValues_ConformanceRules_BannedProperty_RequirementPrecondition = []*ConformanceRules_BannedProperty_RequirementPrecondition{enum_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY, enum_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_WRITE, enum_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_NON_CONSTANT_WRITE, enum_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_READ, enum_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_CALL}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.ConformanceRules.BannedProperty.RequirementPrecondition", []any{enum_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY, enum_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_WRITE, enum_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_NON_CONSTANT_WRITE, enum_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_READ, enum_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_CALL})
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.ConformanceRules$BannedProperty$RequirementPrecondition", []any{enum_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY, enum_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_WRITE, enum_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_NON_CONSTANT_WRITE, enum_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_READ, enum_ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_CALL})
}

func init() {
	j2gort.RegisterEnumConstantsLoader("com.google.javascript.jscomp.ConformanceRules.BannedProperty.RequirementPrecondition", "com.google.javascript.jscomp.ConformanceRules$BannedProperty$RequirementPrecondition", func() {
		clinit_ConformanceRules_BannedProperty_RequirementPrecondition()
	})
}

var enumValues_ConformanceRules_BannedProperty_RequirementPrecondition []*ConformanceRules_BannedProperty_RequirementPrecondition

func ConformanceRules_BannedProperty_RequirementPrecondition_Values() []*ConformanceRules_BannedProperty_RequirementPrecondition {
	clinit_ConformanceRules_BannedProperty_RequirementPrecondition()
	return j2gort.CloneSlice(enumValues_ConformanceRules_BannedProperty_RequirementPrecondition)
}

func ConformanceRules_BannedProperty_RequirementPrecondition_ValueOf(name string) *ConformanceRules_BannedProperty_RequirementPrecondition {
	clinit_ConformanceRules_BannedProperty_RequirementPrecondition()
	for _, c := range enumValues_ConformanceRules_BannedProperty_RequirementPrecondition {
		if c.Name() == name {
			return c
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.ConformanceRules.BannedProperty.RequirementPrecondition." + name))
}

func (this_ *ConformanceRules_BannedProperty_RequirementPrecondition) ShouldCheck(p0_ pkg_rhino.Node) bool {
	if this_.fn_ShouldCheck != nil {
		return this_.fn_ShouldCheck(this_, p0_)
	}
	return this_.baseEnum_ShouldCheck(p0_)
}

func (this_ *ConformanceRules_BannedProperty_RequirementPrecondition) baseEnum_ShouldCheck(p0_ pkg_rhino.Node) bool {
	_ = p0_
	panic("abstract or native method invoked")
}

func (this_ *ConformanceRules_BannedProperty_RequirementPrecondition) IsJavaNil_() bool {
	return this_ == nil
}
func (*ConformanceRules_BannedProperty_RequirementPrecondition) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConformanceRules$BannedProperty$RequirementPrecondition", "RequirementPrecondition"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConformanceRules.BannedProperty.RequirementPrecondition", "com.google.javascript.jscomp.ConformanceRules$BannedProperty$RequirementPrecondition", func(v_ any) bool {
		return j2gort.InstanceOf[*ConformanceRules_BannedProperty_RequirementPrecondition](v_)
	})
}

// Java class: com.google.javascript.jscomp.ConformanceRules.BannedStringRegex
type ConformanceRules_BannedStringRegex struct {
	ConformanceRules_AbstractRuleBase
	stringPattern *j2gort.Ext_java_util_regex_Pattern `gson:"stringPattern"`
}

func (this_ *ConformanceRules_BannedStringRegex) BindSelf_(s any) {
	_ = s
	this_.ConformanceRules_AbstractRuleBase.BindSelf_(s)
}
func (this_ *ConformanceRules_BannedStringRegex) IsJavaNil_() bool { return this_ == nil }
func (*ConformanceRules_BannedStringRegex) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConformanceRules$BannedStringRegex", "BannedStringRegex"
}
func (this_ *ConformanceRules_BannedStringRegex) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ConformanceRules_AbstractRuleBase.IdentityHash_()
}
func (*ConformanceRules_BannedStringRegex) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConformanceRules.BannedStringRegex", "com.google.javascript.jscomp.ConformanceRules$BannedStringRegex", func(v_ any) bool {
		return j2gort.InstanceOf[*ConformanceRules_BannedStringRegex](v_)
	})
}
func (*ConformanceRules_BannedStringRegex) Implements_jscomp_ConformanceRules_AbstractRule() {}
func (*ConformanceRules_BannedStringRegex) Implements_jscomp_CheckConformance_Rule()         {}

func (this_ *ConformanceRules_BannedStringRegex) initFields_ConformanceRules_BannedStringRegex() {
	_ = this_
}

// Java constructor: BannedStringRegex(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.Requirement requirement)
func NewConformanceRules_BannedStringRegex(compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) *ConformanceRules_BannedStringRegex {
	this_ := &ConformanceRules_BannedStringRegex{}
	this_.BindSelf_(this_)
	InitConformanceRules_BannedStringRegex(this_, compiler, requirement)
	return this_
}

func InitConformanceRules_BannedStringRegex(this_ *ConformanceRules_BannedStringRegex, compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) {
	_ = compiler
	_ = requirement
	InitConformanceRules_AbstractRule(&this_.ConformanceRules_AbstractRuleBase, compiler, requirement)
	this_.initFields_ConformanceRules_BannedStringRegex()
	if requirement.GetValueCount() == int32(0) {
		panic(NewCheckConformance_InvalidRequirementSpec_1("missing value"))
	}
	var builder []string = []string{}
	_ = builder
	{
		cur_26 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](requirement.GetValueList())))
		for cur_26.HasNext() {
			value := cur_26.NextStr()
			_ = value
			if j2gort.CheckNotNullString(javalang.StringTrim(j2gort.CheckNotNullString(value))) == "" {
				panic(NewCheckConformance_InvalidRequirementSpec_1("empty strings or whitespace are not allowed"))
			}
			builder = append(builder, j2gort.CheckNotNullString(value))
		}
	}
	var values []string = builder
	_ = values
	if len(values) == 0 {
		panic(NewCheckConformance_InvalidRequirementSpec_1("missing value"))
	}
	var stringRegex *j2gort.Ext_java_util_regex_Pattern = ConformanceRules_BuildPattern(j2gort.ImmListFromTypedSlice(values))
	_ = stringRegex
	this_.stringPattern = stringRegex
}

// Java method: BannedStringRegex.checkConformance(com.google.javascript.jscomp.NodeTraversal traversal, com.google.javascript.rhino.Node node)
func (this_ *ConformanceRules_BannedStringRegex) CheckConformance_2(traversal *NodeTraversal, node pkg_rhino.Node) *ConformanceRules_ConformanceResult {
	_ = traversal
	_ = node
	if node == 0 {
		return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
	}
	if node.IsStringLit() {
		if this_.stringPattern.Matcher(j2gort.BoxString(node.GetString())).Matches() {
			return StaticGet_ConformanceRules_ConformanceResult_VIOLATION()
		}
	} else if node.IsTemplateLitString() {
		var cookedString string = node.GetCookedString()
		_ = cookedString
		if cookedString != j2gort.NullString {
			if this_.stringPattern.Matcher(j2gort.BoxString(cookedString)).Matches() {
				return StaticGet_ConformanceRules_ConformanceResult_VIOLATION()
			}
		} else {
			if !(ConformanceRules_IsTaggedTemplateLitString(node)) {
				j2gort.CheckStateFmt(false, "Found untagged template literal with invalid (uncookable) escape sequence: %s", j2gort.BoxString(node.GetRawString()))
			}
			var rawString string = node.GetRawString()
			_ = rawString
			if this_.stringPattern.Matcher(j2gort.BoxString(rawString)).Matches() {
				return StaticGet_ConformanceRules_ConformanceResult_VIOLATION()
			}
		}
	}
	return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
}

// Java class: com.google.javascript.jscomp.ConformanceRules.ConformanceUtil
type ConformanceRules_ConformanceUtil struct {
	idHash_ atomic.Int32
}

func (this_ *ConformanceRules_ConformanceUtil) BindSelf_(s any) {
	_ = s
}
func (this_ *ConformanceRules_ConformanceUtil) IsJavaNil_() bool { return this_ == nil }
func (*ConformanceRules_ConformanceUtil) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConformanceRules$ConformanceUtil", "ConformanceUtil"
}
func (this_ *ConformanceRules_ConformanceUtil) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ConformanceRules_ConformanceUtil) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConformanceRules.ConformanceUtil", "com.google.javascript.jscomp.ConformanceRules$ConformanceUtil", func(v_ any) bool {
		return j2gort.InstanceOf[*ConformanceRules_ConformanceUtil](v_)
	})
}

func (this_ *ConformanceRules_ConformanceUtil) initFields_ConformanceRules_ConformanceUtil() {
	_ = this_
}

// Java constructor: ConformanceUtil()
func NewConformanceRules_ConformanceUtil() *ConformanceRules_ConformanceUtil {
	this_ := &ConformanceRules_ConformanceUtil{}
	InitConformanceRules_ConformanceUtil(this_)
	return this_
}

func InitConformanceRules_ConformanceUtil(this_ *ConformanceRules_ConformanceUtil) {
	this_.initFields_ConformanceRules_ConformanceUtil()
}

// Java static method: ConformanceUtil.isCallTarget(com.google.javascript.rhino.Node n)
func ConformanceRules_ConformanceUtil_IsCallTarget(n pkg_rhino.Node) bool {
	_ = n
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	return ((parent.IsCall() || parent.IsNew()) && (parent.GetFirstChild() == n))
}

// Java static method: ConformanceUtil.isLooseType(com.google.javascript.rhino.jstype.JSType type)
func ConformanceRules_ConformanceUtil_IsLooseType(type_ pkg_rhino.JSType) bool {
	_ = type_
	return ((type_.IsUnknownType() || type_.IsNoResolvedType()) || type_.IsAllType())
}

// Java static method: ConformanceUtil.evaluateTypeString(com.google.javascript.jscomp.AbstractCompiler compiler, java.lang.String expression)
func ConformanceRules_ConformanceUtil_EvaluateTypeString(compiler AbstractCompiler, expression string) pkg_rhino.JSType {
	_ = compiler
	_ = expression
	var typeNodes pkg_rhino.Node = pkg_parsing.JsDocInfoParser_ParseTypeString(expression)
	_ = typeNodes
	if typeNodes == 0 {
		panic(NewCheckConformance_InvalidRequirementSpec_1("bad type expression"))
	}
	var typeExpr pkg_rhino.JSTypeExpression = pkg_rhino.NewJSTypeExpression(typeNodes, "conformance")
	_ = typeExpr
	return compiler.GetTypeRegistry().EvaluateTypeExpressionInGlobalScope(typeExpr)
}

// Java static method: ConformanceUtil.validateCall(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.rhino.Node callOrNew, com.google.javascript.rhino.jstype.FunctionType functionType, boolean isCallInvocation)
func ConformanceRules_ConformanceUtil_ValidateCall(compiler AbstractCompiler, callOrNew pkg_rhino.Node, functionType *pkg_rhino.FunctionType, isCallInvocation bool) bool {
	_ = compiler
	_ = callOrNew
	_ = functionType
	_ = isCallInvocation
	if !(callOrNew.IsCall() || callOrNew.IsNew()) {
		j2gort.CheckState(false)
	}
	return (ConformanceRules_ConformanceUtil_ValidateParameterList(compiler, callOrNew, functionType, isCallInvocation) && ConformanceRules_ConformanceUtil_ValidateThis(callOrNew, functionType, isCallInvocation))
}

// Java static method: ConformanceUtil.validateThis(com.google.javascript.rhino.Node callOrNew, com.google.javascript.rhino.jstype.FunctionType functionType, boolean isCallInvocation)
func ConformanceRules_ConformanceUtil_ValidateThis(callOrNew pkg_rhino.Node, functionType *pkg_rhino.FunctionType, isCallInvocation bool) bool {
	_ = callOrNew
	_ = functionType
	_ = isCallInvocation
	if callOrNew.IsNew() {
		return true
	}
	var thisType pkg_rhino.JSType = functionType.GetTypeOfThis()
	_ = thisType
	if (thisType == nil || thisType.IsJavaNil_()) || thisType.IsUnknownType() {
		return true
	}
	var thisNode pkg_rhino.Node
	_ = thisNode
	if isCallInvocation {
		thisNode = callOrNew.GetSecondChild()
	} else {
		thisNode = callOrNew.GetFirstFirstChild()
	}
	var thisNodeType pkg_rhino.JSType = thisNode.GetJSType().RestrictByNotNullOrUndefined()
	_ = thisNodeType
	return thisNodeType.IsSubtypeOf_1(thisType)
}

// Java static method: ConformanceUtil.validateParameterList(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.rhino.Node callOrNew, com.google.javascript.rhino.jstype.FunctionType functionType, boolean isCallInvocation)
func ConformanceRules_ConformanceUtil_ValidateParameterList(compiler AbstractCompiler, callOrNew pkg_rhino.Node, functionType *pkg_rhino.FunctionType, isCallInvocation bool) bool {
	_ = compiler
	_ = callOrNew
	_ = functionType
	_ = isCallInvocation
	var arg pkg_rhino.Node = callOrNew.GetSecondChild()
	_ = arg
	if isCallInvocation && (arg != 0) {
		arg = arg.GetNext()
	}
	var argumentTypes []pkg_rhino.JSType = []pkg_rhino.JSType{}
	_ = argumentTypes
	for ; arg != 0; arg = arg.GetNext() {
		var argType pkg_rhino.JSType = arg.GetJSType()
		_ = argType
		if argType == nil || argType.IsJavaNil_() {
			argType = compiler.GetTypeRegistry().GetNativeType(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())
		}
		argumentTypes = append(argumentTypes, j2gort.CheckNotNullElem(argType))
	}
	return functionType.AcceptsArguments(j2gort.ImmListFromTypedSlice(argumentTypes))
}

// Java static method: ConformanceUtil.getPropertyFromDeclarationName(java.lang.String specName)
func ConformanceRules_ConformanceUtil_GetPropertyFromDeclarationName(specName string) string {
	_ = specName
	var parts j2gort.Ext_java_util_List = j2gort.SplitterSplitToListStr(StaticGet_ConformanceRules_ON_PROTOTYPE(), specName)
	_ = parts
	if !((parts.Size() == int32(1)) || (parts.Size() == int32(2))) {
		j2gort.CheckState(false)
	}
	if parts.Size() == int32(2) {
		return j2gort.ListGetStr(parts, int32(1))
	}
	return j2gort.NullString
}

// Java static method: ConformanceUtil.getClassFromDeclarationName(java.lang.String specName)
func ConformanceRules_ConformanceUtil_GetClassFromDeclarationName(specName string) string {
	_ = specName
	var tmp string = specName
	_ = tmp
	var parts j2gort.Ext_java_util_List = j2gort.SplitterSplitToListStr(StaticGet_ConformanceRules_ON_PROTOTYPE(), tmp)
	_ = parts
	if !((parts.Size() == int32(1)) || (parts.Size() == int32(2))) {
		j2gort.CheckState(false)
	}
	if parts.Size() == int32(2) {
		return j2gort.ListGetStr(parts, int32(0))
	}
	return j2gort.NullString
}

// Java static method: ConformanceUtil.removeTypeDecl(java.lang.String specName)
func ConformanceRules_ConformanceUtil_RemoveTypeDecl(specName string) string {
	_ = specName
	var index int32 = j2gort.JStrIndexOfCp(specName, int32(uint16(58)))
	_ = index
	if index < int32(1) {
		panic(NewCheckConformance_InvalidRequirementSpec_1("value should be in the form NAME:TYPE"))
	}
	return j2gort.JStrSub(specName, int32(0), index)
}

// Java static method: ConformanceUtil.getTypeFromValue(java.lang.String specName)
func ConformanceRules_ConformanceUtil_GetTypeFromValue(specName string) string {
	_ = specName
	var index int32 = j2gort.JStrIndexOfCp(specName, int32(uint16(58)))
	_ = index
	if index < int32(1) {
		return j2gort.NullString
	}
	return j2gort.JStrSubFrom(specName, (index + int32(1)))
}

// Java static method: ConformanceUtil.isAnnotatedAsConst(com.google.javascript.rhino.JSDocInfo jsdoc)
func ConformanceRules_ConformanceUtil_IsAnnotatedAsConst(jsdoc *pkg_rhino.JSDocInfo) bool {
	_ = jsdoc
	return ((jsdoc != nil) && jsdoc.IsConstant())
}

// Java static method: ConformanceUtil.getValueFromGlobalName(com.google.javascript.jscomp.@org.jspecify.annotations.Nullable Scope scope, com.google.javascript.jscomp.GlobalNamespace gns, java.lang.String qualifiedName)
func ConformanceRules_ConformanceUtil_GetValueFromGlobalName(scope *Scope, gns *GlobalNamespace, qualifiedName string) string {
	_ = scope
	_ = gns
	_ = qualifiedName
	var gn *GlobalNamespace_Name = j2gort.CastPtr[GlobalNamespace_Name](gns.GetOwnSlot(qualifiedName))
	_ = gn
	if (((gn == nil) || (gn.GetGlobalSets() != int32(1))) || (gn.GetTotalSets() != int32(1))) || !(ConformanceRules_ConformanceUtil_IsAnnotatedAsConst(gn.GetJSDocInfo())) {
		return j2gort.NullString
	}
	return ConformanceRules_ConformanceUtil_InferStringValue(scope, NodeUtil_GetRValueOfLValue(j2gort.CastPtr[GlobalNamespace_Ref](gn.GetDeclaration()).GetNode()), j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
		return j2gort.PtrToAny(gns)
	}))
}

// Java static method: ConformanceUtil.inferStringValue(com.google.javascript.jscomp.@org.jspecify.annotations.Nullable Scope scope, com.google.javascript.rhino.Node node, java.util.function.Supplier<com.google.javascript.jscomp.GlobalNamespace> globalNamespaceSupplier)
func ConformanceRules_ConformanceUtil_InferStringValue(scope *Scope, node pkg_rhino.Node, globalNamespaceSupplier j2gort.Ext_java_util_function_Supplier) string {
	_ = scope
	_ = node
	_ = globalNamespaceSupplier
	if node == 0 {
		return j2gort.NullString
	}
	switch node.TokenSwitchOrd_() {
	case int32(100) /* Token.STRING_KEY */, int32(35) /* Token.STRINGLIT */, int32(131) /* Token.TEMPLATELIT */, int32(133) /* Token.TEMPLATELIT_STRING */ :
		return NodeUtil_GetStringValue(node)
	case int32(32) /* Token.NAME */ :
		if scope == nil {
			return j2gort.NullString
		}
		var name string = node.GetString()
		_ = name
		var var_ *Var = j2gort.CastPtr[Var](scope.GetVar(name))
		_ = var_
		if var_ == nil {
			return j2gort.NullString
		}
		if !(var_.IsConst()) && !(ConformanceRules_ConformanceUtil_IsAnnotatedAsConst(var_.GetJSDocInfo())) {
			return j2gort.NullString
		}
		var initialValue pkg_rhino.Node = var_.GetInitialValue()
		_ = initialValue
		return ConformanceRules_ConformanceUtil_InferStringValue(j2gort.CastPtr[Scope](var_.GetScope()), initialValue, globalNamespaceSupplier)
	case int32(26) /* Token.GETPROP */ :
		var type_ pkg_rhino.JSType = node.GetJSType()
		_ = type_
		if type_ == nil || type_.IsJavaNil_() {
			return j2gort.NullString
		}
		if type_.IsEnumElementType() {
			var enumSource pkg_rhino.Node = type_.ToMaybeEnumElementType().GetEnumType().GetSource()
			_ = enumSource
			if enumSource == 0 {
				return j2gort.NullString
			}
			if !(enumSource.IsObjectLit()) && !(enumSource.IsClassMembers()) {
				return j2gort.NullString
			}
			return ConformanceRules_ConformanceUtil_InferStringValue(nil, NodeUtil_GetFirstPropMatchingKey(enumSource, node.GetString()), globalNamespaceSupplier)
		} else if type_.IsString() {
			var gns *GlobalNamespace = j2gort.CastPtr[GlobalNamespace](globalNamespaceSupplier.Get())
			_ = gns
			if gns == nil {
				return j2gort.NullString
			}
			return ConformanceRules_ConformanceUtil_GetValueFromGlobalName(scope, gns, node.GetQualifiedName())
		}
		return j2gort.NullString
	default:
	}
	return j2gort.NullString
}

// Java static method: ConformanceUtil.isXid(com.google.javascript.rhino.jstype.JSType type)
func ConformanceRules_ConformanceUtil_IsXid(type_ pkg_rhino.JSType) bool {
	_ = type_
	if type_ == nil || type_.IsJavaNil_() {
		return false
	}
	var enumElTy *pkg_rhino.EnumElementType = type_.ToMaybeEnumElementType()
	_ = enumElTy
	if (enumElTy != nil) && (j2gort.CheckNotNullString(enumElTy.GetEnumType().GetReferenceName()) == "enum{xid.String}") {
		return true
	}
	return false
}

// Java static method: ConformanceUtil.isEventHandlerAttrName(java.lang.String attr)
func ConformanceRules_ConformanceUtil_IsEventHandlerAttrName(attr string) bool {
	_ = attr
	return (!(j2gort.CheckNotNullString(attr) == "on") && javalang.StringStartsWith(j2gort.CheckNotNullString(attr), "on"))
}

// Java class: com.google.javascript.jscomp.ConformanceRules.RestrictedNameCall
type ConformanceRules_RestrictedNameCall struct {
	ConformanceRules_AbstractRuleBase
	restrictions []*ConformanceRules_RestrictedNameCall_Restriction `gson:"restrictions"`
}

func (this_ *ConformanceRules_RestrictedNameCall) BindSelf_(s any) {
	_ = s
	this_.ConformanceRules_AbstractRuleBase.BindSelf_(s)
}
func (this_ *ConformanceRules_RestrictedNameCall) IsJavaNil_() bool { return this_ == nil }
func (*ConformanceRules_RestrictedNameCall) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConformanceRules$RestrictedNameCall", "RestrictedNameCall"
}
func (this_ *ConformanceRules_RestrictedNameCall) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ConformanceRules_AbstractRuleBase.IdentityHash_()
}
func (*ConformanceRules_RestrictedNameCall) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConformanceRules.RestrictedNameCall", "com.google.javascript.jscomp.ConformanceRules$RestrictedNameCall", func(v_ any) bool {
		return j2gort.InstanceOf[*ConformanceRules_RestrictedNameCall](v_)
	})
}
func (*ConformanceRules_RestrictedNameCall) Implements_jscomp_ConformanceRules_AbstractRule() {}
func (*ConformanceRules_RestrictedNameCall) Implements_jscomp_CheckConformance_Rule()         {}

func (this_ *ConformanceRules_RestrictedNameCall) initFields_ConformanceRules_RestrictedNameCall() {
	_ = this_
}

// Java constructor: RestrictedNameCall(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.Requirement requirement)
func NewConformanceRules_RestrictedNameCall(compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) *ConformanceRules_RestrictedNameCall {
	this_ := &ConformanceRules_RestrictedNameCall{}
	this_.BindSelf_(this_)
	InitConformanceRules_RestrictedNameCall(this_, compiler, requirement)
	return this_
}

func InitConformanceRules_RestrictedNameCall(this_ *ConformanceRules_RestrictedNameCall, compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) {
	_ = compiler
	_ = requirement
	InitConformanceRules_AbstractRule(&this_.ConformanceRules_AbstractRuleBase, compiler, requirement)
	this_.initFields_ConformanceRules_RestrictedNameCall()
	if requirement.GetValueCount() == int32(0) {
		panic(NewCheckConformance_InvalidRequirementSpec_1("missing value"))
	}
	var builder []*ConformanceRules_RestrictedNameCall_Restriction = []*ConformanceRules_RestrictedNameCall_Restriction{}
	_ = builder
	{
		cur_27 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](requirement.GetValueList())))
		for cur_27.HasNext() {
			value := cur_27.NextStr()
			_ = value
			var name pkg_rhino.Node = NodeUtil_NewQName_2_AbstractCompiler_string(compiler, ConformanceRules_RestrictedNameCall_GetNameFromValue(value))
			_ = name
			var restrictedDecl string = ConformanceRules_ConformanceUtil_GetTypeFromValue(value)
			_ = restrictedDecl
			if (name == 0) || (restrictedDecl == j2gort.NullString) {
				panic(NewCheckConformance_InvalidRequirementSpec_1("bad prop value"))
			}
			var restrictedCallType *pkg_rhino.FunctionType = ConformanceRules_ConformanceUtil_EvaluateTypeString(compiler, restrictedDecl).ToMaybeFunctionType()
			_ = restrictedCallType
			if restrictedCallType == nil {
				panic(NewCheckConformance_InvalidRequirementSpec_1("invalid conformance type"))
			}
			builder = append(builder, NewConformanceRules_RestrictedNameCall_Restriction(name, restrictedCallType))
		}
	}
	this_.restrictions = builder
}

// Java method: RestrictedNameCall.checkConformance(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *ConformanceRules_RestrictedNameCall) CheckConformance_2(t *NodeTraversal, n pkg_rhino.Node) *ConformanceRules_ConformanceResult {
	_ = t
	_ = n
	if ConformanceRules_ConformanceUtil_IsCallTarget(n) && n.IsQualifiedName() {
		{
			var i int32 = int32(0)
			_ = i
			for ; i < int32(len(this_.restrictions)); i = i + 1 {
				var r *ConformanceRules_RestrictedNameCall_Restriction = this_.restrictions[i]
				_ = r
				if n.MatchesQualifiedName_1_Node(r.name) {
					if !(ConformanceRules_ConformanceUtil_ValidateCall(this_.GetConformanceRules_AbstractRuleBase().compiler, n.GetParent(), r.restrictedCallType, false)) {
						return StaticGet_ConformanceRules_ConformanceResult_VIOLATION()
					}
				} else if (n.IsGetProp() && (j2gort.CheckNotNullString(n.GetString()) == "call")) && n.GetFirstChild().MatchesQualifiedName_1_Node(r.name) {
					if !(ConformanceRules_ConformanceUtil_ValidateCall(this_.GetConformanceRules_AbstractRuleBase().compiler, n.GetParent(), r.restrictedCallType, true)) {
						return StaticGet_ConformanceRules_ConformanceResult_VIOLATION()
					}
				}
			}
		}
	}
	return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
}

// Java static method: RestrictedNameCall.getNameFromValue(java.lang.String specName)
func ConformanceRules_RestrictedNameCall_GetNameFromValue(specName string) string {
	_ = specName
	var index int32 = j2gort.JStrIndexOfCp(specName, int32(uint16(58)))
	_ = index
	if index < int32(1) {
		return j2gort.NullString
	}
	return j2gort.JStrSub(specName, int32(0), index)
}

// Java class: com.google.javascript.jscomp.ConformanceRules.RestrictedNameCall.Restriction
type ConformanceRules_RestrictedNameCall_Restriction struct {
	restrictedCallType *pkg_rhino.FunctionType `gson:"restrictedCallType"`
	name               pkg_rhino.Node          `gson:"name"`
	idHash_            atomic.Int32
}

func (this_ *ConformanceRules_RestrictedNameCall_Restriction) BindSelf_(s any) {
	_ = s
}
func (this_ *ConformanceRules_RestrictedNameCall_Restriction) IsJavaNil_() bool { return this_ == nil }
func (*ConformanceRules_RestrictedNameCall_Restriction) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConformanceRules$RestrictedNameCall$Restriction", "Restriction"
}
func (this_ *ConformanceRules_RestrictedNameCall_Restriction) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ConformanceRules_RestrictedNameCall_Restriction) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConformanceRules.RestrictedNameCall.Restriction", "com.google.javascript.jscomp.ConformanceRules$RestrictedNameCall$Restriction", func(v_ any) bool {
		return j2gort.InstanceOf[*ConformanceRules_RestrictedNameCall_Restriction](v_)
	})
}

func (this_ *ConformanceRules_RestrictedNameCall_Restriction) initFields_ConformanceRules_RestrictedNameCall_Restriction() {
	_ = this_
}

// Java constructor: Restriction(com.google.javascript.rhino.Node name, com.google.javascript.rhino.jstype.FunctionType restrictedCallType)
func NewConformanceRules_RestrictedNameCall_Restriction(name pkg_rhino.Node, restrictedCallType *pkg_rhino.FunctionType) *ConformanceRules_RestrictedNameCall_Restriction {
	this_ := &ConformanceRules_RestrictedNameCall_Restriction{}
	InitConformanceRules_RestrictedNameCall_Restriction(this_, name, restrictedCallType)
	return this_
}

func InitConformanceRules_RestrictedNameCall_Restriction(this_ *ConformanceRules_RestrictedNameCall_Restriction, name pkg_rhino.Node, restrictedCallType *pkg_rhino.FunctionType) {
	_ = name
	_ = restrictedCallType
	this_.initFields_ConformanceRules_RestrictedNameCall_Restriction()
	this_.name = name
	this_.restrictedCallType = restrictedCallType
}

// Java class: com.google.javascript.jscomp.ConformanceRules.RestrictedMethodCall
type ConformanceRules_RestrictedMethodCall struct {
	ConformanceRules_AbstractRuleBase
	restrictions []*ConformanceRules_RestrictedMethodCall_Restriction `gson:"restrictions"`
}

func (this_ *ConformanceRules_RestrictedMethodCall) BindSelf_(s any) {
	_ = s
	this_.ConformanceRules_AbstractRuleBase.BindSelf_(s)
}
func (this_ *ConformanceRules_RestrictedMethodCall) IsJavaNil_() bool { return this_ == nil }
func (*ConformanceRules_RestrictedMethodCall) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConformanceRules$RestrictedMethodCall", "RestrictedMethodCall"
}
func (this_ *ConformanceRules_RestrictedMethodCall) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ConformanceRules_AbstractRuleBase.IdentityHash_()
}
func (*ConformanceRules_RestrictedMethodCall) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConformanceRules.RestrictedMethodCall", "com.google.javascript.jscomp.ConformanceRules$RestrictedMethodCall", func(v_ any) bool {
		return j2gort.InstanceOf[*ConformanceRules_RestrictedMethodCall](v_)
	})
}
func (*ConformanceRules_RestrictedMethodCall) Implements_jscomp_ConformanceRules_AbstractRule() {}
func (*ConformanceRules_RestrictedMethodCall) Implements_jscomp_CheckConformance_Rule()         {}

func (this_ *ConformanceRules_RestrictedMethodCall) initFields_ConformanceRules_RestrictedMethodCall() {
	_ = this_
}

// Java constructor: RestrictedMethodCall(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.Requirement requirement)
func NewConformanceRules_RestrictedMethodCall(compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) *ConformanceRules_RestrictedMethodCall {
	this_ := &ConformanceRules_RestrictedMethodCall{}
	this_.BindSelf_(this_)
	InitConformanceRules_RestrictedMethodCall(this_, compiler, requirement)
	return this_
}

func InitConformanceRules_RestrictedMethodCall(this_ *ConformanceRules_RestrictedMethodCall, compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) {
	_ = compiler
	_ = requirement
	InitConformanceRules_AbstractRule(&this_.ConformanceRules_AbstractRuleBase, compiler, requirement)
	this_.initFields_ConformanceRules_RestrictedMethodCall()
	if requirement.GetValueCount() == int32(0) {
		panic(NewCheckConformance_InvalidRequirementSpec_1("missing value"))
	}
	var registry *pkg_rhino.JSTypeRegistry = compiler.GetTypeRegistry()
	_ = registry
	var builder []*ConformanceRules_RestrictedMethodCall_Restriction = []*ConformanceRules_RestrictedMethodCall_Restriction{}
	_ = builder
	{
		cur_28 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](requirement.GetValueList())))
		for cur_28.HasNext() {
			value := cur_28.NextStr()
			_ = value
			var type_ string = ConformanceRules_ConformanceUtil_GetClassFromDeclarationName(ConformanceRules_ConformanceUtil_RemoveTypeDecl(value))
			_ = type_
			var property string = ConformanceRules_ConformanceUtil_GetPropertyFromDeclarationName(ConformanceRules_ConformanceUtil_RemoveTypeDecl(value))
			_ = property
			var restrictedDecl string = ConformanceRules_ConformanceUtil_GetTypeFromValue(value)
			_ = restrictedDecl
			if ((type_ == j2gort.NullString) || (property == j2gort.NullString)) || (restrictedDecl == j2gort.NullString) {
				panic(NewCheckConformance_InvalidRequirementSpec_1("bad prop value"))
			}
			var restrictedCallType *pkg_rhino.FunctionType = ConformanceRules_ConformanceUtil_EvaluateTypeString(compiler, restrictedDecl).ToMaybeFunctionType()
			_ = restrictedCallType
			if restrictedCallType == nil {
				panic(NewCheckConformance_InvalidRequirementSpec_1("invalid conformance type"))
			}
			builder = append(builder, NewConformanceRules_RestrictedMethodCall_Restriction(registry.GetGlobalType(type_), property, restrictedCallType))
		}
	}
	this_.restrictions = builder
}

// Java method: RestrictedMethodCall.checkConformance(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *ConformanceRules_RestrictedMethodCall) CheckConformance_2(t *NodeTraversal, n pkg_rhino.Node) *ConformanceRules_ConformanceResult {
	_ = t
	_ = n
	if !(n.IsGetProp()) || !(ConformanceRules_ConformanceUtil_IsCallTarget(n)) {
		return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
	}
	{
		var i int32 = int32(0)
		_ = i
		for ; i < int32(len(this_.restrictions)); i = i + 1 {
			var r *ConformanceRules_RestrictedMethodCall_Restriction = this_.restrictions[i]
			_ = r
			var result *ConformanceRules_ConformanceResult = StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
			_ = result
			if this_.priv_RestrictedMethodCall_matchesProp_2153(n, r) {
				result = this_.priv_RestrictedMethodCall_checkConformance_2152(t, n, r, false)
			} else if (j2gort.CheckNotNullString(n.GetString()) == "call") && this_.priv_RestrictedMethodCall_matchesProp_2153(n.GetFirstChild(), r) {
				result = this_.priv_RestrictedMethodCall_checkConformance_2152(t, n, r, true)
			}
			if result.level != StaticGet_ConformanceRules_ConformanceLevel_CONFORMANCE() {
				return result
			}
		}
	}
	return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
}

// Java method: RestrictedMethodCall.checkConformance(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.jscomp.ConformanceRules.RestrictedMethodCall.Restriction r, boolean isCallInvocation)
func (this_ *ConformanceRules_RestrictedMethodCall) priv_RestrictedMethodCall_checkConformance_2152(t *NodeTraversal, n pkg_rhino.Node, r *ConformanceRules_RestrictedMethodCall_Restriction, isCallInvocation bool) *ConformanceRules_ConformanceResult {
	_ = t
	_ = n
	_ = r
	_ = isCallInvocation
	var registry *pkg_rhino.JSTypeRegistry = t.GetCompiler().GetTypeRegistry()
	_ = registry
	var methodClassType pkg_rhino.JSType = r.type_
	_ = methodClassType
	var lhs pkg_rhino.Node
	_ = lhs
	if isCallInvocation {
		lhs = n.GetFirstFirstChild()
	} else {
		lhs = n.GetFirstChild()
	}
	if (methodClassType != nil && !methodClassType.IsJavaNil_()) && !(j2gort.IfaceEquals[pkg_rhino.JSType](lhs.GetJSType(), nil)) {
		var targetType pkg_rhino.JSType = lhs.GetJSType().RestrictByNotNullOrUndefined()
		_ = targetType
		if ConformanceRules_ConformanceUtil_IsLooseType(targetType) || ConformanceRules_IsLooseObject(targetType, registry) {
			if this_.GetConformanceRules_AbstractRuleBase().reportLooseTypeViolations && !(ConformanceRules_ConformanceUtil_ValidateCall(this_.GetConformanceRules_AbstractRuleBase().compiler, n.GetParent(), r.restrictedCallType, isCallInvocation)) {
				return StaticGet_ConformanceRules_ConformanceResult_POSSIBLE_VIOLATION_DUE_TO_LOOSE_TYPES()
			}
		} else if targetType.IsSubtypeOf_1(methodClassType) {
			if !(ConformanceRules_ConformanceUtil_ValidateCall(this_.GetConformanceRules_AbstractRuleBase().compiler, n.GetParent(), r.restrictedCallType, isCallInvocation)) {
				return StaticGet_ConformanceRules_ConformanceResult_VIOLATION()
			}
		}
	}
	return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
}

// Java method: RestrictedMethodCall.matchesProp(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.ConformanceRules.RestrictedMethodCall.Restriction r)
func (this_ *ConformanceRules_RestrictedMethodCall) priv_RestrictedMethodCall_matchesProp_2153(n pkg_rhino.Node, r *ConformanceRules_RestrictedMethodCall_Restriction) bool {
	_ = n
	_ = r
	return (n.IsGetProp() && func() bool {
		sr_29 := n.GetString()
		se_30 := r.property
		return j2gort.CheckNotNullString(sr_29) == se_30
	}())
}

// Java class: com.google.javascript.jscomp.ConformanceRules.RestrictedMethodCall.Restriction
type ConformanceRules_RestrictedMethodCall_Restriction struct {
	type_              pkg_rhino.JSType        `gson:"type"`
	property           string                  `gson:"property"`
	restrictedCallType *pkg_rhino.FunctionType `gson:"restrictedCallType"`
	idHash_            atomic.Int32
}

func (this_ *ConformanceRules_RestrictedMethodCall_Restriction) BindSelf_(s any) {
	_ = s
}
func (this_ *ConformanceRules_RestrictedMethodCall_Restriction) IsJavaNil_() bool {
	return this_ == nil
}
func (*ConformanceRules_RestrictedMethodCall_Restriction) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConformanceRules$RestrictedMethodCall$Restriction", "Restriction"
}
func (this_ *ConformanceRules_RestrictedMethodCall_Restriction) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ConformanceRules_RestrictedMethodCall_Restriction) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConformanceRules.RestrictedMethodCall.Restriction", "com.google.javascript.jscomp.ConformanceRules$RestrictedMethodCall$Restriction", func(v_ any) bool {
		return j2gort.InstanceOf[*ConformanceRules_RestrictedMethodCall_Restriction](v_)
	})
}

func (this_ *ConformanceRules_RestrictedMethodCall_Restriction) initFields_ConformanceRules_RestrictedMethodCall_Restriction() {
	_ = this_
	this_.property = j2gort.NullString
}

// Java constructor: Restriction(com.google.javascript.rhino.jstype.JSType type, java.lang.String property, com.google.javascript.rhino.jstype.FunctionType restrictedCallType)
func NewConformanceRules_RestrictedMethodCall_Restriction(type_ pkg_rhino.JSType, property string, restrictedCallType *pkg_rhino.FunctionType) *ConformanceRules_RestrictedMethodCall_Restriction {
	this_ := &ConformanceRules_RestrictedMethodCall_Restriction{}
	InitConformanceRules_RestrictedMethodCall_Restriction(this_, type_, property, restrictedCallType)
	return this_
}

func InitConformanceRules_RestrictedMethodCall_Restriction(this_ *ConformanceRules_RestrictedMethodCall_Restriction, type_ pkg_rhino.JSType, property string, restrictedCallType *pkg_rhino.FunctionType) {
	_ = type_
	_ = property
	_ = restrictedCallType
	this_.initFields_ConformanceRules_RestrictedMethodCall_Restriction()
	this_.type_ = type_
	this_.property = property
	this_.restrictedCallType = restrictedCallType
}

// Java class: com.google.javascript.jscomp.ConformanceRules.RestrictedPropertyWrite
type ConformanceRules_RestrictedPropertyWrite struct {
	ConformanceRules_AbstractRuleBase
	restrictions []*ConformanceRules_RestrictedPropertyWrite_Restriction `gson:"restrictions"`
}

func (this_ *ConformanceRules_RestrictedPropertyWrite) BindSelf_(s any) {
	_ = s
	this_.ConformanceRules_AbstractRuleBase.BindSelf_(s)
}
func (this_ *ConformanceRules_RestrictedPropertyWrite) IsJavaNil_() bool { return this_ == nil }
func (*ConformanceRules_RestrictedPropertyWrite) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConformanceRules$RestrictedPropertyWrite", "RestrictedPropertyWrite"
}
func (this_ *ConformanceRules_RestrictedPropertyWrite) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ConformanceRules_AbstractRuleBase.IdentityHash_()
}
func (*ConformanceRules_RestrictedPropertyWrite) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConformanceRules.RestrictedPropertyWrite", "com.google.javascript.jscomp.ConformanceRules$RestrictedPropertyWrite", func(v_ any) bool {
		return j2gort.InstanceOf[*ConformanceRules_RestrictedPropertyWrite](v_)
	})
}
func (*ConformanceRules_RestrictedPropertyWrite) Implements_jscomp_ConformanceRules_AbstractRule() {}
func (*ConformanceRules_RestrictedPropertyWrite) Implements_jscomp_CheckConformance_Rule()         {}

func (this_ *ConformanceRules_RestrictedPropertyWrite) initFields_ConformanceRules_RestrictedPropertyWrite() {
	_ = this_
}

// Java constructor: RestrictedPropertyWrite(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.Requirement requirement)
func NewConformanceRules_RestrictedPropertyWrite(compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) *ConformanceRules_RestrictedPropertyWrite {
	this_ := &ConformanceRules_RestrictedPropertyWrite{}
	this_.BindSelf_(this_)
	InitConformanceRules_RestrictedPropertyWrite(this_, compiler, requirement)
	return this_
}

func InitConformanceRules_RestrictedPropertyWrite(this_ *ConformanceRules_RestrictedPropertyWrite, compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) {
	_ = compiler
	_ = requirement
	InitConformanceRules_AbstractRule(&this_.ConformanceRules_AbstractRuleBase, compiler, requirement)
	this_.initFields_ConformanceRules_RestrictedPropertyWrite()
	if requirement.GetValueCount() == int32(0) {
		panic(NewCheckConformance_InvalidRequirementSpec_1("missing value"))
	}
	var registry *pkg_rhino.JSTypeRegistry = compiler.GetTypeRegistry()
	_ = registry
	var builder []*ConformanceRules_RestrictedPropertyWrite_Restriction = []*ConformanceRules_RestrictedPropertyWrite_Restriction{}
	_ = builder
	{
		cur_31 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](requirement.GetValueList())))
		for cur_31.HasNext() {
			value := cur_31.NextStr()
			_ = value
			var type_ string = ConformanceRules_ConformanceUtil_GetClassFromDeclarationName(ConformanceRules_ConformanceUtil_RemoveTypeDecl(value))
			_ = type_
			var property string = ConformanceRules_ConformanceUtil_GetPropertyFromDeclarationName(ConformanceRules_ConformanceUtil_RemoveTypeDecl(value))
			_ = property
			var restrictedDecl string = ConformanceRules_ConformanceUtil_GetTypeFromValue(value)
			_ = restrictedDecl
			if ((type_ == j2gort.NullString) || (property == j2gort.NullString)) || (restrictedDecl == j2gort.NullString) {
				panic(NewCheckConformance_InvalidRequirementSpec_1("bad prop value"))
			}
			var restrictedType pkg_rhino.JSType = ConformanceRules_ConformanceUtil_EvaluateTypeString(compiler, restrictedDecl)
			_ = restrictedType
			builder = append(builder, NewConformanceRules_RestrictedPropertyWrite_Restriction(registry.GetGlobalType(type_), property, restrictedType))
		}
	}
	this_.restrictions = builder
}

// Java method: RestrictedPropertyWrite.checkConformance(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *ConformanceRules_RestrictedPropertyWrite) CheckConformance_2(t *NodeTraversal, n pkg_rhino.Node) *ConformanceRules_ConformanceResult {
	_ = t
	_ = n
	if n.IsGetProp() && NodeUtil_IsLhsOfAssign(n) {
		var rhsType pkg_rhino.JSType = n.GetNext().GetJSType()
		_ = rhsType
		var targetType pkg_rhino.JSType = n.GetFirstChild().GetJSType()
		_ = targetType
		if (rhsType != nil && !rhsType.IsJavaNil_()) && (targetType != nil && !targetType.IsJavaNil_()) {
			var targetNotNullType pkg_rhino.JSType = nil
			_ = targetNotNullType
			{

				for _, r := range this_.restrictions {

					_ = r
					{
						var iifeCond_1 bool
						{

							sr_34 := n.GetString()
							se_35 := r.property
							iifeCond_1 = j2gort.CheckNotNullString(sr_34) == se_35
						}
						if iifeCond_1 {
							if !(rhsType.IsSubtypeOf_1(r.restrictedType)) {
								if ConformanceRules_ConformanceUtil_IsLooseType(targetType) {
									if this_.GetConformanceRules_AbstractRuleBase().reportLooseTypeViolations {
										return StaticGet_ConformanceRules_ConformanceResult_POSSIBLE_VIOLATION_DUE_TO_LOOSE_TYPES()
									}
								} else {
									if targetNotNullType == nil || targetNotNullType.IsJavaNil_() {
										targetNotNullType = targetType.RestrictByNotNullOrUndefined()
									}
									if targetNotNullType.IsSubtypeOf_1(r.type_) {
										return StaticGet_ConformanceRules_ConformanceResult_VIOLATION()
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
}

// Java class: com.google.javascript.jscomp.ConformanceRules.RestrictedPropertyWrite.Restriction
type ConformanceRules_RestrictedPropertyWrite_Restriction struct {
	type_          pkg_rhino.JSType `gson:"type"`
	property       string           `gson:"property"`
	restrictedType pkg_rhino.JSType `gson:"restrictedType"`
	idHash_        atomic.Int32
}

func (this_ *ConformanceRules_RestrictedPropertyWrite_Restriction) BindSelf_(s any) {
	_ = s
}
func (this_ *ConformanceRules_RestrictedPropertyWrite_Restriction) IsJavaNil_() bool {
	return this_ == nil
}
func (*ConformanceRules_RestrictedPropertyWrite_Restriction) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConformanceRules$RestrictedPropertyWrite$Restriction", "Restriction"
}
func (this_ *ConformanceRules_RestrictedPropertyWrite_Restriction) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ConformanceRules_RestrictedPropertyWrite_Restriction) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConformanceRules.RestrictedPropertyWrite.Restriction", "com.google.javascript.jscomp.ConformanceRules$RestrictedPropertyWrite$Restriction", func(v_ any) bool {
		return j2gort.InstanceOf[*ConformanceRules_RestrictedPropertyWrite_Restriction](v_)
	})
}

func (this_ *ConformanceRules_RestrictedPropertyWrite_Restriction) initFields_ConformanceRules_RestrictedPropertyWrite_Restriction() {
	_ = this_
	this_.property = j2gort.NullString
}

// Java constructor: Restriction(com.google.javascript.rhino.jstype.JSType type, java.lang.String property, com.google.javascript.rhino.jstype.JSType restrictedType)
func NewConformanceRules_RestrictedPropertyWrite_Restriction(type_ pkg_rhino.JSType, property string, restrictedType pkg_rhino.JSType) *ConformanceRules_RestrictedPropertyWrite_Restriction {
	this_ := &ConformanceRules_RestrictedPropertyWrite_Restriction{}
	InitConformanceRules_RestrictedPropertyWrite_Restriction(this_, type_, property, restrictedType)
	return this_
}

func InitConformanceRules_RestrictedPropertyWrite_Restriction(this_ *ConformanceRules_RestrictedPropertyWrite_Restriction, type_ pkg_rhino.JSType, property string, restrictedType pkg_rhino.JSType) {
	_ = type_
	_ = property
	_ = restrictedType
	this_.initFields_ConformanceRules_RestrictedPropertyWrite_Restriction()
	this_.type_ = type_
	this_.property = property
	this_.restrictedType = restrictedType
}

// Java class: com.google.javascript.jscomp.ConformanceRules.BannedCodePattern
type ConformanceRules_BannedCodePattern struct {
	ConformanceRules_AbstractRuleBase
	restrictions []*TemplateAstMatcher `gson:"restrictions"`
}

func (this_ *ConformanceRules_BannedCodePattern) BindSelf_(s any) {
	_ = s
	this_.ConformanceRules_AbstractRuleBase.BindSelf_(s)
}
func (this_ *ConformanceRules_BannedCodePattern) IsJavaNil_() bool { return this_ == nil }
func (*ConformanceRules_BannedCodePattern) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConformanceRules$BannedCodePattern", "BannedCodePattern"
}
func (this_ *ConformanceRules_BannedCodePattern) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ConformanceRules_AbstractRuleBase.IdentityHash_()
}
func (*ConformanceRules_BannedCodePattern) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConformanceRules.BannedCodePattern", "com.google.javascript.jscomp.ConformanceRules$BannedCodePattern", func(v_ any) bool {
		return j2gort.InstanceOf[*ConformanceRules_BannedCodePattern](v_)
	})
}
func (*ConformanceRules_BannedCodePattern) Implements_jscomp_ConformanceRules_AbstractRule() {}
func (*ConformanceRules_BannedCodePattern) Implements_jscomp_CheckConformance_Rule()         {}

func (this_ *ConformanceRules_BannedCodePattern) initFields_ConformanceRules_BannedCodePattern() {
	_ = this_
}

// Java constructor: BannedCodePattern(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.Requirement requirement)
func NewConformanceRules_BannedCodePattern(compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) *ConformanceRules_BannedCodePattern {
	this_ := &ConformanceRules_BannedCodePattern{}
	this_.BindSelf_(this_)
	InitConformanceRules_BannedCodePattern(this_, compiler, requirement)
	return this_
}

func InitConformanceRules_BannedCodePattern(this_ *ConformanceRules_BannedCodePattern, compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) {
	_ = compiler
	_ = requirement
	InitConformanceRules_AbstractRule(&this_.ConformanceRules_AbstractRuleBase, compiler, requirement)
	this_.initFields_ConformanceRules_BannedCodePattern()
	if requirement.GetValueCount() == int32(0) {
		panic(NewCheckConformance_InvalidRequirementSpec_1("missing value"))
	}
	var builder []*TemplateAstMatcher = []*TemplateAstMatcher{}
	_ = builder
	{
		cur_36 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](requirement.GetValueList())))
		for cur_36.HasNext() {
			value := cur_36.NextStr()
			_ = value
			var parseRoot pkg_rhino.Node = NewCompilerInput_1(SourceFile_FromCode_2("<template>", value)).GetAstRoot(compiler)
			_ = parseRoot
			if !(parseRoot.HasOneChild()) || !(parseRoot.GetFirstChild().IsFunction()) {
				panic(NewCheckConformance_InvalidRequirementSpec_1(("invalid conformance template: " + j2gort.StrValString(value))))
			}
			var templateRoot pkg_rhino.Node = parseRoot.GetFirstChild()
			_ = templateRoot
			var astMatcher *TemplateAstMatcher = NewTemplateAstMatcher(compiler, templateRoot, this_.GetConformanceRules_AbstractRuleBase().typeMatchingStrategy)
			_ = astMatcher
			builder = append(builder, j2gort.CheckNotNullElem(astMatcher))
		}
	}
	this_.restrictions = builder
}

// Java method: BannedCodePattern.checkConformance(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *ConformanceRules_BannedCodePattern) CheckConformance_2(t *NodeTraversal, n pkg_rhino.Node) *ConformanceRules_ConformanceResult {
	_ = t
	_ = n
	var possibleViolation bool = false
	_ = possibleViolation
	{
		var i int32 = int32(0)
		_ = i
		for ; i < int32(len(this_.restrictions)); i = i + 1 {
			var matcher *TemplateAstMatcher = this_.restrictions[i]
			_ = matcher
			if matcher.Matches(n) {
				if matcher.IsLooseMatch() {
					possibleViolation = true
				} else {
					return StaticGet_ConformanceRules_ConformanceResult_VIOLATION()
				}
			}
		}
	}
	if possibleViolation && this_.GetConformanceRules_AbstractRuleBase().reportLooseTypeViolations {
		return StaticGet_ConformanceRules_ConformanceResult_POSSIBLE_VIOLATION_DUE_TO_LOOSE_TYPES()
	} else {
		return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
	}
}

// Java class: com.google.javascript.jscomp.ConformanceRules.CustomRuleProxy
type ConformanceRules_CustomRuleProxy struct {
	customRule CheckConformance_Rule `gson:"customRule"`
	idHash_    atomic.Int32
}

func (this_ *ConformanceRules_CustomRuleProxy) BindSelf_(s any) {
	_ = s
}
func (this_ *ConformanceRules_CustomRuleProxy) IsJavaNil_() bool { return this_ == nil }
func (*ConformanceRules_CustomRuleProxy) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConformanceRules$CustomRuleProxy", "CustomRuleProxy"
}
func (this_ *ConformanceRules_CustomRuleProxy) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ConformanceRules_CustomRuleProxy) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConformanceRules.CustomRuleProxy", "com.google.javascript.jscomp.ConformanceRules$CustomRuleProxy", func(v_ any) bool {
		return j2gort.InstanceOf[*ConformanceRules_CustomRuleProxy](v_)
	})
}
func (*ConformanceRules_CustomRuleProxy) Implements_jscomp_CheckConformance_Rule() {}

var static_ConformanceRules_CustomRuleProxy_CUSTOM_RULE_REGISTRY j2gort.Ext_com_google_common_collect_ImmutableMap

func StaticGet_ConformanceRules_CustomRuleProxy_CUSTOM_RULE_REGISTRY() j2gort.Ext_com_google_common_collect_ImmutableMap {
	clinit_ConformanceRules_CustomRuleProxy()
	return static_ConformanceRules_CustomRuleProxy_CUSTOM_RULE_REGISTRY
}
func StaticSet_ConformanceRules_CustomRuleProxy_CUSTOM_RULE_REGISTRY(v j2gort.Ext_com_google_common_collect_ImmutableMap) {
	clinit_ConformanceRules_CustomRuleProxy()
	static_ConformanceRules_CustomRuleProxy_CUSTOM_RULE_REGISTRY = v
}

var clinitDone_ConformanceRules_CustomRuleProxy uint32
var clinitState_ConformanceRules_CustomRuleProxy j2gort.ClinitState

func clinit_ConformanceRules_CustomRuleProxy() {
	if atomic.LoadUint32(&clinitDone_ConformanceRules_CustomRuleProxy) != 0 {
		return
	}
	clinitSlow_ConformanceRules_CustomRuleProxy()
}

func clinitSlow_ConformanceRules_CustomRuleProxy() {
	if !j2gort.EnterClinit(&clinitState_ConformanceRules_CustomRuleProxy, &clinitDone_ConformanceRules_CustomRuleProxy) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ConformanceRules_CustomRuleProxy, &clinitDone_ConformanceRules_CustomRuleProxy)
	static_ConformanceRules_CustomRuleProxy_CUSTOM_RULE_REGISTRY = ConformanceRules_CustomRuleProxy_CreateCustomRuleRegistry()
}

func (this_ *ConformanceRules_CustomRuleProxy) initFields_ConformanceRules_CustomRuleProxy() {
	_ = this_
}

// Java constructor: CustomRuleProxy(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.Requirement requirement)
func NewConformanceRules_CustomRuleProxy(compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) *ConformanceRules_CustomRuleProxy {
	clinit_ConformanceRules_CustomRuleProxy()
	this_ := &ConformanceRules_CustomRuleProxy{}
	InitConformanceRules_CustomRuleProxy(this_, compiler, requirement)
	return this_
}

func InitConformanceRules_CustomRuleProxy(this_ *ConformanceRules_CustomRuleProxy, compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) {
	_ = compiler
	_ = requirement
	this_.initFields_ConformanceRules_CustomRuleProxy()
	if !(requirement.HasJavaClass()) {
		panic(NewCheckConformance_InvalidRequirementSpec_1("missing java_class"))
	}
	this_.customRule = this_.priv_CustomRuleProxy_createRule_2154(compiler, requirement)
}

// Java method: CustomRuleProxy.getPrecondition()
func (this_ *ConformanceRules_CustomRuleProxy) GetPrecondition() CheckConformance_Precondition {
	return this_.customRule.GetPrecondition()
}

// Java method: CustomRuleProxy.check(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.jscomp.ConformanceConfig.LibraryLevelNonAllowlistedConformanceViolationsBehavior behavior)
func (this_ *ConformanceRules_CustomRuleProxy) Check(t *NodeTraversal, n pkg_rhino.Node, behavior *j2gort.Ext_com_google_javascript_jscomp_ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior) {
	_ = t
	_ = n
	_ = behavior
	this_.customRule.Check(t, n, behavior)
}

// Java method: CustomRuleProxy.createRule(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.Requirement requirement)
func (this_ *ConformanceRules_CustomRuleProxy) priv_CustomRuleProxy_createRule_2154(compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) CheckConformance_Rule {
	_ = compiler
	_ = requirement
	var factory ConformanceRules_CustomRuleProxy_CustomRuleFactory = j2gort.AsIface[ConformanceRules_CustomRuleProxy_CustomRuleFactory](j2gort.MapGetStr(StaticGet_ConformanceRules_CustomRuleProxy_CUSTOM_RULE_REGISTRY(), requirement.GetJavaClass()))
	_ = factory
	if factory == nil || factory.IsJavaNil_() {
		panic(NewCheckConformance_InvalidRequirementSpec_1(("Unknown conformance rule: " + j2gort.StrValString(requirement.GetJavaClass()))))
	}
	return factory.Create(compiler, requirement)
}

// Java static method: CustomRuleProxy.createCustomRuleRegistry()
func ConformanceRules_CustomRuleProxy_CreateCustomRuleRegistry() j2gort.Ext_com_google_common_collect_ImmutableMap {
	clinit_ConformanceRules_CustomRuleProxy()
	var builder j2gort.Ext_com_google_common_collect_ImmutableMap_Builder = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.Ext_com_google_common_collect_ImmutableMap_S_Builder().Put_2("com.google.javascript.jscomp.ConformanceRules$NoOp", lambdaSingleton_conformance_rules_30)).Put_2("com.google.javascript.jscomp.ConformanceRules$InferredConstCheck", lambdaSingleton_conformance_rules_29)).Put_2("com.google.javascript.jscomp.ConformanceRules$BannedDependency", lambdaSingleton_conformance_rules_28)).Put_2("com.google.javascript.jscomp.ConformanceRules$BannedDependencyRegex", lambdaSingleton_conformance_rules_27)).Put_2("com.google.javascript.jscomp.ConformanceRules$BannedName", lambdaSingleton_conformance_rules_26)).Put_2("com.google.javascript.jscomp.ConformanceRules$BannedProperty", lambdaSingleton_conformance_rules_25)).Put_2("com.google.javascript.jscomp.ConformanceRules$BannedStringRegex", lambdaSingleton_conformance_rules_24)).Put_2("com.google.javascript.jscomp.ConformanceRules$RestrictedNameCall", lambdaSingleton_conformance_rules_23)).Put_2("com.google.javascript.jscomp.ConformanceRules$RestrictedMethodCall", lambdaSingleton_conformance_rules_22)).Put_2("com.google.javascript.jscomp.ConformanceRules$RestrictedPropertyWrite", lambdaSingleton_conformance_rules_21)).Put_2("com.google.javascript.jscomp.ConformanceRules$BannedCodePattern", lambdaSingleton_conformance_rules_20)).Put_2("com.google.javascript.jscomp.ConformanceRules$BanForOf", lambdaSingleton_conformance_rules_19)).Put_2("com.google.javascript.jscomp.ConformanceRules$RequireUseStrict", lambdaSingleton_conformance_rules_18)).Put_2("com.google.javascript.jscomp.ConformanceRules$BanThrowOfNonErrorTypes", lambdaSingleton_conformance_rules_17)).Put_2("com.google.javascript.jscomp.ConformanceRules$BanNullDeref", lambdaSingleton_conformance_rules_16)).Put_2("com.google.javascript.jscomp.ConformanceRules$BanUnknownThis", lambdaSingleton_conformance_rules_15)).Put_2("com.google.javascript.jscomp.ConformanceRules$BanUnknownDirectThisPropsReferences", lambdaSingleton_conformance_rules_14)).Put_2("com.google.javascript.jscomp.ConformanceRules$BanUnknownTypedClassPropsReferences", lambdaSingleton_conformance_rules_13)).Put_2("com.google.javascript.jscomp.ConformanceRules$BanUnresolvedType", lambdaSingleton_conformance_rules_12)).Put_2("com.google.javascript.jscomp.ConformanceRules$StrictBanUnresolvedType", lambdaSingleton_conformance_rules_11)).Put_2("com.google.javascript.jscomp.ConformanceRules$BanGlobalVars", lambdaSingleton_conformance_rules_10)).Put_2("com.google.javascript.jscomp.ConformanceRules$BannedEnhance", lambdaSingleton_conformance_rules_9)).Put_2("com.google.javascript.jscomp.ConformanceRules$BannedModsRegex", lambdaSingleton_conformance_rules_8)).Put_2("com.google.javascript.jscomp.ConformanceRules$BanCreateElement", lambdaSingleton_conformance_rules_7)).Put_2("com.google.javascript.jscomp.ConformanceRules$BanCreateDom", lambdaSingleton_conformance_rules_6)).Put_2("com.google.javascript.jscomp.ConformanceRules$BanElementSetAttribute", lambdaSingleton_conformance_rules_5)).Put_2("com.google.javascript.jscomp.ConformanceRules$BanSettingAttributes", lambdaSingleton_conformance_rules_4)).Put_2("com.google.javascript.jscomp.ConformanceRules$BanExecCommand", lambdaSingleton_conformance_rules_3)).Put_2("com.google.javascript.jscomp.ConformanceRules$BanStaticThis", lambdaSingleton_conformance_rules_2))
	_ = builder
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](builder.BuildOrThrow())
}

// Java interface: com.google.javascript.jscomp.ConformanceRules.CustomRuleProxy.CustomRuleFactory
type ConformanceRules_CustomRuleProxy_CustomRuleFactory interface {
	Implements_jscomp_ConformanceRules_CustomRuleProxy_CustomRuleFactory()
	IsJavaNil_() bool
	Create(AbstractCompiler, *j2gort.Ext_com_google_javascript_jscomp_Requirement) CheckConformance_Rule
}

func ToIface_ConformanceRules_CustomRuleProxy_CustomRuleFactory[T interface {
	ConformanceRules_CustomRuleProxy_CustomRuleFactory
	comparable
}](v T) ConformanceRules_CustomRuleProxy_CustomRuleFactory {
	var z T
	if v == z {
		return nil
	}
	return v
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConformanceRules.CustomRuleProxy.CustomRuleFactory", "com.google.javascript.jscomp.ConformanceRules$CustomRuleProxy$CustomRuleFactory", func(v_ any) bool {
		return j2gort.InstanceOf[ConformanceRules_CustomRuleProxy_CustomRuleFactory](v_)
	})
}

type lambda_ConformanceRules_CustomRuleProxy_CustomRuleFactory struct {
	fn func(AbstractCompiler, *j2gort.Ext_com_google_javascript_jscomp_Requirement) CheckConformance_Rule
}

func (l *lambda_ConformanceRules_CustomRuleProxy_CustomRuleFactory) IsJavaNil_() bool {
	return l == nil
}
func NewLambda_ConformanceRules_CustomRuleProxy_CustomRuleFactory(fn func(AbstractCompiler, *j2gort.Ext_com_google_javascript_jscomp_Requirement) CheckConformance_Rule) ConformanceRules_CustomRuleProxy_CustomRuleFactory {
	if fn == nil {
		return nil
	}
	return &lambda_ConformanceRules_CustomRuleProxy_CustomRuleFactory{fn: fn}
}

func (*lambda_ConformanceRules_CustomRuleProxy_CustomRuleFactory) Implements_jscomp_ConformanceRules_CustomRuleProxy_CustomRuleFactory() {
}
func (l *lambda_ConformanceRules_CustomRuleProxy_CustomRuleFactory) Create(p0_ AbstractCompiler, p1_ *j2gort.Ext_com_google_javascript_jscomp_Requirement) CheckConformance_Rule {
	return l.fn(p0_, p1_)
}

// Java class: com.google.javascript.jscomp.ConformanceRules.BanForOf
type ConformanceRules_BanForOf struct {
	ConformanceRules_AbstractRuleBase
}

func (this_ *ConformanceRules_BanForOf) BindSelf_(s any) {
	_ = s
	this_.ConformanceRules_AbstractRuleBase.BindSelf_(s)
}
func (this_ *ConformanceRules_BanForOf) IsJavaNil_() bool { return this_ == nil }
func (*ConformanceRules_BanForOf) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConformanceRules$BanForOf", "BanForOf"
}
func (this_ *ConformanceRules_BanForOf) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ConformanceRules_AbstractRuleBase.IdentityHash_()
}
func (*ConformanceRules_BanForOf) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConformanceRules.BanForOf", "com.google.javascript.jscomp.ConformanceRules$BanForOf", func(v_ any) bool {
		return j2gort.InstanceOf[*ConformanceRules_BanForOf](v_)
	})
}
func (*ConformanceRules_BanForOf) Implements_jscomp_ConformanceRules_AbstractRule() {}
func (*ConformanceRules_BanForOf) Implements_jscomp_CheckConformance_Rule()         {}

func (this_ *ConformanceRules_BanForOf) initFields_ConformanceRules_BanForOf() {
	_ = this_
}

// Java constructor: BanForOf(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.Requirement requirement)
func NewConformanceRules_BanForOf(compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) *ConformanceRules_BanForOf {
	this_ := &ConformanceRules_BanForOf{}
	this_.BindSelf_(this_)
	InitConformanceRules_BanForOf(this_, compiler, requirement)
	return this_
}

func InitConformanceRules_BanForOf(this_ *ConformanceRules_BanForOf, compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) {
	_ = compiler
	_ = requirement
	InitConformanceRules_AbstractRule(&this_.ConformanceRules_AbstractRuleBase, compiler, requirement)
	this_.initFields_ConformanceRules_BanForOf()
}

// Java method: BanForOf.checkConformance(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *ConformanceRules_BanForOf) CheckConformance_2(t *NodeTraversal, n pkg_rhino.Node) *ConformanceRules_ConformanceResult {
	_ = t
	_ = n
	if n.IsForOf() || n.IsForAwaitOf() {
		return StaticGet_ConformanceRules_ConformanceResult_VIOLATION()
	}
	return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
}

// Java class: com.google.javascript.jscomp.ConformanceRules.RequireUseStrict
type ConformanceRules_RequireUseStrict struct {
	ConformanceRules_AbstractRuleBase
}

func (this_ *ConformanceRules_RequireUseStrict) BindSelf_(s any) {
	_ = s
	this_.ConformanceRules_AbstractRuleBase.BindSelf_(s)
}
func (this_ *ConformanceRules_RequireUseStrict) IsJavaNil_() bool { return this_ == nil }
func (*ConformanceRules_RequireUseStrict) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConformanceRules$RequireUseStrict", "RequireUseStrict"
}
func (this_ *ConformanceRules_RequireUseStrict) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ConformanceRules_AbstractRuleBase.IdentityHash_()
}
func (*ConformanceRules_RequireUseStrict) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConformanceRules.RequireUseStrict", "com.google.javascript.jscomp.ConformanceRules$RequireUseStrict", func(v_ any) bool {
		return j2gort.InstanceOf[*ConformanceRules_RequireUseStrict](v_)
	})
}
func (*ConformanceRules_RequireUseStrict) Implements_jscomp_ConformanceRules_AbstractRule() {}
func (*ConformanceRules_RequireUseStrict) Implements_jscomp_CheckConformance_Rule()         {}

func (this_ *ConformanceRules_RequireUseStrict) initFields_ConformanceRules_RequireUseStrict() {
	_ = this_
}

// Java constructor: RequireUseStrict(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.Requirement requirement)
func NewConformanceRules_RequireUseStrict(compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) *ConformanceRules_RequireUseStrict {
	this_ := &ConformanceRules_RequireUseStrict{}
	this_.BindSelf_(this_)
	InitConformanceRules_RequireUseStrict(this_, compiler, requirement)
	return this_
}

func InitConformanceRules_RequireUseStrict(this_ *ConformanceRules_RequireUseStrict, compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) {
	_ = compiler
	_ = requirement
	InitConformanceRules_AbstractRule(&this_.ConformanceRules_AbstractRuleBase, compiler, requirement)
	this_.initFields_ConformanceRules_RequireUseStrict()
	if !(j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](requirement.GetValueList()).IsEmpty()) {
		panic(NewCheckConformance_InvalidRequirementSpec_1("invalid value"))
	}
}

// Java method: RequireUseStrict.checkConformance(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *ConformanceRules_RequireUseStrict) CheckConformance_2(t *NodeTraversal, n pkg_rhino.Node) *ConformanceRules_ConformanceResult {
	_ = t
	_ = n
	if n.IsScript() && !(n.IsUseStrict()) {
		return StaticGet_ConformanceRules_ConformanceResult_VIOLATION()
	} else {
		return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
	}
}

// Java class: com.google.javascript.jscomp.ConformanceRules.BanThrowOfNonErrorTypes
type ConformanceRules_BanThrowOfNonErrorTypes struct {
	ConformanceRules_AbstractRuleBase
	errorObjType pkg_rhino.JSType `gson:"errorObjType"`
}

func (this_ *ConformanceRules_BanThrowOfNonErrorTypes) BindSelf_(s any) {
	_ = s
	this_.ConformanceRules_AbstractRuleBase.BindSelf_(s)
}
func (this_ *ConformanceRules_BanThrowOfNonErrorTypes) IsJavaNil_() bool { return this_ == nil }
func (*ConformanceRules_BanThrowOfNonErrorTypes) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConformanceRules$BanThrowOfNonErrorTypes", "BanThrowOfNonErrorTypes"
}
func (this_ *ConformanceRules_BanThrowOfNonErrorTypes) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ConformanceRules_AbstractRuleBase.IdentityHash_()
}
func (*ConformanceRules_BanThrowOfNonErrorTypes) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConformanceRules.BanThrowOfNonErrorTypes", "com.google.javascript.jscomp.ConformanceRules$BanThrowOfNonErrorTypes", func(v_ any) bool {
		return j2gort.InstanceOf[*ConformanceRules_BanThrowOfNonErrorTypes](v_)
	})
}
func (*ConformanceRules_BanThrowOfNonErrorTypes) Implements_jscomp_ConformanceRules_AbstractRule() {}
func (*ConformanceRules_BanThrowOfNonErrorTypes) Implements_jscomp_CheckConformance_Rule()         {}

func (this_ *ConformanceRules_BanThrowOfNonErrorTypes) initFields_ConformanceRules_BanThrowOfNonErrorTypes() {
	_ = this_
}

// Java constructor: BanThrowOfNonErrorTypes(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.Requirement requirement)
func NewConformanceRules_BanThrowOfNonErrorTypes(compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) *ConformanceRules_BanThrowOfNonErrorTypes {
	this_ := &ConformanceRules_BanThrowOfNonErrorTypes{}
	this_.BindSelf_(this_)
	InitConformanceRules_BanThrowOfNonErrorTypes(this_, compiler, requirement)
	return this_
}

func InitConformanceRules_BanThrowOfNonErrorTypes(this_ *ConformanceRules_BanThrowOfNonErrorTypes, compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) {
	_ = compiler
	_ = requirement
	InitConformanceRules_AbstractRule(&this_.ConformanceRules_AbstractRuleBase, compiler, requirement)
	this_.initFields_ConformanceRules_BanThrowOfNonErrorTypes()
	this_.errorObjType = compiler.GetTypeRegistry().GetGlobalType("Error")
}

// Java method: BanThrowOfNonErrorTypes.checkConformance(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *ConformanceRules_BanThrowOfNonErrorTypes) CheckConformance_2(t *NodeTraversal, n pkg_rhino.Node) *ConformanceRules_ConformanceResult {
	_ = t
	_ = n
	if (this_.errorObjType != nil && !this_.errorObjType.IsJavaNil_()) && n.IsThrow() {
		var thrown pkg_rhino.JSType = n.GetFirstChild().GetJSType()
		_ = thrown
		if thrown != nil && !thrown.IsJavaNil_() {
			if ((!(thrown.IsUnknownType()) && !(thrown.IsAllType())) && !(thrown.IsEmptyType())) && !(thrown.IsSubtypeOf_1(this_.errorObjType)) {
				return StaticGet_ConformanceRules_ConformanceResult_VIOLATION()
			}
		}
	}
	return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
}

// Java class: com.google.javascript.jscomp.ConformanceRules.BanNullDeref
type ConformanceRules_BanNullDeref struct {
	ConformanceRules_AbstractTypeRestrictionRuleBase
}

func (this_ *ConformanceRules_BanNullDeref) BindSelf_(s any) {
	_ = s
	this_.ConformanceRules_AbstractTypeRestrictionRuleBase.BindSelf_(s)
}
func (this_ *ConformanceRules_BanNullDeref) IsJavaNil_() bool { return this_ == nil }
func (*ConformanceRules_BanNullDeref) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConformanceRules$BanNullDeref", "BanNullDeref"
}
func (this_ *ConformanceRules_BanNullDeref) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ConformanceRules_AbstractTypeRestrictionRuleBase.IdentityHash_()
}
func (*ConformanceRules_BanNullDeref) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConformanceRules.BanNullDeref", "com.google.javascript.jscomp.ConformanceRules$BanNullDeref", func(v_ any) bool {
		return j2gort.InstanceOf[*ConformanceRules_BanNullDeref](v_)
	})
}
func (*ConformanceRules_BanNullDeref) Implements_jscomp_ConformanceRules_AbstractTypeRestrictionRule() {
}
func (*ConformanceRules_BanNullDeref) Implements_jscomp_ConformanceRules_AbstractRule() {}
func (*ConformanceRules_BanNullDeref) Implements_jscomp_CheckConformance_Rule()         {}

func (this_ *ConformanceRules_BanNullDeref) initFields_ConformanceRules_BanNullDeref() {
	_ = this_
}

// Java constructor: BanNullDeref(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.Requirement requirement)
func NewConformanceRules_BanNullDeref(compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) *ConformanceRules_BanNullDeref {
	this_ := &ConformanceRules_BanNullDeref{}
	this_.BindSelf_(this_)
	InitConformanceRules_BanNullDeref(this_, compiler, requirement)
	return this_
}

func InitConformanceRules_BanNullDeref(this_ *ConformanceRules_BanNullDeref, compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) {
	_ = compiler
	_ = requirement
	InitConformanceRules_AbstractTypeRestrictionRule(&this_.ConformanceRules_AbstractTypeRestrictionRuleBase, compiler, requirement)
	this_.initFields_ConformanceRules_BanNullDeref()
}

// Java method: BanNullDeref.checkConformance(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *ConformanceRules_BanNullDeref) CheckConformance_2(t *NodeTraversal, n pkg_rhino.Node) *ConformanceRules_ConformanceResult {
	_ = t
	_ = n
	var violation bool = false
	_ = violation
	switch n.TokenSwitchOrd_() {
	case int32(26) /* Token.GETPROP */, int32(27) /* Token.GETELEM */, int32(23) /* Token.NEW */, int32(28) /* Token.CALL */ :
		violation = this_.Report_1(n.GetFirstChild())
	case int32(44) /* Token.IN */ :
		violation = this_.Report_1(n.GetLastChild())
	default:
		violation = false
	}
	if violation {
		return StaticGet_ConformanceRules_ConformanceResult_VIOLATION()
	} else {
		return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
	}
}

// Java method: BanNullDeref.report(com.google.javascript.rhino.Node n)
func (this_ *ConformanceRules_BanNullDeref) Report_1(n pkg_rhino.Node) bool {
	_ = n
	return (((!(j2gort.IfaceEquals[pkg_rhino.JSType](n.GetJSType(), nil)) && ConformanceRules_AbstractTypeRestrictionRule_IsKnown(n)) && this_.priv_BanNullDeref_invalidDeref_2157(n)) && !(this_.IsAllowlistedType(n)))
}

// Java method: BanNullDeref.invalidDeref(com.google.javascript.rhino.Node n)
func (this_ *ConformanceRules_BanNullDeref) priv_BanNullDeref_invalidDeref_2157(n pkg_rhino.Node) bool {
	_ = n
	var type_ pkg_rhino.JSType = n.GetJSType()
	_ = type_
	return (!(type_.IsAllType()) && (type_.IsNullable() || type_.IsVoidable()))
}

// Java class: com.google.javascript.jscomp.ConformanceRules.BanUnknownThis
type ConformanceRules_BanUnknownThis struct {
	ConformanceRules_AbstractTypeRestrictionRuleBase
	reports *pkg_base.LinkedIdentityHashSet[pkg_rhino.Node] `gson:"reports"`
}

func (this_ *ConformanceRules_BanUnknownThis) BindSelf_(s any) {
	_ = s
	this_.ConformanceRules_AbstractTypeRestrictionRuleBase.BindSelf_(s)
}
func (this_ *ConformanceRules_BanUnknownThis) IsJavaNil_() bool { return this_ == nil }
func (*ConformanceRules_BanUnknownThis) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConformanceRules$BanUnknownThis", "BanUnknownThis"
}
func (this_ *ConformanceRules_BanUnknownThis) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ConformanceRules_AbstractTypeRestrictionRuleBase.IdentityHash_()
}
func (*ConformanceRules_BanUnknownThis) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConformanceRules.BanUnknownThis", "com.google.javascript.jscomp.ConformanceRules$BanUnknownThis", func(v_ any) bool {
		return j2gort.InstanceOf[*ConformanceRules_BanUnknownThis](v_)
	})
}
func (*ConformanceRules_BanUnknownThis) Implements_jscomp_ConformanceRules_AbstractTypeRestrictionRule() {
}
func (*ConformanceRules_BanUnknownThis) Implements_jscomp_ConformanceRules_AbstractRule() {}
func (*ConformanceRules_BanUnknownThis) Implements_jscomp_CheckConformance_Rule()         {}

func (this_ *ConformanceRules_BanUnknownThis) initFields_ConformanceRules_BanUnknownThis() {
	_ = this_
	this_.reports = pkg_base.NewLinkedIdentityHashSet[pkg_rhino.Node]()
}

// Java constructor: BanUnknownThis(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.Requirement requirement)
func NewConformanceRules_BanUnknownThis(compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) *ConformanceRules_BanUnknownThis {
	this_ := &ConformanceRules_BanUnknownThis{}
	this_.BindSelf_(this_)
	InitConformanceRules_BanUnknownThis(this_, compiler, requirement)
	return this_
}

func InitConformanceRules_BanUnknownThis(this_ *ConformanceRules_BanUnknownThis, compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) {
	_ = compiler
	_ = requirement
	InitConformanceRules_AbstractTypeRestrictionRule(&this_.ConformanceRules_AbstractTypeRestrictionRuleBase, compiler, requirement)
	this_.initFields_ConformanceRules_BanUnknownThis()
}

// Java method: BanUnknownThis.tsIsAllowlisted()
func (this_ *ConformanceRules_BanUnknownThis) TsIsAllowlisted() bool {
	return true
}

// Java method: BanUnknownThis.checkConformance(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *ConformanceRules_BanUnknownThis) CheckConformance_2(t *NodeTraversal, n pkg_rhino.Node) *ConformanceRules_ConformanceResult {
	_ = t
	_ = n
	if n.IsThis() {
		var type_ pkg_rhino.JSType = n.GetJSType()
		_ = type_
		if ((type_ != nil && !type_.IsJavaNil_()) && type_.IsUnknownType()) && !(this_.IsTypeImmediatelyTightened(n)) {
			var root pkg_rhino.Node = t.GetScopeRoot()
			_ = root
			if !(this_.reports.Contains(root)) {
				this_.reports.Add(root)
				return StaticGet_ConformanceRules_ConformanceResult_VIOLATION()
			}
		}
	}
	return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
}

// Java class: com.google.javascript.jscomp.ConformanceRules.BanUnknownDirectThisPropsReferences
type ConformanceRules_BanUnknownDirectThisPropsReferences struct {
	ConformanceRules_AbstractTypeRestrictionRuleBase
}

func (this_ *ConformanceRules_BanUnknownDirectThisPropsReferences) BindSelf_(s any) {
	_ = s
	this_.ConformanceRules_AbstractTypeRestrictionRuleBase.BindSelf_(s)
}
func (this_ *ConformanceRules_BanUnknownDirectThisPropsReferences) IsJavaNil_() bool {
	return this_ == nil
}
func (*ConformanceRules_BanUnknownDirectThisPropsReferences) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConformanceRules$BanUnknownDirectThisPropsReferences", "BanUnknownDirectThisPropsReferences"
}
func (this_ *ConformanceRules_BanUnknownDirectThisPropsReferences) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ConformanceRules_AbstractTypeRestrictionRuleBase.IdentityHash_()
}
func (*ConformanceRules_BanUnknownDirectThisPropsReferences) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConformanceRules.BanUnknownDirectThisPropsReferences", "com.google.javascript.jscomp.ConformanceRules$BanUnknownDirectThisPropsReferences", func(v_ any) bool {
		return j2gort.InstanceOf[*ConformanceRules_BanUnknownDirectThisPropsReferences](v_)
	})
}
func (*ConformanceRules_BanUnknownDirectThisPropsReferences) Implements_jscomp_ConformanceRules_AbstractTypeRestrictionRule() {
}
func (*ConformanceRules_BanUnknownDirectThisPropsReferences) Implements_jscomp_ConformanceRules_AbstractRule() {
}
func (*ConformanceRules_BanUnknownDirectThisPropsReferences) Implements_jscomp_CheckConformance_Rule() {
}

func (this_ *ConformanceRules_BanUnknownDirectThisPropsReferences) initFields_ConformanceRules_BanUnknownDirectThisPropsReferences() {
	_ = this_
}

// Java constructor: BanUnknownDirectThisPropsReferences(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.Requirement requirement)
func NewConformanceRules_BanUnknownDirectThisPropsReferences(compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) *ConformanceRules_BanUnknownDirectThisPropsReferences {
	this_ := &ConformanceRules_BanUnknownDirectThisPropsReferences{}
	this_.BindSelf_(this_)
	InitConformanceRules_BanUnknownDirectThisPropsReferences(this_, compiler, requirement)
	return this_
}

func InitConformanceRules_BanUnknownDirectThisPropsReferences(this_ *ConformanceRules_BanUnknownDirectThisPropsReferences, compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) {
	_ = compiler
	_ = requirement
	InitConformanceRules_AbstractTypeRestrictionRule(&this_.ConformanceRules_AbstractTypeRestrictionRuleBase, compiler, requirement)
	this_.initFields_ConformanceRules_BanUnknownDirectThisPropsReferences()
}

// Java method: BanUnknownDirectThisPropsReferences.checkConformance(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *ConformanceRules_BanUnknownDirectThisPropsReferences) CheckConformance_2(t *NodeTraversal, n pkg_rhino.Node) *ConformanceRules_ConformanceResult {
	_ = t
	_ = n
	if (((((n.IsGetProp() && ConformanceRules_BanUnknownDirectThisPropsReferences_IsKnownThis(n.GetFirstChild())) && ConformanceRules_AbstractTypeRestrictionRule_IsUnknown(n)) && !(ConformanceRules_AbstractTypeRestrictionRule_IsTemplateType(n))) && this_.IsUsed(n)) && !(this_.IsTypeImmediatelyTightened(n))) && !(ConformanceRules_BanUnknownDirectThisPropsReferences_IsExplicitlyUnknown(n)) {
		return StaticGet_ConformanceRules_ConformanceResult_VIOLATION()
	}
	return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
}

// Java static method: BanUnknownDirectThisPropsReferences.isKnownThis(com.google.javascript.rhino.Node n)
func ConformanceRules_BanUnknownDirectThisPropsReferences_IsKnownThis(n pkg_rhino.Node) bool {
	_ = n
	return (n.IsThis() && !(ConformanceRules_AbstractTypeRestrictionRule_IsUnknown(n)))
}

// Java static method: BanUnknownDirectThisPropsReferences.isExplicitlyUnknown(com.google.javascript.rhino.Node n)
func ConformanceRules_BanUnknownDirectThisPropsReferences_IsExplicitlyUnknown(n pkg_rhino.Node) bool {
	_ = n
	var owner pkg_rhino.ObjectType = pkg_rhino.ObjectType_Cast(n.GetFirstChild().GetJSType())
	_ = owner
	var prop *pkg_rhino.Property
	_ = prop
	if owner != nil && !owner.IsJavaNil_() {
		prop = owner.GetSlot_1_string(n.GetString())
	} else {
		prop = nil
	}
	return ((prop != nil) && !(prop.IsTypeInferred()))
}

// Java class: com.google.javascript.jscomp.ConformanceRules.BanUnknownTypedClassPropsReferences
type ConformanceRules_BanUnknownTypedClassPropsReferences struct {
	ConformanceRules_AbstractTypeRestrictionRuleBase
}

func (this_ *ConformanceRules_BanUnknownTypedClassPropsReferences) BindSelf_(s any) {
	_ = s
	this_.ConformanceRules_AbstractTypeRestrictionRuleBase.BindSelf_(s)
}
func (this_ *ConformanceRules_BanUnknownTypedClassPropsReferences) IsJavaNil_() bool {
	return this_ == nil
}
func (*ConformanceRules_BanUnknownTypedClassPropsReferences) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConformanceRules$BanUnknownTypedClassPropsReferences", "BanUnknownTypedClassPropsReferences"
}
func (this_ *ConformanceRules_BanUnknownTypedClassPropsReferences) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ConformanceRules_AbstractTypeRestrictionRuleBase.IdentityHash_()
}
func (*ConformanceRules_BanUnknownTypedClassPropsReferences) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConformanceRules.BanUnknownTypedClassPropsReferences", "com.google.javascript.jscomp.ConformanceRules$BanUnknownTypedClassPropsReferences", func(v_ any) bool {
		return j2gort.InstanceOf[*ConformanceRules_BanUnknownTypedClassPropsReferences](v_)
	})
}
func (*ConformanceRules_BanUnknownTypedClassPropsReferences) Implements_jscomp_ConformanceRules_AbstractTypeRestrictionRule() {
}
func (*ConformanceRules_BanUnknownTypedClassPropsReferences) Implements_jscomp_ConformanceRules_AbstractRule() {
}
func (*ConformanceRules_BanUnknownTypedClassPropsReferences) Implements_jscomp_CheckConformance_Rule() {
}

func (this_ *ConformanceRules_BanUnknownTypedClassPropsReferences) initFields_ConformanceRules_BanUnknownTypedClassPropsReferences() {
	_ = this_
}

// Java constructor: BanUnknownTypedClassPropsReferences(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.Requirement requirement)
func NewConformanceRules_BanUnknownTypedClassPropsReferences(compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) *ConformanceRules_BanUnknownTypedClassPropsReferences {
	this_ := &ConformanceRules_BanUnknownTypedClassPropsReferences{}
	this_.BindSelf_(this_)
	InitConformanceRules_BanUnknownTypedClassPropsReferences(this_, compiler, requirement)
	return this_
}

func InitConformanceRules_BanUnknownTypedClassPropsReferences(this_ *ConformanceRules_BanUnknownTypedClassPropsReferences, compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) {
	_ = compiler
	_ = requirement
	InitConformanceRules_AbstractTypeRestrictionRule(&this_.ConformanceRules_AbstractTypeRestrictionRuleBase, compiler, requirement)
	this_.initFields_ConformanceRules_BanUnknownTypedClassPropsReferences()
}

// Java method: BanUnknownTypedClassPropsReferences.tsIsAllowlisted()
func (this_ *ConformanceRules_BanUnknownTypedClassPropsReferences) TsIsAllowlisted() bool {
	return true
}

// Java method: BanUnknownTypedClassPropsReferences.checkConformance(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node getprop)
func (this_ *ConformanceRules_BanUnknownTypedClassPropsReferences) CheckConformance_2(t *NodeTraversal, getprop pkg_rhino.Node) *ConformanceRules_ConformanceResult {
	_ = t
	_ = getprop
	if (((((getprop.IsGetProp() && ConformanceRules_AbstractTypeRestrictionRule_IsUnknown(getprop)) && this_.IsUsed(getprop)) && !(this_.IsTypeImmediatelyTightened(getprop))) && this_.priv_BanUnknownTypedClassPropsReferences_isCheckablePropertySource_2158(getprop.GetFirstChild())) && !(ConformanceRules_AbstractTypeRestrictionRule_IsTemplateType(getprop))) && !(this_.priv_BanUnknownTypedClassPropsReferences_isDeclaredUnknown_2160(getprop)) {
		var propName string = getprop.GetString()
		_ = propName
		var typeName string = getprop.GetFirstChild().GetJSType().String()
		_ = typeName
		return NewConformanceRules_ConformanceResult_2(StaticGet_ConformanceRules_ConformanceLevel_VIOLATION(), ("The property \"" + j2gort.StrValString(propName) + "\" on type \"" + j2gort.StrValString(typeName) + "\""))
	}
	return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
}

// Java method: BanUnknownTypedClassPropsReferences.isCheckablePropertySource(com.google.javascript.rhino.Node n)
func (this_ *ConformanceRules_BanUnknownTypedClassPropsReferences) priv_BanUnknownTypedClassPropsReferences_isCheckablePropertySource_2158(n pkg_rhino.Node) bool {
	_ = n
	return ((((ConformanceRules_AbstractTypeRestrictionRule_IsKnown(n) && !(ConformanceRules_AbstractTypeRestrictionRule_IsTop(n))) && this_.priv_BanUnknownTypedClassPropsReferences_isClassType_2159(n)) && !(this_.IsNativeObjectType(n))) && !(this_.IsAllowlistedType(n)))
}

// Java method: BanUnknownTypedClassPropsReferences.isClassType(com.google.javascript.rhino.Node n)
func (this_ *ConformanceRules_BanUnknownTypedClassPropsReferences) priv_BanUnknownTypedClassPropsReferences_isClassType_2159(n pkg_rhino.Node) bool {
	_ = n
	var type_ pkg_rhino.ObjectType = n.GetJSType().RestrictByNotNullOrUndefined().ToMaybeObjectType()
	_ = type_
	if (type_ == nil || type_.IsJavaNil_()) || !(type_.IsInstanceType()) {
		return false
	}
	var ctor *pkg_rhino.FunctionType = type_.GetConstructor()
	_ = ctor
	if ctor == nil {
		return false
	}
	var info *pkg_rhino.JSDocInfo = ctor.GetJSDocInfo()
	_ = info
	var source pkg_rhino.Node = ctor.GetSource()
	_ = source
	return (((info != nil) && info.IsConstructorOrInterface()) || ((source != 0) && source.IsClass()))
}

// Java method: BanUnknownTypedClassPropsReferences.isDeclaredUnknown(com.google.javascript.rhino.Node n)
func (this_ *ConformanceRules_BanUnknownTypedClassPropsReferences) priv_BanUnknownTypedClassPropsReferences_isDeclaredUnknown_2160(n pkg_rhino.Node) bool {
	_ = n
	var target pkg_rhino.Node = n.GetFirstChild()
	_ = target
	var targetType pkg_rhino.ObjectType = target.GetJSType().RestrictByNotNullOrUndefined().ToMaybeObjectType()
	_ = targetType
	if targetType == nil || targetType.IsJavaNil_() {
		return false
	}
	var info *pkg_rhino.JSDocInfo = targetType.GetPropertyJSDocInfo_1_string(n.GetString())
	_ = info
	if (info == nil) || !(info.HasType()) {
		return false
	}
	var expr pkg_rhino.JSTypeExpression = info.GetType()
	_ = expr
	var typeExprNode pkg_rhino.Node = expr.GetRoot()
	_ = typeExprNode
	if (typeExprNode.RawToken_() == uint8(158)) && !(typeExprNode.HasChildren()) {
		return true
	} else if typeExprNode.RawToken_() == uint8(155) {
		{
			var child pkg_rhino.Node = typeExprNode.GetFirstChild()
			_ = child
			for ; child != 0; child = child.GetNext() {
				if child.RawToken_() == uint8(158) {
					return true
				}
			}
		}
	}
	return false
}

// Java class: com.google.javascript.jscomp.ConformanceRules.BanUnresolvedType
type ConformanceRules_BanUnresolvedType struct {
	ConformanceRules_AbstractTypeRestrictionRuleBase
}

func (this_ *ConformanceRules_BanUnresolvedType) BindSelf_(s any) {
	_ = s
	this_.ConformanceRules_AbstractTypeRestrictionRuleBase.BindSelf_(s)
}
func (this_ *ConformanceRules_BanUnresolvedType) IsJavaNil_() bool { return this_ == nil }
func (*ConformanceRules_BanUnresolvedType) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConformanceRules$BanUnresolvedType", "BanUnresolvedType"
}
func (this_ *ConformanceRules_BanUnresolvedType) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ConformanceRules_AbstractTypeRestrictionRuleBase.IdentityHash_()
}
func (*ConformanceRules_BanUnresolvedType) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConformanceRules.BanUnresolvedType", "com.google.javascript.jscomp.ConformanceRules$BanUnresolvedType", func(v_ any) bool {
		return j2gort.InstanceOf[*ConformanceRules_BanUnresolvedType](v_)
	})
}
func (*ConformanceRules_BanUnresolvedType) Implements_jscomp_ConformanceRules_AbstractTypeRestrictionRule() {
}
func (*ConformanceRules_BanUnresolvedType) Implements_jscomp_ConformanceRules_AbstractRule() {}
func (*ConformanceRules_BanUnresolvedType) Implements_jscomp_CheckConformance_Rule()         {}

func (this_ *ConformanceRules_BanUnresolvedType) initFields_ConformanceRules_BanUnresolvedType() {
	_ = this_
}

// Java constructor: BanUnresolvedType(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.Requirement requirement)
func NewConformanceRules_BanUnresolvedType(compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) *ConformanceRules_BanUnresolvedType {
	this_ := &ConformanceRules_BanUnresolvedType{}
	this_.BindSelf_(this_)
	InitConformanceRules_BanUnresolvedType(this_, compiler, requirement)
	return this_
}

func InitConformanceRules_BanUnresolvedType(this_ *ConformanceRules_BanUnresolvedType, compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) {
	_ = compiler
	_ = requirement
	InitConformanceRules_AbstractTypeRestrictionRule(&this_.ConformanceRules_AbstractTypeRestrictionRuleBase, compiler, requirement)
	this_.initFields_ConformanceRules_BanUnresolvedType()
}

// Java method: BanUnresolvedType.checkConformance(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *ConformanceRules_BanUnresolvedType) CheckConformance_2(t *NodeTraversal, n pkg_rhino.Node) *ConformanceRules_ConformanceResult {
	_ = t
	_ = n
	if n.IsGetProp() {
		var target pkg_rhino.Node = n.GetFirstChild()
		_ = target
		var type_ pkg_rhino.JSType = target.GetJSType()
		_ = type_
		var nonConformingPart string = ConformanceRules_BanUnresolvedType_GetNonConformingPart(type_)
		_ = nonConformingPart
		if (nonConformingPart != j2gort.NullString) && !(this_.IsTypeImmediatelyTightened(n)) {
			return NewConformanceRules_ConformanceResult_2(StaticGet_ConformanceRules_ConformanceLevel_VIOLATION(), ("Reference to type '" + j2gort.StrValString(nonConformingPart) + "' never resolved."))
		}
	}
	return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
}

// Java static method: BanUnresolvedType.getNonConformingPart(com.google.javascript.rhino.jstype.JSType type)
func ConformanceRules_BanUnresolvedType_GetNonConformingPart(type_ pkg_rhino.JSType) string {
	_ = type_
	if type_ == nil || type_.IsJavaNil_() {
		return j2gort.NullString
	}
	if type_.IsUnionType() {
		var nonConformingParts j2gort.Ext_java_util_ArrayList = nil
		_ = nonConformingParts
		{
			cur_37 := j2gort.NewIterCursor(j2gort.CheckNotNull(type_.GetUnionMembers()))
			for cur_37.HasNext() {
				rawElem_38 := cur_37.Next()
				part := j2gort.AsIface[pkg_rhino.JSType](rawElem_38)
				_ = part
				var nonConformingPart string = ConformanceRules_BanUnresolvedType_GetNonConformingPart(part)
				_ = nonConformingPart
				if nonConformingPart != j2gort.NullString {
					if nonConformingParts != nil {
						nonConformingParts = nonConformingParts
					} else {
						nonConformingParts = j2gort.NewExt_java_util_ArrayList()
					}
					nonConformingParts.Add_1(j2gort.BoxString(nonConformingPart))
				}
			}
		}
		if nonConformingParts != nil {
			return j2gort.Ext_com_google_common_base_Joiner_On_1_uint16(uint16(124)).Join_1_Ext_java_lang_Iterable(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](nonConformingParts))
		}
	} else if type_.IsNoResolvedType() {
		var noResolvedType pkg_rhino.ObjectType = type_.ToObjectType()
		_ = noResolvedType
		return noResolvedType.GetReferenceName()
	}
	return j2gort.NullString
}

// Java class: com.google.javascript.jscomp.ConformanceRules.StrictBanUnresolvedType
type ConformanceRules_StrictBanUnresolvedType struct {
	ConformanceRules_AbstractTypeRestrictionRuleBase
}

func (this_ *ConformanceRules_StrictBanUnresolvedType) BindSelf_(s any) {
	_ = s
	this_.ConformanceRules_AbstractTypeRestrictionRuleBase.BindSelf_(s)
}
func (this_ *ConformanceRules_StrictBanUnresolvedType) IsJavaNil_() bool { return this_ == nil }
func (*ConformanceRules_StrictBanUnresolvedType) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConformanceRules$StrictBanUnresolvedType", "StrictBanUnresolvedType"
}
func (this_ *ConformanceRules_StrictBanUnresolvedType) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ConformanceRules_AbstractTypeRestrictionRuleBase.IdentityHash_()
}
func (*ConformanceRules_StrictBanUnresolvedType) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConformanceRules.StrictBanUnresolvedType", "com.google.javascript.jscomp.ConformanceRules$StrictBanUnresolvedType", func(v_ any) bool {
		return j2gort.InstanceOf[*ConformanceRules_StrictBanUnresolvedType](v_)
	})
}
func (*ConformanceRules_StrictBanUnresolvedType) Implements_jscomp_ConformanceRules_AbstractTypeRestrictionRule() {
}
func (*ConformanceRules_StrictBanUnresolvedType) Implements_jscomp_ConformanceRules_AbstractRule() {}
func (*ConformanceRules_StrictBanUnresolvedType) Implements_jscomp_CheckConformance_Rule()         {}

func (this_ *ConformanceRules_StrictBanUnresolvedType) initFields_ConformanceRules_StrictBanUnresolvedType() {
	_ = this_
}

// Java constructor: StrictBanUnresolvedType(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.Requirement requirement)
func NewConformanceRules_StrictBanUnresolvedType(compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) *ConformanceRules_StrictBanUnresolvedType {
	this_ := &ConformanceRules_StrictBanUnresolvedType{}
	this_.BindSelf_(this_)
	InitConformanceRules_StrictBanUnresolvedType(this_, compiler, requirement)
	return this_
}

func InitConformanceRules_StrictBanUnresolvedType(this_ *ConformanceRules_StrictBanUnresolvedType, compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) {
	_ = compiler
	_ = requirement
	InitConformanceRules_AbstractTypeRestrictionRule(&this_.ConformanceRules_AbstractTypeRestrictionRuleBase, compiler, requirement)
	this_.initFields_ConformanceRules_StrictBanUnresolvedType()
}

// Java method: StrictBanUnresolvedType.checkConformance(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *ConformanceRules_StrictBanUnresolvedType) CheckConformance_2(t *NodeTraversal, n pkg_rhino.Node) *ConformanceRules_ConformanceResult {
	_ = t
	_ = n
	var nonConformingPart string = ConformanceRules_BanUnresolvedType_GetNonConformingPart(n.GetJSType())
	_ = nonConformingPart
	if (nonConformingPart != j2gort.NullString) && !(this_.IsTypeImmediatelyTightened(n)) {
		return NewConformanceRules_ConformanceResult_2(StaticGet_ConformanceRules_ConformanceLevel_VIOLATION(), ("Reference to type '" + j2gort.StrValString(nonConformingPart) + "' never resolved."))
	}
	return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
}

// Java class: com.google.javascript.jscomp.ConformanceRules.BanGlobalVars
type ConformanceRules_BanGlobalVars struct {
	ConformanceRules_AbstractRuleBase
	allowlistedNames j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"allowlistedNames"`
}

func (this_ *ConformanceRules_BanGlobalVars) BindSelf_(s any) {
	_ = s
	this_.ConformanceRules_AbstractRuleBase.BindSelf_(s)
}
func (this_ *ConformanceRules_BanGlobalVars) IsJavaNil_() bool { return this_ == nil }
func (*ConformanceRules_BanGlobalVars) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConformanceRules$BanGlobalVars", "BanGlobalVars"
}
func (this_ *ConformanceRules_BanGlobalVars) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ConformanceRules_AbstractRuleBase.IdentityHash_()
}
func (*ConformanceRules_BanGlobalVars) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConformanceRules.BanGlobalVars", "com.google.javascript.jscomp.ConformanceRules$BanGlobalVars", func(v_ any) bool {
		return j2gort.InstanceOf[*ConformanceRules_BanGlobalVars](v_)
	})
}
func (*ConformanceRules_BanGlobalVars) Implements_jscomp_ConformanceRules_AbstractRule() {}
func (*ConformanceRules_BanGlobalVars) Implements_jscomp_CheckConformance_Rule()         {}

func (this_ *ConformanceRules_BanGlobalVars) initFields_ConformanceRules_BanGlobalVars() {
	_ = this_
}

// Java constructor: BanGlobalVars(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.Requirement requirement)
func NewConformanceRules_BanGlobalVars(compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) *ConformanceRules_BanGlobalVars {
	this_ := &ConformanceRules_BanGlobalVars{}
	this_.BindSelf_(this_)
	InitConformanceRules_BanGlobalVars(this_, compiler, requirement)
	return this_
}

func InitConformanceRules_BanGlobalVars(this_ *ConformanceRules_BanGlobalVars, compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) {
	_ = compiler
	_ = requirement
	InitConformanceRules_AbstractRule(&this_.ConformanceRules_AbstractRuleBase, compiler, requirement)
	this_.initFields_ConformanceRules_BanGlobalVars()
	this_.allowlistedNames = j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](requirement.GetValueList()))
}

// Java method: BanGlobalVars.checkConformance(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *ConformanceRules_BanGlobalVars) CheckConformance_2(t *NodeTraversal, n pkg_rhino.Node) *ConformanceRules_ConformanceResult {
	_ = t
	_ = n
	if ((t.InGlobalScope() && NodeUtil_IsDeclaration(n)) && !(n.GetBooleanProp(pkg_rhino.StaticGet_Node_IS_NAMESPACE()))) && !(this_.priv_BanGlobalVars_isAllowlisted_2161(n)) {
		var enclosingScript pkg_rhino.Node = NodeUtil_GetEnclosingScript(n)
		_ = enclosingScript
		if (enclosingScript != 0) && ((enclosingScript.GetBooleanProp(pkg_rhino.StaticGet_Node_GOOG_MODULE()) || enclosingScript.GetBooleanProp(pkg_rhino.StaticGet_Node_ES6_MODULE())) || enclosingScript.GetInputId().Equals(j2gort.PtrToAny(this_.GetConformanceRules_AbstractRuleBase().compiler.GetSyntheticCodeInputId()))) {
			return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
		}
		return StaticGet_ConformanceRules_ConformanceResult_VIOLATION()
	}
	return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
}

// Java method: BanGlobalVars.isAllowlisted(com.google.javascript.rhino.Node n)
func (this_ *ConformanceRules_BanGlobalVars) priv_BanGlobalVars_isAllowlisted_2161(n pkg_rhino.Node) bool {
	_ = n
	if n.IsFromExterns() {
		return true
	}
	if n.IsFunction() {
		return this_.priv_BanGlobalVars_isAllowlistedName_2162(n.GetFirstChild().GetString())
	}
	if NodeUtil_IsNameDeclaration(n) {
		var allowlisted []bool = []bool{true}
		_ = allowlisted
		NodeUtil_VisitLhsNodesInNode(n, j2gort.NewLambda_Ext_java_util_function_Consumer(func(lp_39 any) {
			name := j2gort.CastNode[pkg_rhino.Node](lp_39)
			_ = name
			if !(this_.priv_BanGlobalVars_isAllowlistedName_2162(name.GetString())) {
				{
					a_40 := allowlisted
					i_41 := int32(0)
					a_40[i_41] = false
				}
			}
		}))
		return allowlisted[int32(0)]
	}
	return false
}

// Java method: BanGlobalVars.isAllowlistedName(java.lang.String name)
func (this_ *ConformanceRules_BanGlobalVars) priv_BanGlobalVars_isAllowlistedName_2162(name string) bool {
	_ = name
	if j2gort.SetContainsStr(this_.allowlistedNames, name) {
		return true
	}
	return (((((j2gort.CheckNotNullString(name) == "$jscomp") || javalang.StringStartsWith(j2gort.CheckNotNullString(name), "$jscomp$compprop")) || ClosureRewriteModule_IsModuleContent(name)) || ClosureRewriteModule_IsModuleExport(name)) || ScopedAliases_IsScopedAliases(name))
}

// Java class: com.google.javascript.jscomp.ConformanceRules.BannedEnhance
type ConformanceRules_BannedEnhance struct {
	ConformanceRules_AbstractRuleBase
	bannedEnhancedNamespaces j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"bannedEnhancedNamespaces"`
}

func (this_ *ConformanceRules_BannedEnhance) BindSelf_(s any) {
	_ = s
	this_.ConformanceRules_AbstractRuleBase.BindSelf_(s)
}
func (this_ *ConformanceRules_BannedEnhance) IsJavaNil_() bool { return this_ == nil }
func (*ConformanceRules_BannedEnhance) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConformanceRules$BannedEnhance", "BannedEnhance"
}
func (this_ *ConformanceRules_BannedEnhance) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ConformanceRules_AbstractRuleBase.IdentityHash_()
}
func (*ConformanceRules_BannedEnhance) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConformanceRules.BannedEnhance", "com.google.javascript.jscomp.ConformanceRules$BannedEnhance", func(v_ any) bool {
		return j2gort.InstanceOf[*ConformanceRules_BannedEnhance](v_)
	})
}
func (*ConformanceRules_BannedEnhance) Implements_jscomp_ConformanceRules_AbstractRule() {}
func (*ConformanceRules_BannedEnhance) Implements_jscomp_CheckConformance_Rule()         {}

func (this_ *ConformanceRules_BannedEnhance) initFields_ConformanceRules_BannedEnhance() {
	_ = this_
}

// Java constructor: BannedEnhance(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.Requirement requirement)
func NewConformanceRules_BannedEnhance(compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) *ConformanceRules_BannedEnhance {
	this_ := &ConformanceRules_BannedEnhance{}
	this_.BindSelf_(this_)
	InitConformanceRules_BannedEnhance(this_, compiler, requirement)
	return this_
}

func InitConformanceRules_BannedEnhance(this_ *ConformanceRules_BannedEnhance, compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) {
	_ = compiler
	_ = requirement
	InitConformanceRules_AbstractRule(&this_.ConformanceRules_AbstractRuleBase, compiler, requirement)
	this_.initFields_ConformanceRules_BannedEnhance()
	if requirement.GetValueCount() == int32(0) {
		panic(NewCheckConformance_InvalidRequirementSpec_1("missing value"))
	}
	this_.bannedEnhancedNamespaces = j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](requirement.GetValueList()))
}

// Java method: BannedEnhance.checkConformance(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *ConformanceRules_BannedEnhance) CheckConformance_2(t *NodeTraversal, n pkg_rhino.Node) *ConformanceRules_ConformanceResult {
	_ = t
	_ = n
	var docInfo *pkg_rhino.JSDocInfo = n.GetJSDocInfo()
	_ = docInfo
	if (docInfo == nil) || !(docInfo.HasEnhance()) {
		return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
	}
	{
		cur_42 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.bannedEnhancedNamespaces))
		for cur_42.HasNext() {
			banned := cur_42.NextStr()
			_ = banned
			if j2gort.CheckNotNullString(docInfo.GetEnhance()) == banned {
				return NewConformanceRules_ConformanceResult_2(StaticGet_ConformanceRules_ConformanceLevel_VIOLATION(), ("The enhanced namespace \"" + j2gort.StrValString(banned) + "\""))
			}
		}
	}
	return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
}

// Java method: BannedEnhance.getPrecondition()
func (this_ *ConformanceRules_BannedEnhance) GetPrecondition() CheckConformance_Precondition {
	return lambdaSingleton_conformance_rules_31
}

// Java class: com.google.javascript.jscomp.ConformanceRules.BannedModsRegex
type ConformanceRules_BannedModsRegex struct {
	ConformanceRules_AbstractRuleBase
	bannedModsRegex *j2gort.Ext_java_util_regex_Pattern `gson:"bannedModsRegex"`
}

func (this_ *ConformanceRules_BannedModsRegex) BindSelf_(s any) {
	_ = s
	this_.ConformanceRules_AbstractRuleBase.BindSelf_(s)
}
func (this_ *ConformanceRules_BannedModsRegex) IsJavaNil_() bool { return this_ == nil }
func (*ConformanceRules_BannedModsRegex) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConformanceRules$BannedModsRegex", "BannedModsRegex"
}
func (this_ *ConformanceRules_BannedModsRegex) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ConformanceRules_AbstractRuleBase.IdentityHash_()
}
func (*ConformanceRules_BannedModsRegex) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConformanceRules.BannedModsRegex", "com.google.javascript.jscomp.ConformanceRules$BannedModsRegex", func(v_ any) bool {
		return j2gort.InstanceOf[*ConformanceRules_BannedModsRegex](v_)
	})
}
func (*ConformanceRules_BannedModsRegex) Implements_jscomp_ConformanceRules_AbstractRule() {}
func (*ConformanceRules_BannedModsRegex) Implements_jscomp_CheckConformance_Rule()         {}

func (this_ *ConformanceRules_BannedModsRegex) initFields_ConformanceRules_BannedModsRegex() {
	_ = this_
}

// Java constructor: BannedModsRegex(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.Requirement requirement)
func NewConformanceRules_BannedModsRegex(compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) *ConformanceRules_BannedModsRegex {
	this_ := &ConformanceRules_BannedModsRegex{}
	this_.BindSelf_(this_)
	InitConformanceRules_BannedModsRegex(this_, compiler, requirement)
	return this_
}

func InitConformanceRules_BannedModsRegex(this_ *ConformanceRules_BannedModsRegex, compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) {
	_ = compiler
	_ = requirement
	InitConformanceRules_AbstractRule(&this_.ConformanceRules_AbstractRuleBase, compiler, requirement)
	this_.initFields_ConformanceRules_BannedModsRegex()
	var bannedModsRegexList j2gort.Ext_java_util_List = j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](requirement.GetValueList())
	_ = bannedModsRegexList
	if requirement.GetValueCount() == int32(0) {
		panic(NewCheckConformance_InvalidRequirementSpec_1("missing value"))
	}
	this_.bannedModsRegex = ConformanceRules_BuildPattern(bannedModsRegexList)
}

// Java method: BannedModsRegex.checkConformance(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *ConformanceRules_BannedModsRegex) CheckConformance_2(t *NodeTraversal, n pkg_rhino.Node) *ConformanceRules_ConformanceResult {
	_ = t
	_ = n
	var docInfo *pkg_rhino.JSDocInfo = n.GetJSDocInfo()
	_ = docInfo
	if (docInfo == nil) || !(docInfo.HasMods()) {
		return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
	}
	if this_.bannedModsRegex.Matcher(j2gort.BoxString(docInfo.GetMods())).Find_0() {
		return StaticGet_ConformanceRules_ConformanceResult_VIOLATION()
	}
	return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
}

// Java method: BannedModsRegex.getPrecondition()
func (this_ *ConformanceRules_BannedModsRegex) GetPrecondition() CheckConformance_Precondition {
	return lambdaSingleton_conformance_rules_32
}

// Java class: com.google.javascript.jscomp.ConformanceRules.BanCreateElement
type ConformanceRules_BanCreateElement struct {
	ConformanceRules_AbstractRuleBase
	bannedTags    j2gort.OrderedSet[string] `gson:"bannedTags"`
	domHelperType pkg_rhino.JSType          `gson:"domHelperType"`
	documentType  pkg_rhino.JSType          `gson:"documentType"`
}

func (this_ *ConformanceRules_BanCreateElement) BindSelf_(s any) {
	_ = s
	this_.ConformanceRules_AbstractRuleBase.BindSelf_(s)
}
func (this_ *ConformanceRules_BanCreateElement) IsJavaNil_() bool { return this_ == nil }
func (*ConformanceRules_BanCreateElement) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConformanceRules$BanCreateElement", "BanCreateElement"
}
func (this_ *ConformanceRules_BanCreateElement) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ConformanceRules_AbstractRuleBase.IdentityHash_()
}
func (*ConformanceRules_BanCreateElement) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConformanceRules.BanCreateElement", "com.google.javascript.jscomp.ConformanceRules$BanCreateElement", func(v_ any) bool {
		return j2gort.InstanceOf[*ConformanceRules_BanCreateElement](v_)
	})
}
func (*ConformanceRules_BanCreateElement) Implements_jscomp_ConformanceRules_AbstractRule() {}
func (*ConformanceRules_BanCreateElement) Implements_jscomp_CheckConformance_Rule()         {}

func (this_ *ConformanceRules_BanCreateElement) initFields_ConformanceRules_BanCreateElement() {
	_ = this_
}

// Java constructor: BanCreateElement(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.Requirement requirement)
func NewConformanceRules_BanCreateElement(compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) *ConformanceRules_BanCreateElement {
	this_ := &ConformanceRules_BanCreateElement{}
	this_.BindSelf_(this_)
	InitConformanceRules_BanCreateElement(this_, compiler, requirement)
	return this_
}

func InitConformanceRules_BanCreateElement(this_ *ConformanceRules_BanCreateElement, compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) {
	_ = compiler
	_ = requirement
	InitConformanceRules_AbstractRule(&this_.ConformanceRules_AbstractRuleBase, compiler, requirement)
	this_.initFields_ConformanceRules_BanCreateElement()
	this_.bannedTags = j2gort.NewOrderedSet[string]()
	{
		cur_43 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](requirement.GetValueList())))
		for cur_43.HasNext() {
			value := cur_43.NextStr()
			_ = value
			this_.bannedTags.Add(j2gort.Ext_com_google_common_base_Ascii_ToLowerCase_1_string(value))
		}
	}
	if this_.bannedTags.IsEmpty() {
		panic(NewCheckConformance_InvalidRequirementSpec_1("Specify one or more values."))
	}
	this_.domHelperType = compiler.GetTypeRegistry().GetGlobalType("goog.dom.DomHelper")
	this_.documentType = compiler.GetTypeRegistry().GetGlobalType("Document")
}

// Java method: BanCreateElement.checkConformance(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *ConformanceRules_BanCreateElement) CheckConformance_2(t *NodeTraversal, n pkg_rhino.Node) *ConformanceRules_ConformanceResult {
	_ = t
	_ = n
	if n.IsCall() {
		var tag pkg_rhino.Node = n.GetSecondChild()
		_ = tag
		if ((tag != 0) && tag.IsStringLit()) && this_.bannedTags.Contains(j2gort.Ext_com_google_common_base_Ascii_ToLowerCase_1_string(tag.GetString())) {
			return this_.priv_BanCreateElement_checkCreateElement_2163(n)
		}
	}
	return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
}

// Java method: BanCreateElement.checkCreateElement(com.google.javascript.rhino.Node n)
func (this_ *ConformanceRules_BanCreateElement) priv_BanCreateElement_checkCreateElement_2163(n pkg_rhino.Node) *ConformanceRules_ConformanceResult {
	_ = n
	var target pkg_rhino.Node = n.GetFirstChild()
	_ = target
	if !(target.IsGetProp()) {
		return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
	}
	var functionName string = target.GetString()
	_ = functionName
	if !("createElement" == functionName) && !("createDom" == functionName) {
		return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
	}
	var srcObj pkg_rhino.Node = target.GetFirstChild()
	_ = srcObj
	if StaticGet_ConformanceRules_GOOG_DOM().Matches(srcObj) {
		return StaticGet_ConformanceRules_ConformanceResult_VIOLATION()
	}
	var type_ pkg_rhino.JSType = srcObj.GetJSType()
	_ = type_
	if (type_ == nil || type_.IsJavaNil_()) || ConformanceRules_ConformanceUtil_IsLooseType(type_) {
		if this_.GetConformanceRules_AbstractRuleBase().reportLooseTypeViolations {
			return StaticGet_ConformanceRules_ConformanceResult_POSSIBLE_VIOLATION_DUE_TO_LOOSE_TYPES()
		} else {
			return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
		}
	}
	if ((this_.domHelperType != nil && !this_.domHelperType.IsJavaNil_()) && this_.domHelperType.IsSubtypeOf_1(type_)) || ((this_.documentType != nil && !this_.documentType.IsJavaNil_()) && this_.documentType.IsSubtypeOf_1(type_)) {
		return StaticGet_ConformanceRules_ConformanceResult_VIOLATION()
	}
	return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
}

// Java class: com.google.javascript.jscomp.ConformanceRules.BanCreateDom
type ConformanceRules_BanCreateDom struct {
	ConformanceRules_AbstractRuleBase
	bannedTagAttrs [][]string       `gson:"bannedTagAttrs"`
	domHelperType  pkg_rhino.JSType `gson:"domHelperType"`
	classNameTypes pkg_rhino.JSType `gson:"classNameTypes"`
}

func (this_ *ConformanceRules_BanCreateDom) BindSelf_(s any) {
	_ = s
	this_.ConformanceRules_AbstractRuleBase.BindSelf_(s)
}
func (this_ *ConformanceRules_BanCreateDom) IsJavaNil_() bool { return this_ == nil }
func (*ConformanceRules_BanCreateDom) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConformanceRules$BanCreateDom", "BanCreateDom"
}
func (this_ *ConformanceRules_BanCreateDom) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ConformanceRules_AbstractRuleBase.IdentityHash_()
}
func (*ConformanceRules_BanCreateDom) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConformanceRules.BanCreateDom", "com.google.javascript.jscomp.ConformanceRules$BanCreateDom", func(v_ any) bool {
		return j2gort.InstanceOf[*ConformanceRules_BanCreateDom](v_)
	})
}
func (*ConformanceRules_BanCreateDom) Implements_jscomp_ConformanceRules_AbstractRule() {}
func (*ConformanceRules_BanCreateDom) Implements_jscomp_CheckConformance_Rule()         {}

var static_ConformanceRules_BanCreateDom_ELEMENT_TAG_NAMES j2gort.Ext_com_google_common_collect_ImmutableMultimap

func StaticGet_ConformanceRules_BanCreateDom_ELEMENT_TAG_NAMES() j2gort.Ext_com_google_common_collect_ImmutableMultimap {
	clinit_ConformanceRules_BanCreateDom()
	return static_ConformanceRules_BanCreateDom_ELEMENT_TAG_NAMES
}
func StaticSet_ConformanceRules_BanCreateDom_ELEMENT_TAG_NAMES(v j2gort.Ext_com_google_common_collect_ImmutableMultimap) {
	clinit_ConformanceRules_BanCreateDom()
	static_ConformanceRules_BanCreateDom_ELEMENT_TAG_NAMES = v
}

var clinitDone_ConformanceRules_BanCreateDom uint32
var clinitState_ConformanceRules_BanCreateDom j2gort.ClinitState

func clinit_ConformanceRules_BanCreateDom() {
	if atomic.LoadUint32(&clinitDone_ConformanceRules_BanCreateDom) != 0 {
		return
	}
	clinitSlow_ConformanceRules_BanCreateDom()
}

func clinitSlow_ConformanceRules_BanCreateDom() {
	if !j2gort.EnterClinit(&clinitState_ConformanceRules_BanCreateDom, &clinitDone_ConformanceRules_BanCreateDom) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ConformanceRules_BanCreateDom, &clinitDone_ConformanceRules_BanCreateDom)
	static_ConformanceRules_BanCreateDom_ELEMENT_TAG_NAMES = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.MultimapBuilderPutStr(j2gort.Ext_com_google_common_collect_ImmutableMultimap_S_Builder(), "HTMLAnchorElement", "a")), "HTMLAppletElement", "applet")), "HTMLAreaElement", "area")), "HTMLAudioElement", "audio")), "HTMLBRElement", "br")), "HTMLBaseElement", "base")), "HTMLBaseFontElement", "basefont")), "HTMLBodyElement", "body")), "HTMLButtonElement", "button")), "HTMLCanvasElement", "canvas")), "HTMLDListElement", "dl")), "HTMLDataListElement", "datalist")), "HTMLDetailsElement", "details")), "HTMLDialogElement", "dialog")), "HTMLDirectoryElement", "dir")), "HTMLDivElement", "div")), "HTMLEmbedElement", "embed")), "HTMLFieldSetElement", "fieldset")), "HTMLFontElement", "font")), "HTMLFormElement", "form")), "HTMLFrameElement", "frame")), "HTMLFrameSetElement", "frameset")), "HTMLHRElement", "hr")), "HTMLHeadElement", "head")), "HTMLHeadingElement", "h1")), "HTMLHeadingElement", "h2")), "HTMLHeadingElement", "h3")), "HTMLHeadingElement", "h4")), "HTMLHeadingElement", "h5")), "HTMLHeadingElement", "h6")), "HTMLHtmlElement", "html")), "HTMLIFrameElement", "iframe")), "HTMLImageElement", "img")), "HTMLInputElement", "input")), "HTMLIsIndexElement", "isindex")), "HTMLLIElement", "li")), "HTMLLabelElement", "label")), "HTMLLegendElement", "legend")), "HTMLLinkElement", "link")), "HTMLMapElement", "map")), "HTMLMenuElement", "menu")), "HTMLMetaElement", "meta")), "HTMLMeterElement", "meter")), "HTMLModElement", "del")), "HTMLModElement", "ins")), "HTMLOListElement", "ol")), "HTMLObjectElement", "object")), "HTMLOptGroupElement", "optgroup")), "HTMLOptionElement", "option")), "HTMLOutputElement", "output")), "HTMLParagraphElement", "p")), "HTMLParamElement", "param")), "HTMLPreElement", "pre")), "HTMLProgressElement", "progress")), "HTMLQuoteElement", "blockquote")), "HTMLQuoteElement", "q")), "HTMLScriptElement", "script")), "HTMLSelectElement", "select")), "HTMLSourceElement", "source")), "HTMLSpanElement", "span")), "HTMLStyleElement", "style")), "HTMLTableCaptionElement", "caption")), "HTMLTableCellElement", "td")), "HTMLTableCellElement", "th")), "HTMLTableColElement", "col")), "HTMLTableColElement", "colgroup")), "HTMLTableElement", "table")), "HTMLTableRowElement", "tr")), "HTMLTableSectionElement", "tbody")), "HTMLTableSectionElement", "tfoot")), "HTMLTableSectionElement", "thead")), "HTMLTemplateElement", "template")), "HTMLTextAreaElement", "textarea")), "HTMLTitleElement", "title")), "HTMLTrackElement", "track")), "HTMLUListElement", "ul")), "HTMLVideoElement", "video")).Build())
}

func (this_ *ConformanceRules_BanCreateDom) initFields_ConformanceRules_BanCreateDom() {
	_ = this_
}

// Java constructor: BanCreateDom(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.Requirement requirement)
func NewConformanceRules_BanCreateDom(compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) *ConformanceRules_BanCreateDom {
	clinit_ConformanceRules_BanCreateDom()
	this_ := &ConformanceRules_BanCreateDom{}
	this_.BindSelf_(this_)
	InitConformanceRules_BanCreateDom(this_, compiler, requirement)
	return this_
}

func InitConformanceRules_BanCreateDom(this_ *ConformanceRules_BanCreateDom, compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) {
	_ = compiler
	_ = requirement
	InitConformanceRules_AbstractRule(&this_.ConformanceRules_AbstractRuleBase, compiler, requirement)
	this_.initFields_ConformanceRules_BanCreateDom()
	var bannedTagAttrs [][]string = [][]string{}
	_ = bannedTagAttrs
	{
		cur_44 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](requirement.GetValueList())))
		for cur_44.HasNext() {
			value := cur_44.NextStr()
			_ = value
			var tagAttr j2gort.Ext_java_util_List = j2gort.SplitterSplitToListStr(StaticGet_ConformanceRules_ON_DOT(), value)
			_ = tagAttr
			if ((tagAttr.Size() != int32(2)) || (j2gort.CheckNotNullString(j2gort.ListGetStr(tagAttr, int32(0))) == "")) || (j2gort.CheckNotNullString(j2gort.ListGetStr(tagAttr, int32(1))) == "") {
				panic(NewCheckConformance_InvalidRequirementSpec_1("Values must be in the format tagname.attribute."))
			}
			var lowercasedTag string = javalang.StringToLowerCaseASCII(j2gort.ListGetStr(tagAttr, int32(0)))
			_ = lowercasedTag
			bannedTagAttrs = append(bannedTagAttrs, j2gort.CheckNotNullElem([]string{lowercasedTag, j2gort.ListGetStr(tagAttr, int32(1))}))
		}
	}
	this_.bannedTagAttrs = bannedTagAttrs
	if len(this_.bannedTagAttrs) == 0 {
		panic(NewCheckConformance_InvalidRequirementSpec_1("Specify one or more values."))
	}
	this_.domHelperType = compiler.GetTypeRegistry().GetGlobalType("goog.dom.DomHelper")
	this_.classNameTypes = compiler.GetTypeRegistry().CreateUnionType_1_JSType([]pkg_rhino.JSType{compiler.GetTypeRegistry().GetNativeType(pkg_rhino.StaticGet_JSTypeNative_STRING_TYPE()), compiler.GetTypeRegistry().GetNativeType(pkg_rhino.StaticGet_JSTypeNative_ARRAY_TYPE()), compiler.GetTypeRegistry().GetNativeType(pkg_rhino.StaticGet_JSTypeNative_NULL_TYPE()), compiler.GetTypeRegistry().GetNativeType(pkg_rhino.StaticGet_JSTypeNative_VOID_TYPE())})
}

// Java method: BanCreateDom.checkConformance(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *ConformanceRules_BanCreateDom) CheckConformance_2(t *NodeTraversal, n pkg_rhino.Node) *ConformanceRules_ConformanceResult {
	_ = t
	_ = n
	if !(this_.priv_BanCreateDom_isCreateDomCall_2165(n)) {
		return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
	}
	if n.GetChildCount() < int32(3) {
		return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
	}
	var tagNames j2gort.Ext_com_google_common_collect_ImmutableCollection = this_.priv_BanCreateDom_getTagNames_2164(n.GetSecondChild())
	_ = tagNames
	var attrs pkg_rhino.Node = n.GetChildAtIndex(int32(2))
	_ = attrs
	var attrsType pkg_rhino.JSType = attrs.GetJSType()
	_ = attrsType
	if attrsType == nil || attrsType.IsJavaNil_() {
		return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
	}
	var isClassName bool = (!(attrsType.IsUnknownType()) && attrsType.IsSubtypeOf_1(this_.classNameTypes))
	_ = isClassName
	if attrs.IsNull() || ((attrsType != nil && !attrsType.IsJavaNil_()) && attrsType.IsVoidType()) {
		return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
	}
	{

		for _, tagAttr := range this_.bannedTagAttrs {

			_ = tagAttr
			if ((tagNames != nil) && !(j2gort.SetContainsStr(tagNames, tagAttr[int32(0)]))) && !(j2gort.CheckNotNullString(tagAttr[int32(0)]) == "*") {
				continue
			}
			var violation *ConformanceRules_ConformanceResult = nil
			_ = violation
			if (tagNames != nil) || (j2gort.CheckNotNullString(tagAttr[int32(0)]) == "*") {
				violation = StaticGet_ConformanceRules_ConformanceResult_VIOLATION()
			} else if this_.GetConformanceRules_AbstractRuleBase().reportLooseTypeViolations {
				violation = StaticGet_ConformanceRules_ConformanceResult_POSSIBLE_VIOLATION()
			} else {
				violation = StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
			}
			if isClassName {
				if !(j2gort.CheckNotNullString(tagAttr[int32(1)]) == "class") {
					continue
				}
				return violation
			}
			if ((j2gort.CheckNotNullString(tagAttr[int32(1)]) == "textContent") && (n.GetChildCount() > int32(3))) && (violation != StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()) {
				return violation
			}
			if !(attrs.IsObjectLit()) {
				if this_.GetConformanceRules_AbstractRuleBase().reportLooseTypeViolations {
					return StaticGet_ConformanceRules_ConformanceResult_POSSIBLE_VIOLATION()
				} else {
					return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
				}
			}
			var prop pkg_rhino.Node = NodeUtil_GetFirstPropMatchingKey(attrs, tagAttr[int32(1)])
			_ = prop
			if prop != 0 {
				if NodeUtil_IsSomeCompileTimeConstStringValue(prop) {
					continue
				}
				return violation
			}
			{
				var attr pkg_rhino.Node = attrs.GetFirstChild()
				_ = attr
				for ; attr != 0; attr = attr.GetNext() {
					if attr.IsComputedProp() {
						if this_.GetConformanceRules_AbstractRuleBase().reportLooseTypeViolations {
							return StaticGet_ConformanceRules_ConformanceResult_POSSIBLE_VIOLATION_DUE_TO_LOOSE_TYPES()
						} else {
							return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
						}
					}
				}
			}
		}
	}
	return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
}

// Java method: BanCreateDom.getTagNames(com.google.javascript.rhino.Node tag)
func (this_ *ConformanceRules_BanCreateDom) priv_BanCreateDom_getTagNames_2164(tag pkg_rhino.Node) j2gort.Ext_com_google_common_collect_ImmutableCollection {
	_ = tag
	if tag.IsStringLit() {
		return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableCollection](j2gort.ImmSetOfStr1(javalang.StringToLowerCaseASCII(tag.GetString())))
	} else if tag.IsGetProp() && StaticGet_ConformanceRules_GOOG_DOM_TAGNAME().Matches(tag.GetFirstChild()) {
		return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableCollection](j2gort.ImmSetOfStr1(javalang.StringToLowerCaseASCII(tag.GetString())))
	}
	var type_ pkg_rhino.JSType = tag.GetJSType()
	_ = type_
	if (type_ == nil || type_.IsJavaNil_()) || !(type_.IsTemplatizedType()) {
		return nil
	}
	var typeAsObj pkg_rhino.ObjectType = type_.ToMaybeObjectType()
	_ = typeAsObj
	if j2gort.CheckNotNullString(typeAsObj.GetRawType().GetDisplayName()) == "goog.dom.TagName" {
		var tagType pkg_rhino.JSType = j2gort.AsIface[pkg_rhino.JSType](j2gort.Ext_com_google_common_collect_Iterables_GetOnlyElement_1(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](typeAsObj.GetTemplateTypes())))
		_ = tagType
		return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableCollection](j2gort.MultimapGetStr(StaticGet_ConformanceRules_BanCreateDom_ELEMENT_TAG_NAMES(), tagType.GetDisplayName()))
	}
	return nil
}

// Java method: BanCreateDom.isCreateDomCall(com.google.javascript.rhino.Node n)
func (this_ *ConformanceRules_BanCreateDom) priv_BanCreateDom_isCreateDomCall_2165(n pkg_rhino.Node) bool {
	_ = n
	if !(n.IsCall()) {
		return false
	}
	var target pkg_rhino.Node = n.GetFirstChild()
	_ = target
	if !(target.IsGetProp()) {
		return false
	}
	{
		var iifeCond_1 bool
		{

			sr_47 := "createDom"
			se_48 := target.GetString()
			iifeCond_1 = j2gort.CheckNotNullString(sr_47) == se_48
		}
		if !(iifeCond_1) {
			return false
		}
	}
	var srcObj pkg_rhino.Node = target.GetFirstChild()
	_ = srcObj
	if StaticGet_ConformanceRules_GOOG_DOM().Matches(srcObj) {
		return true
	}
	var type_ pkg_rhino.JSType = srcObj.GetJSType()
	_ = type_
	if type_ == nil || type_.IsJavaNil_() {
		return false
	}
	if type_.Equals(this_.domHelperType) {
		return true
	}
	return false
}

// Java class: com.google.javascript.jscomp.ConformanceRules.SecuritySensitiveAttributes
type ConformanceRules_SecuritySensitiveAttributes struct {
	bannedAtrrs             j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"bannedAtrrs"`
	globalNamespaceSupplier j2gort.Ext_java_util_function_Supplier            `gson:"globalNamespaceSupplier"`
	idHash_                 atomic.Int32
}

func (this_ *ConformanceRules_SecuritySensitiveAttributes) BindSelf_(s any) {
	_ = s
}
func (this_ *ConformanceRules_SecuritySensitiveAttributes) IsJavaNil_() bool { return this_ == nil }
func (*ConformanceRules_SecuritySensitiveAttributes) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConformanceRules$SecuritySensitiveAttributes", "SecuritySensitiveAttributes"
}
func (this_ *ConformanceRules_SecuritySensitiveAttributes) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ConformanceRules_SecuritySensitiveAttributes) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConformanceRules.SecuritySensitiveAttributes", "com.google.javascript.jscomp.ConformanceRules$SecuritySensitiveAttributes", func(v_ any) bool {
		return j2gort.InstanceOf[*ConformanceRules_SecuritySensitiveAttributes](v_)
	})
}

var static_ConformanceRules_SecuritySensitiveAttributes_ALL_BANNED_ATTRS j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_ConformanceRules_SecuritySensitiveAttributes_ALL_BANNED_ATTRS() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_ConformanceRules_SecuritySensitiveAttributes()
	return static_ConformanceRules_SecuritySensitiveAttributes_ALL_BANNED_ATTRS
}
func StaticSet_ConformanceRules_SecuritySensitiveAttributes_ALL_BANNED_ATTRS(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_ConformanceRules_SecuritySensitiveAttributes()
	static_ConformanceRules_SecuritySensitiveAttributes_ALL_BANNED_ATTRS = v
}

var clinitDone_ConformanceRules_SecuritySensitiveAttributes uint32
var clinitState_ConformanceRules_SecuritySensitiveAttributes j2gort.ClinitState

func clinit_ConformanceRules_SecuritySensitiveAttributes() {
	if atomic.LoadUint32(&clinitDone_ConformanceRules_SecuritySensitiveAttributes) != 0 {
		return
	}
	clinitSlow_ConformanceRules_SecuritySensitiveAttributes()
}

func clinitSlow_ConformanceRules_SecuritySensitiveAttributes() {
	if !j2gort.EnterClinit(&clinitState_ConformanceRules_SecuritySensitiveAttributes, &clinitDone_ConformanceRules_SecuritySensitiveAttributes) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ConformanceRules_SecuritySensitiveAttributes, &clinitDone_ConformanceRules_SecuritySensitiveAttributes)
	static_ConformanceRules_SecuritySensitiveAttributes_ALL_BANNED_ATTRS = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_7("href", "rel", "src", "srcdoc", "action", "formaction", []any{"sandbox", "icon", "codebase", "data"})
}

func (this_ *ConformanceRules_SecuritySensitiveAttributes) initFields_ConformanceRules_SecuritySensitiveAttributes() {
	_ = this_
}

// Java constructor: SecuritySensitiveAttributes()
func NewConformanceRules_SecuritySensitiveAttributes() *ConformanceRules_SecuritySensitiveAttributes {
	this_ := &ConformanceRules_SecuritySensitiveAttributes{}
	InitConformanceRules_SecuritySensitiveAttributes(this_)
	return this_
}

func InitConformanceRules_SecuritySensitiveAttributes(this_ *ConformanceRules_SecuritySensitiveAttributes) {
	InitConformanceRules_SecuritySensitiveAttributes_2(this_, j2gort.AsIface[j2gort.Ext_java_util_Collection](StaticGet_ConformanceRules_SecuritySensitiveAttributes_ALL_BANNED_ATTRS()), lambdaSingleton_conformance_rules_33)
}

// Java constructor: SecuritySensitiveAttributes(java.util.function.Supplier<com.google.javascript.jscomp.GlobalNamespace> globalNamespaceSupplier)
func NewConformanceRules_SecuritySensitiveAttributes_1_Ext_java_util_function_Supplier(globalNamespaceSupplier j2gort.Ext_java_util_function_Supplier) *ConformanceRules_SecuritySensitiveAttributes {
	this_ := &ConformanceRules_SecuritySensitiveAttributes{}
	InitConformanceRules_SecuritySensitiveAttributes_1_Ext_java_util_function_Supplier(this_, globalNamespaceSupplier)
	return this_
}

func InitConformanceRules_SecuritySensitiveAttributes_1_Ext_java_util_function_Supplier(this_ *ConformanceRules_SecuritySensitiveAttributes, globalNamespaceSupplier j2gort.Ext_java_util_function_Supplier) {
	_ = globalNamespaceSupplier
	InitConformanceRules_SecuritySensitiveAttributes_2(this_, j2gort.AsIface[j2gort.Ext_java_util_Collection](StaticGet_ConformanceRules_SecuritySensitiveAttributes_ALL_BANNED_ATTRS()), globalNamespaceSupplier)
}

// Java constructor: SecuritySensitiveAttributes(java.util.Collection<java.lang.String> bannedAtrrs)
func NewConformanceRules_SecuritySensitiveAttributes_1_Ext_java_util_Collection(bannedAtrrs j2gort.Ext_java_util_Collection) *ConformanceRules_SecuritySensitiveAttributes {
	this_ := &ConformanceRules_SecuritySensitiveAttributes{}
	InitConformanceRules_SecuritySensitiveAttributes_1_Ext_java_util_Collection(this_, bannedAtrrs)
	return this_
}

func InitConformanceRules_SecuritySensitiveAttributes_1_Ext_java_util_Collection(this_ *ConformanceRules_SecuritySensitiveAttributes, bannedAtrrs j2gort.Ext_java_util_Collection) {
	_ = bannedAtrrs
	InitConformanceRules_SecuritySensitiveAttributes_2(this_, bannedAtrrs, lambdaSingleton_conformance_rules_34)
}

// Java constructor: SecuritySensitiveAttributes(java.util.Collection<java.lang.String> bannedAtrrs, java.util.function.Supplier<com.google.javascript.jscomp.GlobalNamespace> globalNamespaceSupplier)
func NewConformanceRules_SecuritySensitiveAttributes_2(bannedAtrrs j2gort.Ext_java_util_Collection, globalNamespaceSupplier j2gort.Ext_java_util_function_Supplier) *ConformanceRules_SecuritySensitiveAttributes {
	this_ := &ConformanceRules_SecuritySensitiveAttributes{}
	InitConformanceRules_SecuritySensitiveAttributes_2(this_, bannedAtrrs, globalNamespaceSupplier)
	return this_
}

func InitConformanceRules_SecuritySensitiveAttributes_2(this_ *ConformanceRules_SecuritySensitiveAttributes, bannedAtrrs j2gort.Ext_java_util_Collection, globalNamespaceSupplier j2gort.Ext_java_util_function_Supplier) {
	_ = bannedAtrrs
	_ = globalNamespaceSupplier
	this_.initFields_ConformanceRules_SecuritySensitiveAttributes()
	this_.bannedAtrrs = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](bannedAtrrs.Stream().Map(lambdaSingleton_conformance_rules_35).Collect_1(j2gort.Ext_com_google_common_collect_ImmutableSet_ToImmutableSet()))
	this_.globalNamespaceSupplier = globalNamespaceSupplier
}

// Java method: SecuritySensitiveAttributes.contains(java.lang.String attributeName)
func (this_ *ConformanceRules_SecuritySensitiveAttributes) Contains(attributeName string) bool {
	_ = attributeName
	return j2gort.SetContainsStr(this_.bannedAtrrs, attributeName)
}

// Java method: SecuritySensitiveAttributes.checkConformanceForAttributeName(com.google.javascript.jscomp.NodeTraversal traversal, com.google.javascript.rhino.Node attrName)
func (this_ *ConformanceRules_SecuritySensitiveAttributes) CheckConformanceForAttributeName(traversal *NodeTraversal, attrName pkg_rhino.Node) *ConformanceRules_ConformanceResult {
	_ = traversal
	_ = attrName
	var literalName string = ConformanceRules_ConformanceUtil_InferStringValue(traversal.GetScope(), attrName, this_.globalNamespaceSupplier)
	_ = literalName
	if literalName == j2gort.NullString {
		if ConformanceRules_ConformanceUtil_IsXid(attrName.GetJSType()) {
			return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
		} else {
			return StaticGet_ConformanceRules_ConformanceResult_VIOLATION()
		}
	}
	literalName = javalang.StringToLowerCaseASCII(literalName)
	if j2gort.SetContainsStr(this_.bannedAtrrs, literalName) || ConformanceRules_ConformanceUtil_IsEventHandlerAttrName(literalName) {
		return StaticGet_ConformanceRules_ConformanceResult_VIOLATION()
	}
	return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
}

// Java method: SecuritySensitiveAttributes.checkConformanceForAttributeNameWithHighConfidence(com.google.javascript.jscomp.NodeTraversal traversal, com.google.javascript.rhino.Node attrName)
func (this_ *ConformanceRules_SecuritySensitiveAttributes) CheckConformanceForAttributeNameWithHighConfidence(traversal *NodeTraversal, attrName pkg_rhino.Node) *ConformanceRules_ConformanceResult {
	_ = traversal
	_ = attrName
	var literalName string = ConformanceRules_ConformanceUtil_InferStringValue(traversal.GetScope(), attrName, this_.globalNamespaceSupplier)
	_ = literalName
	if literalName == j2gort.NullString {
		return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
	}
	if this_.Contains(javalang.StringToLowerCaseASCII(literalName)) {
		return StaticGet_ConformanceRules_ConformanceResult_VIOLATION()
	} else {
		return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
	}
}

// Java class: com.google.javascript.jscomp.ConformanceRules.BanElementSetAttribute
type ConformanceRules_BanElementSetAttribute struct {
	ConformanceRules_AbstractRuleBase
	globalNamespace                  *GlobalNamespace                              `gson:"globalNamespace"`
	elementType                      pkg_rhino.JSType                              `gson:"elementType"`
	securitySensitiveAttributes      *ConformanceRules_SecuritySensitiveAttributes `gson:"securitySensitiveAttributes"`
	defaultDecisionForUncertainCases *ConformanceRules_ConformanceResult           `gson:"defaultDecisionForUncertainCases"`
	performGlobalNamespaceAnalysis   bool                                          `gson:"performGlobalNamespaceAnalysis"`
}

func (this_ *ConformanceRules_BanElementSetAttribute) BindSelf_(s any) {
	_ = s
	this_.ConformanceRules_AbstractRuleBase.BindSelf_(s)
}
func (this_ *ConformanceRules_BanElementSetAttribute) IsJavaNil_() bool { return this_ == nil }
func (*ConformanceRules_BanElementSetAttribute) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConformanceRules$BanElementSetAttribute", "BanElementSetAttribute"
}
func (this_ *ConformanceRules_BanElementSetAttribute) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ConformanceRules_AbstractRuleBase.IdentityHash_()
}
func (*ConformanceRules_BanElementSetAttribute) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConformanceRules.BanElementSetAttribute", "com.google.javascript.jscomp.ConformanceRules$BanElementSetAttribute", func(v_ any) bool {
		return j2gort.InstanceOf[*ConformanceRules_BanElementSetAttribute](v_)
	})
}
func (*ConformanceRules_BanElementSetAttribute) Implements_jscomp_ConformanceRules_AbstractRule() {}
func (*ConformanceRules_BanElementSetAttribute) Implements_jscomp_CheckConformance_Rule()         {}

func StaticGet_ConformanceRules_BanElementSetAttribute_ELEMENT_TYPE_NAME() string { return "Element" }
func StaticGet_ConformanceRules_BanElementSetAttribute_SET_ATTRIBUTE() string     { return "setAttribute" }
func StaticGet_ConformanceRules_BanElementSetAttribute_SET_ATTRIBUTE_NS() string {
	return "setAttributeNS"
}
func StaticGet_ConformanceRules_BanElementSetAttribute_SET_ATTRIBUTE_NODE() string {
	return "setAttributeNode"
}
func StaticGet_ConformanceRules_BanElementSetAttribute_SET_ATTRIBUTE_NODE_NS() string {
	return "setAttributeNodeNS"
}

var static_ConformanceRules_BanElementSetAttribute_BANNED_PROPERTIES j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_ConformanceRules_BanElementSetAttribute_BANNED_PROPERTIES() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_ConformanceRules_BanElementSetAttribute()
	return static_ConformanceRules_BanElementSetAttribute_BANNED_PROPERTIES
}
func StaticSet_ConformanceRules_BanElementSetAttribute_BANNED_PROPERTIES(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_ConformanceRules_BanElementSetAttribute()
	static_ConformanceRules_BanElementSetAttribute_BANNED_PROPERTIES = v
}
func StaticGet_ConformanceRules_BanElementSetAttribute_GLOBAL_NAMESPACE_ANALYSIS_LIMIT() int32 {
	return int32(10000)
}

var clinitDone_ConformanceRules_BanElementSetAttribute uint32
var clinitState_ConformanceRules_BanElementSetAttribute j2gort.ClinitState

func clinit_ConformanceRules_BanElementSetAttribute() {
	if atomic.LoadUint32(&clinitDone_ConformanceRules_BanElementSetAttribute) != 0 {
		return
	}
	clinitSlow_ConformanceRules_BanElementSetAttribute()
}

func clinitSlow_ConformanceRules_BanElementSetAttribute() {
	if !j2gort.EnterClinit(&clinitState_ConformanceRules_BanElementSetAttribute, &clinitDone_ConformanceRules_BanElementSetAttribute) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ConformanceRules_BanElementSetAttribute, &clinitDone_ConformanceRules_BanElementSetAttribute)
	static_ConformanceRules_BanElementSetAttribute_BANNED_PROPERTIES = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_4("setAttribute", "setAttributeNS", "setAttributeNode", "setAttributeNodeNS")
}

func (this_ *ConformanceRules_BanElementSetAttribute) initFields_ConformanceRules_BanElementSetAttribute() {
	_ = this_
	this_.performGlobalNamespaceAnalysis = true
}

// Java constructor: BanElementSetAttribute(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.Requirement requirement)
func NewConformanceRules_BanElementSetAttribute(compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) *ConformanceRules_BanElementSetAttribute {
	this_ := &ConformanceRules_BanElementSetAttribute{}
	this_.BindSelf_(this_)
	InitConformanceRules_BanElementSetAttribute(this_, compiler, requirement)
	return this_
}

func InitConformanceRules_BanElementSetAttribute(this_ *ConformanceRules_BanElementSetAttribute, compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) {
	_ = compiler
	_ = requirement
	InitConformanceRules_AbstractRule(&this_.ConformanceRules_AbstractRuleBase, compiler, requirement)
	this_.initFields_ConformanceRules_BanElementSetAttribute()
	this_.securitySensitiveAttributes = NewConformanceRules_SecuritySensitiveAttributes_2(j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](requirement.GetValueList()), func() j2gort.Ext_java_util_function_Supplier {
		bound_50 := this_
		return j2gort.NewLambda_Ext_java_util_function_Supplier(func() any { return j2gort.PtrToAny(bound_50.GetGlobalNamespace()) })
	}())
	if requirement.GetReportLooseTypeViolations() {
		this_.defaultDecisionForUncertainCases = StaticGet_ConformanceRules_ConformanceResult_VIOLATION()
	} else {
		this_.defaultDecisionForUncertainCases = StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
	}
	this_.elementType = compiler.GetTypeRegistry().GetGlobalType("Element")
}

// Java method: BanElementSetAttribute.getGlobalNamespace()
func (this_ *ConformanceRules_BanElementSetAttribute) GetGlobalNamespace() *GlobalNamespace {
	if this_.performGlobalNamespaceAnalysis && (this_.globalNamespace == nil) {
		if NodeUtil_CountAstSizeUpToLimit(this_.GetConformanceRules_AbstractRuleBase().compiler.GetJsRoot(), int32(10000)) >= int32(10000) {
			this_.performGlobalNamespaceAnalysis = false
		} else {
			this_.globalNamespace = NewGlobalNamespace_2(this_.GetConformanceRules_AbstractRuleBase().compiler, this_.GetConformanceRules_AbstractRuleBase().compiler.GetJsRoot())
		}
	}
	return this_.globalNamespace
}

// Java method: BanElementSetAttribute.checkConformance(com.google.javascript.jscomp.NodeTraversal traversal, com.google.javascript.rhino.Node node)
func (this_ *ConformanceRules_BanElementSetAttribute) CheckConformance_2(traversal *NodeTraversal, node pkg_rhino.Node) *ConformanceRules_ConformanceResult {
	_ = traversal
	_ = node
	if node.IsCall() {
		return this_.priv_BanElementSetAttribute_checkConformanceOnPropertyCall_2170(traversal, node)
	} else if node.IsGetElem() && NodeUtil_IsLValue(node) {
		return this_.priv_BanElementSetAttribute_checkConformanceOnGetElement_2172(traversal, node)
	}
	return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
}

// Java method: BanElementSetAttribute.checkConformanceOnPropertyCall(com.google.javascript.jscomp.NodeTraversal traversal, com.google.javascript.rhino.Node callNode)
func (this_ *ConformanceRules_BanElementSetAttribute) priv_BanElementSetAttribute_checkConformanceOnPropertyCall_2170(traversal *NodeTraversal, callNode pkg_rhino.Node) *ConformanceRules_ConformanceResult {
	_ = traversal
	_ = callNode
	var calledProperty j2gort.Ext_com_google_common_base_Optional = this_.priv_BanElementSetAttribute_getBannedPropertyName_2171(callNode)
	_ = calledProperty
	if !(calledProperty.IsPresent()) {
		return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
	}
	{
		var iifeCond_1 bool
		{

			sr_51 := j2gort.CastString(calledProperty.Get())
			se_52 := "setAttribute"
			iifeCond_1 = j2gort.CheckNotNullString(sr_51) == se_52
		}
		if iifeCond_1 {
			if callNode.GetChildCount() < int32(3) {
				return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
			}
			if this_.GetConformanceRules_AbstractRuleBase().requirement.GetReportLooseTypeViolations() {
				return this_.securitySensitiveAttributes.CheckConformanceForAttributeName(traversal, callNode.GetSecondChild())
			} else {
				return this_.securitySensitiveAttributes.CheckConformanceForAttributeNameWithHighConfidence(traversal, callNode.GetSecondChild())
			}
		}
	}
	{
		var iifeCond_1 bool
		{

			sr_53 := j2gort.CastString(calledProperty.Get())
			se_54 := "setAttributeNS"
			iifeCond_1 = j2gort.CheckNotNullString(sr_53) == se_54
		}
		if iifeCond_1 {
			if callNode.GetChildCount() < int32(4) {
				return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
			}
			if callNode.GetSecondChild().IsNull() {
				if this_.GetConformanceRules_AbstractRuleBase().requirement.GetReportLooseTypeViolations() {
					return this_.securitySensitiveAttributes.CheckConformanceForAttributeName(traversal, callNode.GetChildAtIndex(int32(2)))
				} else {
					return this_.securitySensitiveAttributes.CheckConformanceForAttributeNameWithHighConfidence(traversal, callNode.GetChildAtIndex(int32(2)))
				}
			}
		}
	}
	return this_.defaultDecisionForUncertainCases
}

// Java method: BanElementSetAttribute.getBannedPropertyName(com.google.javascript.rhino.Node callNode)
func (this_ *ConformanceRules_BanElementSetAttribute) priv_BanElementSetAttribute_getBannedPropertyName_2171(callNode pkg_rhino.Node) j2gort.Ext_com_google_common_base_Optional {
	_ = callNode
	var target pkg_rhino.Node = callNode.GetFirstChild()
	_ = target
	if !(target.IsGetProp()) {
		return j2gort.Ext_com_google_common_base_Optional_Absent()
	}
	var propertyName string = target.GetString()
	_ = propertyName
	if !(j2gort.SetContainsStr(StaticGet_ConformanceRules_BanElementSetAttribute_BANNED_PROPERTIES(), propertyName)) {
		return j2gort.Ext_com_google_common_base_Optional_Absent()
	}
	var type_ pkg_rhino.JSType = target.GetFirstChild().GetJSType()
	_ = type_
	if type_ == nil || type_.IsJavaNil_() {
		return j2gort.Ext_com_google_common_base_Optional_Absent()
	}
	if (this_.elementType != nil && !this_.elementType.IsJavaNil_()) && type_.RestrictByNotNullOrUndefined().IsSubtypeOf_1(this_.elementType) {
		return j2gort.Ext_com_google_common_base_Optional_Of(j2gort.BoxString(propertyName))
	}
	return j2gort.Ext_com_google_common_base_Optional_Absent()
}

// Java method: BanElementSetAttribute.checkConformanceOnGetElement(com.google.javascript.jscomp.NodeTraversal traversal, com.google.javascript.rhino.Node getElementNode)
func (this_ *ConformanceRules_BanElementSetAttribute) priv_BanElementSetAttribute_checkConformanceOnGetElement_2172(traversal *NodeTraversal, getElementNode pkg_rhino.Node) *ConformanceRules_ConformanceResult {
	_ = traversal
	_ = getElementNode
	var key pkg_rhino.Node = getElementNode.GetSecondChild()
	_ = key
	var keyName string = ConformanceRules_ConformanceUtil_InferStringValue(traversal.GetScope(), key, func() j2gort.Ext_java_util_function_Supplier {
		bound_55 := this_
		return j2gort.NewLambda_Ext_java_util_function_Supplier(func() any { return j2gort.PtrToAny(bound_55.GetGlobalNamespace()) })
	}())
	_ = keyName
	if keyName != j2gort.NullString {
		if this_.securitySensitiveAttributes.Contains(javalang.StringToLowerCaseASCII(keyName)) || (this_.GetConformanceRules_AbstractRuleBase().requirement.GetReportLooseTypeViolations() && ((j2gort.CheckNotNullString(keyName) == "innerHTML") || (j2gort.CheckNotNullString(keyName) == "outerHTML"))) {
			if this_.priv_BanElementSetAttribute_hasElementType_2173(getElementNode) {
				return StaticGet_ConformanceRules_ConformanceResult_VIOLATION()
			}
		}
		return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
	}
	if !(this_.GetConformanceRules_AbstractRuleBase().requirement.GetReportLooseTypeViolations()) {
		return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
	}
	var keyType pkg_rhino.JSType = key.GetJSType()
	_ = keyType
	if (keyType == nil || keyType.IsJavaNil_()) || ConformanceRules_ConformanceUtil_IsXid(keyType) {
		return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
	}
	if this_.priv_BanElementSetAttribute_hasElementType_2173(getElementNode) {
		if keyType.IsString() {
			return StaticGet_ConformanceRules_ConformanceResult_VIOLATION()
		} else if keyType.IsUnionType() {
			var union *pkg_rhino.UnionType = keyType.ToMaybeUnionType()
			_ = union
			{
				var i int32 = int32(0)
				_ = i
				var n int32 = union.GetAlternatesCount()
				_ = n
				for ; i < n; i = i + 1 {
					if union.GetAlternateAt(i).IsString() {
						return StaticGet_ConformanceRules_ConformanceResult_VIOLATION()
					}
				}
			}
		}
	}
	return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
}

// Java method: BanElementSetAttribute.hasElementType(com.google.javascript.rhino.Node getElementNode)
func (this_ *ConformanceRules_BanElementSetAttribute) priv_BanElementSetAttribute_hasElementType_2173(getElementNode pkg_rhino.Node) bool {
	_ = getElementNode
	var objType pkg_rhino.JSType = getElementNode.GetFirstChild().GetJSType()
	_ = objType
	if (objType == nil || objType.IsJavaNil_()) || (this_.elementType == nil || this_.elementType.IsJavaNil_()) {
		return false
	}
	objType = objType.RestrictByNotNullOrUndefined()
	if (objType.IsUnknownType() || objType.IsEmptyType()) || !(objType.IsSubtypeOf_1(this_.elementType)) {
		return false
	}
	return true
}

// Java class: com.google.javascript.jscomp.ConformanceRules.BanSettingAttributes
type ConformanceRules_BanSettingAttributes struct {
	ConformanceRules_AbstractRuleBase
	restrictions []*ConformanceRules_BanSettingAttributes_AttributeSettingRestriction `gson:"restrictions"`
}

func (this_ *ConformanceRules_BanSettingAttributes) BindSelf_(s any) {
	_ = s
	this_.ConformanceRules_AbstractRuleBase.BindSelf_(s)
}
func (this_ *ConformanceRules_BanSettingAttributes) IsJavaNil_() bool { return this_ == nil }
func (*ConformanceRules_BanSettingAttributes) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConformanceRules$BanSettingAttributes", "BanSettingAttributes"
}
func (this_ *ConformanceRules_BanSettingAttributes) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ConformanceRules_AbstractRuleBase.IdentityHash_()
}
func (*ConformanceRules_BanSettingAttributes) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConformanceRules.BanSettingAttributes", "com.google.javascript.jscomp.ConformanceRules$BanSettingAttributes", func(v_ any) bool {
		return j2gort.InstanceOf[*ConformanceRules_BanSettingAttributes](v_)
	})
}
func (*ConformanceRules_BanSettingAttributes) Implements_jscomp_ConformanceRules_AbstractRule() {}
func (*ConformanceRules_BanSettingAttributes) Implements_jscomp_CheckConformance_Rule()         {}

var static_ConformanceRules_BanSettingAttributes_BANNED_ATTRS *ConformanceRules_SecuritySensitiveAttributes

func StaticGet_ConformanceRules_BanSettingAttributes_BANNED_ATTRS() *ConformanceRules_SecuritySensitiveAttributes {
	clinit_ConformanceRules_BanSettingAttributes()
	return static_ConformanceRules_BanSettingAttributes_BANNED_ATTRS
}
func StaticSet_ConformanceRules_BanSettingAttributes_BANNED_ATTRS(v *ConformanceRules_SecuritySensitiveAttributes) {
	clinit_ConformanceRules_BanSettingAttributes()
	static_ConformanceRules_BanSettingAttributes_BANNED_ATTRS = v
}

var clinitDone_ConformanceRules_BanSettingAttributes uint32
var clinitState_ConformanceRules_BanSettingAttributes j2gort.ClinitState

func clinit_ConformanceRules_BanSettingAttributes() {
	if atomic.LoadUint32(&clinitDone_ConformanceRules_BanSettingAttributes) != 0 {
		return
	}
	clinitSlow_ConformanceRules_BanSettingAttributes()
}

func clinitSlow_ConformanceRules_BanSettingAttributes() {
	if !j2gort.EnterClinit(&clinitState_ConformanceRules_BanSettingAttributes, &clinitDone_ConformanceRules_BanSettingAttributes) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ConformanceRules_BanSettingAttributes, &clinitDone_ConformanceRules_BanSettingAttributes)
	static_ConformanceRules_BanSettingAttributes_BANNED_ATTRS = NewConformanceRules_SecuritySensitiveAttributes()
}

func (this_ *ConformanceRules_BanSettingAttributes) initFields_ConformanceRules_BanSettingAttributes() {
	_ = this_
}

// Java constructor: BanSettingAttributes(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.Requirement requirement)
func NewConformanceRules_BanSettingAttributes(compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) *ConformanceRules_BanSettingAttributes {
	clinit_ConformanceRules_BanSettingAttributes()
	this_ := &ConformanceRules_BanSettingAttributes{}
	this_.BindSelf_(this_)
	InitConformanceRules_BanSettingAttributes(this_, compiler, requirement)
	return this_
}

func InitConformanceRules_BanSettingAttributes(this_ *ConformanceRules_BanSettingAttributes, compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) {
	_ = compiler
	_ = requirement
	InitConformanceRules_AbstractRule(&this_.ConformanceRules_AbstractRuleBase, compiler, requirement)
	this_.initFields_ConformanceRules_BanSettingAttributes()
	if j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](requirement.GetValueList()).IsEmpty() {
		panic(NewCheckConformance_InvalidRequirementSpec_1("missing value"))
	}
	var registry *pkg_rhino.JSTypeRegistry = compiler.GetTypeRegistry()
	_ = registry
	var builder []*ConformanceRules_BanSettingAttributes_AttributeSettingRestriction = []*ConformanceRules_BanSettingAttributes_AttributeSettingRestriction{}
	_ = builder
	{
		cur_56 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](requirement.GetValueList())))
		for cur_56.HasNext() {
			value := cur_56.NextStr()
			_ = value
			var type_ string = ConformanceRules_ConformanceUtil_GetClassFromDeclarationName(value)
			_ = type_
			var property string = ConformanceRules_ConformanceUtil_GetPropertyFromDeclarationName(value)
			_ = property
			if (type_ == j2gort.NullString) || (property == j2gort.NullString) {
				panic(NewCheckConformance_InvalidRequirementSpec_1("bad prop value"))
			}
			builder = append(builder, NewConformanceRules_BanSettingAttributes_AttributeSettingRestriction(registry.GetGlobalType(type_), property))
		}
	}
	this_.restrictions = builder
}

// Java method: BanSettingAttributes.checkConformance(com.google.javascript.jscomp.NodeTraversal traversal, com.google.javascript.rhino.Node node)
func (this_ *ConformanceRules_BanSettingAttributes) CheckConformance_2(traversal *NodeTraversal, node pkg_rhino.Node) *ConformanceRules_ConformanceResult {
	_ = traversal
	_ = node
	if node.IsCall() {
		return this_.priv_BanSettingAttributes_checkConformanceOnPropertyCall_2174(traversal, node)
	}
	return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
}

// Java method: BanSettingAttributes.checkConformanceOnPropertyCall(com.google.javascript.jscomp.NodeTraversal traversal, com.google.javascript.rhino.Node node)
func (this_ *ConformanceRules_BanSettingAttributes) priv_BanSettingAttributes_checkConformanceOnPropertyCall_2174(traversal *NodeTraversal, node pkg_rhino.Node) *ConformanceRules_ConformanceResult {
	_ = traversal
	_ = node
	var calledProperty j2gort.Ext_com_google_common_base_Optional = this_.priv_BanSettingAttributes_getBannedPropertyName_2175(node)
	_ = calledProperty
	if !(calledProperty.IsPresent()) {
		return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
	}
	if node.GetChildCount() < int32(3) {
		return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
	}
	return StaticGet_ConformanceRules_BanSettingAttributes_BANNED_ATTRS().CheckConformanceForAttributeName(traversal, node.GetSecondChild())
}

// Java method: BanSettingAttributes.getBannedPropertyName(com.google.javascript.rhino.Node node)
func (this_ *ConformanceRules_BanSettingAttributes) priv_BanSettingAttributes_getBannedPropertyName_2175(node pkg_rhino.Node) j2gort.Ext_com_google_common_base_Optional {
	_ = node
	var target pkg_rhino.Node = node.GetFirstChild()
	_ = target
	if !(target.IsGetProp()) {
		return j2gort.Ext_com_google_common_base_Optional_Absent()
	}
	var type_ pkg_rhino.JSType = target.GetFirstChild().GetJSType()
	_ = type_
	if type_ == nil || type_.IsJavaNil_() {
		return j2gort.Ext_com_google_common_base_Optional_Absent()
	}
	var propertyName string = target.GetString()
	_ = propertyName
	{

		for _, restricted := range this_.restrictions {

			_ = restricted
			{
				var iifeCond_1 bool
				{

					sr_59 := propertyName
					se_60 := restricted.property
					iifeCond_1 = j2gort.CheckNotNullString(sr_59) == se_60
				}
				if (iifeCond_1 && (restricted.type_ != nil && !restricted.type_.IsJavaNil_())) && type_.IsSubtypeOf_1(restricted.type_) {
					return j2gort.Ext_com_google_common_base_Optional_Of(j2gort.BoxString(propertyName))
				}
			}
		}
	}
	return j2gort.Ext_com_google_common_base_Optional_Absent()
}

// Java class: com.google.javascript.jscomp.ConformanceRules.BanSettingAttributes.AttributeSettingRestriction
type ConformanceRules_BanSettingAttributes_AttributeSettingRestriction struct {
	type_    pkg_rhino.JSType `gson:"type"`
	property string           `gson:"property"`
	idHash_  atomic.Int32
}

func (this_ *ConformanceRules_BanSettingAttributes_AttributeSettingRestriction) BindSelf_(s any) {
	_ = s
}
func (this_ *ConformanceRules_BanSettingAttributes_AttributeSettingRestriction) IsJavaNil_() bool {
	return this_ == nil
}
func (*ConformanceRules_BanSettingAttributes_AttributeSettingRestriction) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConformanceRules$BanSettingAttributes$AttributeSettingRestriction", "AttributeSettingRestriction"
}
func (this_ *ConformanceRules_BanSettingAttributes_AttributeSettingRestriction) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ConformanceRules_BanSettingAttributes_AttributeSettingRestriction) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConformanceRules.BanSettingAttributes.AttributeSettingRestriction", "com.google.javascript.jscomp.ConformanceRules$BanSettingAttributes$AttributeSettingRestriction", func(v_ any) bool {
		return j2gort.InstanceOf[*ConformanceRules_BanSettingAttributes_AttributeSettingRestriction](v_)
	})
}

func (this_ *ConformanceRules_BanSettingAttributes_AttributeSettingRestriction) initFields_ConformanceRules_BanSettingAttributes_AttributeSettingRestriction() {
	_ = this_
	this_.property = j2gort.NullString
}

// Java constructor: AttributeSettingRestriction(com.google.javascript.rhino.jstype.JSType type, java.lang.String property)
func NewConformanceRules_BanSettingAttributes_AttributeSettingRestriction(type_ pkg_rhino.JSType, property string) *ConformanceRules_BanSettingAttributes_AttributeSettingRestriction {
	this_ := &ConformanceRules_BanSettingAttributes_AttributeSettingRestriction{}
	InitConformanceRules_BanSettingAttributes_AttributeSettingRestriction(this_, type_, property)
	return this_
}

func InitConformanceRules_BanSettingAttributes_AttributeSettingRestriction(this_ *ConformanceRules_BanSettingAttributes_AttributeSettingRestriction, type_ pkg_rhino.JSType, property string) {
	_ = type_
	_ = property
	this_.initFields_ConformanceRules_BanSettingAttributes_AttributeSettingRestriction()
	this_.type_ = type_
	this_.property = property
}

// Java class: com.google.javascript.jscomp.ConformanceRules.BanExecCommand
type ConformanceRules_BanExecCommand struct {
	ConformanceRules_AbstractRuleBase
	bannedAttrs []string `gson:"bannedAttrs"`
}

func (this_ *ConformanceRules_BanExecCommand) BindSelf_(s any) {
	_ = s
	this_.ConformanceRules_AbstractRuleBase.BindSelf_(s)
}
func (this_ *ConformanceRules_BanExecCommand) IsJavaNil_() bool { return this_ == nil }
func (*ConformanceRules_BanExecCommand) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConformanceRules$BanExecCommand", "BanExecCommand"
}
func (this_ *ConformanceRules_BanExecCommand) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ConformanceRules_AbstractRuleBase.IdentityHash_()
}
func (*ConformanceRules_BanExecCommand) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConformanceRules.BanExecCommand", "com.google.javascript.jscomp.ConformanceRules$BanExecCommand", func(v_ any) bool {
		return j2gort.InstanceOf[*ConformanceRules_BanExecCommand](v_)
	})
}
func (*ConformanceRules_BanExecCommand) Implements_jscomp_ConformanceRules_AbstractRule() {}
func (*ConformanceRules_BanExecCommand) Implements_jscomp_CheckConformance_Rule()         {}

func (this_ *ConformanceRules_BanExecCommand) initFields_ConformanceRules_BanExecCommand() {
	_ = this_
}

// Java constructor: BanExecCommand(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.Requirement requirement)
func NewConformanceRules_BanExecCommand(compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) *ConformanceRules_BanExecCommand {
	this_ := &ConformanceRules_BanExecCommand{}
	this_.BindSelf_(this_)
	InitConformanceRules_BanExecCommand(this_, compiler, requirement)
	return this_
}

func InitConformanceRules_BanExecCommand(this_ *ConformanceRules_BanExecCommand, compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) {
	_ = compiler
	_ = requirement
	InitConformanceRules_AbstractRule(&this_.ConformanceRules_AbstractRuleBase, compiler, requirement)
	this_.initFields_ConformanceRules_BanExecCommand()
	this_.bannedAttrs = j2gort.ImmListToSliceNullable[string](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableList](j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](requirement.GetValueList()).Stream().Map(lambdaSingleton_conformance_rules_36).Collect_1(j2gort.Ext_com_google_common_collect_ImmutableList_ToImmutableList())))
}

// Java method: BanExecCommand.checkConformance(com.google.javascript.jscomp.NodeTraversal traversal, com.google.javascript.rhino.Node node)
func (this_ *ConformanceRules_BanExecCommand) CheckConformance_2(traversal *NodeTraversal, node pkg_rhino.Node) *ConformanceRules_ConformanceResult {
	_ = traversal
	_ = node
	if !(this_.priv_BanExecCommand_isBannedProperty_2176(node)) {
		return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
	}
	if node.GetChildCount() < int32(2) {
		return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
	}
	var attr pkg_rhino.Node = node.GetSecondChild()
	_ = attr
	if !(attr.IsStringLit()) {
		return StaticGet_ConformanceRules_ConformanceResult_VIOLATION()
	}
	var attrName string = attr.GetString()
	_ = attrName
	if j2gort.SliceContains(this_.bannedAttrs, javalang.StringToLowerCaseASCII(attrName)) {
		return StaticGet_ConformanceRules_ConformanceResult_VIOLATION()
	}
	return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
}

// Java method: BanExecCommand.isBannedProperty(com.google.javascript.rhino.Node node)
func (this_ *ConformanceRules_BanExecCommand) priv_BanExecCommand_isBannedProperty_2176(node pkg_rhino.Node) bool {
	_ = node
	if !(node.IsCall()) {
		return false
	}
	var target pkg_rhino.Node = node.GetFirstChild()
	_ = target
	if !(target.IsGetProp()) {
		return false
	}
	var propertyName string = target.GetString()
	_ = propertyName
	if !(j2gort.CheckNotNullString(propertyName) == "execCommand") {
		return false
	}
	var type_ pkg_rhino.JSType = target.GetFirstChild().GetJSType()
	_ = type_
	if type_ == nil || type_.IsJavaNil_() {
		return false
	}
	var documentType pkg_rhino.JSType = this_.GetConformanceRules_AbstractRuleBase().compiler.GetTypeRegistry().GetGlobalType("Document")
	_ = documentType
	if (documentType == nil || documentType.IsJavaNil_()) || !(type_.RestrictByNotNullOrUndefined().IsSubtypeOf_1(documentType)) {
		return false
	}
	return true
}

// Java class: com.google.javascript.jscomp.ConformanceRules.BanStaticThis
type ConformanceRules_BanStaticThis struct {
	ConformanceRules_AbstractTypeRestrictionRuleBase
}

func (this_ *ConformanceRules_BanStaticThis) BindSelf_(s any) {
	_ = s
	this_.ConformanceRules_AbstractTypeRestrictionRuleBase.BindSelf_(s)
}
func (this_ *ConformanceRules_BanStaticThis) IsJavaNil_() bool { return this_ == nil }
func (*ConformanceRules_BanStaticThis) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConformanceRules$BanStaticThis", "BanStaticThis"
}
func (this_ *ConformanceRules_BanStaticThis) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ConformanceRules_AbstractTypeRestrictionRuleBase.IdentityHash_()
}
func (*ConformanceRules_BanStaticThis) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConformanceRules.BanStaticThis", "com.google.javascript.jscomp.ConformanceRules$BanStaticThis", func(v_ any) bool {
		return j2gort.InstanceOf[*ConformanceRules_BanStaticThis](v_)
	})
}
func (*ConformanceRules_BanStaticThis) Implements_jscomp_ConformanceRules_AbstractTypeRestrictionRule() {
}
func (*ConformanceRules_BanStaticThis) Implements_jscomp_ConformanceRules_AbstractRule() {}
func (*ConformanceRules_BanStaticThis) Implements_jscomp_CheckConformance_Rule()         {}

func (this_ *ConformanceRules_BanStaticThis) initFields_ConformanceRules_BanStaticThis() {
	_ = this_
}

// Java constructor: BanStaticThis(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.Requirement requirement)
func NewConformanceRules_BanStaticThis(compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) *ConformanceRules_BanStaticThis {
	this_ := &ConformanceRules_BanStaticThis{}
	this_.BindSelf_(this_)
	InitConformanceRules_BanStaticThis(this_, compiler, requirement)
	return this_
}

func InitConformanceRules_BanStaticThis(this_ *ConformanceRules_BanStaticThis, compiler AbstractCompiler, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) {
	_ = compiler
	_ = requirement
	InitConformanceRules_AbstractTypeRestrictionRule(&this_.ConformanceRules_AbstractTypeRestrictionRuleBase, compiler, requirement)
	this_.initFields_ConformanceRules_BanStaticThis()
}

// Java method: BanStaticThis.checkConformance(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *ConformanceRules_BanStaticThis) CheckConformance_2(t *NodeTraversal, n pkg_rhino.Node) *ConformanceRules_ConformanceResult {
	_ = t
	_ = n
	if !(n.IsThis()) {
		return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
	}
	var enclosingFunction pkg_rhino.Node = NodeUtil_GetEnclosingNonArrowFunction(n)
	_ = enclosingFunction
	if (enclosingFunction != 0) && ConformanceRules_BanStaticThis_IsStaticMethod(enclosingFunction) {
		return StaticGet_ConformanceRules_ConformanceResult_VIOLATION()
	}
	return StaticGet_ConformanceRules_ConformanceResult_CONFORMANCE()
}

// Java static method: BanStaticThis.isStaticMethod(com.google.javascript.rhino.Node n)
func ConformanceRules_BanStaticThis_IsStaticMethod(n pkg_rhino.Node) bool {
	_ = n
	if !(n.IsFunction()) {
		j2gort.CheckArgument(false)
	}
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	return (((parent != 0) && parent.IsMemberFunctionDef()) && parent.IsStaticMember())
}
