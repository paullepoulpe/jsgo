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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type measuredID int32

const unoptimizedNegID measuredID = -1

type measuredNode struct {
	node         rhino.Node
	children     [3]measuredID
	length       int32
	virtualToken rhino.Token
	numChildren  uint8
	changed      bool
	hasChildren  bool
}

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
	positive measuredID

	// A representation equivalent to the negation of the original
	// condition.
	negative measuredID

	nodes    []measuredNode
	nodesBuf [32]measuredNode
}

var minimizedCondition_NOT_PRECEDENCE = NodeUtil_Precedence(rhino.Token_NOT)

func (c *MinimizedCondition) reset() {
	if c.nodes == nil {
		c.nodes = c.nodesBuf[:0]
	} else {
		c.nodes = c.nodes[:0]
	}
	c.positive = unoptimizedNegID
	c.negative = unoptimizedNegID
}

func (c *MinimizedCondition) allocNode(mn measuredNode) measuredID {
	id := measuredID(len(c.nodes))
	c.nodes = append(c.nodes, mn)
	return id
}

func (c *MinimizedCondition) handle(id measuredID) MinimizedCondition_MeasuredNode {
	mn := &c.nodes[id]
	return MinimizedCondition_MeasuredNode{
		mc:      c,
		id:      id,
		length:  int(mn.length),
		changed: mn.changed,
	}
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
		c.positive, c.negative = c.computeMinimizedCondition(n)
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
func (c *MinimizedCondition) GetMinimized(style MinimizedCondition_MinimizationStyle) MinimizedCondition_MeasuredNode {
	if style == MinimizedCondition_MinimizationStyle_PREFER_UNNEGATED ||
		c.negative == unoptimizedNegID ||
		c.token(c.positive) == rhino.Token_NOT ||
		c.nodes[c.positive].length <= c.nodes[c.negative].length {
		return c.handle(c.positive)
	}
	return c.handle(c.addNot(c.negative))
}

func MinimizedCondition_UnoptimizedInto(n rhino.Node, c *MinimizedCondition) *MinimizedCondition {
	rhino.CheckNotNull(n.GetParent())
	c.reset()
	c.positive = c.forNode(n)
	c.negative = unoptimizedNegID
	return c
}

func (c *MinimizedCondition) pickBest(a, b measuredID) measuredID {
	la, lb := c.nodes[a].length, c.nodes[b].length
	if la == lb {
		if c.nodes[b].changed {
			return a
		}
		return b
	}
	if la < lb {
		return a
	}
	return b
}

func (c *MinimizedCondition) computeMinimizedCondition(n rhino.Node) (measuredID, measuredID) {
	switch n.GetToken() {
	case rhino.Token_NOT:
		subPos, subNeg := c.computeMinimizedCondition(n.GetFirstChild())
		positive := c.pickBest(c.addNode1(n, subPos), subNeg)
		negative := c.pickBest(
			c.negate(subNeg), // since parent node `n` is a NOT, we need to negate the subtree's
			// computed `negative` to obtain the parent `n`'s real negative.
			subPos)
		return positive, c.change(negative)
	case rhino.Token_AND, rhino.Token_OR:
		var complementToken rhino.Token
		if n.IsAnd() {
			complementToken = rhino.Token_OR
		} else {
			complementToken = rhino.Token_AND
		}
		leftPos, leftNeg := c.computeMinimizedCondition(n.GetFirstChild())
		rightPos, rightNeg := c.computeMinimizedCondition(n.GetLastChild())
		posTree := c.addNode2(n, 0, leftPos, rightPos)
		negComplementTree := c.addNode2(n, complementToken, leftNeg, rightNeg)
		positive := c.pickBest(
			posTree,
			c.negate(negComplementTree))
		negative := c.pickBest(
			c.negate(posTree),
			c.change(negComplementTree))
		return positive, c.change(negative)
	case rhino.Token_HOOK:
		cond := n.GetFirstChild()
		thenNode := cond.GetNext()
		elseNode := thenNode.GetNext()
		thenPos, thenNeg := c.computeMinimizedCondition(thenNode)
		elsePos, elseNeg := c.computeMinimizedCondition(elseNode)
		condMeasured := c.forNode(cond)
		positive := c.addNode3(n, condMeasured, thenPos, elsePos)
		negative := c.addNode3(n, condMeasured, thenNeg, elseNeg)
		return positive, c.change(negative)
	case rhino.Token_COMMA:
		lhs := n.GetFirstChild()
		rhsPos, rhsNeg := c.computeMinimizedCondition(lhs.GetNext())
		lhsMeasured := c.forNode(lhs)
		positive := c.addNode2(n, 0, lhsMeasured, rhsPos)
		negative := c.addNode2(n, 0, lhsMeasured, rhsNeg)
		return positive, c.change(negative)
	default:
		pos := c.forNode(n)
		neg := c.negate(pos)
		return pos, c.change(neg)
	}
}

// MinimizedCondition_MeasuredNode is a lightweight handle to a measured node
// inside a MinimizedCondition arena.
type MinimizedCondition_MeasuredNode struct {
	mc      *MinimizedCondition
	id      measuredID
	length  int
	changed bool
}

func (c *MinimizedCondition) token(id measuredID) rhino.Token {
	mn := &c.nodes[id]
	if mn.virtualToken != 0 {
		return mn.virtualToken
	}
	return mn.node.GetToken()
}

func (c *MinimizedCondition) isLowerPrecedenceThan(id measuredID, precedence int) bool {
	return NodeUtil_Precedence(c.token(id)) < precedence
}

func (c *MinimizedCondition) forNode(n rhino.Node) measuredID {
	return c.allocNode(measuredNode{node: n})
}

func (c *MinimizedCondition) change(id measuredID) measuredID {
	if c.nodes[id].changed {
		return id
	}
	mn := c.nodes[id]
	mn.changed = true
	return c.allocNode(mn)
}

func (c *MinimizedCondition) updateToken(id measuredID, token rhino.Token) measuredID {
	mn := c.nodes[id]
	mn.virtualToken = token
	mn.changed = true
	return c.allocNode(mn)
}

func (c *MinimizedCondition) addNot(id measuredID) measuredID {
	mn := c.nodes[id]
	cost := mn.length + 1
	if c.isLowerPrecedenceThan(id, minimizedCondition_NOT_PRECEDENCE) {
		cost += 2
	}
	return c.allocNode(measuredNode{
		node:         mn.node,
		children:     [3]measuredID{id},
		length:       cost,
		virtualToken: rhino.Token_NOT,
		numChildren:  1,
		changed:      true,
		hasChildren:  true,
	})
}

func (c *MinimizedCondition) withoutNot(id measuredID) measuredID {
	jutil.CheckState(c.token(id) == rhino.Token_NOT)
	mn := c.nodes[id]
	if mn.hasChildren {
		return c.change(mn.children[0])
	}
	return c.change(c.forNode(mn.node.GetFirstChild()))
}

func (c *MinimizedCondition) negate(id measuredID) measuredID {
	switch c.token(id) {
	case rhino.Token_EQ:
		return c.updateToken(id, rhino.Token_NE)
	case rhino.Token_NE:
		return c.updateToken(id, rhino.Token_EQ)
	case rhino.Token_SHEQ:
		return c.updateToken(id, rhino.Token_SHNE)
	case rhino.Token_SHNE:
		return c.updateToken(id, rhino.Token_SHEQ)
	case rhino.Token_NOT:
		return c.withoutNot(id)
	default:
		return c.addNot(id)
	}
}

func (c *MinimizedCondition) addNode1(parent rhino.Node, c0 measuredID) measuredID {
	n0 := c.nodes[c0]
	cost := n0.length
	if parent.IsNot() {
		cost++
	}
	parentPrecedence := NodeUtil_Precedence(parent.GetToken())
	if c.isLowerPrecedenceThan(c0, parentPrecedence) {
		cost += 2
	}
	return c.allocNode(measuredNode{
		node:        parent,
		children:    [3]measuredID{c0},
		length:      cost,
		numChildren: 1,
		changed:     n0.changed,
		hasChildren: true,
	})
}

func (c *MinimizedCondition) addNode2(parent rhino.Node, virtualToken rhino.Token, c0, c1 measuredID) measuredID {
	n0, n1 := c.nodes[c0], c.nodes[c1]
	cost := n0.length + n1.length
	tok := virtualToken
	if tok == 0 {
		tok = parent.GetToken()
	}
	if tok == rhino.Token_NOT {
		cost++
	}
	parentPrecedence := NodeUtil_Precedence(tok)
	if c.isLowerPrecedenceThan(c0, parentPrecedence) {
		cost += 2
	}
	if c.isLowerPrecedenceThan(c1, parentPrecedence) {
		cost += 2
	}
	return c.allocNode(measuredNode{
		node:         parent,
		children:     [3]measuredID{c0, c1},
		length:       cost,
		virtualToken: virtualToken,
		numChildren:  2,
		changed:      n0.changed || n1.changed,
		hasChildren:  true,
	})
}

func (c *MinimizedCondition) addNode3(parent rhino.Node, c0, c1, c2 measuredID) measuredID {
	n0, n1, n2 := c.nodes[c0], c.nodes[c1], c.nodes[c2]
	cost := n0.length + n1.length + n2.length
	if parent.IsNot() {
		cost++
	}
	parentPrecedence := NodeUtil_Precedence(parent.GetToken())
	if c.isLowerPrecedenceThan(c0, parentPrecedence) {
		cost += 2
	}
	if c.isLowerPrecedenceThan(c1, parentPrecedence) {
		cost += 2
	}
	if c.isLowerPrecedenceThan(c2, parentPrecedence) {
		cost += 2
	}
	return c.allocNode(measuredNode{
		node:        parent,
		children:    [3]measuredID{c0, c1, c2},
		length:      cost,
		numChildren: 3,
		changed:     n0.changed || n1.changed || n2.changed,
		hasChildren: true,
	})
}

// Java: MinimizedCondition.MeasuredNode.isChanged()
func (m MinimizedCondition_MeasuredNode) IsChanged() bool {
	return m.changed
}

// Java: MinimizedCondition.MeasuredNode.isNot()
func (m MinimizedCondition_MeasuredNode) IsNot() bool {
	return m.mc.token(m.id) == rhino.Token_NOT
}

// Java: MinimizedCondition.MeasuredNode.withoutNot()
func (m MinimizedCondition_MeasuredNode) WithoutNot() MinimizedCondition_MeasuredNode {
	return m.mc.handle(m.mc.withoutNot(m.id))
}

// IsLowerPrecedenceThan reports whether the node type has lower precedence
// than "precedence".
//
// Java: MinimizedCondition.MeasuredNode.isLowerPrecedenceThan(int)
func (m MinimizedCondition_MeasuredNode) IsLowerPrecedenceThan(precedence int) bool {
	return m.mc.isLowerPrecedenceThan(m.id, precedence)
}

// WillChange reports whether the MeasuredNode is a change from the
// original. This can either be a change within the original AST tree or a
// replacement of the original node.
//
// Java: MinimizedCondition.MeasuredNode.willChange(Node)
func (m MinimizedCondition_MeasuredNode) WillChange(original rhino.Node) bool {
	rhino.CheckNotNull(original)
	mn := &m.mc.nodes[m.id]
	return mn.virtualToken != 0 || original != mn.node || mn.changed
}

// ApplyTo updates the AST for the result of this MeasuredNode. This can
// either be a change within the original AST tree or a replacement of the
// original node.
//
// Java: MinimizedCondition.MeasuredNode.applyTo(Node)
func (m MinimizedCondition_MeasuredNode) ApplyTo(original rhino.Node) rhino.Node {
	rhino.CheckNotNull(original)
	jutil.CheckState(m.WillChange(original))
	replacement := m.BuildReplacement()
	if original != replacement {
		minimizedConditionSafeDetach(replacement)
		original.ReplaceWith(replacement)
	}
	return replacement
}

func minimizedConditionSafeDetach(n rhino.Node) rhino.Node {
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
func (m MinimizedCondition_MeasuredNode) BuildReplacement() rhino.Node {
	return m.mc.buildReplacement(m.id)
}

func (c *MinimizedCondition) buildReplacement(id measuredID) rhino.Node {
	mn := c.nodes[id]
	origNode := mn.node
	targetNode := origNode
	if mn.virtualToken != 0 {
		targetNode = rhino.NewNode(mn.virtualToken).Srcref(origNode)
		c.nodes[id].node = targetNode
		c.nodes[id].virtualToken = 0
		if !mn.hasChildren {
			for origNode.HasChildren() {
				targetNode.AddChildToBack(origNode.RemoveFirstChild())
			}
			return targetNode
		}
	}
	if mn.hasChildren {
		targetNode.DetachChildren()
		for i := uint8(0); i < mn.numChildren; i++ {
			targetNode.AddChildToBack(minimizedConditionSafeDetach(c.buildReplacement(mn.children[i])))
		}
	}
	return targetNode
}
