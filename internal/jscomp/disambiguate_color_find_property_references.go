// Port of com/google/javascript/jscomp/disambiguate/ColorFindPropertyReferences.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

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

	propIndex *jutil.LinkedHashMap[string, *PropertyClustering]
	ownerMap  map[*colors.Color]*ColorGraphNode

	colorGraphNodeFactory *ColorGraphNodeFactory
	isPropertyReflector   ColorFindPropertyReferences_IsPropertyReflector
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
	return &ColorFindPropertyReferences{
		propIndex:             jutil.NewLinkedHashMapWithCapacity[string, *PropertyClustering](4096),
		ownerMap:              make(map[*colors.Color]*ColorGraphNode, 512),
		colorGraphNodeFactory: colorGraphNodeFactory,
		isPropertyReflector:   isPropertyReflector,
	}
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
func (r *ColorFindPropertyReferences) GetPropertyIndex() *jutil.LinkedHashMap[string, *PropertyClustering] {
	tmp := r.propIndex
	r.propIndex = nil
	return tmp
}

// Java: ColorFindPropertyReferences.visit(NodeTraversal, Node, Node)
func (r *ColorFindPropertyReferences) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_GETPROP, rhino.Token_OPTCHAIN_GETPROP:
		r.registerPropertyUse(t, n, n.GetFirstChild().GetColor())
	case rhino.Token_OBJECTLIT:
		r.handleObjectLit(t, n)
	case rhino.Token_CALL:
		r.handleCall(t, n)
	case rhino.Token_CLASS:
		r.handleClass(t, n)
	case rhino.Token_OBJECT_PATTERN:
		r.handleObjectPattern(t, n)
	case rhino.Token_FUNCTION:
		r.handleFunction(n)
	default:
	}
}

// Java: ColorFindPropertyReferences.handleObjectLit(NodeTraversal, Node)
func (r *ColorFindPropertyReferences) handleObjectLit(t *NodeTraversal, n rhino.Node) {
	// Object.defineProperties literals are handled at the CALL node.
	if n.GetParent().IsCall() && NodeUtil_IsObjectDefinePropertiesDefinition(n.GetParent()) {
		return
	}

	c := n.GetColor()
	r.traverseObjectlitLike(t, n, func(rhino.Node) *colors.Color { return c })
}

// handleCall examines calls in case they are Object.defineProperties calls
//
// Java: ColorFindPropertyReferences.handleCall(NodeTraversal, Node)
func (r *ColorFindPropertyReferences) handleCall(t *NodeTraversal, call rhino.Node) {
	target := call.GetFirstChild()
	if !target.IsQualifiedName() {
		return
	}

	if r.isPropertyReflector(target) {
		r.handlePropertyReflectorCall(t, call)
	} else if NodeUtil_IsObjectDefinePropertiesDefinition(call) {
		r.handleObjectDefineProperties(t, call)
	}
}

// Java: ColorFindPropertyReferences.handleClass(NodeTraversal, Node)
func (r *ColorFindPropertyReferences) handleClass(t *NodeTraversal, classNode rhino.Node) {
	classType := classNode.GetColor()
	members := NodeUtil_GetClassMembers(classNode)
	jutil.CheckState(members.IsObjectLit() || members.IsObjectPattern() || members.IsClassMembers())

	var instanceOwner *colors.Color
	var protoOwner *colors.Color
	for child := members.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		switch child.GetToken() {
		case rhino.Token_BLOCK, rhino.Token_COMPUTED_PROP, rhino.Token_COMPUTED_FIELD_DEF, rhino.Token_OBJECT_REST, rhino.Token_OBJECT_SPREAD:
			continue
		case rhino.Token_STRING_KEY, rhino.Token_MEMBER_FUNCTION_DEF, rhino.Token_MEMBER_FIELD_DEF, rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF:
			if child.IsQuotedStringKey() {
				continue
			}
			var owner *colors.Color
			if child.IsStaticMember() {
				owner = classType
			} else if child.IsMemberFieldDef() {
				if instanceOwner == nil {
					classInstanceType := classType.GetInstanceColors()
					if classInstanceType.IsEmpty() {
						instanceOwner = colors.StandardColors_UNKNOWN
					} else {
						instanceOwner = colors.Color_CreateUnion(classInstanceType)
					}
				}
				owner = instanceOwner
			} else {
				jutil.CheckStateMsg1(child.IsMemberFunctionDef() || child.IsGetterDef() || child.IsSetterDef(), "%s", child)
				if protoOwner == nil {
					classPrototypeType := classType.GetPrototypes()
					if classPrototypeType.IsEmpty() {
						protoOwner = colors.StandardColors_UNKNOWN
					} else {
						protoOwner = colors.Color_CreateUnion(classPrototypeType)
					}
				}
				owner = protoOwner
			}
			r.registerPropertyUse(t, child, owner)
		default:
			panic(jutil.NewIllegalStateExceptionMsg(
				"Unexpected child of " + members.GetToken().String() + ": " + child.ToStringTree()))
		}
	}
	r.createNode(classType)
}

// Java: ColorFindPropertyReferences.handleFunction(Node)
func (r *ColorFindPropertyReferences) handleFunction(fnNode rhino.Node) {
	// Ensure the flattener knows about the class constructor and instance. Even if we don't see any
	// direct property references off either type, it's possible that a property is referenced
	// via a superclass or implemented interface.
	// e.g. `const /** !FooInterface */ x = new Foo(); x.method();`
	fnType := fnNode.GetColor()
	if fnType != nil && !fnType.GetInstanceColors().IsEmpty() {
		r.createNode(fnType)
	}
}

// Java: ColorFindPropertyReferences.handleObjectPattern(NodeTraversal, Node)
func (r *ColorFindPropertyReferences) handleObjectPattern(t *NodeTraversal, pattern rhino.Node) {
	c := pattern.GetColor()
	r.traverseObjectlitLike(t, pattern, func(rhino.Node) *colors.Color { return c })
}

// Java: ColorFindPropertyReferences.handlePropertyReflectorCall(NodeTraversal, Node)
func (r *ColorFindPropertyReferences) handlePropertyReflectorCall(t *NodeTraversal, call rhino.Node) {
	name := call.GetSecondChild()
	if name.IsNil() || !name.IsStringLit() {
		return
	}

	obj := name.GetNext()
	var objColor *colors.Color
	if !obj.IsNil() {
		objColor = obj.GetColor()
	}
	r.registerPropertyUse(t, name, objColor)
}

// Java: ColorFindPropertyReferences.handleObjectDefineProperties(NodeTraversal, Node)
func (r *ColorFindPropertyReferences) handleObjectDefineProperties(t *NodeTraversal, call rhino.Node) {
	typeObj := call.GetSecondChild()
	objectLiteral := typeObj.GetNext()
	if !objectLiteral.IsObjectLit() {
		return
	}

	c := typeObj.GetColor()
	r.traverseObjectlitLike(t, objectLiteral, func(rhino.Node) *colors.Color { return c })
}

// Java: ColorFindPropertyReferences.traverseObjectlitLike(NodeTraversal, Node, Function<Node, Color>)
func (r *ColorFindPropertyReferences) traverseObjectlitLike(t *NodeTraversal, n rhino.Node, memberOwnerFn func(rhino.Node) *colors.Color) {
	// The keys in an object pattern are r-values, not l-values, but they are still accesses.
	jutil.CheckState(n.IsObjectLit() || n.IsObjectPattern() || n.IsClassMembers())

	for child := n.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		switch child.GetToken() {
		case rhino.Token_BLOCK, rhino.Token_COMPUTED_PROP, rhino.Token_COMPUTED_FIELD_DEF, rhino.Token_OBJECT_REST, rhino.Token_OBJECT_SPREAD:
			continue
		case rhino.Token_STRING_KEY, rhino.Token_MEMBER_FUNCTION_DEF, rhino.Token_MEMBER_FIELD_DEF, rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF:
			if child.IsQuotedStringKey() {
				continue // These won't be renamed due to our assumptions. Ignore them.
			}

			r.registerPropertyUse(t, child, memberOwnerFn(child))
		default:
			panic(jutil.NewIllegalStateExceptionMsg(
				"Unexpected child of " + n.GetToken().String() + ": " + child.ToStringTree()))
		}
	}
}

// registerPropertyUse updates all datastructures as necessary to consider
// property use `site` from type `owner`.
//
// Java: ColorFindPropertyReferences.registerPropertyUse(NodeTraversal, Node, Color)
func (r *ColorFindPropertyReferences) registerPropertyUse(t *NodeTraversal, site rhino.Node, owner *colors.Color) {
	prop := r.propIndex.ComputeIfAbsent(site.GetString(), NewPropertyClustering)
	flatOwner := r.createNode(owner)

	// Set the initial condition for flowing this property along the graph.
	if _, existed := flatOwner.GetAssociatedProps().LoadOrStore(prop, ColorGraphNode_PropAssociation_AST); !existed {
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
