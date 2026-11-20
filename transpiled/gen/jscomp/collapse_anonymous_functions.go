/*
 * Copyright 2008 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/CollapseAnonymousFunctions.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.CollapseAnonymousFunctions
type CollapseAnonymousFunctions struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler AbstractCompiler `gson:"compiler"`
}

func (this_ *CollapseAnonymousFunctions) BindSelf_(s any) {
	_ = s
}
func (this_ *CollapseAnonymousFunctions) IsJavaNil_() bool { return this_ == nil }
func (*CollapseAnonymousFunctions) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CollapseAnonymousFunctions", "CollapseAnonymousFunctions"
}
func (this_ *CollapseAnonymousFunctions) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*CollapseAnonymousFunctions) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CollapseAnonymousFunctions", "com.google.javascript.jscomp.CollapseAnonymousFunctions", func(v_ any) bool {
		return j2gort.InstanceOf[*CollapseAnonymousFunctions](v_)
	})
}
func (*CollapseAnonymousFunctions) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*CollapseAnonymousFunctions) Implements_jscomp_CompilerPass()                            {}
func (*CollapseAnonymousFunctions) Implements_jscomp_NodeTraversal_Callback()                  {}

func (this_ *CollapseAnonymousFunctions) initFields_CollapseAnonymousFunctions() {
	_ = this_
}

// Java constructor: CollapseAnonymousFunctions(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewCollapseAnonymousFunctions(compiler AbstractCompiler) *CollapseAnonymousFunctions {
	this_ := &CollapseAnonymousFunctions{}
	InitCollapseAnonymousFunctions(this_, compiler)
	return this_
}

func InitCollapseAnonymousFunctions(this_ *CollapseAnonymousFunctions, compiler AbstractCompiler) {
	_ = compiler
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_CollapseAnonymousFunctions()
	if !(compiler.GetLifeCycleStage().IsNormalized()) {
		j2gort.CheckArgument(false)
	}
	this_.compiler = compiler
}

// Java method: CollapseAnonymousFunctions.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *CollapseAnonymousFunctions) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, this_)
}

// Java method: CollapseAnonymousFunctions.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CollapseAnonymousFunctions) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if !((n.IsVar() || n.IsLet()) || n.IsConst()) {
		return
	}
	var grandparent pkg_rhino.Node = parent.GetParent()
	_ = grandparent
	if !(parent.IsScript() || (((grandparent != 0) && grandparent.IsFunction()) && parent.IsBlock())) {
		return
	}
	var name pkg_rhino.Node = n.GetOnlyChild()
	_ = name
	if !(name.IsName()) {
		return
	}
	var value pkg_rhino.Node = name.GetFirstChild()
	_ = value
	if (((value != 0) && value.IsFunction()) && !(value.IsArrowFunction())) && !(this_.priv_CollapseAnonymousFunctions_isRecursiveFunction_1394(value)) {
		var fnName pkg_rhino.Node = value.GetFirstChild()
		_ = fnName
		if fnName.GetBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME()) && !(name.GetBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME())) {
			fnName.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME(), false)
		}
		fnName.SetString(name.GetString())
		NodeUtil_CopyNameAnnotations(name, fnName)
		value.Detach()
		n.ReplaceWith(value)
		if !(t.InGlobalScope()) && NodeUtil_IsHoistedFunctionDeclaration(value) {
			parent.AddChildToFront(value.Detach())
		}
		this_.compiler.ReportChangeToChangeScope(value)
		t.ReportCodeChange_0()
	}
}

// Java method: CollapseAnonymousFunctions.isRecursiveFunction(com.google.javascript.rhino.Node function)
func (this_ *CollapseAnonymousFunctions) priv_CollapseAnonymousFunctions_isRecursiveFunction_1394(function pkg_rhino.Node) bool {
	_ = function
	var name pkg_rhino.Node = function.GetFirstChild()
	_ = name
	if j2gort.CheckNotNullString(name.GetString()) == "" {
		return false
	}
	var args pkg_rhino.Node = name.GetNext()
	_ = args
	var body pkg_rhino.Node = args.GetNext()
	_ = body
	return this_.priv_CollapseAnonymousFunctions_containsName_1395(body, name.GetString())
}

// Java method: CollapseAnonymousFunctions.containsName(com.google.javascript.rhino.Node n, java.lang.String name)
func (this_ *CollapseAnonymousFunctions) priv_CollapseAnonymousFunctions_containsName_1395(n pkg_rhino.Node, name string) bool {
	_ = n
	_ = name
	if n.IsName() && (j2gort.CheckNotNullString(n.GetString()) == name) {
		return true
	}
	{
		var child pkg_rhino.Node = n.GetFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			if this_.priv_CollapseAnonymousFunctions_containsName_1395(child, name) {
				return true
			}
		}
	}
	return false
}
func (*CollapseAnonymousFunctions) AlwaysShouldTraverse_()      {}
func (*CollapseAnonymousFunctions) CanReleaseSyntacticScopes_() {}
