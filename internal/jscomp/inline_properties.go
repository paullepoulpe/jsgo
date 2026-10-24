// Port of com/google/javascript/jscomp/InlineProperties.java @ closure-compiler v20261021

/*
 * Copyright 2012 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// InlineProperties attempts to find references to properties that are known
// to be constants and inline the known value.
//
// This pass relies on type information to find these property references
// and properties are assumed to be constant if they are assigned exactly
// once, unconditionally, in either of the following contexts: (1)
// statically on a constructor, or (2) on a class's prototype.
//
// The current implementation only inlines immutable values (as defined by
// NodeUtil.isImmutableValue).
type InlineProperties struct {
	compiler    AbstractCompiler
	registry    *colors.ColorRegistry
	astAnalyzer *AstAnalyzer
	externProps *jutil.LinkedHashSet[string]

	props map[string]*InlineProperties_PropertyInfo
}

// InlineProperties_PropertyInfo is the nested class PropertyInfo.
type InlineProperties_PropertyInfo struct {
	color *colors.Color
	value rhino.Node
}

// Java: InlineProperties.PropertyInfo(Color, Node)
func NewInlineProperties_PropertyInfo(color *colors.Color, value rhino.Node) *InlineProperties_PropertyInfo {
	return &InlineProperties_PropertyInfo{color: color, value: value}
}

var InlineProperties_INVALIDATED = NewInlineProperties_PropertyInfo(nil, rhino.NilNode)

// Java: InlineProperties(AbstractCompiler)
func NewInlineProperties(compiler AbstractCompiler) *InlineProperties {
	p := &InlineProperties{
		compiler:    compiler,
		registry:    compiler.GetColorRegistry(),
		astAnalyzer: compiler.GetAstAnalyzer(),
		props:       make(map[string]*InlineProperties_PropertyInfo, 256),
	}
	p.InvalidateExternProperties()
	return p
}

// Java: InlineProperties.invalidateExternProperties()
func (p *InlineProperties) InvalidateExternProperties() {
	// Invalidate properties defined in externs.
	p.externProps = p.compiler.GetExternProperties()
}

// GetColor gets the JSType from the Node argument and verifies that it is
// present.
//
// Java: InlineProperties.getColor(Node)
func (p *InlineProperties) GetColor(n rhino.Node) *colors.Color {
	color := n.GetColor()

	if color == nil {
		return colors.StandardColors_UNKNOWN
	} else {
		return color
	}
}

// Java: InlineProperties.process(Node, Node)
func (p *InlineProperties) Process(externs rhino.Node, root rhino.Node) {
	// Find and replace the properties in non-extern AST.
	NodeTraversal_Traverse(p.compiler, root, &InlineProperties_GatherCandidates{outer: p})
	NodeTraversal_Traverse(p.compiler, root, NewInlineProperties_ReplaceCandidates(p))
}

// InlineProperties_GatherCandidates is the inner class GatherCandidates.
type InlineProperties_GatherCandidates struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer *InlineProperties
}

// Java: InlineProperties.GatherCandidates.visit(NodeTraversal, Node, Node)
func (g *InlineProperties_GatherCandidates) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	// These are assigned at most once in the branches below
	var invalidatingPropRef bool
	var propName string
	if n.IsGetProp() {
		propName = n.GetString()
		if parent.IsAssign() {
			invalidatingPropRef = !g.MaybeRecordCandidateGetpropDefinition(t, n, parent)
		} else if NodeUtil_IsLValue(n) {
			// Other LValue references invalidate
			// e.g. in an enhanced for loop or a destructuring statement
			invalidatingPropRef = true
		} else if parent.IsDelProp() {
			// Deletes invalidate
			invalidatingPropRef = true
		} else {
			// A property read doesn't invalidate
			invalidatingPropRef = false
		}
	} else if (n.IsStringKey() && !n.GetParent().IsObjectPattern()) ||
		n.IsGetterDef() ||
		n.IsSetterDef() ||
		n.IsMemberFunctionDef() {
		propName = n.GetString()
		// For now, any object literal key invalidates
		// TODO(johnlenz): support prototype properties like:
		//   foo.prototype = { a: 1, b: 2 };
		// TODO(johnlenz): Object.create(), Object.createProperty
		// and getter/setter defs and member functions also invalidate
		// since we do not inline functions in this pass
		// Note that string keys in destructuring patterns are fine, since they just access the prop
		invalidatingPropRef = true
	} else if n.IsMemberFieldDef() {
		propName = n.GetString()
		if n.HasChildren() {
			// class field with initialization
			invalidatingPropRef = !g.MaybeRecordCandidateClassFieldDefinition(n)
		} else {
			// class field with no initialization
			invalidatingPropRef = true
		}
	} else if n.IsCall() &&
		g.outer.compiler.GetCodingConvention().IsPropertyRenameFunction(n.GetFirstChild()) &&
		n.HasMoreThanOneChild() &&
		n.GetSecondChild().IsStringLit() {
		// JSCompiler_renameProperty('foo', obj) and goog.reflect.objectProperty name `foo` for
		// reflective access, so `foo` may be written without a GETPROP, e.g. by
		// `obj[JSCompiler_renameProperty('foo', obj)] = value`.
		propName = n.GetSecondChild().GetString()
		invalidatingPropRef = true
	} else {
		return
	}

	if invalidatingPropRef {
		// checkNotNull(propName): propName is a non-null Java String here.
		g.InvalidateProperty(propName)
	}
}

// MaybeRecordCandidateClassFieldDefinition returns whether this is a valid
// definition for a candidate class field.
//
// Java: InlineProperties.GatherCandidates.maybeRecordCandidateClassFieldDefinition(Node)
func (g *InlineProperties_GatherCandidates) MaybeRecordCandidateClassFieldDefinition(n rhino.Node) bool {
	jutil.CheckStateMsg1(n.IsMemberFieldDef(), "%s", n)
	value := n.GetFirstChild()
	propName := n.GetString()
	classNode := n.GetGrandparent()
	var c *colors.Color

	if n.IsStaticMember() {
		c = g.outer.GetColor(classNode)
	} else {
		possibleInstances := g.outer.GetColor(classNode).GetInstanceColors()
		if possibleInstances.IsEmpty() {
			c = colors.StandardColors_UNKNOWN
		} else {
			c = colors.Color_CreateUnion(possibleInstances)
		}
	}

	return g.MaybeStoreCandidateValue(c, propName, value)
}

// MaybeRecordCandidateGetpropDefinition returns whether this is a valid
// definition for a candidate property.
//
// Java: InlineProperties.GatherCandidates.maybeRecordCandidateGetpropDefinition(NodeTraversal, Node, Node)
func (g *InlineProperties_GatherCandidates) MaybeRecordCandidateGetpropDefinition(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	jutil.CheckStateMsg1(n.IsGetProp() && parent.IsAssign(), "%s", n)
	src := n.GetFirstChild()
	propName := n.GetString()
	value := parent.GetLastChild()

	if src.IsThis() {
		// This is a simple assignment like:
		//    this.foo = 1;
		if g.InConstructor(t) {
			// This may be a valid assignment.
			return g.MaybeStoreCandidateValue(g.outer.GetColor(src), propName, value)
		}
		return false
	} else if t.InGlobalHoistScope() && src.IsGetProp() && src.GetString() == "prototype" {
		// This is a prototype assignment like:
		//    x.prototype.foo = 1;
		instanceType := g.outer.GetColor(src)
		return g.MaybeStoreCandidateValue(instanceType, propName, value)
	} else if t.InGlobalHoistScope() {
		// This is a static assignment like:
		//    x.foo = 1;
		targetType := g.outer.GetColor(src)
		if targetType != nil && targetType.IsConstructor() {
			return g.MaybeStoreCandidateValue(targetType, propName, value)
		}
	}
	return false
}

// Java: InlineProperties.GatherCandidates.invalidateProperty(String)
func (g *InlineProperties_GatherCandidates) InvalidateProperty(propName string) {
	if g.outer.externProps != nil && g.outer.externProps.Contains(propName) {
		return
	}
	g.outer.props[propName] = InlineProperties_INVALIDATED
}

// MaybeStoreCandidateValue adds the candidate property to the map if it
// meets all constness and immutability criteria, and is not already present
// in the map. If the property was already present, it is invalidated.
// Returns true if the property was successfully added.
//
// Java: InlineProperties.GatherCandidates.maybeStoreCandidateValue(Color, String, Node)
func (g *InlineProperties_GatherCandidates) MaybeStoreCandidateValue(color *colors.Color, propName string, value rhino.Node) bool {
	rhino.CheckNotNull(value)
	if g.outer.props[propName] == nil &&
		(g.outer.externProps == nil || !g.outer.externProps.Contains(propName)) &&
		!color.IsInvalidating() &&
		NodeUtil_IsImmutableValue(value) &&
		NodeUtil_IsExecutedExactlyOnce(value) {
		g.outer.props[propName] = NewInlineProperties_PropertyInfo(color, value)
		return true
	}
	return false
}

// InConstructor returns whether the traversal is directly in an ES6 class
// constructor or an @constructor function
//
// This returns false for nested functions inside ctors, including arrow
// functions (even though the `this` is the same). This pass only cares
// about property definitions executed once per ctor invocation, and in
// general we don't know how many times an arrow fn will be executed. In the
// future, we could special case arrow fn IIFEs in this pass if it becomes
// useful.
//
// Java: InlineProperties.GatherCandidates.inConstructor(NodeTraversal)
func (g *InlineProperties_GatherCandidates) InConstructor(t *NodeTraversal) bool {
	root := t.GetEnclosingFunction()
	if root.IsNil() { // we might be in the global scope
		return false
	}
	return NodeUtil_IsConstructor(root)
}

// InlineProperties_ReplaceCandidates is the inner class ReplaceCandidates.
type InlineProperties_ReplaceCandidates struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer *InlineProperties

	hasInSupertypesListSeenSet *jutil.ObjectLinkedHashSet[*colors.Color]
}

// Java: InlineProperties.ReplaceCandidates()
func NewInlineProperties_ReplaceCandidates(outer *InlineProperties) *InlineProperties_ReplaceCandidates {
	return &InlineProperties_ReplaceCandidates{
		outer:                      outer,
		hasInSupertypesListSeenSet: jutil.NewObjectLinkedHashSet[*colors.Color](),
	}
}

// Java: InlineProperties.ReplaceCandidates.visit(NodeTraversal, Node, Node)
func (r *InlineProperties_ReplaceCandidates) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if n.IsGetProp() && !NodeUtil_IsLValue(n) {
		target := n.GetFirstChild()
		propName := n.GetString()
		info := r.outer.props[propName]
		if info != nil && info != InlineProperties_INVALIDATED && r.IsMatchingType(target, info.color) {
			replacement := info.value.CloneTree()
			if r.outer.astAnalyzer.MayHaveSideEffects(n.GetFirstChild()) {
				replacement = rhino.IR_Comma(n.RemoveFirstChild(), replacement).Srcref(n)
			}
			n.ReplaceWith(replacement)
			r.outer.compiler.ReportChangeToEnclosingScope(replacement)
		}
	}
}

// Java: InlineProperties.ReplaceCandidates.isMatchingType(Node, Color)
func (r *InlineProperties_ReplaceCandidates) IsMatchingType(n rhino.Node, src *colors.Color) bool {
	src = InlineProperties_RemoveNullAndUndefinedIfUnion(src)
	dest := InlineProperties_RemoveNullAndUndefinedIfUnion(r.outer.GetColor(n))
	if dest.IsInvalidating() {
		return false
	}
	if dest.IsUnion() || src.IsUnion() {
		return false
	}
	return r.HasInSupertypesList(dest, src)
}

// Java: InlineProperties.ReplaceCandidates.hasInSupertypesList(Color, Color)
func (r *InlineProperties_ReplaceCandidates) HasInSupertypesList(subCtor *colors.Color, superCtor *colors.Color) bool {
	defer r.hasInSupertypesListSeenSet.Remove(subCtor)
	if !r.hasInSupertypesListSeenSet.Add(subCtor) {
		return false
	}
	if subCtor == nil || superCtor == nil {
		return false
	}
	if subCtor.Equals(superCtor) {
		return true
	}

	for immediateSupertype := range r.outer.registry.GetDisambiguationSupertypes(subCtor).All() {
		if !immediateSupertype.IsUnion() && r.HasInSupertypesList(immediateSupertype, superCtor) {
			return true
		}
	}
	return false
}

// Java: InlineProperties.removeNullAndUndefinedIfUnion(Color)
func InlineProperties_RemoveNullAndUndefinedIfUnion(original *colors.Color) *colors.Color {
	if original.IsUnion() {
		return original.SubtractNullOrVoid()
	}
	return original
}
