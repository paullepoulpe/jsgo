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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

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
	markerStack []*rhino.Node

	validMarkers *jutil.ArrayList[*createSyntheticBlocks_Marker]
}

// Java: CreateSyntheticBlocks.UNMATCHED_START_MARKER
var CreateSyntheticBlocks_UNMATCHED_START_MARKER = DiagnosticType_Error("JSC_UNMATCHED_START_MARKER", "Unmatched {0}")

// Java: CreateSyntheticBlocks.UNMATCHED_END_MARKER
var CreateSyntheticBlocks_UNMATCHED_END_MARKER = DiagnosticType_Error("JSC_UNMATCHED_END_MARKER", "Unmatched {1} - {0} not in the same block")

// Java: CreateSyntheticBlocks.INVALID_MARKER_USAGE
var CreateSyntheticBlocks_INVALID_MARKER_USAGE = DiagnosticType_Error(
	"JSC_INVALID_MARKER_USAGE", "Marker {0} can only be used in a simple call expression")

type createSyntheticBlocks_Marker struct {
	startMarker *rhino.Node
	endMarker   *rhino.Node
}

// Java: CreateSyntheticBlocks.Marker(Node, Node)
func newCreateSyntheticBlocks_Marker(startMarker *rhino.Node, endMarker *rhino.Node) *createSyntheticBlocks_Marker {
	return &createSyntheticBlocks_Marker{startMarker: startMarker, endMarker: endMarker}
}

// Java: CreateSyntheticBlocks(AbstractCompiler, String, String)
func NewCreateSyntheticBlocks(compiler AbstractCompiler, startMarkerName string, endMarkerName string) *CreateSyntheticBlocks {
	return &CreateSyntheticBlocks{
		compiler:        compiler,
		startMarkerName: startMarkerName,
		endMarkerName:   endMarkerName,
		validMarkers:    jutil.NewArrayList[*createSyntheticBlocks_Marker](),
	}
}

// Java: CreateSyntheticBlocks.process(Node, Node)
func (p *CreateSyntheticBlocks) Process(externs *rhino.Node, root *rhino.Node) {
	// Find and validate the markers.
	NodeTraversal_Traverse(p.compiler, root, p)

	// Complain about any unmatched markers.
	for i := len(p.markerStack) - 1; i >= 0; i-- {
		node := p.markerStack[i]
		p.compiler.Report(JSError_MakeWithNode(node, CreateSyntheticBlocks_UNMATCHED_START_MARKER, p.startMarkerName))
	}

	// Add the block for the valid marker sets.
	for marker := range p.validMarkers.All() {
		p.AddBlocks(marker)
	}
}

// RewriteFunctionDeclaration rewrites the function declaration from:
// function x() {} FUNCTION NAME x PARAM_LIST BLOCK to: var x = function()
// {}; VAR NAME x FUNCTION NAME (w/ empty string) PARAM_LIST BLOCK
//
// Java: CreateSyntheticBlocks.rewriteFunctionDeclaration(Node)
func (p *CreateSyntheticBlocks) RewriteFunctionDeclaration(n *rhino.Node) {
	// Prepare a variable for the function.
	oldNameNode := n.GetFirstChild()
	fnNameNode := oldNameNode.CloneNode()
	var_ := rhino.IR_Var(fnNameNode).Srcref(n)

	// Prepare the function.
	oldNameNode.SetString("")
	p.compiler.ReportChangeToEnclosingScope(oldNameNode)

	// Move the function to the front of the parent.
	parent := n.GetParent()
	n.Detach()
	parent.AddChildToFront(var_)
	p.compiler.ReportChangeToEnclosingScope(var_)
	fnNameNode.AddChildToFront(n)
}

// Java: CreateSyntheticBlocks.rewriteFunctionDeclarationsInBlock(Node)
func (p *CreateSyntheticBlocks) RewriteFunctionDeclarationsInBlock(block *rhino.Node) {
	jutil.CheckState(block.IsBlock())

	var next *rhino.Node
	for current := block.GetFirstChild(); current != nil; current = next {
		// Get the next node now, because the current node will move and that
		// will change its next sibling.
		next = current.GetNext()
		if NodeUtil_IsFunctionDeclaration(current) {
			p.RewriteFunctionDeclaration(current)
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

	outerBlock := rhino.IR_Block().Srcref(marker.startMarker)
	outerBlock.SetIsSyntheticBlock(true)
	outerBlock.InsertBefore(marker.startMarker)

	innerBlock := rhino.IR_Block().Srcref(marker.startMarker)
	innerBlock.SetIsSyntheticBlock(true)
	// Move everything after the start node up to the end node into the inner
	// block.
	p.MoveSiblingExclusive(innerBlock, marker.startMarker, marker.endMarker)

	// Add the start node.
	outerBlock.AddChildToBack(outerBlock.GetNext().Detach())
	// Add the inner block.
	outerBlock.AddChildToBack(innerBlock)
	// And finally the end node.
	outerBlock.AddChildToBack(outerBlock.GetNext().Detach())

	// NOTE: Moving the code into a block made sense prior to ES2015 when
	// declarations were never block scoped. But with ES2015+ function,
	// class, let and const are all block scoped. Here we are only rewriting
	// functions because this is the behavior that "normalize" previously had
	// and we aren't currently (July 2020) trying to improve this code's
	// behavior. This maintains the status quo and allows the normalization
	// of function to be removed.
	p.RewriteFunctionDeclarationsInBlock(innerBlock)

	jutil.CheckStateMsg1(
		outerBlock.GetParent().IsBlock() || outerBlock.GetParent().IsScript(),
		"%s", outerBlock.GetParent())
	jutil.CheckState(innerBlock.GetParent().IsBlock())
	p.compiler.ReportChangeToEnclosingScope(outerBlock)
}

// MoveSiblingExclusive moves the nodes between start and end from the
// source block to the destination block. If start is null, move the first
// child of the block. If end is null, move the last child of the block.
//
// Java: CreateSyntheticBlocks.moveSiblingExclusive(Node, Node, Node)
func (p *CreateSyntheticBlocks) MoveSiblingExclusive(dest *rhino.Node, start *rhino.Node, end *rhino.Node) {
	jutil.CheckNotNull(start)
	jutil.CheckNotNull(end)
	for start.GetNext() != end {
		child := start.GetNext().Detach()
		dest.AddChildToBack(child)
	}
}

// Java: CreateSyntheticBlocks.visit(NodeTraversal, Node, Node)
func (p *CreateSyntheticBlocks) Visit(unused *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if !n.IsCall() || !n.GetFirstChild().IsName() {
		return
	}

	callTarget := n.GetFirstChild()
	callName := callTarget.GetString()

	if p.startMarkerName == callName {
		if !parent.IsExprResult() {
			p.compiler.Report(JSError_MakeWithNode(n, CreateSyntheticBlocks_INVALID_MARKER_USAGE, p.startMarkerName))
			return
		}
		p.markerStack = append(p.markerStack, parent)
		return
	}

	if p.endMarkerName != callName {
		return
	}

	endMarkerNode := parent
	if !endMarkerNode.IsExprResult() {
		p.compiler.Report(JSError_MakeWithNode(n, CreateSyntheticBlocks_INVALID_MARKER_USAGE, p.endMarkerName))
		return
	}

	if len(p.markerStack) == 0 {
		p.compiler.Report(JSError_MakeWithNode(n, CreateSyntheticBlocks_UNMATCHED_END_MARKER, p.startMarkerName, p.endMarkerName))
		return
	}

	startMarkerNode := p.markerStack[len(p.markerStack)-1]
	p.markerStack = p.markerStack[:len(p.markerStack)-1]
	if endMarkerNode.GetParent() != startMarkerNode.GetParent() {
		// The end marker isn't in the same block as the start marker.
		p.compiler.Report(JSError_MakeWithNode(n, CreateSyntheticBlocks_UNMATCHED_END_MARKER, p.startMarkerName, p.endMarkerName))
		return
	}

	// This is a valid marker, add it to the list of markers to process.
	p.validMarkers.Add(newCreateSyntheticBlocks_Marker(startMarkerNode, endMarkerNode))
}
