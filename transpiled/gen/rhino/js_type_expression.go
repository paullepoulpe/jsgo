/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * The Original Code is Rhino code, released
 * May 6, 1999.
 *
 * The Initial Developer of the Original Code is
 * Netscape Communications Corporation.
 * Portions created by the Initial Developer are Copyright (C) 1997-1999
 * the Initial Developer. All Rights Reserved.
 *
 * Contributor(s):
 *   Bob Jervis
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */
// Port of com/google/javascript/rhino/JSTypeExpression.java (generated mechanically by j2go).

package rhino

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.rhino.JSTypeExpression
type JSTypeExpression struct {
	sourceName string `gson:"sourceName"`
	root       Node   `gson:"root"`
	valid_     bool   `gson:"-"`
}

func (this_ JSTypeExpression) BindSelf_(s any) {
	_ = s
}
func (this_ JSTypeExpression) IsJavaNil_() bool { return !this_.valid_ }
func (JSTypeExpression) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.JSTypeExpression", "JSTypeExpression"
}
func (JSTypeExpression) HasIdentityEquals_() {}
func (this_ JSTypeExpression) IdentityHash_() int32 {
	if !this_.valid_ {
		return 0
	}
	var h_ int32 = 1
	h_ = 31*h_ + int32(this_.root)
	{
		hs_5 := this_.sourceName
		h_ = 31 * h_
		if hs_5 != j2gort.NullString {
			h_ += javalang.StringHashCode(hs_5)
		}
	}
	return h_
}
func (this_ JSTypeExpression) FastValHash_() (int32, bool, bool) {
	if this_.IsJavaNil_() {
		return 0, false, true
	}
	return this_.IdentityHash_(), true, true
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.JSTypeExpression", "com.google.javascript.rhino.JSTypeExpression", func(v_ any) bool {
		return j2gort.InstanceOf[JSTypeExpression](v_)
	})
}

func StaticGet_JSTypeExpression_SerialVersionUID() int64 { return int64(1) }
func StaticGet_JSTypeExpression_IMPLICIT_TEMPLATE_BOUND_SOURCE() string {
	return "<IMPLICIT_TEMPLATE_BOUND>"
}

var static_JSTypeExpression_IMPLICIT_TEMPLATE_BOUND JSTypeExpression

func StaticGet_JSTypeExpression_IMPLICIT_TEMPLATE_BOUND() JSTypeExpression {
	clinit_JSTypeExpression()
	return static_JSTypeExpression_IMPLICIT_TEMPLATE_BOUND
}
func StaticSet_JSTypeExpression_IMPLICIT_TEMPLATE_BOUND(v JSTypeExpression) {
	clinit_JSTypeExpression()
	static_JSTypeExpression_IMPLICIT_TEMPLATE_BOUND = v
}

var clinitDone_JSTypeExpression uint32
var clinitState_JSTypeExpression j2gort.ClinitState

func clinit_JSTypeExpression() {
	if atomic.LoadUint32(&clinitDone_JSTypeExpression) != 0 {
		return
	}
	clinitSlow_JSTypeExpression()
}

func clinitSlow_JSTypeExpression() {
	if !j2gort.EnterClinit(&clinitState_JSTypeExpression, &clinitDone_JSTypeExpression) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_JSTypeExpression, &clinitDone_JSTypeExpression)
	static_JSTypeExpression_IMPLICIT_TEMPLATE_BOUND = NewJSTypeExpression(NewNode_1(StaticGet_Token_QMARK()), "<IMPLICIT_TEMPLATE_BOUND>")
	StaticGet_JSTypeExpression_IMPLICIT_TEMPLATE_BOUND().GetRoot().SetStaticSourceFile(NewSimpleSourceFile("<IMPLICIT_TEMPLATE_BOUND>", StaticGet_StaticSourceFile_SourceKind_STRONG()))
}

func (this_ *JSTypeExpression) initFields_JSTypeExpression() {
	_ = this_
	this_.valid_ = true
	this_.sourceName = j2gort.NullString
}

// Java constructor: JSTypeExpression(com.google.javascript.rhino.Node root, java.lang.String sourceName)
func NewJSTypeExpression(root Node, sourceName string) JSTypeExpression {
	clinit_JSTypeExpression()
	this_ := JSTypeExpression{}
	InitJSTypeExpression(&this_, root, sourceName)
	return this_
}

func InitJSTypeExpression(this_ *JSTypeExpression, root Node, sourceName string) {
	_ = root
	_ = sourceName
	this_.initFields_JSTypeExpression()
	if root.GetParent() != 0 {
		this_.root = root.CloneTree_0()
	} else {
		this_.root = root
	}
	this_.sourceName = sourceName
}

// Java method: JSTypeExpression.replaceNamesWithUnknownType(java.util.Set<java.lang.String> names)
func (this_ JSTypeExpression) ReplaceNamesWithUnknownType(names j2gort.Ext_java_util_Set) JSTypeExpression {
	if !this_.valid_ {
		j2gort.ThrowNPE()
	}
	_ = names
	var oldExprRoot Node = this_.root.CloneTree_0()
	_ = oldExprRoot
	var newExprRoot Node = JSTypeExpression_ReplaceNames(oldExprRoot, names)
	_ = newExprRoot
	var newTypeExpression JSTypeExpression = NewJSTypeExpression(newExprRoot, this_.sourceName)
	_ = newTypeExpression
	return newTypeExpression
}

// Java static method: JSTypeExpression.replaceNames(com.google.javascript.rhino.Node n, java.util.Set<java.lang.String> names)
func JSTypeExpression_ReplaceNames(n Node, names j2gort.Ext_java_util_Set) Node {
	clinit_JSTypeExpression()
	_ = n
	_ = names
	if n == 0 {
		return Node(0)
	}
	{
		var child Node = n.GetFirstChild()
		_ = child
		for child != 0 {
			var next Node = child.GetNext()
			_ = next
			JSTypeExpression_ReplaceNames(child, names)
			child = next
		}
	}
	if n.IsStringLit() && j2gort.SetContainsStr(names, JSTypeExpression_BaseName(n.GetString())) {
		var qMark Node = NewNode_1(StaticGet_Token_QMARK())
		_ = qMark
		if n.HasChildren() && !(n.GetFirstChild().IsBlock()) {
			qMark.AddChildrenToBack(n.RemoveChildren())
		}
		if n.HasParent() {
			n.ReplaceWith(qMark)
		}
		return qMark
	}
	return n
}

// Java static method: JSTypeExpression.baseName(java.lang.String name)
func JSTypeExpression_BaseName(name string) string {
	clinit_JSTypeExpression()
	_ = name
	var dotIndex int32 = j2gort.JStrIndexOfCp(name, int32(uint16(46)))
	_ = dotIndex
	if dotIndex == int32(-1) {
		return name
	} else {
		return j2gort.JStrSub(name, int32(0), dotIndex)
	}
}

// Java method: JSTypeExpression.getAllTypeNodes()
func (this_ JSTypeExpression) GetAllTypeNodes() j2gort.Ext_com_google_common_collect_ImmutableList {
	if !this_.valid_ {
		j2gort.ThrowNPE()
	}
	var builder *j2gort.Ext_com_google_common_collect_ImmutableList_Builder = (&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{})
	_ = builder
	JSTypeExpression_VisitAllTypeNodes(this_.root, func() j2gort.Ext_java_util_function_Consumer {
		bound_1 := builder
		return j2gort.NewLambda_Ext_java_util_function_Consumer(func(mrp0_ any) { bound_1.Add_1_any(mrp0_) })
	}())
	return builder.BuildImmList_()
}

// Java method: JSTypeExpression.getAllTypeNames()
func (this_ JSTypeExpression) GetAllTypeNames() j2gort.Ext_com_google_common_collect_ImmutableSet {
	if !this_.valid_ {
		j2gort.ThrowNPE()
	}
	var builder j2gort.Ext_com_google_common_collect_ImmutableSet_Builder = j2gort.Ext_com_google_common_collect_ImmutableSet_S_Builder()
	_ = builder
	JSTypeExpression_VisitAllTypeNodes(this_.root, j2gort.NewLambda_Ext_java_util_function_Consumer(func(lp_2 any) {
		n := j2gort.CastNode[Node](lp_2)
		_ = n
		j2gort.ImmSetBuilderAddStr(builder, n.GetString())
	}))
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](builder.Build())
}

// Java static method: JSTypeExpression.visitAllTypeNodes(com.google.javascript.rhino.Node n, java.util.function.Consumer<com.google.javascript.rhino.Node> visitor)
func JSTypeExpression_VisitAllTypeNodes(n Node, visitor j2gort.Ext_java_util_function_Consumer) {
	clinit_JSTypeExpression()
	_ = n
	_ = visitor
	if n == 0 {
		return
	}
	{
		var child Node = n.GetFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			JSTypeExpression_VisitAllTypeNodes(child, visitor)
		}
	}
	if n.IsStringLit() {
		visitor.Accept(j2gort.BoxNode(n))
	}
}

// Java static method: JSTypeExpression.makeOptionalArg(com.google.javascript.rhino.JSTypeExpression expr)
func JSTypeExpression_MakeOptionalArg(expr JSTypeExpression) JSTypeExpression {
	clinit_JSTypeExpression()
	_ = expr
	if expr.IsOptionalArg() || expr.IsVarArgs() {
		return expr
	} else {
		var equals Node = NewNode_2(StaticGet_Token_EQUALS(), expr.root.CloneTree_0())
		_ = equals
		equals.ClonePropsFrom(expr.root)
		return NewJSTypeExpression(equals, expr.sourceName)
	}
}

// Java static method: JSTypeExpression.makeVarArgs(com.google.javascript.rhino.JSTypeExpression expr)
func JSTypeExpression_MakeVarArgs(expr JSTypeExpression) JSTypeExpression {
	clinit_JSTypeExpression()
	_ = expr
	if expr.IsOptionalArg() || expr.IsVarArgs() {
		return expr
	} else {
		var equals Node = NewNode_2(StaticGet_Token_ITER_REST(), expr.root.CloneTree_0())
		_ = equals
		equals.ClonePropsFrom(expr.root)
		return NewJSTypeExpression(equals, expr.sourceName)
	}
}

// Java method: JSTypeExpression.isOptionalArg()
func (this_ JSTypeExpression) IsOptionalArg() bool {
	if !this_.valid_ {
		j2gort.ThrowNPE()
	}
	return (this_.root.RawToken_() == uint8(160))
}

// Java method: JSTypeExpression.isVarArgs()
func (this_ JSTypeExpression) IsVarArgs() bool {
	if !this_.valid_ {
		j2gort.ThrowNPE()
	}
	return (this_.root.RawToken_() == uint8(125))
}

// Java method: JSTypeExpression.evaluate(com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable StaticTypedScope scope, com.google.javascript.rhino.jstype.JSTypeRegistry registry)
func (this_ JSTypeExpression) Evaluate(scope StaticTypedScope, registry *JSTypeRegistry) JSType {
	if !this_.valid_ {
		j2gort.ThrowNPE()
	}
	_ = scope
	_ = registry
	var type_ JSType = registry.CreateTypeFromCommentNode_3(this_.root, this_.sourceName, scope)
	_ = type_
	this_.root.SetJSType(type_)
	return type_
}

// Java method: JSTypeExpression.isEquivalentTo(com.google.javascript.rhino.@org.jspecify.annotations.Nullable JSTypeExpression other)
func (this_ JSTypeExpression) IsEquivalentTo(other JSTypeExpression) bool {
	if !this_.valid_ {
		j2gort.ThrowNPE()
	}
	_ = other
	return ((other != (JSTypeExpression{})) && this_.root.IsEquivalentTo_1(other.root))
}

// Java method: JSTypeExpression.getRoot()
func (this_ JSTypeExpression) GetRoot() Node {
	if !this_.valid_ {
		j2gort.ThrowNPE()
	}
	return this_.root
}

// Java method: JSTypeExpression.getSourceName()
func (this_ JSTypeExpression) GetSourceName() string {
	if !this_.valid_ {
		j2gort.ThrowNPE()
	}
	return this_.sourceName
}

// Java method: JSTypeExpression.toString()
func (this_ JSTypeExpression) String() string {
	if !this_.valid_ {
		j2gort.ThrowNPE()
	}
	return ("type: " + j2gort.StrValString(this_.root.ToStringTree()))
}

// Java method: JSTypeExpression.copy()
func (this_ JSTypeExpression) Copy() JSTypeExpression {
	if !this_.valid_ {
		j2gort.ThrowNPE()
	}
	return NewJSTypeExpression(this_.root.CloneTree_0(), this_.sourceName)
}

// Java method: JSTypeExpression.isExplicitUnknownTemplateBound()
func (this_ JSTypeExpression) IsExplicitUnknownTemplateBound() bool {
	if !this_.valid_ {
		j2gort.ThrowNPE()
	}
	return (((this_.root.RawToken_() == uint8(158)) && !(this_.root.HasChildren())) && !(func() bool {
		sr_3 := this_.sourceName
		se_4 := "<IMPLICIT_TEMPLATE_BOUND>"
		return j2gort.CheckNotNullString(sr_3) == se_4
	}()))
}

// Java method: JSTypeExpression.getRecordPropertyNames()
func (this_ JSTypeExpression) GetRecordPropertyNames() j2gort.Ext_com_google_common_collect_ImmutableSet {
	if !this_.valid_ {
		j2gort.ThrowNPE()
	}
	var builder j2gort.Ext_com_google_common_collect_ImmutableSet_Builder = j2gort.Ext_com_google_common_collect_ImmutableSet_S_Builder()
	_ = builder
	JSTypeExpression_GetRecordPropertyNamesRecursive(this_.root, builder)
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](builder.Build())
}

// Java static method: JSTypeExpression.getRecordPropertyNamesRecursive(com.google.javascript.rhino.Node n, com.google.common.collect.ImmutableSet.Builder<java.lang.String> names)
func JSTypeExpression_GetRecordPropertyNamesRecursive(n Node, names j2gort.Ext_com_google_common_collect_ImmutableSet_Builder) {
	clinit_JSTypeExpression()
	_ = n
	_ = names
	if n == 0 {
		return
	}
	{
		var child Node = n.GetFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			JSTypeExpression_GetRecordPropertyNamesRecursive(child, names)
		}
	}
	if n.IsStringKey() {
		j2gort.ImmSetBuilderAddStrVoid(names, n.GetString())
	}
}
