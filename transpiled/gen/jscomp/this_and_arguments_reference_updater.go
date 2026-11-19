/*
 * Copyright 2024 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/ThisAndArgumentsReferenceUpdater.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.ThisAndArgumentsReferenceUpdater
type ThisAndArgumentsReferenceUpdater struct {
	compiler   AbstractCompiler                                          `gson:"compiler"`
	context    *ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext `gson:"context"`
	astFactory *AstFactory                                               `gson:"astFactory"`
	idHash_    atomic.Int32
}

func (this_ *ThisAndArgumentsReferenceUpdater) BindSelf_(s any) {
	_ = s
}
func (this_ *ThisAndArgumentsReferenceUpdater) IsJavaNil_() bool { return this_ == nil }
func (*ThisAndArgumentsReferenceUpdater) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ThisAndArgumentsReferenceUpdater", "ThisAndArgumentsReferenceUpdater"
}
func (this_ *ThisAndArgumentsReferenceUpdater) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ThisAndArgumentsReferenceUpdater) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ThisAndArgumentsReferenceUpdater", "com.google.javascript.jscomp.ThisAndArgumentsReferenceUpdater", func(v_ any) bool {
		return j2gort.InstanceOf[*ThisAndArgumentsReferenceUpdater](v_)
	})
}
func (*ThisAndArgumentsReferenceUpdater) Implements_jscomp_NodeTraversal_Callback() {}

func StaticGet_ThisAndArgumentsReferenceUpdater_ARGUMENTS_VAR() string { return "$jscomp$arguments" }
func StaticGet_ThisAndArgumentsReferenceUpdater_THIS_VAR() string      { return "$jscomp$this" }
func (this_ *ThisAndArgumentsReferenceUpdater) initFields_ThisAndArgumentsReferenceUpdater() {
	_ = this_
}

// Java constructor: ThisAndArgumentsReferenceUpdater(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.ThisAndArgumentsReferenceUpdater.ThisAndArgumentsContext context, com.google.javascript.jscomp.AstFactory astFactory)
func NewThisAndArgumentsReferenceUpdater(compiler AbstractCompiler, context *ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext, astFactory *AstFactory) *ThisAndArgumentsReferenceUpdater {
	this_ := &ThisAndArgumentsReferenceUpdater{}
	InitThisAndArgumentsReferenceUpdater(this_, compiler, context, astFactory)
	return this_
}

func InitThisAndArgumentsReferenceUpdater(this_ *ThisAndArgumentsReferenceUpdater, compiler AbstractCompiler, context *ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext, astFactory *AstFactory) {
	_ = compiler
	_ = context
	_ = astFactory
	this_.initFields_ThisAndArgumentsReferenceUpdater()
	this_.compiler = compiler
	this_.context = context
	this_.astFactory = astFactory
}

// Java method: ThisAndArgumentsReferenceUpdater.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ThisAndArgumentsReferenceUpdater) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsThis() {
		this_.context.SetNeedsThisVarWithType(AstFactory_S_Type_1_Node(n))
		var name pkg_rhino.Node = this_.astFactory.CreateName_2_string_AstFactory_Type(("$jscomp$this" + "$" + j2gort.StrValString(this_.context.uniqueId)), this_.context.GetThisType()).Srcref(n)
		_ = name
		name.MakeNonIndexable()
		if this_.compiler.GetOptions().PreservesDetailedSourceInfo() {
			name.SetOriginalName("this")
		}
		n.ReplaceWith(name)
	} else if n.IsName() && (j2gort.CheckNotNullString(n.GetString()) == "arguments") {
		this_.context.SetNeedsArgumentsVar()
		var name_2 pkg_rhino.Node = this_.astFactory.CreateName_2_string_AstFactory_Type(("$jscomp$arguments" + "$" + j2gort.StrValString(this_.context.uniqueId)), AstFactory_S_Type_1_Node(n)).Srcref(n)
		_ = name_2
		if this_.compiler.GetOptions().PreservesDetailedSourceInfo() {
			name_2.SetOriginalName("arguments")
		}
		n.ReplaceWith(name_2)
	}
}

// Java method: ThisAndArgumentsReferenceUpdater.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ThisAndArgumentsReferenceUpdater) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	if n.IsMemberFieldDef() {
		return false
	}
	if ((parent != 0) && parent.IsComputedFieldDef()) && (n != parent.GetFirstChild()) {
		return false
	}
	return (!(n.IsFunction()) || n.IsArrowFunction())
}

// Java class: com.google.javascript.jscomp.ThisAndArgumentsReferenceUpdater.ThisAndArgumentsContext
type ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext struct {
	uniqueId           string          `gson:"uniqueId"`
	thisType           AstFactory_Type `gson:"thisType"`
	scopeBody          pkg_rhino.Node  `gson:"scopeBody"`
	lastSuperStatement pkg_rhino.Node  `gson:"lastSuperStatement"`
	idHash_            atomic.Int32
	isConstructor      bool `gson:"isConstructor"`
	needsThisVar       bool `gson:"needsThisVar"`
	needsArgumentsVar  bool `gson:"needsArgumentsVar"`
}

func (this_ *ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext) BindSelf_(s any) {
	_ = s
}
func (this_ *ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext) IsJavaNil_() bool {
	return this_ == nil
}
func (*ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ThisAndArgumentsReferenceUpdater$ThisAndArgumentsContext", "ThisAndArgumentsContext"
}
func (this_ *ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ThisAndArgumentsReferenceUpdater.ThisAndArgumentsContext", "com.google.javascript.jscomp.ThisAndArgumentsReferenceUpdater$ThisAndArgumentsContext", func(v_ any) bool {
		return j2gort.InstanceOf[*ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext](v_)
	})
}

func (this_ *ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext) initFields_ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext() {
	_ = this_

}

// Java constructor: ThisAndArgumentsContext(com.google.javascript.rhino.Node scopeBody, boolean isConstructor, java.lang.String uniqueId)
func NewThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext(scopeBody pkg_rhino.Node, isConstructor bool, uniqueId string) *ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext {
	this_ := &ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext{}
	InitThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext(this_, scopeBody, isConstructor, uniqueId)
	return this_
}

func InitThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext(this_ *ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext, scopeBody pkg_rhino.Node, isConstructor bool, uniqueId string) {
	_ = scopeBody
	_ = isConstructor
	_ = uniqueId
	this_.initFields_ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext()
	this_.scopeBody = scopeBody
	this_.isConstructor = isConstructor
	this_.uniqueId = uniqueId
}

// Java method: ThisAndArgumentsContext.getThisType()
func (this_ *ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext) GetThisType() AstFactory_Type {
	return this_.thisType
}

// Java method: ThisAndArgumentsContext.setNeedsThisVarWithType(com.google.javascript.jscomp.AstFactory.Type type)
func (this_ *ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext) SetNeedsThisVarWithType(type_ AstFactory_Type) *ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext {
	_ = type_
	this_.thisType = type_
	this_.needsThisVar = true
	return this_
}

// Java method: ThisAndArgumentsContext.setNeedsArgumentsVar()
func (this_ *ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext) SetNeedsArgumentsVar() *ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext {
	this_.needsArgumentsVar = true
	return this_
}

// Java method: ThisAndArgumentsContext.addVarDeclarations(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.AstFactory astFactory, com.google.javascript.rhino.Node script)
func (this_ *ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext) AddVarDeclarations(compiler AbstractCompiler, astFactory *AstFactory, script pkg_rhino.Node) {
	_ = compiler
	_ = astFactory
	_ = script
	if this_.needsThisVar {
		var name pkg_rhino.Node = astFactory.CreateName_2_string_AstFactory_Type(("$jscomp$this" + "$" + j2gort.StrValString(this_.uniqueId)), this_.GetThisType())
		_ = name
		var thisVar pkg_rhino.Node = pkg_rhino.IR_ConstNode(name, astFactory.CreateThis(this_.GetThisType()))
		_ = thisVar
		NodeUtil_AddFeatureToScript(script, pkg_parser.StaticGet_FeatureSet_Feature_CONST_DECLARATIONS(), compiler)
		thisVar.SrcrefTreeIfMissing(this_.scopeBody)
		ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext_MakeTreeNonIndexable(thisVar)
		if this_.lastSuperStatement == 0 {
			ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext_InsertVarDeclaration(thisVar, this_.scopeBody)
		} else {
			thisVar.InsertAfter(this_.lastSuperStatement)
		}
		compiler.ReportChangeToEnclosingScope(thisVar)
	}
	if this_.needsArgumentsVar {
		var argumentsVar pkg_rhino.Node = astFactory.CreateArgumentsAliasDeclaration(("$jscomp$arguments" + "$" + j2gort.StrValString(this_.uniqueId)))
		_ = argumentsVar
		NodeUtil_AddFeatureToScript(script, pkg_parser.StaticGet_FeatureSet_Feature_CONST_DECLARATIONS(), compiler)
		ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext_InsertVarDeclaration(argumentsVar, this_.scopeBody)
		argumentsVar.SrcrefTreeIfMissing(this_.scopeBody)
		compiler.ReportChangeToEnclosingScope(argumentsVar)
	}
}

// Java static method: ThisAndArgumentsContext.insertVarDeclaration(com.google.javascript.rhino.Node varDeclaration, com.google.javascript.rhino.Node scopeBody)
func ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext_InsertVarDeclaration(varDeclaration pkg_rhino.Node, scopeBody pkg_rhino.Node) {
	_ = varDeclaration
	_ = scopeBody
	if scopeBody.GetParent().IsFunction() {
		var insertBeforePoint pkg_rhino.Node = NodeUtil_GetInsertionPointAfterAllInnerFunctionDeclarations(scopeBody)
		_ = insertBeforePoint
		if insertBeforePoint != 0 {
			varDeclaration.InsertBefore(insertBeforePoint)
		} else {
			scopeBody.AddChildToBack(varDeclaration)
		}
	} else {
		scopeBody.AddChildToFront(varDeclaration)
	}
}

// Java static method: ThisAndArgumentsContext.makeTreeNonIndexable(com.google.javascript.rhino.Node n)
func ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext_MakeTreeNonIndexable(n pkg_rhino.Node) {
	_ = n
	n.MakeNonIndexable()
	{
		var child pkg_rhino.Node = n.GetFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext_MakeTreeNonIndexable(child)
		}
	}
}
func (*ThisAndArgumentsReferenceUpdater) CanReleaseSyntacticScopes_() {}
