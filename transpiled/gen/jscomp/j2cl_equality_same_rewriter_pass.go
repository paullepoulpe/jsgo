/*
 * Copyright 2016 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/J2clEqualitySameRewriterPass.java (generated mechanically by j2go).

package jscomp

import (
	pkg_colors "github.com/paullepoulpe/jsgo/transpiled/gen/colors"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.J2clEqualitySameRewriterPass
type J2clEqualitySameRewriterPass struct {
	AbstractPeepholeOptimizationBase
	useTypes bool `gson:"useTypes"`
	shouldRunJ2clPasses bool `gson:"shouldRunJ2clPasses"`
}

func (this_ *J2clEqualitySameRewriterPass) BindSelf_(s any) {
	_ = s
	this_.AbstractPeepholeOptimizationBase.BindSelf_(s)
}
func (this_ *J2clEqualitySameRewriterPass) IsJavaNil_() bool { return this_ == nil }
func (*J2clEqualitySameRewriterPass) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.J2clEqualitySameRewriterPass", "J2clEqualitySameRewriterPass" }
func (this_ *J2clEqualitySameRewriterPass) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.AbstractPeepholeOptimizationBase.IdentityHash_()
}
func (*J2clEqualitySameRewriterPass) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.J2clEqualitySameRewriterPass", "com.google.javascript.jscomp.J2clEqualitySameRewriterPass", func(v_ any) bool {
		return j2gort.InstanceOf[*J2clEqualitySameRewriterPass](v_)
	})
}
func (*J2clEqualitySameRewriterPass) Implements_jscomp_AbstractPeepholeOptimization() {}

func (this_ *J2clEqualitySameRewriterPass) initFields_J2clEqualitySameRewriterPass() {
	_ = this_
}

// Java constructor: J2clEqualitySameRewriterPass(boolean useTypes)
func NewJ2clEqualitySameRewriterPass(useTypes bool) *J2clEqualitySameRewriterPass {
	this_ := &J2clEqualitySameRewriterPass{}
	this_.BindSelf_(this_)
	InitJ2clEqualitySameRewriterPass(this_, useTypes)
	return this_
}

func InitJ2clEqualitySameRewriterPass(this_ *J2clEqualitySameRewriterPass, useTypes bool) {
	_ = useTypes
	InitAbstractPeepholeOptimization(&this_.AbstractPeepholeOptimizationBase)
	this_.initFields_J2clEqualitySameRewriterPass()
	this_.useTypes = useTypes
}

// Java method: J2clEqualitySameRewriterPass.beginTraversal(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *J2clEqualitySameRewriterPass) BeginTraversal(compiler AbstractCompiler) {
	_ = compiler
	this_.AbstractPeepholeOptimizationBase.BeginTraversal(compiler)
	this_.shouldRunJ2clPasses = J2clSourceFileChecker_ShouldRunJ2clPasses(compiler)
}

// Java method: J2clEqualitySameRewriterPass.optimizeSubtree(com.google.javascript.rhino.Node node)
func (this_ *J2clEqualitySameRewriterPass) OptimizeSubtree(node pkg_rhino.Node) pkg_rhino.Node {
	_ = node
	if (!(this_.shouldRunJ2clPasses)) {
		return node
	}
	if (!(node.HasXChildren(int32(3)))) {
		return node
	}
	var replacement pkg_rhino.Node = pkg_rhino.Node(0)
	_ = replacement
	if (J2clEqualitySameRewriterPass_IsStringEqualsMethod(node)) {
		replacement = this_.priv_J2clEqualitySameRewriterPass_trySubstituteStringEquals_3467(node)
	} else if (J2clEqualitySameRewriterPass_IsEqualitySameCall(node)) {
		replacement = this_.priv_J2clEqualitySameRewriterPass_trySubstituteEqualitySame_3468(node)
	}
	if ((replacement == 0)) {
		return node
	}
	replacement.SrcrefIfMissing(node)
	node.ReplaceWith(replacement)
	this_.ReportChangeToEnclosingScope(replacement)
	return replacement
}

// Java method: J2clEqualitySameRewriterPass.trySubstituteStringEquals(com.google.javascript.rhino.Node callNode)
func (this_ *J2clEqualitySameRewriterPass) priv_J2clEqualitySameRewriterPass_trySubstituteStringEquals_3467(callNode pkg_rhino.Node) pkg_rhino.Node {
	_ = callNode
	var firstExprValue J2clEqualitySameRewriterPass_NodeValue = J2clEqualitySameRewriterPass_GetKnownLiteralValue(callNode.GetSecondChild())
	_ = firstExprValue
	if (((firstExprValue == StaticGet_J2clEqualitySameRewriterPass_NodeValue_UNKNOWN()) || (firstExprValue == StaticGet_J2clEqualitySameRewriterPass_NodeValue_NULL_OR_UNDEFINED()))) {
		return pkg_rhino.Node(0)
	}
	return this_.priv_J2clEqualitySameRewriterPass_trySubstituteEqualitySame_3468(callNode)
}

// Java method: J2clEqualitySameRewriterPass.trySubstituteEqualitySame(com.google.javascript.rhino.Node callNode)
func (this_ *J2clEqualitySameRewriterPass) priv_J2clEqualitySameRewriterPass_trySubstituteEqualitySame_3468(callNode pkg_rhino.Node) pkg_rhino.Node {
	_ = callNode
	var firstExpr pkg_rhino.Node = callNode.GetSecondChild()
	_ = firstExpr
	var firstExprValue J2clEqualitySameRewriterPass_NodeValue = J2clEqualitySameRewriterPass_GetKnownLiteralValue(firstExpr)
	_ = firstExprValue
	var secondExpr pkg_rhino.Node = callNode.GetLastChild()
	_ = secondExpr
	var secondExprValue J2clEqualitySameRewriterPass_NodeValue = J2clEqualitySameRewriterPass_GetKnownLiteralValue(secondExpr)
	_ = secondExprValue
	if (((firstExprValue == StaticGet_J2clEqualitySameRewriterPass_NodeValue_UNKNOWN()) && (secondExprValue == StaticGet_J2clEqualitySameRewriterPass_NodeValue_UNKNOWN()))) {
		return pkg_rhino.Node(0)
	}
	if ((firstExprValue == StaticGet_J2clEqualitySameRewriterPass_NodeValue_NULL_OR_UNDEFINED())) {
		return this_.priv_J2clEqualitySameRewriterPass_rewriteNullCheck_3469(secondExpr, firstExpr)
	}
	if ((secondExprValue == StaticGet_J2clEqualitySameRewriterPass_NodeValue_NULL_OR_UNDEFINED())) {
		return this_.priv_J2clEqualitySameRewriterPass_rewriteNullCheck_3469(firstExpr, secondExpr)
	}
	if (((firstExprValue == StaticGet_J2clEqualitySameRewriterPass_NodeValue_NON_NULL()) || (secondExprValue == StaticGet_J2clEqualitySameRewriterPass_NodeValue_NON_NULL()))) {
		return this_.priv_J2clEqualitySameRewriterPass_rewriteAsStrictEq_3471(firstExpr, secondExpr)
	}
	if !(((firstExprValue == StaticGet_J2clEqualitySameRewriterPass_NodeValue_NUMBER()) || (secondExprValue == StaticGet_J2clEqualitySameRewriterPass_NodeValue_NUMBER()))) { j2gort.CheckState(false) }
	return this_.priv_J2clEqualitySameRewriterPass_rewriteNumberCheck_3472(firstExpr, secondExpr)
}

// Java method: J2clEqualitySameRewriterPass.rewriteNullCheck(com.google.javascript.rhino.Node expr, com.google.javascript.rhino.Node nullExpression)
func (this_ *J2clEqualitySameRewriterPass) priv_J2clEqualitySameRewriterPass_rewriteNullCheck_3469(expr pkg_rhino.Node, nullExpression pkg_rhino.Node) pkg_rhino.Node {
	_ = expr
	_ = nullExpression
	expr.Detach()
	nullExpression.Detach()
	if ((this_.useTypes && this_.priv_J2clEqualitySameRewriterPass_canOnlyBeObject_3470(expr))) {
		return pkg_rhino.IR_Not(expr)
	}
	return pkg_rhino.IR_Eq(expr, nullExpression)
}

// Java method: J2clEqualitySameRewriterPass.canOnlyBeObject(com.google.javascript.rhino.Node n)
func (this_ *J2clEqualitySameRewriterPass) priv_J2clEqualitySameRewriterPass_canOnlyBeObject_3470(n pkg_rhino.Node) bool {
	_ = n
	var color *pkg_colors.Color = n.GetColor()
	_ = color
	if ((color == nil)) {
		return false
	}
	if (color.IsUnion()) {
		color = color.SubtractNullOrVoid()
	}
	return (((!(color.IsUnion()) && !(color.IsPrimitive())) && !(color.Equals(j2gort.PtrToAny(pkg_colors.StaticGet_StandardColors_UNKNOWN())))) && !(color.Equals(j2gort.PtrToAny(pkg_colors.StaticGet_StandardColors_TOP_OBJECT()))))
}

// Java method: J2clEqualitySameRewriterPass.rewriteAsStrictEq(com.google.javascript.rhino.Node firstExpr, com.google.javascript.rhino.Node secondExpr)
func (this_ *J2clEqualitySameRewriterPass) priv_J2clEqualitySameRewriterPass_rewriteAsStrictEq_3471(firstExpr pkg_rhino.Node, secondExpr pkg_rhino.Node) pkg_rhino.Node {
	_ = firstExpr
	_ = secondExpr
	firstExpr.Detach()
	secondExpr.Detach()
	return pkg_rhino.IR_Sheq(firstExpr, secondExpr)
}

// Java method: J2clEqualitySameRewriterPass.rewriteNumberCheck(com.google.javascript.rhino.Node firstExpr, com.google.javascript.rhino.Node secondExpr)
func (this_ *J2clEqualitySameRewriterPass) priv_J2clEqualitySameRewriterPass_rewriteNumberCheck_3472(firstExpr pkg_rhino.Node, secondExpr pkg_rhino.Node) pkg_rhino.Node {
	_ = firstExpr
	_ = secondExpr
	var firstValue *j2gort.BoxDouble = NodeUtil_GetNumberValue(firstExpr)
	_ = firstValue
	var secondValue *j2gort.BoxDouble = NodeUtil_GetNumberValue(secondExpr)
	_ = secondValue
	if (((firstValue != nil) && (secondValue != nil))) {
		firstExpr.Detach()
		secondExpr.Detach()
		return NodeUtil_BooleanNode(firstValue.Equals(j2gort.PtrToAny(secondValue)))
	}
	if ((J2clEqualitySameRewriterPass_IsSafeNumber(firstValue) || J2clEqualitySameRewriterPass_IsSafeNumber(secondValue))) {
		return this_.priv_J2clEqualitySameRewriterPass_rewriteAsStrictEq_3471(firstExpr, secondExpr)
	}
	if ((this_.useTypes && ((this_.priv_J2clEqualitySameRewriterPass_canOnlyBeObject_3470(firstExpr) || this_.priv_J2clEqualitySameRewriterPass_canOnlyBeObject_3470(secondExpr))))) {
		return this_.priv_J2clEqualitySameRewriterPass_rewriteAsStrictEq_3471(firstExpr, secondExpr)
	}
	return pkg_rhino.Node(0)
}

// Java static method: J2clEqualitySameRewriterPass.isSafeNumber(java.lang.Double d)
func J2clEqualitySameRewriterPass_IsSafeNumber(d *j2gort.BoxDouble) bool {
	_ = d
	return (((d != nil) && (j2gort.UnboxDouble(d) != float64(int32(0)))) && !(d.IsNaN()))
}

// Java static method: J2clEqualitySameRewriterPass.getKnownLiteralValue(com.google.javascript.rhino.Node n)
func J2clEqualitySameRewriterPass_GetKnownLiteralValue(n pkg_rhino.Node) J2clEqualitySameRewriterPass_NodeValue {
	_ = n
	switch (NodeUtil_GetKnownValueType(n)).Ordinal() {
	case int32(2) /* ValueType.VOID */:
		return func() J2clEqualitySameRewriterPass_NodeValue { if NodeUtil_CanBeSideEffected_1(n) { return StaticGet_J2clEqualitySameRewriterPass_NodeValue_UNKNOWN() }; return StaticGet_J2clEqualitySameRewriterPass_NodeValue_NULL_OR_UNDEFINED() }()
	case int32(1) /* ValueType.NULL */:
		return StaticGet_J2clEqualitySameRewriterPass_NodeValue_NULL_OR_UNDEFINED()
	case int32(5) /* ValueType.STRING */, int32(6) /* ValueType.BOOLEAN */, int32(7) /* ValueType.OBJECT */, int32(4) /* ValueType.BIGINT */:
		return StaticGet_J2clEqualitySameRewriterPass_NodeValue_NON_NULL()
	case int32(3) /* ValueType.NUMBER */:
		return StaticGet_J2clEqualitySameRewriterPass_NodeValue_NUMBER()
	case int32(0) /* ValueType.UNDETERMINED */:
		return StaticGet_J2clEqualitySameRewriterPass_NodeValue_UNKNOWN()
	}
	panic("unreachable switch expression")
}

// Java static method: J2clEqualitySameRewriterPass.isEqualitySameCall(com.google.javascript.rhino.Node node)
func J2clEqualitySameRewriterPass_IsEqualitySameCall(node pkg_rhino.Node) bool {
	_ = node
	return (node.IsCall() && J2clEqualitySameRewriterPass_HasName(node.GetFirstChild(), "Equality", "$same"))
}

// Java static method: J2clEqualitySameRewriterPass.isStringEqualsMethod(com.google.javascript.rhino.Node node)
func J2clEqualitySameRewriterPass_IsStringEqualsMethod(node pkg_rhino.Node) bool {
	_ = node
	return (node.IsCall() && J2clEqualitySameRewriterPass_HasName(node.GetFirstChild(), "String", "m_equals__java_lang_String__java_lang_Object__boolean"))
}

// Java static method: J2clEqualitySameRewriterPass.hasName(com.google.javascript.rhino.Node fnName, java.lang.String className, java.lang.String methodName)
func J2clEqualitySameRewriterPass_HasName(fnName pkg_rhino.Node, className string, methodName string) bool {
	_ = fnName
	_ = className
	_ = methodName
	if (!(fnName.IsQualifiedName())) {
		return false
	}
	var originalQname string = fnName.GetOriginalQualifiedName()
	_ = originalQname
	return (javalang.StringEndsWith(j2gort.CheckNotNullString(originalQname), j2gort.CheckNotNullString(methodName)) && (j2gort.JStrIndexOf(originalQname, j2gort.CheckNotNullString(className)) >= 0))
}

// Java enum: com.google.javascript.jscomp.J2clEqualitySameRewriterPass.NodeValue
type J2clEqualitySameRewriterPass_NodeValue uint8

const (
	enum_J2clEqualitySameRewriterPass_NodeValue_UNKNOWN J2clEqualitySameRewriterPass_NodeValue = 1
	enum_J2clEqualitySameRewriterPass_NodeValue_NULL_OR_UNDEFINED J2clEqualitySameRewriterPass_NodeValue = 2
	enum_J2clEqualitySameRewriterPass_NodeValue_NUMBER J2clEqualitySameRewriterPass_NodeValue = 3
	enum_J2clEqualitySameRewriterPass_NodeValue_NON_NULL J2clEqualitySameRewriterPass_NodeValue = 4
)

var enumBaseTable_J2clEqualitySameRewriterPass_NodeValue = [...]j2gort.EnumBase{
	{Name_: "UNKNOWN", Ordinal_: 0},
	{Name_: "NULL_OR_UNDEFINED", Ordinal_: 1},
	{Name_: "NUMBER", Ordinal_: 2},
	{Name_: "NON_NULL", Ordinal_: 3},
}
var enumValues_J2clEqualitySameRewriterPass_NodeValue = []J2clEqualitySameRewriterPass_NodeValue{1, 2, 3, 4}

func StaticGet_J2clEqualitySameRewriterPass_NodeValue_UNKNOWN() J2clEqualitySameRewriterPass_NodeValue {
	return enum_J2clEqualitySameRewriterPass_NodeValue_UNKNOWN
}
func StaticGet_J2clEqualitySameRewriterPass_NodeValue_NULL_OR_UNDEFINED() J2clEqualitySameRewriterPass_NodeValue {
	return enum_J2clEqualitySameRewriterPass_NodeValue_NULL_OR_UNDEFINED
}
func StaticGet_J2clEqualitySameRewriterPass_NodeValue_NUMBER() J2clEqualitySameRewriterPass_NodeValue {
	return enum_J2clEqualitySameRewriterPass_NodeValue_NUMBER
}
func StaticGet_J2clEqualitySameRewriterPass_NodeValue_NON_NULL() J2clEqualitySameRewriterPass_NodeValue {
	return enum_J2clEqualitySameRewriterPass_NodeValue_NON_NULL
}

func J2clEqualitySameRewriterPass_NodeValue_Values() []J2clEqualitySameRewriterPass_NodeValue {
	return j2gort.CloneSlice(enumValues_J2clEqualitySameRewriterPass_NodeValue)
}

func J2clEqualitySameRewriterPass_NodeValue_ValueOf(name string) J2clEqualitySameRewriterPass_NodeValue {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_J2clEqualitySameRewriterPass_NodeValue {
		if enumBaseTable_J2clEqualitySameRewriterPass_NodeValue[i_].Name_ == name {
			return J2clEqualitySameRewriterPass_NodeValue(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.J2clEqualitySameRewriterPass.NodeValue." + name))
}

func (this_ J2clEqualitySameRewriterPass_NodeValue) IsJavaNil_() bool { return this_ == 0 }
func (J2clEqualitySameRewriterPass_NodeValue) HasIdentityEquals_() {}
func (J2clEqualitySameRewriterPass_NodeValue) Implements_j2gort_Ext_java_lang_Enum() {}
func (J2clEqualitySameRewriterPass_NodeValue) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ J2clEqualitySameRewriterPass_NodeValue) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase { return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()} }
func (J2clEqualitySameRewriterPass_NodeValue) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.J2clEqualitySameRewriterPass$NodeValue", "NodeValue" }
func (this_ J2clEqualitySameRewriterPass_NodeValue) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_J2clEqualitySameRewriterPass_NodeValue[this_-1]
}
func (this_ J2clEqualitySameRewriterPass_NodeValue) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ J2clEqualitySameRewriterPass_NodeValue) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ J2clEqualitySameRewriterPass_NodeValue) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_J2clEqualitySameRewriterPass_NodeValue[this_-1].Name_
}
func (this_ J2clEqualitySameRewriterPass_NodeValue) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_J2clEqualitySameRewriterPass_NodeValue[this_-1].Name_
}
func (this_ J2clEqualitySameRewriterPass_NodeValue) Equals(other_ any) bool {
	o_, ok_ := other_.(J2clEqualitySameRewriterPass_NodeValue)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ J2clEqualitySameRewriterPass_NodeValue) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_J2clEqualitySameRewriterPass_NodeValue[this_-1].HashCode()
}
func (this_ J2clEqualitySameRewriterPass_NodeValue) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_J2clEqualitySameRewriterPass_NodeValue[this_-1].CompareTo(other_)
}
func (this_ J2clEqualitySameRewriterPass_NodeValue) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_J2clEqualitySameRewriterPass_NodeValue[this_-1].Name_)
}
func init() {
	vals_ := []any{J2clEqualitySameRewriterPass_NodeValue(1), J2clEqualitySameRewriterPass_NodeValue(2), J2clEqualitySameRewriterPass_NodeValue(3), J2clEqualitySameRewriterPass_NodeValue(4)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.J2clEqualitySameRewriterPass.NodeValue", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.J2clEqualitySameRewriterPass$NodeValue", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.J2clEqualitySameRewriterPass.NodeValue", "com.google.javascript.jscomp.J2clEqualitySameRewriterPass$NodeValue", func(v_ any) bool {
		return j2gort.InstanceOf[J2clEqualitySameRewriterPass_NodeValue](v_)
	})
}

