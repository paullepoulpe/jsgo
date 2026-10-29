// Port of com/google/javascript/jscomp/NodeUtil.java @ closure-compiler v20261021

/*
 * Copyright 2004 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// NodeUtil_IsNormalOrOptChainGet: is this a GETPROP, OPTCHAIN_GETPROP,
// GETELEM, or OPTCHAIN_GETELEM?
//
// Java: NodeUtil.isNormalOrOptChainGet(Node)
func NodeUtil_IsNormalOrOptChainGet(n rhino.Node) bool {
	return NodeUtil_IsNormalGet(n) || NodeUtil_IsOptChainGet(n)
}

// NodeUtil_IsNormalOrOptChainGetProp: is this a GETPROP or OPTCHAIN_GETPROP?
//
// Java: NodeUtil.isNormalOrOptChainGetProp(Node)
func NodeUtil_IsNormalOrOptChainGetProp(n rhino.Node) bool {
	return n.IsGetProp() || n.IsOptChainGetProp()
}

// NodeUtil_IsNormalOrOptChainCall: is this a CALL or OPTCHAIN_CALL?
//
// Java: NodeUtil.isNormalOrOptChainCall(Node)
func NodeUtil_IsNormalOrOptChainCall(n rhino.Node) bool {
	return n.IsCall() || n.IsOptChainCall()
}

// NodeUtil_IsNormalGet: is this a GETPROP or GETELEM node?
//
// Java: NodeUtil.isNormalGet(Node)
func NodeUtil_IsNormalGet(n rhino.Node) bool {
	return n.IsGetProp() || n.IsGetElem()
}

// NodeUtil_IsOptChainGet: is this an OPTCHAIN_GETPROP or OPTCHAIN_GETELEM
// node?
//
// Java: NodeUtil.isOptChainGet(Node)
func NodeUtil_IsOptChainGet(n rhino.Node) bool {
	return n.IsOptChainGetProp() || n.IsOptChainGetElem()
}

// NodeUtil_IsOptChainNode: is this a OPTCHAIN_GETPROP, OPTCHAIN_GETELEM,
// OPTCHAIN_CALL node?
//
// Java: NodeUtil.isOptChainNode(Node)
func NodeUtil_IsOptChainNode(n rhino.Node) bool {
	return n.IsOptChainGetProp() || n.IsOptChainGetElem() || n.IsOptChainCall()
}

// NodeUtil_GetStartOfOptChainSegment finds the start of the optional chain.
// E.g Find the `a?. ...` node in `a?.b.c.d` given any other node.
//
// n is a node in an optional chain; returns the start of the optional chain
// that n is part of.
//
// Java: NodeUtil.getStartOfOptChainSegment(Node)
func NodeUtil_GetStartOfOptChainSegment(n rhino.Node) rhino.Node {
	jutil.CheckStateMsg1(NodeUtil_IsOptChainNode(n), "%s", n)
	if n.IsOptionalChainStart() {
		return n
	}
	return NodeUtil_GetStartOfOptChainSegment(n.GetFirstChild())
}

// NodeUtil_GetEndOfOptChainSegment finds the end of an optional chain
// segment.
//
// Each `?.` ends one segment and starts another.
//
// Examples
//
//	a?.b.c.d   // end is the node with children `a?.b.c` and `d`
//	a?.b.c?.d  // given a?.b end is the node with children `a?.b` and `c`
//	a?.b.c?.d  // given a?.b.c end is the node with children `a?.b.c` and `d`
//	(a?.b.c).d // given a?.b end is the node with children `a?.b` and `c`
//
// n is a node in an optional chain; returns the end of the optional chain
// that n is part of.
//
// Java: NodeUtil.getEndOfOptChainSegment(Node)
func NodeUtil_GetEndOfOptChainSegment(n rhino.Node) rhino.Node {
	jutil.CheckStateMsg1(NodeUtil_IsOptChainNode(n), "%s", n)
	if NodeUtil_IsEndOfOptChainSegment(n) {
		return n
	} else {
		return NodeUtil_GetEndOfOptChainSegment(n.GetParent())
	}
}

// NodeUtil_IsEndOfFullOptChain: is this node the final node of a full
// optional chain?
//
// e.g. for `a?.b.c?.d` this method returns true only for the Node with
// children `a?.b.c` and `d`. That node is the end of the whole chain, and
// also represents the whole chain in the AST.
//
// Java: NodeUtil.isEndOfFullOptChain(Node)
func NodeUtil_IsEndOfFullOptChain(n rhino.Node) bool {
	if NodeUtil_IsOptChainNode(n) {
		parent := n.GetParent()
		// the chain continues if this is the first child of another optional chain node
		return !(NodeUtil_IsOptChainNode(parent) && n.IsFirstChildOf(parent))
	} else {
		// not even an optional chain node
		return false
	}
}

// NodeUtil_IsEndOfOptChainSegment: is this node the end of an optional chain
// segment?
//
// Each `?.` begins a new segment and ends the previous one, if any. The end
// of the whole chain is also the end of its final segment.
//
// Java: NodeUtil.isEndOfOptChainSegment(Node)
func NodeUtil_IsEndOfOptChainSegment(n rhino.Node) bool {
	if !NodeUtil_IsOptChainNode(n) {
		return false
	} else {
		parent := n.GetParent()
		// Check for null so this method will work for a disconnected Node.
		if !parent.IsNil() && n.IsFirstChildOf(parent) && NodeUtil_IsOptChainNode(parent) {
			// The parent is a continuation of this optional chain.
			// If it starts a new segment, then this node is the end of a segment
			return parent.IsOptionalChainStart()
		} else {
			// The parent doesn't continue this node's optional chain, though it might be part of a
			// different one.
			// e.g. in `a?.(x?.y.z)` the parent of `x?.y.z` is part of a different optional chain.
			return true
		}
	}
}

// NodeUtil_ConvertToNonOptionalChainSegment: given the end of an optional
// chain segment changes all nodes from the end down to the start into
// non-optional nodes. e.g `({a})?.a.b.c.d()?.x.y.z` gets converted to
// `({a}).a.b.c.d()?.x.y.z` when the passed in endOfOptChainSegment is
// `({a})?.a.b.c.d()`.
//
// Java: NodeUtil.convertToNonOptionalChainSegment(Node)
func NodeUtil_ConvertToNonOptionalChainSegment(endOfOptChainSegment rhino.Node) {
	start := NodeUtil_GetStartOfOptChainSegment(endOfOptChainSegment)
	NodeUtil_ConvertToNonOptionalChainSegmentDownTo(endOfOptChainSegment, start.GetFirstChild())
}

// NodeUtil_ConvertToNonOptionalChainSegmentDownTo: given the end of an
// optional chain segment, changes all nodes from endOfOptChainSegment down
// to (but excluding) stopNode into non-optional nodes.
//
// Java: NodeUtil.convertToNonOptionalChainSegmentDownTo(Node, Node)
func NodeUtil_ConvertToNonOptionalChainSegmentDownTo(endOfOptChainSegment rhino.Node, stopNode rhino.Node) {
	jutil.CheckArgumentMsg1(NodeUtil_IsEndOfOptChainSegment(endOfOptChainSegment), "%s", endOfOptChainSegment)
	var segmentNodes []rhino.Node
	segmentNode := endOfOptChainSegment
	for !segmentNode.IsNil() && segmentNode != stopNode {
		if NodeUtil_IsOptChainNode(segmentNode) {
			segmentNodes = append(segmentNodes, segmentNode)
		}
		segmentNode = segmentNode.GetFirstChild()
	}
	for _, n := range segmentNodes {
		n.SetIsOptionalChainStart(false)
		n.SetToken(NodeUtil_GetNonOptChainToken(n.GetToken()))
	}
}

// Java: NodeUtil.getNonOptChainToken(Token)
func NodeUtil_GetNonOptChainToken(optChainToken rhino.Token) rhino.Token {
	switch optChainToken {
	case rhino.Token_OPTCHAIN_CALL:
		return rhino.Token_CALL
	case rhino.Token_OPTCHAIN_GETELEM:
		return rhino.Token_GETELEM
	case rhino.Token_OPTCHAIN_GETPROP:
		return rhino.Token_GETPROP
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Should be an OPTCHAIN token: " + optChainToken.String()))
	}
}

// NodeUtil_IsNameDeclOrSimpleAssignLhs determines whether this node is
// strictly on the left hand side of an assign or var initialization.
// Notably, this does not include all L-values, only statements where the
// node is used only as an L-value.
//
// n is the node, parent the parent of the node. Returns true if n is the
// left hand of an assign.
//
// Java: NodeUtil.isNameDeclOrSimpleAssignLhs(Node, Node)
func NodeUtil_IsNameDeclOrSimpleAssignLhs(n rhino.Node, parent rhino.Node) bool {
	return ast.RoleIn(n, parent) == ast.RoleAssignTarget || NodeUtil_IsNameDeclaration(parent)
}

// NodeUtil_MayBeObjectLitKey determines whether a node represents a possible
// object literal key (e.g. key1 in {key1: value1, key2: value2}). Computed
// properties are excluded here (see b/111621528). This method does not check
// whether the node is actually in an object literal! it also returns true
// for object pattern keys, and member functions/getters in ES6 classes.
//
// Java: NodeUtil.mayBeObjectLitKey(Node)
//
// TODO(b/189993301): should fields be added to this method?
func NodeUtil_MayBeObjectLitKey(node rhino.Node) bool {
	switch node.GetToken() {
	case rhino.Token_STRING_KEY, rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF, rhino.Token_MEMBER_FUNCTION_DEF:
		return true
	default:
		return false
	}
}

// NodeUtil_IsObjectLitKey determines whether a node represents an object
// literal key (e.g. key1 in {key1: value1, key2: value2}) and is in an
// object literal. Computed properties are excluded here (see b/111621528).
//
// Java: NodeUtil.isObjectLitKey(Node)
func NodeUtil_IsObjectLitKey(node rhino.Node) bool {
	return node.GetParent().IsObjectLit() && NodeUtil_MayBeObjectLitKey(node)
}

// NodeUtil_GetObjectOrClassLitKeyName gets the name of an object literal
// key.
//
// Java: NodeUtil.getObjectOrClassLitKeyName(Node)
func NodeUtil_GetObjectOrClassLitKeyName(key rhino.Node) string {
	keyNode := NodeUtil_GetObjectOrClassLitKeyNode(key)
	if !keyNode.IsNil() {
		return keyNode.GetString()
	}
	panic(jutil.NewIllegalStateExceptionMsg("Unexpected node type: " + key.String()))
}

// NodeUtil_GetObjectOrClassLitKeyNode gets the Node that defines the name of
// an object literal key.
//
// Java: NodeUtil.getObjectOrClassLitKeyNode(Node)
func NodeUtil_GetObjectOrClassLitKeyNode(key rhino.Node) rhino.Node {
	switch key.GetToken() {
	case rhino.Token_STRING_KEY, rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF, rhino.Token_MEMBER_FUNCTION_DEF,
		rhino.Token_MEMBER_FIELD_DEF:
		return key
	case rhino.Token_COMPUTED_PROP, rhino.Token_COMPUTED_FIELD_DEF:
		if key.GetFirstChild().
			IsStringLit() { // TODO(b/189993301): may be an issue with non string lits
			return key.GetFirstChild()
		}
		return rhino.NilNode
	default:
	}
	panic(jutil.NewIllegalStateExceptionMsg("Unexpected node type: " + key.String()))
}

// NodeUtil_IsGetOrSetKey determines whether a node represents an object
// literal get or set key (e.g. key1 in {get key1() {}, set key2(a){}).
//
// Java: NodeUtil.isGetOrSetKey(Node)
func NodeUtil_IsGetOrSetKey(node rhino.Node) bool {
	if _, ok := ast.AsGetterOrSetterDef(node); ok {
		return true
	}
	if comp, ok := ast.AsComputedProp(node); ok {
		return comp.IsGetter() || comp.IsSetter()
	}
	return false
}

// NodeUtil_OpToStr converts an operator's token value (see Token) to a
// string representation.
//
// operator is the operator's token value to convert; returns the string
// representation or null if the token value is not an operator.
//
// Java: NodeUtil.opToStr(Token)
func NodeUtil_OpToStr(operator rhino.Token) *string {
	// Go-only: language-adapter — Java returns interned string constants; the
	// results of the switch below are computed once per token so that the
	// *string result does not allocate on every call (the code printer calls
	// this per operator).
	return nodeUtil_opToStrTable[operator]
}

// NodeUtil_IsObjectDefinePropertiesDefinition returns true if the node is a
// definition with Object.defineProperties.
//
// Java: NodeUtil.isObjectDefinePropertiesDefinition(Node)
func NodeUtil_IsObjectDefinePropertiesDefinition(n rhino.Node) bool {
	// We intentionally don't check optional Object.defineProperties?.() because we expect it to be
	// defined or polyfill-ed, and to avoid changing the current invariant that this method returns
	// true only for non-optional `Object.defineProperties()`.
	call, ok := ast.AsCall(n)
	if !ok || !n.HasXChildren(3) {
		return false
	}
	getprop, ok := ast.AsGetProp(call.Callee().Node())
	return ok &&
		getprop.Property() == "defineProperties" &&
		NodeUtil_IsKnownGlobalObjectReference(getprop.Object().Node())
}

// NodeUtil_IsObjectCreateWithPropertiesDefinition returns true if the node is
// a 2-argument Object.create(proto, descriptors).
//
// Java: NodeUtil.isObjectCreateWithPropertiesDefinition(Node)
func NodeUtil_IsObjectCreateWithPropertiesDefinition(n rhino.Node) bool {
	call, ok := ast.AsCall(n)
	if !ok || !n.HasXChildren(3) {
		return false
	}
	getprop, ok := ast.AsGetProp(call.Callee().Node())
	return ok &&
		getprop.Property() == "create" &&
		NodeUtil_IsKnownGlobalObjectReference(getprop.Object().Node())
}

// Java: NodeUtil.isKnownGlobalObjectReference(Node)
func NodeUtil_IsKnownGlobalObjectReference(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_NAME:
		return n.GetString() == "Object"
	case rhino.Token_GETPROP:
		return NodeUtil_GLOBAL_OBJECT.Matches(n) || NodeUtil_GLOBAL_OBJECT_MANGLED.Matches(n)
	default:
		return false
	}
}

// NodeUtil_IsObjectDefinePropertyDefinition returns true if the node is a
// definition with Object.defineProperty.
//
// Java: NodeUtil.isObjectDefinePropertyDefinition(Node)
func NodeUtil_IsObjectDefinePropertyDefinition(n rhino.Node) bool {
	call, ok := ast.AsCall(n)
	if !ok || !n.HasXChildren(4) {
		return false
	}
	getprop, ok := ast.AsGetProp(call.Callee().Node())
	return ok &&
		getprop.Property() == "defineProperty" &&
		NodeUtil_IsKnownGlobalObjectReference(getprop.Object().Node())
}

// NodeUtil_GetObjectDefinedPropertiesKeys returns a list of STRING_KEY
// properties defined by a Object.defineProperties(o, {...}) call.
//
// Java: NodeUtil.getObjectDefinedPropertiesKeys(Node)
func NodeUtil_GetObjectDefinedPropertiesKeys(definePropertiesCall rhino.Node) []rhino.Node {
	jutil.CheckArgument(NodeUtil_IsObjectDefinePropertiesDefinition(definePropertiesCall))
	var properties []rhino.Node
	if obj, ok := ast.AsObjectLit(definePropertiesCall.GetLastChild()); ok {
		for member := range obj.Members() {
			if key, ok := member.AsStringKey(); ok {
				properties = append(properties, key.Node())
			}
		}
	}
	return properties
}

// NodeUtil_IsPrototypePropertyDeclaration returns true if the node an
// assignment to a prototype property of some constructor.
//
// Java: NodeUtil.isPrototypePropertyDeclaration(Node)
func NodeUtil_IsPrototypePropertyDeclaration(n rhino.Node) bool {
	return NodeUtil_IsExprAssign(n) && NodeUtil_IsPrototypeProperty(n.GetFirstFirstChild())
}

// NodeUtil_IsPrototypeProperty returns whether the node represents a
// qualified prototype property.
//
// Java: NodeUtil.isPrototypeProperty(Node)
func NodeUtil_IsPrototypeProperty(n rhino.Node) bool {
	prop, ok := ast.AsGetProp(n)
	if !ok {
		return false
	}
	recv, ok := ast.AsGetProp(prop.Object().Node())
	return ok && recv.Property() == "prototype"
}

// NodeUtil_IsPrototypeMethod returns whether the node represents a prototype
// method.
//
// Java: NodeUtil.isPrototypeMethod(Node)
func NodeUtil_IsPrototypeMethod(n rhino.Node) bool {
	if !n.IsFunction() {
		return false
	}
	assignNode := n.GetParent()
	if !assignNode.IsAssign() {
		return false
	}
	return NodeUtil_IsPrototypePropertyDeclaration(assignNode.GetParent())
}

// Java: NodeUtil.isPrototypeAssignment(Node)
func NodeUtil_IsPrototypeAssignment(getProp rhino.Node) bool {
	prop, ok := ast.AsGetProp(getProp)
	return ok &&
		ast.RoleOf(getProp) == ast.RoleAssignTarget &&
		prop.Property() == "prototype"
}

// NodeUtil_IsPropertyTest determines whether this node is testing for the
// existence of a property. If true, we will not emit warnings about a
// missing property.
//
// propAccess is the GETPROP or GETELEM being tested.
//
// Java: NodeUtil.isPropertyTest(AbstractCompiler, Node)
func NodeUtil_IsPropertyTest(compiler AbstractCompiler, propAccess rhino.Node) bool {
	parent := propAccess.GetParent()
	switch parent.GetToken() {
	case rhino.Token_CALL:
		return parent.GetFirstChild() != propAccess &&
			compiler.GetCodingConvention().IsPropertyTestFunction(parent)
	case rhino.Token_IF, rhino.Token_WHILE, rhino.Token_DO, rhino.Token_FOR, rhino.Token_FOR_IN:
		return NodeUtil_GetConditionExpression(parent) == propAccess

	case rhino.Token_INSTANCEOF, rhino.Token_TYPEOF, rhino.Token_AND, rhino.Token_OR, rhino.Token_COALESCE,
		rhino.Token_OPTCHAIN_GETPROP:
		return true

	case rhino.Token_NE, rhino.Token_SHNE:
		var other rhino.Node
		if parent.GetFirstChild() == propAccess {
			other = parent.GetSecondChild()
		} else {
			other = parent.GetFirstChild()
		}
		return NodeUtil_IsUndefined(other) || (parent.IsNE() && other.IsNull())
	case rhino.Token_OPTCHAIN_CALL, rhino.Token_OPTCHAIN_GETELEM, rhino.Token_HOOK:
		return parent.GetFirstChild() == propAccess
	case rhino.Token_NOT:
		return parent.GetParent().IsOr() && parent.GetParent().GetFirstChild() == parent
	case rhino.Token_CAST:
		return NodeUtil_IsPropertyTest(compiler, parent)
	default:
		return false
	}
}

// Java: NodeUtil.isPropertyAbsenceTest(Node)
func NodeUtil_IsPropertyAbsenceTest(propAccess rhino.Node) bool {
	parent := propAccess.GetParent()
	switch parent.GetToken() {
	case rhino.Token_EQ, rhino.Token_SHEQ:
		var other rhino.Node
		if parent.GetFirstChild() == propAccess {
			other = parent.GetSecondChild()
		} else {
			other = parent.GetFirstChild()
		}
		return NodeUtil_IsUndefined(other) || (parent.IsEQ() && other.IsNull())
	default:
		return false
	}
}

// NodeUtil_GetPrototypeClassName: qName is a qualified name node
// representing a class prototype, or a property on that prototype, e.g.
// foo.Bar.prototype, or foo.Bar.prototype.toString. Returns the class name
// part of a qualified prototype name, e.g. foo.Bar.
//
// Java: NodeUtil.getPrototypeClassName(Node)
func NodeUtil_GetPrototypeClassName(qName rhino.Node) rhino.Node {
	prop, ok := ast.AsGetProp(qName)
	if !ok {
		return rhino.NilNode
	}
	if prop.Property() == "prototype" {
		return prop.Object().Node()
	}
	if recv, ok := ast.AsGetProp(prop.Object().Node()); ok && recv.Property() == "prototype" {
		return recv.Object().Node()
	}
	return rhino.NilNode
}

// NodeUtil_GetPrototypePropertyName returns the string property name part of
// a qualified prototype name.
//
// Java: NodeUtil.getPrototypePropertyName(Node)
func NodeUtil_GetPrototypePropertyName(qName rhino.Node) string {
	qNameStr, ok := qName.GetQualifiedNameString()
	if !ok {
		panic(jutil.NewNullPointerException())
	}
	prototypeIdx := jutil.StringLastIndexOf(qNameStr, ".prototype.")
	memberIndex := prototypeIdx + jutil.StringLength(".prototype") + 1
	return jutil.StringSubstringFrom(qNameStr, memberIndex)
}

// Java: NodeUtil.isGetterOrSetter(Node)
func NodeUtil_IsGetterOrSetter(propNode rhino.Node) bool {
	if NodeUtil_IsGetOrSetKey(propNode) {
		return true
	}
	key, ok := ast.AsStringKey(propNode)
	if !ok || !key.Value().Node().IsFunction() {
		return false
	}
	keyName := key.Key()
	return keyName == "get" || keyName == "set"
}

// NodeUtil_IsObjLitProperty returns true if the node is a property of an
// object literal.
//
// Java: NodeUtil.isObjLitProperty(Node)
func NodeUtil_IsObjLitProperty(node rhino.Node) bool {
	return node.IsStringKey() ||
		node.IsGetterDef() ||
		node.IsSetterDef() ||
		node.IsMemberFunctionDef() ||
		node.IsComputedProp()
}
