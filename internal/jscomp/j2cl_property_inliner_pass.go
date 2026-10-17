// Port of com/google/javascript/jscomp/J2clPropertyInlinerPass.java @ closure-compiler v20260415

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
}

// Java: J2clPropertyInlinerPass(AbstractCompiler)
func NewJ2clPropertyInlinerPass(compiler AbstractCompiler) *J2clPropertyInlinerPass {
	return &J2clPropertyInlinerPass{compiler: compiler}
}

// Java: J2clPropertyInlinerPass.process(Node, Node)
func (p *J2clPropertyInlinerPass) Process(externs *rhino.Node, root *rhino.Node) {
	if !J2clSourceFileChecker_ShouldRunJ2clPasses(p.compiler) {
		return
	}

	p.newStaticFieldGetterSetterInliner(root).run()
	// This pass may remove getters and setters.
	GatherGetterAndSetterProperties_Update(p.compiler, externs, root)
}

// Java: class J2clPropertyInlinerPass.StaticFieldGetterSetterInliner
type J2clPropertyInlinerPass_StaticFieldGetterSetterInliner struct {
	outer *J2clPropertyInlinerPass
	root  *rhino.Node
}

// Java: J2clPropertyInlinerPass.StaticFieldGetterSetterInliner(Node)
func (p *J2clPropertyInlinerPass) newStaticFieldGetterSetterInliner(root *rhino.Node) *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner {
	return &J2clPropertyInlinerPass_StaticFieldGetterSetterInliner{outer: p, root: root}
}

// Java: J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.run()
func (s *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner) run() {
	gatherer := s.newGatherJ2CLClassGetterSetters()
	NodeTraversal_Traverse(s.outer.compiler, s.root, gatherer)
	result := gatherer.getResults()
	NodeTraversal_Traverse(s.outer.compiler, s.root, NewJ2clPropertyInlinerPass_StaticFieldGetterSetterInliner_DetermineInlinableProperties(result))
	s.newInlinePropertiesPass(result).run()
}

// J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clProperty is the
// abstract J2clProperty.
//
// Java: private abstract static class J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.J2clProperty
type J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clProperty interface {
	// Java: J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.J2clProperty.remove()
	Remove()

	GetJ2clPropertyBase() *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyBase
}

// J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyBase
// holds J2clProperty's fields.
type J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyBase struct {
	getKey         *rhino.Node
	setKey         *rhino.Node
	isSafeToInline bool
}

// Java: J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.J2clProperty(Node, Node)
func (b *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyBase) InitJ2clPropertyBase(getKey *rhino.Node, setKey *rhino.Node) {
	b.getKey = getKey
	b.setKey = setKey
	b.isSafeToInline = true
}

// GetJ2clPropertyBase gives access to the base fields through the interface.
func (b *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyBase) GetJ2clPropertyBase() *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyBase {
	return b
}

// J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs5 is
// a J2CL property with a getter and setter from an Object.defineProperties
// call
//
// Java: private class J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.J2clPropertyEs5 extends J2clProperty
type J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs5 struct {
	J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyBase
	outer *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner
}

// Java: J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.J2clPropertyEs5(Node, Node)
func (s *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner) newJ2clPropertyEs5(getKey *rhino.Node, setKey *rhino.Node) *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs5 {
	p := &J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs5{outer: s}
	p.InitJ2clPropertyBase(getKey, setKey)
	jutil.CheckArgumentMsg1(getKey.IsStringKey() && getKey.GetString() == "get", "%s", getKey)
	jutil.CheckArgumentMsg1(
		setKey == nil || (setKey.IsStringKey() && setKey.GetString() == "set"), "%s", setKey)
	return p
}

// Java: J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.J2clPropertyEs5.remove()
func (p *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs5) Remove() {
	compiler := p.outer.outer.compiler
	nodeToDetach := p.getKey.GetGrandparent()
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

// J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs6 is
// a J2CL property created with a ES6-style static getter and setter
//
// Java: private class J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.J2clPropertyEs6 extends J2clProperty
type J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs6 struct {
	J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyBase
	outer *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner
}

// Java: J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.J2clPropertyEs6(Node, Node)
func (s *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner) newJ2clPropertyEs6(getKey *rhino.Node, setKey *rhino.Node) *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs6 {
	p := &J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs6{outer: s}
	p.InitJ2clPropertyBase(getKey, setKey)
	jutil.CheckArgumentMsg1(getKey.IsGetterDef(), "%s", getKey)
	jutil.CheckArgumentMsg1(setKey == nil || setKey.IsSetterDef(), "%s", setKey)
	return p
}

// Java: J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.J2clPropertyEs6.remove()
func (p *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs6) Remove() {
	compiler := p.outer.outer.compiler
	classMembers := p.getKey.GetParent()
	jutil.CheckStateMsg1(classMembers.IsClassMembers(), "%s", classMembers)
	p.getKey.Detach()
	NodeUtil_MarkFunctionsDeleted(p.getKey, compiler)
	if p.setKey != nil {
		p.setKey.Detach()
		NodeUtil_MarkFunctionsDeleted(p.setKey, compiler)
	}
	compiler.ReportChangeToEnclosingScope(classMembers)
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
func (s *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner) matchesJ2clGetKeySignature(className *string, getKey *rhino.Node) bool {
	if !getKey.GetFirstChild().IsFunction() {
		return false
	}
	getFunction := getKey.GetFirstChild()
	if !getFunction.HasChildren() || !getFunction.GetLastChild().IsBlock() {
		return false
	}
	// `{ return ClassName$$0clinit(), ClassName$prop; }`
	getBlock := getFunction.GetLastChild()
	if !getBlock.HasOneChild() {
		return false
	}
	// `return ClassName$$0clinit(), ClassName$prop;`
	returnStatement := getBlock.GetFirstChild()
	if !returnStatement.IsReturn() {
		return false
	}
	// `ClassName$$0clinit(), ClassName$prop;`
	commaExpr := returnStatement.GetOnlyChild()
	if !commaExpr.IsComma() {
		return false
	}
	// `ClassName$$0clinit()`
	callNode := commaExpr.GetFirstChild()
	if !(callNode.IsCall() && callNode.HasOneChild()) {
		return false
	}
	// `ClassName$$0clinit`
	clinitQname := callNode.GetOnlyChild()
	if !clinitQname.IsQualifiedName() {
		return false
	}
	clinitQnameString := *clinitQname.GetQualifiedName()
	if !s.LooksLikeAJ2clGeneratedClassClinitName(className, clinitQnameString) {
		return false
	}
	// `ClassName$prop`
	propReference := commaExpr.GetSecondChild()
	if !propReference.IsQualifiedName() {
		return false
	}
	propReferenceString := *propReference.GetQualifiedName()

	if !s.LooksLikeAJ2clGeneratedStaticPropertyName(className, propReferenceString) {
		return false
	}
	return true
}

// matchesJ2clSetKeySignature matches J2CL property getters by looking for
// the following signature:
//
//	set: function(value) { (ClassName$$0clinit(), ClassName$$0fieldName = value)};
//	// OR this if CollapseProperties decided not to collapse the names.
//	set: function(value) { (ClassName.$clinit(), ClassName.$fieldNaNme = value)};
//
// Java: J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.matchesJ2clSetKeySignature(String, Node)
func (s *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner) matchesJ2clSetKeySignature(className *string, setKey *rhino.Node) bool {
	if !setKey.GetFirstChild().IsFunction() {
		return false
	}
	setFunction := setKey.GetFirstChild()
	if !setFunction.HasChildren() ||
		!setFunction.GetLastChild().IsBlock() ||
		!setFunction.GetSecondChild().IsParamList() {
		return false
	}
	if !setFunction.GetSecondChild().HasOneChild() {
		// There is a single parameter "value".
		return false
	}
	// `{ ClassName$$0clinit(), ClassName$$0staticProp = initialValue; }`
	setBlock := setFunction.GetLastChild()
	if !setBlock.HasOneChild() {
		return false
	}

	// `ClassName$$0clinit(), ClassName$$0staticProp = initialValue;`
	exprResult := setBlock.GetOnlyChild()
	if !exprResult.IsExprResult() {
		return false
	}
	// `ClassName$$0clinit(), ClassName$$0staticProp = initialValue`
	commaExpr := exprResult.GetOnlyChild()
	if !commaExpr.IsComma() {
		return false
	}
	// `SomeClassName$$0clinit()`
	clinitCall := commaExpr.GetFirstChild()
	if !(clinitCall.IsCall() && clinitCall.HasOneChild()) {
		return false
	}
	// `SomeClassName$$0clinit`
	clinitQname := clinitCall.GetOnlyChild()
	if !clinitQname.IsQualifiedName() {
		return false
	}
	clinitQnameString := *clinitQname.GetQualifiedName()
	if !s.LooksLikeAJ2clGeneratedClassClinitName(className, clinitQnameString) {
		return false
	}
	// `ClassName$$0staticProp = initialValue`
	// Rewriting by CollapseProperties may have changed the LHS
	// from a GETPROP to a NAME.
	assign := commaExpr.GetSecondChild()
	if !assign.IsAssign() {
		return false
	}
	lhsQname := assign.GetFirstChild()
	if !lhsQname.IsQualifiedName() {
		return false
	}
	lhsQnameString := *lhsQname.GetQualifiedName()
	if !s.LooksLikeAJ2clGeneratedStaticPropertyName(className, lhsQnameString) {
		return false
	}
	return true
}

// Java: J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.looksLikeAJ2clGeneratedClassClinitName(String, String)
func (s *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner) LooksLikeAJ2clGeneratedClassClinitName(className *string, candidate string) bool {
	if !strings.HasPrefix(candidate, *className) {
		return false
	}
	candidateAfterClassName := candidate[len(*className):]
	return candidateAfterClassName == "$$0clinit" || // collapsed by CollapseProperties
		candidateAfterClassName == ".$clinit" // not collapsed
}

// Java: J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.looksLikeAJ2clGeneratedStaticPropertyName(String, String)
func (s *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner) LooksLikeAJ2clGeneratedStaticPropertyName(className *string, candidate string) bool {
	if !strings.HasPrefix(candidate, *className) {
		return false
	}
	candidateAfterClassName := candidate[len(*className):]
	return strings.HasPrefix(candidateAfterClassName, "$$0") || // collapsed by CollapseProperties
		strings.HasPrefix(candidateAfterClassName, ".$") // not collapsed
}

// J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_GatherJ2CLClassGetterSetters
// traverses the ast and gathers get and set methods contained in
// Object.defineProperties nodes.
//
// Java: private class J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.GatherJ2CLClassGetterSetters extends AbstractPostOrderCallback
type J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_GatherJ2CLClassGetterSetters struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer                *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner
	j2clPropertiesByName *jutil.LinkedHashMap[string, J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clProperty]
}

// Java: J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.GatherJ2CLClassGetterSetters()
func (s *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner) newGatherJ2CLClassGetterSetters() *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_GatherJ2CLClassGetterSetters {
	return &J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_GatherJ2CLClassGetterSetters{
		outer:                s,
		j2clPropertiesByName: jutil.NewLinkedHashMap[string, J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clProperty](),
	}
}

// Java: J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.GatherJ2CLClassGetterSetters.getResults()
func (g *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_GatherJ2CLClassGetterSetters) getResults() *jutil.LinkedHashMap[string, J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clProperty] {
	return g.j2clPropertiesByName
}

// Java: J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.GatherJ2CLClassGetterSetters.visit(NodeTraversal, Node, Node)
func (g *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_GatherJ2CLClassGetterSetters) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if n.IsClass() {
		g.VisitClass(n)
	} else if NodeUtil_IsObjectDefinePropertiesDefinition(n) {
		g.VisitObjectDefineProperties(n)
	}
}

// Java: J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.GatherJ2CLClassGetterSetters.visitObjectDefineProperties(Node)
func (g *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_GatherJ2CLClassGetterSetters) VisitObjectDefineProperties(n *rhino.Node) {
	className := n.GetSecondChild()
	if !className.IsName() {
		return
	}
	classNameString := className.GetString()
	for _, p := range NodeUtil_GetObjectDefinedPropertiesKeys(n) {
		name := p.GetString()
		propertyLiteral := p.GetFirstChild()
		var getKey *rhino.Node
		var setKey *rhino.Node
		hasSetter := false
		for innerKey := propertyLiteral.GetFirstChild(); innerKey != nil; innerKey = innerKey.GetNext() {
			if !innerKey.IsStringKey() {
				continue
			}
			switch innerKey.GetString() {
			case "get":
				if g.outer.matchesJ2clGetKeySignature(&classNameString, innerKey) {
					getKey = innerKey
				}
			case "set":
				hasSetter = true
				if g.outer.matchesJ2clSetKeySignature(&classNameString, innerKey) {
					setKey = innerKey
				}
			default:
			}
		}
		if getKey != nil && (!hasSetter || setKey != nil) {
			g.j2clPropertiesByName.Put(
				classNameString+"."+name, g.outer.newJ2clPropertyEs5(getKey, setKey))
		}
	}
}

// Java: J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.GatherJ2CLClassGetterSetters.visitClass(Node)
func (g *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_GatherJ2CLClassGetterSetters) VisitClass(classNode *rhino.Node) {
	className := NodeUtil_GetName(classNode)
	classMembers := NodeUtil_GetClassMembers(classNode)

	setterDefByName := jutil.NewLinkedHashMap[string, *rhino.Node]()
	getterDefByName := jutil.NewLinkedHashMap[string, *rhino.Node]()

	// Collect static setters and getters.
	for memberFunction := classMembers.GetFirstChild(); memberFunction != nil; memberFunction = memberFunction.GetNext() {
		if !memberFunction.IsStaticMember() {
			// The only getters and setters we care about are static.
			continue
		}
		switch memberFunction.GetToken() {
		case rhino.Token_GETTER_DEF:
			getterDefByName.Put(memberFunction.GetString(), memberFunction)
		case rhino.Token_SETTER_DEF:
			setterDefByName.Put(memberFunction.GetString(), memberFunction)
		default:
		}
	}

	for propertyName := range getterDefByName.Keys() {
		getterDef := getterDefByName.Get(propertyName)
		setterDef := setterDefByName.Get(propertyName)
		if g.outer.matchesJ2clGetKeySignature(className, getterDef) &&
			(setterDef == nil || g.outer.matchesJ2clSetKeySignature(className, setterDef)) {
			// This is a J2cl static property.
			g.j2clPropertiesByName.Put(
				*className+"."+propertyName, g.outer.newJ2clPropertyEs6(getterDef, setterDef))
		}
	}
}

// Java: private static class J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.DetermineInlinableProperties extends AbstractPostOrderCallback
type J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_DetermineInlinableProperties struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	propertiesByName *jutil.LinkedHashMap[string, J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clProperty]
}

// Java: J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.DetermineInlinableProperties(Map<String, J2clProperty>)
func NewJ2clPropertyInlinerPass_StaticFieldGetterSetterInliner_DetermineInlinableProperties(
	allGetterSetters *jutil.LinkedHashMap[string, J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clProperty]) *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_DetermineInlinableProperties {
	return &J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_DetermineInlinableProperties{propertiesByName: allGetterSetters}
}

// Java: J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.DetermineInlinableProperties.visit(NodeTraversal, Node, Node)
func (d *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_DetermineInlinableProperties) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if NodeUtil_IsCompoundAssignmentOp(n) || n.IsInc() || n.IsDec() {
		assignmentTarget := n.GetFirstChild()
		if assignmentTarget.IsGetProp() {
			accessName := assignmentTarget.GetQualifiedName()
			var prop J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clProperty
			if accessName != nil {
				prop = d.propertiesByName.Get(*accessName)
			}
			if prop != nil {
				prop.GetJ2clPropertyBase().isSafeToInline = false
			}
		}
	}
}

// J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_InlinePropertiesPass
// looks for accesses of j2cl properties and assignments to j2cl properties.
//
// Java: private class J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.InlinePropertiesPass extends AbstractPostOrderCallback
type J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_InlinePropertiesPass struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer            *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner
	propertiesByName *jutil.LinkedHashMap[string, J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clProperty]
}

// Java: J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.InlinePropertiesPass(Map<String, J2clProperty>)
func (s *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner) newInlinePropertiesPass(
	allGetterSetters *jutil.LinkedHashMap[string, J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clProperty]) *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_InlinePropertiesPass {
	return &J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_InlinePropertiesPass{outer: s, propertiesByName: allGetterSetters}
}

// Java: J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.InlinePropertiesPass.run()
func (i *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_InlinePropertiesPass) run() {
	NodeTraversal_Traverse(i.outer.outer.compiler, i.outer.root, i)

	for prop := range i.propertiesByName.Values() {
		if prop.GetJ2clPropertyBase().isSafeToInline {
			prop.Remove()
		}
	}
}

// Java: J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.InlinePropertiesPass.visit(NodeTraversal, Node, Node)
func (i *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_InlinePropertiesPass) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	compiler := i.outer.outer.compiler
	if n.IsGetProp() {
		if parent.IsExprResult() {
			// This is a stub declaration for the type checker. See:
			// ConcretizeStaticInheritanceForInlining
			return
		}
		if NodeUtil_IsAssignmentOp(parent) && parent.GetFirstChild() == n {
			// This case should be handled below.  It needs to be inlined differently.
			return
		}
		accessName := n.GetQualifiedName()
		var prop J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clProperty
		if accessName != nil {
			prop = i.propertiesByName.Get(*accessName)
		}
		if prop != nil && prop.GetJ2clPropertyBase().isSafeToInline {
			injector := NewFunctionInjector_Builder(compiler).
				AssumeStrictThis(true).
				AssumeMinimumCapture(true).
				Build()
			functionCall := rhino.IR_Call(rhino.IR_Name("inlined_j2cl_getter"))
			n.ReplaceWith(functionCall)
			reference := NewFunctionInjector_Reference(functionCall, t.GetScope(), t.GetChunk(), FunctionInjector_InliningMode_DIRECT)
			inlinedCall := injector.Inline(reference, "", prop.GetJ2clPropertyBase().getKey.GetFirstChild())
			t.GetCompiler().ReportChangeToEnclosingScope(inlinedCall)
		}
	}

	if n.IsAssign() {
		assignmentTarget := n.GetFirstChild()
		assignmentValue := n.GetLastChild()
		if assignmentTarget.IsGetProp() {
			accessName := assignmentTarget.GetQualifiedName()
			var prop J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clProperty
			if accessName != nil {
				prop = i.propertiesByName.Get(*accessName)
			}
			if prop != nil && prop.GetJ2clPropertyBase().setKey != nil && prop.GetJ2clPropertyBase().isSafeToInline {
				injector := NewFunctionInjector_Builder(compiler).
					AssumeStrictThis(true).
					AssumeMinimumCapture(true).
					Build()
				assignmentValue.Detach()
				functionCall := rhino.IR_Call(rhino.IR_Name("inlined_j2cl_setter"), assignmentValue)
				n.ReplaceWith(functionCall)
				reference := NewFunctionInjector_Reference(functionCall, t.GetScope(), t.GetChunk(), FunctionInjector_InliningMode_BLOCK)
				injector.MaybePrepareCall(reference)
				inlinedCall := injector.Inline(reference, "", prop.GetJ2clPropertyBase().setKey.GetFirstChild())
				t.GetCompiler().ReportChangeToEnclosingScope(inlinedCall)
			}
		}
	}
}
