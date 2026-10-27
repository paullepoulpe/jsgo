// Port of com/google/javascript/jscomp/J2clPropertyInlinerPass.java @ closure-compiler v20261021

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

package jscomp

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// J2clPropertyInlinerPass targets J2CL output. It looks for static get and
// set methods defined within a class that match the signature of J2CL static
// fields and inlines them at their call sites. This is done for performance
// reasons since getter and setter accesses are slower than regular field
// accesses.
//
// This will be done by looking at all property accesses and determining if
// they have a corresponding get or set method on the property qualifiers
// definition. Some caveats:
//
//   - Avoid inlining if the property is set using compound assignments.
//   - Avoid inlining if the property is incremented using ++ or --
//
// Since this pass only really works after the AggressiveInlineAliases pass
// has run, we have to look for Object.defineProperties instead of es6 get and
// set nodes since es6 transpilation has already occurred if the language out
// is ES5.
//
// Java: public class J2clPropertyInlinerPass implements CompilerPass
type J2clPropertyInlinerPass struct {
	compiler AbstractCompiler
	inliner  *j2clStaticFieldInliner
}

// Java: J2clPropertyInlinerPass(AbstractCompiler)
func NewJ2clPropertyInlinerPass(compiler AbstractCompiler) *J2clPropertyInlinerPass {
	return &J2clPropertyInlinerPass{compiler: compiler}
}

// Java: J2clPropertyInlinerPass.process(Node, Node)
func (p *J2clPropertyInlinerPass) Process(externs rhino.Node, root rhino.Node) {
	if !J2clSourceFileChecker_ShouldRunJ2clPasses(p.compiler) {
		return
	}

	p.getOrCreateInliner(root).run()
	// This pass may remove getters and setters.
	GatherGetterAndSetterProperties_Update(p.compiler, externs, root)
}

// Java: class J2clPropertyInlinerPass.StaticFieldGetterSetterInliner
type j2clStaticFieldInliner struct {
	outer      *J2clPropertyInlinerPass
	root       rhino.Node
	gatherer   j2clGatherer
	determiner j2clDeterminer
	inliner    j2clInliner
}

func (p *J2clPropertyInlinerPass) getOrCreateInliner(root rhino.Node) *j2clStaticFieldInliner {
	s := p.inliner
	if s == nil {
		s = &j2clStaticFieldInliner{outer: p}
		s.gatherer.outer = s
		s.gatherer.j2clPropertiesByName = jutil.NewLinkedHashMap[string, *j2clProperty]()
		s.determiner.propertiesByName = s.gatherer.j2clPropertiesByName
		s.inliner.outer = s
		s.inliner.propertiesByName = s.gatherer.j2clPropertiesByName
		p.inliner = s
	} else {
		s.gatherer.j2clPropertiesByName.Clear()
	}
	s.root = root
	return s
}

// Java: J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.run()
func (s *j2clStaticFieldInliner) run() {
	NodeTraversal_Traverse(s.outer.compiler, s.root, &s.gatherer)
	if s.gatherer.j2clPropertiesByName.IsEmpty() {
		return
	}
	NodeTraversal_Traverse(s.outer.compiler, s.root, &s.determiner)
	s.inliner.run()
}

// j2clProperty represents an ES5 (Object.defineProperties) or ES6 (static
// getter/setter) J2CL static property candidate.
//
// Java: J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.J2clProperty / J2clPropertyEs5 / J2clPropertyEs6
type j2clProperty struct {
	getKey         ast.PlainKeyMember
	setKey         ast.PlainKeyMember
	isSafeToInline bool
	isEs6          bool
}

func newJ2clPropertyEs5(getKey ast.StringKey, setKey ast.StringKey) *j2clProperty {
	jutil.CheckArgumentMsg1(getKey.Key() == "get", "%s", getKey.Node())
	jutil.CheckArgumentMsg1(setKey.Node().IsNil() || setKey.Key() == "set", "%s", setKey.Node())
	return &j2clProperty{
		getKey:         getKey.AsPlainKeyMember(),
		setKey:         setKey.AsPlainKeyMember(),
		isSafeToInline: true,
	}
}

func newJ2clPropertyEs6(getKey ast.GetterDef, setKey ast.SetterDef) *j2clProperty {
	return &j2clProperty{
		getKey:         getKey.AsPlainKeyMember(),
		setKey:         setKey.AsPlainKeyMember(),
		isSafeToInline: true,
		isEs6:          true,
	}
}

func (p *j2clProperty) remove(compiler AbstractCompiler) {
	if p.isEs6 {
		classMembers := p.getKey.Node().GetParent()
		jutil.CheckStateMsg1(classMembers.IsClassMembers(), "%s", classMembers)
		p.getKey.Detach()
		NodeUtil_MarkFunctionsDeleted(p.getKey.Node(), compiler)
		if !p.setKey.Node().IsNil() {
			p.setKey.Detach()
			NodeUtil_MarkFunctionsDeleted(p.setKey.Node(), compiler)
		}
		compiler.ReportChangeToEnclosingScope(classMembers)
		return
	}

	nodeToDetach := p.getKey.Node().GetGrandparent()
	objectLit := nodeToDetach.GetParent()
	jutil.CheckStateMsg1(objectLit.IsObjectLit(), "%s", objectLit)
	nodeToDetach.Detach()
	NodeUtil_MarkFunctionsDeleted(nodeToDetach, compiler)
	compiler.ReportChangeToEnclosingScope(objectLit)
	if !objectLit.HasChildren() {
		// Remove the whole Object.defineProperties call if there are no properties left.
		objectLit.GetGrandparent().Detach()
	}
}

// matchesJ2clGetKeySignature matches J2CL property getters by looking for the
// following signature:
//
//	get: function() { return (ClassName$$0clinit(), ClassName$$0fieldName)};
//	// OR this if CollapseProperties decided not to collapse these names
//	get: function() { return (ClassName.$clinit(), ClassName.$fieldName)};
//
// className is a nullable Java String: it is only dereferenced where Java
// dereferences it.
//
// Java: J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.matchesJ2clGetKeySignature(String, Node)
func matchesJ2clGetKeySignature(className *string, getKey ast.PlainKeyMember) bool {
	getFunction, ok := getKey.Value().AsFunction()
	if !ok {
		return false
	}
	getBlock, ok := getFunction.Body().AsBlock()
	if !ok {
		return false
	}
	// `{ return ClassName$$0clinit(), ClassName$prop; }`
	stmt, ok := getBlock.OnlyStatement()
	if !ok {
		return false
	}
	returnStmt, ok := stmt.AsReturn()
	if !ok {
		return false
	}
	// `ClassName$$0clinit(), ClassName$prop;`
	commaExpr, ok := ast.AsBinaryOf(returnStmt.Node().GetOnlyChild(), rhino.Token_COMMA)
	if !ok {
		return false
	}
	// `ClassName$$0clinit()`
	call, ok := commaExpr.Left().AsCall()
	if !ok || call.HasArgs() {
		return false
	}
	clinitQname, ok := call.Callee().Node().GetQualifiedNameString()
	if !ok || !looksLikeJ2clClinitName(className, clinitQname) {
		return false
	}
	// `ClassName$prop`
	propQname, ok := commaExpr.Right().Node().GetQualifiedNameString()
	return ok && looksLikeJ2clStaticPropName(className, propQname)
}

// matchesJ2clSetKeySignature matches J2CL property setters by looking for
// the following signature:
//
//	set: function(value) { (ClassName$$0clinit(), ClassName$$0fieldName = value)};
//	// OR this if CollapseProperties decided not to collapse the names.
//	set: function(value) { (ClassName.$clinit(), ClassName.$fieldNaNme = value)};
//
// Java: J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.matchesJ2clSetKeySignature(String, Node)
func matchesJ2clSetKeySignature(className *string, setKey ast.PlainKeyMember) bool {
	setFunction, ok := setKey.Value().AsFunction()
	if !ok {
		return false
	}
	setBlock, ok := setFunction.Body().AsBlock()
	if !ok {
		return false
	}
	if _, ok := setFunction.Params().OnlyParam(); !ok {
		// There is a single parameter "value".
		return false
	}
	// `{ ClassName$$0clinit(), ClassName$$0staticProp = initialValue; }`
	stmt, ok := setBlock.OnlyStatement()
	if !ok {
		return false
	}
	exprStmt, ok := stmt.AsExprStmt()
	if !ok {
		return false
	}
	// `ClassName$$0clinit(), ClassName$$0staticProp = initialValue`
	commaExpr, ok := exprStmt.Expr().AsBinaryOf(rhino.Token_COMMA)
	if !ok {
		return false
	}
	// `SomeClassName$$0clinit()`
	clinitCall, ok := commaExpr.Left().AsCall()
	if !ok || clinitCall.HasArgs() {
		return false
	}
	clinitQname, ok := clinitCall.Callee().Node().GetQualifiedNameString()
	if !ok || !looksLikeJ2clClinitName(className, clinitQname) {
		return false
	}
	// `ClassName$$0staticProp = initialValue`
	// Rewriting by CollapseProperties may have changed the LHS
	// from a GETPROP to a NAME.
	assign, ok := commaExpr.Right().AsAssign()
	if !ok {
		return false
	}
	lhsQname, ok := assign.Target().Node().GetQualifiedNameString()
	return ok && looksLikeJ2clStaticPropName(className, lhsQname)
}

// Java: J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.looksLikeAJ2clGeneratedClassClinitName(String, String)
func looksLikeJ2clClinitName(className *string, candidate string) bool {
	if !strings.HasPrefix(candidate, *className) {
		return false
	}
	suffix := candidate[len(*className):]
	return suffix == "$$0clinit" || // collapsed by CollapseProperties
		suffix == ".$clinit" // not collapsed
}

// Java: J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.looksLikeAJ2clGeneratedStaticPropertyName(String, String)
func looksLikeJ2clStaticPropName(className *string, candidate string) bool {
	if !strings.HasPrefix(candidate, *className) {
		return false
	}
	suffix := candidate[len(*className):]
	return strings.HasPrefix(suffix, "$$0") || // collapsed by CollapseProperties
		strings.HasPrefix(suffix, ".$") // not collapsed
}

// j2clGatherer traverses the AST and gathers static J2CL getters and setters
// from classes and Object.defineProperties calls.
//
// Java: private class J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.GatherJ2CLClassGetterSetters
type j2clGatherer struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer                *j2clStaticFieldInliner
	j2clPropertiesByName *jutil.LinkedHashMap[string, *j2clProperty]
	setterDefByName      *jutil.LinkedHashMap[string, ast.SetterDef]
	getterDefByName      *jutil.LinkedHashMap[string, ast.GetterDef]
}

// Java: J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.GatherJ2CLClassGetterSetters.visit(NodeTraversal, Node, Node)
func (g *j2clGatherer) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if class, ok := ast.AsClass(n); ok {
		g.visitClass(class)
	} else if NodeUtil_IsObjectDefinePropertiesDefinition(n) {
		g.visitObjectDefineProperties(ast.MustCall(n))
	}
}

// Java: J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.GatherJ2CLClassGetterSetters.visitObjectDefineProperties(Node)
func (g *j2clGatherer) visitObjectDefineProperties(call ast.Call) {
	firstArg, ok := call.FirstArg()
	if !ok {
		return
	}
	className, ok := firstArg.AsName()
	if !ok {
		return
	}
	classNameString := className.Ident()
	for _, p := range NodeUtil_GetObjectDefinedPropertiesKeys(call.Node()) {
		propKey := ast.MustStringKey(p)
		propLit, ok := propKey.Value().AsObjectLit()
		if !ok {
			continue
		}
		var getKey ast.StringKey
		var setKey ast.StringKey
		hasSetter := false
		for member := range propLit.Members() {
			innerKey, ok := member.AsStringKey()
			if !ok {
				continue
			}
			switch innerKey.Key() {
			case "get":
				if matchesJ2clGetKeySignature(&classNameString, innerKey.AsPlainKeyMember()) {
					getKey = innerKey
				}
			case "set":
				hasSetter = true
				if matchesJ2clSetKeySignature(&classNameString, innerKey.AsPlainKeyMember()) {
					setKey = innerKey
				}
			}
		}
		if !getKey.Node().IsNil() && (!hasSetter || !setKey.Node().IsNil()) {
			g.j2clPropertiesByName.Put(
				classNameString+"."+propKey.Key(), newJ2clPropertyEs5(getKey, setKey))
		}
	}
}

// Java: J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.GatherJ2CLClassGetterSetters.visitClass(Node)
func (g *j2clGatherer) visitClass(class ast.Class) {
	hasStaticGetter := false
	for member := range class.Body().Members() {
		if getter, ok := member.AsGetterDef(); ok && getter.IsStatic() {
			hasStaticGetter = true
			break
		}
	}
	if !hasStaticGetter {
		return
	}

	className := NodeUtil_GetName(class.Node())
	if g.getterDefByName == nil {
		g.getterDefByName = jutil.NewLinkedHashMap[string, ast.GetterDef]()
		g.setterDefByName = jutil.NewLinkedHashMap[string, ast.SetterDef]()
	} else {
		g.getterDefByName.Clear()
		g.setterDefByName.Clear()
	}

	for member := range class.Body().Members() {
		if getter, ok := member.AsGetterDef(); ok && getter.IsStatic() {
			g.getterDefByName.Put(getter.Key(), getter)
		} else if setter, ok := member.AsSetterDef(); ok && setter.IsStatic() {
			g.setterDefByName.Put(setter.Key(), setter)
		}
	}

	for propertyName := range g.getterDefByName.Keys() {
		getterDef := g.getterDefByName.Get(propertyName)
		setterDef := g.setterDefByName.Get(propertyName)
		if matchesJ2clGetKeySignature(className, getterDef.AsPlainKeyMember()) &&
			(setterDef.Node().IsNil() || matchesJ2clSetKeySignature(className, setterDef.AsPlainKeyMember())) {
			g.j2clPropertiesByName.Put(
				*className+"."+propertyName, newJ2clPropertyEs6(getterDef, setterDef))
		}
	}
}

// Java: private static class J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.DetermineInlinableProperties
type j2clDeterminer struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	propertiesByName *jutil.LinkedHashMap[string, *j2clProperty]
}

// Java: J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.DetermineInlinableProperties.visit(NodeTraversal, Node, Node)
func (d *j2clDeterminer) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	var target ast.Target
	if assignOp, ok := ast.AsAssignOp(n); ok {
		target = assignOp.Target()
	} else if incDec, ok := ast.AsIncDec(n); ok {
		target = incDec.Target()
	} else {
		return
	}
	if getProp, ok := target.AsGetProp(); ok {
		if accessName, ok := getProp.Node().GetQualifiedNameString(); ok {
			if prop := d.propertiesByName.Get(accessName); prop != nil {
				prop.isSafeToInline = false
			}
		}
	}
}

// j2clInliner looks for accesses and assignments of J2CL static properties and
// inlines them.
//
// Java: private class J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.InlinePropertiesPass
type j2clInliner struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer            *j2clStaticFieldInliner
	propertiesByName *jutil.LinkedHashMap[string, *j2clProperty]
	injector         *FunctionInjector
}

func (i *j2clInliner) getInjector() *FunctionInjector {
	if i.injector == nil {
		i.injector = NewFunctionInjector_Builder(i.outer.outer.compiler).
			AssumeStrictThis(true).
			Build()
	}
	return i.injector
}

// Java: J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.InlinePropertiesPass.run()
func (i *j2clInliner) run() {
	compiler := i.outer.outer.compiler
	NodeTraversal_Traverse(compiler, i.outer.root, i)

	for prop := range i.propertiesByName.Values() {
		if prop.isSafeToInline {
			prop.remove(compiler)
		}
	}
}

// Java: J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.InlinePropertiesPass.visit(NodeTraversal, Node, Node)
func (i *j2clInliner) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if getProp, ok := ast.AsGetProp(n); ok {
		if parent.IsExprResult() {
			// This is a stub declaration for the type checker. See:
			// ConcretizeStaticInheritanceForInlining
			return
		}
		if assign, ok := ast.AsAssignment(parent); ok && assign.Target().Is(n) {
			// This case should be handled below. It needs to be inlined differently.
			return
		}
		if accessName, ok := getProp.Node().GetQualifiedNameString(); ok {
			if prop := i.propertiesByName.Get(accessName); prop != nil && prop.isSafeToInline {
				injector := i.getInjector()
				functionCall := rhino.IR_Call(rhino.IR_Name("inlined_j2cl_getter"))
				n.ReplaceWith(functionCall)
				reference := NewFunctionInjector_Reference(functionCall, t.GetScope(), t.GetChunk(), FunctionInjector_InliningMode_DIRECT)
				inlinedCall := injector.Inline(reference, "", prop.getKey.Value().Node())
				t.GetCompiler().ReportChangeToEnclosingScope(inlinedCall)
			}
		}
		return
	}

	if assign, ok := ast.AsAssign(n); ok {
		if getProp, ok := assign.Target().AsGetProp(); ok {
			if accessName, ok := getProp.Node().GetQualifiedNameString(); ok {
				prop := i.propertiesByName.Get(accessName)
				if prop != nil && !prop.setKey.Node().IsNil() && prop.isSafeToInline {
					injector := i.getInjector()
					assignmentValue := assign.Value().Detach().Node()
					functionCall := rhino.IR_Call(rhino.IR_Name("inlined_j2cl_setter"), assignmentValue)
					n.ReplaceWith(functionCall)
					reference := NewFunctionInjector_Reference(functionCall, t.GetScope(), t.GetChunk(), FunctionInjector_InliningMode_BLOCK)
					injector.MaybePrepareCall(reference)
					inlinedCall := injector.Inline(reference, "", prop.setKey.Value().Node())
					t.GetCompiler().ReportChangeToEnclosingScope(inlinedCall)
				}
			}
		}
	}
}
