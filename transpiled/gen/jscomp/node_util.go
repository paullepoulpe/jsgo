/*
 * Copyright 2004 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/NodeUtil.java (generated mechanically by j2go).

package jscomp

import (
	"fmt"
	"math"
	"sync/atomic"

	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	pkg_colors "github.com/paullepoulpe/jsgo/transpiled/gen/colors"
	pkg_dtoa "github.com/paullepoulpe/jsgo/transpiled/gen/dtoa"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_parsing "github.com/paullepoulpe/jsgo/transpiled/gen/parsing"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_node_util_1 j2gort.Ext_com_google_common_base_Predicate
var lambdaSingleton_node_util_2 j2gort.Ext_com_google_common_base_Predicate
var lambdaSingleton_node_util_3 j2gort.Ext_com_google_common_base_Predicate
var lambdaSingleton_node_util_4 j2gort.Ext_com_google_common_base_Predicate
var lambdaSingleton_node_util_5 j2gort.Ext_com_google_common_base_Predicate
var lambdaSingleton_node_util_6 j2gort.Ext_com_google_common_base_Predicate
var lambdaSingleton_node_util_7 j2gort.Ext_com_google_common_base_Predicate
var lambdaSingleton_node_util_8 j2gort.Ext_com_google_common_base_Predicate
var lambdaSingleton_node_util_9 j2gort.Ext_com_google_common_base_Predicate
var lambdaSingleton_node_util_10 j2gort.Ext_com_google_common_base_Predicate
var lambdaSingleton_node_util_11 j2gort.Ext_com_google_common_base_Predicate
var lambdaSingleton_node_util_12 j2gort.Ext_com_google_common_base_Predicate
var lambdaSingleton_node_util_13 j2gort.Ext_com_google_common_base_Predicate
var lambdaSingleton_node_util_14 j2gort.Ext_com_google_common_base_Predicate
var lambdaSingleton_node_util_15 j2gort.Ext_com_google_common_base_Predicate
var lambdaSingleton_node_util_16 j2gort.Ext_com_google_common_base_Predicate
var lambdaSingleton_node_util_17 j2gort.Ext_com_google_common_base_Predicate

func init() {
	lambdaSingleton_node_util_1 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_1 uint32) bool {
		n := pkg_rhino.Node(rawNode_1)
		_ = n
		return !(n.IsFunction())
	})
	lambdaSingleton_node_util_2 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_2 uint32) bool {
		n := pkg_rhino.Node(rawNode_2)
		_ = n
		return !(NodeUtil_IsNonArrowFunction(n))
	})
	lambdaSingleton_node_util_3 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_ uint32) bool { return pkg_rhino.Node(rawNode_).IsBlock() })
	lambdaSingleton_node_util_4 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_21 uint32) bool {
		c := pkg_rhino.Node(rawNode_21)
		_ = c
		return true
	})
	lambdaSingleton_node_util_5 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_ uint32) bool { return NodeUtil_IsNonArrowFunction(pkg_rhino.Node(rawNode_)) })
	lambdaSingleton_node_util_6 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_ uint32) bool { return pkg_rhino.Node(rawNode_).IsClass() })
	lambdaSingleton_node_util_7 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_ uint32) bool { return pkg_rhino.Node(rawNode_).IsModuleBody() })
	lambdaSingleton_node_util_8 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_ uint32) bool { return pkg_rhino.Node(rawNode_).IsFunction() })
	lambdaSingleton_node_util_9 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_ uint32) bool { return pkg_rhino.Node(rawNode_).IsScript() })
	lambdaSingleton_node_util_10 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_ uint32) bool { return pkg_rhino.Node(rawNode_).IsBlock() })
	lambdaSingleton_node_util_11 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_ uint32) bool { return NodeUtil_CreatesBlockScope(pkg_rhino.Node(rawNode_)) })
	lambdaSingleton_node_util_12 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_ uint32) bool { return NodeUtil_CreatesScope(pkg_rhino.Node(rawNode_)) })
	lambdaSingleton_node_util_13 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_ uint32) bool { return NodeUtil_IsHoistScopeRoot(pkg_rhino.Node(rawNode_)) })
	lambdaSingleton_node_util_14 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_ uint32) bool { return NodeUtil_IsStatement(pkg_rhino.Node(rawNode_)) })
	lambdaSingleton_node_util_15 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_23 uint32) bool {
		c := pkg_rhino.Node(rawNode_23)
		_ = c
		return (c.IsThis() || c.IsSuper())
	})
	lambdaSingleton_node_util_16 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_ uint32) bool { return pkg_rhino.Node(rawNode_).IsSuper() })
	lambdaSingleton_node_util_17 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_24 uint32) bool {
		node := pkg_rhino.Node(rawNode_24)
		_ = node
		return !(node.IsClass())
	})
}

// Java class: com.google.javascript.jscomp.NodeUtil
type NodeUtil struct {
	idHash_ atomic.Int32
}

func (this_ *NodeUtil) BindSelf_(s any) {
	_ = s
}
func (this_ *NodeUtil) IsJavaNil_() bool { return this_ == nil }
func (*NodeUtil) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.NodeUtil", "NodeUtil"
}
func (this_ *NodeUtil) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*NodeUtil) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.NodeUtil", "com.google.javascript.jscomp.NodeUtil", func(v_ any) bool {
		return j2gort.InstanceOf[*NodeUtil](v_)
	})
}

func StaticGet_NodeUtil_MAX_POSITIVE_INTEGER_NUMBER() int64 { return int64(9007199254740991) }
func StaticGet_NodeUtil_JSC_PROPERTY_NAME_FN() string       { return "JSCompiler_renameProperty" }
func StaticGet_NodeUtil_LARGEST_BASIC_LATIN() uint16        { return uint16(127) }

var static_NodeUtil_GOOG_MODULE_DECLARE_LEGACY_NAMESPACE pkg_rhino.QualifiedName

func StaticGet_NodeUtil_GOOG_MODULE_DECLARE_LEGACY_NAMESPACE() pkg_rhino.QualifiedName {
	clinit_NodeUtil()
	return static_NodeUtil_GOOG_MODULE_DECLARE_LEGACY_NAMESPACE
}
func StaticSet_NodeUtil_GOOG_MODULE_DECLARE_LEGACY_NAMESPACE(v pkg_rhino.QualifiedName) {
	clinit_NodeUtil()
	static_NodeUtil_GOOG_MODULE_DECLARE_LEGACY_NAMESPACE = v
}

var static_NodeUtil_GOOG_SET_TEST_ONLY pkg_rhino.QualifiedName

func StaticGet_NodeUtil_GOOG_SET_TEST_ONLY() pkg_rhino.QualifiedName {
	clinit_NodeUtil()
	return static_NodeUtil_GOOG_SET_TEST_ONLY
}
func StaticSet_NodeUtil_GOOG_SET_TEST_ONLY(v pkg_rhino.QualifiedName) {
	clinit_NodeUtil()
	static_NodeUtil_GOOG_SET_TEST_ONLY = v
}

var static_NodeUtil_GOOG_PROVIDE pkg_rhino.QualifiedName

func StaticGet_NodeUtil_GOOG_PROVIDE() pkg_rhino.QualifiedName {
	clinit_NodeUtil()
	return static_NodeUtil_GOOG_PROVIDE
}
func StaticSet_NodeUtil_GOOG_PROVIDE(v pkg_rhino.QualifiedName) {
	clinit_NodeUtil()
	static_NodeUtil_GOOG_PROVIDE = v
}

var static_NodeUtil_GOOG_MODULE pkg_rhino.QualifiedName

func StaticGet_NodeUtil_GOOG_MODULE() pkg_rhino.QualifiedName {
	clinit_NodeUtil()
	return static_NodeUtil_GOOG_MODULE
}
func StaticSet_NodeUtil_GOOG_MODULE(v pkg_rhino.QualifiedName) {
	clinit_NodeUtil()
	static_NodeUtil_GOOG_MODULE = v
}

var static_NodeUtil_GOOG_MODULE_GET pkg_rhino.QualifiedName

func StaticGet_NodeUtil_GOOG_MODULE_GET() pkg_rhino.QualifiedName {
	clinit_NodeUtil()
	return static_NodeUtil_GOOG_MODULE_GET
}
func StaticSet_NodeUtil_GOOG_MODULE_GET(v pkg_rhino.QualifiedName) {
	clinit_NodeUtil()
	static_NodeUtil_GOOG_MODULE_GET = v
}

var static_NodeUtil_GOOG_REQUIRE pkg_rhino.QualifiedName

func StaticGet_NodeUtil_GOOG_REQUIRE() pkg_rhino.QualifiedName {
	clinit_NodeUtil()
	return static_NodeUtil_GOOG_REQUIRE
}
func StaticSet_NodeUtil_GOOG_REQUIRE(v pkg_rhino.QualifiedName) {
	clinit_NodeUtil()
	static_NodeUtil_GOOG_REQUIRE = v
}

var static_NodeUtil_GOOG_REQUIRE_TYPE pkg_rhino.QualifiedName

func StaticGet_NodeUtil_GOOG_REQUIRE_TYPE() pkg_rhino.QualifiedName {
	clinit_NodeUtil()
	return static_NodeUtil_GOOG_REQUIRE_TYPE
}
func StaticSet_NodeUtil_GOOG_REQUIRE_TYPE(v pkg_rhino.QualifiedName) {
	clinit_NodeUtil()
	static_NodeUtil_GOOG_REQUIRE_TYPE = v
}

var static_NodeUtil_GOOG_FORWARD_DECLARE pkg_rhino.QualifiedName

func StaticGet_NodeUtil_GOOG_FORWARD_DECLARE() pkg_rhino.QualifiedName {
	clinit_NodeUtil()
	return static_NodeUtil_GOOG_FORWARD_DECLARE
}
func StaticSet_NodeUtil_GOOG_FORWARD_DECLARE(v pkg_rhino.QualifiedName) {
	clinit_NodeUtil()
	static_NodeUtil_GOOG_FORWARD_DECLARE = v
}

var static_NodeUtil_GOOG_REQUIRE_DYNAMIC pkg_rhino.QualifiedName

func StaticGet_NodeUtil_GOOG_REQUIRE_DYNAMIC() pkg_rhino.QualifiedName {
	clinit_NodeUtil()
	return static_NodeUtil_GOOG_REQUIRE_DYNAMIC
}
func StaticSet_NodeUtil_GOOG_REQUIRE_DYNAMIC(v pkg_rhino.QualifiedName) {
	clinit_NodeUtil()
	static_NodeUtil_GOOG_REQUIRE_DYNAMIC = v
}

var static_NodeUtil_GOOG_WEAK_USAGE pkg_rhino.QualifiedName

func StaticGet_NodeUtil_GOOG_WEAK_USAGE() pkg_rhino.QualifiedName {
	clinit_NodeUtil()
	return static_NodeUtil_GOOG_WEAK_USAGE
}
func StaticSet_NodeUtil_GOOG_WEAK_USAGE(v pkg_rhino.QualifiedName) {
	clinit_NodeUtil()
	static_NodeUtil_GOOG_WEAK_USAGE = v
}

var static_NodeUtil_GOOG_WEAK_USAGE_MANGLED pkg_rhino.QualifiedName

func StaticGet_NodeUtil_GOOG_WEAK_USAGE_MANGLED() pkg_rhino.QualifiedName {
	clinit_NodeUtil()
	return static_NodeUtil_GOOG_WEAK_USAGE_MANGLED
}
func StaticSet_NodeUtil_GOOG_WEAK_USAGE_MANGLED(v pkg_rhino.QualifiedName) {
	clinit_NodeUtil()
	static_NodeUtil_GOOG_WEAK_USAGE_MANGLED = v
}

var static_NodeUtil_KNOWN_CONSTANTS j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_NodeUtil_KNOWN_CONSTANTS() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_NodeUtil()
	return static_NodeUtil_KNOWN_CONSTANTS
}
func StaticSet_NodeUtil_KNOWN_CONSTANTS(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_NodeUtil()
	static_NodeUtil_KNOWN_CONSTANTS = v
}

var static_NodeUtil_DEFINITE_CFG_ROOTS j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_NodeUtil_DEFINITE_CFG_ROOTS() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_NodeUtil()
	return static_NodeUtil_DEFINITE_CFG_ROOTS
}
func StaticSet_NodeUtil_DEFINITE_CFG_ROOTS(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_NodeUtil()
	static_NodeUtil_DEFINITE_CFG_ROOTS = v
}

var static_NodeUtil_IS_STATEMENT_PARENT j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_NodeUtil_IS_STATEMENT_PARENT() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_NodeUtil()
	return static_NodeUtil_IS_STATEMENT_PARENT
}
func StaticSet_NodeUtil_IS_STATEMENT_PARENT(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_NodeUtil()
	static_NodeUtil_IS_STATEMENT_PARENT = v
}

var static_NodeUtil_GLOBAL_OBJECT pkg_rhino.QualifiedName

func StaticGet_NodeUtil_GLOBAL_OBJECT() pkg_rhino.QualifiedName {
	clinit_NodeUtil()
	return static_NodeUtil_GLOBAL_OBJECT
}
func StaticSet_NodeUtil_GLOBAL_OBJECT(v pkg_rhino.QualifiedName) {
	clinit_NodeUtil()
	static_NodeUtil_GLOBAL_OBJECT = v
}

var static_NodeUtil_GLOBAL_OBJECT_MANGLED pkg_rhino.QualifiedName

func StaticGet_NodeUtil_GLOBAL_OBJECT_MANGLED() pkg_rhino.QualifiedName {
	clinit_NodeUtil()
	return static_NodeUtil_GLOBAL_OBJECT_MANGLED
}
func StaticSet_NodeUtil_GLOBAL_OBJECT_MANGLED(v pkg_rhino.QualifiedName) {
	clinit_NodeUtil()
	static_NodeUtil_GLOBAL_OBJECT_MANGLED = v
}

var static_NodeUtil_MATCH_NOT_FUNCTION j2gort.Ext_com_google_common_base_Predicate

func StaticGet_NodeUtil_MATCH_NOT_FUNCTION() j2gort.Ext_com_google_common_base_Predicate {
	clinit_NodeUtil()
	return static_NodeUtil_MATCH_NOT_FUNCTION
}
func StaticSet_NodeUtil_MATCH_NOT_FUNCTION(v j2gort.Ext_com_google_common_base_Predicate) {
	clinit_NodeUtil()
	static_NodeUtil_MATCH_NOT_FUNCTION = v
}

var static_NodeUtil_MATCH_ANYTHING_BUT_NON_ARROW_FUNCTION j2gort.Ext_com_google_common_base_Predicate

func StaticGet_NodeUtil_MATCH_ANYTHING_BUT_NON_ARROW_FUNCTION() j2gort.Ext_com_google_common_base_Predicate {
	clinit_NodeUtil()
	return static_NodeUtil_MATCH_ANYTHING_BUT_NON_ARROW_FUNCTION
}
func StaticSet_NodeUtil_MATCH_ANYTHING_BUT_NON_ARROW_FUNCTION(v j2gort.Ext_com_google_common_base_Predicate) {
	clinit_NodeUtil()
	static_NodeUtil_MATCH_ANYTHING_BUT_NON_ARROW_FUNCTION = v
}

var static_NodeUtil_NUMBER_NAN pkg_rhino.Node

func StaticGet_NodeUtil_NUMBER_NAN() pkg_rhino.Node {
	clinit_NodeUtil()
	return static_NodeUtil_NUMBER_NAN
}
func StaticSet_NodeUtil_NUMBER_NAN(v pkg_rhino.Node) {
	clinit_NodeUtil()
	static_NodeUtil_NUMBER_NAN = v
}

var static_NodeUtil_GOOG_LOADMODULE pkg_rhino.QualifiedName

func StaticGet_NodeUtil_GOOG_LOADMODULE() pkg_rhino.QualifiedName {
	clinit_NodeUtil()
	return static_NodeUtil_GOOG_LOADMODULE
}
func StaticSet_NodeUtil_GOOG_LOADMODULE(v pkg_rhino.QualifiedName) {
	clinit_NodeUtil()
	static_NodeUtil_GOOG_LOADMODULE = v
}

var clinitDone_NodeUtil uint32
var clinitState_NodeUtil j2gort.ClinitState

func clinit_NodeUtil() {
	if atomic.LoadUint32(&clinitDone_NodeUtil) != 0 {
		return
	}
	clinitSlow_NodeUtil()
}

func clinitSlow_NodeUtil() {
	if !j2gort.EnterClinit(&clinitState_NodeUtil, &clinitDone_NodeUtil) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_NodeUtil, &clinitDone_NodeUtil)
	static_NodeUtil_GOOG_MODULE_DECLARE_LEGACY_NAMESPACE = pkg_rhino.QualifiedName_Of("goog.module.declareLegacyNamespace")
	static_NodeUtil_GOOG_SET_TEST_ONLY = pkg_rhino.QualifiedName_Of("goog.setTestOnly")
	static_NodeUtil_GOOG_PROVIDE = pkg_rhino.QualifiedName_Of("goog.provide")
	static_NodeUtil_GOOG_MODULE = pkg_rhino.QualifiedName_Of("goog.module")
	static_NodeUtil_GOOG_MODULE_GET = pkg_rhino.QualifiedName_Of("goog.module.get")
	static_NodeUtil_GOOG_REQUIRE = pkg_rhino.QualifiedName_Of("goog.require")
	static_NodeUtil_GOOG_REQUIRE_TYPE = pkg_rhino.QualifiedName_Of("goog.requireType")
	static_NodeUtil_GOOG_FORWARD_DECLARE = pkg_rhino.QualifiedName_Of("goog.forwardDeclare")
	static_NodeUtil_GOOG_REQUIRE_DYNAMIC = pkg_rhino.QualifiedName_Of("goog.requireDynamic")
	static_NodeUtil_GOOG_WEAK_USAGE = pkg_rhino.QualifiedName_Of("goog.weakUsage")
	static_NodeUtil_GOOG_WEAK_USAGE_MANGLED = pkg_rhino.QualifiedName_Of("goog$weakUsage")
	static_NodeUtil_KNOWN_CONSTANTS = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_3("undefined", "Infinity", "NaN")
	static_NodeUtil_DEFINITE_CFG_ROOTS = j2gort.Ext_com_google_common_collect_Sets_ImmutableEnumSet_2(j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_FUNCTION()), []j2gort.Ext_java_lang_Enum{j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_SCRIPT()), j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_MODULE_BODY()), j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_ROOT())})
	static_NodeUtil_IS_STATEMENT_PARENT = j2gort.Ext_com_google_common_collect_Sets_ImmutableEnumSet_2(j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_SCRIPT()), []j2gort.Ext_java_lang_Enum{j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_MODULE_BODY()), j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_BLOCK()), j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_LABEL()), j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_NAMESPACE_ELEMENTS()), j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_INTERFACE_MEMBERS())})
	static_NodeUtil_GLOBAL_OBJECT = pkg_rhino.QualifiedName_Of("$jscomp.global.Object")
	static_NodeUtil_GLOBAL_OBJECT_MANGLED = pkg_rhino.QualifiedName_Of("$jscomp$global.Object")
	static_NodeUtil_MATCH_NOT_FUNCTION = lambdaSingleton_node_util_1
	static_NodeUtil_MATCH_ANYTHING_BUT_NON_ARROW_FUNCTION = lambdaSingleton_node_util_2
	static_NodeUtil_NUMBER_NAN = pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Name("Number"), "NaN")
	static_NodeUtil_GOOG_LOADMODULE = pkg_rhino.QualifiedName_Of("goog.loadModule")
}

func (this_ *NodeUtil) initFields_NodeUtil() {
	_ = this_
}

// Java constructor: NodeUtil()
func NewNodeUtil() *NodeUtil {
	this_ := &NodeUtil{}
	InitNodeUtil(this_)
	return this_
}

func InitNodeUtil(this_ *NodeUtil) {
	this_.initFields_NodeUtil()
}

// Java static method: NodeUtil.getBooleanValue(com.google.javascript.rhino.Node n)
func NodeUtil_GetBooleanValue(n pkg_rhino.Node) pkg_base.Tri {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(36) /* Token.NULL */, int32(38) /* Token.FALSE */, int32(87) /* Token.VOID */ :
		return pkg_base.StaticGet_Tri_FALSE()
	case int32(39) /* Token.TRUE */, int32(42) /* Token.REGEXP */, int32(73) /* Token.FUNCTION */, int32(105) /* Token.CLASS */, int32(23) /* Token.NEW */, int32(46) /* Token.ARRAYLIT */, int32(47) /* Token.OBJECTLIT */ :
		return pkg_base.StaticGet_Tri_TRUE()
	case int32(131) /* Token.TEMPLATELIT */ :
		if n.HasOneChild() {
			var templateLitString pkg_rhino.Node = n.GetOnlyChild()
			_ = templateLitString
			if !(templateLitString.IsTemplateLitString()) {
				j2gort.CheckStateMsg(false, j2gort.BoxNode(templateLitString))
			}
			var cookedString string = templateLitString.GetCookedString()
			_ = cookedString
			return pkg_base.Tri_ForBoolean(((cookedString != j2gort.NullString) && !(j2gort.CheckNotNullString(cookedString) == "")))
		} else {
			return pkg_base.StaticGet_Tri_UNKNOWN()
		}
	case int32(35) /* Token.STRINGLIT */ :
		return pkg_base.Tri_ForBoolean((j2gort.JStrLen(n.GetString()) > int32(0)))
	case int32(33) /* Token.NUMBER */ :
		return pkg_base.Tri_ForBoolean((n.GetDouble() != float64(int32(0))))
	case int32(34) /* Token.BIGINT */ :
		return pkg_base.Tri_ForBoolean(!(n.GetBigInt().Equals(j2gort.PtrToAny(j2gort.Ext_java_math_BigInteger_F_ZERO))))
	case int32(19) /* Token.NOT */ :
		return NodeUtil_GetBooleanValue(n.GetLastChild()).Not()
	case int32(32) /* Token.NAME */ :
		{

			switch j2gort.CheckNotNullString(n.GetString()) {
			case "undefined", "NaN":
				return pkg_base.StaticGet_Tri_FALSE()
			case "Infinity":
				return pkg_base.StaticGet_Tri_TRUE()
			default:
				return pkg_base.StaticGet_Tri_UNKNOWN()
			}
			panic("unreachable switch expression")
		}

	case int32(20) /* Token.BITNOT */, int32(21) /* Token.POS */, int32(22) /* Token.NEG */ :
		var doubleVal *j2gort.BoxDouble = NodeUtil_GetNumberValue(n)
		_ = doubleVal
		if doubleVal != nil {
			var isFalsey bool = (doubleVal.IsNaN() || pkg_base.JSCompDoubles_IsEitherZero(j2gort.UnboxDouble(doubleVal)))
			_ = isFalsey
			return pkg_base.Tri_ForBoolean(!(isFalsey))
		}
		var bigintVal *j2gort.Ext_java_math_BigInteger = NodeUtil_GetBigIntValue(n)
		_ = bigintVal
		if bigintVal != nil {
			var isFalsey_2 bool = bigintVal.Equals(j2gort.PtrToAny(j2gort.Ext_java_math_BigInteger_F_ZERO))
			_ = isFalsey_2
			return pkg_base.Tri_ForBoolean(!(isFalsey_2))
		}
		return pkg_base.StaticGet_Tri_UNKNOWN()
	case int32(51) /* Token.ASSIGN */, int32(50) /* Token.COMMA */ :
		return NodeUtil_GetBooleanValue(n.GetLastChild())
	case int32(69) /* Token.AND */, int32(65) /* Token.ASSIGN_AND */ :
		var lhs pkg_base.Tri = NodeUtil_GetBooleanValue(n.GetFirstChild())
		_ = lhs
		var rhs pkg_base.Tri = NodeUtil_GetBooleanValue(n.GetLastChild())
		_ = rhs
		return lhs.And(rhs)
	case int32(68) /* Token.OR */, int32(64) /* Token.ASSIGN_OR */ :
		var lhs_2 pkg_base.Tri = NodeUtil_GetBooleanValue(n.GetFirstChild())
		_ = lhs_2
		var rhs_2 pkg_base.Tri = NodeUtil_GetBooleanValue(n.GetLastChild())
		_ = rhs_2
		return lhs_2.Or(rhs_2)
	case int32(67) /* Token.HOOK */ :
		var trueValue pkg_base.Tri = NodeUtil_GetBooleanValue(n.GetSecondChild())
		_ = trueValue
		var falseValue pkg_base.Tri = NodeUtil_GetBooleanValue(n.GetLastChild())
		_ = falseValue
		if j2gort.CheckNotNullEnum(trueValue) == falseValue {
			return trueValue
		} else {
			return pkg_base.StaticGet_Tri_UNKNOWN()
		}
	case int32(70) /* Token.COALESCE */, int32(66) /* Token.ASSIGN_COALESCE */ :
		var lhs_3 pkg_base.Tri = NodeUtil_GetBooleanValue(n.GetFirstChild())
		_ = lhs_3
		var rhs_3 pkg_base.Tri = NodeUtil_GetBooleanValue(n.GetLastChild())
		_ = rhs_3
		{
			var iifeCond_1 bool
			{

				re_3 := lhs_3
				ae_4 := pkg_base.StaticGet_Tri_TRUE()
				if re_3 == 0 {
					j2gort.ThrowNPE()
				}
				iifeCond_1 = re_3 == ae_4
			}
			if iifeCond_1 || (j2gort.CheckNotNullEnum(lhs_3) == rhs_3) {
				return lhs_3
			} else {
				return pkg_base.StaticGet_Tri_UNKNOWN()
			}
		}
	default:
		return pkg_base.StaticGet_Tri_UNKNOWN()
	}
	panic("unreachable switch expression")
}

// Java static method: NodeUtil.getStringValue(com.google.javascript.rhino.Node n)
func NodeUtil_GetStringValue(n pkg_rhino.Node) string {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(35) /* Token.STRINGLIT */, int32(100) /* Token.STRING_KEY */ :
		return n.GetString()
	case int32(131) /* Token.TEMPLATELIT */ :
		var string_ j2gort.LocalStringBuilder

		{
			var child pkg_rhino.Node = n.GetFirstChild()
			_ = child
			for ; child != 0; child = child.GetNext() {
				var expression pkg_rhino.Node = child
				_ = expression
				if child.IsTemplateLitSub() {
					expression = child.GetFirstChild()
				}
				var expressionString string = NodeUtil_GetStringValue(expression)
				_ = expressionString
				if expressionString == j2gort.NullString {
					return j2gort.NullString
				}
				string_.AppendStr(expressionString)
			}
		}
		return string_.String()
	case int32(133) /* Token.TEMPLATELIT_STRING */ :
		return n.GetCookedString()
	case int32(32) /* Token.NAME */ :
		var name string = n.GetString()
		_ = name
		if (("undefined" == name) || ("Infinity" == name)) || ("NaN" == name) {
			return name
		}
		return j2gort.NullString
	case int32(22) /* Token.NEG */, int32(33) /* Token.NUMBER */ :
		var value *j2gort.BoxDouble = NodeUtil_GetNumberValue(n)
		_ = value
		if value != nil {
			return pkg_dtoa.DToA_NumberToString(value.DoubleValue())
		}
		return j2gort.NullString
	case int32(34) /* Token.BIGINT */ :
		return n.GetBigInt().String()
	case int32(38) /* Token.FALSE */ :
		return "false"
	case int32(39) /* Token.TRUE */ :
		return "true"
	case int32(36) /* Token.NULL */ :
		return "null"
	case int32(87) /* Token.VOID */ :
		return "undefined"
	case int32(19) /* Token.NOT */ :
		var child_2 pkg_base.Tri = NodeUtil_GetBooleanValue(n.GetFirstChild())
		_ = child_2
		if child_2 != pkg_base.StaticGet_Tri_UNKNOWN() {
			{

				if child_2.ToBoolean(true) {
					return "false"
				}
				return "true"
			}

		}
		return j2gort.NullString
	case int32(46) /* Token.ARRAYLIT */ :
		return NodeUtil_ArrayToString(n)
	case int32(47) /* Token.OBJECTLIT */ :
		{

			if n.HasChildren() {
				return j2gort.NullString
			}
			return "[object Object]"
		}

	default:
		return j2gort.NullString
	}
	panic("unreachable switch expression")
}

// Java static method: NodeUtil.getArrayElementStringValue(com.google.javascript.rhino.Node n)
func NodeUtil_GetArrayElementStringValue(n pkg_rhino.Node) string {
	_ = n
	if NodeUtil_IsNullOrUndefined(n) || n.IsEmpty() {
		return ""
	} else {
		return NodeUtil_GetStringValue(n)
	}
}

// Java static method: NodeUtil.arrayToString(com.google.javascript.rhino.Node literal)
func NodeUtil_ArrayToString(literal pkg_rhino.Node) string {
	_ = literal
	var first pkg_rhino.Node = literal.GetFirstChild()
	_ = first
	var result j2gort.LocalStringBuilder

	{
		var n pkg_rhino.Node = first
		_ = n
		for ; n != 0; n = n.GetNext() {
			var childValue string = NodeUtil_GetArrayElementStringValue(n)
			_ = childValue
			if childValue == j2gort.NullString {
				return j2gort.NullString
			}
			if n != first {
				result.AppendChar(uint16(44))
			}
			result.AppendStr(childValue)
		}
	}
	return result.String()
}

// Java static method: NodeUtil.getNumberValue(com.google.javascript.rhino.Node n)
func NodeUtil_GetNumberValue(n pkg_rhino.Node) *j2gort.BoxDouble {
	_ = n
	return NodeUtil_DoGetNumberValue(n, true)
}

// Java static method: NodeUtil.getNumberValueNoConversions(com.google.javascript.rhino.Node n)
func NodeUtil_GetNumberValueNoConversions(n pkg_rhino.Node) *j2gort.BoxDouble {
	_ = n
	return NodeUtil_DoGetNumberValue(n, false)
}

// Java static method: NodeUtil.doGetNumberValue(com.google.javascript.rhino.Node n, boolean numberConversions)
func NodeUtil_DoGetNumberValue(n pkg_rhino.Node, numberConversions bool) *j2gort.BoxDouble {
	_ = n
	_ = numberConversions
	switch n.TokenSwitchOrd_() {
	case int32(33) /* Token.NUMBER */ :
		return j2gort.ValueOfDouble(n.GetDouble())
	case int32(34) /* Token.BIGINT */ :
		return nil
	case int32(87) /* Token.VOID */ :
		if !(numberConversions) {
			return nil
		}
		return j2gort.ValueOfDouble(math.NaN())
	case int32(32) /* Token.NAME */ :
		switch j2gort.CheckNotNullString(n.GetString()) {
		case "undefined":
			{

				if numberConversions {
					return j2gort.ValueOfDouble(math.NaN())
				}
				return nil
			}

		case "NaN":
			return j2gort.ValueOfDouble(math.NaN())
		case "Infinity":
			return j2gort.ValueOfDouble(math.Inf(1))
		default:
			return nil
		}
		panic("unreachable switch expression")
	case int32(21) /* Token.POS */ :
		return NodeUtil_DoGetNumberValue(n.GetOnlyChild(), numberConversions)
	case int32(22) /* Token.NEG */ :
		var val *j2gort.BoxDouble = NodeUtil_DoGetNumberValue(n.GetOnlyChild(), numberConversions)
		_ = val
		if val == nil {
			return nil
		} else {
			return j2gort.ValueOfDouble(-(j2gort.UnboxDouble(val)))
		}
	case int32(20) /* Token.BITNOT */ :
		var val_2 *j2gort.BoxDouble = NodeUtil_DoGetNumberValue(n.GetOnlyChild(), numberConversions)
		_ = val_2
		if val_2 == nil {
			return nil
		} else {
			return j2gort.ValueOfDouble(float64(^(pkg_base.JSCompDoubles_EcmascriptToInt32(j2gort.UnboxDouble(val_2)))))
		}
	case int32(38) /* Token.FALSE */, int32(19) /* Token.NOT */, int32(36) /* Token.NULL */, int32(39) /* Token.TRUE */ :
		if !(numberConversions) {
			return nil
		}
		switch (NodeUtil_GetBooleanValue(n)).Ordinal() {
		case int32(2) /* Tri.TRUE */ :
			return j2gort.ValueOfDouble(float64(1.0))
		case int32(0) /* Tri.FALSE */ :
			return j2gort.ValueOfDouble(float64(0.0))
		case int32(1) /* Tri.UNKNOWN */ :
			return nil
		}
		panic("unreachable switch expression")
	case int32(131) /* Token.TEMPLATELIT */ :
		if !(numberConversions) {
			return nil
		}
		var string_ string = NodeUtil_GetStringValue(n)
		_ = string_
		if string_ == j2gort.NullString {
			return nil
		}
		return NodeUtil_GetStringNumberValue(string_)
	case int32(35) /* Token.STRINGLIT */ :
		if !(numberConversions) {
			return nil
		}
		return NodeUtil_GetStringNumberValue(n.GetString())
	case int32(46) /* Token.ARRAYLIT */, int32(47) /* Token.OBJECTLIT */ :
		if !(numberConversions) {
			return nil
		}
		var value string = NodeUtil_GetStringValue(n)
		_ = value
		if value != j2gort.NullString {
			return NodeUtil_GetStringNumberValue(value)
		} else {
			return nil
		}
	default:
	}
	return nil
}

// Java static method: NodeUtil.getStringNumberValue(java.lang.String rawJsString)
func NodeUtil_GetStringNumberValue(rawJsString string) *j2gort.BoxDouble {
	_ = rawJsString
	if j2gort.JStrIndexOf(rawJsString, "\x0b") >= 0 {
		return nil
	}
	var s string = NodeUtil_TrimJsWhiteSpace(rawJsString)
	_ = s
	if j2gort.CheckNotNullString(s) == "" {
		return j2gort.ValueOfDouble(float64(0.0))
	}
	if (j2gort.JStrLen(s) >= int32(2)) && (int32(j2gort.JStrCharAt(s, int32(0))) == int32(uint16(48))) {
		var radix int32 = int32(0)
		_ = radix
		var prefix uint16 = j2gort.JStrCharAt(s, int32(1))
		_ = prefix
		if (int32(prefix) == int32(uint16(120))) || (int32(prefix) == int32(uint16(88))) {
			radix = int32(16)
		} else if (int32(prefix) == int32(uint16(98))) || (int32(prefix) == int32(uint16(66))) {
			radix = int32(2)
		} else if (int32(prefix) == int32(uint16(111))) || (int32(prefix) == int32(uint16(79))) {
			radix = int32(8)
		}
		if radix != int32(0) {
			var digits string = j2gort.JStrSubFrom(s, int32(2))
			_ = digits
			if j2gort.CheckNotNullString(digits) == "" {
				return j2gort.ValueOfDouble(math.NaN())
			}
			{
				var i int32 = int32(0)
				_ = i
				for ; i < j2gort.JStrLen(digits); i = i + 1 {
					if !(NodeUtil_IsValidDigit(j2gort.JStrCharAt(digits, i), radix)) {
						return j2gort.ValueOfDouble(math.NaN())
					}
				}
			}
			var didRet_5 bool
			var retVal_6 *j2gort.BoxDouble
			func() {
				defer func() {
					if rec_7 := recover(); rec_7 != nil {
						jex_8 := javalang.RecoverException(rec_7)
						_ = jex_8
						if j2gort.InstanceOf[javalang.INumberFormatException](jex_8) {
							e := j2gort.AsIface[javalang.INumberFormatException](jex_8)
							_ = e
							retVal_6 = j2gort.ValueOfDouble(math.NaN())
							didRet_5 = true
							return
						} else {
							panic(rec_7)
						}
					}
				}()
				if j2gort.JStrLen(digits) <= int32(14) {
					retVal_6 = j2gort.ValueOfDouble(float64(j2gort.ParseLong(digits, radix)))
					didRet_5 = true
					return
				} else {
					retVal_6 = j2gort.ValueOfDouble(j2gort.NewExt_java_math_BigInteger_2_string_int32(digits, radix).DoubleValue())
					didRet_5 = true
					return
				}
			}()
			if didRet_5 {
				return retVal_6
			}
		}
	}
	if (((j2gort.JStrLen(s) > int32(3)) && ((int32(j2gort.JStrCharAt(s, int32(0))) == int32(uint16(45))) || (int32(j2gort.JStrCharAt(s, int32(0))) == int32(uint16(43))))) && (int32(j2gort.JStrCharAt(s, int32(1))) == int32(uint16(48)))) && ((int32(j2gort.JStrCharAt(s, int32(2))) == int32(uint16(120))) || (int32(j2gort.JStrCharAt(s, int32(2))) == int32(uint16(88)))) {
		return nil
	}
	if ((j2gort.CheckNotNullString(s) == "infinity") || (j2gort.CheckNotNullString(s) == "-infinity")) || (j2gort.CheckNotNullString(s) == "+infinity") {
		return nil
	}
	if (j2gort.CheckNotNullString(s) == "Infinity") || (j2gort.CheckNotNullString(s) == "+Infinity") {
		return j2gort.ValueOfDouble(math.Inf(1))
	}
	if j2gort.CheckNotNullString(s) == "-Infinity" {
		return j2gort.ValueOfDouble(math.Inf(-1))
	}
	{
		var i_2 int32 = int32(0)
		_ = i_2
		for ; i_2 < j2gort.JStrLen(s); i_2 = i_2 + 1 {
			var c uint16 = j2gort.JStrCharAt(s, i_2)
			_ = c
			if ((((((int32(uint16(48)) <= int32(c)) && (int32(c) <= int32(uint16(57)))) || (int32(c) == int32(uint16(46)))) || (int32(c) == int32(uint16(101)))) || (int32(c) == int32(uint16(69)))) || (int32(c) == int32(uint16(43)))) || (int32(c) == int32(uint16(45))) {
				continue
			}
			return j2gort.ValueOfDouble(math.NaN())
		}
	}
	var didRet_9 bool
	var retVal_10 *j2gort.BoxDouble
	func() {
		defer func() {
			if rec_11 := recover(); rec_11 != nil {
				jex_12 := javalang.RecoverException(rec_11)
				_ = jex_12
				if j2gort.InstanceOf[javalang.INumberFormatException](jex_12) {
					e_2 := j2gort.AsIface[javalang.INumberFormatException](jex_12)
					_ = e_2
					retVal_10 = j2gort.ValueOfDouble(math.NaN())
					didRet_9 = true
					return
				} else {
					panic(rec_11)
				}
			}
		}()
		retVal_10 = j2gort.ValueOfDouble(j2gort.ParseDouble(s))
		didRet_9 = true
		return
	}()
	if didRet_9 {
		return retVal_10
	}
	panic("unreachable")
}

// Java static method: NodeUtil.isValidDigit(char c, int radix)
func NodeUtil_IsValidDigit(c uint16, radix int32) bool {
	_ = c
	_ = radix
	switch radix {
	case int32(2):
		return ((int32(c) == int32(uint16(48))) || (int32(c) == int32(uint16(49))))
	case int32(8):
		return ((int32(uint16(48)) <= int32(c)) && (int32(c) <= int32(uint16(55))))
	case int32(16):
		return ((((int32(uint16(48)) <= int32(c)) && (int32(c) <= int32(uint16(57)))) || ((int32(uint16(97)) <= int32(c)) && (int32(c) <= int32(uint16(102))))) || ((int32(uint16(65)) <= int32(c)) && (int32(c) <= int32(uint16(70)))))
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java static method: NodeUtil.getBigIntValue(com.google.javascript.rhino.Node n)
func NodeUtil_GetBigIntValue(n pkg_rhino.Node) *j2gort.Ext_java_math_BigInteger {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(33) /* Token.NUMBER */ :
		var val float64 = n.GetDouble()
		_ = val
		{

			if pkg_base.JSCompDoubles_IsAtLeastIntegerPrecision(val) && pkg_base.JSCompDoubles_IsExactInt64(val) {
				return j2gort.Ext_java_math_BigInteger_ValueOf(javalang.D2L(val))
			}
			return nil
		}

	case int32(34) /* Token.BIGINT */ :
		return n.GetBigInt()
	case int32(38) /* Token.FALSE */, int32(19) /* Token.NOT */, int32(39) /* Token.TRUE */ :
		{

			switch (NodeUtil_GetBooleanValue(n)).Ordinal() {
			case int32(2) /* Tri.TRUE */ :
				return j2gort.Ext_java_math_BigInteger_F_ONE
			case int32(0) /* Tri.FALSE */ :
				return j2gort.Ext_java_math_BigInteger_F_ZERO
			case int32(1) /* Tri.UNKNOWN */ :
				return nil
			}
			panic("unreachable switch expression")
		}

	case int32(131) /* Token.TEMPLATELIT */ :
		var string_ string = NodeUtil_GetStringValue(n)
		_ = string_
		if string_ == j2gort.NullString {
			return nil
		}
		return NodeUtil_GetStringBigIntValue(string_)
	case int32(35) /* Token.STRINGLIT */ :
		return NodeUtil_GetStringBigIntValue(n.GetString())
	case int32(22) /* Token.NEG */ :
		var result *j2gort.Ext_java_math_BigInteger = NodeUtil_GetBigIntValue(n.GetOnlyChild())
		_ = result
		{

			if result == nil {
				return nil
			}
			return result.Negate()
		}

	case int32(20) /* Token.BITNOT */ :
		var result_2 *j2gort.Ext_java_math_BigInteger = NodeUtil_GetBigIntValue(n.GetOnlyChild())
		_ = result_2
		{

			if result_2 == nil {
				return nil
			}
			return result_2.Not()
		}

	case int32(46) /* Token.ARRAYLIT */, int32(47) /* Token.OBJECTLIT */ :
		var value string = NodeUtil_GetStringValue(n)
		_ = value
		{

			if value != j2gort.NullString {
				return NodeUtil_GetStringBigIntValue(value)
			}
			return nil
		}

	case int32(87) /* Token.VOID */, int32(32) /* Token.NAME */, int32(36) /* Token.NULL */ :
		return nil
	default:
		return nil
	}
	panic("unreachable switch expression")
}

// Java static method: NodeUtil.getStringBigIntValue(java.lang.String rawJsString)
func NodeUtil_GetStringBigIntValue(rawJsString string) *j2gort.Ext_java_math_BigInteger {
	_ = rawJsString
	if j2gort.JStrIndexOf(rawJsString, "\x0b") >= 0 {
		return nil
	}
	var s string = NodeUtil_TrimJsWhiteSpace(rawJsString)
	_ = s
	if j2gort.CheckNotNullString(s) == "" {
		return j2gort.Ext_java_math_BigInteger_F_ZERO
	}
	if (j2gort.JStrLen(s) > int32(2)) && (int32(j2gort.JStrCharAt(s, int32(0))) == int32(uint16(48))) {
		var radix int32 = int32(0)
		_ = radix
		switch j2gort.JStrCharAt(s, int32(1)) {
		case uint16(120), uint16(88):
			radix = int32(16)
		case uint16(111), uint16(79):
			radix = int32(8)
		case uint16(98), uint16(66):
			radix = int32(2)
		default:
			radix = int32(0)
		}
		if radix != int32(0) {
			var didRet_13 bool
			var retVal_14 *j2gort.Ext_java_math_BigInteger
			func() {
				defer func() {
					if rec_15 := recover(); rec_15 != nil {
						jex_16 := javalang.RecoverException(rec_15)
						_ = jex_16
						if j2gort.InstanceOf[javalang.INumberFormatException](jex_16) {
							e := j2gort.AsIface[javalang.INumberFormatException](jex_16)
							_ = e
							retVal_14 = nil
							didRet_13 = true
							return
						} else {
							panic(rec_15)
						}
					}
				}()
				retVal_14 = j2gort.NewExt_java_math_BigInteger_2_string_int32(j2gort.JStrSubFrom(s, int32(2)), radix)
				didRet_13 = true
				return
			}()
			if didRet_13 {
				return retVal_14
			}
		}
	}
	var didRet_17 bool
	var retVal_18 *j2gort.Ext_java_math_BigInteger
	func() {
		defer func() {
			if rec_19 := recover(); rec_19 != nil {
				jex_20 := javalang.RecoverException(rec_19)
				_ = jex_20
				if j2gort.InstanceOf[javalang.INumberFormatException](jex_20) {
					e_2 := j2gort.AsIface[javalang.INumberFormatException](jex_20)
					_ = e_2
					retVal_18 = nil
					didRet_17 = true
					return
				} else {
					panic(rec_19)
				}
			}
		}()
		retVal_18 = j2gort.NewExt_java_math_BigInteger_1_string(s)
		didRet_17 = true
		return
	}()
	if didRet_17 {
		return retVal_18
	}
	panic("unreachable")
}

// Java static method: NodeUtil.trimJsWhiteSpace(java.lang.String s)
func NodeUtil_TrimJsWhiteSpace(s string) string {
	_ = s
	var start int32 = int32(0)
	_ = start
	var end int32 = j2gort.JStrLen(s)
	_ = end
	for (end > int32(0)) && (pkg_rhino.TokenUtil_IsStrWhiteSpaceChar(int32(j2gort.JStrCharAt(s, (end-int32(1))))) == pkg_base.StaticGet_Tri_TRUE()) {
		end = end - 1
	}
	for (start < end) && (pkg_rhino.TokenUtil_IsStrWhiteSpaceChar(int32(j2gort.JStrCharAt(s, start))) == pkg_base.StaticGet_Tri_TRUE()) {
		start = start + 1
	}
	return j2gort.JStrSub(s, start, end)
}

// Java static method: NodeUtil.getName(com.google.javascript.rhino.Node n)
func NodeUtil_GetName(n pkg_rhino.Node) string {
	_ = n
	var nameNode pkg_rhino.Node = NodeUtil_GetNameNode(n)
	_ = nameNode
	if nameNode == 0 {
		return j2gort.NullString
	} else {
		return nameNode.GetQualifiedName()
	}
}

// Java static method: NodeUtil.getNameNode(com.google.javascript.rhino.Node n)
func NodeUtil_GetNameNode(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if !(n.IsFunction() || n.IsClass()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(n))
	}
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	switch parent.TokenSwitchOrd_() {
	case int32(32) /* Token.NAME */ :
		return parent
	case int32(51) /* Token.ASSIGN */ :
		var firstChild pkg_rhino.Node = parent.GetFirstChild()
		_ = firstChild
		{

			if firstChild.IsQualifiedName() {
				return firstChild
			}
			return pkg_rhino.Node(0)
		}

	default:
		var funNameNode pkg_rhino.Node = n.GetFirstChild()
		_ = funNameNode
		{

			if funNameNode.IsEmpty() || (j2gort.CheckNotNullString(funNameNode.GetString()) == "") {
				return pkg_rhino.Node(0)
			}
			return funNameNode
		}

	}
	panic("unreachable switch expression")
}

// Java static method: NodeUtil.removeName(com.google.javascript.rhino.Node n)
func NodeUtil_RemoveName(n pkg_rhino.Node) {
	_ = n
	if !(n.IsFunction() || n.IsClass()) {
		j2gort.CheckState(false)
	}
	var originalName pkg_rhino.Node = n.GetFirstChild()
	_ = originalName
	var emptyName pkg_rhino.Node
	_ = emptyName
	if n.IsFunction() {
		emptyName = pkg_rhino.IR_Name("")
	} else {
		emptyName = pkg_rhino.IR_Empty()
	}
	originalName.ReplaceWith(emptyName.Srcref(originalName))
}

// Java static method: NodeUtil.getNearestFunctionName(com.google.javascript.rhino.Node n)
func NodeUtil_GetNearestFunctionName(n pkg_rhino.Node) string {
	_ = n
	if !(n.IsFunction()) {
		return j2gort.NullString
	}
	var name string = NodeUtil_GetName(n)
	_ = name
	if name != j2gort.NullString {
		return name
	}
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	switch parent.TokenSwitchOrd_() {
	case int32(107) /* Token.MEMBER_FUNCTION_DEF */, int32(96) /* Token.SETTER_DEF */, int32(95) /* Token.GETTER_DEF */, int32(100) /* Token.STRING_KEY */ :
		return parent.GetString()
	case int32(33) /* Token.NUMBER */ :
		return NodeUtil_GetStringValue(parent)
	default:
		return j2gort.NullString
	}
	panic("unreachable switch expression")
}

// Java static method: NodeUtil.getClassMembers(com.google.javascript.rhino.Node n)
func NodeUtil_GetClassMembers(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if !(n.IsClass()) {
		j2gort.CheckArgument(false)
	}
	return n.GetLastChild()
}

// Java static method: NodeUtil.getEs6ClassConstructorMemberFunctionDef(com.google.javascript.rhino.Node classNode)
func NodeUtil_GetEs6ClassConstructorMemberFunctionDef(classNode pkg_rhino.Node) pkg_rhino.Node {
	_ = classNode
	if !(classNode.IsClass()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(classNode))
	}
	var classMembers pkg_rhino.Node = j2gort.CheckNotNullNode(classNode.GetLastChild())
	_ = classMembers
	{
		var memberFunctionDef pkg_rhino.Node = classMembers.GetFirstChild()
		_ = memberFunctionDef
		for ; memberFunctionDef != 0; memberFunctionDef = memberFunctionDef.GetNext() {
			if NodeUtil_IsEs6ConstructorMemberFunctionDef(memberFunctionDef) {
				return memberFunctionDef
			}
		}
	}
	return pkg_rhino.Node(0)
}

// Java static method: NodeUtil.hasInstanceElements(com.google.javascript.rhino.Node classNode)
func NodeUtil_HasInstanceElements(classNode pkg_rhino.Node) bool {
	_ = classNode
	{
		var member pkg_rhino.Node = NodeUtil_GetClassMembers(classNode).GetFirstChild()
		_ = member
		for ; member != 0; member = member.GetNext() {
			if !(member.IsStaticMember()) && ((member.IsMemberFieldDef() || member.IsComputedFieldDef()) || member.IsPrivateIdentifier()) {
				return true
			}
		}
	}
	return false
}

// Java static method: NodeUtil.isImmutableValue(com.google.javascript.rhino.Node n)
func NodeUtil_IsImmutableValue(n pkg_rhino.Node) bool {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(35) /* Token.STRINGLIT */, int32(33) /* Token.NUMBER */, int32(34) /* Token.BIGINT */, int32(36) /* Token.NULL */, int32(39) /* Token.TRUE */, int32(38) /* Token.FALSE */ :
		return true
	case int32(101) /* Token.CAST */, int32(19) /* Token.NOT */, int32(87) /* Token.VOID */, int32(22) /* Token.NEG */ :
		return NodeUtil_IsImmutableValue(n.GetFirstChild())
	case int32(32) /* Token.NAME */ :
		var name string = n.GetString()
		_ = name
		return ((("undefined" == name) || ("Infinity" == name)) || ("NaN" == name))
	case int32(131) /* Token.TEMPLATELIT */ :
		{
			var child pkg_rhino.Node = n.GetFirstChild()
			_ = child
			for ; child != 0; child = child.GetNext() {
				if child.IsTemplateLitSub() {
					if !(NodeUtil_IsImmutableValue(child.GetFirstChild())) {
						return false
					}
				}
			}
		}
		return true
	case int32(133) /* Token.TEMPLATELIT_STRING */ :
		panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Invalid argument "+j2gort.StrVal(n))))
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java static method: NodeUtil.isNumericLiteral(com.google.javascript.rhino.Node n)
func NodeUtil_IsNumericLiteral(n pkg_rhino.Node) bool {
	_ = n
	if n == 0 {
		return false
	}
	return (n.IsNumber() || (((n.IsNeg() || n.IsPos()) && n.HasOneChild()) && n.GetFirstChild().IsNumber()))
}

// Java static method: NodeUtil.getNumericLiteralValue(com.google.javascript.rhino.Node n)
func NodeUtil_GetNumericLiteralValue(n pkg_rhino.Node) float64 {
	_ = n
	if !(NodeUtil_IsNumericLiteral(n)) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(n))
	}
	switch n.TokenSwitchOrd_() {
	case int32(33) /* Token.NUMBER */ :
		return n.GetDouble()
	case int32(21) /* Token.POS */ :
		return n.GetFirstChild().GetDouble()
	case int32(22) /* Token.NEG */ :
		return -(n.GetFirstChild().GetDouble())
	default:
		panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError", "Unreachable: "+j2gort.StrVal(n))))
	}
	panic("unreachable switch expression")
}

// Java static method: NodeUtil.isSymmetricOperation(com.google.javascript.rhino.Node n)
func NodeUtil_IsSymmetricOperation(n pkg_rhino.Node) bool {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(4) /* Token.EQ */, int32(5) /* Token.NE */, int32(40) /* Token.SHEQ */, int32(41) /* Token.SHNE */, int32(15) /* Token.MUL */ :
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java static method: NodeUtil.isRelationalOperation(com.google.javascript.rhino.Node n)
func NodeUtil_IsRelationalOperation(n pkg_rhino.Node) bool {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(8) /* Token.GT */, int32(9) /* Token.GE */, int32(6) /* Token.LT */, int32(7) /* Token.LE */ :
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java static method: NodeUtil.getInverseOperator(com.google.javascript.rhino.Token type)
func NodeUtil_GetInverseOperator(type_ pkg_rhino.Token) pkg_rhino.Token {
	_ = type_
	switch (type_).Ordinal() {
	case int32(8) /* Token.GT */ :
		return pkg_rhino.StaticGet_Token_LT()
	case int32(6) /* Token.LT */ :
		return pkg_rhino.StaticGet_Token_GT()
	case int32(9) /* Token.GE */ :
		return pkg_rhino.StaticGet_Token_LE()
	case int32(7) /* Token.LE */ :
		return pkg_rhino.StaticGet_Token_GE()
	default:
		panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Unexpected token: "+j2gort.StrVal(type_))))
	}
	panic("unreachable switch expression")
}

// Java static method: NodeUtil.isLiteralValue(com.google.javascript.rhino.Node n, boolean includeFunctions)
func NodeUtil_IsLiteralValue(n pkg_rhino.Node, includeFunctions bool) bool {
	_ = n
	_ = includeFunctions
	switch n.TokenSwitchOrd_() {
	case int32(101) /* Token.CAST */ :
		return NodeUtil_IsLiteralValue(n.GetFirstChild(), includeFunctions)
	case int32(46) /* Token.ARRAYLIT */ :
		{
			var child pkg_rhino.Node = n.GetFirstChild()
			_ = child
			for ; child != 0; child = child.GetNext() {
				if (!(child.IsEmpty())) && !(NodeUtil_IsLiteralValue(child, includeFunctions)) {
					return false
				}
			}
		}
		return true
	case int32(42) /* Token.REGEXP */ :
		{
			var child_2 pkg_rhino.Node = n.GetFirstChild()
			_ = child_2
			for ; child_2 != 0; child_2 = child_2.GetNext() {
				if !(NodeUtil_IsLiteralValue(child_2, includeFunctions)) {
					return false
				}
			}
		}
		return true
	case int32(47) /* Token.OBJECTLIT */ :
		{
			var child_3 pkg_rhino.Node = n.GetFirstChild()
			_ = child_3
			for ; child_3 != 0; child_3 = child_3.GetNext() {
				switch child_3.TokenSwitchOrd_() {
				case int32(107) /* Token.MEMBER_FUNCTION_DEF */, int32(95) /* Token.GETTER_DEF */, int32(96) /* Token.SETTER_DEF */ :
					if !(includeFunctions) {
						return false
					}
				case int32(129) /* Token.COMPUTED_PROP */ :
					if !(NodeUtil_IsLiteralValue(child_3.GetFirstChild(), includeFunctions)) || !(NodeUtil_IsLiteralValue(child_3.GetLastChild(), includeFunctions)) {
						return false
					}
				case int32(128) /* Token.OBJECT_SPREAD */ :
					if !(NodeUtil_IsLiteralValue(child_3.GetOnlyChild(), includeFunctions)) {
						return false
					}
				case int32(100) /* Token.STRING_KEY */ :
					if !(NodeUtil_IsLiteralValue(child_3.GetOnlyChild(), includeFunctions)) {
						return false
					}
				default:
					panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Unexpected child of OBJECTLIT: "+j2gort.StrValString(child_3.ToStringTree()))))
				}
			}
		}
		return true
	case int32(73) /* Token.FUNCTION */ :
		return (includeFunctions && !(NodeUtil_IsFunctionDeclaration(n)))
	case int32(131) /* Token.TEMPLATELIT */ :
		{
			var child_4 pkg_rhino.Node = n.GetFirstChild()
			_ = child_4
			for ; child_4 != 0; child_4 = child_4.GetNext() {
				if child_4.IsTemplateLitSub() {
					if !(NodeUtil_IsLiteralValue(child_4.GetFirstChild(), includeFunctions)) {
						return false
					}
				}
			}
		}
		return true
	default:
		return NodeUtil_IsImmutableValue(n)
	}
	panic("unreachable switch expression")
}

// Java static method: NodeUtil.findClosureUnawareScriptRoot(com.google.javascript.rhino.Node script)
func NodeUtil_FindClosureUnawareScriptRoot(script pkg_rhino.Node) pkg_rhino.Node {
	_ = script
	if !(script.IsScript() && script.GetIsInClosureUnawareSubtree()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(script))
	}
	return NodeUtil_FindPreorder(script, lambdaSingleton_node_util_3, lambdaSingleton_node_util_4)
}

// Java static method: NodeUtil.getInsertionPointAfterAllInnerFunctionDeclarations(com.google.javascript.rhino.Node block)
func NodeUtil_GetInsertionPointAfterAllInnerFunctionDeclarations(block pkg_rhino.Node) pkg_rhino.Node {
	_ = block
	if !(block.IsBlock()) {
		j2gort.CheckState(false)
	}
	var current pkg_rhino.Node = block.GetFirstChild()
	_ = current
	for (current != 0) && NodeUtil_IsFunctionDeclaration(current) {
		current = current.GetNext()
	}
	return current
}

// Java static method: NodeUtil.isSomeCompileTimeConstStringValue(com.google.javascript.rhino.Node node)
func NodeUtil_IsSomeCompileTimeConstStringValue(node pkg_rhino.Node) bool {
	_ = node
	if node.IsStringLit() || (node.IsTemplateLit() && node.HasOneChild()) {
		return true
	} else if node.IsAdd() {
		if !(node.HasTwoChildren()) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
		}
		var left pkg_rhino.Node = node.GetFirstChild()
		_ = left
		var right pkg_rhino.Node = node.GetLastChild()
		_ = right
		return (NodeUtil_IsSomeCompileTimeConstStringValue(left) && NodeUtil_IsSomeCompileTimeConstStringValue(right))
	} else if node.IsHook() {
		var left_2 pkg_rhino.Node = node.GetSecondChild()
		_ = left_2
		var right_2 pkg_rhino.Node = node.GetLastChild()
		_ = right_2
		return (NodeUtil_IsSomeCompileTimeConstStringValue(left_2) && NodeUtil_IsSomeCompileTimeConstStringValue(right_2))
	}
	return false
}

// Java static method: NodeUtil.isEmptyBlock(com.google.javascript.rhino.Node block)
func NodeUtil_IsEmptyBlock(block pkg_rhino.Node) bool {
	_ = block
	if !(block.IsBlock()) {
		return false
	}
	{
		var n pkg_rhino.Node = block.GetFirstChild()
		_ = n
		for ; n != 0; n = n.GetNext() {
			if !(n.IsEmpty()) {
				return false
			}
		}
	}
	return true
}

// Java static method: NodeUtil.isBinaryOperator(com.google.javascript.rhino.Node n)
func NodeUtil_IsBinaryOperator(n pkg_rhino.Node) bool {
	_ = n
	return NodeUtil_IsBinaryOperatorType(n.GetToken())
}

// Java static method: NodeUtil.isBinaryOperatorType(com.google.javascript.rhino.Token type)
func NodeUtil_IsBinaryOperatorType(type_ pkg_rhino.Token) bool {
	_ = type_
	switch (type_).Ordinal() {
	case int32(68) /* Token.OR */, int32(69) /* Token.AND */, int32(70) /* Token.COALESCE */, int32(1) /* Token.BITOR */, int32(2) /* Token.BITXOR */, int32(3) /* Token.BITAND */, int32(4) /* Token.EQ */, int32(5) /* Token.NE */, int32(40) /* Token.SHEQ */, int32(41) /* Token.SHNE */, int32(6) /* Token.LT */, int32(8) /* Token.GT */, int32(7) /* Token.LE */, int32(9) /* Token.GE */, int32(45) /* Token.INSTANCEOF */, int32(44) /* Token.IN */, int32(10) /* Token.LSH */, int32(11) /* Token.RSH */, int32(12) /* Token.URSH */, int32(13) /* Token.ADD */, int32(14) /* Token.SUB */, int32(15) /* Token.MUL */, int32(16) /* Token.DIV */, int32(17) /* Token.MOD */, int32(18) /* Token.EXPONENT */ :
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java static method: NodeUtil.isUnaryOperator(com.google.javascript.rhino.Node n)
func NodeUtil_IsUnaryOperator(n pkg_rhino.Node) bool {
	_ = n
	return NodeUtil_IsUnaryOperatorType(n.GetToken())
}

// Java static method: NodeUtil.isUnaryOperatorType(com.google.javascript.rhino.Token type)
func NodeUtil_IsUnaryOperatorType(type_ pkg_rhino.Token) bool {
	_ = type_
	switch (type_).Ordinal() {
	case int32(24) /* Token.DELPROP */, int32(87) /* Token.VOID */, int32(25) /* Token.TYPEOF */, int32(21) /* Token.POS */, int32(22) /* Token.NEG */, int32(20) /* Token.BITNOT */, int32(19) /* Token.NOT */ :
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java static method: NodeUtil.isUpdateOperator(com.google.javascript.rhino.Node n)
func NodeUtil_IsUpdateOperator(n pkg_rhino.Node) bool {
	_ = n
	return NodeUtil_IsUpdateOperatorType(n.GetToken())
}

// Java static method: NodeUtil.isUpdateOperatorType(com.google.javascript.rhino.Token type)
func NodeUtil_IsUpdateOperatorType(type_ pkg_rhino.Token) bool {
	_ = type_
	switch (type_).Ordinal() {
	case int32(71) /* Token.INC */, int32(72) /* Token.DEC */ :
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java static method: NodeUtil.isSimpleOperator(com.google.javascript.rhino.Node n)
func NodeUtil_IsSimpleOperator(n pkg_rhino.Node) bool {
	_ = n
	return NodeUtil_IsSimpleOperatorType(n.GetToken())
}

// Java static method: NodeUtil.isSimpleOperatorType(com.google.javascript.rhino.Token type)
func NodeUtil_IsSimpleOperatorType(type_ pkg_rhino.Token) bool {
	_ = type_
	switch (type_).Ordinal() {
	case int32(13) /* Token.ADD */, int32(3) /* Token.BITAND */, int32(20) /* Token.BITNOT */, int32(1) /* Token.BITOR */, int32(2) /* Token.BITXOR */, int32(50) /* Token.COMMA */, int32(16) /* Token.DIV */, int32(4) /* Token.EQ */, int32(18) /* Token.EXPONENT */, int32(9) /* Token.GE */, int32(8) /* Token.GT */, int32(44) /* Token.IN */, int32(45) /* Token.INSTANCEOF */, int32(7) /* Token.LE */, int32(10) /* Token.LSH */, int32(6) /* Token.LT */, int32(17) /* Token.MOD */, int32(15) /* Token.MUL */, int32(5) /* Token.NE */, int32(19) /* Token.NOT */, int32(11) /* Token.RSH */, int32(40) /* Token.SHEQ */, int32(41) /* Token.SHNE */, int32(14) /* Token.SUB */, int32(25) /* Token.TYPEOF */, int32(87) /* Token.VOID */, int32(21) /* Token.POS */, int32(22) /* Token.NEG */, int32(12) /* Token.URSH */ :
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java static method: NodeUtil.isNamespaceDecl(com.google.javascript.rhino.Node n)
func NodeUtil_IsNamespaceDecl(n pkg_rhino.Node) bool {
	_ = n
	var jsdoc *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(n)
	_ = jsdoc
	if ((jsdoc != nil) && jsdoc.HasTypeNodes()) && !(jsdoc.HasTypedefType()) {
		return false
	}
	var isMarkedConst bool = (n.GetParent().IsConst() || ((jsdoc != nil) && jsdoc.IsConstant()))
	_ = isMarkedConst
	if !(n.IsFromExterns()) && !(isMarkedConst) {
		return false
	}
	var qnameNode pkg_rhino.Node = pkg_rhino.Node(0)
	_ = qnameNode
	var initializer pkg_rhino.Node = pkg_rhino.Node(0)
	_ = initializer
	if NodeUtil_IsNameDeclaration(n.GetParent()) {
		qnameNode = n
		initializer = n.GetFirstChild()
	} else if n.IsExprResult() {
		var expr pkg_rhino.Node = n.GetFirstChild()
		_ = expr
		if !(expr.IsAssign()) || !(expr.GetFirstChild().IsGetProp()) {
			return false
		}
		qnameNode = expr.GetFirstChild()
		initializer = expr.GetLastChild()
	} else if n.IsGetProp() {
		var parent pkg_rhino.Node = n.GetParent()
		_ = parent
		if !(parent.IsAssign()) || !(parent.GetParent().IsExprResult()) {
			return false
		}
		qnameNode = n
		initializer = parent.GetLastChild()
	} else {
		return false
	}
	if (initializer == 0) || (qnameNode == 0) {
		return false
	}
	if initializer.IsObjectLit() {
		return true
	}
	return ((initializer.IsOr() && qnameNode.MatchesQualifiedName_1_Node(initializer.GetFirstChild())) && initializer.GetLastChild().IsObjectLit())
}

// Java static method: NodeUtil.isFromTypeSummary(com.google.javascript.rhino.Node n)
func NodeUtil_IsFromTypeSummary(n pkg_rhino.Node) bool {
	_ = n
	if !(n.IsScript()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(n))
	}
	var info *pkg_rhino.JSDocInfo = n.GetJSDocInfo()
	_ = info
	return ((info != nil) && info.IsTypeSummary())
}

// Java static method: NodeUtil.newExpr(com.google.javascript.rhino.Node child)
func NodeUtil_NewExpr(child pkg_rhino.Node) pkg_rhino.Node {
	_ = child
	return pkg_rhino.IR_ExprResult(child).Srcref(child)
}

// Java static method: NodeUtil.iteratesImpureIterable(com.google.javascript.rhino.Node node)
func NodeUtil_IteratesImpureIterable(node pkg_rhino.Node) bool {
	_ = node
	var parent pkg_rhino.Node = node.GetParent()
	_ = parent
	var iterable pkg_rhino.Node = pkg_rhino.Node(0)
	_ = iterable
	switch node.TokenSwitchOrd_() {
	case int32(127) /* Token.ITER_SPREAD */ :
		iterable = node.GetOnlyChild()
	case int32(114) /* Token.YIELD */ :
		if !(node.IsYieldAll()) {
			return false
		}
		iterable = node.GetOnlyChild()
	case int32(112) /* Token.FOR_OF */, int32(113) /* Token.FOR_AWAIT_OF */ :
		iterable = node.GetSecondChild()
	case int32(125) /* Token.ITER_REST */ :
		switch parent.TokenSwitchOrd_() {
		case int32(49) /* Token.PARAM_LIST */ :
			return false
		case int32(102) /* Token.ARRAY_PATTERN */ :
			return true
		default:
			panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unexpected parent of ITRE_REST: "+j2gort.StrValString(parent.ToStringTree()))))
		}
		panic("unreachable switch expression")
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Expected a kind of node that may trigger iteration: "+j2gort.StrValString(node.ToStringTree()))))
	}
	return !(NodeUtil_IsPureIterable(iterable))
}

// Java static method: NodeUtil.isPureIterable(com.google.javascript.rhino.Node node)
func NodeUtil_IsPureIterable(node pkg_rhino.Node) bool {
	_ = node
	switch node.TokenSwitchOrd_() {
	case int32(46) /* Token.ARRAYLIT */, int32(35) /* Token.STRINGLIT */, int32(131) /* Token.TEMPLATELIT */ :
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java static method: NodeUtil.newHasLocalResult(com.google.javascript.rhino.Node n)
func NodeUtil_NewHasLocalResult(n pkg_rhino.Node) bool {
	_ = n
	if !(n.IsNew()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(n))
	}
	return n.IsOnlyModifiesThisCall()
}

// Java static method: NodeUtil.allArgsUnescapedLocal(com.google.javascript.rhino.Node callOrNew)
func NodeUtil_AllArgsUnescapedLocal(callOrNew pkg_rhino.Node) bool {
	_ = callOrNew
	{
		var arg pkg_rhino.Node = callOrNew.GetSecondChild()
		_ = arg
		for ; arg != 0; arg = arg.GetNext() {
			if !(NodeUtil_EvaluatesToLocalValue(arg)) {
				return false
			}
		}
	}
	return true
}

// Java static method: NodeUtil.canBeSideEffected(com.google.javascript.rhino.Node n)
func NodeUtil_CanBeSideEffected_1(n pkg_rhino.Node) bool {
	_ = n
	return NodeUtil_CanBeSideEffected_3(n, j2gort.AsIface[j2gort.Ext_java_util_Set](StaticGet_NodeUtil_KNOWN_CONSTANTS()), nil)
}

// Java static method: NodeUtil.canBeSideEffected(com.google.javascript.rhino.Node n, java.util.Set<java.lang.String> knownConstants, com.google.javascript.jscomp.@org.jspecify.annotations.Nullable Scope scope)
func NodeUtil_CanBeSideEffected_3(n pkg_rhino.Node, knownConstants j2gort.Ext_java_util_Set, scope *Scope) bool {
	_ = n
	_ = knownConstants
	_ = scope
	switch n.TokenSwitchOrd_() {
	case int32(114) /* Token.YIELD */, int32(28) /* Token.CALL */, int32(31) /* Token.OPTCHAIN_CALL */, int32(23) /* Token.NEW */, int32(130) /* Token.TAGGED_TEMPLATELIT */, int32(115) /* Token.AWAIT */, int32(124) /* Token.DYNAMIC_IMPORT */ :
		return true
	case int32(32) /* Token.NAME */ :
		return (!(NodeUtil_IsConstantVar(n, scope)) && !(j2gort.SetContainsStr(knownConstants, n.GetString())))
	case int32(26) /* Token.GETPROP */ :
		return !(n.GetFirstChild().MatchesName_1_string("Symbol"))
	case int32(27) /* Token.GETELEM */, int32(29) /* Token.OPTCHAIN_GETPROP */, int32(30) /* Token.OPTCHAIN_GETELEM */ :
		return true
	case int32(73) /* Token.FUNCTION */ :
		if !(!(NodeUtil_IsFunctionDeclaration(n))) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(n))
		}
		return false
	default:
	}
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			if NodeUtil_CanBeSideEffected_3(c, knownConstants, scope) {
				return true
			}
		}
	}
	return false
}

// Java static method: NodeUtil.precedence(com.google.javascript.rhino.Token type)
func NodeUtil_Precedence(type_ pkg_rhino.Token) int32 {
	_ = type_
	switch (type_).Ordinal() {
	case int32(50) /* Token.COMMA */ :
		return int32(0)
	case int32(52) /* Token.ASSIGN_BITOR */, int32(53) /* Token.ASSIGN_BITXOR */, int32(54) /* Token.ASSIGN_BITAND */, int32(55) /* Token.ASSIGN_LSH */, int32(56) /* Token.ASSIGN_RSH */, int32(57) /* Token.ASSIGN_URSH */, int32(58) /* Token.ASSIGN_ADD */, int32(59) /* Token.ASSIGN_SUB */, int32(60) /* Token.ASSIGN_MUL */, int32(63) /* Token.ASSIGN_EXPONENT */, int32(61) /* Token.ASSIGN_DIV */, int32(62) /* Token.ASSIGN_MOD */, int32(64) /* Token.ASSIGN_OR */, int32(65) /* Token.ASSIGN_AND */, int32(66) /* Token.ASSIGN_COALESCE */, int32(51) /* Token.ASSIGN */ :
		return int32(1)
	case int32(114) /* Token.YIELD */ :
		return int32(2)
	case int32(67) /* Token.HOOK */ :
		return int32(3)
	case int32(68) /* Token.OR */ :
		return int32(4)
	case int32(69) /* Token.AND */ :
		return int32(5)
	case int32(70) /* Token.COALESCE */ :
		return int32(6)
	case int32(1) /* Token.BITOR */ :
		return int32(7)
	case int32(2) /* Token.BITXOR */ :
		return int32(8)
	case int32(3) /* Token.BITAND */ :
		return int32(9)
	case int32(4) /* Token.EQ */, int32(5) /* Token.NE */, int32(40) /* Token.SHEQ */, int32(41) /* Token.SHNE */ :
		return int32(10)
	case int32(6) /* Token.LT */, int32(8) /* Token.GT */, int32(7) /* Token.LE */, int32(9) /* Token.GE */, int32(45) /* Token.INSTANCEOF */, int32(44) /* Token.IN */ :
		return int32(11)
	case int32(10) /* Token.LSH */, int32(11) /* Token.RSH */, int32(12) /* Token.URSH */ :
		return int32(12)
	case int32(14) /* Token.SUB */, int32(13) /* Token.ADD */ :
		return int32(13)
	case int32(15) /* Token.MUL */, int32(17) /* Token.MOD */, int32(16) /* Token.DIV */ :
		return int32(14)
	case int32(18) /* Token.EXPONENT */ :
		return int32(15)
	case int32(115) /* Token.AWAIT */, int32(23) /* Token.NEW */, int32(24) /* Token.DELPROP */, int32(25) /* Token.TYPEOF */, int32(87) /* Token.VOID */, int32(19) /* Token.NOT */, int32(20) /* Token.BITNOT */, int32(21) /* Token.POS */, int32(22) /* Token.NEG */ :
		return int32(16)
	case int32(71) /* Token.INC */, int32(72) /* Token.DEC */ :
		return int32(17)
	case int32(28) /* Token.CALL */, int32(27) /* Token.GETELEM */, int32(26) /* Token.GETPROP */, int32(31) /* Token.OPTCHAIN_CALL */, int32(30) /* Token.OPTCHAIN_GETELEM */, int32(29) /* Token.OPTCHAIN_GETPROP */, int32(135) /* Token.NEW_TARGET */, int32(136) /* Token.IMPORT_META */, int32(46) /* Token.ARRAYLIT */, int32(102) /* Token.ARRAY_PATTERN */, int32(134) /* Token.DEFAULT_VALUE */, int32(104) /* Token.DESTRUCTURING_LHS */, int32(88) /* Token.EMPTY */, int32(38) /* Token.FALSE */, int32(73) /* Token.FUNCTION */, int32(105) /* Token.CLASS */, int32(164) /* Token.INTERFACE */, int32(32) /* Token.NAME */, int32(36) /* Token.NULL */, int32(33) /* Token.NUMBER */, int32(34) /* Token.BIGINT */, int32(47) /* Token.OBJECTLIT */, int32(103) /* Token.OBJECT_PATTERN */, int32(42) /* Token.REGEXP */, int32(125) /* Token.ITER_REST */, int32(126) /* Token.OBJECT_REST */, int32(127) /* Token.ITER_SPREAD */, int32(128) /* Token.OBJECT_SPREAD */, int32(35) /* Token.STRINGLIT */, int32(100) /* Token.STRING_KEY */, int32(172) /* Token.MEMBER_VARIABLE_DEF */, int32(173) /* Token.INDEX_SIGNATURE */, int32(174) /* Token.CALL_SIGNATURE */, int32(37) /* Token.THIS */, int32(110) /* Token.SUPER */, int32(39) /* Token.TRUE */, int32(130) /* Token.TAGGED_TEMPLATELIT */, int32(131) /* Token.TEMPLATELIT */, int32(124) /* Token.DYNAMIC_IMPORT */, int32(142) /* Token.UNION_TYPE */ :
		return int32(18)
	case int32(140) /* Token.FUNCTION_TYPE */ :
		return int32(19)
	case int32(151) /* Token.ARRAY_TYPE */, int32(141) /* Token.PARAMETERIZED_TYPE */ :
		return int32(20)
	case int32(137) /* Token.STRING_TYPE */, int32(139) /* Token.NUMBER_TYPE */, int32(138) /* Token.BOOLEAN_TYPE */, int32(143) /* Token.ANY_TYPE */, int32(149) /* Token.RECORD_TYPE */, int32(144) /* Token.NULLABLE_TYPE */, int32(147) /* Token.NAMED_TYPE */, int32(150) /* Token.UNDEFINED_TYPE */, int32(145) /* Token.VOID_TYPE */, int32(152) /* Token.GENERIC_TYPE */ :
		return int32(21)
	case int32(101) /* Token.CAST */ :
		return int32(22)
	default:
		if !(type_ != pkg_rhino.StaticGet_Token_TEMPLATELIT_STRING()) {
			j2gort.CheckArgument(false)
		}
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unknown precedence for "+j2gort.StrVal(type_))))
	}
	panic("unreachable switch expression")
}

// Java static method: NodeUtil.isUndefined(com.google.javascript.rhino.Node n)
func NodeUtil_IsUndefined(n pkg_rhino.Node) bool {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(87) /* Token.VOID */ :
		return true
	case int32(32) /* Token.NAME */ :
		return (j2gort.CheckNotNullString(n.GetString()) == "undefined")
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java static method: NodeUtil.isNullOrUndefined(com.google.javascript.rhino.Node n)
func NodeUtil_IsNullOrUndefined(n pkg_rhino.Node) bool {
	_ = n
	return (n.IsNull() || NodeUtil_IsUndefined(n))
}

// Java static method: NodeUtil.getKnownValueType(com.google.javascript.rhino.Node n)
func NodeUtil_GetKnownValueType(n pkg_rhino.Node) NodeUtil_ValueType {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(101) /* Token.CAST */ :
		return NodeUtil_GetKnownValueType(n.GetFirstChild())
	case int32(51) /* Token.ASSIGN */, int32(50) /* Token.COMMA */ :
		return NodeUtil_GetKnownValueType(n.GetLastChild())
	case int32(69) /* Token.AND */, int32(68) /* Token.OR */, int32(70) /* Token.COALESCE */, int32(64) /* Token.ASSIGN_OR */, int32(65) /* Token.ASSIGN_AND */, int32(66) /* Token.ASSIGN_COALESCE */ :
		return NodeUtil_And(NodeUtil_GetKnownValueType(n.GetFirstChild()), NodeUtil_GetKnownValueType(n.GetLastChild()))
	case int32(67) /* Token.HOOK */ :
		return NodeUtil_And(NodeUtil_GetKnownValueType(n.GetSecondChild()), NodeUtil_GetKnownValueType(n.GetLastChild()))
	case int32(13) /* Token.ADD */ :
		var last NodeUtil_ValueType = NodeUtil_GetKnownValueType(n.GetLastChild())
		_ = last
		if last == StaticGet_NodeUtil_ValueType_STRING() {
			return StaticGet_NodeUtil_ValueType_STRING()
		}
		var first NodeUtil_ValueType = NodeUtil_GetKnownValueType(n.GetFirstChild())
		_ = first
		if first == StaticGet_NodeUtil_ValueType_STRING() {
			return StaticGet_NodeUtil_ValueType_STRING()
		}
		if (first == StaticGet_NodeUtil_ValueType_OBJECT()) || (last == StaticGet_NodeUtil_ValueType_OBJECT()) {
			return StaticGet_NodeUtil_ValueType_UNDETERMINED()
		}
		if !(NodeUtil_MayBeString_1_NodeUtil_ValueType(first)) && !(NodeUtil_MayBeString_1_NodeUtil_ValueType(last)) {
			if (first == StaticGet_NodeUtil_ValueType_BIGINT()) || (last == StaticGet_NodeUtil_ValueType_BIGINT()) {
				return StaticGet_NodeUtil_ValueType_BIGINT()
			} else {
				return StaticGet_NodeUtil_ValueType_NUMBER()
			}
		}
		return StaticGet_NodeUtil_ValueType_UNDETERMINED()
	case int32(58) /* Token.ASSIGN_ADD */ :
		var last_2 NodeUtil_ValueType = NodeUtil_GetKnownValueType(n.GetLastChild())
		_ = last_2
		if last_2 == StaticGet_NodeUtil_ValueType_STRING() {
			return StaticGet_NodeUtil_ValueType_STRING()
		}
		return StaticGet_NodeUtil_ValueType_UNDETERMINED()
	case int32(32) /* Token.NAME */ :
		var name string = n.GetString()
		_ = name
		if j2gort.CheckNotNullString(name) == "undefined" {
			return StaticGet_NodeUtil_ValueType_VOID()
		}
		if j2gort.CheckNotNullString(name) == "NaN" {
			return StaticGet_NodeUtil_ValueType_NUMBER()
		}
		if j2gort.CheckNotNullString(name) == "Infinity" {
			return StaticGet_NodeUtil_ValueType_NUMBER()
		}
		return StaticGet_NodeUtil_ValueType_UNDETERMINED()
	case int32(52) /* Token.ASSIGN_BITOR */, int32(53) /* Token.ASSIGN_BITXOR */, int32(54) /* Token.ASSIGN_BITAND */, int32(55) /* Token.ASSIGN_LSH */, int32(56) /* Token.ASSIGN_RSH */, int32(57) /* Token.ASSIGN_URSH */, int32(59) /* Token.ASSIGN_SUB */, int32(60) /* Token.ASSIGN_MUL */, int32(63) /* Token.ASSIGN_EXPONENT */, int32(61) /* Token.ASSIGN_DIV */, int32(62) /* Token.ASSIGN_MOD */ :
		if NodeUtil_GetKnownValueType(n.GetLastChild()) == StaticGet_NodeUtil_ValueType_BIGINT() {
			return StaticGet_NodeUtil_ValueType_BIGINT()
		} else {
			return StaticGet_NodeUtil_ValueType_NUMBER()
		}
	case int32(34) /* Token.BIGINT */ :
		return StaticGet_NodeUtil_ValueType_BIGINT()
	case int32(1) /* Token.BITOR */, int32(2) /* Token.BITXOR */, int32(3) /* Token.BITAND */, int32(10) /* Token.LSH */, int32(11) /* Token.RSH */, int32(14) /* Token.SUB */, int32(15) /* Token.MUL */, int32(17) /* Token.MOD */, int32(16) /* Token.DIV */, int32(18) /* Token.EXPONENT */ :
		var first_2 NodeUtil_ValueType = NodeUtil_GetKnownValueType(n.GetFirstChild())
		_ = first_2
		var last_3 NodeUtil_ValueType = NodeUtil_GetKnownValueType(n.GetLastChild())
		_ = last_3
		if (first_2 == StaticGet_NodeUtil_ValueType_BIGINT()) || (last_3 == StaticGet_NodeUtil_ValueType_BIGINT()) {
			return StaticGet_NodeUtil_ValueType_BIGINT()
		} else {
			return StaticGet_NodeUtil_ValueType_NUMBER()
		}
	case int32(20) /* Token.BITNOT */, int32(22) /* Token.NEG */ :
		if NodeUtil_GetKnownValueType(n.GetOnlyChild()) == StaticGet_NodeUtil_ValueType_BIGINT() {
			return StaticGet_NodeUtil_ValueType_BIGINT()
		} else {
			return StaticGet_NodeUtil_ValueType_NUMBER()
		}
	case int32(71) /* Token.INC */, int32(72) /* Token.DEC */ :
		return StaticGet_NodeUtil_ValueType_NUMBER()
	case int32(12) /* Token.URSH */, int32(21) /* Token.POS */, int32(33) /* Token.NUMBER */ :
		return StaticGet_NodeUtil_ValueType_NUMBER()
	case int32(39) /* Token.TRUE */, int32(38) /* Token.FALSE */, int32(4) /* Token.EQ */, int32(5) /* Token.NE */, int32(40) /* Token.SHEQ */, int32(41) /* Token.SHNE */, int32(6) /* Token.LT */, int32(8) /* Token.GT */, int32(7) /* Token.LE */, int32(9) /* Token.GE */, int32(44) /* Token.IN */, int32(45) /* Token.INSTANCEOF */, int32(19) /* Token.NOT */, int32(24) /* Token.DELPROP */ :
		return StaticGet_NodeUtil_ValueType_BOOLEAN()
	case int32(25) /* Token.TYPEOF */, int32(35) /* Token.STRINGLIT */, int32(131) /* Token.TEMPLATELIT */ :
		return StaticGet_NodeUtil_ValueType_STRING()
	case int32(36) /* Token.NULL */ :
		return StaticGet_NodeUtil_ValueType_NULL()
	case int32(87) /* Token.VOID */ :
		return StaticGet_NodeUtil_ValueType_VOID()
	case int32(73) /* Token.FUNCTION */, int32(23) /* Token.NEW */, int32(46) /* Token.ARRAYLIT */, int32(47) /* Token.OBJECTLIT */, int32(42) /* Token.REGEXP */ :
		return StaticGet_NodeUtil_ValueType_OBJECT()
	default:
		if !(!(n.IsTemplateLitString())) {
			j2gort.CheckArgument(false)
		}
		return StaticGet_NodeUtil_ValueType_UNDETERMINED()
	}
	panic("unreachable")
}

// Java static method: NodeUtil.and(com.google.javascript.jscomp.NodeUtil.ValueType a, com.google.javascript.jscomp.NodeUtil.ValueType b)
func NodeUtil_And(a NodeUtil_ValueType, b NodeUtil_ValueType) NodeUtil_ValueType {
	_ = a
	_ = b
	if a == b {
		return a
	} else {
		return StaticGet_NodeUtil_ValueType_UNDETERMINED()
	}
}

// Java static method: NodeUtil.isNumericResult(com.google.javascript.rhino.Node n)
func NodeUtil_IsNumericResult(n pkg_rhino.Node) bool {
	_ = n
	return (NodeUtil_GetKnownValueType(n) == StaticGet_NodeUtil_ValueType_NUMBER())
}

// Java static method: NodeUtil.isBigIntResult(com.google.javascript.rhino.Node n)
func NodeUtil_IsBigIntResult(n pkg_rhino.Node) bool {
	_ = n
	return (NodeUtil_GetKnownValueType(n) == StaticGet_NodeUtil_ValueType_BIGINT())
}

// Java static method: NodeUtil.isBooleanResult(com.google.javascript.rhino.Node n)
func NodeUtil_IsBooleanResult(n pkg_rhino.Node) bool {
	_ = n
	return (NodeUtil_GetKnownValueType(n) == StaticGet_NodeUtil_ValueType_BOOLEAN())
}

// Java static method: NodeUtil.isStringResult(com.google.javascript.rhino.Node n)
func NodeUtil_IsStringResult(n pkg_rhino.Node) bool {
	_ = n
	return (NodeUtil_GetKnownValueType(n) == StaticGet_NodeUtil_ValueType_STRING())
}

// Java static method: NodeUtil.isObjectResult(com.google.javascript.rhino.Node n)
func NodeUtil_IsObjectResult(n pkg_rhino.Node) bool {
	_ = n
	return (NodeUtil_GetKnownValueType(n) == StaticGet_NodeUtil_ValueType_OBJECT())
}

// Java static method: NodeUtil.mayBeString(com.google.javascript.rhino.Node n)
func NodeUtil_MayBeString_1_Node(n pkg_rhino.Node) bool {
	_ = n
	return NodeUtil_MayBeString_2(n, false)
}

// Java static method: NodeUtil.mayBeString(com.google.javascript.rhino.Node n, boolean useType)
func NodeUtil_MayBeString_2(n pkg_rhino.Node, useType bool) bool {
	_ = n
	_ = useType
	if useType {
		var color *pkg_colors.Color = n.GetColor()
		_ = color
		if color != nil {
			if color.Equals(j2gort.PtrToAny(pkg_colors.StaticGet_StandardColors_STRING())) {
				return true
			} else if ((color.Equals(j2gort.PtrToAny(pkg_colors.StaticGet_StandardColors_NUMBER())) || color.Equals(j2gort.PtrToAny(pkg_colors.StaticGet_StandardColors_BIGINT()))) || color.Equals(j2gort.PtrToAny(pkg_colors.StaticGet_StandardColors_BOOLEAN()))) || color.Equals(j2gort.PtrToAny(pkg_colors.StaticGet_StandardColors_NULL_OR_VOID())) {
				return false
			}
		}
		var type_ pkg_rhino.JSType = n.GetJSType()
		_ = type_
		if type_ != nil && !type_.IsJavaNil_() {
			if type_.GetJSTypeBase().IsStringValueType() {
				return true
			} else if (((type_.GetJSTypeBase().IsNumberValueType() || type_.GetJSTypeBase().IsBigIntValueType()) || type_.GetJSTypeBase().IsBooleanValueType()) || type_.GetJSTypeBase().IsNullType()) || type_.GetJSTypeBase().IsVoidType() {
				return false
			}
		}
	}
	return NodeUtil_MayBeString_1_NodeUtil_ValueType(NodeUtil_GetKnownValueType(n))
}

// Java static method: NodeUtil.mayBeString(com.google.javascript.jscomp.NodeUtil.ValueType type)
func NodeUtil_MayBeString_1_NodeUtil_ValueType(type_ NodeUtil_ValueType) bool {
	_ = type_
	switch (type_).Ordinal() {
	case int32(6) /* ValueType.BOOLEAN */, int32(1) /* ValueType.NULL */, int32(3) /* ValueType.NUMBER */, int32(4) /* ValueType.BIGINT */, int32(2) /* ValueType.VOID */ :
		return false
	case int32(7) /* ValueType.OBJECT */, int32(5) /* ValueType.STRING */, int32(0) /* ValueType.UNDETERMINED */ :
		return true
	}
	panic("unreachable switch expression")
}

// Java static method: NodeUtil.mayBeObject(com.google.javascript.rhino.Node n)
func NodeUtil_MayBeObject_1_Node(n pkg_rhino.Node) bool {
	_ = n
	return NodeUtil_MayBeObject_1_NodeUtil_ValueType(NodeUtil_GetKnownValueType(n))
}

// Java static method: NodeUtil.mayBeObject(com.google.javascript.jscomp.NodeUtil.ValueType type)
func NodeUtil_MayBeObject_1_NodeUtil_ValueType(type_ NodeUtil_ValueType) bool {
	_ = type_
	switch (type_).Ordinal() {
	case int32(6) /* ValueType.BOOLEAN */, int32(1) /* ValueType.NULL */, int32(3) /* ValueType.NUMBER */, int32(4) /* ValueType.BIGINT */, int32(5) /* ValueType.STRING */, int32(2) /* ValueType.VOID */ :
		return false
	case int32(7) /* ValueType.OBJECT */, int32(0) /* ValueType.UNDETERMINED */ :
		return true
	}
	panic("unreachable switch expression")
}

// Java static method: NodeUtil.isAssociative(com.google.javascript.rhino.Token type)
func NodeUtil_IsAssociative(type_ pkg_rhino.Token) bool {
	_ = type_
	switch (type_).Ordinal() {
	case int32(69) /* Token.AND */, int32(68) /* Token.OR */, int32(70) /* Token.COALESCE */, int32(1) /* Token.BITOR */, int32(2) /* Token.BITXOR */, int32(3) /* Token.BITAND */ :
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java static method: NodeUtil.isCommutative(com.google.javascript.rhino.Token type)
func NodeUtil_IsCommutative(type_ pkg_rhino.Token) bool {
	_ = type_
	switch (type_).Ordinal() {
	case int32(15) /* Token.MUL */, int32(1) /* Token.BITOR */, int32(2) /* Token.BITXOR */, int32(3) /* Token.BITAND */ :
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java static method: NodeUtil.isAssignmentOp(com.google.javascript.rhino.Node n)
func NodeUtil_IsAssignmentOp(n pkg_rhino.Node) bool {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(51) /* Token.ASSIGN */, int32(52) /* Token.ASSIGN_BITOR */, int32(53) /* Token.ASSIGN_BITXOR */, int32(54) /* Token.ASSIGN_BITAND */, int32(55) /* Token.ASSIGN_LSH */, int32(56) /* Token.ASSIGN_RSH */, int32(57) /* Token.ASSIGN_URSH */, int32(58) /* Token.ASSIGN_ADD */, int32(59) /* Token.ASSIGN_SUB */, int32(60) /* Token.ASSIGN_MUL */, int32(63) /* Token.ASSIGN_EXPONENT */, int32(61) /* Token.ASSIGN_DIV */, int32(62) /* Token.ASSIGN_MOD */, int32(64) /* Token.ASSIGN_OR */, int32(65) /* Token.ASSIGN_AND */, int32(66) /* Token.ASSIGN_COALESCE */ :
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java static method: NodeUtil.isLogicalAssignmentOp(com.google.javascript.rhino.Node n)
func NodeUtil_IsLogicalAssignmentOp(n pkg_rhino.Node) bool {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(64) /* Token.ASSIGN_OR */, int32(65) /* Token.ASSIGN_AND */, int32(66) /* Token.ASSIGN_COALESCE */ :
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java static method: NodeUtil.isCompoundAssignmentOp(com.google.javascript.rhino.Node n)
func NodeUtil_IsCompoundAssignmentOp(n pkg_rhino.Node) bool {
	_ = n
	return (NodeUtil_IsAssignmentOp(n) && !(n.IsAssign()))
}

// Java static method: NodeUtil.getOpFromAssignmentOp(com.google.javascript.rhino.Node n)
func NodeUtil_GetOpFromAssignmentOp(n pkg_rhino.Node) pkg_rhino.Token {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(52) /* Token.ASSIGN_BITOR */ :
		return pkg_rhino.StaticGet_Token_BITOR()
	case int32(53) /* Token.ASSIGN_BITXOR */ :
		return pkg_rhino.StaticGet_Token_BITXOR()
	case int32(54) /* Token.ASSIGN_BITAND */ :
		return pkg_rhino.StaticGet_Token_BITAND()
	case int32(55) /* Token.ASSIGN_LSH */ :
		return pkg_rhino.StaticGet_Token_LSH()
	case int32(56) /* Token.ASSIGN_RSH */ :
		return pkg_rhino.StaticGet_Token_RSH()
	case int32(57) /* Token.ASSIGN_URSH */ :
		return pkg_rhino.StaticGet_Token_URSH()
	case int32(58) /* Token.ASSIGN_ADD */ :
		return pkg_rhino.StaticGet_Token_ADD()
	case int32(59) /* Token.ASSIGN_SUB */ :
		return pkg_rhino.StaticGet_Token_SUB()
	case int32(60) /* Token.ASSIGN_MUL */ :
		return pkg_rhino.StaticGet_Token_MUL()
	case int32(63) /* Token.ASSIGN_EXPONENT */ :
		return pkg_rhino.StaticGet_Token_EXPONENT()
	case int32(61) /* Token.ASSIGN_DIV */ :
		return pkg_rhino.StaticGet_Token_DIV()
	case int32(62) /* Token.ASSIGN_MOD */ :
		return pkg_rhino.StaticGet_Token_MOD()
	case int32(64) /* Token.ASSIGN_OR */ :
		return pkg_rhino.StaticGet_Token_OR()
	case int32(65) /* Token.ASSIGN_AND */ :
		return pkg_rhino.StaticGet_Token_AND()
	case int32(66) /* Token.ASSIGN_COALESCE */ :
		return pkg_rhino.StaticGet_Token_COALESCE()
	default:
		panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Not an assignment op:"+j2gort.StrVal(n))))
	}
	panic("unreachable switch expression")
}

// Java static method: NodeUtil.getEnclosingType(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Token type)
func NodeUtil_GetEnclosingType(n pkg_rhino.Node, type_ pkg_rhino.Token) pkg_rhino.Node {
	_ = n
	_ = type_
	return NodeUtil_GetEnclosingNode(n, j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_22 uint32) bool {
		n1 := pkg_rhino.Node(rawNode_22)
		_ = n1
		return (n1.GetToken() == type_)
	}))
}

// Java static method: NodeUtil.getEnclosingNonArrowFunction(com.google.javascript.rhino.Node n)
func NodeUtil_GetEnclosingNonArrowFunction(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	return NodeUtil_GetEnclosingNode(n, lambdaSingleton_node_util_5)
}

// Java static method: NodeUtil.getEnclosingClass(com.google.javascript.rhino.Node n)
func NodeUtil_GetEnclosingClass(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	return NodeUtil_GetEnclosingNode(n, lambdaSingleton_node_util_6)
}

// Java static method: NodeUtil.getEnclosingModuleIfPresent(com.google.javascript.rhino.Node n)
func NodeUtil_GetEnclosingModuleIfPresent(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	return NodeUtil_GetEnclosingNode(n, lambdaSingleton_node_util_7)
}

// Java static method: NodeUtil.getEnclosingFunction(com.google.javascript.rhino.Node n)
func NodeUtil_GetEnclosingFunction(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	return NodeUtil_GetEnclosingNode(n, lambdaSingleton_node_util_8)
}

// Java static method: NodeUtil.getEnclosingScript(com.google.javascript.rhino.Node n)
func NodeUtil_GetEnclosingScript(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	return NodeUtil_GetEnclosingNode(n, lambdaSingleton_node_util_9)
}

// Java static method: NodeUtil.getEnclosingBlock(com.google.javascript.rhino.Node n)
func NodeUtil_GetEnclosingBlock(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	return NodeUtil_GetEnclosingNode(n, lambdaSingleton_node_util_10)
}

// Java static method: NodeUtil.getEnclosingBlockScopeRoot(com.google.javascript.rhino.Node n)
func NodeUtil_GetEnclosingBlockScopeRoot(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	return NodeUtil_GetEnclosingNode(n, lambdaSingleton_node_util_11)
}

// Java static method: NodeUtil.getEnclosingScopeRoot(com.google.javascript.rhino.Node n)
func NodeUtil_GetEnclosingScopeRoot(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	return NodeUtil_GetEnclosingNode(n, lambdaSingleton_node_util_12)
}

// Java static method: NodeUtil.getEnclosingHoistScopeRoot(com.google.javascript.rhino.Node n)
func NodeUtil_GetEnclosingHoistScopeRoot(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	return NodeUtil_GetEnclosingNode(n, lambdaSingleton_node_util_13)
}

// Java static method: NodeUtil.isHoistScopeRoot(com.google.javascript.rhino.Node n)
func NodeUtil_IsHoistScopeRoot(n pkg_rhino.Node) bool {
	_ = n
	return ((n.IsFunction() || n.IsModuleBody()) || NodeUtil_IsClassStaticBlock(n))
}

// Java static method: NodeUtil.isInFunction(com.google.javascript.rhino.Node n)
func NodeUtil_IsInFunction(n pkg_rhino.Node) bool {
	_ = n
	return (NodeUtil_GetEnclosingFunction(n) != 0)
}

// Java static method: NodeUtil.getEnclosingStatement(com.google.javascript.rhino.Node n)
func NodeUtil_GetEnclosingStatement(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	return NodeUtil_GetEnclosingNode(n, lambdaSingleton_node_util_14)
}

// Java static method: NodeUtil.getEnclosingNode(com.google.javascript.rhino.Node n, com.google.common.base.Predicate<com.google.javascript.rhino.Node> pred)
func NodeUtil_GetEnclosingNode(n pkg_rhino.Node, pred j2gort.Ext_com_google_common_base_Predicate) pkg_rhino.Node {
	_ = n
	_ = pred
	var curr pkg_rhino.Node = n
	_ = curr
	for (curr != 0) && !(j2gort.PredTestNode(pred, uint32(curr))) {
		curr = curr.GetParent()
	}
	return curr
}

// Java static method: NodeUtil.getFirstPropMatchingKey(com.google.javascript.rhino.Node n, java.lang.String keyName)
func NodeUtil_GetFirstPropMatchingKey(n pkg_rhino.Node, keyName string) pkg_rhino.Node {
	_ = n
	_ = keyName
	if !(n.IsObjectLit() || n.IsClassMembers()) {
		j2gort.CheckState(false)
	}
	{
		var keyNode pkg_rhino.Node = n.GetFirstChild()
		_ = keyNode
		for ; keyNode != 0; keyNode = keyNode.GetNext() {
			if (keyNode.IsStringKey() || keyNode.IsMemberFunctionDef()) && (j2gort.CheckNotNullString(keyNode.GetString()) == keyName) {
				return keyNode.GetFirstChild()
			} else if (keyNode.IsComputedProp() && keyNode.GetFirstChild().IsStringLit()) && (j2gort.CheckNotNullString(keyNode.GetFirstChild().GetString()) == keyName) {
				return keyNode.GetLastChild()
			}
		}
	}
	return pkg_rhino.Node(0)
}

// Java static method: NodeUtil.getFirstGetterMatchingKey(com.google.javascript.rhino.Node n, java.lang.String keyName)
func NodeUtil_GetFirstGetterMatchingKey(n pkg_rhino.Node, keyName string) pkg_rhino.Node {
	_ = n
	_ = keyName
	if !(n.IsClassMembers() || n.IsObjectLit()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(n))
	}
	{
		var keyNode pkg_rhino.Node = n.GetFirstChild()
		_ = keyNode
		for ; keyNode != 0; keyNode = keyNode.GetNext() {
			if keyNode.IsGetterDef() && (j2gort.CheckNotNullString(keyNode.GetString()) == keyName) {
				return keyNode
			}
		}
	}
	return pkg_rhino.Node(0)
}

// Java static method: NodeUtil.referencesOwnReceiver(com.google.javascript.rhino.Node fn)
func NodeUtil_ReferencesOwnReceiver(fn pkg_rhino.Node) bool {
	_ = fn
	if !(fn.IsFunction()) {
		j2gort.CheckState(false)
	}
	if fn.IsArrowFunction() {
		return false
	}
	return (NodeUtil_ReferencesEnclosingReceiver(NodeUtil_GetFunctionParameters(fn)) || NodeUtil_ReferencesEnclosingReceiver(NodeUtil_GetFunctionBody(fn)))
}

// Java static method: NodeUtil.referencesEnclosingReceiver(com.google.javascript.rhino.Node n)
func NodeUtil_ReferencesEnclosingReceiver(n pkg_rhino.Node) bool {
	_ = n
	return NodeUtil_Has(n, lambdaSingleton_node_util_15, StaticGet_NodeUtil_MATCH_ANYTHING_BUT_NON_ARROW_FUNCTION())
}

// Java static method: NodeUtil.referencesSuper(com.google.javascript.rhino.Node n)
func NodeUtil_ReferencesSuper(n pkg_rhino.Node) bool {
	_ = n
	var curr pkg_rhino.Node = n.GetFirstChild()
	_ = curr
	for curr != 0 {
		if NodeUtil_Has(curr, lambdaSingleton_node_util_16, lambdaSingleton_node_util_17) {
			return true
		}
		curr = curr.GetNext()
	}
	return false
}

// Java static method: NodeUtil.isNormalOrOptChainGet(com.google.javascript.rhino.Node n)
func NodeUtil_IsNormalOrOptChainGet(n pkg_rhino.Node) bool {
	_ = n
	return (NodeUtil_IsNormalGet(n) || NodeUtil_IsOptChainGet(n))
}

// Java static method: NodeUtil.isNormalOrOptChainGetProp(com.google.javascript.rhino.Node n)
func NodeUtil_IsNormalOrOptChainGetProp(n pkg_rhino.Node) bool {
	_ = n
	return (n.IsGetProp() || n.IsOptChainGetProp())
}

// Java static method: NodeUtil.isNormalOrOptChainCall(com.google.javascript.rhino.Node n)
func NodeUtil_IsNormalOrOptChainCall(n pkg_rhino.Node) bool {
	_ = n
	return (n.IsCall() || n.IsOptChainCall())
}

// Java static method: NodeUtil.isNormalGet(com.google.javascript.rhino.Node n)
func NodeUtil_IsNormalGet(n pkg_rhino.Node) bool {
	_ = n
	return (n.IsGetProp() || n.IsGetElem())
}

// Java static method: NodeUtil.isOptChainGet(com.google.javascript.rhino.Node n)
func NodeUtil_IsOptChainGet(n pkg_rhino.Node) bool {
	_ = n
	return (n.IsOptChainGetProp() || n.IsOptChainGetElem())
}

// Java static method: NodeUtil.isOptChainNode(com.google.javascript.rhino.Node n)
func NodeUtil_IsOptChainNode(n pkg_rhino.Node) bool {
	_ = n
	return ((n.IsOptChainGetProp() || n.IsOptChainGetElem()) || n.IsOptChainCall())
}

// Java static method: NodeUtil.getStartOfOptChainSegment(com.google.javascript.rhino.Node n)
func NodeUtil_GetStartOfOptChainSegment(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if !(NodeUtil_IsOptChainNode(n)) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(n))
	}
	if n.IsOptionalChainStart() {
		return n
	}
	return NodeUtil_GetStartOfOptChainSegment(n.GetFirstChild())
}

// Java static method: NodeUtil.getEndOfOptChainSegment(com.google.javascript.rhino.Node n)
func NodeUtil_GetEndOfOptChainSegment(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if !(NodeUtil_IsOptChainNode(n)) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(n))
	}
	if NodeUtil_IsEndOfOptChainSegment(n) {
		return n
	} else {
		return NodeUtil_GetEndOfOptChainSegment(n.GetParent())
	}
}

// Java static method: NodeUtil.isEndOfFullOptChain(com.google.javascript.rhino.Node n)
func NodeUtil_IsEndOfFullOptChain(n pkg_rhino.Node) bool {
	_ = n
	if NodeUtil_IsOptChainNode(n) {
		var parent pkg_rhino.Node = n.GetParent()
		_ = parent
		return !(NodeUtil_IsOptChainNode(parent) && n.IsFirstChildOf(parent))
	} else {
		return false
	}
}

// Java static method: NodeUtil.isEndOfOptChainSegment(com.google.javascript.rhino.Node n)
func NodeUtil_IsEndOfOptChainSegment(n pkg_rhino.Node) bool {
	_ = n
	if !(NodeUtil_IsOptChainNode(n)) {
		return false
	} else {
		var parent pkg_rhino.Node = n.GetParent()
		_ = parent
		if ((parent != 0) && n.IsFirstChildOf(parent)) && NodeUtil_IsOptChainNode(parent) {
			return parent.IsOptionalChainStart()
		} else {
			return true
		}
	}
}

// Java static method: NodeUtil.convertToNonOptionalChainSegment(com.google.javascript.rhino.Node endOfOptChainSegment)
func NodeUtil_ConvertToNonOptionalChainSegment(endOfOptChainSegment pkg_rhino.Node) {
	_ = endOfOptChainSegment
	var start pkg_rhino.Node = NodeUtil_GetStartOfOptChainSegment(endOfOptChainSegment)
	_ = start
	NodeUtil_ConvertToNonOptionalChainSegmentDownTo(endOfOptChainSegment, start.GetFirstChild())
}

// Java static method: NodeUtil.convertToNonOptionalChainSegmentDownTo(com.google.javascript.rhino.Node endOfOptChainSegment, com.google.javascript.rhino.Node stopNode)
func NodeUtil_ConvertToNonOptionalChainSegmentDownTo(endOfOptChainSegment pkg_rhino.Node, stopNode pkg_rhino.Node) {
	_ = endOfOptChainSegment
	_ = stopNode
	if !(NodeUtil_IsEndOfOptChainSegment(endOfOptChainSegment)) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(endOfOptChainSegment))
	}
	var segmentNodes j2gort.Deque[pkg_rhino.Node]
	_ = segmentNodes
	var segmentNode pkg_rhino.Node = endOfOptChainSegment
	_ = segmentNode
	for (segmentNode != 0) && !(j2gort.CheckNotNullNode(segmentNode) == stopNode) {
		if NodeUtil_IsOptChainNode(segmentNode) {
			segmentNodes.AddLast(j2gort.CheckNotNullNode(segmentNode))
		}
		segmentNode = segmentNode.GetFirstChild()
	}
	{
		cur_25 := segmentNodes.Iter()
		for cur_25.HasNext() {
			n := cur_25.Next()
			_ = n
			n.SetIsOptionalChainStart(false)
			n.SetToken(NodeUtil_GetNonOptChainToken(n.GetToken()))
		}
	}
}

// Java static method: NodeUtil.getNonOptChainToken(com.google.javascript.rhino.Token optChainToken)
func NodeUtil_GetNonOptChainToken(optChainToken pkg_rhino.Token) pkg_rhino.Token {
	_ = optChainToken
	switch (optChainToken).Ordinal() {
	case int32(31) /* Token.OPTCHAIN_CALL */ :
		return pkg_rhino.StaticGet_Token_CALL()
	case int32(30) /* Token.OPTCHAIN_GETELEM */ :
		return pkg_rhino.StaticGet_Token_GETELEM()
	case int32(29) /* Token.OPTCHAIN_GETPROP */ :
		return pkg_rhino.StaticGet_Token_GETPROP()
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Should be an OPTCHAIN token: "+j2gort.StrVal(optChainToken))))
	}
	panic("unreachable switch expression")
}

// Java static method: NodeUtil.isBlockScopedDeclaration(com.google.javascript.rhino.Node n)
func NodeUtil_IsBlockScopedDeclaration(n pkg_rhino.Node) bool {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(111) /* Token.LET */, int32(97) /* Token.CONST */, int32(86) /* Token.CATCH */ :
		return true
	case int32(105) /* Token.CLASS */ :
		return NodeUtil_IsClassDeclaration(n)
	case int32(73) /* Token.FUNCTION */ :
		return NodeUtil_IsBlockScopedFunctionDeclaration(n)
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java static method: NodeUtil.isNameDeclaration(com.google.javascript.rhino.Node n)
func NodeUtil_IsNameDeclaration(n pkg_rhino.Node) bool {
	_ = n
	return ((n != 0) && ((n.IsVar() || n.IsLet()) || n.IsConst()))
}

// Java static method: NodeUtil.isDestructuringDeclaration(com.google.javascript.rhino.Node n)
func NodeUtil_IsDestructuringDeclaration(n pkg_rhino.Node) bool {
	_ = n
	if NodeUtil_IsNameDeclaration(n) {
		{
			var c pkg_rhino.Node = n.GetFirstChild()
			_ = c
			for ; c != 0; c = c.GetNext() {
				if c.IsDestructuringLhs() {
					return true
				}
			}
		}
	}
	return false
}

// Java static method: NodeUtil.getAssignedValue(com.google.javascript.rhino.Node n)
func NodeUtil_GetAssignedValue(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if !(n.IsName() || n.IsGetProp()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(n))
	}
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	if NodeUtil_IsNameDeclaration(parent) {
		return n.GetFirstChild()
	} else if parent.IsAssign() && (parent.GetFirstChild() == n) {
		return n.GetNext()
	} else {
		return pkg_rhino.Node(0)
	}
}

// Java static method: NodeUtil.isExprAssign(com.google.javascript.rhino.Node n)
func NodeUtil_IsExprAssign(n pkg_rhino.Node) bool {
	_ = n
	return (n.IsExprResult() && n.GetFirstChild().IsAssign())
}

// Java static method: NodeUtil.isExprCall(com.google.javascript.rhino.Node n)
func NodeUtil_IsExprCall(n pkg_rhino.Node) bool {
	_ = n
	return (n.IsExprResult() && n.GetFirstChild().IsCall())
}

// Java static method: NodeUtil.isNonArrowFunction(com.google.javascript.rhino.Node n)
func NodeUtil_IsNonArrowFunction(n pkg_rhino.Node) bool {
	_ = n
	return (n.IsFunction() && !(n.IsArrowFunction()))
}

// Java static method: NodeUtil.isEnhancedFor(com.google.javascript.rhino.Node n)
func NodeUtil_IsEnhancedFor(n pkg_rhino.Node) bool {
	_ = n
	return ((n.IsForOf() || n.IsForAwaitOf()) || n.IsForIn())
}

// Java static method: NodeUtil.isAnyFor(com.google.javascript.rhino.Node n)
func NodeUtil_IsAnyFor(n pkg_rhino.Node) bool {
	_ = n
	return (((n.IsVanillaFor() || n.IsForIn()) || n.IsForOf()) || n.IsForAwaitOf())
}

// Java static method: NodeUtil.isLoopStructure(com.google.javascript.rhino.Node n)
func NodeUtil_IsLoopStructure(n pkg_rhino.Node) bool {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(80) /* Token.FOR */, int32(81) /* Token.FOR_IN */, int32(112) /* Token.FOR_OF */, int32(113) /* Token.FOR_AWAIT_OF */, int32(79) /* Token.DO */, int32(78) /* Token.WHILE */ :
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java static method: NodeUtil.getLoopCodeBlock(com.google.javascript.rhino.Node n)
func NodeUtil_GetLoopCodeBlock(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(80) /* Token.FOR */, int32(81) /* Token.FOR_IN */, int32(112) /* Token.FOR_OF */, int32(113) /* Token.FOR_AWAIT_OF */, int32(78) /* Token.WHILE */ :
		return n.GetLastChild()
	case int32(79) /* Token.DO */ :
		return n.GetFirstChild()
	default:
		return pkg_rhino.Node(0)
	}
	panic("unreachable switch expression")
}

// Java static method: NodeUtil.isWithinLoop(com.google.javascript.rhino.Node n)
func NodeUtil_IsWithinLoop(n pkg_rhino.Node) bool {
	_ = n
	{
		curAnc_ := n
		if curAnc_ == 0 {
			j2gort.ThrowNPE()
		}
		curAnc_ = curAnc_.GetParent()
		for curAnc_ != 0 {
			parent := curAnc_
			curAnc_ = curAnc_.GetParent()

			_ = parent
			if NodeUtil_IsLoopStructure(parent) {
				return true
			}
			if parent.IsFunction() {
				break
			}
		}
	}
	return false
}

// Java static method: NodeUtil.isControlStructure(com.google.javascript.rhino.Node n)
func NodeUtil_IsControlStructure(n pkg_rhino.Node) bool {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(80) /* Token.FOR */, int32(81) /* Token.FOR_IN */, int32(112) /* Token.FOR_OF */, int32(113) /* Token.FOR_AWAIT_OF */, int32(79) /* Token.DO */, int32(78) /* Token.WHILE */, int32(85) /* Token.WITH */, int32(74) /* Token.IF */, int32(92) /* Token.LABEL */, int32(48) /* Token.TRY */, int32(86) /* Token.CATCH */, int32(75) /* Token.SWITCH */, int32(76) /* Token.CASE */, int32(77) /* Token.DEFAULT_CASE */ :
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java static method: NodeUtil.isControlStructureCodeBlock(com.google.javascript.rhino.Node parent, com.google.javascript.rhino.Node n)
func NodeUtil_IsControlStructureCodeBlock(parent pkg_rhino.Node, n pkg_rhino.Node) bool {
	_ = parent
	_ = n
	switch parent.TokenSwitchOrd_() {
	case int32(79) /* Token.DO */ :
		return (parent.GetFirstChild() == n)
	case int32(48) /* Token.TRY */ :
		return ((parent.GetFirstChild() == n) || (parent.GetLastChild() == n))
	case int32(80) /* Token.FOR */, int32(81) /* Token.FOR_IN */, int32(112) /* Token.FOR_OF */, int32(113) /* Token.FOR_AWAIT_OF */, int32(78) /* Token.WHILE */, int32(92) /* Token.LABEL */, int32(85) /* Token.WITH */, int32(86) /* Token.CATCH */ :
		return (parent.GetLastChild() == n)
	case int32(74) /* Token.IF */, int32(75) /* Token.SWITCH */, int32(76) /* Token.CASE */ :
		return (parent.GetFirstChild() != n)
	case int32(77) /* Token.DEFAULT_CASE */ :
		return true
	default:
		if !(NodeUtil_IsControlStructure(parent)) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(parent))
		}
		return false
	}
	panic("unreachable switch expression")
}

// Java static method: NodeUtil.getConditionExpression(com.google.javascript.rhino.Node n)
func NodeUtil_GetConditionExpression(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(74) /* Token.IF */, int32(78) /* Token.WHILE */ :
		return n.GetFirstChild()
	case int32(79) /* Token.DO */ :
		return n.GetLastChild()
	case int32(80) /* Token.FOR */ :
		return n.GetSecondChild()
	case int32(81) /* Token.FOR_IN */, int32(112) /* Token.FOR_OF */, int32(113) /* Token.FOR_AWAIT_OF */, int32(76) /* Token.CASE */ :
		return pkg_rhino.Node(0)
	default:
		panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", j2gort.StrVal(n)+" does not have a condition.")))
	}
	panic("unreachable switch expression")
}

// Java static method: NodeUtil.isStatementBlock(com.google.javascript.rhino.Node n)
func NodeUtil_IsStatementBlock(n pkg_rhino.Node) bool {
	_ = n
	return (((n.IsRoot() || n.IsScript()) || n.IsBlock()) || n.IsModuleBody())
}

// Java static method: NodeUtil.createsBlockScope(com.google.javascript.rhino.Node n)
func NodeUtil_CreatesBlockScope(n pkg_rhino.Node) bool {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(90) /* Token.BLOCK */ :
		var parent pkg_rhino.Node = n.GetParent()
		_ = parent
		if ((parent != 0) && (parent.IsBlock() || parent.IsScript())) && n.IsSyntheticBlock() {
			return false
		}
		return (((parent != 0) && !(NodeUtil_IsSwitchCase(parent))) && !(parent.IsCatch()))
	case int32(80) /* Token.FOR */, int32(81) /* Token.FOR_IN */, int32(112) /* Token.FOR_OF */, int32(113) /* Token.FOR_AWAIT_OF */, int32(91) /* Token.SWITCH_BODY */, int32(105) /* Token.CLASS */ :
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java static method: NodeUtil.createsScope(com.google.javascript.rhino.Node n)
func NodeUtil_CreatesScope(n pkg_rhino.Node) bool {
	_ = n
	return (((((NodeUtil_CreatesBlockScope(n) || n.IsFunction()) || n.IsModuleBody()) || n.IsMemberFieldDef()) || n.IsComputedFieldDef()) || (n.IsRoot() && (n.GetParent() == 0)))
}

// Java static method: NodeUtil.isValidCfgRoot(com.google.javascript.rhino.Node n)
func NodeUtil_IsValidCfgRoot(n pkg_rhino.Node) bool {
	_ = n
	return (j2gort.ImmSetContainsTypedU8_(StaticGet_NodeUtil_DEFINITE_CFG_ROOTS(), uint8(n.GetToken())) || NodeUtil_IsClassStaticBlock(n))
}

// Java static method: NodeUtil.isStatement(com.google.javascript.rhino.Node n)
func NodeUtil_IsStatement(n pkg_rhino.Node) bool {
	_ = n
	return (((!(n.IsModuleBody()) && !(n.IsScript())) && !(n.IsRoot())) && NodeUtil_IsStatementParent(n.GetParent()))
}

// Java static method: NodeUtil.isStatementParent(com.google.javascript.rhino.Node parent)
func NodeUtil_IsStatementParent(parent pkg_rhino.Node) bool {
	_ = parent
	return j2gort.ImmSetContainsTypedU8_(StaticGet_NodeUtil_IS_STATEMENT_PARENT(), uint8(parent.GetToken()))
}

// Java static method: NodeUtil.isDeclarationParent(com.google.javascript.rhino.Node parent)
func NodeUtil_IsDeclarationParent(parent pkg_rhino.Node) bool {
	_ = parent
	switch parent.TokenSwitchOrd_() {
	case int32(171) /* Token.DECLARE */, int32(120) /* Token.EXPORT */ :
		return true
	default:
		return NodeUtil_IsStatementParent(parent)
	}
	panic("unreachable switch expression")
}

// Java static method: NodeUtil.isSwitchCase(com.google.javascript.rhino.Node n)
func NodeUtil_IsSwitchCase(n pkg_rhino.Node) bool {
	_ = n
	return (n.IsCase() || n.IsDefaultCase())
}

// Java static method: NodeUtil.isReferenceName(com.google.javascript.rhino.Node n)
func NodeUtil_IsReferenceName(n pkg_rhino.Node) bool {
	_ = n
	return (n.IsName() && !(j2gort.CheckNotNullString(n.GetString()) == ""))
}

// Java static method: NodeUtil.isNonlocalModuleExportName(com.google.javascript.rhino.Node n)
func NodeUtil_IsNonlocalModuleExportName(n pkg_rhino.Node) bool {
	_ = n
	if !(n.IsName()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(n))
	}
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	if parent.IsImportSpec() && n.IsFirstChildOf(parent) {
		return true
	} else if parent.IsExportSpec() {
		if n.IsFirstChildOf(parent) {
			return NodeUtil_IsExportFrom(parent.GetGrandparent())
		} else {
			return true
		}
	}
	return false
}

// Java static method: NodeUtil.isTryFinallyNode(com.google.javascript.rhino.Node parent, com.google.javascript.rhino.Node child)
func NodeUtil_IsTryFinallyNode(parent pkg_rhino.Node, child pkg_rhino.Node) bool {
	_ = parent
	_ = child
	return ((parent.IsTry() && parent.HasXChildren(int32(3))) && (child == parent.GetLastChild()))
}

// Java static method: NodeUtil.isTryCatchNodeContainer(com.google.javascript.rhino.Node n)
func NodeUtil_IsTryCatchNodeContainer(n pkg_rhino.Node) bool {
	_ = n
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	return (parent.IsTry() && (parent.GetSecondChild() == n))
}

// Java static method: NodeUtil.isInSyntheticScript(com.google.javascript.rhino.Node n)
func NodeUtil_IsInSyntheticScript(n pkg_rhino.Node) bool {
	_ = n
	var sourceFileName string = n.GetSourceFileName()
	_ = sourceFileName
	return ((sourceFileName != j2gort.NullString) && (javalang.StringStartsWith(j2gort.CheckNotNullString(sourceFileName), " [synthetic:") || javalang.StringStartsWith(j2gort.CheckNotNullString(sourceFileName), "src/com/google/javascript/jscomp/js/")))
}

// Java static method: NodeUtil.deleteNode(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.AbstractCompiler compiler)
func NodeUtil_DeleteNode(n pkg_rhino.Node, compiler AbstractCompiler) {
	_ = n
	_ = compiler
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	NodeUtil_MarkFunctionsDeleted_2_Node_AbstractCompiler(n, compiler)
	n.Detach()
	compiler.ReportChangeToEnclosingScope(parent)
}

// Java static method: NodeUtil.deleteFunctionCall(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.AbstractCompiler compiler)
func NodeUtil_DeleteFunctionCall(n pkg_rhino.Node, compiler AbstractCompiler) {
	_ = n
	_ = compiler
	if !(n.IsCall()) {
		j2gort.CheckState(false)
	}
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	if parent.IsExprResult() {
		var grandParent pkg_rhino.Node = parent.GetParent()
		_ = grandParent
		parent.Detach()
		parent = grandParent
	} else {
		n.ReplaceWith(NodeUtil_NewUndefinedNode(n))
	}
	NodeUtil_MarkFunctionsDeleted_2_Node_AbstractCompiler(n, compiler)
	compiler.ReportChangeToEnclosingScope(parent)
}

// Java static method: NodeUtil.deleteChildren(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.AbstractCompiler compiler)
func NodeUtil_DeleteChildren(n pkg_rhino.Node, compiler AbstractCompiler) {
	_ = n
	_ = compiler
	for n.HasChildren() {
		NodeUtil_DeleteNode(n.GetFirstChild(), compiler)
	}
}

// Java static method: NodeUtil.removeChild(com.google.javascript.rhino.Node parent, com.google.javascript.rhino.Node node)
func NodeUtil_RemoveChild(parent pkg_rhino.Node, node pkg_rhino.Node) {
	_ = parent
	_ = node
	if NodeUtil_IsTryFinallyNode(parent, node) {
		if NodeUtil_HasCatchHandler(NodeUtil_GetCatchBlock(parent)) {
			node.Detach()
		} else {
			node.DetachChildren()
		}
	} else if node.IsCatch() {
		var tryNode pkg_rhino.Node = node.GetGrandparent()
		_ = tryNode
		if !(NodeUtil_HasFinally(tryNode)) {
			j2gort.CheckState(false)
		}
		node.Detach()
	} else if NodeUtil_IsTryCatchNodeContainer(node) {
		var tryNode_2 pkg_rhino.Node = node.GetParent()
		_ = tryNode_2
		if !(NodeUtil_HasFinally(tryNode_2)) {
			j2gort.CheckState(false)
		}
		node.DetachChildren()
	} else if node.IsBlock() {
		node.DetachChildren()
	} else if (NodeUtil_IsStatementBlock(parent) || NodeUtil_IsSwitchCase(node)) || node.IsMemberFunctionDef() {
		node.Detach()
	} else if NodeUtil_IsNameDeclaration(parent) || parent.IsExprResult() {
		if parent.HasMoreThanOneChild() {
			node.Detach()
		} else {
			node.Detach()
			NodeUtil_RemoveChild(parent.GetParent(), parent)
		}
	} else if parent.IsLabel() && (node == parent.GetLastChild()) {
		node.Detach()
		NodeUtil_RemoveChild(parent.GetParent(), parent)
	} else if parent.IsVanillaFor() {
		node.ReplaceWith(pkg_rhino.IR_Empty())
	} else if parent.IsObjectPattern() {
		node.Detach()
	} else if parent.IsArrayPattern() {
		if node == parent.GetLastChild() {
			node.Detach()
		} else {
			node.ReplaceWith(pkg_rhino.IR_Empty())
		}
	} else if parent.IsDestructuringLhs() {
		node.Detach()
		if parent.GetParent().HasChildren() {
			NodeUtil_RemoveChild(parent.GetParent(), parent)
		}
	} else if parent.IsRest() {
		parent.Detach()
	} else if parent.IsParamList() {
		node.Detach()
	} else if parent.IsImport() {
		if node == parent.GetFirstChild() {
			node.ReplaceWith(pkg_rhino.IR_Empty())
		} else {
			panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", j2gort.BoxString(("Invalid attempt to remove: " + j2gort.StrVal(node) + " from " + j2gort.StrVal(parent))))))
		}
	} else {
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", j2gort.BoxString(("Invalid attempt to remove node: " + j2gort.StrVal(node) + " of " + j2gort.StrVal(parent))))))
	}
}

// Java static method: NodeUtil.replaceDeclarationChild(com.google.javascript.rhino.Node declChild, com.google.javascript.rhino.Node newStatement)
func NodeUtil_ReplaceDeclarationChild(declChild pkg_rhino.Node, newStatement pkg_rhino.Node) {
	_ = declChild
	_ = newStatement
	if !(NodeUtil_IsNameDeclaration(declChild.GetParent())) {
		j2gort.CheckArgument(false)
	}
	if !(0 == newStatement.GetParent()) {
		j2gort.CheckArgument(false)
	}
	var decl pkg_rhino.Node = declChild.GetParent()
	_ = decl
	if decl.HasOneChild() {
		decl.ReplaceWith(newStatement)
	} else if declChild.GetNext() == 0 {
		declChild.Detach()
		newStatement.InsertAfter(decl)
	} else if declChild.GetPrevious() == 0 {
		declChild.Detach()
		newStatement.InsertBefore(decl)
	} else {
		if !(decl.HasMoreThanOneChild()) {
			j2gort.CheckState(false)
		}
		var newDecl pkg_rhino.Node = pkg_rhino.NewNode_1(decl.GetToken()).Srcref(decl)
		_ = newDecl
		{
			var after pkg_rhino.Node = declChild.GetNext()
			_ = after
			var next pkg_rhino.Node = pkg_rhino.Node(0)
			_ = next
			for ; after != 0; after = next {
				next = after.GetNext()
				newDecl.AddChildToBack(after.Detach())
			}
		}
		declChild.Detach()
		newStatement.InsertAfter(decl)
		newDecl.InsertAfter(newStatement)
	}
}

// Java static method: NodeUtil.maybeAddFinally(com.google.javascript.rhino.Node tryNode)
func NodeUtil_MaybeAddFinally(tryNode pkg_rhino.Node) {
	_ = tryNode
	if !(tryNode.IsTry()) {
		j2gort.CheckState(false)
	}
	if !(NodeUtil_HasFinally(tryNode)) {
		tryNode.AddChildToBack(pkg_rhino.IR_Block_0().Srcref(tryNode))
	}
}

// Java static method: NodeUtil.tryMergeBlock(com.google.javascript.rhino.Node block, boolean ignoreBlockScopedDeclarations)
func NodeUtil_TryMergeBlock(block pkg_rhino.Node, ignoreBlockScopedDeclarations bool) bool {
	_ = block
	_ = ignoreBlockScopedDeclarations
	if !(block.IsBlock()) {
		j2gort.CheckState(false)
	}
	var parent pkg_rhino.Node = block.GetParent()
	_ = parent
	var canMerge bool = (ignoreBlockScopedDeclarations || NodeUtil_CanMergeBlock(block))
	_ = canMerge
	if NodeUtil_IsStatementBlock(parent) && canMerge {
		var previous pkg_rhino.Node = block
		_ = previous
		for block.HasChildren() {
			var child pkg_rhino.Node = block.RemoveFirstChild()
			_ = child
			child.InsertAfter(previous)
			previous = child
		}
		block.Detach()
		return true
	} else {
		return false
	}
}

// Java static method: NodeUtil.canMergeBlock(com.google.javascript.rhino.Node block)
func NodeUtil_CanMergeBlock(block pkg_rhino.Node) bool {
	_ = block
	{
		var c pkg_rhino.Node = block.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			switch c.TokenSwitchOrd_() {
			case int32(92) /* Token.LABEL */ :
				if NodeUtil_CanMergeBlock(c) {
					continue
				} else {
					return false
				}
			case int32(97) /* Token.CONST */, int32(111) /* Token.LET */, int32(105) /* Token.CLASS */, int32(73) /* Token.FUNCTION */ :
				return false
			default:
				continue
			}
		}
	}
	return true
}

// Java static method: NodeUtil.isCallOrNew(com.google.javascript.rhino.Node node)
func NodeUtil_IsCallOrNew(node pkg_rhino.Node) bool {
	_ = node
	return ((node.IsCall() || node.IsNew()) || node.IsOptChainCall())
}

// Java static method: NodeUtil.getFunctionBody(com.google.javascript.rhino.Node fn)
func NodeUtil_GetFunctionBody(fn pkg_rhino.Node) pkg_rhino.Node {
	_ = fn
	if !(fn.IsFunction()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(fn))
	}
	return fn.GetLastChild()
}

// Java static method: NodeUtil.getCallTargetResolvingIndirectCalls(com.google.javascript.rhino.Node call)
func NodeUtil_GetCallTargetResolvingIndirectCalls(call pkg_rhino.Node) pkg_rhino.Node {
	_ = call
	if !((call.IsCall() || call.IsNew()) || call.IsTaggedTemplateLit()) {
		j2gort.CheckArgumentFmt(false, "must be call, new or tagged-template-literal expression, got %s", j2gort.BoxNode(call))
	}
	var target pkg_rhino.Node = call.GetFirstChild()
	_ = target
	if target.IsComma() && target.GetSecondChild().IsQualifiedName() {
		return target.GetSecondChild()
	}
	return target
}

// Java static method: NodeUtil.isDeclaration(com.google.javascript.rhino.Node n)
func NodeUtil_IsDeclaration(n pkg_rhino.Node) bool {
	_ = n
	return ((NodeUtil_IsNameDeclaration(n) || NodeUtil_IsFunctionDeclaration(n)) || NodeUtil_IsClassDeclaration(n))
}

// Java static method: NodeUtil.isNamedExportsLiteral(com.google.javascript.rhino.Node objectLiteral)
func NodeUtil_IsNamedExportsLiteral(objectLiteral pkg_rhino.Node) bool {
	_ = objectLiteral
	if !(objectLiteral.IsObjectLit()) || !(objectLiteral.HasChildren()) {
		return false
	}
	{
		var key pkg_rhino.Node = objectLiteral.GetFirstChild()
		_ = key
		for ; key != 0; key = key.GetNext() {
			if !(key.IsStringKey()) || key.IsQuotedStringKey() {
				return false
			}
			if !(key.GetFirstChild().IsName()) {
				return false
			}
		}
	}
	return true
}

// Java static method: NodeUtil.isFunctionDeclaration(com.google.javascript.rhino.Node n)
func NodeUtil_IsFunctionDeclaration(n pkg_rhino.Node) bool {
	_ = n
	return ((n.IsFunction() && NodeUtil_IsDeclarationParent(n.GetParent())) && NodeUtil_IsNamedFunction(n))
}

// Java static method: NodeUtil.isMethodDeclaration(com.google.javascript.rhino.Node n)
func NodeUtil_IsMethodDeclaration(n pkg_rhino.Node) bool {
	_ = n
	if n.IsFunction() {
		var parent pkg_rhino.Node = n.GetParent()
		_ = parent
		switch parent.TokenSwitchOrd_() {
		case int32(95) /* Token.GETTER_DEF */, int32(96) /* Token.SETTER_DEF */, int32(107) /* Token.MEMBER_FUNCTION_DEF */ :
			return true
		case int32(129) /* Token.COMPUTED_PROP */ :
			return ((parent.GetLastChild() == n) && ((parent.GetBooleanProp(pkg_rhino.StaticGet_Node_COMPUTED_PROP_METHOD()) || parent.GetBooleanProp(pkg_rhino.StaticGet_Node_COMPUTED_PROP_GETTER())) || parent.GetBooleanProp(pkg_rhino.StaticGet_Node_COMPUTED_PROP_SETTER())))
		default:
			return false
		}
		panic("unreachable switch expression")
	} else {
		return false
	}
}

// Java static method: NodeUtil.isClassDeclaration(com.google.javascript.rhino.Node n)
func NodeUtil_IsClassDeclaration(n pkg_rhino.Node) bool {
	_ = n
	return ((n.IsClass() && NodeUtil_IsDeclarationParent(n.GetParent())) && NodeUtil_IsNamedClass(n))
}

// Java static method: NodeUtil.isHoistedFunctionDeclaration(com.google.javascript.rhino.Node n)
func NodeUtil_IsHoistedFunctionDeclaration(n pkg_rhino.Node) bool {
	_ = n
	if NodeUtil_IsFunctionDeclaration(n) {
		var parent pkg_rhino.Node = n.GetParent()
		_ = parent
		return (((parent.IsScript() || parent.IsModuleBody()) || parent.GetParent().IsFunction()) || parent.IsExport())
	}
	return false
}

// Java static method: NodeUtil.isBlockScopedFunctionDeclaration(com.google.javascript.rhino.Node n)
func NodeUtil_IsBlockScopedFunctionDeclaration(n pkg_rhino.Node) bool {
	_ = n
	if !(NodeUtil_IsFunctionDeclaration(n)) {
		return false
	}
	var current pkg_rhino.Node = n.GetParent()
	_ = current
	for current != 0 {
		switch current.TokenSwitchOrd_() {
		case int32(90) /* Token.BLOCK */ :
			return !(current.GetParent().IsFunction())
		case int32(73) /* Token.FUNCTION */, int32(94) /* Token.SCRIPT */, int32(171) /* Token.DECLARE */, int32(120) /* Token.EXPORT */, int32(123) /* Token.MODULE_BODY */ :
			return false
		default:
			if !(current.IsLabel()) {
				j2gort.CheckStateMsg(false, j2gort.BoxNode(current))
			}
			current = current.GetParent()
		}
	}
	return false
}

// Java static method: NodeUtil.isFunctionBlock(com.google.javascript.rhino.Node n)
func NodeUtil_IsFunctionBlock(n pkg_rhino.Node) bool {
	_ = n
	return ((n.IsBlock() && n.HasParent()) && n.GetParent().IsFunction())
}

// Java static method: NodeUtil.isClassStaticBlock(com.google.javascript.rhino.Node n)
func NodeUtil_IsClassStaticBlock(n pkg_rhino.Node) bool {
	_ = n
	return ((n.IsBlock() && n.HasParent()) && n.GetParent().IsClassMembers())
}

// Java static method: NodeUtil.isNonStaticClassFieldInitializer(com.google.javascript.rhino.Node n)
func NodeUtil_IsNonStaticClassFieldInitializer(n pkg_rhino.Node) bool {
	_ = n
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	if parent.IsMemberFieldDef() || (parent.IsComputedFieldDef() && (n == parent.GetSecondChild())) {
		return !(parent.IsStaticMember())
	}
	return false
}

// Java static method: NodeUtil.isFunctionOrNonStaticClassFieldInitializer(com.google.javascript.rhino.Node n)
func NodeUtil_IsFunctionOrNonStaticClassFieldInitializer(n pkg_rhino.Node) bool {
	_ = n
	return (n.IsFunction() || NodeUtil_IsNonStaticClassFieldInitializer(n))
}

// Java static method: NodeUtil.isFunctionExpression(com.google.javascript.rhino.Node n)
func NodeUtil_IsFunctionExpression(n pkg_rhino.Node) bool {
	_ = n
	return ((n.IsFunction() && !(NodeUtil_IsFunctionDeclaration(n))) && !(NodeUtil_IsMethodDeclaration(n)))
}

// Java static method: NodeUtil.isNamedFunctionExpression(com.google.javascript.rhino.Node n)
func NodeUtil_IsNamedFunctionExpression(n pkg_rhino.Node) bool {
	_ = n
	return (NodeUtil_IsFunctionExpression(n) && !(j2gort.CheckNotNullString(n.GetFirstChild().GetString()) == ""))
}

// Java static method: NodeUtil.isClassExpression(com.google.javascript.rhino.Node n)
func NodeUtil_IsClassExpression(n pkg_rhino.Node) bool {
	_ = n
	return (n.IsClass() && (!(NodeUtil_IsNamedClass(n)) || !(NodeUtil_IsDeclarationParent(n.GetParent()))))
}

// Java static method: NodeUtil.isNamedClassExpression(com.google.javascript.rhino.Node n)
func NodeUtil_IsNamedClassExpression(n pkg_rhino.Node) bool {
	_ = n
	return (NodeUtil_IsClassExpression(n) && n.GetFirstChild().IsName())
}

// Java static method: NodeUtil.isNamedFunction(com.google.javascript.rhino.Node n)
func NodeUtil_IsNamedFunction(n pkg_rhino.Node) bool {
	_ = n
	return (n.IsFunction() && NodeUtil_IsReferenceName(n.GetFirstChild()))
}

// Java static method: NodeUtil.isNamedClass(com.google.javascript.rhino.Node n)
func NodeUtil_IsNamedClass(n pkg_rhino.Node) bool {
	_ = n
	return (n.IsClass() && NodeUtil_IsReferenceName(n.GetFirstChild()))
}

// Java static method: NodeUtil.isBleedingFunctionName(com.google.javascript.rhino.Node n)
func NodeUtil_IsBleedingFunctionName(n pkg_rhino.Node) bool {
	_ = n
	if !(n.IsName()) || (j2gort.CheckNotNullString(n.GetString()) == "") {
		return false
	}
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	return (NodeUtil_IsFunctionExpression(parent) && (n == parent.GetFirstChild()))
}

// Java static method: NodeUtil.isEmptyFunctionExpression(com.google.javascript.rhino.Node node)
func NodeUtil_IsEmptyFunctionExpression(node pkg_rhino.Node) bool {
	_ = node
	return (NodeUtil_IsFunctionExpression(node) && NodeUtil_IsEmptyBlock(node.GetLastChild()))
}

// Java static method: NodeUtil.doesFunctionReferenceOwnArgumentsObject(com.google.javascript.rhino.Node fn)
func NodeUtil_DoesFunctionReferenceOwnArgumentsObject(fn pkg_rhino.Node) bool {
	_ = fn
	if !(fn.IsFunction()) {
		j2gort.CheckArgument(false)
	}
	if fn.IsArrowFunction() {
		return false
	}
	return (NodeUtil_ReferencesArgumentsHelper(fn.GetSecondChild()) || NodeUtil_ReferencesArgumentsHelper(fn.GetLastChild()))
}

// Java static method: NodeUtil.referencesArgumentsHelper(com.google.javascript.rhino.Node node)
func NodeUtil_ReferencesArgumentsHelper(node pkg_rhino.Node) bool {
	_ = node
	if node.IsName() && (j2gort.CheckNotNullString(node.GetString()) == "arguments") {
		return true
	}
	if NodeUtil_IsNonArrowFunction(node) {
		return false
	}
	{
		var c pkg_rhino.Node = node.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			if NodeUtil_ReferencesArgumentsHelper(c) {
				return true
			}
		}
	}
	return false
}

// Java static method: NodeUtil.isObjectCallMethod(com.google.javascript.rhino.Node callNode, java.lang.String methodName)
func NodeUtil_IsObjectCallMethod(callNode pkg_rhino.Node, methodName string) bool {
	_ = callNode
	_ = methodName
	if callNode.IsCall() || callNode.IsOptChainCall() {
		var callee pkg_rhino.Node = callNode.GetFirstChild()
		_ = callee
		if NodeUtil_IsNormalOrOptChainGetProp(callee) {
			return (j2gort.CheckNotNullString(callee.GetString()) == methodName)
		} else if NodeUtil_IsNormalOrOptChainGet(callee) {
			var last pkg_rhino.Node = callee.GetLastChild()
			_ = last
			if (last != 0) && last.IsStringLit() {
				var propName string = last.GetString()
				_ = propName
				return (j2gort.CheckNotNullString(propName) == methodName)
			}
		}
	}
	return false
}

// Java static method: NodeUtil.isFunctionObjectCall(com.google.javascript.rhino.Node callNode)
func NodeUtil_IsFunctionObjectCall(callNode pkg_rhino.Node) bool {
	_ = callNode
	return NodeUtil_IsObjectCallMethod(callNode, "call")
}

// Java static method: NodeUtil.isFunctionObjectApply(com.google.javascript.rhino.Node callNode)
func NodeUtil_IsFunctionObjectApply(callNode pkg_rhino.Node) bool {
	_ = callNode
	return NodeUtil_IsObjectCallMethod(callNode, "apply")
}

// Java static method: NodeUtil.isNameDeclOrSimpleAssignLhs(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func NodeUtil_IsNameDeclOrSimpleAssignLhs(n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = n
	_ = parent
	return ((parent.IsAssign() && (parent.GetFirstChild() == n)) || NodeUtil_IsNameDeclaration(parent))
}

// Java static method: NodeUtil.isLValue(com.google.javascript.rhino.Node n)
func NodeUtil_IsLValue(n pkg_rhino.Node) bool {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(32) /* Token.NAME */, int32(26) /* Token.GETPROP */, int32(27) /* Token.GETELEM */ :
	default:
		return false
	}
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	if parent == 0 {
		return false
	}
	switch parent.TokenSwitchOrd_() {
	case int32(118) /* Token.IMPORT_SPEC */ :
		return (parent.GetLastChild() == n)
	case int32(84) /* Token.VAR */, int32(111) /* Token.LET */, int32(97) /* Token.CONST */, int32(125) /* Token.ITER_REST */, int32(126) /* Token.OBJECT_REST */, int32(49) /* Token.PARAM_LIST */, int32(116) /* Token.IMPORT */, int32(71) /* Token.INC */, int32(72) /* Token.DEC */, int32(86) /* Token.CATCH */ :
		return true
	case int32(105) /* Token.CLASS */, int32(73) /* Token.FUNCTION */, int32(134) /* Token.DEFAULT_VALUE */, int32(80) /* Token.FOR */, int32(81) /* Token.FOR_IN */, int32(112) /* Token.FOR_OF */, int32(113) /* Token.FOR_AWAIT_OF */ :
		return (parent.GetFirstChild() == n)
	case int32(102) /* Token.ARRAY_PATTERN */, int32(100) /* Token.STRING_KEY */, int32(129) /* Token.COMPUTED_PROP */ :
		return NodeUtil_IsLhsByDestructuring(n)
	default:
		return (NodeUtil_IsAssignmentOp(parent) && (parent.GetFirstChild() == n))
	}
	panic("unreachable switch expression")
}

// Java static method: NodeUtil.isLValueWrite(com.google.javascript.rhino.Node n)
func NodeUtil_IsLValueWrite(n pkg_rhino.Node) bool {
	_ = n
	if !(NodeUtil_IsLValue(n)) {
		return false
	}
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	var grandparent pkg_rhino.Node = parent.GetParent()
	_ = grandparent
	return ((!(NodeUtil_IsNameDeclaration(parent)) || n.HasChildren()) || ((grandparent != 0) && NodeUtil_IsEnhancedFor(grandparent)))
}

// Java static method: NodeUtil.isDeclarationLValue(com.google.javascript.rhino.Node n)
func NodeUtil_IsDeclarationLValue(n pkg_rhino.Node) bool {
	_ = n
	var isLValue bool = NodeUtil_IsLValue(n)
	_ = isLValue
	if !(isLValue) {
		return false
	}
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	switch parent.TokenSwitchOrd_() {
	case int32(118) /* Token.IMPORT_SPEC */, int32(84) /* Token.VAR */, int32(111) /* Token.LET */, int32(97) /* Token.CONST */, int32(49) /* Token.PARAM_LIST */, int32(116) /* Token.IMPORT */, int32(86) /* Token.CATCH */, int32(105) /* Token.CLASS */, int32(73) /* Token.FUNCTION */ :
		return true
	case int32(100) /* Token.STRING_KEY */, int32(103) /* Token.OBJECT_PATTERN */, int32(102) /* Token.ARRAY_PATTERN */, int32(129) /* Token.COMPUTED_PROP */, int32(134) /* Token.DEFAULT_VALUE */, int32(125) /* Token.ITER_REST */, int32(126) /* Token.OBJECT_REST */ :
		var rootTargetParent pkg_rhino.Node = NodeUtil_GetRootTarget(n).GetParent()
		_ = rootTargetParent
		if rootTargetParent.IsRest() || rootTargetParent.IsDefaultValue() {
			return rootTargetParent.GetParent().IsParamList()
		} else if rootTargetParent.IsDestructuringLhs() {
			return NodeUtil_IsNameDeclaration(rootTargetParent.GetParent())
		}
		return (rootTargetParent.IsParamList() || rootTargetParent.IsCatch())
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java static method: NodeUtil.isLhsOfAssign(com.google.javascript.rhino.Node n)
func NodeUtil_IsLhsOfAssign(n pkg_rhino.Node) bool {
	_ = n
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	return (((parent != 0) && parent.IsAssign()) && (parent.GetFirstChild() == n))
}

// Java static method: NodeUtil.isImportedName(com.google.javascript.rhino.Node n)
func NodeUtil_IsImportedName(n pkg_rhino.Node) bool {
	_ = n
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	return (parent.IsImport() || (parent.IsImportSpec() && (parent.GetLastChild() == n)))
}

// Java static method: NodeUtil.getDeclaringParent(com.google.javascript.rhino.Node targetNode)
func NodeUtil_GetDeclaringParent(targetNode pkg_rhino.Node) pkg_rhino.Node {
	_ = targetNode
	var rootTarget pkg_rhino.Node = NodeUtil_GetRootTarget(targetNode)
	_ = rootTarget
	var parent pkg_rhino.Node = rootTarget.GetParent()
	_ = parent
	if parent.IsRest() || parent.IsDefaultValue() {
		parent = parent.GetParent()
		if !(parent.IsParamList()) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(parent))
		}
	} else if parent.IsDestructuringLhs() {
		parent = parent.GetParent()
		if !(NodeUtil_IsNameDeclaration(parent)) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(parent))
		}
	} else if parent.IsClass() || parent.IsFunction() {
		if !(targetNode == parent.GetFirstChild()) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(targetNode))
		}
	} else if parent.IsImportSpec() {
		if !(targetNode == parent.GetSecondChild()) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(targetNode))
		}
		parent = parent.GetGrandparent()
		if !(parent.IsImport()) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(parent))
		}
	} else {
		if !(((parent.IsParamList() || NodeUtil_IsNameDeclaration(parent)) || parent.IsImport()) || parent.IsCatch()) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(parent))
		}
	}
	return parent
}

// Java static method: NodeUtil.getRootTarget(com.google.javascript.rhino.Node targetNode)
func NodeUtil_GetRootTarget(targetNode pkg_rhino.Node) pkg_rhino.Node {
	_ = targetNode
	var enclosingTarget pkg_rhino.Node = targetNode
	_ = enclosingTarget
	{
		var nextTarget pkg_rhino.Node = NodeUtil_GetEnclosingTarget(enclosingTarget)
		_ = nextTarget
		for ; nextTarget != 0; nextTarget = NodeUtil_GetEnclosingTarget(enclosingTarget) {
			enclosingTarget = nextTarget
		}
	}
	return enclosingTarget
}

// Java static method: NodeUtil.getEnclosingTarget(com.google.javascript.rhino.Node targetNode)
func NodeUtil_GetEnclosingTarget(targetNode pkg_rhino.Node) pkg_rhino.Node {
	_ = targetNode
	if !(j2gort.CheckNotNullNode(targetNode).IsValidAssignmentTarget()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(targetNode))
	}
	var parent pkg_rhino.Node = j2gort.CheckNotNullNode(targetNode.GetParent())
	_ = parent
	var targetIsFirstChild bool = (parent.GetFirstChild() == targetNode)
	_ = targetIsFirstChild
	if parent.IsDefaultValue() || parent.IsRest() {
		if !(targetIsFirstChild) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(parent))
		}
		targetNode = parent
		parent = j2gort.CheckNotNullNode(targetNode.GetParent())
		targetIsFirstChild = (targetNode == parent.GetFirstChild())
	}
	switch parent.TokenSwitchOrd_() {
	case int32(102) /* Token.ARRAY_PATTERN */ :
		return parent
	case int32(103) /* Token.OBJECT_PATTERN */ :
		return parent
	case int32(129) /* Token.COMPUTED_PROP */ :
		if !(!(targetIsFirstChild)) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(parent))
		}
		var grandparent pkg_rhino.Node = j2gort.CheckNotNullNode(parent.GetParent())
		_ = grandparent
		if !(grandparent.IsObjectPattern()) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(grandparent))
		}
		return grandparent
	case int32(100) /* Token.STRING_KEY */ :
		var grandparent_2 pkg_rhino.Node = j2gort.CheckNotNullNode(parent.GetParent())
		_ = grandparent_2
		if !(grandparent_2.IsObjectPattern()) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(grandparent_2))
		}
		return grandparent_2
	case int32(49) /* Token.PARAM_LIST */, int32(111) /* Token.LET */, int32(97) /* Token.CONST */, int32(84) /* Token.VAR */ :
		return pkg_rhino.Node(0)
	case int32(73) /* Token.FUNCTION */, int32(105) /* Token.CLASS */ :
		if !(targetIsFirstChild) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(targetNode))
		}
		return pkg_rhino.Node(0)
	case int32(81) /* Token.FOR_IN */, int32(112) /* Token.FOR_OF */, int32(113) /* Token.FOR_AWAIT_OF */ :
		if !(targetIsFirstChild) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(targetNode))
		}
		return pkg_rhino.Node(0)
	case int32(104) /* Token.DESTRUCTURING_LHS */ :
		if !(targetIsFirstChild) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(targetNode))
		}
		return pkg_rhino.Node(0)
	case int32(116) /* Token.IMPORT */ :
		return pkg_rhino.Node(0)
	case int32(118) /* Token.IMPORT_SPEC */ :
		if !(!(targetIsFirstChild)) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(parent))
		}
		return pkg_rhino.Node(0)
	case int32(86) /* Token.CATCH */ :
		return pkg_rhino.Node(0)
	default:
		if !(NodeUtil_IsAssignmentOp(parent) && targetIsFirstChild) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(parent))
		}
		return pkg_rhino.Node(0)
	}
	panic("unreachable switch expression")
}

// Java static method: NodeUtil.isLhsByDestructuring(com.google.javascript.rhino.Node n)
func NodeUtil_IsLhsByDestructuring(n pkg_rhino.Node) bool {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(32) /* Token.NAME */, int32(26) /* Token.GETPROP */, int32(27) /* Token.GETELEM */ :
		return NodeUtil_IsLhsByDestructuringHelper(n)
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java static method: NodeUtil.isLhsByDestructuringHelper(com.google.javascript.rhino.Node n)
func NodeUtil_IsLhsByDestructuringHelper(n pkg_rhino.Node) bool {
	_ = n
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	var grandparent pkg_rhino.Node = n.GetGrandparent()
	_ = grandparent
	switch parent.TokenSwitchOrd_() {
	case int32(102) /* Token.ARRAY_PATTERN */, int32(125) /* Token.ITER_REST */, int32(126) /* Token.OBJECT_REST */ :
		return true
	case int32(129) /* Token.COMPUTED_PROP */ :
		return (!(n.IsFirstChildOf(parent)) && grandparent.IsObjectPattern())
	case int32(100) /* Token.STRING_KEY */ :
		return grandparent.IsObjectPattern()
	case int32(134) /* Token.DEFAULT_VALUE */ :
		return (n.IsFirstChildOf(parent) && NodeUtil_IsLhsByDestructuringHelper(parent))
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java static method: NodeUtil.mayBeObjectLitKey(com.google.javascript.rhino.Node node)
func NodeUtil_MayBeObjectLitKey(node pkg_rhino.Node) bool {
	_ = node
	switch node.TokenSwitchOrd_() {
	case int32(100) /* Token.STRING_KEY */, int32(95) /* Token.GETTER_DEF */, int32(96) /* Token.SETTER_DEF */, int32(107) /* Token.MEMBER_FUNCTION_DEF */ :
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java static method: NodeUtil.isObjectLitKey(com.google.javascript.rhino.Node node)
func NodeUtil_IsObjectLitKey(node pkg_rhino.Node) bool {
	_ = node
	return (node.GetParent().IsObjectLit() && NodeUtil_MayBeObjectLitKey(node))
}

// Java static method: NodeUtil.getObjectOrClassLitKeyName(com.google.javascript.rhino.Node key)
func NodeUtil_GetObjectOrClassLitKeyName(key pkg_rhino.Node) string {
	_ = key
	var keyNode pkg_rhino.Node = NodeUtil_GetObjectOrClassLitKeyNode(key)
	_ = keyNode
	if keyNode != 0 {
		return keyNode.GetString()
	}
	panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unexpected node type: "+j2gort.StrVal(key))))
}

// Java static method: NodeUtil.getObjectOrClassLitKeyNode(com.google.javascript.rhino.Node key)
func NodeUtil_GetObjectOrClassLitKeyNode(key pkg_rhino.Node) pkg_rhino.Node {
	_ = key
	switch key.TokenSwitchOrd_() {
	case int32(100) /* Token.STRING_KEY */, int32(95) /* Token.GETTER_DEF */, int32(96) /* Token.SETTER_DEF */, int32(107) /* Token.MEMBER_FUNCTION_DEF */, int32(108) /* Token.MEMBER_FIELD_DEF */ :
		return key
	case int32(129) /* Token.COMPUTED_PROP */, int32(109) /* Token.COMPUTED_FIELD_DEF */ :
		if key.GetFirstChild().IsStringLit() {
			return key.GetFirstChild()
		} else {
			return pkg_rhino.Node(0)
		}
	default:
	}
	panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unexpected node type: "+j2gort.StrVal(key))))
}

// Java static method: NodeUtil.isGetOrSetKey(com.google.javascript.rhino.Node node)
func NodeUtil_IsGetOrSetKey(node pkg_rhino.Node) bool {
	_ = node
	switch node.TokenSwitchOrd_() {
	case int32(95) /* Token.GETTER_DEF */, int32(96) /* Token.SETTER_DEF */ :
		return true
	case int32(129) /* Token.COMPUTED_PROP */ :
		return (node.GetBooleanProp(pkg_rhino.StaticGet_Node_COMPUTED_PROP_GETTER()) || node.GetBooleanProp(pkg_rhino.StaticGet_Node_COMPUTED_PROP_SETTER()))
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java static method: NodeUtil.opToStr(com.google.javascript.rhino.Token operator)
func NodeUtil_OpToStr(operator pkg_rhino.Token) string {
	_ = operator
	switch (operator).Ordinal() {
	case int32(70) /* Token.COALESCE */ :
		return "??"
	case int32(1) /* Token.BITOR */ :
		return "|"
	case int32(68) /* Token.OR */ :
		return "||"
	case int32(2) /* Token.BITXOR */ :
		return "^"
	case int32(69) /* Token.AND */ :
		return "&&"
	case int32(3) /* Token.BITAND */ :
		return "&"
	case int32(40) /* Token.SHEQ */ :
		return "==="
	case int32(4) /* Token.EQ */ :
		return "=="
	case int32(19) /* Token.NOT */ :
		return "!"
	case int32(5) /* Token.NE */ :
		return "!="
	case int32(41) /* Token.SHNE */ :
		return "!=="
	case int32(10) /* Token.LSH */ :
		return "<<"
	case int32(44) /* Token.IN */ :
		return "in"
	case int32(7) /* Token.LE */ :
		return "<="
	case int32(6) /* Token.LT */ :
		return "<"
	case int32(12) /* Token.URSH */ :
		return ">>>"
	case int32(11) /* Token.RSH */ :
		return ">>"
	case int32(9) /* Token.GE */ :
		return ">="
	case int32(8) /* Token.GT */ :
		return ">"
	case int32(15) /* Token.MUL */ :
		return "*"
	case int32(16) /* Token.DIV */ :
		return "/"
	case int32(17) /* Token.MOD */ :
		return "%"
	case int32(18) /* Token.EXPONENT */ :
		return "**"
	case int32(20) /* Token.BITNOT */ :
		return "~"
	case int32(13) /* Token.ADD */, int32(21) /* Token.POS */ :
		return "+"
	case int32(14) /* Token.SUB */, int32(22) /* Token.NEG */ :
		return "-"
	case int32(51) /* Token.ASSIGN */ :
		return "="
	case int32(52) /* Token.ASSIGN_BITOR */ :
		return "|="
	case int32(53) /* Token.ASSIGN_BITXOR */ :
		return "^="
	case int32(54) /* Token.ASSIGN_BITAND */ :
		return "&="
	case int32(55) /* Token.ASSIGN_LSH */ :
		return "<<="
	case int32(56) /* Token.ASSIGN_RSH */ :
		return ">>="
	case int32(57) /* Token.ASSIGN_URSH */ :
		return ">>>="
	case int32(58) /* Token.ASSIGN_ADD */ :
		return "+="
	case int32(59) /* Token.ASSIGN_SUB */ :
		return "-="
	case int32(60) /* Token.ASSIGN_MUL */ :
		return "*="
	case int32(63) /* Token.ASSIGN_EXPONENT */ :
		return "**="
	case int32(61) /* Token.ASSIGN_DIV */ :
		return "/="
	case int32(62) /* Token.ASSIGN_MOD */ :
		return "%="
	case int32(64) /* Token.ASSIGN_OR */ :
		return "||="
	case int32(65) /* Token.ASSIGN_AND */ :
		return "&&="
	case int32(66) /* Token.ASSIGN_COALESCE */ :
		return "??="
	case int32(87) /* Token.VOID */ :
		return "void"
	case int32(25) /* Token.TYPEOF */ :
		return "typeof"
	case int32(45) /* Token.INSTANCEOF */ :
		return "instanceof"
	default:
		return j2gort.NullString
	}
	panic("unreachable switch expression")
}

// Java static method: NodeUtil.opToStrNoFail(com.google.javascript.rhino.Token operator)
func NodeUtil_OpToStrNoFail(operator pkg_rhino.Token) string {
	_ = operator
	var res string = NodeUtil_OpToStr(operator)
	_ = res
	if res == j2gort.NullString {
		panic(j2gort.AsIface[javalang.IError](j2gort.ConstructShimThrowable("java.lang.Error", "Unknown op "+j2gort.StrVal(operator))))
	}
	return res
}

// Java static method: NodeUtil.redeclareVarsInsideBranch(com.google.javascript.rhino.Node branch)
func NodeUtil_RedeclareVarsInsideBranch(branch pkg_rhino.Node) {
	_ = branch
	if !(NodeUtil_CanContainHoistedVarsDecls(branch)) {
		return
	}
	var vars j2gort.Ext_java_util_Collection = NodeUtil_GetVarsDeclaredInBranch(branch)
	_ = vars
	if vars.IsEmpty() {
		return
	}
	var parent pkg_rhino.Node = NodeUtil_GetAddingRoot(branch)
	_ = parent
	{
		cur_27 := j2gort.NewIterCursor(j2gort.CheckNotNull(vars))
		for cur_27.HasNext() {
			nameNode := pkg_rhino.Node(cur_27.NextNode())
			_ = nameNode
			var var_ pkg_rhino.Node = pkg_rhino.IR_Var_1(pkg_rhino.IR_Name(nameNode.GetString()).Srcref(nameNode)).Srcref(nameNode)
			_ = var_
			NodeUtil_CopyNameAnnotations(nameNode, var_.GetFirstChild())
			parent.AddChildToFront(var_)
		}
	}
}

// Java static method: NodeUtil.canContainHoistedVarsDecls(com.google.javascript.rhino.Node n)
func NodeUtil_CanContainHoistedVarsDecls(n pkg_rhino.Node) bool {
	_ = n
	if NodeUtil_IsStatement(n) {
		switch n.TokenSwitchOrd_() {
		case int32(93) /* Token.EXPR_RESULT */, int32(0) /* Token.RETURN */, int32(43) /* Token.THROW */, int32(82) /* Token.BREAK */, int32(83) /* Token.CONTINUE */, int32(88) /* Token.EMPTY */, int32(98) /* Token.DEBUGGER */ :
			return false
		default:
			return true
		}
		panic("unreachable switch expression")
	}
	return true
}

// Java static method: NodeUtil.copyNameAnnotations(com.google.javascript.rhino.Node source, com.google.javascript.rhino.Node destination)
func NodeUtil_CopyNameAnnotations(source pkg_rhino.Node, destination pkg_rhino.Node) {
	_ = source
	_ = destination
	if source.GetBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME()) {
		destination.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME(), true)
	}
}

// Java static method: NodeUtil.getAddingRoot(com.google.javascript.rhino.Node n)
func NodeUtil_GetAddingRoot(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	var addingRoot pkg_rhino.Node = pkg_rhino.Node(0)
	_ = addingRoot
	var ancestor pkg_rhino.Node = n
	_ = ancestor
crawl_ancestors:
	for 0 != (func() pkg_rhino.Node {
		var asgn_28 pkg_rhino.Node = ancestor.GetParent()
		ancestor = asgn_28
		return asgn_28
	}()) {
		switch ancestor.TokenSwitchOrd_() {
		case int32(94) /* Token.SCRIPT */, int32(123) /* Token.MODULE_BODY */ :
			addingRoot = ancestor
			break crawl_ancestors
		case int32(73) /* Token.FUNCTION */ :
			addingRoot = ancestor.GetLastChild()
			break crawl_ancestors
		default:
			continue crawl_ancestors
		}
	}
	if !((addingRoot.IsBlock() || addingRoot.IsModuleBody()) || addingRoot.IsScript()) {
		j2gort.CheckState(false)
	}
	if !(!(addingRoot.HasChildren()) || !(addingRoot.GetFirstChild().IsScript())) {
		j2gort.CheckState(false)
	}
	return addingRoot
}

// Java static method: NodeUtil.newDeclaration(com.google.javascript.rhino.Node lhs, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node rhs, com.google.javascript.rhino.Token declarationType)
func NodeUtil_NewDeclaration(lhs pkg_rhino.Node, rhs pkg_rhino.Node, declarationType pkg_rhino.Token) pkg_rhino.Node {
	_ = lhs
	_ = rhs
	_ = declarationType
	if rhs == 0 {
		return pkg_rhino.IR_Declaration_2(lhs, declarationType)
	}
	return pkg_rhino.IR_Declaration_3(lhs, rhs, declarationType)
}

// Java static method: NodeUtil.newQName(com.google.javascript.jscomp.AbstractCompiler compiler, java.lang.String name)
func NodeUtil_NewQName_2_AbstractCompiler_string(compiler AbstractCompiler, name string) pkg_rhino.Node {
	_ = compiler
	_ = name
	return NodeUtil_NewQName_2_CodingConvention_string(compiler.GetCodingConvention(), name)
}

// Java static method: NodeUtil.newQName(com.google.javascript.jscomp.CodingConvention convention, java.lang.String name)
func NodeUtil_NewQName_2_CodingConvention_string(convention CodingConvention, name string) pkg_rhino.Node {
	_ = convention
	_ = name
	var endPos int32 = j2gort.JStrIndexOfCp(name, int32(uint16(46)))
	_ = endPos
	if endPos == int32(-1) {
		endPos = j2gort.JStrLen(name)
	}
	var qname pkg_rhino.Node = pkg_rhino.Node(0)
	_ = qname
	var nodeName string = j2gort.JStrSub(name, int32(0), endPos)
	_ = nodeName
	if "this" == nodeName {
		qname = pkg_rhino.IR_ThisNode()
	} else if "super" == nodeName {
		qname = pkg_rhino.IR_SuperNode()
	} else {
		qname = NodeUtil_NewName_2(convention, nodeName)
	}
	qname.SetLength(endPos)
	{
		var startPos int32 = (endPos + int32(1))
		_ = startPos
		for ; endPos < j2gort.JStrLen(name); startPos = (endPos + int32(1)) {
			endPos = j2gort.JStrIndexOfCpFrom(name, int32(uint16(46)), startPos)
			if endPos == int32(-1) {
				endPos = j2gort.JStrLen(name)
			}
			var part string = j2gort.JStrSub(name, startPos, endPos)
			_ = part
			qname = pkg_rhino.IR_Getprop_2(qname, part)
			if convention.IsConstantKey(part) {
				qname.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME(), true)
			}
			qname.SetLength(j2gort.JStrLen(part))
		}
	}
	return qname
}

// Java static method: NodeUtil.newQName(com.google.javascript.jscomp.AbstractCompiler compiler, java.lang.String name, com.google.javascript.rhino.Node basisNode, java.lang.String originalName)
func NodeUtil_NewQName_4_AbstractCompiler_string_Node_string(compiler AbstractCompiler, name string, basisNode pkg_rhino.Node, originalName string) pkg_rhino.Node {
	_ = compiler
	_ = name
	_ = basisNode
	_ = originalName
	return NodeUtil_NewQName_4_CodingConvention_string_Node_string(compiler.GetCodingConvention(), name, basisNode, originalName)
}

// Java static method: NodeUtil.newQName(com.google.javascript.jscomp.CodingConvention convention, java.lang.String name, com.google.javascript.rhino.Node basisNode, java.lang.String originalName)
func NodeUtil_NewQName_4_CodingConvention_string_Node_string(convention CodingConvention, name string, basisNode pkg_rhino.Node, originalName string) pkg_rhino.Node {
	_ = convention
	_ = name
	_ = basisNode
	_ = originalName
	var node pkg_rhino.Node = NodeUtil_NewQName_2_CodingConvention_string(convention, name)
	_ = node
	node.SrcrefTreeIfMissing(basisNode)
	{
		var iifeCond_1 bool
		{

			sr_29 := originalName
			se_30 := node.GetOriginalName()
			iifeCond_1 = j2gort.CheckNotNullString(sr_29) == se_30
		}
		if !(iifeCond_1) {
			node.SetOriginalName(originalName)
		}
	}
	return node
}

// Java static method: NodeUtil.getDeclarationFromName(com.google.javascript.rhino.Node nameNode, com.google.javascript.rhino.Node value, com.google.javascript.rhino.Token type, com.google.javascript.rhino.JSDocInfo info)
func NodeUtil_GetDeclarationFromName(nameNode pkg_rhino.Node, value pkg_rhino.Node, type_ pkg_rhino.Token, info *pkg_rhino.JSDocInfo) pkg_rhino.Node {
	_ = nameNode
	_ = value
	_ = type_
	_ = info
	var result pkg_rhino.Node = pkg_rhino.Node(0)
	_ = result
	if nameNode.IsName() {
		if value == 0 {
			result = pkg_rhino.IR_Declaration_2(nameNode, type_)
		} else {
			result = pkg_rhino.IR_Declaration_3(nameNode, value, type_)
		}
		result.SetJSDocInfo(info)
	} else if value != 0 {
		result = pkg_rhino.IR_ExprResult(pkg_rhino.IR_Assign(nameNode, value))
		result.GetFirstChild().SetJSDocInfo(info)
	} else {
		result = pkg_rhino.IR_ExprResult(nameNode)
		result.GetFirstChild().SetJSDocInfo(info)
	}
	return result
}

// Java static method: NodeUtil.newPropertyAccess(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.rhino.Node context, java.lang.String name)
func NodeUtil_NewPropertyAccess_3_AbstractCompiler_Node_string(compiler AbstractCompiler, context pkg_rhino.Node, name string) pkg_rhino.Node {
	_ = compiler
	_ = context
	_ = name
	return NodeUtil_NewPropertyAccess_3_CodingConvention_Node_string(compiler.GetCodingConvention(), context, name)
}

// Java static method: NodeUtil.newPropertyAccess(com.google.javascript.jscomp.CodingConvention convention, com.google.javascript.rhino.Node context, java.lang.String name)
func NodeUtil_NewPropertyAccess_3_CodingConvention_Node_string(convention CodingConvention, context pkg_rhino.Node, name string) pkg_rhino.Node {
	_ = convention
	_ = context
	_ = name
	var propNode pkg_rhino.Node = pkg_rhino.IR_Getprop_2(context, name)
	_ = propNode
	if convention.IsConstantKey(name) {
		propNode.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME(), true)
	}
	return propNode
}

// Java static method: NodeUtil.newQNameDeclaration(com.google.javascript.jscomp.AbstractCompiler compiler, java.lang.String name, com.google.javascript.rhino.Node value, com.google.javascript.rhino.JSDocInfo info)
func NodeUtil_NewQNameDeclaration_4_AbstractCompiler_string_Node_JSDocInfo(compiler AbstractCompiler, name string, value pkg_rhino.Node, info *pkg_rhino.JSDocInfo) pkg_rhino.Node {
	_ = compiler
	_ = name
	_ = value
	_ = info
	return NodeUtil_NewQNameDeclaration_4_CodingConvention_string_Node_JSDocInfo(compiler.GetCodingConvention(), name, value, info)
}

// Java static method: NodeUtil.newQNameDeclaration(com.google.javascript.jscomp.CodingConvention convention, java.lang.String name, com.google.javascript.rhino.Node value, com.google.javascript.rhino.JSDocInfo info)
func NodeUtil_NewQNameDeclaration_4_CodingConvention_string_Node_JSDocInfo(convention CodingConvention, name string, value pkg_rhino.Node, info *pkg_rhino.JSDocInfo) pkg_rhino.Node {
	_ = convention
	_ = name
	_ = value
	_ = info
	return NodeUtil_NewQNameDeclaration_5_CodingConvention_string_Node_JSDocInfo_Token(convention, name, value, info, pkg_rhino.StaticGet_Token_VAR())
}

// Java static method: NodeUtil.newQNameDeclaration(com.google.javascript.jscomp.AbstractCompiler compiler, java.lang.String name, com.google.javascript.rhino.Node value, com.google.javascript.rhino.JSDocInfo info, com.google.javascript.rhino.Token type)
func NodeUtil_NewQNameDeclaration_5_AbstractCompiler_string_Node_JSDocInfo_Token(compiler AbstractCompiler, name string, value pkg_rhino.Node, info *pkg_rhino.JSDocInfo, type_ pkg_rhino.Token) pkg_rhino.Node {
	_ = compiler
	_ = name
	_ = value
	_ = info
	_ = type_
	return NodeUtil_NewQNameDeclaration_5_CodingConvention_string_Node_JSDocInfo_Token(compiler.GetCodingConvention(), name, value, info, type_)
}

// Java static method: NodeUtil.newQNameDeclaration(com.google.javascript.jscomp.CodingConvention convention, java.lang.String name, com.google.javascript.rhino.Node value, com.google.javascript.rhino.JSDocInfo info, com.google.javascript.rhino.Token type)
func NodeUtil_NewQNameDeclaration_5_CodingConvention_string_Node_JSDocInfo_Token(convention CodingConvention, name string, value pkg_rhino.Node, info *pkg_rhino.JSDocInfo, type_ pkg_rhino.Token) pkg_rhino.Node {
	_ = convention
	_ = name
	_ = value
	_ = info
	_ = type_
	if !(((type_ == pkg_rhino.StaticGet_Token_VAR()) || (type_ == pkg_rhino.StaticGet_Token_LET())) || (type_ == pkg_rhino.StaticGet_Token_CONST())) {
		j2gort.CheckStateMsg(false, j2gort.BoxNullableVal(type_))
	}
	var nameNode pkg_rhino.Node = NodeUtil_NewQName_2_CodingConvention_string(convention, name)
	_ = nameNode
	return NodeUtil_GetDeclarationFromName(nameNode, value, type_, info)
}

// Java static method: NodeUtil.getRootOfQualifiedName(com.google.javascript.rhino.Node qName)
func NodeUtil_GetRootOfQualifiedName_1_Node(qName pkg_rhino.Node) pkg_rhino.Node {
	_ = qName
	{
		var current pkg_rhino.Node = qName
		_ = current
		for ; true; current = current.GetFirstChild() {
			if (current.IsName() || current.IsThis()) || current.IsSuper() {
				return current
			}
			if !(current.IsGetProp()) {
				j2gort.CheckStateFmt(false, "Not a getprop node:  (%s)", j2gort.BoxNode(current))
			}
		}
	}
	panic("unreachable")
}

// Java static method: NodeUtil.getRootOfQualifiedName(java.lang.String qName)
func NodeUtil_GetRootOfQualifiedName_1_string(qName string) string {
	_ = qName
	var dot int32 = j2gort.JStrIndexOfCp(qName, int32(uint16(46)))
	_ = dot
	if dot == int32(-1) {
		return qName
	}
	return j2gort.JStrSub(qName, int32(0), dot)
}

// Java static method: NodeUtil.newNameWithoutLength(com.google.javascript.jscomp.CodingConvention convention, java.lang.String name)
func NodeUtil_NewNameWithoutLength(convention CodingConvention, name string) pkg_rhino.Node {
	_ = convention
	_ = name
	var nameNode pkg_rhino.Node = pkg_rhino.IR_Name(name)
	_ = nameNode
	if convention.IsConstant(name) {
		nameNode.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME(), true)
	}
	return nameNode
}

// Java static method: NodeUtil.newName(com.google.javascript.jscomp.CodingConvention convention, java.lang.String name)
func NodeUtil_NewName_2(convention CodingConvention, name string) pkg_rhino.Node {
	_ = convention
	_ = name
	var nameNode pkg_rhino.Node = NodeUtil_NewNameWithoutLength(convention, name)
	_ = nameNode
	nameNode.SetLength(j2gort.JStrLen(name))
	return nameNode
}

// Java static method: NodeUtil.newName(com.google.javascript.jscomp.AbstractCompiler compiler, java.lang.String name, com.google.javascript.rhino.Node srcref)
func NodeUtil_NewName_3_AbstractCompiler_string_Node(compiler AbstractCompiler, name string, srcref pkg_rhino.Node) pkg_rhino.Node {
	_ = compiler
	_ = name
	_ = srcref
	return NodeUtil_NewName_3_CodingConvention_string_Node(compiler.GetCodingConvention(), name, srcref)
}

// Java static method: NodeUtil.newName(com.google.javascript.jscomp.CodingConvention convention, java.lang.String name, com.google.javascript.rhino.Node srcref)
func NodeUtil_NewName_3_CodingConvention_string_Node(convention CodingConvention, name string, srcref pkg_rhino.Node) pkg_rhino.Node {
	_ = convention
	_ = name
	_ = srcref
	return NodeUtil_NewNameWithoutLength(convention, name).Srcref(srcref)
}

// Java static method: NodeUtil.newName(com.google.javascript.jscomp.AbstractCompiler compiler, java.lang.String name, com.google.javascript.rhino.Node basisNode, java.lang.String originalName)
func NodeUtil_NewName_4_AbstractCompiler_string_Node_string(compiler AbstractCompiler, name string, basisNode pkg_rhino.Node, originalName string) pkg_rhino.Node {
	_ = compiler
	_ = name
	_ = basisNode
	_ = originalName
	return NodeUtil_NewName_4_CodingConvention_string_Node_string(compiler.GetCodingConvention(), name, basisNode, originalName)
}

// Java static method: NodeUtil.newName(com.google.javascript.jscomp.CodingConvention convention, java.lang.String name, com.google.javascript.rhino.Node basisNode, java.lang.String originalName)
func NodeUtil_NewName_4_CodingConvention_string_Node_string(convention CodingConvention, name string, basisNode pkg_rhino.Node, originalName string) pkg_rhino.Node {
	_ = convention
	_ = name
	_ = basisNode
	_ = originalName
	var nameNode pkg_rhino.Node = NodeUtil_NewNameWithoutLength(convention, name).SrcrefNoOrigName(basisNode)
	_ = nameNode
	nameNode.SetOriginalName(originalName)
	return nameNode
}

// Java static method: NodeUtil.isLatin(java.lang.String s)
func NodeUtil_IsLatin(s string) bool {
	_ = s
	var len_ int32 = j2gort.JStrLen(s)
	_ = len_
	{
		var index int32 = int32(0)
		_ = index
		for ; index < len_; index = index + 1 {
			var c uint16 = j2gort.JStrCharAt(s, index)
			_ = c
			if int32(c) > int32(uint16(127)) {
				return false
			}
		}
	}
	return true
}

// Java static method: NodeUtil.isValidSimpleName(java.lang.String name)
func NodeUtil_IsValidSimpleName(name string) bool {
	_ = name
	return ((pkg_rhino.TokenStream_IsJSIdentifier(name) && !(pkg_rhino.TokenStream_IsKeyword(name))) && NodeUtil_IsLatin(name))
}

// Java static method: NodeUtil.isValidQualifiedName(com.google.javascript.jscomp.CompilerOptions.LanguageMode mode, java.lang.String name)
func NodeUtil_IsValidQualifiedName_2_CompilerOptions_LanguageMode_string(mode CompilerOptions_LanguageMode, name string) bool {
	_ = mode
	_ = name
	return NodeUtil_IsValidQualifiedName_2_FeatureSet_string(mode.ToFeatureSet(), name)
}

// Java static method: NodeUtil.isValidQualifiedName(com.google.javascript.jscomp.parsing.parser.FeatureSet mode, java.lang.String name)
func NodeUtil_IsValidQualifiedName_2_FeatureSet_string(mode *pkg_parser.FeatureSet, name string) bool {
	_ = mode
	_ = name
	var len_ int32 = j2gort.JStrLen(name)
	_ = len_
	if ((len_ == int32(0)) || (int32(j2gort.JStrCharAt(name, int32(0))) == int32(uint16(46)))) || (int32(j2gort.JStrCharAt(name, (len_-int32(1)))) == int32(uint16(46))) {
		return false
	}
	var dot int32 = j2gort.JStrIndexOfCp(name, int32(uint16(46)))
	_ = dot
	if dot == int32(-1) {
		return NodeUtil_IsValidSimpleName(name)
	}
	if !(NodeUtil_IsValidSimpleName(j2gort.JStrSub(name, int32(0), dot))) {
		return false
	}
	var start int32 = (dot + int32(1))
	_ = start
	for {
		var nextDot int32 = j2gort.JStrIndexOfCpFrom(name, int32(uint16(46)), start)
		_ = nextDot
		if nextDot == int32(-1) {
			return NodeUtil_IsValidPropertyName(mode, j2gort.JStrSubFrom(name, start))
		}
		if (nextDot == start) || !(NodeUtil_IsValidPropertyName(mode, j2gort.JStrSub(name, start, nextDot))) {
			return false
		}
		start = (nextDot + int32(1))
	}
	panic("unreachable")
}

// Java static method: NodeUtil.isValidPropertyName(com.google.javascript.jscomp.parsing.parser.FeatureSet mode, java.lang.String name)
func NodeUtil_IsValidPropertyName(mode *pkg_parser.FeatureSet, name string) bool {
	_ = mode
	_ = name
	if NodeUtil_IsValidSimpleName(name) {
		return true
	} else {
		return (mode.Has(pkg_parser.StaticGet_FeatureSet_Feature_KEYWORDS_AS_PROPERTIES()) && pkg_rhino.TokenStream_IsKeyword(name))
	}
}

// Java static method: NodeUtil.getVarsDeclaredInBranch(com.google.javascript.rhino.Node root)
func NodeUtil_GetVarsDeclaredInBranch(root pkg_rhino.Node) j2gort.Ext_java_util_Collection {
	_ = root
	var val_31 NodeUtil_VarCollector
	InitNodeUtil_VarCollector(&val_31)
	var collector *NodeUtil_VarCollector = &val_31
	_ = collector
	NodeUtil_VisitPreOrder_3(root, collector, StaticGet_NodeUtil_MATCH_NOT_FUNCTION())
	return j2gort.AsIface[j2gort.Ext_java_util_Collection](collector.vars.LiveValues())
}

// Java static method: NodeUtil.getLhsNodesHelper(com.google.javascript.rhino.Node n, java.util.function.Consumer<com.google.javascript.rhino.Node> consumer)
func NodeUtil_GetLhsNodesHelper(n pkg_rhino.Node, consumer j2gort.Ext_java_util_function_Consumer) {
	_ = n
	_ = consumer
	switch n.TokenSwitchOrd_() {
	case int32(116) /* Token.IMPORT */ :
		NodeUtil_GetLhsNodesHelper(n.GetFirstChild(), consumer)
		NodeUtil_GetLhsNodesHelper(n.GetSecondChild(), consumer)
		return
	case int32(84) /* Token.VAR */, int32(97) /* Token.CONST */, int32(111) /* Token.LET */, int32(103) /* Token.OBJECT_PATTERN */, int32(102) /* Token.ARRAY_PATTERN */, int32(49) /* Token.PARAM_LIST */, int32(117) /* Token.IMPORT_SPECS */ :
		{
			var child pkg_rhino.Node = n.GetFirstChild()
			_ = child
			for ; child != 0; child = child.GetNext() {
				NodeUtil_GetLhsNodesHelper(child, consumer)
			}
		}
		return
	case int32(104) /* Token.DESTRUCTURING_LHS */, int32(134) /* Token.DEFAULT_VALUE */, int32(86) /* Token.CATCH */, int32(125) /* Token.ITER_REST */, int32(126) /* Token.OBJECT_REST */, int32(101) /* Token.CAST */ :
		NodeUtil_GetLhsNodesHelper(n.GetFirstChild(), consumer)
	case int32(118) /* Token.IMPORT_SPEC */, int32(129) /* Token.COMPUTED_PROP */, int32(100) /* Token.STRING_KEY */ :
		NodeUtil_GetLhsNodesHelper(n.GetLastChild(), consumer)
	case int32(32) /* Token.NAME */, int32(119) /* Token.IMPORT_STAR */ :
		consumer.Accept(j2gort.BoxNode(n))
	case int32(26) /* Token.GETPROP */, int32(27) /* Token.GETELEM */ :
		consumer.Accept(j2gort.BoxNode(n))
	case int32(88) /* Token.EMPTY */ :
	case int32(81) /* Token.FOR_IN */, int32(112) /* Token.FOR_OF */, int32(113) /* Token.FOR_AWAIT_OF */ :
		NodeUtil_GetLhsNodesHelper(n.GetFirstChild(), consumer)
	default:
		if NodeUtil_IsAssignmentOp(n) {
			NodeUtil_GetLhsNodesHelper(n.GetFirstChild(), consumer)
		} else {
			panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Invalid node in lhs: "+j2gort.StrVal(n))))
		}
	}
}

// Java static method: NodeUtil.visitLhsNodesInNode(com.google.javascript.rhino.Node assigningParent, java.util.function.Consumer<com.google.javascript.rhino.Node> consumer)
func NodeUtil_VisitLhsNodesInNode(assigningParent pkg_rhino.Node, consumer j2gort.Ext_java_util_function_Consumer) {
	_ = assigningParent
	_ = consumer
	if !(((((((NodeUtil_IsNameDeclaration(assigningParent) || assigningParent.IsParamList()) || NodeUtil_IsAssignmentOp(assigningParent)) || assigningParent.IsCatch()) || assigningParent.IsDestructuringLhs()) || assigningParent.IsDefaultValue()) || assigningParent.IsImport()) || NodeUtil_IsEnhancedFor(assigningParent)) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(assigningParent))
	}
	NodeUtil_GetLhsNodesHelper(assigningParent, consumer)
}

// Java static method: NodeUtil.visitLhsNodesInDestructuringPattern(com.google.javascript.rhino.Node destructuringPattern, java.util.function.Consumer<com.google.javascript.rhino.Node> consumer)
func NodeUtil_VisitLhsNodesInDestructuringPattern(destructuringPattern pkg_rhino.Node, consumer j2gort.Ext_java_util_function_Consumer) {
	_ = destructuringPattern
	_ = consumer
	if !(destructuringPattern.IsDestructuringPattern()) {
		j2gort.CheckArgumentStr(false, "Must be a destructuring pattern node.")
	}
	NodeUtil_GetLhsNodesHelper(destructuringPattern, consumer)
}

// Java static method: NodeUtil.isObjectDefinePropertiesDefinition(com.google.javascript.rhino.Node n)
func NodeUtil_IsObjectDefinePropertiesDefinition(n pkg_rhino.Node) bool {
	_ = n
	if !(n.IsCall()) || !(n.HasXChildren(int32(3))) {
		return false
	}
	var getprop pkg_rhino.Node = n.GetFirstChild()
	_ = getprop
	if !(getprop.IsGetProp()) {
		return false
	}
	return ((j2gort.CheckNotNullString(getprop.GetString()) == "defineProperties") && NodeUtil_IsKnownGlobalObjectReference(getprop.GetFirstChild()))
}

// Java static method: NodeUtil.isObjectCreateWithPropertiesDefinition(com.google.javascript.rhino.Node n)
func NodeUtil_IsObjectCreateWithPropertiesDefinition(n pkg_rhino.Node) bool {
	_ = n
	if !(n.IsCall()) || !(n.HasXChildren(int32(3))) {
		return false
	}
	var getprop pkg_rhino.Node = n.GetFirstChild()
	_ = getprop
	return ((getprop.IsGetProp() && (j2gort.CheckNotNullString(getprop.GetString()) == "create")) && NodeUtil_IsKnownGlobalObjectReference(getprop.GetFirstChild()))
}

// Java static method: NodeUtil.isKnownGlobalObjectReference(com.google.javascript.rhino.Node n)
func NodeUtil_IsKnownGlobalObjectReference(n pkg_rhino.Node) bool {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(32) /* Token.NAME */ :
		return (j2gort.CheckNotNullString(n.GetString()) == "Object")
	case int32(26) /* Token.GETPROP */ :
		return (StaticGet_NodeUtil_GLOBAL_OBJECT().Matches(n) || StaticGet_NodeUtil_GLOBAL_OBJECT_MANGLED().Matches(n))
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java static method: NodeUtil.isObjectDefinePropertyDefinition(com.google.javascript.rhino.Node n)
func NodeUtil_IsObjectDefinePropertyDefinition(n pkg_rhino.Node) bool {
	_ = n
	if !(n.IsCall()) || !(n.HasXChildren(int32(4))) {
		return false
	}
	var getprop pkg_rhino.Node = n.GetFirstChild()
	_ = getprop
	if !(getprop.IsGetProp()) {
		return false
	}
	return ((j2gort.CheckNotNullString(getprop.GetString()) == "defineProperty") && NodeUtil_IsKnownGlobalObjectReference(getprop.GetFirstChild()))
}

// Java static method: NodeUtil.getObjectDefinedPropertiesKeys(com.google.javascript.rhino.Node definePropertiesCall)
func NodeUtil_GetObjectDefinedPropertiesKeys(definePropertiesCall pkg_rhino.Node) j2gort.Ext_java_lang_Iterable {
	_ = definePropertiesCall
	if !(NodeUtil_IsObjectDefinePropertiesDefinition(definePropertiesCall)) {
		j2gort.CheckArgument(false)
	}
	var properties []pkg_rhino.Node = []pkg_rhino.Node{}
	_ = properties
	var objectLiteral pkg_rhino.Node = definePropertiesCall.GetLastChild()
	_ = objectLiteral
	{
		var key pkg_rhino.Node = objectLiteral.GetFirstChild()
		_ = key
		for ; key != 0; key = key.GetNext() {
			if !(key.IsStringKey()) {
				continue
			}
			properties = append(properties, key)
		}
	}
	return j2gort.WrapTypedArrayList(properties)
}

// Java static method: NodeUtil.isPrototypePropertyDeclaration(com.google.javascript.rhino.Node n)
func NodeUtil_IsPrototypePropertyDeclaration(n pkg_rhino.Node) bool {
	_ = n
	return (NodeUtil_IsExprAssign(n) && NodeUtil_IsPrototypeProperty(n.GetFirstFirstChild()))
}

// Java static method: NodeUtil.isPrototypeProperty(com.google.javascript.rhino.Node n)
func NodeUtil_IsPrototypeProperty(n pkg_rhino.Node) bool {
	_ = n
	if !(n.IsGetProp()) {
		return false
	}
	var recv pkg_rhino.Node = n.GetFirstChild()
	_ = recv
	return (recv.IsGetProp() && (j2gort.CheckNotNullString(recv.GetString()) == "prototype"))
}

// Java static method: NodeUtil.isPrototypeMethod(com.google.javascript.rhino.Node n)
func NodeUtil_IsPrototypeMethod(n pkg_rhino.Node) bool {
	_ = n
	if !(n.IsFunction()) {
		return false
	}
	var assignNode pkg_rhino.Node = n.GetParent()
	_ = assignNode
	if !(assignNode.IsAssign()) {
		return false
	}
	return NodeUtil_IsPrototypePropertyDeclaration(assignNode.GetParent())
}

// Java static method: NodeUtil.isPrototypeAssignment(com.google.javascript.rhino.Node getProp)
func NodeUtil_IsPrototypeAssignment(getProp pkg_rhino.Node) bool {
	_ = getProp
	if !(getProp.IsGetProp()) {
		return false
	}
	var parent pkg_rhino.Node = getProp.GetParent()
	_ = parent
	return ((parent.IsAssign() && getProp.IsFirstChildOf(parent)) && (j2gort.CheckNotNullString(getProp.GetString()) == "prototype"))
}

// Java static method: NodeUtil.isPropertyTest(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.rhino.Node propAccess)
func NodeUtil_IsPropertyTest(compiler AbstractCompiler, propAccess pkg_rhino.Node) bool {
	_ = compiler
	_ = propAccess
	var parent pkg_rhino.Node = propAccess.GetParent()
	_ = parent
	switch parent.TokenSwitchOrd_() {
	case int32(28) /* Token.CALL */ :
		return ((parent.GetFirstChild() != propAccess) && compiler.GetCodingConvention().IsPropertyTestFunction(parent))
	case int32(74) /* Token.IF */, int32(78) /* Token.WHILE */, int32(79) /* Token.DO */, int32(80) /* Token.FOR */, int32(81) /* Token.FOR_IN */ :
		return (NodeUtil_GetConditionExpression(parent) == propAccess)
	case int32(45) /* Token.INSTANCEOF */, int32(25) /* Token.TYPEOF */, int32(69) /* Token.AND */, int32(68) /* Token.OR */, int32(70) /* Token.COALESCE */, int32(29) /* Token.OPTCHAIN_GETPROP */ :
		return true
	case int32(5) /* Token.NE */, int32(41) /* Token.SHNE */ :
		var other pkg_rhino.Node
		_ = other
		if parent.GetFirstChild() == propAccess {
			other = parent.GetSecondChild()
		} else {
			other = parent.GetFirstChild()
		}
		return (NodeUtil_IsUndefined(other) || (parent.IsNE() && other.IsNull()))
	case int32(31) /* Token.OPTCHAIN_CALL */, int32(30) /* Token.OPTCHAIN_GETELEM */, int32(67) /* Token.HOOK */ :
		return (parent.GetFirstChild() == propAccess)
	case int32(19) /* Token.NOT */ :
		return (parent.GetParent().IsOr() && (parent.GetParent().GetFirstChild() == parent))
	case int32(101) /* Token.CAST */ :
		return NodeUtil_IsPropertyTest(compiler, parent)
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java static method: NodeUtil.isPropertyAbsenceTest(com.google.javascript.rhino.Node propAccess)
func NodeUtil_IsPropertyAbsenceTest(propAccess pkg_rhino.Node) bool {
	_ = propAccess
	var parent pkg_rhino.Node = propAccess.GetParent()
	_ = parent
	switch parent.TokenSwitchOrd_() {
	case int32(4) /* Token.EQ */, int32(40) /* Token.SHEQ */ :
		var other pkg_rhino.Node
		_ = other
		if parent.GetFirstChild() == propAccess {
			other = parent.GetSecondChild()
		} else {
			other = parent.GetFirstChild()
		}
		return (NodeUtil_IsUndefined(other) || (parent.IsEQ() && other.IsNull()))
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java static method: NodeUtil.getPrototypeClassName(com.google.javascript.rhino.Node qName)
func NodeUtil_GetPrototypeClassName(qName pkg_rhino.Node) pkg_rhino.Node {
	_ = qName
	if !(qName.IsGetProp()) {
		return pkg_rhino.Node(0)
	}
	if j2gort.CheckNotNullString(qName.GetString()) == "prototype" {
		return qName.GetFirstChild()
	}
	var recv pkg_rhino.Node = qName.GetFirstChild()
	_ = recv
	if recv.IsGetProp() && (j2gort.CheckNotNullString(recv.GetString()) == "prototype") {
		return recv.GetFirstChild()
	}
	return pkg_rhino.Node(0)
}

// Java static method: NodeUtil.getPrototypePropertyName(com.google.javascript.rhino.Node qName)
func NodeUtil_GetPrototypePropertyName(qName pkg_rhino.Node) string {
	_ = qName
	var qNameStr string = qName.GetQualifiedName()
	_ = qNameStr
	var prototypeIdx int32 = j2gort.JStrLastIndexOf(qNameStr, ".prototype.")
	_ = prototypeIdx
	var memberIndex int32 = ((prototypeIdx + j2gort.JStrLen(".prototype")) + int32(1))
	_ = memberIndex
	return j2gort.JStrSubFrom(qNameStr, memberIndex)
}

// Java static method: NodeUtil.newUndefinedNode(com.google.javascript.rhino.Node srcReferenceNode)
func NodeUtil_NewUndefinedNode(srcReferenceNode pkg_rhino.Node) pkg_rhino.Node {
	_ = srcReferenceNode
	var node pkg_rhino.Node = pkg_rhino.IR_VoidNode(pkg_rhino.IR_Number(float64(int32(0))))
	_ = node
	if srcReferenceNode != 0 {
		node.SrcrefTree(srcReferenceNode)
	}
	return node
}

// Java static method: NodeUtil.newVarNode(java.lang.String name, com.google.javascript.rhino.Node value)
func NodeUtil_NewVarNode_2_string_Node(name string, value pkg_rhino.Node) pkg_rhino.Node {
	_ = name
	_ = value
	var lhs pkg_rhino.Node = pkg_rhino.IR_Name(name)
	_ = lhs
	if value != 0 {
		lhs.Srcref(value)
	}
	return NodeUtil_NewVarNode_2_Node_Node(lhs, value)
}

// Java static method: NodeUtil.newVarNode(com.google.javascript.rhino.Node lhs, com.google.javascript.rhino.Node value)
func NodeUtil_NewVarNode_2_Node_Node(lhs pkg_rhino.Node, value pkg_rhino.Node) pkg_rhino.Node {
	_ = lhs
	_ = value
	if lhs.IsDestructuringPattern() {
		j2gort.CheckNotNullNode(value)
		return pkg_rhino.IR_Var_1(pkg_rhino.NewNode_3(pkg_rhino.StaticGet_Token_DESTRUCTURING_LHS(), lhs, value).Srcref(lhs)).Srcref(lhs)
	} else {
		if !(lhs.IsName() && !(lhs.HasChildren())) {
			j2gort.CheckState(false)
		}
		if value != 0 {
			lhs.AddChildToBack(value)
		}
		return pkg_rhino.IR_Var_1(lhs).Srcref(lhs)
	}
}

// Java static method: NodeUtil.emptyFunction()
func NodeUtil_EmptyFunction() pkg_rhino.Node {
	return pkg_rhino.IR_Function(pkg_rhino.IR_Name(""), pkg_rhino.IR_ParamList([]pkg_rhino.Node{}), pkg_rhino.IR_Block_0())
}

// Java static method: NodeUtil.isShallowStatementTree(com.google.javascript.rhino.Node n)
func NodeUtil_IsShallowStatementTree(n pkg_rhino.Node) bool {
	_ = n
	return ((((n == 0) || NodeUtil_IsControlStructure(n)) || NodeUtil_IsStatementBlock(n)) || n.IsSwitchBody())
}

// Java static method: NodeUtil.getNodeTypeReferenceCount(com.google.javascript.rhino.Node node, com.google.javascript.rhino.Token type, com.google.common.base.Predicate<com.google.javascript.rhino.Node> traverseChildrenPred)
func NodeUtil_GetNodeTypeReferenceCount(node pkg_rhino.Node, type_ pkg_rhino.Token, traverseChildrenPred j2gort.Ext_com_google_common_base_Predicate) int32 {
	_ = node
	_ = type_
	_ = traverseChildrenPred
	if j2gort.RefEquals(traverseChildrenPred, j2gort.Ext_com_google_common_base_Predicates_AlwaysTrue()) {
		return NodeUtil_CountNodeTypeAll(node, type_)
	}
	if j2gort.RefEquals(traverseChildrenPred, StaticGet_NodeUtil_MATCH_NOT_FUNCTION()) {
		return NodeUtil_CountNodeTypeNotFunction(node, type_)
	}
	if j2gort.RefEquals(traverseChildrenPred, StaticGet_NodeUtil_MATCH_ANYTHING_BUT_NON_ARROW_FUNCTION()) {
		return NodeUtil_CountNodeTypeExcludingNonArrowFunction(node, type_)
	}
	return NodeUtil_CountNodeTypeWithPred(node, type_, traverseChildrenPred)
}

// Java static method: NodeUtil.countNodeTypeAll(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Token type)
func NodeUtil_CountNodeTypeAll(n pkg_rhino.Node, type_ pkg_rhino.Token) int32 {
	_ = n
	_ = type_
	var total int32
	_ = total
	if n.GetToken() == type_ {
		total = int32(1)
	} else {
		total = int32(0)
	}
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			total = (total + NodeUtil_CountNodeTypeAll(c, type_))
		}
	}
	return total
}

// Java static method: NodeUtil.countNodeTypeNotFunction(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Token type)
func NodeUtil_CountNodeTypeNotFunction(n pkg_rhino.Node, type_ pkg_rhino.Token) int32 {
	_ = n
	_ = type_
	var total int32
	_ = total
	if n.GetToken() == type_ {
		total = int32(1)
	} else {
		total = int32(0)
	}
	if !(n.IsFunction()) {
		{
			var c pkg_rhino.Node = n.GetFirstChild()
			_ = c
			for ; c != 0; c = c.GetNext() {
				total = (total + NodeUtil_CountNodeTypeNotFunction(c, type_))
			}
		}
	}
	return total
}

// Java static method: NodeUtil.countNodeTypeExcludingNonArrowFunction(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Token type)
func NodeUtil_CountNodeTypeExcludingNonArrowFunction(n pkg_rhino.Node, type_ pkg_rhino.Token) int32 {
	_ = n
	_ = type_
	var total int32
	_ = total
	if n.GetToken() == type_ {
		total = int32(1)
	} else {
		total = int32(0)
	}
	if !(NodeUtil_IsNonArrowFunction(n)) {
		{
			var c pkg_rhino.Node = n.GetFirstChild()
			_ = c
			for ; c != 0; c = c.GetNext() {
				total = (total + NodeUtil_CountNodeTypeExcludingNonArrowFunction(c, type_))
			}
		}
	}
	return total
}

// Java static method: NodeUtil.countNodeTypeWithPred(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Token type, com.google.common.base.Predicate<com.google.javascript.rhino.Node> traverseChildrenPred)
func NodeUtil_CountNodeTypeWithPred(n pkg_rhino.Node, type_ pkg_rhino.Token, traverseChildrenPred j2gort.Ext_com_google_common_base_Predicate) int32 {
	_ = n
	_ = type_
	_ = traverseChildrenPred
	var total int32
	_ = total
	if n.GetToken() == type_ {
		total = int32(1)
	} else {
		total = int32(0)
	}
	if j2gort.PredTestNode(traverseChildrenPred, uint32(n)) {
		{
			var c pkg_rhino.Node = n.GetFirstChild()
			_ = c
			for ; c != 0; c = c.GetNext() {
				total = (total + NodeUtil_CountNodeTypeWithPred(c, type_, traverseChildrenPred))
			}
		}
	}
	return total
}

// Java static method: NodeUtil.isNameReferenced(com.google.javascript.rhino.Node node, java.lang.String name, com.google.common.base.Predicate<com.google.javascript.rhino.Node> traverseChildrenPred)
func NodeUtil_IsNameReferenced_3(node pkg_rhino.Node, name string, traverseChildrenPred j2gort.Ext_com_google_common_base_Predicate) bool {
	_ = node
	_ = name
	_ = traverseChildrenPred
	if j2gort.RefEquals(traverseChildrenPred, j2gort.Ext_com_google_common_base_Predicates_AlwaysTrue()) {
		return NodeUtil_IsNameReferenced_2(node, name)
	}
	if j2gort.RefEquals(traverseChildrenPred, StaticGet_NodeUtil_MATCH_ANYTHING_BUT_NON_ARROW_FUNCTION()) {
		return NodeUtil_IsNameReferencedExcludingNonArrowFunction(node, name)
	}
	if j2gort.RefEquals(traverseChildrenPred, StaticGet_NodeUtil_MATCH_NOT_FUNCTION()) {
		return NodeUtil_IsNameReferencedNotFunction(node, name)
	}
	return NodeUtil_IsNameReferencedWithPred(node, name, traverseChildrenPred)
}

// Java static method: NodeUtil.isNameReferencedExcludingNonArrowFunction(com.google.javascript.rhino.Node node, java.lang.String name)
func NodeUtil_IsNameReferencedExcludingNonArrowFunction(node pkg_rhino.Node, name string) bool {
	_ = node
	_ = name
	if node.IsName() && (j2gort.CheckNotNullString(node.GetString()) == name) {
		return true
	}
	if NodeUtil_IsNonArrowFunction(node) {
		return false
	}
	{
		var c pkg_rhino.Node = node.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			if NodeUtil_IsNameReferencedExcludingNonArrowFunction(c, name) {
				return true
			}
		}
	}
	return false
}

// Java static method: NodeUtil.isNameReferencedNotFunction(com.google.javascript.rhino.Node node, java.lang.String name)
func NodeUtil_IsNameReferencedNotFunction(node pkg_rhino.Node, name string) bool {
	_ = node
	_ = name
	if node.IsName() && (j2gort.CheckNotNullString(node.GetString()) == name) {
		return true
	}
	if node.IsFunction() {
		return false
	}
	{
		var c pkg_rhino.Node = node.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			if NodeUtil_IsNameReferencedNotFunction(c, name) {
				return true
			}
		}
	}
	return false
}

// Java static method: NodeUtil.isNameReferencedWithPred(com.google.javascript.rhino.Node node, java.lang.String name, com.google.common.base.Predicate<com.google.javascript.rhino.Node> traverseChildrenPred)
func NodeUtil_IsNameReferencedWithPred(node pkg_rhino.Node, name string, traverseChildrenPred j2gort.Ext_com_google_common_base_Predicate) bool {
	_ = node
	_ = name
	_ = traverseChildrenPred
	if node.IsName() && (j2gort.CheckNotNullString(node.GetString()) == name) {
		return true
	}
	if !(j2gort.PredTestNode(traverseChildrenPred, uint32(node))) {
		return false
	}
	{
		var c pkg_rhino.Node = node.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			if NodeUtil_IsNameReferencedWithPred(c, name, traverseChildrenPred) {
				return true
			}
		}
	}
	return false
}

// Java static method: NodeUtil.isNameReferenced(com.google.javascript.rhino.Node node, java.lang.String name)
func NodeUtil_IsNameReferenced_2(node pkg_rhino.Node, name string) bool {
	_ = node
	_ = name
	if node.IsName() && (j2gort.CheckNotNullString(node.GetString()) == name) {
		return true
	}
	{
		var c pkg_rhino.Node = node.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			if NodeUtil_IsNameReferenced_2(c, name) {
				return true
			}
		}
	}
	return false
}

// Java static method: NodeUtil.getNameReferenceCount(com.google.javascript.rhino.Node node, java.lang.String name)
func NodeUtil_GetNameReferenceCount(node pkg_rhino.Node, name string) int32 {
	_ = node
	_ = name
	var total int32
	_ = total
	if node.IsName() && (j2gort.CheckNotNullString(node.GetString()) == name) {
		total = int32(1)
	} else {
		total = int32(0)
	}
	{
		var c pkg_rhino.Node = node.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			total = (total + NodeUtil_GetNameReferenceCount(c, name))
		}
	}
	return total
}

// Java static method: NodeUtil.has(com.google.javascript.rhino.Node node, com.google.common.base.Predicate<com.google.javascript.rhino.Node> pred, com.google.common.base.Predicate<com.google.javascript.rhino.Node> traverseChildrenPred)
func NodeUtil_Has(node pkg_rhino.Node, pred j2gort.Ext_com_google_common_base_Predicate, traverseChildrenPred j2gort.Ext_com_google_common_base_Predicate) bool {
	_ = node
	_ = pred
	_ = traverseChildrenPred
	if j2gort.PredTestNode(pred, uint32(node)) {
		return true
	}
	if !(j2gort.PredTestNode(traverseChildrenPred, uint32(node))) {
		return false
	}
	{
		var c pkg_rhino.Node = node.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			if NodeUtil_Has(c, pred, traverseChildrenPred) {
				return true
			}
		}
	}
	return false
}
func nodeUtilHasPredAlwaysTrueFn_(node pkg_rhino.Node, pred func(uint32) bool) bool {
	if pred(uint32(node)) {
		return true
	}
	for c := node.GetFirstChild(); c != 0; c = c.GetNext() {
		if nodeUtilHasPredAlwaysTrueFn_(c, pred) {
			return true
		}
	}
	return false
}
func nodeUtilHasPredFn_(node pkg_rhino.Node, pred func(uint32) bool, traverseChildrenPred j2gort.Ext_com_google_common_base_Predicate) bool {
	if pred(uint32(node)) {
		return true
	}
	if !j2gort.PredTestNode(traverseChildrenPred, uint32(node)) {
		return false
	}
	for c := node.GetFirstChild(); c != 0; c = c.GetNext() {
		if nodeUtilHasPredFn_(c, pred, traverseChildrenPred) {
			return true
		}
	}
	return false
}

// Java static method: NodeUtil.findPreorder(com.google.javascript.rhino.Node node, com.google.common.base.Predicate<com.google.javascript.rhino.Node> pred, com.google.common.base.Predicate<com.google.javascript.rhino.Node> traverseChildrenPred)
func NodeUtil_FindPreorder(node pkg_rhino.Node, pred j2gort.Ext_com_google_common_base_Predicate, traverseChildrenPred j2gort.Ext_com_google_common_base_Predicate) pkg_rhino.Node {
	_ = node
	_ = pred
	_ = traverseChildrenPred
	if j2gort.PredTestNode(pred, uint32(node)) {
		return node
	}
	if !(j2gort.PredTestNode(traverseChildrenPred, uint32(node))) {
		return pkg_rhino.Node(0)
	}
	{
		var c pkg_rhino.Node = node.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			var result pkg_rhino.Node = NodeUtil_FindPreorder(c, pred, traverseChildrenPred)
			_ = result
			if result != 0 {
				return result
			}
		}
	}
	return pkg_rhino.Node(0)
}

// Java static method: NodeUtil.getCount(com.google.javascript.rhino.Node n, com.google.common.base.Predicate<com.google.javascript.rhino.Node> pred, com.google.common.base.Predicate<com.google.javascript.rhino.Node> traverseChildrenPred)
func NodeUtil_GetCount(n pkg_rhino.Node, pred j2gort.Ext_com_google_common_base_Predicate, traverseChildrenPred j2gort.Ext_com_google_common_base_Predicate) int32 {
	_ = n
	_ = pred
	_ = traverseChildrenPred
	var total int32 = int32(0)
	_ = total
	if j2gort.PredTestNode(pred, uint32(n)) {
		total = total + 1
	}
	if j2gort.PredTestNode(traverseChildrenPred, uint32(n)) {
		{
			var c pkg_rhino.Node = n.GetFirstChild()
			_ = c
			for ; c != 0; c = c.GetNext() {
				total = (total + NodeUtil_GetCount(c, pred, traverseChildrenPred))
			}
		}
	}
	return total
}

// Java static method: NodeUtil.visitPreOrder(com.google.javascript.rhino.Node node, com.google.javascript.jscomp.NodeUtil.Visitor visitor)
func NodeUtil_VisitPreOrder_2(node pkg_rhino.Node, visitor NodeUtil_Visitor) {
	_ = node
	_ = visitor
	NodeUtil_VisitPreOrder_3(node, visitor, j2gort.Ext_com_google_common_base_Predicates_AlwaysTrue())
}

// Java static method: NodeUtil.visitPreOrder(com.google.javascript.rhino.Node node, com.google.javascript.jscomp.NodeUtil.Visitor visitor, com.google.common.base.Predicate<com.google.javascript.rhino.Node> traverseChildrenPred)
func NodeUtil_VisitPreOrder_3(node pkg_rhino.Node, visitor NodeUtil_Visitor, traverseChildrenPred j2gort.Ext_com_google_common_base_Predicate) {
	_ = node
	_ = visitor
	_ = traverseChildrenPred
	visitor.Visit(node)
	if j2gort.PredTestNode(traverseChildrenPred, uint32(node)) {
		{
			var c pkg_rhino.Node = node.GetFirstChild()
			_ = c
			for ; c != 0; c = c.GetNext() {
				NodeUtil_VisitPreOrder_3(c, visitor, traverseChildrenPred)
			}
		}
	}
}

// Java static method: NodeUtil.visitPostOrder(com.google.javascript.rhino.Node node, com.google.javascript.jscomp.NodeUtil.Visitor visitor)
func NodeUtil_VisitPostOrder_2(node pkg_rhino.Node, visitor NodeUtil_Visitor) {
	_ = node
	_ = visitor
	NodeUtil_VisitPostOrder_3(node, visitor, j2gort.Ext_com_google_common_base_Predicates_AlwaysTrue())
}

// Java static method: NodeUtil.visitPostOrder(com.google.javascript.rhino.Node node, com.google.javascript.jscomp.NodeUtil.Visitor visitor, com.google.common.base.Predicate<com.google.javascript.rhino.Node> traverseChildrenPred)
func NodeUtil_VisitPostOrder_3(node pkg_rhino.Node, visitor NodeUtil_Visitor, traverseChildrenPred j2gort.Ext_com_google_common_base_Predicate) {
	_ = node
	_ = visitor
	_ = traverseChildrenPred
	if j2gort.PredTestNode(traverseChildrenPred, uint32(node)) {
		{
			var c pkg_rhino.Node = node.GetFirstChild()
			_ = c
			for c != 0 {
				var next pkg_rhino.Node = c.GetNext()
				_ = next
				NodeUtil_VisitPostOrder_3(c, visitor, traverseChildrenPred)
				c = next
			}
		}
	}
	visitor.Visit(node)
}

// Java static method: NodeUtil.preOrderIterable(com.google.javascript.rhino.Node root, com.google.common.base.Predicate<com.google.javascript.rhino.Node> travserseNodePredicate)
func NodeUtil_PreOrderIterable_2(root pkg_rhino.Node, travserseNodePredicate j2gort.Ext_com_google_common_base_Predicate) j2gort.Ext_java_lang_Iterable {
	_ = root
	_ = travserseNodePredicate
	return j2gort.NewLambda_Ext_java_lang_Iterable(func() any {
		return NewNodeUtil_PreOrderIterator(root, travserseNodePredicate)
	})
}

// Java static method: NodeUtil.preOrderIterable(com.google.javascript.rhino.Node root)
func NodeUtil_PreOrderIterable_1(root pkg_rhino.Node) j2gort.Ext_java_lang_Iterable {
	_ = root
	return NodeUtil_PreOrderIterable_2(root, j2gort.Ext_com_google_common_base_Predicates_AlwaysTrue())
}

// Java static method: NodeUtil.isExportFrom(com.google.javascript.rhino.Node n)
func NodeUtil_IsExportFrom(n pkg_rhino.Node) bool {
	_ = n
	if !(n.IsExport()) {
		j2gort.CheckArgument(false)
	}
	return n.HasTwoChildren()
}

// Java static method: NodeUtil.hasFinally(com.google.javascript.rhino.Node n)
func NodeUtil_HasFinally(n pkg_rhino.Node) bool {
	_ = n
	if !(n.IsTry()) {
		j2gort.CheckArgument(false)
	}
	return n.HasXChildren(int32(3))
}

// Java static method: NodeUtil.getCatchBlock(com.google.javascript.rhino.Node n)
func NodeUtil_GetCatchBlock(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if !(n.IsTry()) {
		j2gort.CheckArgument(false)
	}
	return n.GetSecondChild()
}

// Java static method: NodeUtil.hasCatchHandler(com.google.javascript.rhino.Node n)
func NodeUtil_HasCatchHandler(n pkg_rhino.Node) bool {
	_ = n
	if !(n.IsBlock()) {
		j2gort.CheckArgument(false)
	}
	return (n.HasChildren() && n.GetFirstChild().IsCatch())
}

// Java static method: NodeUtil.getFunctionParameters(com.google.javascript.rhino.Node fnNode)
func NodeUtil_GetFunctionParameters(fnNode pkg_rhino.Node) pkg_rhino.Node {
	_ = fnNode
	if !(fnNode.IsFunction()) {
		j2gort.CheckArgument(false)
	}
	return fnNode.GetSecondChild()
}

// Java static method: NodeUtil.hasNonSimpleParameters(com.google.javascript.rhino.Node fnNode)
func NodeUtil_HasNonSimpleParameters(fnNode pkg_rhino.Node) bool {
	_ = fnNode
	if !(fnNode.IsFunction()) {
		j2gort.CheckArgument(false)
	}
	var paramList pkg_rhino.Node = NodeUtil_GetFunctionParameters(fnNode)
	_ = paramList
	{
		var param pkg_rhino.Node = paramList.GetFirstChild()
		_ = param
		for ; param != 0; param = param.GetNext() {
			if !(param.IsName()) {
				return true
			}
		}
	}
	return false
}

// Java static method: NodeUtil.functionParametersMayHaveSideEffects(com.google.javascript.rhino.Node fnNode, com.google.javascript.jscomp.AstAnalyzer astAnalyzer)
func NodeUtil_FunctionParametersMayHaveSideEffects(fnNode pkg_rhino.Node, astAnalyzer *AstAnalyzer) bool {
	_ = fnNode
	_ = astAnalyzer
	if !(fnNode.IsFunction()) {
		j2gort.CheckArgument(false)
	}
	var paramList pkg_rhino.Node = NodeUtil_GetFunctionParameters(fnNode)
	_ = paramList
	{
		var param pkg_rhino.Node = paramList.GetFirstChild()
		_ = param
		for ; param != 0; param = param.GetNext() {
			if param.IsName() {
				continue
			}
			if param.IsRest() {
				if param.GetFirstChild().IsName() {
					continue
				}
				return true
			}
			if param.IsDefaultValue() {
				var target pkg_rhino.Node = param.GetFirstChild()
				_ = target
				var initializer pkg_rhino.Node = param.GetSecondChild()
				_ = initializer
				if !(target.IsName()) {
					return true
				}
				if astAnalyzer.MayHaveSideEffects(initializer) {
					return true
				}
				continue
			}
			return true
		}
	}
	return false
}

// Java static method: NodeUtil.isConstantVar(com.google.javascript.rhino.Node node, com.google.javascript.jscomp.@org.jspecify.annotations.Nullable Scope scope)
func NodeUtil_IsConstantVar(node pkg_rhino.Node, scope *Scope) bool {
	_ = node
	_ = scope
	if NodeUtil_IsConstantName(node) {
		return true
	}
	if !(node.IsName()) || (scope == nil) {
		return false
	}
	var var_ *Var = j2gort.CastPtr[Var](scope.GetVar(node.GetString()))
	_ = var_
	return ((var_ != nil) && (var_.IsDeclaredOrInferredConst() || var_.IsConst()))
}

// Java static method: NodeUtil.isConstantName(com.google.javascript.rhino.Node node)
func NodeUtil_IsConstantName(node pkg_rhino.Node) bool {
	_ = node
	return node.GetBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME())
}

// Java static method: NodeUtil.isConstantByConvention(com.google.javascript.jscomp.CodingConvention convention, com.google.javascript.rhino.Node node)
func NodeUtil_IsConstantByConvention(convention CodingConvention, node pkg_rhino.Node) bool {
	_ = convention
	_ = node
	if NodeUtil_IsNormalOrOptChainGetProp(node) {
		return convention.IsConstantKey(node.GetString())
	} else if NodeUtil_MayBeObjectLitKey(node) {
		return convention.IsConstantKey(node.GetString())
	} else if node.IsName() {
		return convention.IsConstant(node.GetString())
	}
	return false
}

// Java static method: NodeUtil.isConstantDeclaration(com.google.javascript.rhino.JSDocInfo info, com.google.javascript.rhino.Node node)
func NodeUtil_IsConstantDeclaration(info *pkg_rhino.JSDocInfo, node pkg_rhino.Node) bool {
	_ = info
	_ = node
	if ((((NodeUtil_IsObjectLitKey(node) || (node.GetParent().IsAssign() && node.IsFirstChildOf(node.GetParent()))) || (node.GetParent().IsExprResult() && NodeUtil_IsNormalGet(node))) || node.IsMemberFieldDef()) || node.IsComputedFieldDef()) || node.IsComputedProp() {
		return ((info != nil) && info.IsConstant())
	}
	if !(node.IsName()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(node))
	}
	var declaringParent pkg_rhino.Node = NodeUtil_GetDeclaringParent(node)
	_ = declaringParent
	if declaringParent.IsConst() {
		return true
	} else if (info != nil) && info.IsConstant() {
		return true
	}
	return node.IsInferredConstantVar()
}

// Java static method: NodeUtil.functionHasInlineJsdocs(com.google.javascript.rhino.Node fn)
func NodeUtil_FunctionHasInlineJsdocs(fn pkg_rhino.Node) bool {
	_ = fn
	if !(fn.IsFunction()) {
		return false
	}
	if fn.GetFirstChild().GetJSDocInfo() != nil {
		return true
	}
	var param pkg_rhino.Node = fn.GetSecondChild().GetFirstChild()
	_ = param
	for param != 0 {
		if param.GetJSDocInfo() != nil {
			return true
		}
		param = param.GetNext()
	}
	return false
}

// Java static method: NodeUtil.getSourceName(com.google.javascript.rhino.Node n)
func NodeUtil_GetSourceName(n pkg_rhino.Node) string {
	_ = n
	var sourceName string = j2gort.NullString
	_ = sourceName
	for (sourceName == j2gort.NullString) && (n != 0) {
		sourceName = n.GetSourceFileName()
		n = n.GetParent()
	}
	return sourceName
}

// Java static method: NodeUtil.getSourceFile(com.google.javascript.rhino.Node n)
func NodeUtil_GetSourceFile(n pkg_rhino.Node) pkg_rhino.StaticSourceFile {
	_ = n
	var sourceName pkg_rhino.StaticSourceFile = nil
	_ = sourceName
	for (sourceName == nil || sourceName.IsJavaNil_()) && (n != 0) {
		sourceName = n.GetStaticSourceFile()
		n = n.GetParent()
	}
	return sourceName
}

// Java static method: NodeUtil.getInputId(com.google.javascript.rhino.Node n)
func NodeUtil_GetInputId(n pkg_rhino.Node) *pkg_rhino.InputId {
	_ = n
	for (n != 0) && !(n.IsScript()) {
		n = n.GetParent()
	}
	if (n != 0) && n.IsScript() {
		return n.GetInputId()
	} else {
		return nil
	}
}

// Java static method: NodeUtil.newCallNode(com.google.javascript.rhino.Node callTarget, com.google.javascript.rhino.Node[] parameters)
func NodeUtil_NewCallNode(callTarget pkg_rhino.Node, parameters []pkg_rhino.Node) pkg_rhino.Node {
	_ = callTarget
	_ = parameters
	var isFreeCall bool = !(NodeUtil_IsNormalGet(callTarget))
	_ = isFreeCall
	var call pkg_rhino.Node = pkg_rhino.IR_Call(callTarget, []pkg_rhino.Node{})
	_ = call
	call.PutBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL(), isFreeCall)
	for _, el_32 := range j2gort.CheckNotNullSlice(parameters) {
		parameter := el_32
		_ = parameter
		call.AddChildToBack(parameter)
	}
	return call
}

// Java static method: NodeUtil.evaluatesToLocalValue(com.google.javascript.rhino.Node value)
func NodeUtil_EvaluatesToLocalValue(value pkg_rhino.Node) bool {
	_ = value
	switch value.TokenSwitchOrd_() {
	case int32(51) /* Token.ASSIGN */ :
		return NodeUtil_IsImmutableValue(value.GetLastChild())
	case int32(50) /* Token.COMMA */ :
		return NodeUtil_EvaluatesToLocalValue(value.GetLastChild())
	case int32(69) /* Token.AND */, int32(68) /* Token.OR */, int32(70) /* Token.COALESCE */ :
		return (NodeUtil_EvaluatesToLocalValue(value.GetFirstChild()) && NodeUtil_EvaluatesToLocalValue(value.GetLastChild()))
	case int32(67) /* Token.HOOK */ :
		return (NodeUtil_EvaluatesToLocalValue(value.GetSecondChild()) && NodeUtil_EvaluatesToLocalValue(value.GetLastChild()))
	case int32(124) /* Token.DYNAMIC_IMPORT */ :
		return true
	case int32(37) /* Token.THIS */, int32(110) /* Token.SUPER */ :
		return false
	case int32(32) /* Token.NAME */ :
		return NodeUtil_IsImmutableValue(value)
	case int32(27) /* Token.GETELEM */, int32(26) /* Token.GETPROP */, int32(30) /* Token.OPTCHAIN_GETELEM */, int32(29) /* Token.OPTCHAIN_GETPROP */ :
		return false
	case int32(28) /* Token.CALL */, int32(31) /* Token.OPTCHAIN_CALL */ :
		return NodeUtil_IsToStringMethodCall(value)
	case int32(130) /* Token.TAGGED_TEMPLATELIT */ :
		return false
	case int32(23) /* Token.NEW */ :
		return NodeUtil_NewHasLocalResult(value)
	case int32(24) /* Token.DELPROP */, int32(71) /* Token.INC */, int32(72) /* Token.DEC */, int32(105) /* Token.CLASS */, int32(73) /* Token.FUNCTION */, int32(42) /* Token.REGEXP */, int32(88) /* Token.EMPTY */, int32(46) /* Token.ARRAYLIT */, int32(47) /* Token.OBJECTLIT */, int32(131) /* Token.TEMPLATELIT */ :
		return true
	case int32(101) /* Token.CAST */ :
		return NodeUtil_EvaluatesToLocalValue(value.GetFirstChild())
	case int32(127) /* Token.ITER_SPREAD */, int32(128) /* Token.OBJECT_SPREAD */, int32(135) /* Token.NEW_TARGET */ :
		return false
	case int32(114) /* Token.YIELD */, int32(115) /* Token.AWAIT */ :
		return false
	default:
		if NodeUtil_IsLogicalAssignmentOp(value) {
			return false
		}
		if (NodeUtil_IsAssignmentOp(value) || NodeUtil_IsSimpleOperator(value)) || NodeUtil_IsImmutableValue(value) {
			return true
		}
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", j2gort.BoxString(("Unexpected expression node: " + j2gort.StrVal(value) + "\n parent:" + j2gort.StrVal(value.GetParent()))))))
	}
	panic("unreachable switch expression")
}

// Java static method: NodeUtil.mayBeUndefined(com.google.javascript.rhino.Node n)
func NodeUtil_MayBeUndefined(n pkg_rhino.Node) bool {
	_ = n
	return !(NodeUtil_IsDefinedValue(n))
}

// Java static method: NodeUtil.isDefinedValue(com.google.javascript.rhino.Node value)
func NodeUtil_IsDefinedValue(value pkg_rhino.Node) bool {
	_ = value
	switch value.TokenSwitchOrd_() {
	case int32(51) /* Token.ASSIGN */, int32(101) /* Token.CAST */, int32(50) /* Token.COMMA */ :
		return NodeUtil_IsDefinedValue(value.GetLastChild())
	case int32(70) /* Token.COALESCE */ :
		return NodeUtil_IsDefinedValue(value.GetSecondChild())
	case int32(69) /* Token.AND */, int32(68) /* Token.OR */ :
		return (NodeUtil_IsDefinedValue(value.GetFirstChild()) && NodeUtil_IsDefinedValue(value.GetLastChild()))
	case int32(67) /* Token.HOOK */ :
		return (NodeUtil_IsDefinedValue(value.GetSecondChild()) && NodeUtil_IsDefinedValue(value.GetLastChild()))
	case int32(28) /* Token.CALL */, int32(31) /* Token.OPTCHAIN_CALL */, int32(27) /* Token.GETELEM */, int32(26) /* Token.GETPROP */, int32(30) /* Token.OPTCHAIN_GETELEM */, int32(29) /* Token.OPTCHAIN_GETPROP */, int32(130) /* Token.TAGGED_TEMPLATELIT */, int32(37) /* Token.THIS */, int32(114) /* Token.YIELD */, int32(115) /* Token.AWAIT */, int32(87) /* Token.VOID */ :
		return false
	case int32(24) /* Token.DELPROP */, int32(71) /* Token.INC */, int32(72) /* Token.DEC */, int32(105) /* Token.CLASS */, int32(73) /* Token.FUNCTION */, int32(42) /* Token.REGEXP */, int32(88) /* Token.EMPTY */, int32(46) /* Token.ARRAYLIT */, int32(47) /* Token.OBJECTLIT */, int32(131) /* Token.TEMPLATELIT */, int32(35) /* Token.STRINGLIT */, int32(33) /* Token.NUMBER */, int32(34) /* Token.BIGINT */, int32(36) /* Token.NULL */, int32(39) /* Token.TRUE */, int32(38) /* Token.FALSE */, int32(23) /* Token.NEW */ :
		return true
	case int32(133) /* Token.TEMPLATELIT_STRING */ :
		return (value.GetCookedString() != j2gort.NullString)
	case int32(32) /* Token.NAME */ :
		var name string = value.GetString()
		_ = name
		return (("Infinity" == name) || ("NaN" == name))
	default:
		if NodeUtil_IsAssignmentOp(value) || NodeUtil_IsSimpleOperator(value) {
			return true
		}
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", j2gort.BoxString(("Unexpected expression node: " + j2gort.StrVal(value) + "\n parent:" + j2gort.StrVal(value.GetParent()))))))
	}
	panic("unreachable switch expression")
}

// Java static method: NodeUtil.getNthSibling(com.google.javascript.rhino.Node first, int index)
func NodeUtil_GetNthSibling(first pkg_rhino.Node, index int32) pkg_rhino.Node {
	_ = first
	_ = index
	var sibling pkg_rhino.Node = first
	_ = sibling
	for (index != int32(0)) && (sibling != 0) {
		sibling = sibling.GetNext()
		index = index - 1
	}
	return sibling
}

// Java static method: NodeUtil.getArgumentForFunction(com.google.javascript.rhino.Node function, int index)
func NodeUtil_GetArgumentForFunction(function pkg_rhino.Node, index int32) pkg_rhino.Node {
	_ = function
	_ = index
	if !(function.IsFunction()) {
		j2gort.CheckState(false)
	}
	return NodeUtil_GetNthSibling(function.GetSecondChild().GetFirstChild(), index)
}

// Java static method: NodeUtil.getArgumentForCallOrNew(com.google.javascript.rhino.Node call, int index)
func NodeUtil_GetArgumentForCallOrNew(call pkg_rhino.Node, index int32) pkg_rhino.Node {
	_ = call
	_ = index
	if !(NodeUtil_IsCallOrNew(call)) {
		j2gort.CheckState(false)
	}
	return NodeUtil_GetNthSibling(call.GetSecondChild(), index)
}

// Java static method: NodeUtil.isInvocationTarget(com.google.javascript.rhino.Node n)
func NodeUtil_IsInvocationTarget(n pkg_rhino.Node) bool {
	_ = n
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	return (((parent != 0) && (NodeUtil_IsCallOrNew(parent) || parent.IsTaggedTemplateLit())) && (parent.GetFirstChild() == n))
}

// Java static method: NodeUtil.isInvocation(com.google.javascript.rhino.Node n)
func NodeUtil_IsInvocation(n pkg_rhino.Node) bool {
	_ = n
	return (NodeUtil_IsCallOrNew(n) || n.IsTaggedTemplateLit())
}

// Java static method: NodeUtil.isCallOrNewArgument(com.google.javascript.rhino.Node n)
func NodeUtil_IsCallOrNewArgument(n pkg_rhino.Node) bool {
	_ = n
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	return (((parent != 0) && NodeUtil_IsCallOrNew(parent)) && (parent.GetFirstChild() != n))
}

// Java static method: NodeUtil.isToStringMethodCall(com.google.javascript.rhino.Node call)
func NodeUtil_IsToStringMethodCall(call pkg_rhino.Node) bool {
	_ = call
	var getNode pkg_rhino.Node = call.GetFirstChild()
	_ = getNode
	var name string = j2gort.NullString
	_ = name
	if NodeUtil_IsNormalOrOptChainGetProp(getNode) {
		name = getNode.GetString()
	} else if NodeUtil_IsNormalOrOptChainGet(getNode) {
		var propNode pkg_rhino.Node = getNode.GetLastChild()
		_ = propNode
		if !(propNode.IsStringLit()) {
			return false
		}
		name = propNode.GetString()
	} else {
		return false
	}
	return ("toString" == name)
}

// Java static method: NodeUtil.getDeclaredTypeExpression(com.google.javascript.rhino.Node declaration)
func NodeUtil_GetDeclaredTypeExpression(declaration pkg_rhino.Node) pkg_rhino.JSTypeExpression {
	_ = declaration
	if !(declaration.IsName() || declaration.IsStringKey()) {
		j2gort.CheckArgument(false)
	}
	var nameJsdoc *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(declaration)
	_ = nameJsdoc
	if nameJsdoc != nil {
		return nameJsdoc.GetType()
	}
	var parent pkg_rhino.Node = declaration.GetParent()
	_ = parent
	if parent.IsRest() || parent.IsDefaultValue() {
		parent = parent.GetParent()
	}
	if parent.IsParamList() {
		var functionJsdoc *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(parent.GetParent())
		_ = functionJsdoc
		if functionJsdoc != nil {
			return functionJsdoc.GetParameterType(declaration.GetString())
		}
	}
	return (pkg_rhino.JSTypeExpression{})
}

// Java static method: NodeUtil.getBestJSDocInfo(com.google.javascript.rhino.Node n)
func NodeUtil_GetBestJSDocInfo(n pkg_rhino.Node) *pkg_rhino.JSDocInfo {
	_ = n
	var jsdocNode pkg_rhino.Node = NodeUtil_GetBestJSDocInfoNode(n)
	_ = jsdocNode
	if jsdocNode == 0 {
		return nil
	} else {
		return jsdocNode.GetJSDocInfo()
	}
}

// Java static method: NodeUtil.getBestJsDocInfoNodeStrict(com.google.javascript.rhino.Node n)
func NodeUtil_GetBestJsDocInfoNodeStrict(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	var isDeclaredName bool = (n.IsName() && NodeUtil_IsNameDeclaration(n.GetParent()))
	_ = isDeclaredName
	if isDeclaredName {
		return NodeUtil_GetBestJsDocInfoNodeInternal(n)
	}
	var isDeclaredClass bool = NodeUtil_IsClassDeclaration(n)
	_ = isDeclaredClass
	if isDeclaredClass {
		return NodeUtil_GetBestJsDocInfoNodeInternal(n)
	}
	var isDeclaredFunction bool = NodeUtil_IsFunctionDeclaration(n)
	_ = isDeclaredFunction
	if isDeclaredFunction {
		return NodeUtil_GetBestJsDocInfoNodeInternal(n)
	}
	var isRhsClass bool = NodeUtil_IsRhsClass(n)
	_ = isRhsClass
	if isRhsClass {
		return NodeUtil_GetBestJsDocInfoNodeInternal(n)
	}
	var isRhsFunction bool = NodeUtil_IsRhsFunction(n)
	_ = isRhsFunction
	if isRhsFunction {
		return NodeUtil_GetBestJsDocInfoNodeInternal(n)
	}
	var isRhsClassName bool = (n.IsName() && NodeUtil_IsRhsClass(n.GetParent()))
	_ = isRhsClassName
	if isRhsClassName {
		return NodeUtil_GetBestJsDocInfoNodeInternal(n)
	}
	var isRhsFunctionName bool = (n.IsName() && NodeUtil_IsRhsFunction(n.GetParent()))
	_ = isRhsFunctionName
	if isRhsFunctionName {
		return NodeUtil_GetBestJsDocInfoNodeInternal(n)
	}
	var isExport bool = n.IsExport()
	_ = isExport
	if isExport {
		return NodeUtil_GetBestJsDocInfoNodeInternal(n)
	}
	panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Not allowed to get JSDocInfo node for node: "+j2gort.StrVal(n))))
}

// Java static method: NodeUtil.isRhsClass(com.google.javascript.rhino.Node n)
func NodeUtil_IsRhsClass(n pkg_rhino.Node) bool {
	_ = n
	return (NodeUtil_IsNamedClassExpression(n) && (n.GetParent().IsName() || n.GetParent().IsAssign()))
}

// Java static method: NodeUtil.isRhsFunction(com.google.javascript.rhino.Node n)
func NodeUtil_IsRhsFunction(n pkg_rhino.Node) bool {
	_ = n
	return (NodeUtil_IsNamedFunctionExpression(n) && (n.GetParent().IsName() || n.GetParent().IsAssign()))
}

// Java static method: NodeUtil.getBestJSDocInfoNode(com.google.javascript.rhino.Node n)
func NodeUtil_GetBestJSDocInfoNode(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	return NodeUtil_GetBestJsDocInfoNodeInternal(n)
}

// Java static method: NodeUtil.getBestJsDocInfoNodeInternal(com.google.javascript.rhino.Node n)
func NodeUtil_GetBestJsDocInfoNodeInternal(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if n.IsExprResult() {
		return NodeUtil_GetBestJSDocInfoNode(n.GetFirstChild())
	}
	var info *pkg_rhino.JSDocInfo = n.GetJSDocInfo()
	_ = info
	if info == nil {
		var parent pkg_rhino.Node = n.GetParent()
		_ = parent
		if (parent == 0) || n.IsExprResult() {
			return pkg_rhino.Node(0)
		}
		if ((((parent.IsName() || parent.IsExport()) || parent.IsVar()) || parent.IsLet()) || parent.IsConst()) || parent.IsDeclare() {
			return NodeUtil_GetBestJSDocInfoNode(parent)
		} else if parent.IsAssign() {
			return NodeUtil_GetBestJSDocInfoNode(parent)
		} else if NodeUtil_MayBeObjectLitKey(parent) || parent.IsComputedProp() {
			return parent
		} else if (parent.IsFunction() || parent.IsClass()) && (n == parent.GetFirstChild()) {
			return NodeUtil_GetBestJSDocInfoNode(parent)
		} else if NodeUtil_IsNameDeclaration(parent) && parent.HasOneChild() {
			return parent
		} else if (((parent.IsHook() && (parent.GetFirstChild() != n)) || parent.IsOr()) || parent.IsAnd()) || (parent.IsComma() && (parent.GetFirstChild() != n)) {
			return NodeUtil_GetBestJSDocInfoNode(parent)
		}
	}
	return n
}

// Java static method: NodeUtil.getBestLValue(com.google.javascript.rhino.Node n)
func NodeUtil_GetBestLValue(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	if NodeUtil_IsFunctionDeclaration(n) || NodeUtil_IsClassDeclaration(n) {
		return n.GetFirstChild()
	} else if n.IsClassMembers() {
		return NodeUtil_GetBestLValue(parent)
	} else if parent.IsName() {
		return parent
	} else if parent.IsAssign() {
		return parent.GetFirstChild()
	} else if NodeUtil_MayBeObjectLitKey(parent) || parent.IsComputedProp() {
		return parent
	} else if (((parent.IsHook() && (parent.GetFirstChild() != n)) || parent.IsOr()) || parent.IsAnd()) || (parent.IsComma() && (parent.GetFirstChild() != n)) {
		return NodeUtil_GetBestLValue(parent)
	} else if parent.IsCast() {
		return NodeUtil_GetBestLValue(parent)
	}
	return pkg_rhino.Node(0)
}

// Java static method: NodeUtil.getRValueOfLValue(com.google.javascript.rhino.Node n)
func NodeUtil_GetRValueOfLValue(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	switch parent.TokenSwitchOrd_() {
	case int32(51) /* Token.ASSIGN */, int32(52) /* Token.ASSIGN_BITOR */, int32(53) /* Token.ASSIGN_BITXOR */, int32(54) /* Token.ASSIGN_BITAND */, int32(55) /* Token.ASSIGN_LSH */, int32(56) /* Token.ASSIGN_RSH */, int32(57) /* Token.ASSIGN_URSH */, int32(58) /* Token.ASSIGN_ADD */, int32(59) /* Token.ASSIGN_SUB */, int32(60) /* Token.ASSIGN_MUL */, int32(63) /* Token.ASSIGN_EXPONENT */, int32(61) /* Token.ASSIGN_DIV */, int32(62) /* Token.ASSIGN_MOD */, int32(64) /* Token.ASSIGN_OR */, int32(65) /* Token.ASSIGN_AND */, int32(66) /* Token.ASSIGN_COALESCE */, int32(104) /* Token.DESTRUCTURING_LHS */ :
		return n.GetNext()
	case int32(84) /* Token.VAR */, int32(111) /* Token.LET */, int32(97) /* Token.CONST */ :
		return n.GetLastChild()
	case int32(47) /* Token.OBJECTLIT */ :
		if !((((n.IsStringKey() || n.IsComputedProp()) || n.IsMemberFunctionDef()) || n.IsGetterDef()) || n.IsSetterDef()) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(n))
		}
		return n.GetLastChild()
	case int32(106) /* Token.CLASS_MEMBERS */ :
		if !((((n.IsMemberFunctionDef() || n.IsMemberFieldDef()) || n.IsComputedFieldDef()) || n.IsGetterDef()) || n.IsSetterDef()) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(n))
		}
		return n.GetLastChild()
	case int32(73) /* Token.FUNCTION */, int32(105) /* Token.CLASS */ :
		return parent
	default:
	}
	return pkg_rhino.Node(0)
}

// Java static method: NodeUtil.getBestLValueOwner(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node lValue)
func NodeUtil_GetBestLValueOwner(lValue pkg_rhino.Node) pkg_rhino.Node {
	_ = lValue
	if (lValue == 0) || (lValue.GetParent() == 0) {
		return pkg_rhino.Node(0)
	}
	if NodeUtil_MayBeObjectLitKey(lValue) || lValue.IsComputedProp() {
		return NodeUtil_GetBestLValue(lValue.GetParent())
	} else if NodeUtil_IsNormalGet(lValue) {
		return lValue.GetFirstChild()
	}
	return pkg_rhino.Node(0)
}

// Java static method: NodeUtil.getBestLValueName(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node lValue)
func NodeUtil_GetBestLValueName(lValue pkg_rhino.Node) string {
	_ = lValue
	if (lValue == 0) || (lValue.GetParent() == 0) {
		return j2gort.NullString
	}
	if lValue.GetParent().IsClassMembers() && !(lValue.IsComputedProp()) {
		var className string = NodeUtil_GetName(lValue.GetGrandparent())
		_ = className
		if className == j2gort.NullString {
			return j2gort.NullString
		}
		var methodName string = lValue.GetString()
		_ = methodName
		var maybePrototype string
		_ = maybePrototype
		if lValue.IsStaticMember() {
			maybePrototype = "."
		} else {
			maybePrototype = ".prototype."
		}
		return (j2gort.StrValString(className) + j2gort.StrValString(maybePrototype) + j2gort.StrValString(methodName))
	}
	if NodeUtil_MayBeObjectLitKey(lValue) {
		var owner pkg_rhino.Node = NodeUtil_GetBestLValue(lValue.GetParent())
		_ = owner
		if owner != 0 {
			var ownerName string = NodeUtil_GetBestLValueName(owner)
			_ = ownerName
			if ownerName != j2gort.NullString {
				var key string = NodeUtil_GetObjectOrClassLitKeyName(lValue)
				_ = key
				if pkg_rhino.TokenStream_IsJSIdentifier(key) {
					return (j2gort.StrValString(ownerName) + "." + j2gort.StrValString(key))
				} else {
					return j2gort.NullString
				}
			}
		}
		return j2gort.NullString
	}
	return lValue.GetQualifiedName()
}

// Java static method: NodeUtil.getBestLValueRoot(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node lValue)
func NodeUtil_GetBestLValueRoot(lValue pkg_rhino.Node) pkg_rhino.Node {
	_ = lValue
	if lValue == 0 {
		return pkg_rhino.Node(0)
	}
	switch lValue.TokenSwitchOrd_() {
	case int32(100) /* Token.STRING_KEY */ :
		return NodeUtil_GetBestLValueRoot(NodeUtil_GetBestLValue(lValue.GetParent()))
	case int32(26) /* Token.GETPROP */, int32(27) /* Token.GETELEM */ :
		return NodeUtil_GetBestLValueRoot(lValue.GetFirstChild())
	case int32(37) /* Token.THIS */, int32(110) /* Token.SUPER */, int32(32) /* Token.NAME */ :
		return lValue
	default:
		return pkg_rhino.Node(0)
	}
	panic("unreachable switch expression")
}

// Java static method: NodeUtil.isExpressionResultUsed(com.google.javascript.rhino.Node expr)
func NodeUtil_IsExpressionResultUsed(expr pkg_rhino.Node) bool {
	_ = expr
	var parent pkg_rhino.Node = expr.GetParent()
	_ = parent
	switch parent.TokenSwitchOrd_() {
	case int32(90) /* Token.BLOCK */, int32(93) /* Token.EXPR_RESULT */ :
		return false
	case int32(101) /* Token.CAST */ :
		return NodeUtil_IsExpressionResultUsed(parent)
	case int32(67) /* Token.HOOK */, int32(69) /* Token.AND */, int32(68) /* Token.OR */, int32(70) /* Token.COALESCE */ :
		return ((expr == parent.GetFirstChild()) || NodeUtil_IsExpressionResultUsed(parent))
	case int32(50) /* Token.COMMA */ :
		var grandparent pkg_rhino.Node = parent.GetParent()
		_ = grandparent
		if (((grandparent.IsCall() || grandparent.IsTaggedTemplateLit()) && (parent == grandparent.GetFirstChild())) && (expr == parent.GetFirstChild())) && parent.HasTwoChildren() {
			var calledFn pkg_rhino.Node = parent.GetLastChild()
			_ = calledFn
			if calledFn.IsName() && calledFn.MatchesName_1_string("eval") {
				return true
			}
			if NodeUtil_IsNormalOrOptChainGet(calledFn) {
				return true
			}
		}
		return ((expr != parent.GetFirstChild()) && NodeUtil_IsExpressionResultUsed(parent))
	case int32(80) /* Token.FOR */ :
		return (parent.GetSecondChild() == expr)
	default:
	}
	return true
}

// Java static method: NodeUtil.isExecutedExactlyOnce(com.google.javascript.rhino.Node n)
func NodeUtil_IsExecutedExactlyOnce(n pkg_rhino.Node) bool {
	_ = n
inspect:
	for doFirst_33 := true; doFirst_33 || ((func() pkg_rhino.Node { var asgn_34 pkg_rhino.Node = n.GetParent(); n = asgn_34; return asgn_34 }()) != 0); doFirst_33 = false {
		var parent pkg_rhino.Node = n.GetParent()
		_ = parent
		switch parent.TokenSwitchOrd_() {
		case int32(74) /* Token.IF */, int32(67) /* Token.HOOK */, int32(69) /* Token.AND */, int32(68) /* Token.OR */, int32(70) /* Token.COALESCE */, int32(64) /* Token.ASSIGN_OR */, int32(65) /* Token.ASSIGN_AND */, int32(66) /* Token.ASSIGN_COALESCE */, int32(30) /* Token.OPTCHAIN_GETELEM */, int32(31) /* Token.OPTCHAIN_CALL */, int32(134) /* Token.DEFAULT_VALUE */ :
			if parent.GetFirstChild() != n {
				return false
			}
			continue inspect
		case int32(80) /* Token.FOR */, int32(81) /* Token.FOR_IN */, int32(112) /* Token.FOR_OF */, int32(113) /* Token.FOR_AWAIT_OF */ :
			if NodeUtil_IsEnhancedFor(parent) {
				if parent.GetSecondChild() != n {
					return false
				}
			} else {
				if parent.GetFirstChild() != n {
					return false
				}
			}
			continue inspect
		case int32(48) /* Token.TRY */ :
			if !(NodeUtil_HasFinally(parent)) || (parent.GetLastChild() != n) {
				return false
			}
			continue inspect
		case int32(78) /* Token.WHILE */, int32(79) /* Token.DO */, int32(76) /* Token.CASE */, int32(77) /* Token.DEFAULT_CASE */ :
			return false
		case int32(94) /* Token.SCRIPT */, int32(73) /* Token.FUNCTION */ :
			break inspect
		default:
		}
	}
	return true
}

// Java static method: NodeUtil.booleanNode(boolean value)
func NodeUtil_BooleanNode(value bool) pkg_rhino.Node {
	_ = value
	if value {
		return pkg_rhino.IR_TrueNode()
	} else {
		return pkg_rhino.IR_FalseNode()
	}
}

// Java static method: NodeUtil.numberNode(double value, com.google.javascript.rhino.Node srcref)
func NodeUtil_NumberNode(value float64, srcref pkg_rhino.Node) pkg_rhino.Node {
	_ = value
	_ = srcref
	var result pkg_rhino.Node = pkg_rhino.Node(0)
	_ = result
	if math.IsNaN(value) {
		result = pkg_rhino.IR_Name("NaN")
		result.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME(), true)
	} else {
		if math.IsInf(value, 0) {
			result = pkg_rhino.IR_Name("Infinity")
			result.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME(), true)
		} else {
			result = pkg_rhino.IR_Number(math.Abs(value))
		}
		if pkg_base.JSCompDoubles_IsNegative(value) {
			result = pkg_rhino.IR_Neg(result)
		}
	}
	if srcref != 0 {
		result.SrcrefTree(srcref)
	}
	return result
}

// Java static method: NodeUtil.isNaN(com.google.javascript.rhino.Node n)
func NodeUtil_IsNaN(n pkg_rhino.Node) bool {
	_ = n
	return (((n.IsName() && (j2gort.CheckNotNullString(n.GetString()) == "NaN")) || (((((n.RawToken_() == uint8(16)) && n.GetFirstChild().IsNumber()) && (n.GetFirstChild().GetDouble() == float64(int32(0)))) && n.GetLastChild().IsNumber()) && (n.GetLastChild().GetDouble() == float64(int32(0))))) || n.MatchesQualifiedName_1_Node(StaticGet_NodeUtil_NUMBER_NAN()))
}

// Java static method: NodeUtil.countAstSizeUpToLimit(com.google.javascript.rhino.Node n, int limit)
func NodeUtil_CountAstSizeUpToLimit(n pkg_rhino.Node, limit int32) int32 {
	_ = n
	_ = limit
	var wrappedSize []int32 = []int32{int32(0)}
	_ = wrappedSize
	NodeUtil_VisitPreOrder_3(n, NewLambda_NodeUtil_Visitor(func(n12 pkg_rhino.Node) {
		_ = n12
		a_35 := wrappedSize
		i_36 := int32(0)
		a_35[i_36] = a_35[i_36] + 1
	}), j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_37 uint32) bool {
		n1 := pkg_rhino.Node(rawNode_37)
		_ = n1
		return (wrappedSize[int32(0)] < limit)
	}))
	return wrappedSize[int32(0)]
}

// Java static method: NodeUtil.countAstSize(com.google.javascript.rhino.Node n)
func NodeUtil_CountAstSize(n pkg_rhino.Node) int32 {
	_ = n
	var count int32 = int32(1)
	_ = count
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			count = (count + NodeUtil_CountAstSize(c))
		}
	}
	return count
}

// Java static method: NodeUtil.createConstantJsDoc()
func NodeUtil_CreateConstantJsDoc() *pkg_rhino.JSDocInfo {
	var builder *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_S_Builder()
	_ = builder
	builder.RecordConstancy()
	return builder.Build_0()
}

// Java static method: NodeUtil.isGoogProvideCall(com.google.javascript.rhino.Node n)
func NodeUtil_IsGoogProvideCall(n pkg_rhino.Node) bool {
	_ = n
	if NodeUtil_IsExprCall(n) {
		var target pkg_rhino.Node = n.GetFirstFirstChild()
		_ = target
		return StaticGet_NodeUtil_GOOG_PROVIDE().Matches(target)
	}
	return false
}

// Java static method: NodeUtil.isGoogModuleCall(com.google.javascript.rhino.Node n)
func NodeUtil_IsGoogModuleCall(n pkg_rhino.Node) bool {
	_ = n
	if NodeUtil_IsExprCall(n) {
		var target pkg_rhino.Node = n.GetFirstFirstChild()
		_ = target
		return StaticGet_NodeUtil_GOOG_MODULE().Matches(target)
	}
	return false
}

// Java static method: NodeUtil.isGoogModuleGetCall(com.google.javascript.rhino.Node callNode)
func NodeUtil_IsGoogModuleGetCall(callNode pkg_rhino.Node) bool {
	_ = callNode
	if !(callNode.IsCall()) {
		return false
	}
	return ((StaticGet_NodeUtil_GOOG_MODULE_GET().Matches(callNode.GetFirstChild()) && callNode.HasTwoChildren()) && callNode.GetSecondChild().IsStringLit())
}

// Java static method: NodeUtil.isGoogRequireCall(com.google.javascript.rhino.Node call)
func NodeUtil_IsGoogRequireCall(call pkg_rhino.Node) bool {
	_ = call
	if call.IsCall() {
		var target pkg_rhino.Node = call.GetFirstChild()
		_ = target
		return StaticGet_NodeUtil_GOOG_REQUIRE().Matches(target)
	}
	return false
}

// Java static method: NodeUtil.isGoogRequireTypeCall(com.google.javascript.rhino.Node call)
func NodeUtil_IsGoogRequireTypeCall(call pkg_rhino.Node) bool {
	_ = call
	if call.IsCall() {
		var target pkg_rhino.Node = call.GetFirstChild()
		_ = target
		return StaticGet_NodeUtil_GOOG_REQUIRE_TYPE().Matches(target)
	}
	return false
}

// Java static method: NodeUtil.isGoogForwardDeclareCall(com.google.javascript.rhino.Node call)
func NodeUtil_IsGoogForwardDeclareCall(call pkg_rhino.Node) bool {
	_ = call
	if call.IsCall() {
		var target pkg_rhino.Node = call.GetFirstChild()
		_ = target
		return StaticGet_NodeUtil_GOOG_FORWARD_DECLARE().Matches(target)
	}
	return false
}

// Java static method: NodeUtil.isGoogRequireDynamicCall(com.google.javascript.rhino.Node call)
func NodeUtil_IsGoogRequireDynamicCall(call pkg_rhino.Node) bool {
	_ = call
	if call.IsCall() {
		var target pkg_rhino.Node = call.GetFirstChild()
		_ = target
		return StaticGet_NodeUtil_GOOG_REQUIRE_DYNAMIC().Matches(target)
	}
	return false
}

// Java static method: NodeUtil.isGoogWeakUsageCall(com.google.javascript.rhino.Node call)
func NodeUtil_IsGoogWeakUsageCall(call pkg_rhino.Node) bool {
	_ = call
	if call.IsCall() {
		var target pkg_rhino.Node = call.GetFirstChild()
		_ = target
		return (StaticGet_NodeUtil_GOOG_WEAK_USAGE().Matches(target) || StaticGet_NodeUtil_GOOG_WEAK_USAGE_MANGLED().Matches(target))
	}
	return false
}

// Java static method: NodeUtil.isModuleScopeRoot(com.google.javascript.rhino.Node n)
func NodeUtil_IsModuleScopeRoot(n pkg_rhino.Node) bool {
	_ = n
	return (n.IsModuleBody() || NodeUtil_IsBundledGoogModuleScopeRoot(n))
}

// Java static method: NodeUtil.isGoogModuleExportsReference(com.google.javascript.jscomp.Scope scope, com.google.javascript.rhino.Node possibleName)
func NodeUtil_IsGoogModuleExportsReference(scope *Scope, possibleName pkg_rhino.Node) bool {
	_ = scope
	_ = possibleName
	if !(possibleName.IsName()) || !(j2gort.CheckNotNullString(possibleName.GetString()) == "exports") {
		return false
	}
	var var_ *Var = j2gort.CastPtr[Var](scope.GetVar(possibleName.GetString()))
	_ = var_
	var scopeRoot pkg_rhino.Node
	_ = scopeRoot
	if var_ != nil {
		scopeRoot = var_.GetScopeRoot()
	} else {
		scopeRoot = pkg_rhino.Node(0)
	}
	return ((scopeRoot != 0) && (scopeRoot.IsModuleBody() || (scopeRoot.IsFunction() && NodeUtil_IsBundledGoogModuleScopeRoot(NodeUtil_GetFunctionBody(scopeRoot)))))
}

// Java static method: NodeUtil.isBundledGoogModuleCall(com.google.javascript.rhino.Node n)
func NodeUtil_IsBundledGoogModuleCall(n pkg_rhino.Node) bool {
	_ = n
	if !((n.IsCall() && n.HasTwoChildren()) && StaticGet_NodeUtil_GOOG_LOADMODULE().Matches(n.GetFirstChild())) {
		return false
	}
	return (((n.HasParent() && n.GetParent().IsExprResult()) && (n.GetGrandparent() != 0)) && n.GetGrandparent().IsScript())
}

// Java static method: NodeUtil.isBundledGoogModuleScopeRoot(com.google.javascript.rhino.Node n)
func NodeUtil_IsBundledGoogModuleScopeRoot(n pkg_rhino.Node) bool {
	_ = n
	if (!(n.IsBlock()) || !(n.HasChildren())) || !(NodeUtil_IsGoogModuleCall(n.GetFirstChild())) {
		return false
	}
	var function pkg_rhino.Node = n.GetParent()
	_ = function
	if (((function == 0) || !(function.IsFunction())) || !(NodeUtil_GetFunctionParameters(function).HasOneChild())) || !(NodeUtil_GetFunctionParameters(function).GetFirstChild().MatchesName_1_string("exports")) {
		return false
	}
	var call pkg_rhino.Node = function.GetParent()
	_ = call
	if (!(call.IsCall()) || !(call.HasTwoChildren())) || !(StaticGet_NodeUtil_GOOG_LOADMODULE().Matches(call.GetFirstChild())) {
		return false
	}
	return (call.GetParent().IsExprResult() && call.GetGrandparent().IsScript())
}

// Java static method: NodeUtil.isGoogModuleDeclareLegacyNamespaceCall(com.google.javascript.rhino.Node n)
func NodeUtil_IsGoogModuleDeclareLegacyNamespaceCall(n pkg_rhino.Node) bool {
	_ = n
	if NodeUtil_IsExprCall(n) {
		var target pkg_rhino.Node = n.GetFirstFirstChild()
		_ = target
		return StaticGet_NodeUtil_GOOG_MODULE_DECLARE_LEGACY_NAMESPACE().Matches(target)
	}
	return false
}

// Java static method: NodeUtil.isGoogSetTestOnlyCall(com.google.javascript.rhino.Node n)
func NodeUtil_IsGoogSetTestOnlyCall(n pkg_rhino.Node) bool {
	_ = n
	if NodeUtil_IsExprCall(n) {
		var target pkg_rhino.Node = n.GetFirstFirstChild()
		_ = target
		return StaticGet_NodeUtil_GOOG_SET_TEST_ONLY().Matches(target)
	}
	return false
}

// Java static method: NodeUtil.isTopLevel(com.google.javascript.rhino.Node n)
func NodeUtil_IsTopLevel(n pkg_rhino.Node) bool {
	_ = n
	return (n.IsScript() || n.IsModuleBody())
}

// Java static method: NodeUtil.isGoogModuleFile(com.google.javascript.rhino.Node n)
func NodeUtil_IsGoogModuleFile(n pkg_rhino.Node) bool {
	_ = n
	return (((n.IsScript() && n.HasChildren()) && n.GetFirstChild().IsModuleBody()) && NodeUtil_IsGoogModuleCall(n.GetFirstFirstChild()))
}

// Java static method: NodeUtil.isLegacyGoogModuleFile(com.google.javascript.rhino.Node n)
func NodeUtil_IsLegacyGoogModuleFile(n pkg_rhino.Node) bool {
	_ = n
	return (NodeUtil_IsGoogModuleFile(n) && NodeUtil_IsGoogModuleDeclareLegacyNamespaceCall(n.GetFirstChild().GetSecondChild()))
}

// Java static method: NodeUtil.isConstructor(com.google.javascript.rhino.Node fnNode)
func NodeUtil_IsConstructor(fnNode pkg_rhino.Node) bool {
	_ = fnNode
	if (fnNode == 0) || !(fnNode.IsFunction()) {
		return false
	}
	var type_ pkg_rhino.JSType = fnNode.GetJSType()
	_ = type_
	var jsDocInfo *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(fnNode)
	_ = jsDocInfo
	var color *pkg_colors.Color = fnNode.GetColor()
	_ = color
	return (((((type_ != nil && !type_.IsJavaNil_()) && type_.GetJSTypeBase().IsConstructor()) || ((jsDocInfo != nil) && jsDocInfo.IsConstructor())) || ((color != nil) && color.IsConstructor())) || NodeUtil_IsEs6Constructor(fnNode))
}

// Java static method: NodeUtil.isEs6ConstructorMemberFunctionDef(com.google.javascript.rhino.Node memberFunctionDef)
func NodeUtil_IsEs6ConstructorMemberFunctionDef(memberFunctionDef pkg_rhino.Node) bool {
	_ = memberFunctionDef
	if !(memberFunctionDef.IsMemberFunctionDef()) {
		return false
	}
	return ((memberFunctionDef.GetParent().IsClassMembers() && !(memberFunctionDef.IsStaticMember())) && (j2gort.CheckNotNullString(memberFunctionDef.GetString()) == "constructor"))
}

// Java static method: NodeUtil.isEs6Constructor(com.google.javascript.rhino.Node fnNode)
func NodeUtil_IsEs6Constructor(fnNode pkg_rhino.Node) bool {
	_ = fnNode
	if !(fnNode.IsFunction()) {
		return false
	}
	var memberFunctionDef pkg_rhino.Node = fnNode.GetParent()
	_ = memberFunctionDef
	return ((memberFunctionDef != 0) && NodeUtil_IsEs6ConstructorMemberFunctionDef(memberFunctionDef))
}

// Java static method: NodeUtil.isGetterOrSetter(com.google.javascript.rhino.Node propNode)
func NodeUtil_IsGetterOrSetter(propNode pkg_rhino.Node) bool {
	_ = propNode
	if NodeUtil_IsGetOrSetKey(propNode) {
		return true
	}
	if !(propNode.IsStringKey()) || !(propNode.GetFirstChild().IsFunction()) {
		return false
	}
	var keyName string = propNode.GetString()
	_ = keyName
	return ((j2gort.CheckNotNullString(keyName) == "get") || (j2gort.CheckNotNullString(keyName) == "set"))
}

// Java static method: NodeUtil.isCallTo(com.google.javascript.rhino.Node n, java.lang.String qualifiedName)
func NodeUtil_IsCallTo_2_Node_string(n pkg_rhino.Node, qualifiedName string) bool {
	_ = n
	_ = qualifiedName
	return (n.IsCall() && n.GetFirstChild().MatchesQualifiedName_1_string(qualifiedName))
}

// Java static method: NodeUtil.isCallTo(com.google.javascript.rhino.Node n, com.google.javascript.rhino.QualifiedName qualifiedName)
func NodeUtil_IsCallTo_2_Node_QualifiedName(n pkg_rhino.Node, qualifiedName pkg_rhino.QualifiedName) bool {
	_ = n
	_ = qualifiedName
	return (n.IsCall() && qualifiedName.Matches(n.GetFirstChild()))
}

// Java static method: NodeUtil.isCallTo(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node targetMethod)
func NodeUtil_IsCallTo_2_Node_Node(n pkg_rhino.Node, targetMethod pkg_rhino.Node) bool {
	_ = n
	_ = targetMethod
	if !(n.IsCall()) {
		return false
	}
	return n.GetFirstChild().MatchesQualifiedName_1_Node(targetMethod)
}

// Java static method: NodeUtil.collectExternVariableNames(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.rhino.Node externs)
func NodeUtil_CollectExternVariableNames(compiler AbstractCompiler, externs pkg_rhino.Node) j2gort.Ext_com_google_common_collect_ImmutableSet {
	_ = compiler
	_ = externs
	clinit_ReferenceCollector()
	var val_38 ReferenceCollector
	InitReferenceCollector_3(&val_38, compiler, StaticGet_ReferenceCollector_DO_NOTHING_BEHAVIOR(), NewSyntacticScopeCreator_1(compiler))
	var externsRefs *ReferenceCollector = &val_38
	_ = externsRefs
	externsRefs.Process_1(externs)
	var externsNames j2gort.Ext_com_google_common_collect_ImmutableSet_Builder = j2gort.Ext_com_google_common_collect_ImmutableSet_S_Builder()
	_ = externsNames
	{
		cur_39 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](externsRefs.GetAllSymbols())))
		for cur_39.HasNext() {
			rawElem_40 := cur_39.Next()
			v := j2gort.CastPtr[Var](rawElem_40)
			_ = v
			if !(v.IsParam()) {
				j2gort.ImmSetBuilderAddStrVoid(externsNames, v.GetName())
			}
		}
	}
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](externsNames.Build())
}

// Java static method: NodeUtil.createSynthesizedExternsSymbol(com.google.javascript.jscomp.AbstractCompiler compiler, java.lang.String nameToAdd)
func NodeUtil_CreateSynthesizedExternsSymbol(compiler AbstractCompiler, nameToAdd string) {
	_ = compiler
	_ = nameToAdd
	var name pkg_rhino.Node = pkg_rhino.IR_Name(nameToAdd)
	_ = name
	name.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME(), true)
	var var_ pkg_rhino.Node = pkg_rhino.IR_Var_1(name)
	_ = var_
	var input *CompilerInput = compiler.GetSynthesizedExternsInput()
	_ = input
	var root pkg_rhino.Node = input.GetAstRoot(compiler)
	_ = root
	name.SetStaticSourceFileFrom(root)
	var_.SetStaticSourceFileFrom(root)
	root.AddChildToBack(var_)
	compiler.ReportChangeToEnclosingScope(var_)
}

// Java static method: NodeUtil.markNewScopesChanged(com.google.javascript.rhino.Node node, com.google.javascript.jscomp.AbstractCompiler compiler)
func NodeUtil_MarkNewScopesChanged_2_Node_AbstractCompiler(node pkg_rhino.Node, compiler AbstractCompiler) {
	_ = node
	_ = compiler
	NodeUtil_MarkNewScopesChanged_2_Node_ChangeTracker(node, compiler.GetChangeTracker())
}

// Java static method: NodeUtil.markNewScopesChanged(com.google.javascript.rhino.Node node, com.google.javascript.jscomp.ChangeTracker changeTracker)
func NodeUtil_MarkNewScopesChanged_2_Node_ChangeTracker(node pkg_rhino.Node, changeTracker *ChangeTracker) {
	_ = node
	_ = changeTracker
	changeTracker.MarkNewScopesChanged(node)
}

// Java static method: NodeUtil.markFunctionsDeleted(com.google.javascript.rhino.Node node, com.google.javascript.jscomp.AbstractCompiler compiler)
func NodeUtil_MarkFunctionsDeleted_2_Node_AbstractCompiler(node pkg_rhino.Node, compiler AbstractCompiler) {
	_ = node
	_ = compiler
	NodeUtil_MarkFunctionsDeleted_2_Node_ChangeTracker(node, compiler.GetChangeTracker())
}

// Java static method: NodeUtil.markFunctionsDeleted(com.google.javascript.rhino.Node node, com.google.javascript.jscomp.ChangeTracker changeTracker)
func NodeUtil_MarkFunctionsDeleted_2_Node_ChangeTracker(node pkg_rhino.Node, changeTracker *ChangeTracker) {
	_ = node
	_ = changeTracker
	if node.IsFunction() {
		changeTracker.ReportFunctionDeleted(node)
	}
	{
		var child pkg_rhino.Node = node.GetFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			NodeUtil_MarkFunctionsDeleted_2_Node_ChangeTracker(child, changeTracker)
		}
	}
	var shadowed pkg_rhino.Node = node.GetClosureUnawareShadow()
	_ = shadowed
	if shadowed != 0 {
		NodeUtil_MarkFunctionsDeleted_2_Node_ChangeTracker(shadowed, changeTracker)
	}
}

// Java static method: NodeUtil.getParentChangeScopeNodes(java.util.List<com.google.javascript.rhino.Node> scopeNodes)
func NodeUtil_GetParentChangeScopeNodes(scopeNodes j2gort.Ext_java_util_List) j2gort.Ext_java_util_List {
	_ = scopeNodes
	var parentScopeNodes j2gort.Ext_java_util_Set = j2gort.NewExt_java_util_LinkedHashSet_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](scopeNodes))
	_ = parentScopeNodes
	{
		cur_41 := j2gort.NewIterCursor(j2gort.CheckNotNull(scopeNodes))
		for cur_41.HasNext() {
			scopeNode := pkg_rhino.Node(cur_41.NextNode())
			_ = scopeNode
			j2gort.SetAddNode(parentScopeNodes, ChangeTracker_GetEnclosingChangeScopeRoot(scopeNode))
		}
	}
	return j2gort.NewExt_java_util_ArrayList_1_Ext_java_util_Collection(parentScopeNodes)
}

// Java static method: NodeUtil.removeNestedChangeScopeNodes(java.util.List<com.google.javascript.rhino.Node> scopeNodes)
func NodeUtil_RemoveNestedChangeScopeNodes(scopeNodes j2gort.Ext_java_util_List) j2gort.Ext_java_util_List {
	_ = scopeNodes
	var uniqueScopeNodes j2gort.Ext_java_util_Set = j2gort.NewExt_java_util_LinkedHashSet_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](scopeNodes))
	_ = uniqueScopeNodes
	{
		cur_42 := j2gort.NewIterCursor(j2gort.CheckNotNull(scopeNodes))
		for cur_42.HasNext() {
			scopeNode := pkg_rhino.Node(cur_42.NextNode())
			_ = scopeNode
			{
				var ancestor pkg_rhino.Node = scopeNode.GetParent()
				_ = ancestor
				for ; ancestor != 0; ancestor = ancestor.GetParent() {
					if ChangeTracker_IsChangeScopeRoot(ancestor) && j2gort.SetContainsNode(uniqueScopeNodes, ancestor) {
						j2gort.SetRemoveNode(uniqueScopeNodes, scopeNode)
						break
					}
				}
			}
		}
	}
	return j2gort.NewExt_java_util_ArrayList_1_Ext_java_util_Collection(uniqueScopeNodes)
}

// Java static method: NodeUtil.firstInvocationArgCursor(com.google.javascript.rhino.Node invocation, boolean isTaggedTemplate)
func NodeUtil_FirstInvocationArgCursor(invocation pkg_rhino.Node, isTaggedTemplate bool) pkg_rhino.Node {
	_ = invocation
	_ = isTaggedTemplate
	if isTaggedTemplate {
		var templateLit pkg_rhino.Node = invocation.GetLastChild()
		_ = templateLit
		if !(templateLit.IsTemplateLit()) {
			j2gort.CheckState(false)
		}
		return NodeUtil_AdvanceToTemplateLitSub(templateLit.GetFirstChild())
	}
	if !(NodeUtil_IsCallOrNew(invocation)) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(invocation))
	}
	return invocation.GetSecondChild()
}

// Java static method: NodeUtil.invocationArgFromCursor(com.google.javascript.rhino.Node cursor, boolean isTaggedTemplate)
func NodeUtil_InvocationArgFromCursor(cursor pkg_rhino.Node, isTaggedTemplate bool) pkg_rhino.Node {
	_ = cursor
	_ = isTaggedTemplate
	if isTaggedTemplate {
		return cursor.GetFirstChild()
	} else {
		return cursor
	}
}

// Java static method: NodeUtil.nextInvocationArgCursor(com.google.javascript.rhino.Node cursor, boolean isTaggedTemplate)
func NodeUtil_NextInvocationArgCursor(cursor pkg_rhino.Node, isTaggedTemplate bool) pkg_rhino.Node {
	_ = cursor
	_ = isTaggedTemplate
	var next pkg_rhino.Node = cursor.GetNext()
	_ = next
	if isTaggedTemplate {
		return NodeUtil_AdvanceToTemplateLitSub(next)
	} else {
		return next
	}
}

// Java static method: NodeUtil.advanceToTemplateLitSub(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node child)
func NodeUtil_AdvanceToTemplateLitSub(child pkg_rhino.Node) pkg_rhino.Node {
	_ = child
	for (child != 0) && !(child.IsTemplateLitSub()) {
		child = child.GetNext()
	}
	return child
}

// Java static method: NodeUtil.getInvocationArgsCount(com.google.javascript.rhino.Node invocation)
func NodeUtil_GetInvocationArgsCount(invocation pkg_rhino.Node) int32 {
	_ = invocation
	if invocation.IsTaggedTemplateLit() {
		var templateLit pkg_rhino.Node = invocation.GetLastChild()
		_ = templateLit
		if !(templateLit.IsTemplateLit()) {
			j2gort.CheckState(false)
		}
		var count int32 = int32(1)
		_ = count
		{
			var child pkg_rhino.Node = templateLit.GetFirstChild()
			_ = child
			for ; child != 0; child = child.GetNext() {
				if child.IsTemplateLitSub() {
					count = count + 1
				}
			}
		}
		return count
	} else {
		return (invocation.GetChildCount() - int32(1))
	}
}

// Java static method: NodeUtil.getAllVarNamesDeclaredInModule(com.google.javascript.rhino.Node moduleNode, com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.ScopeCreator scopeCreator, com.google.javascript.jscomp.Scope globalScope)
func NodeUtil_GetAllVarNamesDeclaredInModule(moduleNode pkg_rhino.Node, compiler AbstractCompiler, scopeCreator ScopeCreator, globalScope *Scope) j2gort.Ext_java_util_Set {
	_ = moduleNode
	_ = compiler
	_ = scopeCreator
	_ = globalScope
	if !(moduleNode.IsModuleBody()) {
		j2gort.CheckStateStr(false, "getAllVarsDeclaredInModule expects a module body node")
	}
	if !(globalScope.IsGlobal()) {
		j2gort.CheckStateMsg(false, j2gort.PtrToAny(globalScope))
	}
	var nameVars j2gort.Ext_java_util_Set = j2gort.NewExt_java_util_LinkedHashSet()
	_ = nameVars
	var finder NodeTraversal_ScopedCallback = NewNodeUtil__Anon_66(nameVars)
	_ = finder
	{
		clinit_NodeTraversal()
		var ntb_ NodeTraversal_Builder
		InitNodeTraversal_Builder(&ntb_)
		(&ntb_).SetCompiler(compiler)
		(&ntb_).SetCallback_1_NodeTraversal_Callback(Upcast_NodeTraversal_Callback(finder))
		(&ntb_).SetScopeCreator(scopeCreator)
		(&ntb_).TraverseWithScope(moduleNode, ToIface_AbstractScope(globalScope))
	}
	return nameVars
}

// Java static method: NodeUtil.getAllVarsDeclaredInFunction(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.ScopeCreator scopeCreator, com.google.javascript.jscomp.Scope scope)
func NodeUtil_GetAllVarsDeclaredInFunction(compiler AbstractCompiler, scopeCreator ScopeCreator, scope *Scope) *NodeUtil_AllVarsDeclaredInFunction {
	_ = compiler
	_ = scopeCreator
	_ = scope
	if !(scope.IsFunctionScope()) {
		j2gort.CheckStateMsg(false, j2gort.PtrToAny(scope))
	}
	var nameVarMap j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_LinkedHashMap()
	_ = nameVarMap
	var orderedVars j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList()
	_ = orderedVars
	var finder NodeTraversal_ScopedCallback = NewNodeUtil__Anon_67(nameVarMap, orderedVars, scope)
	_ = finder
	{
		clinit_NodeTraversal()
		var ntb_ NodeTraversal_Builder
		InitNodeTraversal_Builder(&ntb_)
		(&ntb_).SetCompiler(compiler)
		(&ntb_).SetCallback_1_NodeTraversal_Callback(Upcast_NodeTraversal_Callback(finder))
		(&ntb_).SetScopeCreator(scopeCreator)
		(&ntb_).TraverseAtScope(ToIface_AbstractScope(scope))
	}
	return NewNodeUtil_AllVarsDeclaredInFunction(nameVarMap, orderedVars)
}

// Java static method: NodeUtil.isObjLitProperty(com.google.javascript.rhino.Node node)
func NodeUtil_IsObjLitProperty(node pkg_rhino.Node) bool {
	_ = node
	return ((((node.IsStringKey() || node.IsGetterDef()) || node.IsSetterDef()) || node.IsMemberFunctionDef()) || node.IsComputedProp())
}

// Java static method: NodeUtil.isBlocklessArrowFunctionResult(com.google.javascript.rhino.Node n)
func NodeUtil_IsBlocklessArrowFunctionResult(n pkg_rhino.Node) bool {
	_ = n
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	return ((((parent != 0) && parent.IsFunction()) && (n == parent.GetLastChild())) && !(n.IsBlock()))
}

// Java static method: NodeUtil.getFeatureSetOfScript(com.google.javascript.rhino.Node scriptNode)
func NodeUtil_GetFeatureSetOfScript(scriptNode pkg_rhino.Node) *pkg_parser.FeatureSet {
	_ = scriptNode
	if !(scriptNode.IsScript()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(scriptNode))
	}
	return j2gort.CastPtr[pkg_parser.FeatureSet](scriptNode.GetFeatureSet())
}

// Java static method: NodeUtil.addFeatureToScript(com.google.javascript.rhino.Node scriptNode, com.google.javascript.jscomp.parsing.parser.FeatureSet.Feature feature, com.google.javascript.jscomp.AbstractCompiler compiler)
func NodeUtil_AddFeatureToScript(scriptNode pkg_rhino.Node, feature *pkg_parser.FeatureSet_Feature, compiler AbstractCompiler) {
	_ = scriptNode
	_ = feature
	_ = compiler
	if !(scriptNode.IsScript()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(scriptNode))
	}
	var currentFeatures *pkg_parser.FeatureSet = NodeUtil_GetFeatureSetOfScript(scriptNode)
	_ = currentFeatures
	var newFeatures *pkg_parser.FeatureSet
	_ = newFeatures
	if currentFeatures != nil {
		newFeatures = currentFeatures.With_1_FeatureSet_Feature(feature)
	} else {
		newFeatures = pkg_parser.StaticGet_FeatureSet_BARE_MINIMUM().With_1_FeatureSet_Feature(feature)
	}
	if newFeatures != currentFeatures {
		scriptNode.SetFeatureSet(j2gort.PtrToAny(newFeatures))
	}
	if feature != pkg_parser.StaticGet_FeatureSet_Feature_MODULES() {
		if !(compiler.GetAllowableFeatures().Contains_1_FeatureSet_Feature(feature)) {
			j2gort.CheckStateFmt(false, "Cannot add feature: %s. It is not supported in the output language, and either 1) its corresponding transpilation pass has already run or 2) transpilation of this feature is unsupported entirely", j2gort.PtrToAny(feature))
		}
	} else {
		compiler.SetAllowableFeatures(compiler.GetAllowableFeatures().With_1_FeatureSet_Feature(feature))
	}
}

// Java static method: NodeUtil.removeFeatureFromScript(com.google.javascript.rhino.Node scriptNode, com.google.javascript.jscomp.parsing.parser.FeatureSet.Feature feature)
func NodeUtil_RemoveFeatureFromScript(scriptNode pkg_rhino.Node, feature *pkg_parser.FeatureSet_Feature) {
	_ = scriptNode
	_ = feature
	var currentFeatures *pkg_parser.FeatureSet = NodeUtil_GetFeatureSetOfScript(scriptNode)
	_ = currentFeatures
	if currentFeatures != nil {
		var updated *pkg_parser.FeatureSet = currentFeatures.Without_2(feature, []*pkg_parser.FeatureSet_Feature{})
		_ = updated
		if updated != currentFeatures {
			scriptNode.SetFeatureSet(j2gort.PtrToAny(updated))
		}
	}
}

// Java static method: NodeUtil.removeFeaturesFromScript(com.google.javascript.rhino.Node scriptNode, com.google.javascript.jscomp.parsing.parser.FeatureSet featureSet)
func NodeUtil_RemoveFeaturesFromScript(scriptNode pkg_rhino.Node, featureSet *pkg_parser.FeatureSet) {
	_ = scriptNode
	_ = featureSet
	var currentFeatures *pkg_parser.FeatureSet = NodeUtil_GetFeatureSetOfScript(scriptNode)
	_ = currentFeatures
	if currentFeatures != nil {
		var updated *pkg_parser.FeatureSet = currentFeatures.Without_1(featureSet)
		_ = updated
		if updated != currentFeatures {
			scriptNode.SetFeatureSet(j2gort.PtrToAny(updated))
		}
	}
}

// Java static method: NodeUtil.removeFeatureFromAllScripts(com.google.javascript.rhino.Node root, com.google.javascript.jscomp.parsing.parser.FeatureSet.Feature feature, com.google.javascript.jscomp.AbstractCompiler compiler)
func NodeUtil_RemoveFeatureFromAllScripts(root pkg_rhino.Node, feature *pkg_parser.FeatureSet_Feature, compiler AbstractCompiler) {
	_ = root
	_ = feature
	_ = compiler
	if !(root.IsRoot()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(root))
	}
	{
		var childNode pkg_rhino.Node = root.GetFirstChild()
		_ = childNode
		for ; childNode != 0; childNode = childNode.GetNext() {
			if !(childNode.IsScript()) {
				j2gort.CheckState(false)
			}
			NodeUtil_RemoveFeatureFromScript(childNode, feature)
		}
	}
	compiler.MarkFeatureNotAllowed(feature)
}

// Java static method: NodeUtil.removeFeaturesFromAllScripts(com.google.javascript.rhino.Node root, com.google.javascript.jscomp.parsing.parser.FeatureSet featureSet, com.google.javascript.jscomp.AbstractCompiler compiler)
func NodeUtil_RemoveFeaturesFromAllScripts(root pkg_rhino.Node, featureSet *pkg_parser.FeatureSet, compiler AbstractCompiler) {
	_ = root
	_ = featureSet
	_ = compiler
	if !(root.IsRoot()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(root))
	}
	{
		var childNode pkg_rhino.Node = root.GetFirstChild()
		_ = childNode
		for ; childNode != 0; childNode = childNode.GetNext() {
			if !(childNode.IsScript()) {
				j2gort.CheckState(false)
			}
			NodeUtil_RemoveFeaturesFromScript(childNode, featureSet)
		}
	}
	compiler.MarkFeatureSetNotAllowed(featureSet)
}

// Java static method: NodeUtil.addFeatureToAllScripts(com.google.javascript.rhino.Node root, com.google.javascript.jscomp.parsing.parser.FeatureSet.Feature feature, com.google.javascript.jscomp.AbstractCompiler compiler)
func NodeUtil_AddFeatureToAllScripts(root pkg_rhino.Node, feature *pkg_parser.FeatureSet_Feature, compiler AbstractCompiler) {
	_ = root
	_ = feature
	_ = compiler
	if !(root.IsRoot()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(root))
	}
	{
		var childNode pkg_rhino.Node = root.GetFirstChild()
		_ = childNode
		for ; childNode != 0; childNode = childNode.GetNext() {
			if !(childNode.IsScript()) {
				j2gort.CheckState(false)
			}
			NodeUtil_AddFeatureToScript(childNode, feature, compiler)
		}
	}
}

// Java static method: NodeUtil.addFeaturesToScript(com.google.javascript.rhino.Node scriptNode, com.google.javascript.jscomp.parsing.parser.FeatureSet features, com.google.javascript.jscomp.AbstractCompiler compiler)
func NodeUtil_AddFeaturesToScript(scriptNode pkg_rhino.Node, features *pkg_parser.FeatureSet, compiler AbstractCompiler) {
	_ = scriptNode
	_ = features
	_ = compiler
	if !(scriptNode.IsScript()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(scriptNode))
	}
	{
		cur_43 := j2gort.NewIterCursor(j2gort.CheckNotNull(features.GetFeatures()))
		for cur_43.HasNext() {
			rawElem_44 := cur_43.Next()
			feature := j2gort.CastPtr[pkg_parser.FeatureSet_Feature](rawElem_44)
			_ = feature
			NodeUtil_AddFeatureToScript(scriptNode, feature, compiler)
		}
	}
}

// Java static method: NodeUtil.getParamOrPatternNames(com.google.javascript.rhino.Node n, java.util.function.Consumer<com.google.javascript.rhino.Node> cb)
func NodeUtil_GetParamOrPatternNames(n pkg_rhino.Node, cb j2gort.Ext_java_util_function_Consumer) {
	_ = n
	_ = cb
	pkg_parsing.ParsingUtil_GetParamOrPatternNames(n, cb)
}

// Java static method: NodeUtil.getGoogRequireInfo(java.lang.String name, com.google.javascript.jscomp.Scope scope)
func NodeUtil_GetGoogRequireInfo_2(name string, scope *Scope) *NodeUtil_GoogRequire {
	_ = name
	_ = scope
	var var_ *Var = j2gort.CastPtr[Var](scope.GetVar(name))
	_ = var_
	if var_ == nil {
		return nil
	}
	return NodeUtil_GetGoogRequireInfo_1_Var(var_)
}

// Java static method: NodeUtil.getGoogRequireInfo(com.google.javascript.jscomp.Var var)
func NodeUtil_GetGoogRequireInfo_1_Var(var_ *Var) *NodeUtil_GoogRequire {
	_ = var_
	j2gort.CheckNotNullPtr(var_)
	if !(var_.GetScopeRoot().IsModuleBody()) || (var_.GetNameNode() == 0) {
		return nil
	}
	return NodeUtil_GetGoogRequireInfo_1_Node(var_.GetNameNode())
}

// Java static method: NodeUtil.getGoogRequireInfo(com.google.javascript.rhino.Node nameNode)
func NodeUtil_GetGoogRequireInfo_1_Node(nameNode pkg_rhino.Node) *NodeUtil_GoogRequire {
	_ = nameNode
	if nameNode.IsImportStar() {
		return nil
	}
	if !(nameNode.IsName()) {
		j2gort.CheckStateFmt(false, "unexpected node type: %s", j2gort.BoxNode(nameNode))
	}
	if NodeUtil_IsNameDeclaration(nameNode.GetParent()) {
		var requireCall pkg_rhino.Node = nameNode.GetFirstChild()
		_ = requireCall
		if requireCall == 0 {
			return nil
		}
		var isStrongRequire bool = NodeUtil_IsGoogRequireCall(requireCall)
		_ = isStrongRequire
		if !(isStrongRequire || NodeUtil_IsGoogRequireTypeCall(requireCall)) {
			return nil
		}
		var namespace string = requireCall.GetSecondChild().GetString()
		_ = namespace
		return NodeUtil_GoogRequire_FromNamespace(namespace, isStrongRequire)
	} else if nameNode.GetParent().IsStringKey() && nameNode.GetGrandparent().IsObjectPattern() {
		var requireCall_2 pkg_rhino.Node = nameNode.GetGrandparent().GetNext()
		_ = requireCall_2
		if requireCall_2 == 0 {
			return nil
		}
		var isStrongRequire_2 bool = NodeUtil_IsGoogRequireCall(requireCall_2)
		_ = isStrongRequire_2
		if !(isStrongRequire_2 || NodeUtil_IsGoogRequireTypeCall(requireCall_2)) {
			return nil
		}
		var property string = nameNode.GetParent().GetString()
		_ = property
		var namespace_2 string = requireCall_2.GetSecondChild().GetString()
		_ = namespace_2
		return NodeUtil_GoogRequire_FromNamespaceAndProperty(namespace_2, property, isStrongRequire_2)
	}
	return nil
}

// Java static method: NodeUtil.estimateNumLines(com.google.javascript.rhino.Node scriptNode)
func NodeUtil_EstimateNumLines(scriptNode pkg_rhino.Node) int32 {
	_ = scriptNode
	if !(scriptNode.IsScript()) {
		j2gort.CheckArgument(false)
	}
	var current pkg_rhino.Node = scriptNode
	_ = current
	for current.HasChildren() {
		current = current.GetLastChild()
	}
	return (current.GetLineno() + int32(1))
}

// Java static method: NodeUtil.getOriginalName(java.lang.String name)
func NodeUtil_GetOriginalName(name string) string {
	_ = name
	name = j2gort.JStrReplaceAll(name, "\\$\\$\\d+$", "")
	name = MakeDeclaredNamesUnique_ContextualRenameInverter_GetOriginalName(name)
	if javalang.StringStartsWith(j2gort.CheckNotNullString(name), "module$exports$") {
		var lastDollar int32 = j2gort.JStrLastIndexOfCp(name, int32(uint16(36)))
		_ = lastDollar
		if lastDollar != int32(-1) {
			name = j2gort.JStrSubFrom(name, (lastDollar + int32(1)))
		}
	} else if javalang.StringStartsWith(j2gort.CheckNotNullString(name), "module$contents$") {
		var lastUnderscore int32 = j2gort.JStrLastIndexOfCp(name, int32(uint16(95)))
		_ = lastUnderscore
		if lastUnderscore != int32(-1) {
			name = j2gort.JStrSubFrom(name, (lastUnderscore + int32(1)))
		}
	}
	return name
}

// Java enum: com.google.javascript.jscomp.NodeUtil.ValueType
type NodeUtil_ValueType uint8

const (
	enum_NodeUtil_ValueType_UNDETERMINED NodeUtil_ValueType = 1
	enum_NodeUtil_ValueType_NULL         NodeUtil_ValueType = 2
	enum_NodeUtil_ValueType_VOID         NodeUtil_ValueType = 3
	enum_NodeUtil_ValueType_NUMBER       NodeUtil_ValueType = 4
	enum_NodeUtil_ValueType_BIGINT       NodeUtil_ValueType = 5
	enum_NodeUtil_ValueType_STRING       NodeUtil_ValueType = 6
	enum_NodeUtil_ValueType_BOOLEAN      NodeUtil_ValueType = 7
	enum_NodeUtil_ValueType_OBJECT       NodeUtil_ValueType = 8
)

var enumBaseTable_NodeUtil_ValueType = [...]j2gort.EnumBase{
	{Name_: "UNDETERMINED", Ordinal_: 0},
	{Name_: "NULL", Ordinal_: 1},
	{Name_: "VOID", Ordinal_: 2},
	{Name_: "NUMBER", Ordinal_: 3},
	{Name_: "BIGINT", Ordinal_: 4},
	{Name_: "STRING", Ordinal_: 5},
	{Name_: "BOOLEAN", Ordinal_: 6},
	{Name_: "OBJECT", Ordinal_: 7},
}
var enumValues_NodeUtil_ValueType = []NodeUtil_ValueType{1, 2, 3, 4, 5, 6, 7, 8}

func StaticGet_NodeUtil_ValueType_UNDETERMINED() NodeUtil_ValueType {
	return enum_NodeUtil_ValueType_UNDETERMINED
}
func StaticGet_NodeUtil_ValueType_NULL() NodeUtil_ValueType {
	return enum_NodeUtil_ValueType_NULL
}
func StaticGet_NodeUtil_ValueType_VOID() NodeUtil_ValueType {
	return enum_NodeUtil_ValueType_VOID
}
func StaticGet_NodeUtil_ValueType_NUMBER() NodeUtil_ValueType {
	return enum_NodeUtil_ValueType_NUMBER
}
func StaticGet_NodeUtil_ValueType_BIGINT() NodeUtil_ValueType {
	return enum_NodeUtil_ValueType_BIGINT
}
func StaticGet_NodeUtil_ValueType_STRING() NodeUtil_ValueType {
	return enum_NodeUtil_ValueType_STRING
}
func StaticGet_NodeUtil_ValueType_BOOLEAN() NodeUtil_ValueType {
	return enum_NodeUtil_ValueType_BOOLEAN
}
func StaticGet_NodeUtil_ValueType_OBJECT() NodeUtil_ValueType {
	return enum_NodeUtil_ValueType_OBJECT
}

func NodeUtil_ValueType_Values() []NodeUtil_ValueType {
	return j2gort.CloneSlice(enumValues_NodeUtil_ValueType)
}

func NodeUtil_ValueType_ValueOf(name string) NodeUtil_ValueType {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_NodeUtil_ValueType {
		if enumBaseTable_NodeUtil_ValueType[i_].Name_ == name {
			return NodeUtil_ValueType(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.NodeUtil.ValueType." + name))
}

func (this_ NodeUtil_ValueType) IsJavaNil_() bool                              { return this_ == 0 }
func (NodeUtil_ValueType) HasIdentityEquals_()                                 {}
func (NodeUtil_ValueType) Implements_j2gort_Ext_java_lang_Enum()               {}
func (NodeUtil_ValueType) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ NodeUtil_ValueType) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (NodeUtil_ValueType) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.NodeUtil$ValueType", "ValueType"
}
func (this_ NodeUtil_ValueType) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_NodeUtil_ValueType[this_-1]
}
func (this_ NodeUtil_ValueType) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ NodeUtil_ValueType) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ NodeUtil_ValueType) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_NodeUtil_ValueType[this_-1].Name_
}
func (this_ NodeUtil_ValueType) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_NodeUtil_ValueType[this_-1].Name_
}
func (this_ NodeUtil_ValueType) Equals(other_ any) bool {
	o_, ok_ := other_.(NodeUtil_ValueType)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ NodeUtil_ValueType) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_NodeUtil_ValueType[this_-1].HashCode()
}
func (this_ NodeUtil_ValueType) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_NodeUtil_ValueType[this_-1].CompareTo(other_)
}
func (this_ NodeUtil_ValueType) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_NodeUtil_ValueType[this_-1].Name_)
}
func init() {
	vals_ := []any{NodeUtil_ValueType(1), NodeUtil_ValueType(2), NodeUtil_ValueType(3), NodeUtil_ValueType(4), NodeUtil_ValueType(5), NodeUtil_ValueType(6), NodeUtil_ValueType(7), NodeUtil_ValueType(8)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.NodeUtil.ValueType", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.NodeUtil$ValueType", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.NodeUtil.ValueType", "com.google.javascript.jscomp.NodeUtil$ValueType", func(v_ any) bool {
		return j2gort.InstanceOf[NodeUtil_ValueType](v_)
	})
}

// Java class: com.google.javascript.jscomp.NodeUtil.VarCollector
type NodeUtil_VarCollector struct {
	vars    j2gort.OrderedMap[string, pkg_rhino.Node] `gson:"vars"`
	idHash_ atomic.Int32
}

func (this_ *NodeUtil_VarCollector) BindSelf_(s any) {
	_ = s
}
func (this_ *NodeUtil_VarCollector) IsJavaNil_() bool { return this_ == nil }
func (*NodeUtil_VarCollector) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.NodeUtil$VarCollector", "VarCollector"
}
func (this_ *NodeUtil_VarCollector) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*NodeUtil_VarCollector) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.NodeUtil.VarCollector", "com.google.javascript.jscomp.NodeUtil$VarCollector", func(v_ any) bool {
		return j2gort.InstanceOf[*NodeUtil_VarCollector](v_)
	})
}
func (*NodeUtil_VarCollector) Implements_jscomp_NodeUtil_Visitor() {}

func (this_ *NodeUtil_VarCollector) initFields_NodeUtil_VarCollector() {
	_ = this_
	this_.vars = j2gort.NewOrderedMap[string, pkg_rhino.Node]()
}

// Java constructor: VarCollector()
func NewNodeUtil_VarCollector() *NodeUtil_VarCollector {
	this_ := &NodeUtil_VarCollector{}
	InitNodeUtil_VarCollector(this_)
	return this_
}

func InitNodeUtil_VarCollector(this_ *NodeUtil_VarCollector) {
	this_.initFields_NodeUtil_VarCollector()
}

// Java method: VarCollector.visit(com.google.javascript.rhino.Node n)
func (this_ *NodeUtil_VarCollector) Visit(n pkg_rhino.Node) {
	_ = n
	if n.IsVar() {
		NodeUtil_VisitLhsNodesInNode(n, j2gort.NewLambda_Ext_java_util_function_Consumer(func(lp_45 any) {
			name := j2gort.CastNode[pkg_rhino.Node](lp_45)
			_ = name
			this_.vars.PutIfAbsent(name.GetString(), name)
		}))
	}
}

// Java class: com.google.javascript.jscomp.NodeUtil.MatchNameNode
type NodeUtil_MatchNameNode struct {
	name    string `gson:"name"`
	idHash_ atomic.Int32
}

func (this_ *NodeUtil_MatchNameNode) BindSelf_(s any) {
	_ = s
}
func (this_ *NodeUtil_MatchNameNode) IsJavaNil_() bool { return this_ == nil }
func (*NodeUtil_MatchNameNode) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.NodeUtil$MatchNameNode", "MatchNameNode"
}
func (this_ *NodeUtil_MatchNameNode) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*NodeUtil_MatchNameNode) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.NodeUtil.MatchNameNode", "com.google.javascript.jscomp.NodeUtil$MatchNameNode", func(v_ any) bool {
		return j2gort.InstanceOf[*NodeUtil_MatchNameNode](v_)
	})
}
func (*NodeUtil_MatchNameNode) Implements_j2gort_Ext_com_google_common_base_Predicate() {}
func (*NodeUtil_MatchNameNode) Implements_j2gort_Ext_java_util_function_Predicate()     {}

func (this_ *NodeUtil_MatchNameNode) initFields_NodeUtil_MatchNameNode() {
	_ = this_
	this_.name = j2gort.NullString
}

// Java constructor: MatchNameNode(java.lang.String name)
func NewNodeUtil_MatchNameNode(name string) *NodeUtil_MatchNameNode {
	this_ := &NodeUtil_MatchNameNode{}
	InitNodeUtil_MatchNameNode(this_, name)
	return this_
}

func InitNodeUtil_MatchNameNode(this_ *NodeUtil_MatchNameNode, name string) {
	_ = name
	this_.initFields_NodeUtil_MatchNameNode()
	this_.name = name
}

// Java method: MatchNameNode.apply(com.google.javascript.rhino.Node n)
func (this_ *NodeUtil_MatchNameNode) TestNode_(rawNode_46 uint32) bool {
	if this_ == nil {
		j2gort.ThrowNPE()
	}
	n := pkg_rhino.Node(rawNode_46)
	_ = n
	return (n.IsName() && func() bool {
		sr_47 := n.GetString()
		se_48 := this_.name
		return j2gort.CheckNotNullString(sr_47) == se_48
	}())
}

func (this_ *NodeUtil_MatchNameNode) Apply(arg_ any) bool {
	if this_ == nil {
		j2gort.ThrowNPE()
	}
	return this_.TestNode_(j2gort.CastNode[uint32](arg_))
}

func (this_ *NodeUtil_MatchNameNode) Equals(p0_ any) bool {
	return j2gort.RefEquals(this_, p0_)
}

func (this_ *NodeUtil_MatchNameNode) Test(p0_ any) bool {
	return this_.Apply(p0_)
}

// Java class: com.google.javascript.jscomp.NodeUtil.MatchNodeType
type NodeUtil_MatchNodeType struct {
	idHash_ atomic.Int32
	type_   pkg_rhino.Token `gson:"type"`
}

func (this_ *NodeUtil_MatchNodeType) BindSelf_(s any) {
	_ = s
}
func (this_ *NodeUtil_MatchNodeType) IsJavaNil_() bool { return this_ == nil }
func (*NodeUtil_MatchNodeType) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.NodeUtil$MatchNodeType", "MatchNodeType"
}
func (this_ *NodeUtil_MatchNodeType) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*NodeUtil_MatchNodeType) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.NodeUtil.MatchNodeType", "com.google.javascript.jscomp.NodeUtil$MatchNodeType", func(v_ any) bool {
		return j2gort.InstanceOf[*NodeUtil_MatchNodeType](v_)
	})
}
func (*NodeUtil_MatchNodeType) Implements_j2gort_Ext_com_google_common_base_Predicate() {}
func (*NodeUtil_MatchNodeType) Implements_j2gort_Ext_java_util_function_Predicate()     {}

func (this_ *NodeUtil_MatchNodeType) initFields_NodeUtil_MatchNodeType() {
	_ = this_
}

// Java constructor: MatchNodeType(com.google.javascript.rhino.Token type)
func NewNodeUtil_MatchNodeType(type_ pkg_rhino.Token) *NodeUtil_MatchNodeType {
	this_ := &NodeUtil_MatchNodeType{}
	InitNodeUtil_MatchNodeType(this_, type_)
	return this_
}

func InitNodeUtil_MatchNodeType(this_ *NodeUtil_MatchNodeType, type_ pkg_rhino.Token) {
	_ = type_
	this_.initFields_NodeUtil_MatchNodeType()
	this_.type_ = type_
}

// Java method: MatchNodeType.apply(com.google.javascript.rhino.Node n)
func (this_ *NodeUtil_MatchNodeType) TestNode_(rawNode_49 uint32) bool {
	if this_ == nil {
		j2gort.ThrowNPE()
	}
	n := pkg_rhino.Node(rawNode_49)
	_ = n
	return (n.GetToken() == this_.type_)
}

func (this_ *NodeUtil_MatchNodeType) Apply(arg_ any) bool {
	if this_ == nil {
		j2gort.ThrowNPE()
	}
	return this_.TestNode_(j2gort.CastNode[uint32](arg_))
}

func (this_ *NodeUtil_MatchNodeType) Equals(p0_ any) bool {
	return j2gort.RefEquals(this_, p0_)
}

func (this_ *NodeUtil_MatchNodeType) Test(p0_ any) bool {
	return this_.Apply(p0_)
}

// Java class: com.google.javascript.jscomp.NodeUtil.MatchDeclaration
type NodeUtil_MatchDeclaration struct {
	idHash_ atomic.Int32
}

func (this_ *NodeUtil_MatchDeclaration) BindSelf_(s any) {
	_ = s
}
func (this_ *NodeUtil_MatchDeclaration) IsJavaNil_() bool { return this_ == nil }
func (*NodeUtil_MatchDeclaration) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.NodeUtil$MatchDeclaration", "MatchDeclaration"
}
func (this_ *NodeUtil_MatchDeclaration) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*NodeUtil_MatchDeclaration) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.NodeUtil.MatchDeclaration", "com.google.javascript.jscomp.NodeUtil$MatchDeclaration", func(v_ any) bool {
		return j2gort.InstanceOf[*NodeUtil_MatchDeclaration](v_)
	})
}
func (*NodeUtil_MatchDeclaration) Implements_j2gort_Ext_com_google_common_base_Predicate() {}
func (*NodeUtil_MatchDeclaration) Implements_j2gort_Ext_java_util_function_Predicate()     {}

func (this_ *NodeUtil_MatchDeclaration) initFields_NodeUtil_MatchDeclaration() {
	_ = this_
}

// Java constructor: MatchDeclaration()
func NewNodeUtil_MatchDeclaration() *NodeUtil_MatchDeclaration {
	this_ := &NodeUtil_MatchDeclaration{}
	InitNodeUtil_MatchDeclaration(this_)
	return this_
}

func InitNodeUtil_MatchDeclaration(this_ *NodeUtil_MatchDeclaration) {
	this_.initFields_NodeUtil_MatchDeclaration()
}

// Java method: MatchDeclaration.apply(com.google.javascript.rhino.Node n)
func (this_ *NodeUtil_MatchDeclaration) TestNode_(rawNode_50 uint32) bool {
	if this_ == nil {
		j2gort.ThrowNPE()
	}
	n := pkg_rhino.Node(rawNode_50)
	_ = n
	return NodeUtil_IsDeclaration(n)
}

func (this_ *NodeUtil_MatchDeclaration) Apply(arg_ any) bool {
	if this_ == nil {
		j2gort.ThrowNPE()
	}
	return this_.TestNode_(j2gort.CastNode[uint32](arg_))
}

func (this_ *NodeUtil_MatchDeclaration) Equals(p0_ any) bool {
	return j2gort.RefEquals(this_, p0_)
}

func (this_ *NodeUtil_MatchDeclaration) Test(p0_ any) bool {
	return this_.Apply(p0_)
}

// Java class: com.google.javascript.jscomp.NodeUtil.MatchShallowStatement
type NodeUtil_MatchShallowStatement struct {
	idHash_ atomic.Int32
}

func (this_ *NodeUtil_MatchShallowStatement) BindSelf_(s any) {
	_ = s
}
func (this_ *NodeUtil_MatchShallowStatement) IsJavaNil_() bool { return this_ == nil }
func (*NodeUtil_MatchShallowStatement) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.NodeUtil$MatchShallowStatement", "MatchShallowStatement"
}
func (this_ *NodeUtil_MatchShallowStatement) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*NodeUtil_MatchShallowStatement) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.NodeUtil.MatchShallowStatement", "com.google.javascript.jscomp.NodeUtil$MatchShallowStatement", func(v_ any) bool {
		return j2gort.InstanceOf[*NodeUtil_MatchShallowStatement](v_)
	})
}
func (*NodeUtil_MatchShallowStatement) Implements_j2gort_Ext_com_google_common_base_Predicate() {}
func (*NodeUtil_MatchShallowStatement) Implements_j2gort_Ext_java_util_function_Predicate()     {}

func (this_ *NodeUtil_MatchShallowStatement) initFields_NodeUtil_MatchShallowStatement() {
	_ = this_
}

// Java constructor: MatchShallowStatement()
func NewNodeUtil_MatchShallowStatement() *NodeUtil_MatchShallowStatement {
	this_ := &NodeUtil_MatchShallowStatement{}
	InitNodeUtil_MatchShallowStatement(this_)
	return this_
}

func InitNodeUtil_MatchShallowStatement(this_ *NodeUtil_MatchShallowStatement) {
	this_.initFields_NodeUtil_MatchShallowStatement()
}

// Java method: MatchShallowStatement.apply(com.google.javascript.rhino.Node n)
func (this_ *NodeUtil_MatchShallowStatement) TestNode_(rawNode_51 uint32) bool {
	if this_ == nil {
		j2gort.ThrowNPE()
	}
	n := pkg_rhino.Node(rawNode_51)
	_ = n
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	return ((n.IsRoot() || n.IsBlock()) || (!(n.IsFunction()) && NodeUtil_IsShallowStatementTree(parent)))
}

func (this_ *NodeUtil_MatchShallowStatement) Apply(arg_ any) bool {
	if this_ == nil {
		j2gort.ThrowNPE()
	}
	return this_.TestNode_(j2gort.CastNode[uint32](arg_))
}

func (this_ *NodeUtil_MatchShallowStatement) Equals(p0_ any) bool {
	return j2gort.RefEquals(this_, p0_)
}

func (this_ *NodeUtil_MatchShallowStatement) Test(p0_ any) bool {
	return this_.Apply(p0_)
}

// Java interface: com.google.javascript.jscomp.NodeUtil.Visitor
type NodeUtil_Visitor interface {
	Implements_jscomp_NodeUtil_Visitor()
	IsJavaNil_() bool
	Visit(pkg_rhino.Node)
}

func ToIface_NodeUtil_Visitor[T interface {
	NodeUtil_Visitor
	comparable
}](v T) NodeUtil_Visitor {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_NodeUtil_Visitor[T NodeUtil_Visitor](v T) NodeUtil_Visitor {
	var vi NodeUtil_Visitor = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.NodeUtil.Visitor", "com.google.javascript.jscomp.NodeUtil$Visitor", func(v_ any) bool {
		return j2gort.InstanceOf[NodeUtil_Visitor](v_)
	})
}

type lambda_NodeUtil_Visitor struct {
	fn func(pkg_rhino.Node)
}

func (l *lambda_NodeUtil_Visitor) IsJavaNil_() bool { return l == nil }
func NewLambda_NodeUtil_Visitor(fn func(pkg_rhino.Node)) NodeUtil_Visitor {
	if fn == nil {
		return nil
	}
	return &lambda_NodeUtil_Visitor{fn: fn}
}

func (*lambda_NodeUtil_Visitor) Implements_jscomp_NodeUtil_Visitor() {}
func (l *lambda_NodeUtil_Visitor) Visit(p0_ pkg_rhino.Node) {
	l.fn(p0_)
}

// Java class: com.google.javascript.jscomp.NodeUtil.PreOrderIterator
type NodeUtil_PreOrderIterator struct {
	j2gort.Ext_com_google_common_collect_AbstractIteratorBase
	traverseNodePredicate j2gort.Ext_com_google_common_base_Predicate `gson:"traverseNodePredicate"`
	current               pkg_rhino.Node                              `gson:"current"`
	idHash_               atomic.Int32
}

func (this_ *NodeUtil_PreOrderIterator) BindSelf_(s any) {
	_ = s
	this_.Ext_com_google_common_collect_AbstractIteratorBase.BindExtSelf_(s)
}
func (this_ *NodeUtil_PreOrderIterator) IsJavaNil_() bool { return this_ == nil }
func (*NodeUtil_PreOrderIterator) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.NodeUtil$PreOrderIterator", "PreOrderIterator"
}
func (this_ *NodeUtil_PreOrderIterator) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.NodeUtil.PreOrderIterator", "com.google.javascript.jscomp.NodeUtil$PreOrderIterator", func(v_ any) bool {
		return j2gort.InstanceOf[*NodeUtil_PreOrderIterator](v_)
	})
}
func (*NodeUtil_PreOrderIterator) Implements_j2gort_Ext_com_google_common_collect_AbstractIterator() {
}
func (*NodeUtil_PreOrderIterator) Implements_j2gort_Ext_com_google_common_collect_UnmodifiableIterator() {
}
func (*NodeUtil_PreOrderIterator) Implements_j2gort_Ext_java_util_Iterator() {}

func (this_ *NodeUtil_PreOrderIterator) initFields_NodeUtil_PreOrderIterator() {
	_ = this_
}

// Java constructor: PreOrderIterator(com.google.javascript.rhino.Node root, com.google.common.base.Predicate<com.google.javascript.rhino.Node> traverseNodePredicate)
func NewNodeUtil_PreOrderIterator(root pkg_rhino.Node, traverseNodePredicate j2gort.Ext_com_google_common_base_Predicate) *NodeUtil_PreOrderIterator {
	this_ := &NodeUtil_PreOrderIterator{}
	this_.BindSelf_(this_)
	InitNodeUtil_PreOrderIterator(this_, root, traverseNodePredicate)
	return this_
}

func InitNodeUtil_PreOrderIterator(this_ *NodeUtil_PreOrderIterator, root pkg_rhino.Node, traverseNodePredicate j2gort.Ext_com_google_common_base_Predicate) {
	_ = root
	_ = traverseNodePredicate
	j2gort.InitExt_com_google_common_collect_AbstractIterator(&this_.Ext_com_google_common_collect_AbstractIteratorBase)
	this_.initFields_NodeUtil_PreOrderIterator()
	j2gort.CheckNotNullNode(root)
	this_.traverseNodePredicate = traverseNodePredicate
	this_.current = root
}

// Java method: PreOrderIterator.computeNext()
func (this_ *NodeUtil_PreOrderIterator) ComputeNext() any {
	if this_.current == 0 {
		return j2gort.BoxNode(j2gort.CastNode[pkg_rhino.Node](this_.EndOfData()))
	}
	var returnValue pkg_rhino.Node = this_.current
	_ = returnValue
	this_.current = this_.priv_PreOrderIterator_calculateNextNode_4034(returnValue)
	return j2gort.BoxNode(returnValue)
}

// Java method: PreOrderIterator.calculateNextNode(com.google.javascript.rhino.Node currentNode)
func (this_ *NodeUtil_PreOrderIterator) priv_PreOrderIterator_calculateNextNode_4034(currentNode pkg_rhino.Node) pkg_rhino.Node {
	_ = currentNode
	j2gort.CheckNotNullNode(currentNode)
	if j2gort.PredTestNode(this_.traverseNodePredicate, uint32(currentNode)) {
		if currentNode.HasChildren() {
			return currentNode.GetFirstChild()
		}
	}
	for currentNode != 0 {
		var next pkg_rhino.Node = currentNode.GetNext()
		_ = next
		if next != 0 {
			return next
		}
		currentNode = currentNode.GetParent()
	}
	return pkg_rhino.Node(0)
}

// Java class: com.google.javascript.jscomp.NodeUtil.AllVarsDeclaredInFunction
type NodeUtil_AllVarsDeclaredInFunction struct {
	allVarsInFn j2gort.Ext_java_util_Map  `gson:"allVarsInFn"`
	orderedVars j2gort.Ext_java_util_List `gson:"orderedVars"`
	idHash_     atomic.Int32
}

func (this_ *NodeUtil_AllVarsDeclaredInFunction) BindSelf_(s any) {
	_ = s
}
func (this_ *NodeUtil_AllVarsDeclaredInFunction) IsJavaNil_() bool { return this_ == nil }
func (*NodeUtil_AllVarsDeclaredInFunction) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.NodeUtil$AllVarsDeclaredInFunction", "AllVarsDeclaredInFunction"
}
func (this_ *NodeUtil_AllVarsDeclaredInFunction) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*NodeUtil_AllVarsDeclaredInFunction) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.NodeUtil.AllVarsDeclaredInFunction", "com.google.javascript.jscomp.NodeUtil$AllVarsDeclaredInFunction", func(v_ any) bool {
		return j2gort.InstanceOf[*NodeUtil_AllVarsDeclaredInFunction](v_)
	})
}

func (this_ *NodeUtil_AllVarsDeclaredInFunction) initFields_NodeUtil_AllVarsDeclaredInFunction() {
	_ = this_
}

// Java constructor: AllVarsDeclaredInFunction(java.util.Map<java.lang.String,com.google.javascript.jscomp.Var> allVarsInFn, java.util.List<com.google.javascript.jscomp.Var> orderedVars)
func NewNodeUtil_AllVarsDeclaredInFunction(allVarsInFn j2gort.Ext_java_util_Map, orderedVars j2gort.Ext_java_util_List) *NodeUtil_AllVarsDeclaredInFunction {
	this_ := &NodeUtil_AllVarsDeclaredInFunction{}
	InitNodeUtil_AllVarsDeclaredInFunction(this_, allVarsInFn, orderedVars)
	return this_
}

func InitNodeUtil_AllVarsDeclaredInFunction(this_ *NodeUtil_AllVarsDeclaredInFunction, allVarsInFn j2gort.Ext_java_util_Map, orderedVars j2gort.Ext_java_util_List) {
	_ = allVarsInFn
	_ = orderedVars
	this_.initFields_NodeUtil_AllVarsDeclaredInFunction()
	if !(allVarsInFn.IsEmpty() == orderedVars.IsEmpty()) {
		j2gort.CheckState(false)
	}
	this_.allVarsInFn = allVarsInFn
	this_.orderedVars = orderedVars
}

// Java method: AllVarsDeclaredInFunction.getAllVariables()
func (this_ *NodeUtil_AllVarsDeclaredInFunction) GetAllVariables() j2gort.Ext_java_util_Map {
	return this_.allVarsInFn
}

// Java method: AllVarsDeclaredInFunction.getAllVariablesInOrder()
func (this_ *NodeUtil_AllVarsDeclaredInFunction) GetAllVariablesInOrder() j2gort.Ext_java_util_List {
	return this_.orderedVars
}

// Java record: com.google.javascript.jscomp.NodeUtil.GoogRequire
type NodeUtil_GoogRequire struct {
	namespace       string `gson:"namespace"`
	property        string `gson:"property"`
	isStrongRequire bool   `gson:"isStrongRequire"`
}

func (this_ *NodeUtil_GoogRequire) BindSelf_(s any) {
	_ = s
}
func (this_ *NodeUtil_GoogRequire) IsJavaNil_() bool { return this_ == nil }
func (*NodeUtil_GoogRequire) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.NodeUtil$GoogRequire", "GoogRequire"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.NodeUtil.GoogRequire", "com.google.javascript.jscomp.NodeUtil$GoogRequire", func(v_ any) bool {
		return j2gort.InstanceOf[*NodeUtil_GoogRequire](v_)
	})
}
func (*NodeUtil_GoogRequire) Implements_j2gort_Ext_java_lang_Record() {}

func (this_ *NodeUtil_GoogRequire) initFields_NodeUtil_GoogRequire() {
	_ = this_
}

// Java constructor: GoogRequire(java.lang.String namespace, java.lang.@org.jspecify.annotations.Nullable String property, boolean isStrongRequire)
func NewNodeUtil_GoogRequire(namespace string, property string, isStrongRequire bool) *NodeUtil_GoogRequire {
	this_ := &NodeUtil_GoogRequire{}
	InitNodeUtil_GoogRequire(this_, namespace, property, isStrongRequire)
	return this_
}

func InitNodeUtil_GoogRequire(this_ *NodeUtil_GoogRequire, namespace string, property string, isStrongRequire bool) {
	_ = namespace
	_ = property
	_ = isStrongRequire
	this_.initFields_NodeUtil_GoogRequire()
	j2gort.CheckNotNullString(namespace)
	this_.namespace = namespace
	this_.property = property
	this_.isStrongRequire = isStrongRequire
}

// Java static method: GoogRequire.fromNamespace(java.lang.String namespace, boolean isStrongRequire)
func NodeUtil_GoogRequire_FromNamespace(namespace string, isStrongRequire bool) *NodeUtil_GoogRequire {
	_ = namespace
	_ = isStrongRequire
	return NewNodeUtil_GoogRequire(namespace, j2gort.NullString, isStrongRequire)
}

// Java static method: GoogRequire.fromNamespaceAndProperty(java.lang.String namespace, java.lang.String property, boolean isStrongRequire)
func NodeUtil_GoogRequire_FromNamespaceAndProperty(namespace string, property string, isStrongRequire bool) *NodeUtil_GoogRequire {
	_ = namespace
	_ = property
	_ = isStrongRequire
	return NewNodeUtil_GoogRequire(namespace, property, isStrongRequire)
}

func (this_ *NodeUtil_GoogRequire) Namespace() string     { return this_.namespace }
func (this_ *NodeUtil_GoogRequire) Property() string      { return this_.property }
func (this_ *NodeUtil_GoogRequire) IsStrongRequire() bool { return this_.isStrongRequire }
func (this_ *NodeUtil_GoogRequire) Equals(other_ any) bool {
	o_, ok_ := other_.(*NodeUtil_GoogRequire)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if this_.namespace != o_.namespace {
		return false
	}
	if this_.property != o_.property {
		return false
	}
	if this_.isStrongRequire != o_.isStrongRequire {
		return false
	}
	return true
}
func (this_ *NodeUtil_GoogRequire) HashCode() int32 {
	var h_ int32 = 0
	{
		hs_52 := this_.namespace
		h_ = 31 * h_
		if hs_52 != j2gort.NullString {
			h_ += javalang.StringHashCode(hs_52)
		}
	}
	{
		hs_53 := this_.property
		h_ = 31 * h_
		if hs_53 != j2gort.NullString {
			h_ += javalang.StringHashCode(hs_53)
		}
	}
	if this_.isStrongRequire {
		h_ = 31*h_ + 1231
	} else {
		h_ = 31*h_ + 1237
	}
	return h_
}
func (this_ *NodeUtil_GoogRequire) String() string {
	_ = fmt.Sprint
	return "GoogRequire[" + "namespace=" + j2gort.StrVal(j2gort.BoxString(this_.namespace)) + ", " + "property=" + j2gort.StrVal(j2gort.BoxString(this_.property)) + ", " + "isStrongRequire=" + j2gort.StrVal(j2gort.ValueOfBoolean(this_.isStrongRequire)) + "]"
}

// Java class:
type NodeUtil__Anon_66 struct {
	cap_0_nameVars j2gort.Ext_java_util_Set
	idHash_        atomic.Int32
}

func (this_ *NodeUtil__Anon_66) BindSelf_(s any) {
	_ = s
}
func (this_ *NodeUtil__Anon_66) IsJavaNil_() bool { return this_ == nil }
func (*NodeUtil__Anon_66) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.NodeUtil$1", ""
}
func (this_ *NodeUtil__Anon_66) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*NodeUtil__Anon_66) HasIdentityEquals_()                             {}
func (*NodeUtil__Anon_66) Implements_jscomp_NodeTraversal_ScopedCallback() {}
func (*NodeUtil__Anon_66) Implements_jscomp_NodeTraversal_Callback()       {}

func (this_ *NodeUtil__Anon_66) initFields_NodeUtil__Anon_66() {
	_ = this_
}

func NewNodeUtil__Anon_66(cap_p_nameVars j2gort.Ext_java_util_Set) *NodeUtil__Anon_66 {
	this_ := &NodeUtil__Anon_66{}
	this_.cap_0_nameVars = cap_p_nameVars
	this_.initFields_NodeUtil__Anon_66()
	return this_
}

// Java method: .enterScope(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *NodeUtil__Anon_66) EnterScope(t *NodeTraversal) {
	_ = t
	var currentScope *Scope = t.GetScope()
	_ = currentScope
	if currentScope.IsModuleScope() {
		{
			cur_52 := (currentScope).iterVars_()
			for cur_52.HasNext() {

				v := cur_52.Next()

				_ = v
				j2gort.SetAddStr(this_.cap_0_nameVars, v.GetName())
			}
		}
	}
}

// Java method: .exitScope(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *NodeUtil__Anon_66) ExitScope(t *NodeTraversal) {
	_ = t
}

// Java method: .shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *NodeUtil__Anon_66) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	return n.IsModuleBody()
}

// Java method: .visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *NodeUtil__Anon_66) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
}

// Java class:
type NodeUtil__Anon_67 struct {
	cap_0_nameVarMap  j2gort.Ext_java_util_Map
	cap_1_orderedVars j2gort.Ext_java_util_List
	cap_2_scope       *Scope
	idHash_           atomic.Int32
}

func (this_ *NodeUtil__Anon_67) BindSelf_(s any) {
	_ = s
}
func (this_ *NodeUtil__Anon_67) IsJavaNil_() bool { return this_ == nil }
func (*NodeUtil__Anon_67) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.NodeUtil$2", ""
}
func (this_ *NodeUtil__Anon_67) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*NodeUtil__Anon_67) HasIdentityEquals_()                             {}
func (*NodeUtil__Anon_67) Implements_jscomp_NodeTraversal_ScopedCallback() {}
func (*NodeUtil__Anon_67) Implements_jscomp_NodeTraversal_Callback()       {}

func (this_ *NodeUtil__Anon_67) initFields_NodeUtil__Anon_67() {
	_ = this_
}

func NewNodeUtil__Anon_67(cap_p_nameVarMap j2gort.Ext_java_util_Map, cap_p_orderedVars j2gort.Ext_java_util_List, cap_p_scope *Scope) *NodeUtil__Anon_67 {
	this_ := &NodeUtil__Anon_67{}
	this_.cap_0_nameVarMap = cap_p_nameVarMap
	this_.cap_1_orderedVars = cap_p_orderedVars
	this_.cap_2_scope = cap_p_scope
	this_.initFields_NodeUtil__Anon_67()
	return this_
}

// Java method: .enterScope(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *NodeUtil__Anon_67) EnterScope(t *NodeTraversal) {
	_ = t
	var currentScope *Scope = t.GetScope()
	_ = currentScope
	{
		cur_54 := (currentScope).iterVars_()
		for cur_54.HasNext() {

			v := cur_54.Next()

			_ = v
			j2gort.MapPutStr(this_.cap_0_nameVarMap, v.GetName(), j2gort.PtrToAny(v))
			this_.cap_1_orderedVars.Add_1(j2gort.PtrToAny(v))
		}
	}
}

// Java method: .exitScope(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *NodeUtil__Anon_67) ExitScope(t *NodeTraversal) {
	_ = t
}

// Java method: .shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *NodeUtil__Anon_67) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	return (!(n.IsFunction()) || (n == this_.cap_2_scope.GetRootNode_0()))
}

// Java method: .visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *NodeUtil__Anon_67) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
}
func (*NodeUtil__Anon_66) NoOpVisit_()                 {}
func (*NodeUtil__Anon_66) CanReleaseSyntacticScopes_() {}
func (*NodeUtil__Anon_67) NoOpVisit_()                 {}
