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
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
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
	prototypeID uint32
	props       map[uint32]inlinePropInfo
}

type inlinePropInfo struct {
	color       *colors.Color
	value       rhino.Node
	invalidated bool
}

var (
	inlinePropsGatherTokens = func() ast.TokenSet {
		var s ast.TokenSet
		s.Add(rhino.Token_GETPROP)
		s.Add(rhino.Token_STRING_KEY)
		s.Add(rhino.Token_GETTER_DEF)
		s.Add(rhino.Token_SETTER_DEF)
		s.Add(rhino.Token_MEMBER_FUNCTION_DEF)
		s.Add(rhino.Token_MEMBER_FIELD_DEF)
		s.Add(rhino.Token_CALL)
		return s
	}()
	inlinePropsReplaceTokens = func() ast.TokenSet {
		var s ast.TokenSet
		s.Add(rhino.Token_GETPROP)
		return s
	}()
)

// Java: InlineProperties(AbstractCompiler)
func NewInlineProperties(compiler AbstractCompiler) *InlineProperties {
	return &InlineProperties{
		compiler:    compiler,
		registry:    compiler.GetColorRegistry(),
		astAnalyzer: compiler.GetAstAnalyzer(),
		props:       make(map[uint32]inlinePropInfo, 256),
	}
}

func inlinePropColor(n rhino.Node) *colors.Color {
	if c := n.GetColor(); c != nil {
		return c
	}
	return colors.StandardColors_UNKNOWN
}

// Java: InlineProperties.process(Node, Node)
func (p *InlineProperties) Process(externs rhino.Node, root rhino.Node) {
	p.prototypeID = uint32(strpool.Intern("prototype"))
	if externProps := p.compiler.GetExternProperties(); externProps != nil {
		for name := range externProps.All() {
			if id := strpool.Lookup(name); id != 0 {
				p.props[uint32(id)] = inlinePropInfo{invalidated: true}
			}
		}
	}

	// Find and replace the properties in non-extern AST.
	NodeTraversal_Traverse(p.compiler, root, &inlineGatherCandidates{outer: p})
	replacer := &inlineReplaceCandidates{outer: p}
	replacer.seen = replacer.seenBuf[:0]
	NodeTraversal_Traverse(p.compiler, root, replacer)
}

type inlineGatherCandidates struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer *InlineProperties
}

func (g *inlineGatherCandidates) visitTokens() *ast.TokenSet {
	return &inlinePropsGatherTokens
}

// Java: InlineProperties.GatherCandidates.visit(NodeTraversal, Node, Node)
func (g *inlineGatherCandidates) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	var invalidatingPropRef bool
	var propID uint32

	switch n.GetToken() {
	case rhino.Token_GETPROP:
		gp := ast.MustGetProp(n)
		propID = n.GetStringID()
		if assign, ok := ast.AsAssign(parent); ok {
			invalidatingPropRef = !g.maybeRecordCandidateGetpropDefinition(t, gp, assign)
		} else if NodeUtil_IsLValue(n) {
			// Other LValue references invalidate
			// e.g. in an enhanced for loop or a destructuring statement
			invalidatingPropRef = true
		} else if parent.IsDelProp() {
			// Deletes invalidate
			invalidatingPropRef = true
		}
	case rhino.Token_STRING_KEY:
		if parent.IsObjectPattern() {
			// String keys in destructuring patterns are fine, since they just access the prop.
			return
		}
		propID = n.GetStringID()
		invalidatingPropRef = true
	case rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF, rhino.Token_MEMBER_FUNCTION_DEF:
		// For now, any object literal key, getter/setter def, or member function invalidates
		// since we do not inline functions in this pass.
		propID = n.GetStringID()
		invalidatingPropRef = true
	case rhino.Token_MEMBER_FIELD_DEF:
		fieldDef := ast.MustMemberFieldDef(n)
		propID = n.GetStringID()
		if val, ok := fieldDef.Value(); ok {
			invalidatingPropRef = !g.maybeRecordCandidateClassFieldDefinition(fieldDef, propID, val.Node())
		} else {
			invalidatingPropRef = true
		}
	case rhino.Token_CALL:
		call := ast.MustCall(n)
		if !g.outer.compiler.GetCodingConvention().IsPropertyRenameFunction(call.Callee().Node()) {
			return
		}
		firstArg, ok := call.FirstArg()
		if !ok {
			return
		}
		strLit, ok := firstArg.AsStringLit()
		if !ok {
			return
		}
		// JSCompiler_renameProperty('foo', obj) and goog.reflect.objectProperty name `foo` for
		// reflective access, so `foo` may be written without a GETPROP, e.g. by
		// `obj[JSCompiler_renameProperty('foo', obj)] = value`.
		propID = strLit.Node().GetStringID()
		invalidatingPropRef = true
	default:
		return
	}

	if invalidatingPropRef {
		g.outer.props[propID] = inlinePropInfo{invalidated: true}
	}
}

func (g *inlineGatherCandidates) maybeRecordCandidateClassFieldDefinition(fieldDef ast.MemberFieldDef, propID uint32, value rhino.Node) bool {
	classNode := fieldDef.Node().GetGrandparent()
	var c *colors.Color
	if fieldDef.IsStatic() {
		c = inlinePropColor(classNode)
	} else {
		possibleInstances := inlinePropColor(classNode).GetInstanceColors()
		if len(possibleInstances) == 0 {
			c = colors.StandardColors_UNKNOWN
		} else {
			c = colors.Color_CreateUnionSlice(possibleInstances)
		}
	}
	return g.maybeStoreCandidateValue(c, propID, value)
}

func (g *inlineGatherCandidates) maybeRecordCandidateGetpropDefinition(t *NodeTraversal, gp ast.GetProp, assign ast.Assign) bool {
	src := gp.Object().Node()
	propID := gp.Node().GetStringID()
	value := assign.Value().Node()

	if src.IsThis() {
		// This is a simple assignment like:
		//    this.foo = 1;
		if inlineInConstructor(t) {
			return g.maybeStoreCandidateValue(inlinePropColor(src), propID, value)
		}
		return false
	} else if t.InGlobalHoistScope() {
		if srcGetProp, ok := ast.AsGetProp(src); ok && srcGetProp.Node().GetStringID() == g.outer.prototypeID {
			// This is a prototype assignment like:
			//    x.prototype.foo = 1;
			return g.maybeStoreCandidateValue(inlinePropColor(src), propID, value)
		}
		// This is a static assignment like:
		//    x.foo = 1;
		targetType := inlinePropColor(src)
		if targetType != nil && targetType.IsConstructor() {
			return g.maybeStoreCandidateValue(targetType, propID, value)
		}
	}
	return false
}

func (g *inlineGatherCandidates) maybeStoreCandidateValue(color *colors.Color, propID uint32, value rhino.Node) bool {
	rhino.CheckNotNull(value)
	if _, exists := g.outer.props[propID]; !exists &&
		!color.IsInvalidating() &&
		NodeUtil_IsImmutableValue(value) &&
		NodeUtil_IsExecutedExactlyOnce(value) {
		g.outer.props[propID] = inlinePropInfo{color: color, value: value}
		return true
	}
	return false
}

// inlineInConstructor returns whether the traversal is directly in an ES6 class
// constructor or an @constructor function.
//
// This returns false for nested functions inside ctors, including arrow
// functions (even though the `this` is the same). This pass only cares
// about property definitions executed once per ctor invocation, and in
// general we don't know how many times an arrow fn will be executed.
func inlineInConstructor(t *NodeTraversal) bool {
	root := t.GetEnclosingFunction()
	if root.IsNil() {
		return false
	}
	return NodeUtil_IsConstructor(root)
}

type inlineReplaceCandidates struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer   *InlineProperties
	seen    []*colors.Color
	seenBuf [16]*colors.Color
}

func (r *inlineReplaceCandidates) visitTokens() *ast.TokenSet {
	return &inlinePropsReplaceTokens
}

// Java: InlineProperties.ReplaceCandidates.visit(NodeTraversal, Node, Node)
func (r *inlineReplaceCandidates) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if NodeUtil_IsLValue(n) {
		return
	}
	gp, ok := ast.AsGetProp(n)
	if !ok {
		return
	}
	info, ok := r.outer.props[n.GetStringID()]
	if !ok || info.invalidated {
		return
	}
	target := gp.Object().Node()
	if r.isMatchingType(target, info.color) {
		replacement := info.value.CloneTree()
		if r.outer.astAnalyzer.MayHaveSideEffects(target) {
			replacement = rhino.IR_Comma(gp.Object().Detach().Node(), replacement).Srcref(n)
		}
		n.ReplaceWith(replacement)
		r.outer.compiler.ReportChangeToEnclosingScope(replacement)
	}
}

func (r *inlineReplaceCandidates) isMatchingType(n rhino.Node, src *colors.Color) bool {
	src = inlineRemoveNullAndUndefinedIfUnion(src)
	dest := inlineRemoveNullAndUndefinedIfUnion(inlinePropColor(n))
	if dest.IsInvalidating() || dest.IsUnion() || src.IsUnion() {
		return false
	}
	if src == nil {
		return false
	}
	r.seen = r.seen[:0]
	return r.hasInSupertypesList(dest, src)
}

func (r *inlineReplaceCandidates) hasInSupertypesList(subCtor *colors.Color, superCtor *colors.Color) bool {
	if subCtor == nil {
		return false
	}
	for _, s := range r.seen {
		if s == subCtor || s.Equals(subCtor) {
			return false
		}
	}
	r.seen = append(r.seen, subCtor)
	if subCtor == superCtor || subCtor.Equals(superCtor) {
		return true
	}

	for immediateSupertype := range r.outer.registry.GetDisambiguationSupertypes(subCtor).All() {
		if !immediateSupertype.IsUnion() && r.hasInSupertypesList(immediateSupertype, superCtor) {
			return true
		}
	}
	return false
}

func inlineRemoveNullAndUndefinedIfUnion(original *colors.Color) *colors.Color {
	if original.IsUnion() {
		return original.SubtractNullOrVoid()
	}
	return original
}
