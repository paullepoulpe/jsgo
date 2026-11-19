/*
 * Copyright 2019 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/AstAnalyzer.java (generated mechanically by j2go).

package jscomp

import (
	"fmt"
	"sync/atomic"

	pkg_colors "github.com/paullepoulpe/jsgo/transpiled/gen/colors"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.AstAnalyzer
type AstAnalyzer struct {
	typeRegistry                 *pkg_rhino.JSTypeRegistry `gson:"typeRegistry"`
	accessorSummary              *AccessorSummary          `gson:"accessorSummary"`
	idHash_                      atomic.Int32
	useTypesForLocalOptimization bool `gson:"useTypesForLocalOptimization"`
	assumeGettersArePure         bool `gson:"assumeGettersArePure"`
	assumeKnownBuiltinsArePure   bool `gson:"assumeKnownBuiltinsArePure"`
	hasRegexpGlobalReferences    bool `gson:"hasRegexpGlobalReferences"`
}

func (this_ *AstAnalyzer) BindSelf_(s any) {
	_ = s
}
func (this_ *AstAnalyzer) IsJavaNil_() bool { return this_ == nil }
func (*AstAnalyzer) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AstAnalyzer", "AstAnalyzer"
}
func (this_ *AstAnalyzer) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*AstAnalyzer) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AstAnalyzer", "com.google.javascript.jscomp.AstAnalyzer", func(v_ any) bool {
		return j2gort.InstanceOf[*AstAnalyzer](v_)
	})
}

var static_AstAnalyzer_CONSTRUCTORS_WITHOUT_SIDE_EFFECTS j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_AstAnalyzer_CONSTRUCTORS_WITHOUT_SIDE_EFFECTS() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_AstAnalyzer()
	return static_AstAnalyzer_CONSTRUCTORS_WITHOUT_SIDE_EFFECTS
}
func StaticSet_AstAnalyzer_CONSTRUCTORS_WITHOUT_SIDE_EFFECTS(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_AstAnalyzer()
	static_AstAnalyzer_CONSTRUCTORS_WITHOUT_SIDE_EFFECTS = v
}

var static_AstAnalyzer_BUILTIN_FUNCTIONS_WITHOUT_SIDEEFFECTS j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_AstAnalyzer_BUILTIN_FUNCTIONS_WITHOUT_SIDEEFFECTS() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_AstAnalyzer()
	return static_AstAnalyzer_BUILTIN_FUNCTIONS_WITHOUT_SIDEEFFECTS
}
func StaticSet_AstAnalyzer_BUILTIN_FUNCTIONS_WITHOUT_SIDEEFFECTS(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_AstAnalyzer()
	static_AstAnalyzer_BUILTIN_FUNCTIONS_WITHOUT_SIDEEFFECTS = v
}

var static_AstAnalyzer_REGEXP_METHODS j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_AstAnalyzer_REGEXP_METHODS() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_AstAnalyzer()
	return static_AstAnalyzer_REGEXP_METHODS
}
func StaticSet_AstAnalyzer_REGEXP_METHODS(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_AstAnalyzer()
	static_AstAnalyzer_REGEXP_METHODS = v
}

var static_AstAnalyzer_STRING_REGEXP_METHODS j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_AstAnalyzer_STRING_REGEXP_METHODS() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_AstAnalyzer()
	return static_AstAnalyzer_STRING_REGEXP_METHODS
}
func StaticSet_AstAnalyzer_STRING_REGEXP_METHODS(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_AstAnalyzer()
	static_AstAnalyzer_STRING_REGEXP_METHODS = v
}

var clinitDone_AstAnalyzer uint32
var clinitState_AstAnalyzer j2gort.ClinitState

func clinit_AstAnalyzer() {
	if atomic.LoadUint32(&clinitDone_AstAnalyzer) != 0 {
		return
	}
	clinitSlow_AstAnalyzer()
}

func clinitSlow_AstAnalyzer() {
	if !j2gort.EnterClinit(&clinitState_AstAnalyzer, &clinitDone_AstAnalyzer) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_AstAnalyzer, &clinitDone_AstAnalyzer)
	static_AstAnalyzer_CONSTRUCTORS_WITHOUT_SIDE_EFFECTS = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_7("Array", "Date", "Error", "Object", "RegExp", "XMLHttpRequest", []any{})
	static_AstAnalyzer_BUILTIN_FUNCTIONS_WITHOUT_SIDEEFFECTS = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_7("Object", "Array", "String", "Number", "BigInt", "Boolean", []any{"RegExp", "Error", "Symbol"})
	static_AstAnalyzer_REGEXP_METHODS = j2gort.ImmSetOfStr2("test", "exec")
	static_AstAnalyzer_STRING_REGEXP_METHODS = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_4("match", "replace", "search", "split")
}

func (this_ *AstAnalyzer) initFields_AstAnalyzer() {
	_ = this_
}

// Java constructor: AstAnalyzer(com.google.javascript.jscomp.AstAnalyzer.Options options, com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable JSTypeRegistry typeRegistry, com.google.javascript.jscomp.@org.jspecify.annotations.Nullable AccessorSummary accessorSummary)
func NewAstAnalyzer(options *AstAnalyzer_Options, typeRegistry *pkg_rhino.JSTypeRegistry, accessorSummary *AccessorSummary) *AstAnalyzer {
	this_ := &AstAnalyzer{}
	InitAstAnalyzer(this_, options, typeRegistry, accessorSummary)
	return this_
}

func InitAstAnalyzer(this_ *AstAnalyzer, options *AstAnalyzer_Options, typeRegistry *pkg_rhino.JSTypeRegistry, accessorSummary *AccessorSummary) {
	_ = options
	_ = typeRegistry
	_ = accessorSummary
	this_.initFields_AstAnalyzer()
	if !(options.assumeGettersArePure || (accessorSummary != nil)) {
		j2gort.CheckArgumentStr(false, "accessorSummary must be provided if assumeGettersArePure is false")
	}
	this_.typeRegistry = typeRegistry
	this_.accessorSummary = accessorSummary
	this_.useTypesForLocalOptimization = options.useTypesForLocalOptimization
	this_.assumeGettersArePure = options.assumeGettersArePure
	this_.hasRegexpGlobalReferences = options.hasRegexpGlobalReferences
	this_.assumeKnownBuiltinsArePure = options.assumeKnownBuiltinsArePure
}

// Java method: AstAnalyzer.mayEffectMutableState(com.google.javascript.rhino.Node n)
func (this_ *AstAnalyzer) MayEffectMutableState(n pkg_rhino.Node) bool {
	_ = n
	return this_.priv_AstAnalyzer_checkForStateChangeHelper_604(n, true)
}

// Java method: AstAnalyzer.mayHaveSideEffects(com.google.javascript.rhino.Node n)
func (this_ *AstAnalyzer) MayHaveSideEffects(n pkg_rhino.Node) bool {
	_ = n
	return this_.priv_AstAnalyzer_checkForStateChangeHelper_604(n, false)
}

// Java method: AstAnalyzer.functionCallHasSideEffects(com.google.javascript.rhino.Node callNode)
func (this_ *AstAnalyzer) FunctionCallHasSideEffects(callNode pkg_rhino.Node) bool {
	_ = callNode
	if !((callNode.IsCall() || callNode.IsTaggedTemplateLit()) || callNode.IsOptChainCall()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(callNode))
	}
	if callNode.IsNoSideEffectsCall() {
		return false
	}
	if callNode.IsOnlyModifiesArgumentsCall() && NodeUtil_AllArgsUnescapedLocal(callNode) {
		return false
	}
	var callee pkg_rhino.Node = callNode.GetFirstChild()
	_ = callee
	if callee.IsName() {
		var name string = callee.GetString()
		_ = name
		if this_.assumeKnownBuiltinsArePure && j2gort.SetContainsStr(StaticGet_AstAnalyzer_BUILTIN_FUNCTIONS_WITHOUT_SIDEEFFECTS(), name) {
			return false
		}
	} else if callee.IsGetProp() || callee.IsOptChainGetProp() {
		var method string = callee.GetString()
		_ = method
		if this_.assumeKnownBuiltinsArePure && (func() bool {
			if "valueOf" == method {
				return callNode.HasOneChild()
			}
			return (("toString" == method) && (callNode.HasOneChild() || callNode.HasTwoChildren()))
		}()) {
			return false
		}
		if callNode.IsOnlyModifiesThisCall() && NodeUtil_EvaluatesToLocalValue(callee.GetFirstChild()) {
			return false
		}
		if ((this_.assumeKnownBuiltinsArePure && callee.GetFirstChild().IsName()) && callee.IsQualifiedName()) && (j2gort.CheckNotNullString(callee.GetFirstChild().GetString()) == "Math") {
			switch j2gort.CheckNotNullString(method) {
			case "abs", "acos", "acosh", "asin", "asinh", "atan", "atanh", "atan2", "cbrt", "ceil", "cos", "cosh", "exp", "expm1", "floor", "hypot", "log", "log10", "log1p", "log2", "max", "min", "pow", "round", "sign", "sin", "sinh", "sqrt", "tan", "tanh", "trunc":
				return false
			case "random":
				return !(callNode.HasOneChild())
			default:
			}
		}
		if !(this_.hasRegexpGlobalReferences) && this_.assumeKnownBuiltinsArePure {
			if callee.GetFirstChild().IsRegExp() && j2gort.SetContainsStr(StaticGet_AstAnalyzer_REGEXP_METHODS(), method) {
				return false
			} else if this_.priv_AstAnalyzer_isTypedAsString_603(callee.GetFirstChild()) {
				var param pkg_rhino.Node = callee.GetNext()
				_ = param
				if param != 0 {
					if param.IsStringLit() {
						if j2gort.SetContainsStr(StaticGet_AstAnalyzer_STRING_REGEXP_METHODS(), method) {
							return false
						}
					} else if param.IsRegExp() {
						if "replace" == method {
							return !(param.GetNext().IsStringLit())
						} else if j2gort.SetContainsStr(StaticGet_AstAnalyzer_STRING_REGEXP_METHODS(), method) {
							return false
						}
					}
				}
			}
		}
	}
	return true
}

// Java method: AstAnalyzer.isTypedAsString(com.google.javascript.rhino.Node n)
func (this_ *AstAnalyzer) priv_AstAnalyzer_isTypedAsString_603(n pkg_rhino.Node) bool {
	_ = n
	if n.IsStringLit() {
		return true
	}
	if this_.useTypesForLocalOptimization {
		var color *pkg_colors.Color = n.GetColor()
		_ = color
		if color != nil {
			return color.Equals(j2gort.PtrToAny(pkg_colors.StaticGet_StandardColors_STRING()))
		}
		var type_ pkg_rhino.JSType = n.GetJSType()
		_ = type_
		if type_ != nil && !type_.IsJavaNil_() {
			var nativeStringType pkg_rhino.JSType = this_.typeRegistry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_STRING_TYPE())
			_ = nativeStringType
			if type_.GetJSTypeBase().EqualsJSType(nativeStringType) {
				return true
			}
		}
	}
	return false
}

// Java method: AstAnalyzer.checkForStateChangeHelper(com.google.javascript.rhino.Node n, boolean checkForNewObjects)
func (this_ *AstAnalyzer) priv_AstAnalyzer_checkForStateChangeHelper_604(n pkg_rhino.Node, checkForNewObjects bool) bool {
	_ = n
	_ = checkForNewObjects
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(43) /* Token.THROW */, int32(114) /* Token.YIELD */, int32(115) /* Token.AWAIT */, int32(113) /* Token.FOR_AWAIT_OF */, int32(112) /* Token.FOR_OF */, int32(81) /* Token.FOR_IN */, int32(84) /* Token.VAR */, int32(111) /* Token.LET */, int32(97) /* Token.CONST */, int32(120) /* Token.EXPORT */, int32(124) /* Token.DYNAMIC_IMPORT */ :
		return true
	case int32(110) /* Token.SUPER */ :
		return false
	case int32(47) /* Token.OBJECTLIT */, int32(46) /* Token.ARRAYLIT */, int32(42) /* Token.REGEXP */ :
		return (checkForNewObjects || this_.priv_AstAnalyzer_checkForChangeStateInChildren_605(n, checkForNewObjects))
	case int32(126) /* Token.OBJECT_REST */, int32(128) /* Token.OBJECT_SPREAD */ :
		return (!(this_.assumeGettersArePure) || this_.priv_AstAnalyzer_checkForChangeStateInChildren_605(n, checkForNewObjects))
	case int32(125) /* Token.ITER_REST */, int32(127) /* Token.ITER_SPREAD */ :
		return (NodeUtil_IteratesImpureIterable(n) || this_.priv_AstAnalyzer_checkForChangeStateInChildren_605(n, checkForNewObjects))
	case int32(32) /* Token.NAME */ :
		return n.HasChildren()
	case int32(73) /* Token.FUNCTION */ :
		return (checkForNewObjects || NodeUtil_IsFunctionDeclaration(n))
	case int32(95) /* Token.GETTER_DEF */, int32(96) /* Token.SETTER_DEF */, int32(107) /* Token.MEMBER_FUNCTION_DEF */ :
		return false
	case int32(129) /* Token.COMPUTED_PROP */ :
		return func() bool {
			switch parent.TokenSwitchOrd_() {
			case int32(106) /* Token.CLASS_MEMBERS */ :
				return this_.priv_AstAnalyzer_checkForStateChangeHelper_604(n.GetFirstChild(), checkForNewObjects)
			case int32(103) /* Token.OBJECT_PATTERN */ :
				return (parent.GetLastChild().IsObjectRest() || this_.priv_AstAnalyzer_checkForChangeStateInChildren_605(n, checkForNewObjects))
			case int32(47) /* Token.OBJECTLIT */ :
				return this_.priv_AstAnalyzer_checkForChangeStateInChildren_605(n, checkForNewObjects)
			default:
				panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Illegal COMPUTED_PROP parent "+j2gort.StrVal(parent))))
			}
			panic("unreachable switch expression")
		}()
	case int32(108) /* Token.MEMBER_FIELD_DEF */ :
		return ((n.IsStaticMember() && n.HasChildren()) && this_.priv_AstAnalyzer_checkForStateChangeHelper_604(n.GetFirstChild(), checkForNewObjects))
	case int32(109) /* Token.COMPUTED_FIELD_DEF */ :
		return (this_.priv_AstAnalyzer_checkForStateChangeHelper_604(n.GetFirstChild(), checkForNewObjects) || ((n.IsStaticMember() && (n.GetSecondChild() != 0)) && this_.priv_AstAnalyzer_checkForStateChangeHelper_604(n.GetSecondChild(), checkForNewObjects)))
	case int32(105) /* Token.CLASS */ :
		return (((checkForNewObjects || NodeUtil_IsClassDeclaration(n)) || this_.priv_AstAnalyzer_checkForStateChangeHelper_604(n.GetSecondChild(), checkForNewObjects)) || this_.priv_AstAnalyzer_checkForStateChangeHelper_604(n.GetLastChild(), checkForNewObjects))
	case int32(23) /* Token.NEW */ :
		return ((checkForNewObjects || this_.ConstructorCallHasSideEffects(n)) || this_.priv_AstAnalyzer_checkForChangeStateInChildren_605(n, checkForNewObjects))
	case int32(28) /* Token.CALL */, int32(31) /* Token.OPTCHAIN_CALL */, int32(130) /* Token.TAGGED_TEMPLATELIT */ :
		return (this_.FunctionCallHasSideEffects(n) || this_.priv_AstAnalyzer_checkForChangeStateInChildren_605(n, checkForNewObjects))
	case int32(101) /* Token.CAST */, int32(69) /* Token.AND */, int32(90) /* Token.BLOCK */, int32(89) /* Token.ROOT */, int32(93) /* Token.EXPR_RESULT */, int32(67) /* Token.HOOK */, int32(74) /* Token.IF */, int32(106) /* Token.CLASS_MEMBERS */, int32(49) /* Token.PARAM_LIST */, int32(134) /* Token.DEFAULT_VALUE */, int32(33) /* Token.NUMBER */, int32(34) /* Token.BIGINT */, int32(68) /* Token.OR */, int32(70) /* Token.COALESCE */, int32(37) /* Token.THIS */, int32(39) /* Token.TRUE */, int32(38) /* Token.FALSE */, int32(36) /* Token.NULL */, int32(35) /* Token.STRINGLIT */, int32(75) /* Token.SWITCH */, int32(132) /* Token.TEMPLATELIT_SUB */, int32(48) /* Token.TRY */, int32(88) /* Token.EMPTY */, int32(131) /* Token.TEMPLATELIT */, int32(133) /* Token.TEMPLATELIT_STRING */ :
		return this_.priv_AstAnalyzer_checkForChangeStateInChildren_605(n, checkForNewObjects)
	case int32(100) /* Token.STRING_KEY */ :
		if parent.IsObjectPattern() {
			if this_.priv_AstAnalyzer_getPropertyKind_609(n.GetString()).HasGetter() {
				return true
			} else if parent.GetLastChild().IsObjectRest() {
				return true
			}
		}
		return this_.priv_AstAnalyzer_checkForChangeStateInChildren_605(n, checkForNewObjects)
	case int32(27) /* Token.GETELEM */, int32(30) /* Token.OPTCHAIN_GETELEM */ :
		return this_.priv_AstAnalyzer_checkForChangeStateInChildren_605(n, checkForNewObjects)
	case int32(26) /* Token.GETPROP */, int32(29) /* Token.OPTCHAIN_GETPROP */ :
		return (this_.priv_AstAnalyzer_getPropertyKind_609(n.GetString()).HasGetterOrSetter() || this_.priv_AstAnalyzer_checkForChangeStateInChildren_605(n, checkForNewObjects))
	default:
		if NodeUtil_IsSimpleOperator(n) {
			return this_.priv_AstAnalyzer_checkForChangeStateInChildren_605(n, checkForNewObjects)
		}
		if NodeUtil_IsAssignmentOp(n) {
			return this_.priv_AstAnalyzer_checkAssignmentForChangeState_606(n, checkForNewObjects)
		}
		return true
	}
	panic("unreachable switch expression")
}

// Java method: AstAnalyzer.checkForChangeStateInChildren(com.google.javascript.rhino.Node n, boolean checkForNewObjects)
func (this_ *AstAnalyzer) priv_AstAnalyzer_checkForChangeStateInChildren_605(n pkg_rhino.Node, checkForNewObjects bool) bool {
	_ = n
	_ = checkForNewObjects
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			if this_.priv_AstAnalyzer_checkForStateChangeHelper_604(c, checkForNewObjects) {
				return true
			}
		}
	}
	return false
}

// Java method: AstAnalyzer.checkAssignmentForChangeState(com.google.javascript.rhino.Node assignOp, boolean checkForNewObjects)
func (this_ *AstAnalyzer) priv_AstAnalyzer_checkAssignmentForChangeState_606(assignOp pkg_rhino.Node, checkForNewObjects bool) bool {
	_ = assignOp
	_ = checkForNewObjects
	var assignTarget pkg_rhino.Node = assignOp.GetFirstChild()
	_ = assignTarget
	if assignTarget.IsName() {
		return true
	}
	var assignValue pkg_rhino.Node = assignOp.GetLastChild()
	_ = assignValue
	if this_.priv_AstAnalyzer_checkForStateChangeHelper_604(assignTarget, checkForNewObjects) || this_.priv_AstAnalyzer_checkForStateChangeHelper_604(assignValue, checkForNewObjects) {
		return true
	}
	if !(NodeUtil_IsNormalGet(assignTarget)) {
		return !(NodeUtil_IsLiteralValue(assignTarget, true))
	}
	var current pkg_rhino.Node = assignTarget.GetFirstChild()
	_ = current
	if NodeUtil_EvaluatesToLocalValue(current) {
		return false
	}
	for NodeUtil_IsNormalGet(current) {
		current = current.GetFirstChild()
	}
	return !(NodeUtil_IsLiteralValue(current, true))
}

// Java method: AstAnalyzer.constructorCallHasSideEffects(com.google.javascript.rhino.Node newNode)
func (this_ *AstAnalyzer) ConstructorCallHasSideEffects(newNode pkg_rhino.Node) bool {
	_ = newNode
	if !(newNode.IsNew()) {
		j2gort.CheckArgumentFmt(false, "Expected NEW node, got %s", j2gort.BoxNullableVal(newNode.GetToken()))
	}
	if newNode.IsNoSideEffectsCall() {
		return false
	}
	if newNode.IsOnlyModifiesArgumentsCall() && NodeUtil_AllArgsUnescapedLocal(newNode) {
		return false
	}
	var nameNode pkg_rhino.Node = newNode.GetFirstChild()
	_ = nameNode
	return ((!(nameNode.IsName()) || !(this_.assumeKnownBuiltinsArePure)) || !(j2gort.SetContainsStr(StaticGet_AstAnalyzer_CONSTRUCTORS_WITHOUT_SIDE_EFFECTS(), nameNode.GetString())))
}

// Java method: AstAnalyzer.nodeTypeMayHaveSideEffects(com.google.javascript.rhino.Node n)
func (this_ *AstAnalyzer) NodeTypeMayHaveSideEffects(n pkg_rhino.Node) bool {
	_ = n
	if NodeUtil_IsAssignmentOp(n) {
		return true
	}
	switch n.TokenSwitchOrd_() {
	case int32(24) /* Token.DELPROP */, int32(72) /* Token.DEC */, int32(71) /* Token.INC */, int32(114) /* Token.YIELD */, int32(43) /* Token.THROW */, int32(115) /* Token.AWAIT */, int32(124) /* Token.DYNAMIC_IMPORT */ :
		return true
	case int32(81) /* Token.FOR_IN */, int32(112) /* Token.FOR_OF */, int32(113) /* Token.FOR_AWAIT_OF */ :
		return true
	case int32(31) /* Token.OPTCHAIN_CALL */, int32(28) /* Token.CALL */, int32(130) /* Token.TAGGED_TEMPLATELIT */ :
		return this_.FunctionCallHasSideEffects(n)
	case int32(23) /* Token.NEW */ :
		return this_.ConstructorCallHasSideEffects(n)
	case int32(32) /* Token.NAME */ :
		return n.HasChildren()
	case int32(104) /* Token.DESTRUCTURING_LHS */ :
		return true
	case int32(126) /* Token.OBJECT_REST */, int32(128) /* Token.OBJECT_SPREAD */ :
		return !(this_.assumeGettersArePure)
	case int32(125) /* Token.ITER_REST */, int32(127) /* Token.ITER_SPREAD */ :
		return NodeUtil_IteratesImpureIterable(n)
	case int32(100) /* Token.STRING_KEY */ :
		if n.GetParent().IsObjectPattern() {
			return this_.priv_AstAnalyzer_getPropertyKind_609(n.GetString()).HasGetter()
		}
		return false
	case int32(26) /* Token.GETPROP */, int32(29) /* Token.OPTCHAIN_GETPROP */ :
		return this_.priv_AstAnalyzer_getPropertyKind_609(n.GetString()).HasGetterOrSetter()
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java method: AstAnalyzer.getPropertyKind(java.lang.String name)
func (this_ *AstAnalyzer) priv_AstAnalyzer_getPropertyKind_609(name string) *AccessorSummary_PropertyAccessKind {
	_ = name
	if this_.assumeGettersArePure {
		return StaticGet_AccessorSummary_PropertyAccessKind_NORMAL()
	} else {
		return this_.accessorSummary.GetKind(name)
	}
}

// Java record: com.google.javascript.jscomp.AstAnalyzer.Options
type AstAnalyzer_Options struct {
	useTypesForLocalOptimization bool `gson:"useTypesForLocalOptimization"`
	assumeGettersArePure         bool `gson:"assumeGettersArePure"`
	hasRegexpGlobalReferences    bool `gson:"hasRegexpGlobalReferences"`
	assumeKnownBuiltinsArePure   bool `gson:"assumeKnownBuiltinsArePure"`
}

func (this_ *AstAnalyzer_Options) BindSelf_(s any) {
	_ = s
}
func (this_ *AstAnalyzer_Options) IsJavaNil_() bool { return this_ == nil }
func (*AstAnalyzer_Options) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AstAnalyzer$Options", "Options"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AstAnalyzer.Options", "com.google.javascript.jscomp.AstAnalyzer$Options", func(v_ any) bool {
		return j2gort.InstanceOf[*AstAnalyzer_Options](v_)
	})
}
func (*AstAnalyzer_Options) Implements_j2gort_Ext_java_lang_Record() {}

func (this_ *AstAnalyzer_Options) initFields_AstAnalyzer_Options() {
	_ = this_
}

// Java constructor: Options(boolean useTypesForLocalOptimization, boolean assumeGettersArePure, boolean hasRegexpGlobalReferences, boolean assumeKnownBuiltinsArePure)
func NewAstAnalyzer_Options(useTypesForLocalOptimization bool, assumeGettersArePure bool, hasRegexpGlobalReferences bool, assumeKnownBuiltinsArePure bool) *AstAnalyzer_Options {
	this_ := &AstAnalyzer_Options{}
	InitAstAnalyzer_Options(this_, useTypesForLocalOptimization, assumeGettersArePure, hasRegexpGlobalReferences, assumeKnownBuiltinsArePure)
	return this_
}

func InitAstAnalyzer_Options(this_ *AstAnalyzer_Options, useTypesForLocalOptimization bool, assumeGettersArePure bool, hasRegexpGlobalReferences bool, assumeKnownBuiltinsArePure bool) {
	_ = useTypesForLocalOptimization
	_ = assumeGettersArePure
	_ = hasRegexpGlobalReferences
	_ = assumeKnownBuiltinsArePure
	this_.initFields_AstAnalyzer_Options()
	this_.useTypesForLocalOptimization = useTypesForLocalOptimization
	this_.assumeGettersArePure = assumeGettersArePure
	this_.hasRegexpGlobalReferences = hasRegexpGlobalReferences
	this_.assumeKnownBuiltinsArePure = assumeKnownBuiltinsArePure
}

// Java static method: Options.builder()
func AstAnalyzer_Options_S_Builder() AstAnalyzer_Options_Builder {
	return NewAutoBuilder_AstAnalyzer_Options_Builder()
}

func (this_ *AstAnalyzer_Options) UseTypesForLocalOptimization() bool {
	return this_.useTypesForLocalOptimization
}
func (this_ *AstAnalyzer_Options) AssumeGettersArePure() bool { return this_.assumeGettersArePure }
func (this_ *AstAnalyzer_Options) HasRegexpGlobalReferences() bool {
	return this_.hasRegexpGlobalReferences
}
func (this_ *AstAnalyzer_Options) AssumeKnownBuiltinsArePure() bool {
	return this_.assumeKnownBuiltinsArePure
}
func (this_ *AstAnalyzer_Options) Equals(other_ any) bool {
	o_, ok_ := other_.(*AstAnalyzer_Options)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if this_.useTypesForLocalOptimization != o_.useTypesForLocalOptimization {
		return false
	}
	if this_.assumeGettersArePure != o_.assumeGettersArePure {
		return false
	}
	if this_.hasRegexpGlobalReferences != o_.hasRegexpGlobalReferences {
		return false
	}
	if this_.assumeKnownBuiltinsArePure != o_.assumeKnownBuiltinsArePure {
		return false
	}
	return true
}
func (this_ *AstAnalyzer_Options) HashCode() int32 {
	var h_ int32 = 0
	if this_.useTypesForLocalOptimization {
		h_ = 31*h_ + 1231
	} else {
		h_ = 31*h_ + 1237
	}
	if this_.assumeGettersArePure {
		h_ = 31*h_ + 1231
	} else {
		h_ = 31*h_ + 1237
	}
	if this_.hasRegexpGlobalReferences {
		h_ = 31*h_ + 1231
	} else {
		h_ = 31*h_ + 1237
	}
	if this_.assumeKnownBuiltinsArePure {
		h_ = 31*h_ + 1231
	} else {
		h_ = 31*h_ + 1237
	}
	return h_
}
func (this_ *AstAnalyzer_Options) String() string {
	_ = fmt.Sprint
	return "Options[" + "useTypesForLocalOptimization=" + j2gort.StrVal(j2gort.ValueOfBoolean(this_.useTypesForLocalOptimization)) + ", " + "assumeGettersArePure=" + j2gort.StrVal(j2gort.ValueOfBoolean(this_.assumeGettersArePure)) + ", " + "hasRegexpGlobalReferences=" + j2gort.StrVal(j2gort.ValueOfBoolean(this_.hasRegexpGlobalReferences)) + ", " + "assumeKnownBuiltinsArePure=" + j2gort.StrVal(j2gort.ValueOfBoolean(this_.assumeKnownBuiltinsArePure)) + "]"
}

// Java interface: com.google.javascript.jscomp.AstAnalyzer.Options.Builder
type AstAnalyzer_Options_Builder interface {
	Implements_jscomp_AstAnalyzer_Options_Builder()
	IsJavaNil_() bool
	SetUseTypesForLocalOptimization(bool) AstAnalyzer_Options_Builder
	SetAssumeGettersArePure(bool) AstAnalyzer_Options_Builder
	SetHasRegexpGlobalReferences(bool) AstAnalyzer_Options_Builder
	SetAssumeKnownBuiltinsArePure(bool) AstAnalyzer_Options_Builder
	Build() *AstAnalyzer_Options
}

func ToIface_AstAnalyzer_Options_Builder[T interface {
	AstAnalyzer_Options_Builder
	comparable
}](v T) AstAnalyzer_Options_Builder {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_AstAnalyzer_Options_Builder[T AstAnalyzer_Options_Builder](v T) AstAnalyzer_Options_Builder {
	var vi AstAnalyzer_Options_Builder = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AstAnalyzer.Options.Builder", "com.google.javascript.jscomp.AstAnalyzer$Options$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[AstAnalyzer_Options_Builder](v_)
	})
}
