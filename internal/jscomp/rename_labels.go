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
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

var renameLabelsVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_LABEL)
	s.Add(rhino.Token_BREAK)
	s.Add(rhino.Token_CONTINUE)
	return s
}()

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
	defaultGen   *DefaultNameGenerator
	removeUnused bool
	markChanges  bool
}

// Java: RenameLabels(AbstractCompiler)
func NewRenameLabels(compiler AbstractCompiler) *RenameLabels {
	return &RenameLabels{
		compiler:     compiler,
		defaultGen:   NewDefaultNameGenerator(),
		removeUnused: true,
		markChanges:  true,
	}
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

type labelStackEntry struct {
	origID     strpool.ID
	depth      int32
	referenced bool
}

// RenameLabels_ProcessLabels iterates through the nodes, renaming all the
// labels.
type RenameLabels_ProcessLabels struct {
	outer *RenameLabels

	markChanges bool

	// Flat lexical label stack across all nested scopes, partitioned by
	// function-scope start offsets in scopeStarts.
	entries        []labelStackEntry
	entriesBuf     [16]labelStackEntry
	scopeStarts    []int
	scopeStartsBuf [8]int

	// Pre-interned generated short label IDs indexed by (depth - 1).
	names    []strpool.ID
	namesBuf [8]strpool.ID
}

// Java: RenameLabels.ProcessLabels(boolean)
func NewRenameLabels_ProcessLabels(outer *RenameLabels, markChanges bool) *RenameLabels_ProcessLabels {
	p := &RenameLabels_ProcessLabels{outer: outer, markChanges: markChanges}
	p.entries = p.entriesBuf[:0]
	p.scopeStarts = p.scopeStartsBuf[:1]
	p.scopeStarts[0] = 0
	p.names = p.namesBuf[:0]
	return p
}

func (*RenameLabels_ProcessLabels) visitTokens() *ast.TokenSet {
	return &renameLabelsVisitTokens
}

func (p *RenameLabels_ProcessLabels) nextLabelNameID() strpool.ID {
	if p.outer.nameSupplier != nil {
		return strpool.Intern(p.outer.nameSupplier())
	}
	return strpool.Intern(p.outer.defaultGen.GenerateNextName())
}

// Java: RenameLabels.ProcessLabels.enterScope(NodeTraversal)
func (p *RenameLabels_ProcessLabels) EnterScope(nodeTraversal *NodeTraversal) {
	// Start a new namespace for label names.
	if _, ok := ast.AsFunction(nodeTraversal.GetScopeRoot()); ok {
		p.scopeStarts = append(p.scopeStarts, len(p.entries))
	}
}

// Java: RenameLabels.ProcessLabels.exitScope(NodeTraversal)
func (p *RenameLabels_ProcessLabels) ExitScope(nodeTraversal *NodeTraversal) {
	if _, ok := ast.AsFunction(nodeTraversal.GetScopeRoot()); ok {
		top := len(p.scopeStarts) - 1
		p.entries = p.entries[:p.scopeStarts[top]]
		p.scopeStarts = p.scopeStarts[:top]
	}
}

// ShouldTraverse is call when descending into the Node tree, so it is used
// here to build the context for label renames.
//
// Java: RenameLabels.ProcessLabels.shouldTraverse(NodeTraversal, Node, Node)
func (p *RenameLabels_ProcessLabels) ShouldTraverse(nodeTraversal *NodeTraversal, node rhino.Node, parent rhino.Node) bool {
	if lbl, ok := ast.AsLabel(node); ok {
		scopeStart := p.scopeStarts[len(p.scopeStarts)-1]
		origID := strpool.ID(lbl.Name().Node().GetStringID())
		for i := scopeStart; i < len(p.entries); i++ {
			jutil.CheckState(p.entries[i].origID != origID)
		}
		currentDepth := len(p.entries) - scopeStart + 1
		p.entries = append(p.entries, labelStackEntry{
			origID: origID,
			depth:  int32(currentDepth),
		})

		// Create a new name, if needed, for this depth.
		if len(p.names) < currentDepth {
			p.names = append(p.names, p.nextLabelNameID())
		}
	}

	return true
}

// Visit delegates the actual processing of the node to visitLabel and
// visitBreakOrContinue.
//
// Java: RenameLabels.ProcessLabels.visit(NodeTraversal, Node, Node)
func (p *RenameLabels_ProcessLabels) Visit(t *NodeTraversal, node rhino.Node, parent rhino.Node) {
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
func (p *RenameLabels_ProcessLabels) VisitBreakOrContinue(t *NodeTraversal, node rhino.Node) {
	var labelName ast.LabelName
	var hasLabel bool
	if brk, ok := ast.AsBreak(node); ok {
		labelName, hasLabel = brk.Label()
	} else {
		labelName, hasLabel = ast.MustContinue(node).Label()
	}
	if !hasLabel {
		return
	}
	// This is a named break or continue.
	nameNode := labelName.Node()
	origID := strpool.ID(nameNode.GetStringID())
	jutil.CheckState(origID > strpool.ID(rhino.RhinoStringPool_EmptyID))
	scopeStart := p.scopeStarts[len(p.scopeStarts)-1]
	for i := len(p.entries) - 1; i >= scopeStart; i-- {
		if p.entries[i].origID == origID {
			p.entries[i].referenced = true
			newID := p.names[p.entries[i].depth-1]
			if origID != newID {
				setNodeStringID(nameNode, newID)
				if p.markChanges {
					t.ReportCodeChange()
				}
			}
			break
		}
	}
}

// VisitLabel renames or removes labels.
//
// node is the label node.
//
// Java: RenameLabels.ProcessLabels.visitLabel(NodeTraversal, Node)
func (p *RenameLabels_ProcessLabels) VisitLabel(t *NodeTraversal, node rhino.Node) {
	lbl := ast.MustLabel(node)
	nameNode := lbl.Name().Node()
	origID := strpool.ID(nameNode.GetStringID())
	topIdx := len(p.entries) - 1
	scopeStart := p.scopeStarts[len(p.scopeStarts)-1]
	jutil.CheckState(topIdx >= scopeStart && p.entries[topIdx].origID == origID)
	entry := p.entries[topIdx]
	p.entries = p.entries[:topIdx]

	// This is a label...
	if entry.referenced || !p.outer.removeUnused {
		newID := p.names[entry.depth-1]
		if origID != newID {
			// ... and it is used, give it the short name.
			setNodeStringID(nameNode, newID)
			if p.markChanges {
				t.ReportCodeChange()
			}
		}
	} else {
		// ... and it is not referenced, just remove it.
		newChild := lbl.Body().Node()
		newChild.Detach()
		node.ReplaceWith(newChild)
		if newChild.IsBlock() {
			NodeUtil_TryMergeBlock(newChild, false)
		}
		if p.markChanges {
			t.ReportCodeChange()
		}
	}
}

// Java: RenameLabels.process(Node, Node)
func (r *RenameLabels) Process(externs rhino.Node, root rhino.Node) {
	p := RenameLabels_ProcessLabels{outer: r, markChanges: r.markChanges}
	p.entries = p.entriesBuf[:0]
	p.scopeStarts = p.scopeStartsBuf[:1]
	p.scopeStarts[0] = 0
	p.names = p.namesBuf[:0]
	NodeTraversal_Traverse(r.compiler, root, &p)
}
