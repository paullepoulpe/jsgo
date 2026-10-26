// Port of com/google/javascript/jscomp/CreateSyntheticBlocks.java @ closure-compiler v20261021

/*
 * Copyright 2009 The Closure Compiler Authors.
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

// Go-only: this pass uses the typed AST view (internal/ast): a marker is the
// expression statement (ast.ExprStmt) of a call whose callee is a name
// (call.Callee().AsName()); AddBlocks moves the two marker statements into
// the outer block by name instead of taking the outer block's next sibling
// twice; the synthetic blocks are ast.Blocks; and a function declaration is
// rewritten through its typed name (fn.Name(), fn.RemoveName()).
// MoveSiblingExclusive stays on rhino.Nodes: it moves whatever siblings lie
// between two nodes.

// CreateSyntheticBlocks creates synthetic blocks to simplify preventing
// optimizations from moving code past "markers" in the source.
//
// This supports a serving infrastructure that selectively removes parts
// that the code, that are tagged using the start/end section markers. This
// removal happens at serving time not compilation time. The conventional
// markers are "_START_TEMPLATE_SECTION" and "_END_TEMPLATE_SECTION"
// respectively.
type CreateSyntheticBlocks struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	compiler AbstractCompiler

	// Name of the start marker.
	startMarkerName string

	// Name of the end marker.
	endMarkerName string

	// Markers can be nested.
	//
	// Go-only: language-adapter — the ArrayDeque used as a stack
	// (push/pop/iteration at the front): the top of the stack is the last
	// element of the slice.
	markerStack []ast.ExprStmt

	validMarkers []*createSyntheticBlocks_Marker
}

// Java: CreateSyntheticBlocks.UNMATCHED_START_MARKER
var CreateSyntheticBlocks_UNMATCHED_START_MARKER = DiagnosticType_Error("JSC_UNMATCHED_START_MARKER", "Unmatched {0}")

// Java: CreateSyntheticBlocks.UNMATCHED_END_MARKER
var CreateSyntheticBlocks_UNMATCHED_END_MARKER = DiagnosticType_Error("JSC_UNMATCHED_END_MARKER", "Unmatched {1} - {0} not in the same block")

// Java: CreateSyntheticBlocks.INVALID_MARKER_USAGE
var CreateSyntheticBlocks_INVALID_MARKER_USAGE = DiagnosticType_Error(
	"JSC_INVALID_MARKER_USAGE", "Marker {0} can only be used in a simple call expression")

type createSyntheticBlocks_Marker struct {
	startMarker ast.ExprStmt
	endMarker   ast.ExprStmt
}

// Java: CreateSyntheticBlocks.Marker(Node, Node)
func newCreateSyntheticBlocks_Marker(startMarker ast.ExprStmt, endMarker ast.ExprStmt) *createSyntheticBlocks_Marker {
	return &createSyntheticBlocks_Marker{startMarker: startMarker, endMarker: endMarker}
}

// Java: CreateSyntheticBlocks(AbstractCompiler, String, String)
func NewCreateSyntheticBlocks(compiler AbstractCompiler, startMarkerName string, endMarkerName string) *CreateSyntheticBlocks {
	return &CreateSyntheticBlocks{
		compiler:        compiler,
		startMarkerName: startMarkerName,
		endMarkerName:   endMarkerName,
	}
}

// Java: CreateSyntheticBlocks.process(Node, Node)
func (p *CreateSyntheticBlocks) Process(externs rhino.Node, root rhino.Node) {
	// Find and validate the markers.
	NodeTraversal_Traverse(p.compiler, root, p)

	// Complain about any unmatched markers.
	for i := len(p.markerStack) - 1; i >= 0; i-- {
		node := p.markerStack[i]
		p.compiler.Report(JSError_MakeWithNode(node.Node(), CreateSyntheticBlocks_UNMATCHED_START_MARKER, p.startMarkerName))
	}

	// Add the block for the valid marker sets.
	for _, marker := range p.validMarkers {
		p.AddBlocks(marker)
	}
}

// RewriteFunctionDeclaration rewrites the function declaration from:
// function x() {} FUNCTION NAME x PARAM_LIST BLOCK to: var x = function()
// {}; VAR NAME x FUNCTION NAME (w/ empty string) PARAM_LIST BLOCK
//
// Java: CreateSyntheticBlocks.rewriteFunctionDeclaration(Node)
func (p *CreateSyntheticBlocks) RewriteFunctionDeclaration(n ast.Function) {
	// Prepare a variable for the function.
	oldNameNode, ok := n.Name()
	jutil.CheckState(ok) // a function declaration has a name
	fnNameNode := oldNameNode.Node().CloneNode()
	var_ := rhino.IR_Var(fnNameNode).Srcref(n.Node())

	// Prepare the function.
	n.RemoveName()
	p.compiler.ReportChangeToEnclosingScope(oldNameNode.Node())

	// Move the function to the front of the parent.
	parent := n.Node().GetParent()
	n.Detach()
	parent.AddChildToFront(var_)
	p.compiler.ReportChangeToEnclosingScope(var_)
	fnNameNode.AddChildToFront(n.Node())
}

// Java: CreateSyntheticBlocks.rewriteFunctionDeclarationsInBlock(Node)
func (p *CreateSyntheticBlocks) RewriteFunctionDeclarationsInBlock(block ast.Block) {
	// ReadAhead: the current statement may move.
	for current := range ast.ReadAhead(block.FirstStatement()) {
		if fn, ok := current.AsFunction(); ok && NodeUtil_IsFunctionDeclaration(fn.Node()) {
			p.RewriteFunctionDeclaration(fn)
		}
	}
}

// AddBlocks adds synthetic blocks for the given marker.
//
// Java: CreateSyntheticBlocks.addBlocks(Marker)
func (p *CreateSyntheticBlocks) AddBlocks(marker *createSyntheticBlocks_Marker) {
	// Add a block around the template section so that
	//   START
	//   BODY
	//   END
	// is transformed to:
	//   BLOCK (synthetic)
	//     START
	//       BLOCK (synthetic)
	//         BODY
	//     END
	// This prevents the start or end markers from mingling with the code in
	// the block body.

	start, end := marker.startMarker, marker.endMarker

	outerBlock := createSyntheticBlocks_newSyntheticBlock(start)
	outerBlock.Node().InsertBefore(start.Node())

	innerBlock := createSyntheticBlocks_newSyntheticBlock(start)
	// Move everything after the start node up to the end node into the inner
	// block.
	p.MoveSiblingExclusive(innerBlock.Node(), start.Node(), end.Node())

	// Add the start node.
	outerBlock.AppendStatement(start.Detach().AsStatement())
	// Add the inner block.
	outerBlock.AppendStatement(innerBlock.AsStatement())
	// And finally the end node.
	outerBlock.AppendStatement(end.Detach().AsStatement())

	// NOTE: Moving the code into a block made sense prior to ES2015 when
	// declarations were never block scoped. But with ES2015+ function,
	// class, let and const are all block scoped. Here we are only rewriting
	// functions because this is the behavior that "normalize" previously had
	// and we aren't currently (July 2020) trying to improve this code's
	// behavior. This maintains the status quo and allows the normalization
	// of function to be removed.
	p.RewriteFunctionDeclarationsInBlock(innerBlock)

	outerParent := outerBlock.Node().GetParent()
	jutil.CheckStateMsg1(outerParent.IsBlock() || outerParent.IsScript(), "%s", outerParent)
	jutil.CheckState(innerBlock.Node().GetParent().IsBlock())
	p.compiler.ReportChangeToEnclosingScope(outerBlock.Node())
}

// createSyntheticBlocks_newSyntheticBlock returns an empty synthetic block at
// the position of marker.
//
// Go-only: Java builds the two blocks of addBlocks inline.
func createSyntheticBlocks_newSyntheticBlock(marker ast.ExprStmt) ast.Block {
	block := ast.NewBlock()
	block.Node().Srcref(marker.Node())
	block.Node().SetIsSyntheticBlock(true)
	return block
}

// MoveSiblingExclusive moves the nodes between start and end from the
// source block to the destination block. If start is null, move the first
// child of the block. If end is null, move the last child of the block.
//
// Java: CreateSyntheticBlocks.moveSiblingExclusive(Node, Node, Node)
func (p *CreateSyntheticBlocks) MoveSiblingExclusive(dest rhino.Node, start rhino.Node, end rhino.Node) {
	rhino.CheckNotNull(start)
	rhino.CheckNotNull(end)
	for start.GetNext() != end {
		child := start.GetNext().Detach()
		dest.AddChildToBack(child)
	}
}

// Java: CreateSyntheticBlocks.visit(NodeTraversal, Node, Node)
func (p *CreateSyntheticBlocks) Visit(unused *NodeTraversal, n rhino.Node, parent rhino.Node) {
	call, ok := ast.AsCall(n)
	if !ok {
		return
	}
	callTarget, ok := call.Callee().AsName()
	if !ok {
		return
	}

	callName := callTarget.Ident()

	if p.startMarkerName == callName {
		startMarkerNode, ok := ast.AsExprStmt(parent)
		if !ok {
			p.compiler.Report(JSError_MakeWithNode(n, CreateSyntheticBlocks_INVALID_MARKER_USAGE, p.startMarkerName))
			return
		}
		p.markerStack = append(p.markerStack, startMarkerNode)
		return
	}

	if p.endMarkerName != callName {
		return
	}

	endMarkerNode, ok := ast.AsExprStmt(parent)
	if !ok {
		p.compiler.Report(JSError_MakeWithNode(n, CreateSyntheticBlocks_INVALID_MARKER_USAGE, p.endMarkerName))
		return
	}

	if len(p.markerStack) == 0 {
		p.compiler.Report(JSError_MakeWithNode(n, CreateSyntheticBlocks_UNMATCHED_END_MARKER, p.startMarkerName, p.endMarkerName))
		return
	}

	startMarkerNode := p.markerStack[len(p.markerStack)-1]
	p.markerStack = p.markerStack[:len(p.markerStack)-1]
	if endMarkerNode.Node().GetParent() != startMarkerNode.Node().GetParent() {
		// The end marker isn't in the same block as the start marker.
		p.compiler.Report(JSError_MakeWithNode(n, CreateSyntheticBlocks_UNMATCHED_END_MARKER, p.startMarkerName, p.endMarkerName))
		return
	}

	// This is a valid marker, add it to the list of markers to process.
	p.validMarkers = append(p.validMarkers, newCreateSyntheticBlocks_Marker(startMarkerNode, endMarkerNode))
}
