// Port of com/google/javascript/jscomp/RenameLabels.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// RenameLabels renames all the labels so that they have short names, to
// reduce code size and also to obfuscate the code.
//
// Label names have a unique namespace, so variable or function names clashes
// are not a concern, but keywords clashes are.
//
// Additionally, labels names are only within the statements include in the
// label and do not cross function boundaries. This means that it is possible
// to create one label name that is used for labels at any given depth of
// label nesting. Typically, the name "a" will be used for all top-level
// labels, "b" for the next nested label, and so on. For example:
//
//	function bar() {
//	  a: {
//	    b: {
//	      foo();
//	    }
//	  }
//
//	  a: {
//	    b: break a;
//	  }
//	}
//
// The general processes is as follows: process() is the entry point for the
// CompilerPass, and from there a standard "ScopedCallback" traversal is done,
// where "shouldTraverse" is called when descending the tree, and the "visit"
// is called in a depth first manner. The name for the label is selected
// during the decent in "shouldTraverse", and the references to the label
// name are renamed as they are encountered during the "visit". This means
// that if the label is unreferenced, it is known when the label node is
// visited, and, if so, can be safely removed.
type RenameLabels struct {
	compiler     AbstractCompiler
	nameSupplier func() string
	removeUnused bool
	markChanges  bool
}

// Java: RenameLabels(AbstractCompiler)
func NewRenameLabels(compiler AbstractCompiler) *RenameLabels {
	return NewRenameLabelsWithSupplier(compiler, NewRenameLabels_DefaultNameSupplier().Get, true, true)
}

// Java: RenameLabels(AbstractCompiler, Supplier<String>, boolean, boolean)
func NewRenameLabelsWithSupplier(
	compiler AbstractCompiler,
	supplier func() string,
	removeUnused bool,
	markChanges bool) *RenameLabels {
	return &RenameLabels{
		compiler:     compiler,
		nameSupplier: supplier,
		removeUnused: removeUnused,
		markChanges:  markChanges,
	}
}

// RenameLabels_DefaultNameSupplier is the nested class DefaultNameSupplier.
type RenameLabels_DefaultNameSupplier struct {
	// DefaultNameGenerator is used to create safe label names.
	nameGenerator NameGenerator
}

// Java: RenameLabels.DefaultNameSupplier()
func NewRenameLabels_DefaultNameSupplier() *RenameLabels_DefaultNameSupplier {
	return &RenameLabels_DefaultNameSupplier{
		nameGenerator: NewDefaultNameGeneratorWithReserved(
			jutil.NewLinkedHashSet[string](), "", jutil.NewLinkedHashSet[uint16]()),
	}
}

// Java: RenameLabels.DefaultNameSupplier.get()
func (s *RenameLabels_DefaultNameSupplier) Get() string {
	return s.nameGenerator.GenerateNextName()
}

// RenameLabels_ProcessLabels iterates through the nodes, renaming all the
// labels.
type RenameLabels_ProcessLabels struct {
	outer *RenameLabels

	markChanges bool

	// A stack of labels namespaces. Labels in an outer scope aren't part of
	// an inner scope, so a new namespace is created each time a scope is
	// entered.
	//
	// Go-only: language-adapter — the ArrayDeque used as a stack
	// (push/peek/pop at the head) is a slice whose top is its last element.
	namespaceStack    []RenameLabels_LabelNamespace
	namespaceStackBuf [8]RenameLabels_LabelNamespace

	// The list of generated names. Typically, the first name will be "a",
	// the second "b", etc.
	names    []string
	namesBuf [4]string
}

// Java: RenameLabels.ProcessLabels(boolean)
func NewRenameLabels_ProcessLabels(outer *RenameLabels, markChanges bool) *RenameLabels_ProcessLabels {
	p := &RenameLabels_ProcessLabels{outer: outer, markChanges: markChanges}
	p.namespaceStack = p.namespaceStackBuf[:1]
	p.names = p.namesBuf[:0]
	return p
}

// peek is Deque.peek() on namespaceStack.
//
// Go-only: language-adapter — see namespaceStack.
func (p *RenameLabels_ProcessLabels) peek() *RenameLabels_LabelNamespace {
	return &p.namespaceStack[len(p.namespaceStack)-1]
}

// Java: RenameLabels.ProcessLabels.enterScope(NodeTraversal)
func (p *RenameLabels_ProcessLabels) EnterScope(nodeTraversal *NodeTraversal) {
	// Start a new namespace for label names.
	if nodeTraversal.GetScopeRoot().IsFunction() {
		p.namespaceStack = append(p.namespaceStack, RenameLabels_LabelNamespace{})
	}
}

// Java: RenameLabels.ProcessLabels.exitScope(NodeTraversal)
func (p *RenameLabels_ProcessLabels) ExitScope(nodeTraversal *NodeTraversal) {
	if nodeTraversal.GetScopeRoot().IsFunction() {
		top := len(p.namespaceStack) - 1
		p.namespaceStack[top].renameMap = nil
		p.namespaceStack = p.namespaceStack[:top]
	}
}

// ShouldTraverse is call when descending into the Node tree, so it is used
// here to build the context for label renames.
//
// Java: RenameLabels.ProcessLabels.shouldTraverse(NodeTraversal, Node, Node)
func (p *RenameLabels_ProcessLabels) ShouldTraverse(nodeTraversal *NodeTraversal, node *rhino.Node, parent *rhino.Node) bool {
	if node.IsLabel() {
		// Determine the new name for this label.
		current := p.peek()
		if current.renameMap == nil {
			current.renameMap = jutil.NewLinkedHashMap[string, *RenameLabels_LabelInfo]()
		}
		currentDepth := current.renameMap.Size() + 1
		name := node.GetFirstChild().GetString()

		// Store the context for this label name.
		li := NewRenameLabels_LabelInfo(currentDepth)
		jutil.CheckState(!current.renameMap.ContainsKey(name))
		current.renameMap.Put(name, li)

		// Create a new name, if needed, for this depth.
		if len(p.names) < currentDepth {
			p.names = append(p.names, p.outer.nameSupplier())
		}
	}

	return true
}

// Visit delegates the actual processing of the node to visitLabel and
// visitBreakOrContinue.
//
// Java: RenameLabels.ProcessLabels.visit(NodeTraversal, Node, Node)
func (p *RenameLabels_ProcessLabels) Visit(t *NodeTraversal, node *rhino.Node, parent *rhino.Node) {
	switch node.GetToken() {
	case rhino.Token_LABEL:
		p.VisitLabel(t, node)
	case rhino.Token_BREAK, rhino.Token_CONTINUE:
		p.VisitBreakOrContinue(t, node)
	default:
	}
}

// VisitBreakOrContinue renames label references in breaks and continues.
//
// node is the break or continue node.
//
// Java: RenameLabels.ProcessLabels.visitBreakOrContinue(NodeTraversal, Node)
func (p *RenameLabels_ProcessLabels) VisitBreakOrContinue(t *NodeTraversal, node *rhino.Node) {
	nameNode := node.GetFirstChild()
	if nameNode != nil {
		// This is a named break or continue;
		name := nameNode.GetString()
		jutil.CheckState(name != "")
		li := p.GetLabelInfo(name)
		if li != nil {
			newName := p.GetNameForId(li.id)
			// Mark the label as referenced so it isn't removed.
			li.referenced = true
			if name != newName {
				// Give it the short name.
				nameNode.SetString(newName)
				if p.markChanges {
					t.ReportCodeChange()
				}
			}
		}
	}
}

// VisitLabel renames or removes labels.
//
// node is the label node.
//
// Java: RenameLabels.ProcessLabels.visitLabel(NodeTraversal, Node)
func (p *RenameLabels_ProcessLabels) VisitLabel(t *NodeTraversal, node *rhino.Node) {
	nameNode := node.GetFirstChild()
	jutil.CheckState(nameNode != nil)
	name := nameNode.GetString()
	li := p.GetLabelInfo(name)
	// This is a label...
	if li.referenced || !p.outer.removeUnused {
		newName := p.GetNameForId(li.id)
		if name != newName {
			// ... and it is used, give it the short name.
			nameNode.SetString(newName)
			if p.markChanges {
				t.ReportCodeChange()
			}
		}
	} else {
		// ... and it is not referenced, just remove it.
		newChild := node.GetLastChild()
		newChild.Detach()
		node.ReplaceWith(newChild)
		if newChild.IsBlock() {
			NodeUtil_TryMergeBlock(newChild, false)
		}
		if p.markChanges {
			t.ReportCodeChange()
		}
	}

	// Remove the label from the current stack of labels.
	p.peek().renameMap.Remove(name)
}

// GetNameForId returns the short name of the identified label.
//
// id is the id, which is the depth of the label in the current context, for
// which to get a short name.
//
// Java: RenameLabels.ProcessLabels.getNameForId(int)
func (p *RenameLabels_ProcessLabels) GetNameForId(id int) string {
	return p.names[id-1]
}

// GetLabelInfo returns the structure representing the name in the current
// context.
//
// name is the name to retrieve information about.
//
// Java: RenameLabels.ProcessLabels.getLabelInfo(String)
func (p *RenameLabels_ProcessLabels) GetLabelInfo(name string) *RenameLabels_LabelInfo {
	m := p.peek().renameMap
	if m == nil {
		return nil
	}
	return m.Get(name)
}

// Java: RenameLabels.process(Node, Node)
func (r *RenameLabels) Process(externs *rhino.Node, root *rhino.Node) {
	// Do variable reference counting.
	NodeTraversal_Traverse(r.compiler, root, NewRenameLabels_ProcessLabels(r, r.markChanges))
}

// RenameLabels_LabelInfo is the nested class LabelInfo.
type RenameLabels_LabelInfo struct {
	referenced bool
	id         int
}

// Java: RenameLabels.LabelInfo(int)
func NewRenameLabels_LabelInfo(id int) *RenameLabels_LabelInfo {
	return &RenameLabels_LabelInfo{id: id}
}

// RenameLabels_LabelNamespace is the nested class LabelNamespace.
type RenameLabels_LabelNamespace struct {
	renameMap *jutil.LinkedHashMap[string, *RenameLabels_LabelInfo]
}

// Java: RenameLabels.LabelNamespace()
func NewRenameLabels_LabelNamespace() *RenameLabels_LabelNamespace {
	return &RenameLabels_LabelNamespace{renameMap: jutil.NewLinkedHashMap[string, *RenameLabels_LabelInfo]()}
}
