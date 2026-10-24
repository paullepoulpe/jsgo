// Port of com/google/javascript/jscomp/MinimizedCondition.java @ closure-compiler v20261021

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

	nodes     [32]MinimizedCondition_MeasuredNode
	nodesUsed int
	overflow  []MinimizedCondition_MeasuredNode
}

var minimizedCondition_unoptimizedNeg = &MinimizedCondition_MeasuredNode{
	length:  math.MaxInt32,
	changed: true,
}

var minimizedCondition_NOT_PRECEDENCE = NodeUtil_Precedence(rhino.Token_NOT)

func (c *MinimizedCondition) reset() {
	c.positive = nil
	c.negative = nil
	c.nodesUsed = 0
	if c.overflow != nil {
		c.overflow = c.overflow[:0]
	}
}

func (c *MinimizedCondition) allocNode() *MinimizedCondition_MeasuredNode {
	if c == nil {
		return &MinimizedCondition_MeasuredNode{}
	}
	if c.nodesUsed < len(c.nodes) {
		m := &c.nodes[c.nodesUsed]
		c.nodesUsed++
		*m = MinimizedCondition_MeasuredNode{mc: c}
		return m
	}
	if len(c.overflow) == cap(c.overflow) {
		newCap := cap(c.overflow) * 2
		if newCap < 64 {
			newCap = 64
		}
		c.overflow = make([]MinimizedCondition_MeasuredNode, 0, newCap)
	}
	i := len(c.overflow)
	c.overflow = c.overflow[:i+1]
	m := &c.overflow[i]
	*m = MinimizedCondition_MeasuredNode{mc: c}
	return m
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
func MinimizedCondition_FromConditionNode(n rhino.Node) *MinimizedCondition {
	return MinimizedCondition_FromConditionNodeInto(n, &MinimizedCondition{})
}

func MinimizedCondition_FromConditionNodeInto(n rhino.Node, c *MinimizedCondition) *MinimizedCondition {
	jutil.CheckState(n.HasParent())
	switch n.GetToken() {
	case rhino.Token_NOT, rhino.Token_AND, rhino.Token_OR, rhino.Token_HOOK, rhino.Token_COMMA:
		c.reset()
		pos, neg := minimizedCondition_computeMinimizedConditionInto(n, c)
		c.positive = pos
		c.negative = neg
		return c
	default:
		return MinimizedCondition_UnoptimizedInto(n, c)
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
		c.positive.IsNot() ||
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
func MinimizedCondition_Unoptimized(n rhino.Node) *MinimizedCondition {
	return MinimizedCondition_UnoptimizedInto(n, &MinimizedCondition{})
}

func MinimizedCondition_UnoptimizedInto(n rhino.Node, c *MinimizedCondition) *MinimizedCondition {
	rhino.CheckNotNull(n.GetParent())
	c.reset()
	c.negative = minimizedCondition_unoptimizedNeg
	pos := &c.nodes[0]
	c.nodesUsed = 1
	*pos = MinimizedCondition_MeasuredNode{
		node: n,
		mc:   c,
	}
	c.positive = pos
	return c
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
func minimizedCondition_computeMinimizedCondition(n rhino.Node) *MinimizedCondition {
	c := &MinimizedCondition{}
	pos, neg := minimizedCondition_computeMinimizedConditionInto(n, c)
	c.positive = pos
	c.negative = neg
	return c
}

func minimizedCondition_computeMinimizedConditionInto(n rhino.Node, c *MinimizedCondition) (*MinimizedCondition_MeasuredNode, *MinimizedCondition_MeasuredNode) {
	switch n.GetToken() {
	case rhino.Token_NOT:
		subPos, subNeg := minimizedCondition_computeMinimizedConditionInto(n.GetFirstChild(), c)
		positive := MinimizedCondition_PickBest(minimizedCondition_MeasuredNode_addNode1(c, n, subPos), subNeg)
		negative := MinimizedCondition_PickBest(
			subNeg.negateIn(c), // since parent node `n` is a NOT, we need to negate the subtree's
			// computed `negative` to obtain the parent `n`'s real negative.
			subPos)
		return positive, negative.changeIn(c)
	case rhino.Token_AND, rhino.Token_OR:
		var complementToken rhino.Token
		if n.IsAnd() {
			complementToken = rhino.Token_OR
		} else {
			complementToken = rhino.Token_AND
		}
		leftPos, leftNeg := minimizedCondition_computeMinimizedConditionInto(n.GetFirstChild(), c)
		rightPos, rightNeg := minimizedCondition_computeMinimizedConditionInto(n.GetLastChild(), c)
		posTree := minimizedCondition_MeasuredNode_addNode2(c, n, 0, leftPos, rightPos)
		negComplementTree := minimizedCondition_MeasuredNode_addNode2(c, n, complementToken, leftNeg, rightNeg)
		positive := MinimizedCondition_PickBest(
			posTree,
			negComplementTree.negateIn(c))
		negative := MinimizedCondition_PickBest(
			posTree.negateIn(c),
			negComplementTree.changeIn(c))
		return positive, negative.changeIn(c)
	case rhino.Token_HOOK:
		cond := n.GetFirstChild()
		thenNode := cond.GetNext()
		elseNode := thenNode.GetNext()
		thenPos, thenNeg := minimizedCondition_computeMinimizedConditionInto(thenNode, c)
		elsePos, elseNeg := minimizedCondition_computeMinimizedConditionInto(elseNode, c)
		condMeasured := minimizedCondition_MeasuredNode_forNodeIn(c, cond)
		positive := minimizedCondition_MeasuredNode_addNode3(
			c, n, condMeasured, thenPos, elsePos)
		negative := minimizedCondition_MeasuredNode_addNode3(
			c, n, condMeasured, thenNeg, elseNeg)
		return positive, negative.changeIn(c)
	case rhino.Token_COMMA:
		lhs := n.GetFirstChild()
		rhsPos, rhsNeg := minimizedCondition_computeMinimizedConditionInto(lhs.GetNext(), c)
		lhsMeasured := minimizedCondition_MeasuredNode_forNodeIn(c, lhs)
		positive := minimizedCondition_MeasuredNode_addNode2(c, n, 0, lhsMeasured, rhsPos)
		negative := minimizedCondition_MeasuredNode_addNode2(c, n, 0, lhsMeasured, rhsNeg)
		return positive, negative.changeIn(c)
	default:
		pos := minimizedCondition_MeasuredNode_forNodeIn(c, n)
		neg := pos.negateIn(c)
		return pos, neg.changeIn(c)
	}
}

// MinimizedCondition_MeasuredNode is an AST-node along with some
// additional metadata.
type MinimizedCondition_MeasuredNode struct {
	node          rhino.Node
	mc            *MinimizedCondition
	child0        *MinimizedCondition_MeasuredNode
	child1        *MinimizedCondition_MeasuredNode
	child2        *MinimizedCondition_MeasuredNode
	extraChildren []*MinimizedCondition_MeasuredNode
	length        int
	numChildren   int32
	virtualToken  rhino.Token
	changed       bool
	hasChildren   bool
}

func (m *MinimizedCondition_MeasuredNode) token() rhino.Token {
	if m.virtualToken != 0 {
		return m.virtualToken
	}
	return m.node.GetToken()
}

// Java: MinimizedCondition.MeasuredNode(Node, MeasuredNode[], int, boolean)
func NewMinimizedCondition_MeasuredNode(n rhino.Node, children []*MinimizedCondition_MeasuredNode, len_ int, ch bool) *MinimizedCondition_MeasuredNode {
	return newMinimizedCondition_MeasuredNodeIn(nil, n, children, len_, ch)
}

func newMinimizedCondition_MeasuredNodeIn(c *MinimizedCondition, n rhino.Node, children []*MinimizedCondition_MeasuredNode, len_ int, ch bool) *MinimizedCondition_MeasuredNode {
	m := c.allocNode()
	m.node = n
	m.length = len_
	m.changed = ch
	if children != nil {
		m.hasChildren = true
		m.numChildren = int32(len(children))
		switch len(children) {
		case 0:
		case 1:
			m.child0 = children[0]
		case 2:
			m.child0 = children[0]
			m.child1 = children[1]
		case 3:
			m.child0 = children[0]
			m.child1 = children[1]
			m.child2 = children[2]
		default:
			m.extraChildren = children
		}
	}
	return m
}

// Java: MinimizedCondition.MeasuredNode.isChanged()
func (m *MinimizedCondition_MeasuredNode) IsChanged() bool {
	return m.changed
}

// Java: MinimizedCondition.MeasuredNode.isNot()
func (m *MinimizedCondition_MeasuredNode) IsNot() bool {
	return m.token() == rhino.Token_NOT
}

// Java: MinimizedCondition.MeasuredNode.withoutNot()
func (m *MinimizedCondition_MeasuredNode) WithoutNot() *MinimizedCondition_MeasuredNode {
	return m.withoutNotIn(m.mc)
}

func (m *MinimizedCondition_MeasuredNode) withoutNotIn(c *MinimizedCondition) *MinimizedCondition_MeasuredNode {
	jutil.CheckState(m.IsNot())
	if m.hasChildren {
		return m.child0.changeIn(c)
	}
	return minimizedCondition_MeasuredNode_forNodeIn(c, m.node.GetFirstChild()).changeIn(c)
}

// Java: MinimizedCondition.MeasuredNode.negate()
func (m *MinimizedCondition_MeasuredNode) Negate() *MinimizedCondition_MeasuredNode {
	return m.negateIn(m.mc)
}

func (m *MinimizedCondition_MeasuredNode) negateIn(c *MinimizedCondition) *MinimizedCondition_MeasuredNode {
	switch m.token() {
	case rhino.Token_EQ:
		return m.updateTokenIn(c, rhino.Token_NE)
	case rhino.Token_NE:
		return m.updateTokenIn(c, rhino.Token_EQ)
	case rhino.Token_SHEQ:
		return m.updateTokenIn(c, rhino.Token_SHNE)
	case rhino.Token_SHNE:
		return m.updateTokenIn(c, rhino.Token_SHEQ)
	case rhino.Token_NOT:
		return m.withoutNotIn(c)
	default:
		return m.addNotIn(c)
	}
}

// Java: MinimizedCondition.MeasuredNode.normalizeChildren(Node, MeasuredNode[])
func MinimizedCondition_MeasuredNode_NormalizeChildren(node rhino.Node, children []*MinimizedCondition_MeasuredNode) []*MinimizedCondition_MeasuredNode {
	if children != nil || !node.HasChildren() {
		return children
	} else {
		measuredChildren := make([]*MinimizedCondition_MeasuredNode, node.GetChildCount())
		child := 0
		for c := node.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			measuredChildren[child] = minimizedCondition_MeasuredNode_forNode(c)
			child++
		}
		return measuredChildren
	}
}

// Java: MinimizedCondition.MeasuredNode.updateToken(Token)
func (m *MinimizedCondition_MeasuredNode) UpdateToken(token rhino.Token) *MinimizedCondition_MeasuredNode {
	return m.updateTokenIn(m.mc, token)
}

func (m *MinimizedCondition_MeasuredNode) updateTokenIn(c *MinimizedCondition, token rhino.Token) *MinimizedCondition_MeasuredNode {
	out := c.allocNode()
	*out = *m
	out.mc = c
	out.virtualToken = token
	out.changed = true
	return out
}

// Java: MinimizedCondition.MeasuredNode.addNot()
func (m *MinimizedCondition_MeasuredNode) AddNot() *MinimizedCondition_MeasuredNode {
	return m.addNotIn(m.mc)
}

func (m *MinimizedCondition_MeasuredNode) addNotIn(c *MinimizedCondition) *MinimizedCondition_MeasuredNode {
	cost := m.length + 1
	if m.IsLowerPrecedenceThan(minimizedCondition_NOT_PRECEDENCE) {
		cost += 2
	}
	out := c.allocNode()
	out.node = m.node
	out.virtualToken = rhino.Token_NOT
	out.length = cost
	out.changed = true
	out.hasChildren = true
	out.numChildren = 1
	out.child0 = m
	return out
}

// Java: MinimizedCondition.MeasuredNode.change()
func (m *MinimizedCondition_MeasuredNode) Change() *MinimizedCondition_MeasuredNode {
	return m.changeIn(m.mc)
}

func (m *MinimizedCondition_MeasuredNode) changeIn(c *MinimizedCondition) *MinimizedCondition_MeasuredNode {
	if m.IsChanged() {
		return m
	}
	out := c.allocNode()
	*out = *m
	out.mc = c
	out.changed = true
	return out
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
func minimizedCondition_MeasuredNode_estimateCostOneLevel(n rhino.Node, children ...*MinimizedCondition_MeasuredNode) int {
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
	return NodeUtil_Precedence(m.token()) < precedence
}

// minimizedCondition_MeasuredNode_addNode: the returned MeasuredNode is
// only marked as changed if the children are marked as changed.
//
// Java: MinimizedCondition.MeasuredNode.addNode(Node, MeasuredNode...)
func minimizedCondition_MeasuredNode_addNode(parent rhino.Node, children ...*MinimizedCondition_MeasuredNode) *MinimizedCondition_MeasuredNode {
	cost := 0
	changed := false
	for _, child := range children {
		cost += child.length
		changed = changed || child.changed
	}
	cost += minimizedCondition_MeasuredNode_estimateCostOneLevel(parent, children...)
	return NewMinimizedCondition_MeasuredNode(parent, children, cost, changed)
}

func minimizedCondition_MeasuredNode_addNode1(c *MinimizedCondition, parent rhino.Node, c0 *MinimizedCondition_MeasuredNode) *MinimizedCondition_MeasuredNode {
	cost := c0.length
	if parent.IsNot() {
		cost++
	}
	parentPrecedence := NodeUtil_Precedence(parent.GetToken())
	if c0.IsLowerPrecedenceThan(parentPrecedence) {
		cost += 2
	}
	m := c.allocNode()
	m.node = parent
	m.length = cost
	m.changed = c0.changed
	m.hasChildren = true
	m.numChildren = 1
	m.child0 = c0
	return m
}

func minimizedCondition_MeasuredNode_addNode2(c *MinimizedCondition, parent rhino.Node, virtualToken rhino.Token, c0, c1 *MinimizedCondition_MeasuredNode) *MinimizedCondition_MeasuredNode {
	cost := c0.length + c1.length
	tok := virtualToken
	if tok == 0 {
		tok = parent.GetToken()
	}
	if tok == rhino.Token_NOT {
		cost++
	}
	parentPrecedence := NodeUtil_Precedence(tok)
	if c0.IsLowerPrecedenceThan(parentPrecedence) {
		cost += 2
	}
	if c1.IsLowerPrecedenceThan(parentPrecedence) {
		cost += 2
	}
	m := c.allocNode()
	m.node = parent
	m.virtualToken = virtualToken
	m.length = cost
	m.changed = c0.changed || c1.changed
	m.hasChildren = true
	m.numChildren = 2
	m.child0 = c0
	m.child1 = c1
	return m
}

func minimizedCondition_MeasuredNode_addNode3(c *MinimizedCondition, parent rhino.Node, c0, c1, c2 *MinimizedCondition_MeasuredNode) *MinimizedCondition_MeasuredNode {
	cost := c0.length + c1.length + c2.length
	if parent.IsNot() {
		cost++
	}
	parentPrecedence := NodeUtil_Precedence(parent.GetToken())
	if c0.IsLowerPrecedenceThan(parentPrecedence) {
		cost += 2
	}
	if c1.IsLowerPrecedenceThan(parentPrecedence) {
		cost += 2
	}
	if c2.IsLowerPrecedenceThan(parentPrecedence) {
		cost += 2
	}
	m := c.allocNode()
	m.node = parent
	m.length = cost
	m.changed = c0.changed || c1.changed || c2.changed
	m.hasChildren = true
	m.numChildren = 3
	m.child0 = c0
	m.child1 = c1
	m.child2 = c2
	return m
}

// minimizedCondition_MeasuredNode_forNode returns a MeasuredNode for a
// non-participating AST Node. This is used for leaf expression nodes.
//
// Java: MinimizedCondition.MeasuredNode.forNode(Node)
func minimizedCondition_MeasuredNode_forNode(n rhino.Node) *MinimizedCondition_MeasuredNode {
	return minimizedCondition_MeasuredNode_forNodeIn(nil, n)
}

func minimizedCondition_MeasuredNode_forNodeIn(c *MinimizedCondition, n rhino.Node) *MinimizedCondition_MeasuredNode {
	m := c.allocNode()
	m.node = n
	return m
}

// WillChange reports whether the MeasuredNode is a change from the
// original. This can either be a change within the original AST tree or a
// replacement of the original node.
//
// Java: MinimizedCondition.MeasuredNode.willChange(Node)
func (m *MinimizedCondition_MeasuredNode) WillChange(original rhino.Node) bool {
	rhino.CheckNotNull(original)
	return m.virtualToken != 0 || original != m.node || m.IsChanged()
}

// ApplyTo updates the AST for the result of this MeasuredNode. This can
// either be a change within the original AST tree or a replacement of the
// original node.
//
// Java: MinimizedCondition.MeasuredNode.applyTo(Node)
func (m *MinimizedCondition_MeasuredNode) ApplyTo(original rhino.Node) rhino.Node {
	rhino.CheckNotNull(original)
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
func (m *MinimizedCondition_MeasuredNode) SafeDetach(n rhino.Node) rhino.Node {
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
func (m *MinimizedCondition_MeasuredNode) BuildReplacement() rhino.Node {
	origNode := m.node
	targetNode := origNode
	if m.virtualToken != 0 {
		targetNode = rhino.NewNode(m.virtualToken).Srcref(origNode)
		m.node = targetNode
		m.virtualToken = 0
		if !m.hasChildren {
			for origNode.HasChildren() {
				targetNode.AddChildToBack(origNode.RemoveFirstChild())
			}
			return targetNode
		}
	}
	if m.hasChildren {
		targetNode.DetachChildren()
		switch m.numChildren {
		case 0:
		case 1:
			targetNode.AddChildToBack(m.SafeDetach(m.child0.BuildReplacement()))
		case 2:
			targetNode.AddChildToBack(m.SafeDetach(m.child0.BuildReplacement()))
			targetNode.AddChildToBack(m.SafeDetach(m.child1.BuildReplacement()))
		case 3:
			targetNode.AddChildToBack(m.SafeDetach(m.child0.BuildReplacement()))
			targetNode.AddChildToBack(m.SafeDetach(m.child1.BuildReplacement()))
			targetNode.AddChildToBack(m.SafeDetach(m.child2.BuildReplacement()))
		default:
			for _, child := range m.extraChildren {
				replacementChild := m.SafeDetach(child.BuildReplacement())
				targetNode.AddChildToBack(replacementChild)
			}
		}
	}
	return targetNode
}
