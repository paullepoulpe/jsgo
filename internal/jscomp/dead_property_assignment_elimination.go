// Port of com/google/javascript/jscomp/DeadPropertyAssignmentElimination.java @ closure-compiler v20261021

/*
 * Copyright 2016 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

// DeadPropertyAssignmentElimination is an optimization pass that finds and
// removes dead property assignments within functions and classes.
//
// This pass does not currently use the control-flow graph. It makes the
// following assumptions:
//
//   - Functions with inner functions or non-static class field initializers
//     are not processed.
//   - All properties are read whenever entering a block node. Dead
//     assignments within a block are processed.
//   - Hook nodes are not processed (it's assumed they read everything)
//   - Switch blocks are not processed (it's assumed they read everything)
//   - Any reference to a property getter/setter is treated like a call that
//     escapes all props.
//   - If there's an Object.definePropert{y,ies} call where the object or
//     property name is aliased then the optimization does not run at all.
//   - Properties names defined in externs will not be pruned.
type DeadPropertyAssignmentElimination struct {
	compiler AbstractCompiler
}

// Java: DeadPropertyAssignmentElimination(AbstractCompiler)
func NewDeadPropertyAssignmentElimination(compiler AbstractCompiler) *DeadPropertyAssignmentElimination {
	return &DeadPropertyAssignmentElimination{compiler: compiler}
}

// Java: DeadPropertyAssignmentElimination.process(Node, Node)
func (d *DeadPropertyAssignmentElimination) Process(externs rhino.Node, root rhino.Node) {
	// GatherExternProperties must be enabled for this pass to safely know what property writes are
	// eligible for removal.
	externProps := d.compiler.GetExternProperties()
	accessorSummary := d.compiler.GetAccessorSummary()
	if externProps == nil || accessorSummary == nil {
		return
	}

	accessorIDs := accessorSummary.AccessorIDs()
	skiplisted := make(map[uint32]struct{}, len(accessorIDs)+externProps.Size())
	for _, id := range accessorIDs {
		skiplisted[id] = struct{}{}
	}
	for name := range externProps.All() {
		if id := strpool.Lookup(name); id != 0 {
			skiplisted[uint32(id)] = struct{}{}
		}
	}

	visitor := &dpaeFunctionVisitor{
		traversal: dpaeTraversal{
			skiplisted:  skiplisted,
			thisID:      uint32(strpool.Intern("this")),
			superID:     uint32(strpool.Intern("super")),
			propertyMap: make(map[uint32]*deadProp, 32),
		},
	}
	visitor.InitNodeTraversal_AbstractChangedScopeCallbackBase(visitor)
	NodeTraversal_Traverse(d.compiler, d.compiler.GetJsRoot(), visitor)
}

type dpaeFunctionVisitor struct {
	NodeTraversal_AbstractChangedScopeCallbackBase
	traversal dpaeTraversal
}

// Java: DeadPropertyAssignmentElimination.FunctionVisitor.enterChangedScopeRoot(AbstractCompiler, Node)
func (v *dpaeFunctionVisitor) EnterChangedScopeRoot(compiler AbstractCompiler, root rhino.Node) {
	fn, ok := ast.AsFunction(root)
	if !ok {
		return
	}

	body := fn.Body().Node()
	if dpaeShouldSkipBody(body) {
		return
	}

	t := &v.traversal
	t.reset()
	NodeTraversal_Traverse(compiler, body, t)

	// Any candidate property assignment can have a write removed if that write is never read
	// and it's written to at least one more time on the same qualified name.
	for _, prop := range t.writtenProps {
		writes := prop.writes
		if len(writes) <= 1 {
			continue
		}
		for i := 0; i+1 < len(writes); i++ {
			if writes[i].isRead || !writes[i].assignedAt.MatchesQualifiedNameNode(writes[i+1].assignedAt) {
				continue
			}
			lhs := writes[i].assignedAt
			assign := ast.MustAssignment(lhs.GetParent())
			if plainAssign, ok := assign.AsAssign(); ok {
				// replace "a.b.c = <expr>" with "<expr>"
				rhs := plainAssign.Value().Detach().Node()
				plainAssign.Node().ReplaceWith(rhs)
				compiler.ReportChangeToEnclosingScope(rhs)
			} else {
				// replace "a.b.c += <expr>" with "a.b.c + expr"
				assignNode := assign.Node()
				assignNode.SetToken(NodeUtil_GetOpFromAssignmentOp(assignNode))
				compiler.ReportChangeToEnclosingScope(assignNode)
			}
		}
	}
}

func dpaeShouldSkipBody(body rhino.Node) bool {
	if !body.HasChildren() {
		return true
	}
	assignCount := 0
	if dpaeHasBailoutAndCountAssigns(body, &assignCount) {
		return true
	}
	return assignCount < 2
}

func dpaeHasBailoutAndCountAssigns(n rhino.Node, assignCount *int) bool {
	if NodeUtil_IsFunctionOrNonStaticClassFieldInitializer(n) {
		return true
	}
	if assign, ok := ast.AsAssignment(n); ok {
		if _, isGetProp := assign.Target().AsGetProp(); isGetProp {
			*assignCount++
		}
	}
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if dpaeHasBailoutAndCountAssigns(c, assignCount) {
			return true
		}
	}
	return false
}

type deadPropWrite struct {
	assignedAt rhino.Node
	isRead     bool
}

type deadProp struct {
	id            uint32
	writes        []deadPropWrite
	children      []*deadProp
	lastParent    *deadProp
	lastVisited   uint32
	enqueuedWrite bool
}

func (p *deadProp) markLastWriteRead() {
	if n := len(p.writes); n > 0 {
		p.writes[n-1].isRead = true
	}
}

func (p *deadProp) addChild(child *deadProp, newlyCreated bool) {
	if newlyCreated {
		p.children = append(p.children, child)
		child.lastParent = p
		return
	}
	if child.lastParent == p {
		return
	}
	for _, existing := range p.children {
		if existing == child {
			child.lastParent = p
			return
		}
	}
	p.children = append(p.children, child)
	child.lastParent = p
}

// dpaeTraversal collects candidate property assignments within a single function body.
type dpaeTraversal struct {
	// A set of property name IDs that are potentially unsafe to remove duplicate writes to.
	skiplisted map[uint32]struct{}
	thisID     uint32
	superID    uint32

	// A map of property name IDs to their nodes.
	// Note: the references a.b and c.b assume that it's the same b,
	// because a and c may be aliased, and we don't track aliasing.
	propertyMap  map[uint32]*deadProp
	props        []*deadProp
	numProps     int
	writtenProps []*deadProp
	bfsQueue     []*deadProp
	visitStamp   uint32
}

func (f *dpaeTraversal) reset() {
	clear(f.propertyMap)
	f.numProps = 0
	f.writtenProps = f.writtenProps[:0]
	f.visitStamp = 0
}

func (f *dpaeTraversal) propForID(id uint32) (*deadProp, bool) {
	if p := f.propertyMap[id]; p != nil {
		return p, false
	}
	var p *deadProp
	if f.numProps < len(f.props) {
		p = f.props[f.numProps]
		p.id = id
		p.writes = p.writes[:0]
		p.children = p.children[:0]
		p.lastParent = nil
		p.lastVisited = 0
		p.enqueuedWrite = false
	} else {
		p = &deadProp{id: id}
		f.props = append(f.props, p)
	}
	f.numProps++
	f.propertyMap[id] = p
	return p, true
}

func (f *dpaeTraversal) addWrite(p *deadProp, lhs rhino.Node) {
	if !p.enqueuedWrite {
		p.enqueuedWrite = true
		f.writtenProps = append(f.writtenProps, p)
	}
	p.writes = append(p.writes, deadPropWrite{assignedAt: lhs})
}

// getOrCreateProperty returns the deadProp referenced by propNode, or nil if
// propNode is not a qualified name.
func (f *dpaeTraversal) getOrCreateProperty(propNode rhino.Node) *deadProp {
	switch propNode.GetToken() {
	case rhino.Token_GETPROP:
		/* Using the GETPROP chain, build out the tree of children properties.

		 For example, from a.b.c and a.c we can build:
		         a
		        / \
		       b   c
		      /
		     c

		Note: c is the same Property in this tree.
		*/
		parentProp := f.getOrCreateProperty(ast.MustGetProp(propNode).Object().Node())
		if parentProp == nil {
			return nil
		}
		prop, created := f.propForID(propNode.GetStringID())
		parentProp.addChild(prop, created)
		return prop
	case rhino.Token_NAME:
		id := propNode.GetStringID()
		if id <= 1 {
			return nil
		}
		prop, _ := f.propForID(id)
		return prop
	case rhino.Token_THIS:
		prop, _ := f.propForID(f.thisID)
		return prop
	case rhino.Token_SUPER:
		prop, _ := f.propForID(f.superID)
		return prop
	default:
		return nil
	}
}

func (f *dpaeTraversal) markChildrenRead(p *deadProp) {
	children := p.children
	if len(children) == 0 {
		return
	}
	hasGrandchildren := false
	for _, child := range children {
		if len(child.children) != 0 {
			hasGrandchildren = true
			break
		}
	}
	if !hasGrandchildren {
		for _, child := range children {
			child.markLastWriteRead()
		}
		return
	}

	f.visitStamp++
	stamp := f.visitStamp
	if stamp == 0 {
		for i := 0; i < f.numProps; i++ {
			f.props[i].lastVisited = 0
		}
		f.visitStamp = 1
		stamp = 1
	}

	queue := f.bfsQueue[:0]
	for _, child := range children {
		child.lastVisited = stamp
		queue = append(queue, child)
	}
	// Ensure we don't re-enqueue ourselves from grandchildren.
	p.lastVisited = stamp

	for head := 0; head < len(queue); head++ {
		cur := queue[head]
		cur.markLastWriteRead()
		for _, grandchild := range cur.children {
			if grandchild.lastVisited != stamp {
				grandchild.lastVisited = stamp
				queue = append(queue, grandchild)
			}
		}
	}
	f.bfsQueue = queue[:0]
}

// Java: DeadPropertyAssignmentElimination.FindCandidateAssignmentTraversal.shouldTraverse(NodeTraversal, Node, Node)
func (f *dpaeTraversal) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	return f.visitNode(n, parent)
}

// Java: DeadPropertyAssignmentElimination.FindCandidateAssignmentTraversal.visit(NodeTraversal, Node, Node)
func (f *dpaeTraversal) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	// Visit the LHS of an assignment in post-order
	if assign, ok := ast.AsAssignment(n); ok {
		f.visitAssignmentLhs(assign)
	}

	// Assume that all properties may be read when control flow leaves the function
	if NodeUtil_IsInvocation(n) || n.IsYield() || n.IsAwait() {
		f.markAllPropsRead()
	}

	// Mark all properties as read when leaving a block since we haven't proven that the block
	// will execute.
	if n.IsBlock() && n.HasChildren() {
		f.markAllPropsRead()
	}
}

func dpaeIsConditionalExpression(tok rhino.Token) bool {
	switch tok {
	case rhino.Token_AND, rhino.Token_OR, rhino.Token_HOOK, rhino.Token_COALESCE, rhino.Token_OPTCHAIN_CALL,
		rhino.Token_OPTCHAIN_GETELEM, rhino.Token_OPTCHAIN_GETPROP:
		return true
	default:
		return false
	}
}

func (f *dpaeTraversal) visitAssignmentLhs(assign ast.Assignment) {
	lhs := assign.Target().Node()
	prop := f.getOrCreateProperty(lhs)
	if prop == nil {
		return
	}

	lhsGetProp, isGetProp := ast.AsGetProp(lhs)
	if !isGetProp {
		prop.markLastWriteRead()
		f.markChildrenRead(prop)
		return
	}

	// If it's mutating assignment (+=, *=, etc.) then mark the last assignment read first.
	if _, isPlainAssign := assign.AsAssign(); !isPlainAssign {
		prop.markLastWriteRead()
	}

	// Reassignment of a qualified name prefix might change what child properties are referenced
	// later on, so consider children properties as read.
	// Ex. a.b.c = 10; a.b = other; a.b.c = 20;
	f.markChildrenRead(prop)
	f.addWrite(prop, lhs)

	// Now we need to go up the prop chain and mark those as read.
	for child := lhsGetProp.Object().Node(); !child.IsNil(); {
		var childID uint32
		var next rhino.Node
		switch child.GetToken() {
		case rhino.Token_GETPROP:
			childID = child.GetStringID()
			next = ast.MustGetProp(child).Object().Node()
		case rhino.Token_NAME:
			childID = child.GetStringID()
		case rhino.Token_THIS:
			childID = f.thisID
		case rhino.Token_SUPER:
			childID = f.superID
		default:
			return
		}
		childProp := f.propertyMap[childID]
		if childProp == nil {
			return
		}
		childProp.markLastWriteRead()
		child = next
	}
}

func (f *dpaeTraversal) visitNode(n rhino.Node, parent rhino.Node) bool {
	tok := n.GetToken()
	switch tok {
	case rhino.Token_GETPROP:
		// Handle potential getters/setters.
		if _, isSkiplisted := f.skiplisted[n.GetStringID()]; isSkiplisted {
			// We treat getters/setters as if they were a call, thus we mark all properties as read.
			f.markAllPropsRead()
			return true
		}

		if parentAssign, ok := ast.AsAssignment(parent); ok && parentAssign.Target().Is(n) {
			// This is a write to the property, so skip the read handling below.
			// We'll handle this write in the visit() method.
			// We need to continue traversing, because there could be a function call child.
			return true
		}
		if prop := f.getOrCreateProperty(n); prop != nil {
			prop.markLastWriteRead()

			// Only mark children properties as read if we're at the end of the referenced
			// property chain.
			// Ex. A read of "a.b.c" should mark a, a.b, a.b.c, and a.b.c.* as read, but not a.d
			if !parent.IsGetProp() {
				f.markChildrenRead(prop)
			}
		}
		return true
	case rhino.Token_THIS, rhino.Token_NAME:
		if nameProp := f.getOrCreateProperty(n); nameProp != nil {
			nameProp.markLastWriteRead()
			if !parent.IsGetProp() {
				f.markChildrenRead(nameProp)
			}
		}
		return true
	case rhino.Token_THROW, rhino.Token_FOR, rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF, rhino.Token_ITER_SPREAD, rhino.Token_ARRAY_PATTERN, rhino.Token_SWITCH:
		// Loops and switch statements may execute out of order, and implicit iteration operations
		// (for-of, for-await-of, spread, array destructuring) invoke arbitrary JS iterator
		// methods which can read properties.
		f.markAllPropsRead()
		return false
	case rhino.Token_OBJECT_PATTERN:
		for patProp := range ast.MustObjectPattern(n).Props() {
			if sk, ok := patProp.AsStringKey(); ok {
				if prop := f.propertyMap[sk.Node().GetStringID()]; prop != nil {
					prop.markLastWriteRead()
					f.markChildrenRead(prop)
				}
			} else {
				f.markAllPropsRead()
				break
			}
		}
		return true
	case rhino.Token_BLOCK:
		// We don't do flow analysis yet so we're going to assume everything written up to this
		// block is read.
		if n.HasChildren() {
			f.markAllPropsRead()
		}
		return true
	default:
		if dpaeIsConditionalExpression(tok) {
			f.markAllPropsRead()
			return false
		}
		return true
	}
}

func (f *dpaeTraversal) markAllPropsRead() {
	for _, prop := range f.writtenProps {
		prop.markLastWriteRead()
	}
}
