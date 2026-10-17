// Port of com/google/javascript/jscomp/ConcretizeStaticInheritanceForInlining.java @ closure-compiler v20261015

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

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ConcretizeStaticInheritanceForInlining rewrites static inheritance to
// explicitly copy inherited properties from superclass to subclass so that
// the optimizer knows the subclass has those properties.
//
// For example, the main transpilation passes will convert this ES6 code:
//
//	class Foo { static f() {} }
//	class Bar extends Foo {}
//
// to this ES3 code:
//
//	function Foo() {}
//	Foo.f = function() {};
//	function Bar() {}
//	$jscomp.inherits(Bar, Foo);
//
// and then this class will convert that to
//
//	function Foo() {}
//	Foo.f = function() {};
//	function Bar() {}
//	$jscomp.inherits(Bar, Foo);
//	Bar.f = Foo.f;
//
// Additionally, there are getter and setter fields which are transpiled
// from:
//
//	class Foo { static get prop() { return 1; } }
//	class Bar extends Foo {}
//
// to:
//
//	var Foo = function() {};
//	Foo.prop; // stub declaration so that the optimizer knows about prop
//	Object.defineProperties(Foo, {prop:{get:function() { return 1; }}});
//
//	var Bar = function() {};
//	$jscomp.inherits(Bar, Foo);
//
// The stub declaration of Foo.prop needs to be duplicated for Bar so that
// the optimizer knows that Bar also has this property. (ES5 classes don't
// have class-side inheritance).
//
//	var Bar = function() {};
//	Bar.prop;
//	$jscomp.inherits(Bar, Foo);
//
// In order to gather the stub declarations, this pass gathers all GETPROPs
// on a class. In order to determine which of these are the stub
// declarations it filters them based on names discovered in
// Object.defineProperties. Unfortunately, we cannot simply gather the
// defined properties because they don't have the JSDoc, which may include
// optimization-relevant annotations like @nocollapse.
//
// TODO(tdeegan): In the future the JSDoc for getter/setter properties could
// be stored in the defineProperties functions. It would reduce the
// complexity of this pass significantly.
//
// NOTE: currently this pass only exists to prevent property collapsing from
// breaking some simple class-side inheritance cases when transpiling.
type ConcretizeStaticInheritanceForInlining struct {
	duplicateClassNames *jutil.LinkedHashSet[string]

	inherits RuntimeJsLibManager_JsLibField

	compiler AbstractCompiler

	classByAlias *jutil.LinkedHashMap[string, *ConcretizeStaticInheritanceForInlining_JavascriptClass]
}

var ConcretizeStaticInheritanceForInlining_DUPLICATE_CLASS = DiagnosticType_Error("DUPLICATE_CLASS", "Multiple classes cannot share the same name: {0}")

// Property names that may cause issues if they are concretized.
var ConcretizeStaticInheritanceForInlining_BANNED_PROP_NAMES = jutil.NewLinkedHashSet("prototype", "getInstance")

// ConcretizeStaticInheritanceForInlining_JavascriptClass is the nested class
// JavascriptClass.
type ConcretizeStaticInheritanceForInlining_JavascriptClass struct {
	// All static members to the class including get set properties.
	staticMembers *jutil.LinkedHashSet[*rhino.Node]
	// Keep updated the set of static member names to avoid O(n^2) searches.
	staticMemberNames *jutil.LinkedHashSet[string]
	// Collect all get set properties as defined by Object.defineProperties(...)
	definedProperties *jutil.LinkedHashSet[string]
}

// Java: ConcretizeStaticInheritanceForInlining.JavascriptClass()
func NewConcretizeStaticInheritanceForInlining_JavascriptClass() *ConcretizeStaticInheritanceForInlining_JavascriptClass {
	return &ConcretizeStaticInheritanceForInlining_JavascriptClass{
		staticMembers:     jutil.NewLinkedHashSet[*rhino.Node](),
		staticMemberNames: jutil.NewLinkedHashSet[string](),
		definedProperties: jutil.NewLinkedHashSet[string](),
	}
}

// Java: ConcretizeStaticInheritanceForInlining.JavascriptClass.addStaticMember(Node)
func (c *ConcretizeStaticInheritanceForInlining_JavascriptClass) AddStaticMember(node *rhino.Node) {
	c.staticMembers.Add(node)
	c.staticMemberNames.Add(node.GetFirstChild().GetString())
}

// Java: ConcretizeStaticInheritanceForInlining(AbstractCompiler)
func NewConcretizeStaticInheritanceForInlining(compiler AbstractCompiler) *ConcretizeStaticInheritanceForInlining {
	c := &ConcretizeStaticInheritanceForInlining{
		duplicateClassNames: jutil.NewLinkedHashSet[string](),
		classByAlias:        jutil.NewLinkedHashMap[string, *ConcretizeStaticInheritanceForInlining_JavascriptClass](),
	}
	c.compiler = compiler
	c.inherits = compiler.GetRuntimeJsLibManager().GetJsLibField("$jscomp.inherits")
	return c
}

// classByAliasGet is
// classByAlias.get(qname) for a nullable qualified name: a null key is never
// present.
//
// Go-only: language-adapter — LinkedHashMap.get(null).
func (c *ConcretizeStaticInheritanceForInlining) classByAliasGet(qname *string) *ConcretizeStaticInheritanceForInlining_JavascriptClass {
	if qname == nil {
		return nil
	}
	return c.classByAlias.Get(*qname)
}

// classByAliasContainsKey is classByAlias.containsKey(qname) for a nullable
// qualified name.
//
// Go-only: language-adapter — LinkedHashMap.containsKey(null).
func (c *ConcretizeStaticInheritanceForInlining) classByAliasContainsKey(qname *string) bool {
	return qname != nil && c.classByAlias.ContainsKey(*qname)
}

// Java: ConcretizeStaticInheritanceForInlining.process(Node, Node)
func (c *ConcretizeStaticInheritanceForInlining) Process(externs *rhino.Node, root *rhino.Node) {
	findStaticMembers := NewConcretizeStaticInheritanceForInlining_FindStaticMembers(c)
	// Since this is an optimization pass, running on externs would be unnecessary
	NodeTraversal_Traverse(c.compiler, root, findStaticMembers)
	c.ProcessInherits(findStaticMembers)
}

// Java: ConcretizeStaticInheritanceForInlining.processInherits(FindStaticMembers)
func (c *ConcretizeStaticInheritanceForInlining) ProcessInherits(findStaticMembers *ConcretizeStaticInheritanceForInlining_FindStaticMembers) {
	for _, inheritsCall := range findStaticMembers.inheritsCalls.Slice() {
		superclassNameNode := inheritsCall.GetLastChild()
		superclassQname := superclassNameNode.GetQualifiedName()
		subclassNameNode := superclassNameNode.GetPrevious()
		subclassQname := subclassNameNode.GetQualifiedName()
		superClass := c.classByAliasGet(superclassQname)
		subClass := c.classByAliasGet(subclassQname)
		if superclassQname != nil && c.duplicateClassNames.Contains(*superclassQname) {
			c.compiler.Report(JSError_MakeWithNode(inheritsCall, ConcretizeStaticInheritanceForInlining_DUPLICATE_CLASS, *superclassQname))
			return
		}
		if superClass == nil || subClass == nil {
			continue
		}
		c.CopyStaticMembers(superClass, subClass, inheritsCall, findStaticMembers)
	}
}

// Java: ConcretizeStaticInheritanceForInlining.copyStaticMembers(JavascriptClass, JavascriptClass, Node, FindStaticMembers)
func (c *ConcretizeStaticInheritanceForInlining) CopyStaticMembers(
	superClass *ConcretizeStaticInheritanceForInlining_JavascriptClass,
	subClass *ConcretizeStaticInheritanceForInlining_JavascriptClass,
	inheritsCall *rhino.Node,
	findStaticMembers *ConcretizeStaticInheritanceForInlining_FindStaticMembers) {
	for staticMember := range superClass.staticMembers.All() {
		jutil.CheckStateMsg1(staticMember.IsAssign(), "%s", staticMember)
		memberName := staticMember.GetFirstChild().GetString()
		if superClass.definedProperties.Contains(memberName) {
			continue
		}
		if c.IsOverriden(subClass, memberName) {
			continue
		}
		if findStaticMembers.IsBefore(inheritsCall, staticMember) {
			// Don't copy members that are defined after the $jscomp.inherits call,
			// since they will not work correctly in IE<11, where static inheritance
			// is done by copying, rather than prototype manipulation.
			continue
		}

		info := rhino.JSDocInfo_Builder_MaybeCopyFrom(staticMember.GetJSDocInfo())
		function := staticMember.GetLastChild()
		sourceInfoNode := function
		if function.IsFunction() {
			sourceInfoNode = function.GetFirstChild()
		}

		subclassNameNode := inheritsCall.GetSecondChild()
		superclassNameNode := subclassNameNode.GetNext()
		assign := rhino.IR_Assign(
			rhino.IR_Getprop(subclassNameNode.CloneTree(), memberName),
			rhino.IR_Getprop(superclassNameNode.CloneTree(), memberName))
		assign.SetJSDocInfo(info.Build())
		exprResult := rhino.IR_ExprResult(assign)
		exprResult.SrcrefTreeIfMissing(sourceInfoNode)
		inheritsExpressionResult := inheritsCall.GetParent()
		exprResult.InsertAfter(inheritsExpressionResult)
		c.compiler.ReportChangeToEnclosingScope(inheritsExpressionResult)

		// Add the static member to the subclass so that subclasses also copy this member.
		subClass.AddStaticMember(assign)
	}
}

// Java: ConcretizeStaticInheritanceForInlining.isOverriden(JavascriptClass, String)
func (c *ConcretizeStaticInheritanceForInlining) IsOverriden(subClass *ConcretizeStaticInheritanceForInlining_JavascriptClass, memberName string) bool {
	if subClass.staticMemberNames.Contains(memberName) {
		// This subclass overrides the static method, so there is no need to copy the
		// method from the base class.
		return true
	}
	if subClass.definedProperties.Contains(memberName) {
		return true
	}
	return false
}

// Java: ConcretizeStaticInheritanceForInlining.isReferenceToClass(NodeTraversal, Node)
func (c *ConcretizeStaticInheritanceForInlining) IsReferenceToClass(t *NodeTraversal, n *rhino.Node) bool {
	className := n.GetQualifiedName()
	if !c.classByAliasContainsKey(className) {
		return false
	}

	if !n.IsName() {
		return true
	}

	var_ := t.GetScope().GetVar(*className)
	return var_ == nil || !var_.IsLocal()
}

// ConcretizeStaticInheritanceForInlining_FindStaticMembers is the inner
// class FindStaticMembers.
type ConcretizeStaticInheritanceForInlining_FindStaticMembers struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer *ConcretizeStaticInheritanceForInlining

	inheritsCalls *jutil.ArrayList[*rhino.Node]
	// Store the order we find class definitions and static fields.  Copied statics must occur
	// after both the namespace and the copied property are defined.
	nodeOrder *jutil.LinkedHashMap[*rhino.Node, int]
}

// Java: ConcretizeStaticInheritanceForInlining.FindStaticMembers()
func NewConcretizeStaticInheritanceForInlining_FindStaticMembers(outer *ConcretizeStaticInheritanceForInlining) *ConcretizeStaticInheritanceForInlining_FindStaticMembers {
	return &ConcretizeStaticInheritanceForInlining_FindStaticMembers{
		outer:         outer,
		inheritsCalls: jutil.NewArrayList[*rhino.Node](),
		nodeOrder:     jutil.NewLinkedHashMap[*rhino.Node, int](),
	}
}

// Java: ConcretizeStaticInheritanceForInlining.FindStaticMembers.visit(NodeTraversal, Node, Node)
func (f *ConcretizeStaticInheritanceForInlining_FindStaticMembers) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_CALL:
		if f.outer.inherits.Matches(n.GetFirstChild()) {
			f.inheritsCalls.Add(n)
			f.nodeOrder.Put(n, f.nodeOrder.Size())
		}
		if NodeUtil_IsObjectDefinePropertiesDefinition(n) {
			f.VisitDefinedPropertiesCall(t, n)
		}
	case rhino.Token_CONST, rhino.Token_LET, rhino.Token_VAR:
		f.VisitVariableDeclaration(n)
	case rhino.Token_ASSIGN:
		f.VisitAssign(t, n)
	case rhino.Token_FUNCTION:
		f.VisitFunctionClassDef(n)
	default:
	}
}

// Java: ConcretizeStaticInheritanceForInlining.FindStaticMembers.visitDefinedPropertiesCall(NodeTraversal, Node)
func (f *ConcretizeStaticInheritanceForInlining_FindStaticMembers) VisitDefinedPropertiesCall(t *NodeTraversal, definePropertiesCall *rhino.Node) {
	object := definePropertiesCall.GetSecondChild()
	if f.outer.IsReferenceToClass(t, object) {
		className := object.GetQualifiedName()
		c := f.outer.classByAliasGet(className)
		for _, prop := range NodeUtil_GetObjectDefinedPropertiesKeys(definePropertiesCall) {
			c.definedProperties.Add(prop.GetString())
		}
	}
}

// Java: ConcretizeStaticInheritanceForInlining.FindStaticMembers.visitFunctionClassDef(Node)
func (f *ConcretizeStaticInheritanceForInlining_FindStaticMembers) VisitFunctionClassDef(n *rhino.Node) {
	classInfo := NodeUtil_GetBestJSDocInfo(n)
	if classInfo != nil && classInfo.IsConstructor() {
		name := NodeUtil_GetName(n)
		if name == nil {
			return
		}
		if f.outer.classByAlias.ContainsKey(*name) {
			f.outer.duplicateClassNames.Add(*name)
		} else {
			f.outer.classByAlias.Put(*name, NewConcretizeStaticInheritanceForInlining_JavascriptClass())
		}
	}
}

// Java: ConcretizeStaticInheritanceForInlining.FindStaticMembers.setAlias(String, String)
func (f *ConcretizeStaticInheritanceForInlining_FindStaticMembers) SetAlias(original *string, alias *string) {
	jutil.CheckNotNullMsg(original, "original is null")
	jutil.CheckNotNullMsg(alias, "alias is null")
	jutil.CheckArgument(f.outer.classByAlias.ContainsKey(*original))
	f.outer.classByAlias.Put(*alias, f.outer.classByAlias.Get(*original))
}

// Java: ConcretizeStaticInheritanceForInlining.FindStaticMembers.visitAssign(NodeTraversal, Node)
func (f *ConcretizeStaticInheritanceForInlining_FindStaticMembers) VisitAssign(t *NodeTraversal, n *rhino.Node) {
	// Alias for classes. We assume that the alias appears after the class declaration.
	existingClassQname := n.GetLastChild().GetQualifiedName()
	if existingClassQname != nil && f.outer.classByAlias.ContainsKey(*existingClassQname) {
		alias := n.GetFirstChild().GetQualifiedName()
		if alias != nil {
			f.SetAlias(existingClassQname, alias)
		}
	} else if n.GetFirstChild().IsGetProp() {
		getProp := n.GetFirstChild()
		classNode := getProp.GetFirstChild()
		if f.outer.IsReferenceToClass(t, classNode) && !ConcretizeStaticInheritanceForInlining_BANNED_PROP_NAMES.Contains(getProp.GetString()) {
			f.outer.classByAliasGet(classNode.GetQualifiedName()).AddStaticMember(n)
			f.nodeOrder.Put(n, f.nodeOrder.Size())
		}
	}
}

// Java: ConcretizeStaticInheritanceForInlining.FindStaticMembers.visitVariableDeclaration(Node)
func (f *ConcretizeStaticInheritanceForInlining_FindStaticMembers) VisitVariableDeclaration(n *rhino.Node) {
	child := n.GetFirstChild()
	if !child.HasChildren() {
		return
	}
	maybeOriginalName := child.GetFirstChild().GetQualifiedName()
	if maybeOriginalName != nil && f.outer.classByAlias.ContainsKey(*maybeOriginalName) {
		maybeAlias := child.GetQualifiedName()
		if maybeAlias != nil {
			f.SetAlias(maybeOriginalName, maybeAlias)
		}
	}
}

// Java: ConcretizeStaticInheritanceForInlining.FindStaticMembers.isBefore(Node, Node)
func (f *ConcretizeStaticInheritanceForInlining_FindStaticMembers) IsBefore(earlier *rhino.Node, later *rhino.Node) bool {
	earlierPosition, earlierOk := f.nodeOrder.Lookup(earlier)
	laterPosition, laterOk := f.nodeOrder.Lookup(later)
	return earlierOk && laterOk && earlierPosition < laterPosition
}
