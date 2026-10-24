// Port of com/google/javascript/jscomp/DeadPropertyAssignmentElimination.java @ closure-compiler v20261015

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

	traversal DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal
}

// Java: DeadPropertyAssignmentElimination.FunctionVisitor(Set<String>)
func NewDeadPropertyAssignmentElimination_FunctionVisitor(skiplistedPropNames jutil.Set[string]) *DeadPropertyAssignmentElimination_FunctionVisitor {
	v := &DeadPropertyAssignmentElimination_FunctionVisitor{}
	v.InitNodeTraversal_AbstractChangedScopeCallbackBase(v)
	v.skiplistedPropNames = skiplistedPropNames
	return v
}

// Java: DeadPropertyAssignmentElimination.FunctionVisitor.enterChangedScopeRoot(AbstractCompiler, Node)
func (v *DeadPropertyAssignmentElimination_FunctionVisitor) EnterChangedScopeRoot(compiler AbstractCompiler, root rhino.Node) {
	if !root.IsFunction() {
		return
	}

	body := NodeUtil_GetFunctionBody(root)
	if dpaeShouldSkipBody(body) {
		return
	}

	traversal := &v.traversal
	traversal.reset(v.skiplistedPropNames)
	traversal.traverse(compiler, body)

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
	if NodeUtil_IsAssignmentOp(n) && n.GetFirstChild().IsGetProp() {
		*assignCount++
	}
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if dpaeHasBailoutAndCountAssigns(c, assignCount) {
			return true
		}
	}
	return false
}

// deadPropertyAssignmentElimination_alwaysTrue is Guava's
// Predicates.alwaysTrue().
func deadPropertyAssignmentElimination_alwaysTrue(rhino.Node) bool { return true }

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
	writes    jutil.ArrayList[*DeadPropertyAssignmentElimination_PropertyWrite]
	writesBuf [2]*DeadPropertyAssignmentElimination_PropertyWrite

	children jutil.LinkedHashSet[*DeadPropertyAssignmentElimination_Property]
}

// Java: DeadPropertyAssignmentElimination.Property(String)
func NewDeadPropertyAssignmentElimination_Property(name string) *DeadPropertyAssignmentElimination_Property {
	p := &DeadPropertyAssignmentElimination_Property{
		name: name,
	}
	p.writes.InitWithBuffer(p.writesBuf[:])
	return p
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
	if p.children.IsEmpty() {
		return
	}
	hasGrandchildren := false
	for childProperty := range p.children.All() {
		if !childProperty.children.IsEmpty() {
			hasGrandchildren = true
			break
		}
	}
	if !hasGrandchildren {
		for childProperty := range p.children.All() {
			childProperty.MarkLastWriteRead()
		}
		return
	}

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
func (p *DeadPropertyAssignmentElimination_Property) AddWrite(lhs rhino.Node) {
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
	assignedAt    rhino.Node
	isRead        bool
	qualifiedName *string
}

// Java: DeadPropertyAssignmentElimination.PropertyWrite(Node)
func NewDeadPropertyAssignmentElimination_PropertyWrite(assignedAt rhino.Node) *DeadPropertyAssignmentElimination_PropertyWrite {
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
	propertyMap      *jutil.LinkedHashMap[string, *DeadPropertyAssignmentElimination_Property]
	propertyMapStore jutil.LinkedHashMap[string, *DeadPropertyAssignmentElimination_Property]
	propertyMapSlots [16]jutil.LhmSlot[string, *DeadPropertyAssignmentElimination_Property]

	propSlab  []DeadPropertyAssignmentElimination_Property
	propBuf   [16]DeadPropertyAssignmentElimination_Property
	writeSlab []DeadPropertyAssignmentElimination_PropertyWrite
	writeBuf  [16]DeadPropertyAssignmentElimination_PropertyWrite

	// A set of properties names that are potentially unsafe to remove
	// duplicate writes to.
	skiplistedPropNames jutil.Set[string]
}

// Java: DeadPropertyAssignmentElimination.FindCandidateAssignmentTraversal(Set<String>)
func NewDeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal(skiplistedPropNames jutil.Set[string]) *DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal {
	f := &DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal{}
	f.reset(skiplistedPropNames)
	return f
}

func (f *DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal) reset(skiplistedPropNames jutil.Set[string]) {
	f.skiplistedPropNames = skiplistedPropNames
	if f.propertyMap == nil {
		f.propertyMapStore.InitWithBuffer(f.propertyMapSlots[:])
		f.propertyMap = &f.propertyMapStore
	} else {
		f.propertyMap.Clear()
	}
	clear(f.propBuf[:])
	f.propSlab = f.propBuf[:0:len(f.propBuf)]
	clear(f.writeBuf[:])
	f.writeSlab = f.writeBuf[:0:len(f.writeBuf)]
}

func (f *DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal) allocProperty(name string) *DeadPropertyAssignmentElimination_Property {
	if len(f.propSlab) == cap(f.propSlab) {
		nextCap := min(max(cap(f.propSlab)*2, 16), 256)
		f.propSlab = make([]DeadPropertyAssignmentElimination_Property, 0, nextCap)
	}
	idx := len(f.propSlab)
	f.propSlab = f.propSlab[:idx+1]
	p := &f.propSlab[idx]
	p.name = name
	p.writes.InitWithBuffer(p.writesBuf[:])
	return p
}

func (f *DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal) addWrite(p *DeadPropertyAssignmentElimination_Property, lhs rhino.Node) {
	jutil.CheckArgument(lhs.IsQualifiedName())
	if len(f.writeSlab) == cap(f.writeSlab) {
		nextCap := min(max(cap(f.writeSlab)*2, 16), 256)
		f.writeSlab = make([]DeadPropertyAssignmentElimination_PropertyWrite, 0, nextCap)
	}
	idx := len(f.writeSlab)
	f.writeSlab = f.writeSlab[:idx+1]
	w := &f.writeSlab[idx]
	w.assignedAt = lhs
	w.qualifiedName = lhs.GetQualifiedName()
	p.writes.Add(w)
}

func (f *DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal) traverse(compiler AbstractCompiler, root rhino.Node) {
	jutil.CheckNotNull(compiler)
	if root.IsNil() {
		b := NodeTraversal_Builder{compiler: compiler, callback: f}
		newNodeTraversal(&b).Traverse(root)
		return
	}
	var currentNode rhino.Node
	var currentScript rhino.Node
	jutil.Try(func() {
		if Platform_IsThreadInterrupted() {
			panic(jutil.NewRuntimeExceptionWithCause(jutil.NewInterruptedException()))
		}
		currentNode = root
		f.walkBranch(root, rhino.Node{}, &currentNode, &currentScript)
	}, func(unexpectedException jutil.Throwable) {
		b := NodeTraversal_Builder{compiler: compiler, callback: f}
		t := newNodeTraversal(&b)
		t.InitTraversal(root)
		if !currentScript.IsNil() {
			t.currentScript = currentScript
		}
		t.currentNode = currentNode
		t.ThrowUnexpectedException(unexpectedException)
	})
}

func (f *DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal) walkBranch(n, parent rhino.Node, currentNode, currentScript *rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_SCRIPT:
		if Platform_IsThreadInterrupted() {
			panic(jutil.NewRuntimeExceptionWithCause(jutil.NewInterruptedException()))
		}
		*currentNode = n
		*currentScript = n
		if f.VisitNode(n, parent) {
			for child := n.GetFirstChild(); !child.IsNil(); {
				next := child.GetNext()
				f.walkBranch(child, n, currentNode, currentScript)
				child = next
			}
			*currentNode = n
			f.Visit(nil, n, parent)
		}
		return
	case rhino.Token_FUNCTION:
		*currentNode = n
		if f.VisitNode(n, parent) {
			fnName := n.GetFirstChild()
			if !parent.IsNil() {
				_ = NodeUtil_IsFunctionDeclaration(n)
			}
			args := fnName.GetNext()
			body := args.GetNext()
			f.walkBranch(fnName, n, currentNode, currentScript)
			f.walkBranch(args, n, currentNode, currentScript)
			f.walkBranch(body, n, currentNode, currentScript)
			*currentNode = n
			f.Visit(nil, n, parent)
		}
		return
	case rhino.Token_CLASS:
		*currentNode = n
		if f.VisitNode(n, parent) {
			className := n.GetFirstChild()
			extendsClause := className.GetNext()
			body := extendsClause.GetNext()
			_ = NodeUtil_IsClassExpression(n)
			f.walkBranch(extendsClause, n, currentNode, currentScript)
			f.walkBranch(className, n, currentNode, currentScript)
			f.walkBranch(body, n, currentNode, currentScript)
			*currentNode = n
			f.Visit(nil, n, parent)
		}
		return
	case rhino.Token_CLASS_MEMBERS:
		*currentNode = n
		if !f.VisitNode(n, parent) {
			return
		}
		for child := n.GetFirstChild(); !child.IsNil(); {
			next := child.GetNext()
			if child.IsComputedProp() || child.IsComputedFieldDef() {
				f.walkBranch(child.GetFirstChild(), child, currentNode, currentScript)
			}
			child = next
		}
		for child := n.GetFirstChild(); !child.IsNil(); {
			next := child.GetNext()
			switch child.GetToken() {
			case rhino.Token_COMPUTED_PROP:
				*currentNode = n
				if f.VisitNode(child, n) {
					f.walkBranch(child.GetLastChild(), child, currentNode, currentScript)
					*currentNode = n
					f.Visit(nil, child, n)
				}
			case rhino.Token_COMPUTED_FIELD_DEF:
				*currentNode = n
				if f.VisitNode(child, n) {
					if child.HasTwoChildren() {
						f.walkBranch(child.GetLastChild(), child, currentNode, currentScript)
					}
					*currentNode = n
					f.Visit(nil, child, n)
				}
			case rhino.Token_MEMBER_FIELD_DEF, rhino.Token_BLOCK, rhino.Token_MEMBER_FUNCTION_DEF,
				rhino.Token_MEMBER_VARIABLE_DEF, rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF:
				f.walkBranch(child, n, currentNode, currentScript)
			default:
				panic(jutil.NewIllegalStateExceptionMsg("Invalid class member: " + child.GetToken().String()))
			}
			child = next
		}
		*currentNode = n
		f.Visit(nil, n, parent)
		return
	default:
	}

	*currentNode = n
	if !f.VisitNode(n, parent) {
		return
	}
	for child := n.GetFirstChild(); !child.IsNil(); {
		next := child.GetNext()
		f.walkBranch(child, n, currentNode, currentScript)
		child = next
	}
	*currentNode = n
	f.Visit(nil, n, parent)
}

// GetOrCreateProperty gets a Property given the node that references it;
// the Property is created if it does not already exist.
//
// Returns a Property, or null if the provided node is not a qualified name.
//
// Java: DeadPropertyAssignmentElimination.FindCandidateAssignmentTraversal.getOrCreateProperty(Node)
func (f *DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal) GetOrCreateProperty(propNode rhino.Node) *DeadPropertyAssignmentElimination_Property {
	var propName string
	switch propNode.GetToken() {
	case rhino.Token_GETPROP:
		if !propNode.GetFirstChild().IsQualifiedName() {
			return nil
		}
		propName = propNode.GetString()
	case rhino.Token_NAME:
		propName = propNode.GetString()
		if propName == "" {
			return nil
		}
	case rhino.Token_THIS:
		propName = "this"
	case rhino.Token_SUPER:
		propName = "super"
	default:
		return nil
	}

	property := f.propertyMap.Get(propName)
	if property == nil {
		property = f.allocProperty(propName)
		f.propertyMap.PutNew(propName, property)
	}

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
func (f *DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal) ShouldTraverse(nodeTraversal *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	return f.VisitNode(n, parent)
}

// Java: DeadPropertyAssignmentElimination.FindCandidateAssignmentTraversal.visit(NodeTraversal, Node, Node)
func (f *DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
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
func (f *DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal) VisitBlock(blockNode rhino.Node) {
	jutil.CheckArgument(blockNode.IsBlock())

	// We don't do flow analysis yet so we're going to assume everything written up to this
	// block is read.
	if blockNode.HasChildren() {
		f.MarkAllPropsRead()
	}
}

// Java: DeadPropertyAssignmentElimination.FindCandidateAssignmentTraversal.isConditionalExpression(Node)
func DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal_IsConditionalExpression(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_AND, rhino.Token_OR, rhino.Token_HOOK, rhino.Token_COALESCE, rhino.Token_OPTCHAIN_CALL,
		rhino.Token_OPTCHAIN_GETELEM, rhino.Token_OPTCHAIN_GETPROP:
		return true
	default:
		return false
	}
}

// Java: DeadPropertyAssignmentElimination.FindCandidateAssignmentTraversal.visitAssignmentLhs(Node)
func (f *DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal) VisitAssignmentLhs(lhs rhino.Node) {
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
	f.addWrite(property, lhs)

	// Now we need to go up the prop chain and mark those as read.
	child := lhs.GetFirstChild()
	for !child.IsNil() {
		var childName string
		if child.IsGetProp() || child.IsName() {
			childName = child.GetString()
		} else if child.IsThis() {
			childName = "this"
		} else if child.IsSuper() {
			childName = "super"
		} else {
			break
		}
		childProperty := f.propertyMap.Get(childName)
		if childProperty == nil {
			break
		}
		childProperty.MarkLastWriteRead()
		child = child.GetFirstChild()
	}
}

// Java: DeadPropertyAssignmentElimination.FindCandidateAssignmentTraversal.visitNode(Node, Node)
func (f *DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal) VisitNode(n rhino.Node, parent rhino.Node) bool {
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
	case rhino.Token_THROW, rhino.Token_FOR, rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF, rhino.Token_ITER_SPREAD, rhino.Token_ARRAY_PATTERN, rhino.Token_SWITCH:
		// Loops and switch statements may execute out of order, and implicit iteration operations
		// (for-of, for-await-of, spread, array destructuring) invoke arbitrary JS iterator
		// methods which can read properties.
		f.MarkAllPropsRead()
		return false
	case rhino.Token_OBJECT_PATTERN:
		for child := n.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
			if child.IsStringKey() {
				propName := child.GetString()
				property := f.propertyMap.Get(propName)
				if property != nil {
					property.MarkLastWriteRead()
					property.MarkChildrenRead()
				}
			} else {
				f.MarkAllPropsRead()
				break
			}
		}
		return true
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
