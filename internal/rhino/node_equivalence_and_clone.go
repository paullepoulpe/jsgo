// Port of com/google/javascript/rhino/Node.java @ closure-compiler v20261021

/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * The Original Code is Rhino code, released
 * May 6, 1999.
 *
 * The Initial Developer of the Original Code is
 * Netscape Communications Corporation.
 * Portions created by the Initial Developer are Copyright (C) 1997-1999
 * the Initial Developer. All Rights Reserved.
 *
 * Contributor(s):
 *   Norris Boyd
 *   Roger Lawrence
 *   Mike McCabe
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

package rhino

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Node_RecursionMode is the recursion mode for isEquivalentTo.
//
// Java: Node.RecursionMode
type Node_RecursionMode int32

const (
	// Do not compare children or shadow AST.
	Node_RecursionMode_SHALLOW Node_RecursionMode = iota
	// Recursively compare children but not shadow AST.
	Node_RecursionMode_DEEP_NO_SHADOW
	// Recursively compare children and shadow AST.
	Node_RecursionMode_DEEP
)

var node_RecursionMode_names = [...]string{"SHALLOW", "DEEP_NO_SHADOW", "DEEP"}

// String is Enum.name().
func (m Node_RecursionMode) String() string { return node_RecursionMode_names[m] }

// Java: Node.RecursionMode.recurseShadowAst()
func (m Node_RecursionMode) RecurseShadowAst() bool {
	return m == Node_RecursionMode_DEEP
}

// Java: Node.RecursionMode.recurseChildren()
func (m Node_RecursionMode) RecurseChildren() bool {
	return m == Node_RecursionMode_DEEP || m == Node_RecursionMode_DEEP_NO_SHADOW
}

// Node_TypeComparison is the type comparison mode for isEquivalentTo.
//
// Java: Node.TypeComparison
type Node_TypeComparison int32

const (
	// Ignore JSTypes/Colors.
	Node_TypeComparison_IGNORE Node_TypeComparison = iota
	// Compare JSTypes/Colors.
	Node_TypeComparison_COMPARE
)

// String is Enum.name().
func (c Node_TypeComparison) String() string { return [...]string{"IGNORE", "COMPARE"}[c] }

// Node_JsDocComparison is the JsDoc comparison mode for isEquivalentTo.
//
// Java: Node.JsDocComparison
type Node_JsDocComparison int32

const (
	// Ignore JSDoc.
	Node_JsDocComparison_IGNORE Node_JsDocComparison = iota
	// Compare JSDoc.
	Node_JsDocComparison_COMPARE
)

// String is Enum.name().
func (c Node_JsDocComparison) String() string { return [...]string{"IGNORE", "COMPARE"}[c] }

// Node_SideEffectComparison is the side effect comparison mode for
// isEquivalentTo.
//
// Java: Node.SideEffectComparison
type Node_SideEffectComparison int32

const (
	// Ignore side effect flags.
	Node_SideEffectComparison_IGNORE Node_SideEffectComparison = iota
	// Compare side effect flags.
	Node_SideEffectComparison_COMPARE
)

// String is Enum.name().
func (c Node_SideEffectComparison) String() string { return [...]string{"IGNORE", "COMPARE"}[c] }

// IsEquivalentToShallow checks equivalence without going into child nodes.
//
// Java: Node.isEquivalentToShallow(Node)
func (n Node) IsEquivalentToShallow(node Node) bool {
	return n.IsEquivalentToWithModes(
		node,
		Node_RecursionMode_SHALLOW,
		Node_TypeComparison_IGNORE,
		Node_JsDocComparison_IGNORE,
		Node_SideEffectComparison_IGNORE)
}

// IsEquivalentWithSideEffectsTo returns true if this node is equivalent
// semantically to another including side effects.
//
// Java: Node.isEquivalentWithSideEffectsTo(Node)
func (n Node) IsEquivalentWithSideEffectsTo(node Node) bool {
	return n.IsEquivalentToWithModes(
		node,
		Node_RecursionMode_DEEP,
		Node_TypeComparison_IGNORE,
		Node_JsDocComparison_IGNORE,
		Node_SideEffectComparison_COMPARE)
}

// IsEquivalentWithSideEffectsToShallow returns true if this node is
// equivalent semantically to another including side effects.
//
// Java: Node.isEquivalentWithSideEffectsToShallow(Node)
func (n Node) IsEquivalentWithSideEffectsToShallow(node Node) bool {
	return n.IsEquivalentToWithModes(
		node,
		Node_RecursionMode_SHALLOW,
		Node_TypeComparison_IGNORE,
		Node_JsDocComparison_IGNORE,
		Node_SideEffectComparison_COMPARE)
}

// IsEquivalentToTyped returns true if this node is equivalent semantically
// to another and the types are equivalent.
//
// Java: Node.isEquivalentToTyped(Node)
func (n Node) IsEquivalentToTyped(node Node) bool {
	return n.IsEquivalentToWithModes(
		node,
		Node_RecursionMode_DEEP,
		Node_TypeComparison_COMPARE,
		Node_JsDocComparison_COMPARE,
		Node_SideEffectComparison_IGNORE)
}

// IsEquivalentTo returns true if this node is equivalent semantically to
// another.
//
// Java: Node.isEquivalentTo(Node)
func (n Node) IsEquivalentTo(node Node) bool {
	return n.IsEquivalentToWithModes(
		node,
		Node_RecursionMode_DEEP,
		Node_TypeComparison_IGNORE,
		Node_JsDocComparison_IGNORE,
		Node_SideEffectComparison_IGNORE)
}

// IsEquivalentToWithFlags returns whether this node is equivalent
// semantically to the provided node.
//
// compareType: whether to compare the JSTypes of the nodes. recurse: whether
// to compare the children of the current node. If not, only the count of
// the children are compared. jsDoc: whether to check that the JsDoc of the
// nodes are equivalent.
//
// Deprecated: Use IsEquivalentToWithModes instead.
//
// Java: Node.isEquivalentTo(Node, boolean, boolean, boolean)
func (n Node) IsEquivalentToWithFlags(node Node, compareType bool, recurse bool, jsDoc bool) bool {
	return n.IsEquivalentToWithFlagsAndSideEffect(node, compareType, recurse, jsDoc, false)
}

// IsEquivalentToWithFlagsAndSideEffect returns whether this node is
// equivalent semantically to the provided node.
//
// sideEffect: whether to check that the side-effect flags of the nodes are
// equivalent.
//
// Deprecated: Use IsEquivalentToWithModes instead.
//
// Java: Node.isEquivalentTo(Node, boolean, boolean, boolean, boolean)
func (n Node) IsEquivalentToWithFlagsAndSideEffect(node Node, compareType bool, recurse bool, jsDoc bool, sideEffect bool) bool {
	recursion := Node_RecursionMode_SHALLOW
	if recurse {
		recursion = Node_RecursionMode_DEEP
	}
	typeConfig := Node_TypeComparison_IGNORE
	if compareType {
		typeConfig = Node_TypeComparison_COMPARE
	}
	jsDocConfig := Node_JsDocComparison_IGNORE
	if jsDoc {
		jsDocConfig = Node_JsDocComparison_COMPARE
	}
	sideEffectConfig := Node_SideEffectComparison_IGNORE
	if sideEffect {
		sideEffectConfig = Node_SideEffectComparison_COMPARE
	}
	return n.IsEquivalentToWithModes(node, recursion, typeConfig, jsDocConfig, sideEffectConfig)
}

// IsEquivalentToWithModes returns whether this node is equivalent
// semantically to the provided node.
//
// Go-only: language-adapter — dispatches to the overrides of the folded
// subclasses (NumberNode, BigIntNode, StringNode,
// TemplateLiteralSubstringNode), each of which calls the base
// implementation first.
//
// Java: Node.isEquivalentTo(Node, RecursionMode, TypeComparison, JsDocComparison, SideEffectComparison)
func (n Node) IsEquivalentToWithModes(
	node Node,
	recursion Node_RecursionMode,
	typeConfig Node_TypeComparison,
	jsDocConfig Node_JsDocComparison,
	sideEffectConfig Node_SideEffectComparison,
) bool {
	switch n.nodeKind() {
	case node_kindNumber:
		// Java: NumberNode.isEquivalentTo
		return n.isEquivalentToBase(node, recursion, typeConfig, jsDocConfig, sideEffectConfig) &&
			(n.node_number() == node.node_number()) // -0.0 and NaN are forbidden.
	case node_kindBigInt:
		// Java: BigIntNode.isEquivalentTo
		return n.isEquivalentToBase(node, recursion, typeConfig, jsDocConfig, sideEffectConfig) &&
			jutil.BigIntegerEquals(n.GetBigInt(), node.GetBigInt())
	case node_kindString:
		// Java: StringNode.isEquivalentTo
		return n.isEquivalentToBase(node, recursion, typeConfig, jsDocConfig, sideEffectConfig) &&
			uint32(n.rec().payload) == uint32(node.rec().payload)
	case node_kindTemplateLiteralSubstring:
		// Java: TemplateLiteralSubstringNode.isEquivalentTo
		if !n.isEquivalentToBase(node, recursion, typeConfig, jsDocConfig, sideEffectConfig) {
			return false
		}
		return n.rec().payload == node.rec().payload && n.node_cookedID() == node.node_cookedID()
	}
	return n.isEquivalentToBase(node, recursion, typeConfig, jsDocConfig, sideEffectConfig)
}

// isEquivalentToBase is Node's own isEquivalentTo (the super call of the
// folded subclasses).
func (n Node) isEquivalentToBase(
	node Node,
	recursion Node_RecursionMode,
	typeConfig Node_TypeComparison,
	jsDocConfig Node_JsDocComparison,
	sideEffectConfig Node_SideEffectComparison,
) bool {
	if n.rec().token != node.rec().token ||
		n.GetChildCount() != node.GetChildCount() ||
		n.nodeKind() != node.nodeKind() {
		return false
	}

	if typeConfig == Node_TypeComparison_COMPARE &&
		!jutil.ObjectsEquals(n.typeOrColorValue(), node.typeOrColorValue()) {
		return false
	}

	if jsDocConfig == Node_JsDocComparison_COMPARE &&
		!JSDocInfo_AreEquivalent(n.GetJSDocInfo(), node.GetJSDocInfo()) {
		return false
	}

	thisDte := n.GetDeclaredTypeExpression()
	thatDte := node.GetDeclaredTypeExpression()
	if thisDte == thatDte {
		// Do nothing
	} else if thisDte.IsNil() || thatDte.IsNil() {
		return false
	} else if !thisDte.IsEquivalentToWithModes(
		thatDte, recursion, typeConfig, jsDocConfig, sideEffectConfig) {
		return false
	}

	if (n.props() & node_propEqualityMask) != (node.props() & node_propEqualityMask) {
		return false
	}
	if ((n.rec().kindFlags | node.rec().kindFlags) & kindFlagHasExtra) != 0 {
		if n.GetIntProp(Node_Prop_EXPORT_DEFAULT) != node.GetIntProp(Node_Prop_EXPORT_DEFAULT) ||
			n.GetIntProp(Node_Prop_EXPORT_ALL_FROM) != node.GetIntProp(Node_Prop_EXPORT_ALL_FROM) ||
			n.GetIntProp(Node_Prop_INCRDECR) != node.GetIntProp(Node_Prop_INCRDECR) ||
			n.GetIntProp(Node_Prop_QUOTED) != node.GetIntProp(Node_Prop_QUOTED) {
			return false
		}
	}

	if sideEffectConfig == Node_SideEffectComparison_COMPARE {
		if n.GetSideEffectFlags() != node.GetSideEffectFlags() {
			return false
		}

		if n.IsUnusedParameter() != node.IsUnusedParameter() {
			return false
		}
	}

	thisIsShadowHost := !n.GetClosureUnawareShadow().IsNil()
	nodeIsShadowHost := !node.GetClosureUnawareShadow().IsNil()
	if thisIsShadowHost != nodeIsShadowHost {
		return false
	}

	if recursion.RecurseShadowAst() && thisIsShadowHost && nodeIsShadowHost {
		thisShadow := n.GetClosureUnawareShadow()
		nodeShadow := node.GetClosureUnawareShadow()
		if !thisShadow.IsEquivalentToWithModes(
			nodeShadow, recursion, typeConfig, jsDocConfig, sideEffectConfig) {
			return false
		}
	}

	if recursion.RecurseChildren() {
		for c, c2 := n.firstNode(), node.firstNode(); !c.IsNil(); c, c2 = c.nextNode(), c2.nextNode() {
			if !c.IsEquivalentToWithModes(c2, recursion, typeConfig, jsDocConfig, sideEffectConfig) {
				return false
			}
		}
	}

	return true
}

const node_propEqualityMask uint64 = propFlag_ARROW_FN |
	propFlag_ASYNC_FN |
	propFlag_GENERATOR_FN |
	propFlag_START_OF_OPT_CHAIN |
	propFlag_STATIC_MEMBER |
	propFlag_YIELD_ALL |
	propFlag_EXPORT_DEFAULT |
	propFlag_EXPORT_ALL_FROM |
	propFlag_INCRDECR |
	propFlag_QUOTED |
	propFlag_FREE_CALL |
	propFlag_DIRECT_EVAL |
	propFlag_COMPUTED_PROP_METHOD |
	propFlag_COMPUTED_PROP_GETTER |
	propFlag_COMPUTED_PROP_SETTER

// GetQualifiedName takes a set of GETPROP nodes and produces a string that
// is each property separated by dots. If the node ultimately under the left
// sub-tree is not a simple name, this is not a valid qualified name.
//
// Returns null if this is not a qualified name, or a dot-separated string of
// the name and properties.
//
// Java: Node.getQualifiedName()
//
// Go-only: performance — the body is in getQualifiedNameString, so that this
// wrapper inlines: a caller that only reads the result keeps the string
// header on its stack instead of the heap (these headers were 6% of
// checks-only's allocations).
func (n Node) GetQualifiedName() *string {
	s, ok := n.getQualifiedNameString()
	if !ok {
		return nil
	}
	return &s
}

// GetQualifiedNameString is the body of GetQualifiedName; ok is false
// where Java returns null.
func (n Node) GetQualifiedNameString() (string, bool) {
	if n.IsNil() {
		return "", false
	}
	return n.getQualifiedNameString()
}

// getQualifiedNameString is the body of getQualifiedName; ok is false
// where Java returns null.
func (n Node) getQualifiedNameString() (string, bool) {
	switch n.rec().token {
	case Token_NAME:
		name := n.GetString()
		if name == "" {
			return "", false
		}
		return name, true
	case Token_GETPROP:
		first := n.firstNode()
		if !first.IsGetProp() {
			left, ok := first.getQualifiedNameString()
			if !ok {
				return "", false
			}
			return left + "." + n.GetString(), true
		}
		var partsBuf [8]string
		parts := partsBuf[:0]
		totalLen := 0
		cur := n
		for cur.IsGetProp() {
			p := cur.GetString()
			parts = append(parts, p)
			totalLen += 1 + len(p)
			cur = cur.firstNode()
		}
		left, ok := cur.getQualifiedNameString()
		if !ok {
			return "", false
		}
		var builder strings.Builder
		builder.Grow(len(left) + totalLen)
		builder.WriteString(left)
		for i := len(parts) - 1; i >= 0; i-- {
			builder.WriteByte('.')
			builder.WriteString(parts[i])
		}
		return builder.String(), true
	case Token_THIS:
		return "this", true
	case Token_SUPER:
		return "super", true
	default:
		return "", false
	}
}

// Java: Node.getQualifiedNameObject()
func (n Node) GetQualifiedNameObject() QualifiedName {
	if n.IsQualifiedName() {
		return NewQualifiedName_NodeQname(n)
	}
	return nil
}

// GetOriginalQualifiedName takes a set of GETPROP nodes and produces a
// string that is each property separated by dots. If the node ultimately
// under the left sub-tree is not a simple name, this is not a valid
// qualified name. This method returns the original name of each segment
// rather than the renamed version.
//
// Deprecated: "original name" is poorly defined. See #getOriginalName
//
// Java: Node.getOriginalQualifiedName()
func (n Node) GetOriginalQualifiedName() *string {
	if n.rec().token == Token_NAME {
		var name string
		if n.GetOriginalName() == nil {
			name = n.GetString()
		} else {
			name = *n.GetOriginalName()
		}
		if name == "" {
			return nil
		}
		return &name
	} else if n.rec().token == Token_GETPROP {
		left := n.GetFirstChild().GetOriginalQualifiedName()
		if left == nil {
			return nil
		}

		var right string
		if n.GetOriginalName() == nil {
			right = n.GetString()
		} else {
			right = *n.GetOriginalName()
		}

		s := *left + "." + right
		return &s
	} else if n.rec().token == Token_THIS {
		s := "this"
		return &s
	} else if n.rec().token == Token_SUPER {
		s := "super"
		return &s
	} else {
		return nil
	}
}

// IsQualifiedName returns whether a node corresponds to a simple or a
// qualified name, such as `x` or `a.b.c` or `this.a`.
//
// Java: Node.isQualifiedName()
func (n Node) IsQualifiedName() bool {
	for {
		switch n.rec().token {
		case Token_NAME:
			return n.GetStringID() > 1
		case Token_THIS, Token_SUPER:
			return true
		case Token_GETPROP:
			n = n.firstNode()
		default:
			return false
		}
	}
}

// MatchesName returns whether a node matches a simple name, such as `x`,
// returns false if this is not a NAME node.
//
// Java: Node.matchesName(String)
func (n Node) MatchesName(name string) bool {
	if n.rec().token != Token_NAME {
		return false
	}
	internalString := n.GetString()
	return internalString != "" && name == internalString
}

// MatchesNameNode checks that if two NAME node match, returns false if
// either node is not a NAME node. As a empty string is not considered a
// valid Name (it is an AST placeholder), empty strings are never considered
// to be matches.
//
// Java: Node.matchesName(Node)
func (n Node) MatchesNameNode(other Node) bool {
	if n.rec().token != Token_NAME || other.rec().token != Token_NAME {
		return false
	}

	// Compare interned string IDs directly (id > 1 means non-empty!).
	id := rhinoStringID(n.rec().payload)
	return id > 1 && id == rhinoStringID(other.rec().payload)
}

// MatchesQualifiedName returns whether a node matches a simple or a
// qualified name, such as `x` or `a.b.c` or `this.a`.
//
// Java: Node.matchesQualifiedName(String)
func (n Node) MatchesQualifiedName(qname string) bool {
	for {
		switch n.rec().token {
		case Token_NAME, Token_IMPORT_STAR:
			name := n.GetString()
			return name != "" && name == qname && strings.IndexByte(name, '.') < 0
		case Token_THIS:
			return qname == "this"
		case Token_SUPER:
			return qname == "super"
		case Token_GETPROP:
			prop := n.GetString()
			dot := len(qname) - len(prop) - 1
			if dot < 1 || qname[dot] != '.' || qname[dot+1:] != prop || strings.IndexByte(prop, '.') >= 0 {
				return false
			}
			qname = qname[:dot]
			n = n.firstNode()
		default:
			return false
		}
	}
}

// MatchesQualifiedNameWithEnd returns whether a node matches a simple or a
// qualified name, such as `x` or `a.b.c` or `this.a`.
//
// Java: Node.matchesQualifiedName(String, int)
func (n Node) MatchesQualifiedNameWithEnd(qname string, endIndex int) bool {
	start := jutil.StringLastIndexOfCharFrom(qname, '.', endIndex-1) + 1

	switch n.GetToken() {
	case Token_NAME, Token_IMPORT_STAR:
		name := n.GetString()
		return start == 0 && name != "" && jutil.StringLength(name) == endIndex && jutil.StringStartsWith(qname, name)
	case Token_THIS:
		return start == 0 && endIndex == 4 && jutil.StringStartsWith(qname, "this")
	case Token_SUPER:
		return start == 0 && endIndex == 5 && jutil.StringStartsWith(qname, "super")
	case Token_GETPROP:
		prop := n.GetString()
		return start > 1 &&
			jutil.StringLength(prop) == endIndex-start &&
			jutil.StringRegionMatches(prop, 0, qname, start, endIndex-start) &&
			n.GetFirstChild().MatchesQualifiedNameWithEnd(qname, start-1)

	case Token_MEMBER_FUNCTION_DEF:
		// These are explicitly *not* qualified name components.
		return false
	default:
		return false
	}
}

// MatchesQualifiedNameNode returns whether a node matches a simple or a
// qualified name, such as `x` or `a.b.c` or `this.a`.
//
// Java: Node.matchesQualifiedName(Node)
func (n Node) MatchesQualifiedNameNode(other Node) bool {
	if other.rec().token != n.rec().token {
		return false
	}

	switch n.rec().token {
	case Token_NAME:
		return n.MatchesNameNode(other)
	case Token_THIS, Token_SUPER:
		return true
	case Token_GETPROP:
		return uint32(n.rec().payload) == uint32(other.rec().payload) &&
			n.GetFirstChild().MatchesQualifiedNameNode(other.GetFirstChild())

	case Token_MEMBER_FUNCTION_DEF:
		// These are explicitly *not* qualified name components.
		return false
	default:
		return false
	}
}

// IsUnscopedQualifiedName returns whether a node corresponds to a simple or
// a qualified name without a "this" reference, such as `a.b.c`, but not
// `this.a`.
//
// Java: Node.isUnscopedQualifiedName()
func (n Node) IsUnscopedQualifiedName() bool {
	for {
		switch n.rec().token {
		case Token_NAME:
			return n.GetStringID() > 1
		case Token_GETPROP:
			n = n.firstNode()
		default:
			return false
		}
	}
}

// Java: Node.isValidAssignmentTarget()
func (n Node) IsValidAssignmentTarget() bool {
	switch n.GetToken() {
	case Token_NAME, Token_GETPROP, Token_GETELEM, Token_ARRAY_PATTERN, Token_OBJECT_PATTERN:
		return true
	default:
		return false
	}
}

// Clone always throws: use CloneNode.
//
// Java: Node.clone()
func (n Node) Clone() any {
	panic(jutil.NewUnsupportedOperationExceptionMsg("Did you mean cloneNode?"))
}

// CloneNode returns a detached clone of the Node, specifically excluding its
// children.
//
// Java: Node.cloneNode()
func (n Node) CloneNode() Node {
	return n.CloneNodeWithTypeExprs(false)
}

// CloneNodeWithTypeExprs returns a detached clone of the Node, specifically
// excluding its children.
//
// Go-only: language-adapter — dispatches to the cloneNode overrides of the
// folded subclasses.
//
// Java: Node.cloneNode(boolean)
func (n Node) CloneNodeWithTypeExprs(cloneTypeExprs bool) Node {
	clone := allocNode()
	clone.rec().token = n.rec().token
	clone.rec().kindFlags = n.rec().kindFlags & kindMask
	clone.rec().payload = n.rec().payload
	Node_CopyBaseNodeFields(n, clone, cloneTypeExprs)
	return clone
}

// Java: Node.copyBaseNodeFields(Node, Node, boolean)
func Node_CopyBaseNodeFields(source Node, dest Node, cloneTypeExprs bool) {
	dest.rec().linenoCharno = source.rec().linenoCharno
	dest.setTC(dest.store().translateTC(source.store(), source.tc(), source.rec().kindFlags&kindFlagTypeIsColor != 0))
	dest.rec().kindFlags = (dest.rec().kindFlags &^ (kindFlagTypeIsColor | kindFlagHasExtBox | kindFlagHasExtra)) |
		(source.rec().kindFlags & (kindFlagTypeIsColor | kindFlagHasExtBox))
	if srcBox := source.extBoxOrNil(); srcBox != nil {
		srcBox.propListHead.markShared()
		boxCopy := *srcBox
		boxCopy.typeBeforeCast = dest.store().translateTypeHandle(source.store(), boxCopy.typeBeforeCast)
		dest.setAux(nodeAux{ext: &boxCopy})
	} else {
		dest.setAux(source.auxValue())
	}
	if source.rec().kindFlags&kindFlagHasExtra != 0 && source.coldPropListHead() != nil {
		dest.rec().kindFlags |= kindFlagHasExtra
	}
	dest.copyLengthFrom(source)
	dest.node_setOriginalNameID(source.node_originalNameID())
	dest.setPropsAndSourceFile(source)

	// TODO(johnlenz): Remove this once JSTypeExpression are immutable
	if cloneTypeExprs {
		info := source.GetJSDocInfo()
		if info != nil {
			dest.SetJSDocInfo(info.CloneWithTypeNodes(true))
		}
	}
}

// CloneTree returns a detached clone of the Node and all its children.
//
// Java: Node.cloneTree()
func (n Node) CloneTree() Node {
	return n.CloneTreeWithTypeExprs(false)
}

// Java: Node.cloneTree(boolean)
func (n Node) CloneTreeWithTypeExprs(cloneTypeExprs bool) Node {
	result := n.CloneNodeWithTypeExprs(cloneTypeExprs)
	if n.rec().first != 0 {
		var firstChild Node
		var lastChild Node
		for n2 := n.firstNode(); !n2.IsNil(); n2 = n2.nextNode() {
			n2clone := n2.CloneTreeWithTypeExprs(cloneTypeExprs)
			n2clone.rec().parent = n2clone.ref(result)
			if firstChild.IsNil() {
				firstChild = n2clone
				lastChild = firstChild
			} else {
				lastChild.rec().next = lastChild.ref(n2clone)
				n2clone.rec().previous = n2clone.ref(lastChild)
				lastChild = n2clone
			}
		}
		firstChild.rec().previous = firstChild.ref(lastChild)
		lastChild.rec().next = 0
		result.rec().first = result.ref(firstChild)
	}

	if shadow := n.GetClosureUnawareShadow(); !shadow.IsNil() {
		result.SetClosureUnawareShadow(shadow.CloneTreeWithTypeExprs(cloneTypeExprs))
	}

	return result
}

// Srcref copies the source info from `other` onto `this`.
//
// Java: Node.srcref(Node)
func (n Node) Srcref(other Node) Node {
	n.SetStaticSourceFileFrom(other)
	n.node_setOriginalNameID(other.node_originalNameID())
	n.rec().linenoCharno = other.rec().linenoCharno
	n.copyLengthFrom(other)
	return n
}

// SrcrefTree copies, for all Nodes in the subtree of `this`, the source info
// from `other`.
//
// Java: Node.srcrefTree(Node)
func (n Node) SrcrefTree(other Node) Node {
	n.Srcref(other)
	for child := n.firstNode(); !child.IsNil(); child = child.nextNode() {
		child.SrcrefTree(other)
	}
	return n
}

// SrcrefIfMissing copies the source info from `other` iff source info is not
// set on `this`.
//
// Java: Node.srcrefIfMissing(Node)
func (n Node) SrcrefIfMissing(other Node) Node {
	if n.rec().sourceFile == 0 {
		n.SetStaticSourceFileFrom(other)
		n.rec().linenoCharno = other.rec().linenoCharno
		n.copyLengthFrom(other)
	}

	// TODO(lharker): should this be inside the above if condition?
	// If the node already has a source file, it seems strange to
	// go ahead and set the original name anyway.
	if n.node_originalNameID() == 0 {
		n.node_setOriginalNameID(other.node_originalNameID())
	}

	return n
}

// SrcrefTreeIfMissing copies, for all Nodes in the subtree of `this`, iff
// source info is not set, the source info from `other`.
//
// Java: Node.srcrefTreeIfMissing(Node)
func (n Node) SrcrefTreeIfMissing(other Node) Node {
	n.SrcrefIfMissing(other)
	for child := n.firstNode(); !child.IsNil(); child = child.nextNode() {
		child.SrcrefTreeIfMissing(other)
	}
	return n
}
