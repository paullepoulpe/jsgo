// Port of com/google/javascript/jscomp/disambiguate/ColorFindPropertyReferences.java @ closure-compiler v20261028

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

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// isInObjectDefinePropertiesCall reports whether n is an argument of an
// Object.defineProperties (or coding-convention equivalent) call.
func isInObjectDefinePropertiesCall(n rhino.Node) bool {
	parent := n.GetParent()
	return parent.IsCall() && NodeUtil_IsObjectDefinePropertiesDefinition(parent)
}

// unpackObjectDefineProperties unpacks Object.defineProperties(typeObj, props)
// when call is a definition call whose second argument is an object literal.
func unpackObjectDefineProperties(call ast.Call) (typeObj ast.Arg, props ast.ObjectLit, ok bool) {
	if !NodeUtil_IsObjectDefinePropertiesDefinition(call.Node()) {
		return ast.Arg{}, ast.ObjectLit{}, false
	}
	firstArg, hasFirst := call.FirstArg()
	if !hasFirst {
		return ast.Arg{}, ast.ObjectLit{}, false
	}
	secondArg, hasSecond := firstArg.Next()
	if !hasSecond {
		return ast.Arg{}, ast.ObjectLit{}, false
	}
	props, ok = secondArg.AsObjectLit()
	return firstArg, props, ok
}

// classMemberOwners lazily computes and caches the owner Colors for static,
// instance-field, and prototype-method members of an ES6 class.
type classMemberOwners struct {
	ctor     *colors.Color
	instance *colors.Color
	proto    *colors.Color
}

func colorUnionOrUnknown(cs []*colors.Color) *colors.Color {
	if len(cs) == 0 {
		return colors.StandardColors_UNKNOWN
	}
	return colors.Color_CreateUnionSlice(cs)
}

func (c *classMemberOwners) ownerOf(member ast.ClassMember) *colors.Color {
	n := member.Node()
	if n.IsStaticMember() {
		return c.ctor
	}
	if _, ok := member.AsMemberFieldDef(); ok {
		if c.instance == nil {
			if c.ctor == nil {
				c.instance = colors.StandardColors_UNKNOWN
			} else {
				c.instance = colorUnionOrUnknown(c.ctor.GetInstanceColors())
			}
		}
		return c.instance
	}
	jutil.CheckStateMsg1(n.IsMemberFunctionDef() || n.IsGetterDef() || n.IsSetterDef(), "%s", n)
	if c.proto == nil {
		if c.ctor == nil {
			c.proto = colors.StandardColors_UNKNOWN
		} else {
			c.proto = colorUnionOrUnknown(c.ctor.GetPrototypes())
		}
	}
	return c.proto
}

// ColorFindPropertyReferences traverses the AST, collecting connections
// between JSTypes, property accesses, and their accociated Nodes.
//
// Also collects declarations of constructors even if no property accesses
// are visible.
//
// This callback is intended for both source and externs.
//
// Java: final class ColorFindPropertyReferences extends AbstractPostOrderCallback
type ColorFindPropertyReferences struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	propIndex  *propertyClusteringIndex
	ownerMap   map[*colors.Color]*ColorGraphNode
	propAssocs map[propAssocKey]ColorGraphNode_PropAssociation

	colorGraphNodeFactory *ColorGraphNodeFactory
	isPropertyReflector   ColorFindPropertyReferences_IsPropertyReflector
	visitor               ast.Visitor[*NodeTraversal]
	kinds                 ast.TokenSet
}

// ColorFindPropertyReferences_IsPropertyReflector tests whether the named JS
// function is a "property reflector"; a function that treats a string
// literal as a property name.
//
// Such a function is assumed to have the following signature:
//
//  1. A string literal that is the name being reflected.
//  2. An object on which the property is referenced. (Optional)
//  3. Var args.
//
// Go-only: language-adapter — a @FunctionalInterface, so a Go func type.
//
// Java: ColorFindPropertyReferences.IsPropertyReflector.test(Node)
type ColorFindPropertyReferences_IsPropertyReflector func(nameNode rhino.Node) bool

// Java: ColorFindPropertyReferences(ColorGraphNodeFactory, IsPropertyReflector)
func NewColorFindPropertyReferences(
	colorGraphNodeFactory *ColorGraphNodeFactory, isPropertyReflector ColorFindPropertyReferences_IsPropertyReflector) *ColorFindPropertyReferences {
	r := &ColorFindPropertyReferences{
		propIndex:             newPropertyClusteringIndex(4096),
		ownerMap:              make(map[*colors.Color]*ColorGraphNode, 512),
		colorGraphNodeFactory: colorGraphNodeFactory,
		isPropertyReflector:   isPropertyReflector,
	}
	r.visitor = ast.Visitor[*NodeTraversal]{
		GetProp: func(t *NodeTraversal, gp ast.GetProp) {
			r.registerPropertyUse(t, gp.Node(), gp.Object().Node().GetColor())
		},
		OptChainGetProp: func(t *NodeTraversal, gp ast.OptChainGetProp) {
			r.registerPropertyUse(t, gp.Node(), gp.Object().Node().GetColor())
		},
		ObjectLit: func(t *NodeTraversal, lit ast.ObjectLit) {
			// Object.defineProperties literals are handled at the CALL node.
			if !isInObjectDefinePropertiesCall(lit.Node()) {
				r.traverseObjectlitLike(t, lit.Node(), lit.Node().GetColor())
			}
		},
		ObjectPattern: func(t *NodeTraversal, pat ast.ObjectPattern) {
			r.traverseObjectlitLike(t, pat.Node(), pat.Node().GetColor())
		},
		Call:     r.handleCall,
		Class:    r.handleClass,
		Function: func(_ *NodeTraversal, fn ast.Function) { r.handleFunction(fn) },
	}
	r.kinds = r.visitor.Tokens()
	return r
}

func (r *ColorFindPropertyReferences) createNode(owner *colors.Color) *ColorGraphNode {
	if owner == nil {
		owner = colors.StandardColors_UNKNOWN
	}
	if flat := r.ownerMap[owner]; flat != nil {
		return flat
	}
	flat := r.colorGraphNodeFactory.CreateNode(owner)
	r.ownerMap[owner] = flat
	return flat
}

// Java: ColorFindPropertyReferences.getPropertyIndex()
func (r *ColorFindPropertyReferences) GetPropertyIndex() *propertyClusteringIndex {
	tmp := r.propIndex
	r.propIndex = nil
	return tmp
}

var (
	_ nodeTraversal_noVisitLeaves = (*ColorFindPropertyReferences)(nil)
	_ nodeTraversal_visitTokens   = (*ColorFindPropertyReferences)(nil)
)

func (*ColorFindPropertyReferences) noVisitLeaves() {}

func (r *ColorFindPropertyReferences) visitTokens() *ast.TokenSet {
	return &r.kinds
}

// Java: ColorFindPropertyReferences.visit(NodeTraversal, Node, Node)
func (r *ColorFindPropertyReferences) Visit(t *NodeTraversal, n rhino.Node, _ rhino.Node) {
	r.visitor.Visit(t, n)
}

// handleCall examines calls in case they are Object.defineProperties calls
//
// Java: ColorFindPropertyReferences.handleCall(NodeTraversal, Node)
func (r *ColorFindPropertyReferences) handleCall(t *NodeTraversal, call ast.Call) {
	target := call.Callee().Node()
	if !target.IsQualifiedName() {
		return
	}

	if r.isPropertyReflector(target) {
		r.handlePropertyReflectorCall(t, call)
	} else if typeObj, objectLit, ok := unpackObjectDefineProperties(call); ok {
		r.traverseObjectlitLike(t, objectLit.Node(), typeObj.Node().GetColor())
	}
}

// Java: ColorFindPropertyReferences.handleClass(NodeTraversal, Node)
func (r *ColorFindPropertyReferences) handleClass(t *NodeTraversal, cls ast.Class) {
	classType := cls.Node().GetColor()
	body := cls.Body()
	owners := classMemberOwners{ctor: classType}
	for member := range body.Members() {
		child := member.Node()
		if isNamedPropertyMember(body.Node(), child) {
			r.registerPropertyUse(t, child, owners.ownerOf(member))
		}
	}
	r.createNode(classType)
}

// Java: ColorFindPropertyReferences.handleFunction(Node)
func (r *ColorFindPropertyReferences) handleFunction(fn ast.Function) {
	// Ensure the flattener knows about the class constructor and instance. Even if we don't see any
	// direct property references off either type, it's possible that a property is referenced
	// via a superclass or implemented interface.
	// e.g. `const /** !FooInterface */ x = new Foo(); x.method();`
	fnType := fn.Node().GetColor()
	if fnType != nil && len(fnType.GetInstanceColors()) != 0 {
		r.createNode(fnType)
	}
}

// Java: ColorFindPropertyReferences.handlePropertyReflectorCall(NodeTraversal, Node)
func (r *ColorFindPropertyReferences) handlePropertyReflectorCall(t *NodeTraversal, call ast.Call) {
	firstArg, ok := call.FirstArg()
	if !ok {
		return
	}
	name, ok := firstArg.AsStringLit()
	if !ok {
		return
	}

	var objColor *colors.Color
	if obj, ok := firstArg.Next(); ok {
		objColor = obj.Node().GetColor()
	}
	r.registerPropertyUse(t, name.Node(), objColor)
}

// Java: ColorFindPropertyReferences.traverseObjectlitLike(NodeTraversal, Node, Function<Node, Color>)
func (r *ColorFindPropertyReferences) traverseObjectlitLike(t *NodeTraversal, n rhino.Node, owner *colors.Color) {
	// The keys in an object pattern are r-values, not l-values, but they are still accesses.
	jutil.CheckState(n.IsObjectLit() || n.IsObjectPattern() || n.IsClassMembers())

	for child := n.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		if isNamedPropertyMember(n, child) {
			r.registerPropertyUse(t, child, owner)
		}
	}
}

func isNamedPropertyMember(container, child rhino.Node) bool {
	switch child.GetToken() {
	case rhino.Token_BLOCK, rhino.Token_COMPUTED_PROP, rhino.Token_COMPUTED_FIELD_DEF, rhino.Token_OBJECT_REST, rhino.Token_OBJECT_SPREAD:
		return false
	case rhino.Token_STRING_KEY, rhino.Token_MEMBER_FUNCTION_DEF, rhino.Token_MEMBER_FIELD_DEF, rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF:
		// Quoted keys won't be renamed due to our assumptions. Ignore them.
		return !child.IsQuotedStringKey()
	default:
		panic(jutil.NewIllegalStateExceptionMsg(
			"Unexpected child of " + container.GetToken().String() + ": " + child.ToStringTree()))
	}
}

// registerPropertyUse updates all datastructures as necessary to consider
// property use `site` from type `owner`.
//
// Java: ColorFindPropertyReferences.registerPropertyUse(NodeTraversal, Node, Color)
func (r *ColorFindPropertyReferences) registerPropertyUse(t *NodeTraversal, site rhino.Node, owner *colors.Color) {
	prop := r.propIndex.getOrCreateByID(site.GetStringID())
	flatOwner := r.createNode(owner)

	// Set the initial condition for flowing this property along the graph.
	if flatOwner.AddProp(prop) {
		if r.propAssocs != nil {
			r.propAssocs[propAssocKey{node: flatOwner, prop: prop}] = ColorGraphNode_PropAssociation_AST
		}
		// Make sure there's a cluster for this name/type combination.
		prop.GetClusters().Add(flatOwner)
	}
	// Record the site to rename once clusters are found. If it's an extern, we won't rename anyway.
	prop.GetUseSites().PutNew(site, flatOwner)

	// Track the cluster of types whose properties must keep their original name after
	// disambiguation. Note: an "enum type" is the type of an enum object like "{STOP: 0, GO: 1}".
	// NOTE: we can't use site.isFromExterns() because sometimes nodes have source file information
	// that doesn't match the containing script. This could lead to renaming properties that are
	// in externs if the property node didn't have an externs source file. Related: b/186056977.
	if t.GetCurrentScript().IsFromExterns() ||
		(owner != nil && owner.GetPropertiesKeepOriginalName()) {
		prop.RegisterOriginalNameType(flatOwner)
	}
}
