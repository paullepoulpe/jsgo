// Port of com/google/javascript/jscomp/DeadPropertyAssignmentElimination.java @ closure-compiler v20260415

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
	"iter"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// DeadPropertyAssignmentElimination is an optimization pass that finds and
// removes dead property assignments within functions and classes.
//
// This pass does not currently use the control-flow graph. It makes the
// following assumptions:
//
//   - Functions with inner functions are not processed.
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
func (d *DeadPropertyAssignmentElimination) Process(externs *rhino.Node, root *rhino.Node) {
	// GatherExternProperties must be enabled for this pass to safely know what property writes are
	// eligible for removal.
	if d.compiler.GetExternProperties() == nil || d.compiler.GetAccessorSummary() == nil {
		return
	}

	var skiplistedPropNames jutil.Set[string] = &deadPropertyAssignmentElimination_union{
		set1: d.compiler.GetAccessorSummary().GetAccessors(), set2: d.compiler.GetExternProperties()}

	NodeTraversal_Traverse(
		d.compiler, d.compiler.GetJsRoot(), NewDeadPropertyAssignmentElimination_FunctionVisitor(skiplistedPropNames))
}

// deadPropertyAssignmentElimination_union is the view returned by
// Sets.union(accessors.keySet(), externProperties).
//
// Go-only: language-adapter — Guava's Sets.union view (the first set is a
// map's key set).
type deadPropertyAssignmentElimination_union struct {
	set1 *jutil.LinkedHashMap[string, AccessorSummary_PropertyAccessKind]
	set2 *jutil.LinkedHashSet[string]
}

func (u *deadPropertyAssignmentElimination_union) Size() int {
	size := u.set1.Size()
	for e := range u.set2.All() {
		if !u.set1.ContainsKey(e) {
			size++
		}
	}
	return size
}

func (u *deadPropertyAssignmentElimination_union) IsEmpty() bool {
	return u.set1.IsEmpty() && u.set2.IsEmpty()
}

func (u *deadPropertyAssignmentElimination_union) Contains(e string) bool {
	return u.set1.ContainsKey(e) || u.set2.Contains(e)
}

func (u *deadPropertyAssignmentElimination_union) All() iter.Seq[string] {
	return func(yield func(string) bool) {
		for e := range u.set1.Keys() {
			if !yield(e) {
				return
			}
		}
		for e := range u.set2.All() {
			if !u.set1.ContainsKey(e) && !yield(e) {
				return
			}
		}
	}
}

// DeadPropertyAssignmentElimination_FunctionVisitor is the nested class
// FunctionVisitor.
type DeadPropertyAssignmentElimination_FunctionVisitor struct {
	NodeTraversal_AbstractChangedScopeCallbackBase

	// A set of properties names that are potentially unsafe to remove
	// duplicate writes to.
	skiplistedPropNames jutil.Set[string]
}

// Java: DeadPropertyAssignmentElimination.FunctionVisitor(Set<String>)
func NewDeadPropertyAssignmentElimination_FunctionVisitor(skiplistedPropNames jutil.Set[string]) *DeadPropertyAssignmentElimination_FunctionVisitor {
	v := &DeadPropertyAssignmentElimination_FunctionVisitor{}
	v.InitNodeTraversal_AbstractChangedScopeCallbackBase(v)
	v.skiplistedPropNames = skiplistedPropNames
	return v
}

// Java: DeadPropertyAssignmentElimination.FunctionVisitor.enterChangedScopeRoot(AbstractCompiler, Node)
func (v *DeadPropertyAssignmentElimination_FunctionVisitor) EnterChangedScopeRoot(compiler AbstractCompiler, root *rhino.Node) {
	if !root.IsFunction() {
		return
	}

	body := NodeUtil_GetFunctionBody(root)
	if !body.HasChildren() || NodeUtil_Has(body, (*rhino.Node).IsFunction, deadPropertyAssignmentElimination_alwaysTrue) {
		return
	}

	traversal := NewDeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal(v.skiplistedPropNames)
	NodeTraversal_Traverse(compiler, body, traversal)

	// Any candidate property assignment can have a write removed if that write is never read
	// and it's written to at least one more time.
	for property := range traversal.propertyMap.Values() {
		if property.writes.Size() <= 1 {
			continue
		}
		// Iterators.peekingIterator(property.writes.iterator())
		writes := property.writes.Slice()
		for i := 0; i < len(writes); i++ {
			propertyWrite := writes[i]
			if i+1 < len(writes) && propertyWrite.IsSafeToRemove(writes[i+1]) {
				lhs := propertyWrite.assignedAt
				rhs := lhs.GetNext()
				assignNode := lhs.GetParent()
				if assignNode.IsAssign() {
					// replace "a.b.c = <expr>" with "<expr>"
					rhs.Detach()
					assignNode.ReplaceWith(rhs)
					compiler.ReportChangeToEnclosingScope(rhs)
				} else {
					jutil.CheckState(NodeUtil_IsAssignmentOp(assignNode))
					// replace "a.b.c += <expr>" with "a.b.c + expr"
					opType := NodeUtil_GetOpFromAssignmentOp(assignNode)
					assignNode.SetToken(opType)
					compiler.ReportChangeToEnclosingScope(assignNode)
				}
			}
		}
	}
}

// deadPropertyAssignmentElimination_alwaysTrue is Guava's
// Predicates.alwaysTrue().
func deadPropertyAssignmentElimination_alwaysTrue(*rhino.Node) bool { return true }

// DeadPropertyAssignmentElimination_Property is the nested class Property.
type DeadPropertyAssignmentElimination_Property struct {
	name string

	// This pass doesn't use a control-flow graph; this field contains a rough
	// approximation of the control flow. For writes in the same block, they
	// appear in this list in program-execution order.
	// All writes in a list are to the same property name, but the full
	// qualified names may differ, eg, a.b.c and e.d.c can be in the list.
	// Consecutive writes to the same qname may mean that the first write can
	// be removed (see isSafeToRemove).
	//
	// Go-only: language-adapter — the ArrayDeque is only appended to
	// (addLast) and read from the back (getLast) or front to back, so it is
	// an ArrayList.
	writes *jutil.ArrayList[*DeadPropertyAssignmentElimination_PropertyWrite]

	children *jutil.LinkedHashSet[*DeadPropertyAssignmentElimination_Property]
}

// Java: DeadPropertyAssignmentElimination.Property(String)
func NewDeadPropertyAssignmentElimination_Property(name string) *DeadPropertyAssignmentElimination_Property {
	return &DeadPropertyAssignmentElimination_Property{
		name:     name,
		writes:   jutil.NewArrayList[*DeadPropertyAssignmentElimination_PropertyWrite](),
		children: jutil.NewLinkedHashSet[*DeadPropertyAssignmentElimination_Property](),
	}
}

// Java: DeadPropertyAssignmentElimination.Property.markLastWriteRead()
func (p *DeadPropertyAssignmentElimination_Property) MarkLastWriteRead() {
	if !p.writes.IsEmpty() {
		p.writes.GetLast().MarkRead()
	}
}

// MarkChildrenRead marks all children of this property as read.
//
// Java: DeadPropertyAssignmentElimination.Property.markChildrenRead()
func (p *DeadPropertyAssignmentElimination_Property) MarkChildrenRead() {
	// If a property is in propertiesSet, it has been added to the queue and processed,
	// it will not be added to the queue again.
	propertiesSet := jutil.NewLinkedHashSet(p.children.ToSlice()...)
	// Go-only: language-adapter — the ArrayDeque used as a FIFO queue is a
	// slice with a head index.
	propertyQueue := propertiesSet.ToSlice()

	// Ensure we don't process ourselves.
	propertiesSet.Add(p)

	for head := 0; head < len(propertyQueue); head++ {
		childProperty := propertyQueue[head]
		childProperty.MarkLastWriteRead()
		for grandchildProperty := range childProperty.children.All() {
			if propertiesSet.Add(grandchildProperty) {
				propertyQueue = append(propertyQueue, grandchildProperty)
			}
		}
	}
}

// Java: DeadPropertyAssignmentElimination.Property.addWrite(Node)
func (p *DeadPropertyAssignmentElimination_Property) AddWrite(lhs *rhino.Node) {
	jutil.CheckArgument(lhs.IsQualifiedName())
	p.writes.Add(NewDeadPropertyAssignmentElimination_PropertyWrite(lhs))
}

// Java: DeadPropertyAssignmentElimination.Property.toString()
func (p *DeadPropertyAssignmentElimination_Property) String() string {
	return "Property " + p.name
}

// DeadPropertyAssignmentElimination_PropertyWrite is the nested class
// PropertyWrite.
type DeadPropertyAssignmentElimination_PropertyWrite struct {
	assignedAt    *rhino.Node
	isRead        bool
	qualifiedName *string
}

// Java: DeadPropertyAssignmentElimination.PropertyWrite(Node)
func NewDeadPropertyAssignmentElimination_PropertyWrite(assignedAt *rhino.Node) *DeadPropertyAssignmentElimination_PropertyWrite {
	w := &DeadPropertyAssignmentElimination_PropertyWrite{}
	jutil.CheckArgument(assignedAt.IsQualifiedName())
	w.assignedAt = assignedAt
	w.qualifiedName = assignedAt.GetQualifiedName()
	return w
}

// Java: DeadPropertyAssignmentElimination.PropertyWrite.isSafeToRemove(PropertyWrite)
func (w *DeadPropertyAssignmentElimination_PropertyWrite) IsSafeToRemove(nextWrite *DeadPropertyAssignmentElimination_PropertyWrite) bool {
	return !w.isRead && nextWrite != nil && deadPropertyAssignmentElimination_stringsEqual(w.qualifiedName, nextWrite.qualifiedName)
}

// deadPropertyAssignmentElimination_stringsEqual is Objects.equals on two
// nullable Strings.
func deadPropertyAssignmentElimination_stringsEqual(a *string, b *string) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

// Java: DeadPropertyAssignmentElimination.PropertyWrite.markRead()
func (w *DeadPropertyAssignmentElimination_PropertyWrite) MarkRead() {
	w.isRead = true
}

// DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal is a
// NodeTraversal that operates within a function block and collects
// candidate properties assignments.
type DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal struct {
	// A map of property names to their nodes.
	//
	// Note: the references a.b and c.b will assume that it's the same b,
	// because a and c may be aliased, and we don't track aliasing.
	propertyMap *jutil.LinkedHashMap[string, *DeadPropertyAssignmentElimination_Property]

	// A set of properties names that are potentially unsafe to remove
	// duplicate writes to.
	skiplistedPropNames jutil.Set[string]
}

// Java: DeadPropertyAssignmentElimination.FindCandidateAssignmentTraversal(Set<String>)
func NewDeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal(skiplistedPropNames jutil.Set[string]) *DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal {
	return &DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal{
		propertyMap:         jutil.NewLinkedHashMap[string, *DeadPropertyAssignmentElimination_Property](),
		skiplistedPropNames: skiplistedPropNames,
	}
}

// GetOrCreateProperty gets a Property given the node that references it;
// the Property is created if it does not already exist.
//
// Returns a Property, or null if the provided node is not a qualified name.
//
// Java: DeadPropertyAssignmentElimination.FindCandidateAssignmentTraversal.getOrCreateProperty(Node)
func (f *DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal) GetOrCreateProperty(propNode *rhino.Node) *DeadPropertyAssignmentElimination_Property {
	if !propNode.IsQualifiedName() {
		return nil
	}

	var propName string
	if propNode.IsGetProp() {
		propName = propNode.GetString()
	} else {
		propName = *propNode.GetQualifiedName()
	}

	property := f.propertyMap.ComputeIfAbsent(propName, NewDeadPropertyAssignmentElimination_Property)

	/* Using the GETPROP chain, build out the tree of children properties.

	 For example, from a.b.c and a.c we can build:
	         a
	        / \
	       b   c
	      /
	     c

	Note: c is the same Property in this tree.
	*/
	if propNode.IsGetProp() {
		parentProperty := f.GetOrCreateProperty(propNode.GetFirstChild())
		if parentProperty != nil {
			parentProperty.children.Add(property)
		}
	}

	return property
}

// Java: DeadPropertyAssignmentElimination.FindCandidateAssignmentTraversal.shouldTraverse(NodeTraversal, Node, Node)
func (f *DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal) ShouldTraverse(nodeTraversal *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	return f.VisitNode(n, parent)
}

// Java: DeadPropertyAssignmentElimination.FindCandidateAssignmentTraversal.visit(NodeTraversal, Node, Node)
func (f *DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	// Visit the LHS of an assignment in post-order
	if NodeUtil_IsAssignmentOp(n) {
		f.VisitAssignmentLhs(n.GetFirstChild())
	}

	// Assume that all properties may be read when control flow leaves the function
	if NodeUtil_IsInvocation(n) || n.IsYield() || n.IsAwait() {
		f.MarkAllPropsRead()
	}

	// Mark all properties as read when leaving a block since we haven't proven that the block
	// will execute.
	if n.IsBlock() {
		f.VisitBlock(n)
	}
}

// Java: DeadPropertyAssignmentElimination.FindCandidateAssignmentTraversal.visitBlock(Node)
func (f *DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal) VisitBlock(blockNode *rhino.Node) {
	jutil.CheckArgument(blockNode.IsBlock())

	// We don't do flow analysis yet so we're going to assume everything written up to this
	// block is read.
	if blockNode.HasChildren() {
		f.MarkAllPropsRead()
	}
}

// Java: DeadPropertyAssignmentElimination.FindCandidateAssignmentTraversal.isConditionalExpression(Node)
func DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal_IsConditionalExpression(n *rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_AND, rhino.Token_OR, rhino.Token_HOOK, rhino.Token_COALESCE, rhino.Token_OPTCHAIN_CALL,
		rhino.Token_OPTCHAIN_GETELEM, rhino.Token_OPTCHAIN_GETPROP:
		return true
	default:
		return false
	}
}

// Java: DeadPropertyAssignmentElimination.FindCandidateAssignmentTraversal.visitAssignmentLhs(Node)
func (f *DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal) VisitAssignmentLhs(lhs *rhino.Node) {
	property := f.GetOrCreateProperty(lhs)

	if property == nil {
		return
	}

	if !lhs.IsGetProp() {
		property.MarkLastWriteRead()
		property.MarkChildrenRead()
		return
	}

	assignNode := lhs.GetParent()

	// If it's mutating assignment (+=, *=, etc.) then mark the last assignment read first.
	if !assignNode.IsAssign() {
		property.MarkLastWriteRead()
	}

	// Reassignment of a qualified name prefix might change what child properties are referenced
	// later on, so consider children properties as read.
	// Ex. a.b.c = 10; a.b = other; a.b.c = 20;
	property.MarkChildrenRead()
	property.AddWrite(lhs)

	// Now we need to go up the prop chain and mark those as read.
	child := lhs.GetFirstChild()
	for child != nil {
		childProperty := f.GetOrCreateProperty(child)
		if childProperty == nil {
			break
		}
		childProperty.MarkLastWriteRead()
		child = child.GetFirstChild()
	}
}

// Java: DeadPropertyAssignmentElimination.FindCandidateAssignmentTraversal.visitNode(Node, Node)
func (f *DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal) VisitNode(n *rhino.Node, parent *rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_GETPROP:
		// Handle potential getters/setters.
		if n.IsGetProp() && f.skiplistedPropNames.Contains(n.GetString()) {
			// We treat getters/setters as if they were a call, thus we mark all properties as read.
			f.MarkAllPropsRead()
			return true
		}

		if NodeUtil_IsAssignmentOp(parent) && parent.GetFirstChild() == n {
			// This is a write to the property, so skip the read handling below.
			// We'll handle this write in the visit() method
			// We need to continue traversing, because there could be a function call
			// child.
			return true
		}
		property := f.GetOrCreateProperty(n)
		if property != nil {
			// Mark all children properties as read.
			property.MarkLastWriteRead()

			// Only mark children properties as read if we're at the end of the referenced
			// property chain.
			// Ex. A read of "a.b.c" should mark a, a.b, a.b.c, and a.b.c.* as read, but not a.d
			if !parent.IsGetProp() {
				property.MarkChildrenRead()
			}
		}
		return true
	case rhino.Token_THIS, rhino.Token_NAME:
		nameProp := jutil.CheckNotNull(f.GetOrCreateProperty(n))
		nameProp.MarkLastWriteRead()
		if !parent.IsGetProp() {
			nameProp.MarkChildrenRead()
		}
		return true
	case rhino.Token_THROW, rhino.Token_FOR, rhino.Token_FOR_IN, rhino.Token_SWITCH:
		// TODO(kevinoconnor): Switch/for statements need special consideration since they may
		// execute out of order.
		f.MarkAllPropsRead()
		return false
	case rhino.Token_BLOCK:
		f.VisitBlock(n)
		return true
	default:
		if DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal_IsConditionalExpression(n) {
			f.MarkAllPropsRead()
			return false
		}
		return true
	}
}

// Java: DeadPropertyAssignmentElimination.FindCandidateAssignmentTraversal.markAllPropsRead()
func (f *DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal) MarkAllPropsRead() {
	for property := range f.propertyMap.Values() {
		if property.writes.IsEmpty() {
			continue
		}

		property.MarkLastWriteRead()
	}
}
