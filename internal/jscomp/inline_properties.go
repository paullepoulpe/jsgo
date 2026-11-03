// Port of com/google/javascript/jscomp/InlineProperties.java @ closure-compiler v20261028

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
	props       map[uint32]inlinePropInfo
	// Go-only: the properties whose entry is invalidated (an entry never
	// becomes a candidate again), to skip the map on repeated invalidations.
	invalidated smallBitSet
}

// invalidate records an invalidating reference to propID.
func (p *InlineProperties) invalidate(propID uint32) {
	if p.invalidated.Get(int(propID)) {
		return
	}
	p.invalidated.Set(int(propID))
	p.props[propID] = inlinePropInfo{invalidated: true}
}

// hasEntry reports whether propID has an entry (invalidated or a candidate).
func (p *InlineProperties) hasEntry(propID uint32) bool {
	if p.invalidated.Get(int(propID)) {
		return true
	}
	_, ok := p.props[propID]
	return ok
}

type inlinePropInfo struct {
	color       *colors.Color
	value       rhino.Node
	invalidated bool
}

var (
	inlinePropsGatherTokens = makeVisitTokenSet(
		rhino.Token_GETPROP,
		rhino.Token_STRING_KEY,
		rhino.Token_GETTER_DEF,
		rhino.Token_SETTER_DEF,
		rhino.Token_MEMBER_FUNCTION_DEF,
		rhino.Token_MEMBER_FIELD_DEF,
		rhino.Token_CALL,
	)
	inlinePropsReplaceTokens = makeVisitTokenSet(rhino.Token_GETPROP)
)

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

func (p *InlineProperties) Process(externs rhino.Node, root rhino.Node) {
	if externProps := p.compiler.GetExternPropSet(); externProps != nil {
		for _, id := range externProps.IDs() {
			p.invalidate(uint32(id))
		}
	}

	if ix := propUseIndexFor(p.compiler, externs, root); ix != nil {
		p.processIndexed(ix, root)
		return
	}

	// Find and replace the properties in non-extern AST.
	NodeTraversal_Traverse(p.compiler, root, &inlineGatherCandidates{outer: p})
	replacer := &inlineReplaceCandidates{outer: p}
	replacer.seen = replacer.seenBuf[:0]
	NodeTraversal_Traverse(p.compiler, root, replacer)
}

// processIndexed is Process with the property-use index (prop_use_index.go)
// in place of the two traversals.
//
// Go-only: GatherCandidates' callback is replayed in the traversal's order:
// the index classified each node as the callback does (inlineClassifyGather);
// an invalidating reference is applied as is, a definition is evaluated with
// the traversal state the callback reads (whether it is in the global hoist
// scope, recorded when the node was indexed, and the enclosing function, the
// node's change scope when that is a function). A definition of a property
// that already has an entry is invalidating without evaluating it, as
// maybeStoreCandidateValue decides last and every check before it is free of
// side effects.
// ReplaceCandidates' callback is replayed on the GETPROPs of the candidate
// properties, in the traversal's post-order: replacing one never creates a
// GETPROP and changes nothing after it in that order but its parent, which the
// replay visits after it, as the traversal does.
func (p *InlineProperties) processIndexed(ix *propUseIndex, root rhino.Node) {
	rp := piReplay{compiler: p.compiler}
	defer rp.done()
	g := &inlineGatherCandidates{outer: p}
	ix.each(root, piListIP, nil, func(s *piScope, it piItem) {
		if it.bits&piIPInv != 0 {
			p.invalidate(it.prop)
			return
		}
		// A definition: it can only be recorded if the property has no
		// entry yet (maybeStoreCandidateValue), and invalidates it otherwise.
		if p.hasEntry(it.prop) {
			p.invalidate(it.prop)
			return
		}
		rp.cur = it.n
		kind := inlineGatherFieldDef
		if it.n.IsGetProp() {
			kind = inlineGatherGetpropDef
		}
		g.apply(inlineTraversalState{inGlobalHoistScope: it.bits&piHoist != 0, enclosingFunction: s.enclosingFunction()},
			it.n, it.n.GetParent(), it.prop, kind)
	})

	var candidates smallBitSet
	// order-independent: builds a set
	for propID, info := range p.props {
		if !info.invalidated {
			candidates.Set(int(propID))
		}
	}
	if candidates.IsEmpty() {
		return
	}
	replacer := &inlineReplaceCandidates{outer: p}
	replacer.seen = replacer.seenBuf[:0]
	ix.each(root, piListUses, &candidates, func(_ *piScope, it piItem) {
		rp.cur = it.n
		replacer.Visit(nil, it.n, it.n.GetParent())
	})
}

// inlineTraversalState is the state of the traversal GatherCandidates reads.
type inlineTraversalState struct {
	inGlobalHoistScope bool
	enclosingFunction  rhino.Node
}

type inlineGatherCandidates struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer *InlineProperties
}

func (g *inlineGatherCandidates) visitTokens() *ast.TokenSet {
	return &inlinePropsGatherTokens
}

func (g *inlineGatherCandidates) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	g.visit(inlineTraversalState{inGlobalHoistScope: t.InGlobalHoistScope(), enclosingFunction: t.GetEnclosingFunction()}, n, parent)
}

func (g *inlineGatherCandidates) visit(t inlineTraversalState, n rhino.Node, parent rhino.Node) {
	propID, kind := inlineClassifyGather(g.outer.compiler, n, parent)
	g.apply(t, n, parent, propID, kind)
}

// inlineGatherKind is what GatherCandidates.visit does with a node.
type inlineGatherKind uint8

const (
	inlineGatherNone inlineGatherKind = iota
	// inlineGatherInvalidate: an invalidating reference.
	inlineGatherInvalidate
	// inlineGatherGetpropDef: a GETPROP in an assignment, a candidate
	// definition or invalidating.
	inlineGatherGetpropDef
	// inlineGatherFieldDef: a class field with a value, likewise.
	inlineGatherFieldDef
)

// inlineClassifyGather says what GatherCandidates.visit does with n (it reads
// n, its children and its ancestors up to the enclosing function only), and
// for which property.
func inlineClassifyGather(compiler AbstractCompiler, n rhino.Node, parent rhino.Node) (uint32, inlineGatherKind) {
	switch n.GetToken() {
	case rhino.Token_GETPROP:
		if parent.IsAssign() {
			return n.GetStringID(), inlineGatherGetpropDef
		} else if NodeUtil_IsLValue(n) {
			// Other LValue references invalidate
			// e.g. in an enhanced for loop or a destructuring statement
			return n.GetStringID(), inlineGatherInvalidate
		} else if parent.IsDelProp() {
			// Deletes invalidate
			return n.GetStringID(), inlineGatherInvalidate
		}
	case rhino.Token_STRING_KEY:
		if parent.IsObjectPattern() {
			// String keys in destructuring patterns are fine, since they just access the prop.
			return 0, inlineGatherNone
		}
		return n.GetStringID(), inlineGatherInvalidate
	case rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF, rhino.Token_MEMBER_FUNCTION_DEF:
		// For now, any object literal key, getter/setter def, or member function invalidates
		// since we do not inline functions in this pass.
		return n.GetStringID(), inlineGatherInvalidate
	case rhino.Token_MEMBER_FIELD_DEF:
		if _, ok := ast.MustMemberFieldDef(n).Value(); ok {
			return n.GetStringID(), inlineGatherFieldDef
		}
		return n.GetStringID(), inlineGatherInvalidate
	case rhino.Token_CALL:
		call := ast.MustCall(n)
		if !compiler.GetCodingConvention().IsPropertyRenameFunction(call.Callee().Node()) {
			return 0, inlineGatherNone
		}
		firstArg, ok := call.FirstArg()
		if !ok {
			return 0, inlineGatherNone
		}
		strLit, ok := firstArg.AsStringLit()
		if !ok {
			return 0, inlineGatherNone
		}
		// JSCompiler_renameProperty('foo', obj) and goog.reflect.objectProperty name `foo` for
		// reflective access, so `foo` may be written without a GETPROP, e.g. by
		// `obj[JSCompiler_renameProperty('foo', obj)] = value`.
		return strLit.Node().GetStringID(), inlineGatherInvalidate
	}
	return 0, inlineGatherNone
}

// apply does what GatherCandidates.visit does with n, classified as kind.
func (g *inlineGatherCandidates) apply(t inlineTraversalState, n rhino.Node, parent rhino.Node, propID uint32, kind inlineGatherKind) {
	var invalidatingPropRef bool
	switch kind {
	case inlineGatherNone:
		return
	case inlineGatherInvalidate:
		invalidatingPropRef = true
	case inlineGatherGetpropDef:
		invalidatingPropRef = !g.maybeRecordCandidateGetpropDefinition(t, ast.MustGetProp(n), ast.MustAssign(parent))
	case inlineGatherFieldDef:
		fieldDef := ast.MustMemberFieldDef(n)
		val, _ := fieldDef.Value()
		invalidatingPropRef = !g.maybeRecordCandidateClassFieldDefinition(fieldDef, propID, val.Node())
	}
	if invalidatingPropRef {
		g.outer.invalidate(propID)
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

func (g *inlineGatherCandidates) maybeRecordCandidateGetpropDefinition(t inlineTraversalState, gp ast.GetProp, assign ast.Assign) bool {
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
	} else if t.inGlobalHoistScope {
		if srcGetProp, ok := ast.AsGetProp(src); ok && nodeStringID(srcGetProp.Node()) == idPrototype {
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
	if !g.outer.hasEntry(propID) &&
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
func inlineInConstructor(t inlineTraversalState) bool {
	root := t.enclosingFunction
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

	for _, immediateSupertype := range r.outer.registry.GetDisambiguationSupertypesSlice(subCtor) {
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
