// Port of com/google/javascript/jscomp/instrumentation/BranchCoverageInstrumentationCallback.java @ closure-compiler v20260415

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// BranchCoverageInstrumentationCallback instruments branch coverage for
// javascript.
//
// Java: public class BranchCoverageInstrumentationCallback extends NodeTraversal.AbstractCfgCallback
type BranchCoverageInstrumentationCallback struct {
	NodeTraversal_AbstractCfgCallbackBase

	compiler            AbstractCompiler
	instrumentationData *jutil.LinkedHashMap[string, *FileInstrumentationData]
}

// Java: BranchCoverageInstrumentationCallback.BRANCH_ARRAY_NAME_PREFIX
const branchCoverageInstrumentationCallback_BRANCH_ARRAY_NAME_PREFIX = "JSCompiler_lcov_branch_data_"

// BranchCoverageInstrumentationCallback_CreateArrayName returns a string that
// can be used for the branch coverage data.
//
// Java: BranchCoverageInstrumentationCallback.createArrayName(NodeTraversal)
func BranchCoverageInstrumentationCallback_CreateArrayName(traversal *NodeTraversal) string {
	return branchCoverageInstrumentationCallback_BRANCH_ARRAY_NAME_PREFIX +
		CoverageUtil_CreateIdentifierFromText(CoverageInstrumentationCallback_GetFileName(traversal))
}

// Java: BranchCoverageInstrumentationCallback(AbstractCompiler, Map<String, FileInstrumentationData>)
func NewBranchCoverageInstrumentationCallback(
	compiler AbstractCompiler, instrumentationData *jutil.LinkedHashMap[string, *FileInstrumentationData]) *BranchCoverageInstrumentationCallback {
	c := &BranchCoverageInstrumentationCallback{}
	c.InitNodeTraversal_AbstractCfgCallbackBase(c)
	c.compiler = compiler
	c.instrumentationData = instrumentationData
	return c
}

// branchCoverageInstrumentationCallback_stringsEqual is
// Objects.equals(String, String) on nullable strings.
func branchCoverageInstrumentationCallback_stringsEqual(a, b *string) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

// Java: BranchCoverageInstrumentationCallback.visit(NodeTraversal, Node, Node)
func (c *BranchCoverageInstrumentationCallback) Visit(traversal *NodeTraversal, node *rhino.Node, parent *rhino.Node) {
	fileNamePtr := traversal.GetSourceName()

	// If origin of node is not from sourceFile, do not instrument. This typically occurs when
	// polyfill code is injected into the sourceFile AST and this check avoids instrumenting it. We
	// avoid instrumentation as this callback does not distinguish between sourceFile code and
	// injected code and can result in an error.
	if !branchCoverageInstrumentationCallback_stringsEqual(fileNamePtr, node.GetSourceFileName()) {
		return
	}
	// Go-only: language-adapter — the source name as a map key (see
	// CoverageInstrumentationCallback_GetFileName); a null one reaches the map
	// only through get, where Java's map.get(null) is null.
	fileName := ""
	if fileNamePtr != nil {
		fileName = *fileNamePtr
	}

	if node.IsScript() {
		if fileNamePtr != nil && c.instrumentationData.Get(fileName) != nil {
			toAddTo := node
			if node.HasChildren() && node.GetFirstChild().IsModuleBody() {
				toAddTo = node.GetFirstChild()
			}
			// Add instrumentation code
			toAddTo.AddChildrenToFront(c.NewHeaderNode(traversal, toAddTo).RemoveChildren())
			c.compiler.ReportChangeToEnclosingScope(node)
			c.InstrumentBranchCoverage(traversal, c.instrumentationData.Get(fileName))
		}
	}

	if node.IsIf() {
		cfg := c.GetControlFlowGraph(c.compiler)
		hasDefaultBlock := false
		for _, outEdge := range cfg.GetOutEdges(node) {
			if outEdge.GetValue() == ControlFlowGraph_Branch_ON_FALSE {
				destination := outEdge.GetDestination().GetValue()
				if destination != nil &&
					destination.IsBlock() &&
					destination.GetParent() != nil &&
					destination.GetParent().IsIf() {
					hasDefaultBlock = true
				}
				break
			}
		}
		if !hasDefaultBlock {
			c.AddDefaultBlock(node)
		}
		c.instrumentationData.ComputeIfAbsent(
			fileName, func(k string) *FileInstrumentationData {
				return NewFileInstrumentationData(k, BranchCoverageInstrumentationCallback_CreateArrayName(traversal))
			})
		c.ProcessBranchInfo(node, c.instrumentationData.Get(fileName), c.GetChildrenBlocks(node))
	} else if NodeUtil_IsLoopStructure(node) {
		blocks := c.GetChildrenBlocks(node)
		cfg := c.GetControlFlowGraph(c.compiler)
		for _, outEdge := range cfg.GetOutEdges(node) {
			if outEdge.GetValue() == ControlFlowGraph_Branch_ON_FALSE {
				destination := outEdge.GetDestination().GetValue()
				if destination != nil && destination.IsBlock() {
					blocks.Add(destination)
				} else {
					exitBlock := rhino.IR_Block()
					if destination != nil && destination.GetParent().IsBlock() {
						// When the destination of an outEdge of the CFG is not null and the source node's
						// parent is a block. If parent is not a block it may result in an illegal state
						// exception
						exitBlock.InsertBefore(destination)
					} else {
						// When the destination of an outEdge of the CFG is null then we need to add an empty
						// block directly after the loop structure that we can later instrument
						exitBlock.InsertAfter(outEdge.GetSource().GetValue())
					}
					blocks.Add(exitBlock)
				}
			}
		}
		c.instrumentationData.ComputeIfAbsent(
			fileName, func(k string) *FileInstrumentationData {
				return NewFileInstrumentationData(k, BranchCoverageInstrumentationCallback_CreateArrayName(traversal))
			})
		c.ProcessBranchInfo(node, c.instrumentationData.Get(fileName), blocks)
	}
}

// Java: BranchCoverageInstrumentationCallback.getChildrenBlocks(Node)
func (c *BranchCoverageInstrumentationCallback) GetChildrenBlocks(node *rhino.Node) *jutil.ArrayList[*rhino.Node] {
	blocks := jutil.NewArrayList[*rhino.Node]()
	for child := node.GetFirstChild(); child != nil; child = child.GetNext() {
		if child.IsBlock() {
			blocks.Add(child)
		}
	}
	return blocks
}

// InstrumentBranchCoverage adds instrumentation code for branch coverage. For
// each block that correspond to a branch, insert an assignment of the branch
// coverage data to the front of the block.
//
// Java: BranchCoverageInstrumentationCallback.instrumentBranchCoverage(NodeTraversal, FileInstrumentationData)
func (c *BranchCoverageInstrumentationCallback) InstrumentBranchCoverage(traversal *NodeTraversal, data *FileInstrumentationData) {
	maxLine := data.MaxBranchPresentLine()
	branchCoverageOffset := 0
	for lineIdx := 1; lineIdx <= maxLine; lineIdx++ {
		// getNumBranches never returns null (0 for a line without branches).
		numBranches := data.GetNumBranches(lineIdx)
		for branchIdx := 1; branchIdx <= numBranches; branchIdx++ {
			block := data.GetBranchNode(lineIdx, branchIdx)
			block.AddChildToFront(
				c.NewBranchInstrumentationNode(traversal, block, branchCoverageOffset+branchIdx-1))
			c.compiler.ReportChangeToEnclosingScope(block)
		}
		branchCoverageOffset += numBranches
	}
}

// NewBranchInstrumentationNode creates an assignment to the branch coverage
// data for the given index into the array.
//
// Returns the newly constructed assignment node.
//
// Java: BranchCoverageInstrumentationCallback.newBranchInstrumentationNode(NodeTraversal, Node, int)
func (c *BranchCoverageInstrumentationCallback) NewBranchInstrumentationNode(traversal *NodeTraversal, node *rhino.Node, idx int) *rhino.Node {
	arrayName := BranchCoverageInstrumentationCallback_CreateArrayName(traversal)

	// Create instrumentation Node
	getElemNode := rhino.IR_Getelem(rhino.IR_Name(arrayName), rhino.IR_Number(float64(idx))) // Make line number 0-based
	exprNode := rhino.IR_ExprResult(rhino.IR_Assign(getElemNode, rhino.IR_TrueNode()))

	// Note line as instrumented
	fileName := CoverageInstrumentationCallback_GetFileName(traversal)
	c.instrumentationData.ComputeIfAbsent(
		fileName, func(k string) *FileInstrumentationData { return NewFileInstrumentationData(k, arrayName) })
	return exprNode.SrcrefTreeIfMissing(node)
}

// ProcessBranchInfo adds branch instrumentation information for each block.
//
// Java: BranchCoverageInstrumentationCallback.processBranchInfo(Node, FileInstrumentationData, List<Node>)
func (c *BranchCoverageInstrumentationCallback) ProcessBranchInfo(branchNode *rhino.Node, data *FileInstrumentationData, blocks *jutil.ArrayList[*rhino.Node]) {
	lineNumber := branchNode.GetLineno()
	data.SetBranchPresent(lineNumber)

	// Instrument for each block
	numBranches := 0
	for child := range blocks.All() {
		data.PutBranchNode(lineNumber, numBranches+1, child)
		numBranches++
	}
	data.AddBranches(lineNumber, numBranches)
}

// AddDefaultBlock adds a default block for conditional statements, e.g., If,
// Switch.
//
// Java: BranchCoverageInstrumentationCallback.addDefaultBlock(Node)
func (c *BranchCoverageInstrumentationCallback) AddDefaultBlock(node *rhino.Node) *rhino.Node {
	defaultBlock := rhino.IR_Block()
	node.AddChildToBack(defaultBlock)
	return defaultBlock.SrcrefTreeIfMissing(node)
}

// Java: BranchCoverageInstrumentationCallback.newHeaderNode(NodeTraversal, Node)
func (c *BranchCoverageInstrumentationCallback) NewHeaderNode(traversal *NodeTraversal, srcref *rhino.Node) *rhino.Node {
	fileName := CoverageInstrumentationCallback_GetFileName(traversal)
	data := c.instrumentationData.Get(fileName)
	jutil.CheckNotNull(data)

	// var JSCompiler_lcov_branch_data_xx = [];
	// __jscov['branchesTaken'].push(JSCompiler_lcov_branch_data_xx);
	objName := CoverageInstrumentationPass_JS_INSTRUMENTATION_OBJECT_NAME
	var nodes []*rhino.Node
	nodes = append(nodes, c.NewArrayDeclarationNode(traversal))
	nodes = append(nodes,
		rhino.IR_ExprResult(
			rhino.IR_Call(
				rhino.IR_Getprop(rhino.IR_Getelem(rhino.IR_Name(objName), rhino.IR_String("branchesTaken")), "push"),
				rhino.IR_Name(BranchCoverageInstrumentationCallback_CreateArrayName(traversal)))))
	// __jscov['branchPresent'].push(hex-data);
	nodes = append(nodes,
		rhino.IR_ExprResult(
			rhino.IR_Call(
				rhino.IR_Getprop(rhino.IR_Getelem(rhino.IR_Name(objName), rhino.IR_String("branchPresent")), "push"),
				rhino.IR_String(data.GetBranchPresentAsHexString()))))
	nodes = append(nodes, c.NewBranchesInLineNode("JSCompiler_lcov_branchesInLine", data))
	// __jscov['branchesInLine'].push(JSCompiler_lcov_branchesInLine);
	nodes = append(nodes,
		rhino.IR_ExprResult(
			rhino.IR_Call(
				rhino.IR_Getprop(rhino.IR_Getelem(rhino.IR_Name(objName), rhino.IR_String("branchesInLine")), "push"),
				rhino.IR_Name("JSCompiler_lcov_branchesInLine"))))
	// __jscov['fileNames'].push(filename);
	nodes = append(nodes,
		rhino.IR_ExprResult(
			rhino.IR_Call(
				rhino.IR_Getprop(rhino.IR_Getelem(rhino.IR_Name(objName), rhino.IR_String("fileNames")), "push"),
				rhino.IR_String(fileName))))
	return rhino.IR_BlockOfList(nodes).SrcrefTreeIfMissing(srcref)
}

// Java: BranchCoverageInstrumentationCallback.newArrayDeclarationNode(NodeTraversal)
func (c *BranchCoverageInstrumentationCallback) NewArrayDeclarationNode(traversal *NodeTraversal) *rhino.Node {
	return rhino.IR_VarWithValue(rhino.IR_Name(BranchCoverageInstrumentationCallback_CreateArrayName(traversal)), rhino.IR_Arraylit())
}

// Java: BranchCoverageInstrumentationCallback.newBranchesInLineNode(String, FileInstrumentationData)
func (c *BranchCoverageInstrumentationCallback) NewBranchesInLineNode(name string, data *FileInstrumentationData) *rhino.Node {
	var assignments []*rhino.Node
	// var JSCompiler_lcov_branchesInLine = [];
	assignments = append(assignments, rhino.IR_VarWithValue(rhino.IR_Name(name), rhino.IR_Arraylit()))
	lineWithBranch := 0
	for lineIdx := 1; lineIdx <= data.MaxBranchPresentLine(); lineIdx++ {
		numBranches := data.GetNumBranches(lineIdx)
		if numBranches > 0 {
			// JSCompiler_lcov_branchesInLine[<branch-index>] = 2;
			assignment := rhino.IR_ExprResult(
				rhino.IR_Assign(
					rhino.IR_Getelem(rhino.IR_Name(name), rhino.IR_Number(float64(lineWithBranch))),
					rhino.IR_Number(float64(numBranches))))
			lineWithBranch++
			assignments = append(assignments, assignment.SrcrefTreeIfMissing(assignment))
		}
	}
	return rhino.IR_BlockOfList(assignments)
}
