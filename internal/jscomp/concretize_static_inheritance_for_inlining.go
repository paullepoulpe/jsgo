// Port of com/google/javascript/jscomp/ConcretizeStaticInheritanceForInlining.java @ closure-compiler v20261021

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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
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
	duplicateClassNames map[string]bool

	inherits RuntimeJsLibManager_JsLibField

	compiler AbstractCompiler

	classByAlias    map[string]*ConcretizeStaticInheritanceForInlining_JavascriptClass
	classAliasLeafs map[string]struct{}
}

var ConcretizeStaticInheritanceForInlining_DUPLICATE_CLASS = DiagnosticType_Error("DUPLICATE_CLASS", "Multiple classes cannot share the same name: {0}")

// isBannedStaticInheritancePropName reports property names that may cause
// issues if they are concretized.
func isBannedStaticInheritancePropName(name string) bool {
	return name == "prototype" || name == "getInstance"
}

// ConcretizeStaticInheritanceForInlining_JavascriptClass is the nested class
// JavascriptClass.
type ConcretizeStaticInheritanceForInlining_JavascriptClass struct {
	// All static members to the class including get set properties, in insertion order.
	staticMembers   []rhino.Node
	staticMemberSet map[rhino.Node]bool
	// Keep updated the set of static member names to avoid O(n^2) searches.
	staticMemberNames map[string]bool
	// Collect all get set properties as defined by Object.defineProperties(...)
	definedProperties map[string]bool
}

// Java: ConcretizeStaticInheritanceForInlining.JavascriptClass()
func NewConcretizeStaticInheritanceForInlining_JavascriptClass() *ConcretizeStaticInheritanceForInlining_JavascriptClass {
	return &ConcretizeStaticInheritanceForInlining_JavascriptClass{}
}

// Java: ConcretizeStaticInheritanceForInlining.JavascriptClass.addStaticMember(Node)
func (c *ConcretizeStaticInheritanceForInlining_JavascriptClass) AddStaticMember(node rhino.Node) {
	if c.staticMemberSet == nil {
		c.staticMemberSet = make(map[rhino.Node]bool)
		c.staticMemberNames = make(map[string]bool)
	}
	if !c.staticMemberSet[node] {
		c.staticMemberSet[node] = true
		c.staticMembers = append(c.staticMembers, node)
	}
	c.staticMemberNames[node.GetFirstChild().GetString()] = true
}

func (c *ConcretizeStaticInheritanceForInlining_JavascriptClass) addDefinedProperty(name string) {
	if c.definedProperties == nil {
		c.definedProperties = make(map[string]bool)
	}
	c.definedProperties[name] = true
}

// Java: ConcretizeStaticInheritanceForInlining(AbstractCompiler)
func NewConcretizeStaticInheritanceForInlining(compiler AbstractCompiler) *ConcretizeStaticInheritanceForInlining {
	return &ConcretizeStaticInheritanceForInlining{
		compiler:        compiler,
		inherits:        compiler.GetRuntimeJsLibManager().GetJsLibField("$jscomp.inherits"),
		classByAlias:    make(map[string]*ConcretizeStaticInheritanceForInlining_JavascriptClass, 64),
		classAliasLeafs: make(map[string]struct{}, 64),
	}
}

func (c *ConcretizeStaticInheritanceForInlining) putClassByAlias(name string, cls *ConcretizeStaticInheritanceForInlining_JavascriptClass) {
	c.classByAlias[name] = cls
	leaf := name
	if dot := strings.LastIndexByte(name, '.'); dot >= 0 {
		leaf = name[dot+1:]
	}
	c.classAliasLeafs[leaf] = struct{}{}
}

func (c *ConcretizeStaticInheritanceForInlining) addDuplicateClassName(name string) {
	if c.duplicateClassNames == nil {
		c.duplicateClassNames = make(map[string]bool)
	}
	c.duplicateClassNames[name] = true
}

func (c *ConcretizeStaticInheritanceForInlining) lookupClassByNode(n rhino.Node) (string, *ConcretizeStaticInheritanceForInlining_JavascriptClass) {
	switch n.GetToken() {
	case rhino.Token_NAME:
		name := n.GetString()
		if name == "" {
			return "", nil
		}
		if _, ok := c.classAliasLeafs[name]; !ok {
			return "", nil
		}
		return name, c.classByAlias[name]
	case rhino.Token_GETPROP:
		if _, ok := c.classAliasLeafs[n.GetString()]; !ok {
			return "", nil
		}
		qname, ok := n.GetQualifiedNameString()
		if !ok {
			return "", nil
		}
		return qname, c.classByAlias[qname]
	case rhino.Token_THIS, rhino.Token_SUPER:
		qname, ok := n.GetQualifiedNameString()
		if !ok {
			return "", nil
		}
		return qname, c.classByAlias[qname]
	default:
		return "", nil
	}
}

func (c *ConcretizeStaticInheritanceForInlining) resolveReferencedClass(t *NodeTraversal, n rhino.Node) *ConcretizeStaticInheritanceForInlining_JavascriptClass {
	className, cls := c.lookupClassByNode(n)
	if cls == nil {
		return nil
	}
	if !n.IsName() {
		return cls
	}
	var_ := t.GetScope().GetVar(className)
	if var_ == nil || !var_.IsLocal() {
		return cls
	}
	return nil
}

// Java: ConcretizeStaticInheritanceForInlining.process(Node, Node)
func (c *ConcretizeStaticInheritanceForInlining) Process(externs rhino.Node, root rhino.Node) {
	findStaticMembers := NewConcretizeStaticInheritanceForInlining_FindStaticMembers(c)
	// Since this is an optimization pass, running on externs would be unnecessary
	NodeTraversal_Traverse(c.compiler, root, findStaticMembers)
	c.ProcessInherits(findStaticMembers)
}

// Java: ConcretizeStaticInheritanceForInlining.processInherits(FindStaticMembers)
func (c *ConcretizeStaticInheritanceForInlining) ProcessInherits(findStaticMembers *ConcretizeStaticInheritanceForInlining_FindStaticMembers) {
	for _, inheritsCall := range findStaticMembers.inheritsCalls {
		superclassNameNode := inheritsCall.GetLastChild()
		superclassQname, superOk := superclassNameNode.GetQualifiedNameString()
		subclassNameNode := superclassNameNode.GetPrevious()
		subclassQname, subOk := subclassNameNode.GetQualifiedNameString()
		var superClass, subClass *ConcretizeStaticInheritanceForInlining_JavascriptClass
		if superOk {
			superClass = c.classByAlias[superclassQname]
		}
		if subOk {
			subClass = c.classByAlias[subclassQname]
		}
		if superOk && c.duplicateClassNames[superclassQname] {
			c.compiler.Report(JSError_MakeWithNode(inheritsCall, ConcretizeStaticInheritanceForInlining_DUPLICATE_CLASS, superclassQname))
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
	inheritsCall rhino.Node,
	findStaticMembers *ConcretizeStaticInheritanceForInlining_FindStaticMembers) {
	for _, staticMember := range superClass.staticMembers {
		assignMember := ast.MustAssign(staticMember)
		memberName := assignMember.Target().Node().GetString()
		if superClass.definedProperties[memberName] {
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
		function := assignMember.Value().Node()
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
	// If this subclass overrides the static method or defines a getter/setter property,
	// there is no need to copy the method from the base class.
	return subClass.staticMemberNames[memberName] || subClass.definedProperties[memberName]
}

// Java: ConcretizeStaticInheritanceForInlining.isReferenceToClass(NodeTraversal, Node)
func (c *ConcretizeStaticInheritanceForInlining) IsReferenceToClass(t *NodeTraversal, n rhino.Node) bool {
	className, ok := n.GetQualifiedNameString()
	if !ok {
		return false
	}
	if _, exists := c.classByAlias[className]; !exists {
		return false
	}

	if !n.IsName() {
		return true
	}

	var_ := t.GetScope().GetVar(className)
	return var_ == nil || !var_.IsLocal()
}

// ConcretizeStaticInheritanceForInlining_FindStaticMembers is the inner
// class FindStaticMembers.
type ConcretizeStaticInheritanceForInlining_FindStaticMembers struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer *ConcretizeStaticInheritanceForInlining

	inheritsCalls []rhino.Node
	// Store the order we find class definitions and static fields.  Copied statics must occur
	// after both the namespace and the copied property are defined.
	nodeOrder map[rhino.Node]int
}

// Java: ConcretizeStaticInheritanceForInlining.FindStaticMembers()
func NewConcretizeStaticInheritanceForInlining_FindStaticMembers(outer *ConcretizeStaticInheritanceForInlining) *ConcretizeStaticInheritanceForInlining_FindStaticMembers {
	return &ConcretizeStaticInheritanceForInlining_FindStaticMembers{
		outer:     outer,
		nodeOrder: make(map[rhino.Node]int),
	}
}

// Java: ConcretizeStaticInheritanceForInlining.FindStaticMembers.visit(NodeTraversal, Node, Node)
func (f *ConcretizeStaticInheritanceForInlining_FindStaticMembers) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_CALL:
		call := ast.MustCall(n)
		if f.outer.inherits.Matches(call.Callee().Node()) {
			f.inheritsCalls = append(f.inheritsCalls, n)
			f.nodeOrder[n] = len(f.nodeOrder)
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
func (f *ConcretizeStaticInheritanceForInlining_FindStaticMembers) VisitDefinedPropertiesCall(t *NodeTraversal, definePropertiesCall rhino.Node) {
	call := ast.MustCall(definePropertiesCall)
	if firstArg, ok := call.FirstArg(); ok {
		if c := f.outer.resolveReferencedClass(t, firstArg.Node()); c != nil {
			for _, prop := range NodeUtil_GetObjectDefinedPropertiesKeys(definePropertiesCall) {
				c.addDefinedProperty(prop.GetString())
			}
		}
	}
}

// Java: ConcretizeStaticInheritanceForInlining.FindStaticMembers.visitFunctionClassDef(Node)
func (f *ConcretizeStaticInheritanceForInlining_FindStaticMembers) VisitFunctionClassDef(n rhino.Node) {
	classInfo := NodeUtil_GetBestJSDocInfo(n)
	if classInfo != nil && classInfo.IsConstructor() {
		name, ok := NodeUtil_GetNameString(n)
		if !ok {
			return
		}
		if _, exists := f.outer.classByAlias[name]; exists {
			f.outer.addDuplicateClassName(name)
		} else {
			f.outer.putClassByAlias(name, NewConcretizeStaticInheritanceForInlining_JavascriptClass())
		}
	}
}

// Java: ConcretizeStaticInheritanceForInlining.FindStaticMembers.setAlias(String, String)
func (f *ConcretizeStaticInheritanceForInlining_FindStaticMembers) SetAlias(original string, alias string) {
	existing, exists := f.outer.classByAlias[original]
	jutil.CheckArgument(exists)
	f.outer.putClassByAlias(alias, existing)
}

// Java: ConcretizeStaticInheritanceForInlining.FindStaticMembers.visitAssign(NodeTraversal, Node)
func (f *ConcretizeStaticInheritanceForInlining_FindStaticMembers) VisitAssign(t *NodeTraversal, n rhino.Node) {
	assign := ast.MustAssign(n)
	lhs := assign.Target().Node()
	rhs := assign.Value().Node()
	// Alias for classes. We assume that the alias appears after the class declaration.
	if _, existingClass := f.outer.lookupClassByNode(rhs); existingClass != nil {
		if alias, ok := lhs.GetQualifiedNameString(); ok {
			f.outer.putClassByAlias(alias, existingClass)
		}
	} else if lhs.IsGetProp() {
		classNode := lhs.GetFirstChild()
		if cls := f.outer.resolveReferencedClass(t, classNode); cls != nil && !isBannedStaticInheritancePropName(lhs.GetString()) {
			cls.AddStaticMember(n)
			f.nodeOrder[n] = len(f.nodeOrder)
		}
	}
}

// Java: ConcretizeStaticInheritanceForInlining.FindStaticMembers.visitVariableDeclaration(Node)
func (f *ConcretizeStaticInheritanceForInlining_FindStaticMembers) VisitVariableDeclaration(n rhino.Node) {
	decl := ast.MustNameDeclaration(n)
	firstDecl, ok := decl.FirstDeclarator()
	if !ok {
		return
	}
	child := firstDecl.Node()
	if !child.HasChildren() {
		return
	}
	if _, existingClass := f.outer.lookupClassByNode(child.GetFirstChild()); existingClass != nil {
		if maybeAlias, ok := child.GetQualifiedNameString(); ok {
			f.outer.putClassByAlias(maybeAlias, existingClass)
		}
	}
}

// Java: ConcretizeStaticInheritanceForInlining.FindStaticMembers.isBefore(Node, Node)
func (f *ConcretizeStaticInheritanceForInlining_FindStaticMembers) IsBefore(earlier rhino.Node, later rhino.Node) bool {
	earlierPosition, earlierOk := f.nodeOrder[earlier]
	laterPosition, laterOk := f.nodeOrder[later]
	return earlierOk && laterOk && earlierPosition < laterPosition
}
