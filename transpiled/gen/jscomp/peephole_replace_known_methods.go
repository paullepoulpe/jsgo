/*
 * Copyright 2011 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/PeepholeReplaceKnownMethods.java (generated mechanically by j2go).

package jscomp

import (
	"math"
	"sync/atomic"

	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	pkg_colors "github.com/paullepoulpe/jsgo/transpiled/gen/colors"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.PeepholeReplaceKnownMethods
type PeepholeReplaceKnownMethods struct {
	AbstractPeepholeOptimizationBase
	late     bool `gson:"late"`
	useTypes bool `gson:"useTypes"`
}

func (this_ *PeepholeReplaceKnownMethods) BindSelf_(s any) {
	_ = s
	this_.AbstractPeepholeOptimizationBase.BindSelf_(s)
}
func (this_ *PeepholeReplaceKnownMethods) IsJavaNil_() bool { return this_ == nil }
func (*PeepholeReplaceKnownMethods) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PeepholeReplaceKnownMethods", "PeepholeReplaceKnownMethods"
}
func (this_ *PeepholeReplaceKnownMethods) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.AbstractPeepholeOptimizationBase.IdentityHash_()
}
func (*PeepholeReplaceKnownMethods) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PeepholeReplaceKnownMethods", "com.google.javascript.jscomp.PeepholeReplaceKnownMethods", func(v_ any) bool {
		return j2gort.InstanceOf[*PeepholeReplaceKnownMethods](v_)
	})
}
func (*PeepholeReplaceKnownMethods) Implements_jscomp_AbstractPeepholeOptimization() {}

func StaticGet_PeepholeReplaceKnownMethods_ECMA_WHITESPACE() string {
	return "[ \t\n-\r\\u00A0\\u1680\\u2000-\\u200A\\u2028\\u2029\\u202F\\u205F\\u3000\\uFEFF]+"
}

var static_PeepholeReplaceKnownMethods_TRIM_BOTH_PATTERN *j2gort.Ext_java_util_regex_Pattern

func StaticGet_PeepholeReplaceKnownMethods_TRIM_BOTH_PATTERN() *j2gort.Ext_java_util_regex_Pattern {
	clinit_PeepholeReplaceKnownMethods()
	return static_PeepholeReplaceKnownMethods_TRIM_BOTH_PATTERN
}
func StaticSet_PeepholeReplaceKnownMethods_TRIM_BOTH_PATTERN(v *j2gort.Ext_java_util_regex_Pattern) {
	clinit_PeepholeReplaceKnownMethods()
	static_PeepholeReplaceKnownMethods_TRIM_BOTH_PATTERN = v
}

var static_PeepholeReplaceKnownMethods_TRIM_START_PATTERN *j2gort.Ext_java_util_regex_Pattern

func StaticGet_PeepholeReplaceKnownMethods_TRIM_START_PATTERN() *j2gort.Ext_java_util_regex_Pattern {
	clinit_PeepholeReplaceKnownMethods()
	return static_PeepholeReplaceKnownMethods_TRIM_START_PATTERN
}
func StaticSet_PeepholeReplaceKnownMethods_TRIM_START_PATTERN(v *j2gort.Ext_java_util_regex_Pattern) {
	clinit_PeepholeReplaceKnownMethods()
	static_PeepholeReplaceKnownMethods_TRIM_START_PATTERN = v
}

var static_PeepholeReplaceKnownMethods_TRIM_END_PATTERN *j2gort.Ext_java_util_regex_Pattern

func StaticGet_PeepholeReplaceKnownMethods_TRIM_END_PATTERN() *j2gort.Ext_java_util_regex_Pattern {
	clinit_PeepholeReplaceKnownMethods()
	return static_PeepholeReplaceKnownMethods_TRIM_END_PATTERN
}
func StaticSet_PeepholeReplaceKnownMethods_TRIM_END_PATTERN(v *j2gort.Ext_java_util_regex_Pattern) {
	clinit_PeepholeReplaceKnownMethods()
	static_PeepholeReplaceKnownMethods_TRIM_END_PATTERN = v
}

var clinitDone_PeepholeReplaceKnownMethods uint32
var clinitState_PeepholeReplaceKnownMethods j2gort.ClinitState

func clinit_PeepholeReplaceKnownMethods() {
	if atomic.LoadUint32(&clinitDone_PeepholeReplaceKnownMethods) != 0 {
		return
	}
	clinitSlow_PeepholeReplaceKnownMethods()
}

func clinitSlow_PeepholeReplaceKnownMethods() {
	if !j2gort.EnterClinit(&clinitState_PeepholeReplaceKnownMethods, &clinitDone_PeepholeReplaceKnownMethods) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_PeepholeReplaceKnownMethods, &clinitDone_PeepholeReplaceKnownMethods)
	static_PeepholeReplaceKnownMethods_TRIM_BOTH_PATTERN = j2gort.Ext_java_util_regex_Pattern_Compile_1(("^" + "[ \t\n-\r\\u00A0\\u1680\\u2000-\\u200A\\u2028\\u2029\\u202F\\u205F\\u3000\\uFEFF]+" + "|" + "[ \t\n-\r\\u00A0\\u1680\\u2000-\\u200A\\u2028\\u2029\\u202F\\u205F\\u3000\\uFEFF]+" + "\\z"))
	static_PeepholeReplaceKnownMethods_TRIM_START_PATTERN = j2gort.Ext_java_util_regex_Pattern_Compile_1(("^" + "[ \t\n-\r\\u00A0\\u1680\\u2000-\\u200A\\u2028\\u2029\\u202F\\u205F\\u3000\\uFEFF]+"))
	static_PeepholeReplaceKnownMethods_TRIM_END_PATTERN = j2gort.Ext_java_util_regex_Pattern_Compile_1(("[ \t\n-\r\\u00A0\\u1680\\u2000-\\u200A\\u2028\\u2029\\u202F\\u205F\\u3000\\uFEFF]+" + "\\z"))
}

func (this_ *PeepholeReplaceKnownMethods) initFields_PeepholeReplaceKnownMethods() {
	_ = this_
}

// Java constructor: PeepholeReplaceKnownMethods(boolean late, boolean useTypes)
func NewPeepholeReplaceKnownMethods(late bool, useTypes bool) *PeepholeReplaceKnownMethods {
	clinit_PeepholeReplaceKnownMethods()
	this_ := &PeepholeReplaceKnownMethods{}
	this_.BindSelf_(this_)
	InitPeepholeReplaceKnownMethods(this_, late, useTypes)
	return this_
}

func InitPeepholeReplaceKnownMethods(this_ *PeepholeReplaceKnownMethods, late bool, useTypes bool) {
	_ = late
	_ = useTypes
	InitAbstractPeepholeOptimization(&this_.AbstractPeepholeOptimizationBase)
	this_.initFields_PeepholeReplaceKnownMethods()
	this_.late = late
	this_.useTypes = useTypes
}

// Java method: PeepholeReplaceKnownMethods.optimizeSubtree(com.google.javascript.rhino.Node subtree)
func (this_ *PeepholeReplaceKnownMethods) OptimizeSubtree(subtree pkg_rhino.Node) pkg_rhino.Node {
	_ = subtree
	if NodeUtil_IsGoogWeakUsageCall(subtree) && this_.late {
		var name pkg_rhino.Node = subtree.GetSecondChild().Detach()
		_ = name
		subtree.ReplaceWith(name)
		this_.ReportChangeToEnclosingScope(name)
		return name
	} else if subtree.IsCall() {
		return this_.priv_PeepholeReplaceKnownMethods_tryFoldKnownMethods_4253(subtree)
	}
	return subtree
}

// Java method: PeepholeReplaceKnownMethods.tryFoldKnownMethods(com.google.javascript.rhino.Node subtree)
func (this_ *PeepholeReplaceKnownMethods) priv_PeepholeReplaceKnownMethods_tryFoldKnownMethods_4253(subtree pkg_rhino.Node) pkg_rhino.Node {
	_ = subtree
	if !(subtree.IsCall()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(subtree))
	}
	subtree = this_.priv_PeepholeReplaceKnownMethods_tryFoldArrayJoin_4269(subtree)
	if subtree.IsCall() {
		subtree = this_.priv_PeepholeReplaceKnownMethods_tryToFoldArrayConcat_4279(subtree)
		if !(subtree.IsCall()) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(subtree))
		}
		var callTarget pkg_rhino.Node = j2gort.CheckNotNullNode(subtree.GetFirstChild())
		_ = callTarget
		if callTarget.IsGetProp() {
			if this_.IsASTNormalized() && callTarget.GetFirstChild().IsQualifiedName() {
				switch j2gort.CheckNotNullString(callTarget.GetFirstChild().GetQualifiedName()) {
				case "Array":
					return this_.priv_PeepholeReplaceKnownMethods_tryFoldKnownArrayMethods_4254(subtree, callTarget)
				case "Math":
					return this_.priv_PeepholeReplaceKnownMethods_tryFoldKnownMathMethods_4257(subtree, callTarget)
				case "Number":
					return this_.priv_PeepholeReplaceKnownMethods_tryFoldKnownNumberMethods_4256(subtree, callTarget)
				default:
				}
			}
			subtree = this_.priv_PeepholeReplaceKnownMethods_tryFoldKnownStringMethods_4258(subtree, callTarget)
		} else if callTarget.IsName() {
			subtree = this_.priv_PeepholeReplaceKnownMethods_tryFoldKnownNumericMethods_4259(subtree, callTarget)
		}
	}
	return subtree
}

// Java method: PeepholeReplaceKnownMethods.tryFoldKnownArrayMethods(com.google.javascript.rhino.Node subtree, com.google.javascript.rhino.Node callTarget)
func (this_ *PeepholeReplaceKnownMethods) priv_PeepholeReplaceKnownMethods_tryFoldKnownArrayMethods_4254(subtree pkg_rhino.Node, callTarget pkg_rhino.Node) pkg_rhino.Node {
	_ = subtree
	_ = callTarget
	if !(subtree.IsCall() && callTarget.IsGetProp()) {
		j2gort.CheckArgument(false)
	}
	var methodName string = callTarget.GetString()
	_ = methodName
	switch j2gort.CheckNotNullString(methodName) {
	case "of":
		subtree.RemoveFirstChild()
		var arraylit pkg_rhino.Node = pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_ARRAYLIT())
		_ = arraylit
		arraylit.AddChildrenToBack(subtree.RemoveChildren())
		subtree.ReplaceWith(arraylit)
		this_.ReportChangeToEnclosingScope(arraylit)
		return arraylit
	case "isArray":
		return this_.priv_PeepholeReplaceKnownMethods_tryFoldArrayIsArray_4255(subtree, callTarget)
	default:
		return subtree
	}
	panic("unreachable switch expression")
}

// Java method: PeepholeReplaceKnownMethods.tryFoldArrayIsArray(com.google.javascript.rhino.Node subtree, com.google.javascript.rhino.Node callTarget)
func (this_ *PeepholeReplaceKnownMethods) priv_PeepholeReplaceKnownMethods_tryFoldArrayIsArray_4255(subtree pkg_rhino.Node, callTarget pkg_rhino.Node) pkg_rhino.Node {
	_ = subtree
	_ = callTarget
	if !(subtree.IsCall() && callTarget.IsGetProp()) {
		j2gort.CheckArgument(false)
	}
	var firstArg pkg_rhino.Node = callTarget.GetNext()
	_ = firstArg
	if firstArg == 0 {
		var newNode pkg_rhino.Node = NodeUtil_BooleanNode(false).Srcref(subtree)
		_ = newNode
		subtree.ReplaceWith(newNode)
		this_.ReportChangeToEnclosingScope(newNode)
		return newNode
	}
	if firstArg.GetNext() != 0 {
		return subtree
	}
	if this_.MayHaveSideEffects(firstArg) {
		return subtree
	}
	var isArray bool = false
	_ = isArray
	if firstArg.IsArrayLit() {
		isArray = true
	} else if (((NodeUtil_IsImmutableValue(firstArg) || firstArg.IsObjectLit()) || firstArg.IsRegExp()) || firstArg.IsFunction()) || firstArg.IsClass() {
		isArray = false
	} else {
		return subtree
	}
	this_.MarkFunctionsDeleted(subtree)
	var newNode_2 pkg_rhino.Node = NodeUtil_BooleanNode(isArray).Srcref(subtree)
	_ = newNode_2
	subtree.ReplaceWith(newNode_2)
	this_.ReportChangeToEnclosingScope(newNode_2)
	return newNode_2
}

// Java method: PeepholeReplaceKnownMethods.tryFoldKnownNumberMethods(com.google.javascript.rhino.Node subtree, com.google.javascript.rhino.Node callTarget)
func (this_ *PeepholeReplaceKnownMethods) priv_PeepholeReplaceKnownMethods_tryFoldKnownNumberMethods_4256(subtree pkg_rhino.Node, callTarget pkg_rhino.Node) pkg_rhino.Node {
	_ = subtree
	_ = callTarget
	if !(subtree.IsCall() && callTarget.IsGetProp()) {
		j2gort.CheckArgument(false)
	}
	var methodName string = callTarget.GetString()
	_ = methodName
	switch j2gort.CheckNotNullString(methodName) {
	case "parseInt", "parseFloat":
		var firstArg pkg_rhino.Node = callTarget.GetNext()
		_ = firstArg
		if ((firstArg != 0) && firstArg.IsStringLit()) || NodeUtil_IsNumericLiteral(firstArg) {
			return this_.priv_PeepholeReplaceKnownMethods_tryFoldParseNumber_4263(subtree, methodName, firstArg)
		}
		return subtree
	default:
	}
	var valueNode pkg_rhino.Node = callTarget.GetNext()
	_ = valueNode
	if (valueNode == 0) || (valueNode.GetNext() != 0) {
		return subtree
	}
	var onlyArg *j2gort.BoxDouble = this_.GetSideEffectFreeNumberValueNoConversion(valueNode)
	_ = onlyArg
	var replacement pkg_rhino.Node = pkg_rhino.Node(0)
	_ = replacement
	if onlyArg != nil {
		switch j2gort.CheckNotNullString(methodName) {
		case "isFinite":
			replacement = NodeUtil_BooleanNode(func() bool { df_1 := j2gort.UnboxDouble(onlyArg); return !math.IsNaN(df_1) && !math.IsInf(df_1, 0) }())
		case "isNaN":
			replacement = NodeUtil_BooleanNode(math.IsNaN(j2gort.UnboxDouble(onlyArg)))
		case "isSafeInteger":
			replacement = NodeUtil_BooleanNode((((j2gort.UnboxDouble(onlyArg) < math.Pow(float64(int32(2)), float64(int32(53)))) && (j2gort.UnboxDouble(onlyArg) >= (float64(int32(-1)) * (math.Pow(float64(int32(2)), float64(int32(53))) - float64(int32(1)))))) && (float64(onlyArg.LongValue()) == onlyArg.DoubleValue())))
		default:
		}
	} else if !(this_.MayHaveSideEffects(valueNode)) && (NodeUtil_IsImmutableValue(valueNode) || NodeUtil_IsLiteralValue(valueNode, false)) {
		switch j2gort.CheckNotNullString(methodName) {
		case "isFinite", "isNaN", "isSafeInteger":
			replacement = NodeUtil_BooleanNode(false)
		default:
		}
	}
	if replacement != 0 {
		subtree.ReplaceWith(replacement)
		this_.ReportChangeToEnclosingScope(replacement)
		return replacement
	}
	return subtree
}

// Java method: PeepholeReplaceKnownMethods.tryFoldKnownMathMethods(com.google.javascript.rhino.Node subtree, com.google.javascript.rhino.Node callTarget)
func (this_ *PeepholeReplaceKnownMethods) priv_PeepholeReplaceKnownMethods_tryFoldKnownMathMethods_4257(subtree pkg_rhino.Node, callTarget pkg_rhino.Node) pkg_rhino.Node {
	_ = subtree
	_ = callTarget
	if !(subtree.IsCall() && callTarget.IsGetProp()) {
		j2gort.CheckArgument(false)
	}
	var args j2gort.Ext_java_util_List = j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())
	_ = args
	{
		var arg pkg_rhino.Node = callTarget.GetNext()
		_ = arg
		for ; arg != 0; arg = arg.GetNext() {
			var d *j2gort.BoxDouble = this_.GetSideEffectFreeNumberValue(arg)
			_ = d
			if d != nil {
				if args.IsEmpty() {
					args = j2gort.NewExt_java_util_ArrayList()
				}
				args.Add_1(j2gort.PtrToAny(d))
			} else {
				return subtree
			}
		}
	}
	var replacement *j2gort.BoxDouble = nil
	_ = replacement
	var methodName string = callTarget.GetString()
	_ = methodName
	if args.Size() == int32(1) {
		var arg_2 float64 = j2gort.UnboxDouble(j2gort.CastPtr[j2gort.BoxDouble](args.Get(int32(0))))
		_ = arg_2
		switch j2gort.CheckNotNullString(methodName) {
		case "abs":
			replacement = j2gort.ValueOfDouble(math.Abs(arg_2))
		case "ceil":
			replacement = j2gort.ValueOfDouble(math.Ceil(arg_2))
		case "floor":
			replacement = j2gort.ValueOfDouble(math.Floor(arg_2))
		case "fround":
			if (math.IsNaN(arg_2) || math.IsInf(arg_2, 0)) || (arg_2 == float64(int32(0))) {
				replacement = j2gort.ValueOfDouble(arg_2)
			} else if float64(float32(arg_2)) == arg_2 {
				replacement = j2gort.ValueOfDouble(float64(float32(arg_2)))
			} else {
				replacement = nil
			}
		case "round":
			if math.IsNaN(arg_2) || math.IsInf(arg_2, 0) {
				replacement = j2gort.ValueOfDouble(arg_2)
			} else if (arg_2 >= -(float64(0.5))) && ((arg_2 < float64(0.0)) || ((arg_2 == float64(0.0)) && (int64(math.Float64bits(arg_2)) < int64(int32(0))))) {
				replacement = j2gort.ValueOfDouble(-(float64(0.0)))
			} else if arg_2 == math.Floor(arg_2) {
				replacement = j2gort.ValueOfDouble(arg_2)
			} else {
				replacement = j2gort.ValueOfDouble(float64(j2gort.Ext_java_lang_Math_Round_1_float64(arg_2)))
			}
		case "sign":
			replacement = j2gort.ValueOfDouble(j2gort.Ext_java_lang_Math_Signum_1_float64(arg_2))
		case "trunc":
			if math.IsNaN(arg_2) || math.IsInf(arg_2, 0) {
				replacement = j2gort.ValueOfDouble(arg_2)
			} else {
				replacement = j2gort.ValueOfDouble((j2gort.Ext_java_lang_Math_Signum_1_float64(arg_2) * math.Floor(math.Abs(arg_2))))
			}
		case "clz32":
			replacement = j2gort.ValueOfDouble(float64(j2gort.IntLeadingZeros(pkg_base.JSCompDoubles_EcmascriptToUint32(arg_2))))
		default:
		}
	}
	if replacement == nil {
		switch j2gort.CheckNotNullString(methodName) {
		case "max":
			var result float64 = math.Inf(-1)
			_ = result
			{
				cur_2 := j2gort.NewIterCursor(j2gort.CheckNotNull(args))
				for cur_2.HasNext() {
					rawElem_3 := cur_2.Next()
					d_2 := j2gort.CastPtr[j2gort.BoxDouble](rawElem_3)
					_ = d_2
					result = math.Max(result, j2gort.UnboxDouble(d_2))
				}
			}
			replacement = j2gort.ValueOfDouble(result)
		case "min":
			var result_2 float64 = math.Inf(1)
			_ = result_2
			{
				cur_4 := j2gort.NewIterCursor(j2gort.CheckNotNull(args))
				for cur_4.HasNext() {
					rawElem_5 := cur_4.Next()
					d_3 := j2gort.CastPtr[j2gort.BoxDouble](rawElem_5)
					_ = d_3
					result_2 = math.Min(result_2, j2gort.UnboxDouble(d_3))
				}
			}
			replacement = j2gort.ValueOfDouble(result_2)
		case "imul":
			if args.Size() < int32(2) {
				replacement = j2gort.ValueOfDouble(float64(0.0))
			} else {
				replacement = j2gort.ValueOfDouble(float64((pkg_base.JSCompDoubles_EcmascriptToInt32(j2gort.UnboxDouble(j2gort.CastPtr[j2gort.BoxDouble](args.Get(int32(0))))) * pkg_base.JSCompDoubles_EcmascriptToInt32(j2gort.UnboxDouble(j2gort.CastPtr[j2gort.BoxDouble](args.Get(int32(1))))))))
			}
		case "pow":
			if args.Size() == int32(2) {
				replacement = j2gort.ValueOfDouble(math.Pow(j2gort.UnboxDouble(j2gort.CastPtr[j2gort.BoxDouble](args.Get(int32(0)))), j2gort.UnboxDouble(j2gort.CastPtr[j2gort.BoxDouble](args.Get(int32(1))))))
			}
		default:
		}
	}
	if replacement != nil {
		var numberNode pkg_rhino.Node = NodeUtil_NumberNode(j2gort.UnboxDouble(replacement), subtree)
		_ = numberNode
		subtree.ReplaceWith(numberNode)
		this_.ReportChangeToEnclosingScope(numberNode)
		return numberNode
	}
	return subtree
}

// Java method: PeepholeReplaceKnownMethods.tryFoldKnownStringMethods(com.google.javascript.rhino.Node subtree, com.google.javascript.rhino.Node callTarget)
func (this_ *PeepholeReplaceKnownMethods) priv_PeepholeReplaceKnownMethods_tryFoldKnownStringMethods_4258(subtree pkg_rhino.Node, callTarget pkg_rhino.Node) pkg_rhino.Node {
	_ = subtree
	_ = callTarget
	if !(subtree.IsCall() && callTarget.IsGetProp()) {
		j2gort.CheckArgument(false)
	}
	var stringNode pkg_rhino.Node = callTarget.GetFirstChild()
	_ = stringNode
	var isStringLiteral bool = stringNode.IsStringLit()
	_ = isStringLiteral
	var functionNameString string = callTarget.GetString()
	_ = functionNameString
	var firstArg pkg_rhino.Node = callTarget.GetNext()
	_ = firstArg
	if isStringLiteral {
		if j2gort.CheckNotNullString(functionNameString) == "split" {
			return this_.priv_PeepholeReplaceKnownMethods_tryFoldStringSplit_4278(subtree, stringNode, firstArg)
		} else if firstArg == 0 {
			switch j2gort.CheckNotNullString(functionNameString) {
			case "toLowerCase":
				return this_.priv_PeepholeReplaceKnownMethods_tryFoldStringToLowerCase_4260(subtree, stringNode)
			case "toUpperCase":
				return this_.priv_PeepholeReplaceKnownMethods_tryFoldStringToUpperCase_4261(subtree, stringNode)
			case "trim", "trimStart", "trimEnd", "trimLeft", "trimRight":
				return this_.priv_PeepholeReplaceKnownMethods_tryFoldStringTrim_4262(subtree, stringNode, functionNameString)
			case "at":
				return this_.priv_PeepholeReplaceKnownMethods_tryFoldStringAt_4268(subtree, stringNode, pkg_rhino.Node(0))
			default:
			}
		} else {
			if NodeUtil_IsImmutableValue(firstArg) {
				switch j2gort.CheckNotNullString(functionNameString) {
				case "indexOf", "lastIndexOf":
					return this_.priv_PeepholeReplaceKnownMethods_tryFoldStringIndexOf_4264(subtree, functionNameString, stringNode, firstArg)
				case "includes":
					return this_.priv_PeepholeReplaceKnownMethods_tryFoldStringIncludes_4265(subtree, stringNode, firstArg)
				case "startsWith":
					return this_.priv_PeepholeReplaceKnownMethods_tryFoldStringStartsWith_4266(subtree, stringNode, firstArg)
				case "endsWith":
					return this_.priv_PeepholeReplaceKnownMethods_tryFoldStringEndsWith_4267(subtree, stringNode, firstArg)
				case "at":
					return this_.priv_PeepholeReplaceKnownMethods_tryFoldStringAt_4268(subtree, stringNode, firstArg)
				case "substr":
					return this_.priv_PeepholeReplaceKnownMethods_tryFoldStringSubstr_4270(subtree, stringNode, firstArg)
				case "substring", "slice":
					return this_.priv_PeepholeReplaceKnownMethods_tryFoldStringSubstringOrSlice_4271(subtree, stringNode, firstArg)
				case "charAt":
					return this_.priv_PeepholeReplaceKnownMethods_tryFoldStringCharAt_4273(subtree, stringNode, firstArg)
				case "charCodeAt":
					return this_.priv_PeepholeReplaceKnownMethods_tryFoldStringCharCodeAt_4276(subtree, stringNode, firstArg)
				case "replace":
					return this_.priv_PeepholeReplaceKnownMethods_tryFoldStringReplace_4274(subtree, stringNode, firstArg)
				case "replaceAll":
					return this_.priv_PeepholeReplaceKnownMethods_tryFoldStringReplaceAll_4275(subtree, stringNode, firstArg)
				default:
				}
			}
		}
	}
	if (this_.useTypes && (firstArg != 0)) && (isStringLiteral || pkg_colors.StaticGet_StandardColors_STRING().Equals(j2gort.PtrToAny(stringNode.GetColor()))) {
		if subtree.HasXChildren(int32(3)) {
			var maybeStart *j2gort.BoxDouble = this_.GetSideEffectFreeNumberValue(firstArg)
			_ = maybeStart
			if maybeStart != nil {
				var start int32 = maybeStart.IntValue()
				_ = start
				var maybeLengthOrEnd *j2gort.BoxDouble = this_.GetSideEffectFreeNumberValue(firstArg.GetNext())
				_ = maybeLengthOrEnd
				if maybeLengthOrEnd != nil {
					switch j2gort.CheckNotNullString(functionNameString) {
					case "substr":
						var length int32 = maybeLengthOrEnd.IntValue()
						_ = length
						if (start >= int32(0)) && (length == int32(1)) {
							return this_.priv_PeepholeReplaceKnownMethods_replaceWithCharAt_4272(subtree, callTarget, firstArg)
						}
					case "substring", "slice":
						var end int32 = maybeLengthOrEnd.IntValue()
						_ = end
						if (start >= int32(0)) && ((end - start) == int32(1)) {
							return this_.priv_PeepholeReplaceKnownMethods_replaceWithCharAt_4272(subtree, callTarget, firstArg)
						}
					default:
					}
				}
			}
		}
	}
	return subtree
}

// Java method: PeepholeReplaceKnownMethods.tryFoldKnownNumericMethods(com.google.javascript.rhino.Node subtree, com.google.javascript.rhino.Node callTarget)
func (this_ *PeepholeReplaceKnownMethods) priv_PeepholeReplaceKnownMethods_tryFoldKnownNumericMethods_4259(subtree pkg_rhino.Node, callTarget pkg_rhino.Node) pkg_rhino.Node {
	_ = subtree
	_ = callTarget
	if !(subtree.IsCall()) {
		j2gort.CheckArgument(false)
	}
	if this_.IsASTNormalized() {
		var functionNameString string = callTarget.GetString()
		_ = functionNameString
		var firstArgument pkg_rhino.Node = callTarget.GetNext()
		_ = firstArgument
		if ((firstArgument != 0) && (firstArgument.IsStringLit() || NodeUtil_IsNumericLiteral(firstArgument))) && ((j2gort.CheckNotNullString(functionNameString) == "parseInt") || (j2gort.CheckNotNullString(functionNameString) == "parseFloat")) {
			subtree = this_.priv_PeepholeReplaceKnownMethods_tryFoldParseNumber_4263(subtree, functionNameString, firstArgument)
		}
	}
	return subtree
}

// Java method: PeepholeReplaceKnownMethods.tryFoldStringToLowerCase(com.google.javascript.rhino.Node subtree, com.google.javascript.rhino.Node stringNode)
func (this_ *PeepholeReplaceKnownMethods) priv_PeepholeReplaceKnownMethods_tryFoldStringToLowerCase_4260(subtree pkg_rhino.Node, stringNode pkg_rhino.Node) pkg_rhino.Node {
	_ = subtree
	_ = stringNode
	var lowered string = javalang.StringToLowerCaseASCII(stringNode.GetString())
	_ = lowered
	var replacement pkg_rhino.Node = pkg_rhino.IR_String(lowered)
	_ = replacement
	subtree.ReplaceWith(replacement)
	this_.ReportChangeToEnclosingScope(replacement)
	return replacement
}

// Java method: PeepholeReplaceKnownMethods.tryFoldStringToUpperCase(com.google.javascript.rhino.Node subtree, com.google.javascript.rhino.Node stringNode)
func (this_ *PeepholeReplaceKnownMethods) priv_PeepholeReplaceKnownMethods_tryFoldStringToUpperCase_4261(subtree pkg_rhino.Node, stringNode pkg_rhino.Node) pkg_rhino.Node {
	_ = subtree
	_ = stringNode
	var upped string = javalang.StringToUpperCaseASCII(stringNode.GetString())
	_ = upped
	var replacement pkg_rhino.Node = pkg_rhino.IR_String(upped)
	_ = replacement
	subtree.ReplaceWith(replacement)
	this_.ReportChangeToEnclosingScope(replacement)
	return replacement
}

// Java method: PeepholeReplaceKnownMethods.tryFoldStringTrim(com.google.javascript.rhino.Node subtree, com.google.javascript.rhino.Node stringNode, java.lang.String functionName)
func (this_ *PeepholeReplaceKnownMethods) priv_PeepholeReplaceKnownMethods_tryFoldStringTrim_4262(subtree pkg_rhino.Node, stringNode pkg_rhino.Node, functionName string) pkg_rhino.Node {
	_ = subtree
	_ = stringNode
	_ = functionName
	var result string = j2gort.NullString
	_ = result
	switch j2gort.CheckNotNullString(functionName) {
	case "trim":
		result = StaticGet_PeepholeReplaceKnownMethods_TRIM_BOTH_PATTERN().Matcher(j2gort.BoxString(stringNode.GetString())).ReplaceAll_1_string("")
	case "trimStart", "trimLeft":
		result = StaticGet_PeepholeReplaceKnownMethods_TRIM_START_PATTERN().Matcher(j2gort.BoxString(stringNode.GetString())).ReplaceAll_1_string("")
	case "trimEnd", "trimRight":
		result = StaticGet_PeepholeReplaceKnownMethods_TRIM_END_PATTERN().Matcher(j2gort.BoxString(stringNode.GetString())).ReplaceAll_1_string("")
	default:
		result = j2gort.NullString
	}
	if result == j2gort.NullString {
		return subtree
	}
	var replacement pkg_rhino.Node = pkg_rhino.IR_String(result).Srcref(subtree)
	_ = replacement
	subtree.ReplaceWith(replacement)
	this_.ReportChangeToEnclosingScope(replacement)
	return replacement
}

// Java static method: PeepholeReplaceKnownMethods.normalizeNumericString(java.lang.String input)
func PeepholeReplaceKnownMethods_NormalizeNumericString(input string) string {
	clinit_PeepholeReplaceKnownMethods()
	_ = input
	if j2gort.Ext_com_google_common_base_Strings_IsNullOrEmpty(input) {
		return input
	}
	var startIndex int32 = int32(0)
	_ = startIndex
	var endIndex int32 = (j2gort.JStrLen(input) - int32(1))
	_ = endIndex
	for ((startIndex < j2gort.JStrLen(input)) && (int32(j2gort.JStrCharAt(input, startIndex)) == int32(uint16(48)))) && (int32(j2gort.JStrCharAt(input, startIndex)) != int32(uint16(46))) {
		startIndex = startIndex + 1
	}
	if j2gort.JStrIndexOfCp(input, int32(uint16(46))) >= int32(0) {
		for (endIndex >= int32(0)) && (int32(j2gort.JStrCharAt(input, endIndex)) == int32(uint16(48))) {
			endIndex = endIndex - 1
		}
		if int32(j2gort.JStrCharAt(input, endIndex)) == int32(uint16(46)) {
			endIndex = endIndex - 1
		}
	}
	if startIndex >= endIndex {
		return input
	}
	return j2gort.JStrSub(input, startIndex, (endIndex + int32(1)))
}

// Java method: PeepholeReplaceKnownMethods.tryFoldParseNumber(com.google.javascript.rhino.Node n, java.lang.String functionName, com.google.javascript.rhino.Node firstArg)
func (this_ *PeepholeReplaceKnownMethods) priv_PeepholeReplaceKnownMethods_tryFoldParseNumber_4263(n pkg_rhino.Node, functionName string, firstArg pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	_ = functionName
	_ = firstArg
	if !(n.IsCall()) {
		j2gort.CheckArgument(false)
	}
	var isParseInt bool = (j2gort.CheckNotNullString(functionName) == "parseInt")
	_ = isParseInt
	var secondArg pkg_rhino.Node = firstArg.GetNext()
	_ = secondArg
	var radix int32 = int32(0)
	_ = radix
	if secondArg != 0 {
		if !(isParseInt) {
			return n
		}
		if (secondArg.GetNext() != 0) || !(secondArg.IsNumber()) {
			return n
		} else {
			var tmpRadix float64 = secondArg.GetDouble()
			_ = tmpRadix
			if tmpRadix != float64(javalang.D2I(tmpRadix)) {
				return n
			}
			radix = javalang.D2I(tmpRadix)
			if ((radix < int32(0)) || (radix == int32(1))) || (radix > int32(36)) {
				return n
			}
		}
	}
	if !(isParseInt) {
		var stringVal string = j2gort.NullString
		_ = stringVal
		if NodeUtil_IsNumericLiteral(firstArg) {
			var checkVal *j2gort.BoxDouble = this_.GetSideEffectFreeNumberValue(firstArg)
			_ = checkVal
			var numericNode pkg_rhino.Node = NodeUtil_NumberNode(j2gort.UnboxDouble(checkVal), n)
			_ = numericNode
			n.ReplaceWith(numericNode)
			this_.ReportChangeToEnclosingScope(numericNode)
			return numericNode
		} else {
			stringVal = this_.GetSideEffectFreeStringValue(firstArg)
			if stringVal == j2gort.NullString {
				return n
			}
			var checkVal_2 *j2gort.BoxDouble = NodeUtil_GetStringNumberValue(stringVal)
			_ = checkVal_2
			if checkVal_2 == nil {
				return n
			}
			stringVal = NodeUtil_TrimJsWhiteSpace(stringVal)
			if j2gort.CheckNotNullString(stringVal) == "" {
				return n
			}
			if j2gort.CheckNotNullString(stringVal) == "0" {
				var newNode pkg_rhino.Node = pkg_rhino.IR_Number(float64(int32(0)))
				_ = newNode
				n.ReplaceWith(newNode)
				this_.ReportChangeToEnclosingScope(newNode)
				return newNode
			}
			var normalizedNewVal string = "0"
			_ = normalizedNewVal
			var didRet_6 bool
			var retVal_7 pkg_rhino.Node
			func() {
				defer func() {
					if rec_8 := recover(); rec_8 != nil {
						jex_9 := javalang.RecoverException(rec_8)
						_ = jex_9
						if j2gort.InstanceOf[javalang.INumberFormatException](jex_9) {
							e := j2gort.AsIface[javalang.INumberFormatException](jex_9)
							_ = e
							retVal_7 = n
							didRet_6 = true
							return
						} else {
							panic(rec_8)
						}
					}
				}()
				var newVal float64 = j2gort.ParseDouble(stringVal)
				_ = newVal
				var newNode_2 pkg_rhino.Node = NodeUtil_NumberNode(newVal, n)
				_ = newNode_2
				normalizedNewVal = PeepholeReplaceKnownMethods_NormalizeNumericString(javalang.DoubleToString(newVal))
				if !(j2gort.CheckNotNullString(PeepholeReplaceKnownMethods_NormalizeNumericString(stringVal)) == normalizedNewVal) {
					retVal_7 = n
					didRet_6 = true
					return
				}
				n.ReplaceWith(newNode_2)
				this_.ReportChangeToEnclosingScope(newNode_2)
				retVal_7 = newNode_2
				didRet_6 = true
				return
			}()
			if didRet_6 {
				return retVal_7
			}
		}
	}
	var stringVal_2 string = this_.GetSideEffectFreeStringValue(firstArg)
	_ = stringVal_2
	if stringVal_2 == j2gort.NullString {
		return n
	}
	stringVal_2 = NodeUtil_TrimJsWhiteSpace(stringVal_2)
	if j2gort.CheckNotNullString(stringVal_2) == "" {
		return n
	}
	var isNegative bool = false
	_ = isNegative
	if javalang.StringStartsWith(j2gort.CheckNotNullString(stringVal_2), "-") {
		isNegative = true
		stringVal_2 = j2gort.JStrSubFrom(stringVal_2, int32(1))
	} else if javalang.StringStartsWith(j2gort.CheckNotNullString(stringVal_2), "+") {
		stringVal_2 = j2gort.JStrSubFrom(stringVal_2, int32(1))
	}
	if (radix == int32(0)) || (radix == int32(16)) {
		if (j2gort.JStrLen(stringVal_2) > int32(1)) && j2gort.Ext_com_google_common_base_Ascii_EqualsIgnoreCase(j2gort.BoxString(j2gort.JStrSub(stringVal_2, int32(0), int32(2))), "0x") {
			radix = int32(16)
			stringVal_2 = j2gort.JStrSubFrom(stringVal_2, int32(2))
		} else if radix == int32(0) {
			if !(this_.IsEcmaScript5OrGreater()) && javalang.StringStartsWith(j2gort.CheckNotNullString(stringVal_2), "0") {
				return n
			}
			radix = int32(10)
		}
	}
	var endDigitIndex int32 = int32(0)
	_ = endDigitIndex
	for endDigitIndex < j2gort.JStrLen(stringVal_2) {
		var c uint16 = j2gort.JStrCharAt(stringVal_2, endDigitIndex)
		_ = c
		if PeepholeReplaceKnownMethods_GetRadixDigit(c, radix) < int32(0) {
			break
		}
		endDigitIndex = endDigitIndex + 1
	}
	if endDigitIndex == int32(0) {
		return n
	}
	var digits string = j2gort.JStrSub(stringVal_2, int32(0), endDigitIndex)
	_ = digits
	var newVal_2 float64 = float64(0)
	_ = newVal_2
	var didRet_10 bool
	var retVal_11 pkg_rhino.Node
	func() {
		defer func() {
			if rec_12 := recover(); rec_12 != nil {
				jex_13 := javalang.RecoverException(rec_12)
				_ = jex_13
				if j2gort.InstanceOf[javalang.INumberFormatException](jex_13) {
					e_2 := j2gort.AsIface[javalang.INumberFormatException](jex_13)
					_ = e_2
					retVal_11 = n
					didRet_10 = true
					return
				} else {
					panic(rec_12)
				}
			}
		}()
		var bi *j2gort.Ext_java_math_BigInteger = j2gort.NewExt_java_math_BigInteger_2_string_int32(digits, radix)
		_ = bi
		newVal_2 = bi.DoubleValue()
	}()
	if didRet_10 {
		return retVal_11
	}
	if isNegative {
		newVal_2 = -(newVal_2)
	}
	var newNode_3 pkg_rhino.Node = NodeUtil_NumberNode(newVal_2, n)
	_ = newNode_3
	n.ReplaceWith(newNode_3)
	this_.ReportChangeToEnclosingScope(newNode_3)
	return newNode_3
}

// Java static method: PeepholeReplaceKnownMethods.getRadixDigit(char c, int radix)
func PeepholeReplaceKnownMethods_GetRadixDigit(c uint16, radix int32) int32 {
	clinit_PeepholeReplaceKnownMethods()
	_ = c
	_ = radix
	if (int32(c) >= int32(uint16(48))) && (int32(c) <= int32(uint16(57))) {
		var d int32 = (int32(c) - int32(uint16(48)))
		_ = d
		if d < radix {
			return d
		} else {
			return int32(-1)
		}
	}
	if (int32(c) >= int32(uint16(97))) && (int32(c) <= int32(uint16(122))) {
		var d_2 int32 = ((int32(c) - int32(uint16(97))) + int32(10))
		_ = d_2
		if d_2 < radix {
			return d_2
		} else {
			return int32(-1)
		}
	}
	if (int32(c) >= int32(uint16(65))) && (int32(c) <= int32(uint16(90))) {
		var d_3 int32 = ((int32(c) - int32(uint16(65))) + int32(10))
		_ = d_3
		if d_3 < radix {
			return d_3
		} else {
			return int32(-1)
		}
	}
	return int32(-1)
}

// Java method: PeepholeReplaceKnownMethods.tryFoldStringIndexOf(com.google.javascript.rhino.Node n, java.lang.String functionName, com.google.javascript.rhino.Node lstringNode, com.google.javascript.rhino.Node firstArg)
func (this_ *PeepholeReplaceKnownMethods) priv_PeepholeReplaceKnownMethods_tryFoldStringIndexOf_4264(n pkg_rhino.Node, functionName string, lstringNode pkg_rhino.Node, firstArg pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	_ = functionName
	_ = lstringNode
	_ = firstArg
	if !(n.IsCall()) {
		j2gort.CheckArgument(false)
	}
	if !(lstringNode.IsStringLit()) {
		j2gort.CheckArgument(false)
	}
	var lstring string = lstringNode.GetString()
	_ = lstring
	var isIndexOf bool = (j2gort.CheckNotNullString(functionName) == "indexOf")
	_ = isIndexOf
	var secondArg pkg_rhino.Node = firstArg.GetNext()
	_ = secondArg
	var searchValue string = this_.GetSideEffectFreeStringValue(firstArg)
	_ = searchValue
	if searchValue == j2gort.NullString {
		return n
	}
	var fromIndex int32
	_ = fromIndex
	if isIndexOf {
		fromIndex = int32(0)
	} else {
		fromIndex = j2gort.JStrLen(lstring)
	}
	if secondArg != 0 {
		if (secondArg.GetNext() != 0) || !(secondArg.IsNumber()) {
			return n
		} else {
			fromIndex = javalang.D2I(secondArg.GetDouble())
		}
	}
	var indexVal int32
	_ = indexVal
	if isIndexOf {
		indexVal = j2gort.JStrIndexOfFrom(lstring, searchValue, fromIndex)
	} else {
		indexVal = j2gort.JStrLastIndexOfFrom(lstring, searchValue, fromIndex)
	}
	var newNode pkg_rhino.Node = NodeUtil_NumberNode(float64(indexVal), n)
	_ = newNode
	n.ReplaceWith(newNode)
	this_.ReportChangeToEnclosingScope(newNode)
	return newNode
}

// Java method: PeepholeReplaceKnownMethods.tryFoldStringIncludes(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node lstringNode, com.google.javascript.rhino.Node firstArg)
func (this_ *PeepholeReplaceKnownMethods) priv_PeepholeReplaceKnownMethods_tryFoldStringIncludes_4265(n pkg_rhino.Node, lstringNode pkg_rhino.Node, firstArg pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	_ = lstringNode
	_ = firstArg
	if !(n.IsCall()) {
		j2gort.CheckArgument(false)
	}
	if !(lstringNode.IsStringLit()) {
		j2gort.CheckArgument(false)
	}
	if firstArg.IsRegExp() {
		return n
	}
	var lstring string = lstringNode.GetString()
	_ = lstring
	var searchValue string = this_.GetSideEffectFreeStringValue(firstArg)
	_ = searchValue
	if searchValue == j2gort.NullString {
		return n
	}
	var secondArg pkg_rhino.Node = firstArg.GetNext()
	_ = secondArg
	var fromIndex int32 = int32(0)
	_ = fromIndex
	if secondArg != 0 {
		if (secondArg.GetNext() != 0) || this_.MayHaveSideEffects(secondArg) {
			return n
		}
		if NodeUtil_IsNumericLiteral(secondArg) {
			var posVal float64 = NodeUtil_GetNumericLiteralValue(secondArg)
			_ = posVal
			fromIndex = javalang.D2I(math.Min(math.Max(posVal, float64(int32(0))), float64(j2gort.JStrLen(lstring))))
		} else if !(NodeUtil_IsNullOrUndefined(secondArg)) {
			return n
		}
	}
	var contains bool = (j2gort.JStrIndexOfFrom(lstring, searchValue, fromIndex) != int32(-1))
	_ = contains
	var newNode pkg_rhino.Node = NodeUtil_BooleanNode(contains).Srcref(n)
	_ = newNode
	n.ReplaceWith(newNode)
	this_.ReportChangeToEnclosingScope(newNode)
	return newNode
}

// Java method: PeepholeReplaceKnownMethods.tryFoldStringStartsWith(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node lstringNode, com.google.javascript.rhino.Node firstArg)
func (this_ *PeepholeReplaceKnownMethods) priv_PeepholeReplaceKnownMethods_tryFoldStringStartsWith_4266(n pkg_rhino.Node, lstringNode pkg_rhino.Node, firstArg pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	_ = lstringNode
	_ = firstArg
	if !(n.IsCall()) {
		j2gort.CheckArgument(false)
	}
	if !(lstringNode.IsStringLit()) {
		j2gort.CheckArgument(false)
	}
	if firstArg.IsRegExp() {
		return n
	}
	var lstring string = lstringNode.GetString()
	_ = lstring
	var searchValue string = this_.GetSideEffectFreeStringValue(firstArg)
	_ = searchValue
	if searchValue == j2gort.NullString {
		return n
	}
	var secondArg pkg_rhino.Node = firstArg.GetNext()
	_ = secondArg
	var start int32 = int32(0)
	_ = start
	if secondArg != 0 {
		if (secondArg.GetNext() != 0) || this_.MayHaveSideEffects(secondArg) {
			return n
		}
		if NodeUtil_IsNumericLiteral(secondArg) {
			var posVal float64 = NodeUtil_GetNumericLiteralValue(secondArg)
			_ = posVal
			start = javalang.D2I(math.Min(math.Max(posVal, float64(int32(0))), float64(j2gort.JStrLen(lstring))))
		} else if !(NodeUtil_IsNullOrUndefined(secondArg)) {
			return n
		}
	}
	var result bool = j2gort.JStrStartsWithOffset(lstring, searchValue, start)
	_ = result
	var newNode pkg_rhino.Node = NodeUtil_BooleanNode(result).Srcref(n)
	_ = newNode
	n.ReplaceWith(newNode)
	this_.ReportChangeToEnclosingScope(newNode)
	return newNode
}

// Java method: PeepholeReplaceKnownMethods.tryFoldStringEndsWith(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node lstringNode, com.google.javascript.rhino.Node firstArg)
func (this_ *PeepholeReplaceKnownMethods) priv_PeepholeReplaceKnownMethods_tryFoldStringEndsWith_4267(n pkg_rhino.Node, lstringNode pkg_rhino.Node, firstArg pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	_ = lstringNode
	_ = firstArg
	if !(n.IsCall()) {
		j2gort.CheckArgument(false)
	}
	if !(lstringNode.IsStringLit()) {
		j2gort.CheckArgument(false)
	}
	if firstArg.IsRegExp() {
		return n
	}
	var lstring string = lstringNode.GetString()
	_ = lstring
	var searchValue string = this_.GetSideEffectFreeStringValue(firstArg)
	_ = searchValue
	if searchValue == j2gort.NullString {
		return n
	}
	var secondArg pkg_rhino.Node = firstArg.GetNext()
	_ = secondArg
	var end int32 = j2gort.JStrLen(lstring)
	_ = end
	if secondArg != 0 {
		if (secondArg.GetNext() != 0) || this_.MayHaveSideEffects(secondArg) {
			return n
		}
		if NodeUtil_IsNumericLiteral(secondArg) {
			var posVal float64 = NodeUtil_GetNumericLiteralValue(secondArg)
			_ = posVal
			end = javalang.D2I(math.Min(math.Max(posVal, float64(int32(0))), float64(j2gort.JStrLen(lstring))))
		} else if secondArg.IsNull() {
			end = int32(0)
		} else if NodeUtil_IsUndefined(secondArg) {
			end = j2gort.JStrLen(lstring)
		} else {
			return n
		}
	}
	var searchLength int32 = j2gort.JStrLen(searchValue)
	_ = searchLength
	var start int32 = (end - searchLength)
	_ = start
	var result bool = ((start >= int32(0)) && j2gort.JStrStartsWithOffset(lstring, searchValue, start))
	_ = result
	var newNode pkg_rhino.Node = NodeUtil_BooleanNode(result).Srcref(n)
	_ = newNode
	n.ReplaceWith(newNode)
	this_.ReportChangeToEnclosingScope(newNode)
	return newNode
}

// Java method: PeepholeReplaceKnownMethods.tryFoldStringAt(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node stringNode, com.google.javascript.rhino.Node firstArg)
func (this_ *PeepholeReplaceKnownMethods) priv_PeepholeReplaceKnownMethods_tryFoldStringAt_4268(n pkg_rhino.Node, stringNode pkg_rhino.Node, firstArg pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	_ = stringNode
	_ = firstArg
	if !(n.IsCall()) {
		j2gort.CheckArgument(false)
	}
	if !(stringNode.IsStringLit()) {
		j2gort.CheckArgument(false)
	}
	var posVal float64 = float64(0)
	_ = posVal
	if firstArg == 0 {
		posVal = float64(0.0)
	} else {
		if firstArg.GetNext() != 0 {
			return n
		}
		var maybePos *j2gort.BoxDouble = this_.GetSideEffectFreeNumberValue(firstArg)
		_ = maybePos
		if maybePos == nil {
			return n
		}
		if maybePos.IsNaN() {
			posVal = float64(0.0)
		} else {
			posVal = j2gort.UnboxDouble(maybePos)
		}
	}
	var stringVal string = stringNode.GetString()
	_ = stringVal
	var len_ int32 = j2gort.JStrLen(stringVal)
	_ = len_
	var replacement pkg_rhino.Node = pkg_rhino.Node(0)
	_ = replacement
	if (math.IsInf(posVal, 0) || (posVal >= float64(len_))) || (posVal <= float64(-(len_ + int32(1)))) {
		replacement = NodeUtil_NewUndefinedNode(n)
	} else {
		var intPos int64 = javalang.D2L(posVal)
		_ = intPos
		var k int64
		_ = k
		if intPos >= int64(int32(0)) {
			k = intPos
		} else {
			k = (int64(len_) + intPos)
		}
		if (k < int64(int32(0))) || (k >= int64(len_)) {
			replacement = NodeUtil_NewUndefinedNode(n)
		} else {
			replacement = pkg_rhino.IR_String(javalang.StringFromChar(j2gort.JStrCharAt(stringVal, int32(k)))).Srcref(n)
		}
	}
	n.ReplaceWith(replacement)
	this_.ReportChangeToEnclosingScope(replacement)
	return replacement
}

// Java method: PeepholeReplaceKnownMethods.tryFoldArrayJoin(com.google.javascript.rhino.Node n)
func (this_ *PeepholeReplaceKnownMethods) priv_PeepholeReplaceKnownMethods_tryFoldArrayJoin_4269(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if !(n.IsCall()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(n))
	}
	var callTarget pkg_rhino.Node = n.GetFirstChild()
	_ = callTarget
	if (callTarget == 0) || !(callTarget.IsGetProp()) {
		return n
	}
	var right pkg_rhino.Node = callTarget.GetNext()
	_ = right
	if right != 0 {
		if (right.GetNext() != 0) || !(NodeUtil_IsImmutableValue(right)) {
			return n
		}
	}
	var arrayNode pkg_rhino.Node = callTarget.GetFirstChild()
	_ = arrayNode
	if !(arrayNode.IsArrayLit()) || !(j2gort.CheckNotNullString(callTarget.GetString()) == "join") {
		return n
	}
	if (right != 0) && (NodeUtil_IsUndefined(right) || (right.IsStringLit() && func() bool {
		sr_14 := ","
		se_15 := right.GetString()
		return j2gort.CheckNotNullString(sr_14) == se_15
	}())) {
		right.Detach()
		this_.ReportChangeToEnclosingScope(n)
	}
	var joinString string
	_ = joinString
	if (right == 0) || NodeUtil_IsUndefined(right) {
		joinString = ","
	} else {
		joinString = NodeUtil_GetStringValue(right)
	}
	var arrayFoldedChildren []pkg_rhino.Node = []pkg_rhino.Node{}
	_ = arrayFoldedChildren
	var sb *j2gort.Ext_java_lang_StringBuilder = nil
	_ = sb
	var foldedSize int32 = int32(0)
	_ = foldedSize
	var prev pkg_rhino.Node = pkg_rhino.Node(0)
	_ = prev
	var elem pkg_rhino.Node = arrayNode.GetFirstChild()
	_ = elem
	for elem != 0 {
		if NodeUtil_IsImmutableValue(elem) || elem.IsEmpty() {
			if sb == nil {
				sb = j2gort.NewExt_java_lang_StringBuilder()
			} else {
				sb.Append_1_string(joinString)
			}
			var elementStr string = NodeUtil_GetArrayElementStringValue(elem)
			_ = elementStr
			if elementStr == j2gort.NullString {
				return n
			}
			sb.Append_1_string(elementStr)
		} else {
			if sb != nil {
				j2gort.CheckNotNullNode(prev)
				foldedSize = (foldedSize + (sb.Length() + int32(2)))
				arrayFoldedChildren = append(arrayFoldedChildren, pkg_rhino.IR_String(sb.String()).SrcrefIfMissing(prev))
				sb = nil
			}
			foldedSize = (foldedSize + InlineCostEstimator_GetCost_1(elem))
			arrayFoldedChildren = append(arrayFoldedChildren, elem)
		}
		prev = elem
		elem = elem.GetNext()
	}
	if sb != nil {
		j2gort.CheckNotNullNode(prev)
		foldedSize = (foldedSize + (sb.Length() + int32(2)))
		arrayFoldedChildren = append(arrayFoldedChildren, pkg_rhino.IR_String(sb.String()).SrcrefIfMissing(prev))
	}
	foldedSize = (foldedSize + (int32(len(arrayFoldedChildren)) - int32(1)))
	var originalSize int32 = InlineCostEstimator_GetCost_1(n)
	_ = originalSize
	switch int32(len(arrayFoldedChildren)) {
	case int32(0):
		var emptyStringNode pkg_rhino.Node = pkg_rhino.IR_String("")
		_ = emptyStringNode
		n.ReplaceWith(emptyStringNode)
		this_.ReportChangeToEnclosingScope(emptyStringNode)
		return emptyStringNode
	case int32(1):
		var foldedStringNode pkg_rhino.Node = j2gort.SliceRemoveAtElem(&arrayFoldedChildren, int32(0))
		_ = foldedStringNode
		if foldedStringNode.IsSpread() || (foldedSize > originalSize) {
			return n
		}
		if foldedStringNode.IsStringLit() {
			arrayNode.DetachChildren()
			n.ReplaceWith(foldedStringNode)
			this_.ReportChangeToEnclosingScope(foldedStringNode)
			return foldedStringNode
		} else {
			return n
		}
	default:
		if arrayNode.HasXChildren(int32(len(arrayFoldedChildren))) {
			return n
		}
		var kJoinOverhead int32 = j2gort.JStrLen("[].join()")
		_ = kJoinOverhead
		foldedSize = (foldedSize + kJoinOverhead)
		foldedSize = (foldedSize + func() int32 {
			if right != 0 {
				return InlineCostEstimator_GetCost_1(right)
			}
			return int32(0)
		}())
		if foldedSize > originalSize {
			return n
		}
		arrayNode.DetachChildren()
		{
			for _, node := range arrayFoldedChildren {
				_ = node
				arrayNode.AddChildToBack(node)
			}
		}
		this_.ReportChangeToEnclosingScope(arrayNode)
	}
	return n
}

// Java method: PeepholeReplaceKnownMethods.tryFoldStringSubstr(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node stringNode, com.google.javascript.rhino.Node arg1)
func (this_ *PeepholeReplaceKnownMethods) priv_PeepholeReplaceKnownMethods_tryFoldStringSubstr_4270(n pkg_rhino.Node, stringNode pkg_rhino.Node, arg1 pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	_ = stringNode
	_ = arg1
	if !(n.IsCall()) {
		j2gort.CheckArgument(false)
	}
	if !(stringNode.IsStringLit()) {
		j2gort.CheckArgument(false)
	}
	if !(arg1 != 0) {
		j2gort.CheckArgument(false)
	}
	var start int32 = int32(0)
	_ = start
	var length int32 = int32(0)
	_ = length
	var stringAsString string = stringNode.GetString()
	_ = stringAsString
	var maybeStart *j2gort.BoxDouble = this_.GetSideEffectFreeNumberValue(arg1)
	_ = maybeStart
	if maybeStart != nil {
		start = maybeStart.IntValue()
	} else {
		return n
	}
	var arg2 pkg_rhino.Node = arg1.GetNext()
	_ = arg2
	if arg2 != 0 {
		var maybeLength *j2gort.BoxDouble = this_.GetSideEffectFreeNumberValue(arg2)
		_ = maybeLength
		if maybeLength != nil {
			length = maybeLength.IntValue()
		} else {
			return n
		}
		if arg2.GetNext() != 0 {
			return n
		}
	} else {
		length = (j2gort.JStrLen(stringAsString) - start)
	}
	if (((start + length) > j2gort.JStrLen(stringAsString)) || (length < int32(0))) || (start < int32(0)) {
		return n
	}
	var result string = j2gort.JStrSub(stringAsString, start, (start + length))
	_ = result
	var resultNode pkg_rhino.Node = pkg_rhino.IR_String(result)
	_ = resultNode
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	n.ReplaceWith(resultNode)
	this_.ReportChangeToEnclosingScope(parent)
	return resultNode
}

// Java method: PeepholeReplaceKnownMethods.tryFoldStringSubstringOrSlice(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node stringNode, com.google.javascript.rhino.Node arg1)
func (this_ *PeepholeReplaceKnownMethods) priv_PeepholeReplaceKnownMethods_tryFoldStringSubstringOrSlice_4271(n pkg_rhino.Node, stringNode pkg_rhino.Node, arg1 pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	_ = stringNode
	_ = arg1
	if !(n.IsCall()) {
		j2gort.CheckArgument(false)
	}
	if !(stringNode.IsStringLit()) {
		j2gort.CheckArgument(false)
	}
	if !(arg1 != 0) {
		j2gort.CheckArgument(false)
	}
	var start int32 = int32(0)
	_ = start
	var end int32 = int32(0)
	_ = end
	var stringAsString string = stringNode.GetString()
	_ = stringAsString
	var maybeStart *j2gort.BoxDouble = this_.GetSideEffectFreeNumberValue(arg1)
	_ = maybeStart
	if maybeStart != nil {
		start = maybeStart.IntValue()
	} else {
		return n
	}
	var arg2 pkg_rhino.Node = arg1.GetNext()
	_ = arg2
	if arg2 != 0 {
		var maybeEnd *j2gort.BoxDouble = this_.GetSideEffectFreeNumberValue(arg2)
		_ = maybeEnd
		if maybeEnd != nil {
			end = maybeEnd.IntValue()
		} else {
			return n
		}
		if arg2.GetNext() != 0 {
			return n
		}
	} else {
		end = j2gort.JStrLen(stringAsString)
	}
	if ((((end > j2gort.JStrLen(stringAsString)) || (start > j2gort.JStrLen(stringAsString))) || (start < int32(0))) || (end < int32(0))) || (start > end) {
		return n
	}
	var result string = j2gort.JStrSub(stringAsString, start, end)
	_ = result
	var resultNode pkg_rhino.Node = pkg_rhino.IR_String(result)
	_ = resultNode
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	n.ReplaceWith(resultNode)
	this_.ReportChangeToEnclosingScope(parent)
	return resultNode
}

// Java method: PeepholeReplaceKnownMethods.replaceWithCharAt(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node callTarget, com.google.javascript.rhino.Node firstArg)
func (this_ *PeepholeReplaceKnownMethods) priv_PeepholeReplaceKnownMethods_replaceWithCharAt_4272(n pkg_rhino.Node, callTarget pkg_rhino.Node, firstArg pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	_ = callTarget
	_ = firstArg
	callTarget.SetString("charAt")
	firstArg.GetNext().Detach()
	this_.ReportChangeToEnclosingScope(firstArg)
	return n
}

// Java method: PeepholeReplaceKnownMethods.tryFoldStringCharAt(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node stringNode, com.google.javascript.rhino.Node arg1)
func (this_ *PeepholeReplaceKnownMethods) priv_PeepholeReplaceKnownMethods_tryFoldStringCharAt_4273(n pkg_rhino.Node, stringNode pkg_rhino.Node, arg1 pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	_ = stringNode
	_ = arg1
	if !(n.IsCall()) {
		j2gort.CheckArgument(false)
	}
	if !(stringNode.IsStringLit()) {
		j2gort.CheckArgument(false)
	}
	var index int32 = int32(0)
	_ = index
	var stringAsString string = stringNode.GetString()
	_ = stringAsString
	if ((arg1 != 0) && arg1.IsNumber()) && (arg1.GetNext() == 0) {
		index = javalang.D2I(arg1.GetDouble())
	} else {
		return n
	}
	if (index < int32(0)) || (j2gort.JStrLen(stringAsString) <= index) {
		return n
	}
	var resultNode pkg_rhino.Node = pkg_rhino.IR_String(j2gort.JStrSub(stringAsString, index, (index + int32(1))))
	_ = resultNode
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	n.ReplaceWith(resultNode)
	this_.ReportChangeToEnclosingScope(parent)
	return resultNode
}

// Java method: PeepholeReplaceKnownMethods.tryFoldStringReplace(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node stringNode, com.google.javascript.rhino.Node arg1)
func (this_ *PeepholeReplaceKnownMethods) priv_PeepholeReplaceKnownMethods_tryFoldStringReplace_4274(n pkg_rhino.Node, stringNode pkg_rhino.Node, arg1 pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	_ = stringNode
	_ = arg1
	if !(n.IsCall()) {
		j2gort.CheckArgument(false)
	}
	if !(stringNode.IsStringLit()) {
		j2gort.CheckArgument(false)
	}
	var arg2 pkg_rhino.Node = arg1.GetNext()
	_ = arg2
	if (arg2 == 0) || (arg2.GetNext() != 0) {
		return n
	}
	if !(arg1.IsStringLit()) || !(arg2.IsStringLit()) {
		return n
	}
	var lookForPattern string = arg1.GetString()
	_ = lookForPattern
	var replacementPattern string = arg2.GetString()
	_ = replacementPattern
	if j2gort.JStrIndexOf(replacementPattern, "$") >= 0 {
		return n
	}
	var original string = stringNode.GetString()
	_ = original
	var index int32 = j2gort.JStrIndexOf(original, lookForPattern)
	_ = index
	if index == int32(-1) {
		return n
	}
	var newString string = (j2gort.StrValString(j2gort.JStrSub(original, int32(0), index)) + j2gort.StrValString(replacementPattern) + j2gort.StrValString(j2gort.JStrSubFrom(original, (index+j2gort.JStrLen(lookForPattern)))))
	_ = newString
	var resultNode pkg_rhino.Node = pkg_rhino.IR_String(newString).Srcref(stringNode)
	_ = resultNode
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	n.ReplaceWith(resultNode)
	this_.ReportChangeToEnclosingScope(parent)
	return resultNode
}

// Java method: PeepholeReplaceKnownMethods.tryFoldStringReplaceAll(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node stringNode, com.google.javascript.rhino.Node arg1)
func (this_ *PeepholeReplaceKnownMethods) priv_PeepholeReplaceKnownMethods_tryFoldStringReplaceAll_4275(n pkg_rhino.Node, stringNode pkg_rhino.Node, arg1 pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	_ = stringNode
	_ = arg1
	if !(n.IsCall()) {
		j2gort.CheckArgument(false)
	}
	if !(stringNode.IsStringLit()) {
		j2gort.CheckArgument(false)
	}
	var arg2 pkg_rhino.Node = arg1.GetNext()
	_ = arg2
	if (arg2 == 0) || (arg2.GetNext() != 0) {
		return n
	}
	if !(arg1.IsStringLit()) || !(arg2.IsStringLit()) {
		return n
	}
	var replacementPattern string = arg2.GetString()
	_ = replacementPattern
	if j2gort.JStrIndexOf(replacementPattern, "$") >= 0 {
		return n
	}
	var original string = stringNode.GetString()
	_ = original
	var newString string = javalang.StringReplace(original, j2gort.CheckNotNullString(arg1.GetString()), j2gort.CheckNotNullString(replacementPattern))
	_ = newString
	var resultNode pkg_rhino.Node = pkg_rhino.IR_String(newString).Srcref(stringNode)
	_ = resultNode
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	n.ReplaceWith(resultNode)
	this_.ReportChangeToEnclosingScope(parent)
	return resultNode
}

// Java method: PeepholeReplaceKnownMethods.tryFoldStringCharCodeAt(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node stringNode, com.google.javascript.rhino.Node arg1)
func (this_ *PeepholeReplaceKnownMethods) priv_PeepholeReplaceKnownMethods_tryFoldStringCharCodeAt_4276(n pkg_rhino.Node, stringNode pkg_rhino.Node, arg1 pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	_ = stringNode
	_ = arg1
	if !(n.IsCall()) {
		j2gort.CheckArgument(false)
	}
	if !(stringNode.IsStringLit()) {
		j2gort.CheckArgument(false)
	}
	var index int32 = int32(0)
	_ = index
	var stringAsString string = stringNode.GetString()
	_ = stringAsString
	if ((arg1 != 0) && arg1.IsNumber()) && (arg1.GetNext() == 0) {
		index = javalang.D2I(arg1.GetDouble())
	} else {
		return n
	}
	if (index < int32(0)) || (j2gort.JStrLen(stringAsString) <= index) {
		return n
	}
	var resultNode pkg_rhino.Node = pkg_rhino.IR_Number(float64(j2gort.JStrCharAt(stringAsString, index)))
	_ = resultNode
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	n.ReplaceWith(resultNode)
	this_.ReportChangeToEnclosingScope(parent)
	return resultNode
}

// Java static method: PeepholeReplaceKnownMethods.jsSplitMatch(java.lang.String stringValue, int startIndex, java.lang.String separator)
func PeepholeReplaceKnownMethods_JsSplitMatch(stringValue string, startIndex int32, separator string) int32 {
	clinit_PeepholeReplaceKnownMethods()
	_ = stringValue
	_ = startIndex
	_ = separator
	if (startIndex + j2gort.JStrLen(separator)) > j2gort.JStrLen(stringValue) {
		return int32(-1)
	}
	var matchIndex int32 = j2gort.JStrIndexOfFrom(stringValue, separator, startIndex)
	_ = matchIndex
	if matchIndex < int32(0) {
		return int32(-1)
	}
	return matchIndex
}

// Java method: PeepholeReplaceKnownMethods.jsSplit(java.lang.String stringValue, java.lang.String separator, int limit)
func (this_ *PeepholeReplaceKnownMethods) priv_PeepholeReplaceKnownMethods_jsSplit_4277(stringValue string, separator string, limit int32) []string {
	_ = stringValue
	_ = separator
	_ = limit
	if !(limit >= int32(0)) {
		j2gort.CheckArgument(false)
	}
	if !(stringValue != j2gort.NullString) {
		j2gort.CheckArgument(false)
	}
	if limit == int32(0) {
		return j2gort.NewStringArray(int32(0))
	}
	if separator == j2gort.NullString {
		return []string{stringValue}
	}
	var splitStrings j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList()
	_ = splitStrings
	if j2gort.CheckNotNullString(separator) == "" {
		{
			var i int32 = int32(0)
			_ = i
			for ; (i < j2gort.JStrLen(stringValue)) && (i < limit); i = i + 1 {
				splitStrings.Add_1(j2gort.BoxString(j2gort.JStrSub(stringValue, i, (i + int32(1)))))
			}
		}
	} else {
		var startIndex int32 = int32(0)
		_ = startIndex
		var matchIndex int32 = int32(0)
		_ = matchIndex
		for ((func() int32 {
			var asgn_17 int32 = PeepholeReplaceKnownMethods_JsSplitMatch(stringValue, startIndex, separator)
			matchIndex = asgn_17
			return asgn_17
		}()) >= int32(0)) && (splitStrings.Size() < limit) {
			splitStrings.Add_1(j2gort.BoxString(j2gort.JStrSub(stringValue, startIndex, matchIndex)))
			startIndex = (matchIndex + j2gort.JStrLen(separator))
		}
		if splitStrings.Size() < limit {
			if startIndex < j2gort.JStrLen(stringValue) {
				splitStrings.Add_1(j2gort.BoxString(j2gort.JStrSubFrom(stringValue, startIndex)))
			} else {
				splitStrings.Add_1("")
			}
		}
	}
	return j2gort.ColToNewStringArray(splitStrings)
}

// Java method: PeepholeReplaceKnownMethods.tryFoldStringSplit(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node stringNode, com.google.javascript.rhino.Node arg1)
func (this_ *PeepholeReplaceKnownMethods) priv_PeepholeReplaceKnownMethods_tryFoldStringSplit_4278(n pkg_rhino.Node, stringNode pkg_rhino.Node, arg1 pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	_ = stringNode
	_ = arg1
	if this_.late {
		return n
	}
	if !(n.IsCall()) {
		j2gort.CheckArgument(false)
	}
	if !(stringNode.IsStringLit()) {
		j2gort.CheckArgument(false)
	}
	var separator string = j2gort.NullString
	_ = separator
	var stringValue string = stringNode.GetString()
	_ = stringValue
	var limit int32 = (j2gort.JStrLen(stringValue) + int32(1))
	_ = limit
	if arg1 != 0 {
		if arg1.IsStringLit() {
			separator = arg1.GetString()
		} else if !(arg1.IsNull()) {
			return n
		}
		var arg2 pkg_rhino.Node = arg1.GetNext()
		_ = arg2
		if arg2 != 0 {
			if NodeUtil_IsNumericLiteral(arg2) {
				var uint32Limit int64 = j2gort.BoxInt_ToUnsignedLong(pkg_base.JSCompDoubles_EcmascriptToUint32(NodeUtil_GetNumericLiteralValue(arg2)))
				_ = uint32Limit
				limit = int32(min(uint32Limit, int64(limit)))
			} else {
				return n
			}
		}
	}
	var stringArray []string = this_.priv_PeepholeReplaceKnownMethods_jsSplit_4277(stringValue, separator, limit)
	_ = stringArray
	var arrayOfStrings pkg_rhino.Node = pkg_rhino.IR_Arraylit_1_Node([]pkg_rhino.Node{})
	_ = arrayOfStrings
	for _, el_18 := range j2gort.CheckNotNullSlice(stringArray) {
		element := el_18
		_ = element
		arrayOfStrings.AddChildToBack(pkg_rhino.IR_String(element).Srcref(stringNode))
	}
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	n.ReplaceWith(arrayOfStrings)
	this_.ReportChangeToEnclosingScope(parent)
	return arrayOfStrings
}

// Java method: PeepholeReplaceKnownMethods.tryToFoldArrayConcat(com.google.javascript.rhino.Node n)
func (this_ *PeepholeReplaceKnownMethods) priv_PeepholeReplaceKnownMethods_tryToFoldArrayConcat_4279(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if !(n.IsCall()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(n))
	}
	if !(this_.IsASTNormalized()) || !(this_.useTypes) {
		return n
	}
	var concatFunctionCall PeepholeReplaceKnownMethods_ConcatFunctionCall = PeepholeReplaceKnownMethods_CreateConcatFunctionCallForNode(n)
	_ = concatFunctionCall
	if concatFunctionCall == nil || concatFunctionCall.IsJavaNil_() {
		return n
	}
	concatFunctionCall = this_.priv_PeepholeReplaceKnownMethods_tryToRemoveArrayLiteralFromFrontOfConcat_4280(concatFunctionCall)
	j2gort.CheckNotNullIface(concatFunctionCall)
	return this_.priv_PeepholeReplaceKnownMethods_tryToFoldConcatChaining_4281(concatFunctionCall)
}

// Java method: PeepholeReplaceKnownMethods.tryToRemoveArrayLiteralFromFrontOfConcat(com.google.javascript.jscomp.PeepholeReplaceKnownMethods.ConcatFunctionCall concatFunctionCall)
func (this_ *PeepholeReplaceKnownMethods) priv_PeepholeReplaceKnownMethods_tryToRemoveArrayLiteralFromFrontOfConcat_4280(concatFunctionCall PeepholeReplaceKnownMethods_ConcatFunctionCall) PeepholeReplaceKnownMethods_ConcatFunctionCall {
	_ = concatFunctionCall
	j2gort.CheckNotNullIface(concatFunctionCall)
	var callNode pkg_rhino.Node = concatFunctionCall.GetPeepholeReplaceKnownMethods_ConcatFunctionCallBase().callNode
	_ = callNode
	var arrayLiteralToRemove pkg_rhino.Node = concatFunctionCall.GetPeepholeReplaceKnownMethods_ConcatFunctionCallBase().calleeNode
	_ = arrayLiteralToRemove
	if !(arrayLiteralToRemove.IsArrayLit()) || arrayLiteralToRemove.HasChildren() {
		return concatFunctionCall
	}
	var firstArg pkg_rhino.Node = concatFunctionCall.GetPeepholeReplaceKnownMethods_ConcatFunctionCallBase().firstArgumentNode
	_ = firstArg
	if !(PeepholeReplaceKnownMethods_ContainsExactlyArray(firstArg)) {
		return concatFunctionCall
	}
	firstArg.Detach()
	arrayLiteralToRemove.ReplaceWith(firstArg)
	this_.ReportChangeToEnclosingScope(callNode)
	return PeepholeReplaceKnownMethods_CreateConcatFunctionCallForNode(callNode)
}

// Java method: PeepholeReplaceKnownMethods.tryToFoldConcatChaining(com.google.javascript.jscomp.PeepholeReplaceKnownMethods.ConcatFunctionCall concatFunctionCall)
func (this_ *PeepholeReplaceKnownMethods) priv_PeepholeReplaceKnownMethods_tryToFoldConcatChaining_4281(concatFunctionCall PeepholeReplaceKnownMethods_ConcatFunctionCall) pkg_rhino.Node {
	_ = concatFunctionCall
	j2gort.CheckNotNullIface(concatFunctionCall)
	var concatCallNode pkg_rhino.Node = concatFunctionCall.GetPeepholeReplaceKnownMethods_ConcatFunctionCallBase().callNode
	_ = concatCallNode
	var maybeFunctionCall pkg_rhino.Node = concatFunctionCall.GetPeepholeReplaceKnownMethods_ConcatFunctionCallBase().calleeNode
	_ = maybeFunctionCall
	if !(maybeFunctionCall.IsCall()) {
		return concatCallNode
	}
	var previousConcatFunctionCall PeepholeReplaceKnownMethods_ConcatFunctionCall = PeepholeReplaceKnownMethods_CreateConcatFunctionCallForNode(maybeFunctionCall)
	_ = previousConcatFunctionCall
	if previousConcatFunctionCall == nil || previousConcatFunctionCall.IsJavaNil_() {
		return concatCallNode
	}
	{
		var arg pkg_rhino.Node = concatFunctionCall.GetPeepholeReplaceKnownMethods_ConcatFunctionCallBase().firstArgumentNode
		_ = arg
		for ; arg != 0; arg = arg.GetNext() {
			if this_.MayHaveSideEffects(arg) {
				return concatCallNode
			}
		}
	}
	var previousConcatCallNode pkg_rhino.Node = previousConcatFunctionCall.GetPeepholeReplaceKnownMethods_ConcatFunctionCallBase().callNode
	_ = previousConcatCallNode
	var arg_2 pkg_rhino.Node = concatFunctionCall.GetPeepholeReplaceKnownMethods_ConcatFunctionCallBase().firstArgumentNode
	_ = arg_2
	for arg_2 != 0 {
		var currentArg pkg_rhino.Node = arg_2
		_ = currentArg
		arg_2 = arg_2.GetNext()
		previousConcatCallNode.AddChildToBack(currentArg.Detach())
	}
	concatCallNode.ReplaceWith(previousConcatCallNode.Detach())
	this_.ReportChangeToEnclosingScope(previousConcatCallNode)
	return previousConcatCallNode
}

// Java static method: PeepholeReplaceKnownMethods.createConcatFunctionCallForNode(com.google.javascript.rhino.Node n)
func PeepholeReplaceKnownMethods_CreateConcatFunctionCallForNode(n pkg_rhino.Node) PeepholeReplaceKnownMethods_ConcatFunctionCall {
	clinit_PeepholeReplaceKnownMethods()
	_ = n
	if !(n.IsCall()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(n))
	}
	var callTarget pkg_rhino.Node = j2gort.CheckNotNullNode(n.GetFirstChild())
	_ = callTarget
	if !(callTarget.IsGetProp()) || !(j2gort.CheckNotNullString(callTarget.GetString()) == "concat") {
		return nil
	}
	var calleeNode pkg_rhino.Node = callTarget.GetFirstChild()
	_ = calleeNode
	if !(PeepholeReplaceKnownMethods_ContainsExactlyArray(calleeNode)) {
		return nil
	}
	var firstArgumentNode pkg_rhino.Node = n.GetSecondChild()
	_ = firstArgumentNode
	return NewPeepholeReplaceKnownMethods__Anon_69(n, calleeNode, firstArgumentNode)
}

// Java static method: PeepholeReplaceKnownMethods.containsExactlyArray(com.google.javascript.rhino.Node n)
func PeepholeReplaceKnownMethods_ContainsExactlyArray(n pkg_rhino.Node) bool {
	clinit_PeepholeReplaceKnownMethods()
	_ = n
	if n == 0 {
		return false
	}
	if n.IsArrayLit() {
		return true
	}
	if !(n.IsCall()) {
		return false
	}
	var callee pkg_rhino.Node = n.GetFirstChild()
	_ = callee
	return ((callee.IsGetProp() && (j2gort.CheckNotNullString(callee.GetString()) == "concat")) && PeepholeReplaceKnownMethods_ContainsExactlyArray(callee.GetFirstChild()))
}

// Java polymorphic class interface: com.google.javascript.jscomp.PeepholeReplaceKnownMethods.ConcatFunctionCall
type PeepholeReplaceKnownMethods_ConcatFunctionCall interface {
	Implements_jscomp_PeepholeReplaceKnownMethods_ConcatFunctionCall()
	IsJavaNil_() bool
	GetPeepholeReplaceKnownMethods_ConcatFunctionCallBase() *PeepholeReplaceKnownMethods_ConcatFunctionCallBase
}

func ToIface_PeepholeReplaceKnownMethods_ConcatFunctionCall[T interface {
	PeepholeReplaceKnownMethods_ConcatFunctionCall
	comparable
}](v T) PeepholeReplaceKnownMethods_ConcatFunctionCall {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_PeepholeReplaceKnownMethods_ConcatFunctionCall[T PeepholeReplaceKnownMethods_ConcatFunctionCall](v T) PeepholeReplaceKnownMethods_ConcatFunctionCall {
	var vi PeepholeReplaceKnownMethods_ConcatFunctionCall = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.PeepholeReplaceKnownMethods.ConcatFunctionCall
type PeepholeReplaceKnownMethods_ConcatFunctionCallBase struct {
	callNode          pkg_rhino.Node `gson:"callNode"`
	calleeNode        pkg_rhino.Node `gson:"calleeNode"`
	firstArgumentNode pkg_rhino.Node `gson:"firstArgumentNode"`
	idHash_           atomic.Int32
}

func (this_ *PeepholeReplaceKnownMethods_ConcatFunctionCallBase) BindSelf_(s any) {
	_ = s
}
func (this_ *PeepholeReplaceKnownMethods_ConcatFunctionCallBase) IsJavaNil_() bool {
	return this_ == nil
}
func (*PeepholeReplaceKnownMethods_ConcatFunctionCallBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PeepholeReplaceKnownMethods$ConcatFunctionCall", "ConcatFunctionCall"
}
func (this_ *PeepholeReplaceKnownMethods_ConcatFunctionCallBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PeepholeReplaceKnownMethods.ConcatFunctionCall", "com.google.javascript.jscomp.PeepholeReplaceKnownMethods$ConcatFunctionCall", func(v_ any) bool {
		return j2gort.InstanceOf[PeepholeReplaceKnownMethods_ConcatFunctionCall](v_)
	})
}
func (b *PeepholeReplaceKnownMethods_ConcatFunctionCallBase) GetPeepholeReplaceKnownMethods_ConcatFunctionCallBase() *PeepholeReplaceKnownMethods_ConcatFunctionCallBase {
	return b
}
func (*PeepholeReplaceKnownMethods_ConcatFunctionCallBase) Implements_jscomp_PeepholeReplaceKnownMethods_ConcatFunctionCall() {
}

func (this_ *PeepholeReplaceKnownMethods_ConcatFunctionCallBase) initFields_PeepholeReplaceKnownMethods_ConcatFunctionCall() {
	_ = this_
}

func InitPeepholeReplaceKnownMethods_ConcatFunctionCall(this_ *PeepholeReplaceKnownMethods_ConcatFunctionCallBase, callNode pkg_rhino.Node, calleeNode pkg_rhino.Node, firstArgumentNode pkg_rhino.Node) {
	_ = callNode
	_ = calleeNode
	_ = firstArgumentNode
	this_.initFields_PeepholeReplaceKnownMethods_ConcatFunctionCall()
	this_.callNode = j2gort.CheckNotNullNode(callNode)
	this_.calleeNode = j2gort.CheckNotNullNode(calleeNode)
	this_.firstArgumentNode = firstArgumentNode
}

// Java class:
type PeepholeReplaceKnownMethods__Anon_69 struct {
	PeepholeReplaceKnownMethods_ConcatFunctionCallBase
}

func (this_ *PeepholeReplaceKnownMethods__Anon_69) BindSelf_(s any) {
	_ = s
}
func (this_ *PeepholeReplaceKnownMethods__Anon_69) IsJavaNil_() bool { return this_ == nil }
func (*PeepholeReplaceKnownMethods__Anon_69) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PeepholeReplaceKnownMethods$1", ""
}
func (this_ *PeepholeReplaceKnownMethods__Anon_69) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.PeepholeReplaceKnownMethods_ConcatFunctionCallBase.IdentityHash_()
}
func (*PeepholeReplaceKnownMethods__Anon_69) HasIdentityEquals_() {}
func (*PeepholeReplaceKnownMethods__Anon_69) Implements_jscomp_PeepholeReplaceKnownMethods_ConcatFunctionCall() {
}

func (this_ *PeepholeReplaceKnownMethods__Anon_69) initFields_PeepholeReplaceKnownMethods__Anon_69() {
	_ = this_
}

func NewPeepholeReplaceKnownMethods__Anon_69(sp0_ pkg_rhino.Node, sp1_ pkg_rhino.Node, sp2_ pkg_rhino.Node) *PeepholeReplaceKnownMethods__Anon_69 {
	this_ := &PeepholeReplaceKnownMethods__Anon_69{}
	InitPeepholeReplaceKnownMethods_ConcatFunctionCall(&this_.PeepholeReplaceKnownMethods_ConcatFunctionCallBase, sp0_, sp1_, sp2_)
	this_.initFields_PeepholeReplaceKnownMethods__Anon_69()
	return this_
}
