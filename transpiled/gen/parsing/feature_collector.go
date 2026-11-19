/*
 * Copyright 2025 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/parsing/FeatureCollector.java (generated mechanically by j2go).

package parsing

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.parsing.FeatureCollector
type FeatureCollector struct {
	scriptFeatures *pkg_parser.FeatureSet `gson:"scriptFeatures"`
	idHash_        atomic.Int32
}

func (this_ *FeatureCollector) BindSelf_(s any) {
	_ = s
}
func (this_ *FeatureCollector) IsJavaNil_() bool { return this_ == nil }
func (*FeatureCollector) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.FeatureCollector", "FeatureCollector"
}
func (this_ *FeatureCollector) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*FeatureCollector) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.FeatureCollector", "com.google.javascript.jscomp.parsing.FeatureCollector", func(v_ any) bool {
		return j2gort.InstanceOf[*FeatureCollector](v_)
	})
}

func (this_ *FeatureCollector) initFields_FeatureCollector() {
	_ = this_
	this_.scriptFeatures = pkg_parser.StaticGet_FeatureSet_BARE_MINIMUM()
}

// Java constructor: FeatureCollector()
func NewFeatureCollector() *FeatureCollector {
	this_ := &FeatureCollector{}
	InitFeatureCollector(this_)
	return this_
}

func InitFeatureCollector(this_ *FeatureCollector) {
	this_.initFields_FeatureCollector()
}

// Java method: FeatureCollector.allFeatures()
func (this_ *FeatureCollector) AllFeatures() *pkg_parser.FeatureSet {
	return this_.scriptFeatures
}

// Java method: FeatureCollector.visitScript(com.google.javascript.rhino.Node script)
func (this_ *FeatureCollector) VisitScript(script pkg_rhino.Node) {
	_ = script
	var emptyContext j2gort.Ext_java_util_ArrayList = j2gort.NewExt_java_util_ArrayList()
	_ = emptyContext
	emptyContext.Add_1(j2gort.BoxNullableVal(StaticGet_FeatureCollector_FeatureContext_NONE()))
	this_.priv_FeatureCollector_visitRecursive_6775(emptyContext, script)
}

// Java method: FeatureCollector.visitRecursive(java.util.List<com.google.javascript.jscomp.parsing.FeatureCollector.FeatureContext> contextStack, com.google.javascript.rhino.Node node)
func (this_ *FeatureCollector) priv_FeatureCollector_visitRecursive_6775(contextStack j2gort.Ext_java_util_List, node pkg_rhino.Node) {
	_ = contextStack
	_ = node
	var newContext FeatureCollector_FeatureContext = this_.VisitSingleNode(contextStack, node)
	_ = newContext
	contextStack.Add_1(j2gort.BoxNullableVal(newContext))
	{
		var child pkg_rhino.Node = node.GetFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			this_.priv_FeatureCollector_visitRecursive_6775(contextStack, child)
		}
	}
	contextStack.Remove_1_int32((contextStack.Size() - int32(1)))
}

// Java method: FeatureCollector.visitSingleNode(java.util.List<com.google.javascript.jscomp.parsing.FeatureCollector.FeatureContext> contextStack, com.google.javascript.rhino.Node node)
func (this_ *FeatureCollector) VisitSingleNode(contextStack j2gort.Ext_java_util_List, node pkg_rhino.Node) FeatureCollector_FeatureContext {
	_ = contextStack
	_ = node
	if !(contextStack.Size() >= int32(1)) {
		j2gort.CheckStateStr(false, "The provided feature context stack was empty. It must have at least one element.")
	}
	this_.priv_FeatureCollector_recordScriptFeatures_6777(contextStack, node)
	return FeatureCollector_ContextFor(contextStack, node)
}

// Java method: FeatureCollector.recordScriptFeatures(java.util.List<com.google.javascript.jscomp.parsing.FeatureCollector.FeatureContext> contextStack, com.google.javascript.rhino.Node node)
func (this_ *FeatureCollector) priv_FeatureCollector_recordScriptFeatures_6777(contextStack j2gort.Ext_java_util_List, node pkg_rhino.Node) {
	_ = contextStack
	_ = node
	var parentContext FeatureCollector_FeatureContext = j2gort.Cast[FeatureCollector_FeatureContext](contextStack.GetLast())
	_ = parentContext
	switch node.TokenSwitchOrd_() {
	case int32(73) /* Token.FUNCTION */ :
		if node.IsAsyncGeneratorFunction() {
			this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_ASYNC_GENERATORS())
		}
		if node.IsArrowFunction() {
			this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_ARROW_FUNCTIONS())
		}
		if node.IsAsyncFunction() {
			this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_ASYNC_FUNCTIONS())
		}
		if node.IsGeneratorFunction() {
			this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_GENERATORS())
		}
		{
			var iifeCond_1 bool
			{

				re_1 := parentContext
				ae_2 := StaticGet_FeatureCollector_FeatureContext_BLOCK_SCOPE()
				if re_1 == 0 {
					j2gort.ThrowNPE()
				}
				iifeCond_1 = re_1 == ae_2
			}
			if iifeCond_1 {
				this_.scriptFeatures = this_.scriptFeatures.With_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_BLOCK_SCOPED_FUNCTION_DECLARATION())
			}
		}
	case int32(49) /* Token.PARAM_LIST */, int32(28) /* Token.CALL */, int32(23) /* Token.NEW */ :
	case int32(100) /* Token.STRING_KEY */ :
		if node.IsShorthandProperty() && func() bool {
			re_3 := parentContext
			ae_4 := StaticGet_FeatureCollector_FeatureContext_OBJECT_LITERAL()
			if re_3 == 0 {
				j2gort.ThrowNPE()
			}
			return re_3 == ae_4
		}() {
			this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_SHORTHAND_OBJECT_PROPERTIES())
		}
	case int32(134) /* Token.DEFAULT_VALUE */ :
		{
			var iifeCond_1 bool
			{

				re_5 := parentContext
				ae_6 := StaticGet_FeatureCollector_FeatureContext_PARAM_LIST()
				if re_5 == 0 {
					j2gort.ThrowNPE()
				}
				iifeCond_1 = re_5 == ae_6
			}
			if iifeCond_1 {
				this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_DEFAULT_PARAMETERS())
			}
		}
	case int32(95) /* Token.GETTER_DEF */ :
		this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_GETTER())
		{
			var iifeCond_1 bool
			{

				re_7 := parentContext
				ae_8 := StaticGet_FeatureCollector_FeatureContext_CLASS_MEMBERS()
				if re_7 == 0 {
					j2gort.ThrowNPE()
				}
				iifeCond_1 = re_7 == ae_8
			}
			if iifeCond_1 {
				this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_CLASS_GETTER_SETTER())
			}
		}
		if node.IsPrivateIdentifier() {
			this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_PRIVATE_ELEMENTS())
		}
	case int32(42) /* Token.REGEXP */ :
		this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_REGEXP_SYNTAX())
	case int32(96) /* Token.SETTER_DEF */ :
		this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_SETTER())
		{
			var iifeCond_1 bool
			{

				re_9 := parentContext
				ae_10 := StaticGet_FeatureCollector_FeatureContext_CLASS_MEMBERS()
				if re_9 == 0 {
					j2gort.ThrowNPE()
				}
				iifeCond_1 = re_9 == ae_10
			}
			if iifeCond_1 {
				this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_CLASS_GETTER_SETTER())
			}
		}
		if node.IsPrivateIdentifier() {
			this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_PRIVATE_ELEMENTS())
		}
	case int32(90) /* Token.BLOCK */ :
		{
			var iifeCond_1 bool
			{

				re_11 := parentContext
				ae_12 := StaticGet_FeatureCollector_FeatureContext_CLASS_MEMBERS()
				if re_11 == 0 {
					j2gort.ThrowNPE()
				}
				iifeCond_1 = re_11 == ae_12
			}
			if iifeCond_1 {
				this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_CLASS_STATIC_BLOCK())
			}
		}
	case int32(88) /* Token.EMPTY */ :
		{
			var iifeCond_1 bool
			{

				re_13 := parentContext
				ae_14 := StaticGet_FeatureCollector_FeatureContext_CATCH()
				if re_13 == 0 {
					j2gort.ThrowNPE()
				}
				iifeCond_1 = re_13 == ae_14
			}
			if iifeCond_1 {
				this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_OPTIONAL_CATCH_BINDING())
			}
		}
	case int32(125) /* Token.ITER_REST */ :
		{
			var iifeCond_1 bool
			{

				re_15 := parentContext
				ae_16 := StaticGet_FeatureCollector_FeatureContext_PARAM_LIST()
				if re_15 == 0 {
					j2gort.ThrowNPE()
				}
				iifeCond_1 = re_15 == ae_16
			}
			if iifeCond_1 {
				this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_REST_PARAMETERS())
			} else {
				this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_ARRAY_PATTERN_REST())
			}
		}
	case int32(127) /* Token.ITER_SPREAD */ :
		this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_SPREAD_EXPRESSIONS())
	case int32(126) /* Token.OBJECT_REST */ :
		this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_OBJECT_PATTERN_REST())
	case int32(128) /* Token.OBJECT_SPREAD */ :
		this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_OBJECT_LITERALS_WITH_SPREAD())
	case int32(34) /* Token.BIGINT */ :
		this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_BIGINT())
	case int32(18) /* Token.EXPONENT */, int32(63) /* Token.ASSIGN_EXPONENT */ :
		this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_EXPONENT_OP())
	case int32(130) /* Token.TAGGED_TEMPLATELIT */, int32(131) /* Token.TEMPLATELIT */ :
		this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_TEMPLATE_LITERALS())
	case int32(135) /* Token.NEW_TARGET */ :
		this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_NEW_TARGET())
	case int32(129) /* Token.COMPUTED_PROP */ :
		this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_COMPUTED_PROPERTIES())
		var isGetter bool = node.GetBooleanProp(pkg_rhino.StaticGet_Node_COMPUTED_PROP_GETTER())
		_ = isGetter
		var isSetter bool = node.GetBooleanProp(pkg_rhino.StaticGet_Node_COMPUTED_PROP_SETTER())
		_ = isSetter
		var isClassMember bool = func() bool {
			re_17 := parentContext
			ae_18 := StaticGet_FeatureCollector_FeatureContext_CLASS_MEMBERS()
			if re_17 == 0 {
				j2gort.ThrowNPE()
			}
			return re_17 == ae_18
		}()
		_ = isClassMember
		if isGetter {
			this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_GETTER())
		} else if isSetter {
			this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_SETTER())
		}
		if (isGetter || isSetter) && isClassMember {
			this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_CLASS_GETTER_SETTER())
		}
	case int32(29) /* Token.OPTCHAIN_GETPROP */, int32(31) /* Token.OPTCHAIN_CALL */, int32(30) /* Token.OPTCHAIN_GETELEM */ :
		this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_OPTIONAL_CHAINING())
	case int32(70) /* Token.COALESCE */ :
		this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_NULL_COALESCE_OP())
	case int32(124) /* Token.DYNAMIC_IMPORT */ :
		this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_DYNAMIC_IMPORT())
	case int32(64) /* Token.ASSIGN_OR */, int32(65) /* Token.ASSIGN_AND */ :
		this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_LOGICAL_ASSIGNMENT())
	case int32(66) /* Token.ASSIGN_COALESCE */ :
		this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_NULL_COALESCE_OP())
		this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_LOGICAL_ASSIGNMENT())
	case int32(112) /* Token.FOR_OF */ :
		this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_FOR_OF())
	case int32(113) /* Token.FOR_AWAIT_OF */ :
		this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_FOR_AWAIT_OF())
		this_.priv_FeatureCollector_maybeAddTopLevelAwait_6778(contextStack)
	case int32(116) /* Token.IMPORT */, int32(120) /* Token.EXPORT */ :
		this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_MODULES())
	case int32(97) /* Token.CONST */ :
		this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_CONST_DECLARATIONS())
	case int32(111) /* Token.LET */ :
		this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_LET_DECLARATIONS())
	case int32(105) /* Token.CLASS */ :
		this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_CLASSES())
	case int32(107) /* Token.MEMBER_FUNCTION_DEF */ :
		this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_MEMBER_DECLARATIONS())
		if node.IsPrivateIdentifier() {
			this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_PRIVATE_ELEMENTS())
		}
	case int32(108) /* Token.MEMBER_FIELD_DEF */ :
		if node.IsPrivateIdentifier() {
			this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_PRIVATE_ELEMENTS())
		} else {
			this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_PUBLIC_CLASS_FIELDS())
		}
	case int32(44) /* Token.IN */ :
		if (node.GetFirstChild() != 0) && node.GetFirstChild().IsPrivateIdentifier() {
			this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_PRIVATE_ELEMENTS())
		}
	case int32(109) /* Token.COMPUTED_FIELD_DEF */ :
		this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_PUBLIC_CLASS_FIELDS())
	case int32(110) /* Token.SUPER */ :
		this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_SUPER())
	case int32(102) /* Token.ARRAY_PATTERN */ :
		this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_ARRAY_DESTRUCTURING())
	case int32(103) /* Token.OBJECT_PATTERN */ :
		this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_OBJECT_DESTRUCTURING())
	case int32(115) /* Token.AWAIT */ :
		this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_ASYNC_FUNCTIONS())
		this_.priv_FeatureCollector_maybeAddTopLevelAwait_6778(contextStack)
	default:
	}
}

// Java method: FeatureCollector.maybeAddTopLevelAwait(java.util.List<com.google.javascript.jscomp.parsing.FeatureCollector.FeatureContext> contextStack)
func (this_ *FeatureCollector) priv_FeatureCollector_maybeAddTopLevelAwait_6778(contextStack j2gort.Ext_java_util_List) {
	_ = contextStack
	if !(contextStack.Contains(j2gort.BoxNullableVal(StaticGet_FeatureCollector_FeatureContext_FUNCTION()))) {
		this_.priv_FeatureCollector_addScriptFeature_6779(pkg_parser.StaticGet_FeatureSet_Feature_TOP_LEVEL_AWAIT())
	}
}

// Java method: FeatureCollector.addScriptFeature(com.google.javascript.jscomp.parsing.parser.FeatureSet.Feature x)
func (this_ *FeatureCollector) priv_FeatureCollector_addScriptFeature_6779(x *pkg_parser.FeatureSet_Feature) {
	_ = x
	this_.scriptFeatures = this_.scriptFeatures.With_1_FeatureSet_Feature(x)
}

// Java static method: FeatureCollector.contextFor(java.util.List<com.google.javascript.jscomp.parsing.FeatureCollector.FeatureContext> contextStack, com.google.javascript.rhino.Node node)
func FeatureCollector_ContextFor(contextStack j2gort.Ext_java_util_List, node pkg_rhino.Node) FeatureCollector_FeatureContext {
	_ = contextStack
	_ = node
	switch node.TokenSwitchOrd_() {
	case int32(49) /* Token.PARAM_LIST */ :
		return StaticGet_FeatureCollector_FeatureContext_PARAM_LIST()
	case int32(106) /* Token.CLASS_MEMBERS */ :
		return StaticGet_FeatureCollector_FeatureContext_CLASS_MEMBERS()
	case int32(105) /* Token.CLASS */ :
		return StaticGet_FeatureCollector_FeatureContext_CLASS()
	case int32(86) /* Token.CATCH */ :
		return StaticGet_FeatureCollector_FeatureContext_CATCH()
	case int32(90) /* Token.BLOCK */ :
		return func() FeatureCollector_FeatureContext {
			{
				var iifeCond_1 bool
				{

					re_19 := j2gort.Cast[FeatureCollector_FeatureContext](contextStack.GetLast())
					ae_20 := StaticGet_FeatureCollector_FeatureContext_FUNCTION()
					if re_19 == 0 {
						j2gort.ThrowNPE()
					}
					iifeCond_1 = re_19 == ae_20
				}
				if iifeCond_1 {
					return StaticGet_FeatureCollector_FeatureContext_NONE()
				}
			}
			return StaticGet_FeatureCollector_FeatureContext_BLOCK_SCOPE()
		}()
	case int32(73) /* Token.FUNCTION */ :
		return StaticGet_FeatureCollector_FeatureContext_FUNCTION()
	case int32(47) /* Token.OBJECTLIT */ :
		return StaticGet_FeatureCollector_FeatureContext_OBJECT_LITERAL()
	default:
		return StaticGet_FeatureCollector_FeatureContext_NONE()
	}
	panic("unreachable switch expression")
}

// Java enum: com.google.javascript.jscomp.parsing.FeatureCollector.FeatureContext
type FeatureCollector_FeatureContext uint8

const (
	enum_FeatureCollector_FeatureContext_PARAM_LIST     FeatureCollector_FeatureContext = 1
	enum_FeatureCollector_FeatureContext_CLASS_MEMBERS  FeatureCollector_FeatureContext = 2
	enum_FeatureCollector_FeatureContext_CLASS          FeatureCollector_FeatureContext = 3
	enum_FeatureCollector_FeatureContext_CATCH          FeatureCollector_FeatureContext = 4
	enum_FeatureCollector_FeatureContext_BLOCK_SCOPE    FeatureCollector_FeatureContext = 5
	enum_FeatureCollector_FeatureContext_FUNCTION       FeatureCollector_FeatureContext = 6
	enum_FeatureCollector_FeatureContext_OBJECT_LITERAL FeatureCollector_FeatureContext = 7
	enum_FeatureCollector_FeatureContext_NONE           FeatureCollector_FeatureContext = 8
)

var enumBaseTable_FeatureCollector_FeatureContext = [...]j2gort.EnumBase{
	{Name_: "PARAM_LIST", Ordinal_: 0},
	{Name_: "CLASS_MEMBERS", Ordinal_: 1},
	{Name_: "CLASS", Ordinal_: 2},
	{Name_: "CATCH", Ordinal_: 3},
	{Name_: "BLOCK_SCOPE", Ordinal_: 4},
	{Name_: "FUNCTION", Ordinal_: 5},
	{Name_: "OBJECT_LITERAL", Ordinal_: 6},
	{Name_: "NONE", Ordinal_: 7},
}
var enumValues_FeatureCollector_FeatureContext = []FeatureCollector_FeatureContext{1, 2, 3, 4, 5, 6, 7, 8}

func StaticGet_FeatureCollector_FeatureContext_PARAM_LIST() FeatureCollector_FeatureContext {
	return enum_FeatureCollector_FeatureContext_PARAM_LIST
}
func StaticGet_FeatureCollector_FeatureContext_CLASS_MEMBERS() FeatureCollector_FeatureContext {
	return enum_FeatureCollector_FeatureContext_CLASS_MEMBERS
}
func StaticGet_FeatureCollector_FeatureContext_CLASS() FeatureCollector_FeatureContext {
	return enum_FeatureCollector_FeatureContext_CLASS
}
func StaticGet_FeatureCollector_FeatureContext_CATCH() FeatureCollector_FeatureContext {
	return enum_FeatureCollector_FeatureContext_CATCH
}
func StaticGet_FeatureCollector_FeatureContext_BLOCK_SCOPE() FeatureCollector_FeatureContext {
	return enum_FeatureCollector_FeatureContext_BLOCK_SCOPE
}
func StaticGet_FeatureCollector_FeatureContext_FUNCTION() FeatureCollector_FeatureContext {
	return enum_FeatureCollector_FeatureContext_FUNCTION
}
func StaticGet_FeatureCollector_FeatureContext_OBJECT_LITERAL() FeatureCollector_FeatureContext {
	return enum_FeatureCollector_FeatureContext_OBJECT_LITERAL
}
func StaticGet_FeatureCollector_FeatureContext_NONE() FeatureCollector_FeatureContext {
	return enum_FeatureCollector_FeatureContext_NONE
}

func FeatureCollector_FeatureContext_Values() []FeatureCollector_FeatureContext {
	return j2gort.CloneSlice(enumValues_FeatureCollector_FeatureContext)
}

func FeatureCollector_FeatureContext_ValueOf(name string) FeatureCollector_FeatureContext {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_FeatureCollector_FeatureContext {
		if enumBaseTable_FeatureCollector_FeatureContext[i_].Name_ == name {
			return FeatureCollector_FeatureContext(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.parsing.FeatureCollector.FeatureContext." + name))
}

func (this_ FeatureCollector_FeatureContext) IsJavaNil_() bool                              { return this_ == 0 }
func (FeatureCollector_FeatureContext) HasIdentityEquals_()                                 {}
func (FeatureCollector_FeatureContext) Implements_j2gort_Ext_java_lang_Enum()               {}
func (FeatureCollector_FeatureContext) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ FeatureCollector_FeatureContext) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (FeatureCollector_FeatureContext) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.FeatureCollector$FeatureContext", "FeatureContext"
}
func (this_ FeatureCollector_FeatureContext) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_FeatureCollector_FeatureContext[this_-1]
}
func (this_ FeatureCollector_FeatureContext) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ FeatureCollector_FeatureContext) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ FeatureCollector_FeatureContext) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_FeatureCollector_FeatureContext[this_-1].Name_
}
func (this_ FeatureCollector_FeatureContext) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_FeatureCollector_FeatureContext[this_-1].Name_
}
func (this_ FeatureCollector_FeatureContext) Equals(other_ any) bool {
	o_, ok_ := other_.(FeatureCollector_FeatureContext)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ FeatureCollector_FeatureContext) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_FeatureCollector_FeatureContext[this_-1].HashCode()
}
func (this_ FeatureCollector_FeatureContext) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_FeatureCollector_FeatureContext[this_-1].CompareTo(other_)
}
func (this_ FeatureCollector_FeatureContext) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_FeatureCollector_FeatureContext[this_-1].Name_)
}
func init() {
	vals_ := []any{FeatureCollector_FeatureContext(1), FeatureCollector_FeatureContext(2), FeatureCollector_FeatureContext(3), FeatureCollector_FeatureContext(4), FeatureCollector_FeatureContext(5), FeatureCollector_FeatureContext(6), FeatureCollector_FeatureContext(7), FeatureCollector_FeatureContext(8)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.parsing.FeatureCollector.FeatureContext", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.parsing.FeatureCollector$FeatureContext", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.FeatureCollector.FeatureContext", "com.google.javascript.jscomp.parsing.FeatureCollector$FeatureContext", func(v_ any) bool {
		return j2gort.InstanceOf[FeatureCollector_FeatureContext](v_)
	})
}
