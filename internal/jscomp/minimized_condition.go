// Port of com/google/javascript/jscomp/MinimizedCondition.java @ closure-compiler v20260415

/*
 * Copyright 2013 The Closure Compiler Authors.
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
	"math"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// MinimizedCondition is a class that represents a minimized conditional
// expression. This is a conditional expression that has been massaged
// according to DeMorgan's laws in order to minimize the length of the
// source representation.
//
// Depending on the context, a leading NOT node in front of the conditional
// may or may not be counted as a cost, so this class provides ways to
// access minimized versions of both of those abstract syntax trees (ASTs).
type MinimizedCondition struct {
	// A representation equivalent to the original condition.
	positive *MinimizedCondition_MeasuredNode

	// A representation equivalent to the negation of the original
	// condition.
	negative *MinimizedCondition_MeasuredNode
}

// MinimizedCondition_MinimizationStyle: definitions of the style of
// minimization preferred.
type MinimizedCondition_MinimizationStyle int32

const (
	// MinimizedCondition_MinimizationStyle_PREFER_UNNEGATED: compute the
	// length of the minimized condition as including any leading NOT node,
	// if present.
	MinimizedCondition_MinimizationStyle_PREFER_UNNEGATED MinimizedCondition_MinimizationStyle = iota
	// MinimizedCondition_MinimizationStyle_ALLOW_LEADING_NOT: compute the
	// length of the minimized condition without penalizing a leading NOT
	// node, if present.
	MinimizedCondition_MinimizationStyle_ALLOW_LEADING_NOT
)

var minimizedCondition_MinimizationStyle_names = [...]string{"PREFER_UNNEGATED", "ALLOW_LEADING_NOT"}

// MinimizedCondition_MinimizationStyle_values is MinimizationStyle.values().
func MinimizedCondition_MinimizationStyle_values() []MinimizedCondition_MinimizationStyle {
	return []MinimizedCondition_MinimizationStyle{
		MinimizedCondition_MinimizationStyle_PREFER_UNNEGATED,
		MinimizedCondition_MinimizationStyle_ALLOW_LEADING_NOT,
	}
}

// Name is Enum.name().
func (s MinimizedCondition_MinimizationStyle) Name() string {
	return minimizedCondition_MinimizationStyle_names[s]
}

// String is Enum.toString().
func (s MinimizedCondition_MinimizationStyle) String() string { return s.Name() }

// Ordinal is Enum.ordinal().
func (s MinimizedCondition_MinimizationStyle) Ordinal() int { return int(s) }

// Java: MinimizedCondition(MeasuredNode, MeasuredNode)
func newMinimizedCondition(p *MinimizedCondition_MeasuredNode, n *MinimizedCondition_MeasuredNode) *MinimizedCondition {
	return &MinimizedCondition{positive: p, negative: n.Change()}
}

// MinimizedCondition_FromConditionNode returns a MinimizedCondition that
// represents the condition node after minimization.
//
// Java: MinimizedCondition.fromConditionNode(Node)
func MinimizedCondition_FromConditionNode(n *rhino.Node) *MinimizedCondition {
	jutil.CheckState(n.HasParent())
	switch n.GetToken() {
	case rhino.Token_NOT, rhino.Token_AND, rhino.Token_OR, rhino.Token_HOOK, rhino.Token_COMMA:
		return minimizedCondition_computeMinimizedCondition(n)
	default:
		return MinimizedCondition_Unoptimized(n)
	}
}

// GetMinimized returns the shorter representation of the original
// condition node.
//
// Depending on the context, this may require to either penalize or not the
// existence of a leading NOT node.
//
//   - When style is PREFER_UNNEGATED, simply try to minimize the total
//     length of the conditional.
//   - When style is ALLOW_LEADING_NOT, prefer the right side in cases such
//     as: !x || !y || z  ==>  !(x && y && !z). This is useful in contexts
//     such as IFs or HOOKs where subsequent optimizations can efficiently
//     deal with leading NOTs.
//
// Returns the minimized condition MeasuredNode, with equivalent semantics
// to that passed to FromConditionNode.
//
// Java: MinimizedCondition.getMinimized(MinimizationStyle)
func (c *MinimizedCondition) GetMinimized(style MinimizedCondition_MinimizationStyle) *MinimizedCondition_MeasuredNode {
	if style == MinimizedCondition_MinimizationStyle_PREFER_UNNEGATED ||
		c.positive.node.IsNot() ||
		c.positive.length <= c.negative.length {
		return c.positive
	} else {
		return c.negative.AddNot()
	}
}

// MinimizedCondition_Unoptimized returns a MeasuredNode of the given
// condition node, without minimizing the result.
//
// Since a MinimizedCondition necessarily must contain two trees, this
// method sets the negative side to a invalid node with an unreasonably high
// length so that it will never be chosen by GetMinimized.
//
// n is the conditional expression tree; returns a MinimizedCondition object
// representing that tree.
//
// Java: MinimizedCondition.unoptimized(Node)
func MinimizedCondition_Unoptimized(n *rhino.Node) *MinimizedCondition {
	jutil.CheckNotNull(n.GetParent())
	pos := NewMinimizedCondition_MeasuredNode(n, nil, 0, false)
	neg := NewMinimizedCondition_MeasuredNode(nil, nil, math.MaxInt32, true)
	return newMinimizedCondition(pos, neg)
}

// MinimizedCondition_PickBest returns the best, prefer unchanged.
//
// Java: MinimizedCondition.pickBest(MeasuredNode, MeasuredNode)
func MinimizedCondition_PickBest(a *MinimizedCondition_MeasuredNode, b *MinimizedCondition_MeasuredNode) *MinimizedCondition_MeasuredNode {
	if a.length == b.length {
		if b.IsChanged() {
			return a
		}
		return b
	}

	if a.length < b.length {
		return a
	}
	return b
}

// minimizedCondition_computeMinimizedCondition minimizes the condition at
// the given node.
//
// n is the conditional expression tree to minimize; returns a
// MinimizedCondition object representing that tree.
//
// Java: MinimizedCondition.computeMinimizedCondition(Node)
func minimizedCondition_computeMinimizedCondition(n *rhino.Node) *MinimizedCondition {
	switch n.GetToken() {
	case rhino.Token_NOT:
		subtree := minimizedCondition_computeMinimizedCondition(n.GetFirstChild())
		positive := MinimizedCondition_PickBest(minimizedCondition_MeasuredNode_addNode(n, subtree.positive), subtree.negative)
		negative := MinimizedCondition_PickBest(
			subtree.negative.
				Negate(), // since parent node `n` is a NOT, we need to negate the subtree's
			// computed `negative` to obtain the parent `n`'s real negative.
			subtree.positive)
		return newMinimizedCondition(positive, negative)
	case rhino.Token_AND, rhino.Token_OR:
		var complementToken rhino.Token
		if n.IsAnd() {
			complementToken = rhino.Token_OR
		} else {
			complementToken = rhino.Token_AND
		}
		complementNode := rhino.NewNode(complementToken).Srcref(n)
		leftSubtree := minimizedCondition_computeMinimizedCondition(n.GetFirstChild())
		rightSubtree := minimizedCondition_computeMinimizedCondition(n.GetLastChild())
		positive := MinimizedCondition_PickBest(
			minimizedCondition_MeasuredNode_addNode(n, leftSubtree.positive, rightSubtree.positive),
			minimizedCondition_MeasuredNode_addNode(complementNode, leftSubtree.negative, rightSubtree.negative).
				Negate())
		negative := MinimizedCondition_PickBest(
			minimizedCondition_MeasuredNode_addNode(n, leftSubtree.positive, rightSubtree.positive).Negate(),
			minimizedCondition_MeasuredNode_addNode(complementNode, leftSubtree.negative, rightSubtree.negative).
				Change())
		return newMinimizedCondition(positive, negative)
	case rhino.Token_HOOK:
		cond := n.GetFirstChild()
		thenNode := cond.GetNext()
		elseNode := thenNode.GetNext()
		thenSubtree := minimizedCondition_computeMinimizedCondition(thenNode)
		elseSubtree := minimizedCondition_computeMinimizedCondition(elseNode)
		positive := minimizedCondition_MeasuredNode_addNode(
			n, minimizedCondition_MeasuredNode_forNode(cond), thenSubtree.positive, elseSubtree.positive)
		negative := minimizedCondition_MeasuredNode_addNode(
			n, minimizedCondition_MeasuredNode_forNode(cond), thenSubtree.negative, elseSubtree.negative)
		return newMinimizedCondition(positive, negative)
	case rhino.Token_COMMA:
		lhs := n.GetFirstChild()
		rhsSubtree := minimizedCondition_computeMinimizedCondition(lhs.GetNext())
		positive := minimizedCondition_MeasuredNode_addNode(n, minimizedCondition_MeasuredNode_forNode(lhs), rhsSubtree.positive)
		negative := minimizedCondition_MeasuredNode_addNode(n, minimizedCondition_MeasuredNode_forNode(lhs), rhsSubtree.negative)
		return newMinimizedCondition(positive, negative)
	default:
		pos := minimizedCondition_MeasuredNode_forNode(n)
		neg := pos.Negate()
		return newMinimizedCondition(pos, neg)
	}
}

// MinimizedCondition_MeasuredNode is an AST-node along with some
// additional metadata.
type MinimizedCondition_MeasuredNode struct {
	node     *rhino.Node
	length   int
	changed  bool
	children []*MinimizedCondition_MeasuredNode
}

// Java: MinimizedCondition.MeasuredNode(Node, MeasuredNode[], int, boolean)
func NewMinimizedCondition_MeasuredNode(n *rhino.Node, children []*MinimizedCondition_MeasuredNode, len_ int, ch bool) *MinimizedCondition_MeasuredNode {
	return &MinimizedCondition_MeasuredNode{
		node:     n,
		children: children,
		length:   len_,
		changed:  ch,
	}
}

// Java: MinimizedCondition.MeasuredNode.isChanged()
func (m *MinimizedCondition_MeasuredNode) IsChanged() bool {
	return m.changed
}

// Java: MinimizedCondition.MeasuredNode.isNot()
func (m *MinimizedCondition_MeasuredNode) IsNot() bool {
	return m.node.IsNot()
}

// Java: MinimizedCondition.MeasuredNode.withoutNot()
func (m *MinimizedCondition_MeasuredNode) WithoutNot() *MinimizedCondition_MeasuredNode {
	jutil.CheckState(m.IsNot())
	return MinimizedCondition_MeasuredNode_NormalizeChildren(m.node, m.children)[0].Change()
}

// Java: MinimizedCondition.MeasuredNode.negate()
func (m *MinimizedCondition_MeasuredNode) Negate() *MinimizedCondition_MeasuredNode {
	switch m.node.GetToken() {
	case rhino.Token_EQ:
		return m.UpdateToken(rhino.Token_NE)
	case rhino.Token_NE:
		return m.UpdateToken(rhino.Token_EQ)
	case rhino.Token_SHEQ:
		return m.UpdateToken(rhino.Token_SHNE)
	case rhino.Token_SHNE:
		return m.UpdateToken(rhino.Token_SHEQ)
	case rhino.Token_NOT:
		return m.WithoutNot()
	default:
		return m.AddNot()
	}
}

// Java: MinimizedCondition.MeasuredNode.normalizeChildren(Node, MeasuredNode[])
func MinimizedCondition_MeasuredNode_NormalizeChildren(node *rhino.Node, children []*MinimizedCondition_MeasuredNode) []*MinimizedCondition_MeasuredNode {
	if children != nil || !node.HasChildren() {
		return children
	} else {
		measuredChildren := make([]*MinimizedCondition_MeasuredNode, node.GetChildCount())
		child := 0
		for c := node.GetFirstChild(); c != nil; c = c.GetNext() {
			measuredChildren[child] = minimizedCondition_MeasuredNode_forNode(c)
			child++
		}
		return measuredChildren
	}
}

// Java: MinimizedCondition.MeasuredNode.updateToken(Token)
func (m *MinimizedCondition_MeasuredNode) UpdateToken(token rhino.Token) *MinimizedCondition_MeasuredNode {
	return NewMinimizedCondition_MeasuredNode(
		rhino.NewNode(token).Srcref(m.node), MinimizedCondition_MeasuredNode_NormalizeChildren(m.node, m.children), m.length, true)
}

// Java: MinimizedCondition.MeasuredNode.addNot()
func (m *MinimizedCondition_MeasuredNode) AddNot() *MinimizedCondition_MeasuredNode {
	return minimizedCondition_MeasuredNode_addNode(
		rhino.NewNode(rhino.Token_NOT).Srcref(m.node), m).Change()
}

// Java: MinimizedCondition.MeasuredNode.change()
func (m *MinimizedCondition_MeasuredNode) Change() *MinimizedCondition_MeasuredNode {
	if m.IsChanged() {
		return m
	}
	return NewMinimizedCondition_MeasuredNode(m.node, m.children, m.length, true)
}

// minimizedCondition_MeasuredNode_estimateCostOneLevel estimates the number
// of characters in the textual representation of the given node and that
// will be devoted to negation or parentheses. Since these are the only
// characters that flipping a condition according to De Morgan's rule can
// affect, these are the only ones we count. Not nodes are counted by the
// NOT node itself, whereas parentheses around an expression are counted by
// the parent node.
//
// n is the node to be checked; returns the number of negations and
// parentheses in the node.
//
// Java: MinimizedCondition.MeasuredNode.estimateCostOneLevel(Node, MeasuredNode...)
func minimizedCondition_MeasuredNode_estimateCostOneLevel(n *rhino.Node, children ...*MinimizedCondition_MeasuredNode) int {
	cost := 0
	if n.IsNot() {
		cost++ // A negation is needed.
	}
	parentPrecedence := NodeUtil_Precedence(n.GetToken())
	for _, child := range children {
		if child.IsLowerPrecedenceThan(parentPrecedence) {
			cost += 2 // A pair of parenthesis is needed.
		}
	}
	return cost
}

// IsLowerPrecedenceThan reports whether the node type has lower precedence
// than "precedence".
//
// Java: MinimizedCondition.MeasuredNode.isLowerPrecedenceThan(int)
func (m *MinimizedCondition_MeasuredNode) IsLowerPrecedenceThan(precedence int) bool {
	return NodeUtil_Precedence(m.node.GetToken()) < precedence
}

// minimizedCondition_MeasuredNode_addNode: the returned MeasuredNode is
// only marked as changed if the children are marked as changed.
//
// Java: MinimizedCondition.MeasuredNode.addNode(Node, MeasuredNode...)
func minimizedCondition_MeasuredNode_addNode(parent *rhino.Node, children ...*MinimizedCondition_MeasuredNode) *MinimizedCondition_MeasuredNode {
	cost := 0
	changed := false
	for _, child := range children {
		cost += child.length
		changed = changed || child.changed
	}
	cost += minimizedCondition_MeasuredNode_estimateCostOneLevel(parent, children...)
	return NewMinimizedCondition_MeasuredNode(parent, children, cost, changed)
}

// minimizedCondition_MeasuredNode_forNode returns a MeasuredNode for a
// non-participating AST Node. This is used for leaf expression nodes.
//
// Java: MinimizedCondition.MeasuredNode.forNode(Node)
func minimizedCondition_MeasuredNode_forNode(n *rhino.Node) *MinimizedCondition_MeasuredNode {
	return NewMinimizedCondition_MeasuredNode(n, nil, 0, false)
}

// WillChange reports whether the MeasuredNode is a change from the
// original. This can either be a change within the original AST tree or a
// replacement of the original node.
//
// Java: MinimizedCondition.MeasuredNode.willChange(Node)
func (m *MinimizedCondition_MeasuredNode) WillChange(original *rhino.Node) bool {
	jutil.CheckNotNull(original)
	return original != m.node || m.IsChanged()
}

// ApplyTo updates the AST for the result of this MeasuredNode. This can
// either be a change within the original AST tree or a replacement of the
// original node.
//
// Java: MinimizedCondition.MeasuredNode.applyTo(Node)
func (m *MinimizedCondition_MeasuredNode) ApplyTo(original *rhino.Node) *rhino.Node {
	jutil.CheckNotNull(original)
	jutil.CheckState(m.WillChange(original))
	replacement := m.BuildReplacement()
	if original != replacement {
		m.SafeDetach(replacement)
		original.ReplaceWith(replacement)
	}
	return replacement
}

// SafeDetach detaches a node only IIF it is in the tree.
//
// Java: MinimizedCondition.MeasuredNode.safeDetach(Node)
func (m *MinimizedCondition_MeasuredNode) SafeDetach(n *rhino.Node) *rhino.Node {
	if n.HasParent() {
		return n.Detach()
	}
	return n
}

// BuildReplacement builds the final AST structure, detaching component
// Nodes as necessary from the original AST. The root Node, if currently
// attached is left attached to avoid the need to keep track of its
// position.
//
// Java: MinimizedCondition.MeasuredNode.buildReplacement()
func (m *MinimizedCondition_MeasuredNode) BuildReplacement() *rhino.Node {
	if m.children != nil {
		m.node.DetachChildren()
		for _, child := range m.children {
			replacementChild := m.SafeDetach(child.BuildReplacement())
			m.node.AddChildToBack(replacementChild)
		}
	}
	return m.node
}
